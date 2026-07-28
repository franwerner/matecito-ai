package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/franwerner/matecito-ai/apps/api/internal/store"
)

func TestEnsureMachineID_GeneratesAndReuses(t *testing.T) {
	dir := t.TempDir()

	first, err := store.EnsureMachineID(dir)
	if err != nil {
		t.Fatalf("first EnsureMachineID: %v", err)
	}
	if first == "" {
		t.Fatal("expected a non-empty machine id")
	}

	second, err := store.EnsureMachineID(dir)
	if err != nil {
		t.Fatalf("second EnsureMachineID: %v", err)
	}
	if second != first {
		t.Fatalf("machine id changed across calls: %q != %q", first, second)
	}
}

func TestEnsureMachineID_RejectsEmptyFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "machine-id"), nil, 0o644); err != nil {
		t.Fatalf("seed empty machine-id file: %v", err)
	}

	if _, err := store.EnsureMachineID(dir); err == nil {
		t.Fatal("expected an error for an empty machine-id file, got nil")
	}
}
