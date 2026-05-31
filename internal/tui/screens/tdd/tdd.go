package tdd

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	"github.com/franwerner/matecito-ai/internal/tui/nav"
	"github.com/franwerner/matecito-ai/internal/tui/styles"
)

// TddModel gestiona el toggle de strictTdd para el config activo según el scope.
// Opera sobre la ruta que recibe al construirse: puede ser el config global
// (~/.matecito-ai/config.json) o el per-proyecto (<repo>/.matecito-ai/config.json).
// AppModel.buildChild resuelve la ruta con ConfigPathForScope antes de llamar a New.
type TddModel struct {
	configPath string
	// cfg es el config cargado; si no existía, es un Config vacío.
	cfg agentmodel.Config
	// original guarda el valor inicial para detectar cambios.
	original *bool
	// current es el valor que el usuario está editando.
	current bool
}

// New construye un TddModel que opera sobre configPath.
// configPath ya está resuelto por AppModel según el scope activo.
func New(configPath string) TddModel {
	m := TddModel{configPath: configPath}

	cfg, err := agentmodel.Load(configPath)
	if err == nil && cfg != nil {
		m.cfg = *cfg
	}

	if m.cfg.StrictTdd != nil {
		m.current = *m.cfg.StrictTdd
		v := *m.cfg.StrictTdd
		m.original = &v
	}

	return m
}

func (m TddModel) Init() tea.Cmd { return nil }

func (m TddModel) Update(msg tea.Msg) (nav.ChildModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			m.current = !m.current
		case "q", "esc":
			// persistir y volver
			return m, m.saveAndBack()
		case "ctrl+c":
			// descartar sin guardar
			return m, func() tea.Msg { return nav.QuitMsg{} }
		}
	}
	return m, nil
}

func (m TddModel) View() string {
	var sb strings.Builder

	sb.WriteString(styles.Title.Render("TDD (este proyecto)") + "\n\n")

	state := "deshabilitado"
	if m.current {
		state = "habilitado"
	}

	marker := "  [ ] "
	if m.current {
		marker = styles.Selected.Render("  [x] ")
	}
	sb.WriteString(marker + "Strict TDD — " + styles.Dimmed.Render(state) + "\n\n")

	if m.original == nil {
		sb.WriteString(styles.Dimmed.Render("  (hereda del config global / default: false)") + "\n\n")
	}

	sb.WriteString(styles.Footer.Render("enter/space toggle  q/esc guardar  ctrl+c descartar"))
	return sb.String()
}

// saveAndBack persiste el config en configPath y emite BackMsg.
func (m TddModel) saveAndBack() tea.Cmd {
	return func() tea.Msg {
		m.cfg.StrictTdd = &m.current
		_ = agentmodel.Save(m.configPath, &m.cfg)
		return nav.BackMsg{}
	}
}
