package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/setup/deploy"
	pkgsync "github.com/franwerner/matecito-ai/internal/setup/sync"
	"github.com/franwerner/matecito-ai/internal/tui/nav"
)

// TestUpdate_DoneMsg_SelfReplace covers the guarantee behind Requirement "Re-exec
// is CLI-only; TUI never self-execs" (spec #951): on a doneMsg carrying
// SelfReplaced, the TUI must stay on the restart-message screen and must NOT
// emit any command that could trigger a re-exec — the only safe outcome is a
// nil tea.Cmd. Without SelfReplaced it must behave as before: return to the
// menu via nav.BackMsg. It also covers the "marca / sin marca" half of the
// pending-sync Requirement: SelfReplaced=true must persist the mark;
// SelfReplaced=false must leave the state file untouched (ilegible is its
// own test below, injecting a corrupt statePath).
func TestUpdate_DoneMsg_SelfReplace(t *testing.T) {
	cases := []struct {
		name             string
		result           pkgsync.Result
		wantSelfReplaced bool
		wantNilCmd       bool
	}{
		{
			name:             "self-replaced stays on restart screen with no cmd",
			result:           pkgsync.Result{SelfReplaced: true},
			wantSelfReplaced: true,
			wantNilCmd:       true,
		},
		{
			name:             "no self-replace returns to the menu",
			result:           pkgsync.Result{SelfReplaced: false},
			wantSelfReplaced: false,
			wantNilCmd:       false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			statePath := filepath.Join(t.TempDir(), "sync-state.json")
			m := SyncModel{statePath: statePath}

			gotModel, cmd := m.Update(doneMsg{result: tc.result})

			sm, ok := gotModel.(SyncModel)
			if !ok {
				t.Fatalf("Update returned %T, want SyncModel", gotModel)
			}
			if sm.selfReplaced != tc.wantSelfReplaced {
				t.Fatalf("selfReplaced = %v, want %v", sm.selfReplaced, tc.wantSelfReplaced)
			}
			if !sm.done {
				t.Fatal("expected done to be true after doneMsg")
			}
			if pending := pkgsync.LoadPendingSync(statePath); pending != tc.wantSelfReplaced {
				t.Fatalf("pending-sync mark = %v, want %v", pending, tc.wantSelfReplaced)
			}

			if tc.wantNilCmd {
				if cmd != nil {
					t.Fatalf("expected a nil tea.Cmd (no re-exec) when SelfReplaced=%v, got non-nil", tc.result.SelfReplaced)
				}
				return
			}

			if cmd == nil {
				t.Fatal("expected a non-nil cmd to navigate back when SelfReplaced=false")
			}
			if _, ok := cmd().(nav.BackMsg); !ok {
				t.Fatalf("expected cmd to yield nav.BackMsg, got %T", cmd())
			}
		})
	}
}

// TestUpdate_DoneMsg_SelfReplace_UnreadableState covers the "ilegible" half
// of the pending-sync Requirement reaching the screen: a corrupt state file
// must not be overwritten (no-clobber), yet the restart aviso still shows
// exactly as when the write succeeds — Requirement "estado ilegible o
// corrupto — no se pisa" says the screen shows its aviso regardless.
func TestUpdate_DoneMsg_SelfReplace_UnreadableState(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "sync-state.json")
	corrupt := []byte("not json{{{")
	if err := os.WriteFile(statePath, corrupt, 0o644); err != nil {
		t.Fatal(err)
	}

	m := SyncModel{statePath: statePath}
	gotModel, cmd := m.Update(doneMsg{result: pkgsync.Result{SelfReplaced: true}})

	sm, ok := gotModel.(SyncModel)
	if !ok {
		t.Fatalf("Update returned %T, want SyncModel", gotModel)
	}
	if !sm.selfReplaced {
		t.Fatal("expected the aviso to show even when the mark can't be written")
	}
	if cmd != nil {
		t.Fatal("expected a nil cmd (no re-exec) even when the mark write fails")
	}

	got, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != string(corrupt) {
		t.Fatalf("expected the corrupt state file untouched, got:\n%s", got)
	}
}

