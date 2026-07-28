package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// RecordBranchState holds the schema definition for the RecordBranchState
// entity: the per-(record, branch) vigency index, per
// process/index-decision-records. Naturally keyed by (record_id, branch);
// carries its own surrogate id for the same reason as ContentObject (see
// that schema's doc comment and ROADMAP-2 anexo "Nota — claves compuestas").
type RecordBranchState struct {
	ent.Schema
}

func (RecordBranchState) Mixin() []ent.Mixin {
	return []ent.Mixin{idMixin{}, timeMixin{}}
}

func (RecordBranchState) Fields() []ent.Field {
	return []ent.Field{
		field.String("project_id"),
		field.String("record_id"),
		field.String("branch"),
		field.String("status"),
		field.String("current_version_id"),
		field.Int8("malformed").
			Default(0),
	}
}

func (RecordBranchState) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project", Project.Type).
			Field("project_id").
			Unique().
			Required(),
		edge.To("record", Record.Type).
			Field("record_id").
			Unique().
			Required(),
		edge.To("current_version", RecordVersion.Type).
			Field("current_version_id").
			Unique().
			Required(),
	}
}

func (RecordBranchState) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("record_id", "branch").Unique(),
	}
}

func (RecordBranchState) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			// Ent's inflector pluralizes "RecordBranchState" to
			// "record_branch_states"; the contract fixes the table name as
			// the singular "record_branch_state" (Anexo A).
			Table: "record_branch_state",
			Checks: map[string]string{
				"status_valid": "status IN ('active','deleted')",
			},
		},
	}
}
