package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ProjectPath holds the schema definition for the ProjectPath entity: the
// per-machine lookup from a canonical root path to a project-id, per
// register-project and storage-sync-model. Never the identity itself.
type ProjectPath struct {
	ent.Schema
}

func (ProjectPath) Mixin() []ent.Mixin {
	return []ent.Mixin{idMixin{}, timeMixin{}}
}

func (ProjectPath) Fields() []ent.Field {
	return []ent.Field{
		field.String("project_id"),
		field.String("machine_id"),
		field.String("root_path"),
		field.String("status").
			Default("active"),
		field.String("last_seen_at"),
	}
}

func (ProjectPath) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project", Project.Type).
			Field("project_id").
			Unique().
			Required(),
	}
}

func (ProjectPath) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("machine_id", "root_path").Unique(),
	}
}

func (ProjectPath) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Checks: map[string]string{
				"status_valid": "status IN ('active','stale')",
			},
		},
	}
}
