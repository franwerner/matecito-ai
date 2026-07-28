package store

import "errors"

// Sentinels for the write-side cases a caller must distinguish with
// errors.Is (runtime/error-handling: "cada paquete expone sus propios
// sentinels"). Read-side lookups that find nothing all resolve to
// ErrNotFound; write-side conflicts get their own sentinel per case so the
// transport boundary can translate each to its own public code.
var (
	// ErrNotFound is returned when a lookup by identity finds no row.
	ErrNotFound = errors.New("store: not found")

	// ErrChangeNameTaken covers both ways a change name conflicts
	// (flow/start-change's change_name_conflict): the branch already maps
	// to a different name (StartChange refuses to re-map it), or the
	// (project_id, name) unique index rejects a new change whose name is
	// already used by another branch in the same project (R5).
	ErrChangeNameTaken = errors.New("store: change name already used")

	// ErrBranchTaken is the defensive translation of the (project_id,
	// branch) unique index. StartChange always looks up the existing row
	// before creating one, inside the same write-transaction, so this
	// should not surface in normal operation — it exists so a constraint
	// violation is never reported as a bare internal error.
	ErrBranchTaken = errors.New("store: branch already mapped")

	// ErrPathAlreadyBound is returned when a (machine_id, root_path) is
	// already bound to a different project.
	ErrPathAlreadyBound = errors.New("store: path bound elsewhere")
)
