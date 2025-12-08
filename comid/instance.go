// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
)

// Instance stores an instance identity. The supported formats are UUID, UEID and variable-length opaque bytes.
type Instance struct {
	Value IInstanceValue
}

// NewInstance creates a new instance with the value of the specified type
// populated using the provided value.
func NewInstance(val any, typ string) (*Instance, error) {
	factory, ok := instanceValueRegister[typ]
	if !ok {
		return nil, fmt.Errorf("unknown instance type: %s", typ)
	}

	return factory(val)
}

// Valid checks for the validity of given instance
func (o Instance) Valid() error {
	if o.String() == "" {
		return fmt.Errorf("invalid instance id")
	}
	return nil
}

// String returns a printable string of the Instance value.  UUIDs use the
// canonical 8-4-4-4-12 format, UEIDs are hex encoded.
func (o Instance) String() string {
	if o.Value == nil {
		return ""
	}

	return o.Value.String()
}

// Type returns a string naming the type of the underlying Instance value.
func (o Instance) Type() string {
	return o.Value.Type()
}

// Bytes returns a []byte containing the bytes of the underlying Instance
// value.
func (o Instance) Bytes() []byte {
	return o.Value.Bytes()
}

// MarshalCBOR serializes the target instance to CBOR
func (o Instance) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.Value)
}

func (o *Instance) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &o.Value)
}

// UnmarshalJSON deserializes the supplied JSON object into the target Instance
// The instance object must have the following shape:
//
//	{
//	  "type": "<INSTANCE_TYPE>",
//	  "value": <INSTANCE_VALUE>
//	}
//
// where <INSTANCE_TYPE> must be one of the known IInstanceValue implementation
// type names (available in the base implementation: "ueid" and "uuid"), and
// <INSTANCE_VALUE> is the JSON encoding of the instance value. The exact
// encoding is <INSTANCE_TYPE> dependent. For the base implmentation types it is
//
//	ueid: base64-encoded bytes, e.g. "YWNtZS1pbXBsZW1lbnRhdGlvbi1pZC0wMDAwMDAwMDE="
//	uuid: standard UUID string representation, e.g. "550e8400-e29b-41d4-a716-446655440000"
//	bytes: a variable-length opaque byte string, example {0x07, 0x12, 0x34}

//nolint:dupl
func (o *Instance) UnmarshalJSON(data []byte) error {
	var tnv encoding.TypeAndValue

	if err := json.Unmarshal(data, &tnv); err != nil {
		return fmt.Errorf("instance decoding failure: %w", err)
	}

	decoded, err := NewInstance(nil, tnv.Type)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(tnv.Value, &decoded.Value); err != nil {
		return fmt.Errorf(
			"cannot unmarshal instance: %w",
			err,
		)
	}

	if err := decoded.Value.Valid(); err != nil {
		return fmt.Errorf("invalid %s: %w", tnv.Type, err)
	}

	o.Value = decoded.Value

	return nil
}

// MarshalJSON serializes the Instance into a JSON object.
func (o Instance) MarshalJSON() ([]byte, error) {
	valueBytes, err := json.Marshal(o.Value)
	if err != nil {
		return nil, err
	}

	value := encoding.TypeAndValue{
		Type:  o.Value.Type(),
		Value: valueBytes,
	}

	return json.Marshal(value)
}

// SetUEID sets the identity of the target instance to the supplied UEID
func (o *Instance) SetUEID(val eat.UEID) *Instance {
	if o != nil {
		if val.Validate() != nil {
			return nil
		}
		o.Value = TaggedUEID(val)
	}
	return o
}

// SetUUID sets the identity of the target instance to the supplied UUID
func (o *Instance) SetUUID(val uuid.UUID) *Instance {
	if o != nil {
		o.Value = TaggedUUID(val)
	}
	return o
}

