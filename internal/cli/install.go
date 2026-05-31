package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/franwerner/matecito-ai/internal/setup/deploy"
	"github.com/franwerner/matecito-ai/internal/setup/install"
	"github.com/franwerner/matecito-ai/internal/setup/sync"
)

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
			backupDir, err := deploy.BackupDir()
			if err != nil {
				return err
			}

			syncOpts := sync.Options{
				SelfVersion: version,
				Stdin:       os.Stdin,
				Stdout:      os.Stdout,
				Stderr:      os.Stderr,
				BackupDir:   backupDir,
			}

			installOpts := install.Options{
				DryRun:    dryRun,
				Yes:       yes,
				Stdin:     os.Stdin,
				Stdout:    os.Stdout,
				Stderr:    os.Stderr,
				BackupDir: backupDir,
			}

			// Detectar el estado de binarios + deploy una sola vez.
			states, _ := sync.Detect(syncOpts)
			syncActions := sync.PlanSync(states)

			activeSyncActions := make([]sync.SyncAction, 0, len(syncActions))
			for _, a := range syncActions {
				if a.Kind != sync.ActionSkip {
					activeSyncActions = append(activeSyncActions, a)
				}
			}

			// Evaluar cuáles pasos de MCP/config necesitan ejecutarse.
			allInstallSteps := install.AllSteps(installOpts)
			activeInstallSteps := make([]install.Step, 0, len(allInstallSteps))
			for _, s := range allInstallSteps {
				if s.Check() {
					activeInstallSteps = append(activeInstallSteps, s)
				}
			}

			// Nada para hacer: ambas fases están al día.
			if len(activeSyncActions) == 0 && len(activeInstallSteps) == 0 {
				fmt.Fprintln(os.Stdout, "Nada para hacer — todo está instalado y actualizado.")
				return nil
			}

			// Mostrar plan combinado.
			n := 1
			fmt.Fprintln(os.Stdout, "Plan:")
			for _, a := range activeSyncActions {
				verb := "instalar"
				if a.Kind == sync.ActionUpdate {
					verb = "actualizar"
				}
				fmt.Fprintf(os.Stdout, "  %d. %s — %s\n", n, a.Component, verb)
				n++
			}
			for _, s := range activeInstallSteps {
				fmt.Fprintf(os.Stdout, "  %d. %s\n     %s\n", n, s.Name, s.Plan)
				n++
			}

			if dryRun {
				fmt.Fprintln(os.Stdout, "\n(dry-run) no se ejecutó nada.")
				return nil
			}

			if !yes {
				if !confirmInstall(os.Stdin, os.Stdout, "\n¿Ejecutar? [y/N]: ") {
					fmt.Fprintln(os.Stdout, "Cancelado.")
					return nil
				}
			}

			// Ejecutar fase 1 (binarios + deploy) sin prompt interior.
			if len(activeSyncActions) > 0 {
				syncOpts.Yes = true
				syncOpts.PreDetected = states
				syncResult := sync.Sync(syncOpts)
				if syncResult.HasErrors() {
					// los errores ya fueron impresos por Sync; continuar con MCPs.
					_ = syncResult
				}
			}

			// Ejecutar fase 2 (MCPs y configuración) sin prompt interior.
			if len(activeInstallSteps) > 0 {
				installOpts.Yes = true
				if err := install.Run(installOpts); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Muestra el plan sin ejecutar nada")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "No pedir confirmación interactiva")
	return cmd
}

func confirmInstall(in *os.File, out *os.File, prompt string) bool {
	fmt.Fprint(out, prompt)
	sc := bufio.NewScanner(in)
	if !sc.Scan() {
		return false
	}
	ans := strings.ToLower(strings.TrimSpace(sc.Text()))
	return ans == "y" || ans == "yes" || ans == "s" || ans == "si" || ans == "sí"
}
