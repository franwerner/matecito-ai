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
		"unknown":   "mcp__unknown__*",              // not in registry → convention
	}
	for name, want := range cases {
		if got := permissionPattern(name); got != want {
			t.Errorf("permissionPattern(%q) = %q, want %q", name, got, want)
		}
	}
}
