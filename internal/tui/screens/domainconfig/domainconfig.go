// Package domainconfig renders a domain's manifest config schema generically.
// `agent-models` fields open the model-per-agent screen; `bool`/`enum` are edited
// inline. Values persist under domainConfig[domain] in the config at configPath.
// This is the renderer that makes a domain "auto-configure" from its manifest.
package domainconfig

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	"github.com/franwerner/matecito-ai/internal/manifest"
	"github.com/franwerner/matecito-ai/internal/tui/nav"
	"github.com/franwerner/matecito-ai/internal/tui/styles"
)

type DomainConfigModel struct {
	domain     string
	label      string
	configPath string
	fields     []manifest.ConfigField
	cfg        agentmodel.Config
	cursor     int
}

func New(domain, label, configPath string, fields []manifest.ConfigField) DomainConfigModel {
	m := DomainConfigModel{domain: domain, label: label, configPath: configPath, fields: fields}
	if cfg, err := agentmodel.Load(configPath); err == nil && cfg != nil {
		m.cfg = *cfg
	}
	return m
}

func (m DomainConfigModel) Init() tea.Cmd { return nil }

func (m DomainConfigModel) Update(msg tea.Msg) (nav.ChildModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.fields)-1 {
				m.cursor++
			}
		case "enter", " ":
			return m.activate()
		case "q", "esc":
			return m, m.saveAndBack()
		case "ctrl+c":
			return m, func() tea.Msg { return nav.QuitMsg{} }
		}
	}
	return m, nil
}

func (m DomainConfigModel) activate() (nav.ChildModel, tea.Cmd) {
	if len(m.fields) == 0 {
		return m, nil
	}
	f := m.fields[m.cursor]
	switch f.Type {
	case "agent-models":
		// persist any inline edits first, then open the model screen for this domain
		cfg, path, domain := m.cfg, m.configPath, m.domain
		return m, func() tea.Msg {
			_ = agentmodel.Save(path, &cfg)
			return nav.OpenModelsMsg{Domain: domain}
		}
	case "bool":
		m.toggleBool(f)
	case "enum":
		m.cycleEnum(f)
	}
	return m, nil
}

// typedBool returns the persisted pointer for a typed bool field (strictTdd /
// flagDecisionGaps) and whether the key is a typed one. Non-typed bools live in
// the generic Settings map instead.
func (m DomainConfigModel) typedBool(key string) (*bool, bool) {
	switch key {
	case "strictTdd":
		return m.cfg.DomainStrictTdd(m.domain), true
	case "flagDecisionGaps":
		return m.cfg.DomainFlagDecisionGaps(m.domain), true
	}
	return nil, false
}

func (m *DomainConfigModel) setTypedBool(key string, v bool) bool {
	switch key {
	case "strictTdd":
		m.cfg.SetDomainStrictTdd(m.domain, &v)
		return true
	case "flagDecisionGaps":
		m.cfg.SetDomainFlagDecisionGaps(m.domain, &v)
		return true
	}
	return false
}

func (m *DomainConfigModel) toggleBool(f manifest.ConfigField) {
	v := !m.boolValue(f)
	if m.setTypedBool(f.Key, v) {
		return
	}
	m.cfg.SetDomainSetting(m.domain, f.Key, v)
}

func (m DomainConfigModel) boolValue(f manifest.ConfigField) bool {
	if p, ok := m.typedBool(f.Key); ok {
		if p != nil {
			return *p
		}
	} else if v, ok := m.cfg.DomainSetting(m.domain, f.Key).(bool); ok {
		return v
	}
	if d, ok := f.Default.(bool); ok {
		return d
	}
	return false
}

func (m *DomainConfigModel) cycleEnum(f manifest.ConfigField) {
	if len(f.Options) == 0 {
		return
	}
	cur := m.enumValue(f)
	idx := 0
	for i, o := range f.Options {
		if o == cur {
			idx = i
			break
		}
	}
	m.cfg.SetDomainSetting(m.domain, f.Key, f.Options[(idx+1)%len(f.Options)])
}

func (m DomainConfigModel) enumValue(f manifest.ConfigField) string {
	if v, ok := m.cfg.DomainSetting(m.domain, f.Key).(string); ok {
		return v
	}
	if d, ok := f.Default.(string); ok {
		return d
	}
	if len(f.Options) > 0 {
		return f.Options[0]
	}
	return ""
}

func (m DomainConfigModel) saveAndBack() tea.Cmd {
	cfg, path := m.cfg, m.configPath
	return func() tea.Msg {
		_ = agentmodel.Save(path, &cfg)
		return nav.BackMsg{}
	}
}

func (m DomainConfigModel) View() string {
	var sb strings.Builder
	sb.WriteString(styles.Title.Render("Config — "+m.label) + "\n\n")

	if len(m.fields) == 0 {
		sb.WriteString(styles.Dimmed.Render("  este dominio no declara configuración") + "\n\n")
	}
	for i, f := range m.fields {
		var val string
		switch f.Type {
		case "agent-models":
			val = styles.Dimmed.Render("→ enter")
		case "bool":
			if m.boolValue(f) {
				val = "[x]"
			} else {
				val = "[ ]"
			}
		case "enum":
			val = m.enumValue(f)
		}
		line := fmt.Sprintf("%-22s %s", f.Label, val)
		if i == m.cursor {
			sb.WriteString(styles.Selected.Render("> "+line) + "\n")
		} else {
			sb.WriteString("  " + line + "\n")
		}
	}

	sb.WriteString("\n" + styles.Footer.Render("↑↓ campo  enter abrir/toggle  q/esc guardar  ctrl+c descartar"))
	return sb.String()
}
