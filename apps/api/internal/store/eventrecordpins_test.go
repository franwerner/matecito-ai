package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/franwerner/matecito-ai/apps/api/internal/store"
)

// TestPinEventToRecordVersion_FreezesCurrentVersion covers "el pin congela
// la versión actual": the pin resolves the branch's current version
// implicitly and freezes it.
func TestPinEventToRecordVersion_FreezesCurrentVersion(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "main", "add-feature")

	v := writeRecordVersion(ctx, t, st, p.ID, "api", "data-modeling", "main", []byte("v1"))
	rec, err := st.GetRecord(ctx, p.ID, "api", "data-modeling")
	if err != nil {
		t.Fatalf("GetRecord: %v", err)
	}
	ev, err := st.AppendEvent(ctx, store.AppendEventCmd{ProjectID: p.ID, ChangeID: c.ID, Type: "submit", Payload: `{"n":1}`})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}

	pin, err := st.PinEventToRecordVersion(ctx, store.PinEventToRecordVersionCmd{
		ProjectID: p.ID, EventID: ev.ID, RecordID: rec.ID, Branch: "main",
	})
	if err != nil {
		t.Fatalf("PinEventToRecordVersion: %v", err)
	}
	if pin.RecordVersionID != v.ID {
		t.Fatalf("pin.RecordVersionID = %q, want the branch's current version %q", pin.RecordVersionID, v.ID)
	}

	frozen, err := st.GetRecordVersion(ctx, v.ID)
	if err != nil {
		t.Fatalf("GetRecordVersion: %v", err)
	}
	if frozen.Status != store.VersionFrozen {
		t.Fatalf("Status = %q, want %q", frozen.Status, store.VersionFrozen)
	}
	if frozen.FrozenAt.IsZero() {
		t.Fatal("FrozenAt is zero, want it stamped")
	}
}

// TestPinEventToRecordVersion_IdempotentOnRepeat covers R8's "un evento
// MUST NOT pinear dos veces la misma versión": repeating the same
// (event, version) pair is idempotent, following this package's established
// idiom for every other natural-key write (RegisterProject, StartChange,
// AppendEvent, PutContentObject) instead of a new conflict error.
func TestPinEventToRecordVersion_IdempotentOnRepeat(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "main", "add-feature")

	writeRecordVersion(ctx, t, st, p.ID, "api", "data-modeling", "main", []byte("v1"))
	rec, err := st.GetRecord(ctx, p.ID, "api", "data-modeling")
	if err != nil {
		t.Fatalf("GetRecord: %v", err)
	}
	ev, err := st.AppendEvent(ctx, store.AppendEventCmd{ProjectID: p.ID, ChangeID: c.ID, Type: "submit", Payload: `{"n":1}`})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}

	cmd := store.PinEventToRecordVersionCmd{ProjectID: p.ID, EventID: ev.ID, RecordID: rec.ID, Branch: "main"}
	first, err := st.PinEventToRecordVersion(ctx, cmd)
	if err != nil {
		t.Fatalf("first PinEventToRecordVersion: %v", err)
	}
	second, err := st.PinEventToRecordVersion(ctx, cmd)
	if err != nil {
		t.Fatalf("second PinEventToRecordVersion: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("second.ID = %q, want %q (idempotent, not a duplicate)", second.ID, first.ID)
	}
}

// TestEventRecordPinUnique_RawInsertRejected exercises R8's unique index
// directly against the database, bypassing the Repository and the
// generated client (delivery/testing-strategy: constraint rejections are
// exercised with raw SQL). event_record_pins is the third and last table
// enmienda 7 moved off a composite primary key to a surrogate id +
// UNIQUE(event_id, record_version_id); this is the demonstration —
// alongside content_objects and record_branch_state — that the surrogate
// id does not enable duplicating that natural key. The only existing pin
// test (TestPinEventToRecordVersion_IdempotentOnRepeat) goes through the
// Repository's get-or-create path, which cannot tell a real index apart
// from a deleted one — this test is what actually proves the base rejects
// the duplicate.
func TestEventRecordPinUnique_RawInsertRejected(t *testing.T) {
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

	if _, err := raw.ExecContext(ctx, `INSERT INTO event_record_pins
		(id, created_at, updated_at, project_id, event_id, record_version_id)
		VALUES ('pin-1', ?, ?, ?, ?, ?)`,
		ts, ts, fx.projectID, fx.eventID, fx.recordVersionID); err != nil {
		t.Fatalf("seed event_record_pins: %v", err)
	}

	// Negative: a second pin for the same (event_id, record_version_id),
	// distinct id.
	_, err = raw.ExecContext(ctx, `INSERT INTO event_record_pins
		(id, created_at, updated_at, project_id, event_id, record_version_id)
		VALUES ('pin-1-duplicate', ?, ?, ?, ?, ?)`,
		ts, ts, fx.projectID, fx.eventID, fx.recordVersionID)
	if err == nil {
		t.Fatal("expected a UNIQUE(event_id, record_version_id) violation, got nil error")
	}

	// Positive control: the same version pinned by a DIFFERENT event is
	// accepted — what's forbidden is the repeated pair, not the event or
	// the version on its own.
	if _, err := raw.ExecContext(ctx, `INSERT INTO events
		(id, created_at, updated_at, position, type, payload, occurred_at, idempotency_hash, project_id, change_id)
		VALUES ('event-2', ?, ?, 2, 'submit', '{}', ?, 'hash-2', ?, ?)`,
		ts, ts, ts, fx.projectID, fx.changeID); err != nil {
		t.Fatalf("seed second event: %v", err)
	}
	if _, err := raw.ExecContext(ctx, `INSERT INTO event_record_pins
		(id, created_at, updated_at, project_id, event_id, record_version_id)
		VALUES ('pin-2', ?, ?, ?, 'event-2', ?)`,
		ts, ts, fx.projectID, fx.recordVersionID); err != nil {
		t.Fatalf("same version pinned by a different event should succeed: %v", err)
	}
}

// TestPinEventToRecordVersion_NotFoundWhenBranchNeverIndexed covers the
// precondition: pinning against a (record, branch) with no indexed state
// yet has nothing to resolve.
func TestPinEventToRecordVersion_NotFoundWhenBranchNeverIndexed(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "main", "add-feature")

	ev, err := st.AppendEvent(ctx, store.AppendEventCmd{ProjectID: p.ID, ChangeID: c.ID, Type: "submit", Payload: `{"n":1}`})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}

	_, err = st.PinEventToRecordVersion(ctx, store.PinEventToRecordVersionCmd{
		ProjectID: p.ID, EventID: ev.ID, RecordID: "no-such-record", Branch: "main",
	})
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
