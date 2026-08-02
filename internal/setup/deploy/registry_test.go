package deploy_test

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/franwerner/matecito-ai/internal/setup/deploy"
)

// TestPlan_FourStatuses verifies the plan classifies into all four states:
// new (never on disk), changed (differs from what's on disk), same (matches),
// and removed (registered but the payload no longer produces it).
func TestPlan_FourStatuses(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                        {Data: []byte("# core\n")},
		"domains/development/agents/new.md":     {Data: []byte("new content")},
		"domains/development/agents/changed.md": {Data: []byte("new bytes")},
		"domains/development/agents/same.md":    {Data: []byte("same bytes")},
	}
	claudeHome := newClaudeHome(t)
	stateDir := newStateDir(t)

	// Pre-seed disk state: "changed.md" exists with different content,
	// "same.md" exists with identical content. "new.md" is absent.
	mustWrite(t, filepath.Join(claudeHome, "agents", "changed.md"), []byte("old bytes"))
	mustWrite(t, filepath.Join(claudeHome, "agents", "same.md"), []byte("same bytes"))
	// "gone.md" must still be present on disk for Plan to flag it as Removed —
	// a registered destination absent from disk is filtered out of the plan
	// entirely (see TestPlan_RegistryEntry_AbsentFromDisk_NotInPlan).
	mustWrite(t, filepath.Join(claudeHome, "agents", "gone.md"), []byte("still on disk"))

	// Seed a registry entry for a destination this payload no longer produces.
	reg := deploy.Registry{Version: 1, Entries: []deploy.Entry{
		{Path: "agents/gone.md", Hash: "irrelevant-for-plan", Origin: "domain:development"},
	}}
	if err := deploy.SaveRegistry(stateDir, reg); err != nil {
		t.Fatalf("SaveRegistry: %v", err)
	}

	ops, err := deploy.Plan(payload, claudeHome, stateDir, nil)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}

	statusOf := func(suffix string) (deploy.FileStatus, bool) {
		for _, op := range ops {
			if filepath.Base(op.Target) == suffix {
				return op.Status, true
			}
		}
		return 0, false
	}

	if s, ok := statusOf("new.md"); !ok || s != deploy.StatusNew {
		t.Errorf("new.md: got status %v (found=%v), want StatusNew", s, ok)
	}
	if s, ok := statusOf("changed.md"); !ok || s != deploy.StatusChanged {
		t.Errorf("changed.md: got status %v (found=%v), want StatusChanged", s, ok)
	}
	if s, ok := statusOf("same.md"); !ok || s != deploy.StatusSame {
		t.Errorf("same.md: got status %v (found=%v), want StatusSame", s, ok)
	}
	if s, ok := statusOf("gone.md"); !ok || s != deploy.StatusRemoved {
		t.Errorf("gone.md: got status %v (found=%v), want StatusRemoved", s, ok)
	}
}

// TestPlan_RemovedCoversRenameAndDeactivatedDomain verifies "qué queda a
// borrar": a plain stale registry entry AND a registry entry belonging to a
// domain that is no longer active both resolve to StatusRemoved with no
// special-casing — a rename or a deactivated domain is just "registered, not
// produced today", same as any other removal.
func TestPlan_RemovedCoversRenameAndDeactivatedDomain(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                    {Data: []byte("# core\n")},
		"domains/development/agents/new.md": {Data: []byte("content")},
		"domains/design/agents/designer.md": {Data: []byte("design content")},
	}
	claudeHome := newClaudeHome(t)
	stateDir := newStateDir(t)

	// Both registered destinations must still be present on disk for Plan to
	// flag them as Removed — see TestPlan_RegistryEntry_AbsentFromDisk_NotInPlan.
	mustWrite(t, filepath.Join(claudeHome, "agents", "renamed-away.md"), []byte("stale content"))
	mustWrite(t, filepath.Join(claudeHome, "agents", "designer.md"), []byte("design content"))

	reg := deploy.Registry{Version: 1, Entries: []deploy.Entry{
		{Path: "agents/renamed-away.md", Hash: "h1", Origin: "domain:development"},
		{Path: "agents/designer.md", Hash: "h2", Origin: "domain:design"},
	}}
	if err := deploy.SaveRegistry(stateDir, reg); err != nil {
		t.Fatalf("SaveRegistry: %v", err)
	}

	// design is NOT in the active set — its registered entry must show up as
	// Removed, exactly like the plain stale rename, with no special treatment.
	ops, err := deploy.Plan(payload, claudeHome, stateDir, []string{"development"})
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}

	var foundRename, foundDeactivated bool
	for _, op := range ops {
		if op.Status != deploy.StatusRemoved {
			continue
		}
		switch filepath.Base(op.Target) {
		case "renamed-away.md":
			foundRename = true
		case "designer.md":
			foundDeactivated = true
		}
	}
	if !foundRename {
		t.Error("expected the stale rename entry to be StatusRemoved")
	}
	if !foundDeactivated {
		t.Error("expected the deactivated-domain entry to be StatusRemoved")
	}
}

