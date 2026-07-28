package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RecordVersionRef holds the schema definition for the RecordVersionRef
// entity: one row per "## Relacionados" / "## Referencias" line of a
// record's projection. target_slug/target_owning_root are deliberately
// plain text, NOT an edge/relation to the target Record — targets get
// deleted, reappear and get renamed, and the contract explicitly forbids a
// foreign key here (see ROADMAP-2 gotcha #2): with text, that churn only
// recomputes the dangling flag instead of needing to repoint a reference.
type RecordVersionRef struct {
	ent.Schema
}

func (RecordVersionRef) Mixin() []ent.Mixin {
	return []ent.Mixin{idMixin{}, timeMixin{}}
}

func (RecordVersionRef) Fields() []ent.Field {
	return []ent.Field{
		field.String("project_id"),
		field.String("record_version_id"),
		field.String("relation"),
		field.String("target_slug"),
		field.String("target_owning_root").
			Optional(),
		field.Int8("dangling").
			Default(0),
	}
}

func (RecordVersionRef) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project", Project.Type).
			Field("project_id").
			Unique().
			Required(),
		edge.To("record_version", RecordVersion.Type).
			Field("record_version_id").
			Unique().
			Required(),
	}
}
