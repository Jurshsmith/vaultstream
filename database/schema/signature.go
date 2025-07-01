package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Signature holds the schema definition for the Signature entity.
type Signature struct {
	ent.Schema
}

// Fields of the Signature.
func (Signature) Fields() []ent.Field {
	return []ent.Field{
		// Use record_id as the primary key.
		field.Int("record_id").
			Unique().
			Immutable().
			StructTag(`json:"record_id"`),
		// The key ID is a positive integer.
		field.Int("key_id").
			Positive().
			StructTag(`json:"key_id"`),
		// The signature value as non-empty text.
		field.String("value").
			NotEmpty().
			StructTag(`json:"value"`),
		// Creation timestamp.
		field.Time("inserted_at").
			Default(time.Now).
			Immutable().
			StructTag(`json:"inserted_at"`),
	}
}

// Edges of the Signature.
func (Signature) Edges() []ent.Edge {
	return []ent.Edge{
		// Each signature belongs to one record.
		edge.From("record", Record.Type).
			Ref("signature").
			Unique().
			Required().
			Immutable(). // <-- Added to match the immutable field.
			Field("record_id"),
	}
}

// Indexes of the Signature.
func (Signature) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("value").Unique(),
	}
}