// TestPlan_RegistryEntry_AbsentFromDisk_NotInPlan verifies the preview-
// inflation bug fix on the registry branch: a registered destination the
// person already deleted by hand must not surface as a StatusRemoved op at
// all. Apply already treats an absent target as a silent no-op; the point of
// this filter is that Plan's preview must not list (or count) a deletion
// that will never happen — otherwise the person confirming a migration with
// a stale, mostly-empty registry sees a wildly inflated "a borrar" count.
func TestPlan_RegistryEntry_AbsentFromDisk_NotInPlan(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                    {Data: []byte("# core\n")},
		"domains/development/agents/new.md": {Data: []byte("content")},
	}
	claudeHome := newClaudeHome(t)
	stateDir := newStateDir(t)

	// "gone.md" is registered but was never written to claudeHome in this
	// test — simulating the person deleting it by hand after it was deployed.
	reg := deploy.Registry{Version: 1, Entries: []deploy.Entry{
		{Path: "agents/gone.md", Hash: "irrelevant", Origin: "domain:development"},
	}}
	if err := deploy.SaveRegistry(stateDir, reg); err != nil {
		t.Fatalf("SaveRegistry: %v", err)
	}

	ops, err := deploy.Plan(payload, claudeHome, stateDir, nil)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}

	for _, op := range ops {
		if op.Status == deploy.StatusRemoved {
			t.Errorf("expected no StatusRemoved op for a registered destination absent from disk, got %+v", op)
		}
	}
}

// TestPlan_ReactivatedDomain_NotRemoved covers the open question from the
// design ("chequear que un Removed que un dominio inactivo repondría al
// reactivarse se comporta como se espera"): a registered destination that a
// domain reactivated this run reproduces is not StatusRemoved — it is just
// New/Changed/Same again, same as any other set-difference outcome.
func TestPlan_ReactivatedDomain_NotRemoved(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                    {Data: []byte("# core\n")},
		"domains/design/agents/designer.md": {Data: []byte("design content")},
	}
	claudeHome := newClaudeHome(t)
	stateDir := newStateDir(t)

	reg := deploy.Registry{Version: 1, Entries: []deploy.Entry{
		{Path: "agents/designer.md", Hash: "h1", Origin: "domain:design"},
	}}
	if err := deploy.SaveRegistry(stateDir, reg); err != nil {
		t.Fatalf("SaveRegistry: %v", err)
	}

	// design is active again this run — its own component reproduces the
	// exact registered target, so the diff must not flag it as Removed.
	ops, err := deploy.Plan(payload, claudeHome, stateDir, []string{"design"})
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}

	for _, op := range ops {
		if filepath.Base(op.Target) == "designer.md" && op.Status == deploy.StatusRemoved {
			t.Errorf("expected the reactivated domain's target NOT to be StatusRemoved, got op: %+v", op)
		}
	}
}

