package engram

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/check"
)

// writeBin writes a minimal shell script executable to dir/name that exits
// with exitCode and optionally prints a line to stdout. Local to this
// package: internal/setup/sync/binary_detection_test.go's helper of the same
// name lives in package sync and is unreachable from here.
func writeBin(t *testing.T, dir, name string, exitCode int, stdoutLine string) {
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
		t.Fatalf("writeBin %s: %v", name, err)
	}
}

// isolatedPATH replaces the process PATH with only the given dirs for the
// duration of the test. The original PATH is restored on Cleanup.
func isolatedPATH(t *testing.T, dirs ...string) {
	t.Helper()
	orig := os.Getenv("PATH")
	os.Setenv("PATH", strings.Join(dirs, string(os.PathListSeparator)))
	t.Cleanup(func() { os.Setenv("PATH", orig) })
}

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "engram-check-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

// TestDetectBinary mirrors the five spec scenarios of
// process/binary-install-step-idempotency through the resolveBinary seam.
// Case names identical to internal/setup/sync/binary_detection_test.go's
// engram matrix (task 3.3), so a divergence between the two surfaces fails on
// exactly one side.
func TestDetectBinary(t *testing.T) {
	orig := resolveBinary
	t.Cleanup(func() { resolveBinary = orig })

	t.Run("canonical-present", func(t *testing.T) {
		d := tempDir(t)
		writeBin(t, d, "engram", 0, "1.20.0")
		path := filepath.Join(d, "engram")
		resolveBinary = func() (string, error) { return path, nil }

		r := detectBinary()
		if r.Status != check.StatusOK {
			t.Errorf("Status = %v, want StatusOK", r.Status)
		}
		if r.Detail != path {
			t.Errorf("Detail = %q, want the resolved path %q", r.Detail, path)
		}
	})

	t.Run("canonical-absent", func(t *testing.T) {
		d := tempDir(t)
		path := filepath.Join(d, "engram") // never written
		resolveBinary = func() (string, error) { return path, nil }

		r := detectBinary()
		if r.Status != check.StatusMissing {
			t.Errorf("Status = %v, want StatusMissing", r.Status)
		}
	})

	t.Run("canonical-but-broken", func(t *testing.T) {
		d := tempDir(t)
		writeBin(t, d, "engram", 1, "")
		path := filepath.Join(d, "engram")
		resolveBinary = func() (string, error) { return path, nil }

		r := detectBinary()
		if r.Status != check.StatusMissing {
			t.Errorf("Status = %v, want StatusMissing", r.Status)
		}
	})

	t.Run("path-not-yet-on-session-PATH", func(t *testing.T) {
		d := tempDir(t)
		writeBin(t, d, "engram", 0, "1.20.0")
		path := filepath.Join(d, "engram")
		isolatedPATH(t, tempDir(t)) // deliberately excludes d
		resolveBinary = func() (string, error) { return path, nil }

		r := detectBinary()
		if r.Status != check.StatusOK {
			t.Errorf("Status = %v, want StatusOK (absolute path bypasses PATH lookup)", r.Status)
		}
	})

	t.Run("resolver-fails", func(t *testing.T) {
		resolveBinary = func() (string, error) { return "", errors.New("cannot resolve install location") }

		r := detectBinary()
		if r.Status != check.StatusMissing {
			t.Errorf("Status = %v, want StatusMissing", r.Status)
		}
		if r.FixHint == "" {
			t.Error("StatusMissing result must carry a non-empty FixHint")
		}
	})
}
