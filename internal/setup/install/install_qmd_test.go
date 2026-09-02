package install_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/mcp"
	"github.com/franwerner/matecito-ai/internal/setup/install"
)

// qmdStepName is the Step.Name qmdMCPStep registers under, used to look it up
// from install.AllSteps like every other MCP step test in this package.
const qmdStepName = "qmd MCP (record search)"

// findQmdStep locates the qmd MCP step in AllSteps, failing the test if absent.
func findQmdStep(t *testing.T, opts install.Options) install.Step {
	t.Helper()
	for _, s := range install.AllSteps(opts) {
		if s.Name == qmdStepName {
			return s
		}
	}
	t.Fatal("qmd MCP step not found in AllSteps")
	return install.Step{}
}

// releasePayload builds a minimal valid GitHub "latest release" API response
// with the given tag and assets, mirroring releasedl_test.go's apiPayload.
func releasePayload(tag string, assets map[string]string) []byte {
	list := make([]map[string]any, 0, len(assets))
	for name, url := range assets {
		list = append(list, map[string]any{"name": name, "browser_download_url": url})
	}
	payload := map[string]any{"tag_name": tag, "assets": list}
	b, _ := json.Marshal(payload)
	return b
}

// githubAPIRedirect is an http.RoundTripper that redirects any request bound
// for api.github.com to a local httptest.Server. qmdLatestTarballURL's
// production call path (Run calls it with no apiBaseURL override) hardcodes
// https://api.github.com, so this is what lets Run's network step be tested
// against a local server instead of the real network, without touching
// production code.
type githubAPIRedirect struct {
	base   http.RoundTripper
	target *url.URL
}

func (r *githubAPIRedirect) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Host == "api.github.com" {
		req = req.Clone(req.Context())
		req.URL.Scheme = r.target.Scheme
		req.URL.Host = r.target.Host
	}
	return r.base.RoundTrip(req)
}

// stubGitHubReleaseServer stands up a local server serving the given handler
// and redirects api.github.com requests to it for the test's duration.
func stubGitHubReleaseServer(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	target, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("parse httptest server URL: %v", err)
	}
	orig := http.DefaultTransport
	http.DefaultTransport = &githubAPIRedirect{base: orig, target: target}
	t.Cleanup(func() { http.DefaultTransport = orig })
}

// qmdReleaseHandler serves a valid latest-release payload with a qmd.tgz asset.
func qmdReleaseHandler(tag, assetURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(releasePayload(tag, map[string]string{
			"qmd.tgz": assetURL,
		}))
	}
}

// isolateMCPFind points the mcp package's HOME-based JSON lookup at a fresh,
// isolated directory and, when registered is true, seeds it with a ~/.claude.json
// that declares "qmd" — so mcp.Find("qmd") reflects only what the test wants,
// never the real machine's own registration. mcp.Find checks this JSON source
// before falling back to `claude mcp list`, so a "claude" binary absent from
// the isolated PATH (the tests below never write one for Check-only cases)
// simply makes the CLI fallback report not-found, which is what "registered =
// false" needs anyway.
func isolateMCPFind(t *testing.T, registered bool) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	if registered {
		doc := []byte(`{"mcpServers":{"qmd":{"type":"stdio","command":"qmd","args":["mcp"]}}}`)
		if err := os.WriteFile(filepath.Join(home, ".claude.json"), doc, 0o600); err != nil {
			t.Fatalf("write .claude.json: %v", err)
		}
	}
	mcp.InvalidateCLICache()
	t.Cleanup(func() { mcp.InvalidateCLICache() })
}

// --- Check semantics ---

// TestQmdMCPStep_Check_RegistrationAbsent verifies Check reports pending when
// the registration is absent, regardless of whether the binary is on PATH.
func TestQmdMCPStep_Check_RegistrationAbsent(t *testing.T) {
	isolateMCPFind(t, false) // registration absent → mcp.Find("qmd") reports absent

	d := tempDir(t)
	writeBin(t, d, "qmd", 0, "") // binary present; registration is what's missing
	isolatedPATH(t, d)

	opts := install.Options{Yes: true}
	step := findQmdStep(t, opts)
	if !step.Check() {
		t.Error("expected Check=true when registration is absent, got false")
	}
}

