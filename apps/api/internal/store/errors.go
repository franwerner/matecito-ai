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

	// ErrInvalidContentHash is returned when a hash string does not carry
	// the algorithm tag this store understands, shaped like a real digest
	// (R7: "sha256:<hex minúscula>", constructed and validated in one
	// place).
	ErrInvalidContentHash = errors.New("store: invalid content hash")

	// ErrInvalidNoteTransition is returned when a caller asks for an
	// iteration-note status change its current status does not allow
	// (R10: pending -> delivered -> resolved, never backwards, never
	// skipped). This is the write-edge enforcement the spec deferred to
	// Batch 3 instead of a trigger.
	ErrInvalidNoteTransition = errors.New("store: invalid iteration note transition")

	// ErrNoteNotDeletable is returned when a physical delete is attempted
	// on an iteration note that is not `pending` (R11: the only physical
	// DELETE this store exposes). Same write-edge enforcement as
	// ErrInvalidNoteTransition.
	ErrNoteNotDeletable = errors.New("store: iteration note not deletable")
)
