package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/franwerner/matecito-ai/internal/setup/deploy"
	"github.com/franwerner/matecito-ai/internal/setup/sync"
)

// installRunner holds NewInstallCmd's RunE logic behind injectable seams:
// Stdin/Stdout/Stderr (RunE hardcoded os.Stdin/os.Stdout directly, a
// pre-existing limitation this file's own comments already documented) and
// Detect/Sync (sync.Detect's network round-trip and sync.Sync's real
// execution). newInstallRunner wires the real defaults; tests substitute
// fakes instead, with no change to the command's observable behavior.
type installRunner struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	Detect func(sync.Options) ([]sync.ComponentState, error)
	Sync   func(sync.Options) sync.Result
}

func newInstallRunner() installRunner {
	return installRunner{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Detect: sync.Detect,
		Sync:   sync.Sync,
	}
}

func NewInstallCmd() *cobra.Command {
	var dryRun bool
	var yes bool

	cmd := &cobra.Command{
		Use:     "install",
		GroupID: "setup",
		Short:   "Instala/actualiza lo que falte (binarios, MCPs, fork) y deploya el payload",
		Long: `install detecta qué binarios faltan o están desactualizados, registra los MCPs
necesarios, deploya el fork a ~/.claude/ y hace backup de ~/.claude.json.
Muestra el plan combinado antes de ejecutar. Se continúa ante errores de binarios.`,
		Example: `  # Preview del plan, no ejecuta nada
  matecito-ai install --dry-run

  # Instalación sin prompts (CI-friendly)
  matecito-ai install --yes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return newInstallRunner().run(version, dryRun, yes)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Muestra el plan sin ejecutar nada")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "No pedir confirmación interactiva")
	return cmd
}

// run is NewInstallCmd's RunE body, extracted onto installRunner so it reads
// its injected seams instead of the process's real stdio and sync.Detect/
// sync.Sync — identical behavior, testable without a live network round-trip.
func (r installRunner) run(selfVersion string, dryRun, yes bool) error {
	backupDir, err := deploy.BackupDir()
	if err != nil {
		return err
	}

	syncOpts := sync.Options{
		SelfVersion: selfVersion,
		Resume:      sync.ResumeRequested(),
		Stdin:       r.Stdin,
		Stdout:      r.Stdout,
		Stderr:      r.Stderr,
		BackupDir:   backupDir,
	}

	// Detectar el estado de binarios + deploy + config una sola vez.
	states, _ := r.Detect(syncOpts)
	syncActions := sync.PlanSync(states)

	activeSyncActions := make([]sync.SyncAction, 0, len(syncActions))
	for _, a := range syncActions {
		if a.Kind != sync.ActionSkip {
			activeSyncActions = append(activeSyncActions, a)
		}
	}

	// Nada para hacer.
	if len(activeSyncActions) == 0 {
		fmt.Fprintln(r.Stdout, "Nada para hacer — todo está instalado y actualizado.")
		return nil
	}

	// A resumed run already confirmed once in the original invocation —
	// skip this command's bespoke plan/dry-run/confirm UI entirely.
	if !syncOpts.Resume {
		// Build payload-source and payload-summary lookups for the plan display.
		sourceByComponent := make(map[string]string, len(states))
		summaryByComponent := make(map[string]deploy.Summary, len(states))
		for _, s := range states {
			if s.PayloadSource != "" {
				sourceByComponent[s.Name] = s.PayloadSource
			}
			if s.Name == "deploy" {
				summaryByComponent[s.Name] = s.PayloadSummary
			}
		}

		// Mostrar plan combinado.
		n := 1
		fmt.Fprintln(r.Stdout, "Plan:")
		for _, a := range activeSyncActions {
			verb := "instalar"
			if a.Kind == sync.ActionUpdate {
				verb = "actualizar"
			}
			fmt.Fprintf(r.Stdout, "  %d. %s — %s\n", n, a.Component, verb)
			if src, ok := sourceByComponent[a.Component]; ok {
				fmt.Fprintf(r.Stdout, "     payload: %s\n", src)
			}
			if s, ok := summaryByComponent[a.Component]; ok {
				for _, line := range sync.PayloadPlanLines(s) {
					fmt.Fprintln(r.Stdout, line)
				}
			}
			n++
		}

		if dryRun {
			fmt.Fprintln(r.Stdout, "\n(dry-run) no se ejecutó nada.")
			return nil
		}

		if !yes {
			if !confirmInstall(r.Stdin, r.Stdout, "\n¿Ejecutar? [y/N]: ") {
				fmt.Fprintln(r.Stdout, "Cancelado.")
				return nil
			}
		}

		// Ya mostramos el plan acá — el motor no debe repetirlo.
		syncOpts.PlanShown = true
	}

	// Ejecutar binarios + deploy + config (sin prompt interior; ya confirmamos).
	syncOpts.Yes = true
	syncOpts.PreDetected = states
	syncResult := r.Sync(syncOpts)

	// CLI-only trigger: the engine never re-execs itself, so this is
	// the one place a self-replace hands off to the new binary.
	sync.FinishSelfReplace(r.Stdout, r.Stderr, syncResult.SelfReplaced)

	if err := surfaceCodegraphError(syncResult); err != nil {
		return err
	}

	return nil
}

// surfaceCodegraphError returns a wrapped error when the sync result carries a
// codegraph failure, or nil when codegraph succeeded. Keeping this as a named
// function makes the decision directly testable without running the full command.
func surfaceCodegraphError(result sync.Result) error {
	if cgErr := result.Errors["codegraph"]; cgErr != nil {
		return fmt.Errorf("codegraph: %w", cgErr)
	}
	return nil
}

func confirmInstall(in io.Reader, out io.Writer, prompt string) bool {
	fmt.Fprint(out, prompt)
	sc := bufio.NewScanner(in)
	if !sc.Scan() {
		return false
	}
	ans := strings.ToLower(strings.TrimSpace(sc.Text()))
	return ans == "y" || ans == "yes" || ans == "s" || ans == "si" || ans == "sí"
}
