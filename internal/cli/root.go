package cli

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
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
		Example: `  # Reportar estado del entorno
  matecito-ai verify

  # Instalar todo lo que falte (prereqs detectados auto)
  matecito-ai install --dry-run
  matecito-ai install`,
	}

	root.CompletionOptions.DisableDefaultCmd = true

	root.AddGroup(
		&cobra.Group{ID: "setup", Title: "Setup:"},
		&cobra.Group{ID: "status", Title: "Diagnóstico:"},
	)

	root.AddCommand(NewVerifyCmd())
	root.AddCommand(NewInitCmd())
	root.AddCommand(NewInstallCmd())
	return root
}
