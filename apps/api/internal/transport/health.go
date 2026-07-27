package transport

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Process string `json:"process"`
	Store   string `json:"store"`
}

// StoreStatus is the Phase 1 stub for the persistence health field (spec
// Stub C): there is no real store yet — Phase 2 wires SQLite + migrations —
// so it always reports "ok" and health never depends on a real store.
func StoreStatus() string {
	return "ok"
}

// HealthHandler reports the process as up (the handler only runs while the
// process is alive) plus storeStatus(), injected so Phase 2 can swap in the
// real store's status without changing this handler.
func HealthHandler(storeStatus func() string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(healthResponse{
			Process: "ok",
			Store:   storeStatus(),
		})
	}
}
