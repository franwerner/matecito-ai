package store

import (
	"context"
	"fmt"

	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent"
	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent/record"
	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent/recordversionref"
)

// GetRecordVersionRef looks up a ref row by id.
func (e *Engine) GetRecordVersionRef(ctx context.Context, id string) (RecordVersionRef, error) {
	row, err := e.client.RecordVersionRef.Get(ctx, id)
	if ent.IsNotFound(err) {
		return RecordVersionRef{}, ErrNotFound
	}
	if err != nil {
		return RecordVersionRef{}, fmt.Errorf("store: get record version ref %s: %w", id, err)
	}
	return toRecordVersionRef(row)
}

// CreateRecordVersionRef adds one relation row to cmd.RecordVersionID.
// Dangling is computed once, here, by checking whether the target
// currently resolves to a record — never supplied by the caller (R9): a
// missing target does not block persisting the ref, it only flags it.
func (e *Engine) CreateRecordVersionRef(ctx context.Context, cmd CreateRecordVersionRefCmd) (RecordVersionRef, error) {
	return submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (RecordVersionRef, error) {
		version, err := tx.RecordVersion.Get(ctx, cmd.RecordVersionID)
		if err != nil {
			if ent.IsNotFound(err) {
				return RecordVersionRef{}, ErrNotFound
			}
			return RecordVersionRef{}, fmt.Errorf("store: get record version %s: %w", cmd.RecordVersionID, err)
		}

		dangling, err := danglingTarget(ctx, tx, version.ProjectID, cmd.TargetOwningRoot, cmd.TargetSlug)
		if err != nil {
			return RecordVersionRef{}, err
		}

		created, err := tx.RecordVersionRef.Create().
			SetProjectID(version.ProjectID).
			SetRecordVersionID(cmd.RecordVersionID).
			SetRelation(cmd.Relation).
			SetTargetSlug(cmd.TargetSlug).
			SetTargetOwningRoot(cmd.TargetOwningRoot).
			SetDangling(boolToInt8(dangling)).
			Save(ctx)
		if err != nil {
			return RecordVersionRef{}, fmt.Errorf("store: create record version ref: %w", err)
		}
		return toRecordVersionRef(created)
	})
}

// RecalculateDanglingRefs recomputes the dangling flag for every ref that
// targets (owningRoot, slug) in projectID. It is the primitive the future
// indexing process (Phase 5, out of Phase 2's scope) calls when it observes
// that target appear or disappear (process/index-decision-records); this
// store does not itself watch the record tree for that change, it only
// exposes the recompute.
func (e *Engine) RecalculateDanglingRefs(ctx context.Context, projectID, owningRoot, slug string) error {
	_, err := submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (struct{}, error) {
		dangling, err := danglingTarget(ctx, tx, projectID, owningRoot, slug)
		if err != nil {
			return struct{}{}, err
		}

		if _, err := tx.RecordVersionRef.Update().
			Where(
				recordversionref.ProjectID(projectID),
				recordversionref.TargetOwningRoot(owningRoot),
				recordversionref.TargetSlug(slug),
			).
			SetDangling(boolToInt8(dangling)).
			Save(ctx); err != nil {
			return struct{}{}, fmt.Errorf("store: recalculate dangling refs for %s/%s: %w", owningRoot, slug, err)
		}
		return struct{}{}, nil
	})
	return err
}

// danglingTarget reports whether (owningRoot, slug) does NOT currently
// resolve to a record in projectID — the store's one definition of
// "dangling", branch-independent like records itself (records has no
// physical delete, so a target that once existed never becomes dangling
// again through this store alone).
func danglingTarget(ctx context.Context, tx *ent.Tx, projectID, owningRoot, slug string) (bool, error) {
	exists, err := tx.Record.Query().
		Where(record.ProjectID(projectID), record.OwningRoot(owningRoot), record.Slug(slug)).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("store: check record target %s/%s exists: %w", owningRoot, slug, err)
	}
	return !exists, nil
}

func boolToInt8(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

func toRecordVersionRef(row *ent.RecordVersionRef) (RecordVersionRef, error) {
	createdAt, err := parseStoreTime(row.CreatedAt)
	if err != nil {
		return RecordVersionRef{}, err
	}
	updatedAt, err := parseStoreTime(row.UpdatedAt)
	if err != nil {
		return RecordVersionRef{}, err
	}
	return RecordVersionRef{
		ID:               row.ID,
		ProjectID:        row.ProjectID,
		RecordVersionID:  row.RecordVersionID,
		Relation:         row.Relation,
		TargetSlug:       row.TargetSlug,
		TargetOwningRoot: row.TargetOwningRoot,
		Dangling:         row.Dangling != 0,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}, nil
}
