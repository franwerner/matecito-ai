package sync

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// statePath returns a path to a not-yet-existing sync-state.json inside a
// fresh t.TempDir(), so every test starts from an absent file.
func statePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "sync-state.json")
}

// TestSetPendingSync_CreatesFreshStateWhenAbsent covers the Requirement
// scenario "no hay estado persistido todavía": setting the mark on a path
// with no file yet creates one with the mark set.
func TestSetPendingSync_CreatesFreshStateWhenAbsent(t *testing.T) {
	path := statePath(t)

	if err := SetPendingSync(path, true); err != nil {
		t.Fatalf("SetPendingSync: %v", err)
	}
	if !LoadPendingSync(path) {
		t.Fatal("expected the mark to be set after creating a fresh state")
	}
}

// TestSetPendingSync_PreservesLastCheck verifies read-modify-write: setting
// the mark must not lose an existing lastCheck (design's Hallazgo C —
// SaveSyncState used to clobber the whole file on every write).
func TestSetPendingSync_PreservesLastCheck(t *testing.T) {
	path := statePath(t)
	now := time.Now()
	if err := SaveSyncState(path, now); err != nil {
		t.Fatalf("SaveSyncState: %v", err)
	}

	if err := SetPendingSync(path, true); err != nil {
		t.Fatalf("SetPendingSync: %v", err)
	}

	got, err := LoadSyncState(path)
	if err != nil {
		t.Fatalf("LoadSyncState: %v", err)
	}
	if got.UTC().Format(time.RFC3339) != now.UTC().Format(time.RFC3339) {
		t.Fatalf("LoadSyncState after SetPendingSync = %v, want %v", got, now)
	}
	if !LoadPendingSync(path) {
		t.Fatal("expected the mark to still be set")
	}
}

// TestSaveSyncState_PreservesPendingSync is the inverse of the test above:
// writing lastCheck must not clear an already-set mark.
func TestSaveSyncState_PreservesPendingSync(t *testing.T) {
	path := statePath(t)
	if err := SetPendingSync(path, true); err != nil {
		t.Fatalf("SetPendingSync: %v", err)
	}

	if err := SaveSyncState(path, time.Now()); err != nil {
		t.Fatalf("SaveSyncState: %v", err)
	}

	if !LoadPendingSync(path) {
		t.Fatal("expected the mark to survive a SaveSyncState write (RMW)")
	}
}

// TestLoadPendingSync_AbsentFile covers the "sin marca" baseline: no file at
// all reads as no mark.
func TestLoadPendingSync_AbsentFile(t *testing.T) {
	path := statePath(t)
	if LoadPendingSync(path) {
		t.Fatal("expected no mark on a path with no file")
	}
}

// TestLoadPendingSync_OldFileWithoutField covers a state file written before
// this change existed: no pendingSync key at all must read as false, not
// error.
func TestLoadPendingSync_OldFileWithoutField(t *testing.T) {
	path := statePath(t)
	if err := os.WriteFile(path, []byte(`{"lastCheck":"2024-01-01T00:00:00Z"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if LoadPendingSync(path) {
		t.Fatal("expected an old-format file without pendingSync to read as false")
	}
}

// TestSetPendingSync_UnreadableFile_NoClobber covers the Requirement
// scenario "estado ilegible o corrupto — no se pisa": SetPendingSync must
// refuse to write over a file it can't parse, leaving it byte-for-byte
// untouched — losing the mark is preferable to clobbering the rest of the
// state.
func TestSetPendingSync_UnreadableFile_NoClobber(t *testing.T) {
	path := statePath(t)
	corrupt := []byte("not json{{{")
	if err := os.WriteFile(path, corrupt, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := SetPendingSync(path, true); err == nil {
		t.Fatal("expected SetPendingSync to refuse writing over a corrupt state file")
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != string(corrupt) {
		t.Fatalf("expected the corrupt file untouched, got:\n%s", got)
	}
}

// TestLoadSyncState_TreatsCorruptAsAbsent is a regression for LoadSyncState's
// pre-existing documented contract (corrupt JSON reads as zero time, no
// error) now that it delegates to the shared loadStateFile/mutateStateFile
// pair.
func TestLoadSyncState_TreatsCorruptAsAbsent(t *testing.T) {
	path := statePath(t)
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := LoadSyncState(path)
	if err != nil {
		t.Fatalf("LoadSyncState: %v", err)
	}
	if !got.IsZero() {
		t.Fatalf("expected zero time for corrupt state, got %v", got)
	}
}
