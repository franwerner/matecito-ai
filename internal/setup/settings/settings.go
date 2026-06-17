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

// HookEntry is the settings.json representation of one hook handler under a
// specific event/matcher group. Id is the matecito identity marker; when set
// it is written as "matecitoId" in settings.json and used for reconciliation.
type HookEntry struct {
	Event   string
	Matcher string
	Command string
	If      string
	Type    string
	Timeout int
	Id      string // identity marker; maps to "matecitoId" in settings.json
}

// HookList extracts all hook handlers from doc["hooks"] as a flat slice of
// HookEntry values. The settings.json hooks structure is:
//
//	{"hooks": {"<Event>": [{"matcher": "...", "hooks": [{"type":"command","command":"..."}]}]}}
func HookList(doc map[string]any) []HookEntry {
	hooksRaw, ok := doc["hooks"].(map[string]any)
	if !ok {
		return nil
	}
	var entries []HookEntry
	for event, groupsRaw := range hooksRaw {
		groups, ok := groupsRaw.([]any)
		if !ok {
			continue
		}
		for _, gRaw := range groups {
			g, ok := gRaw.(map[string]any)
			if !ok {
				continue
			}
			matcher, _ := g["matcher"].(string)
			handlersRaw, ok := g["hooks"].([]any)
			if !ok {
				continue
			}
			for _, hRaw := range handlersRaw {
				h, ok := hRaw.(map[string]any)
				if !ok {
					continue
				}
				cmd, _ := h["command"].(string)
				if cmd == "" {
					continue
				}
				ifVal, _ := h["if"].(string)
				timeout := 0
				if tv, ok := h["timeout"].(float64); ok {
					timeout = int(tv)
				}
				mid, _ := h["matecitoId"].(string)
				entries = append(entries, HookEntry{
					Event:   event,
					Matcher: matcher,
					Command: cmd,
					If:      ifVal,
					Timeout: timeout,
					Id:      mid,
				})
			}
		}
	}
	return entries
}


