package check

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFakeExecutable writes a minimal shell script to dir/name that exits
// with exitCode and optionally prints a version line to stdout, returning its
// absolute path.
func writeFakeExecutable(t *testing.T, dir, name string, exitCode int, stdoutLine string) string {
	t.Helper()
	var sb strings.Builder
	sb.WriteString("#!/bin/sh\n")
	if stdoutLine != "" {
		sb.WriteString("echo ")
		sb.WriteString(stdoutLine)
		sb.WriteString("\n")
	}
	if exitCode != 0 {
		sb.WriteString("exit 1\n")
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(sb.String()), 0o755); err != nil {
		t.Fatalf("writeFakeExecutable %s: %v", name, err)
	}
	return path
}

// TestProbeAt covers the five scenarios spec process/binary-install-step-idempotency
// contracts for presence observed at a canonical location: resolver error,
// absent at the resolved path, present-but-unhealthy, healthy present, and
// healthy present while its directory is absent from the running session's
// PATH (proving the probe never consults PATH).
func TestProbeAt(t *testing.T) {
	t.Run("resolver error reports Missing without calling RunVersion", func(t *testing.T) {
		wantErr := errors.New("cannot resolve install location")
		resolve := func() (string, error) { return "", wantErr }

		r := ProbeAt("fake", resolve, nil, true, "install it")
		if r.Status != StatusMissing {
			t.Errorf("Status = %v, want StatusMissing", r.Status)
		}
		if r.Version != "" {
			t.Errorf("Version = %q, want empty (RunVersion must not run)", r.Version)
		}
		if r.Detail == "" {
			t.Error("Detail must not be empty on a resolver error")
		}
		if r.FixHint != "install it" {
			t.Errorf("FixHint = %q, want %q", r.FixHint, "install it")
		}
	})

	t.Run("absent at the resolved path reports Missing", func(t *testing.T) {
		dir := t.TempDir()
		absent := filepath.Join(dir, "does-not-exist")
		resolve := func() (string, error) { return absent, nil }

		r := ProbeAt("fake", resolve, nil, true, "")
		if r.Status != StatusMissing {
			t.Errorf("Status = %v, want StatusMissing", r.Status)
		}
		if r.Version != "" {
			t.Errorf("Version = %q, want empty", r.Version)
		}
	})

	t.Run("present but exits non-zero reports Missing", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFakeExecutable(t, dir, "fake", 1, "")
		resolve := func() (string, error) { return path, nil }

		r := ProbeAt("fake", resolve, nil, true, "")
		if r.Status != StatusMissing {
			t.Errorf("Status = %v, want StatusMissing", r.Status)
		}
		if r.Version != "" {
			t.Errorf("Version = %q, want empty", r.Version)
		}
	})

	t.Run("healthy present reports OK with the parsed version and the resolved path as Detail", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFakeExecutable(t, dir, "fake", 0, "1.2.3")
		resolve := func() (string, error) { return path, nil }

		r := ProbeAt("fake", resolve, nil, true, "")
		if r.Status != StatusOK {
			t.Errorf("Status = %v, want StatusOK", r.Status)
		}
		if r.Version != "1.2.3" {
			t.Errorf("Version = %q, want %q", r.Version, "1.2.3")
		}
		if r.Detail != path {
			t.Errorf("Detail = %q, want the resolved path %q", r.Detail, path)
		}
		if r.Detail == notFoundInPATHDetail {
			t.Errorf("Detail must never be %q for a healthy probe", notFoundInPATHDetail)
		}
	})

	t.Run("healthy present even when its directory is absent from the session's PATH", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFakeExecutable(t, dir, "fake", 0, "1.0.0")
		resolve := func() (string, error) { return path, nil }

		origPath := os.Getenv("PATH")
		os.Setenv("PATH", "")
		t.Cleanup(func() { os.Setenv("PATH", origPath) })

		r := ProbeAt("fake", resolve, nil, true, "")
		if r.Status != StatusOK {
			t.Errorf("Status = %v, want StatusOK even though the resolved dir is absent from PATH", r.Status)
		}
		if r.Version != "1.0.0" {
			t.Errorf("Version = %q, want %q", r.Version, "1.0.0")
		}
	})
}

// TestRunVersion_Unaffected verifies ProbeAt's addition leaves RunVersion's
// own bare-name behavior byte-identical — the const hoist must not change
// what callers like internal/checks/prereqs observe.
func TestRunVersion_Unaffected(t *testing.T) {
	r := RunVersion("nonexistent-binary-xyz", "nonexistent-binary-xyz", nil, true, "install it")
	if r.Status != StatusMissing {
		t.Errorf("Status = %v, want StatusMissing", r.Status)
	}
	if r.Detail != notFoundInPATHDetail {
		t.Errorf("Detail = %q, want %q", r.Detail, notFoundInPATHDetail)
	}
}