// TestQmdMCPStep_Check_RegistrationPresentBinaryAbsent verifies Check reports
// pending when the registration exists but the qmd binary is not on PATH —
// the case this step's widened Check exists to close (design decision
// structure/mcp-step-guards-both-artifacts).
func TestQmdMCPStep_Check_RegistrationPresentBinaryAbsent(t *testing.T) {
	isolateMCPFind(t, true)

	d := tempDir(t)
	// qmd is intentionally NOT written to d → LookPath("qmd") fails.
	isolatedPATH(t, d)

	opts := install.Options{Yes: true}
	step := findQmdStep(t, opts)
	if !step.Check() {
		t.Error("expected Check=true when registration is present but the binary is missing, got false")
	}
}

// TestQmdMCPStep_Check_BothPresent verifies Check reports nothing pending when
// the registration is present and the qmd binary resolves inside the
// canonical npm bin dir and runs — the good case under the widened,
// provenance-aware Check (design decision structure/mcp-step-guards-both-artifacts,
// revised). The fake qmd lives in the canonical bin dir (mirroring the npm
// stub's "npm config get prefix" answer) rather than an arbitrary temp dir,
// so it is not mistaken for a foreign executable.
func TestQmdMCPStep_Check_BothPresent(t *testing.T) {
	isolateMCPFind(t, true)

	npmDir := tempDir(t)
	writeBin(t, npmDir, "npm", 0, "/tmp/npm-global")
	canonicalBinDir := filepath.Join("/tmp/npm-global", "bin")
	if err := os.MkdirAll(canonicalBinDir, 0o755); err != nil {
		t.Fatalf("mkdir canonical bin dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll("/tmp/npm-global") })
	writeBin(t, canonicalBinDir, "qmd", 0, "")
	isolatedPATH(t, canonicalBinDir, npmDir)

	opts := install.Options{Yes: true}
	step := findQmdStep(t, opts)
	if step.Check() {
		t.Error("expected Check=false when registration is present and the binary is in the canonical dir and runs, got true")
	}
}

// TestQmdMCPStep_Check_ForeignLocation verifies Check reports pending when the
// registration is present and the qmd winning on PATH runs fine but lives
// outside the canonical npm bin dir — the incident this widened Check exists
// to close (design decision structure/mcp-step-guards-both-artifacts, revised).
func TestQmdMCPStep_Check_ForeignLocation(t *testing.T) {
	isolateMCPFind(t, true)

	npmDir := tempDir(t)
	writeBin(t, npmDir, "npm", 0, "/tmp/npm-global")
	foreignDir := tempDir(t)
	writeBin(t, foreignDir, "qmd", 0, "") // runs fine, but not in the canonical dir
	isolatedPATH(t, foreignDir, npmDir)

	opts := install.Options{Yes: true}
	step := findQmdStep(t, opts)
	if !step.Check() {
		t.Error("expected Check=true when the winning qmd lives outside the canonical bin dir, got false")
	}
}

// TestQmdMCPStep_Check_CanonicalButBroken verifies Check reports pending when
// the qmd in the canonical bin dir is present but fails to run — the health
// clause, catching the incident's "Cannot find package 'fast-glob'" symptom.
func TestQmdMCPStep_Check_CanonicalButBroken(t *testing.T) {
	isolateMCPFind(t, true)

	npmDir := tempDir(t)
	writeBin(t, npmDir, "npm", 0, "/tmp/npm-global")
	canonicalBinDir := filepath.Join("/tmp/npm-global", "bin")
	if err := os.MkdirAll(canonicalBinDir, 0o755); err != nil {
		t.Fatalf("mkdir canonical bin dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll("/tmp/npm-global") })
	writeBin(t, canonicalBinDir, "qmd", 1, "") // in the canonical dir, but exits non-zero
	isolatedPATH(t, canonicalBinDir, npmDir)

	opts := install.Options{Yes: true}
	step := findQmdStep(t, opts)
	if !step.Check() {
		t.Error("expected Check=true when the canonical qmd fails to run, got false")
	}
}

// TestQmdMCPStep_Check_CanonicalPresentButDirNotOnPATH verifies Check reports
// nothing pending when the registration is present and a healthy qmd sits in
// the canonical npm bin dir, even though that dir is absent from the running
// session's PATH — the fresh-HOME defect this change fixes (design decision
// structure/mcp-step-guards-both-artifacts, revised): what a prior run
// installed and persisted for future sessions does not depend on the current
// session's PATH already having picked it up.
func TestQmdMCPStep_Check_CanonicalPresentButDirNotOnPATH(t *testing.T) {
	isolateMCPFind(t, true)

	npmDir := tempDir(t)
	writeBin(t, npmDir, "npm", 0, "/tmp/npm-global")
	canonicalBinDir := filepath.Join("/tmp/npm-global", "bin")
	if err := os.MkdirAll(canonicalBinDir, 0o755); err != nil {
		t.Fatalf("mkdir canonical bin dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll("/tmp/npm-global") })
	writeBin(t, canonicalBinDir, "qmd", 0, "") // healthy, but its dir is deliberately not on PATH
	isolatedPATH(t, npmDir)                    // canonical dir excluded on purpose

	opts := install.Options{Yes: true}
	step := findQmdStep(t, opts)
	if step.Check() {
		t.Error("expected Check=false when a healthy canonical qmd's dir is absent from PATH, got true")
	}
}

// TestQmdMCPStep_Check_CanonicalAbsentWithResolverAvailable verifies Check
// still reports pending when the registration and resolver both work but no
// qmd was ever installed into the canonical dir — the widened presence probe
// must not be read as "never pending" (design decision
// structure/mcp-step-guards-both-artifacts, revised).
func TestQmdMCPStep_Check_CanonicalAbsentWithResolverAvailable(t *testing.T) {
	isolateMCPFind(t, true)

	npmDir := tempDir(t)
	writeBin(t, npmDir, "npm", 0, "/tmp/npm-global")
	canonicalBinDir := filepath.Join("/tmp/npm-global", "bin")
	if err := os.MkdirAll(canonicalBinDir, 0o755); err != nil {
		t.Fatalf("mkdir canonical bin dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll("/tmp/npm-global") })
	// canonicalBinDir left empty — no qmd binary ever written into it.
	isolatedPATH(t, npmDir)

	opts := install.Options{Yes: true}
	step := findQmdStep(t, opts)
	if !step.Check() {
		t.Error("expected Check=true when the canonical dir has no qmd binary, got false")
	}
}

// --- Run failure paths ---

// TestQmdMCPStep_Run_NpmAbsent verifies Run returns a named error when npm is
// not on PATH, before any network call is attempted.
func TestQmdMCPStep_Run_NpmAbsent(t *testing.T) {
	d := tempDir(t)
	// npm is intentionally NOT written → absent from isolated PATH.
	isolatedPATH(t, d)

	var out bytes.Buffer
	opts := install.Options{Stdout: &out, Stderr: &out, Yes: true}
	step := findQmdStep(t, opts)

	err := step.Run()
	if err == nil {
		t.Fatal("expected error when npm is absent, got nil")
	}
	if !strings.Contains(err.Error(), "npm") {
		t.Errorf("error should mention 'npm'; got: %v", err)
	}
}

// TestQmdMCPStep_Run_BinaryAbsentAfterNpmSuccess verifies Run returns a named
// error when npm exits 0 (the release resolved and "installed") but the qmd
// binary does not land on PATH.
func TestQmdMCPStep_Run_BinaryAbsentAfterNpmSuccess(t *testing.T) {
	stubGitHubReleaseServer(t, qmdReleaseHandler("v2.8.3-mate.6", "https://example.com/qmd.tgz"))

	d := tempDir(t)
	writeBin(t, d, "npm", 0, "/tmp/npm-global")
	// qmd is intentionally NOT written → absent from isolated PATH after "install".
	isolatedPATH(t, d)

	var out bytes.Buffer
	opts := install.Options{Stdout: &out, Stderr: &out, Yes: true}
	step := findQmdStep(t, opts)

	err := step.Run()
	if err == nil {
		t.Fatal("expected error when qmd binary is absent after npm install, got nil")
	}
	if !strings.Contains(err.Error(), "qmd") {
		t.Errorf("error should mention 'qmd'; got: %v", err)
	}
	if !strings.Contains(err.Error(), "PATH") {
		t.Errorf("error should mention 'PATH'; got: %v", err)
	}
}

// TestQmdMCPStep_Run_ClaudeAbsent verifies Run returns a named error when
// claude is not on PATH, reached only after the release resolved and the
// binary landed — i.e. asset resolution succeeded first.
func TestQmdMCPStep_Run_ClaudeAbsent(t *testing.T) {
	isolateMCPFind(t, false) // registration absent → Run reaches the claude mcp add branch
	stubGitHubReleaseServer(t, qmdReleaseHandler("v2.8.3-mate.6", "https://example.com/qmd.tgz"))

	d := tempDir(t)
	writeBin(t, d, "npm", 0, "/tmp/npm-global")
	writeBin(t, d, "qmd", 0, "")
	// claude is intentionally NOT written → absent from isolated PATH.
	isolatedPATH(t, d)

	var out bytes.Buffer
	opts := install.Options{Stdout: &out, Stderr: &out, Yes: true}
	step := findQmdStep(t, opts)

	err := step.Run()
	if err == nil {
		t.Fatal("expected error when claude is absent, got nil")
	}
	if !strings.Contains(err.Error(), "claude") {
		t.Errorf("error should mention 'claude'; got: %v", err)
	}
}

// TestQmdMCPStep_Run_AlreadyRegistered_SkipsClaudeMCPAdd verifies a repair run
// (binary missing, registration already present) reinstalls the binary and
// returns nil without requiring claude on PATH at all — Run must not attempt
// a duplicate `claude mcp add`. The pre-existing qmd lives in the canonical
// npm bin dir (mirroring the npm stub's prefix answer) — not an arbitrary
// temp dir — so decision 2's shadowing-executable check does not fire here;
// that check has its own dedicated cases below.
func TestQmdMCPStep_Run_AlreadyRegistered_SkipsClaudeMCPAdd(t *testing.T) {
	isolateMCPFind(t, true)
	stubGitHubReleaseServer(t, qmdReleaseHandler("v2.8.3-mate.6", "https://example.com/qmd.tgz"))

	d := tempDir(t)
	writeBin(t, d, "npm", 0, "/tmp/npm-global")
	canonicalBinDir := filepath.Join("/tmp/npm-global", "bin")
	if err := os.MkdirAll(canonicalBinDir, 0o755); err != nil {
		t.Fatalf("mkdir canonical bin dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll("/tmp/npm-global") })
	writeBin(t, canonicalBinDir, "qmd", 0, "")
	// claude is intentionally NOT written — if Run tried to register again it
	// would fail on this LookPath; its absence proves the branch was skipped.
	isolatedPATH(t, canonicalBinDir, d)

	var out bytes.Buffer
	opts := install.Options{Stdout: &out, Stderr: &out, Yes: true}
	step := findQmdStep(t, opts)

	if err := step.Run(); err != nil {
		t.Fatalf("expected nil on a repair run with registration already present, got: %v", err)
	}
}

// --- Asset resolution (exercised through Run, against a local GitHub API stub) ---

// TestQmdMCPStep_Run_AssetAbsent verifies Run fails naming the release tag and
// listing the available assets when the latest release has no qmd.tgz asset —
// this is the "asset resolution: qmd.tgz absent" scenario.
func TestQmdMCPStep_Run_AssetAbsent(t *testing.T) {
	stubGitHubReleaseServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(releasePayload("v9.9.9", map[string]string{
			"tobilu-qmd-9.9.9.tgz": "https://example.com/tobilu-qmd-9.9.9.tgz",
		}))
	})

	d := tempDir(t)
	writeBin(t, d, "npm", 0, "/tmp/npm-global")
	isolatedPATH(t, d)

	var out bytes.Buffer
	opts := install.Options{Stdout: &out, Stderr: &out, Yes: true}
	step := findQmdStep(t, opts)

	err := step.Run()
	if err == nil {
		t.Fatal("expected error when the latest release has no qmd.tgz asset, got nil")
	}
	if !strings.Contains(err.Error(), "v9.9.9") {
		t.Errorf("error should mention the release tag %q; got: %v", "v9.9.9", err)
	}
	if !strings.Contains(err.Error(), "tobilu-qmd-9.9.9.tgz") {
		t.Errorf("error should list the available asset name; got: %v", err)
	}
}

// TestQmdMCPStep_Run_AssetFound_ReachesInstall verifies the happy path of
// asset resolution: when qmd.tgz IS present, Run proceeds past resolution and
// into the npm install step (observed here by npm exiting 0 and the binary
// landing, i.e. no asset-resolution error at all). The pre-existing qmd lives
// in the canonical npm bin dir so decision 2's shadowing-executable check
// does not fire — that behavior has its own dedicated cases below.
func TestQmdMCPStep_Run_AssetFound_ReachesInstall(t *testing.T) {
	isolateMCPFind(t, true) // already registered → no claude needed
	stubGitHubReleaseServer(t, qmdReleaseHandler("v2.8.3-mate.6", "https://example.com/qmd.tgz"))

	d := tempDir(t)
	writeBin(t, d, "npm", 0, "/tmp/npm-global")
	canonicalBinDir := filepath.Join("/tmp/npm-global", "bin")
	if err := os.MkdirAll(canonicalBinDir, 0o755); err != nil {
		t.Fatalf("mkdir canonical bin dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll("/tmp/npm-global") })
	writeBin(t, canonicalBinDir, "qmd", 0, "")
	isolatedPATH(t, canonicalBinDir, d)

	var out bytes.Buffer
	opts := install.Options{Stdout: &out, Stderr: &out, Yes: true}
	step := findQmdStep(t, opts)

	if err := step.Run(); err != nil {
		t.Fatalf("expected nil once qmd.tgz resolves and npm install succeeds, got: %v", err)
	}
}

// TestQmdMCPStep_Run_FailsNamingShadowingExecutable verifies decision 2
// (structure/qmd-run-fails-on-shadowing-executable): when a qmd this step did
// not install still wins on PATH before ensureUserNpmPrefix mutates it, Run
// still installs and registers its own copy into the canonical bin dir, then
// fails naming the exact foreign absolute path — not the canonical one.
func TestQmdMCPStep_Run_FailsNamingShadowingExecutable(t *testing.T) {
	isolateMCPFind(t, true) // already registered → registration branch is not what's under test
	stubGitHubReleaseServer(t, qmdReleaseHandler("v2.8.3-mate.6", "https://example.com/qmd.tgz"))

	npmDir := tempDir(t)
	writeBin(t, npmDir, "npm", 0, "/tmp/npm-global")
	foreignDir := tempDir(t)
	writeBin(t, foreignDir, "qmd", 0, "") // the shadowing executable, ahead in PATH before mutation
	canonicalBinDir := filepath.Join("/tmp/npm-global", "bin")
	if err := os.MkdirAll(canonicalBinDir, 0o755); err != nil {
		t.Fatalf("mkdir canonical bin dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll("/tmp/npm-global") })
	writeBin(t, canonicalBinDir, "qmd", 0, "") // what "npm install" is standing in for
	// canonicalBinDir is deliberately absent from the PATH set here: Run's own
	// ensureUserNpmPrefix call is what prepends it, mirroring production.
	isolatedPATH(t, foreignDir, npmDir)

	var out bytes.Buffer
	opts := install.Options{Stdout: &out, Stderr: &out, Yes: true}
	step := findQmdStep(t, opts)

	err := step.Run()
	if err == nil {
		t.Fatal("expected error when a foreign qmd shadows the canonical install, got nil")
	}
	wantPath := filepath.Join(foreignDir, "qmd")
	if !strings.Contains(err.Error(), wantPath) {
		t.Errorf("error should name the foreign absolute path %q; got: %v", wantPath, err)
	}
	if strings.Contains(err.Error(), canonicalBinDir) {
		t.Errorf("error should not name the canonical path %q; got: %v", canonicalBinDir, err)
	}
}

// --- Regression: a legitimate install does not start reporting pending forever ---

// TestQmdMCPStep_Check_SecondRunIdempotent mirrors the spec's "segunda corrida
// idempotente" scenario: on a machine where the canonical qmd is the only qmd
// on PATH and the registration is present, Check reports not-pending, and
// calling it again changes nothing — the guard against the widened Check
// itself becoming a source of permanent false pendings.
func TestQmdMCPStep_Check_SecondRunIdempotent(t *testing.T) {
	isolateMCPFind(t, true)

	npmDir := tempDir(t)
	writeBin(t, npmDir, "npm", 0, "/tmp/npm-global")
	canonicalBinDir := filepath.Join("/tmp/npm-global", "bin")
	if err := os.MkdirAll(canonicalBinDir, 0o755); err != nil {
		t.Fatalf("mkdir canonical bin dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll("/tmp/npm-global") })
	writeBin(t, canonicalBinDir, "qmd", 0, "")
	isolatedPATH(t, canonicalBinDir, npmDir)

	opts := install.Options{Yes: true}
	step := findQmdStep(t, opts)
	if step.Check() {
		t.Fatal("expected Check=false on the first run over a canonical-only machine")
	}
	if step.Check() {
		t.Error("expected Check=false again on the second run — nothing should have changed")
	}
}

// TestQmdMCPStep_Check_NotPendingRightAfterRun mirrors the spec's "lo que este
// paso acaba de instalar no vuelve a reportarse pendiente" scenario: right
// after a successful Run, Check reports not-pending — what this step just
// installed and registered counts as installed by its own Check.
func TestQmdMCPStep_Check_NotPendingRightAfterRun(t *testing.T) {
	isolateMCPFind(t, true) // already registered → Run's registration branch is a no-op
	stubGitHubReleaseServer(t, qmdReleaseHandler("v2.8.3-mate.6", "https://example.com/qmd.tgz"))

	npmDir := tempDir(t)
	writeBin(t, npmDir, "npm", 0, "/tmp/npm-global")
	canonicalBinDir := filepath.Join("/tmp/npm-global", "bin")
	if err := os.MkdirAll(canonicalBinDir, 0o755); err != nil {
		t.Fatalf("mkdir canonical bin dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll("/tmp/npm-global") })
	writeBin(t, canonicalBinDir, "qmd", 0, "") // what "npm install" is standing in for
	isolatedPATH(t, canonicalBinDir, npmDir)   // no foreign qmd anywhere on PATH

	var out bytes.Buffer
	opts := install.Options{Stdout: &out, Stderr: &out, Yes: true}
	step := findQmdStep(t, opts)

	if err := step.Run(); err != nil {
		t.Fatalf("expected Run to succeed on a clean machine, got: %v", err)
	}
	if step.Check() {
		t.Error("expected Check=false immediately after a successful Run, got true (pending)")
	}
}

// TestQmdMCPStep_Run_GitHubAPINon200 verifies Run fails when the GitHub API
// responds with a non-200 status while resolving the release.
func TestQmdMCPStep_Run_GitHubAPINon200(t *testing.T) {
	stubGitHubReleaseServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	d := tempDir(t)
	writeBin(t, d, "npm", 0, "/tmp/npm-global")
	isolatedPATH(t, d)

	var out bytes.Buffer
	opts := install.Options{Stdout: &out, Stderr: &out, Yes: true}
	step := findQmdStep(t, opts)

	err := step.Run()
	if err == nil {
		t.Fatal("expected error on a non-200 GitHub API response, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should mention the status code; got: %v", err)
	}
}
