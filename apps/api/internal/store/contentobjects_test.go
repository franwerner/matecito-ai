package store_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/apps/api/internal/store"
)

// TestPutContentObject_DedupsByHash covers R7's "dedup por hash": storing
// the same content twice resolves to the same row, never a duplicate.
func TestPutContentObject_DedupsByHash(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	first, err := st.PutContentObject(ctx, store.PutContentObjectCmd{ProjectID: p.ID, Bytes: []byte("hello world")})
	if err != nil {
		t.Fatalf("first PutContentObject: %v", err)
	}

	second, err := st.PutContentObject(ctx, store.PutContentObjectCmd{ProjectID: p.ID, Bytes: []byte("hello world")})
	if err != nil {
		t.Fatalf("second PutContentObject: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("second.ID = %q, want %q (same content, not a duplicate)", second.ID, first.ID)
	}
}

// TestPutContentObject_HashCarriesAlgorithm covers "el hash porta su
// algoritmo": the stored key is shaped `sha256:<hex>` and matches the
// content's real SHA-256.
func TestPutContentObject_HashCarriesAlgorithm(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	obj, err := st.PutContentObject(ctx, store.PutContentObjectCmd{ProjectID: p.ID, Bytes: []byte("abc")})
	if err != nil {
		t.Fatalf("PutContentObject: %v", err)
	}
	// SHA-256("abc"), well known test vector.
	const want = "sha256:ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if obj.Hash != want {
		t.Fatalf("Hash = %q, want %q", obj.Hash, want)
	}
}

// TestPutContentObject_RoundTripsArbitraryBytes covers "round-trip sin
// corrupción": multibyte UTF-8, CRLF and null bytes all come back
// identical.
func TestPutContentObject_RoundTripsArbitraryBytes(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	original := []byte("héllo\r\nwörld\x00\x01\x02")
	put, err := st.PutContentObject(ctx, store.PutContentObjectCmd{ProjectID: p.ID, Bytes: original})
	if err != nil {
		t.Fatalf("PutContentObject: %v", err)
	}

	got, err := st.GetContentObject(ctx, p.ID, put.Hash)
	if err != nil {
		t.Fatalf("GetContentObject: %v", err)
	}
	if string(got.Bytes) != string(original) {
		t.Fatalf("round-tripped bytes = %q, want %q", got.Bytes, original)
	}
	if got.SizeBytes != int64(len(original)) {
		t.Fatalf("SizeBytes = %d, want %d", got.SizeBytes, len(original))
	}
}

// TestGetContentObject_NotFound covers the Reader's ErrNotFound contract.
func TestGetContentObject_NotFound(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	_, err := st.GetContentObject(ctx, p.ID, "sha256:"+strings.Repeat("0", 64))
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

// TestGetContentObject_InvalidHashRejected covers the store's own hash-shape
// validation: a malformed key is rejected before ever touching the base.
func TestGetContentObject_InvalidHashRejected(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	_, err := st.GetContentObject(ctx, p.ID, "not-a-hash")
	if !errors.Is(err, store.ErrInvalidContentHash) {
		t.Fatalf("err = %v, want ErrInvalidContentHash", err)
	}
}

// TestContentObjectHashUnique_RawInsertRejected exercises R7's unique index
// directly against the database, bypassing the Repository and the
// generated client (delivery/testing-strategy: constraint rejections are
// exercised with raw SQL, not through the Repository). This is the
// demonstration behind enmienda 7's own justification for moving
// content_objects off a composite primary key: "el índice único impide los
// duplicados exactamente igual que lo haría la clave primaria" — a
// surrogate id must not make the natural key duplicable.
func TestContentObjectHashUnique_RawInsertRejected(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	st, err := store.Open(ctx, dir, testMachineID)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() {
		if err := st.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})

	raw := openRaw(t, dir)
	fx := seedFixtures(ctx, t, raw)

	// Negative: same (project_id, hash) as the seeded content object, but a
	// distinct id — the surrogate id must not enable duplicating the
	// natural key.
	_, err = raw.ExecContext(ctx, `INSERT INTO content_objects
		(id, created_at, updated_at, project_id, hash, size_bytes, bytes)
		VALUES ('duplicate-content', ?, ?, ?, 'sha256:abc', 3, ?)`,
		ts, ts, fx.projectID, []byte("abc"))
	if err == nil {
		t.Fatal("expected a UNIQUE(project_id, hash) violation, got nil error")
	}

	// Positive control: the same hash under a DIFFERENT project is
	// accepted — the dedup is between branches within a project, never
	// across projects.
	if _, err := raw.ExecContext(ctx, `INSERT INTO projects
		(id, created_at, updated_at, name, status)
		VALUES ('project-2', ?, ?, 'beta', 'active')`, ts, ts); err != nil {
		t.Fatalf("seed second project: %v", err)
	}
	if _, err := raw.ExecContext(ctx, `INSERT INTO content_objects
		(id, created_at, updated_at, project_id, hash, size_bytes, bytes)
		VALUES ('content-in-project-2', ?, ?, 'project-2', 'sha256:abc', 3, ?)`,
		ts, ts, []byte("abc")); err != nil {
		t.Fatalf("same hash under a different project should succeed: %v", err)
	}
}
