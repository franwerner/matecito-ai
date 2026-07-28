package store_test

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"slices"
	"testing"

	"github.com/franwerner/matecito-ai/apps/api/internal/store"

	_ "modernc.org/sqlite"
)

// ts is a fixed, valid RFC 3339-with-milliseconds timestamp for fixture rows
// — the exact convention created_at/updated_at/occurred_at must follow.
const ts = "2026-01-01T00:00:00.000Z"

// testMachineID stands in for the daemon's own persisted machine_id
// (store.EnsureMachineID) — these tests only care about schema/migration
// behavior, not the value itself.
const testMachineID = "test-machine"

// openRaw opens a second, completely independent connection to the same
// database file store.Open just migrated — no Ent, no store package code —
// so tests exercise the database itself (delivery/testing-strategy: SQLite
// real, no mocks; ROADMAP-2 gotcha #3: CHECKs must reject a raw INSERT that
// never touches the persistence edge or the generated client).
func openRaw(t *testing.T, persistDir string) *sql.DB {
	t.Helper()
	dbPath := filepath.Join(persistDir, "broker.db")
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_pragma=foreign_keys(1)", dbPath))
	if err != nil {
		t.Fatalf("openRaw: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func queryNames(t *testing.T, db *sql.DB, query string) []string {
	t.Helper()
	rows, err := db.Query(query)
	if err != nil {
		t.Fatalf("query %q: %v", query, err)
	}
	defer rows.Close()

	var got []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan: %v", err)
		}
		got = append(got, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows: %v", err)
	}
	return got
}

func TestOpen_CreatesSchemaFromScratch(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()

	st, err := store.Open(ctx, dir, testMachineID)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer st.Close()

	raw := openRaw(t, dir)

	wantTables := []string{
		"changes", "content_objects", "event_record_pins", "events",
		"iteration_notes", "project_paths", "projects", "record_branch_state",
		"record_version_refs", "record_versions", "records",
	}
	gotTables := queryNames(t, raw, `SELECT name FROM sqlite_master
		WHERE type = 'table' AND name NOT LIKE 'sqlite_%' AND name != 'atlas_schema_revisions'
		ORDER BY name`)
	slices.Sort(wantTables)
	if !slices.Equal(gotTables, wantTables) {
		t.Fatalf("tables = %v, want %v", gotTables, wantTables)
	}

	wantIndexes := []string{
		"change_project_id_branch",
		"change_project_id_name",
		"contentobject_project_id_hash",
		"event_change_id_position",
		"event_idempotency_hash",
		"eventrecordpin_event_id_record_version_id",
		"projectpath_machine_id_root_path",
		"record_project_id_owning_root_slug",
		"recordbranchstate_record_id_branch",
	}
	gotIndexes := queryNames(t, raw, `SELECT name FROM sqlite_master
		WHERE type = 'index' AND name NOT LIKE 'sqlite_%'
		ORDER BY name`)
	slices.Sort(wantIndexes)
	if !slices.Equal(gotIndexes, wantIndexes) {
		t.Fatalf("unique indexes = %v, want %v", gotIndexes, wantIndexes)
	}
}

// TestOpen_IdempotentOnAlreadyMigrated covers the R1 "arranque idempotente"
// scenario: re-opening an already-migrated database applies nothing new.
func TestOpen_IdempotentOnAlreadyMigrated(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()

	st1, err := store.Open(ctx, dir, testMachineID)
	if err != nil {
		t.Fatalf("first Open: %v", err)
	}
	if err := st1.Close(); err != nil {
		t.Fatalf("close first: %v", err)
	}

	st2, err := store.Open(ctx, dir, testMachineID)
	if err != nil {
		t.Fatalf("second Open: %v", err)
	}
	defer st2.Close()

	raw := openRaw(t, dir)
	var applied int
	if err := raw.QueryRow(`SELECT COUNT(*) FROM atlas_schema_revisions`).Scan(&applied); err != nil {
		t.Fatalf("count revisions: %v", err)
	}
	if applied != 1 {
		t.Fatalf("applied revisions = %d, want 1 (migration file count)", applied)
	}
}

// TestForeignKeysEnforced covers R2's "las FKs se enforcen" scenario and
// ROADMAP-2 gotcha #1: SQLite ignores foreign keys unless the DSN turns
// them on, and the generated client does not do that for us.
func TestForeignKeysEnforced(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	st, err := store.Open(ctx, dir, testMachineID)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer st.Close()

	raw := openRaw(t, dir)
	_, err = raw.ExecContext(ctx, `INSERT INTO changes
		(id, created_at, updated_at, branch, name, status, project_id)
		VALUES ('change-1', ?, ?, 'main', 'orphan-change', 'active', 'no-such-project')`, ts, ts)
	if err == nil {
		t.Fatal("expected a foreign key violation inserting a change for a nonexistent project, got nil error")
	}
}

// fixtures holds the minimal valid parent rows CHECK-constraint cases below
// build on: a project, a change, a record, a content object, a record
// version and an event, each satisfying every FK and CHECK on its own row.
type fixtures struct {
	projectID       string
	changeID        string
	recordID        string
	contentObjectID string
	recordVersionID string
	eventID         string
}

func seedFixtures(ctx context.Context, t *testing.T, db *sql.DB) fixtures {
	t.Helper()
	fx := fixtures{
		projectID:       "project-1",
		changeID:        "change-1",
		recordID:        "record-1",
		contentObjectID: "content-1",
		recordVersionID: "version-1",
		eventID:         "event-1",
	}

	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := db.ExecContext(ctx, query, args...); err != nil {
			t.Fatalf("seed fixture %q: %v", query, err)
		}
	}

	exec(`INSERT INTO projects (id, created_at, updated_at, name, status)
		VALUES (?, ?, ?, 'acme', 'active')`, fx.projectID, ts, ts)
	exec(`INSERT INTO changes (id, created_at, updated_at, branch, name, status, project_id)
		VALUES (?, ?, ?, 'main', 'acme-change', 'active', ?)`, fx.changeID, ts, ts, fx.projectID)
	exec(`INSERT INTO records (id, created_at, updated_at, owning_root, slug, project_id)
		VALUES (?, ?, ?, 'apps/api', 'my-record', ?)`, fx.recordID, ts, ts, fx.projectID)
	exec(`INSERT INTO content_objects (id, created_at, updated_at, hash, size_bytes, bytes, project_id)
		VALUES (?, ?, ?, 'sha256:abc', 3, ?, ?)`, fx.contentObjectID, ts, ts, []byte("abc"), fx.projectID)
	exec(`INSERT INTO record_versions (id, created_at, updated_at, status, project_id, record_id, content_object_id)
		VALUES (?, ?, ?, 'current', ?, ?, ?)`, fx.recordVersionID, ts, ts, fx.projectID, fx.recordID, fx.contentObjectID)
	exec(`INSERT INTO events (id, created_at, updated_at, position, type, payload, occurred_at, idempotency_hash, project_id, change_id)
		VALUES (?, ?, ?, 1, 'submit', '{}', ?, 'hash-1', ?, ?)`, fx.eventID, ts, ts, ts, fx.projectID, fx.changeID)

	return fx
}

