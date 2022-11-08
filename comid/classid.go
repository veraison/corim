// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// ClassID represents a $class-id-type-choice, which can be one of TaggedUUID,
// TaggedOID, or TaggedImplID (PSA-specific extension)
type ClassID struct {
	val interface{}
}

type ClassIDType uint16

const (
	ClassIDTypeUUID = ClassIDType(iota)
	ClassIDTypeImplID
	ClassIDTypeOID

	ClassIDTypeUnknown = ^ClassIDType(0)
)

// SetUUID sets the value of the targed ClassID to the supplied UUID
func (o *ClassID) SetUUID(uuid UUID) *ClassID {
	if o != nil {
		o.val = TaggedUUID(uuid)
	}
	return o
}

type ImplID [32]byte
type TaggedImplID ImplID

func (o ImplID) MarshalJSON() ([]byte, error) {
	return json.Marshal(o[:])
}

func (o *ImplID) UnmarshalJSON(data []byte) error {
	var b []byte

	if err := json.Unmarshal(data, &b); err != nil {
		return fmt.Errorf("bad ImplID: %w", err)
	}

	if nb := len(b); nb != 32 {
		return fmt.Errorf("bad ImplID format: got %d bytes, want 32", nb)
	}

	copy(o[:], b)

	return nil
}

type TaggedOID OID

// SetImplID sets the value of the targed ClassID to the supplied PSA
// Implementation ID (see Section 3.2.2 of draft-tschofenig-rats-psa-token)
func (o *ClassID) SetImplID(implID ImplID) *ClassID {
	if o != nil {
		o.val = TaggedImplID(implID)
	}
	return o
}

func (o ClassID) GetImplID() (ImplID, error) {
	switch t := o.val.(type) {
	case TaggedImplID:
		return ImplID(t), nil
	default:
		return ImplID{}, fmt.Errorf("class-id type is: %T", t)
	}
}

// SetOID sets the value of the targed ClassID to the supplied OID.
// The OID is a string in dotted-decimal notation
func (o *ClassID) SetOID(s string) *ClassID {
	if o != nil {
		var berOID OID
		if berOID.FromString(s) != nil {
			return nil
		}
		o.val = TaggedOID(berOID)
	}
	return o
}

// MarshalCBOR serializes the target ClassID to CBOR
func (o ClassID) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

// UnmarshalCBOR deserializes the supplied CBOR buffer into the target ClassID.
// It is undefined behavior to try and inspect the target ClassID in case this
// method returns an error.
func (o *ClassID) UnmarshalCBOR(data []byte) error {
	var implID TaggedImplID

	if dm.Unmarshal(data, &implID) == nil {
		o.val = implID
		return nil
	}

	var uuid TaggedUUID

	if dm.Unmarshal(data, &uuid) == nil {
		o.val = uuid
		return nil
	}

	var oid TaggedOID

	if dm.Unmarshal(data, &oid) == nil {
		o.val = oid
		return nil
	}

	return fmt.Errorf("unknown class id (CBOR: %x)", data)
}

// UnmarshalJSON deserializes the supplied JSON object into the target ClassID
// The class id object must have one of the following shapes:
//
// UUID:
//
//	{
//	  "type": "uuid",
//	  "value": "69E027B2-7157-4758-BCB4-D9F167FE49EA"
//	}
//
// OID:
//
//	{
//	  "type": "oid",
//	  "value": "2.16.840.1.113741.1.15.4.2"
//	}
//
// PSA Implementation ID:
//
//	{
//	  "type": "psa.impl-id",
//	  "value": "YWNtZS1pbXBsZW1lbnRhdGlvbi1pZC0wMDAwMDAwMDE="
//	}
func (o *ClassID) UnmarshalJSON(data []byte) error {
	var v tnv

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Type {
	case "uuid": // nolint: goconst
		var x UUID
		if err := x.UnmarshalJSON(v.Value); err != nil {
			return err
		}
		o.val = TaggedUUID(x)
	case "oid":
		var x OID
		if err := x.UnmarshalJSON(v.Value); err != nil {
			return err
		}
		o.val = TaggedOID(x)
	case "psa.impl-id":
		var x ImplID
		if err := x.UnmarshalJSON(v.Value); err != nil {
			return err
		}
		o.val = TaggedImplID(x)
	default:
		return fmt.Errorf("unknown type '%s' for class id", v.Type)
	}

	return nil
}

// MarshalJSON serializes the target ClassID to JSON
func (o ClassID) MarshalJSON() ([]byte, error) {
	var (
		v   tnv
		b   []byte
		err error
	)

	switch t := o.val.(type) {
	case TaggedUUID:
		b, err = UUID(t).MarshalJSON()
		if err != nil {
			return nil, err
		}
		v = tnv{Type: "uuid", Value: b}
	case TaggedOID:
		b, err = OID(t).MarshalJSON()
		if err != nil {
			return nil, err
		}
		v = tnv{Type: "oid", Value: b}
	case TaggedImplID:
		b, err = ImplID(t).MarshalJSON()
		if err != nil {
			return nil, err
		}
		v = tnv{Type: "psa.impl-id", Value: b}
	default:
		return nil, fmt.Errorf("unknown type %T for class-id", t)
	}

	return json.Marshal(v)
}

// Type returns the type of the target ClassID, i.e., one of UUID, OID or PSA
// Implementation ID
func (o ClassID) Type() ClassIDType {
	switch o.val.(type) {
	case TaggedUUID:
		return ClassIDTypeUUID
	case TaggedImplID:
		return ClassIDTypeImplID
	case TaggedOID:
		return ClassIDTypeOID
	}
	return ClassIDTypeUnknown
}

// String returns a printable string of the ClassID value. UUIDs use the
// canonical 8-4-4-4-12 format, PSA Implementation IDs are base64 encoded.
// OIDs are output in dotted-decimal notation.
func (o ClassID) String() string {
	switch t := o.val.(type) {
	case TaggedUUID:
		return UUID(t).String()
	case TaggedImplID:
		b := [32]byte(t)
		return base64.StdEncoding.EncodeToString(b[:])
	case TaggedOID:
		return OID(t).String()
	default:
		return ""
	}
}

// Unset tests whether the target ClassID has been initialized
func (o ClassID) Unset() bool {
	return o.val == nil || o.Type() == ClassIDTypeUnknown
}
