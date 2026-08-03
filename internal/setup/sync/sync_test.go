package sync

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/setup/deploy"
)

// readTracker is an io.Reader that records whether Read was ever called, so
// tests can assert stdin was never touched (resume must not prompt).
type readTracker struct {
	read bool
}

func (r *readTracker) Read(p []byte) (int, error) {
	r.read = true
	return 0, io.EOF
}

// TestSync_Resume verifies the resume-run contract from design #952 / spec
// "Resumed run skips update and prompt": the self (matecito-ai) action is
// excluded from the active plan — guaranteeing termination and
// SelfReplaced==false — the "Plan:" print and the interactive confirm are
// suppressed, and other actions still execute with visible per-action
// progress output.
//
// The non-self fixture component ("widget-test") is an unrecognized name:
// Sync's execution switch has no case for it, so it runs the exact same
// active/print/execute path as a real component but resolves to a no-op
// (runErr stays nil, "✓ OK" is printed) without touching the filesystem or
// network — keeping this a true unit test.
func TestSync_Resume(t *testing.T) {
	states := []ComponentState{
		{Name: "matecito-ai", Present: false}, // would plan an install if not excluded
		{Name: "widget-test", Present: false}, // plans an install; unrecognized by the switch
	}

	stdin := &readTracker{}
	var out, errOut strings.Builder
	opts := Options{
		Stdin:       stdin,
		Stdout:      &out,
		Stderr:      &errOut,
		PreDetected: states,
		Resume:      true,
		Yes:         false,
		BackupDir:   t.TempDir(),
	}

	result := Sync(opts)

	if result.SelfReplaced {
		t.Fatal("expected SelfReplaced=false: the self action must be excluded on resume (termination guarantee)")
	}

	got := out.String()
	if strings.Contains(got, "Plan:") {
		t.Fatalf("expected no \"Plan:\" print on resume, got:\n%s", got)
	}
	if stdin.read {
		t.Fatal("expected stdin to never be read on resume (Yes is forced true)")
	}
	if strings.Contains(got, "matecito-ai") {
		t.Fatalf("expected no per-action progress line for the excluded self component, got:\n%s", got)
	}
	if !strings.Contains(got, "✓ OK") {
		t.Fatalf("expected per-action progress output (\"✓ OK\") for the non-self action, got:\n%s", got)
	}
}

// TestSync_NoResume_ShowsPlanAndReadsStdin is the control case: without
// Resume, the plan is printed and a "no" answer on stdin cancels — proving
// the resume behavior above is a deliberate suppression, not an accident of
// the fixture.
func TestSync_NoResume_ShowsPlanAndReadsStdin(t *testing.T) {
	states := []ComponentState{
		{Name: "widget-test", Present: false},
	}

	var out, errOut strings.Builder
	opts := Options{
		Stdin:       strings.NewReader("n\n"),
		Stdout:      &out,
		Stderr:      &errOut,
		PreDetected: states,
		Resume:      false,
		Yes:         false,
	}

	result := Sync(opts)

	got := out.String()
	if !strings.Contains(got, "Plan:") {
		t.Fatalf("expected \"Plan:\" print without resume, got:\n%s", got)
	}
	if !strings.Contains(got, "Cancelado.") {
		t.Fatalf("expected the run to be cancelled by the \"n\" stdin answer, got:\n%s", got)
	}
	if result.SelfReplaced {
		t.Fatal("expected SelfReplaced=false: nothing was executed after cancelling")
	}
}

// TestSync_PlanShown_SuppressesPrint verifies the new gate: a caller that
// already showed its own plan (install's CLI preview, the TUI's confirmation
// screen) sets PlanShown, and the engine's own "Plan:" print must not repeat
// it — distinct from Resume, and NOT gated by Yes (see Options.PlanShown doc:
// `update --yes` must keep seeing the engine's print, since it has none of its
// own).
func TestSync_PlanShown_SuppressesPrint(t *testing.T) {
	states := []ComponentState{
		{Name: "widget-test", Present: false},
	}

	var out, errOut strings.Builder
	opts := Options{
		Stdin:       strings.NewReader(""),
		Stdout:      &out,
		Stderr:      &errOut,
		PreDetected: states,
		Resume:      false,
		PlanShown:   true,
		Yes:         true,
	}

	Sync(opts)

	got := out.String()
	if strings.Contains(got, "Plan:") {
		t.Fatalf("expected no \"Plan:\" print when PlanShown=true, got:\n%s", got)
	}
	if !strings.Contains(got, "✓ OK") {
		t.Fatalf("expected the action to still execute despite PlanShown, got:\n%s", got)
	}
}

