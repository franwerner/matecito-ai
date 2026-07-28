package store

import (
	"context"
	"fmt"
	"time"

	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent"
	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent/record"
	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent/recordbranchstate"
)

// GetRecord looks up a record by its natural key: slug + owning_root
// (process/index-decision-records), branch-independent.
func (e *Engine) GetRecord(ctx context.Context, projectID, owningRoot, slug string) (Record, error) {
	row, err := e.client.Record.Query().
		Where(record.ProjectID(projectID), record.OwningRoot(owningRoot), record.Slug(slug)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return Record{}, ErrNotFound
	}
	if err != nil {
		return Record{}, fmt.Errorf("store: get record %s/%s: %w", owningRoot, slug, err)
	}
	return toRecord(row)
}

// UpsertRecord resolves or creates the record for (cmd.ProjectID,
// cmd.OwningRoot, cmd.Slug) — the same identity always resolves to the same
// row, never a duplicate.
func (e *Engine) UpsertRecord(ctx context.Context, cmd UpsertRecordCmd) (Record, error) {
	return submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (Record, error) {
		return upsertRecord(ctx, tx, cmd.ProjectID, cmd.OwningRoot, cmd.Slug)
	})
}

func upsertRecord(ctx context.Context, tx *ent.Tx, projectID, owningRoot, slug string) (Record, error) {
	existing, err := tx.Record.Query().
		Where(record.ProjectID(projectID), record.OwningRoot(owningRoot), record.Slug(slug)).
		Only(ctx)
	switch {
	case err == nil:
		return toRecord(existing)
	case !ent.IsNotFound(err):
		return Record{}, fmt.Errorf("store: lookup record %s/%s: %w", owningRoot, slug, err)
	}

	created, err := tx.Record.Create().
		SetProjectID(projectID).
		SetOwningRoot(owningRoot).
		SetSlug(slug).
		Save(ctx)
	if err != nil {
		return Record{}, fmt.Errorf("store: create record %s/%s: %w", owningRoot, slug, err)
	}
	return toRecord(created)
}

func toRecord(row *ent.Record) (Record, error) {
	createdAt, err := parseStoreTime(row.CreatedAt)
	if err != nil {
		return Record{}, err
	}
	updatedAt, err := parseStoreTime(row.UpdatedAt)
	if err != nil {
		return Record{}, err
	}
	return Record{
		ID:         row.ID,
		ProjectID:  row.ProjectID,
		OwningRoot: row.OwningRoot,
		Slug:       row.Slug,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

// GetRecordVersion looks up a version by id.
func (e *Engine) GetRecordVersion(ctx context.Context, id string) (RecordVersion, error) {
	row, err := e.client.RecordVersion.Get(ctx, id)
	if ent.IsNotFound(err) {
		return RecordVersion{}, ErrNotFound
	}
	if err != nil {
		return RecordVersion{}, fmt.Errorf("store: get record version %s: %w", id, err)
	}
	return toRecordVersion(row)
}

// WriteRecordVersion applies cmd's content to (cmd.OwningRoot, cmd.Slug) on
// cmd.Branch — the copy-on-write lazy at the heart of
// lifecycle/record-version. Update-in-place proceeds unless the branch's
// current version is frozen (pinned by an event), in which case it forks a
// new current version for this branch instead. Versions are never shared
// across branches — only the content they reference is deduplicated by
// hash.
func (e *Engine) WriteRecordVersion(ctx context.Context, cmd WriteRecordVersionCmd) (RecordVersion, error) {
	return submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (RecordVersion, error) {
		rec, err := upsertRecord(ctx, tx, cmd.ProjectID, cmd.OwningRoot, cmd.Slug)
		if err != nil {
			return RecordVersion{}, err
		}

		content, err := putContentObjectTx(ctx, tx, cmd.ProjectID, cmd.Bytes)
		if err != nil {
			return RecordVersion{}, err
		}

		branchState, err := tx.RecordBranchState.Query().
			Where(recordbranchstate.RecordID(rec.ID), recordbranchstate.Branch(cmd.Branch)).
			Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return RecordVersion{}, fmt.Errorf("store: lookup record branch state %s@%s: %w", rec.ID, cmd.Branch, err)
		}
		// branchState is nil here exactly when the lookup above returned
		// ent.IsNotFound — this branch has never indexed this record before.

		version, err := resolveWrittenVersion(ctx, tx, cmd, rec.ID, content.ID, branchState)
		if err != nil {
			return RecordVersion{}, err
		}

		if branchState == nil {
			if _, err := tx.RecordBranchState.Create().
				SetProjectID(cmd.ProjectID).
				SetRecordID(rec.ID).
				SetBranch(cmd.Branch).
				SetStatus(string(RecordBranchActive)).
				SetCurrentVersionID(version.ID).
				Save(ctx); err != nil {
				return RecordVersion{}, fmt.Errorf("store: create record branch state %s@%s: %w", rec.ID, cmd.Branch, err)
			}
		} else if err := tx.RecordBranchState.UpdateOneID(branchState.ID).
			SetStatus(string(RecordBranchActive)).
			SetCurrentVersionID(version.ID).
			SetMalformed(0).
			Exec(ctx); err != nil {
			return RecordVersion{}, fmt.Errorf("store: update record branch state %s@%s: %w", rec.ID, cmd.Branch, err)
		}

		return toRecordVersion(version)
	})
}

