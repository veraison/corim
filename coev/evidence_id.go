// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
)

// EvidenceID stores an evidence identity. The supported formats as of now are UUID.
type EvidenceID struct {
	Value IEvidenceValue
}

// NewEvidenceID creates a new evidence with the value of the specified type
// populated using the provided value.
func NewEvidenceID(val any, typ string) (*EvidenceID, error) {
	factory, ok := evidenceIDValueRegister[typ]
	if !ok {
		return nil, fmt.Errorf("unknown EvidenceID type: %s", typ)
	}

	return factory(val)
}

// Valid checks for the validity of given EvidenceID
func (o EvidenceID) Valid() error {
	if o.String() == "" {
		return errors.New("no EvidenceID")
	}
	_, err := o.GetUUID()
	if err != nil {
		return fmt.Errorf("unable to fetch valid UUID: %w", err)
	}
	return nil
}

// String returns a printable string of the EvidenceID value.  UUIDs use the
// canonical 8-4-4-4-12 format.
func (o EvidenceID) String() string {
	if o.Value == nil {
		return ""
	}

	return o.Value.String()
}

// Type returns a string naming the type of the underlying EvidenceID value.
func (o EvidenceID) Type() string {
	return o.Value.Type()
}

// Bytes returns a []byte containing the bytes of the underlying EvidenceID
// value.
func (o EvidenceID) Bytes() []byte {
	return o.Value.Bytes()
}

// MarshalCBOR serializes the target EvidenceID to CBOR
func (o EvidenceID) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.Value)
}

func (o *EvidenceID) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &o.Value)
}

// UnmarshalJSON deserializes the supplied JSON object into the target EvidenceID
// The evidence object must have the following shape:
//
//	{
//	  "type": "<EVIDENCE_TYPE>",
//	  "value": <EVIDENCE_VALUE>
//	}
//
// where <EVIDENCE_TYPE> must be one of the known IEvidenceValue implementation
// type names (available in the base implementation: "uuid"), and
// <EVIDENCE_VALUE> is the JSON encoding of the evidence value. The exact
// encoding is <EVIDENCE_TYPE> dependent. For the base implmentation types it is
//
//	uuid: standard UUID string representation, e.g. "550e8400-e29b-41d4-a716-446655440000"

//nolint:dupl
func (o *EvidenceID) UnmarshalJSON(data []byte) error {
	var tnv encoding.TypeAndValue

	if err := json.Unmarshal(data, &tnv); err != nil {
		return fmt.Errorf("EvidenceID decoding failure: %w", err)
	}

	decoded, err := NewEvidenceID(nil, tnv.Type)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(tnv.Value, &decoded.Value); err != nil {
		return fmt.Errorf(
			"cannot unmarshal EvidenceID: %w",
			err,
		)
	}

	if err := decoded.Value.Valid(); err != nil {
		return fmt.Errorf("invalid %s: %w", tnv.Type, err)
	}

	o.Value = decoded.Value

	return nil
}

// MarshalJSON serializes the EvidenceID into a JSON object.
func (o EvidenceID) MarshalJSON() ([]byte, error) {
	valueBytes, err := json.Marshal(o.Value)
	if err != nil {
		return nil, err
	}

	value := encoding.TypeAndValue{
		Type:  o.Value.Type(),
		Value: valueBytes,
	}

	return json.Marshal(value)
}

// SetUUID sets the identity of the target Evidence to the supplied UUID
func (o *EvidenceID) SetUUID(val uuid.UUID) *EvidenceID {
	if o != nil {
		o.Value = comid.TaggedUUID(val)
	}
	return o
}

func (o EvidenceID) GetUUID() (comid.UUID, error) {
	switch t := o.Value.(type) {
	case *comid.TaggedUUID:
		return comid.UUID(*t), nil
	case comid.TaggedUUID:
		return comid.UUID(t), nil
	default:
		return comid.UUID{}, fmt.Errorf("evidence-id type is: %T", t)
	}
}

// IEvidenceValue is the interface implemented by all EvidenceID value
// implementations.
type IEvidenceValue interface {
	extensions.ITypeChoiceValue

	Bytes() []byte
}

// NewUUIDEvidenceID instantiates a new EvidenceID with the supplied UUID identity
func NewUUIDEvidenceID(val any) (*EvidenceID, error) {
	if val == nil {
		return &EvidenceID{&comid.TaggedUUID{}}, nil
	}

	ret, err := comid.NewTaggedUUID(val)
	if err != nil {
		return nil, err
	}

	return &EvidenceID{ret}, nil
}

// MustNewUUIDEvidenceID is like NewUUIDEvidenceID execept it does not return an
// error, assuming that the provided value is valid. It panics if that isn't
// the case.
func MustNewUUIDEvidenceID(val any) *EvidenceID {
	ret, err := NewUUIDEvidenceID(val)
	if err != nil {
		panic(err)
	}

	return ret
}

// IEvidenceValueFactory defines the signature for the factory functions that may be
// registered using RegisterEvidenceType to provide a new implementation of the
// corresponding type choice. The factory function should create a new *EvidenceID
// with the underlying value created based on the provided input. The range of
// valid inputs is up to the specific type choice implementation, however it
// _must_ accept nil as one of the inputs, and return the Zero value for
// implemented type.
// See also https://go.dev/ref/spec#The_zero_value
type IEvidenceFactory func(any) (*EvidenceID, error)

var evidenceIDValueRegister = map[string]IEvidenceFactory{
	comid.UUIDType: NewUUIDEvidenceID,
}
