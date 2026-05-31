package menu

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/tui/nav"
	"github.com/franwerner/matecito-ai/internal/tui/styles"
)

var items = []struct {
	label  string
	screen nav.Screen
	quit   bool
}{
	{"Instalar", nav.ScreenInstall, false},
	{"Verificar", nav.ScreenVerify, false},
	{"Configuración", nav.ScreenConfig, false},
	{"Salir", nav.ScreenMenu, true},
}

// MenuModel implementa el menú principal.
type MenuModel struct {
	cursor int
}

func New() MenuModel { return MenuModel{} }

func (m MenuModel) Init() tea.Cmd { return nil }

func (m MenuModel) Update(msg tea.Msg) (nav.ChildModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(items)-1 {
				m.cursor++
			}
		case "enter":
			if items[m.cursor].quit {
				return m, func() tea.Msg { return nav.QuitMsg{} }
			}
			return m, func() tea.Msg { return nav.NavigateMsg{To: items[m.cursor].screen} }
		case "q", "ctrl+c":
			return m, func() tea.Msg { return nav.QuitMsg{} }
		}
	}
	return m, nil
}

func (m MenuModel) View() string {
	var s string
	for i, item := range items {
		line := "  " + item.label
		if i == m.cursor {
			line = styles.Selected.Render("> " + item.label)
		}
		s += line + "\n"
	}
	s += "\n" + styles.Footer.Render("↑/↓ navegar  enter seleccionar  q salir")
	return s
}
