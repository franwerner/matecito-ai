package schema

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// idMixin gives every entity a UUID v7 primary key, per data-modeling: IDs
// are TEXT, canonical-with-hyphens, ordered by creation time.
type idMixin struct {
	mixin.Schema
}

func (idMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(newUUIDv7).
			Immutable().
			Unique(),
	}
}

func newUUIDv7() string {
	id, err := uuid.NewV7()
	if err != nil {
		// crypto/rand failure is not recoverable at the call site; every
		// caller of a Default field func has no error return to propagate to.
		panic(fmt.Sprintf("schema: uuid v7 generation failed: %v", err))
	}
	return id.String()
}

// timeMixin gives every entity created_at/updated_at, per data-modeling: TEXT
// RFC 3339 UTC with milliseconds, NOT NULL.
type timeMixin struct {
	mixin.Schema
}

func (timeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("created_at").
			Immutable().
			DefaultFunc(nowRFC3339Milli),
		field.String("updated_at").
			DefaultFunc(nowRFC3339Milli),
	}
}

// Hooks refreshes updated_at on every update. field.String has no
// UpdateDefault (that builder method only exists for numeric/time field
// types), so the refresh is done generically via Mutation.SetField instead.
// Scoped to update ops directly against ent.Op — the hook.On helper lives in
// the generated package, which doesn't exist yet the first time this runs.
func (timeMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if m.Op()&(ent.OpUpdate|ent.OpUpdateOne) != 0 {
					if err := m.SetField("updated_at", nowRFC3339Milli()); err != nil {
						return nil, err
					}
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}

// rfc3339Milli is RFC 3339 with millisecond precision and a literal "Z" —
// values are always normalized to UTC before formatting, so the offset is
// never anything else.
const rfc3339Milli = "2006-01-02T15:04:05.000Z"

func nowRFC3339Milli() string {
	return time.Now().UTC().Format(rfc3339Milli)
}
