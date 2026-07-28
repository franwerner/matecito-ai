package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/franwerner/matecito-ai/apps/api/internal/store"
)

// TestStartChange_CreatesForNewBranch covers flow/start-change's "alta con
// rama".
func TestStartChange_CreatesForNewBranch(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	c, err := st.StartChange(ctx, store.StartChangeCmd{ProjectID: p.ID, Branch: "feat/x", Name: "add-feature"})
	if err != nil {
		t.Fatalf("StartChange: %v", err)
	}
	if c.Branch != "feat/x" || c.Name != "add-feature" {
		t.Fatalf("got Change{Branch: %q, Name: %q}, want {feat/x, add-feature}", c.Branch, c.Name)
	}
	if c.Status != store.ChangeActive {
		t.Fatalf("Status = %q, want %q", c.Status, store.ChangeActive)
	}
}

// TestStartChange_ContinuesCrossSession covers "continúa cross-session":
// starting the same (project, branch) again returns the same change, never
// a duplicate.
func TestStartChange_ContinuesCrossSession(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	first, err := st.StartChange(ctx, store.StartChangeCmd{ProjectID: p.ID, Branch: "feat/x", Name: "add-feature"})
	if err != nil {
		t.Fatalf("first StartChange: %v", err)
	}

	second, err := st.StartChange(ctx, store.StartChangeCmd{ProjectID: p.ID, Branch: "feat/x", Name: "add-feature"})
	if err != nil {
		t.Fatalf("second StartChange: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("second.ID = %q, want %q (same change, not a duplicate)", second.ID, first.ID)
	}

	found, err := st.FindChangeByBranch(ctx, p.ID, "feat/x")
	if err != nil {
		t.Fatalf("FindChangeByBranch: %v", err)
	}
	if found.ID != first.ID {
		t.Fatalf("FindChangeByBranch.ID = %q, want %q", found.ID, first.ID)
	}
}

// TestStartChange_BranchRemapConflict covers the "conflicto de nombre"
// scenario: a second StartChange on the same branch with a different name
// is rejected, not applied as a rename.
func TestStartChange_BranchRemapConflict(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	if _, err := st.StartChange(ctx, store.StartChangeCmd{ProjectID: p.ID, Branch: "feat/x", Name: "add-feature"}); err != nil {
		t.Fatalf("first StartChange: %v", err)
	}

	_, err := st.StartChange(ctx, store.StartChangeCmd{ProjectID: p.ID, Branch: "feat/x", Name: "different-name"})
	if !errors.Is(err, store.ErrChangeNameTaken) {
		t.Fatalf("err = %v, want ErrChangeNameTaken", err)
	}

	// Nothing changed: the branch still maps to its original name.
	unchanged, err := st.FindChangeByBranch(ctx, p.ID, "feat/x")
	if err != nil {
		t.Fatalf("FindChangeByBranch: %v", err)
	}
	if unchanged.Name != "add-feature" {
		t.Fatalf("Name = %q, want unchanged %q", unchanged.Name, "add-feature")
	}
}

// TestStartChange_NameConflictAcrossBranches covers R5's "nombre único por
// proyecto": a distinct branch cannot claim a name already used by another
// branch in the same project.
func TestStartChange_NameConflictAcrossBranches(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	if _, err := st.StartChange(ctx, store.StartChangeCmd{ProjectID: p.ID, Branch: "feat/a", Name: "shared-name"}); err != nil {
		t.Fatalf("first StartChange: %v", err)
	}

	_, err := st.StartChange(ctx, store.StartChangeCmd{ProjectID: p.ID, Branch: "feat/b", Name: "shared-name"})
	if !errors.Is(err, store.ErrChangeNameTaken) {
		t.Fatalf("err = %v, want ErrChangeNameTaken", err)
	}
}

// TestStartChange_SameBranchDifferentProjectsDoNotCross covers "dos ramas
// del mismo proyecto no se cruzan" transposed to two different projects
// sharing a branch name — each resolves independently.
func TestStartChange_SameBranchDifferentProjectsDoNotCross(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	a := registerProject(ctx, t, st, "project-a", "acme")
	b := registerProject(ctx, t, st, "project-b", "beta")

	ca, err := st.StartChange(ctx, store.StartChangeCmd{ProjectID: a.ID, Branch: "main", Name: "change-a"})
	if err != nil {
		t.Fatalf("StartChange a: %v", err)
	}
	cb, err := st.StartChange(ctx, store.StartChangeCmd{ProjectID: b.ID, Branch: "main", Name: "change-b"})
	if err != nil {
		t.Fatalf("StartChange b: %v", err)
	}
	if ca.ID == cb.ID {
		t.Fatalf("expected distinct changes, got the same ID %q", ca.ID)
	}
}

// TestSetChangeStatus_CloseAndReopen covers lifecycle/change's close/reopen
// transitions.
func TestSetChangeStatus_CloseAndReopen(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	c, err := st.StartChange(ctx, store.StartChangeCmd{ProjectID: p.ID, Branch: "feat/x", Name: "add-feature"})
	if err != nil {
		t.Fatalf("StartChange: %v", err)
	}

	if err := st.SetChangeStatus(ctx, c.ID, store.ChangeClosed); err != nil {
		t.Fatalf("SetChangeStatus(closed): %v", err)
	}
	closed, err := st.FindChangeByBranch(ctx, p.ID, "feat/x")
	if err != nil {
		t.Fatalf("FindChangeByBranch after close: %v", err)
	}
	if closed.Status != store.ChangeClosed {
		t.Fatalf("Status = %q, want %q", closed.Status, store.ChangeClosed)
	}

	if err := st.SetChangeStatus(ctx, c.ID, store.ChangeActive); err != nil {
		t.Fatalf("SetChangeStatus(active): %v", err)
	}
	reopened, err := st.FindChangeByBranch(ctx, p.ID, "feat/x")
	if err != nil {
		t.Fatalf("FindChangeByBranch after reopen: %v", err)
	}
	if reopened.Status != store.ChangeActive {
		t.Fatalf("Status = %q, want %q", reopened.Status, store.ChangeActive)
	}
}

// TestSetChangeStatus_NotFound covers the Writer's ErrNotFound contract.
func TestSetChangeStatus_NotFound(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)

	err := st.SetChangeStatus(ctx, "no-such-change", store.ChangeClosed)
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

// TestFindChangeByBranch_NotFound covers the Reader's ErrNotFound contract.
func TestFindChangeByBranch_NotFound(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)
	p := registerProject(ctx, t, st, "project-1", "acme")

	_, err := st.FindChangeByBranch(ctx, p.ID, "no-such-branch")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

// TestChangeNameUnique_RawInsertRejected exercises R5's unique constraint
// directly against the database, bypassing the Repository and the
// generated client (delivery/testing-strategy: constraint rejections are
// exercised with raw SQL).
func TestChangeNameUnique_RawInsertRejected(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	st, err := store.Open(ctx, dir, testMachineID)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() {
		if err := st.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})

	raw := openRaw(t, dir)
	fx := seedFixtures(ctx, t, raw)

	_, err = raw.ExecContext(ctx, `INSERT INTO changes
		(id, created_at, updated_at, branch, name, status, project_id)
		VALUES ('change-2', ?, ?, 'feat/other', 'acme-change', 'active', ?)`, ts, ts, fx.projectID)
	if err == nil {
		t.Fatal("expected a UNIQUE(project_id, name) violation, got nil error")
	}
}