// TestApply_RegistryReflectsApplied verifies "el registro refleja lo
// aplicado": after a run that writes files and deletes a registered one (hash
// matches), the rewritten registry has one entry per deployed file and none
// for the deleted one.
func TestApply_RegistryReflectsApplied(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                  {Data: []byte("# core\n")},
		"domains/development/agents/a.md": {Data: []byte("a content")},
	}
	claudeHome := newClaudeHome(t)
	stateDir := newStateDir(t)
	backupDir := newBackupDir(t)

	goneTarget := filepath.Join(claudeHome, "agents", "gone.md")
	mustWrite(t, goneTarget, []byte("gone content"))
	goneHash := sha256Hex(t, []byte("gone content"))

	reg := deploy.Registry{Version: 1, Entries: []deploy.Entry{
		{Path: "agents/gone.md", Hash: goneHash, Origin: "domain:development"},
	}}
	if err := deploy.SaveRegistry(stateDir, reg); err != nil {
		t.Fatalf("SaveRegistry: %v", err)
	}

	ops, err := deploy.Plan(payload, claudeHome, stateDir, nil)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	result, err := deploy.Apply(payload, ops, claudeHome, stateDir, backupDir)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(result.Preserved) != 0 {
		t.Errorf("expected nothing preserved, got %v", result.Preserved)
	}
	if _, err := os.Stat(goneTarget); !os.IsNotExist(err) {
		t.Errorf("expected gone.md removed from disk, stat err=%v", err)
	}

	final, found, err := deploy.LoadRegistry(stateDir)
	if err != nil || !found {
		t.Fatalf("LoadRegistry after Apply: found=%v err=%v", found, err)
	}
	var hasGone, hasAgent bool
	for _, e := range final.Entries {
		if e.Path == "agents/gone.md" {
			hasGone = true
		}
		if e.Path == "agents/a.md" {
			hasAgent = true
		}
	}
	if hasGone {
		t.Error("expected no registry entry for the deleted file")
	}
	if !hasAgent {
		t.Error("expected a registry entry for the deployed file")
	}
}

// TestApply_RegistryEditOverwritten verifies the new per-source deletion
// criterion at the Apply() level: a registry-origin StatusRemoved entry whose
// disk content the person edited after we wrote it is still backed up and
// deleted — registry membership alone is the proof of ownership — and Apply
// reports it as an overwritten edit, with the backup location, instead of
// preserving it.
func TestApply_RegistryEditOverwritten(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                  {Data: []byte("# core\n")},
		"domains/development/agents/a.md": {Data: []byte("a content")},
	}
	claudeHome := newClaudeHome(t)
	stateDir := newStateDir(t)
	backupDir := newBackupDir(t)

	editedTarget := filepath.Join(claudeHome, "agents", "edited.md")
	mustWrite(t, editedTarget, []byte("edited by the person"))
	originalHash := sha256Hex(t, []byte("original deployed content"))

	reg := deploy.Registry{Version: 1, Entries: []deploy.Entry{
		{Path: "agents/edited.md", Hash: originalHash, Origin: "domain:development"},
	}}
	if err := deploy.SaveRegistry(stateDir, reg); err != nil {
		t.Fatalf("SaveRegistry: %v", err)
	}

	ops, err := deploy.Plan(payload, claudeHome, stateDir, nil)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	result, err := deploy.Apply(payload, ops, claudeHome, stateDir, backupDir)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if len(result.Preserved) != 0 {
		t.Errorf("a registry-origin entry must never be preserved, got %v", result.Preserved)
	}
	if _, err := os.Stat(editedTarget); !os.IsNotExist(err) {
		t.Errorf("expected edited.md removed from disk despite the mismatch, stat err=%v", err)
	}
	if len(result.OverwrittenEdits) != 1 {
		t.Fatalf("expected exactly one OverwrittenEdits entry, got %v", result.OverwrittenEdits)
	}
	edit := result.OverwrittenEdits[0]
	if edit.Target != editedTarget {
		t.Errorf("edit.Target = %q, want %q", edit.Target, editedTarget)
	}
	got, err := os.ReadFile(edit.BackupPath)
	if err != nil {
		t.Fatalf("expected backup at edit.BackupPath=%s: %v", edit.BackupPath, err)
	}
	if string(got) != "edited by the person" {
		t.Errorf("backup content mismatch: got %q, want the person's edited content", got)
	}
}

