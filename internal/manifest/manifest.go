// Package manifest loads the per-domain area contract (domains/<id>/manifest.json)
// from the payload. The manifest is the machine-readable half of the two-consumer
// contract: the Go side (deploy, MCP registration, check gating) reads it, while
// the agent reads the domain's CLAUDE.md fragment.
package manifest

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
	"github.com/franwerner/matecito-ai/internal/setup/deploy"
)

// knownEvents is the set of Claude Code hook event names as of 2026-06.
// Events not in this set emit a warning at install time but are not rejected
// so that manifests remain forward-compatible with new events.
var knownEvents = map[string]bool{
	"SessionStart": true, "Setup": true, "UserPromptSubmit": true,
	"UserPromptExpansion": true, "PreToolUse": true, "PermissionRequest": true,
	"PermissionDenied": true, "PostToolUse": true, "PostToolUseFailure": true,
	"PostToolBatch": true, "Notification": true, "MessageDisplay": true,
	"SubagentStart": true, "SubagentStop": true, "TaskCreated": true,
	"TaskCompleted": true, "Stop": true, "StopFailure": true, "TeammateIdle": true,
	"InstructionsLoaded": true, "ConfigChange": true, "CwdChanged": true,
	"FileChanged": true, "WorktreeCreate": true, "WorktreeRemove": true,
	"PreCompact": true, "PostCompact": true, "Elicitation": true,
	"ElicitationResult": true, "SessionEnd": true,
}

