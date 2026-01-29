// Copyright 2021-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fortio.org/safecast"
	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
)

// ClassID identifies the environment via a well-known identifier. This can be
// an OID, a UUID, variable-length opaque bytes or a profile-defined extension type.
type ClassID struct {
	Value IClassIDValue
}

// NewClassID creates a new ClassID of the specified type using the specified value.
func NewClassID(val any, typ string) (*ClassID, error) {
	factory, ok := classIDValueRegister[typ]
	if !ok {
		return nil, fmt.Errorf("unknown class id type: %s", typ)
	}

	return factory(val)
}

// Valid returns nil if the ClassID is valid, or an error describing the
// problem, if it is not.
func (o ClassID) Valid() error {
	if o.Value == nil {
		return errors.New("nil value")
	}

	return o.Value.Valid()
}

// Type returns the type of the ClassID
func (o ClassID) Type() string {
	if o.Value == nil {
		return ""
	}

	return o.Value.Type()
}

// Bytes returns a []byte containing the raw bytes of the class id value
func (o ClassID) Bytes() []byte {
	if o.Value == nil {
		return []byte{}
	}
	return o.Value.Bytes()
}

// IsSet returns true iff the underlying class id value has been set (is not nil)
func (o ClassID) IsSet() bool {
	return o.Value != nil
}

// MarshalCBOR serializes the target ClassID to CBOR
func (o ClassID) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.Value)
}

// UnmarshalCBOR deserializes the supplied CBOR buffer into the target ClassID.
// It is undefined behavior to try and inspect the target ClassID in case this
// method returns an error.
func (o *ClassID) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &o.Value)
}

// UnmarshalJSON deserializes the supplied JSON object into the target ClassID
// The class id object must have the following shape:
//
//	{
//	  "type": "<CLASS_ID_TYPE>",
//	  "value": <CLASS_ID_VALUE>
//	}
//
// where <CLASS_ID_TYPE> must be one of the known IClassIDValue implementation
// type names (available in this implementation: "uuid", "oid", "int", "bytes"),
// and <CLASS_ID_VALUE> is the JSON encoding of the underlying class id value.
// The exact encoding is <CLASS_ID_TYPE> dependent.
//		oid: dot-separated integers, e.g. "1.2.3.4"
//		uuid: standard UUID string representation, e.g. "550e8400-e29b-41d4-a716-446655440000"
//		int: an integer value, e.g. 7
//	 bytes: a variable length opaque bytes, example {0x07, 0x12, 0x34}

//nolint:dupl
func (o *ClassID) UnmarshalJSON(data []byte) error {
	var tnv encoding.TypeAndValue

	if err := json.Unmarshal(data, &tnv); err != nil {
		return fmt.Errorf("class id decoding failure: %w", err)
	}

	decoded, err := NewClassID(nil, tnv.Type)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(tnv.Value, &decoded.Value); err != nil {
		return fmt.Errorf(
			"cannot unmarshal class id: %w",
			err,
		)
	}

	if err := decoded.Value.Valid(); err != nil {
		return fmt.Errorf("invalid %s: %w", tnv.Type, err)
	}

	o.Value = decoded.Value

	return nil
}

// MarshalJSON serializes the target ClassID to JSON
func (o ClassID) MarshalJSON() ([]byte, error) {
	return extensions.TypeChoiceValueMarshalJSON(o.Value)
}

// String returns a printable string of the ClassID value. UUIDs use the
// canonical 8-4-4-4-12 format, OIDs are output in dotted-decimal notation.
func (o ClassID) String() string {
	if o.Value == nil {
		return ""
	}

	return o.Value.String()
}

type IClassIDValue interface {
	extensions.ITypeChoiceValue

	Bytes() []byte
}

// SetUUID sets the value of the target ClassID to the supplied UUID
func (o *ClassID) SetUUID(uuid UUID) *ClassID {
	if o != nil {
		o.Value = TaggedUUID(uuid)
	}
	return o
}

func (o ClassID) GetUUID() (UUID, error) {
	switch t := o.Value.(type) {
	case *TaggedUUID:
		return UUID(*t), nil
	case TaggedUUID:
		return UUID(t), nil
	default:
		return UUID{}, fmt.Errorf("class-id type is: %T", t)
	}
}

// SetOID sets the value of the target ClassID to the supplied OID.
// The OID is a string in dotted-decimal notation
func (o *ClassID) SetOID(s string) error {
	if o != nil {
		var berOID OID
		if err := berOID.FromString(s); err != nil {
			return err
		}
		o.Value = TaggedOID(berOID)
	}
	return nil
}

// GetOID gets the value of the OID in a string dotted-decimal notation
func (o ClassID) GetOID() (string, error) {
	switch t := o.Value.(type) {
	case *TaggedOID:
		oid := OID(*t)
		stringOID := oid.String()
		return stringOID, nil
	case TaggedOID:
		oid := OID(t)
		stringOID := oid.String()
		return stringOID, nil
	default:
		return "", fmt.Errorf("class-id type is: %T", t)
	}
}

