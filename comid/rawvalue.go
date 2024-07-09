// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
)

// RawValue models a $raw-value-type-choice.  For now, the only available type is bytes.
type RawValue struct {
	val interface{}
}

func NewRawValue() *RawValue {
	return &RawValue{}
}

func (o *RawValue) SetBytes(val []byte) *RawValue {
	if o != nil {
		v, err := NewBytes(val)
		if err != nil {
			return nil
		}
		o.val = *v
	}
	return o
}

func (o RawValue) GetBytes() ([]byte, error) {
	if o.val == nil {
		return nil, fmt.Errorf("raw value is not set")
	}

	switch t := o.val.(type) {
	case TaggedBytes:
		return []byte(t), nil
	default:
		return nil, fmt.Errorf("unknown type %T for $raw-value-type-choice", t)
	}
}

func (o RawValue) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

func (o *RawValue) UnmarshalCBOR(data []byte) error {
	var rawValue TaggedBytes

	if dm.Unmarshal(data, &rawValue) == nil {
		o.val = rawValue
		return nil
	}

	return fmt.Errorf("unknown raw-value (CBOR: %x)", data)
}

// UnmarshalJSON deserializes the type'n'value JSON object into the target RawValue.
// The only supported type is BytesType with value
func (o *RawValue) UnmarshalJSON(data []byte) error {
	var v tnv

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Type {
	case BytesType:
		var x []byte
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal $raw-value-type-choice of type bytes: %w",
				err,
			)
		}
		o.val = TaggedBytes(x)
	default:
		return fmt.Errorf("unknown type %s for $raw-value-type-choice", v.Type)
	}

	return nil
}

func (o RawValue) MarshalJSON() ([]byte, error) {
	var (
		v   tnv
		b   []byte
		err error
	)

	switch t := o.val.(type) {
	case TaggedBytes:
		b, err = json.Marshal(o.val)
		if err != nil {
			return nil, err
		}
		v = tnv{Type: BytesType, Value: b}
	default:
		return nil, fmt.Errorf("unknown type %T for raw-value-type-choice", t)
	}

	return json.Marshal(v)
}
