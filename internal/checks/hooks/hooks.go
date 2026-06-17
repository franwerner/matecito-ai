package hooks

import (
	"fmt"

	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/hook"
	"github.com/franwerner/matecito-ai/internal/setup/settings"
)

// resolveHooks is the resolution seam so tests can substitute the real
// hook.ActiveHooks without depending on host environment state.
var resolveHooks = func() ([]hook.Hook, error) {
	return hook.ActiveHooks()
}

// loadSettings is the settings load seam for the same reason.
var loadSettings = func() (map[string]any, error) {
	return settings.Load()
}

// All returns one check.Result per active declared hook: whether settings.json
// contains the handler for that hook. Script presence is not checked — hooks
// are now Go subcommands invoked by name, not deployed shell scripts.
func All() []check.Result {
	hooks, err := resolveHooks()
	if err != nil {
		return []check.Result{{
			Name:    "hooks",
			Status:  check.StatusMissing,
			Detail:  fmt.Sprintf("no se pudo resolver hooks activos: %v", err),
			FixHint: "Corré `matecito-ai install` para registrar los hooks",
		}}
	}
	if len(hooks) == 0 {
		return nil
	}

	doc, docErr := loadSettings()

	results := make([]check.Result, 0, len(hooks))
	for _, h := range hooks {
		results = append(results, checkSettings(h, doc, docErr))
	}
	return results
}

// checkSettings reports whether settings.json contains a handler whose
// matecitoId matches the declared hook's Id and whose command also matches.
// Every resolved hook carries a non-empty Id after resolution.
func checkSettings(h hook.Hook, doc map[string]any, docErr error) check.Result {
	name := fmt.Sprintf("hook handler: %s/%s", h.Event, h.Command())
	r := check.Result{
		Name:     name,
		Required: false,
		FixHint:  "Corré `matecito-ai install` para registrar el hook en settings.json",
	}
	if docErr != nil {
		r.Status = check.StatusMissing
		r.Detail = "no se pudo leer ~/.claude/settings.json"
		return r
	}
	existing := settings.HookList(doc)
	found := false
	// Identity-based check: match by matecitoId and command.
	// Every resolved hook carries a non-empty Id after resolution.
	for _, e := range existing {
		if e.Id == h.Id && e.Command == h.Command() {
			found = true
			break
		}
	}
	if found {
		r.Status = check.StatusOK
		r.Detail = "en settings.json"
		return r
	}
	r.Status = check.StatusMissing
	r.Detail = "falta en settings.json"
	return r
}
