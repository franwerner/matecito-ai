package header

import (
	"fmt"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	"github.com/franwerner/matecito-ai/internal/tui/styles"
)

// Header holds the state required to render the persistent top bar.
// Version and ProjectName are set at construction; LatestTag is filled
// asynchronously by the release-check tea.Cmd.
// Scope and InProject govern the scope indicator appended after the project name.
type Header struct {
	Version     string
	LatestTag   string
	ProjectName string
	Scope       agentmodel.Scope
	InProject   bool
}

// Render produces a pure string for the header line:
//
//	matecito-ai vX.Y.Z  [update available: <tag>]  <project>  Scope: Global  [s: toggle]
//
// The badge is included only when agentmodel.ShouldShowBadge returns true.
// The project name is appended only when non-empty.
// The scope indicator is always shown; the toggle hint appears only when InProject is true.
func (h Header) Render() string {
	line := styles.Wordmark() + " " + styles.Dimmed.Render(h.Version)

	if agentmodel.ShouldShowBadge(h.Version, h.LatestTag) {
		line += styles.Accent.Render(fmt.Sprintf("  [update available: %s]", h.LatestTag))
	}

	if h.ProjectName != "" {
		line += fmt.Sprintf("  %s", h.ProjectName)
	}

	scopeLabel := "Global"
	if h.Scope == agentmodel.ScopeProject && h.ProjectName != "" {
		scopeLabel = h.ProjectName
	}
	line += "  Scope: " + styles.Accent.Render(scopeLabel)

	if h.InProject {
		line += styles.Dimmed.Render("  [s: toggle]")
	}

	return line
}
