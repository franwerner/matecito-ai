package verify

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

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
	"github.com/franwerner/matecito-ai/internal/setup/install"
	"github.com/franwerner/matecito-ai/internal/tui/nav"
	"github.com/franwerner/matecito-ai/internal/tui/styles"
)

type cluster struct {
	name    string
	results []check.Result
}

type verifyDoneMsg struct{ clusters []cluster }

// VerifyModel runs the same check clusters as the verify CLI command and
// renders their results inside the TUI. It does not modify checks/* in any way.
type VerifyModel struct {
	running  bool
	clusters []cluster
}

func New() VerifyModel {
	return VerifyModel{running: true}
}

func (m VerifyModel) Init() tea.Cmd {
	return runChecks
}

func (m VerifyModel) Update(msg tea.Msg) (nav.ChildModel, tea.Cmd) {
	switch msg := msg.(type) {
	case verifyDoneMsg:
		m.running = false
		m.clusters = msg.clusters
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace", "b":
			return m, func() tea.Msg { return nav.BackMsg{} }
		}
	}
	return m, nil
}

func (m VerifyModel) View() string {
	if m.running {
		return styles.Dimmed.Render("  ejecutando verificaciones…") + "\n\n" +
			styles.Footer.Render("esc volver")
	}

	var sb strings.Builder
	for _, cl := range m.clusters {
		sb.WriteString(styles.Title.Render(cl.name) + "\n")
		for _, r := range cl.results {
			sb.WriteString(renderResult(r) + "\n")
		}
		sb.WriteString("\n")
	}
	sb.WriteString(styles.Footer.Render("esc volver"))
	return sb.String()
}

func renderResult(r check.Result) string {
	icon := statusIcon(r.Status)
	line := fmt.Sprintf("  %s  %s", icon, r.Name)
	if r.Version != "" {
		line += fmt.Sprintf(" (%s)", r.Version)
	}
	if r.Detail != "" {
		line += "  " + styles.Dimmed.Render(r.Detail)
	}
	return line
}

func statusIcon(s check.Status) string {
	switch s {
	case check.StatusOK:
		return styles.Success.Render("✓")
	case check.StatusMissing:
		return styles.Error.Render("✗")
	case check.StatusOutdated:
		return styles.Warn.Render("!")
	default:
		return styles.Dimmed.Render("?")
	}
}

// runChecks is the tea.Cmd that calls each check cluster synchronously inside a
// goroutine and returns the results as a verifyDoneMsg. It mirrors the same
// cluster order used by the verify CLI command.
func runChecks() tea.Msg {
	sddDir := defaultSDDDir()

	// Gate each cluster by the active domains, mirroring the verify CLI command:
	// a check only shows when some active domain declares its MCP/binary. On
	// resolution error, default to showing it (legacy behavior).
	activeMCP, mcpErr := manifest.ActiveMCPFromEnv()
	mcpActive := func(name string) bool { return mcpErr != nil || slices.Contains(activeMCP, name) }
	activeBins, binErr := manifest.ActiveBinariesFromEnv()
	binActive := func(name string) bool { return binErr != nil || slices.Contains(activeBins, name) }
	activeIDs, _, idErr := manifest.ResolveFromEnv()
	devActive := idErr != nil || slices.Contains(activeIDs, "development")
	activeHooks, hooksErr := manifest.ActiveHooksFromEnv()
	hooksActive := hooksErr == nil && len(activeHooks) > 0

	clusters := []cluster{{"Prerequisites", prereqs.All()}}
	if mcpActive("engram") {
		clusters = append(clusters, cluster{"Engram", engram.All()})
	}
	if mcpActive("codegraph") {
		clusters = append(clusters, cluster{"CodeGraph", codegraph.All()})
	}
	if mcpActive("context7") {
		clusters = append(clusters, cluster{"context7", context7.All()})
	}
	if mcpActive("drawio") {
		clusters = append(clusters, cluster{"drawio", drawio.All()})
	}
	if mcpActive("debugger") {
		clusters = append(clusters, cluster{"debugger", debugger.All()})
	}
	if binActive("proofshot") {
		clusters = append(clusters, cluster{"proofshot", proofshot.All()})
	}
	clusters = append(clusters,
		cluster{"Integración con Claude Code", claudemd.All()},
		cluster{"Auto-aprobación de tools (settings.json)", permissions.All(install.ActiveMCPPatterns())},
	)
	if hooksActive {
		clusters = append(clusters, cluster{"Hooks de dominios activos", hooks.All()})
	}
	if devActive {
		clusters = append(clusters, cluster{"Cross-check SDD ↔ MCP (" + sddDir + ")", sdd.CrossCheck(sddDir)})
	}
	return verifyDoneMsg{clusters: clusters}
}

func defaultSDDDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".claude/agents"
	}
	return filepath.Join(home, ".claude", "agents")
}
