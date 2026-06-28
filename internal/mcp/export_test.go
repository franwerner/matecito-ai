package mcp

// SetRunMCPList replaces the CLI runner for tests. Call ResetRunMCPList in
// t.Cleanup to restore the default.
func SetRunMCPList(fn func() ([]byte, error)) {
	runMCPList = fn
}

// ResetRunMCPList restores the default CLI runner and invalidates the cache.
func ResetRunMCPList() {
	runMCPList = defaultRunMCPList
	InvalidateCLICache()
}
