package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/franwerner/matecito-ai/internal/install"
)

func NewUpdateCmd() *cobra.Command {
	var only string

	cmd := &cobra.Command{
		Use:     "update",
		GroupID: "setup",
		Short:   "Actualiza matecito-ai y Engram a sus últimas releases de GitHub",
		Long: `update descarga las últimas releases de matecito-ai y Engram desde GitHub,
verifica el SHA256 de cada asset y reemplaza los binarios existentes.

Por default actualiza ambos. Usá --only para acotar a uno.`,
		Example: `  # Actualizar matecito-ai y Engram
  matecito-ai update

  # Solo el binario de matecito-ai
  matecito-ai update --only=self

  # Solo Engram
  matecito-ai update --only=engram`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := install.Options{
				Stdin:  os.Stdin,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}

			doSelf := only == "" || only == "self"
			doEngram := only == "" || only == "engram"

			if !doSelf && !doEngram {
				return fmt.Errorf("--only debe ser 'self' o 'engram' (recibí %q)", only)
			}

			if doEngram {
				fmt.Fprintln(opts.Stdout, "Actualizando Engram…")
				if err := install.InstallEngram(opts); err != nil {
					return err
				}
			}
			if doSelf {
				fmt.Fprintln(opts.Stdout, "Actualizando matecito-ai…")
				if err := install.InstallSelf(opts); err != nil {
					return err
				}
			}

			fmt.Fprintln(opts.Stdout, "Listo.")
			return nil
		},
	}

	cmd.Flags().StringVar(&only, "only", "", "Actualizar solo un target: self | engram")
	return cmd
}
