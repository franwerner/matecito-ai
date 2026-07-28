package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Record holds the schema definition for the Record entity: identity by
// slug + owning_root, branch-independent, per process/index-decision-records.
type Record struct {
	ent.Schema
}

func (Record) Mixin() []ent.Mixin {
	return []ent.Mixin{idMixin{}, timeMixin{}}
}

func (Record) Fields() []ent.Field {
	return []ent.Field{
		field.String("project_id"),
		field.String("owning_root"),
		field.String("slug"),
	}
}

func (Record) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project", Project.Type).
			Field("project_id").
			Unique().
			Required(),
	}
}

func (Record) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "owning_root", "slug").Unique(),
	}
}
