package hook_test

import (
	"testing"

	"github.com/franwerner/matecito-ai/internal/hook"
)

// snapshotRegistry saves the current registry state and restores it on cleanup,
// preventing cross-test bleed via the package-level registry.
func snapshotRegistry(t *testing.T) {
	t.Helper()
	saved := hook.Registered()
	orig := make([]hook.Hook, len(saved))
	copy(orig, saved)
	t.Cleanup(func() {
		// Reset registry to the state captured before this test.
		// We reach into the registry by re-registering only original hooks;
		// the exported API lacks a Reset, so we exploit that Register appends
		// and Registered returns the slice — a deliberate package-level state
		// that tests must bracket. We use an internal test-only reset instead
		// via the ResetRegistry helper exposed from export_test.go.
		hook.ResetRegistry(orig)
	})
}

// TestForDomains_IncludesSharedHook_EmptyActive verifies that a hook with
// Domain == SharedDomain is returned even when the active set is empty.
func TestForDomains_IncludesSharedHook_EmptyActive(t *testing.T) {
	snapshotRegistry(t)

	hook.Register(hook.Hook{
		Domain:     hook.SharedDomain,
		Subcommand: "shared-sentinel",
		Event:      "PreToolUse",
		Run:        func(_ []byte) hook.Result { return hook.Result{} },
	})

	got := hook.ForDomains([]string{})
	if len(got) != 1 || got[0].Subcommand != "shared-sentinel" {
		t.Errorf("expected shared hook in result, got: %v", got)
	}
}

// TestForDomains_IncludesSharedHook_PartialActive verifies that both a
// SharedDomain hook and an active-domain hook are returned together.
func TestForDomains_IncludesSharedHook_PartialActive(t *testing.T) {
	snapshotRegistry(t)

	hook.Register(hook.Hook{
		Domain:     hook.SharedDomain,
		Subcommand: "shared-sentinel",
		Event:      "PreToolUse",
		Run:        func(_ []byte) hook.Result { return hook.Result{} },
	})
	hook.Register(hook.Hook{
		Domain:     "development",
		Subcommand: "dev-hook",
		Event:      "PreToolUse",
		Run:        func(_ []byte) hook.Result { return hook.Result{} },
	})

	got := hook.ForDomains([]string{"development"})
	if len(got) != 2 {
		t.Fatalf("expected 2 hooks, got %d: %v", len(got), got)
	}
	var foundShared, foundDev bool
	for _, h := range got {
		switch h.Subcommand {
		case "shared-sentinel":
			foundShared = true
		case "dev-hook":
			foundDev = true
		}
	}
	if !foundShared {
		t.Error("SharedDomain hook missing from result")
	}
	if !foundDev {
		t.Error("development hook missing from result")
	}
}

// TestForDomains_ExcludesNonSharedInactiveHook verifies that a hook whose
// Domain is neither in the active set nor SharedDomain is excluded.
func TestForDomains_ExcludesNonSharedInactiveHook(t *testing.T) {
	snapshotRegistry(t)

	hook.Register(hook.Hook{
		Domain:     "other",
		Subcommand: "other-hook",
		Event:      "PreToolUse",
		Run:        func(_ []byte) hook.Result { return hook.Result{} },
	})

	for _, active := range [][]string{{}, {"development"}} {
		got := hook.ForDomains(active)
		for _, h := range got {
			if h.Subcommand == "other-hook" {
				t.Errorf("unexpected hook %q in result for active=%v", h.Subcommand, active)
			}
		}
	}
}
