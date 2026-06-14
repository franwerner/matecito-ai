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

// DefaultDomain is the domain that pre-M7 (flat) config keys migrate into.
const DefaultDomain = "development"

// DomainConfig holds the per-domain settings (M7): model-per-agent overrides
// (scoped to that domain's agents) and the domain's guards (e.g. strictTdd for
// development). Shared settings (domains, flagDecisionGaps) stay top-level.
type DomainConfig struct {
	Models           map[string]string `json:"models,omitempty"`
	StrictTdd        *bool             `json:"strictTdd,omitempty"`
	FlagDecisionGaps *bool             `json:"flagDecisionGaps,omitempty"`
	// Settings holds generic manifest-declared fields (enum/bool beyond the typed
	// ones above), keyed by the field key. Lets a domain add config without a Go
	// struct change.
	Settings map[string]any `json:"settings,omitempty"`
}

// Config holds the persisted configuration for matecito-ai.
// StrictTdd and FlagDecisionGaps use pointers so that key-absent (nil) is
// distinct from false — necessary for the per-project-vs-global precedence
// resolution (spec R5.6/S8.5).
type Config struct {
	// Shared (cross-domain).
	// Domains lists the area domains installed for this scope. Empty/absent
	// means "all domains present in the payload" (compat shim — keeps the
	// development-only default behaving as before the multi-domain split).
	Domains []string `json:"domains,omitempty"`
	// DomainConfig holds per-domain settings keyed by domain id (M7).
	DomainConfig map[string]*DomainConfig `json:"domainConfig,omitempty"`

	// Legacy top-level keys (pre-M7 / pre-per-domain-flag). Read for backward
	// compatibility and folded into DomainConfig[DefaultDomain] by normalize();
	// never written back.
	Models           map[string]string `json:"models,omitempty"`
	StrictTdd        *bool             `json:"strictTdd,omitempty"`
	FlagDecisionGaps *bool             `json:"flagDecisionGaps,omitempty"`
}

// normalize folds legacy top-level Models/StrictTdd into DomainConfig[DefaultDomain]
// and clears them, so consumers always see the nested (M7) shape and Save writes
// only the new form. Idempotent.
func (c *Config) normalize() {
	if len(c.Models) == 0 && c.StrictTdd == nil && c.FlagDecisionGaps == nil {
		return
	}
	dev := c.ensureDomain(DefaultDomain)
	if len(c.Models) > 0 && dev.Models == nil {
		dev.Models = c.Models
	}
	if c.StrictTdd != nil && dev.StrictTdd == nil {
		dev.StrictTdd = c.StrictTdd
	}
	if c.FlagDecisionGaps != nil && dev.FlagDecisionGaps == nil {
		dev.FlagDecisionGaps = c.FlagDecisionGaps
	}
	c.Models = nil
	c.StrictTdd = nil
	c.FlagDecisionGaps = nil
}

func (c *Config) ensureDomain(domain string) *DomainConfig {
	if c.DomainConfig == nil {
		c.DomainConfig = map[string]*DomainConfig{}
	}
	dc := c.DomainConfig[domain]
	if dc == nil {
		dc = &DomainConfig{}
		c.DomainConfig[domain] = dc
	}
	return dc
}

// DomainModels returns the model-per-agent map for domain (nil if unset).
func (c *Config) DomainModels(domain string) map[string]string {
	if c == nil || c.DomainConfig == nil {
		return nil
	}
	if dc := c.DomainConfig[domain]; dc != nil {
		return dc.Models
	}
	return nil
}

// DomainStrictTdd returns the strictTdd pointer for domain (nil if unset).
func (c *Config) DomainStrictTdd(domain string) *bool {
	if c == nil || c.DomainConfig == nil {
		return nil
	}
	if dc := c.DomainConfig[domain]; dc != nil {
		return dc.StrictTdd
	}
	return nil
}

// SetDomainModelOverride sets (or clears, when model=="") a single agent override.
func (c *Config) SetDomainModelOverride(domain, agent, model string) {
	dc := c.ensureDomain(domain)
	if dc.Models == nil {
		dc.Models = map[string]string{}
	}
	if model == "" {
		delete(dc.Models, agent)
	} else {
		dc.Models[agent] = model
	}
}

// SetDomainStrictTdd sets the strictTdd flag for domain.
func (c *Config) SetDomainStrictTdd(domain string, v *bool) {
	c.ensureDomain(domain).StrictTdd = v
}

// DomainFlagDecisionGaps returns the flagDecisionGaps pointer for domain (nil if unset).
func (c *Config) DomainFlagDecisionGaps(domain string) *bool {
	if c == nil || c.DomainConfig == nil {
		return nil
	}
	if dc := c.DomainConfig[domain]; dc != nil {
		return dc.FlagDecisionGaps
	}
	return nil
}

// SetDomainFlagDecisionGaps sets the flagDecisionGaps flag for domain.
func (c *Config) SetDomainFlagDecisionGaps(domain string, v *bool) {
	c.ensureDomain(domain).FlagDecisionGaps = v
}

// DomainSetting returns a generic per-domain manifest field value (nil if unset).
func (c *Config) DomainSetting(domain, key string) any {
	if c == nil || c.DomainConfig == nil {
		return nil
	}
	if dc := c.DomainConfig[domain]; dc != nil && dc.Settings != nil {
		return dc.Settings[key]
	}
	return nil
}

// SetDomainSetting sets a generic per-domain manifest field value.
func (c *Config) SetDomainSetting(domain, key string, value any) {
	dc := c.ensureDomain(domain)
	if dc.Settings == nil {
		dc.Settings = map[string]any{}
	}
	dc.Settings[key] = value
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
	cfg.normalize()
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
	cfg.normalize()

	if err := Save(configPath, cfg); err != nil {
		return nil, fmt.Errorf("agentmodel: write config.json during migration: %w", err)
	}
	if err := os.Remove(modelsPath); err != nil {
		return nil, fmt.Errorf("agentmodel: remove models.json after migration: %w", err)
	}
	return cfg, nil
}

// Validate checks that all model values are in ValidModels. Agent keys are NOT
// constrained to a hardcoded roster: each domain's agents are discovered from the
// payload (domains/<domain>/agents/*.md), so a config may name any agent the
// active domains ship — validating names here would reject new/renamed agents.
// StrictTdd is unconstrained (any *bool is valid).
func Validate(cfg *Config) error {
	if cfg == nil {
		return nil
	}
	cfg.normalize()
	for _, dc := range cfg.DomainConfig {
		if dc == nil {
			continue
		}
		for key, val := range dc.Models {
			if !IsValidModel(val) {
				return fmt.Errorf("agentmodel: invalid model %q for agent %q (valid: %v)", val, key, ValidModels)
			}
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

// DefaultsForDomain iterates domains/<domain>/agents/*.md from payloadFS, calls
// ReadModel on each, and returns a map of agent-name → model. Files with no
// model: line are skipped.
func DefaultsForDomain(payloadFS fs.FS, domain string) (map[string]string, error) {
	result := make(map[string]string)

	entries, err := fs.Glob(payloadFS, "domains/"+domain+"/agents/*.md")
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
		name := strings.TrimSuffix(filepath.Base(entry), ".md")
		result[name] = model
	}
	return result, nil
}
