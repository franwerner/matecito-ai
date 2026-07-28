package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ContentObject holds the schema definition for the ContentObject entity:
// content-addressable storage shared by record versions (Phase 5) and code
// snapshots (Phase 8), per storage-sync-model. Naturally keyed by
// (project_id, hash); carries its own surrogate id because Ent v0.14.6 only
// supports a composite primary key without a surrogate on edge-schema join
// tables, not on a regular entity (verified gate, see ROADMAP-2 anexo "Nota
// — claves compuestas"). The UNIQUE index below gives the same integrity
// guarantee as the composite PK would have.
type ContentObject struct {
	ent.Schema
}

func (ContentObject) Mixin() []ent.Mixin {
	return []ent.Mixin{idMixin{}, timeMixin{}}
}

func (ContentObject) Fields() []ent.Field {
	return []ent.Field{
		field.String("project_id"),
		field.String("hash"),
		field.Int64("size_bytes"),
		field.Bytes("bytes"),
	}
}

func (ContentObject) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project", Project.Type).
			Field("project_id").
			Unique().
			Required(),
	}
}

func (ContentObject) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "hash").Unique(),
	}
}
