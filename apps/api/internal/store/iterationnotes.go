package store

import (
	"context"
	"fmt"
	"time"

	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent"
)

// CreateIterationNote anchors a new note in `pending` (flow/iteration-notes).
// The anchoring shape itself — a range needs a file, a file needs an event —
// is a database CHECK (R10); the store passes cmd through as given and lets
// the constraint reject an invalid combination.
func (e *Engine) CreateIterationNote(ctx context.Context, cmd CreateIterationNoteCmd) (IterationNote, error) {
	return submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (IterationNote, error) {
		create := tx.IterationNote.Create().
			SetProjectID(cmd.ProjectID).
			SetChangeID(cmd.ChangeID).
			SetBody(cmd.Body).
			SetStatus(string(NotePending))
		if cmd.EventID != "" {
			create = create.SetEventID(cmd.EventID)
		}
		if cmd.FilePath != "" {
			create = create.SetFilePath(cmd.FilePath)
		}
		if cmd.LineFrom != 0 {
			create = create.SetLineFrom(cmd.LineFrom)
		}
		if cmd.LineTo != 0 {
			create = create.SetLineTo(cmd.LineTo)
		}

		created, err := create.Save(ctx)
		if err != nil {
			return IterationNote{}, fmt.Errorf("store: create iteration note: %w", err)
		}
		return toIterationNote(created)
	})
}

// GetIterationNote looks up a note by id.
func (e *Engine) GetIterationNote(ctx context.Context, id string) (IterationNote, error) {
	row, err := e.client.IterationNote.Get(ctx, id)
	if ent.IsNotFound(err) {
		return IterationNote{}, ErrNotFound
	}
	if err != nil {
		return IterationNote{}, fmt.Errorf("store: get iteration note %s: %w", id, err)
	}
	return toIterationNote(row)
}

// DeliverIterationNote moves a note from `pending` to `delivered` — the AI
// read it via the query tool, with the user's confirmation
// (flow/iteration-notes).
func (e *Engine) DeliverIterationNote(ctx context.Context, id string) (IterationNote, error) {
	return e.transitionIterationNote(ctx, id, NotePending, NoteDelivered,
		func(u *ent.IterationNoteUpdateOne, now string) *ent.IterationNoteUpdateOne {
			return u.SetDeliveredAt(now)
		})
}

// ResolveIterationNote moves a note from `delivered` to `resolved` — the AI
// marked it resolved via the resolution tool, traced in the event-log
// (flow/iteration-notes).
func (e *Engine) ResolveIterationNote(ctx context.Context, id string) (IterationNote, error) {
	return e.transitionIterationNote(ctx, id, NoteDelivered, NoteResolved,
		func(u *ent.IterationNoteUpdateOne, now string) *ent.IterationNoteUpdateOne {
			return u.SetResolvedAt(now)
		})
}

// transitionIterationNote reads the note's current status inside the write
// transaction and rejects the move unless it starts from `from` — R10's
// value-dependent transitions are enforced here, at the Repository's write
// edge, not by a trigger (Batch 3 decision: Ent does not declare triggers,
// and R1 forbids hand-editing the derived migrations to add one).
func (e *Engine) transitionIterationNote(
	ctx context.Context,
	id string,
	from, to IterationNoteStatus,
	stamp func(*ent.IterationNoteUpdateOne, string) *ent.IterationNoteUpdateOne,
) (IterationNote, error) {
	return submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (IterationNote, error) {
		current, err := tx.IterationNote.Get(ctx, id)
		if ent.IsNotFound(err) {
			return IterationNote{}, ErrNotFound
		}
		if err != nil {
			return IterationNote{}, fmt.Errorf("store: get iteration note %s: %w", id, err)
		}
		if IterationNoteStatus(current.Status) != from {
			return IterationNote{}, fmt.Errorf("%w: note %s is %q, not %q", ErrInvalidNoteTransition, id, current.Status, from)
		}

		update := stamp(tx.IterationNote.UpdateOneID(id).SetStatus(string(to)), nowStoreTime())
		updated, err := update.Save(ctx)
		if err != nil {
			return IterationNote{}, fmt.Errorf("store: transition iteration note %s: %w", id, err)
		}
		return toIterationNote(updated)
	})
}

// DeleteIterationNote physically removes a note — the only physical delete
// this store exposes (R11), and only while it is still `pending`: a
// `delivered` or `resolved` note is history and stays.
func (e *Engine) DeleteIterationNote(ctx context.Context, id string) error {
	_, err := submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (struct{}, error) {
		current, err := tx.IterationNote.Get(ctx, id)
		if ent.IsNotFound(err) {
			return struct{}{}, ErrNotFound
		}
		if err != nil {
			return struct{}{}, fmt.Errorf("store: get iteration note %s: %w", id, err)
		}
		if IterationNoteStatus(current.Status) != NotePending {
			return struct{}{}, fmt.Errorf("%w: note %s is %q", ErrNoteNotDeletable, id, current.Status)
		}

		if err := tx.IterationNote.DeleteOneID(id).Exec(ctx); err != nil {
			return struct{}{}, fmt.Errorf("store: delete iteration note %s: %w", id, err)
		}
		return struct{}{}, nil
	})
	return err
}

func toIterationNote(row *ent.IterationNote) (IterationNote, error) {
	createdAt, err := parseStoreTime(row.CreatedAt)
	if err != nil {
		return IterationNote{}, err
	}
	updatedAt, err := parseStoreTime(row.UpdatedAt)
	if err != nil {
		return IterationNote{}, err
	}

	var deliveredAt time.Time
	if row.DeliveredAt != "" {
		if deliveredAt, err = parseStoreTime(row.DeliveredAt); err != nil {
			return IterationNote{}, err
		}
	}
	var resolvedAt time.Time
	if row.ResolvedAt != "" {
		if resolvedAt, err = parseStoreTime(row.ResolvedAt); err != nil {
			return IterationNote{}, err
		}
	}

	return IterationNote{
		ID:          row.ID,
		ProjectID:   row.ProjectID,
		ChangeID:    row.ChangeID,
		EventID:     row.EventID,
		FilePath:    row.FilePath,
		LineFrom:    row.LineFrom,
		LineTo:      row.LineTo,
		Body:        row.Body,
		Status:      IterationNoteStatus(row.Status),
		DeliveredAt: deliveredAt,
		ResolvedAt:  resolvedAt,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}
