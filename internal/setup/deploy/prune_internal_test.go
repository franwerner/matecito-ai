package deploy

import (
	"os"
	"path/filepath"
	"testing"
)

// TestPruneEmptyDirsAfterDelete_ClimbsNestedDirsButStopsAtRoot verifies the
// climb goes through however many nested levels a skill folder had, and
// stops exactly at (never past) the component root — the root itself must
// still exist afterward, per "la raíz del componente sobrevive vacía".
func TestPruneEmptyDirsAfterDelete_ClimbsNestedDirsButStopsAtRoot(t *testing.T) {
	claudeHome := t.TempDir()
	nested := filepath.Join(claudeHome, "skills", "design-foo", "refs", "notes.md")
	mustWriteFile(t, nested, []byte("content"))
	if err := os.Remove(nested); err != nil {
		t.Fatalf("os.Remove(nested): %v", err)
	}

	pruneEmptyDirsAfterDelete(nested, claudeHome)

	for _, gone := range []string{
		filepath.Join(claudeHome, "skills", "design-foo", "refs"),
		filepath.Join(claudeHome, "skills", "design-foo"),
	} {
		if _, err := os.Stat(gone); !os.IsNotExist(err) {
			t.Errorf("expected %s removed, stat err=%v", gone, err)
		}
	}

	root := filepath.Join(claudeHome, "skills")
	info, err := os.Stat(root)
	if err != nil {
		t.Fatalf("expected the component root %s to survive: %v", root, err)
	}
	if !info.IsDir() {
		t.Errorf("expected %s to still be a directory", root)
	}
}

// TestPruneEmptyDirsAfterDelete_StopsAtDirectoryWithForeignFile verifies "un
// directorio con algo ajeno queda intacto": a directory that still holds
// something after our file is gone — ours or not — is left standing, and the
// climb does not proceed past it either.
func TestPruneEmptyDirsAfterDelete_StopsAtDirectoryWithForeignFile(t *testing.T) {
	claudeHome := t.TempDir()
	ours := filepath.Join(claudeHome, "skills", "mixed", "ours.md")
	theirs := filepath.Join(claudeHome, "skills", "mixed", "theirs.md")
	mustWriteFile(t, ours, []byte("ours"))
	mustWriteFile(t, theirs, []byte("theirs"))
	if err := os.Remove(ours); err != nil {
		t.Fatalf("os.Remove(ours): %v", err)
	}

	pruneEmptyDirsAfterDelete(ours, claudeHome)

	got, err := os.ReadFile(theirs)
	if err != nil {
		t.Fatalf("expected the foreign file to survive untouched: %v", err)
	}
	if string(got) != "theirs" {
		t.Errorf("foreign file content changed: got %q", got)
	}
	if _, err := os.Stat(filepath.Join(claudeHome, "skills", "mixed")); err != nil {
		t.Errorf("expected the directory holding a foreign file to survive: %v", err)
	}
}

// TestPruneEmptyDirsAfterDelete_ComponentRootNeverRemoved_EvenGenuinelyEmpty
// is the explicit, non-incidental version of "la raíz del componente
// sobrevive vacía": nothing else exists anywhere under the root once the one
// file directly inside it is gone (no sibling, no nested dir left standing
// to "accidentally" save the root) — the root must still be a real,
// statable, empty directory, never removed by the climb.
func TestPruneEmptyDirsAfterDelete_ComponentRootNeverRemoved_EvenGenuinelyEmpty(t *testing.T) {
	claudeHome := t.TempDir()
	root := filepath.Join(claudeHome, "skills")
	lonely := filepath.Join(root, "lonely.md")
	mustWriteFile(t, lonely, []byte("content"))
	if err := os.Remove(lonely); err != nil {
		t.Fatalf("os.Remove(lonely): %v", err)
	}

	pruneEmptyDirsAfterDelete(lonely, claudeHome)

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("expected the component root to still exist and be readable: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected the root to be genuinely empty, found %v", entries)
	}
}

// TestPruneEmptyDirsAfterDelete_OutsideComponentRoots_NeverTouchesDisk covers
// the legacy root:"state" case: a target that isn't under any of
// componentRoots (the abandoned deployed-manifest.json, loose under
// stateDir, not claudeHome at all) must not trigger any cleanup — not even an
// attempt to remove its own parent directory.
func TestPruneEmptyDirsAfterDelete_OutsideComponentRoots_NeverTouchesDisk(t *testing.T) {
	claudeHome := t.TempDir()
	stateDir := t.TempDir()
	target := filepath.Join(stateDir, "deployed-manifest.json")
	mustWriteFile(t, target, []byte("content"))
	if err := os.Remove(target); err != nil {
		t.Fatalf("os.Remove(target): %v", err)
	}

	pruneEmptyDirsAfterDelete(target, claudeHome)

	if _, err := os.Stat(stateDir); err != nil {
		t.Errorf("expected stateDir itself untouched: %v", err)
	}
}

// TestComponentRootFor covers componentRootFor directly: every one of the
// five recognized roots resolves, a nested path resolves to its top root, and
// anything outside claudeHome (or matching no known root name) is rejected.
func TestComponentRootFor(t *testing.T) {
	claudeHome := t.TempDir()

	for _, name := range componentRoots {
		target := filepath.Join(claudeHome, name, "sub", "file.md")
		root, ok := componentRootFor(target, claudeHome)
		if !ok {
			t.Errorf("%s: expected a component root match", name)
			continue
		}
		if want := filepath.Join(claudeHome, name); root != want {
			t.Errorf("%s: root = %q, want %q", name, root, want)
		}
	}

	if _, ok := componentRootFor(filepath.Join(claudeHome, "matecito-ai.md"), claudeHome); ok {
		t.Error("expected the top-level composed matecito-ai.md file to match no component root")
	}
	if _, ok := componentRootFor(filepath.Join(t.TempDir(), "deployed-manifest.json"), claudeHome); ok {
		t.Error("expected a path outside claudeHome to match no component root")
	}
}
