package store_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/franwerner/matecito-ai/apps/api/internal/store"
)

func startChange(ctx context.Context, t *testing.T, st store.Repository, projectID, branch, name string) store.Change {
	t.Helper()
	c, err := st.StartChange(ctx, store.StartChangeCmd{ProjectID: projectID, Branch: branch, Name: name})
	if err != nil {
		t.Fatalf("StartChange(%s): %v", branch, err)
	}
	return c
}

// TestAppendEvent_PositionsAreMonotonicAndGapless covers R6's "posiciones
// consecutivas por change".
func TestAppendEvent_PositionsAreMonotonicAndGapless(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "feat/x", "add-feature")

	var positions []int64
	for i := range 3 {
		ev, err := st.AppendEvent(ctx, store.AppendEventCmd{
			ProjectID: p.ID,
			ChangeID:  c.ID,
			Type:      "submit",
			Payload:   fmt.Sprintf(`{"n":%d}`, i),
		})
		if err != nil {
			t.Fatalf("AppendEvent #%d: %v", i, err)
		}
		positions = append(positions, ev.Position)
	}

	want := []int64{1, 2, 3}
	for i, p := range positions {
		if p != want[i] {
			t.Fatalf("positions = %v, want %v", positions, want)
		}
	}
}

// TestAppendEvent_IndependentSequencesPerChange covers "secuencias
// independientes entre changes": two changes of the same project each keep
// their own position sequence from 1.
func TestAppendEvent_IndependentSequencesPerChange(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	ca := startChange(ctx, t, st, p.ID, "feat/a", "change-a")
	cb := startChange(ctx, t, st, p.ID, "feat/b", "change-b")

	evA, err := st.AppendEvent(ctx, store.AppendEventCmd{ProjectID: p.ID, ChangeID: ca.ID, Type: "submit", Payload: `{"a":1}`})
	if err != nil {
		t.Fatalf("AppendEvent a: %v", err)
	}
	evB, err := st.AppendEvent(ctx, store.AppendEventCmd{ProjectID: p.ID, ChangeID: cb.ID, Type: "submit", Payload: `{"b":1}`})
	if err != nil {
		t.Fatalf("AppendEvent b: %v", err)
	}
	if evA.Position != 1 || evB.Position != 1 {
		t.Fatalf("positions = (%d, %d), want (1, 1) — independent sequences", evA.Position, evB.Position)
	}
}

// TestAppendEvent_IdempotentReappend covers "re-append idempotente": the
// same payload twice returns the same event, no new row and no new
// position consumed.
func TestAppendEvent_IdempotentReappend(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "feat/x", "add-feature")

	cmd := store.AppendEventCmd{ProjectID: p.ID, ChangeID: c.ID, Type: "submit", Payload: `{"anchor":"change/spec"}`}
	first, err := st.AppendEvent(ctx, cmd)
	if err != nil {
		t.Fatalf("first AppendEvent: %v", err)
	}

	second, err := st.AppendEvent(ctx, cmd)
	if err != nil {
		t.Fatalf("second AppendEvent: %v", err)
	}
	if second.ID != first.ID || second.Position != first.Position {
		t.Fatalf("second = {ID: %q, Position: %d}, want the same as first {ID: %q, Position: %d}",
			second.ID, second.Position, first.ID, first.Position)
	}

	// A third, distinct payload must still get the next position — the
	// idempotent replay above must not have consumed one.
	third, err := st.AppendEvent(ctx, store.AppendEventCmd{ProjectID: p.ID, ChangeID: c.ID, Type: "submit", Payload: `{"anchor":"change/design"}`})
	if err != nil {
		t.Fatalf("third AppendEvent: %v", err)
	}
	if third.Position != 2 {
		t.Fatalf("third.Position = %d, want 2 (no position consumed by the idempotent replay)", third.Position)
	}
}

// TestAppendEvent_PayloadHashRoundTrips covers "el payload conserva los
// bytes hasheados": re-hashing the stored payload reproduces the stored
// idempotency_hash, and the payload itself comes back byte-for-byte.
func TestAppendEvent_PayloadHashRoundTrips(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "feat/x", "add-feature")

	const payload = `{"anchor":{"change":"feat/x","phase":"spec"},"body":"héllo\nwörld"}`
	ev, err := st.AppendEvent(ctx, store.AppendEventCmd{ProjectID: p.ID, ChangeID: c.ID, Type: "submit", Payload: payload})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	if ev.Payload != payload {
		t.Fatalf("stored payload = %q, want %q (verbatim, never re-serialized)", ev.Payload, payload)
	}

	events, err := st.ListEvents(ctx, c.ID)
	if err != nil {
		t.Fatalf("ListEvents: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("len(events) = %d, want 1", len(events))
	}
	if events[0].IdempotencyHash != ev.IdempotencyHash {
		t.Fatalf("re-fetched idempotency_hash = %q, want %q", events[0].IdempotencyHash, ev.IdempotencyHash)
	}
}