// TestCheckConstraintsRejectRawInserts is the hardened gotcha #3 criterion:
// every CHECK the schema declares must be enforced by the database itself,
// not only by the Ent client — verified here with raw INSERT statements
// that go through neither the generated client nor a future persistence
// edge.
func TestCheckConstraintsRejectRawInserts(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	st, err := store.Open(ctx, dir, testMachineID)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer st.Close()

	raw := openRaw(t, dir)
	fx := seedFixtures(ctx, t, raw)

	cases := []struct {
		name  string
		query string
		args  []any
	}{
		{
			name:  "projects.status outside enum",
			query: `INSERT INTO projects (id, created_at, updated_at, name, status) VALUES ('bad-project', ?, ?, 'acme', 'bogus')`,
			args:  []any{ts, ts},
		},
		{
			name:  "project_paths.status outside enum",
			query: `INSERT INTO project_paths (id, created_at, updated_at, machine_id, root_path, status, last_seen_at, project_id) VALUES ('bad-path', ?, ?, 'machine-1', '/tmp/x', 'bogus', ?, ?)`,
			args:  []any{ts, ts, ts, fx.projectID},
		},
		{
			name:  "changes.status outside enum",
			query: `INSERT INTO changes (id, created_at, updated_at, branch, name, status, project_id) VALUES ('bad-change', ?, ?, 'feat/x', 'bad-change', 'bogus', ?)`,
			args:  []any{ts, ts, fx.projectID},
		},
		{
			name:  "events.agent_kind outside enum",
			query: `INSERT INTO events (id, created_at, updated_at, position, type, payload, occurred_at, agent_kind, idempotency_hash, project_id, change_id) VALUES ('bad-event', ?, ?, 2, 'submit', '{}', ?, 'bogus', 'hash-2', ?, ?)`,
			args:  []any{ts, ts, ts, fx.projectID, fx.changeID},
		},
		{
			name:  "record_versions.status outside enum",
			query: `INSERT INTO record_versions (id, created_at, updated_at, status, project_id, record_id, content_object_id) VALUES ('bad-version', ?, ?, 'bogus', ?, ?, ?)`,
			args:  []any{ts, ts, fx.projectID, fx.recordID, fx.contentObjectID},
		},
		{
			name:  "record_branch_state.status outside enum",
			query: `INSERT INTO record_branch_state (id, created_at, updated_at, branch, status, project_id, record_id, current_version_id) VALUES ('bad-branch-state', ?, ?, 'main', 'bogus', ?, ?, ?)`,
			args:  []any{ts, ts, fx.projectID, fx.recordID, fx.recordVersionID},
		},
		{
			name:  "iteration_notes.status outside enum",
			query: `INSERT INTO iteration_notes (id, created_at, updated_at, body, status, project_id, change_id) VALUES ('bad-note-status', ?, ?, 'body', 'bogus', ?, ?)`,
			args:  []any{ts, ts, fx.projectID, fx.changeID},
		},
		{
			name:  "iteration_notes file_anchor: file_path without event_id",
			query: `INSERT INTO iteration_notes (id, created_at, updated_at, body, status, file_path, project_id, change_id) VALUES ('bad-note-anchor', ?, ?, 'body', 'pending', 'foo.go', ?, ?)`,
			args:  []any{ts, ts, fx.projectID, fx.changeID},
		},
		{
			name:  "iteration_notes range_anchor: line_from without file_path",
			query: `INSERT INTO iteration_notes (id, created_at, updated_at, body, status, line_from, line_to, project_id, change_id) VALUES ('bad-note-range', ?, ?, 'body', 'pending', 5, 10, ?, ?)`,
			args:  []any{ts, ts, fx.projectID, fx.changeID},
		},
		{
			name:  "iteration_notes range_anchor: line_from > line_to",
			query: `INSERT INTO iteration_notes (id, created_at, updated_at, body, status, file_path, line_from, line_to, event_id, project_id, change_id) VALUES ('bad-note-inverted', ?, ?, 'body', 'pending', 'foo.go', 10, 5, ?, ?, ?)`,
			args:  []any{ts, ts, fx.eventID, fx.projectID, fx.changeID},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := raw.ExecContext(ctx, tc.query, tc.args...)
			if err == nil {
				t.Fatalf("expected a CHECK violation, got nil error")
			}
		})
	}
}

