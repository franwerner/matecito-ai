package sync

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/setup/install"
)

// writeBin writes a minimal shell script executable to dir/name that exits
// with exitCode and optionally prints a line to stdout. Local to this
// package: the existing helpers of the same name in
// internal/setup/install/install_codegraph_test.go live in package
// install_test and are unreachable from here.
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
	d, err := os.MkdirTemp("", "sync-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

// binaryDetectionCase is the twin of the scenario table walked in
// internal/checks/*/..._internal_test.go — same names, same five cases from
// spec process/binary-install-step-idempotency. A divergence between the
// install surface and the verify surface fails on exactly one side.
type binaryDetectionCase struct {
	name        string
	setup       func(t *testing.T) func() (string, error) // returns the resolver to probe with
	wantPresent bool
}

func canonicalPresentCase(binName, versionLine string) binaryDetectionCase {
	return binaryDetectionCase{
		name: "canonical-present",
		setup: func(t *testing.T) func() (string, error) {
			d := tempDir(t)
			writeBin(t, d, binName, 0, versionLine)
			path := filepath.Join(d, binName)
			return func() (string, error) { return path, nil }
		},
		wantPresent: true,
	}
}

func canonicalAbsentCase(binName string) binaryDetectionCase {
	return binaryDetectionCase{
		name: "canonical-absent",
		setup: func(t *testing.T) func() (string, error) {
			d := tempDir(t)
			path := filepath.Join(d, binName) // never written
			return func() (string, error) { return path, nil }
		},
		wantPresent: false,
	}
}

func canonicalButBrokenCase(binName string) binaryDetectionCase {
	return binaryDetectionCase{
		name: "canonical-but-broken",
		setup: func(t *testing.T) func() (string, error) {
			d := tempDir(t)
			writeBin(t, d, binName, 1, "")
			path := filepath.Join(d, binName)
			return func() (string, error) { return path, nil }
		},
		wantPresent: false,
	}
}

func pathNotOnSessionPATHCase(binName, versionLine string) binaryDetectionCase {
	return binaryDetectionCase{
		name: "path-not-yet-on-session-PATH",
		setup: func(t *testing.T) func() (string, error) {
			d := tempDir(t)
			writeBin(t, d, binName, 0, versionLine)
			path := filepath.Join(d, binName)
			isolatedPATH(t, tempDir(t)) // deliberately excludes d
			return func() (string, error) { return path, nil }
		},
		wantPresent: true,
	}
}

func resolverFailsCase() binaryDetectionCase {
	return binaryDetectionCase{
		name: "resolver-fails",
		setup: func(t *testing.T) func() (string, error) {
			return func() (string, error) { return "", errors.New("cannot resolve install location") }
		},
		wantPresent: false,
	}
}

// runBinaryDetectionMatrix exercises probeInstalledBinary over the five
// canonical scenarios for one binary name.
func runBinaryDetectionMatrix(t *testing.T, binName, versionLine string, versionArgs []string) {
	cases := []binaryDetectionCase{
		canonicalPresentCase(binName, versionLine),
		canonicalAbsentCase(binName),
		canonicalButBrokenCase(binName),
		pathNotOnSessionPATHCase(binName, versionLine),
		resolverFailsCase(),
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resolve := tc.setup(t)
			present, _ := probeInstalledBinary(binName, resolve, versionArgs)
			if present != tc.wantPresent {
				t.Errorf("probeInstalledBinary(%s) present = %v, want %v", tc.name, present, tc.wantPresent)
			}
		})
	}
}

// TestProbeInstalledBinary_Engram mirrors the five spec scenarios through the
// engram detection path. Case names identical to
// internal/checks/engram/engram_internal_test.go's matrix (task 6.4).
func TestProbeInstalledBinary_Engram(t *testing.T) {
	runBinaryDetectionMatrix(t, "engram", "1.20.0", []string{"version"})
}

