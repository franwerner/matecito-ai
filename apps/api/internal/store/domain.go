package store

import "time"

// ProjectStatus is the persisted lifecycle state of a project (data-modeling
// R4): active by default, inactive after a logical deregister — never
// deleted physically.
type ProjectStatus string

// The two values of ProjectStatus.
const (
	ProjectActive   ProjectStatus = "active"
	ProjectInactive ProjectStatus = "inactive"
)

// Project is the store's domain view of a projects row: time.Time and a
// typed status, never the Ent-generated *ent.Project
// (data-access-entity-framework: the narrow Repository boundary does not
// leak generated types).
type Project struct {
	ID        string
	Name      string
	Status    ProjectStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RegisterProjectCmd registers or (idempotently) re-resolves a project. ID
// is optional: empty generates a new UUID v7 (the CREATE case in
// flow/register-project — no `.id` and no name match). A caller-supplied ID
// models the LINK case, where an existing project-id was already resolved
// (from a committed `.id`, or a confirmed find_project match) before
// calling the store. Registering the same ID twice is idempotent: the
// existing row is returned unchanged, never duplicated.
type RegisterProjectCmd struct {
	ID   string
	Name string
}

// ProjectPathStatus is the lifecycle of a per-machine root-path lookup row
// (R4): active while it resolves a checkout, stale once the checkout moved
// away from that path — never deleted.
type ProjectPathStatus string

// The two values of ProjectPathStatus.
const (
	ProjectPathActive ProjectPathStatus = "active"
	ProjectPathStale  ProjectPathStatus = "stale"
)

// ProjectPath is the store's domain view of a project_paths row: a
// per-machine lookup from a canonical root path to a project, never the
// project's identity itself (register-project, storage-sync-model).
type ProjectPath struct {
	ID         string
	ProjectID  string
	MachineID  string
	RootPath   string
	Status     ProjectPathStatus
	LastSeenAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// BindProjectPathCmd binds RootPath — this machine's canonical checkout
// root — to ProjectID. Binding a second root for a project it already has
// is a legitimate worktree: both stay `active` (R4).
type BindProjectPathCmd struct {
	ProjectID string
	RootPath  string
}

// MoveProjectPathCmd relocates a project's checkout on this machine: ToRootPath
// is bound and FromRootPath transitions to `stale` in the same write (R4 "el
// repo se mueve"). Deciding that a checkout moved — rather than that a second
// worktree appeared — needs the filesystem, so it is the caller's call; the
// store only executes the transition it is told about.
type MoveProjectPathCmd struct {
	ProjectID    string
	FromRootPath string
	ToRootPath   string
}

// ChangeStatus is the lifecycle of a change (lifecycle/change): active
// while it accepts events, closed once done. Reversible only via
// SetChangeStatus — never a side effect of another operation.
type ChangeStatus string

// The two values of ChangeStatus.
const (
	ChangeActive ChangeStatus = "active"
	ChangeClosed ChangeStatus = "closed"
)

// Change is the store's domain view of a changes row.
type Change struct {
	ID        string
	ProjectID string
	Branch    string
	Name      string
	Status    ChangeStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// StartChangeCmd starts or continues the change for (ProjectID, Branch) —
// flow/start-change: the same (project, branch) always resolves to the same
// change, never a duplicate.
type StartChangeCmd struct {
	ProjectID string
	Branch    string
	Name      string
}

// AgentKind labels which kind of actor produced an event — a denormalized
// tag on the event row, never a relation to an "agents" table (data-modeling
// enmienda 5: no attributes, no lifecycle, an open and per-project set of
// names).
type AgentKind string

// The three values of AgentKind.
const (
	AgentPhase        AgentKind = "phase"
	AgentOrchestrator AgentKind = "orchestrator"
	AgentSubagent     AgentKind = "subagent"
)

// Event is the store's domain view of an events row: the append-only log
// entry and its envelope (data-modeling: `{type, payload}`). ID and Position
// are two distinct identifiers (R6) — ID is the row's stable identity,
// Position is the change's own monotonic read cursor.
type Event struct {
	ID              string
	ProjectID       string
	ChangeID        string
	Position        int64
	Type            string
	Payload         string
	OccurredAt      time.Time
	AgentKind       AgentKind
	AgentName       string
	SessionID       string
	IdempotencyHash string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// AppendEventCmd appends one event to (ProjectID, ChangeID)'s log. Payload
// MUST already be the exact canonical bytes the caller wants persisted and
// hashed for idempotency — the store neither re-serializes nor validates its
// shape (Batch 3 decision: the caller folds the anchor — change, phase —
// into Payload before calling; the payload stays opaque to the store).
// OccurredAt lets a late-arriving event (spool reconciliation) keep its
// true, earlier timestamp instead of the append time (R6 "el orden
// semántico es occurred_at"); a zero value defaults to now. AgentKind,
// AgentName and SessionID are optional labels — empty means omitted, same
// as every other "empty means absent" Cmd field in this package.
type AppendEventCmd struct {
	ProjectID  string
	ChangeID   string
	Type       string
	Payload    string
	OccurredAt time.Time
	AgentKind  AgentKind
	AgentName  string
	SessionID  string
}

// ContentObject is the store's domain view of a content_objects row:
// content-addressable storage that later phases share between record
// versions and code snapshots (storage-sync-model). Consumers reference it
// by ID, not by its natural key (R7 enmienda 8).
type ContentObject struct {
	ID        string
	ProjectID string
	Hash      string
	SizeBytes int64
	Bytes     []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PutContentObjectCmd stores Bytes content-addressed under ProjectID. The
// store computes the hash itself — no caller ever supplies one (R7 "ningún
// productor genera claves propias").
type PutContentObjectCmd struct {
	ProjectID string
	Bytes     []byte
}

// IterationNoteStatus is the lifecycle of an iteration note
// (flow/iteration-notes): pending -> delivered -> resolved, never
// backwards, never skipped.
type IterationNoteStatus string

// The three values of IterationNoteStatus.
const (
	NotePending   IterationNoteStatus = "pending"
	NoteDelivered IterationNoteStatus = "delivered"
	NoteResolved  IterationNoteStatus = "resolved"
)

// IterationNote is the store's domain view of an iteration_notes row. No
// author field: the system is single-user (data-modeling: no
// multitenancy). EventID/FilePath/LineFrom/LineTo empty/zero means that part
// of the anchor is absent (flow/iteration-notes' "nota general").
type IterationNote struct {
	ID          string
	ProjectID   string
	ChangeID    string
	EventID     string
	FilePath    string
	LineFrom    int
	LineTo      int
	Body        string
	Status      IterationNoteStatus
	DeliveredAt time.Time
	ResolvedAt  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateIterationNoteCmd anchors a new note: ChangeID is mandatory, EventID
// optional, FilePath/LineFrom/LineTo optional and nested — a range needs a
// file, a file needs an event. The store passes the shape through as given;
// the database CHECK (R10) is what rejects an invalid combination.
type CreateIterationNoteCmd struct {
	ProjectID string
	ChangeID  string
	EventID   string
	FilePath  string
	LineFrom  int
	LineTo    int
	Body      string
}

// RecordVersionStatus is the lifecycle of a record version's content
// (lifecycle/record-version): current while it can still be updated
// in-place, frozen forever once an event pins it — no transition back.
type RecordVersionStatus string

// The two values of RecordVersionStatus.
const (
	VersionCurrent RecordVersionStatus = "current"
	VersionFrozen  RecordVersionStatus = "frozen"
)

// Record is the store's domain view of a records row: identity by
// (owning_root, slug), branch-independent
// (process/index-decision-records).
type Record struct {
	ID         string
	ProjectID  string
	OwningRoot string
	Slug       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// UpsertRecordCmd resolves or creates the record for
// (ProjectID, OwningRoot, Slug) — the same identity always resolves to the
// same row, never a duplicate.
type UpsertRecordCmd struct {
	ProjectID  string
	OwningRoot string
	Slug       string
}

// RecordVersion is the store's domain view of a record_versions row: not
// scoped by branch (lifecycle/record-version) — Kind/DocStatus are the
// projection extracted from the record's content, frozen together with it.
// FrozenAt is the zero time while Status is current.
type RecordVersion struct {
	ID              string
	ProjectID       string
	RecordID        string
	ContentObjectID string
	Status          RecordVersionStatus
	Kind            string
	DocStatus       string
	FrozenAt        time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// WriteRecordVersionCmd applies a record's on-disk content to the store for
// (Branch) — the copy-on-write lazy at the heart of
// lifecycle/record-version: the branch's current version forks into a new
// one instead of updating in place only if it is frozen. Versions are never
// shared across branches; only the content they reference is deduplicated
// by hash. Kind/DocStatus are the projection extracted alongside Bytes;
// both are written every call, since a version snapshots content and
// projection together, never one without the other.
type WriteRecordVersionCmd struct {
	ProjectID  string
	OwningRoot string
	Slug       string
	Branch     string
	Bytes      []byte
	Kind       string
	DocStatus  string
}

// RecordBranchStatus is a record's soft-delete state within one
// (project, branch) — never physically removed
// (process/index-decision-records: "borrado = soft-delete").
type RecordBranchStatus string

// The two values of RecordBranchStatus.
const (
	RecordBranchActive  RecordBranchStatus = "active"
	RecordBranchDeleted RecordBranchStatus = "deleted"
)

// RecordBranchState is the store's domain view of a record_branch_state
// row: the per-(record, branch) vigency index. Malformed flags a `.md` that
// failed to parse in this branch, without losing CurrentVersionID's last
// good version.
type RecordBranchState struct {
	ID               string
	ProjectID        string
	RecordID         string
	Branch           string
	Status           RecordBranchStatus
	CurrentVersionID string
	Malformed        bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// RecordVersionRef is the store's domain view of a record_version_refs row:
// one relation line from a version's projection. TargetSlug/TargetOwningRoot
// are plain text, never a foreign key (R9) — Dangling is recomputed on
// demand by the store, never enforced by a database constraint.
type RecordVersionRef struct {
	ID               string
	ProjectID        string
	RecordVersionID  string
	Relation         string
	TargetSlug       string
	TargetOwningRoot string
	Dangling         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// CreateRecordVersionRefCmd adds one relation row to RecordVersionID.
// Dangling is computed by the store at creation time from whether the
// target currently resolves to a record — never supplied by the caller
// (R9).
type CreateRecordVersionRefCmd struct {
	ProjectID        string
	RecordVersionID  string
	Relation         string
	TargetSlug       string
	TargetOwningRoot string
}

// EventRecordPin is the store's domain view of an event_record_pins row:
// the pin that freezes a record version the first time an event resolves
// it (lifecycle/record-version).
type EventRecordPin struct {
	ID              string
	ProjectID       string
	EventID         string
	RecordVersionID string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// PinEventToRecordVersionCmd pins EventID to (RecordID, Branch)'s current
// version, resolved implicitly by the store — the caller never names a
// version id (lifecycle/record-version: "el agente nunca manda una
// versión"). The pin freezes that version if this is its first pin;
// repeating the same (event, version) pair is idempotent and returns the
// existing pin unchanged.
type PinEventToRecordVersionCmd struct {
	ProjectID string
	EventID   string
	RecordID  string
	Branch    string
}
