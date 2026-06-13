package sddmodel

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	"github.com/franwerner/matecito-ai/internal/setup/deploy"
	"github.com/franwerner/matecito-ai/internal/tui/nav"
	"github.com/franwerner/matecito-ai/internal/tui/styles"
)

// GlobalSentinel es el valor especial que indica "sin override por proyecto —
// heredar del config global". Solo existe durante la edición en memoria;
// nunca se persiste en config.json.
const GlobalSentinel = "(global)"

// AgentModelModel permite seleccionar el modelo (opus/sonnet/haiku/fable) para cada
// uno de los 10 agentes SDD. En Project scope ofrece además el estado "(global)"
// que elimina el override per-proyecto al guardar.
// q/Esc → guarda y vuelve. Ctrl+C → descarta sin guardar.
type AgentModelModel struct {
	configPath string
	scope      agentmodel.Scope
	// domain is the area domain whose agents are being configured (M7); agents is
	// that domain's agent list (the rows). Persisted under domainConfig[domain].
	domain string
	agents []string
	// globalCfg y projectCfg son los configs cargados al construirse la pantalla;
	// se usan para mostrar el modelo efectivo cuando el valor es GlobalSentinel.
	globalCfg  *agentmodel.Config
	projectCfg *agentmodel.Config
	// models es el mapa editable de agente → modelo durante la sesión.
	models  map[string]string
	cursor  int
	saveErr error
}

// saveErrMsg reporta que el guardado falló; mantiene la pantalla abierta
// para que el usuario vea el error en vez de volver al menú silenciosamente.
type saveErrMsg struct{ err error }

// New construye un AgentModelModel cargando el config del scope activo.
// globalCfg y projectCfg son pre-cargados por AppModel.buildChild para
// evitar I/O duplicado y para tener ambos disponibles en render.
func New(globalCfg *agentmodel.Config, projectCfg *agentmodel.Config, configPath string, scope agentmodel.Scope, domain string, agents []string) AgentModelModel {
	m := AgentModelModel{
		configPath: configPath,
		scope:      scope,
		domain:     domain,
		agents:     agents,
		globalCfg:  globalCfg,
		projectCfg: projectCfg,
		models:     make(map[string]string),
	}

	// cargar modelos desde el config del scope activo
	var activeCfg *agentmodel.Config
	if scope == agentmodel.ScopeProject {
		activeCfg = projectCfg
	} else {
		activeCfg = globalCfg
	}

	if activeCfg != nil && len(activeCfg.DomainModels(domain)) > 0 {
		for k, v := range activeCfg.DomainModels(domain) {
			m.models[k] = v
		}
	} else if scope == agentmodel.ScopeGlobal {
		// en scope global sin config, sembrar desde defaults del payload
		payloadFS, _, fsErr := deploy.ResolvePayloadFS()
		if fsErr == nil {
			defaults, defErr := agentmodel.DefaultsForDomain(payloadFS, domain)
			if defErr == nil {
				for k, v := range defaults {
					m.models[k] = v
				}
			}
		}
	}

	// en Project scope: agentes sin override → sentinel "(global)"
	// en Global scope: agentes sin entrada → sonnet como fallback
	for _, agent := range m.agents {
		if _, ok := m.models[agent]; !ok {
			if scope == agentmodel.ScopeProject {
				m.models[agent] = GlobalSentinel
			} else {
				m.models[agent] = agentmodel.ValidModels[2] // sonnet
			}
		}
	}

	return m
}

func (m AgentModelModel) Init() tea.Cmd { return nil }

func (m AgentModelModel) Update(msg tea.Msg) (nav.ChildModel, tea.Cmd) {
	switch msg := msg.(type) {
	case saveErrMsg:
		m.saveErr = msg.err
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.agents)-1 {
				m.cursor++
			}

		// ciclar modelo
		case "right", "l":
			m.cycleModel(1)
		case "left", "h":
			m.cycleModel(-1)

		// selección directa por índice de ValidModels (1-based: 1=fable … 4=haiku)
		case "1", "2", "3", "4":
			if idx := int(msg.String()[0] - '1'); idx < len(agentmodel.ValidModels) {
				m.models[m.agents[m.cursor]] = agentmodel.ValidModels[idx]
			}

		case "q", "esc":
			return m, m.saveAndBack()

		case "ctrl+c":
			return m, func() tea.Msg { return nav.QuitMsg{} }
		}
	}
	return m, nil
}

