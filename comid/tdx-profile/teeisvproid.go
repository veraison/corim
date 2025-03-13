// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

// nolint:dupl
package tdx

import (
	"encoding/json"
	"fmt"

	"github.com/veraison/corim/encoding"
)

// TeeISVProdID stores an ISV Product Identifier. The supported formats are uint and variable-length bytes.
type TeeISVProdID struct {
	val interface{}
}

// NewTeeISVProdID creates a new TeeISVProdID from the
// supplied interface and return a pointer to TeeISVProdID
// Supported values are positive integers and byte array
func NewTeeISVProdID(val interface{}) (*TeeISVProdID, error) {
	switch t := val.(type) {
	case uint, uint64:
		return &TeeISVProdID{val: t}, nil
	case []byte:
		return &TeeISVProdID{val: t}, nil
	case int:
		if t < 0 {
			return nil, fmt.Errorf("negative integer %d for TeeISVProdID", t)
		}
		return &TeeISVProdID{val: t}, nil
	default:
		return nil, fmt.Errorf("unsupported TeeISVProdID type: %T", t)
	}
}

// SetTeeISVProdID sets the supplied value of TeeISVProdID from the interface
// Supported values are either positive integers or byte array
func (o *TeeISVProdID) SetTeeISVProdID(val interface{}) error {
	switch t := val.(type) {
	case uint, uint64:
		o.val = val
	case []byte:
		o.val = val
	case int:
		if t < 0 {
			return fmt.Errorf("unsupported negative TeeISVProdID: %d", t)
		}
		o.val = val
	default:
		return fmt.Errorf("unsupported TeeISVProdID type: %T", t)
	}
	return nil
}

// Valid checks for validity of TeeISVProdID and returns an error if Invalid
func (o TeeISVProdID) Valid() error {
	if o.val == nil {
		return fmt.Errorf("empty TeeISVProdID")
	}
	switch t := o.val.(type) {
	case uint, uint64:
		return nil
	case []byte:
		if len(t) == 0 {
			return fmt.Errorf("empty TeeISVProdID")
		}
	case int:
		if t < 0 {
			return fmt.Errorf("unsupported negative TeeISVProdID: %d", t)
		}
	default:
		return fmt.Errorf("unsupported TeeISVProdID type: %T", t)
	}
	return nil
}

// GetUint returns a uint TeeISVProdID
func (o TeeISVProdID) GetUint() (uint, error) {
	switch t := o.val.(type) {
	case uint64:
		return uint(t), nil
	case uint:
		return t, nil
	default:
		return 0, fmt.Errorf("TeeISVProdID type is: %T", t)
	}
}

// GetBytes returns a []byte TeeISVProdID
func (o TeeISVProdID) GetBytes() ([]byte, error) {
	switch t := o.val.(type) {
	case []byte:
		if len(t) == 0 {
			return nil, fmt.Errorf("TeeISVProdID type is of zero length")
		}
		return t, nil
	default:
		return nil, fmt.Errorf("TeeIsvProdID type is: %T", t)
	}
}

// IsBytes returns true if TeeISVProdID is a byte array
func (o TeeISVProdID) IsBytes() bool {
	return isType[[]byte](o.val)
}

// IsUint returns true if TeeISVProdID is a positive integer
func (o TeeISVProdID) IsUint() bool {
	return isType[uint64](o.val) || isType[uint](o.val)
}

// MarshalJSON Marshals TeeISVProdID to JSON
func (o TeeISVProdID) MarshalJSON() ([]byte, error) {
	if o.Valid() != nil {
		return nil, fmt.Errorf("invalid TeeISVProdID")
	}
	var (
		v   encoding.TypeAndValue
		b   []byte
		err error
	)
	switch t := o.val.(type) {
	case uint, uint64, int:
		b, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		v = encoding.TypeAndValue{Type: UintType, Value: b}
	case []byte:
		b, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		v = encoding.TypeAndValue{Type: BytesType, Value: b}
	default:
		return nil, fmt.Errorf("unknown type %T for TeeISVProdID", t)
	}
	return json.Marshal(v)
}

// UnmarshalJSON UnMarshals supplied JSON buffer to TeeISVProdID
func (o *TeeISVProdID) UnmarshalJSON(data []byte) error {
	var v encoding.TypeAndValue

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Type {
	case UintType:
		var x uint
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal TeeISVProdID of type uint: %w", err)
		}
		o.val = x
	case BytesType:
		var x []byte
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal TeeISVProdID of type bytes: %w", err)
		}
		o.val = x
	}
	return nil
}

// MarshalCBOR Marshals TeeISVProdID to CBOR bytes
func (o TeeISVProdID) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

// UnmarshalCBOR UnMarshals supplied CBOR bytes to TeeISVProdID
func (o *TeeISVProdID) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &o.val)
}
