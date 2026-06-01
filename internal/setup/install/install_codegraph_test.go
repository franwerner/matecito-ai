package install_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/setup/install"
)

// writeBin writes a minimal shell script executable to dir/name that exits
// with exitCode and optionally prints a line to stdout.
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
	d, err := os.MkdirTemp("", "install-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

// TestInstallCodegraph_BinaryPresentAfterInstall verifies the happy path: npm
// exits 0 and codegraph resolves in PATH → InstallCodegraph returns nil.
func TestInstallCodegraph_BinaryPresentAfterInstall(t *testing.T) {
	d := tempDir(t)
	// npm config get prefix returns a non-system path → ensureUserNpmPrefix skips.
	// npm install -g ... succeeds.
	// Both npm subcommands are handled by the same stub; it just exits 0.
	writeBin(t, d, "npm", 0, "/tmp/npm-global")
	writeBin(t, d, "codegraph", 0, "")
	isolatedPATH(t, d)

	var out bytes.Buffer
	err := install.InstallCodegraph(install.Options{Stdout: &out, Stderr: &out})
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
}

// TestInstallCodegraph_BinaryAbsentAfterNpmSuccess verifies that when npm exits
// 0 but the codegraph binary is absent from PATH, InstallCodegraph returns a
// non-nil error mentioning "codegraph" and "PATH".
func TestInstallCodegraph_BinaryAbsentAfterNpmSuccess(t *testing.T) {
	d := tempDir(t)
	writeBin(t, d, "npm", 0, "/tmp/npm-global")
	// codegraph is intentionally NOT written to d.
	isolatedPATH(t, d) // only d is in PATH → exec.LookPath("codegraph") must fail

	var out bytes.Buffer
	err := install.InstallCodegraph(install.Options{Stdout: &out, Stderr: &out})
	if err == nil {
		t.Fatal("expected error when codegraph binary is absent after npm install, got nil")
	}
	if !strings.Contains(err.Error(), "codegraph") {
		t.Errorf("error should mention 'codegraph'; got: %v", err)
	}
	if !strings.Contains(err.Error(), "PATH") {
		t.Errorf("error should mention 'PATH'; got: %v", err)
	}
}

// TestInstallCodegraph_NpmFails verifies that when npm install exits non-zero,
// InstallCodegraph returns an error before reaching the LookPath check.
func TestInstallCodegraph_NpmFails(t *testing.T) {
	d := tempDir(t)
	writeBin(t, d, "npm", 1, "")
	// codegraph is present — if LookPath were reached it would succeed.
	writeBin(t, d, "codegraph", 0, "")
	isolatedPATH(t, d)

	var out bytes.Buffer
	err := install.InstallCodegraph(install.Options{Stdout: &out, Stderr: &out})
	if err == nil {
		t.Fatal("expected error when npm fails, got nil")
	}
	// The error must NOT be the LookPath guard message.
	if strings.Contains(err.Error(), "PATH") {
		t.Errorf("error should be from npm, not LookPath guard; got: %v", err)
	}
}

// TestCodegraphMCPStep_BinaryAbsent verifies that codegraphMCPStep.Run returns
// an error naming codegraph when the binary is absent, without invoking
// 'claude mcp add'.
func TestCodegraphMCPStep_BinaryAbsent(t *testing.T) {
	// Skip if codegraph can't be excluded (would require removing system binaries).
	origPath, err := exec.LookPath("codegraph")
	_ = origPath

	d := tempDir(t)
	writeBin(t, d, "claude", 0, "")
	// codegraph NOT written → absent from isolated PATH
	isolatedPATH(t, d)

	if err == nil {
		// codegraph was found before we isolated PATH; verify isolation works.
		if _, lookErr := exec.LookPath("codegraph"); lookErr == nil {
			t.Skip("PATH isolation did not hide codegraph; skipping test")
		}
	}

	var out bytes.Buffer
	opts := install.Options{Stdout: &out, Stderr: &out, Yes: true}
	steps := install.AllSteps(opts)

	var cgMCPRun func() error
	for _, s := range steps {
		if s.Name == "CodeGraph MCP" {
			cgMCPRun = s.Run
			break
		}
	}
	if cgMCPRun == nil {
		t.Fatal("CodeGraph MCP step not found in AllSteps")
	}

	runErr := cgMCPRun()
	if runErr == nil {
		t.Fatal("expected error when codegraph binary is absent, got nil")
	}
	if !strings.Contains(runErr.Error(), "codegraph") {
		t.Errorf("error should mention 'codegraph'; got: %v", runErr)
	}
}
