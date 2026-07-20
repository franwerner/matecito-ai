package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/setup/sync"
)

// TestInstallCmd_ResumeWiring documents (and guards) the precondition behind
// NewInstallCmd's "Resume: sync.ResumeRequested()" line and its
// "if !syncOpts.Resume { ...plan/dry-run/confirm... }" short-circuit: both
// read the resume flag from the same env var ResumeRequested() checks.
//
// RunE itself is not unit-tested end-to-end here: it hardcodes os.Stdin/
// os.Stdout (not injectable), and reaching the resume branch for real would
// require Sync's full Detect() round-trip (network calls to check binary/
// deploy versions) and deploy.BackupDir() (real filesystem paths) — none of
// which are seams exposed by the current design. TestSync_Resume in
// internal/setup/sync/sync_test.go covers the actual skip/no-plan/no-prompt
// behavior at the engine level once Resume is set; this test covers the one
// piece install.go's RunE contributes: deriving that flag from the env.
func TestInstallCmd_ResumeWiring(t *testing.T) {
	t.Run("ResumeRequested reflects MATECITO_RESUME=1", func(t *testing.T) {
		t.Setenv("MATECITO_RESUME", "1")
		if !sync.ResumeRequested() {
			t.Fatal("expected sync.ResumeRequested() to be true with MATECITO_RESUME=1, matching the value install.go assigns to syncOpts.Resume")
		}
	})

	t.Run("ResumeRequested is false without the env var", func(t *testing.T) {
		t.Setenv("MATECITO_RESUME", "")
		if sync.ResumeRequested() {
			t.Fatal("expected sync.ResumeRequested() to be false without MATECITO_RESUME=1")
		}
	})
}

// TestInstallCmd_CodegraphErrorSurfaced verifies that when syncResult carries a
// codegraph error, the install command returns a non-nil error and does NOT
// silently continue to MCP registration (requirement: CLI reports real
// codegraph install state).
//
// The test reaches the codegraph-error branch by invoking the logic that was
// previously "_ = syncResult" directly through a unit-level helper that
// mirrors what RunE does: if syncResult.Errors["codegraph"] is set, the
// command must return an error wrapping it.
func TestInstallCmd_CodegraphErrorSurfaced(t *testing.T) {
	cgErr := errors.New("npm install failed")
	result := sync.Result{
		Errors: map[string]error{
			"codegraph": cgErr,
		},
	}

	// Replicate the decision logic introduced in RunE:
	// if cgErr := syncResult.Errors["codegraph"]; cgErr != nil { return fmt.Errorf(...) }
	reportedErr := surfaceCodegraphError(result)
	if reportedErr == nil {
		t.Fatal("expected error to be surfaced when syncResult.Errors['codegraph'] is set, got nil")
	}
	if !strings.Contains(reportedErr.Error(), "codegraph") {
		t.Errorf("returned error should mention 'codegraph'; got: %v", reportedErr)
	}
	if !errors.Is(reportedErr, cgErr) {
		t.Errorf("returned error should wrap the original; got: %v", reportedErr)
	}
}

// TestInstallCmd_CodegraphSuccessNoError verifies that when syncResult has no
// codegraph error, surfaceCodegraphError returns nil (normal path continues).
func TestInstallCmd_CodegraphSuccessNoError(t *testing.T) {
	result := sync.Result{}
	if err := surfaceCodegraphError(result); err != nil {
		t.Fatalf("expected nil when codegraph succeeded, got: %v", err)
	}
}
