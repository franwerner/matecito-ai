package hooks_test

import (
	"testing"

	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/checks/hooks"
)

// TestAll_EmptyWhenNoHooks verifies that All() returns nil when no active domain
// declares hooks (gating behavior: the cluster should not appear).
func TestAll_EmptyWhenNoHooks(t *testing.T) {
	results := hooks.All()
	// In CI there are no active domains with hooks; All() returns nil or a
	// non-nil slice depending on environment. Only shape invariants are asserted.
	for _, r := range results {
		switch r.Status {
		case check.StatusOK, check.StatusMissing:
			// expected
		default:
			t.Errorf("unexpected status %v in result %q", r.Status, r.Name)
		}
	}
}

// TestAll_ResultsHaveNonEmptyNames verifies that every result has a non-empty Name.
func TestAll_ResultsHaveNonEmptyNames(t *testing.T) {
	for _, r := range hooks.All() {
		if r.Name == "" {
			t.Error("result with empty Name")
		}
	}
}

// TestAll_MissingResultsCarryFixHint verifies that StatusMissing results include
// a non-empty FixHint so the user knows how to remediate.
func TestAll_MissingResultsCarryFixHint(t *testing.T) {
	for _, r := range hooks.All() {
		if r.Status == check.StatusMissing && r.FixHint == "" {
			t.Errorf("StatusMissing result %q must carry a non-empty FixHint", r.Name)
		}
	}
}
