package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/franwerner/matecito-ai/apps/api/internal/store"
)

func openStore(t *testing.T) store.Repository {
	t.Helper()
	ctx := context.Background()
	st, err := store.Open(ctx, t.TempDir(), testMachineID)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() {
		if err := st.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})
	return st
}

// TestRegisterProject_CreateGeneratesID covers the CREATE case
// (flow/register-project): an empty ID gets a freshly generated project-id.
func TestRegisterProject_CreateGeneratesID(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)

	p, err := st.RegisterProject(ctx, store.RegisterProjectCmd{Name: "acme"})
	if err != nil {
		t.Fatalf("RegisterProject: %v", err)
	}
	if p.ID == "" {
		t.Fatal("expected a generated project id, got empty string")
	}
	if p.Name != "acme" {
		t.Fatalf("Name = %q, want %q", p.Name, "acme")
	}
	if p.Status != store.ProjectActive {
		t.Fatalf("Status = %q, want %q", p.Status, store.ProjectActive)
	}
}

// TestRegisterProject_IdempotentByID covers register-project's "idempotente"
// scenario: registering the same project-id twice returns the same row,
// never a duplicate.
func TestRegisterProject_IdempotentByID(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)

	first, err := st.RegisterProject(ctx, store.RegisterProjectCmd{ID: "project-fixed", Name: "acme"})
	if err != nil {
		t.Fatalf("first RegisterProject: %v", err)
	}

	second, err := st.RegisterProject(ctx, store.RegisterProjectCmd{ID: "project-fixed", Name: "acme-renamed"})
	if err != nil {
		t.Fatalf("second RegisterProject: %v", err)
	}

	if second.ID != first.ID {
		t.Fatalf("second.ID = %q, want %q (same project, not a duplicate)", second.ID, first.ID)
	}
	// The idempotent re-register must not silently rename the existing row.
	if second.Name != first.Name {
		t.Fatalf("second.Name = %q, want unchanged %q", second.Name, first.Name)
	}

	matches, err := st.FindProjectsByName(ctx, "acme")
	if err != nil {
		t.Fatalf("FindProjectsByName: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("len(matches) = %d, want 1 (no duplicate row)", len(matches))
	}
}

// TestFindProjectsByName_MultipleMatches covers "name sin unicidad": two
// distinct projects with the same name both come back for the caller to
// confirm.
func TestFindProjectsByName_MultipleMatches(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)

	if _, err := st.RegisterProject(ctx, store.RegisterProjectCmd{ID: "project-a", Name: "shared-name"}); err != nil {
		t.Fatalf("register a: %v", err)
	}
	if _, err := st.RegisterProject(ctx, store.RegisterProjectCmd{ID: "project-b", Name: "shared-name"}); err != nil {
		t.Fatalf("register b: %v", err)
	}

	matches, err := st.FindProjectsByName(ctx, "shared-name")
	if err != nil {
		t.Fatalf("FindProjectsByName: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("len(matches) = %d, want 2", len(matches))
	}
}

// TestFindProjectByID_NotFound covers the Reader's ErrNotFound contract.
func TestFindProjectByID_NotFound(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)

	_, err := st.FindProjectByID(ctx, "no-such-project")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

// TestDeactivateProject_LogicalOnly covers R4's "desregistrar no borra": the
// project stays resolvable, only its status flips.
func TestDeactivateProject_LogicalOnly(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)

	p, err := st.RegisterProject(ctx, store.RegisterProjectCmd{Name: "acme"})
	if err != nil {
		t.Fatalf("RegisterProject: %v", err)
	}

	if err := st.DeactivateProject(ctx, p.ID); err != nil {
		t.Fatalf("DeactivateProject: %v", err)
	}

	got, err := st.FindProjectByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("FindProjectByID after deactivate: %v", err)
	}
	if got.Status != store.ProjectInactive {
		t.Fatalf("Status = %q, want %q", got.Status, store.ProjectInactive)
	}
}

// TestDeactivateProject_NotFound covers the Writer's ErrNotFound contract on
// a target that does not exist.
func TestDeactivateProject_NotFound(t *testing.T) {
	ctx := context.Background()
	st := openStore(t)

	err := st.DeactivateProject(ctx, "no-such-project")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
