package deploy

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLegacyRemovedOpsFrom_RootRouting verifies a "state" root entry resolves
// under stateDir while a "claude" (or default) root resolves under claudeHome —
// the only two roots the embedded catalog uses (claude for the 102 payload
// destinations, state for the abandoned deployed-manifest.json).
func TestLegacyRemovedOpsFrom_RootRouting(t *testing.T) {
	claudeHome := t.TempDir()
	stateDir := t.TempDir()
	catalog := []LegacyEntry{
		{Root: "claude", Path: "skills/old-skill/SKILL.md", Hashes: []string{"abc"}},
		{Root: "state", Path: "deployed-manifest.json", Hashes: []string{"def"}},
	}
	// Both targets must exist on disk: legacyRemovedOpsFrom now skips a
	// catalog entry absent from this host (see TestLegacyRemovedOpsFrom_SkipsAbsentTarget),
	// which this test is not exercising — it's exercising Root routing.
	mustWriteFile(t, filepath.Join(claudeHome, "skills", "old-skill", "SKILL.md"), []byte("irrelevant"))
	mustWriteFile(t, filepath.Join(stateDir, "deployed-manifest.json"), []byte("irrelevant"))

	ops := legacyRemovedOpsFrom(catalog, map[string]bool{}, claudeHome, stateDir)
	if len(ops) != 2 {
		t.Fatalf("expected 2 ops, got %d: %v", len(ops), ops)
	}

	want := map[string]string{
		filepath.Join(claudeHome, "skills", "old-skill", "SKILL.md"): "abc",
		filepath.Join(stateDir, "deployed-manifest.json"):            "def",
	}
	for _, op := range ops {
		hashes, ok := want[op.Target]
		if !ok {
			t.Errorf("unexpected target %q", op.Target)
			continue
		}
		if len(op.LegacyHashes) != 1 || op.LegacyHashes[0] != hashes {
			t.Errorf("target %q: got hashes %v, want [%s]", op.Target, op.LegacyHashes, hashes)
		}
		if op.Status != StatusRemoved {
			t.Errorf("target %q: expected StatusRemoved", op.Target)
		}
	}
}

// TestLegacyRemovedOpsFrom_SkipsCurrentTarget verifies the defensive skip: an
// entry whose target the current plan already reproduces (e.g. a reactivated
// domain) is never turned into a removal.
func TestLegacyRemovedOpsFrom_SkipsCurrentTarget(t *testing.T) {
	claudeHome := t.TempDir()
	stateDir := t.TempDir()
	target := filepath.Join(claudeHome, "skills", "back-again", "SKILL.md")
	catalog := []LegacyEntry{{Root: "claude", Path: "skills/back-again/SKILL.md", Hashes: []string{"abc"}}}

	ops := legacyRemovedOpsFrom(catalog, map[string]bool{target: true}, claudeHome, stateDir)
	if len(ops) != 0 {
		t.Fatalf("expected the reactivated target to be skipped, got %v", ops)
	}
}

// TestLegacyRemovedOpsFrom_SkipsAbsentTarget verifies the preview-inflation
// bug fix: a catalog entry whose target is not actually present on this host
// must not produce a StatusRemoved op at all — the overwhelming common case,
// since the embedded catalog covers every orphan across every past version,
// and a given host only ever accumulated a few of them. Apply already treats
// an absent target as a silent no-op; the point of this filter is that Plan's
// preview must not list (or count) a deletion that will never happen.
func TestLegacyRemovedOpsFrom_SkipsAbsentTarget(t *testing.T) {
	claudeHome := t.TempDir()
	stateDir := t.TempDir()
	catalog := []LegacyEntry{
		{Root: "claude", Path: "skills/never-existed-here/SKILL.md", Hashes: []string{"abc"}},
		{Root: "state", Path: "deployed-manifest.json", Hashes: []string{"def"}},
	}

	ops := legacyRemovedOpsFrom(catalog, map[string]bool{}, claudeHome, stateDir)
	if len(ops) != 0 {
		t.Fatalf("expected no ops for catalog entries absent from disk, got %v", ops)
	}
}

