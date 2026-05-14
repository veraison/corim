// Copyright 2021-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
)

const MaskedType = "masked"

// IRawValueValue is the interface implemented by concrete RawValue value
// types.
type IRawValueValue interface {
	extensions.ITypeChoiceValue

	Bytes() []byte
}

// NewRawValue returns the pointer to a new RawValue of the specified type,
// constructed using the provided value b. The type of b depends on the
// specified raw value type. For a bytes value, it must be a []byte, for
// masked bytes, it must be a [2][]byte, [][]byte of length 2, or a []byte (in
// which case the mast will be a value of the same length with all bits set).
func NewRawValue(b any, typ string) (*RawValue, error) {
	factory, ok := rawValueValueRegister[typ]
	if !ok {
		return nil, fmt.Errorf("unexpected RawValue type: %s", typ)
	}

	return factory(b)
}

// MustNewRawValue is the same as NewRawValues, but it panics on error.
func MustNewRawValue(b any, typ string) *RawValue {
	ret, err := NewRawValue(b, typ)
	if err != nil {
		panic(err)
	}

	return ret
}

// NewRawValuesFromBytes returns a pointer to a new RawValue with the specified
// bytes as the underlying value.
func NewRawValueFromBytes(b []byte) *RawValue {
	ret, err := NewBytesRawValue(b)
	if err != nil {
		// cannot happen
		panic(err)
	}

	return ret
}

// NewRawValueWithMask returns a pointer to a new RawValues with a masked
// underlying values created from the specified value and mask. An error is
// returned if the length of the mask does not match that of the value.
func NewRawValueWithMask(val []byte, mask []byte) (*RawValue, error) {
	return NewMaskedRawValue([2][]byte{val, mask})
}

// MustNewRawValueWithMask is like NewRawValueWithMask but panics if the value
// and mask are of different lengths.
func MustNewRawValueWithMask(val []byte, mask []byte) *RawValue {
	ret, err := NewRawValueWithMask(val, mask)
	if err != nil {
		panic(err)
	}

	return ret
}

// RawValue models a $raw-value-type-choice.  For now, the only available type is bytes.
type RawValue struct {
	Value IRawValueValue
}

// Type returns the type of the underlying value
func (o RawValue) Type() string {
	return o.Value.Type()
}

// Bytes returns the bytes of the underlying value
func (o RawValue) Bytes() []byte {
	return o.Value.Bytes()
}

// Mask returns the mask of the underlying value. If the underlying value is
// not masked, nil is returned.
func (o RawValue) Mask() []byte {
	masked, ok := o.Value.(interface{ Mask() []byte })
	if ok {
		return masked.Mask()
	}

	return nil
}

// Equal returns true if both raw values are equal, or false if they are not or
// of one or more of RawValues is invalid. If both values specify masks and
// masks differ, false is returned (even if bytes masked with respective masks
// equal).
func (o RawValue) Equal(other *RawValue) bool {
	ret, _ := maskedEqual(o.Bytes(), o.Mask(), other.Bytes(), other.Mask())
	return ret
}

// CompareAgainstReference checks if a RawValue object matches with a reference
//
// (draft-ietf-rats-corim §9.4.6.1.4).
func (o RawValue) CompareAgainstReference(ref []byte, mask []byte) bool {
	ret, _ := maskedEqual(o.Bytes(), o.Mask(), ref, mask)
	return ret
}

func (o RawValue) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.Value)
}

func (o *RawValue) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &o.Value)
}

