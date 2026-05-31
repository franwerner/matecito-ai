package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	"github.com/franwerner/matecito-ai/internal/setup/releasedl"
	"github.com/franwerner/matecito-ai/internal/tui/header"
	"github.com/franwerner/matecito-ai/internal/tui/screens/config"
	"github.com/franwerner/matecito-ai/internal/tui/screens/install"
	"github.com/franwerner/matecito-ai/internal/tui/screens/menu"
	"github.com/franwerner/matecito-ai/internal/tui/screens/sddmodel"
	"github.com/franwerner/matecito-ai/internal/tui/screens/tdd"
	"github.com/franwerner/matecito-ai/internal/tui/screens/update"
	"github.com/franwerner/matecito-ai/internal/tui/screens/verify"
)

const releaseCheckTimeout = 5 * time.Second

// releaseCheckMsg carries the result of the async GitHub release check.
type releaseCheckMsg struct {
	tag string
	err error
}

// AppModel is the top-level bubbletea model. It owns the screen router,
// the active child model, the persistent header, and the project context.
type AppModel struct {
	screen           Screen
	child            ChildModel
	hdr              header.Header
	ctx              ProjectContext
	globalConfigPath string
	// scope is the active config scope; resets to ScopeGlobal on each TUI open.
	scope agentmodel.Scope
}

// NewAppModel builds the initial AppModel on the menu screen.
// version and ctx are set at startup; the header LatestTag is filled
// asynchronously by the release-check cmd returned from Init.
func NewAppModel(version, globalConfigPath string, ctx ProjectContext) AppModel {
	return AppModel{
		screen:           ScreenMenu,
		child:            menu.New(),
		hdr:              header.Header{Version: version, ProjectName: ctx.Name, InProject: ctx.InProject},
		ctx:              ctx,
		globalConfigPath: globalConfigPath,
		scope:            agentmodel.ScopeGlobal,
	}
}

// Init starts the menu child and fires the async release check in parallel.
func (m AppModel) Init() tea.Cmd {
	return tea.Batch(m.child.Init(), releaseCheckCmd(m.hdr.Version))
}

// releaseCheckCmd returns a tea.Cmd that queries the latest GitHub release
// with a short timeout. On any error it returns an empty tag so the header
// badge stays hidden (R7.5/R7.7 — never blocks TUI startup).
func releaseCheckCmd(currentVersion string) tea.Cmd {
	return func() tea.Msg {
		p, err := releasedl.Detect()
		if err != nil {
			return releaseCheckMsg{err: err}
		}
		rel, err := releasedl.LatestReleaseWithTimeout(releasedl.MatecitoRepo, p, releaseCheckTimeout)
		if err != nil {
			return releaseCheckMsg{err: err}
		}
		return releaseCheckMsg{tag: rel.Tag}
	}
}

// Update intercepts router-level messages before delegating to the active child.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case releaseCheckMsg:
		if msg.err == nil {
			m.hdr.LatestTag = agentmodel.NormalizeVersion(msg.tag)
		}
		return m, nil

	case NavigateMsg:
		child, cmd := m.buildChild(msg.To)
		m.screen = msg.To
		m.child = child
		return m, tea.Batch(child.Init(), cmd)

	case BackMsg:
		m.screen = ScreenMenu
		m.child = menu.New()
		return m, m.child.Init()

	case QuitMsg:
		return m, tea.Quit

	case ToggleScopeMsg:
		// solo cambia el scope cuando el TUI está dentro de un proyecto
		if m.ctx.InProject {
			if m.scope == agentmodel.ScopeGlobal {
				m.scope = agentmodel.ScopeProject
			} else {
				m.scope = agentmodel.ScopeGlobal
			}
			m.hdr.Scope = m.scope
			// reconstruir la pantalla activa cuando es scope-aware
			if m.screen == ScreenSddModel || m.screen == ScreenTdd || m.screen == ScreenConfig {
				child, cmd := m.buildChild(m.screen)
				m.child = child
				return m, tea.Batch(child.Init(), cmd)
			}
		}
		return m, nil

	case tea.KeyMsg:
		// ctrl+c es la escotilla de escape universal: corta desde cualquier
		// pantalla, incluso durante un install/update en curso.
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "s" {
			return m.Update(ToggleScopeMsg{})
		}
	}

	newChild, cmd := m.child.Update(msg)
	m.child = newChild
	return m, cmd
}

// View renders the persistent header above the active child screen.
func (m AppModel) View() string {
	return m.hdr.Render() + "\n" + m.child.View()
}

// buildChild constructs the concrete screen model for the requested Screen.
// It passes globalConfigPath, scope, and ProjectContext where the screen needs them.
func (m AppModel) buildChild(s Screen) (ChildModel, tea.Cmd) {
	switch s {
	case ScreenInstall:
		return install.New(), nil
	case ScreenUpdate:
		return update.New(), nil
	case ScreenVerify:
		return verify.New(), nil
	case ScreenConfig:
		return config.New(config.ProjectContext{
			InProject: m.ctx.InProject,
			RepoRoot:  m.ctx.RepoRoot,
		}, m.scope), nil
	case ScreenSddModel:
		configPath, _ := agentmodel.ConfigPathForScope(m.scope, m.ctx.RepoRoot)
		globalCfg, _ := agentmodel.Load(m.globalConfigPath)
		var projectCfg *agentmodel.Config
		if m.ctx.InProject {
			projectCfg, _ = agentmodel.Load(agentmodel.ProjectConfigPath(m.ctx.RepoRoot))
		}
		return sddmodel.New(globalCfg, projectCfg, configPath, m.scope), nil
	case ScreenTdd:
		configPath, _ := agentmodel.ConfigPathForScope(m.scope, m.ctx.RepoRoot)
		return tdd.New(configPath), nil
	default:
		return menu.New(), nil
	}
}
