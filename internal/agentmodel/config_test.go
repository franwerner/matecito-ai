package agentmodel_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
)

// --- Load ---

func TestLoad_ENOENT(t *testing.T) {
	// S5.1: non-existent file → zero-value Config, nil error
	cfg, err := agentmodel.Load(filepath.Join(t.TempDir(), "does_not_exist.json"))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if cfg != nil {
		t.Errorf("expected nil Config, got %+v", cfg)
	}
}

func TestLoad_BadJSON(t *testing.T) {
	// S5.2: invalid JSON → error
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte("{invalid"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := agentmodel.Load(path)
	if err == nil {
		t.Fatal("expected error for bad JSON")
	}
	if cfg != nil {
		t.Errorf("expected nil Config on error, got %+v", cfg)
	}
}

func TestLoad_ConfigPresent_IgnoresModelsJSON(t *testing.T) {
	// S5.8: config.json exists → models.json is ignored even if present
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	modelsPath := filepath.Join(dir, "models.json")

	configData := map[string]interface{}{
		"models": map[string]string{"sdd-apply": "haiku"},
	}
	writeJSON(t, configPath, configData)
	writeJSON(t, modelsPath, map[string]string{"sdd-apply": "sonnet"})

	cfg, err := agentmodel.Load(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil Config")
	}
	if cfg.Models["sdd-apply"] != "haiku" {
		t.Errorf("expected haiku from config.json, got %q", cfg.Models["sdd-apply"])
	}
	// models.json should still exist (not deleted when config.json is present)
	if _, err := os.Stat(modelsPath); os.IsNotExist(err) {
		t.Error("models.json should not be removed when config.json exists")
	}
}

func TestLoad_Migration(t *testing.T) {
	// S5.7: models.json present, config.json absent → migrate
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	modelsPath := filepath.Join(dir, "models.json")

	writeJSON(t, modelsPath, map[string]string{"sdd-apply": "sonnet"})

	cfg, err := agentmodel.Load(configPath)
	if err != nil {
		t.Fatalf("unexpected error during migration: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil Config after migration")
	}
	if cfg.Models["sdd-apply"] != "sonnet" {
		t.Errorf("expected sonnet, got %q", cfg.Models["sdd-apply"])
	}
	// strictTdd should be written as explicit false (spec S5.7)
	if cfg.StrictTdd == nil {
		t.Error("expected strictTdd to be non-nil (false) after migration")
	} else if *cfg.StrictTdd {
		t.Error("expected strictTdd=false after migration")
	}

	// config.json must have been created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config.json should have been created by migration")
	}
	// models.json must have been removed
	if _, err := os.Stat(modelsPath); !os.IsNotExist(err) {
		t.Error("models.json should have been removed after migration")
	}

	// second Load should read config.json and NOT recreate models.json
	cfg2, err := agentmodel.Load(configPath)
	if err != nil {
		t.Fatalf("second load error: %v", err)
	}
	if cfg2.Models["sdd-apply"] != "sonnet" {
		t.Errorf("second load: expected sonnet, got %q", cfg2.Models["sdd-apply"])
	}
}

// --- Validate ---

func TestValidate_UnknownAgentKey(t *testing.T) {
	// S5.3
	cfg := &agentmodel.Config{
		Models: map[string]string{"sdd-unknown": "sonnet"},
	}
	err := agentmodel.Validate(cfg)
	if err == nil {
		t.Fatal("expected error for unknown agent key")
	}
	if !strings.Contains(err.Error(), "sdd-unknown") {
		t.Errorf("error should mention the bad key, got: %v", err)
	}
}

func TestValidate_BadModelValue(t *testing.T) {
	// S5.4
	cfg := &agentmodel.Config{
		Models: map[string]string{"sdd-apply": "gpt-4"},
	}
	err := agentmodel.Validate(cfg)
	if err == nil {
		t.Fatal("expected error for invalid model value")
	}
	if !strings.Contains(err.Error(), "gpt-4") {
		t.Errorf("error should mention the bad value, got: %v", err)
	}
}

func TestValidate_Valid(t *testing.T) {
	// S5.5
	tru := true
	cfg := &agentmodel.Config{
		Models:    map[string]string{"sdd-apply": "sonnet"},
		StrictTdd: &tru,
	}
	if err := agentmodel.Validate(cfg); err != nil {
		t.Errorf("unexpected error for valid config: %v", err)
	}
}

func TestValidate_NilConfig(t *testing.T) {
	// nil config should not panic
	if err := agentmodel.Validate(nil); err != nil {
		t.Errorf("nil config should pass validation, got: %v", err)
	}
}

// --- Save ---

