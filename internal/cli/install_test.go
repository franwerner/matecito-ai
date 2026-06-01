package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/setup/sync"
)

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
