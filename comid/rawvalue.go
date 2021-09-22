// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
)

type RawValue struct {
	val interface{}
}

type TaggedRawValueBytes []byte

func NewRawValue() *RawValue {
	return &RawValue{}
}

func (o *RawValue) SetBytes(val []byte) *RawValue {
	if o != nil {
		o.val = TaggedRawValueBytes(val)
	}
	return o
}

func (o RawValue) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

func (o *RawValue) UnmarshalCBOR(data []byte) error {
	var rawValue TaggedRawValueBytes

	if dm.Unmarshal(data, &rawValue) == nil {
		o.val = rawValue
		return nil
	}

	return fmt.Errorf("unknown raw-value (CBOR: %x)", data)
}

// UnmarshalJSON deserializes the type'n'value JSON object into the target RawValue.
// The only supported type is "bytes" with value
func (o *RawValue) UnmarshalJSON(data []byte) error {
	var v tnv

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Type {
	case "bytes":
		var x []byte
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal $raw-value-type-choice of type bytes: %w",
				err,
			)
		}
		o.val = TaggedRawValueBytes(x)
	default:
		return fmt.Errorf("unknown type %s for $raw-value-type-choice", v.Type)
	}

	return nil
}
