package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent"
	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent/change"
)

// FindChangeByBranch resolves the change for (projectID, branch) — the
// only key a change ever has (rule/event-scoping, model M1).
func (e *Engine) FindChangeByBranch(ctx context.Context, projectID, branch string) (Change, error) {
	row, err := e.client.Change.Query().
		Where(change.ProjectID(projectID), change.Branch(branch)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return Change{}, ErrNotFound
	}
	if err != nil {
		return Change{}, fmt.Errorf("store: find change by branch %q: %w", branch, err)
	}
	return toChange(row)
}

// StartChange starts or continues the change for (cmd.ProjectID,
// cmd.Branch) (flow/start-change). An existing row for that key is always
// returned as-is — a second call with a different cmd.Name is the
// "conflicto de nombre" scenario, not a rename.
func (e *Engine) StartChange(ctx context.Context, cmd StartChangeCmd) (Change, error) {
	return submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (Change, error) {
		existing, err := tx.Change.Query().
			Where(change.ProjectID(cmd.ProjectID), change.Branch(cmd.Branch)).
			Only(ctx)
		switch {
		case err == nil:
			if existing.Name != cmd.Name {
				return Change{}, fmt.Errorf("%w: branch %q already mapped to change %q", ErrChangeNameTaken, cmd.Branch, existing.Name)
			}
			return toChange(existing)
		case ent.IsNotFound(err):
			created, err := tx.Change.Create().
				SetProjectID(cmd.ProjectID).
				SetBranch(cmd.Branch).
				SetName(cmd.Name).
				Save(ctx)
			if err != nil {
				return Change{}, translateChangeConstraint(err)
			}
			return toChange(created)
		default:
			return Change{}, fmt.Errorf("store: lookup change (project %s, branch %q): %w", cmd.ProjectID, cmd.Branch, err)
		}
	})
}

// SetChangeStatus is the only way a change moves between active and closed
// (lifecycle/change: never a side effect of another operation).
func (e *Engine) SetChangeStatus(ctx context.Context, id string, s ChangeStatus) error {
	_, err := submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (struct{}, error) {
		if err := tx.Change.UpdateOneID(id).
			SetStatus(string(s)).
			Exec(ctx); err != nil {
			if ent.IsNotFound(err) {
				return struct{}{}, ErrNotFound
			}
			return struct{}{}, fmt.Errorf("store: set change %s status: %w", id, err)
		}
		return struct{}{}, nil
	})
	return err
}

// translateChangeConstraint maps a failed Change insert's unique-index
// violation to its specific sentinel — changes has two independent unique
// indexes (project_id+branch, project_id+name) and the caller needs to
// know which one fired, not just "some constraint failed".
func translateChangeConstraint(err error) error {
	if !ent.IsConstraintError(err) {
		return fmt.Errorf("store: create change: %w", err)
	}

	msg := err.Error()
	switch {
	case strings.Contains(msg, "changes.branch") || strings.Contains(msg, "change_project_id_branch"):
		return fmt.Errorf("%w: %w", ErrBranchTaken, err)
	case strings.Contains(msg, "changes.name") || strings.Contains(msg, "change_project_id_name"):
		return fmt.Errorf("%w: %w", ErrChangeNameTaken, err)
	default:
		return fmt.Errorf("store: create change: %w", err)
	}
}

func toChange(row *ent.Change) (Change, error) {
	createdAt, err := parseStoreTime(row.CreatedAt)
	if err != nil {
		return Change{}, err
	}
	updatedAt, err := parseStoreTime(row.UpdatedAt)
	if err != nil {
		return Change{}, err
	}
	return Change{
		ID:        row.ID,
		ProjectID: row.ProjectID,
		Branch:    row.Branch,
		Name:      row.Name,
		Status:    ChangeStatus(row.Status),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
