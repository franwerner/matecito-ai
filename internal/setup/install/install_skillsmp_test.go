package install_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/setup/install"
)

// TestSkillsmpMCPStep_ClaudeAbsent verifies that skillsmpMCPStep.Run returns an
// error mentioning "claude" when the claude binary is not in PATH, without
// attempting to register the remote MCP.
func TestSkillsmpMCPStep_ClaudeAbsent(t *testing.T) {
	d := tempDir(t)
	// claude is intentionally NOT written → absent from isolated PATH.
	isolatedPATH(t, d)

	var out bytes.Buffer
	opts := install.Options{Stdout: &out, Stderr: &out, Yes: true}
	steps := install.AllSteps(opts)

	var skillsmpRun func() error
	for _, s := range steps {
		if s.Name == "skillsmp MCP (remote)" {
			skillsmpRun = s.Run
			break
		}
	}
	if skillsmpRun == nil {
		t.Fatal("skillsmp MCP step not found in AllSteps")
	}

	runErr := skillsmpRun()
	if runErr == nil {
		t.Fatal("expected error when claude binary is absent, got nil")
	}
	if !strings.Contains(runErr.Error(), "claude") {
		t.Errorf("error should mention 'claude'; got: %v", runErr)
	}
}

// TestSkillsmpMCPStep_RemoteHTTPPlan verifies skillsmp is registered as a
// remote-HTTP MCP (same pattern as figma/canva), keyless — no --scope, --env,
// or --header — matching the anonymous tier by default.
func TestSkillsmpMCPStep_RemoteHTTPPlan(t *testing.T) {
	var out bytes.Buffer
	opts := install.Options{Stdout: &out, Stderr: &out, Yes: true}
	steps := install.AllSteps(opts)

	var skillsmpPlan string
	found := false
	for _, s := range steps {
		if s.Name == "skillsmp MCP (remote)" {
			skillsmpPlan = s.Plan
			found = true
			break
		}
	}
	if !found {
		t.Fatal("skillsmp MCP step not found in AllSteps")
	}

	wantPlan := "claude mcp add --transport http skillsmp https://skillsmp.com/mcp"
	if skillsmpPlan != wantPlan {
		t.Errorf("skillsmp Plan = %q, want %q", skillsmpPlan, wantPlan)
	}
	for _, forbidden := range []string{"--scope", "--env", "--header", "npx"} {
		if strings.Contains(skillsmpPlan, forbidden) {
			t.Errorf("skillsmp Plan %q must not contain %q (keyless anonymous tier)", skillsmpPlan, forbidden)
		}
	}
}