// TestAppendEvent_AgentLabelRoundTrips covers "la etiqueta del agente vive
// en el evento": the three label fields round-trip on the same row, no
// agent lookup involved.
func TestAppendEvent_AgentLabelRoundTrips(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "feat/x", "add-feature")

	ev, err := st.AppendEvent(ctx, store.AppendEventCmd{
		ProjectID: p.ID,
		ChangeID:  c.ID,
		Type:      "submit",
		Payload:   `{"phase":"apply"}`,
		AgentKind: store.AgentPhase,
		AgentName: "sdd-apply",
		SessionID: "session-123",
	})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	if ev.AgentKind != store.AgentPhase || ev.AgentName != "sdd-apply" || ev.SessionID != "session-123" {
		t.Fatalf("agent label = {%q, %q, %q}, want {%q, %q, %q}",
			ev.AgentKind, ev.AgentName, ev.SessionID, store.AgentPhase, "sdd-apply", "session-123")
	}

	// The three fields are also individually optional — empty is valid.
	bare, err := st.AppendEvent(ctx, store.AppendEventCmd{ProjectID: p.ID, ChangeID: c.ID, Type: "submit", Payload: `{"phase":"bare"}`})
	if err != nil {
		t.Fatalf("AppendEvent (bare): %v", err)
	}
	if bare.AgentKind != "" || bare.AgentName != "" || bare.SessionID != "" {
		t.Fatalf("bare agent label = {%q, %q, %q}, want all empty", bare.AgentKind, bare.AgentName, bare.SessionID)
	}
}

// TestAppendEvent_KeepsGivenOccurredAt covers "el orden semántico es
// occurred_at": a late-arriving event keeps its own, earlier occurred_at
// while still receiving the next log position.
func TestAppendEvent_KeepsGivenOccurredAt(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "feat/x", "add-feature")

	if _, err := st.AppendEvent(ctx, store.AppendEventCmd{ProjectID: p.ID, ChangeID: c.ID, Type: "submit", Payload: `{"n":1}`}); err != nil {
		t.Fatalf("AppendEvent #1: %v", err)
	}

	late := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	ev, err := st.AppendEvent(ctx, store.AppendEventCmd{
		ProjectID:  p.ID,
		ChangeID:   c.ID,
		Type:       "submit",
		Payload:    `{"n":2,"late":true}`,
		OccurredAt: late,
	})
	if err != nil {
		t.Fatalf("AppendEvent (late): %v", err)
	}
	if ev.Position != 2 {
		t.Fatalf("late event Position = %d, want 2 (still the next slot in the log)", ev.Position)
	}
	if !ev.OccurredAt.Equal(late) {
		t.Fatalf("OccurredAt = %v, want %v (its own timestamp, not append time)", ev.OccurredAt, late)
	}
}

// TestListEventsFrom_ResumesAtPosition covers R6's read side: a resume read
// from a given position returns only events at or after it.
func TestListEventsFrom_ResumesAtPosition(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")
	c := startChange(ctx, t, st, p.ID, "feat/x", "add-feature")

	for i := range 3 {
		if _, err := st.AppendEvent(ctx, store.AppendEventCmd{
			ProjectID: p.ID, ChangeID: c.ID, Type: "submit", Payload: fmt.Sprintf(`{"n":%d}`, i),
		}); err != nil {
			t.Fatalf("AppendEvent #%d: %v", i, err)
		}
	}

	fromTwo, err := st.ListEventsFrom(ctx, c.ID, 2)
	if err != nil {
		t.Fatalf("ListEventsFrom(2): %v", err)
	}
	if len(fromTwo) != 2 {
		t.Fatalf("len(fromTwo) = %d, want 2", len(fromTwo))
	}
	if fromTwo[0].Position != 2 || fromTwo[1].Position != 3 {
		t.Fatalf("fromTwo positions = [%d, %d], want [2, 3]", fromTwo[0].Position, fromTwo[1].Position)
	}

	full, err := st.ListEvents(ctx, c.ID)
	if err != nil {
		t.Fatalf("ListEvents: %v", err)
	}
	if len(full) != 3 {
		t.Fatalf("len(full snapshot) = %d, want 3", len(full))
	}
}