// cycleModel avanza o retrocede el modelo del agente actual.
// En Project scope el ciclo incluye GlobalSentinel como primera opción.
func (m *AgentModelModel) cycleModel(delta int) {
	agent := m.agents[m.cursor]
	current := m.models[agent]

	options := m.cycleOptions()
	idx := 0
	for i, v := range options {
		if v == current {
			idx = i
			break
		}
	}
	n := len(options)
	idx = ((idx+delta)%n + n) % n
	m.models[agent] = options[idx]
}

// cycleOptions retorna los valores posibles al ciclar, según el scope activo.
// En Project scope se antepone GlobalSentinel para permitir limpiar el override.
func (m *AgentModelModel) cycleOptions() []string {
	if m.scope == agentmodel.ScopeProject {
		return append([]string{GlobalSentinel}, agentmodel.ValidModels...)
	}
	return agentmodel.ValidModels
}

func (m AgentModelModel) View() string {
	var sb strings.Builder

	title := "Modelos por agente — " + m.domain
	if m.scope == agentmodel.ScopeProject {
		title += " — scope: proyecto"
	}
	sb.WriteString(styles.Title.Render(title) + "\n\n")

	for i, agent := range m.agents {
		model := m.models[agent]
		modelStr := m.renderModelPills(agent, model)

		row := fmt.Sprintf("  %-14s  %s", agent, modelStr)
		if i == m.cursor {
			row = styles.Selected.Render(fmt.Sprintf("> %-14s  %s", agent, modelStr))
		}
		sb.WriteString(row + "\n")
	}

	if m.saveErr != nil {
		sb.WriteString("\n" + styles.Error.Render("Error al guardar: "+m.saveErr.Error()) + "\n")
	}

	sb.WriteString("\n" + styles.Footer.Render(
		"↑/↓ agente  ←/→ modelo  1-4 directo  q/esc guardar  ctrl+c descartar",
	))
	return sb.String()
}

// renderModelPills muestra los modelos válidos con el activo resaltado.
// Cuando el valor es GlobalSentinel, muestra el modelo efectivo resuelto como referencia.
func (m AgentModelModel) renderModelPills(agent, active string) string {
	if active == GlobalSentinel {
		resolved, _ := agentmodel.ResolveModel(m.globalCfg, nil, m.domain, agent)
		ref := "(hereda global)"
		if resolved != "" {
			ref = fmt.Sprintf("(hereda: %s)", resolved)
		}
		return styles.Dimmed.Render(GlobalSentinel + "  " + ref)
	}

	parts := make([]string, len(agentmodel.ValidModels))
	for i, v := range agentmodel.ValidModels {
		if v == active {
			parts[i] = styles.Selected.Render("[" + v + "]")
		} else {
			parts[i] = styles.Dimmed.Render(" " + v + " ")
		}
	}
	return strings.Join(parts, " ")
}

// saveAndBack persiste el config del scope activo con los modelos editados y emite BackMsg.
// Las entradas con valor GlobalSentinel se eliminan del mapa (no se persisten).
func (m AgentModelModel) saveAndBack() tea.Cmd {
	return func() tea.Msg {
		cfg, err := agentmodel.Load(m.configPath)
		if err != nil || cfg == nil {
			cfg = &agentmodel.Config{}
		}

		for agent, model := range m.models {
			if model == GlobalSentinel {
				// limpiar el override per-proyecto para este agente
				cfg.SetDomainModelOverride(m.domain, agent, "")
			} else {
				cfg.SetDomainModelOverride(m.domain, agent, model)
			}
		}

		if err := agentmodel.Save(m.configPath, cfg); err != nil {
			return saveErrMsg{err: err}
		}
		return nav.BackMsg{}
	}
}
