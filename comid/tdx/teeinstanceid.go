// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

// nolint:dupl
package tdx

import (
	"encoding/json"
	"fmt"

	"github.com/veraison/corim/encoding"
)

// TeeInstanceID stores an TEE Instance Identifier. The supported formats are uint and variable-length bytes.
type TeeInstanceID struct {
	val interface{}
}

// NewTeeInstanceID creates a new InstanceID from the
// supplied interface. The supported types are positive integers and
// byte array
func NewTeeInstanceID(val any) (*TeeInstanceID, error) {
	switch t := val.(type) {
	case uint, uint64:
		return &TeeInstanceID{val: t}, nil
	case []byte:
		return &TeeInstanceID{val: t}, nil
	case int:
		if t < 0 {
			return nil, fmt.Errorf("unsupported negative %d TeeInstanceID", t)
		}
		return &TeeInstanceID{val: t}, nil
	default:
		return nil, fmt.Errorf("unsupported TeeInstanceID type: %T", t)
	}
}

// SetTeeInstanceID sets the supplied value of Instance ID
func (o *TeeInstanceID) SetTeeInstanceID(val any) error {
	switch t := val.(type) {
	case uint, uint64:
		o.val = val
	case []byte:
		o.val = val
	case int:
		if t < 0 {
			return fmt.Errorf("unsupported negative TeeInstanceID: %d", t)
		}
		o.val = val
	default:
		return fmt.Errorf("unsupported TeeInstanceID type: %T", t)
	}
	return nil
}

// valid checks for validity of TeeInstanceID and
// returns an error if Invalid
func (o TeeInstanceID) Valid() error {
	if o.val == nil {
		return fmt.Errorf("empty TeeInstanceID")
	}
	switch t := o.val.(type) {
	case uint, uint64:
		return nil
	case []byte:
		if len(t) == 0 {
			return fmt.Errorf("empty TeeInstanceID")
		}
	case int:
		if t < 0 {
			return fmt.Errorf("unsupported negative TeeInstanceID: %d", t)
		}
	default:
		return fmt.Errorf("unsupported TeeInstanceID type: %T", t)
	}
	return nil
}

// GetUint returns unsigned integer TeeInstanceID
func (o TeeInstanceID) GetUint() (uint, error) {
	switch t := o.val.(type) {
	case uint64:
		return uint(t), nil
	case uint:
		return t, nil
	default:
		return 0, fmt.Errorf("TeeInstanceID type is: %T", t)
	}
}

// GetBytes returns the bytes TeeInstanceID
func (o TeeInstanceID) GetBytes() ([]byte, error) {
	switch t := o.val.(type) {
	case []byte:
		if len(t) == 0 {
			return nil, fmt.Errorf("TeeInstanceID type is of zero length")
		}
		return t, nil
	default:
		return nil, fmt.Errorf("TeeInstanceID type is: %T", t)
	}
}

// IsBytes returns true if TeeInstanceID is of type []byte array
func (o TeeInstanceID) IsBytes() bool {
	return isType[[]byte](o.val)
}

// IsUnit returns true if TeeInstanceID is of type unsigned integer
func (o TeeInstanceID) IsUint() bool {
	return isType[uint64](o.val) || isType[uint](o.val)
}

// MarshalJSON Marshals TeeInstanceID to JSON
func (o TeeInstanceID) MarshalJSON() ([]byte, error) {
	if o.Valid() != nil {
		return nil, fmt.Errorf("invalid TeeInstanceID")
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
		v = encoding.TypeAndValue{Type: "uint", Value: b}
	case []byte:
		b, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		v = encoding.TypeAndValue{Type: "bytes", Value: b}
	default:
		return nil, fmt.Errorf("unknown type %T for TeeInstanceID", t)
	}
	return json.Marshal(v)
}

// UnmarshalJSON UnMarshals supplied JSON bytes to TeeInstanceID
func (o *TeeInstanceID) UnmarshalJSON(data []byte) error {
	var v encoding.TypeAndValue

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Type {
	case UintType:
		var x uint
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal TeeInstanceID of type uint: %w", err)
		}
		o.val = x
	case BytesType:
		var x []byte
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal TeeInstanceID of type bytes: %w", err)
		}
		o.val = x
	}
	return nil
}

// MarshalCBOR Marshals TeeInstanceID to CBOR
func (o TeeInstanceID) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

// UnmarshalCBOR UnMarshals supplied CBOR bytes to TeeInstanceID
func (o *TeeInstanceID) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &o.val)
}
