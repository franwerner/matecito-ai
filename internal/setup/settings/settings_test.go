package settings

import (
	"reflect"
	"testing"
)

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
