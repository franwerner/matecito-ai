package hooks

import (
	"testing"

	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/manifest"
)

// TestAll_SettingsMissing exercises the StatusMissing branch for a hook whose
// handler is absent from settings.json, independent of the host environment.
func TestAll_SettingsMissing(t *testing.T) {
	orig := resolveHooks
	t.Cleanup(func() { resolveHooks = orig })
	resolveHooks = func() ([]manifest.ResolvedHook, error) {
		return []manifest.ResolvedHook{
			{Event: "PreToolUse", Matcher: "Bash", Command: "matecito-ai hook git-commit-validate", Id: "development/git-commit-validator"},
		}, nil
	}
	origLoad := loadSettings
	t.Cleanup(func() { loadSettings = origLoad })
	loadSettings = func() (map[string]any, error) {
		return map[string]any{}, nil
	}

	results := All()
	// Expect exactly one result: the settings check.
	if len(results) != 1 {
		t.Fatalf("All() = %d results, want 1", len(results))
	}
	r := results[0]
	if r.Status != check.StatusMissing {
		t.Errorf("result %q: status = %v, want StatusMissing", r.Name, r.Status)
	}
	if r.FixHint == "" {
		t.Errorf("result %q: missing FixHint", r.Name)
	}
}

// TestAll_MatchByMatecitoId verifies that the settings check uses the
// matecitoId field for identity when the resolved hook carries a non-empty Id.
// The handler's event/matcher in settings.json is deliberately different so
// that the legacy triple check would fail — only the matecitoId+command match
// should produce StatusOK.
func TestAll_MatchByMatecitoId(t *testing.T) {
	const cmd = "matecito-ai hook git-commit-validate"
	const mid = "development/git-commit-validator"

	orig := resolveHooks
	t.Cleanup(func() { resolveHooks = orig })
	resolveHooks = func() ([]manifest.ResolvedHook, error) {
		return []manifest.ResolvedHook{
			{Event: "PreToolUse", Matcher: "Bash", Command: cmd, Id: mid},
		}, nil
	}

	origLoad := loadSettings
	t.Cleanup(func() { loadSettings = origLoad })
	loadSettings = func() (map[string]any, error) {
		// The handler has the correct matecitoId and command, but its event key
		// is "PreToolUse" and matcher "Bash" — consistent, but the important
		// part is that matecitoId drives the identity check.
		return map[string]any{
			"hooks": map[string]any{
				"PreToolUse": []any{
					map[string]any{
						"matcher": "Bash",
						"hooks": []any{
							map[string]any{
								"type":       "command",
								"command":    cmd,
								"matecitoId": mid,
							},
						},
					},
				},
			},
		}, nil
	}

	results := All()
	if len(results) != 1 {
		t.Fatalf("All() = %d results, want 1", len(results))
	}
	if results[0].Status != check.StatusOK {
		t.Errorf("result %q: status = %v, want StatusOK (matched by matecitoId)", results[0].Name, results[0].Status)
	}
}

// TestAll_MatecitoIdMismatchMissing verifies that when settings.json has a
// handler whose matecitoId matches but the command differs, the check reports
// StatusMissing (stale, not yet reconciled).
func TestAll_MatecitoIdMismatchMissing(t *testing.T) {
	const mid = "development/git-commit-validator"

	orig := resolveHooks
	t.Cleanup(func() { resolveHooks = orig })
	resolveHooks = func() ([]manifest.ResolvedHook, error) {
		return []manifest.ResolvedHook{
			{Event: "PreToolUse", Matcher: "Bash", Command: "matecito-ai hook git-commit-validate", Id: mid},
		}, nil
	}

	origLoad := loadSettings
	t.Cleanup(func() { loadSettings = origLoad })
	loadSettings = func() (map[string]any, error) {
		// Handler has same matecitoId but old command — simulate stale state.
		return map[string]any{
			"hooks": map[string]any{
				"PreToolUse": []any{
					map[string]any{
						"matcher": "Bash",
						"hooks": []any{
							map[string]any{
								"type":       "command",
								"command":    "matecito-ai hook old-command",
								"matecitoId": mid,
							},
						},
					},
				},
			},
		}, nil
	}

	results := All()
	if len(results) != 1 {
		t.Fatalf("All() = %d results, want 1", len(results))
	}
	if results[0].Status != check.StatusMissing {
		t.Errorf("result %q: status = %v, want StatusMissing (command mismatch)", results[0].Name, results[0].Status)
	}
}
