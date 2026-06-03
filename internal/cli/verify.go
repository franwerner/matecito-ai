package cli

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/checks/claudemd"
	"github.com/franwerner/matecito-ai/internal/checks/codegraph"
	"github.com/franwerner/matecito-ai/internal/checks/context7"
	"github.com/franwerner/matecito-ai/internal/checks/drawio"
	"github.com/franwerner/matecito-ai/internal/checks/engram"
	"github.com/franwerner/matecito-ai/internal/checks/permissions"
	"github.com/franwerner/matecito-ai/internal/checks/prereqs"
	"github.com/franwerner/matecito-ai/internal/checks/proofshot"
	"github.com/franwerner/matecito-ai/internal/checks/sdd"
	"github.com/franwerner/matecito-ai/internal/render"
)

func NewVerifyCmd() *cobra.Command {
	var sddDir string

	cmd := &cobra.Command{
		Use:     "verify",
		GroupID: "status",
		Short:   "Reporta el estado del entorno (prereqs + Engram + CodeGraph + context7 + drawio + proofshot + SDD)",
		Long:    "verify chequea prerequisites del sistema, el estado de los componentes\nregistrados y la coherencia entre el SDD forkeado y los MCP reales.",
		Example: `  # Estado del entorno (default sdd-dir: ~/.claude/agents)
  matecito-ai verify

  # Cross-check contra el fork local antes de install
  matecito-ai verify --sdd-dir ./payload/agents`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pre := prereqs.All()
			eng := engram.All()
			cg := codegraph.All()
			c7 := context7.All()
			dr := drawio.All()
			ps := proofshot.All()
			integ := claudemd.All()
			perm := permissions.All()
			sx := sdd.CrossCheck(sddDir)

			render.Section(os.Stdout, "Prerequisites", pre)
			render.Section(os.Stdout, "Engram", eng)
			render.Section(os.Stdout, "CodeGraph", cg)
			render.Section(os.Stdout, "context7", c7)
			render.Section(os.Stdout, "drawio", dr)
			render.Section(os.Stdout, "proofshot", ps)
			render.Section(os.Stdout, "Integración con Claude Code", integ)
			render.Section(os.Stdout, "Auto-aprobación de tools (settings.json)", perm)
			render.Section(os.Stdout, "Cross-check SDD ↔ MCP ("+sddDir+")", sx)

			all := make([]check.Result, 0, len(pre)+len(eng)+len(cg)+len(c7)+len(dr)+len(ps)+len(integ)+len(perm)+len(sx))
			all = append(all, pre...)
			all = append(all, eng...)
			all = append(all, cg...)
			all = append(all, c7...)
			all = append(all, dr...)
			all = append(all, ps...)
			all = append(all, integ...)
			all = append(all, perm...)
			all = append(all, sx...)

			if code := render.Summary(os.Stdout, all); code != 0 {
				os.Exit(code)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&sddDir, "sdd-dir", defaultSDDDir(),
		"Directorio donde viven los agentes del SDD (sdd-*.md)")
	return cmd
}

func defaultSDDDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".claude/agents"
	}
	return filepath.Join(home, ".claude", "agents")
}