// TestView_SelfReplaced_ShowsDeferredSyncMessage covers the MODIFIED
// Requirement "Desenlace del auto-reemplazo en la interfaz interactiva": the
// restart screen must tell the person the sync completes on its own next
// time, not to manually re-run the binary.
func TestView_SelfReplaced_ShowsDeferredSyncMessage(t *testing.T) {
	m := SyncModel{done: true, selfReplaced: true}

	got := m.View()
	if !strings.Contains(got, "matecito-ai actualizado.") {
		t.Fatalf("expected the update aviso, got:\n%s", got)
	}
	if !strings.Contains(got, "próxima vez") {
		t.Fatalf("expected a message about the deferred sync completing on the next open, got:\n%s", got)
	}
	if strings.Contains(got, "Re-ejecutá") {
		t.Fatalf("expected the old manual-restart message to be gone, got:\n%s", got)
	}
}

// TestStartSync_SuppressesEnginePlanAndStreamsCompletion drives startSync()
// for real — the one seam no other test in this package (or in
// internal/setup/sync) reaches: it is the only place that sets
// opts.PlanShown = true right before invoking the real pkgsync.Sync in a
// goroutine and streaming its output line by line. A no-op, unrecognized
// component name (same fixture pattern as internal/setup/sync's own
// TestSync_Resume) keeps this a hermetic unit test — no network, no real
// ~/.claude — while still exercising the real Sync() call and its actual
// "Plan:" gate.
func TestStartSync_SuppressesEnginePlanAndStreamsCompletion(t *testing.T) {
	m := SyncModel{
		opts: pkgsync.Options{
			PreDetected: []pkgsync.ComponentState{{Name: "widget-test", Present: false}},
		},
	}

	cmd := m.startSync()
	msg := cmd()

	var lines []string
	for {
		switch mm := msg.(type) {
		case outputLineMsg:
			lines = append(lines, mm.line)
			msg = waitForLine(mm)()
		case doneMsg:
			if mm.err != nil {
				t.Fatalf("unexpected doneMsg error: %v", mm.err)
			}
			got := strings.Join(lines, "\n")
			if strings.Contains(got, "Plan:") {
				t.Fatalf("expected startSync to suppress the engine's own \"Plan:\" print (PlanShown=true), got:\n%s", got)
			}
			if !strings.Contains(got, "widget-test") || !strings.Contains(got, "✓ OK") {
				t.Fatalf("expected the streamed run to actually execute the action, got:\n%s", got)
			}
			return
		default:
			t.Fatalf("unexpected message type %T", mm)
		}
	}
}

// TestView_AwaitingConfirm_ShowsPayloadBreakdown verifies the confirmation
// screen shows the payload's New/Changed/Removed/Same breakdown and the full
// list of destinos a borrar before the user confirms — Requirement "el plan
// anticipa los borrados" (spec: "Mostrar el plan y esperar confirmación").
func TestView_AwaitingConfirm_ShowsPayloadBreakdown(t *testing.T) {
	m := SyncModel{
		awaitingConfirm: true,
		planActions: []pkgsync.SyncAction{
			{Component: "deploy", Kind: pkgsync.ActionUpdate},
		},
		planStates: []pkgsync.ComponentState{
			{
				Name:    "deploy",
				Present: true,
				PayloadSummary: deploy.Summary{
					New: 1, Changed: 0, Same: 4, Removed: 1,
					Removals: []string{"skills/old-thing/SKILL.md"},
				},
			},
		},
	}

	got := m.View()
	for _, want := range []string{
		"1 nuevos, 0 cambiados, 1 a borrar, 4 iguales",
		"skills/old-thing/SKILL.md",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected View() to contain %q, got:\n%s", want, got)
		}
	}
}
