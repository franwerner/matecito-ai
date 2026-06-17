package manifest_test

import (
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/franwerner/matecito-ai/internal/manifest"
)

func payload() fstest.MapFS {
	return fstest.MapFS{
		"domains/development/manifest.json": {Data: []byte(`{
			"id": "development",
			"label": "Development",
			"summary": "build software with SDD",
			"workspace": "repository",
			"alignmentArtifact": "spec",
			"decisionRecord": { "term": "ADR", "dir": ".matecito-ai/adr" },
			"canonicalCatalog": "design-patterns",
			"phases": ["intake", "spec", "apply", "verify", "archive"],
			"guards": ["strict-tdd", "review-workload"],
			"explorationTool": "codegraph",
			"mcp": ["codegraph"],
			"binaries": ["engram", "codegraph", "proofshot"]
		}`)},
		"domains/design/manifest.json": {Data: []byte(`{
			"id": "design",
			"label": "Design",
			"workspace": "folder",
			"alignmentArtifact": "brief",
			"decisionRecord": { "term": "DDR", "dir": ".matecito-ai/ddr" },
			"canonicalCatalog": "design-principles",
			"phases": ["intake", "brief", "produce", "verify", "archive"],
			"guards": ["visual-accessibility"],
			"mcp": ["figma", "canva"],
			"binaries": ["engram"]
		}`)},
		// a domain dir without a manifest must be ignored by discovery
		"domains/empty/.keep": {Data: []byte{}},
	}
}

func TestLoad(t *testing.T) {
	m, err := manifest.Load(payload(), "development")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if m.AlignmentArtifact != "spec" || m.DecisionRecord.Term != "ADR" || m.DecisionRecord.Dir != ".matecito-ai/adr" {
		t.Errorf("unexpected manifest: %+v", m)
	}
	if m.Summary != "build software with SDD" {
		t.Errorf("Summary = %q", m.Summary)
	}
	if !reflect.DeepEqual(m.MCP, []string{"codegraph"}) {
		t.Errorf("MCP = %v", m.MCP)
	}
}

func TestDiscoverIDs_SkipsManifestless(t *testing.T) {
	ids, err := manifest.DiscoverIDs(payload())
	if err != nil {
		t.Fatalf("DiscoverIDs: %v", err)
	}
	want := []string{"design", "development"} // sorted; "empty" has no manifest
	if !reflect.DeepEqual(ids, want) {
		t.Errorf("ids = %v, want %v", ids, want)
	}
}

func TestActiveIDs_ShimWhenUnconfigured(t *testing.T) {
	ids, err := manifest.ActiveIDs(nil, payload())
	if err != nil {
		t.Fatalf("ActiveIDs: %v", err)
	}
	want := []string{"design", "development"}
	if !reflect.DeepEqual(ids, want) {
		t.Errorf("shim ids = %v, want %v", ids, want)
	}
}

func TestActiveIDs_FiltersConfigured(t *testing.T) {
	ids, err := manifest.ActiveIDs([]string{"development", "ghost"}, payload())
	if err != nil {
		t.Fatalf("ActiveIDs: %v", err)
	}
	want := []string{"development"} // "ghost" has no manifest → dropped
	if !reflect.DeepEqual(ids, want) {
		t.Errorf("configured ids = %v, want %v", ids, want)
	}
}

func TestLoad_ConfigSchema(t *testing.T) {
	fsys := fstest.MapFS{
		"domains/dev/manifest.json": {Data: []byte(`{
			"id": "dev",
			"config": [
				{"key":"models","type":"agent-models","label":"Models per agent"},
				{"key":"strictTdd","type":"bool","label":"Strict TDD","default":false}
			]
		}`)},
	}
	m, err := manifest.Load(fsys, "dev")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(m.Config) != 2 {
		t.Fatalf("config len = %d, want 2", len(m.Config))
	}
	if m.Config[0].Key != "models" || m.Config[0].Type != "agent-models" {
		t.Errorf("field0 = %+v", m.Config[0])
	}
	if m.Config[1].Type != "bool" || m.Config[1].Label != "Strict TDD" {
		t.Errorf("field1 = %+v", m.Config[1])
	}
}