func NewOIDClassID(val any) (*ClassID, error) {
	ret, err := NewTaggedOID(val)
	if err != nil {
		return nil, err
	}

	return &ClassID{ret}, nil
}

func MustNewOIDClassID(val any) *ClassID {
	ret, err := NewOIDClassID(val)
	if err != nil {
		panic(err)
	}

	return ret
}

func NewUUIDClassID(val any) (*ClassID, error) {
	if val == nil {
		return &ClassID{&TaggedUUID{}}, nil
	}

	ret, err := NewTaggedUUID(val)
	if err != nil {
		return nil, err
	}

	return &ClassID{ret}, nil
}

func MustNewUUIDClassID(val any) *ClassID {
	ret, err := NewUUIDClassID(val)
	if err != nil {
		panic(err)
	}

	return ret
}

const IntType = "int"

type TaggedInt int

func NewIntClassID(val any) (*ClassID, error) {
	if val == nil {
		zeroVal := TaggedInt(0)
		return &ClassID{&zeroVal}, nil
	}

	var ret TaggedInt

	switch t := val.(type) {
	case string:
		i, err := strconv.Atoi(t)
		if err != nil {
			return nil, fmt.Errorf("bad int: %w", err)
		}
		ret = TaggedInt(i)
	case []byte:
		if len(t) != 8 {
			return nil, fmt.Errorf("bad int: want 8 bytes, got %d bytes", len(t))
		}
		ti, err := safecast.Convert[int, uint64](binary.BigEndian.Uint64(t))
		if err != nil {
			return nil, err
		}
		ret = TaggedInt(ti)
	case int:
		ret = TaggedInt(t)
	case *int:
		ret = TaggedInt(*t)
	case int64:
		ti, err := safecast.Convert[int, int64](t)
		if err != nil {
			return nil, err
		}
		ret = TaggedInt(ti)
	case *int64:
		ti, err := safecast.Convert[int, int64](*t)
		if err != nil {
			return nil, err
		}
		ret = TaggedInt(ti)
	case uint64:
		ti, err := safecast.Convert[int, uint64](t)
		if err != nil {
			return nil, err
		}
		ret = TaggedInt(ti)
	case *uint64:
		ti, err := safecast.Convert[int, uint64](*t)
		if err != nil {
			return nil, err
		}
		ret = TaggedInt(ti)
	default:
		return nil, fmt.Errorf("unexpected type for int: %T", t)
	}

	if err := ret.Valid(); err != nil {
		return nil, err
	}

	return &ClassID{&ret}, nil
}

func (o TaggedInt) String() string {
	return fmt.Sprint(int(o))
}

func (o TaggedInt) Valid() error {
	return nil
}

func (o TaggedInt) Type() string {
	return "int"
}

func (o TaggedInt) Bytes() []byte {
	var ret [8]byte
	var uo uint64
	io := int(o) // Needed for gosec flow typing.
	// Use 2's complement for negative values since this can't return an error.
	if io < 0 {
		uo = ^uint64(0) - uint64(-io) + 1
	} else {
		uo = uint64(io)
	}
	binary.BigEndian.PutUint64(ret[:], uo)
	return ret[:]
}

// MustNewBytesClassID is like NewBytesClassID except it does not return an
// error, assuming that the provided value is valid. It panics if that isn't
// the case.
func MustNewBytesClassID(val any) *ClassID {
	ret, err := NewBytesClassID(val)
	if err != nil {
		panic(err)
	}

	return ret
}

// NewBytesClassID creates a New ClassID of type bytes
// The supplied interface parameter could be
// a byte slice, a pointer to a byte slice or a string
func NewBytesClassID(val any) (*ClassID, error) {
	ret, err := NewBytes(val)
	if err != nil {
		return nil, err
	}
	return &ClassID{ret}, nil
}

// IClassIDFactory defines the signature for the factory functions that may be
// registred using RegisterClassIDType to provide a new implementation of the
// corresponding type choice. The factory function should create a new *ClassID
// with the underlying value created based on the provided input. The range of
// valid inputs is up to the specific type choice implementation, however it
// _must_ accept nil as one of the inputs, and return the Zero value for
// implemented type.
// See also https://go.dev/ref/spec#The_zero_value
type IClassIDFactory func(any) (*ClassID, error)

var classIDValueRegister = map[string]IClassIDFactory{
	OIDType:   NewOIDClassID,
	UUIDType:  NewUUIDClassID,
	IntType:   NewIntClassID,
	BytesType: NewBytesClassID,
}

// RegisterClassIDType registers a new IClassIDValue implementation (created
// by the provided IClassIDFactory) under the specified CBOR tag.
func RegisterClassIDType(tag uint64, factory IClassIDFactory) error {
	nilVal, err := factory(nil)
	if err != nil {
		return err
	}

	typ := nilVal.Type()
	if _, exists := classIDValueRegister[typ]; exists {
		return fmt.Errorf("class ID type with name %q already exists", typ)
	}

	if err := registerCOMIDTag(tag, nilVal.Value); err != nil {
		return err
	}

	classIDValueRegister[typ] = factory

	return nil
}
