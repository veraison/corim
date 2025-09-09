// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"encoding/json"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/corim/encoding"
)

// NumericType represents the abstraction of a numeric type, allowed values are int/uint/float
type NumericType struct {
	val interface{}
}

// NewNumericType creates a new NumericType from the
// supplied interface. The supported types are unsigned integers, signed integers and
// float64
func NewNumericType(val any) (*NumericType, error) {
	switch t := val.(type) {
	case uint, uint64:
		return &NumericType{val: t}, nil
	case float64:
		return &NumericType{val: t}, nil
	case int:
		return &NumericType{val: t}, nil
	default:
		return nil, fmt.Errorf("unsupported NumericType type: %T", t)
	}
}

func (o *NumericType) SetNumericType(val any) error {
	switch t := val.(type) {
	case uint, uint64:
		o.val = val
	case float64:
		o.val = val
	case int:
		o.val = val
	default:
		return fmt.Errorf("unsupported NumericType type: %T", t)
	}
	return nil
}

// IsFloat returns true if NumericType is of type float64 array
func (o NumericType) IsFloat() bool {
	return isType[float64](o.val)
}

// IsUnit returns true if NumericType is of type unsigned integer
func (o NumericType) IsUint() bool {
	return isType[uint64](o.val) || isType[uint](o.val)
}

// IsInt returns true if NumericType is of type integer
func (o NumericType) IsInt() bool {
	return isType[int](o.val)
}

// GetUint returns unsigned integer NumericType
func (o NumericType) GetUint() (uint, error) {
	switch t := o.val.(type) {
	case uint64:
		return uint(t), nil
	case uint:
		return t, nil
	default:
		return 0, fmt.Errorf("NumericType type is: %T", t)
	}
}

// GetInt returns integer NumericType
func (o NumericType) GetInt() (int, error) {
	switch t := o.val.(type) {
	case int:
		return t, nil
	default:
		return 0, fmt.Errorf("NumericType type is: %T", t)
	}
}

// GetFloat returns float NumericType
func (o NumericType) GetFloat() (float64, error) {
	switch t := o.val.(type) {
	case float64:
		return t, nil
	default:
		return 0, fmt.Errorf("NumericType type is: %T", t)
	}
}

// MarshalCBOR Marshals NumericType to CBOR bytes
func (o NumericType) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(o.val)
}

// UnmarshalCBOR UnMarshals supplied CBOR bytes to NumericType
func (o *NumericType) UnmarshalCBOR(data []byte) error {
	return cbor.Unmarshal(data, &o.val)
}

// Valid checks for validity of NumericType and returns an error if Invalid
func (o NumericType) Valid() error {
	if o.val == nil {
		return fmt.Errorf("empty NumericType")
	}
	switch t := o.val.(type) {
	case uint, uint64, float64, int:
		return nil
	default:
		return fmt.Errorf("unsupported NumericType type: %T", t)
	}
}

// MarshalJSON Marshals NumericType to JSON
func (o NumericType) MarshalJSON() ([]byte, error) {
	if o.Valid() != nil {
		return nil, fmt.Errorf("invalid NumericType")
	}
	var (
		v   encoding.TypeAndValue
		b   []byte
		err error
	)
	switch t := o.val.(type) {
	case uint, uint64:
		b, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		v = encoding.TypeAndValue{Type: UintType, Value: b}
	case int:
		b, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		v = encoding.TypeAndValue{Type: IntType, Value: b}
	case float64:
		b, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		v = encoding.TypeAndValue{Type: FloatType, Value: b}
	default:
		return nil, fmt.Errorf("unknown type %T for NumericType", t)
	}
	return json.Marshal(v)
}

// UnmarshalJSON UnMarshals supplied JSON buffer to NumericType
func (o *NumericType) UnmarshalJSON(data []byte) error {
	var v encoding.TypeAndValue

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Type {
	case UintType:
		var x uint
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal NumericType of type uint: %w", err)
		}
		o.val = x
	case IntType:
		var x int
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal NumericType of type int: %w", err)
		}
		o.val = x
	case FloatType:
		var x float64
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal NumericType of type float: %w", err)
		}
		o.val = x
	}
	return nil
}