func (o RawValue) MarshalJSON() ([]byte, error) {
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

func (o *RawValue) UnmarshalJSON(data []byte) error {
	var value encoding.TypeAndValue
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	decoded, err := NewRawValue(nil, value.Type)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(value.Value, &decoded.Value); err != nil {
		return err
	}

	o.Value = decoded.Value

	return nil
}

// NewBytesRawValue creates a new RawValue with the underlying value being
// TaggedBytes constructed from the provided input.
func NewBytesRawValue(val any) (*RawValue, error) {
	if val == nil {
		return &RawValue{&TaggedBytes{}}, nil
	}

	switch t := val.(type) {
	case []byte:
		return &RawValue{(*TaggedBytes)(&t)}, nil
	case TaggedBytes:
		return &RawValue{&t}, nil
	case *TaggedBytes:
		return &RawValue{t}, nil
	default:
		return nil, fmt.Errorf("value must be a byte slice; found %T", t)
	}
}

// TaggedMaskedRawValue is a RawValue type that combines a value with a mask.
// The mask must be of the same length as the value. When this is being
// compared with an unmasked RawValue, only the masked bits of both underlying
// values are compared.
type TaggedMaskedRawValue struct {
	_         struct{} `cbor:",toarray"`
	Value     []byte   `json:"value"`
	MaskBytes []byte   `json:"mask"`
}

func (o TaggedMaskedRawValue) Type() string {
	return MaskedType
}

func (o TaggedMaskedRawValue) Bytes() []byte {
	return o.Value
}

func (o TaggedMaskedRawValue) Mask() []byte {
	return o.MaskBytes
}

func (o TaggedMaskedRawValue) String() string {
	return string(o.Value)
}

func (o TaggedMaskedRawValue) Valid() error {
	if len(o.Value) == 0 {
		return errors.New("value not set")
	}

	if len(o.MaskBytes) != len(o.Value) {
		return errors.New("mask and value lengths differ")
	}

	return nil
}

// NewMaskedRawValue returns a new RawValue with the underlying value and mask
// constructed from the provided input. If the input is a [2][]byte or a
// [][]byte of length 2, then the first element is used as the value, and the
// second element as the mask. If the input is a []byte, then that is used as
// the value and the mask is a []byte of the same length with all bits set.
func NewMaskedRawValue(val any) (*RawValue, error) {
	if val == nil {
		return &RawValue{&TaggedMaskedRawValue{}}, nil
	}

	switch t := val.(type) {
	case []byte:
		return &RawValue{&TaggedMaskedRawValue{Value: t, MaskBytes: allSet(len(t))}}, nil
	case TaggedBytes:
		return &RawValue{&TaggedMaskedRawValue{Value: t, MaskBytes: allSet(len(t))}}, nil
	case *TaggedBytes:
		return &RawValue{&TaggedMaskedRawValue{Value: *t, MaskBytes: allSet(len(*t))}}, nil
	case [2][]byte:
		return &RawValue{&TaggedMaskedRawValue{Value: t[0], MaskBytes: t[1]}}, nil
	case [][]byte:
		if len(t) != 2 {
			return nil, errors.New("[][]byte must contain exactly two elements")
		}
		return &RawValue{&TaggedMaskedRawValue{Value: t[0], MaskBytes: t[1]}}, nil
	default:
		return nil, fmt.Errorf("value must be a byte slice; found %T", t)
	}
}

// MustNewMaskedRawValue is like NewMaskedRawValue but panics on error.
func MustNewMaskedRawValue(val any) *RawValue {
	ret, err := NewMaskedRawValue(val)
	if err != nil {
		panic(err)
	}

	return ret
}

// IRawValueFactory defines the signature for the factory functions that may be
// registred using RegisterRawValueType to provide a new implementation of the
// corresponding type choice. The factory function should create a new
// *RawValue with the underlying value created based on the provided input. The
// range of valid inputs is up to the specific type choice implementation,
// however it _must_ accept nil as one of the inputs, and return the Zero value
// for the implemented type.
// See also https://go.dev/ref/spec#The_zero_value
type IRawValueFactory func(any) (*RawValue, error)

var rawValueValueRegister = map[string]IRawValueFactory{
	BytesType:  NewBytesRawValue,
	MaskedType: NewMaskedRawValue,
}

// RegisterRawValueType registers a new IRawValueValue implementation
// (created by the provided IRawValueFactory) under the specified type name
// and CBOR tag.
func RegisterRawValueType(tag uint64, factory IRawValueFactory) error {
	nilVal, err := factory(nil)
	if err != nil {
		return err
	}

	typ := nilVal.Type()
	if _, exists := rawValueValueRegister[typ]; exists {
		return fmt.Errorf("raw value type with name %q already exists", typ)
	}

	if err := registerCOMIDTag(tag, nilVal.Value); err != nil {
		return err
	}

	rawValueValueRegister[typ] = factory

	return nil
}

func allSet(n int) []byte {
	ret := make([]byte, n)
	for i := range ret {
		ret[i] = 0xff
	}

	return ret
}

func maskedEqual(lhs []byte, lhsMask []byte, rhs []byte, rhsMask []byte) (bool, error) {
	if lhsMask != nil && rhsMask != nil && !bytes.Equal(lhsMask, rhsMask) {
		return false, errors.New("LHS and RHS masks are both set and are unequal")
	}

	if lhsMask != nil && len(lhsMask) != len(lhs) {
		return false, errors.New("LHS mask and value lengths differ")
	}

	if rhsMask != nil && len(rhsMask) != len(rhs) {
		return false, errors.New("RHS mask and value lengths differ")
	}

	if len(lhs) != len(rhs) {
		return false, nil
	}

	mask := lhsMask
	if rhsMask != nil {
		mask = rhsMask
	}

	if mask != nil {
		lhs = applyMask(lhs, mask)
		rhs = applyMask(rhs, mask)
	}

	return bytes.Equal(lhs, rhs), nil
}

func applyMask(val, mask []byte) []byte {
	ret := make([]byte, len(val))
	for i, b := range val {
		ret[i] = b & mask[i]
	}

	return ret
}
