// Package ent holds the code ent generates from the schema definitions, plus
// the directive that regenerates it. Nothing here is written by hand.
package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/versioned-migration ./schema
