package sync

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/manifest"
	"github.com/franwerner/matecito-ai/internal/render"
	"github.com/franwerner/matecito-ai/internal/setup/deploy"
	"github.com/franwerner/matecito-ai/internal/setup/install"
)

// defaultTimeout es el timeout por defecto para las consultas de versiones remotas.
// Coincide con releaseCheckTimeout de la TUI (~5s) para mantener coherencia.
const defaultTimeout = 5 * time.Second

// Options configura el comportamiento de Detect y Sync.
type Options struct {
	// SelfVersion es la versión actual del binario matecito-ai, inyectada por cli
	// para evitar un ciclo de importación (cli → sync → cli).
	SelfVersion string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	// Timeout controla cada petición de versión remota (GitHub / npm).
	// Si es cero se usa defaultTimeout.
	Timeout time.Duration

	// DryRun imprime el plan sin ejecutar ninguna acción.
	DryRun bool

	// Yes omite la confirmación interactiva antes de ejecutar.
	Yes bool

	// BackupDir es la carpeta donde deploy.Apply respalda archivos modificados.
	// Si está vacío, Sync resuelve uno vía deploy.BackupDir().
	BackupDir string

	// PreDetected permite inyectar estados ya detectados (e.g. desde la TUI)
	// para que Sync no vuelva a llamar a Detect y evitar el doble round-trip de red.
	// Si es nil, Sync llama a Detect internamente como siempre.
	PreDetected []ComponentState
}

func (o Options) timeout() time.Duration {
	if o.Timeout > 0 {
		return o.Timeout
	}
	return defaultTimeout
}

func (o Options) stdin() io.Reader {
	if o.Stdin != nil {
		return o.Stdin
	}
	return os.Stdin
}

func (o Options) stdout() io.Writer {
	if o.Stdout != nil {
		return o.Stdout
	}
	return os.Stdout
}

func (o Options) stderr() io.Writer {
	if o.Stderr != nil {
		return o.Stderr
	}
	return os.Stderr
}

// Result resume el resultado de una ejecución de Sync.
type Result struct {
	// Actions son las acciones planificadas que se intentaron ejecutar.
	Actions []SyncAction
	// Errors acumula los errores por componente (nil cuando ninguno falló).
	Errors map[string]error
	// SelfReplaced es true cuando el binario matecito-ai fue reemplazado con éxito.
	SelfReplaced bool
}

// HasErrors devuelve true si algún componente falló durante Sync.
func (r Result) HasErrors() bool {
	return len(r.Errors) > 0
}

// ComponentKind distinguishes what a component represents.
type ComponentKind int

const (
	KindSelf         ComponentKind = iota // el propio binario matecito-ai
	KindDepEngram                         // plugin MCP de Engram
	KindDepCodeGraph                      // CLI de CodeGraph
	KindDeploy                            // payload de archivos estáticos
)

// Component identifica un componente gestionado por el sincronizador.
type Component struct {
	Name string
	Kind ComponentKind
}

// ComponentState describe el estado observado de un componente en el sistema.
// Populated by Detect(); consumed (purely) by PlanSync().
//
// CurrentVersion y LatestVersion se usan para self + deps; vacíos para deploy.
// Unknown=true significa que la versión latest no pudo obtenerse; PlanSync lo
// trata como skip para no forzar reinstalaciones cuando la red no está disponible.
type ComponentState struct {
	Name           string
	Present        bool
	CurrentVersion string
	LatestVersion  string
	PayloadChanged bool   // solo relevante para KindDeploy
	PayloadSource  string // "embedded" o ruta local; solo relevante para KindDeploy
	Pending        bool   // reconciliación pendiente sin semántica de versión (config ecosistema)
	Unknown        bool   // true cuando latest no pudo obtenerse (offline / error)
}

// configComponent es el componente coarse-grained que reconcilia la config del
// ecosistema (registro de MCPs + permissions.allow + referencia @matecito-ai.md),
// modelado igual que deploy para que el update lo muestre en el plan y lo aplique.
const configComponent = "config ecosistema"

