package tui

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	pkgsync "github.com/franwerner/matecito-ai/internal/setup/sync"
)

// captureStdout redirects os.Stdout for the duration of fn and returns
// whatever was written to it.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	var sb strings.Builder
	buf := make([]byte, 4096)
	for {
		n, readErr := r.Read(buf)
		sb.Write(buf[:n])
		if readErr != nil {
			break
		}
	}
	return sb.String()
}

// TestRunDeferredSync_NoPendingMark_DoesNotCallSync covers the Requirement
// scenario "sin marca no corre nada": with no pending-sync mark, the
// injected sync function must never be invoked.
func TestRunDeferredSync_NoPendingMark_DoesNotCallSync(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "sync-state.json")

	called := false
	sync := func(pkgsync.Options) pkgsync.Result {
		called = true
		return pkgsync.Result{}
	}

	runDeferredSync(sync, statePath, "1.0.0")

	if called {
		t.Fatal("expected sync to never be called without a pending mark")
	}
}

// TestRunDeferredSync_PendingMark_RunsResumeSyncAndClearsOnSuccess covers
// "la marca dispara el sync antes de abrir la interfaz" and its cleanup:
// with the mark set, sync must run with Resume=true and
// DeferPayloadOnSelfReplace=false (design decision "Opciones del
// diferido"), and a successful result clears the mark.
func TestRunDeferredSync_PendingMark_RunsResumeSyncAndClearsOnSuccess(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "sync-state.json")
	if err := pkgsync.SetPendingSync(statePath, true); err != nil {
		t.Fatalf("SetPendingSync: %v", err)
	}

	var gotOpts pkgsync.Options
	sync := func(opts pkgsync.Options) pkgsync.Result {
		gotOpts = opts
		return pkgsync.Result{}
	}

	runDeferredSync(sync, statePath, "1.2.3")

	if !gotOpts.Resume {
		t.Error("expected Resume=true")
	}
	if gotOpts.DeferPayloadOnSelfReplace {
		t.Error("expected DeferPayloadOnSelfReplace=false")
	}
	if gotOpts.SelfVersion != "1.2.3" {
		t.Errorf("SelfVersion = %q, want %q", gotOpts.SelfVersion, "1.2.3")
	}
	if pkgsync.LoadPendingSync(statePath) {
		t.Fatal("expected the mark to be cleared after a successful deferred sync")
	}
}

// TestRunDeferredSync_Failure_LeavesMarkAndWarnsOnStdout covers "el sync
// diferido falla, avisa y se reintenta": a failing result must NOT clear the
// mark, and the failure must be reported on stdout — the spec's explicit
// requirement, not stderr — so the next startup retries.
func TestRunDeferredSync_Failure_LeavesMarkAndWarnsOnStdout(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "sync-state.json")
	if err := pkgsync.SetPendingSync(statePath, true); err != nil {
		t.Fatalf("SetPendingSync: %v", err)
	}

	sync := func(pkgsync.Options) pkgsync.Result {
		return pkgsync.Result{Errors: map[string]error{"deploy": errors.New("boom")}}
	}

	stdout := captureStdout(t, func() {
		runDeferredSync(sync, statePath, "1.0.0")
	})

	if !pkgsync.LoadPendingSync(statePath) {
		t.Fatal("expected the mark to stay set after a failed deferred sync")
	}
	if stdout == "" {
		t.Fatal("expected a failure notice on stdout")
	}
}
