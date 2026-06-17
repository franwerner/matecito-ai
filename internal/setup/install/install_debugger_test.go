package install_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/setup/install"
)

// TestDebuggerMCPStep_ClaudeAbsent verifies that debuggerMCPStep.Run returns an
// error mentioning "claude" when the claude binary is not in PATH, without
// attempting to run npx.
func TestDebuggerMCPStep_ClaudeAbsent(t *testing.T) {
	d := tempDir(t)
	// claude is intentionally NOT written → absent from isolated PATH.
	// npx is also absent so we confirm claude is checked first.
	isolatedPATH(t, d)

	var out bytes.Buffer
	opts := install.Options{Stdout: &out, Stderr: &out, Yes: true}
	steps := install.AllSteps(opts)

	var debuggerRun func() error
	for _, s := range steps {
		if s.Name == "debugger MCP (mcp-debugger)" {
			debuggerRun = s.Run
			break
		}
	}
	if debuggerRun == nil {
		t.Fatal("debugger MCP step not found in AllSteps")
	}

	runErr := debuggerRun()
	if runErr == nil {
		t.Fatal("expected error when claude binary is absent, got nil")
	}
	if !strings.Contains(runErr.Error(), "claude") {
		t.Errorf("error should mention 'claude'; got: %v", runErr)
	}
}
