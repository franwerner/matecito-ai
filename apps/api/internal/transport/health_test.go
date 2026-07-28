package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// alwaysOK stands in for a store whose persistence is open and migrated —
// the one Store value HealthHandler maps to 200 (R14).
func alwaysOK() string { return "ok" }

func TestHealthHandler(t *testing.T) {
	handler := HealthHandler(alwaysOK)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var body healthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}
	if body.Process != "ok" {
		t.Errorf("Process = %q, want ok", body.Process)
	}
	if body.Store != "ok" {
		t.Errorf("Store = %q, want ok", body.Store)
	}
}

func TestHealthHandler_NoExternalDependency(t *testing.T) {
	// storeStatus is injected and, in Phase 1, never touches a real
	// dependency (edr/observability/health-checks: local check only) — the
	// endpoint still responds 200 without reaching out anywhere.
	called := false
	handler := HealthHandler(func() string {
		called = true
		return "ok"
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if !called {
		t.Error("expected storeStatus to be called")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestNewMux_HealthRoute(t *testing.T) {
	mux := NewMux(alwaysOK)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

// TestHealthHandler_NotOKReports503 covers R14's "store no listo": any
// Store value other than "ok" — unavailable, migrating, outdated — reports
// 503, not 200.
func TestHealthHandler_NotOKReports503(t *testing.T) {
	for _, state := range []string{"unavailable", "migrating", "outdated"} {
		t.Run(state, func(t *testing.T) {
			handler := HealthHandler(func() string { return state })

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()
			handler(rec, req)

			if rec.Code != http.StatusServiceUnavailable {
				t.Fatalf("status = %d, want 503", rec.Code)
			}
			var body healthResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("invalid JSON body: %v", err)
			}
			if body.Store != state {
				t.Errorf("Store = %q, want %q", body.Store, state)
			}
			if body.Process != "ok" {
				t.Errorf("Process = %q, want ok (the process itself is alive)", body.Process)
			}
		})
	}
}
