package settings

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude", "settings.json"), nil
}

func Load() (map[string]any, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]any{}, nil
	}
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(string(data)) == "" {
		return map[string]any{}, nil
	}
	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	if doc == nil {
		doc = map[string]any{}
	}
	return doc, nil
}

func AllowList(doc map[string]any) []string {
	perms, ok := doc["permissions"].(map[string]any)
	if !ok {
		return nil
	}
	raw, ok := perms["allow"].([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// MissingPatterns returns the expected patterns not present in allow. The
// caller supplies expected (derived from the active domains); this package
// stays MCP-agnostic.
func MissingPatterns(allow, expected []string) []string {
	have := make(map[string]struct{}, len(allow))
	for _, a := range allow {
		have[a] = struct{}{}
	}
	var missing []string
	for _, p := range expected {
		if _, ok := have[p]; !ok {
			missing = append(missing, p)
		}
	}
	return missing
}

// Merge adds the expected patterns missing from doc's permissions.allow. It only
// adds; it never removes patterns already present.
func Merge(doc map[string]any, expected []string) bool {
	missing := MissingPatterns(AllowList(doc), expected)
	if len(missing) == 0 {
		return false
	}
	perms, ok := doc["permissions"].(map[string]any)
	if !ok {
		perms = map[string]any{}
		doc["permissions"] = perms
	}
	allow, _ := perms["allow"].([]any)
	for _, p := range missing {
		allow = append(allow, p)
	}
	perms["allow"] = allow
	return true
}
