package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/franwerner/matecito-ai/internal/setup/deploy"
	"github.com/franwerner/matecito-ai/internal/setup/sync"
)

func NewUpdateCmd() *cobra.Command {
	var dryRun bool
	var yes bool

	cmd := &cobra.Command{
		Use:     "update",
		GroupID: "setup",
		Short:   "Actualiza los componentes del ecosistema (binarios, payload y config de MCPs)",
		Long: `update detecta qué binarios están desactualizados o faltan, actualiza el
payload de ~/.claude/ y reconcilia la config del ecosistema (registro de MCPs +
permissions.allow + referencia @matecito-ai.md). Muestra el plan antes de
ejecutar. Se continúa ante errores por componente.`,
		Example: `  # Preview del plan, no ejecuta nada
  matecito-ai update --dry-run

  # Actualizar sin prompts (CI-friendly)
  matecito-ai update --yes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			backupDir, err := deploy.BackupDir()
			if err != nil {
				return err
			}

			opts := sync.Options{
				SelfVersion: version,
				DryRun:      dryRun,
				Yes:         yes,
				Resume:      sync.ResumeRequested(),
				Stdin:       os.Stdin,
				Stdout:      os.Stdout,
				Stderr:      os.Stderr,
				BackupDir:   backupDir,
			}

			result := sync.Sync(opts)

			// CLI-only trigger: the engine never re-execs itself, so this is
			// the one place a self-replace hands off to the new binary.
			sync.FinishSelfReplace(os.Stdout, os.Stderr, result.SelfReplaced)

			if result.HasErrors() {
				return fmt.Errorf("update: uno o más componentes fallaron durante la sincronización")
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Muestra el plan sin ejecutar nada")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "No pedir confirmación interactiva")
	return cmd
}
