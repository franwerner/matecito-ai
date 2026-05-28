package cli

import (
	"github.com/spf13/cobra"

	"github.com/franwerner/matecito-ai/internal/install"
)

func NewInstallCmd() *cobra.Command {
	opts := install.Options{}

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Instala/registra lo que falte (Engram, CodeGraph, MCPs)",
		Long:  "install detecta qué falta, imprime un plan, hace backup de ~/.claude.json,\ny ejecuta los pasos necesarios en orden. Se detiene al primer error.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return install.Run(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Muestra el plan sin ejecutar nada")
	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "No pedir confirmación interactiva")
	return cmd
}
