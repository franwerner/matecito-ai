package transport

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Process string `json:"process"`
	Store   string `json:"store"`
}

// healthOK is the one Store value that reports the broker ready (R14):
// anything else — unavailable, migrating, outdated — is a 503.
const healthOK = "ok"

// HealthHandler reports the process as up (the handler only runs while the
// process is alive) plus storeStatus(), injected so the caller supplies the
// real persistence status without this handler depending on the store
// package directly (Stub C, observability/health-checks: local check only).
// Responds 200 only when storeStatus() is "ok"; any other value is 503, per
// R14.
func HealthHandler(storeStatus func() string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		status := storeStatus()
		w.Header().Set("Content-Type", "application/json")
		if status != healthOK {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		_ = json.NewEncoder(w).Encode(healthResponse{
			Process: "ok",
			Store:   status,
		})
	}
}
