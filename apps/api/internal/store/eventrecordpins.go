package store

import (
	"context"
	"fmt"

	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent"
	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent/eventrecordpin"
	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent/recordbranchstate"
)

// PinEventToRecordVersion pins cmd.EventID to (cmd.RecordID, cmd.Branch)'s
// current version, resolved implicitly here — the caller never names a
// version id (lifecycle/record-version: "el agente nunca manda una
// versión"). The pin freezes that version if this is its first pin;
// repeating the same (event, version) pair is idempotent and returns the
// existing pin unchanged — the same idiom every other natural-key write in
// this package follows (RegisterProject, StartChange, AppendEvent,
// PutContentObject), rather than a new conflict sentinel with no scenario
// that needs one.
func (e *Engine) PinEventToRecordVersion(ctx context.Context, cmd PinEventToRecordVersionCmd) (EventRecordPin, error) {
	return submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (EventRecordPin, error) {
		branchState, err := tx.RecordBranchState.Query().
			Where(recordbranchstate.RecordID(cmd.RecordID), recordbranchstate.Branch(cmd.Branch)).
			Only(ctx)
		if ent.IsNotFound(err) {
			return EventRecordPin{}, ErrNotFound
		}
		if err != nil {
			return EventRecordPin{}, fmt.Errorf("store: lookup record branch state %s@%s: %w", cmd.RecordID, cmd.Branch, err)
		}

		existing, err := tx.EventRecordPin.Query().
			Where(eventrecordpin.EventID(cmd.EventID), eventrecordpin.RecordVersionID(branchState.CurrentVersionID)).
			Only(ctx)
		switch {
		case err == nil:
			return toEventRecordPin(existing)
		case !ent.IsNotFound(err):
			return EventRecordPin{}, fmt.Errorf("store: lookup event record pin: %w", err)
		}

		if err := freezeRecordVersion(ctx, tx, branchState.CurrentVersionID); err != nil {
			return EventRecordPin{}, err
		}

		created, err := tx.EventRecordPin.Create().
			SetProjectID(cmd.ProjectID).
			SetEventID(cmd.EventID).
			SetRecordVersionID(branchState.CurrentVersionID).
			Save(ctx)
		if err != nil {
			return EventRecordPin{}, fmt.Errorf("store: create event record pin: %w", err)
		}
		return toEventRecordPin(created)
	})
}

// freezeRecordVersion moves versionID from current to frozen — the
// one-way transition lifecycle/record-version's pin triggers. A version
// that is already frozen (a second event pinning the same version) is left
// untouched: freezing is idempotent, never re-stamped.
func freezeRecordVersion(ctx context.Context, tx *ent.Tx, versionID string) error {
	version, err := tx.RecordVersion.Get(ctx, versionID)
	if err != nil {
		return fmt.Errorf("store: get record version %s: %w", versionID, err)
	}
	if RecordVersionStatus(version.Status) == VersionFrozen {
		return nil
	}

	if err := tx.RecordVersion.UpdateOneID(versionID).
		SetStatus(string(VersionFrozen)).
		SetFrozenAt(nowStoreTime()).
		Exec(ctx); err != nil {
		return fmt.Errorf("store: freeze record version %s: %w", versionID, err)
	}
	return nil
}

func toEventRecordPin(row *ent.EventRecordPin) (EventRecordPin, error) {
	createdAt, err := parseStoreTime(row.CreatedAt)
	if err != nil {
		return EventRecordPin{}, err
	}
	updatedAt, err := parseStoreTime(row.UpdatedAt)
	if err != nil {
		return EventRecordPin{}, err
	}
	return EventRecordPin{
		ID:              row.ID,
		ProjectID:       row.ProjectID,
		EventID:         row.EventID,
		RecordVersionID: row.RecordVersionID,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}
