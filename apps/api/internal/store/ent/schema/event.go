package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Event holds the schema definition for the Event entity: the append-only
// event-log, per data-modeling and submit-phase-artifact. The agent is a
// denormalized label on the event (agent_kind/agent_name/session_id) — NOT a
// relation to any "agents" table, which does not exist (see data-modeling
// enmienda 5 and the ROADMAP-2 anexo).
type Event struct {
	ent.Schema
}

func (Event) Mixin() []ent.Mixin {
	return []ent.Mixin{idMixin{}, timeMixin{}}
}

func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.String("project_id"),
		field.String("change_id"),
		field.Int64("position"),
		field.String("type"),
		field.String("payload"),
		field.String("occurred_at").
			DefaultFunc(nowRFC3339Milli),
		field.String("agent_kind").
			Optional(),
		field.String("agent_name").
			Optional(),
		field.String("session_id").
			Optional(),
		field.String("idempotency_hash"),
	}
}

func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project", Project.Type).
			Field("project_id").
			Unique().
			Required(),
		edge.To("change", Change.Type).
			Field("change_id").
			Unique().
			Required(),
	}
}

func (Event) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("change_id", "position").Unique(),
		index.Fields("idempotency_hash").Unique(),
	}
}

func (Event) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Checks: map[string]string{
				"agent_kind_valid": "agent_kind IN ('phase','orchestrator','subagent')",
			},
		},
	}
}
