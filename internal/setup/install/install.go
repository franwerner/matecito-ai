package install

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/franwerner/matecito-ai/internal/hook"
	"github.com/franwerner/matecito-ai/internal/manifest"
	"github.com/franwerner/matecito-ai/internal/mcp"
	"github.com/franwerner/matecito-ai/internal/platform"
	"github.com/franwerner/matecito-ai/internal/setup/deploy"
	"github.com/franwerner/matecito-ai/internal/setup/releasedl"
	"github.com/franwerner/matecito-ai/internal/setup/settings"
)

type Step struct {
	Name  string
	Plan  string
	Check func() bool
	Run   func() error
}

type Options struct {
	DryRun bool
	Yes    bool
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	// BackupDir es la carpeta donde se respaldan todos los archivos que
	// la corrida de install modifique. La setea Run una sola vez para
	// que todas las funciones de backup compartan el mismo timestamp.
	BackupDir string

	// SelfVersion es la versión actual del binario matecito-ai. Se pasa a
	// sync.Options para que el planificador detecte si el binario está desactualizado
	// sin necesitar importar el paquete cli (evita ciclo de importación).
	SelfVersion string
}

func Run(opts Options) error {
	if opts.Stdin == nil {
		opts.Stdin = os.Stdin
	}
	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}
	if opts.Stderr == nil {
		opts.Stderr = os.Stderr
	}
	if opts.BackupDir == "" {
		bd, err := deploy.BackupDir()
		if err != nil {
			return fmt.Errorf("resolviendo backup dir: %w", err)
		}
		opts.BackupDir = bd
	}

	steps := AllSteps(opts)
	plan := make([]Step, 0, len(steps))
	for _, s := range steps {
		if s.Check() {
			plan = append(plan, s)
		}
	}

	if len(plan) == 0 {
		fmt.Fprintln(opts.Stdout, "Nada para hacer — todo está instalado y actualizado.")
		return nil
	}

	fmt.Fprintln(opts.Stdout, "Plan:")
	for i, s := range plan {
		fmt.Fprintf(opts.Stdout, "  %d. %s\n     %s\n", i+1, s.Name, s.Plan)
	}

	if opts.DryRun {
		fmt.Fprintln(opts.Stdout, "\n(dry-run) no se ejecutó nada.")
		return nil
	}

	if !opts.Yes {
		if !confirm(opts.Stdin, opts.Stdout, "\n¿Ejecutar? [y/N]: ") {
			fmt.Fprintln(opts.Stdout, "Cancelado.")
			return nil
		}
	}

	if err := backupClaudeJSON(opts.BackupDir); err != nil {
		return fmt.Errorf("falló el backup de ~/.claude.json: %w", err)
	}

	for i, s := range plan {
		fmt.Fprintf(opts.Stdout, "\n[%d/%d] %s\n", i+1, len(plan), s.Name)
		if err := s.Run(); err != nil {
			fmt.Fprintf(opts.Stderr, "✗ falló: %v\n", err)
			if hasBackup(opts.BackupDir) {
				fmt.Fprintf(opts.Stderr, "Backup intacto en %s\n", opts.BackupDir)
			}
			return fmt.Errorf("install detenido en %q", s.Name)
		}
		fmt.Fprintln(opts.Stdout, "✓ OK")
	}

	if hasBackup(opts.BackupDir) {
		fmt.Fprintf(opts.Stdout, "\nBackup: %s\n", opts.BackupDir)
	}
	fmt.Fprintln(opts.Stdout, "\nListo. Verificá con: matecito-ai verify")
	return nil
}

// hasBackup devuelve true si la carpeta de backup existe en disco — i.e. si
// alguno de los pasos de la corrida tuvo algo que respaldar.
func hasBackup(backupDir string) bool {
	if backupDir == "" {
		return false
	}
	info, err := os.Stat(backupDir)
	return err == nil && info.IsDir()
}

// mcpDef is the single source of truth for one MCP server: how to register it
// (step) and the permission tool pattern it needs auto-approved (toolPattern).
// An empty toolPattern means "derive by convention" (mcp__<name>__*); only
// servers that break the convention (engram, a Claude Code plugin) set it.
type mcpDef struct {
	step        func(Options) Step
	toolPattern string
}

