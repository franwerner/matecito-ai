package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// ConnectionState represents whether the host was asked about a registered
// integration's connectivity, and, if it was, what it answered.
// ConnectionUnknown (the zero value) means the host was never consulted — it
// is distinct from ConnectionDown, which means the host was consulted and
// reported the integration not connected.
type ConnectionState int

const (
	ConnectionUnknown ConnectionState = iota
	ConnectionUp
	ConnectionDown
)

type Found struct {
	Name       string
	Connection ConnectionState
	Source     string
}

func defaultRunMCPList() ([]byte, error) {
	if _, err := exec.LookPath("claude"); err != nil {
		return nil, err
	}
	return exec.Command("claude", "mcp", "list").CombinedOutput()
}

// runMCPList is the indirection used to invoke "claude mcp list". Tests
// replace this var with a stub to avoid requiring a real claude binary.
var runMCPList func() ([]byte, error) = defaultRunMCPList

// cliCache holds the memoized result of a single "claude mcp list" call.
// Access is guarded by cliCacheMu / cliCacheOnce so concurrent goroutines
// (e.g. the TUI Sync goroutine) never race.
var (
	cliCacheMu   sync.Mutex
	cliCacheOnce *sync.Once
	cliCacheData []byte
	cliCacheErr  error
)

func init() {
	cliCacheOnce = &sync.Once{}
}

// InvalidateCLICache discards the cached "claude mcp list" result so the next
// Find or ListAll call re-invokes the runner. Call this after registering a new
// MCP server (e.g. after ApplyConfigSteps) to avoid stale reads.
func InvalidateCLICache() {
	cliCacheMu.Lock()
	defer cliCacheMu.Unlock()
	cliCacheOnce = &sync.Once{}
	cliCacheData = nil
	cliCacheErr = nil
}

// cachedCLIOutput returns the memoized output of runMCPList, invoking it at
// most once per command invocation (or since the last InvalidateCLICache call).
func cachedCLIOutput() ([]byte, error) {
	cliCacheMu.Lock()
	once := cliCacheOnce
	cliCacheMu.Unlock()

	once.Do(func() {
		data, err := runMCPList()
		cliCacheMu.Lock()
		cliCacheData = data
		cliCacheErr = err
		cliCacheMu.Unlock()
	})

	cliCacheMu.Lock()
	defer cliCacheMu.Unlock()
	return cliCacheData, cliCacheErr
}

func Find(needle string) (Found, bool) {
	if f, ok := findInJSON(needle); ok {
		f.Connection = connectionFromCLI(f.Name)
		return f, true
	}
	if f, ok := findViaCLI(needle); ok {
		return f, true
	}
	return Found{}, false
}

func (f Found) Describe() string {
	switch f.Connection {
	case ConnectionUp:
		return fmt.Sprintf("%q (conectado)", f.Name)
	case ConnectionDown:
		return fmt.Sprintf("%q (registrado, no conectado)", f.Name)
	default:
		if f.Source == "json" {
			return fmt.Sprintf("%q en ~/.claude.json", f.Name)
		}
		return fmt.Sprintf("%q", f.Name)
	}
}

// connectionFromCLI reports the host's own connectivity for a registered
// integration by name, read from the memoized "claude mcp list" output. It
// stays ConnectionUnknown when the CLI cannot be consulted at all (no claude
// in PATH, a failing or unparseable listing) or when the listing has no line
// for this name — the honest "not asked" answer, never a negative one by
// default.
func connectionFromCLI(name string) ConnectionState {
	out, err := cachedCLIOutput()
	if err != nil {
		return ConnectionUnknown
	}
	lo := strings.ToLower(name)
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, ": ")
		if idx <= 0 || strings.ToLower(line[:idx]) != lo {
			continue
		}
		if strings.Contains(line, "✓ Connected") {
			return ConnectionUp
		}
		return ConnectionDown
	}
	return ConnectionUnknown
}

func findViaCLI(needle string) (Found, bool) {
	out, err := cachedCLIOutput()
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
		state := ConnectionDown
		if strings.Contains(line, "✓ Connected") {
			state = ConnectionUp
		}
		return Found{Name: name, Connection: state, Source: "cli"}, true
	}
	return Found{}, false
}

func ListAll() []string {
	set := map[string]struct{}{}

	if out, err := cachedCLIOutput(); err == nil {
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