func TestSave_Atomic(t *testing.T) {
	// S5.6: atomic write; .tmp gone after; dir created; round-trip equal
	dir := t.TempDir()
	subDir := filepath.Join(dir, "newdir")
	path := filepath.Join(subDir, "config.json")

	tru := true
	cfg := &agentmodel.Config{
		Models:    map[string]string{"sdd-apply": "sonnet"},
		StrictTdd: &tru,
	}

	if err := agentmodel.Save(path, cfg); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// .tmp must be gone
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Error(".tmp file should not exist after successful Save")
	}

	// directory must exist
	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		t.Error("directory should have been created by Save")
	}

	// round-trip: Load should return equivalent config
	got, err := agentmodel.Load(path)
	if err != nil {
		t.Fatalf("Load after Save error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil Config after round-trip")
	}
	if !reflect.DeepEqual(got.Models, cfg.Models) {
		t.Errorf("Models mismatch: got %v, want %v", got.Models, cfg.Models)
	}
	if got.StrictTdd == nil || *got.StrictTdd != *cfg.StrictTdd {
		t.Error("StrictTdd mismatch after round-trip")
	}
}

// --- Defaults ---

func TestDefaults(t *testing.T) {
	// S4.9: fstest.MapFS with 9 agent files each with a model: line → 9 entries
	mapFS := fstest.MapFS{}
	agents := []string{
		"sdd-intake", "sdd-explore", "sdd-propose", "sdd-spec",
		"sdd-design", "sdd-tasks", "sdd-apply", "sdd-verify", "sdd-archive",
	}
	models := []string{
		"opus", "sonnet", "haiku", "opus",
		"sonnet", "haiku", "sonnet", "haiku", "opus",
	}
	for i, a := range agents {
		content := "---\nname: " + a + "\nmodel: " + models[i] + "\n---\nbody\n"
		mapFS["agents/"+a+".md"] = &fstest.MapFile{Data: []byte(content)}
	}

	defaults, err := agentmodel.Defaults(mapFS)
	if err != nil {
		t.Fatalf("Defaults() error: %v", err)
	}
	if len(defaults) != 9 {
		t.Errorf("expected 9 entries, got %d", len(defaults))
	}
	for i, a := range agents {
		if v, ok := defaults[a]; !ok {
			t.Errorf("missing key %q", a)
		} else if v != models[i] {
			t.Errorf("defaults[%q] = %q, want %q", a, v, models[i])
		}
	}
}

func TestDefaults_EmptyFS(t *testing.T) {
	// empty FS → empty map, no error
	defaults, err := agentmodel.Defaults(fstest.MapFS{})
	if err != nil {
		t.Fatalf("Defaults() on empty FS error: %v", err)
	}
	if len(defaults) != 0 {
		t.Errorf("expected empty map, got %v", defaults)
	}
}

// --- ConfigPath / ProjectConfigPath ---

func TestConfigPath(t *testing.T) {
	path, err := agentmodel.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error: %v", err)
	}
	if !strings.HasSuffix(path, "/.matecito-ai/config.json") {
		t.Errorf("ConfigPath() = %q, want suffix /.matecito-ai/config.json", path)
	}
	// must be absolute
	if !filepath.IsAbs(path) {
		t.Errorf("ConfigPath() = %q, want absolute path", path)
	}
}

func TestProjectConfigPath(t *testing.T) {
	got := agentmodel.ProjectConfigPath("/home/user/myrepo")
	want := "/home/user/myrepo/.matecito-ai/config.json"
	if got != want {
		t.Errorf("ProjectConfigPath() = %q, want %q", got, want)
	}
}

// --- ConfigPathForScope ---

func TestConfigPathForScope(t *testing.T) {
	tests := []struct {
		name       string
		scope      agentmodel.Scope
		repoRoot   string
		homeEnv    string
		wantSuffix string
		wantExact  string
		wantErr    bool
	}{
		{
			// Domain 4: scope=Global → ~/.matecito-ai/config.json
			name:      "global_scope",
			scope:     agentmodel.ScopeGlobal,
			repoRoot:  "",
			homeEnv:   "/fakehome",
			wantExact: "/fakehome/.matecito-ai/config.json",
		},
		{
			// Domain 4: scope=Project → <repoRoot>/.matecito-ai/config.json
			name:      "project_scope",
			scope:     agentmodel.ScopeProject,
			repoRoot:  "/path/to/repo",
			wantExact: "/path/to/repo/.matecito-ai/config.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.homeEnv != "" {
				t.Setenv("HOME", tc.homeEnv)
			}
			got, err := agentmodel.ConfigPathForScope(tc.scope, tc.repoRoot)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantExact != "" && got != tc.wantExact {
				t.Errorf("ConfigPathForScope() = %q, want %q", got, tc.wantExact)
			}
			if tc.wantSuffix != "" && !strings.HasSuffix(got, tc.wantSuffix) {
				t.Errorf("ConfigPathForScope() = %q, want suffix %q", got, tc.wantSuffix)
			}
		})
	}
}

// --- helper ---

func writeJSON(t *testing.T, path string, v interface{}) {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