// HookSpec is the schema of a co-located hook.json file. It declares a single
// hook handler for the enclosing domain. Command is the exact command string
// that Claude Code will invoke (e.g. "matecito-ai hook git-commit-validate").
// Id is optional; when absent it is derived as "<domainId>/<hookFolderName>".
type HookSpec struct {
	Event   string `json:"event"`
	Type    string `json:"type"`
	Command string `json:"command"`
	Matcher string `json:"matcher,omitempty"`
	If      string `json:"if,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
	Id      string `json:"id,omitempty"`
}

// ResolvedHook is a parsed hook.json with Command taken AS-IS (no path
// resolution). Install and verify consume this. Id carries the resolved
// identity: explicit from hook.json when set, derived as
// "<domainId>/<hookFolderName>" otherwise.
type ResolvedHook struct {
	Event   string
	Matcher string
	If      string
	Type    string
	Timeout int
	Command string // exact command string from hook.json, taken AS-IS (not a path)
	Id      string // identity for reconciliation; never empty after resolution
}

// DecisionRecord names the domain's decision-record type and where records live
// (ADR under .matecito-ai/adr for development; DDR under .matecito-ai/ddr for design).
type DecisionRecord struct {
	Term string `json:"term"`
	Dir  string `json:"dir"`
}

// ConfigField declares one user-configurable setting for the domain. The TUI
// renders a widget per Type and persists the value under the config's
// domainConfig[<id>] entry. Types: "bool", "enum", "agent-models".
type ConfigField struct {
	Key     string   `json:"key"`
	Type    string   `json:"type"`
	Label   string   `json:"label"`
	Default any      `json:"default,omitempty"`
	Options []string `json:"options,omitempty"` // for type "enum"
}

// Manifest is the area contract a domain implements to plug into the kernel.
type Manifest struct {
	ID                string         `json:"id"`
	Label             string         `json:"label"`
	Summary           string         `json:"summary"`
	Workspace         string         `json:"workspace"`
	AlignmentArtifact string         `json:"alignmentArtifact"`
	DecisionRecord    DecisionRecord `json:"decisionRecord"`
	CanonicalCatalog  string         `json:"canonicalCatalog"`
	Phases            []string       `json:"phases"`
	Guards            []string       `json:"guards"`
	ExplorationTool   string         `json:"explorationTool,omitempty"`
	MCP      []string      `json:"mcp"`
	Binaries []string      `json:"binaries"`
	Config   []ConfigField `json:"config,omitempty"`
}

// RelPath is the manifest location within the payload for a domain id.
func RelPath(id string) string {
	return path.Join("domains", id, "manifest.json")
}

// Load reads and parses a single domain manifest from payloadFS.
func Load(payloadFS fs.FS, id string) (*Manifest, error) {
	rel := RelPath(id)
	data, err := fs.ReadFile(payloadFS, rel)
	if err != nil {
		return nil, fmt.Errorf("manifest: read %s: %w", rel, err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("manifest: parse %s: %w", rel, err)
	}
	if m.ID == "" {
		m.ID = id
	}
	return &m, nil
}

// DomainAgents lists a domain's phase-agent names (the basenames of its
// domains/<id>/agents/*.md files), sorted. Returns nil when the domain ships no
// agents/ directory (a valid case — a domain may omit the component).
func DomainAgents(payloadFS fs.FS, id string) []string {
	dir := path.Join("domains", id, "agents")
	entries, err := fs.ReadDir(payloadFS, dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		names = append(names, strings.TrimSuffix(e.Name(), ".md"))
	}
	sort.Strings(names)
	return names
}

// DiscoverIDs lists every domain under domains/ that ships a manifest.json,
// in sorted (deterministic) order.
func DiscoverIDs(payloadFS fs.FS) ([]string, error) {
	entries, err := fs.ReadDir(payloadFS, "domains")
	if err != nil {
		return nil, fmt.Errorf("manifest: read domains/: %w", err)
	}
	var ids []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := fs.Stat(payloadFS, RelPath(e.Name())); err == nil {
			ids = append(ids, e.Name())
		}
	}
	sort.Strings(ids)
	return ids, nil
}

// ActiveIDs resolves which domains are active: the configured list when non-empty,
// otherwise every domain present in the payload (compat shim — preserves the
// development-only default). Configured ids without a manifest are dropped.
func ActiveIDs(configured []string, payloadFS fs.FS) ([]string, error) {
	discovered, err := DiscoverIDs(payloadFS)
	if err != nil {
		return nil, err
	}
	if len(configured) == 0 {
		return discovered, nil
	}
	present := make(map[string]bool, len(discovered))
	for _, id := range discovered {
		present[id] = true
	}
	var active []string
	for _, id := range configured {
		if present[id] {
			active = append(active, id)
		}
	}
	return active, nil
}

func loadIDs(payloadFS fs.FS, ids []string) ([]Manifest, error) {
	manifests := make([]Manifest, 0, len(ids))
	for _, id := range ids {
		m, err := Load(payloadFS, id)
		if err != nil {
			return nil, err
		}
		manifests = append(manifests, *m)
	}
	return manifests, nil
}

// LoadActive loads the manifests of every active domain.
func LoadActive(configured []string, payloadFS fs.FS) ([]Manifest, error) {
	ids, err := ActiveIDs(configured, payloadFS)
	if err != nil {
		return nil, err
	}
	return loadIDs(payloadFS, ids)
}

// MCPNames returns the deduplicated MCP server names declared across the given
// manifests, in first-seen order.
func MCPNames(manifests []Manifest) []string {
	seen := map[string]bool{}
	var names []string
	for _, m := range manifests {
		for _, name := range m.MCP {
			if !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	return names
}

// BinaryNames returns the deduplicated binary/CLI names declared across the
// given manifests, in first-seen order.
func BinaryNames(manifests []Manifest) []string {
	seen := map[string]bool{}
	var names []string
	for _, m := range manifests {
		for _, name := range m.Binaries {
			if !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	return names
}

// ActiveMCP returns the deduplicated set of MCP server names declared by the
// active domains, in first-seen order across sorted domain ids.
func ActiveMCP(configured []string, payloadFS fs.FS) ([]string, error) {
	manifests, err := LoadActive(configured, payloadFS)
	if err != nil {
		return nil, err
	}
	return MCPNames(manifests), nil
}

// ActiveBinaries returns the deduplicated set of binary names declared by the
// active domains, in first-seen order across sorted domain ids.
func ActiveBinaries(configured []string, payloadFS fs.FS) ([]string, error) {
	manifests, err := LoadActive(configured, payloadFS)
	if err != nil {
		return nil, err
	}
	return BinaryNames(manifests), nil
}

// ResolveFromEnv resolves the active domain ids from the global config's
// `domains` (compat shim when empty) against the embedded/local payload.
// It is the disk-backed entry point for install and verify.
func ResolveFromEnv() (ids []string, payloadFS fs.FS, err error) {
	payloadFS, _, err = deploy.ResolvePayloadFS()
	if err != nil {
		return nil, nil, err
	}
	var configured []string
	if p, e := agentmodel.ConfigPath(); e == nil {
		if cfg, e := agentmodel.Load(p); e == nil && cfg != nil {
			configured = cfg.Domains
		}
	}
	ids, err = ActiveIDs(configured, payloadFS)
	if err != nil {
		return nil, nil, err
	}
	return ids, payloadFS, nil
}

// ActiveIDsFromEnv resolves just the active domain ids from the environment
// (global config `domains` + payload, compat shim when empty). Returns nil on
// resolution error so callers fall back to the all-domains shim.
func ActiveIDsFromEnv() ([]string, error) {
	ids, _, err := ResolveFromEnv()
	return ids, err
}

// ActiveMCPFromEnv resolves the active MCP server names from the environment.
func ActiveMCPFromEnv() ([]string, error) {
	ids, payloadFS, err := ResolveFromEnv()
	if err != nil {
		return nil, err
	}
	manifests, err := loadIDs(payloadFS, ids)
	if err != nil {
		return nil, err
	}
	return MCPNames(manifests), nil
}

// ActiveBinariesFromEnv resolves the active binary names from the environment.
func ActiveBinariesFromEnv() ([]string, error) {
	ids, payloadFS, err := ResolveFromEnv()
	if err != nil {
		return nil, err
	}
	manifests, err := loadIDs(payloadFS, ids)
	if err != nil {
		return nil, err
	}
	return BinaryNames(manifests), nil
}

// IsDomainActive reports whether id is among the environment's active domains.
func IsDomainActive(id string) bool {
	ids, _, err := ResolveFromEnv()
	if err != nil {
		return false
	}
	for _, x := range ids {
		if x == id {
			return true
		}
	}
	return false
}

// ResolveHooksFromFS scans the given active domain ids for co-located hook.json
// files under domains/<id>/hooks/<hookName>/hook.json in payloadFS, parses each
// HookSpec, and returns a ResolvedHook with Command taken AS-IS from hook.json
// (no script-path resolution). Unknown events warn to stderr but do not block.
func ResolveHooksFromFS(ids []string, payloadFS fs.FS) ([]ResolvedHook, error) {
	var resolved []ResolvedHook
	for _, id := range ids {
		hooksDir := path.Join("domains", id, "hooks")
		entries, err := fs.ReadDir(payloadFS, hooksDir)
		if err != nil {
			// domain ships no hooks/ tree — skip silently
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			hookName := e.Name()
			specPath := path.Join(hooksDir, hookName, "hook.json")
			data, err := fs.ReadFile(payloadFS, specPath)
			if err != nil {
				// hook folder without hook.json — skip
				continue
			}
			var spec HookSpec
			if err := json.Unmarshal(data, &spec); err != nil {
				fmt.Fprintf(os.Stderr, "warning: %s: invalid hook.json: %v (skipping)\n", specPath, err)
				continue
			}
			if !knownEvents[spec.Event] {
				fmt.Fprintf(os.Stderr, "warning: hook %s/%s declares unknown event %q (continuing)\n", id, hookName, spec.Event)
			}
			// Derive the identity when hook.json omits an explicit id.
			resolvedID := spec.Id
			if resolvedID == "" {
				resolvedID = id + "/" + hookName
			}
			resolved = append(resolved, ResolvedHook{
				Event:   spec.Event,
				Matcher: spec.Matcher,
				If:      spec.If,
				Type:    spec.Type,
				Timeout: spec.Timeout,
				Command: spec.Command,
				Id:      resolvedID,
			})
		}
	}
	return resolved, nil
}

// ActiveHooksFromEnv scans each active domain's hooks/ tree in the payload FS
// for co-located hook.json files, taking each hook's Command AS-IS. Mirrors
// ActiveMCPFromEnv.
func ActiveHooksFromEnv() ([]ResolvedHook, error) {
	ids, payloadFS, err := ResolveFromEnv()
	if err != nil {
		return nil, err
	}
	return ResolveHooksFromFS(ids, payloadFS)
}
