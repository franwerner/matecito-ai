package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/franwerner/matecito-ai/apps/api/internal/store"
)

func registerProject(ctx context.Context, t *testing.T, st store.Repository, id, name string) store.Project {
	t.Helper()
	p, err := st.RegisterProject(ctx, store.RegisterProjectCmd{ID: id, Name: name})
	if err != nil {
		t.Fatalf("RegisterProject(%s): %v", id, err)
	}
	return p
}

// TestBindProjectPath_CreatesAndResolves covers R4's "el path resuelve al
// project-id".
func TestBindProjectPath_CreatesAndResolves(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	bound, err := st.BindProjectPath(ctx, store.BindProjectPathCmd{ProjectID: p.ID, RootPath: "/repo/acme"})
	if err != nil {
		t.Fatalf("BindProjectPath: %v", err)
	}
	if bound.Status != store.ProjectPathActive {
		t.Fatalf("Status = %q, want %q", bound.Status, store.ProjectPathActive)
	}

	resolved, err := st.ResolveProjectPath(ctx, "/repo/acme")
	if err != nil {
		t.Fatalf("ResolveProjectPath: %v", err)
	}
	if resolved.ID != p.ID {
		t.Fatalf("resolved.ID = %q, want %q", resolved.ID, p.ID)
	}
}

// TestBindProjectPath_TwoWorktreesBothActive covers R4's "dos checkouts del
// mismo proyecto en la misma máquina": both paths stay active, no move is
// implied without an explicit MoveProjectPath.
func TestBindProjectPath_TwoWorktreesBothActive(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	if _, err := st.BindProjectPath(ctx, store.BindProjectPathCmd{ProjectID: p.ID, RootPath: "/repo/main"}); err != nil {
		t.Fatalf("bind main worktree: %v", err)
	}
	if _, err := st.BindProjectPath(ctx, store.BindProjectPathCmd{ProjectID: p.ID, RootPath: "/repo/feature-wt"}); err != nil {
		t.Fatalf("bind second worktree: %v", err)
	}

	for _, root := range []string{"/repo/main", "/repo/feature-wt"} {
		resolved, err := st.ResolveProjectPath(ctx, root)
		if err != nil {
			t.Fatalf("ResolveProjectPath(%s): %v", root, err)
		}
		if resolved.ID != p.ID {
			t.Fatalf("resolved.ID for %s = %q, want %q", root, resolved.ID, p.ID)
		}
	}
}

// TestMoveProjectPath_RepoMoved covers R4's "el repo se mueve": the old row
// transitions to stale and the new root resolves to the same project-id.
func TestMoveProjectPath_RepoMoved(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	if _, err := st.BindProjectPath(ctx, store.BindProjectPathCmd{ProjectID: p.ID, RootPath: "/old/path"}); err != nil {
		t.Fatalf("bind old path: %v", err)
	}

	moved, err := st.MoveProjectPath(ctx, store.MoveProjectPathCmd{
		ProjectID:    p.ID,
		FromRootPath: "/old/path",
		ToRootPath:   "/new/path",
	})
	if err != nil {
		t.Fatalf("MoveProjectPath: %v", err)
	}
	if moved.Status != store.ProjectPathActive {
		t.Fatalf("new path Status = %q, want %q", moved.Status, store.ProjectPathActive)
	}

	resolved, err := st.ResolveProjectPath(ctx, "/new/path")
	if err != nil {
		t.Fatalf("ResolveProjectPath(new): %v", err)
	}
	if resolved.ID != p.ID {
		t.Fatalf("resolved.ID = %q, want %q", resolved.ID, p.ID)
	}

	// The old path is stale, not deleted, and no longer resolves.
	_, err = st.ResolveProjectPath(ctx, "/old/path")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("ResolveProjectPath(old) err = %v, want ErrNotFound (stale rows don't resolve)", err)
	}
}

// TestBindProjectPath_PathBoundToAnotherProject covers the ErrPathAlreadyBound
// sentinel: the same (machine, root_path) cannot flip to a different
// project.
func TestBindProjectPath_PathBoundToAnotherProject(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	a := registerProject(ctx, t, st, "project-a", "acme")
	b := registerProject(ctx, t, st, "project-b", "beta")

	if _, err := st.BindProjectPath(ctx, store.BindProjectPathCmd{ProjectID: a.ID, RootPath: "/repo/shared"}); err != nil {
		t.Fatalf("bind for a: %v", err)
	}

	_, err := st.BindProjectPath(ctx, store.BindProjectPathCmd{ProjectID: b.ID, RootPath: "/repo/shared"})
	if !errors.Is(err, store.ErrPathAlreadyBound) {
		t.Fatalf("err = %v, want ErrPathAlreadyBound", err)
	}
}

// TestResolveProjectPath_NotFound covers the Reader's ErrNotFound contract.
func TestResolveProjectPath_NotFound(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)

	_, err := st.ResolveProjectPath(ctx, "/never/bound")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
