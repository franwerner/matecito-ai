package sync

import (
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/setup/deploy"
)

// TestPayloadPlanLines verifies the single formatter shared by the three
// independent "Plan:" previews: one breakdown line with all four category
// counts, plus one "- borrar: <path>" line per removal, in order.
func TestPayloadPlanLines(t *testing.T) {
	s := deploy.Summary{
		New: 2, Changed: 1, Same: 5, Removed: 2,
		Removals: []string{"agents/gone-one.md", "skills/gone-two/SKILL.md"},
	}

	lines := PayloadPlanLines(s)

	if len(lines) != 3 {
		t.Fatalf("expected 1 breakdown line + 2 removal lines, got %d: %v", len(lines), lines)
	}
	if !strings.Contains(lines[0], "2 nuevos, 1 cambiados, 2 a borrar, 5 iguales") {
		t.Errorf("unexpected breakdown line: %q", lines[0])
	}
	for i, want := range s.Removals {
		if !strings.Contains(lines[i+1], want) {
			t.Errorf("removal line %d = %q, want to contain %q", i+1, lines[i+1], want)
		}
	}
}

// TestPayloadPlanLines_NoRemovals verifies the removal lines are omitted
// entirely (not even an empty header) when nothing is pending deletion.
func TestPayloadPlanLines_NoRemovals(t *testing.T) {
	lines := PayloadPlanLines(deploy.Summary{New: 1, Same: 3})
	if len(lines) != 1 {
		t.Fatalf("expected exactly the breakdown line with no removal lines, got %v", lines)
	}
	if !strings.Contains(lines[0], "1 nuevos, 0 cambiados, 0 a borrar, 3 iguales") {
		t.Errorf("unexpected breakdown line: %q", lines[0])
	}
}
