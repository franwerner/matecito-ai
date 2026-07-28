package store_test

import (
	"context"
	"testing"

	"github.com/franwerner/matecito-ai/apps/api/internal/store"
)

func writeRecordVersion(ctx context.Context, t *testing.T, st store.Repository, projectID, owningRoot, slug, branch string, bytes []byte) store.RecordVersion {
	t.Helper()
	v, err := st.WriteRecordVersion(ctx, store.WriteRecordVersionCmd{
		ProjectID:  projectID,
		OwningRoot: owningRoot,
		Slug:       slug,
		Branch:     branch,
		Bytes:      bytes,
		Kind:       "edr",
		DocStatus:  "Accepted",
	})
	if err != nil {
		t.Fatalf("WriteRecordVersion(%s@%s): %v", slug, branch, err)
	}
	return v
}

// TestUpsertRecord_Idempotent covers process/index-decision-records'
// identity: the same (project, owning_root, slug) always resolves to the
// same row, never a duplicate.
func TestUpsertRecord_Idempotent(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	first, err := st.UpsertRecord(ctx, store.UpsertRecordCmd{ProjectID: p.ID, OwningRoot: "api", Slug: "data-modeling"})
	if err != nil {
		t.Fatalf("first UpsertRecord: %v", err)
	}
	second, err := st.UpsertRecord(ctx, store.UpsertRecordCmd{ProjectID: p.ID, OwningRoot: "api", Slug: "data-modeling"})
	if err != nil {
		t.Fatalf("second UpsertRecord: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("second.ID = %q, want %q (same record, not a duplicate)", second.ID, first.ID)
	}
}

// TestWriteRecordVersion_FirstWrite covers "un `.md` nuevo aparece": the
// record's first version becomes the branch's current version.
func TestWriteRecordVersion_FirstWrite(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	v := writeRecordVersion(ctx, t, st, p.ID, "api", "data-modeling", "main", []byte("v1"))
	if v.Status != store.VersionCurrent {
		t.Fatalf("Status = %q, want %q", v.Status, store.VersionCurrent)
	}

	rec, err := st.GetRecord(ctx, p.ID, "api", "data-modeling")
	if err != nil {
		t.Fatalf("GetRecord: %v", err)
	}
	bs, err := st.GetRecordBranchState(ctx, rec.ID, "main")
	if err != nil {
		t.Fatalf("GetRecordBranchState: %v", err)
	}
	if bs.CurrentVersionID != v.ID {
		t.Fatalf("CurrentVersionID = %q, want %q", bs.CurrentVersionID, v.ID)
	}
	if bs.Status != store.RecordBranchActive {
		t.Fatalf("Status = %q, want %q", bs.Status, store.RecordBranchActive)
	}
}

// TestWriteRecordVersion_UpdateInPlace covers the "update-in-place sin
// referencia" scenario: a version that is not frozen updates its own row
// instead of forking.
func TestWriteRecordVersion_UpdateInPlace(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	first := writeRecordVersion(ctx, t, st, p.ID, "api", "data-modeling", "main", []byte("v1"))
	second := writeRecordVersion(ctx, t, st, p.ID, "api", "data-modeling", "main", []byte("v2"))

	if second.ID != first.ID {
		t.Fatalf("second.ID = %q, want %q (update-in-place, not a fork)", second.ID, first.ID)
	}
	if second.ContentObjectID == first.ContentObjectID {
		t.Fatal("ContentObjectID unchanged, want it to point at the new content")
	}
}

// TestWriteRecordVersion_ForkOnFrozenVersion covers "fork al cambiar una
// versión congelada": once an event pins the current version, the next
// write forks instead of touching it, and the pin keeps resolving the
// frozen version with its original content.
func TestWriteRecordVersion_ForkOnFrozenVersion(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "main", "add-feature")

	frozen := writeRecordVersion(ctx, t, st, p.ID, "api", "data-modeling", "main", []byte("v1"))

	rec, err := st.GetRecord(ctx, p.ID, "api", "data-modeling")
	if err != nil {
		t.Fatalf("GetRecord: %v", err)
	}
	ev, err := st.AppendEvent(ctx, store.AppendEventCmd{ProjectID: p.ID, ChangeID: c.ID, Type: "submit", Payload: `{"n":1}`})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	if _, err := st.PinEventToRecordVersion(ctx, store.PinEventToRecordVersionCmd{
		ProjectID: p.ID, EventID: ev.ID, RecordID: rec.ID, Branch: "main",
	}); err != nil {
		t.Fatalf("PinEventToRecordVersion: %v", err)
	}

	forked := writeRecordVersion(ctx, t, st, p.ID, "api", "data-modeling", "main", []byte("v2"))
	if forked.ID == frozen.ID {
		t.Fatal("forked.ID == frozen.ID, want a new version (frozen must fork, not update in place)")
	}

	stillFrozen, err := st.GetRecordVersion(ctx, frozen.ID)
	if err != nil {
		t.Fatalf("GetRecordVersion(frozen): %v", err)
	}
	if stillFrozen.Status != store.VersionFrozen {
		t.Fatalf("frozen version Status = %q, want %q", stillFrozen.Status, store.VersionFrozen)
	}
	if stillFrozen.ContentObjectID != frozen.ContentObjectID {
		t.Fatal("frozen version's content changed, want it intact")
	}
}

