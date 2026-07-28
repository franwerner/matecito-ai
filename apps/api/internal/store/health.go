package store

import (
	"context"

	atlas "ariga.io/atlas/sql/migrate"
)

// HealthState is the granular persistence status the health endpoint
// reports (R14, observability/health-checks).
type HealthState string

// The values of HealthState.
const (
	// HealthOK: the persistence is open and migrated to the latest version.
	HealthOK HealthState = "ok"
	// HealthUnavailable: the database connection itself is not answering.
	HealthUnavailable HealthState = "unavailable"
	// HealthMigrating: a migration is currently running. Declared here for
	// R14's model, but unreachable through StoreStatus today — this daemon
	// always applies every pending migration synchronously before it ever
	// binds the transport (R1, cmd/broker/lifecycle.go), so nothing can
	// query health while a migration of THIS daemon's own start is in
	// flight. It becomes reachable once a synced, already-running instance
	// needs to migrate a database live (storage-sync-model's deployable,
	// shareable database, Phase 5+); nothing in this batch produces that
	// path, so nothing computes this value yet.
	HealthMigrating HealthState = "migrating"
	// HealthOutdated: the database is open, but its applied revisions do
	// not match every migration this binary embeds (pending or failed).
	HealthOutdated HealthState = "outdated"
)

// StoreStatus is a local, on-demand check of the persistence's own state
// (R14, observability/health-checks: local check, no external dependency)
// — it never reaches outside this process's own SQLite connection. Safe to
// call repeatedly; it does no writes.
func (e *Engine) StoreStatus(ctx context.Context) HealthState {
	if err := e.db.PingContext(ctx); err != nil {
		return HealthUnavailable
	}

	outdated, err := e.hasOutdatedMigrations(ctx)
	if err != nil || outdated {
		return HealthOutdated
	}

	return HealthOK
}

// hasOutdatedMigrations reports whether this binary's embedded migrations
// and the database's applied revisions have drifted apart: an embedded file
// with no matching applied revision (pending), or one whose revision
// recorded an error or an incomplete statement count (failed/partial).
func (e *Engine) hasOutdatedMigrations(ctx context.Context) (bool, error) {
	dir, err := embeddedDir()
	if err != nil {
		return false, err
	}
	files, err := dir.Files()
	if err != nil {
		return false, err
	}

	applied, err := newRevisions(e.db).ReadRevisions(ctx)
	if err != nil {
		return false, err
	}
	byVersion := make(map[string]*atlas.Revision, len(applied))
	for _, rev := range applied {
		byVersion[rev.Version] = rev
	}

	for _, f := range files {
		rev, ok := byVersion[f.Version()]
		if !ok || rev.Error != "" || rev.Applied < rev.Total {
			return true, nil
		}
	}
	return false, nil
}
