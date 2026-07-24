package install

import "testing"

func TestPermissionPattern_ConventionAndOverride(t *testing.T) {
	cases := map[string]string{
		"engram":    "mcp__plugin_engram_engram__*", // override: installed as a Claude Code plugin
		"context7":  "mcp__context7__*",             // convention
		"codegraph": "mcp__codegraph__*",            // convention
		"drawio":    "mcp__drawio__*",               // convention
		"debugger":  "mcp__debugger__*",             // convention
		"figma":     "mcp__figma__*",                // convention
		"skillsmp":  "mcp__skillsmp__*",             // convention
		"unknown":   "mcp__unknown__*",              // not in registry → convention
	}
	for name, want := range cases {
		if got := permissionPattern(name); got != want {
			t.Errorf("permissionPattern(%q) = %q, want %q", name, got, want)
		}
	}
}

// TestMCPRegistry_Skillsmp verifies skillsmp is registered as a first-class MCP
// (step present) and included in the safety-net defaultMCP set.
func TestMCPRegistry_Skillsmp(t *testing.T) {
	def, ok := mcpRegistry["skillsmp"]
	if !ok {
		t.Fatal("mcpRegistry missing \"skillsmp\" entry")
	}
	if def.step == nil {
		t.Error("mcpRegistry[\"skillsmp\"].step is nil")
	}
	found := false
	for _, name := range defaultMCP {
		if name == "skillsmp" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("defaultMCP = %v, want it to contain \"skillsmp\"", defaultMCP)
	}
}