// TestWriteRecordVersion_IdenticalContentDoesNotShareVersion covers the
// spec's "dos ramas con contenido idéntico no comparten versión": each
// branch's first write always creates its own version row, even when the
// `.md` is byte-for-byte identical to another branch's — only the content
// they reference is deduplicated (TestWriteRecordVersion_ContentSharedAcrossBranchesByHash
// below). Editing one branch's version must never touch the other's.
func TestWriteRecordVersion_IdenticalContentDoesNotShareVersion(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	onMain := writeRecordVersion(ctx, t, st, p.ID, "api", "data-modeling", "main", []byte("identical"))
	onFeature := writeRecordVersion(ctx, t, st, p.ID, "api", "data-modeling", "feat/x", []byte("identical"))
	if onFeature.ID == onMain.ID {
		t.Fatalf("feat/x's first version = %q, want its own version distinct from main's %q (versions are never shared)", onFeature.ID, onMain.ID)
	}
	if onFeature.ContentObjectID != onMain.ContentObjectID {
		t.Fatalf("ContentObjectID = %q, want %q (identical content still dedups to one row)", onFeature.ContentObjectID, onMain.ContentObjectID)
	}

	changed := writeRecordVersion(ctx, t, st, p.ID, "api", "data-modeling", "feat/x", []byte("changed on feat/x"))
	if changed.ID != onFeature.ID {
		t.Fatalf("changed.ID = %q, want %q (update-in-place, not a fork — feat/x's version was never frozen or shared)", changed.ID, onFeature.ID)
	}

	original, err := st.GetRecordVersion(ctx, onMain.ID)
	if err != nil {
		t.Fatalf("GetRecordVersion(main's version): %v", err)
	}
	if original.ContentObjectID != onMain.ContentObjectID {
		t.Fatal("main's version content changed, want it untouched by feat/x's edit")
	}
}

// TestWriteRecordVersion_ContentSharedAcrossBranchesByHash covers "contenido
// compartido por hash entre ramas": two branches' versions, even with
// different projections, resolve to one deduplicated content_objects row —
// the version row is never shared, but the content it references is.
func TestWriteRecordVersion_ContentSharedAcrossBranchesByHash(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	onMain, err := st.WriteRecordVersion(ctx, store.WriteRecordVersionCmd{
		ProjectID: p.ID, OwningRoot: "api", Slug: "data-modeling", Branch: "main",
		Bytes: []byte("same bytes"), Kind: "edr", DocStatus: "Accepted",
	})
	if err != nil {
		t.Fatalf("WriteRecordVersion(main): %v", err)
	}
	onFeature, err := st.WriteRecordVersion(ctx, store.WriteRecordVersionCmd{
		ProjectID: p.ID, OwningRoot: "api", Slug: "data-modeling", Branch: "feat/x",
		Bytes: []byte("same bytes"), Kind: "edr", DocStatus: "Draft",
	})
	if err != nil {
		t.Fatalf("WriteRecordVersion(feat/x): %v", err)
	}

	if onFeature.ID == onMain.ID {
		t.Fatal("versions share the same row, want distinct versions (versions are never shared across branches)")
	}
	if onFeature.ContentObjectID != onMain.ContentObjectID {
		t.Fatalf("ContentObjectID = %q, want %q (same bytes dedup to one content row)", onFeature.ContentObjectID, onMain.ContentObjectID)
	}
}

