package sync

import (
	"testing"

	pkgsync "github.com/franwerner/matecito-ai/internal/setup/sync"
	"github.com/franwerner/matecito-ai/internal/tui/nav"
)

// TestUpdate_DoneMsg_SelfReplace covers the guarantee behind Requirement "Re-exec
// is CLI-only; TUI never self-execs" (spec #951): on a doneMsg carrying
// SelfReplaced, the TUI must stay on the restart-message screen and must NOT
// emit any command that could trigger a re-exec — the only safe outcome is a
// nil tea.Cmd. Without SelfReplaced it must behave as before: return to the
// menu via nav.BackMsg.
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
			m := SyncModel{}

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
