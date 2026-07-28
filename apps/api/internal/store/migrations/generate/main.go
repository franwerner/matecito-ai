//go:build ignore

// Command generate derives a new versioned migration file from the current
// Ent schema definition, per data-access-entity-framework: migrations are
// never edited by hand, only regenerated from the schema in code. Invoke via
// `go generate ./internal/store/migrations` (see generate.go), which passes
// a migration name as the sole argument.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent/migrate"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("migration name is required. Use: go generate ./internal/store/migrations -- <name>")
	}

	ctx := context.Background()
	dir, err := atlas.NewLocalDir("internal/store/migrations")
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}

	// Diffed against a throwaway, empty dev database: what matters is the
	// Ent schema (desired state), not any persisted data.
	devDB := filepath.Join(os.TempDir(), fmt.Sprintf("broker-migrate-dev-%d.db", os.Getpid()))
	defer os.Remove(devDB)
	db, err := sql.Open("sqlite", "file:"+devDB+"?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("failed opening dev database: %v", err)
	}
	defer db.Close()

	drv := entsql.OpenDB(dialect.SQLite, db)
	m, err := schema.NewMigrate(drv,
		schema.WithDir(dir),
		schema.WithDialect(dialect.SQLite),
		schema.WithFormatter(atlas.DefaultFormatter),
	)
	if err != nil {
		log.Fatalf("failed creating migrate engine: %v", err)
	}
	if err := m.NamedDiff(ctx, os.Args[1], migrate.Tables...); err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
}
