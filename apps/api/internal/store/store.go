// Package store is the broker's persistence boundary: the Ent-defined
// schema, its derived and embedded migrations, and (from Batch 2 on) the
// narrow Repository — per structure/folder-structure and
// data/data-access-entity-framework.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent"
	"github.com/franwerner/matecito-ai/apps/api/internal/store/migrations"

	atlas "ariga.io/atlas/sql/migrate"
	atsqlite "ariga.io/atlas/sql/sqlite"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite" // registers driver "sqlite" (pure Go, no cgo)
)

// dbFileName is the single global SQLite database file inside the daemon's
// persist-dir (storage-sync-model: one self-contained, deployable store).
const dbFileName = "broker.db"

// Store owns the broker's single SQLite connection. Batch 1 only opens and
// migrates it; the narrow Repository boundary (write side / read side) that
// wraps Client is Batch 2 — see data/data-access-entity-framework.
type Store struct {
	Client *ent.Client
	db     *sql.DB
}

// Open opens (creating if absent) the database under persistDir, enables
// foreign-key enforcement — SQLite ignores FKs by default and the generated
// client does not change that (ROADMAP-2 gotcha #1) — and applies every
// pending embedded migration before returning. Returns a wrapped error
// without serving on any failure (R1, R14): the daemon must not run with an
// unmigrated or partially migrated schema.
func Open(ctx context.Context, persistDir string) (*Store, error) {
	dbPath := filepath.Join(persistDir, dbFileName)
	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)", dbPath)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("store: open %s: %w", dbPath, err)
	}
	// The daemon is the single writer for this database; one connection
	// avoids SQLITE_BUSY from the pool racing itself over the same file.
	db.SetMaxOpenConns(1)

	if err := migrateUp(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("store: migrate %s: %w", dbPath, err)
	}

	drv := entsql.OpenDB(dialect.SQLite, db)
	return &Store{Client: ent.NewClient(ent.Driver(drv)), db: db}, nil
}

// Close releases the underlying connection (closing the Ent client closes
// the wrapped *sql.DB too).
func (s *Store) Close() error {
	return s.Client.Close()
}

// migrateUp applies every pending file from the embedded migration
// directory through Atlas's pure-Go executor — no atlas CLI, no cgo driver,
// per data-access-entity-framework's "binario autocontenido" and
// modernc-sqlite's "sin cgo".
func migrateUp(ctx context.Context, db *sql.DB) error {
	rrw := newRevisions(db)
	if err := rrw.ensureTable(ctx); err != nil {
		return fmt.Errorf("ensure revisions table: %w", err)
	}

	dir, err := embeddedDir()
	if err != nil {
		return fmt.Errorf("read embedded migrations: %w", err)
	}

	drv, err := atsqlite.Open(db)
	if err != nil {
		return fmt.Errorf("open atlas sqlite driver: %w", err)
	}

	ex, err := atlas.NewExecutor(drv, dir, rrw)
	if err != nil {
		return fmt.Errorf("create migration executor: %w", err)
	}
	// n=0 applies every pending file. A clean, already-migrated database
	// has none pending — Atlas reports that as ErrNoPendingFiles rather
	// than a no-op success, which is exactly the idempotent-restart case
	// (R1), not a failure.
	if err := ex.ExecuteN(ctx, 0); err != nil && !errors.Is(err, atlas.ErrNoPendingFiles) {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

// embeddedDir bridges the embed.FS into an atlas migrate.Dir. MemDir.
// CopyFiles recomputes the checksum file from exactly these files, so the
// statically embedded atlas.sum is never consulted at runtime — there is no
// risk of it drifting from the embedded .sql files.
func embeddedDir() (*atlas.MemDir, error) {
	entries, err := migrations.FS.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var files []atlas.File
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		data, err := migrations.FS.ReadFile(entry.Name())
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", entry.Name(), err)
		}
		files = append(files, atlas.NewLocalFile(entry.Name(), data))
	}

	dir := &atlas.MemDir{}
	if err := dir.CopyFiles(files); err != nil {
		return nil, err
	}
	return dir, nil
}