// mustWriteFile writes data at path, creating parent directories as needed.
func mustWriteFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestApplyRemoval_LegacyHashMatch verifies the migration path: a disk file
// whose content hashes to one of the entry's historical hashes is backed up
// and deleted.
func TestApplyRemoval_LegacyHashMatch(t *testing.T) {
	claudeHome := t.TempDir()
	stateDir := t.TempDir()
	backupDir := t.TempDir()

	target := filepath.Join(claudeHome, "skills", "old", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	content := []byte("old skill body")
	if err := os.WriteFile(target, content, 0o644); err != nil {
		t.Fatal(err)
	}
	h := hashBytes(content)

	op := FileOp{Target: target, Status: StatusRemoved, LegacyHashes: []string{"unrelated", h}}
	deleted, preserved, edit, err := applyRemoval(op, claudeHome, stateDir, backupDir)
	if err != nil {
		t.Fatalf("applyRemoval: %v", err)
	}
	if !deleted || preserved {
		t.Fatalf("expected deleted=true preserved=false, got deleted=%v preserved=%v", deleted, preserved)
	}
	if edit != nil {
		t.Errorf("legacy path never reports an OverwrittenEdit, got %+v", edit)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Errorf("expected target removed from disk, stat err=%v", err)
	}
	bk := filepath.Join(backupDir, "skills", "old", "SKILL.md")
	got, err := os.ReadFile(bk)
	if err != nil {
		t.Fatalf("expected backup at %s: %v", bk, err)
	}
	if string(got) != string(content) {
		t.Errorf("backup content mismatch: got %q want %q", got, content)
	}
}

// TestApplyRemoval_LegacyHashMismatch_Preserved verifies the "working tree
// sucio" case: content that matches none of the historical hashes (Hashes may
// even be empty — "vacío = nunca se borra") is preserved and reported, and the
// run does not fail.
func TestApplyRemoval_LegacyHashMismatch_Preserved(t *testing.T) {
	claudeHome := t.TempDir()
	stateDir := t.TempDir()
	backupDir := t.TempDir()

	target := filepath.Join(claudeHome, "skills", "dirty", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("dirty-tree content"), 0o644); err != nil {
		t.Fatal(err)
	}

	op := FileOp{Target: target, Status: StatusRemoved, LegacyHashes: nil}
	deleted, preserved, edit, err := applyRemoval(op, claudeHome, stateDir, backupDir)
	if err != nil {
		t.Fatalf("applyRemoval: %v", err)
	}
	if deleted || !preserved {
		t.Fatalf("expected deleted=false preserved=true, got deleted=%v preserved=%v", deleted, preserved)
	}
	if edit != nil {
		t.Errorf("a preserved legacy entry must not report an OverwrittenEdit, got %+v", edit)
	}
	if _, err := os.Stat(target); err != nil {
		t.Errorf("expected target left untouched on disk: %v", err)
	}
}

// TestApplyRemoval_AbsentFromDisk_SilentNoError verifies "uno que ya no está
// en disco sale del registro sin error": neither deleted nor preserved, no
// error — the overwhelming common case for the 103 embedded legacy entries on
// any given host.
func TestApplyRemoval_AbsentFromDisk_SilentNoError(t *testing.T) {
	claudeHome := t.TempDir()
	stateDir := t.TempDir()
	backupDir := t.TempDir()

	target := filepath.Join(claudeHome, "skills", "never-existed-here", "SKILL.md")
	op := FileOp{Target: target, Status: StatusRemoved, LegacyHashes: []string{"whatever"}}
	deleted, preserved, edit, err := applyRemoval(op, claudeHome, stateDir, backupDir)
	if err != nil {
		t.Fatalf("applyRemoval: %v", err)
	}
	if deleted || preserved {
		t.Fatalf("expected deleted=false preserved=false for an absent file, got deleted=%v preserved=%v", deleted, preserved)
	}
	if edit != nil {
		t.Errorf("an absent file must not report an OverwrittenEdit, got %+v", edit)
	}
}

// TestApplyRemoval_RegisteredHash_MatchAndMismatch covers the registry-diff
// path (RegisteredHash set): the registry's ownership always wins now — a
// registered entry is deleted whether or not the disk hash still matches. The
// hash only decides what gets reported: a clean delete when it matches, an
// OverwrittenEdit (with the backup location) when the person edited the file
// after we wrote it.
func TestApplyRemoval_RegisteredHash_MatchAndMismatch(t *testing.T) {
	claudeHome := t.TempDir()
	stateDir := t.TempDir()
	backupDir := t.TempDir()

	t.Run("match deletes cleanly", func(t *testing.T) {
		target := filepath.Join(claudeHome, "agents", "gone.md")
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatal(err)
		}
		content := []byte("registered content")
		if err := os.WriteFile(target, content, 0o644); err != nil {
			t.Fatal(err)
		}
		op := FileOp{Target: target, Status: StatusRemoved, RegisteredHash: hashBytes(content)}
		deleted, preserved, edit, err := applyRemoval(op, claudeHome, stateDir, backupDir)
		if err != nil {
			t.Fatalf("applyRemoval: %v", err)
		}
		if !deleted || preserved {
			t.Fatalf("expected deleted=true preserved=false, got deleted=%v preserved=%v", deleted, preserved)
		}
		if edit != nil {
			t.Errorf("a matching hash must not report an OverwrittenEdit, got %+v", edit)
		}
	})

	t.Run("mismatch deletes anyway and reports the overwritten edit", func(t *testing.T) {
		target := filepath.Join(claudeHome, "agents", "edited.md")
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, []byte("user-edited content"), 0o644); err != nil {
			t.Fatal(err)
		}
		op := FileOp{Target: target, Status: StatusRemoved, RegisteredHash: hashBytes([]byte("original content"))}
		deleted, preserved, edit, err := applyRemoval(op, claudeHome, stateDir, backupDir)
		if err != nil {
			t.Fatalf("applyRemoval: %v", err)
		}
		if !deleted || preserved {
			t.Fatalf("expected deleted=true preserved=false (registry ownership always deletes), got deleted=%v preserved=%v", deleted, preserved)
		}
		if _, err := os.Stat(target); !os.IsNotExist(err) {
			t.Errorf("expected target removed from disk despite the mismatch, stat err=%v", err)
		}
		if edit == nil {
			t.Fatal("expected an OverwrittenEdit report for the hash mismatch")
		}
		if edit.Target != target {
			t.Errorf("edit.Target = %q, want %q", edit.Target, target)
		}
		got, err := os.ReadFile(edit.BackupPath)
		if err != nil {
			t.Fatalf("expected backup at edit.BackupPath=%s: %v", edit.BackupPath, err)
		}
		if string(got) != "user-edited content" {
			t.Errorf("backup content mismatch: got %q, want the person's edited content", got)
		}
	})
}

