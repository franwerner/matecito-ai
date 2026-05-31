package permissions

import (
	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/settings"
)

func All() []check.Result {
	doc, err := settings.Load()
	if err != nil {
		return []check.Result{{
			Name:     "settings.json",
			Required: false,
			Status:   check.StatusMissing,
			Detail:   "no se pudo leer ~/.claude/settings.json",
			FixHint:  "Corré `matecito-ai install` para configurar permissions.allow",
		}}
	}

	allow := settings.AllowList(doc)
	have := make(map[string]struct{}, len(allow))
	for _, a := range allow {
		have[a] = struct{}{}
	}

	results := make([]check.Result, 0, len(settings.EcosystemPatterns))
	for _, p := range settings.EcosystemPatterns {
		r := check.Result{Name: p, Required: false}
		if _, ok := have[p]; ok {
			r.Status = check.StatusOK
			r.Detail = "en permissions.allow"
		} else {
			r.Status = check.StatusMissing
			r.Detail = "falta en permissions.allow"
			r.FixHint = "Corré `matecito-ai install` para agregarlo"
		}
		results = append(results, r)
	}
	return results
}
