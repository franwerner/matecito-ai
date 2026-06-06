package decisiongaps

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	"github.com/franwerner/matecito-ai/internal/tui/nav"
	"github.com/franwerner/matecito-ai/internal/tui/styles"
)

// DecisionGapsModel gestiona el toggle de flagDecisionGaps para el config activo
// según el scope. Opera sobre la ruta que recibe al construirse: global
// (~/.matecito-ai/config.json) o per-proyecto (<repo>/.matecito-ai/config.json).
// AppModel.buildChild resuelve la ruta con ConfigPathForScope antes de llamar a New.
type DecisionGapsModel struct {
	configPath string
	// cfg es el config cargado; si no existía, es un Config vacío.
	cfg agentmodel.Config
	// original guarda el valor inicial para detectar cambios.
	original *bool
	// current es el valor que el usuario está editando.
	current bool
}

// New construye un DecisionGapsModel que opera sobre configPath.
// configPath ya está resuelto por AppModel según el scope activo.
func New(configPath string) DecisionGapsModel {
	m := DecisionGapsModel{configPath: configPath}

	cfg, err := agentmodel.Load(configPath)
	if err == nil && cfg != nil {
		m.cfg = *cfg
	}

	if m.cfg.FlagDecisionGaps != nil {
		m.current = *m.cfg.FlagDecisionGaps
		v := *m.cfg.FlagDecisionGaps
		m.original = &v
	}

	return m
}

func (m DecisionGapsModel) Init() tea.Cmd { return nil }

func (m DecisionGapsModel) Update(msg tea.Msg) (nav.ChildModel, tea.Cmd) {
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

func (m DecisionGapsModel) View() string {
	var sb strings.Builder

	sb.WriteString(styles.Title.Render("Auto-mine ADR (este proyecto)") + "\n\n")

	state := "deshabilitado"
	if m.current {
		state = "habilitado"
	}

	marker := "  [ ] "
	if m.current {
		marker = styles.Selected.Render("  [x] ")
	}
	sb.WriteString(marker + "flagDecisionGaps — " + styles.Dimmed.Render(state) + "\n\n")

	sb.WriteString(styles.Dimmed.Render("  Detecta huecos de decisión en el flujo SDD y ofrece minarlos como ADRs Inferred.") + "\n")
	sb.WriteString(styles.Dimmed.Render("  Solo actúa si además existe .matecito-ai/adr/ con contenido.") + "\n\n")

	if m.original == nil {
		sb.WriteString(styles.Dimmed.Render("  (hereda del config global / default: false)") + "\n\n")
	}

	sb.WriteString(styles.Footer.Render("enter/space toggle  q/esc guardar  ctrl+c descartar"))
	return sb.String()
}

// saveAndBack persiste el config en configPath y emite BackMsg.
func (m DecisionGapsModel) saveAndBack() tea.Cmd {
	return func() tea.Msg {
		m.cfg.FlagDecisionGaps = &m.current
		_ = agentmodel.Save(m.configPath, &m.cfg)
		return nav.BackMsg{}
	}
}
