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
