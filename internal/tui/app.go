package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	"github.com/franwerner/matecito-ai/internal/manifest"
	pkgsync "github.com/franwerner/matecito-ai/internal/setup/sync"
	"github.com/franwerner/matecito-ai/internal/tui/header"
	"github.com/franwerner/matecito-ai/internal/tui/screens/config"
	"github.com/franwerner/matecito-ai/internal/tui/screens/domainconfig"
	"github.com/franwerner/matecito-ai/internal/tui/screens/domains"
	"github.com/franwerner/matecito-ai/internal/tui/screens/menu"
	"github.com/franwerner/matecito-ai/internal/tui/screens/sddmodel"
	tuisync "github.com/franwerner/matecito-ai/internal/tui/screens/sync"
	"github.com/franwerner/matecito-ai/internal/tui/screens/verify"
)

const (
	releaseCheckTimeout = 5 * time.Second
	syncCheckInterval   = 24 * time.Hour
)

// syncCheckMsg carries the result of the startup sync-detect + plan step.
// Se emite solo cuando ShouldCheck fue true; matecitoTag alimenta el badge del
// header y actions determina si navegar a ScreenSync. states se pasa a SyncModel
// vía PreDetected para evitar un segundo Detect() al ejecutar.
type syncCheckMsg struct {
	matecitoTag string // tag latest de matecito-ai para el badge (vacío si offline)
	states      []pkgsync.ComponentState
	actions     []pkgsync.SyncAction
	err         error
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
	scope           agentmodel.Scope
	syncOpts        pkgsync.Options
	updateAvailable bool
	// reexecRequested se setea cuando el binario fue auto-actualizado y hay que
	// re-ejecutar. tui.Run lo lee del modelo final que devuelve p.Run() y hace
	// el ReExec recién después, con la terminal ya restaurada por bubbletea.
	reexecRequested bool
}

// NewAppModel builds the initial AppModel on the menu screen.
// version and ctx are set at startup; the header LatestTag is filled
// asynchronously by the release-check cmd returned from Init.
func NewAppModel(version, globalConfigPath string, ctx ProjectContext) AppModel {
	return AppModel{
		screen:           ScreenMenu,
		child:            menu.New(false),
		hdr:              header.Header{Version: version, ProjectName: ctx.Name, InProject: ctx.InProject},
		ctx:              ctx,
		globalConfigPath: globalConfigPath,
		scope:            agentmodel.ScopeGlobal,
		syncOpts:         pkgsync.Options{SelfVersion: version, Timeout: releaseCheckTimeout},
	}
}

// Init inicia el child del menú y dispara el chequeo de versiones/sync según
// el throttle. Si el throttle indica que no es momento (< 24h desde el último
// check), solo arranca el menú sin tráfico de red. Si es momento, lanza UN
// cmd unificado que obtiene el tag de matecito-ai (para el badge) y el plan
// de sync (para decidir si navegar a ScreenSync).
func (m AppModel) Init() tea.Cmd {
	statePath, err := pkgsync.SyncStatePath()
	if err != nil {
		// si no se puede obtener la ruta, salteamos el check esta vez
		return m.child.Init()
	}

	lastCheck, _ := pkgsync.LoadSyncState(statePath)
	if !pkgsync.ShouldCheck(time.Now(), lastCheck, syncCheckInterval) {
		// throttle activo: arrancar directo al menú sin red
		return m.child.Init()
	}

	// throttle vencido: lanzar chequeo unificado (badge + plan de sync)
	return tea.Batch(m.child.Init(), unifiedCheckCmd(m.syncOpts, statePath))
}

// unifiedCheckCmd obtiene en paralelo el tag latest de matecito-ai (para el
// badge del header) y el plan de sync (para decidir si navegar a ScreenSync).
// Al finalizar persiste el timestamp del check en statePath.
// Cualquier error de red → tag vacío + plan vacío; nunca bloquea el arranque.
func unifiedCheckCmd(opts pkgsync.Options, statePath string) tea.Cmd {
	return func() tea.Msg {
		// Detect hace las tres llamadas de red (self, engram, codegraph) y
		// el diff de deploy; cada fuente falla de forma independiente.
		states, _ := pkgsync.Detect(opts)
		actions := pkgsync.PlanSync(states)

		// Persistir el timestamp del check independientemente del resultado.
		_ = pkgsync.SaveSyncState(statePath, time.Now())

		// Extraer el tag latest de matecito-ai para alimentar el badge del header.
		var matecitoTag string
		for _, s := range states {
			if s.Name == "matecito-ai" && !s.Unknown {
				matecitoTag = s.LatestVersion
				break
			}
		}

		// Filtrar solo las acciones no-skip para decidir si hace falta sync.
		active := make([]pkgsync.SyncAction, 0, len(actions))
		for _, a := range actions {
			if a.Kind != pkgsync.ActionSkip {
				active = append(active, a)
			}
		}

		return syncCheckMsg{matecitoTag: matecitoTag, states: states, actions: active}
	}
}

