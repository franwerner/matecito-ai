package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	handler := HealthHandler(func() string { return "ok" })

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
	mux := NewMux(StoreStatus)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}
