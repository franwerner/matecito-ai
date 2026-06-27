// Package hook is the compiled-in registry of matecito-ai hook handlers. Each
// hook co-locates its declaration (event, matcher, identity) with its
// implementation (the Run function) in a domain subpackage, which self-registers
// via init(). Install and verify read this registry instead of scanning payload
// files.
package hook

import "github.com/franwerner/matecito-ai/internal/manifest"

// Result is returned by a hook's Run function. Code matches the Claude Code hook
// exit-code contract: 0 = allow, 2 = block. Message is the text to write to
// stderr (empty when no message is needed).
type Result struct {
	Code    int
	Message string
}

// Hook declares one hook handler and its implementation. Event/Matcher/If map
// to the Claude Code settings.json hook structure; Id is the matecito identity
// marker used for reconciliation. Run receives the raw hook payload JSON from
// stdin and returns the exit code + stderr message.
type Hook struct {
	Domain     string
	Id         string
	Subcommand string
	Event      string
	Matcher    string
	If         string
	Timeout    int
	Run        func(payloadJSON []byte) Result
}

// Command returns the exact command string Claude Code invokes for this hook.
func (h Hook) Command() string {
	return "matecito-ai hook " + h.Subcommand
}

// SharedDomain is the sentinel domain value for hooks that must be active
// regardless of which domains the user has enabled. Hooks registered with
// Domain == SharedDomain are included by ForDomains for any active set.
const SharedDomain = "shared"

var registry []Hook

// Register adds a hook to the compiled-in registry, deriving Id as
// "<Domain>/<Subcommand>" when left empty. Domain subpackages call this from
// init().
func Register(h Hook) {
	if h.Id == "" {
		h.Id = h.Domain + "/" + h.Subcommand
	}
	registry = append(registry, h)
}

// Registered returns all hooks in the registry.
func Registered() []Hook {
	return registry
}

// ForDomains returns the registered hooks whose Domain is among ids or whose
// Domain is SharedDomain (always included regardless of active set).
func ForDomains(ids []string) []Hook {
	want := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		want[id] = struct{}{}
	}
	var out []Hook
	for _, h := range registry {
		if _, ok := want[h.Domain]; ok || h.Domain == SharedDomain {
			out = append(out, h)
		}
	}
	return out
}

// BySubcommand returns the registered hook for the given subcommand name.
func BySubcommand(name string) (Hook, bool) {
	for _, h := range registry {
		if h.Subcommand == name {
			return h, true
		}
	}
	return Hook{}, false
}

// ActiveHooks returns the registered hooks for the environment's active domains.
func ActiveHooks() ([]Hook, error) {
	ids, err := manifest.ActiveIDsFromEnv()
	if err != nil {
		return nil, err
	}
	return ForDomains(ids), nil
}
