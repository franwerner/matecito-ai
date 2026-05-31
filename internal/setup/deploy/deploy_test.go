package deploy_test

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/franwerner/matecito-ai/internal/setup/deploy"
)

// makePayloadFS builds a minimal fstest.MapFS that satisfies Plan:
// it must contain all Mappings roots (CLAUDE.md, agents/, skills/, references/).
func makePayloadFS(agentFiles map[string][]byte) fstest.MapFS {
	m := fstest.MapFS{
		"CLAUDE.md":             {Data: []byte("# root\n")},
		"skills/.keep":          {Data: []byte{}},
		"references/.keep":      {Data: []byte{}},
		"agents/sdd-apply.md":   {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		"agents/sdd-design.md":  {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		"agents/sdd-archive.md": {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		"agents/sdd-explore.md": {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		"agents/sdd-intake.md":  {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		"agents/sdd-propose.md": {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		"agents/sdd-spec.md":    {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		"agents/sdd-tasks.md":   {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
		"agents/sdd-verify.md":  {Data: []byte("---\nmodel: sonnet\n---\nbody\n")},
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

	ops, err := deploy.Plan(payload, claudeHome)
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
		want := payload["agents/sdd-apply.md"].Data
		if !bytes.Equal(got, want) {
			t.Errorf("agent file mismatch:\ngot:  %q\nwant: %q", got, want)
		}
	})

	// non-agent file (CLAUDE.md → matecito-ai.md) must also be byte-identical
	t.Run("non_agent_file_verbatim", func(t *testing.T) {
		got, err := os.ReadFile(filepath.Join(claudeHome, "matecito-ai.md"))
		if err != nil {
			t.Fatalf("reading matecito-ai.md: %v", err)
		}
		want := payload["CLAUDE.md"].Data
		if !bytes.Equal(got, want) {
			t.Errorf("non-agent file mismatch:\ngot:  %q\nwant: %q", got, want)
		}
	})
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

	// embedded payload must have CLAUDE.md at root
	_, statErr := fs.Stat(fsys, "CLAUDE.md")
	if statErr != nil {
		t.Errorf("CLAUDE.md not found in resolved FS (%s): %v", source, statErr)
	}
}