// TestSync_DryRun_ShowsPayloadBreakdown verifies the engine's own "Plan:"
// preview includes the payload New/Changed/Removed/Same breakdown and the
// full list of destinos a borrar — DryRun is used so the run returns before
// ever reaching the real "deploy" execution branch (this must stay a pure
// unit test: it must never touch this host's actual ~/.claude).
func TestSync_DryRun_ShowsPayloadBreakdown(t *testing.T) {
	states := []ComponentState{
		{
			Name:    "deploy",
			Present: true,
			PayloadSummary: deploy.Summary{
				New: 1, Changed: 2, Same: 3, Removed: 2,
				Removals: []string{"agents/gone-one.md", "skills/gone-two/SKILL.md"},
			},
			PayloadChanged: true,
		},
	}

	var out, errOut strings.Builder
	opts := Options{
		Stdout:      &out,
		Stderr:      &errOut,
		PreDetected: states,
		DryRun:      true,
	}

	Sync(opts)

	got := out.String()
	for _, want := range []string{
		"1 nuevos, 2 cambiados, 2 a borrar, 3 iguales",
		"agents/gone-one.md",
		"skills/gone-two/SKILL.md",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected plan output to contain %q, got:\n%s", want, got)
		}
	}
}

// mustWriteFile writes data at path, creating parent directories as needed —
// used to build a real, on-disk local payload/ tree (not an fstest.MapFS),
// since the real "deploy" branch below resolves the payload via
// deploy.ResolvePayloadFS(), which walks the real filesystem.
func mustWriteFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestSync_Deploy_RealBranch_PreservesLegacyAndOverwritesEditedRegistryEntry
// exercises the real `case "deploy":` execution branch — the one place that
// calls the real deploy.Plan/deploy.Apply and prints both avisos — which no
// other test in this package reaches (the others stop at DryRun or an
// unrecognized no-op component). It is hermetic: HOME and the payload lookup
// are both redirected to temp directories for the whole test, so it never
// touches this host's actual ~/.claude or ~/.matecito-ai.
//
// Two Sync() calls model the two halves of the Requirement "preservar no es
// fallar" (update-ecosystem): the first is a migration run (no registry yet)
// that sweeps the real embedded legacy catalog and preserves one known
// orphan seeded with content that matches none of its historical hashes; the
// second is a normal run where the registry (written by the first) has an
// entry the payload no longer produces and whose disk content the person
// edited afterwards — deleted anyway per the new per-source criterion, and
// reported as an overwritten edit. Neither run reports a "deploy" error.
func TestSync_Deploy_RealBranch_PreservesLegacyAndOverwritesEditedRegistryEntry(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	payloadRoot := t.TempDir()
	t.Chdir(payloadRoot)

	payloadDir := filepath.Join(payloadRoot, "payload")
	mustWriteFile(t, filepath.Join(payloadDir, "core", "CLAUDE.md"), []byte("# core\n"))
	mustWriteFile(t, filepath.Join(payloadDir, "domains", "widgets", "agents", "foo.md"), []byte("content-v1"))

	claudeHome := filepath.Join(home, ".claude")

	// A real embedded legacy-catalog target ("claude" root), seeded with
	// content that matches none of its historical hashes: the migration run
	// below must preserve it, not delete it.
	legacyTarget := filepath.Join(claudeHome, "agents", "project-decisions-mine.md")
	mustWriteFile(t, legacyTarget, []byte("not a known historical version"))

	deployState := []ComponentState{{Name: "deploy", Present: true, PayloadChanged: true}}

	backupDir1 := t.TempDir()
	var out1 strings.Builder
	result1 := Sync(Options{
		Stdout:      &out1,
		Stderr:      &out1,
		PreDetected: deployState,
		Yes:         true,
		BackupDir:   backupDir1,
	})

	if err := result1.Errors["deploy"]; err != nil {
		t.Fatalf("migration run: unexpected deploy error: %v", err)
	}
	got1 := out1.String()
	if !strings.Contains(got1, "preservado (editado por el usuario): "+legacyTarget) {
		t.Errorf("expected the legacy orphan to be reported preserved, got:\n%s", got1)
	}
	if !strings.Contains(got1, "podés borrarlo vos") {
		t.Errorf("expected the preserved aviso to suggest the person's own action, got:\n%s", got1)
	}
	if _, err := os.Stat(legacyTarget); err != nil {
		t.Errorf("expected the preserved legacy file untouched on disk: %v", err)
	}

	fooTarget := filepath.Join(claudeHome, "agents", "foo.md")
	if _, err := os.Stat(fooTarget); err != nil {
		t.Fatalf("expected foo.md deployed by the migration run: %v", err)
	}

	// Second run: the payload no longer produces foo.md (simulates a rename/
	// removal), and the person edited foo.md's content on disk after we
	// deployed it — the registry-origin removal must delete it anyway and
	// report the overwritten edit.
	if err := os.RemoveAll(filepath.Join(payloadDir, "domains", "widgets")); err != nil {
		t.Fatalf("removing widgets domain from payload: %v", err)
	}
	if err := os.WriteFile(fooTarget, []byte("content-v1-edited-by-the-person"), 0o644); err != nil {
		t.Fatalf("simulating a user edit on foo.md: %v", err)
	}

	backupDir2 := t.TempDir()
	var out2 strings.Builder
	result2 := Sync(Options{
		Stdout:      &out2,
		Stderr:      &out2,
		PreDetected: deployState,
		Yes:         true,
		BackupDir:   backupDir2,
	})

	if err := result2.Errors["deploy"]; err != nil {
		t.Fatalf("update run: unexpected deploy error: %v", err)
	}
	got2 := out2.String()
	if !strings.Contains(got2, "edición pisada al borrar: "+fooTarget) {
		t.Errorf("expected the edited registry entry to be reported overwritten, got:\n%s", got2)
	}
	if !strings.Contains(got2, "tu versión anterior quedó respaldada en "+backupDir2) {
		t.Errorf("expected the overwritten-edit aviso to name a backup under %s, got:\n%s", backupDir2, got2)
	}
	if _, err := os.Stat(fooTarget); !os.IsNotExist(err) {
		t.Errorf("expected foo.md removed from disk despite the mismatch, stat err=%v", err)
	}
}