// TestCheckConstraintsAllowValidRows is the positive control for the table
// above: the same anchoring shapes, with values inside the CHECKs, must
// succeed — confirming the constraints reject exactly the invalid shape and
// nothing more.
func TestCheckConstraintsAllowValidRows(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	st, err := store.Open(ctx, dir, testMachineID)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer st.Close()

	raw := openRaw(t, dir)
	fx := seedFixtures(ctx, t, raw)

	// A note anchored to change only (no event, no file, no range) —
	// the least anchored valid shape.
	if _, err := raw.ExecContext(ctx, `INSERT INTO iteration_notes
		(id, created_at, updated_at, body, status, project_id, change_id)
		VALUES ('good-note-bare', ?, ?, 'body', 'pending', ?, ?)`,
		ts, ts, fx.projectID, fx.changeID); err != nil {
		t.Fatalf("bare note insert should succeed: %v", err)
	}

	// A note anchored to event + file + a valid ascending range.
	if _, err := raw.ExecContext(ctx, `INSERT INTO iteration_notes
		(id, created_at, updated_at, body, status, file_path, line_from, line_to, event_id, project_id, change_id)
		VALUES ('good-note-ranged', ?, ?, 'body', 'pending', 'foo.go', 5, 10, ?, ?, ?)`,
		ts, ts, fx.eventID, fx.projectID, fx.changeID); err != nil {
		t.Fatalf("ranged note insert should succeed: %v", err)
	}

	if _, err := raw.ExecContext(ctx, `INSERT INTO record_branch_state
		(id, created_at, updated_at, branch, status, project_id, record_id, current_version_id)
		VALUES ('good-branch-state', ?, ?, 'main', 'active', ?, ?, ?)`,
		ts, ts, fx.projectID, fx.recordID, fx.recordVersionID); err != nil {
		t.Fatalf("valid record_branch_state insert should succeed: %v", err)
	}
}