// mcpRegistry maps a manifest MCP name to its descriptor. Nothing is installed
// or auto-approved unless an active domain's manifest (domains/<id>/manifest.json)
// declares the name — there is no global/base MCP.
var mcpRegistry = map[string]mcpDef{
	"engram":    {step: engramMCPStep, toolPattern: "mcp__plugin_engram_engram__*"},
	"context7":  {step: context7MCPStep},
	"codegraph": {step: codegraphMCPStep},
	"drawio":    {step: drawioMCPStep},
	"debugger":  {step: debuggerMCPStep},
	"figma":     {step: figmaMCPStep},
	"canva":     {step: canvaMCPStep},
}

// defaultMCP is the development set, used only as a safety net when the active
// MCP cannot be resolved from the environment.
var defaultMCP = []string{"engram", "context7", "codegraph", "drawio"}

// permissionPattern returns the auto-approve tool pattern for an MCP name: the
// descriptor's override when set, otherwise the mcp__<name>__* convention.
func permissionPattern(name string) string {
	if def, ok := mcpRegistry[name]; ok && def.toolPattern != "" {
		return def.toolPattern
	}
	return "mcp__" + name + "__*"
}

// AllSteps devuelve los pasos de registro y configuración que install.Run gestiona.
// Los pasos de binarios (engram, codegraph, matecito-ai) y deploy son responsabilidad
// de sync.Sync; AllSteps solo cubre los pasos de MCP y configuración de ~/.claude/.
// No hay MCP base; todos salen de los manifests de los dominios activos.
func AllSteps(opts Options) []Step {
	steps := domainMCPSteps(opts)
	steps = append(steps, claudeMdReferenceStep(opts), mcpPermissionsStep(opts), hooksDeployStep(opts))
	return steps
}

// domainMCPSteps builds the MCP install steps contributed by the active domains'
// manifests. On any resolution error it falls back to defaultMCP so a broken
// config never silently drops MCP setup.
func domainMCPSteps(opts Options) []Step {
	names, err := manifest.ActiveMCPFromEnv()
	if err != nil {
		names = defaultMCP
	}
	steps := make([]Step, 0, len(names))
	for _, name := range names {
		if def, ok := mcpRegistry[name]; ok {
			steps = append(steps, def.step(opts))
		}
	}
	return steps
}

// ActiveMCPPatterns derives the permission patterns to auto-approve from the
// active domains: one per declared MCP (via permissionPattern) plus the always-on
// "Skill" pattern. It is the single place that turns the MCP registry into
// permissions; settings.go stays MCP-agnostic.
func ActiveMCPPatterns() []string {
	names, err := manifest.ActiveMCPFromEnv()
	if err != nil {
		names = defaultMCP
	}
	patterns := make([]string, 0, len(names)+1)
	for _, name := range names {
		patterns = append(patterns, permissionPattern(name))
	}
	return append(patterns, "Skill")
}

// AllBinarySteps devuelve los pasos de binarios para usos que no pasen por sync.Sync.
// No se usa en el flujo normal del comando install, pero se conserva para compatibilidad.
func AllBinarySteps(opts Options) []Step {
	return []Step{
		engramBinaryStep(opts),
		codegraphBinaryStep(opts),
		deployStep(opts),
	}
}

// ConfigStepsPending reporta si algún paso de config del ecosistema (registro de
// MCPs, permissions.allow, referencia @matecito-ai.md) todavía tiene trabajo
// pendiente. Lo usa sync para decidir si el componente "config ecosistema" del
// update necesita reconciliarse.
func ConfigStepsPending(opts Options) bool {
	for _, s := range AllSteps(opts) {
		if s.Check() {
			return true
		}
	}
	return false
}