// TestChangeBranchNotNull_RawInsertRejected covers R5's "no hay change sin
// rama": the database itself rejects a change with a null or absent branch
// — no sentinel, no admitted null — exercised with raw SQL against the
// migrated database, bypassing the Repository and the generated client
// (R13: constraint rejections are exercised with raw SQL, not through the
// Repository, since the guarantee lives in the database, not in this
// binary).
func TestChangeBranchNotNull_RawInsertRejected(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	st, err := store.Open(ctx, dir, testMachineID)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() {
		if err := st.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})

	raw := openRaw(t, dir)
	fx := seedFixtures(ctx, t, raw)

	cases := []struct {
		name  string
		query string
	}{
		{
			name: "branch explicitly null",
			query: `INSERT INTO changes
				(id, created_at, updated_at, branch, name, status, project_id)
				VALUES ('bad-change-null-branch', ?, ?, NULL, 'null-branch-change', 'active', ?)`,
		},
		{
			name: "branch column omitted",
			query: `INSERT INTO changes
				(id, created_at, updated_at, name, status, project_id)
				VALUES ('bad-change-no-branch', ?, ?, 'no-branch-change', 'active', ?)`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := raw.ExecContext(ctx, tc.query, ts, ts, fx.projectID)
			if err == nil {
				t.Fatal("expected a NOT NULL violation on changes.branch, got nil error")
			}
		})
	}
}