// TestApply_DeactivatedDomain_PrunesNestedDirsButKeepsComponentRoot verifies
// the Requirement "Un borrado no deja directorios vacíos", scenario
// "desactivar un dominio no deja carpetas fantasma": a deactivated domain's
// own, some of them nested, skill folders disappear entirely once their last
// file is gone, while skills/ itself — the component root, shared with the
// person's own files and other plugins' — keeps existing.
func TestApply_DeactivatedDomain_PrunesNestedDirsButKeepsComponentRoot(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                  {Data: []byte("# core\n")},
		"domains/development/agents/a.md": {Data: []byte("a content")},
	}
	claudeHome := newClaudeHome(t)
	stateDir := newStateDir(t)
	backupDir := newBackupDir(t)

	skillFile := []byte("skill body")
	skillHash := sha256Hex(t, skillFile)
	nestedFile := []byte("nested notes")
	nestedHash := sha256Hex(t, nestedFile)
	otherSkillFile := []byte("other skill body")
	otherSkillHash := sha256Hex(t, otherSkillFile)

	mustWrite(t, filepath.Join(claudeHome, "skills", "design-a", "SKILL.md"), skillFile)
	mustWrite(t, filepath.Join(claudeHome, "skills", "design-a", "refs", "notes.md"), nestedFile)
	mustWrite(t, filepath.Join(claudeHome, "skills", "design-b", "SKILL.md"), otherSkillFile)

	reg := deploy.Registry{Version: 1, Entries: []deploy.Entry{
		{Path: "skills/design-a/SKILL.md", Hash: skillHash, Origin: "domain:design"},
		{Path: "skills/design-a/refs/notes.md", Hash: nestedHash, Origin: "domain:design"},
		{Path: "skills/design-b/SKILL.md", Hash: otherSkillHash, Origin: "domain:design"},
	}}
	if err := deploy.SaveRegistry(stateDir, reg); err != nil {
		t.Fatalf("SaveRegistry: %v", err)
	}

	// design is not in the active set: none of its registered skill files are
	// reproduced by this run's payload, so all three resolve to StatusRemoved.
	ops, err := deploy.Plan(payload, claudeHome, stateDir, []string{"development"})
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if _, err := deploy.Apply(payload, ops, claudeHome, stateDir, backupDir); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	for _, gone := range []string{
		filepath.Join(claudeHome, "skills", "design-a", "refs"),
		filepath.Join(claudeHome, "skills", "design-a"),
		filepath.Join(claudeHome, "skills", "design-b"),
	} {
		if _, err := os.Stat(gone); !os.IsNotExist(err) {
			t.Errorf("expected %s removed, stat err=%v", gone, err)
		}
	}

	if _, err := os.Stat(filepath.Join(claudeHome, "skills")); err != nil {
		t.Errorf("expected skills/ (the component root) to survive: %v", err)
	}
}

// TestApply_DirectoryWithForeignFile_SurvivesPrune verifies the Requirement
// "Un borrado no deja directorios vacíos", scenario "un directorio con algo
// ajeno queda intacto": deleting our file out of a directory that also holds
// something that isn't ours must leave that directory (and the foreign file)
// standing.
func TestApply_DirectoryWithForeignFile_SurvivesPrune(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                  {Data: []byte("# core\n")},
		"domains/development/agents/a.md": {Data: []byte("a content")},
	}
	claudeHome := newClaudeHome(t)
	stateDir := newStateDir(t)
	backupDir := newBackupDir(t)

	ourFile := []byte("our content")
	ourHash := sha256Hex(t, ourFile)
	mustWrite(t, filepath.Join(claudeHome, "skills", "mixed", "ours.md"), ourFile)
	// A file in the same directory that the registry never tracked — not ours.
	foreignTarget := filepath.Join(claudeHome, "skills", "mixed", "theirs.md")
	mustWrite(t, foreignTarget, []byte("foreign content"))

	reg := deploy.Registry{Version: 1, Entries: []deploy.Entry{
		{Path: "skills/mixed/ours.md", Hash: ourHash, Origin: "domain:development"},
	}}
	if err := deploy.SaveRegistry(stateDir, reg); err != nil {
		t.Fatalf("SaveRegistry: %v", err)
	}

	ops, err := deploy.Plan(payload, claudeHome, stateDir, nil)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if _, err := deploy.Apply(payload, ops, claudeHome, stateDir, backupDir); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if _, err := os.Stat(filepath.Join(claudeHome, "skills", "mixed", "ours.md")); !os.IsNotExist(err) {
		t.Errorf("expected ours.md removed, stat err=%v", err)
	}
	got, err := os.ReadFile(foreignTarget)
	if err != nil {
		t.Fatalf("expected the foreign file to survive: %v", err)
	}
	if string(got) != "foreign content" {
		t.Errorf("foreign file content changed: got %q", got)
	}
	if _, err := os.Stat(filepath.Join(claudeHome, "skills", "mixed")); err != nil {
		t.Errorf("expected the directory holding the foreign file to survive: %v", err)
	}
}

