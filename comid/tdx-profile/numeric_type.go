// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// NumericType represents the abstraction of a numeric type, allowed values are int/uint/float
// Define an interface for numeric types
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
