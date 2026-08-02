package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/setup/deploy"
	"github.com/franwerner/matecito-ai/internal/setup/sync"
)

// TestInstallCmd_ResumeWiring documents (and guards) the precondition behind
// NewInstallCmd's "Resume: sync.ResumeRequested()" line and its
// "if !syncOpts.Resume { ...plan/dry-run/confirm... }" short-circuit: both
// read the resume flag from the same env var ResumeRequested() checks.
//
// installRunner.run (the RunE body) is exercised end to end by
// TestInstallRunner_* below via injected Stdin/Stdout and Detect/Sync, so
// this test stays narrowly about the env-var wiring itself; TestSync_Resume
// in internal/setup/sync/sync_test.go covers the actual skip/no-plan/
// no-prompt behavior at the engine level once Resume is set.
func TestInstallCmd_ResumeWiring(t *testing.T) {
	t.Run("ResumeRequested reflects MATECITO_RESUME=1", func(t *testing.T) {
		t.Setenv("MATECITO_RESUME", "1")
		if !sync.ResumeRequested() {
			t.Fatal("expected sync.ResumeRequested() to be true with MATECITO_RESUME=1, matching the value install.go assigns to syncOpts.Resume")
		}
	})

	t.Run("ResumeRequested is false without the env var", func(t *testing.T) {
		t.Setenv("MATECITO_RESUME", "")
		if sync.ResumeRequested() {
			t.Fatal("expected sync.ResumeRequested() to be false without MATECITO_RESUME=1")
		}
	})
}

// TestInstallCmd_CodegraphErrorSurfaced verifies that when syncResult carries a
// codegraph error, the install command returns a non-nil error and does NOT
// silently continue to MCP registration (requirement: CLI reports real
// codegraph install state).
func TestInstallCmd_CodegraphErrorSurfaced(t *testing.T) {
	cgErr := errors.New("npm install failed")
	result := sync.Result{
		Errors: map[string]error{
			"codegraph": cgErr,
		},
	}

	reportedErr := surfaceComponentErrors(result)
	if reportedErr == nil {
		t.Fatal("expected error to be surfaced when syncResult.Errors['codegraph'] is set, got nil")
	}
	if !strings.Contains(reportedErr.Error(), "codegraph") {
		t.Errorf("returned error should mention 'codegraph'; got: %v", reportedErr)
	}
	if !errors.Is(reportedErr, cgErr) {
		t.Errorf("returned error should wrap the original; got: %v", reportedErr)
	}
}

// TestInstallCmd_CodegraphSuccessNoError verifies that when syncResult has no
// component error, surfaceComponentErrors returns nil (normal path continues).
func TestInstallCmd_CodegraphSuccessNoError(t *testing.T) {
	result := sync.Result{}
	if err := surfaceComponentErrors(result); err != nil {
		t.Fatalf("expected nil when every component succeeded, got: %v", err)
	}
}

// TestInstallCmd_NonCodegraphFailureSurfaced covers the asymmetry install used
// to have against update: only codegraph reached the exit code, so a failed
// payload deploy printed its error and the command still exited 0.
func TestInstallCmd_NonCodegraphFailureSurfaced(t *testing.T) {
	deployErr := errors.New("no se pudo escribir el destino")
	binErr := errors.New("descarga interrumpida")
	result := sync.Result{
		Errors: map[string]error{
			"deploy": deployErr,
			"engram": binErr,
		},
	}

	reportedErr := surfaceComponentErrors(result)
	if reportedErr == nil {
		t.Fatal("expected a non-codegraph component failure to reach the exit code, got nil")
	}
	for _, want := range []string{"deploy", "engram"} {
		if !strings.Contains(reportedErr.Error(), want) {
			t.Errorf("returned error should name the failed component %q; got: %v", want, reportedErr)
		}
	}
	if !errors.Is(reportedErr, deployErr) || !errors.Is(reportedErr, binErr) {
		t.Errorf("returned error should wrap every cause; got: %v", reportedErr)
	}
}

// TestInstallCmd_NilComponentErrorIsNotAFailure guards the map-with-nil-value
// case: sync.Result.HasErrors only counts keys, so a nil entry would otherwise
// fail a run where nothing actually failed.
func TestInstallCmd_NilComponentErrorIsNotAFailure(t *testing.T) {
	result := sync.Result{Errors: map[string]error{"codegraph": nil}}
	if err := surfaceComponentErrors(result); err != nil {
		t.Fatalf("expected nil when the only recorded error is nil, got: %v", err)
	}
}

// TestInstallRunner_Run_ShowsBreakdownInOwnPreview verifies install's own
// "Plan:" preview (a renderer distinct from the sync engine's) includes the
// payload breakdown and the destinos a borrar, closing the "desglose en su
// preview propio" gap. Detect is faked to avoid the real network round-trip;
// dryRun=true stops the run before Sync would ever be invoked.
func TestInstallRunner_Run_ShowsBreakdownInOwnPreview(t *testing.T) {
	var out strings.Builder
	r := installRunner{
		Stdin:  strings.NewReader(""),
		Stdout: &out,
		Stderr: io.Discard,
		Detect: func(sync.Options) ([]sync.ComponentState, error) {
			return []sync.ComponentState{
				{
					Name:           "deploy",
					Present:        true,
					PayloadChanged: true,
					PayloadSummary: deploy.Summary{
						New: 1, Changed: 2, Same: 3, Removed: 1,
						Removals: []string{"agents/gone.md"},
					},
				},
			}, nil
		},
		Sync: func(sync.Options) sync.Result {
			t.Fatal("Sync must not be invoked on a dry-run")
			return sync.Result{}
		},
	}

	if err := r.run("v-test", true /* dryRun */, false /* yes */); err != nil {
		t.Fatalf("run: %v", err)
	}

	got := out.String()
	for _, want := range []string{
		"1 nuevos, 2 cambiados, 1 a borrar, 3 iguales",
		"agents/gone.md",
		"(dry-run)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected install's own preview to contain %q, got:\n%s", want, got)
		}
	}
}

// TestInstallRunner_Run_PlanShownOnce verifies "el plan se presenta una sola
// vez por corrida": after showing its own preview and confirming, run must
// set syncOpts.PlanShown = true before invoking Sync — the one piece of
// install.go's own wiring that could previously only be checked by manual
// code review. The fake Sync simulates the engine's own gate (it prints
// "Plan:" itself only if PlanShown is still false), so a passing test proves
// the combined output never shows the plan twice.
func TestInstallRunner_Run_PlanShownOnce(t *testing.T) {
	var out strings.Builder
	syncCalled := false
	r := installRunner{
		Stdin:  strings.NewReader(""),
		Stdout: &out,
		Stderr: io.Discard,
		Detect: func(sync.Options) ([]sync.ComponentState, error) {
			return []sync.ComponentState{{Name: "deploy", Present: true, PayloadChanged: true}}, nil
		},
		Sync: func(opts sync.Options) sync.Result {
			syncCalled = true
			if !opts.PlanShown {
				t.Error("expected syncOpts.PlanShown=true before Sync is invoked")
				fmt.Fprintln(opts.Stdout, "Plan:") // simulates the engine's own print if not suppressed
			}
			return sync.Result{}
		},
	}

	if err := r.run("v-test", false /* dryRun */, true /* yes */); err != nil {
		t.Fatalf("run: %v", err)
	}
	if !syncCalled {
		t.Fatal("expected Sync to be invoked")
	}
	if n := strings.Count(out.String(), "Plan:"); n != 1 {
		t.Fatalf("expected exactly one \"Plan:\" print across the combined output, got %d in:\n%s", n, out.String())
	}
}
