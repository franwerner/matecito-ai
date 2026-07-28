package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/franwerner/matecito-ai/apps/api/internal/store"
)

// TestCreateIterationNote_AnchoredToLines covers "nota anclada a líneas de
// un diff".
func TestCreateIterationNote_AnchoredToLines(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "feat/x", "add-feature")
	ev, err := st.AppendEvent(ctx, store.AppendEventCmd{ProjectID: p.ID, ChangeID: c.ID, Type: "apply", Payload: `{"batch":1}`})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}

	note, err := st.CreateIterationNote(ctx, store.CreateIterationNoteCmd{
		ProjectID: p.ID,
		ChangeID:  c.ID,
		EventID:   ev.ID,
		FilePath:  "internal/store/events.go",
		LineFrom:  10,
		LineTo:    20,
		Body:      "revisar el manejo de la posición",
	})
	if err != nil {
		t.Fatalf("CreateIterationNote: %v", err)
	}
	if note.Status != store.NotePending {
		t.Fatalf("Status = %q, want %q", note.Status, store.NotePending)
	}
	if note.EventID != ev.ID || note.FilePath != "internal/store/events.go" || note.LineFrom != 10 || note.LineTo != 20 {
		t.Fatalf("anchor = {event:%q file:%q from:%d to:%d}, want {%q, internal/store/events.go, 10, 20}",
			note.EventID, note.FilePath, note.LineFrom, note.LineTo, ev.ID)
	}
}

// TestCreateIterationNote_GeneralOnChange covers "nota general del change":
// no event, no file, no range.
func TestCreateIterationNote_GeneralOnChange(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "feat/x", "add-feature")

	note, err := st.CreateIterationNote(ctx, store.CreateIterationNoteCmd{
		ProjectID: p.ID,
		ChangeID:  c.ID,
		Body:      "revisar el enfoque general",
	})
	if err != nil {
		t.Fatalf("CreateIterationNote: %v", err)
	}
	if note.Status != store.NotePending {
		t.Fatalf("Status = %q, want %q", note.Status, store.NotePending)
	}
	if note.EventID != "" || note.FilePath != "" || note.LineFrom != 0 || note.LineTo != 0 {
		t.Fatalf("expected an unanchored note, got {event:%q file:%q from:%d to:%d}",
			note.EventID, note.FilePath, note.LineFrom, note.LineTo)
	}
}

// TestIterationNoteLifecycle_DeliverThenResolve covers "lectura confirmada
// entrega las notas" and "resolución trazada": pending -> delivered ->
// resolved, each stamped.
func TestIterationNoteLifecycle_DeliverThenResolve(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "feat/x", "add-feature")
	note, err := st.CreateIterationNote(ctx, store.CreateIterationNoteCmd{ProjectID: p.ID, ChangeID: c.ID, Body: "iterar"})
	if err != nil {
		t.Fatalf("CreateIterationNote: %v", err)
	}

	delivered, err := st.DeliverIterationNote(ctx, note.ID)
	if err != nil {
		t.Fatalf("DeliverIterationNote: %v", err)
	}
	if delivered.Status != store.NoteDelivered {
		t.Fatalf("Status = %q, want %q", delivered.Status, store.NoteDelivered)
	}
	if delivered.DeliveredAt.IsZero() {
		t.Fatal("DeliveredAt is zero, want it stamped")
	}

	resolved, err := st.ResolveIterationNote(ctx, note.ID)
	if err != nil {
		t.Fatalf("ResolveIterationNote: %v", err)
	}
	if resolved.Status != store.NoteResolved {
		t.Fatalf("Status = %q, want %q", resolved.Status, store.NoteResolved)
	}
	if resolved.ResolvedAt.IsZero() {
		t.Fatal("ResolvedAt is zero, want it stamped")
	}

	got, err := st.GetIterationNote(ctx, note.ID)
	if err != nil {
		t.Fatalf("GetIterationNote: %v", err)
	}
	if got.Status != store.NoteResolved {
		t.Fatalf("re-fetched Status = %q, want %q", got.Status, store.NoteResolved)
	}
}