// Update intercepts router-level messages before delegating to the active child.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case syncCheckMsg:
		// alimentar el badge con el tag latest de matecito-ai
		if msg.matecitoTag != "" {
			m.hdr.LatestTag = agentmodel.NormalizeVersion(msg.matecitoTag)
		}
		// inyectar los estados ya detectados para que SyncModel no llame a Detect de nuevo
		// cuando el usuario entre a "Actualizar" manualmente.
		m.syncOpts.PreDetected = msg.states
		m.updateAvailable = len(msg.actions) > 0
		// Si hay acciones pendientes y seguimos en el menú, reconstruir el child
		// para que la vista refleje el aviso de actualización disponible.
		if m.screen == ScreenMenu {
			m.child = menu.New(m.updateAvailable)
		}
		return m, nil

	case NavigateMsg:
		child, cmd := m.buildChild(msg.To)
		m.screen = msg.To
		m.child = child
		return m, tea.Batch(child.Init(), cmd)

	case OpenDomainConfigMsg:
		child := m.buildDomainConfig(msg.Domain)
		m.child = child
		return m, child.Init()

	case OpenModelsMsg:
		child := m.buildModels(msg.Domain)
		m.child = child
		return m, child.Init()

	case BackMsg:
		m.screen = ScreenMenu
		m.child = menu.New(m.updateAvailable)
		return m, m.child.Init()

	case QuitMsg:
		return m, tea.Quit

	case ReExecMsg:
		// Registrar el re-exec y salir limpio. El ReExec real corre en tui.Run
		// después de p.Run(), cuando bubbletea ya restauró la terminal.
		m.reexecRequested = true
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
			if m.screen == ScreenSddModel || m.screen == ScreenConfig {
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
		return tuisync.New(m.syncOpts), nil
	case ScreenSync:
		return tuisync.New(m.syncOpts), nil
	case ScreenVerify:
		return verify.New(), nil
	case ScreenConfig:
		return config.New(config.ProjectContext{
			InProject: m.ctx.InProject,
			RepoRoot:  m.ctx.RepoRoot,
		}, m.scope), nil
	case ScreenSddModel:
		return m.buildModels(agentmodel.DefaultDomain), nil
	case ScreenDomains:
		// domains are a per-user install concept → always the global config.
		return domains.New(m.globalConfigPath), nil
	default:
		return menu.New(m.updateAvailable), nil
	}
}

// buildDomainConfig constructs the generic per-domain config screen, rendered
// from the domain's manifest config schema.
func (m AppModel) buildDomainConfig(domain string) ChildModel {
	label := domain
	var fields []manifest.ConfigField
	if _, payloadFS, err := manifest.ResolveFromEnv(); err == nil {
		if mf, e := manifest.Load(payloadFS, domain); e == nil {
			fields = mf.Config
			if mf.Label != "" {
				label = mf.Label
			}
		}
	}
	return domainconfig.New(domain, label, m.globalConfigPath, fields)
}

// buildModels constructs the model-per-agent screen scoped to a domain's agents.
func (m AppModel) buildModels(domain string) ChildModel {
	configPath, _ := agentmodel.ConfigPathForScope(m.scope, m.ctx.RepoRoot)
	globalCfg, _ := agentmodel.Load(m.globalConfigPath)
	var projectCfg *agentmodel.Config
	if m.ctx.InProject {
		projectCfg, _ = agentmodel.Load(agentmodel.ProjectConfigPath(m.ctx.RepoRoot))
	}
	return sddmodel.New(globalCfg, projectCfg, configPath, m.scope, domain, m.domainAgents(domain))
}

// domainAgents returns a domain's agents for its model config, discovered from
// the payload (every domains/<domain>/agents/*.md) — development included, so the
// list stays in sync with the deployed agents without a hardcoded roster.
func (m AppModel) domainAgents(domain string) []string {
	if _, payloadFS, err := manifest.ResolveFromEnv(); err == nil {
		return manifest.DomainAgents(payloadFS, domain)
	}
	return nil
}
