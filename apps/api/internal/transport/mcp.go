package transport

import "context"

// CheckMCPAvailability is the Phase 1 stub for the startup MCP availability
// gate (contracts/mcp-server: the daemon verifies the MCP is up before
// serving). There is no MCP surface yet — Phase 3 builds it in this same
// package — so this always succeeds trivially and the daemon never aborts
// on this check in Phase 1 (spec Stub B).
func CheckMCPAvailability(_ context.Context) error {
	return nil
}