// resolveWrittenVersion decides update-in-place vs fork for a
// WriteRecordVersion call. branchState is nil when cmd.Branch has never
// indexed recordID before: there is nothing on this branch to fork from, so
// this branch's first version is always created outright — versions are
// never shared across branches (lifecycle/record-version), regardless of
// whether some other branch already has identical content.
func resolveWrittenVersion(
	ctx context.Context,
	tx *ent.Tx,
	cmd WriteRecordVersionCmd,
	recordID, contentObjectID string,
	branchState *ent.RecordBranchState,
) (*ent.RecordVersion, error) {
	if branchState == nil {
		return createRecordVersion(ctx, tx, cmd.ProjectID, recordID, contentObjectID, cmd.Kind, cmd.DocStatus)
	}

	current, err := tx.RecordVersion.Get(ctx, branchState.CurrentVersionID)
	if err != nil {
		return nil, fmt.Errorf("store: get current record version %s: %w", branchState.CurrentVersionID, err)
	}

	// lifecycle/record-version's one fork condition: frozen (pinned by an
	// event). Forks a new version for this branch instead of touching the
	// frozen row, which stays exactly as it was for whoever references it.
	if RecordVersionStatus(current.Status) == VersionFrozen {
		return createRecordVersion(ctx, tx, cmd.ProjectID, recordID, contentObjectID, cmd.Kind, cmd.DocStatus)
	}

	updated, err := tx.RecordVersion.UpdateOneID(current.ID).
		SetContentObjectID(contentObjectID).
		SetKind(cmd.Kind).
		SetDocStatus(cmd.DocStatus).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("store: update record version %s in place: %w", current.ID, err)
	}
	return updated, nil
}

func createRecordVersion(ctx context.Context, tx *ent.Tx, projectID, recordID, contentObjectID, kind, docStatus string) (*ent.RecordVersion, error) {
	created, err := tx.RecordVersion.Create().
		SetProjectID(projectID).
		SetRecordID(recordID).
		SetContentObjectID(contentObjectID).
		SetStatus(string(VersionCurrent)).
		SetKind(kind).
		SetDocStatus(docStatus).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("store: create record version for record %s: %w", recordID, err)
	}
	return created, nil
}

