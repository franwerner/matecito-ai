package store

import (
	"context"
	"fmt"

	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent"
	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent/projectpath"
)

// ResolveProjectPath resolves this machine's canonical root to the project
// it currently checks out (register-project, R4). Only `active` rows
// resolve — a `stale` row is history, not a live mapping.
func (e *Engine) ResolveProjectPath(ctx context.Context, root string) (Project, error) {
	row, err := e.client.ProjectPath.Query().
		Where(
			projectpath.MachineID(e.machineID),
			projectpath.RootPath(root),
			projectpath.Status(string(ProjectPathActive)),
		).
		Only(ctx)
	if ent.IsNotFound(err) {
		return Project{}, ErrNotFound
	}
	if err != nil {
		return Project{}, fmt.Errorf("store: resolve project path %q: %w", root, err)
	}
	return e.FindProjectByID(ctx, row.ProjectID)
}

// BindProjectPath binds cmd.RootPath to cmd.ProjectID on this machine.
func (e *Engine) BindProjectPath(ctx context.Context, cmd BindProjectPathCmd) (ProjectPath, error) {
	return submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (ProjectPath, error) {
		return bindPath(ctx, tx, e.machineID, cmd.ProjectID, cmd.RootPath)
	})
}

// MoveProjectPath binds cmd.ToRootPath and retires cmd.FromRootPath in one
// write, so a relocated checkout never leaves two `active` rows behind.
func (e *Engine) MoveProjectPath(ctx context.Context, cmd MoveProjectPathCmd) (ProjectPath, error) {
	return submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (ProjectPath, error) {
		bound, err := bindPath(ctx, tx, e.machineID, cmd.ProjectID, cmd.ToRootPath)
		if err != nil {
			return ProjectPath{}, err
		}

		if cmd.FromRootPath != cmd.ToRootPath {
			if err := staleOldPath(ctx, tx, e.machineID, cmd.ProjectID, cmd.FromRootPath); err != nil {
				return ProjectPath{}, err
			}
		}

		return bound, nil
	})
}

func bindPath(ctx context.Context, tx *ent.Tx, machineID, projectID, rootPath string) (ProjectPath, error) {
	now := nowStoreTime()

	existing, err := tx.ProjectPath.Query().
		Where(projectpath.MachineID(machineID), projectpath.RootPath(rootPath)).
		Only(ctx)
	switch {
	case ent.IsNotFound(err):
		created, err := tx.ProjectPath.Create().
			SetProjectID(projectID).
			SetMachineID(machineID).
			SetRootPath(rootPath).
			SetStatus(string(ProjectPathActive)).
			SetLastSeenAt(now).
			Save(ctx)
		if err != nil {
			return ProjectPath{}, fmt.Errorf("store: bind project path: %w", err)
		}
		return toProjectPath(created)
	case err != nil:
		return ProjectPath{}, fmt.Errorf("store: lookup project path %q: %w", rootPath, err)
	}

	if existing.ProjectID != projectID {
		return ProjectPath{}, ErrPathAlreadyBound
	}

	updated, err := tx.ProjectPath.UpdateOneID(existing.ID).
		SetStatus(string(ProjectPathActive)).
		SetLastSeenAt(now).
		Save(ctx)
	if err != nil {
		return ProjectPath{}, fmt.Errorf("store: rebind project path %q: %w", rootPath, err)
	}
	return toProjectPath(updated)
}

// staleOldPath marks the previous checkout root `stale`, if it still
// exists, still belongs to this project and is still active. A caller that
// names a previous path which no longer applies is a silent no-op, not an
// error — the transition is best-effort bookkeeping, never a precondition
// of the bind itself.
func staleOldPath(ctx context.Context, tx *ent.Tx, machineID, projectID, previousRootPath string) error {
	old, err := tx.ProjectPath.Query().
		Where(
			projectpath.MachineID(machineID),
			projectpath.RootPath(previousRootPath),
			projectpath.ProjectID(projectID),
			projectpath.Status(string(ProjectPathActive)),
		).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("store: lookup previous project path %q: %w", previousRootPath, err)
	}

	if err := tx.ProjectPath.UpdateOneID(old.ID).
		SetStatus(string(ProjectPathStale)).
		Exec(ctx); err != nil {
		return fmt.Errorf("store: stale project path %q: %w", previousRootPath, err)
	}
	return nil
}

func toProjectPath(row *ent.ProjectPath) (ProjectPath, error) {
	lastSeenAt, err := parseStoreTime(row.LastSeenAt)
	if err != nil {
		return ProjectPath{}, err
	}
	createdAt, err := parseStoreTime(row.CreatedAt)
	if err != nil {
		return ProjectPath{}, err
	}
	updatedAt, err := parseStoreTime(row.UpdatedAt)
	if err != nil {
		return ProjectPath{}, err
	}
	return ProjectPath{
		ID:         row.ID,
		ProjectID:  row.ProjectID,
		MachineID:  row.MachineID,
		RootPath:   row.RootPath,
		Status:     ProjectPathStatus(row.Status),
		LastSeenAt: lastSeenAt,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}
