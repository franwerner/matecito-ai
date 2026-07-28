// Package migrations embeds the broker's versioned migration files, derived
// from the Ent schema (never edited by hand — see
// data/data-access-entity-framework and generate/main.go).
package migrations

import "embed"

// FS embeds every generated migration file plus the Atlas checksum file, so
// the daemon binary carries and applies them without any external file or
// tool at runtime.
//
//go:embed *.sql atlas.sum
var FS embed.FS
