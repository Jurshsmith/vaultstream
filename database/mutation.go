// Code generated by ent, DO NOT EDIT.

package database

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/jurshsmith/vaultstream/database/predicate"
	"github.com/jurshsmith/vaultstream/database/record"
	"github.com/jurshsmith/vaultstream/database/signature"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeRecord    = "Record"
	TypeSignature = "Signature"
)

// RecordMutation represents an operation that mutates the Record nodes in the graph.
type RecordMutation struct {
	config
	op               Op
	typ              string
	id               *int
	inserted_at      *time.Time
	clearedFields    map[string]struct{}
	signature        *int
	clearedsignature bool
	done             bool
	oldValue         func(context.Context) (*Record, error)
	predicates       []predicate.Record
}

var _ ent.Mutation = (*RecordMutation)(nil)

// recordOption allows management of the mutation configuration using functional options.
type recordOption func(*RecordMutation)

// newRecordMutation creates new mutation for the Record entity.
func newRecordMutation(c config, op Op, opts ...recordOption) *RecordMutation {
	m := &RecordMutation{
		config:        c,
		op:            op,
		typ:           TypeRecord,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withRecordID sets the ID field of the mutation.
func withRecordID(id int) recordOption {
	return func(m *RecordMutation) {
		var (
			err   error
			once  sync.Once
			value *Record
		)
		m.oldValue = func(ctx context.Context) (*Record, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Record.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withRecord sets the old Record of the mutation.
func withRecord(node *Record) recordOption {
	return func(m *RecordMutation) {
		m.oldValue = func(context.Context) (*Record, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m RecordMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m RecordMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("database: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Record entities.
func (m *RecordMutation) SetID(id int) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *RecordMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *RecordMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Record.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetInsertedAt sets the "inserted_at" field.
func (m *RecordMutation) SetInsertedAt(t time.Time) {
	m.inserted_at = &t
}

// InsertedAt returns the value of the "inserted_at" field in the mutation.
func (m *RecordMutation) InsertedAt() (r time.Time, exists bool) {
	v := m.inserted_at
	if v == nil {
		return
	}
	return *v, true
}

// OldInsertedAt returns the old "inserted_at" field's value of the Record entity.
// If the Record object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RecordMutation) OldInsertedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldInsertedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldInsertedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldInsertedAt: %w", err)
	}
	return oldValue.InsertedAt, nil
}

// ResetInsertedAt resets all changes to the "inserted_at" field.
func (m *RecordMutation) ResetInsertedAt() {
	m.inserted_at = nil
}

// SetSignatureID sets the "signature" edge to the Signature entity by id.
func (m *RecordMutation) SetSignatureID(id int) {
	m.signature = &id
}

// ClearSignature clears the "signature" edge to the Signature entity.
func (m *RecordMutation) ClearSignature() {
	m.clearedsignature = true
}

// SignatureCleared reports if the "signature" edge to the Signature entity was cleared.
func (m *RecordMutation) SignatureCleared() bool {
	return m.clearedsignature
}

// SignatureID returns the "signature" edge ID in the mutation.
func (m *RecordMutation) SignatureID() (id int, exists bool) {
	if m.signature != nil {
		return *m.signature, true
	}
	return
}

// SignatureIDs returns the "signature" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// SignatureID instead. It exists only for internal usage by the builders.
func (m *RecordMutation) SignatureIDs() (ids []int) {
	if id := m.signature; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetSignature resets all changes to the "signature" edge.
func (m *RecordMutation) ResetSignature() {
	m.signature = nil
	m.clearedsignature = false
}

// Where appends a list predicates to the RecordMutation builder.
func (m *RecordMutation) Where(ps ...predicate.Record) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the RecordMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *RecordMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Record, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *RecordMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *RecordMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Record).
func (m *RecordMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *RecordMutation) Fields() []string {
	fields := make([]string, 0, 1)
	if m.inserted_at != nil {
		fields = append(fields, record.FieldInsertedAt)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *RecordMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case record.FieldInsertedAt:
		return m.InsertedAt()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *RecordMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case record.FieldInsertedAt:
		return m.OldInsertedAt(ctx)
	}
	return nil, fmt.Errorf("unknown Record field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *RecordMutation) SetField(name string, value ent.Value) error {
	switch name {
	case record.FieldInsertedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetInsertedAt(v)
		return nil
	}
	return fmt.Errorf("unknown Record field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *RecordMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *RecordMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *RecordMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Record numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *RecordMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *RecordMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *RecordMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Record nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *RecordMutation) ResetField(name string) error {
	switch name {
	case record.FieldInsertedAt:
		m.ResetInsertedAt()
		return nil
	}
	return fmt.Errorf("unknown Record field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *RecordMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.signature != nil {
		edges = append(edges, record.EdgeSignature)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *RecordMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case record.EdgeSignature:
		if id := m.signature; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *RecordMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *RecordMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *RecordMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedsignature {
		edges = append(edges, record.EdgeSignature)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *RecordMutation) EdgeCleared(name string) bool {
	switch name {
	case record.EdgeSignature:
		return m.clearedsignature
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *RecordMutation) ClearEdge(name string) error {
	switch name {
	case record.EdgeSignature:
		m.ClearSignature()
		return nil
	}
	return fmt.Errorf("unknown Record unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *RecordMutation) ResetEdge(name string) error {
	switch name {
	case record.EdgeSignature:
		m.ResetSignature()
		return nil
	}
	return fmt.Errorf("unknown Record edge %s", name)
}

// SignatureMutation represents an operation that mutates the Signature nodes in the graph.
type SignatureMutation struct {
	config
	op            Op
	typ           string
	id            *int
	key_id        *int
	addkey_id     *int
	value         *string
	inserted_at   *time.Time
	clearedFields map[string]struct{}
	record        *int
	clearedrecord bool
	done          bool
	oldValue      func(context.Context) (*Signature, error)
	predicates    []predicate.Signature
}

var _ ent.Mutation = (*SignatureMutation)(nil)

// signatureOption allows management of the mutation configuration using functional options.
type signatureOption func(*SignatureMutation)

// newSignatureMutation creates new mutation for the Signature entity.
func newSignatureMutation(c config, op Op, opts ...signatureOption) *SignatureMutation {
	m := &SignatureMutation{
		config:        c,
		op:            op,
		typ:           TypeSignature,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withSignatureID sets the ID field of the mutation.
func withSignatureID(id int) signatureOption {
	return func(m *SignatureMutation) {
		var (
			err   error
			once  sync.Once
			value *Signature
		)
		m.oldValue = func(ctx context.Context) (*Signature, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Signature.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withSignature sets the old Signature of the mutation.
func withSignature(node *Signature) signatureOption {
	return func(m *SignatureMutation) {
		m.oldValue = func(context.Context) (*Signature, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SignatureMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SignatureMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("database: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *SignatureMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *SignatureMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Signature.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetRecordID sets the "record_id" field.
func (m *SignatureMutation) SetRecordID(i int) {
	m.record = &i
}

// RecordID returns the value of the "record_id" field in the mutation.
func (m *SignatureMutation) RecordID() (r int, exists bool) {
	v := m.record
	if v == nil {
		return
	}
	return *v, true
}

// OldRecordID returns the old "record_id" field's value of the Signature entity.
// If the Signature object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SignatureMutation) OldRecordID(ctx context.Context) (v int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldRecordID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldRecordID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldRecordID: %w", err)
	}
	return oldValue.RecordID, nil
}

// ResetRecordID resets all changes to the "record_id" field.
func (m *SignatureMutation) ResetRecordID() {
	m.record = nil
}

// SetKeyID sets the "key_id" field.
func (m *SignatureMutation) SetKeyID(i int) {
	m.key_id = &i
	m.addkey_id = nil
}

// KeyID returns the value of the "key_id" field in the mutation.
func (m *SignatureMutation) KeyID() (r int, exists bool) {
	v := m.key_id
	if v == nil {
		return
	}
	return *v, true
}

// OldKeyID returns the old "key_id" field's value of the Signature entity.
// If the Signature object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SignatureMutation) OldKeyID(ctx context.Context) (v int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldKeyID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldKeyID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldKeyID: %w", err)
	}
	return oldValue.KeyID, nil
}

// AddKeyID adds i to the "key_id" field.
func (m *SignatureMutation) AddKeyID(i int) {
	if m.addkey_id != nil {
		*m.addkey_id += i
	} else {
		m.addkey_id = &i
	}
}

// AddedKeyID returns the value that was added to the "key_id" field in this mutation.
func (m *SignatureMutation) AddedKeyID() (r int, exists bool) {
	v := m.addkey_id
	if v == nil {
		return
	}
	return *v, true
}

// ResetKeyID resets all changes to the "key_id" field.
func (m *SignatureMutation) ResetKeyID() {
	m.key_id = nil
	m.addkey_id = nil
}

// SetValue sets the "value" field.
func (m *SignatureMutation) SetValue(s string) {
	m.value = &s
}

// Value returns the value of the "value" field in the mutation.
func (m *SignatureMutation) Value() (r string, exists bool) {
	v := m.value
	if v == nil {
		return
	}
	return *v, true
}

// OldValue returns the old "value" field's value of the Signature entity.
// If the Signature object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SignatureMutation) OldValue(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldValue is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldValue requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldValue: %w", err)
	}
	return oldValue.Value, nil
}

// ResetValue resets all changes to the "value" field.
func (m *SignatureMutation) ResetValue() {
	m.value = nil
}

// SetInsertedAt sets the "inserted_at" field.
func (m *SignatureMutation) SetInsertedAt(t time.Time) {
	m.inserted_at = &t
}

// InsertedAt returns the value of the "inserted_at" field in the mutation.
func (m *SignatureMutation) InsertedAt() (r time.Time, exists bool) {
	v := m.inserted_at
	if v == nil {
		return
	}
	return *v, true
}

// OldInsertedAt returns the old "inserted_at" field's value of the Signature entity.
// If the Signature object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SignatureMutation) OldInsertedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldInsertedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldInsertedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldInsertedAt: %w", err)
	}
	return oldValue.InsertedAt, nil
}

// ResetInsertedAt resets all changes to the "inserted_at" field.
func (m *SignatureMutation) ResetInsertedAt() {
	m.inserted_at = nil
}

// ClearRecord clears the "record" edge to the Record entity.
func (m *SignatureMutation) ClearRecord() {
	m.clearedrecord = true
	m.clearedFields[signature.FieldRecordID] = struct{}{}
}

// RecordCleared reports if the "record" edge to the Record entity was cleared.
func (m *SignatureMutation) RecordCleared() bool {
	return m.clearedrecord
}

// RecordIDs returns the "record" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// RecordID instead. It exists only for internal usage by the builders.
func (m *SignatureMutation) RecordIDs() (ids []int) {
	if id := m.record; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetRecord resets all changes to the "record" edge.
func (m *SignatureMutation) ResetRecord() {
	m.record = nil
	m.clearedrecord = false
}

// Where appends a list predicates to the SignatureMutation builder.
func (m *SignatureMutation) Where(ps ...predicate.Signature) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the SignatureMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *SignatureMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Signature, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *SignatureMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *SignatureMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Signature).
func (m *SignatureMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *SignatureMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.record != nil {
		fields = append(fields, signature.FieldRecordID)
	}
	if m.key_id != nil {
		fields = append(fields, signature.FieldKeyID)
	}
	if m.value != nil {
		fields = append(fields, signature.FieldValue)
	}
	if m.inserted_at != nil {
		fields = append(fields, signature.FieldInsertedAt)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *SignatureMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case signature.FieldRecordID:
		return m.RecordID()
	case signature.FieldKeyID:
		return m.KeyID()
	case signature.FieldValue:
		return m.Value()
	case signature.FieldInsertedAt:
		return m.InsertedAt()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *SignatureMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case signature.FieldRecordID:
		return m.OldRecordID(ctx)
	case signature.FieldKeyID:
		return m.OldKeyID(ctx)
	case signature.FieldValue:
		return m.OldValue(ctx)
	case signature.FieldInsertedAt:
		return m.OldInsertedAt(ctx)
	}
	return nil, fmt.Errorf("unknown Signature field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SignatureMutation) SetField(name string, value ent.Value) error {
	switch name {
	case signature.FieldRecordID:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRecordID(v)
		return nil
	case signature.FieldKeyID:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetKeyID(v)
		return nil
	case signature.FieldValue:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetValue(v)
		return nil
	case signature.FieldInsertedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetInsertedAt(v)
		return nil
	}
	return fmt.Errorf("unknown Signature field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *SignatureMutation) AddedFields() []string {
	var fields []string
	if m.addkey_id != nil {
		fields = append(fields, signature.FieldKeyID)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *SignatureMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case signature.FieldKeyID:
		return m.AddedKeyID()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SignatureMutation) AddField(name string, value ent.Value) error {
	switch name {
	case signature.FieldKeyID:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddKeyID(v)
		return nil
	}
	return fmt.Errorf("unknown Signature numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *SignatureMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *SignatureMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *SignatureMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Signature nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *SignatureMutation) ResetField(name string) error {
	switch name {
	case signature.FieldRecordID:
		m.ResetRecordID()
		return nil
	case signature.FieldKeyID:
		m.ResetKeyID()
		return nil
	case signature.FieldValue:
		m.ResetValue()
		return nil
	case signature.FieldInsertedAt:
		m.ResetInsertedAt()
		return nil
	}
	return fmt.Errorf("unknown Signature field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *SignatureMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.record != nil {
		edges = append(edges, signature.EdgeRecord)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *SignatureMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case signature.EdgeRecord:
		if id := m.record; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *SignatureMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *SignatureMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *SignatureMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedrecord {
		edges = append(edges, signature.EdgeRecord)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *SignatureMutation) EdgeCleared(name string) bool {
	switch name {
	case signature.EdgeRecord:
		return m.clearedrecord
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *SignatureMutation) ClearEdge(name string) error {
	switch name {
	case signature.EdgeRecord:
		m.ClearRecord()
		return nil
	}
	return fmt.Errorf("unknown Signature unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *SignatureMutation) ResetEdge(name string) error {
	switch name {
	case signature.EdgeRecord:
		m.ResetRecord()
		return nil
	}
	return fmt.Errorf("unknown Signature edge %s", name)
}
