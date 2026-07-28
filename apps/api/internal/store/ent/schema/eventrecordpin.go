package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// EventRecordPin holds the schema definition for the EventRecordPin entity:
// pins an event to the record version it froze, per lifecycle/record-version.
// Naturally keyed by (event_id, record_version_id); carries its own
// surrogate id for the same reason as ContentObject (see that schema's doc
// comment and ROADMAP-2 anexo "Nota — claves compuestas").
type EventRecordPin struct {
	ent.Schema
}

func (EventRecordPin) Mixin() []ent.Mixin {
	return []ent.Mixin{idMixin{}, timeMixin{}}
}

func (EventRecordPin) Fields() []ent.Field {
	return []ent.Field{
		field.String("project_id"),
		field.String("event_id"),
		field.String("record_version_id"),
	}
}

func (EventRecordPin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project", Project.Type).
			Field("project_id").
			Unique().
			Required(),
		edge.To("event", Event.Type).
			Field("event_id").
			Unique().
			Required(),
		edge.To("record_version", RecordVersion.Type).
			Field("record_version_id").
			Unique().
			Required(),
	}
}

func (EventRecordPin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("event_id", "record_version_id").Unique(),
	}
}
