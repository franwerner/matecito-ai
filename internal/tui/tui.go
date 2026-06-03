package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	pkgsync "github.com/franwerner/matecito-ai/internal/setup/sync"
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
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	// Re-exec recién acá: p.Run() ya retornó, así que bubbletea restauró la
	// terminal (raw mode, alt screen, bracketed paste). En Unix syscall.Exec
	// reemplaza el proceso; en Windows ReExec() es no-op (devuelve nil) y la
	// vista ya mostró el aviso de reinicio manual.
	if app, ok := finalModel.(AppModel); ok && app.reexecRequested {
		return pkgsync.ReExec()
	}
	return nil
}
