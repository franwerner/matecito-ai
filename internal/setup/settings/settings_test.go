package settings

import (
	"reflect"
	"testing"
)

// hookDocWithID returns a minimal settings.json document containing a single
// hook handler with the given matecitoId marker.
func hookDocWithID(event, matcher, command, mid string) map[string]any {
	handler := map[string]any{"type": "command", "command": command}
	if mid != "" {
		handler["matecitoId"] = mid
	}
	group := map[string]any{
		"hooks": []any{handler},
	}
	if matcher != "" {
		group["matcher"] = matcher
	}
	return map[string]any{
		"hooks": map[string]any{
			event: []any{group},
		},
	}
}

func TestMissingPatterns_AgainstExpected(t *testing.T) {
	allow := []string{"mcp__codegraph__*", "Skill"}
	expected := []string{"mcp__codegraph__*", "mcp__context7__*", "Skill"}
	got := MissingPatterns(allow, expected)
	want := []string{"mcp__context7__*"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MissingPatterns = %v, want %v", got, want)
	}
}

func TestMerge_OnlyAddsMissing(t *testing.T) {
	doc := map[string]any{
		"permissions": map[string]any{
			"allow": []any{"mcp__figma__*"},
		},
	}
	expected := []string{"mcp__figma__*", "Skill"}
	if !Merge(doc, expected) {
		t.Fatal("Merge should report a change when a pattern is missing")
	}
	got := AllowList(doc)
	// existing entry preserved, missing one appended; nothing removed.
	want := []string{"mcp__figma__*", "Skill"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("allow = %v, want %v", got, want)
	}
}

func TestMerge_NoopWhenComplete(t *testing.T) {
	doc := map[string]any{
		"permissions": map[string]any{
			"allow": []any{"Skill"},
		},
	}
	if Merge(doc, []string{"Skill"}) {
		t.Error("Merge should report no change when all patterns are present")
	}
}

// --- ReconcileHooks tests ---

// TestReconcileHooks_Idempotent verifies that reconciling when the declared set
// is already exactly present returns false (no change).
func TestReconcileHooks_Idempotent(t *testing.T) {
	const mid = "development/git-commit-validator"
	doc := hookDocWithID("PreToolUse", "Bash", "matecito-ai hook git-commit-validate", mid)
	expected := []HookEntry{
		{Event: "PreToolUse", Matcher: "Bash", Command: "matecito-ai hook git-commit-validate", If: "", Type: "command", Id: mid},
	}
	if ReconcileHooks(doc, expected) {
		t.Error("ReconcileHooks must return false (no change) when declared set is already present")
	}
}

// TestReconcileHooks_RemovesStaleAndAddsNew verifies that a handler whose
// command changed is removed (same matecitoId) and the updated one is added.
func TestReconcileHooks_RemovesStaleAndAddsNew(t *testing.T) {
	const mid = "development/git-commit-validator"
	oldCmd := "matecito-ai hook git-commit-validate-v1"
	newCmd := "matecito-ai hook git-commit-validate"

	// Doc has the old command under the same matecitoId.
	doc := hookDocWithID("PreToolUse", "Bash", oldCmd, mid)
	expected := []HookEntry{
		{Event: "PreToolUse", Matcher: "Bash", Command: newCmd, Type: "command", Id: mid},
	}
	if !ReconcileHooks(doc, expected) {
		t.Fatal("ReconcileHooks must return true (changed) when a stale handler is replaced")
	}
	entries := HookList(doc)
	if len(entries) != 1 {
		t.Fatalf("HookList after reconcile = %d entries, want 1", len(entries))
	}
	if entries[0].Command != newCmd {
		t.Errorf("command = %q, want %q", entries[0].Command, newCmd)
	}
	if entries[0].Id != mid {
		t.Errorf("matecitoId = %q, want %q", entries[0].Id, mid)
	}
}

// TestReconcileHooks_PreservesUserHandler verifies that a handler without a
// matecitoId in the same matcher group is never removed.
func TestReconcileHooks_PreservesUserHandler(t *testing.T) {
	const mid = "development/git-commit-validator"
	userCmd := "/home/u/my-hook.sh"
	matecitoCmd := "matecito-ai hook git-commit-validate"

	// Doc has a user handler (no matecitoId) and the correct matecito handler.
	doc := map[string]any{
		"hooks": map[string]any{
			"PreToolUse": []any{
				map[string]any{
					"matcher": "Bash",
					"hooks": []any{
						map[string]any{"type": "command", "command": userCmd},
						map[string]any{"type": "command", "command": matecitoCmd, "matecitoId": mid},
					},
				},
			},
		},
	}
	expected := []HookEntry{
		{Event: "PreToolUse", Matcher: "Bash", Command: matecitoCmd, Type: "command", Id: mid},
	}
	if ReconcileHooks(doc, expected) {
		t.Error("ReconcileHooks must return false (no change) when matecito handler is correct")
	}
	entries := HookList(doc)
	if len(entries) != 2 {
		t.Fatalf("HookList = %d entries, want 2 (user + matecito)", len(entries))
	}
	commands := map[string]bool{}
	for _, e := range entries {
		commands[e.Command] = true
	}
	if !commands[userCmd] {
		t.Error("user handler must be preserved")
	}
	if !commands[matecitoCmd] {
		t.Error("matecito handler must be present")
	}
}

// TestReconcileHooks_RemovesDeletedHook verifies that when a declared hook is
// removed from expected (the domain no longer ships it), its matecito handler
// is removed from settings.json and the group is cleaned up.
func TestReconcileHooks_RemovesDeletedHook(t *testing.T) {
	const mid = "development/git-commit-validator"
	cmd := "matecito-ai hook git-commit-validate"

	// Doc has the hook registered, but expected is empty (hook was deleted).
	doc := hookDocWithID("PreToolUse", "Bash", cmd, mid)
	expected := []HookEntry{} // no hooks declared
	if !ReconcileHooks(doc, expected) {
		t.Fatal("ReconcileHooks must return true (changed) when a removed hook is cleaned up")
	}
	entries := HookList(doc)
	if len(entries) != 0 {
		t.Errorf("HookList after reconcile = %d entries, want 0", len(entries))
	}
}

// TestReconcileHooks_EmptyGroupDropped verifies that a group whose only handler
// was a matecito one (now removed) does not leave an empty "hooks": [] entry.
func TestReconcileHooks_EmptyGroupDropped(t *testing.T) {
	const mid = "development/git-commit-validator"
	cmd := "matecito-ai hook git-commit-validate"

	doc := hookDocWithID("PreToolUse", "Bash", cmd, mid)
	ReconcileHooks(doc, []HookEntry{})
	hooksMap, _ := doc["hooks"].(map[string]any)
	if groups, ok := hooksMap["PreToolUse"]; ok {
		gs, _ := groups.([]any)
		if len(gs) != 0 {
			t.Errorf("expected empty or no PreToolUse group, got %v", gs)
		}
	}
	// Key must also be absent or point to empty slice.
	if _, exists := hooksMap["PreToolUse"]; exists {
		t.Error("PreToolUse key should be absent after its only group was emptied")
	}
}
