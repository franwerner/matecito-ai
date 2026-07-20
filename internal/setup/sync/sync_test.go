package sync

import (
	"io"
	"strings"
	"testing"
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