func (o Instance) GetUEID() (eat.UEID, error) {
	switch t := o.Value.(type) {
	case TaggedUEID:
		return eat.UEID(t), nil
	case *TaggedUEID:
		return eat.UEID(*t), nil
	default:
		return eat.UEID{}, fmt.Errorf("instance-id type is: %T", t)
	}
}

func (o Instance) GetUUID() (UUID, error) {
	switch t := o.Value.(type) {
	case *TaggedUUID:
		return UUID(*t), nil
	case TaggedUUID:
		return UUID(t), nil
	default:
		return UUID{}, fmt.Errorf("instance-id type is: %T", t)
	}
}

// IInstanceValue is the interface implemented by all Instance value
// implementations.
type IInstanceValue interface {
	extensions.ITypeChoiceValue

	Bytes() []byte
}

// NewUEIDInstance instantiates a new instance with the supplied UEID identity.
func NewUEIDInstance(val any) (*Instance, error) {
	if val == nil {
		return &Instance{&TaggedUEID{}}, nil
	}

	ret, err := NewTaggedUEID(val)
	if err != nil {
		return nil, err
	}
	return &Instance{ret}, nil
}

// MustNewBytesInstance is like NewBytesInstance except it does not return an
// error, assuming that the provided value is valid. It panics if that isn't
// the case.
func MustNewBytesInstance(val any) *Instance {
	ret, err := NewBytesInstance(val)
	if err != nil {
		panic(err)
	}

	return ret
}

// MustNewUEIDInstance is like NewUEIDInstance execept it does not return an
// error, assuming that the provided value is valid. It panics if that isn't
// the case.
func MustNewUEIDInstance(val any) *Instance {
	ret, err := NewUEIDInstance(val)
	if err != nil {
		panic(err)
	}
	return ret
}

// NewUUIDInstance instantiates a new instance with the supplied UUID identity
func NewUUIDInstance(val any) (*Instance, error) {
	if val == nil {
		return &Instance{&TaggedUUID{}}, nil
	}

	ret, err := NewTaggedUUID(val)
	if err != nil {
		return nil, err
	}

	return &Instance{ret}, nil
}

// NewBytesInstance creates a new instance of type bytes
// The supplied interface parameter could be
// a byte slice, a pointer to a byte slice or a string
func NewBytesInstance(val any) (*Instance, error) {
	ret, err := NewBytes(val)
	if err != nil {
		return nil, err
	}
	return &Instance{ret}, nil
}

// MustNewUUIDInstance is like NewUUIDInstance execept it does not return an
// error, assuming that the provided value is valid. It panics if that isn't
// the case.
func MustNewUUIDInstance(val any) *Instance {
	ret, err := NewUUIDInstance(val)
	if err != nil {
		panic(err)
	}

	return ret
}

// IInstanceFactory defines the signature for the factory functions that may be
// registered using RegisterInstanceType to provide a new implementation of the
// corresponding type choice. The factory function should create a new *Instance
// with the underlying value created based on the provided input. The range of
// valid inputs is up to the specific type choice implementation, however it
// _must_ accept nil as one of the inputs, and return the Zero value for
// implemented type.
// See also https://go.dev/ref/spec#The_zero_value
type IInstanceFactory func(any) (*Instance, error)

var instanceValueRegister = map[string]IInstanceFactory{
	UEIDType:  NewUEIDInstance,
	UUIDType:  NewUUIDInstance,
	BytesType: NewBytesInstance,
}

// RegisterInstanceType registers a new IInstanceValue implementation (created
// by the provided IInstanceFactory) under the specified CBOR tag.
func RegisterInstanceType(tag uint64, factory IInstanceFactory) error {
	nilVal, err := factory(nil)
	if err != nil {
		return err
	}

	typ := nilVal.Type()
	if _, exists := instanceValueRegister[typ]; exists {
		return fmt.Errorf("class ID type with name %q already exists", typ)
	}

	if err := registerCOMIDTag(tag, nilVal.Value); err != nil {
		return err
	}

	instanceValueRegister[typ] = factory

	return nil
}