func TestDomainAgents(t *testing.T) {
	fsys := fstest.MapFS{
		"domains/dev/manifest.json":        {Data: []byte(`{"id":"dev"}`)},
		"domains/dev/agents/sdd-apply.md":  {Data: []byte("x")},
		"domains/dev/agents/sdd-verify.md": {Data: []byte("x")},
		"domains/dev/agents/notes.txt":     {Data: []byte("x")},
		"domains/noagents/manifest.json":   {Data: []byte(`{"id":"noagents"}`)},
	}
	got := manifest.DomainAgents(fsys, "dev")
	want := []string{"sdd-apply", "sdd-verify"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("agents = %v, want %v", got, want)
	}
	if manifest.DomainAgents(fsys, "noagents") != nil {
		t.Error("domain without agents/ should return nil")
	}
}

func TestActiveMCP_DedupAcrossDomains(t *testing.T) {
	names, err := manifest.ActiveMCP(nil, payload())
	if err != nil {
		t.Fatalf("ActiveMCP: %v", err)
	}
	// sorted domain order: design (figma, canva) then development (codegraph)
	want := []string{"figma", "canva", "codegraph"}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("mcp = %v, want %v", names, want)
	}
}

func TestLoad_Binaries(t *testing.T) {
	m, err := manifest.Load(payload(), "development")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !reflect.DeepEqual(m.Binaries, []string{"engram", "codegraph", "proofshot"}) {
		t.Errorf("Binaries = %v", m.Binaries)
	}
}

func TestActiveBinaries_DedupAcrossDomains(t *testing.T) {
	names, err := manifest.ActiveBinaries(nil, payload())
	if err != nil {
		t.Fatalf("ActiveBinaries: %v", err)
	}
	// sorted domain order: design (engram) then development (engram, codegraph, proofshot)
	want := []string{"engram", "codegraph", "proofshot"}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("binaries = %v, want %v", names, want)
	}
}

// payloadWithHookJSON returns a MapFS with one domain that ships a co-located
// hook.json and one that ships no hooks/ tree, to exercise ActiveHooksFromEnv
// without depending on the host environment.
func payloadWithHookJSON() fstest.MapFS {
	return fstest.MapFS{
		"domains/dev/manifest.json": {Data: []byte(`{"id":"dev"}`)},
		"domains/dev/hooks/git-commit-validator/hook.json": {Data: []byte(`{
			"event": "PreToolUse",
			"type": "command",
			"matcher": "Bash",
			"if": "Bash(git commit *)",
			"command": "matecito-ai hook git-commit-validate"
		}`)},
		"domains/nohooks/manifest.json": {Data: []byte(`{"id":"nohooks"}`)},
	}
}

// payloadWithHookJSONTimeout mirrors payloadWithHookJSON but the hook.json
// includes a timeout field, verifying the optional field parses correctly.
func payloadWithHookJSONTimeout() fstest.MapFS {
	return fstest.MapFS{
		"domains/dev/manifest.json": {Data: []byte(`{"id":"dev"}`)},
		"domains/dev/hooks/timed-hook/hook.json": {Data: []byte(`{
			"event": "PostToolUse",
			"type": "command",
			"command": "matecito-ai hook run",
			"timeout": 30
		}`)},
	}
}

// TestResolveHooksFromFS_HookJSONParsed verifies that a domain with a co-located
// hook.json produces a resolved hook with all fields populated correctly.
func TestResolveHooksFromFS_HookJSONParsed(t *testing.T) {
	hooks, err := manifest.ResolveHooksFromFS([]string{"dev"}, payloadWithHookJSON())
	if err != nil {
		t.Fatalf("ResolveHooksFromFS: %v", err)
	}
	if len(hooks) != 1 {
		t.Fatalf("hooks = %d, want 1", len(hooks))
	}
	h := hooks[0]
	if h.Event != "PreToolUse" {
		t.Errorf("Event = %q, want PreToolUse", h.Event)
	}
	if h.Matcher != "Bash" {
		t.Errorf("Matcher = %q, want Bash", h.Matcher)
	}
	if h.If != "Bash(git commit *)" {
		t.Errorf("If = %q", h.If)
	}
	if h.Type != "command" {
		t.Errorf("Type = %q, want command", h.Type)
	}
	wantCmd := "matecito-ai hook git-commit-validate"
	if h.Command != wantCmd {
		t.Errorf("Command = %q, want %q", h.Command, wantCmd)
	}
	if h.Timeout != 0 {
		t.Errorf("Timeout = %d, want 0 (no timeout in spec)", h.Timeout)
	}
}

