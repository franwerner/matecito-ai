package debugger_test

import (
	"testing"

	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/checks/debugger"
)

// TestAll_ReturnsOneResult verifies that All() always returns exactly one
// check.Result regardless of MCP registration state.
func TestAll_ReturnsOneResult(t *testing.T) {
	results := debugger.All()
	if len(results) != 1 {
		t.Fatalf("All() returned %d results, want 1", len(results))
	}
}

// TestAll_ResultName verifies the result name identifies the debugger MCP.
func TestAll_ResultName(t *testing.T) {
	results := debugger.All()
	if results[0].Name != "debugger MCP" {
		t.Errorf("Name = %q, want %q", results[0].Name, "debugger MCP")
	}
}

// TestAll_StatusIsKnown verifies the result carries a recognized status.
func TestAll_StatusIsKnown(t *testing.T) {
	results := debugger.All()
	r := results[0]
	switch r.Status {
	case check.StatusOK, check.StatusMissing:
		// expected
	default:
		t.Errorf("unexpected status %v; want StatusOK or StatusMissing", r.Status)
	}
}

// TestAll_MissingCarriesFixHint verifies that when the MCP is not registered
// the result includes a non-empty FixHint so the user knows how to register it.
func TestAll_MissingCarriesFixHint(t *testing.T) {
	results := debugger.All()
	r := results[0]
	if r.Status == check.StatusMissing && r.FixHint == "" {
		t.Error("StatusMissing result must carry a non-empty FixHint")
	}
}