// ApplyConfigSteps corre los pasos de config del ecosistema que tengan trabajo
// pendiente (gateados por su Check), respaldando ~/.claude.json antes. Es la vía
// por la que sync.Sync reconcilia la config en el flujo de update, igual que lo
// hace install — sin reimprimir plan ni pedir confirmación (sync ya lo hizo).
func ApplyConfigSteps(opts Options) error {
	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}
	if opts.BackupDir == "" {
		bd, err := deploy.BackupDir()
		if err != nil {
			return err
		}
		opts.BackupDir = bd
	}
	if err := backupClaudeJSON(opts.BackupDir); err != nil {
		return err
	}
	for _, s := range AllSteps(opts) {
		if !s.Check() {
			continue
		}
		fmt.Fprintf(opts.Stdout, "  %s\n", s.Name)
		if err := s.Run(); err != nil {
			return fmt.Errorf("%s: %w", s.Name, err)
		}
	}
	return nil
}

func mcpPermissionsStep(opts Options) Step {
	return Step{
		Name: "Auto-aprobación de tools del ecosistema (settings.json)",
		Plan: "agregar patrones del ecosistema (MCP + Skill) a permissions.allow en ~/.claude/settings.json (no toca defaultMode ni Bash/Write/Edit)",
		Check: func() bool {
			doc, err := settings.Load()
			if err != nil {
				return false
			}
			return len(settings.MissingPatterns(settings.AllowList(doc), ActiveMCPPatterns())) > 0
		},
		Run: func() error {
			path, err := settings.Path()
			if err != nil {
				return err
			}
			doc, err := settings.Load()
			if err != nil {
				return err
			}
			if !settings.Merge(doc, ActiveMCPPatterns()) {
				return nil
			}
			if err := backupSettings(path, opts.BackupDir); err != nil {
				return err
			}
			out, err := json.MarshalIndent(doc, "", "  ")
			if err != nil {
				return err
			}
			out = append(out, '\n')
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			return os.WriteFile(path, out, 0o644)
		},
	}
}

// activeHookEntries converts the active domains' registered hooks into the
// HookEntry type that settings.ReconcileHooks consumes. Id is forwarded so
// matecito-owned handlers carry their identity marker.
func activeHookEntries() []settings.HookEntry {
	hooks, err := hook.ActiveHooks()
	if err != nil {
		return nil
	}
	entries := make([]settings.HookEntry, 0, len(hooks))
	for _, h := range hooks {
		entries = append(entries, settings.HookEntry{
			Event:   h.Event,
			Matcher: h.Matcher,
			Command: h.Command(),
			If:      h.If,
			Timeout: h.Timeout,
			Id:      h.Id,
		})
	}
	return entries
}

// hooksDeployStep registers the active domains' hooks in ~/.claude/settings.json
// using identity-based reconciliation. Check returns true when ReconcileHooks
// would change settings.json (computed against the loaded document). Run backs
// up settings.json then reconciles and writes.
func hooksDeployStep(opts Options) Step {
	return Step{
		Name: "Hooks de dominios activos (settings.json)",
		Plan: "reconciliar handlers de hooks por matecitoId en ~/.claude/settings.json (reemplaza stale handlers, preserva handlers de usuario)",
		Check: func() bool {
			doc, err := settings.Load()
			if err != nil {
				return false
			}
			return settings.ReconcileHooks(doc, activeHookEntries())
		},
		Run: func() error {
			path, err := settings.Path()
			if err != nil {
				return err
			}
			doc, err := settings.Load()
			if err != nil {
				return err
			}
			if !settings.ReconcileHooks(doc, activeHookEntries()) {
				return nil
			}
			if err := backupSettings(path, opts.BackupDir); err != nil {
				return err
			}
			out, err := json.MarshalIndent(doc, "", "  ")
			if err != nil {
				return err
			}
			out = append(out, '\n')
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			return os.WriteFile(path, out, 0o644)
		},
	}
}

func backupSettings(settingsPath, backupDir string) error {
	data, err := os.ReadFile(settingsPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(backupDir, "settings.json"), data, 0o644)
}

type deployPrep struct {
	source    string
	ops       []deploy.FileOp
	summary   string
	shouldRun bool
	err       error
}

