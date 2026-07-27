package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteError_NotFound(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteError(rec, http.StatusNotFound, ErrNotFound)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	var body publicError
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}
	if body.Code != "not_found" {
		t.Errorf("Code = %q, want not_found", body.Code)
	}
}

func TestWriteError_NeverLeaksInternals(t *testing.T) {
	internal := errors.New("sqlite: disk I/O error at /home/user/.matecito-ai/broker.db")
	wrapped := fmt.Errorf("store: query failed: %w", internal)

	rec := httptest.NewRecorder()
	WriteError(rec, http.StatusInternalServerError, wrapped)

	body := rec.Body.String()
	if strings.Contains(body, "/home/user") || strings.Contains(body, "sqlite") {
		t.Fatalf("response leaked internal details: %s", body)
	}

	var parsed publicError
	if err := json.Unmarshal(rec.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}
	if parsed.Code != "internal_error" {
		t.Errorf("Code = %q, want internal_error", parsed.Code)
	}
}

func TestNewMux_UnmatchedRouteUsesErrorBoundary(t *testing.T) {
	mux := NewMux(StoreStatus)
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	var body publicError
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}
	if body.Code != "not_found" {
		t.Errorf("Code = %q, want not_found", body.Code)
	}
}