func toRecordVersion(row *ent.RecordVersion) (RecordVersion, error) {
	createdAt, err := parseStoreTime(row.CreatedAt)
	if err != nil {
		return RecordVersion{}, err
	}
	updatedAt, err := parseStoreTime(row.UpdatedAt)
	if err != nil {
		return RecordVersion{}, err
	}

	var frozenAt time.Time
	if row.FrozenAt != "" {
		if frozenAt, err = parseStoreTime(row.FrozenAt); err != nil {
			return RecordVersion{}, err
		}
	}

	return RecordVersion{
		ID:              row.ID,
		ProjectID:       row.ProjectID,
		RecordID:        row.RecordID,
		ContentObjectID: row.ContentObjectID,
		Status:          RecordVersionStatus(row.Status),
		Kind:            row.Kind,
		DocStatus:       row.DocStatus,
		FrozenAt:        frozenAt,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}

// GetRecordBranchState looks up recordID's vigency entry on branch.
func (e *Engine) GetRecordBranchState(ctx context.Context, recordID, branch string) (RecordBranchState, error) {
	row, err := e.client.RecordBranchState.Query().
		Where(recordbranchstate.RecordID(recordID), recordbranchstate.Branch(branch)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return RecordBranchState{}, ErrNotFound
	}
	if err != nil {
		return RecordBranchState{}, fmt.Errorf("store: get record branch state %s@%s: %w", recordID, branch, err)
	}
	return toRecordBranchState(row)
}

// SetRecordDeleted soft-deletes recordID's entry on branch — its `.md`
// disappeared from that branch's tree. The entry stays and
// CurrentVersionID is untouched: any frozen version it points to keeps
// resolving (R9, R11).
func (e *Engine) SetRecordDeleted(ctx context.Context, recordID, branch string) (RecordBranchState, error) {
	return submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (RecordBranchState, error) {
		existing, err := tx.RecordBranchState.Query().
			Where(recordbranchstate.RecordID(recordID), recordbranchstate.Branch(branch)).
			Only(ctx)
		if ent.IsNotFound(err) {
			return RecordBranchState{}, ErrNotFound
		}
		if err != nil {
			return RecordBranchState{}, fmt.Errorf("store: lookup record branch state %s@%s: %w", recordID, branch, err)
		}

		updated, err := tx.RecordBranchState.UpdateOneID(existing.ID).
			SetStatus(string(RecordBranchDeleted)).
			Save(ctx)
		if err != nil {
			return RecordBranchState{}, fmt.Errorf("store: soft-delete record branch state %s@%s: %w", recordID, branch, err)
		}
		return toRecordBranchState(updated)
	})
}

// SetRecordMalformed flags recordID's `.md` as unparseable on branch,
// without touching CurrentVersionID — the entry keeps resolving to the
// last version that did parse (R9: "flaggeado ... sin perder la última
// versión sana").
func (e *Engine) SetRecordMalformed(ctx context.Context, recordID, branch string, malformed bool) (RecordBranchState, error) {
	return submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (RecordBranchState, error) {
		existing, err := tx.RecordBranchState.Query().
			Where(recordbranchstate.RecordID(recordID), recordbranchstate.Branch(branch)).
			Only(ctx)
		if ent.IsNotFound(err) {
			return RecordBranchState{}, ErrNotFound
		}
		if err != nil {
			return RecordBranchState{}, fmt.Errorf("store: lookup record branch state %s@%s: %w", recordID, branch, err)
		}

		updated, err := tx.RecordBranchState.UpdateOneID(existing.ID).
			SetMalformed(boolToInt8(malformed)).
			Save(ctx)
		if err != nil {
			return RecordBranchState{}, fmt.Errorf("store: set record branch state %s@%s malformed: %w", recordID, branch, err)
		}
		return toRecordBranchState(updated)
	})
}

func toRecordBranchState(row *ent.RecordBranchState) (RecordBranchState, error) {
	createdAt, err := parseStoreTime(row.CreatedAt)
	if err != nil {
		return RecordBranchState{}, err
	}
	updatedAt, err := parseStoreTime(row.UpdatedAt)
	if err != nil {
		return RecordBranchState{}, err
	}
	return RecordBranchState{
		ID:               row.ID,
		ProjectID:        row.ProjectID,
		RecordID:         row.RecordID,
		Branch:           row.Branch,
		Status:           RecordBranchStatus(row.Status),
		CurrentVersionID: row.CurrentVersionID,
		Malformed:        row.Malformed != 0,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}, nil
}
