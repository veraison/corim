// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
)

// RawInt describes an integer value that can be compared with linear order in
// the target environment. It follows the type-choice pattern  widely
// used in CoRIM and implements the extensions.ITypeChoiceValue interface
// https://ietf-rats-wg.github.io/draft-ietf-rats-corim/draft-ietf-rats-corim.html#name-raw-int
type RawInt struct {
	Value extensions.ITypeChoiceValue
}

// NewRawInt returns a *RawInt of the specified type
func NewRawInt(val any, typ string) (*RawInt, error) {
	factory, ok := rawIntValueRegister[typ]
	if !ok {
		return nil, fmt.Errorf("unknown type: %s", typ)
	}

	return factory(val)
}

// IsSet confirms if RawInt has a value or if it's empty
func (o RawInt) IsSet() bool { return o.Value != nil }

// Type returns the type of RawInt
func (o RawInt) Type() string {
	if o.IsSet() {
		return o.Value.Type()
	}

	return ""
}

// Valid checks if the RawInt is valid
func (o RawInt) Valid() error {
	if !o.IsSet() {
		return errors.New("RawInt value unset")
	}

	return o.Value.Valid()
}

// String returns the RawInt in a string format
func (o RawInt) String() string {
	if !o.IsSet() {
		return ""
	}

	return o.Value.String()
}

// MarshalJSON serializes RawInt into JSON
func (o RawInt) MarshalJSON() ([]byte, error) {
	valueBytes, err := json.Marshal(o.Value)
	if err != nil {
		return nil, err
	}

	value := encoding.TypeAndValue{Type: o.Value.Type(), Value: valueBytes}

	return json.Marshal(value)
}

// UnmarshalJSON de-serializes input JSON data into RawInt
func (o *RawInt) UnmarshalJSON(data []byte) error {
	var tnv encoding.TypeAndValue

	if err := json.Unmarshal(data, &tnv); err != nil {
		return err
	}

	decoded, err := NewRawInt(nil, tnv.Type)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(tnv.Value, decoded.Value); err != nil {
		return err
	}

	if err := decoded.Value.Valid(); err != nil {
		return fmt.Errorf("invalid %s: %w", tnv.Type, err)
	}

	o.Value = decoded.Value

	return nil
}

// MarshalCBOR serializes RawInt to CBOR
func (o RawInt) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.Value)
}

// UnmarshalCBOR de-serializes input CBOR into RawInt
func (o *RawInt) UnmarshalCBOR(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty input")
	}

	majorType := (data[0] & 0xe0) >> 5
	switch majorType {
	case 0, 1:
		rawIntInteger := new(RawIntInteger)
		if err := dm.Unmarshal(data, rawIntInteger); err != nil {
			return err
		}
		o.Value = rawIntInteger
	case 6:
		rawIntRange := new(TaggedRawIntRange)
		if err := dm.Unmarshal(data, rawIntRange); err != nil {
			return err
		}
		o.Value = rawIntRange
	default:
		return fmt.Errorf("RawInt: Error: unknown major type: %d", majorType)
	}

	return nil
}

// RawIntIntegerType is the name of the Integer version of RawInt
const RawIntIntegerType = "rawIntInteger"

// RawIntInteger implements the Integer version of RawInt. It can
// hold both positive and negative integer values
type RawIntInteger int64

// NewRawIntInteger creates a RawIntInteger from the given input value
func NewRawIntInteger(val any) (*RawIntInteger, error) {
	var ret RawIntInteger

	if val == nil {
		return &ret, nil
	}

	switch v := val.(type) {
	case RawIntInteger:
		ret = v
	case *RawIntInteger:
		ret = *v
	case int64:
		ret = RawIntInteger(v)
	default:
		return nil, fmt.Errorf("unexpected type for RawIntInteger: %T", v)
	}

	return &ret, nil
}

// Valid is a no-op for RawIntInteger and always returns nil
func (o RawIntInteger) Valid() error { return nil }

// String converts RawIntInteger into a string
func (o RawIntInteger) String() string { return strconv.FormatInt(int64(o), 10) }

// Type returns the name/type of RawIntInteger
func (o RawIntInteger) Type() string {
	return RawIntIntegerType
}

// CompareAgainstRefInteger compares RawIntInteger against another RawIntInteger reference value
func (o RawIntInteger) CompareAgainstRefInteger(ref RawIntInteger) bool {
	return o == ref
}