// TestResolveHooksFromFS_TimeoutField verifies that the optional timeout field
// in hook.json is parsed and carried through.
func TestResolveHooksFromFS_TimeoutField(t *testing.T) {
	hooks, err := manifest.ResolveHooksFromFS([]string{"dev"}, payloadWithHookJSONTimeout())
	if err != nil {
		t.Fatalf("ResolveHooksFromFS: %v", err)
	}
	if len(hooks) != 1 {
		t.Fatalf("hooks = %d, want 1", len(hooks))
	}
	if hooks[0].Timeout != 30 {
		t.Errorf("Timeout = %d, want 30", hooks[0].Timeout)
	}
}

// TestResolveHooksFromFS_InactiveDomainExcluded verifies that passing only
// "dev" excludes "nohooks" even though it has a manifest.
func TestResolveHooksFromFS_InactiveDomainExcluded(t *testing.T) {
	// only pass "nohooks" — dev's hook must not appear
	hooks, err := manifest.ResolveHooksFromFS([]string{"nohooks"}, payloadWithHookJSON())
	if err != nil {
		t.Fatalf("ResolveHooksFromFS: %v", err)
	}
	if len(hooks) != 0 {
		t.Errorf("hooks = %v, want empty (nohooks ships no hooks/ tree)", hooks)
	}
}

// TestResolveHooksFromFS_EmptyWhenNoDomains verifies that an empty id list
// returns an empty slice without error.
func TestResolveHooksFromFS_EmptyWhenNoDomains(t *testing.T) {
	hooks, err := manifest.ResolveHooksFromFS(nil, payloadWithHookJSON())
	if err != nil {
		t.Fatalf("ResolveHooksFromFS: %v", err)
	}
	if len(hooks) != 0 {
		t.Errorf("hooks = %v, want empty", hooks)
	}
}

// TestResolveHooksFromFS_IdAutoDerived verifies that when hook.json contains no
// "id" field, the resolved hook's Id is derived as "<domainId>/<folderName>".
func TestResolveHooksFromFS_IdAutoDerived(t *testing.T) {
	hooks, err := manifest.ResolveHooksFromFS([]string{"dev"}, payloadWithHookJSON())
	if err != nil {
		t.Fatalf("ResolveHooksFromFS: %v", err)
	}
	if len(hooks) != 1 {
		t.Fatalf("hooks = %d, want 1", len(hooks))
	}
	wantID := "dev/git-commit-validator"
	if hooks[0].Id != wantID {
		t.Errorf("Id = %q, want %q (auto-derived)", hooks[0].Id, wantID)
	}
}

// payloadWithExplicitHookID returns a MapFS whose hook.json carries an explicit
// "id" field; ResolveHooksFromFS must honor it unchanged.
func payloadWithExplicitHookID() fstest.MapFS {
	return fstest.MapFS{
		"domains/dev/manifest.json": {Data: []byte(`{"id":"dev"}`)},
		"domains/dev/hooks/my-hook/hook.json": {Data: []byte(`{
			"event": "PostToolUse",
			"type": "command",
			"command": "matecito-ai hook run",
			"id": "custom/explicit-id"
		}`)},
	}
}

// TestResolveHooksFromFS_ExplicitIdHonored verifies that an explicit "id" in
// hook.json is used as-is and is not overridden by the auto-derive logic.
func TestResolveHooksFromFS_ExplicitIdHonored(t *testing.T) {
	hooks, err := manifest.ResolveHooksFromFS([]string{"dev"}, payloadWithExplicitHookID())
	if err != nil {
		t.Fatalf("ResolveHooksFromFS: %v", err)
	}
	if len(hooks) != 1 {
		t.Fatalf("hooks = %d, want 1", len(hooks))
	}
	wantID := "custom/explicit-id"
	if hooks[0].Id != wantID {
		t.Errorf("Id = %q, want %q (explicit from hook.json)", hooks[0].Id, wantID)
	}
}
