// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
)

// SVN is the Security Version Number. This typically changes only when a
// security relevant change is needed to the measured environment.
type SVN struct {
	Value ISVNValue
}

// NewSVN creates a new SVN of the specified and value. The type must be one of
// the strings defined by the spec ("exact-value", "min-value"), or has been
// registered with RegisterSVNType().
func NewSVN(val any, typ string) (*SVN, error) {
	factory, ok := svnValueRegister[typ]
	if !ok {
		return nil, fmt.Errorf("unknown SVN type: %s", typ)
	}

	return factory(val)
}

// MustNewSVN is like NewSVN but does not return an error, assuming that the
// provided value is valid. It panics if this is not the case.
func MustNewSVN(val any, typ string) *SVN {
	ret, err := NewSVN(val, typ)
	if err != nil {
		panic(err)
	}

	return ret
}

// MarshalCBOR returns the CBOR encoding of the SVN.
func (o SVN) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.Value)
}

// UnmarshalCBOR populates the SVN form the provided CBOR bytes.
func (o *SVN) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &o.Value)
}

// UnmarshalJSON deserializes the supplied JSON object into the target SVN
// The SVN object must have the following shape:
//
//	{
//	  "type": "<SVN_TYPE>",
//	  "value": <SVN_VALUE>
//	}
//
// where <SVN_TYPE> must be one of the known ISVNValue implementation
// type names (available in the base implementation: "exact-value",
// "min-value"), and <SVN_VALUE> is the JSON encoding of the underlying
// class id value. The exact encoding is <SVN_TYPE> dependent. For both base
// types, it is an integer (JSON number).
func (o *SVN) UnmarshalJSON(data []byte) error {
	var tnv encoding.TypeAndValue

	if err := json.Unmarshal(data, &tnv); err != nil {
		return fmt.Errorf("SVN decoding failure: %w", err)
	}

	decoded, err := NewSVN(nil, tnv.Type)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(tnv.Value, &decoded.Value); err != nil {
		return fmt.Errorf("invalid SVN %s: %w", tnv.Type, err)
	}

	if err := decoded.Value.Valid(); err != nil {
		return fmt.Errorf("invalid SVN %s: %w", tnv.Type, err)
	}

	o.Value = decoded.Value

	return nil
}

// MarshalJSON serializes the SVN int a JSON object
func (o SVN) MarshalJSON() ([]byte, error) {
	return extensions.TypeChoiceValueMarshalJSON(o.Value)
}

// ISVNValue is the interface that must be implemented by all SVN values.
type ISVNValue interface {
	extensions.ITypeChoiceValue
}

const (
	ExactValueType = "exact-value"
	MinValueType   = "min-value"
)

type TaggedSVN uint64

func NewTaggedSVN(val any) (*SVN, error) {
	var ret TaggedSVN

	if val == nil {
		return &SVN{&ret}, nil
	}

	u, err := convertToSVNUint64(val)
	if err != nil {
		return nil, err
	}
	ret = TaggedSVN(u)

	return &SVN{&ret}, nil
}

func MustNewTaggedSVN(val any) *SVN {
	ret, err := NewTaggedSVN(val)
	if err != nil {
		panic(err)
	}

	return ret
}

func (o TaggedSVN) String() string {
	return fmt.Sprint(uint64(o))
}

func (o TaggedSVN) Type() string {
	return ExactValueType
}

func (o TaggedSVN) Valid() error {
	return nil
}

type TaggedMinSVN uint64

func NewTaggedMinSVN(val any) (*SVN, error) {
	var ret TaggedMinSVN

	if val == nil {
		return &SVN{&ret}, nil
	}

	u, err := convertToSVNUint64(val)
	if err != nil {
		return nil, err
	}
	ret = TaggedMinSVN(u)

	return &SVN{&ret}, nil
}

func MustNewTaggedMinSVN(val any) *SVN {
	ret, err := NewTaggedMinSVN(val)
	if err != nil {
		panic(err)
	}

	return ret
}

func (o TaggedMinSVN) String() string {
	return fmt.Sprint(uint64(o))
}

func (o TaggedMinSVN) Type() string {
	return MinValueType
}

func (o TaggedMinSVN) Valid() error {
	return nil
}

// This function converts various types to uint64 for SVN.
// convertToSVNUint64 converts various types to uint64 for SVN
func convertToSVNUint64(val any) (uint64, error) {
    switch t := val.(type) {
    case string:
        u, err := strconv.ParseUint(t, 10, 64)
        if err != nil {
            return 0, err
        }
        return u, nil
    case uint64:
        return t, nil
    case uint:
        return uint64(t), nil
    case int:
        if t < 0 {
            return 0, fmt.Errorf("SVN cannot be negative: %d", t)
        }
        return uint64(t), nil
    case int64:
        if t < 0 {
            return 0, fmt.Errorf("SVN cannot be negative: %d", t)
        }
        return uint64(t), nil
    case TaggedSVN:
        return uint64(t), nil
    case *TaggedSVN:
        return uint64(*t), nil
    case TaggedMinSVN:
        return uint64(t), nil
    case *TaggedMinSVN:
        return uint64(*t), nil
    default:
        return 0, fmt.Errorf("unexpected type for SVN: %T", t)
    }
}

// ISVNFactory defines the signature for the factory functions that may be
// registred using RegisterSVNType to provide a new implementation of the
// corresponding type choice. The factory function should create a new *SVN
// with the underlying value created based on the provided input. The range of
// valid inputs is up to the specific type choice implementation, however it
// _must_ accept nil as one of the inputs, and return the Zero value for
// implemented type.
// See also https://go.dev/ref/spec#The_zero_value
type ISVNFactory func(any) (*SVN, error)

var svnValueRegister = map[string]ISVNFactory{
	ExactValueType: NewTaggedSVN,
	MinValueType:   NewTaggedMinSVN,
}

// RegisterSVNType registers a new ISVNValue implementation
// (created by the provided ISVNFactory) under the specified CBOR tag.
func RegisterSVNType(tag uint64, factory ISVNFactory) error {

	nilVal, err := factory(nil)
	if err != nil {
		return err
	}

	typ := nilVal.Value.Type()
	if _, exists := svnValueRegister[typ]; exists {
		return fmt.Errorf("SVN type with name %q already exists", typ)
	}

	if err := registerCOMIDTag(tag, nilVal.Value); err != nil {
		return err
	}

	svnValueRegister[typ] = factory

	return nil
}
