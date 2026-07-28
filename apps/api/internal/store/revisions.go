package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	atlas "ariga.io/atlas/sql/migrate"
)

// revisionsTable tracks which embedded migration files were already applied,
// so a restart against an already-migrated database is a no-op (R1
// "arranque idempotente"). It is bootstrapping/tooling state, not part of
// the 11-table domain contract in ROADMAP-2 Anexo A, so it is created by raw
// SQL here rather than through the Ent schema.
const revisionsTable = "atlas_schema_revisions"

const createRevisionsTableSQL = `CREATE TABLE IF NOT EXISTS ` + revisionsTable + ` (
	version text NOT NULL PRIMARY KEY,
	description text NOT NULL,
	type integer NOT NULL,
	applied integer NOT NULL,
	total integer NOT NULL,
	executed_at text NOT NULL,
	execution_time_ns integer NOT NULL,
	error text NOT NULL DEFAULT '',
	error_stmt text NOT NULL DEFAULT '',
	hash text NOT NULL DEFAULT '',
	partial_hashes text NOT NULL DEFAULT '',
	operator_version text NOT NULL DEFAULT ''
)`

// revisions implements atlas's migrate.RevisionReadWriter against the raw
// *sql.DB — not the Ent client, so migration bookkeeping never depends on
// the generated code it precedes.
type revisions struct {
	db *sql.DB
}

func newRevisions(db *sql.DB) *revisions {
	return &revisions{db: db}
}

func (r *revisions) ensureTable(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, createRevisionsTableSQL)
	return err
}

func (r *revisions) Ident() *atlas.TableIdent {
	return &atlas.TableIdent{Name: revisionsTable}
}

const revisionColumns = `version, description, type, applied, total, executed_at, execution_time_ns, error, error_stmt, hash, partial_hashes, operator_version`

func (r *revisions) ReadRevisions(ctx context.Context) ([]*atlas.Revision, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+revisionColumns+` FROM `+revisionsTable+` ORDER BY version`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*atlas.Revision
	for rows.Next() {
		rev, err := scanRevision(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, rev)
	}
	return out, rows.Err()
}

func (r *revisions) ReadRevision(ctx context.Context, version string) (*atlas.Revision, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+revisionColumns+` FROM `+revisionsTable+` WHERE version = ?`, version)
	rev, err := scanRevision(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, atlas.ErrRevisionNotExist
	}
	if err != nil {
		return nil, err
	}
	return rev, nil
}

func (r *revisions) WriteRevision(ctx context.Context, rev *atlas.Revision) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO `+revisionsTable+` (`+revisionColumns+`)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (version) DO UPDATE SET
			description = excluded.description,
			type = excluded.type,
			applied = excluded.applied,
			total = excluded.total,
			executed_at = excluded.executed_at,
			execution_time_ns = excluded.execution_time_ns,
			error = excluded.error,
			error_stmt = excluded.error_stmt,
			hash = excluded.hash,
			partial_hashes = excluded.partial_hashes,
			operator_version = excluded.operator_version`,
		rev.Version, rev.Description, uint(rev.Type), rev.Applied, rev.Total,
		rev.ExecutedAt.UTC().Format(time.RFC3339Nano), rev.ExecutionTime.Nanoseconds(),
		rev.Error, rev.ErrorStmt, rev.Hash, strings.Join(rev.PartialHashes, ","), rev.OperatorVersion,
	)
	return err
}

func (r *revisions) DeleteRevision(ctx context.Context, version string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM `+revisionsTable+` WHERE version = ?`, version)
	return err
}

// rowScanner is satisfied by both *sql.Row and *sql.Rows.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanRevision(row rowScanner) (*atlas.Revision, error) {
	var (
		rev              atlas.Revision
		revType          uint
		executedAt       string
		executionTimeNs  int64
		partialHashesRaw string
	)
	if err := row.Scan(
		&rev.Version, &rev.Description, &revType, &rev.Applied, &rev.Total,
		&executedAt, &executionTimeNs, &rev.Error, &rev.ErrorStmt, &rev.Hash,
		&partialHashesRaw, &rev.OperatorVersion,
	); err != nil {
		return nil, err
	}
	rev.Type = atlas.RevisionType(revType)
	rev.ExecutionTime = time.Duration(executionTimeNs)
	if partialHashesRaw != "" {
		rev.PartialHashes = strings.Split(partialHashesRaw, ",")
	}
	t, err := time.Parse(time.RFC3339Nano, executedAt)
	if err != nil {
		return nil, fmt.Errorf("parse executed_at %q: %w", executedAt, err)
	}
	rev.ExecutedAt = t
	return &rev, nil
}
