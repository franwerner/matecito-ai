package hook

// ResetRegistry replaces the package-level registry with the provided slice.
// Only for use in tests (via export_test.go) to prevent cross-test bleed.
func ResetRegistry(hooks []Hook) {
	registry = make([]Hook, len(hooks))
	copy(registry, hooks)
}
