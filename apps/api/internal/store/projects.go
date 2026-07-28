package store

import (
	"context"
	"fmt"

	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent"
	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent/project"
	"github.com/google/uuid"
)

// FindProjectByID is the Reader lookup by project-id (R4/R3).
func (e *Engine) FindProjectByID(ctx context.Context, id string) (Project, error) {
	row, err := e.client.Project.Get(ctx, id)
	if ent.IsNotFound(err) {
		return Project{}, ErrNotFound
	}
	if err != nil {
		return Project{}, fmt.Errorf("store: find project %s: %w", id, err)
	}
	return toProject(row)
}

// FindProjectsByName supports find_project(name): name is NOT unique, so
// several matches can come back for the caller to confirm
// (register-project).
func (e *Engine) FindProjectsByName(ctx context.Context, name string) ([]Project, error) {
	rows, err := e.client.Project.Query().
		Where(project.Name(name)).
		Order(project.ByCreatedAt()).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("store: find projects by name %q: %w", name, err)
	}

	out := make([]Project, 0, len(rows))
	for _, row := range rows {
		p, err := toProject(row)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

// RegisterProject is idempotent by ID: registering the same project-id
// twice returns the existing row unchanged, never a duplicate
// (register-project's "idempotente" scenario). An empty cmd.ID models the
// CREATE case and gets a freshly generated UUID v7.
func (e *Engine) RegisterProject(ctx context.Context, cmd RegisterProjectCmd) (Project, error) {
	return submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (Project, error) {
		id := cmd.ID
		if id == "" {
			generated, err := uuid.NewV7()
			if err != nil {
				return Project{}, fmt.Errorf("store: generate project id: %w", err)
			}
			id = generated.String()
		}

		existing, err := tx.Project.Get(ctx, id)
		switch {
		case err == nil:
			return toProject(existing)
		case ent.IsNotFound(err):
			created, err := tx.Project.Create().
				SetID(id).
				SetName(cmd.Name).
				Save(ctx)
			if err != nil {
				return Project{}, fmt.Errorf("store: register project: %w", err)
			}
			return toProject(created)
		default:
			return Project{}, fmt.Errorf("store: lookup project %s: %w", id, err)
		}
	})
}

// DeactivateProject is the logical deregister (R4): the project and its
// history (changes, events, records) stay intact, only status flips to
// inactive.
func (e *Engine) DeactivateProject(ctx context.Context, id string) error {
	_, err := submit(ctx, e.writer, func(ctx context.Context, tx *ent.Tx) (struct{}, error) {
		if err := tx.Project.UpdateOneID(id).
			SetStatus(string(ProjectInactive)).
			Exec(ctx); err != nil {
			if ent.IsNotFound(err) {
				return struct{}{}, ErrNotFound
			}
			return struct{}{}, fmt.Errorf("store: deactivate project %s: %w", id, err)
		}
		return struct{}{}, nil
	})
	return err
}

func toProject(row *ent.Project) (Project, error) {
	createdAt, err := parseStoreTime(row.CreatedAt)
	if err != nil {
		return Project{}, err
	}
	updatedAt, err := parseStoreTime(row.UpdatedAt)
	if err != nil {
		return Project{}, err
	}
	return Project{
		ID:        row.ID,
		Name:      row.Name,
		Status:    ProjectStatus(row.Status),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
