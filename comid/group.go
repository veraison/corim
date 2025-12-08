// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
)

// Group stores a group identity. The supported formats are UUID and variable-length opaque bytes.
type Group struct {
	Value IGroupValue
}

// NewGroup instantiates an empty group
func NewGroup(val any, typ string) (*Group, error) {
	factory, ok := groupValueRegister[typ]
	if !ok {
		return nil, fmt.Errorf("unknown group type: %s", typ)
	}

	return factory(val)
}

// Valid checks for the validity of given group
func (o Group) Valid() error {
	if o.Value == nil {
		return errors.New("no value set")
	}

	return o.Value.Valid()
}

// String returns a printable string of the Group value.  UUIDs use the
// canonical 8-4-4-4-12 format, UEIDs are hex encoded.
func (o Group) String() string {
	return o.Value.String()
}

// Type returns the type of the Group
func (o Group) Type() string {
	if o.Value == nil {
		return ""
	}

	return o.Value.Type()
}

// Bytes returns a []byte containing the raw bytes of the group value
func (o Group) Bytes() []byte {
	if o.Value == nil {
		return []byte{}
	}
	return o.Value.Bytes()
}

// MarshalCBOR serializes the target group to CBOR
func (o Group) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.Value)
}

// UnmarshalCBOR deserializes the supplied CBOR into the target group
func (o *Group) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &o.Value)
}

// UnmarshalJSON deserializes the supplied JSON type/value object into the Group
// target.  The following formats are supported:
//
//	 (a) UUID, e.g.:
//		{
//		  "type": "uuid",
//		  "value": "69E027B2-7157-4758-BCB4-D9F167FE49EA"
//		}
//
// (b) Tagged bytes, e.g. :
//
//	{
//	  "type": "bytes",
//	  "value": "MTIzNDU2Nzg5"
//	}

//nolint:dupl
func (o *Group) UnmarshalJSON(data []byte) error {
	var tnv encoding.TypeAndValue

	if err := json.Unmarshal(data, &tnv); err != nil {
		return fmt.Errorf("group decoding failure: %w", err)
	}

	decoded, err := NewGroup(nil, tnv.Type)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(tnv.Value, &decoded.Value); err != nil {
		return fmt.Errorf(
			"cannot unmarshal group: %w",
			err,
		)
	}

	if err := decoded.Value.Valid(); err != nil {
		return fmt.Errorf("invalid %s: %w", tnv.Type, err)
	}

	o.Value = decoded.Value

	return nil
}

func (o Group) MarshalJSON() ([]byte, error) {
	return extensions.TypeChoiceValueMarshalJSON(o.Value)
}

type IGroupValue interface {
	extensions.ITypeChoiceValue

	Bytes() []byte
}

func NewUUIDGroup(val any) (*Group, error) {
	if val == nil {
		return &Group{&TaggedUUID{}}, nil
	}

	u, err := NewTaggedUUID(val)
	if err != nil {
		return nil, err
	}

	return &Group{u}, nil
}

func MustNewUUIDGroup(val any) *Group {
	ret, err := NewUUIDGroup(val)
	if err != nil {
		panic(err)
	}

	return ret
}

// MustNewBytesGroup is like NewBytesGroup except it does not return an
// error, assuming that the provided value is valid. It panics if that isn't
// the case.
func MustNewBytesGroup(val any) *Group {
	ret, err := NewBytesGroup(val)
	if err != nil {
		panic(err)
	}

	return ret
}

// NewBytesGroup creates a New Group of type bytes
// The supplied interface parameter could be
// a byte slice, a pointer to a byte slice or a string
func NewBytesGroup(val any) (*Group, error) {
	ret, err := NewBytes(val)
	if err != nil {
		return nil, err
	}
	return &Group{ret}, nil
}

// IGroupFactory defines the signature for the factory functions that may be
// registered using RegisterGroupType to provide a new implementation of the
// corresponding type choice. The factory function should create a new *Group
// with the underlying value created based on the provided input. The range of
// valid inputs is up to the specific type choice implementation, however it
// _must_ accept nil as one of the inputs, and return the Zero value for
// implemented type.
// See also https://go.dev/ref/spec#The_zero_value
type IGroupFactory func(any) (*Group, error)

var groupValueRegister = map[string]IGroupFactory{
	UUIDType:  NewUUIDGroup,
	BytesType: NewBytesGroup,
}

// RegisterGroupType registers a new IGroupValue implementation
// (created by the provided IGroupFactory) under the specified type name
// and CBOR tag.
func RegisterGroupType(tag uint64, factory IGroupFactory) error {

	nilVal, err := factory(nil)
	if err != nil {
		return err
	}

	typ := nilVal.Value.Type()
	if _, exists := groupValueRegister[typ]; exists {
		return fmt.Errorf("Group type with name %q already exists", typ)
	}

	if err := registerCOMIDTag(tag, nilVal.Value); err != nil {
		return err
	}

	groupValueRegister[typ] = factory

	return nil
}