func prepareDeploy() deployPrep {
	var p deployPrep
	payloadFS, source, err := deploy.ResolvePayloadFS()
	if err != nil {
		p.summary = "no se pudo resolver payload: " + err.Error()
		p.err = err
		return p
	}
	p.source = source

	claudeHome, err := deploy.ClaudeHome()
	if err != nil {
		p.err = err
		return p
	}
	active, _ := manifest.ActiveIDsFromEnv()
	p.ops, err = deploy.Plan(payloadFS, claudeHome, active)
	if err != nil {
		p.err = err
		p.summary = "error planeando deploy: " + err.Error()
		return p
	}
	s := deploy.Summarize(p.ops)
	p.summary = fmt.Sprintf("desde %s: %d nuevos, %d cambiados (%d sin cambio)",
		source, s.New, s.Changed, s.Same)
	p.shouldRun = s.New+s.Changed > 0
	return p
}

func deployStep(opts Options) Step {
	prep := prepareDeploy()

	return Step{
		Name:  "Deploy del fork (payload → ~/.claude/)",
		Plan:  prep.summary,
		Check: func() bool { return prep.shouldRun },
		Run: func() error {
			if prep.err != nil {
				return prep.err
			}

			// re-resolver payload para el Run (el fs.FS de prepareDeploy es
			// un valor de interface no almacenado en el struct)
			payloadFS, _, err := deploy.ResolvePayloadFS()
			if err != nil {
				return fmt.Errorf("no se pudo resolver payload: %w", err)
			}

			claudeHome, err := deploy.ClaudeHome()
			if err != nil {
				return err
			}
			_, err = deploy.Apply(payloadFS, prep.ops, claudeHome, opts.BackupDir)
			if err != nil {
				return err
			}
			s := deploy.Summarize(prep.ops)
			fmt.Fprintf(opts.Stdout, "  %d nuevos, %d cambiados, %d sin cambio\n", s.New, s.Changed, s.Same)
			return nil
		},
	}
}

func engramBinaryStep(opts Options) Step {
	return Step{
		Name: "Engram (binario)",
		Plan: "descargar última release de GitHub → ~/.local/bin/engram (SHA256 verificado)",
		Check: func() bool {
			_, err := exec.LookPath("engram")
			return err != nil
		},
		Run: func() error {
			return InstallEngram(opts)
		},
	}
}

// InstallSelf descarga la última release de matecito-ai desde GitHub,
// verifica SHA256 y reemplaza el binario actualmente en ejecución por la
// nueva versión. Usa os.Executable para resolver dónde escribir.
//
// En Linux y macOS reemplazar el binario en uso es seguro: el kernel mantiene
// la inode abierta y el proceso en curso sigue usando la versión vieja hasta
// que termine; la siguiente invocación ya usa la nueva.
func InstallSelf(opts Options) error {
	plat, err := releasedl.Detect()
	if err != nil {
		return err
	}
	rel, err := releasedl.LatestRelease(releasedl.MatecitoRepo, plat)
	if err != nil {
		return err
	}
	dest, err := os.Executable()
	if err != nil {
		return err
	}
	return releasedl.Download(releasedl.MatecitoRepo, rel, dest, opts.Stdout)
}

// InstallEngram descarga la última release de Engram desde GitHub, verifica
// el checksum SHA256, instala el binario y asegura que la carpeta destino
// esté en PATH. Usada por sync.Sync en el flujo unificado de install/sync.
func InstallEngram(opts Options) error {
	plat, err := releasedl.Detect()
	if err != nil {
		return err
	}
	rel, err := releasedl.LatestRelease(releasedl.EngramRepo, plat)
	if err != nil {
		return err
	}
	dest, err := releasedl.DefaultBinaryPath(releasedl.EngramRepo)
	if err != nil {
		return err
	}
	if err := releasedl.Download(releasedl.EngramRepo, rel, dest, opts.Stdout); err != nil {
		return err
	}
	_, err = platform.Detect().EnsurePathInShell(filepath.Dir(dest), opts.Stdout)
	return err
}

