package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/franwerner/matecito-ai/internal/install"
)

func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		GroupID: "setup",
		Short:   "Actualiza Engram a la última release de GitHub",
		Long:    "update fuerza la descarga de la última versión de Engram desde GitHub Releases\ny la instala sobre la existente. Verifica SHA256 antes de reemplazar el binario.",
		Example: `  matecito-ai update`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := install.Options{
				Stdin:  os.Stdin,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}
			fmt.Fprintln(opts.Stdout, "Actualizando Engram…")
			if err := install.InstallEngram(opts); err != nil {
				return err
			}
			fmt.Fprintln(opts.Stdout, "Listo.")
			return nil
		},
	}
	return cmd
}
