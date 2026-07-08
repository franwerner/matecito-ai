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
			"decisionRecord": { "term": "EDR", "dir": ".matecito-ai/edr" },
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
	if m.AlignmentArtifact != "spec" || m.DecisionRecord.Term != "EDR" || m.DecisionRecord.Dir != ".matecito-ai/edr" {
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