func InstallCodegraph(opts Options) error {
	if _, err := exec.LookPath("npm"); err != nil {
		return errors.New("npm no está instalado")
	}
	if err := ensureUserNpmPrefix(opts); err != nil {
		return err
	}
	if err := runIO(opts, "npm", "install", "-g", "@colbymchenry/codegraph"); err != nil {
		return err
	}
	if _, err := exec.LookPath("codegraph"); err != nil {
		return errors.New("codegraph: npm install terminó pero el binario no quedó en PATH")
	}
	return nil
}

// proofShotPostStepTimeout es el tiempo máximo para el post-step `proofshot install`.
// El post-step descarga un browser; con este cap evitamos bloqueos en modo no-atendido.
const proofShotPostStepTimeout = 5 * time.Minute

// InstallProofshot instala el paquete npm proofshot de forma global y ejecuta
// el post-step `proofshot install` (descarga browser) de manera no-interactiva.
//
// El post-step es best-effort: si falla o se supera proofShotPostStepTimeout,
// se imprime una instrucción manual y se retorna nil (el binario ya quedó instalado).
// Esto refleja la política "se continúa ante errores de binarios" del instalador.
func InstallProofshot(opts Options) error {
	if _, err := exec.LookPath("npm"); err != nil {
		return errors.New("npm no está instalado")
	}
	if err := ensureUserNpmPrefix(opts); err != nil {
		return err
	}
	if err := runIO(opts, "npm", "install", "-g", "proofshot"); err != nil {
		return err
	}
	if _, err := exec.LookPath("proofshot"); err != nil {
		return errors.New("proofshot: npm install terminó pero el binario no quedó en PATH")
	}

	// Post-step: proofshot install descarga el browser del agente. Se ejecuta
	// con stdin cerrado y un timeout acotado para no bloquear en modo --yes.
	ctx, cancel := context.WithTimeout(context.Background(), proofShotPostStepTimeout)
	defer cancel()

	postCmd := exec.CommandContext(ctx, "proofshot", "install")
	postCmd.Stdin = strings.NewReader("") // stdin cerrado → ningún prompt interactivo bloquea
	postCmd.Stdout = opts.Stdout
	postCmd.Stderr = opts.Stderr

	if err := postCmd.Run(); err != nil {
		// El post-step falló o expiró; continuar y avisar para instalación manual.
		fmt.Fprintf(opts.Stdout,
			"proofshot instalado; corré `proofshot install` manualmente para completar el setup del browser\n")
	}
	return nil
}

func UpdateEngramPlugin(opts Options) error {
	if _, err := exec.LookPath("claude"); err != nil {
		return errors.New("claude no está en PATH")
	}
	if err := runIO(opts, "claude", "plugin", "marketplace", "update", "engram"); err != nil {
		return err
	}
	// `plugin update` requires the installed plugin id (plugin@marketplace), not
	// the short name — `claude plugin list` reports it as `engram@engram`.
	return runIO(opts, "claude", "plugin", "update", "engram@engram")
}

func codegraphBinaryStep(opts Options) Step {
	return Step{
		Name: "CodeGraph (binario)",
		Plan: "npm install -g @colbymchenry/codegraph (configura ~/.npm-global si hace falta)",
		Check: func() bool {
			_, err := exec.LookPath("codegraph")
			return err != nil
		},
		Run: func() error {
			return InstallCodegraph(opts)
		},
	}
}

func isSystemPath(p string) bool {
	switch p {
	case "/usr", "/usr/local", "/opt":
		return true
	}
	return strings.HasPrefix(p, "/usr/") || strings.HasPrefix(p, "/opt/")
}

func ensureUserNpmPrefix(opts Options) error {
	out, err := exec.Command("npm", "config", "get", "prefix").CombinedOutput()
	if err != nil {
		return err
	}
	prefix := strings.TrimSpace(string(out))

	if isSystemPath(prefix) {
		// El prefix apunta a una ruta del sistema; reconfiguramos a ~/.npm-global
		// para que `npm install -g` no requiera sudo.
		fmt.Fprintf(opts.Stdout, "  npm prefix actual = %q (system-owned) → reconfigurando a ~/.npm-global\n", prefix)
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		prefix = filepath.Join(home, ".npm-global")
		if err := os.MkdirAll(prefix, 0o755); err != nil {
			return err
		}
		if err := runIO(opts, "npm", "config", "set", "prefix", prefix); err != nil {
			return err
		}
	}

	// Siempre aseguramos que el bin dir del prefix (sea el original user-local o el
	// recién configurado ~/.npm-global) esté en PATH. EnsurePathInShell es idempotente:
	// no agrega entradas duplicadas si el dir ya está en el shell RC.
	binDir := filepath.Join(prefix, "bin")
	os.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	_, err = platform.Detect().EnsurePathInShell(binDir, opts.Stdout)
	return err
}

