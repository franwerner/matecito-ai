package config

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	"github.com/franwerner/matecito-ai/internal/manifest"
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
	// domain, when set, makes the entry open the per-domain config screen
	// (rendered from that domain's manifest schema) instead of a fixed screen.
	domain string
}

// ConfigMenuModel es el submenú de configuración: lo compartido (General) más
// una entrada por dominio activo (M7).
type ConfigMenuModel struct {
	entries []menuEntry
	cursor  int
}

func New(ctx ProjectContext, scope agentmodel.Scope) ConfigMenuModel {
	// Shared (cross-domain) entries. Auto-mine (flagDecisionGaps) is now per-domain
	// (rendered inside each domain's config screen), not a top-level entry.
	entries := []menuEntry{
		{label: "Dominios (global)", screen: nav.ScreenDomains},
	}
	// One entry per active domain → its generic config screen.
	if ids, payloadFS, err := manifest.ResolveFromEnv(); err == nil {
		for _, id := range ids {
			label := id
			if mf, e := manifest.Load(payloadFS, id); e == nil && mf.Label != "" {
				label = mf.Label
			}
			entries = append(entries, menuEntry{label: label + " (config)", domain: id})
		}
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
				e := m.entries[m.cursor]
				if e.domain != "" {
					domain := e.domain
					return m, func() tea.Msg { return nav.OpenDomainConfigMsg{Domain: domain} }
				}
				screen := e.screen
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
