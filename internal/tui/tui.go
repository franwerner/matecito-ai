package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
)

// RunOpts carries the parameters the caller must supply to start the TUI.
type RunOpts struct {
	Version string
}

// Run builds and launches the full-screen bubbletea TUI.
// It detects the project context and resolves the global config path before
// handing off to tea.NewProgram; it returns any program error directly.
func Run(opts RunOpts) error {
	ctx := DetectProject()

	globalConfigPath, err := agentmodel.ConfigPath()
	if err != nil {
		globalConfigPath = ""
	}

	appModel := NewAppModel(opts.Version, globalConfigPath, ctx)

	p := tea.NewProgram(appModel, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
