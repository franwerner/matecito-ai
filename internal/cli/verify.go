package cli

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/checks/claudemd"
	"github.com/franwerner/matecito-ai/internal/checks/codegraph"
	"github.com/franwerner/matecito-ai/internal/checks/context7"
	"github.com/franwerner/matecito-ai/internal/checks/debugger"
	"github.com/franwerner/matecito-ai/internal/checks/drawio"
	"github.com/franwerner/matecito-ai/internal/checks/engram"
	"github.com/franwerner/matecito-ai/internal/checks/hooks"
	"github.com/franwerner/matecito-ai/internal/checks/permissions"
	"github.com/franwerner/matecito-ai/internal/checks/prereqs"
	"github.com/franwerner/matecito-ai/internal/checks/proofshot"
	"github.com/franwerner/matecito-ai/internal/checks/sdd"
	"github.com/franwerner/matecito-ai/internal/manifest"
	"github.com/franwerner/matecito-ai/internal/render"
	"github.com/franwerner/matecito-ai/internal/setup/install"
)

func NewVerifyCmd() *cobra.Command {
	var sddDir string

	cmd := &cobra.Command{
		Use:     "verify",
		GroupID: "status",
		Short:   "Reporta el estado del entorno (prereqs + Engram + CodeGraph + context7 + drawio + debugger + proofshot + hooks + SDD)",
		Long:    "verify chequea prerequisites del sistema, el estado de los componentes\nregistrados y la coherencia entre el SDD forkeado y los MCP reales.",
		Example: `  # Estado del entorno (default sdd-dir: ~/.claude/agents)
  matecito-ai verify

  # Cross-check contra el fork local antes de install
  matecito-ai verify --sdd-dir ./payload/agents`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Gate domain-specific sections by the active domains. On any
			// resolution error, default to showing them (legacy behavior).
			activeMCP, mcpErr := manifest.ActiveMCPFromEnv()
			engramActive := mcpErr != nil || containsString(activeMCP, "engram")
			context7Active := mcpErr != nil || containsString(activeMCP, "context7")
			codegraphActive := mcpErr != nil || containsString(activeMCP, "codegraph")
			drawioActive := mcpErr != nil || containsString(activeMCP, "drawio")
			debuggerActive := mcpErr != nil || containsString(activeMCP, "debugger")
			activeBins, binErr := manifest.ActiveBinariesFromEnv()
			proofshotActive := binErr != nil || containsString(activeBins, "proofshot")
			activeIDs, _, idErr := manifest.ResolveFromEnv()
			devActive := idErr != nil || containsString(activeIDs, "development")
			activeHooks, hooksErr := manifest.ActiveHooksFromEnv()
			hooksActive := hooksErr == nil && len(activeHooks) > 0

			pre := prereqs.All()
			integ := claudemd.All()
			perm := permissions.All(install.ActiveMCPPatterns())

			var eng []check.Result
			if engramActive {
				eng = engram.All()
			}
			var c7 []check.Result
			if context7Active {
				c7 = context7.All()
			}
			var ps []check.Result
			if proofshotActive {
				ps = proofshot.All()
			}
			var cg []check.Result
			if codegraphActive {
				cg = codegraph.All()
			}
			var dr []check.Result
			if drawioActive {
				dr = drawio.All()
			}
			var dbg []check.Result
			if debuggerActive {
				dbg = debugger.All()
			}
			var sx []check.Result
			if devActive {
				sx = sdd.CrossCheck(sddDir)
			}
			var hk []check.Result
			if hooksActive {
				hk = hooks.All()
			}

			render.Section(os.Stdout, "Prerequisites", pre)
			if engramActive {
				render.Section(os.Stdout, "Engram", eng)
			}
			if codegraphActive {
				render.Section(os.Stdout, "CodeGraph", cg)
			}
			if context7Active {
				render.Section(os.Stdout, "context7", c7)
			}
			if drawioActive {
				render.Section(os.Stdout, "drawio", dr)
			}
			if debuggerActive {
				render.Section(os.Stdout, "debugger", dbg)
			}
			if proofshotActive {
				render.Section(os.Stdout, "proofshot", ps)
			}
			render.Section(os.Stdout, "Integración con Claude Code", integ)
			render.Section(os.Stdout, "Auto-aprobación de tools (settings.json)", perm)
			if hooksActive {
				render.Section(os.Stdout, "Hooks de dominios activos", hk)
			}
			if devActive {
				render.Section(os.Stdout, "Cross-check SDD ↔ MCP ("+sddDir+")", sx)
			}

			all := make([]check.Result, 0, len(pre)+len(eng)+len(cg)+len(c7)+len(dr)+len(dbg)+len(ps)+len(integ)+len(perm)+len(hk)+len(sx))
			all = append(all, pre...)
			all = append(all, eng...)
			all = append(all, cg...)
			all = append(all, c7...)
			all = append(all, dr...)
			all = append(all, dbg...)
			all = append(all, ps...)
			all = append(all, integ...)
			all = append(all, perm...)
			all = append(all, hk...)
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

func containsString(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
