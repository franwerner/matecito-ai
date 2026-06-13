// Package domains is the TUI screen that selects the user's active area domains.
// It persists agentmodel.Config.Domains in the global config. The selection only
// governs MCP registration and the verify display immediately; the deployed
// domain payloads (CLAUDE.md fragments, agents, skills, references) take effect
// after a redeploy via the "Actualizar" (sync) flow — Option 1: save + instruct.
package domains

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	"github.com/franwerner/matecito-ai/internal/manifest"
	"github.com/franwerner/matecito-ai/internal/setup/deploy"
	"github.com/franwerner/matecito-ai/internal/tui/nav"
	"github.com/franwerner/matecito-ai/internal/tui/styles"
)

type DomainsModel struct {
	configPath string
	cfg        agentmodel.Config
	ids        []string          // discovered domain ids, sorted (deterministic)
	labels     map[string]string // id → display label
	selected   map[string]bool
	cursor     int
	loadErr    string
}

// New builds the screen over configPath (the global config). It discovers the
// domains shipped in the payload and seeds the selection from cfg.Domains, or
// from "all" when unconfigured (the compat shim).
func New(configPath string) DomainsModel {
	m := DomainsModel{
		configPath: configPath,
		labels:     map[string]string{},
		selected:   map[string]bool{},
	}

	if cfg, err := agentmodel.Load(configPath); err == nil && cfg != nil {
		m.cfg = *cfg
	}

	payloadFS, _, err := deploy.ResolvePayloadFS()
	if err != nil {
		m.loadErr = "no se pudo resolver el payload: " + err.Error()
		return m
	}
	ids, err := manifest.DiscoverIDs(payloadFS)
	if err != nil {
		m.loadErr = "no se pudieron descubrir dominios: " + err.Error()
		return m
	}
	m.ids = ids

	for _, id := range ids {
		if mf, err := manifest.Load(payloadFS, id); err == nil && mf.Label != "" {
			m.labels[id] = mf.Label
		} else {
			m.labels[id] = id
		}
	}

	if len(m.cfg.Domains) == 0 {
		for _, id := range ids {
			m.selected[id] = true
		}
	} else {
		for _, id := range m.cfg.Domains {
			m.selected[id] = true
		}
	}
	return m
}

func (m DomainsModel) Init() tea.Cmd { return nil }

func (m DomainsModel) Update(msg tea.Msg) (nav.ChildModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.ids)-1 {
				m.cursor++
			}
		case " ":
			if len(m.ids) > 0 {
				id := m.ids[m.cursor]
				m.selected[id] = !m.selected[id]
			}
		case "enter", "q", "esc":
			// An empty list means "all" via the compat shim, so refuse to save a
			// zero selection (it would silently re-enable every domain).
			if m.countSelected() == 0 {
				return m, nil
			}
			return m, m.saveAndBack()
		case "ctrl+c":
			return m, func() tea.Msg { return nav.QuitMsg{} }
		}
	}
	return m, nil
}

func (m DomainsModel) countSelected() int {
	n := 0
	for _, id := range m.ids {
		if m.selected[id] {
			n++
		}
	}
	return n
}

// saveAndBack persists the selection and returns to the menu. When every domain
// is selected it stores nil so the shim keeps auto-including future domains;
// otherwise it stores the explicit subset.
func (m DomainsModel) saveAndBack() tea.Cmd {
	return func() tea.Msg {
		var chosen []string
		for _, id := range m.ids {
			if m.selected[id] {
				chosen = append(chosen, id)
			}
		}
		if len(chosen) == len(m.ids) {
			m.cfg.Domains = nil
		} else {
			m.cfg.Domains = chosen
		}
		_ = agentmodel.Save(m.configPath, &m.cfg)
		return nav.BackMsg{}
	}
}

func (m DomainsModel) View() string {
	var sb strings.Builder
	sb.WriteString(styles.Title.Render("Dominios (global)") + "\n\n")

	if m.loadErr != "" {
		sb.WriteString(styles.Dimmed.Render("  "+m.loadErr) + "\n\n")
		sb.WriteString(styles.Footer.Render("esc volver  ctrl+c salir"))
		return sb.String()
	}

	for i, id := range m.ids {
		marker := "  [ ] "
		if m.selected[id] {
			marker = "  [x] "
		}
		line := marker + m.labels[id] + " " + styles.Dimmed.Render("("+id+")")
		if i == m.cursor {
			line = styles.Selected.Render(line)
		}
		sb.WriteString(line + "\n")
	}
	sb.WriteString("\n")

	if m.countSelected() == 0 {
		sb.WriteString(styles.Dimmed.Render("  seleccioná al menos un dominio para guardar") + "\n\n")
	} else {
		sb.WriteString(styles.Dimmed.Render("  tras guardar, corré «Actualizar» para aplicar el cambio") + "\n\n")
	}

	sb.WriteString(styles.Footer.Render("↑↓ mover  space toggle  enter/q guardar  ctrl+c descartar"))
	return sb.String()
}
