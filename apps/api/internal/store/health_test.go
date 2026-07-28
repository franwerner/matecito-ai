package store_test

import (
	"context"
	"testing"

	"github.com/franwerner/matecito-ai/apps/api/internal/store"
)

// TestStoreStatus_OK covers R14's "store listo": an open, fully migrated
// database reports ok.
func TestStoreStatus_OK(t *testing.T) {
	ctx := context.Background()
	st, err := store.Open(ctx, t.TempDir(), testMachineID)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer st.Close()

	if got := st.StoreStatus(ctx); got != store.HealthOK {
		t.Fatalf("StoreStatus = %q, want %q", got, store.HealthOK)
	}
}

// TestStoreStatus_UnavailableWhenConnectionClosed covers the "base no
// abierta" case: StoreStatus is a local check, no mock needed — a real,
// closed connection reports unavailable.
func TestStoreStatus_UnavailableWhenConnectionClosed(t *testing.T) {
	ctx := context.Background()
	st, err := store.Open(ctx, t.TempDir(), testMachineID)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	if got := st.StoreStatus(ctx); got != store.HealthUnavailable {
		t.Fatalf("StoreStatus = %q, want %q", got, store.HealthUnavailable)
	}
}

// TestStoreStatus_OutdatedWhenRevisionMissing covers "migraciones
// pendientes o fallidas": an embedded migration file with no matching
// applied revision reports outdated. The gap is manufactured with a raw
// DELETE against the real, migrated database — same pattern
// store_test.go's CHECK-constraint tests use to exercise the database
// itself, no mocks (delivery/testing-strategy).
func TestStoreStatus_OutdatedWhenRevisionMissing(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	st, err := store.Open(ctx, dir, testMachineID)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer st.Close()

	raw := openRaw(t, dir)
	if _, err := raw.ExecContext(ctx, `DELETE FROM atlas_schema_revisions`); err != nil {
		t.Fatalf("delete revisions: %v", err)
	}

	if got := st.StoreStatus(ctx); got != store.HealthOutdated {
		t.Fatalf("StoreStatus = %q, want %q", got, store.HealthOutdated)
	}
}
