package sync

import (
	"fmt"

	"github.com/franwerner/matecito-ai/internal/setup/deploy"
)

// PayloadPlanLines renders the deploy summary breakdown — the single
// formatter shared by the three independent "Plan:" previews (install's own
// preview, the sync engine's own preview, and the TUI's confirmation screen)
// so the payload breakdown never drifts between them. Indentation matches the
// two-space-per-level convention each of those three callers already uses for
// their own "verb — component" lines.
func PayloadPlanLines(s deploy.Summary) []string {
	lines := []string{fmt.Sprintf("     archivos: %d nuevos, %d cambiados, %d a borrar, %d iguales", s.New, s.Changed, s.Removed, s.Same)}
	for _, r := range s.Removals {
		lines = append(lines, fmt.Sprintf("       - borrar: %s", r))
	}
	return lines
}
