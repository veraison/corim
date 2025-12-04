// Copyright 2021-2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
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
		o.val = TaggedBytes(val)
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

// Equal confirms if the RawValue instances are equal
func (o RawValue) Equal(r RawValue) bool {
	return reflect.DeepEqual(o, r)
}

// CompareAgainstReference checks if a RawValue object matches with a reference
//
// See section-8.9.6.1.4 in the IETF CoRIM spec for the rules to compare a
// RawValue object against a reference.
func (o RawValue) CompareAgainstReference(ref []byte, mask *[]byte) bool {
	claim, err := o.GetBytes()
	if err != nil {
		log.Printf("RawValue:CompareAgainstReference: Error: %v", err)
		return false
	}

	if mask != nil {
		if len(claim) != len(ref) {
			return false
		}

		if len(*mask) != len(claim) {
			log.Printf("RawValue:CompareAgainstReference: Error: mask length")
			return false
		}

		for i := range *mask {
			claim[i] = (*mask)[i] & claim[i]
			ref[i] = (*mask)[i] & ref[i]
		}
	}

	return bytes.Equal(claim, ref)
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
