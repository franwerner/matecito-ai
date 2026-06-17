package cli

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"

	"github.com/franwerner/matecito-ai/internal/render"
	"github.com/franwerner/matecito-ai/internal/tui"
)

var (
	version = "0.1.0-dev"
	commit  = "unknown"
	date    = "unknown"
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			if len(s.Value) >= 7 {
				commit = s.Value[:7]
			}
		case "vcs.time":
			date = s.Value
		case "vcs.modified":
			if s.Value == "true" {
				commit += "-dirty"
			}
		}
	}
}

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "matecito-ai",
		Short:         "CLI de setup del ecosistema matecito-ai",
		Long:          "matecito-ai verifica, inicia e instala las dependencias del ecosistema\n(Engram, CodeGraph, context7) sobre Claude Code, y deploya el fork\ndel SDD a ~/.claude/.",
		Version:       fmt.Sprintf("%s (commit %s, built %s)", version, commit, date),
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		Example: `  # Reportar estado del entorno
  matecito-ai verify

  # Instalar todo lo que falte (prereqs detectados auto)
  matecito-ai install --dry-run
  matecito-ai install`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !render.IsTTY(os.Stdout) {
				return cmd.Help()
			}
			return tui.Run(tui.RunOpts{Version: version})
		},
	}

	root.CompletionOptions.DisableDefaultCmd = true

	root.AddGroup(
		&cobra.Group{ID: "setup", Title: "Setup:"},
		&cobra.Group{ID: "status", Title: "Diagnóstico:"},
	)

	root.AddCommand(NewVerifyCmd())
	root.AddCommand(NewInstallCmd())
	root.AddCommand(NewUpdateCmd())
	root.AddCommand(NewHookCmd())
	return root
}