// TestProbeInstalledBinary_Codegraph mirrors the five spec scenarios through
// the codegraph detection path, plus a real-resolver case using a fake npm on
// an isolated PATH — install.CodegraphBinaryPath resolves through it. Case
// names identical to internal/checks/codegraph/codegraph_internal_test.go's
// matrix (task 6.5).
func TestProbeInstalledBinary_Codegraph(t *testing.T) {
	runBinaryDetectionMatrix(t, "codegraph", "2.0.0", []string{"--version"})

	t.Run("real-resolver", func(t *testing.T) {
		npmDir := tempDir(t)
		prefix := tempDir(t)
		writeBin(t, npmDir, "npm", 0, prefix)
		canonicalBinDir := filepath.Join(prefix, "bin")
		if err := os.MkdirAll(canonicalBinDir, 0o755); err != nil {
			t.Fatalf("mkdir canonical bin dir: %v", err)
		}
		writeBin(t, canonicalBinDir, "codegraph", 0, "2.0.0")
		isolatedPATH(t, canonicalBinDir, npmDir)

		present, version := probeInstalledBinary("codegraph", install.CodegraphBinaryPath, []string{"--version"})
		if !present {
			t.Error("expected present=true through the real npm resolver")
		}
		if version != "2.0.0" {
			t.Errorf("version = %q, want %q", version, "2.0.0")
		}
	})
}

// TestProbeInstalledBinary_Proofshot mirrors the five spec scenarios through
// the proofshot detection path, plus a real-resolver case, mirroring the
// codegraph one. Case names identical to
// internal/checks/proofshot/proofshot_internal_test.go's matrix (task 6.6).
func TestProbeInstalledBinary_Proofshot(t *testing.T) {
	runBinaryDetectionMatrix(t, "proofshot", "1.0.0", []string{"--version"})

	t.Run("real-resolver", func(t *testing.T) {
		npmDir := tempDir(t)
		prefix := tempDir(t)
		writeBin(t, npmDir, "npm", 0, prefix)
		canonicalBinDir := filepath.Join(prefix, "bin")
		if err := os.MkdirAll(canonicalBinDir, 0o755); err != nil {
			t.Fatalf("mkdir canonical bin dir: %v", err)
		}
		writeBin(t, canonicalBinDir, "proofshot", 0, "1.0.0")
		isolatedPATH(t, canonicalBinDir, npmDir)

		present, version := probeInstalledBinary("proofshot", install.ProofshotBinaryPath, []string{"--version"})
		if !present {
			t.Error("expected present=true through the real npm resolver")
		}
		if version != "1.0.0" {
			t.Errorf("version = %q, want %q", version, "1.0.0")
		}
	})
}

// TestProbeInstalledBinary_Independence verifies one binary's unresolvable
// path leaves the other two verdicts intact — requirement "El lugar de
// instalación es propio de cada paso".
func TestProbeInstalledBinary_Independence(t *testing.T) {
	engramDir := tempDir(t)
	writeBin(t, engramDir, "engram", 0, "1.0.0")
	engramPath := filepath.Join(engramDir, "engram")

	proofshotDir := tempDir(t)
	writeBin(t, proofshotDir, "proofshot", 0, "1.0.0")
	proofshotPath := filepath.Join(proofshotDir, "proofshot")

	engramPresent, _ := probeInstalledBinary("engram", func() (string, error) { return engramPath, nil }, []string{"version"})
	codegraphPresent, _ := probeInstalledBinary("codegraph", func() (string, error) { return "", errors.New("cannot resolve") }, []string{"--version"})
	proofshotPresent, _ := probeInstalledBinary("proofshot", func() (string, error) { return proofshotPath, nil }, []string{"--version"})

	if !engramPresent {
		t.Error("expected engram present, independent of codegraph's failure")
	}
	if codegraphPresent {
		t.Error("expected codegraph missing: its resolver failed")
	}
	if !proofshotPresent {
		t.Error("expected proofshot present, independent of codegraph's failure")
	}
}

// TestProbeInstalledBinary_SecondRunIdempotent verifies two consecutive calls
// over the same healthy canonical binary both report present, with no state
// mutation between calls.
func TestProbeInstalledBinary_SecondRunIdempotent(t *testing.T) {
	d := tempDir(t)
	writeBin(t, d, "engram", 0, "1.0.0")
	path := filepath.Join(d, "engram")
	resolve := func() (string, error) { return path, nil }

	present1, version1 := probeInstalledBinary("engram", resolve, []string{"version"})
	present2, version2 := probeInstalledBinary("engram", resolve, []string{"version"})

	if !present1 || !present2 {
		t.Errorf("expected present on both calls, got present1=%v present2=%v", present1, present2)
	}
	if version1 != version2 {
		t.Errorf("version drifted between calls: %q then %q", version1, version2)
	}
}
