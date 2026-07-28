package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// IterationNote holds the schema definition for the IterationNote entity,
// per flow/iteration-notes. No author column: the system is single-user
// (data-modeling: no multitenancy). Anchoring is validated declaratively via
// the CHECK constraints below, not in application code.
type IterationNote struct {
	ent.Schema
}

func (IterationNote) Mixin() []ent.Mixin {
	return []ent.Mixin{idMixin{}, timeMixin{}}
}

func (IterationNote) Fields() []ent.Field {
	return []ent.Field{
		field.String("project_id"),
		field.String("change_id"),
		field.String("event_id").
			Optional(),
		field.String("file_path").
			Optional(),
		field.Int("line_from").
			Optional(),
		field.Int("line_to").
			Optional(),
		field.String("body"),
		field.String("status").
			Default("pending"),
		field.String("delivered_at").
			Optional(),
		field.String("resolved_at").
			Optional(),
	}
}

func (IterationNote) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project", Project.Type).
			Field("project_id").
			Unique().
			Required(),
		edge.To("change", Change.Type).
			Field("change_id").
			Unique().
			Required(),
		edge.To("event", Event.Type).
			Field("event_id").
			Unique(),
	}
}

func (IterationNote) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Checks: map[string]string{
				"status_valid": "status IN ('pending','delivered','resolved')",
				"file_anchor":  "file_path IS NULL OR event_id IS NOT NULL",
				// Wrapped in one extra outer paren: the check-expr renderer
				// only auto-wraps a raw expression in "CHECK (...)" when it
				// doesn't already look parenthesized (starts with "(", ends
				// with ")"); ours already does (two OR'd groups), so without
				// this explicit outer wrap the OR spills outside the CHECK
				// clause and the migration DDL fails to parse — verified
				// against a raw CREATE TABLE before fixing.
				"range_anchor": "((line_from IS NULL AND line_to IS NULL) OR (file_path IS NOT NULL AND line_from <= line_to))",
			},
		},
	}
}
