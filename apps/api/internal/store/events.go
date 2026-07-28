package store

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent"
	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent/event"
)

// AppendEvent appends one event to (cmd.ProjectID, cmd.ChangeID)'s log
// (submit-phase-artifact, R6). Idempotent by idempotency_hash — computed by
// contentHash over cmd.Payload's exact bytes: a retry with the same payload
// returns the event already on the log, unchanged, instead of consuming a
// new position.
func (e *Engine) AppendEvent(ctx context.Context, cmd AppendEventCmd) (Event, error) {
	return submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (Event, error) {
		hash := contentHash([]byte(cmd.Payload))

		existing, err := tx.Event.Query().Where(event.IdempotencyHash(hash)).Only(ctx)
		switch {
		case err == nil:
			return toEvent(existing)
		case !ent.IsNotFound(err):
			return Event{}, fmt.Errorf("store: lookup event by idempotency hash: %w", err)
		}

		position, err := nextEventPosition(ctx, tx, cmd.ChangeID)
		if err != nil {
			return Event{}, err
		}

		occurredAt := cmd.OccurredAt
		if occurredAt.IsZero() {
			occurredAt = time.Now()
		}

		create := tx.Event.Create().
			SetProjectID(cmd.ProjectID).
			SetChangeID(cmd.ChangeID).
			SetPosition(position).
			SetType(cmd.Type).
			SetPayload(cmd.Payload).
			SetOccurredAt(formatStoreTime(occurredAt)).
			SetIdempotencyHash(hash)
		if cmd.AgentKind != "" {
			create = create.SetAgentKind(string(cmd.AgentKind))
		}
		if cmd.AgentName != "" {
			create = create.SetAgentName(cmd.AgentName)
		}
		if cmd.SessionID != "" {
			create = create.SetSessionID(cmd.SessionID)
		}

		created, err := create.Save(ctx)
		if err != nil {
			return Event{}, fmt.Errorf("store: append event: %w", err)
		}
		return toEvent(created)
	})
}

// nextEventPosition returns changeID's next monotonic position: 1 for a
// change with no events yet, otherwise one past the highest position
// already on its log. Safe without its own locking because it always runs
// inside a submit'd transaction on the single writer goroutine — no other
// write can be in flight concurrently (runtime/concurrency-async).
func nextEventPosition(ctx context.Context, tx *ent.Tx, changeID string) (int64, error) {
	last, err := tx.Event.Query().
		Where(event.ChangeID(changeID)).
		Order(event.ByPosition(sql.OrderDesc())).
		First(ctx)
	switch {
	case ent.IsNotFound(err):
		return 1, nil
	case err != nil:
		return 0, fmt.Errorf("store: find last event position for change %s: %w", changeID, err)
	}
	return last.Position + 1, nil
}

// ListEvents returns changeID's full event-log in position order — the
// resume snapshot (submit-phase-artifact).
func (e *Engine) ListEvents(ctx context.Context, changeID string) ([]Event, error) {
	return e.listEventsFrom(ctx, changeID, 0)
}

// ListEventsFrom returns changeID's event-log starting at fromPosition
// (inclusive) — the cursor-resume read (R6 "lectura desde una posición
// dada").
func (e *Engine) ListEventsFrom(ctx context.Context, changeID string, fromPosition int64) ([]Event, error) {
	return e.listEventsFrom(ctx, changeID, fromPosition)
}

func (e *Engine) listEventsFrom(ctx context.Context, changeID string, fromPosition int64) ([]Event, error) {
	rows, err := e.client.Event.Query().
		Where(event.ChangeID(changeID), event.PositionGTE(fromPosition)).
		Order(event.ByPosition(sql.OrderAsc())).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("store: list events for change %s from position %d: %w", changeID, fromPosition, err)
	}

	events := make([]Event, 0, len(rows))
	for _, row := range rows {
		ev, err := toEvent(row)
		if err != nil {
			return nil, err
		}
		events = append(events, ev)
	}
	return events, nil
}

func toEvent(row *ent.Event) (Event, error) {
	occurredAt, err := parseStoreTime(row.OccurredAt)
	if err != nil {
		return Event{}, err
	}
	createdAt, err := parseStoreTime(row.CreatedAt)
	if err != nil {
		return Event{}, err
	}
	updatedAt, err := parseStoreTime(row.UpdatedAt)
	if err != nil {
		return Event{}, err
	}
	return Event{
		ID:              row.ID,
		ProjectID:       row.ProjectID,
		ChangeID:        row.ChangeID,
		Position:        row.Position,
		Type:            row.Type,
		Payload:         row.Payload,
		OccurredAt:      occurredAt,
		AgentKind:       AgentKind(row.AgentKind),
		AgentName:       row.AgentName,
		SessionID:       row.SessionID,
		IdempotencyHash: row.IdempotencyHash,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}
