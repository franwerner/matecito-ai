package transport

import (
	"context"
	"testing"
)

func TestCheckMCPAvailability_NeverAbortsInPhase1(t *testing.T) {
	if err := CheckMCPAvailability(context.Background()); err != nil {
		t.Fatalf("stub must always succeed in Phase 1, got: %v", err)
	}
}