// ActionKind is the action PlanSync assigns to a component.
type ActionKind int

const (
	ActionInstall ActionKind = iota
	ActionUpdate
	ActionSkip
)

// SyncAction pairs a component name with the planned action.
type SyncAction struct {
	Component string
	Kind      ActionKind
}

// PlanSync devuelve las acciones recomendadas para una lista de estados de
// componentes. Es una función pura: sin I/O, sin efectos secundarios.
//
// Reglas de decisión (en orden de precedencia):
//  1. Unknown=true → skip (no se conoce la versión latest; no bloquear al usuario).
//  2. !Present → install.
//  3. PayloadChanged (deploy) → update.
//  4. LatestVersion vacío → skip (sin referencia para comparar).
//  5. NormalizeVersion(Current) != NormalizeVersion(Latest) → update (componente desactualizado).
//  6. Resto → skip (ya está actualizado).
func PlanSync(states []ComponentState) []SyncAction {
	actions := make([]SyncAction, 0, len(states))
	for _, s := range states {
		actions = append(actions, SyncAction{
			Component: s.Name,
			Kind:      decide(s),
		})
	}
	return actions
}

// Detect realiza todas las operaciones de I/O necesarias para construir el
// estado observado de cada componente gestionado por el sincronizador.
//
// Cada componente se detecta de forma independiente: un error en uno no cancela
// los demás. Cuando la versión latest de un componente no puede obtenerse, ese
// componente queda con Unknown=true y PlanSync lo omitirá (skip).
func Detect(opts Options) ([]ComponentState, error) {
	t := opts.timeout()
	states := make([]ComponentState, 0, 4)

	// --- matecito-ai (self) ---
	selfState := ComponentState{
		Name:           "matecito-ai",
		Present:        true,
		CurrentVersion: opts.SelfVersion,
	}
	latestSelf, err := fetchLatestMatecito(t)
	if err != nil {
		selfState.Unknown = true
	} else {
		selfState.LatestVersion = latestSelf
	}
	states = append(states, selfState)

	// Binaries are gated by the active domains' manifests: a binary is detected
	// (and thus installed) only when some active domain declares it. On resolution
	// error, fall back to detecting all (legacy behavior) so a broken config never
	// silently drops a binary.
	activeBins, binErr := manifest.ActiveBinariesFromEnv()
	wantBinary := func(name string) bool {
		return binErr != nil || slices.Contains(activeBins, name)
	}

	// --- Engram ---
	if wantBinary("engram") {
		engramResult := check.RunVersion("engram", "engram", []string{"version"}, false, "")
		engramState := ComponentState{
			Name:           "engram",
			Present:        engramResult.Status != check.StatusMissing,
			CurrentVersion: engramResult.Version,
		}
		latestEngram, err := fetchLatestEngram(t)
		if err != nil {
			engramState.Unknown = true
		} else {
			engramState.LatestVersion = latestEngram
		}
		states = append(states, engramState)
	}

	// --- CodeGraph ---
	if wantBinary("codegraph") {
		cgResult := check.RunVersion("codegraph", "codegraph", []string{"--version"}, false, "")
		cgState := ComponentState{
			Name:           "codegraph",
			Present:        cgResult.Status != check.StatusMissing,
			CurrentVersion: cgResult.Version,
		}
		latestCG, err := fetchLatestCodeGraph(t)
		if err != nil {
			cgState.Unknown = true
		} else {
			cgState.LatestVersion = latestCG
		}
		states = append(states, cgState)
	}

	// --- ProofShot ---
	if wantBinary("proofshot") {
		psResult := check.RunVersion("proofshot", "proofshot", []string{"--version"}, false, "")
		psState := ComponentState{
			Name:           "proofshot",
			Present:        psResult.Status != check.StatusMissing,
			CurrentVersion: psResult.Version,
		}
		latestPS, err := fetchLatestProofshot(t)
		if err != nil {
			psState.Unknown = true
		} else {
			psState.LatestVersion = latestPS
		}
		states = append(states, psState)
	}

	// --- Deploy (payload) ---
	deployState := ComponentState{Name: "deploy"}
	payloadFS, payloadSource, deployErr := deploy.ResolvePayloadFS()
	if deployErr == nil {
		deployState.PayloadSource = payloadSource
		claudeHome, homeErr := deploy.ClaudeHome()
		if homeErr == nil {
			active, _ := manifest.ActiveIDsFromEnv()
			ops, planErr := deploy.Plan(payloadFS, claudeHome, active)
			if planErr == nil {
				s := deploy.Summarize(ops)
				deployState.Present = true
				deployState.PayloadChanged = s.New+s.Changed > 0
			}
		}
	}
	// deploy no tiene semántica de versión; Unknown queda false (se evalúa solo por PayloadChanged)
	states = append(states, deployState)

	// --- Config del ecosistema (MCPs + permisos + referencia CLAUDE.md) ---
	// Reconcilia los pasos de install.AllSteps. Va último para que los binarios y
	// el deploy ya estén aplicados cuando se registren los MCPs y se ajusten los
	// permisos. Esto hace que `update` deje el entorno consistente al sumar un MCP
	// nuevo (p.ej. drawio), no solo `install`.
	configState := ComponentState{
		Name:    configComponent,
		Present: true,
		Pending: install.ConfigStepsPending(install.Options{}),
	}
	states = append(states, configState)

	return states, nil
}