// ReconcileHooks performs identity-based reconciliation of hook handlers.
// Matecito owns every handler that carries a non-empty "matecitoId" field in
// settings.json.
//
// Algorithm:
//  1. Remove every handler with a non-empty matecitoId from all event groups.
//     Drop any group whose "hooks" array becomes empty (but keep groups that
//     still hold user handlers, i.e. handlers without matecitoId). Never touch
//     handlers without a matecitoId.
//  2. Add each expected hook (with its Id written as matecitoId) into the
//     matching event/matcher group, creating the group when none exists.
//
// Returns true when the document was modified. Idempotent: if the declared set
// is already exactly present (same matecitoId, event, matcher, command,
// if, timeout) the function returns false without touching the document.
func ReconcileHooks(doc map[string]any, expected []HookEntry) bool {
	// Build the desired state index keyed by id for quick lookup.
	type desiredEntry struct {
		HookEntry
		seen bool
	}
	desired := make(map[string]*desiredEntry, len(expected))
	for i := range expected {
		e := expected[i]
		if e.Id != "" {
			desired[e.Id] = &desiredEntry{HookEntry: e}
		}
	}

	hooksMap, _ := doc["hooks"].(map[string]any)
	if hooksMap == nil {
		hooksMap = map[string]any{}
	}

	changed := false

	// Step 1 — scan existing handlers; check if current matecito handlers match
	// desired and mark stale ones for removal.
	// We also use this pass to detect whether any change is needed at all.
	for event, groupsRaw := range hooksMap {
		groups, ok := groupsRaw.([]any)
		if !ok {
			continue
		}
		for gi, gRaw := range groups {
			g, ok := gRaw.(map[string]any)
			if !ok {
				continue
			}
			handlersRaw, ok := g["hooks"].([]any)
			if !ok {
				continue
			}
			matcher, _ := g["matcher"].(string)
			kept := make([]any, 0, len(handlersRaw))
			for _, hRaw := range handlersRaw {
				h, ok := hRaw.(map[string]any)
				if !ok {
					kept = append(kept, hRaw)
					continue
				}
				mid, _ := h["matecitoId"].(string)
				if mid == "" {
					// user handler — never touch
					kept = append(kept, hRaw)
					continue
				}
				// matecito-owned handler: check if it matches the desired state
				d, declared := desired[mid]
				if declared {
					cmd, _ := h["command"].(string)
					ifVal, _ := h["if"].(string)
					timeout := 0
					if tv, ok := h["timeout"].(float64); ok {
						timeout = int(tv)
					}
					typ, _ := h["type"].(string)
					matchesEvent := event == d.Event
					matchesMatcher := matcher == d.Matcher
					matchesCmd := cmd == d.Command
					matchesIf := ifVal == d.If
					matchesTimeout := timeout == d.Timeout
					matchesType := typ == d.Type || (typ == "command" && d.Type == "")
					if matchesEvent && matchesMatcher && matchesCmd && matchesIf && matchesTimeout && matchesType {
						// already correct — keep it and mark as seen
						d.seen = true
						kept = append(kept, hRaw)
						continue
					}
				}
				// stale or undeclared matecito handler — remove it
				changed = true
			}
			g["hooks"] = kept
			groups[gi] = g
		}
		hooksMap[event] = groups
	}

	// Drop groups whose hooks array is now empty.
	for event, groupsRaw := range hooksMap {
		groups, ok := groupsRaw.([]any)
		if !ok {
			continue
		}
		kept := make([]any, 0, len(groups))
		for _, gRaw := range groups {
			g, ok := gRaw.(map[string]any)
			if !ok {
				kept = append(kept, gRaw)
				continue
			}
			handlersRaw, _ := g["hooks"].([]any)
			if len(handlersRaw) > 0 {
				kept = append(kept, gRaw)
			}
			// empty group dropped (changed already flagged above when handlers were removed)
		}
		if len(kept) == 0 {
			delete(hooksMap, event)
		} else {
			hooksMap[event] = kept
		}
	}

	// Step 2 — add expected handlers that were not seen (new or stale-replaced).
	for _, d := range desired {
		if d.seen {
			continue
		}
		changed = true
		groupsRaw, _ := hooksMap[d.Event].([]any)
		appended := false
		for i, gRaw := range groupsRaw {
			g, ok := gRaw.(map[string]any)
			if !ok {
				continue
			}
			matcher, _ := g["matcher"].(string)
			if matcher != d.Matcher {
				continue
			}
			handlersRaw, _ := g["hooks"].([]any)
			handlersRaw = append(handlersRaw, buildHandlerObj(d.HookEntry))
			g["hooks"] = handlersRaw
			groupsRaw[i] = g
			appended = true
			break
		}
		if !appended {
			newGroup := map[string]any{
				"hooks": []any{buildHandlerObj(d.HookEntry)},
			}
			if d.Matcher != "" {
				newGroup["matcher"] = d.Matcher
			}
			groupsRaw = append(groupsRaw, newGroup)
		}
		hooksMap[d.Event] = groupsRaw
	}

	if changed {
		doc["hooks"] = hooksMap
	}
	return changed
}

// buildHandlerObj constructs the JSON object for a hook handler. Type is taken
// from the hook spec; timeout is omitted when unset so Claude Code applies its
// own default. When h.Id is non-empty, "matecitoId" is written so the handler
// is owned by matecito and can be reconciled later.
func buildHandlerObj(h HookEntry) map[string]any {
	t := h.Type
	if t == "" {
		t = "command"
	}
	obj := map[string]any{
		"type":    t,
		"command": h.Command,
	}
	if h.If != "" {
		obj["if"] = h.If
	}
	if h.Timeout > 0 {
		obj["timeout"] = h.Timeout
	}
	if h.Id != "" {
		obj["matecitoId"] = h.Id
	}
	return obj
}
