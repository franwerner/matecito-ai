package transport

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/franwerner/matecito-ai/apps/api/internal/store"
)

// ErrNotFound is this package's sentinel for unmatched routes.
var ErrNotFound = errors.New("transport: route not found")

type publicError struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

// errorCode maps an internal error to the stable public message/code pair.
// Unrecognized errors fall back to a generic message so internals never leak
// (runtime/error-handling: transport is the only boundary that translates).
//
// store.ErrBranchTaken has no code of its own: it is the defensive
// translation of the same (project_id, branch) unique index StartChange
// already checks before writing, so in practice it reports the same
// conflict as store.ErrChangeNameTaken — change_name_conflict, named in
// flow/start-change.
func errorCode(err error) (message, code string) {
	switch {
	case errors.Is(err, ErrNotFound):
		return "resource not found", "not_found"
	case errors.Is(err, store.ErrNotFound):
		return "resource not found", "not_found"
	case errors.Is(err, store.ErrChangeNameTaken), errors.Is(err, store.ErrBranchTaken):
		return "change name already used", "change_name_conflict"
	case errors.Is(err, store.ErrPathAlreadyBound):
		return "path already bound to another project", "path_already_bound"
	default:
		return "internal error", "internal_error"
	}
}

// WriteError is the transport package's single translation point from an
// internal error to the public {error, code} response.
func WriteError(w http.ResponseWriter, status int, err error) {
	message, code := errorCode(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(publicError{Error: message, Code: code})
}
