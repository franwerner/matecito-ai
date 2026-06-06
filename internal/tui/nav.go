package tui

// Re-export nav types so package tui consumers don't need to change import paths.
// Screen packages must import internal/tui/nav directly to avoid the import cycle
// that would arise from: tui/app.go → screens/* → tui.
import "github.com/franwerner/matecito-ai/internal/tui/nav"

// Screen identifies each top-level screen in the TUI router.
type Screen = nav.Screen

const (
	ScreenMenu     = nav.ScreenMenu
	ScreenInstall  = nav.ScreenInstall
	ScreenSync     = nav.ScreenSync
	ScreenVerify   = nav.ScreenVerify
	ScreenConfig   = nav.ScreenConfig
	ScreenSddModel     = nav.ScreenSddModel
	ScreenTdd          = nav.ScreenTdd
	ScreenDecisionGaps = nav.ScreenDecisionGaps
)

// NavigateMsg asks the AppModel router to switch to the given screen.
type NavigateMsg = nav.NavigateMsg

// BackMsg asks the router to return to the main menu.
type BackMsg = nav.BackMsg

// QuitMsg terminates the bubbletea program.
type QuitMsg = nav.QuitMsg

// ReExecMsg asks the program to quit so the binary can be re-exec'd after the
// terminal is restored.
type ReExecMsg = nav.ReExecMsg

// ToggleScopeMsg asks AppModel to flip the active config scope (Global ↔ Project).
type ToggleScopeMsg = nav.ToggleScopeMsg

// ChildModel is the interface every screen model must satisfy so AppModel can
// route uniformly without knowing concrete screen types.
type ChildModel = nav.ChildModel
