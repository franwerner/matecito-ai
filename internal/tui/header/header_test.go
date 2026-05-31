package header_test

import (
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	"github.com/franwerner/matecito-ai/internal/tui/header"
)

func TestRender(t *testing.T) {
	cases := []struct {
		name        string
		version     string
		latestTag   string
		projectName string
		wantBadge   bool
		wantProject bool
	}{
		{
			name:      "badge shown when newer version available",
			version:   "0.2.0",
			latestTag: "v0.3.0",
			wantBadge: true,
		},
		{
			name:      "badge hidden when versions equal after normalization",
			version:   "0.2.0",
			latestTag: "v0.2.0",
			wantBadge: false,
		},
		{
			name:      "badge hidden when latest is empty",
			version:   "0.2.0",
			latestTag: "",
			wantBadge: false,
		},
		{
			name:      "badge hidden for dev build",
			version:   "0.1.0-dev",
			latestTag: "v0.3.0",
			wantBadge: false,
		},
		{
			name:        "project name shown when non-empty",
			version:     "0.2.0",
			projectName: "matecito-ai",
			wantProject: true,
		},
		{
			name:        "project name omitted when empty",
			version:     "0.2.0",
			projectName: "",
			wantProject: false,
		},
		{
			name:        "badge and project shown together",
			version:     "0.2.0",
			latestTag:   "v0.3.0",
			projectName: "my-project",
			wantBadge:   true,
			wantProject: true,
		},
		{
			name:        "dev build: badge hidden, project still shown",
			version:     "0.1.0-dev",
			latestTag:   "v0.3.0",
			projectName: "my-project",
			wantBadge:   false,
			wantProject: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := header.Header{
				Version:     tc.version,
				LatestTag:   tc.latestTag,
				ProjectName: tc.projectName,
			}
			rendered := h.Render()

			hasBadge := strings.Contains(rendered, "update available")
			if hasBadge != tc.wantBadge {
				t.Errorf("badge: got %v, want %v\nrendered: %q", hasBadge, tc.wantBadge, rendered)
			}

			hasProject := tc.projectName != "" && strings.Contains(rendered, tc.projectName)
			if hasProject != tc.wantProject {
				t.Errorf("project name: got %v, want %v\nrendered: %q", hasProject, tc.wantProject, rendered)
			}

			if !strings.Contains(rendered, tc.version) {
				t.Errorf("version %q not found in rendered output: %q", tc.version, rendered)
			}

			// el indicador de scope siempre aparece en el header
			if !strings.Contains(rendered, "Scope:") {
				t.Errorf("scope indicator missing in rendered output: %q", rendered)
			}
		})
	}
}

func TestRenderScope(t *testing.T) {
	cases := []struct {
		name        string
		scope       agentmodel.Scope
		inProject   bool
		projectName string
		wantScope   string
		wantToggle  bool
	}{
		{
			name:       "global scope outside repo — no toggle hint",
			scope:      agentmodel.ScopeGlobal,
			inProject:  false,
			wantScope:  "Scope: Global",
			wantToggle: false,
		},
		{
			name:        "global scope inside repo — toggle hint visible",
			scope:       agentmodel.ScopeGlobal,
			inProject:   true,
			projectName: "my-project",
			wantScope:   "Scope: Global",
			wantToggle:  true,
		},
		{
			name:        "project scope inside repo — shows project name",
			scope:       agentmodel.ScopeProject,
			inProject:   true,
			projectName: "my-project",
			wantScope:   "Scope: my-project",
			wantToggle:  true,
		},
		{
			name:       "project scope but no project name — falls back to Global label",
			scope:      agentmodel.ScopeProject,
			inProject:  true,
			wantScope:  "Scope: Global",
			wantToggle: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := header.Header{
				Version:     "0.2.0",
				Scope:       tc.scope,
				InProject:   tc.inProject,
				ProjectName: tc.projectName,
			}
			rendered := h.Render()

			if !strings.Contains(rendered, tc.wantScope) {
				t.Errorf("scope label: want %q in %q", tc.wantScope, rendered)
			}

			hasToggle := strings.Contains(rendered, "[s: toggle]")
			if hasToggle != tc.wantToggle {
				t.Errorf("toggle hint: got %v, want %v\nrendered: %q", hasToggle, tc.wantToggle, rendered)
			}
		})
	}
}