func engramMCPStep(opts Options) Step {
	return Step{
		Name: "Engram MCP (plugin)",
		Plan: "claude plugin marketplace add Gentleman-Programming/engram && claude plugin install engram@engram",
		Check: func() bool {
			_, ok := mcp.Find("engram")
			return !ok
		},
		Run: func() error {
			if _, err := exec.LookPath("claude"); err != nil {
				return errors.New("claude no está en PATH")
			}
			if err := runIO(opts, "claude", "plugin", "marketplace", "add", "Gentleman-Programming/engram"); err != nil {
				return err
			}
			// Install with the full plugin@marketplace id so the install matches the
			// id `plugin update`/`plugin list` use later (`engram@engram`).
			return runIO(opts, "claude", "plugin", "install", "engram@engram")
		},
	}
}

func codegraphMCPStep(opts Options) Step {
	return Step{
		Name: "CodeGraph MCP",
		Plan: "claude mcp add --scope user codegraph -- codegraph serve --mcp",
		Check: func() bool {
			_, ok := mcp.Find("codegraph")
			return !ok
		},
		Run: func() error {
			if _, err := exec.LookPath("claude"); err != nil {
				return errors.New("claude no está en PATH")
			}
			if _, err := exec.LookPath("codegraph"); err != nil {
				return errors.New("codegraph: binario no encontrado en PATH; instalá codegraph antes de registrar el MCP")
			}
			return runIO(opts, "claude", "mcp", "add", "--scope", "user", "codegraph", "--", "codegraph", "serve", "--mcp")
		},
	}
}

// figmaMCPStep registers the Figma remote MCP (http). Contributed by the design
// domain manifest (mcp: ["figma"]).
func figmaMCPStep(opts Options) Step {
	return Step{
		Name: "Figma MCP (remote)",
		Plan: "claude mcp add --transport http figma https://mcp.figma.com/mcp",
		Check: func() bool {
			_, ok := mcp.Find("figma")
			return !ok
		},
		Run: func() error {
			if _, err := exec.LookPath("claude"); err != nil {
				return errors.New("claude no está en PATH")
			}
			return runIO(opts, "claude", "mcp", "add", "--transport", "http", "figma", "https://mcp.figma.com/mcp")
		},
	}
}

// canvaMCPStep registers the Canva remote MCP (http). Contributed by the design
// domain manifest (mcp: ["canva"]). Canva ships an official hosted MCP server at
// mcp.canva.com/mcp (OAuth per user) — the same remote pattern as Figma. This is
// NOT the @canva/cli MCP (which is for building Canva apps).
func canvaMCPStep(opts Options) Step {
	return Step{
		Name: "Canva MCP (remote)",
		Plan: "claude mcp add --transport http canva https://mcp.canva.com/mcp",
		Check: func() bool {
			_, ok := mcp.Find("canva")
			return !ok
		},
		Run: func() error {
			if _, err := exec.LookPath("claude"); err != nil {
				return errors.New("claude no está en PATH")
			}
			return runIO(opts, "claude", "mcp", "add", "--transport", "http", "canva", "https://mcp.canva.com/mcp")
		},
	}
}

func context7MCPStep(opts Options) Step {
	return Step{
		Name: "context7 MCP",
		Plan: "claude mcp add --scope user context7 -- npx -y @upstash/context7-mcp@latest",
		Check: func() bool {
			_, ok := mcp.Find("context7")
			return !ok
		},
		Run: func() error {
			if _, err := exec.LookPath("claude"); err != nil {
				return errors.New("claude no está en PATH")
			}
			return runIO(opts, "claude", "mcp", "add", "--scope", "user", "context7", "--", "npx", "-y", "@upstash/context7-mcp@latest")
		},
	}
}

