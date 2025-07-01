package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Record holds the schema definition for the Record entity.
type Record struct {
	ent.Schema
}

// Fields of the Record.
func (Record) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Unique().
			Immutable().
			StructTag(`json:"id"`),
		field.Time("inserted_at").
			Default(time.Now).
			StructTag(`json:"inserted_at"`),
	}
}

// Edges of the Record.
func (Record) Edges() []ent.Edge {
	return []ent.Edge{
		// The inverse of Signature's "record" edge.
		edge.To("signature", Signature.Type).Unique(),
	}
}