// Sync ejecuta el ciclo completo: Detect → PlanSync → ejecutar acciones
// continue-on-error. Un fallo en un componente no cancela los demás.
//
// Si el binario matecito-ai fue reemplazado con éxito y el proceso corre en
// una TTY Unix, intenta exec(2) para re-lanzarse con el nuevo binario.
// En entornos no-TTY o Windows imprime un aviso y retorna sin exec.
func Sync(opts Options) Result {
	out := opts.stdout()
	errOut := opts.stderr()

	var states []ComponentState
	if len(opts.PreDetected) > 0 {
		// Reutilizar los estados ya detectados para evitar un segundo round-trip de red.
		states = opts.PreDetected
	} else {
		var err error
		states, err = Detect(opts)
		if err != nil {
			// Detect nunca retorna error hoy, pero respetamos el contrato.
			fmt.Fprintf(errOut, "sync: error detectando estado: %v\n", err)
		}
	}

	actions := PlanSync(states)
	result := Result{Actions: actions}

	active := make([]SyncAction, 0, len(actions))
	for _, a := range actions {
		if a.Kind != ActionSkip {
			active = append(active, a)
		}
	}

	if len(active) == 0 {
		fmt.Fprintln(out, "Nada para hacer — todo está instalado y actualizado.")
		return result
	}

	// Build a source lookup so the deploy entry can show its payload origin.
	sourceByComponent := make(map[string]string, len(states))
	for _, s := range states {
		if s.PayloadSource != "" {
			sourceByComponent[s.Name] = s.PayloadSource
		}
	}

	// Mostrar el plan antes de ejecutar (dry-run lo imprime y sale).
	fmt.Fprintln(out, "Plan:")
	for i, a := range active {
		verb := "instalar"
		if a.Kind == ActionUpdate {
			verb = "actualizar"
		}
		fmt.Fprintf(out, "  %d. %s — %s\n", i+1, a.Component, verb)
		if src, ok := sourceByComponent[a.Component]; ok {
			fmt.Fprintf(out, "     payload: %s\n", src)
		}
	}

	if opts.DryRun {
		fmt.Fprintln(out, "\n(dry-run) no se ejecutó nada.")
		return result
	}

	if !opts.Yes {
		if !confirmSync(opts.stdin(), out, "\n¿Ejecutar? [y/N]: ") {
			fmt.Fprintln(out, "Cancelado.")
			return result
		}
	}

	backupDir := opts.BackupDir
	if backupDir == "" {
		bd, bdErr := deploy.BackupDir()
		if bdErr != nil {
			fmt.Fprintf(errOut, "sync: resolviendo backup dir: %v\n", bdErr)
		} else {
			backupDir = bd
		}
	}

	installOpts := install.Options{
		Stdout: out,
		Stderr: errOut,
	}

	for _, a := range active {
		fmt.Fprintf(out, "\n[%s] %s\n", actionLabel(a.Kind), a.Component)
		var runErr error
		switch a.Component {
		case "matecito-ai":
			runErr = install.InstallSelf(installOpts)
			if runErr == nil {
				result.SelfReplaced = true
			}
		case "engram":
			if a.Kind == ActionInstall {
				runErr = install.InstallEngram(installOpts)
			} else {
				runErr = install.UpdateEngramPlugin(installOpts)
			}
		case "codegraph":
			runErr = install.InstallCodegraph(installOpts)
		case "proofshot":
			runErr = install.InstallProofshot(installOpts)
		case "deploy":
			payloadFS, _, fsErr := deploy.ResolvePayloadFS()
			if fsErr != nil {
				runErr = fsErr
				break
			}
			claudeHome, homeErr := deploy.ClaudeHome()
			if homeErr != nil {
				runErr = homeErr
				break
			}
			active, _ := manifest.ActiveIDsFromEnv()
			ops, planErr := deploy.Plan(payloadFS, claudeHome, active)
			if planErr != nil {
				runErr = planErr
				break
			}
			_, runErr = deploy.Apply(payloadFS, ops, claudeHome, backupDir)
		case configComponent:
			runErr = install.ApplyConfigSteps(install.Options{
				Stdout:    out,
				Stderr:    errOut,
				BackupDir: backupDir,
			})
		}

		if runErr != nil {
			if result.Errors == nil {
				result.Errors = make(map[string]error)
			}
			result.Errors[a.Component] = runErr
			fmt.Fprintf(errOut, "✗ falló %s: %v\n", a.Component, runErr)
		} else {
			fmt.Fprintln(out, "✓ OK")
		}
	}

	// Re-exec solo si el binario fue reemplazado, el proceso es un TTY Unix y
	// no hubo error en el propio paso de self-update.
	if result.SelfReplaced && runtime.GOOS != "windows" && render.IsTTY(out) {
		if execErr := ReExec(); execErr != nil {
			fmt.Fprintf(errOut, "sync: re-exec falló: %v\n", execErr)
		}
		// Si ReExec retorna (no debería en Unix salvo error), continuamos.
	} else if result.SelfReplaced {
		// Entorno no-TTY o Windows: avisar al usuario que reinicie manualmente.
		self, _ := exec.LookPath(os.Args[0])
		if self == "" {
			self = "matecito-ai"
		}
		fmt.Fprintf(out, "matecito-ai actualizado — re-ejecutá %s para usar la nueva versión.\n", self)
	}

	return result
}

func actionLabel(k ActionKind) string {
	switch k {
	case ActionInstall:
		return "instalar"
	case ActionUpdate:
		return "actualizar"
	}
	return "skip"
}

func confirmSync(in io.Reader, out io.Writer, prompt string) bool {
	fmt.Fprint(out, prompt)
	sc := bufio.NewScanner(in)
	if !sc.Scan() {
		return false
	}
	ans := strings.ToLower(strings.TrimSpace(sc.Text()))
	return ans == "y" || ans == "yes" || ans == "s" || ans == "si" || ans == "sí"
}

func decide(s ComponentState) ActionKind {
	if s.Unknown {
		return ActionSkip
	}
	if !s.Present {
		return ActionInstall
	}
	if s.PayloadChanged {
		return ActionUpdate
	}
	if s.Pending {
		return ActionUpdate
	}
	if s.LatestVersion == "" {
		return ActionSkip
	}
	if agentmodel.NormalizeVersion(s.CurrentVersion) != agentmodel.NormalizeVersion(s.LatestVersion) {
		return ActionUpdate
	}
	return ActionSkip
}