func drawioMCPStep(opts Options) Step {
	return Step{
		Name: "draw.io MCP (next-ai-draw-io)",
		Plan: "claude mcp add --scope user drawio -- npx -y @next-ai-drawio/mcp-server@latest",
		Check: func() bool {
			_, ok := mcp.Find("drawio")
			return !ok
		},
		Run: func() error {
			if _, err := exec.LookPath("claude"); err != nil {
				return errors.New("claude no está en PATH")
			}
			return runIO(opts, "claude", "mcp", "add", "--scope", "user", "drawio", "--", "npx", "-y", "@next-ai-drawio/mcp-server@latest")
		},
	}
}

func debuggerMCPStep(opts Options) Step {
	return Step{
		Name: "debugger MCP (mcp-debugger)",
		Plan: "claude mcp add --scope user debugger -- npx -y @debugmcp/mcp-debugger@latest stdio",
		Check: func() bool {
			_, ok := mcp.Find("debugger")
			return !ok
		},
		Run: func() error {
			if _, err := exec.LookPath("claude"); err != nil {
				return errors.New("claude no está en PATH")
			}
			return runIO(opts, "claude", "mcp", "add", "--scope", "user", "debugger", "--", "npx", "-y", "@debugmcp/mcp-debugger@latest", "stdio")
		},
	}
}

// claudeMdMarker es la línea que importa el matecito-ai.md en el CLAUDE.md
// del usuario. Claude Code lee @<archivo>.md como import nativo.
const claudeMdMarker = "@matecito-ai.md"

func claudeMdReferenceStep(opts Options) Step {
	home, _ := os.UserHomeDir()
	claudeMdPath := filepath.Join(home, ".claude", "CLAUDE.md")

	return Step{
		Name: "Referencia en ~/.claude/CLAUDE.md",
		Plan: fmt.Sprintf("prependear `%s` (crea el archivo si no existe; backup automático si ya hay contenido)", claudeMdMarker),
		Check: func() bool {
			data, err := os.ReadFile(claudeMdPath)
			if errors.Is(err, os.ErrNotExist) {
				return true
			}
			if err != nil {
				return false
			}
			return !strings.Contains(string(data), claudeMdMarker)
		},
		Run: func() error {
			existing, err := os.ReadFile(claudeMdPath)
			if errors.Is(err, os.ErrNotExist) {
				if err := os.MkdirAll(filepath.Dir(claudeMdPath), 0o755); err != nil {
					return err
				}
				return os.WriteFile(claudeMdPath, []byte(claudeMdMarker+"\n"), 0o644)
			}
			if err != nil {
				return err
			}
			if err := backupClaudeMd(claudeMdPath, opts.BackupDir); err != nil {
				return err
			}
			newContent := claudeMdMarker + "\n\n" + string(existing)
			return os.WriteFile(claudeMdPath, []byte(newContent), 0o644)
		},
	}
}

func backupClaudeMd(claudeMdPath, backupDir string) error {
	data, err := os.ReadFile(claudeMdPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(backupDir, "CLAUDE.md"), data, 0o644)
}

func runIO(opts Options, bin string, args ...string) error {
	c := exec.Command(bin, args...)
	c.Stdin = opts.Stdin
	c.Stdout = opts.Stdout
	c.Stderr = opts.Stderr
	return c.Run()
}

func backupClaudeJSON(backupDir string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	src := filepath.Join(home, ".claude.json")
	info, err := os.Stat(src)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("%s es un directorio", src)
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return err
	}
	dst := filepath.Join(backupDir, ".claude.json")
	return os.WriteFile(dst, data, 0600)
}

func confirm(in io.Reader, out io.Writer, prompt string) bool {
	fmt.Fprint(out, prompt)
	sc := bufio.NewScanner(in)
	if !sc.Scan() {
		return false
	}
	ans := strings.ToLower(strings.TrimSpace(sc.Text()))
	return ans == "y" || ans == "yes" || ans == "s" || ans == "si" || ans == "sí"
}