// TestSummarize_RelativizesRemovals_IncludingStateRootFallback closes the
// coverage gap on Summarize itself: every other test that exercises the
// breakdown display builds a deploy.Summary{} literal directly, never calling
// Summarize() on real []FileOp. This one feeds it real StatusRemoved ops built
// by the actual production path (legacyRemovedOpsFrom), including the one
// root:"state" entry whose Target falls outside claudeHome and must fall back
// to an absolute display path instead of a confusing "..".
func TestSummarize_RelativizesRemovals_IncludingStateRootFallback(t *testing.T) {
	claudeHome := t.TempDir()
	stateDir := t.TempDir()

	catalog := []LegacyEntry{
		{Root: "claude", Path: "skills/old-skill/SKILL.md", Hashes: []string{"abc"}},
		{Root: "state", Path: "deployed-manifest.json", Hashes: []string{"def"}},
	}
	// Both targets must exist on disk for legacyRemovedOpsFrom to emit them —
	// this test is about Summarize's counting/relativization, not about the
	// absent-target filter (see TestLegacyRemovedOpsFrom_SkipsAbsentTarget).
	mustWriteFile(t, filepath.Join(claudeHome, "skills", "old-skill", "SKILL.md"), []byte("irrelevant"))
	mustWriteFile(t, filepath.Join(stateDir, "deployed-manifest.json"), []byte("irrelevant"))
	removed := legacyRemovedOpsFrom(catalog, map[string]bool{}, claudeHome, stateDir)

	ops := append([]FileOp{
		{Target: filepath.Join(claudeHome, "agents", "new.md"), Status: StatusNew},
		{Target: filepath.Join(claudeHome, "agents", "changed.md"), Status: StatusChanged},
		{Target: filepath.Join(claudeHome, "agents", "same.md"), Status: StatusSame},
	}, removed...)

	s := Summarize(ops, claudeHome)

	if s.New != 1 || s.Changed != 1 || s.Same != 1 || s.Removed != 2 {
		t.Fatalf("unexpected category counts: %+v", s)
	}
	if len(s.Removals) != 2 {
		t.Fatalf("expected 2 removal display entries, got %v", s.Removals)
	}

	wantClaudeRelative := "skills/old-skill/SKILL.md"
	wantStateAbsolute := filepath.Join(stateDir, "deployed-manifest.json")
	var gotClaudeRelative, gotStateAbsolute bool
	for _, r := range s.Removals {
		if r == wantClaudeRelative {
			gotClaudeRelative = true
		}
		if r == wantStateAbsolute {
			gotStateAbsolute = true
		}
	}
	if !gotClaudeRelative {
		t.Errorf("expected claudeHome-relative removal %q in %v", wantClaudeRelative, s.Removals)
	}
	if !gotStateAbsolute {
		t.Errorf("expected the root:state entry to fall back to its absolute path %q in %v", wantStateAbsolute, s.Removals)
	}
}
