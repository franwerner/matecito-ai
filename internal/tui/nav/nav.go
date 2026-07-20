// Package nav defines the navigation message types and screen enum used across
// the TUI layer. Screen packages import this package (not the parent tui package)
// to avoid import cycles: tui/app.go imports both nav and the screen packages,
// while screen packages only import nav.
package nav

import tea "github.com/charmbracelet/bubbletea"

// Screen identifies each top-level screen in the TUI router.
type Screen int

const (
	ScreenMenu Screen = iota
	ScreenInstall
	ScreenSync
	ScreenVerify
	ScreenConfig
	ScreenSddModel
	ScreenDomains
)

// NavigateMsg asks the AppModel router to switch to the given screen.
type NavigateMsg struct{ To Screen }

// OpenDomainConfigMsg asks AppModel to open the generic per-domain config screen
// (rendered from the domain's manifest config schema).
type OpenDomainConfigMsg struct{ Domain string }

// OpenModelsMsg asks AppModel to open the model-per-agent screen for a domain.
type OpenModelsMsg struct{ Domain string }

// BackMsg asks the router to return to the main menu.
type BackMsg struct{}

// QuitMsg terminates the bubbletea program.
type QuitMsg struct{}

// ToggleScopeMsg asks AppModel to flip the active config scope (Global ↔ Project).
// AppModel ignores it when InProject is false — scope is locked to Global outside a repo.
type ToggleScopeMsg struct{}

// ChildModel is the interface every screen model must satisfy so AppModel can
// route uniformly without knowing concrete screen types.
type ChildModel interface {
	Init() tea.Cmd
	Update(tea.Msg) (ChildModel, tea.Cmd)
	View() string
}
