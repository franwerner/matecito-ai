package cli

import "github.com/spf13/cobra"

const version = "0.1.0-dev"

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "matecito",
		Short:         "CLI de setup del ecosistema matecito-ai",
		Long:          "matecito verifica, inicia e instala las dependencias del ecosistema matecito-ai\n(Engram, CodeGraph, context7) sobre Claude Code.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(NewVerifyCmd())
	root.AddCommand(NewInitCmd())
	root.AddCommand(NewInstallCmd())
	root.AddCommand(NewDoctorCmd())
	return root
}
