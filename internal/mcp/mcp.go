package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type Found struct {
	Name      string
	Connected bool
	Source    string
}

func Find(needle string) (Found, bool) {
	if f, ok := findViaCLI(needle); ok {
		return f, true
	}
	if f, ok := findInJSON(needle); ok {
		return f, true
	}
	return Found{}, false
}

func (f Found) Describe() string {
	switch {
	case f.Source == "cli" && f.Connected:
		return fmt.Sprintf("%q (conectado)", f.Name)
	case f.Source == "cli":
		return fmt.Sprintf("%q (registrado, no conectado)", f.Name)
	case f.Source == "json":
		return fmt.Sprintf("%q en ~/.claude.json", f.Name)
	default:
		return fmt.Sprintf("%q", f.Name)
	}
}

func findViaCLI(needle string) (Found, bool) {
	if _, err := exec.LookPath("claude"); err != nil {
		return Found{}, false
	}
	out, err := exec.Command("claude", "mcp", "list").CombinedOutput()
	if err != nil {
		return Found{}, false
	}
	lo := strings.ToLower(needle)
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(strings.ToLower(line), lo) {
			continue
		}
		name := line
		if idx := strings.Index(line, ": "); idx > 0 {
			name = line[:idx]
		}
		connected := strings.Contains(line, "✓ Connected")
		return Found{Name: name, Connected: connected, Source: "cli"}, true
	}
	return Found{}, false
}

func ListAll() []string {
	set := map[string]struct{}{}
	if _, err := exec.LookPath("claude"); err == nil {
		if out, err := exec.Command("claude", "mcp", "list").CombinedOutput(); err == nil {
			for _, line := range strings.Split(string(out), "\n") {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				idx := strings.Index(line, ": ")
				if idx <= 0 {
					continue
				}
				set[line[:idx]] = struct{}{}
			}
		}
	}
	if home, err := os.UserHomeDir(); err == nil {
		if data, err := os.ReadFile(filepath.Join(home, ".claude.json")); err == nil {
			var doc struct {
				McpServers map[string]json.RawMessage `json:"mcpServers"`
			}
			if json.Unmarshal(data, &doc) == nil {
				for name := range doc.McpServers {
					set[name] = struct{}{}
				}
			}
		}
	}
	out := make([]string, 0, len(set))
	for name := range set {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func findInJSON(needle string) (Found, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Found{}, false
	}
	data, err := os.ReadFile(filepath.Join(home, ".claude.json"))
	if err != nil {
		return Found{}, false
	}
	var doc struct {
		McpServers map[string]json.RawMessage `json:"mcpServers"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return Found{}, false
	}
	lo := strings.ToLower(needle)
	for name := range doc.McpServers {
		if strings.Contains(strings.ToLower(name), lo) {
			return Found{Name: name, Source: "json"}, true
		}
	}
	return Found{}, false
}
