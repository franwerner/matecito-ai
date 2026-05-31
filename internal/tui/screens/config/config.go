package config

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	"github.com/franwerner/matecito-ai/internal/tui/nav"
	"github.com/franwerner/matecito-ai/internal/tui/styles"
)

// ProjectContext describe si el TUI se lanzó desde dentro de un proyecto.
// Se pasa al constructor para decidir si mostrar la entrada TDD.
type ProjectContext struct {
	InProject bool
	RepoRoot  string
}

type menuEntry struct {
	label  string
	screen nav.Screen
}

// ConfigMenuModel es el submenú de configuración.
// "TDD (este proyecto)" sólo aparece cuando InProject==true (R6.4/S6.4-S6.5).
type ConfigMenuModel struct {
	entries []menuEntry
	cursor  int
}

func New(ctx ProjectContext, scope agentmodel.Scope) ConfigMenuModel {
	// el label de TDD refleja el scope activo: en Global edita el config global,
	// en Project el del repo. La ruta efectiva la resuelve AppModel con ConfigPathForScope.
	tddLabel := "TDD (este proyecto)"
	if scope == agentmodel.ScopeGlobal {
		tddLabel = "TDD (global)"
	}
	entries := []menuEntry{
		{"Modelos por agente (sdd-model)", nav.ScreenSddModel},
		{tddLabel, nav.ScreenTdd},
	}
	return ConfigMenuModel{entries: entries}
}

func (m ConfigMenuModel) Init() tea.Cmd { return nil }

func (m ConfigMenuModel) Update(msg tea.Msg) (nav.ChildModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.entries) > 0 {
				screen := m.entries[m.cursor].screen
				return m, func() tea.Msg { return nav.NavigateMsg{To: screen} }
			}
		case "esc", "backspace", "b":
			return m, func() tea.Msg { return nav.BackMsg{} }
		case "q", "ctrl+c":
			return m, func() tea.Msg { return nav.QuitMsg{} }
		}
	}
	return m, nil
}

func (m ConfigMenuModel) View() string {
	var sb strings.Builder

	sb.WriteString(styles.Title.Render("Configuración") + "\n\n")

	for i, e := range m.entries {
		line := "  " + e.label
		if i == m.cursor {
			line = styles.Selected.Render("> " + e.label)
		}
		sb.WriteString(line + "\n")
	}

	sb.WriteString("\n" + styles.Footer.Render("↑/↓ navegar  enter seleccionar  esc volver  q salir"))
	return sb.String()
}
