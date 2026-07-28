package store

import "context"

// Reader is the store's read side: the narrow set of lookups the rest of
// the broker needs, never the generated Ent client or database/sql
// (data-access-entity-framework).
type Reader interface {
	FindProjectByID(ctx context.Context, id string) (Project, error)
	FindProjectsByName(ctx context.Context, name string) ([]Project, error)
	ResolveProjectPath(ctx context.Context, root string) (Project, error)
	FindChangeByBranch(ctx context.Context, projectID, branch string) (Change, error)
	ListEvents(ctx context.Context, changeID string) ([]Event, error)
	ListEventsFrom(ctx context.Context, changeID string, fromPosition int64) ([]Event, error)
	GetContentObject(ctx context.Context, projectID, hash string) (ContentObject, error)
	GetIterationNote(ctx context.Context, id string) (IterationNote, error)
	GetRecord(ctx context.Context, projectID, owningRoot, slug string) (Record, error)
	GetRecordVersion(ctx context.Context, id string) (RecordVersion, error)
	GetRecordBranchState(ctx context.Context, recordID, branch string) (RecordBranchState, error)
	GetRecordVersionRef(ctx context.Context, id string) (RecordVersionRef, error)
}

// Writer is the store's write side. Every method is synchronous to its
// caller — it enqueues a command on the single writer goroutine and blocks
// for the result, so at most one write transaction is ever in flight
// (runtime/concurrency-async: single-writer fed by channel).
type Writer interface {
	RegisterProject(ctx context.Context, cmd RegisterProjectCmd) (Project, error)
	DeactivateProject(ctx context.Context, id string) error
	BindProjectPath(ctx context.Context, cmd BindProjectPathCmd) (ProjectPath, error)
	MoveProjectPath(ctx context.Context, cmd MoveProjectPathCmd) (ProjectPath, error)
	StartChange(ctx context.Context, cmd StartChangeCmd) (Change, error)
	SetChangeStatus(ctx context.Context, id string, s ChangeStatus) error
	AppendEvent(ctx context.Context, cmd AppendEventCmd) (Event, error)
	PutContentObject(ctx context.Context, cmd PutContentObjectCmd) (ContentObject, error)
	CreateIterationNote(ctx context.Context, cmd CreateIterationNoteCmd) (IterationNote, error)
	DeliverIterationNote(ctx context.Context, id string) (IterationNote, error)
	ResolveIterationNote(ctx context.Context, id string) (IterationNote, error)
	DeleteIterationNote(ctx context.Context, id string) error
	UpsertRecord(ctx context.Context, cmd UpsertRecordCmd) (Record, error)
	WriteRecordVersion(ctx context.Context, cmd WriteRecordVersionCmd) (RecordVersion, error)
	SetRecordDeleted(ctx context.Context, recordID, branch string) (RecordBranchState, error)
	SetRecordMalformed(ctx context.Context, recordID, branch string, malformed bool) (RecordBranchState, error)
	CreateRecordVersionRef(ctx context.Context, cmd CreateRecordVersionRefCmd) (RecordVersionRef, error)
	RecalculateDanglingRefs(ctx context.Context, projectID, owningRoot, slug string) error
	PinEventToRecordVersion(ctx context.Context, cmd PinEventToRecordVersionCmd) (EventRecordPin, error)
}

// Repository is the persistence boundary the rest of the broker consumes: a
// narrow read side plus a narrow write side — never repositories per
// entity, never the generated client (data-access-entity-framework).
type Repository interface {
	Reader
	Writer
}