// TestSetRecordDeleted_ThenWriteReactivates covers "soft-delete y
// reaparición por rama": the entry is reused, not recreated, and a frozen
// version survives the delete untouched.
func TestSetRecordDeleted_ThenWriteReactivates(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	v1 := writeRecordVersion(ctx, t, st, p.ID, "api", "data-modeling", "main", []byte("v1"))
	rec, err := st.GetRecord(ctx, p.ID, "api", "data-modeling")
	if err != nil {
		t.Fatalf("GetRecord: %v", err)
	}

	deleted, err := st.SetRecordDeleted(ctx, rec.ID, "main")
	if err != nil {
		t.Fatalf("SetRecordDeleted: %v", err)
	}
	if deleted.Status != store.RecordBranchDeleted {
		t.Fatalf("Status = %q, want %q", deleted.Status, store.RecordBranchDeleted)
	}
	if deleted.CurrentVersionID != v1.ID {
		t.Fatalf("CurrentVersionID = %q, want unchanged %q", deleted.CurrentVersionID, v1.ID)
	}

	v2 := writeRecordVersion(ctx, t, st, p.ID, "api", "data-modeling", "main", []byte("v2"))
	reactivated, err := st.GetRecordBranchState(ctx, rec.ID, "main")
	if err != nil {
		t.Fatalf("GetRecordBranchState: %v", err)
	}
	if reactivated.ID != deleted.ID {
		t.Fatalf("reactivated entry ID = %q, want the same entry %q reused", reactivated.ID, deleted.ID)
	}
	if reactivated.Status != store.RecordBranchActive {
		t.Fatalf("Status = %q, want %q", reactivated.Status, store.RecordBranchActive)
	}
	if reactivated.CurrentVersionID != v2.ID {
		t.Fatalf("CurrentVersionID = %q, want the reindexed version %q", reactivated.CurrentVersionID, v2.ID)
	}
}

// TestRecordBranchStateUnique_RawInsertRejected exercises R9's unique index
// directly against the database, bypassing the Repository and the
// generated client (delivery/testing-strategy: constraint rejections are
// exercised with raw SQL). record_branch_state is one of the three tables
// enmienda 7 moved off a composite primary key to a surrogate id +
// UNIQUE(record_id, branch); this is the demonstration that the surrogate
// id does not enable duplicating that natural key.
func TestRecordBranchStateUnique_RawInsertRejected(t *testing.T) {
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

	if _, err := raw.ExecContext(ctx, `INSERT INTO record_branch_state
		(id, created_at, updated_at, branch, status, project_id, record_id, current_version_id)
		VALUES ('branch-state-main', ?, ?, 'main', 'active', ?, ?, ?)`,
		ts, ts, fx.projectID, fx.recordID, fx.recordVersionID); err != nil {
		t.Fatalf("seed record_branch_state on main: %v", err)
	}

	// Negative: a second entry for the same (record_id, branch), distinct id.
	_, err = raw.ExecContext(ctx, `INSERT INTO record_branch_state
		(id, created_at, updated_at, branch, status, project_id, record_id, current_version_id)
		VALUES ('branch-state-main-duplicate', ?, ?, 'main', 'active', ?, ?, ?)`,
		ts, ts, fx.projectID, fx.recordID, fx.recordVersionID)
	if err == nil {
		t.Fatal("expected a UNIQUE(record_id, branch) violation, got nil error")
	}

	// Positive control: the same record_id on a DIFFERENT branch is
	// accepted — this is the per-branch vigency the model needs.
	if _, err := raw.ExecContext(ctx, `INSERT INTO record_branch_state
		(id, created_at, updated_at, branch, status, project_id, record_id, current_version_id)
		VALUES ('branch-state-feat-x', ?, ?, 'feat/x', 'active', ?, ?, ?)`,
		ts, ts, fx.projectID, fx.recordID, fx.recordVersionID); err != nil {
		t.Fatalf("same record_id on a different branch should succeed: %v", err)
	}
}

// TestSetRecordMalformed_KeepsCurrentVersion covers "`.md` malformado
// flaggeado por rama sin perder la última versión sana": the flag flips
// without touching CurrentVersionID.
func TestSetRecordMalformed_KeepsCurrentVersion(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	v := writeRecordVersion(ctx, t, st, p.ID, "api", "data-modeling", "main", []byte("v1"))
	rec, err := st.GetRecord(ctx, p.ID, "api", "data-modeling")
	if err != nil {
		t.Fatalf("GetRecord: %v", err)
	}

	flagged, err := st.SetRecordMalformed(ctx, rec.ID, "main", true)
	if err != nil {
		t.Fatalf("SetRecordMalformed: %v", err)
	}
	if !flagged.Malformed {
		t.Fatal("Malformed = false, want true")
	}
	if flagged.CurrentVersionID != v.ID {
		t.Fatalf("CurrentVersionID = %q, want unchanged %q", flagged.CurrentVersionID, v.ID)
	}
	if flagged.Status != store.RecordBranchActive {
		t.Fatalf("Status = %q, want unaffected %q", flagged.Status, store.RecordBranchActive)
	}
}
