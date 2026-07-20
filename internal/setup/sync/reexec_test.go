package sync

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"testing"
)

// TestResumeEnv covers the pure env-injection helper: it must add the resume
// flag, leave the rest of the env untouched, never mutate its input, and
// never duplicate the flag when it is already present.
func TestResumeEnv(t *testing.T) {
	t.Run("nil env gets the flag", func(t *testing.T) {
		got := resumeEnv(nil)
		want := []string{resumeEnvVar + "=1"}
		if !slices.Equal(got, want) {
			t.Fatalf("resumeEnv(nil) = %v, want %v", got, want)
		}
	})

	t.Run("preserves the rest of the env", func(t *testing.T) {
		in := []string{"PATH=/usr/bin", "HOME=/home/user"}
		got := resumeEnv(in)
		want := []string{"PATH=/usr/bin", "HOME=/home/user", resumeEnvVar + "=1"}
		if !slices.Equal(got, want) {
			t.Fatalf("resumeEnv(%v) = %v, want %v", in, got, want)
		}
		// The input slice must not be mutated in place.
		if !slices.Equal(in, []string{"PATH=/usr/bin", "HOME=/home/user"}) {
			t.Fatalf("resumeEnv mutated its input, got %v", in)
		}
	})

	t.Run("does not duplicate the flag when already present", func(t *testing.T) {
		in := []string{"PATH=/usr/bin", resumeEnvVar + "=1"}
		got := resumeEnv(in)
		want := []string{"PATH=/usr/bin", resumeEnvVar + "=1"}
		if !slices.Equal(got, want) {
			t.Fatalf("resumeEnv(%v) = %v, want %v (no duplicate flag)", in, got, want)
		}
	})
}

// TestResumeRequested covers detection of a resumed invocation via the
// resume env var, in both directions.
func TestResumeRequested(t *testing.T) {
	t.Run("true when the flag is set to 1", func(t *testing.T) {
		t.Setenv(resumeEnvVar, "1")
		if !ResumeRequested() {
			t.Fatal("expected ResumeRequested() to be true when the env var is 1")
		}
	})

	t.Run("false when the flag is unset", func(t *testing.T) {
		t.Setenv(resumeEnvVar, "")
		if ResumeRequested() {
			t.Fatal("expected ResumeRequested() to be false when the env var is empty/unset")
		}
	})

	t.Run("false when the flag holds an unexpected value", func(t *testing.T) {
		t.Setenv(resumeEnvVar, "0")
		if ResumeRequested() {
			t.Fatal("expected ResumeRequested() to be false when the env var is not \"1\"")
		}
	})
}

// stubReExecFn swaps the reExecFn seam for the duration of a test and returns
// a restore func; real syscall.Exec/spawn are never invoked by these tests.
func stubReExecFn(fake func() error) func() {
	orig := reExecFn
	reExecFn = fake
	return func() { reExecFn = orig }
}

// TestFinishSelfReplace covers the three branches spelled out in design #952:
// no-op when nothing was replaced, invoke the seam when it was, and degrade
// to the manual-rerun message (without a hard error) when the seam fails.
func TestFinishSelfReplace(t *testing.T) {
	t.Run("no-op when replaced is false", func(t *testing.T) {
		called := false
		defer stubReExecFn(func() error { called = true; return nil })()

		var out, errOut bytes.Buffer
		if err := FinishSelfReplace(&out, &errOut, false); err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if called {
			t.Fatal("expected reExecFn NOT to be called when replaced=false")
		}
		if out.Len() != 0 {
			t.Fatalf("expected no output when replaced=false, got %q", out.String())
		}
	})

	t.Run("invokes reExecFn when replaced is true", func(t *testing.T) {
		called := false
		defer stubReExecFn(func() error { called = true; return nil })()

		var out, errOut bytes.Buffer
		if err := FinishSelfReplace(&out, &errOut, true); err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !called {
			t.Fatal("expected reExecFn to be called when replaced=true")
		}
	})

	t.Run("falls back to the manual-rerun message on reExecFn error", func(t *testing.T) {
		defer stubReExecFn(func() error { return errors.New("exec failed") })()

		var out, errOut bytes.Buffer
		err := FinishSelfReplace(&out, &errOut, true)
		if err != nil {
			t.Fatalf("expected nil error even when reExecFn fails (graceful fallback), got %v", err)
		}
		if !strings.Contains(out.String(), "re-ejecut") {
			t.Fatalf("expected the manual-rerun fallback message in out, got %q", out.String())
		}
	})
}
