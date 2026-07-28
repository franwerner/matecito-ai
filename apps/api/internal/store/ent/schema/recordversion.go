package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RecordVersion holds the schema definition for the RecordVersion entity:
// not scoped by branch, per lifecycle/record-version — a version is a
// snapshot of the content AND the projection (kind, doc_status), so both
// freeze together. Deliberately has NO unique index on "one current version
// per record": two branches can have two simultaneously-current versions;
// vigency is expressed per-branch by RecordBranchState.current_version_id.
//
// content_object_id references ContentObject.id (a plain relation) rather
// than a composite (project_id, content_hash) foreign key, because
// ContentObject no longer has that composite as its primary key (see the
// ContentObject schema doc and ROADMAP-2 anexo "Nota — claves compuestas").
// The content hash is a join away through that relation; dedup is
// unaffected since ContentObject still enforces UNIQUE(project_id, hash).
type RecordVersion struct {
	ent.Schema
}

func (RecordVersion) Mixin() []ent.Mixin {
	return []ent.Mixin{idMixin{}, timeMixin{}}
}

func (RecordVersion) Fields() []ent.Field {
	return []ent.Field{
		field.String("project_id"),
		field.String("record_id"),
		field.String("content_object_id"),
		field.String("status"),
		field.String("kind").
			Optional(),
		field.String("doc_status").
			Optional(),
		field.String("frozen_at").
			Optional(),
	}
}

func (RecordVersion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project", Project.Type).
			Field("project_id").
			Unique().
			Required(),
		edge.To("record", Record.Type).
			Field("record_id").
			Unique().
			Required(),
		edge.To("content_object", ContentObject.Type).
			Field("content_object_id").
			Unique().
			Required(),
	}
}

func (RecordVersion) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Checks: map[string]string{
				"status_valid": "status IN ('current','frozen')",
			},
		},
	}
}
