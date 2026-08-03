package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	pkgsync "github.com/franwerner/matecito-ai/internal/setup/sync"
)

// RunOpts carries the parameters the caller must supply to start the TUI.
type RunOpts struct {
	Version string
}

// Run builds and launches the full-screen bubbletea TUI.
// It first consumes any pending-sync mark left by a previous run's
// self-replace (see runDeferredSync) — synchronously, before tea.NewProgram,
// so its stdout aviso is visible and never fights the alt screen. It then
// detects the project context and resolves the global config path before
// handing off to tea.NewProgram; it returns any program error directly.
func Run(opts RunOpts) error {
	if statePath, err := pkgsync.SyncStatePath(); err == nil {
		runDeferredSync(pkgsync.Sync, statePath, opts.Version)
	}

	ctx := DetectProject()

	globalConfigPath, err := agentmodel.ConfigPath()
	if err != nil {
		globalConfigPath = ""
	}

	appModel := NewAppModel(opts.Version, globalConfigPath, ctx)

	p := tea.NewProgram(appModel, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

// runDeferredSync consumes the pending-sync mark a previous TUI run may have
// left after replacing the matecito-ai binary (screens/sync's doneMsg
// handling sets it on every self-replace). With the mark set, it runs sync
// — with the executable already updated — before the TUI opens, and clears
// the mark only on success. A failure never blocks the TUI: it's reported
// on stdout (not stderr — the spec's explicit requirement) and the mark
// stays put, so the next startup retries; that's safe because Sync is
// idempotent (design's Hallazgo B).
//
// Resume: true excludes matecito-ai itself from the plan (see sync.Sync) and
// DeferPayloadOnSelfReplace: false — this run already has the new binary's
// own embedded payload, so deploy and config ecosistema run for real instead
// of deferring again (design decision "Opciones del diferido").
//
// sync and statePath are injected so this stays a pure function tui_test.go
// can exercise without spawning a real bubbletea program or a real network
// self-install — internal/setup/sync is already an internal/tui dependency
// (see app.go), so this adds no new package dependency.
func runDeferredSync(sync func(pkgsync.Options) pkgsync.Result, statePath, version string) {
	if !pkgsync.LoadPendingSync(statePath) {
		return
	}

	result := sync(pkgsync.Options{
		SelfVersion:               version,
		Resume:                    true,
		DeferPayloadOnSelfReplace: false,
	})
	if result.HasErrors() {
		fmt.Fprintln(os.Stdout, "matecito-ai: la sincronización diferida falló, se reintentará en el próximo arranque.")
		return
	}
	_ = pkgsync.SetPendingSync(statePath, false)
}
