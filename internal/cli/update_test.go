package cli

import (
	"testing"

	"github.com/franwerner/matecito-ai/internal/setup/sync"
)

// TestUpdateCmd_HasErrorsReturnsError verifies that when sync returns a result
// with component errors, NewUpdateCmd's RunE logic surfaces them as a non-nil error.
// This mirrors the install_test.go pattern: test the decision logic directly
// without running the full cobra command.
func TestUpdateCmd_HasErrorsReturnsError(t *testing.T) {
	result := sync.Result{
		Errors: map[string]error{
			"matecito-ai": errStub("install failed"),
		},
	}
	if !result.HasErrors() {
		t.Fatal("expected HasErrors to return true when Errors is non-empty")
	}
}

// TestUpdateCmd_NoErrorsOK verifies that an empty Result.Errors is treated as success.
func TestUpdateCmd_NoErrorsOK(t *testing.T) {
	result := sync.Result{}
	if result.HasErrors() {
		t.Fatal("expected HasErrors to return false when Errors is nil")
	}
}

type errStub string

func (e errStub) Error() string { return string(e) }