// CompareAgainstRefRange compares RawIntInteger against a TaggedRawIntRange
func (o RawIntInteger) CompareAgainstRefRange(ref TaggedRawIntRange) bool {
	obj, err := convertRawIntIntegerToInt64(o)
	if err != nil {
		log.Printf("RawIntInteger:CompareAgainstRefRange: Error: %v", err)
		return false
	}

	if ref.Min != nil && obj < *ref.Min {
		return false
	}

	if ref.Max != nil && obj > *ref.Max {
		return false
	}

	return true
}

func convertRawIntIntegerToInt64(val any) (int64, error) {
	switch v := val.(type) {
	case RawIntInteger:
		return int64(v), nil
	case *RawIntInteger:
		return int64(*v), nil
	default:
		return 0, fmt.Errorf("unexpected type for RawIntInteger: %T", v)
	}
}

// TaggedRawIntRangeType is the name of the Range version of RawInt
const TaggedRawIntRangeType = "rawIntRange"

// TaggedRawIntRange implements the Range version of RawInt. The range is
// made of minimum and maximum values. If the minimum is nil, it's assumed
// to be negative infinity. If the maximum is nil, it's assumed to be
// positive infinity.
type TaggedRawIntRange struct {
	Min *int64
	Max *int64
}

// NewRawIntRange creates a TaggedRawIntRange with the input value
func NewRawIntRange(val any) (*TaggedRawIntRange, error) {
	var ret TaggedRawIntRange

	if val == nil {
		return &ret, nil
	}

	switch v := val.(type) {
	case TaggedRawIntRange:
		ret = v
	case *TaggedRawIntRange:
		ret = *v
	default:
		return nil, fmt.Errorf("unexpected type for TaggedRawIntRange: %T", v)
	}

	return &ret, nil
}

// Valid checks if TaggedRawIntRange is valid
func (o TaggedRawIntRange) Valid() error {
	if o.Min != nil && o.Max != nil && *o.Min > *o.Max {
		return fmt.Errorf("TaggedRawIntRange: Invalid Range, Min: %d Max: %d", *o.Min, *o.Max)
	}

	return nil
}

// String converts TaggedRawIntRange to a string
func (o TaggedRawIntRange) String() string {
	var rangeMin, rangeMax string

	if o.Min != nil {
		rangeMin = fmt.Sprintf("[%d", *o.Min)
	} else {
		rangeMin = "(-inf"
	}

	if o.Max != nil {
		rangeMax = fmt.Sprintf("%d]", *o.Max)
	} else {
		rangeMax = "inf)"
	}

	return fmt.Sprintf("%s:%s", rangeMin, rangeMax)
}

// Type returns the name/type of TaggedRawIntRange
func (o TaggedRawIntRange) Type() string {
	return TaggedRawIntRangeType
}

// CompareAgainstRefInteger compares TaggedRawIntRange against a given
// RawIntInteger reference
func (o TaggedRawIntRange) CompareAgainstRefInteger(ref RawIntInteger) bool {
	refVal, err := convertRawIntIntegerToInt64(ref)
	if err != nil {
		log.Printf("TaggedRawIntRange:CompareAgainstRefInteger: Error: %v", err)
		return false
	}

	if o.Min != nil && refVal == *o.Min && o.Max != nil && refVal == *o.Max {
		return true
	}

	return false
}

// CompareAgainstRefRange compares TaggedRawIntRange against another
// TaggedRawIntRange reference
func (o TaggedRawIntRange) CompareAgainstRefRange(ref TaggedRawIntRange) bool {
	if o.Min != nil && ref.Min != nil && *o.Min < *ref.Min {
		return false
	}

	if o.Max != nil && ref.Max != nil && *o.Max > *ref.Max {
		return false
	}

	return true
}

// NewRawIntIntegerType returns a new RawIntInteger from given value
func NewRawIntIntegerType(val any) (*RawInt, error) {
	ret, err := NewRawIntInteger(val)
	if err != nil {
		return nil, err
	}

	return &RawInt{ret}, nil
}

// NewRawIntRangeType returns a new TaggedRawIntRange from given value
func NewRawIntRangeType(val any) (*RawInt, error) {
	ret, err := NewRawIntRange(val)
	if err != nil {
		return nil, err
	}

	return &RawInt{ret}, nil
}

// IRawIntFactory type defines a factory pattern for RawInt
type IRawIntFactory = func(val any) (*RawInt, error)

var rawIntValueRegister = map[string]IRawIntFactory{
	RawIntIntegerType:     NewRawIntIntegerType,
	TaggedRawIntRangeType: NewRawIntRangeType,
}