// TestDecide_MatecitoAIDevBuild verifies the Requirement "Precedencia de
// decisión por componente": a matecito-ai dev build is skipped regardless of
// how it differs from the latest release, the skip is scoped to matecito-ai
// only, and it takes precedence over PayloadChanged/Pending so no other
// update reason can override it.
func TestDecide_MatecitoAIDevBuild(t *testing.T) {
	cases := []struct {
		name string
		s    ComponentState
		want ActionKind
	}{
		{
			name: "dev build skipped even though latest differs",
			s:    ComponentState{Name: "matecito-ai", Present: true, CurrentVersion: "0.1.0-dev", LatestVersion: "v0.2.0"},
			want: ActionSkip,
		},
		{
			name: "non-dev build still updates normally (existing behavior intact)",
			s:    ComponentState{Name: "matecito-ai", Present: true, CurrentVersion: "0.1.0", LatestVersion: "v0.2.0"},
			want: ActionUpdate,
		},
		{
			name: "v-prefix-only difference still treated as equal and skipped",
			s:    ComponentState{Name: "matecito-ai", Present: true, CurrentVersion: "v0.2.0", LatestVersion: "0.2.0"},
			want: ActionSkip,
		},
		{
			name: "unknown latest version skips regardless of dev build",
			s:    ComponentState{Name: "matecito-ai", Present: true, CurrentVersion: "0.1.0", Unknown: true},
			want: ActionSkip,
		},
		{
			name: "dev build skip has precedence over PayloadChanged/Pending",
			s:    ComponentState{Name: "matecito-ai", Present: true, CurrentVersion: "0.1.0-dev", LatestVersion: "v0.2.0", PayloadChanged: true, Pending: true},
			want: ActionSkip,
		},
		{
			name: "the skip is scoped to matecito-ai: another component with a dev-build-shaped version updates normally",
			s:    ComponentState{Name: "engram", Present: true, CurrentVersion: "0.1.0-dev", LatestVersion: "v0.2.0"},
			want: ActionUpdate,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := decide(tc.s); got != tc.want {
				t.Errorf("decide(%+v) = %v, want %v", tc.s, got, tc.want)
			}
		})
	}
}

// TestPlanSync_MatecitoAIDevBuild_DoesNotBlockOtherComponents verifies the
// Requirement scenario "el dev build no frena a los demás": a matecito-ai dev
// build is planned as skip while another out-of-date component in the same
// run is still planned as update.
func TestPlanSync_MatecitoAIDevBuild_DoesNotBlockOtherComponents(t *testing.T) {
	states := []ComponentState{
		{Name: "matecito-ai", Present: true, CurrentVersion: "0.1.0-dev", LatestVersion: "v0.2.0"},
		{Name: "engram", Present: true, CurrentVersion: "0.1.0", LatestVersion: "v0.2.0"},
	}

	actions := PlanSync(states)

	want := map[string]ActionKind{"matecito-ai": ActionSkip, "engram": ActionUpdate}
	for _, a := range actions {
		if a.Kind != want[a.Component] {
			t.Errorf("PlanSync: %s = %v, want %v", a.Component, a.Kind, want[a.Component])
		}
	}
}

// TestPayloadChanged verifies the Requirement "El plan informa si hay trabajo
// pendiente": a payload with ONLY orphans left to sweep (New=Changed=0,
// Removed>0) counts as changed exactly like one with new or changed files.
func TestPayloadChanged(t *testing.T) {
	cases := []struct {
		name string
		s    deploy.Summary
		want bool
	}{
		{"all same, nothing pending", deploy.Summary{Same: 5}, false},
		{"only new", deploy.Summary{New: 1}, true},
		{"only changed", deploy.Summary{Changed: 1}, true},
		{"only removed (orphans only)", deploy.Summary{Same: 5, Removed: 1}, true},
		{"nothing at all", deploy.Summary{}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := payloadChanged(tc.s); got != tc.want {
				t.Errorf("payloadChanged(%+v) = %v, want %v", tc.s, got, tc.want)
			}
		})
	}
}
