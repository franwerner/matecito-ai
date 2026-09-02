package mcp_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/franwerner/matecito-ai/internal/mcp"
)

// withHome overrides HOME for the duration of the test so findInJSON reads a
// controlled ~/.claude.json. It also invalidates the CLI cache on entry and
// exit to ensure a clean slate between tests.
func withHome(t *testing.T, dir string) {
	t.Helper()
	prev := os.Getenv("HOME")
	t.Setenv("HOME", dir)
	// Force the mcp package to re-read os.UserHomeDir on the next call.
	mcp.InvalidateCLICache()
	t.Cleanup(func() {
		os.Setenv("HOME", prev)
		mcp.InvalidateCLICache()
	})
}

// writeClaudeJSON writes a minimal ~/.claude.json with the given mcpServers map
// inside dir (which becomes the test's HOME).
func writeClaudeJSON(t *testing.T, dir string, servers map[string]any) {
	t.Helper()
	doc := map[string]any{"mcpServers": servers}
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal claude.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".claude.json"), data, 0o600); err != nil {
		t.Fatalf("write .claude.json: %v", err)
	}
}

// stubRunner replaces the CLI runner with one that returns the given output and
// records how many times it was called.
func stubRunner(t *testing.T, output string) *int {
	t.Helper()
	calls := 0
	mcp.SetRunMCPList(func() ([]byte, error) {
		calls++
		return []byte(output), nil
	})
	t.Cleanup(func() { mcp.ResetRunMCPList() })
	return &calls
}

// TestFindJSONFirst verifies that Find resolves presence from the JSON source
// (never from the CLI listing) while still consulting the CLI runner once, to
// enrich the hit with the host's own connectivity — a deliberate reversal:
// before the tri-state Connection field existed, a JSON hit never invoked the
// CLI runner at all.
func TestFindJSONFirst(t *testing.T) {
	home := t.TempDir()
	writeClaudeJSON(t, home, map[string]any{
		"context7": map[string]any{"type": "stdio"},
	})
	withHome(t, home)

	calls := stubRunner(t, "context7: npx -y @upstash/context7-mcp - ✓ Connected\n")

	f, ok := mcp.Find("context7")
	if !ok {
		t.Fatal("Find: expected found=true")
	}
	if f.Source != "json" {
		t.Fatalf("Find: expected source=json, got %q", f.Source)
	}
	if *calls != 1 {
		t.Fatalf("Find: CLI runner invoked %d time(s), expected 1 (enrichment)", *calls)
	}
}

// TestFindCLIFallback verifies that when the name is absent from JSON, Find
// falls back to the CLI runner and detects the server there.
func TestFindCLIFallback(t *testing.T) {
	home := t.TempDir()
	// ~/.claude.json exists but does not contain "engram" under mcpServers.
	writeClaudeJSON(t, home, map[string]any{
		"context7": map[string]any{"type": "stdio"},
	})
	withHome(t, home)

	calls := stubRunner(t, "engram: npx -y engram-mcp\n")

	f, ok := mcp.Find("engram")
	if !ok {
		t.Fatal("Find: expected found=true via CLI fallback")
	}
	if f.Source != "cli" {
		t.Fatalf("Find: expected source=cli, got %q", f.Source)
	}
	if *calls != 1 {
		t.Fatalf("Find: CLI runner invoked %d time(s), expected 1", *calls)
	}
}

// TestFindEnrichesJSONHitConnectionUp verifies scenario 1: a registered
// integration the host launches successfully is found (from JSON, for
// presence) and reports ConnectionUp (from the CLI listing, for connectivity).
func TestFindEnrichesJSONHitConnectionUp(t *testing.T) {
	home := t.TempDir()
	writeClaudeJSON(t, home, map[string]any{
		"context7": map[string]any{"type": "stdio"},
	})
	withHome(t, home)
	stubRunner(t, "context7: npx -y @upstash/context7-mcp - ✓ Connected\n")

	f, ok := mcp.Find("context7")
	if !ok {
		t.Fatal("Find: expected found=true")
	}
	if f.Connection != mcp.ConnectionUp {
		t.Fatalf("Find: expected Connection=ConnectionUp, got %v", f.Connection)
	}
}

