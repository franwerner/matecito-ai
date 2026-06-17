package debugger

import (
	"testing"

	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/mcp"
)

// TestDetectMCP_RegisteredReportsOK exercises the StatusOK branch by overriding
// the find seam to report the MCP as registered, which the host environment
// cannot guarantee deterministically.
func TestDetectMCP_RegisteredReportsOK(t *testing.T) {
	orig := find
	t.Cleanup(func() { find = orig })
	find = func(string) (mcp.Found, bool) {
		return mcp.Found{Name: "debugger", Connected: true, Source: "cli"}, true
	}

	r := detectMCP()
	if r.Status != check.StatusOK {
		t.Errorf("Status = %v, want StatusOK", r.Status)
	}
	if r.Detail == "" {
		t.Error("StatusOK result must carry a Detail from Describe()")
	}
}

// TestDetectMCP_UnregisteredReportsMissing exercises the StatusMissing branch
// deterministically, independent of whether the host has the MCP registered.
func TestDetectMCP_UnregisteredReportsMissing(t *testing.T) {
	orig := find
	t.Cleanup(func() { find = orig })
	find = func(string) (mcp.Found, bool) {
		return mcp.Found{}, false
	}

	r := detectMCP()
	if r.Status != check.StatusMissing {
		t.Errorf("Status = %v, want StatusMissing", r.Status)
	}
	if r.FixHint == "" {
		t.Error("StatusMissing result must carry a non-empty FixHint")
	}
}
