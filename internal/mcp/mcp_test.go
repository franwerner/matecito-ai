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

// TestFindJSONFirst verifies that Find returns from the JSON source without
// invoking the CLI runner when the name is present in ~/.claude.json mcpServers.
func TestFindJSONFirst(t *testing.T) {
	home := t.TempDir()
	writeClaudeJSON(t, home, map[string]any{
		"context7": map[string]any{"type": "stdio"},
	})
	withHome(t, home)

	calls := stubRunner(t, "") // runner must not be called

	f, ok := mcp.Find("context7")
	if !ok {
		t.Fatal("Find: expected found=true")
	}
	if f.Source != "json" {
		t.Fatalf("Find: expected source=json, got %q", f.Source)
	}
	if *calls != 0 {
		t.Fatalf("Find: CLI runner invoked %d time(s), expected 0", *calls)
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
