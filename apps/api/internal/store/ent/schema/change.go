package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Change holds the schema definition for the Change entity. branch is
// mandatory — per the event-scoping/start-change amendment there is no
// change without a branch, and no fallback sentinel.
type Change struct {
	ent.Schema
}

func (Change) Mixin() []ent.Mixin {
	return []ent.Mixin{idMixin{}, timeMixin{}}
}

func (Change) Fields() []ent.Field {
	return []ent.Field{
		field.String("project_id"),
		field.String("branch"),
		field.String("name"),
		field.String("status").
			Default("active"),
	}
}

func (Change) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project", Project.Type).
			Field("project_id").
			Unique().
			Required(),
	}
}

func (Change) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "branch").Unique(),
		index.Fields("project_id", "name").Unique(),
	}
}

func (Change) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Checks: map[string]string{
				"status_valid": "status IN ('active','closed')",
			},
		},
	}
}