// TestApply_ComponentRootSurvivesEmpty verifies the Requirement "Un borrado no
// deja directorios vacíos", scenario "la raíz del componente sobrevive
// vacía" — explicitly and non-incidentally: the ONLY registered file lives
// directly inside the component root (no sibling, no nested directory left
// standing to save it by accident), so once it is deleted the root itself is
// genuinely, verifiably empty — and it must still exist.
func TestApply_ComponentRootSurvivesEmpty(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                  {Data: []byte("# core\n")},
		"domains/development/agents/a.md": {Data: []byte("a content")},
	}
	claudeHome := newClaudeHome(t)
	stateDir := newStateDir(t)
	backupDir := newBackupDir(t)

	content := []byte("lonely content")
	mustWrite(t, filepath.Join(claudeHome, "skills", "lonely.md"), content)

	reg := deploy.Registry{Version: 1, Entries: []deploy.Entry{
		{Path: "skills/lonely.md", Hash: sha256Hex(t, content), Origin: "domain:development"},
	}}
	if err := deploy.SaveRegistry(stateDir, reg); err != nil {
		t.Fatalf("SaveRegistry: %v", err)
	}

	ops, err := deploy.Plan(payload, claudeHome, stateDir, nil)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if _, err := deploy.Apply(payload, ops, claudeHome, stateDir, backupDir); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	root := filepath.Join(claudeHome, "skills")
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("expected the component root to survive and be readable: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected the root to be genuinely empty, found %v", entries)
	}
}

// TestApply_Migration_NoRegistry_DoesNotDeleteByAbsence verifies "Migración en
// la primera corrida": with no registry, Apply writes the fresh registry and
// does not treat "nothing registered" as "everything gone" — no deletions
// beyond what the embedded legacy catalog itself would match by hash (which,
// against a fresh claudeHome with none of those legacy paths present, matches
// nothing).
func TestApply_Migration_NoRegistry_DoesNotDeleteByAbsence(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                  {Data: []byte("# core\n")},
		"domains/development/agents/a.md": {Data: []byte("a content")},
	}
	claudeHome := newClaudeHome(t)
	stateDir := newStateDir(t) // empty: no registry exists yet
	backupDir := newBackupDir(t)

	ops, err := deploy.Plan(payload, claudeHome, stateDir, nil)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	result, err := deploy.Apply(payload, ops, claudeHome, stateDir, backupDir)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if _, err := os.Stat(filepath.Join(claudeHome, "agents", "a.md")); err != nil {
		t.Errorf("expected a.md written: %v", err)
	}

	final, found, err := deploy.LoadRegistry(stateDir)
	if err != nil || !found {
		t.Fatalf("expected a registry to exist after the migration run: found=%v err=%v", found, err)
	}
	var hasAgent bool
	for _, e := range final.Entries {
		if e.Path == "agents/a.md" {
			hasAgent = true
		}
	}
	if !hasAgent {
		t.Error("expected the migration run's registry to track the deployed file")
	}
	_ = result // BackupCreated/Preserved are incidental here; covered elsewhere.
}

func mustWrite(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// sha256Hex hashes with the same primitive deploy uses internally
// (crypto/sha256), so a seeded registry hash can be computed before any file
// exists on disk.
func sha256Hex(t *testing.T, data []byte) string {
	t.Helper()
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