// TestIterationNote_DeliveredCannotGoBackToPending covers R10's "transiciones
// sin vuelta atrás": a delivered note rejects being moved back to pending —
// there is no such method, so this is exercised as "delivering it again"
// being refused once it is no longer pending.
func TestIterationNote_DeliveredCannotGoBackToPending(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "feat/x", "add-feature")
	note, err := st.CreateIterationNote(ctx, store.CreateIterationNoteCmd{ProjectID: p.ID, ChangeID: c.ID, Body: "iterar"})
	if err != nil {
		t.Fatalf("CreateIterationNote: %v", err)
	}
	if _, err := st.DeliverIterationNote(ctx, note.ID); err != nil {
		t.Fatalf("DeliverIterationNote: %v", err)
	}

	_, err = st.DeliverIterationNote(ctx, note.ID)
	if !errors.Is(err, store.ErrInvalidNoteTransition) {
		t.Fatalf("second DeliverIterationNote err = %v, want ErrInvalidNoteTransition", err)
	}

	unchanged, err := st.GetIterationNote(ctx, note.ID)
	if err != nil {
		t.Fatalf("GetIterationNote: %v", err)
	}
	if unchanged.Status != store.NoteDelivered {
		t.Fatalf("Status = %q, want unchanged %q", unchanged.Status, store.NoteDelivered)
	}
}

// TestIterationNote_ResolveRequiresDelivered covers the same "sin vuelta
// atrás"/no-skip rule from the other direction: resolving a still-pending
// note is refused.
func TestIterationNote_ResolveRequiresDelivered(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "feat/x", "add-feature")
	note, err := st.CreateIterationNote(ctx, store.CreateIterationNoteCmd{ProjectID: p.ID, ChangeID: c.ID, Body: "iterar"})
	if err != nil {
		t.Fatalf("CreateIterationNote: %v", err)
	}

	_, err = st.ResolveIterationNote(ctx, note.ID)
	if !errors.Is(err, store.ErrInvalidNoteTransition) {
		t.Fatalf("err = %v, want ErrInvalidNoteTransition", err)
	}
}

// TestDeleteIterationNote_PendingIsPhysicallyDeleted covers "la nota
// pendiente se borra de verdad".
func TestDeleteIterationNote_PendingIsPhysicallyDeleted(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "feat/x", "add-feature")
	note, err := st.CreateIterationNote(ctx, store.CreateIterationNoteCmd{ProjectID: p.ID, ChangeID: c.ID, Body: "descartar"})
	if err != nil {
		t.Fatalf("CreateIterationNote: %v", err)
	}

	if err := st.DeleteIterationNote(ctx, note.ID); err != nil {
		t.Fatalf("DeleteIterationNote: %v", err)
	}

	_, err = st.GetIterationNote(ctx, note.ID)
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("GetIterationNote after delete err = %v, want ErrNotFound", err)
	}
}

// TestDeleteIterationNote_DeliveredOrResolvedNotDeletable covers "la nota
// entregada o resuelta no se borra": R11's only exception is `pending`.
func TestDeleteIterationNote_DeliveredOrResolvedNotDeletable(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "feat/x", "add-feature")

	delivered, err := st.CreateIterationNote(ctx, store.CreateIterationNoteCmd{ProjectID: p.ID, ChangeID: c.ID, Body: "a"})
	if err != nil {
		t.Fatalf("CreateIterationNote: %v", err)
	}
	if _, err := st.DeliverIterationNote(ctx, delivered.ID); err != nil {
		t.Fatalf("DeliverIterationNote: %v", err)
	}
	if err := st.DeleteIterationNote(ctx, delivered.ID); !errors.Is(err, store.ErrNoteNotDeletable) {
		t.Fatalf("delete delivered note err = %v, want ErrNoteNotDeletable", err)
	}

	resolved, err := st.CreateIterationNote(ctx, store.CreateIterationNoteCmd{ProjectID: p.ID, ChangeID: c.ID, Body: "b"})
	if err != nil {
		t.Fatalf("CreateIterationNote: %v", err)
	}
	if _, err := st.DeliverIterationNote(ctx, resolved.ID); err != nil {
		t.Fatalf("DeliverIterationNote: %v", err)
	}
	if _, err := st.ResolveIterationNote(ctx, resolved.ID); err != nil {
		t.Fatalf("ResolveIterationNote: %v", err)
	}
	if err := st.DeleteIterationNote(ctx, resolved.ID); !errors.Is(err, store.ErrNoteNotDeletable) {
		t.Fatalf("delete resolved note err = %v, want ErrNoteNotDeletable", err)
	}
}

// TestGetIterationNote_NotFound covers the Reader's ErrNotFound contract.
func TestGetIterationNote_NotFound(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)

	_, err := st.GetIterationNote(ctx, "no-such-note")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
