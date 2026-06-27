package deploy_test

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/franwerner/matecito-ai/internal/setup/deploy"
)

// makePayloadFS builds a minimal fstest.MapFS that satisfies Plan: the
// domain-agnostic core/CLAUDE.md plus one domain (development) shipping the
// agents/, skills/ and references/ components.
func makePayloadFS(agentFiles map[string][]byte) fstest.MapFS {
	const dev = "domains/development/"
	m := fstest.MapFS{
		"core/CLAUDE.md":              {Data: []byte("# root\n")},
		dev + "skills/.keep":          {Data: []byte{}},
		dev + "references/.keep":      {Data: []byte{}},
		dev + "agents/sdd-apply.md":   {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		dev + "agents/sdd-design.md":  {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		dev + "agents/sdd-archive.md": {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		dev + "agents/sdd-explore.md": {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		dev + "agents/sdd-intake.md":  {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		dev + "agents/sdd-propose.md": {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		dev + "agents/sdd-spec.md":    {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		dev + "agents/sdd-tasks.md":   {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		dev + "agents/sdd-verify.md":  {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
	}
	for k, v := range agentFiles {
		m[k] = &fstest.MapFile{Data: v}
	}
	return m
}

func newClaudeHome(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "claude-home-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func newBackupDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "backup-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

// TestApply_Verbatim verifies that Apply copies files byte-for-byte from the
// payload without any modification — both agent files and non-agent files.
func TestApply_Verbatim(t *testing.T) {
	payload := makePayloadFS(nil)
	claudeHome := newClaudeHome(t)
	backupDir := newBackupDir(t)

	ops, err := deploy.Plan(payload, claudeHome, nil)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}

	_, err = deploy.Apply(payload, ops, claudeHome, backupDir)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	// agent file must be byte-identical to payload
	t.Run("agent_file_verbatim", func(t *testing.T) {
		got, err := os.ReadFile(filepath.Join(claudeHome, "agents", "sdd-apply.md"))
		if err != nil {
			t.Fatalf("reading written file: %v", err)
		}
		want := payload["domains/development/agents/sdd-apply.md"].Data
		if !bytes.Equal(got, want) {
			t.Errorf("agent file mismatch:\ngot:  %q\nwant: %q", got, want)
		}
	})

	// matecito-ai.md = kernel core + generated domains index (no domain bodies).
	t.Run("matecito_md_is_core_plus_index", func(t *testing.T) {
		got, err := os.ReadFile(filepath.Join(claudeHome, "matecito-ai.md"))
		if err != nil {
			t.Fatalf("reading matecito-ai.md: %v", err)
		}
		if !strings.Contains(string(got), "# root") {
			t.Errorf("matecito-ai.md missing kernel core, got:\n%s", got)
		}
		if !strings.Contains(string(got), "Active domains — load on demand") {
			t.Errorf("matecito-ai.md missing domains index, got:\n%s", got)
		}
	})
}

// TestCompose_ClaudeMd verifies matecito-ai.md = kernel core + a generated index
// of the active domains (with each domain's label, summary and fragment path),
// and that the domain CLAUDE.md *body* is deployed standalone (NOT inlined).
func TestCompose_ClaudeMd(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                    {Data: []byte("# core\n")},
		"domains/development/manifest.json": {Data: []byte(`{"id":"development","label":"Development","summary":"build software"}`)},
		"domains/development/CLAUDE.md":     {Data: []byte("# development body\n")},
	}
	claudeHome := newClaudeHome(t)
	backupDir := newBackupDir(t)

	ops, err := deploy.Plan(payload, claudeHome, nil)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if _, err := deploy.Apply(payload, ops, claudeHome, backupDir); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	idx, err := os.ReadFile(filepath.Join(claudeHome, "matecito-ai.md"))
	if err != nil {
		t.Fatalf("reading matecito-ai.md: %v", err)
	}
	for _, want := range []string{"# core", "Active domains — load on demand", "Development", "build software", "matecito-ai/domains/development.md"} {
		if !strings.Contains(string(idx), want) {
			t.Errorf("matecito-ai.md missing %q, got:\n%s", want, idx)
		}
	}
	if strings.Contains(string(idx), "# development body") {
		t.Errorf("domain body must NOT be inlined into matecito-ai.md, got:\n%s", idx)
	}

	// the domain body is deployed standalone, byte-identical
	frag, err := os.ReadFile(filepath.Join(claudeHome, "matecito-ai", "domains", "development.md"))
	if err != nil {
		t.Fatalf("reading domain fragment: %v", err)
	}
	if want := payload["domains/development/CLAUDE.md"].Data; !bytes.Equal(frag, want) {
		t.Errorf("fragment mismatch:\ngot:  %q\nwant: %q", frag, want)
	}
}

// TestPlan_SkillClashAcrossDomains verifies the deploy guard rejects two domains
// that expose a skill folder of the same name, with a domain-aware message
// (M4 Option 1: detect, do not auto-prefix).
func TestPlan_SkillClashAcrossDomains(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md": {Data: []byte("# core\n")},
		"domains/development/skills/grp/audit/SKILL.md": {Data: []byte("dev")},
		"domains/design/skills/grp/audit/SKILL.md":      {Data: []byte("design")},
	}
	claudeHome := newClaudeHome(t)

	_, err := deploy.Plan(payload, claudeHome, nil)
	if err == nil {
		t.Fatal("expected clash error, got nil")
	}
	for _, want := range []string{"development", "design", "audit"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("clash error %q missing %q", err.Error(), want)
		}
	}
}

// TestPlan_FiltersInactiveDomains verifies that an explicit active set excludes
// non-active domains from both the composed matecito-ai.md and the component ops.
func TestPlan_FiltersInactiveDomains(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                                   {Data: []byte("# core\n")},
		"domains/development/CLAUDE.md":                    {Data: []byte("# development\n")},
		"domains/development/skills/grp/devskill/SKILL.md": {Data: []byte("dev")},
		"domains/design/CLAUDE.md":                         {Data: []byte("# design\n")},
		"domains/design/skills/grp/designskill/SKILL.md":   {Data: []byte("design")},
	}
	claudeHome := newClaudeHome(t)
	backupDir := newBackupDir(t)

	ops, err := deploy.Plan(payload, claudeHome, []string{"development"})
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if _, err := deploy.Apply(payload, ops, claudeHome, backupDir); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	// matecito-ai.md index must list development but NOT design
	got, err := os.ReadFile(filepath.Join(claudeHome, "matecito-ai.md"))
	if err != nil {
		t.Fatalf("reading matecito-ai.md: %v", err)
	}
	if !strings.Contains(string(got), "domains/development.md") || strings.Contains(string(got), "domains/design.md") {
		t.Errorf("index should list development only, got:\n%s", got)
	}

	// development fragment deployed standalone, design fragment not
	if _, err := os.Stat(filepath.Join(claudeHome, "matecito-ai", "domains", "development.md")); err != nil {
		t.Errorf("expected development fragment deployed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(claudeHome, "matecito-ai", "domains", "design.md")); err == nil {
		t.Error("design fragment should NOT be deployed when design is inactive")
	}

	// development skill deployed, design skill not
	if _, err := os.Stat(filepath.Join(claudeHome, "skills", "devskill", "SKILL.md")); err != nil {
		t.Errorf("expected development skill deployed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(claudeHome, "skills", "designskill", "SKILL.md")); err == nil {
		t.Error("design skill should NOT be deployed when design is inactive")
	}
}

// TestPlan_RealPayload_IndexesBothDomains is an end-to-end proof of symmetry:
// the real payload (local or embedded) indexes development AND design in
// matecito-ai.md and deploys each body standalone, without a clash.
func TestPlan_RealPayload_IndexesBothDomains(t *testing.T) {
	payloadFS, _, err := deploy.ResolvePayloadFS()
	if err != nil {
		t.Fatalf("ResolvePayloadFS: %v", err)
	}
	claudeHome := newClaudeHome(t)
	backupDir := newBackupDir(t)

	ops, err := deploy.Plan(payloadFS, claudeHome, nil)
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if _, err := deploy.Apply(payloadFS, ops, claudeHome, backupDir); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	idx, err := os.ReadFile(filepath.Join(claudeHome, "matecito-ai.md"))
	if err != nil {
		t.Fatalf("reading matecito-ai.md: %v", err)
	}
	for _, want := range []string{"domains/development.md", "domains/design.md"} {
		if !strings.Contains(string(idx), want) {
			t.Errorf("index missing %q", want)
		}
	}

	// each domain body is deployed standalone, carrying its title
	for id, title := range map[string]string{"development": "Development domain", "design": "Design domain"} {
		frag, err := os.ReadFile(filepath.Join(claudeHome, "matecito-ai", "domains", id+".md"))
		if err != nil {
			t.Fatalf("reading %s fragment: %v", id, err)
		}
		if !strings.Contains(string(frag), title) {
			t.Errorf("%s fragment missing %q", id, title)
		}
	}
}

// TestPlan_SharedDeploys_ZeroActiveDomains verifies that shared components
// deploy even when the active-domain set is empty.
func TestPlan_SharedDeploys_ZeroActiveDomains(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                              {Data: []byte("# core\n")},
		"domains/development/agents/noop.md":          {Data: []byte("noop")},
		"shared/skills/grp/myskill/SKILL.md":          {Data: []byte("shared skill")},
	}
	claudeHome := newClaudeHome(t)

	ops, err := deploy.Plan(payload, claudeHome, []string{})
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}

	for _, op := range ops {
		if strings.HasSuffix(op.Target, filepath.Join("skills", "myskill", "SKILL.md")) {
			return
		}
	}
	t.Errorf("expected shared skill target in ops, got: %v", ops)
}

// TestPlan_SharedDeploys_AlongsideActiveDomain verifies that shared components
// deploy alongside active-domain components.
func TestPlan_SharedDeploys_AlongsideActiveDomain(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                                    {Data: []byte("# core\n")},
		"domains/development/agents/agent.md":               {Data: []byte("agent")},
		"shared/skills/grp/myskill/SKILL.md":                {Data: []byte("shared skill")},
	}
	claudeHome := newClaudeHome(t)

	ops, err := deploy.Plan(payload, claudeHome, []string{"development"})
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}

	var foundShared, foundAgent bool
	for _, op := range ops {
		if strings.HasSuffix(op.Target, filepath.Join("skills", "myskill", "SKILL.md")) {
			foundShared = true
		}
		if strings.HasSuffix(op.Target, filepath.Join("agents", "agent.md")) {
			foundAgent = true
		}
	}
	if !foundShared {
		t.Error("expected shared skill target in ops")
	}
	if !foundAgent {
		t.Error("expected development agent target in ops")
	}
}

// TestPlan_SharedSkill_Flattens verifies that a shared skill under
// shared/skills/<group>/<skill>/ is flattened to skills/<skill>/, matching
// the domain-skill flattening rule.
func TestPlan_SharedSkill_Flattens(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                     {Data: []byte("# core\n")},
		"domains/.keep":                       {Data: []byte{}},
		"shared/skills/grp/myskill/SKILL.md": {Data: []byte("shared skill")},
	}
	claudeHome := newClaudeHome(t)

	ops, err := deploy.Plan(payload, claudeHome, []string{})
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}

	want := filepath.Join(claudeHome, "skills", "myskill", "SKILL.md")
	for _, op := range ops {
		if op.Target == want {
			return
		}
	}
	t.Errorf("expected target %q, ops: %v", want, ops)
}

// TestPlan_SharedVsDomain_SameTarget_Clashes verifies that a shared component
// and a domain component resolving to the same target still produce a clashError.
func TestPlan_SharedVsDomain_SameTarget_Clashes(t *testing.T) {
	payload := fstest.MapFS{
		"core/CLAUDE.md":                                    {Data: []byte("# core\n")},
		"shared/skills/grp/audit/SKILL.md":                  {Data: []byte("shared")},
		"domains/development/skills/grp/audit/SKILL.md":     {Data: []byte("dev")},
	}
	claudeHome := newClaudeHome(t)

	_, err := deploy.Plan(payload, claudeHome, nil)
	if err == nil {
		t.Fatal("expected clash error, got nil")
	}
}

// TestResolvePayloadFS_Embedded verifies that ResolvePayloadFS falls back to the
// embedded payload when no local payload/ dir is present (R4.7).
func TestResolvePayloadFS_Embedded(t *testing.T) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp, err := os.MkdirTemp("", "no-payload-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	defer os.Chdir(orig) //nolint:errcheck
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	fsys, source, err := deploy.ResolvePayloadFS()
	if err != nil {
		t.Fatalf("ResolvePayloadFS: %v", err)
	}
	if source == "" {
		t.Error("source should not be empty")
	}

	// embedded payload must have the core CLAUDE.md
	_, statErr := fs.Stat(fsys, "core/CLAUDE.md")
	if statErr != nil {
		t.Errorf("core/CLAUDE.md not found in resolved FS (%s): %v", source, statErr)
	}
}

// TestWalkDir_GitkeepFiltered verifies that .gitkeep placeholder files are never
// emitted as FileOps, while real files in the same subtree are still included.
func TestWalkDir_GitkeepFiltered(t *testing.T) {
	// A subtree that contains only .gitkeep files — simulates empty shared/agents/
	// and shared/references/ directories as they exist in the repository.
	onlyGitkeep := fstest.MapFS{
		"core/CLAUDE.md":             {Data: []byte("# core\n")},
		"domains/.keep":              {Data: []byte{}},
		"shared/agents/.gitkeep":     {Data: []byte{}},
		"shared/references/.gitkeep": {Data: []byte{}},
		"shared/skills/.gitkeep":     {Data: []byte{}},
	}
	claudeHome := newClaudeHome(t)

	ops, err := deploy.Plan(onlyGitkeep, claudeHome, []string{})
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	for _, op := range ops {
		if filepath.Base(op.Target) == ".gitkeep" {
			t.Errorf("unexpected .gitkeep in ops: %s", op.Target)
		}
	}

	// A real file alongside a .gitkeep must still be deployed.
	withReal := fstest.MapFS{
		"core/CLAUDE.md":           {Data: []byte("# core\n")},
		"domains/.keep":            {Data: []byte{}},
		"shared/agents/.gitkeep":   {Data: []byte{}},
		"shared/agents/myagent.md": {Data: []byte("agent")},
	}
	ops2, err := deploy.Plan(withReal, claudeHome, []string{})
	if err != nil {
		t.Fatalf("Plan (withReal): %v", err)
	}
	var foundAgent bool
	for _, op := range ops2 {
		if filepath.Base(op.Target) == ".gitkeep" {
			t.Errorf("unexpected .gitkeep in ops2: %s", op.Target)
		}
		if strings.HasSuffix(op.Target, filepath.Join("agents", "myagent.md")) {
			foundAgent = true
		}
	}
	if !foundAgent {
		t.Error("expected myagent.md in ops2, not found")
	}
}