// TestFindEnrichesJSONHitConnectionDown verifies scenario 2: a registered
// integration the host fails to launch is still found (presence unaffected)
// and reports ConnectionDown, not simply absent.
func TestFindEnrichesJSONHitConnectionDown(t *testing.T) {
	home := t.TempDir()
	writeClaudeJSON(t, home, map[string]any{
		"qmd": map[string]any{"type": "stdio"},
	})
	withHome(t, home)
	stubRunner(t, "qmd: qmd mcp - ✗ Failed to connect\n")

	f, ok := mcp.Find("qmd")
	if !ok {
		t.Fatal("Find: expected found=true")
	}
	if f.Connection != mcp.ConnectionDown {
		t.Fatalf("Find: expected Connection=ConnectionDown, got %v", f.Connection)
	}
}

// TestFindUnregisteredNotConfusedWithDown verifies scenario 3: an integration
// absent from both the JSON registration and the CLI listing is reported as
// not found — never as a hit with a negative connection state.
func TestFindUnregisteredNotConfusedWithDown(t *testing.T) {
	home := t.TempDir()
	writeClaudeJSON(t, home, map[string]any{
		"context7": map[string]any{"type": "stdio"},
	})
	withHome(t, home)
	stubRunner(t, "context7: npx -y @upstash/context7-mcp - ✓ Connected\n")

	_, ok := mcp.Find("never-registered")
	if ok {
		t.Fatal("Find: expected found=false for an unregistered integration")
	}
}

// TestFindOtherPresenceUnchangedByConnectivity verifies scenario 4: enriching
// one integration's connectivity does not change whether every other
// registered (or unregistered) integration is found.
func TestFindOtherPresenceUnchangedByConnectivity(t *testing.T) {
	home := t.TempDir()
	writeClaudeJSON(t, home, map[string]any{
		"context7": map[string]any{"type": "stdio"},
		"qmd":      map[string]any{"type": "stdio"},
	})
	withHome(t, home)
	stubRunner(t, "context7: npx -y @upstash/context7-mcp - ✓ Connected\nqmd: qmd mcp - ✗ Failed to connect\n")

	if _, ok := mcp.Find("context7"); !ok {
		t.Error("Find(context7): expected found=true")
	}
	if _, ok := mcp.Find("qmd"); !ok {
		t.Error("Find(qmd): expected found=true")
	}
	if _, ok := mcp.Find("not-registered-anywhere"); ok {
		t.Error("Find(not-registered-anywhere): expected found=false")
	}
}

// TestCLIRunnerCalledOnce verifies that the CLI runner is invoked at most once
// even when Find is called multiple times for names absent from JSON.
func TestCLIRunnerCalledOnce(t *testing.T) {
	home := t.TempDir()
	// No .claude.json so every lookup falls back to CLI.
	withHome(t, home)

	calls := stubRunner(t, "alpha: cmd-a\nbeta: cmd-b\ngamma: cmd-c\n")

	for _, name := range []string{"alpha", "beta", "gamma"} {
		if _, ok := mcp.Find(name); !ok {
			t.Fatalf("Find(%q): expected found", name)
		}
	}
	if *calls != 1 {
		t.Fatalf("CLI runner invoked %d time(s), expected exactly 1", *calls)
	}
}

// TestInvalidateForcesRerun verifies that InvalidateCLICache causes the runner
// to be invoked again on the next Find call.
func TestInvalidateForcesRerun(t *testing.T) {
	home := t.TempDir()
	withHome(t, home)

	calls := stubRunner(t, "alpha: cmd-a\n")

	mcp.Find("alpha") // populates cache
	mcp.InvalidateCLICache()
	mcp.Find("alpha") // must re-invoke runner

	if *calls != 2 {
		t.Fatalf("CLI runner invoked %d time(s) after invalidate, expected 2", *calls)
	}
}

// TestConcurrentAccess checks that concurrent calls to Find do not race.
// Run with: go test -race ./internal/mcp/...
func TestConcurrentAccess(t *testing.T) {
	home := t.TempDir()
	writeClaudeJSON(t, home, map[string]any{
		"context7": map[string]any{"type": "stdio"},
	})
	withHome(t, home)
	stubRunner(t, "engram: npx engram-mcp\n")

	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			mcp.Find("context7")
			mcp.Find("engram")
		}()
	}
	wg.Wait()
}
