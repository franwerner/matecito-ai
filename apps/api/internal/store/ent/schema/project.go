package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Project holds the schema definition for the Project entity: identity by
// project-id, per data-modeling and storage-sync-model. name is NOT unique —
// find_project(name) can return several matches and the caller confirms.
type Project struct {
	ent.Schema
}

func (Project) Mixin() []ent.Mixin {
	return []ent.Mixin{idMixin{}, timeMixin{}}
}

func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.String("status").
			Default("active"),
	}
}

func (Project) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Checks: map[string]string{
				"status_valid": "status IN ('active','inactive')",
			},
		},
	}
}
