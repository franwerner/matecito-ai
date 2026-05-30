package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/franwerner/matecito-ai/internal/install"
)

func NewUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "update",
		GroupID: "setup",
		Short:   "Actualiza todas las deps del ecosistema a sus últimas versiones",
		Long: `update actualiza todas las dependencias del ecosistema a su última versión:
matecito-ai, el binario de Engram, el plugin de Engram y CodeGraph.
context7 corre con npx @latest en cada sesión, así que no requiere actualización.`,
		Example: `  # Actualizar todo
  matecito-ai update`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := install.Options{
				Stdin:  os.Stdin,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}

			tasks := []struct {
				name string
				fn   func(install.Options) error
			}{
				{"matecito-ai", install.InstallSelf},
				{"engram (binario)", install.InstallEngram},
				{"engram (plugin)", install.UpdateEngramPlugin},
				{"codegraph", install.InstallCodegraph},
			}

			var failed []string
			for _, t := range tasks {
				fmt.Fprintf(opts.Stdout, "Actualizando %s…\n", t.name)
				if err := t.fn(opts); err != nil {
					fmt.Fprintf(opts.Stderr, "  ✗ %s: %v\n", t.name, err)
					failed = append(failed, t.name)
				}
			}

			fmt.Fprintln(opts.Stdout, "context7 corre con npx @latest en cada sesión — nada que actualizar.")

			if len(failed) > 0 {
				return fmt.Errorf("falló la actualización de: %s", strings.Join(failed, ", "))
			}
			fmt.Fprintln(opts.Stdout, "Listo.")
			return nil
		},
	}
}
