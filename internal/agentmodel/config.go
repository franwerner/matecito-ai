package agentmodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	// ConfigRelPath is the config file path relative to a repo root.
	ConfigRelPath = ".matecito-ai/config.json"

	legacyModelsFileName = "models.json"
)

// Scope identifies the active configuration scope in the TUI.
type Scope int

const (
	ScopeGlobal Scope = iota
	ScopeProject
)

// ConfigPathForScope returns the config.json path for the given scope.
// ScopeGlobal: ~/.matecito-ai/config.json (may error if HOME unavailable).
// ScopeProject: <repoRoot>/.matecito-ai/config.json (no error possible).
func ConfigPathForScope(scope Scope, repoRoot string) (string, error) {
	if scope == ScopeProject {
		return filepath.Join(repoRoot, ConfigRelPath), nil
	}
	return ConfigPath()
}

// Config holds the persisted configuration for matecito-ai.
// StrictTdd and FlagDecisionGaps use pointers so that key-absent (nil) is
// distinct from false — necessary for the per-project-vs-global precedence
// resolution (spec R5.6/S8.5).
type Config struct {
	Models           map[string]string `json:"models,omitempty"`
	StrictTdd        *bool             `json:"strictTdd,omitempty"`
	FlagDecisionGaps *bool             `json:"flagDecisionGaps,omitempty"`
}

// ConfigPath returns the global config file path: ~/.matecito-ai/config.json.
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("agentmodel: cannot determine home dir: %w", err)
	}
	return filepath.Join(home, ".matecito-ai", "config.json"), nil
}

// ProjectConfigPath returns the per-project config path relative to repoRoot.
func ProjectConfigPath(repoRoot string) string {
	return filepath.Join(repoRoot, ConfigRelPath)
}

// Load reads a Config from path.
// ENOENT → (nil, nil). Bad JSON → (nil, error).
// When path is the global config path and config.json is absent but models.json
// exists in the same directory, the one-time migration runs automatically.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// attempt migration only when the sibling models.json exists
			dir := filepath.Dir(path)
			modelsPath := filepath.Join(dir, legacyModelsFileName)
			return migrate(path, modelsPath)
		}
		return nil, fmt.Errorf("agentmodel: read %s: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("agentmodel: parse %s: %w", path, err)
	}
	return &cfg, nil
}

// migrate reads models.json, builds a Config with explicit strictTdd:false (spec S5.7),
// writes config.json, removes models.json, and returns the Config.
// Returns (nil, nil) when models.json also does not exist.
func migrate(configPath, modelsPath string) (*Config, error) {
	modelsData, err := os.ReadFile(modelsPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("agentmodel: read models.json for migration: %w", err)
	}

	var models map[string]string
	if err := json.Unmarshal(modelsData, &models); err != nil {
		return nil, fmt.Errorf("agentmodel: parse models.json: %w", err)
	}

	fals := false
	cfg := &Config{
		Models:    models,
		StrictTdd: &fals,
	}

	if err := Save(configPath, cfg); err != nil {
		return nil, fmt.Errorf("agentmodel: write config.json during migration: %w", err)
	}
	if err := os.Remove(modelsPath); err != nil {
		return nil, fmt.Errorf("agentmodel: remove models.json after migration: %w", err)
	}
	return cfg, nil
}

// Validate checks that all agent keys are in the canonical list and all model
// values are in ValidModels. StrictTdd is unconstrained (any *bool is valid).
func Validate(cfg *Config) error {
	if cfg == nil {
		return nil
	}
	for key, val := range cfg.Models {
		if !IsValidAgent(key) {
			return fmt.Errorf("agentmodel: unknown agent %q in config", key)
		}
		if !IsValidModel(val) {
			return fmt.Errorf("agentmodel: invalid model %q for agent %q (valid: %v)", val, key, ValidModels)
		}
	}
	return nil
}

// Save writes cfg to path atomically: MkdirAll the directory, write to <path>.tmp,
// then os.Rename. File mode 0o644.
func Save(path string, cfg *Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("agentmodel: mkdir %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("agentmodel: marshal config: %w", err)
	}
	data = append(data, '\n')

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("agentmodel: write tmp %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("agentmodel: rename %s → %s: %w", tmp, path, err)
	}
	return nil
}

// Defaults iterates agents/sdd-*.md files from payloadFS, calls ReadModel on each,
// and returns a map of agent-name → model. Files with no model: line are skipped.
func Defaults(payloadFS fs.FS) (map[string]string, error) {
	result := make(map[string]string)

	entries, err := fs.Glob(payloadFS, "agents/sdd-*.md")
	if err != nil {
		return nil, fmt.Errorf("agentmodel: glob agents: %w", err)
	}

	for _, entry := range entries {
		data, err := fs.ReadFile(payloadFS, entry)
		if err != nil {
			return nil, fmt.Errorf("agentmodel: read %s: %w", entry, err)
		}
		model, err := ReadModel(data)
		if err != nil {
			return nil, fmt.Errorf("agentmodel: ReadModel %s: %w", entry, err)
		}
		if model == "" {
			continue
		}
		// derive agent name from filename: "agents/sdd-apply.md" → "sdd-apply"
		base := filepath.Base(entry)
		name := strings.TrimSuffix(base, ".md")
		result[name] = model
	}
	return result, nil
}
