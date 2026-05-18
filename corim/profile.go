// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
)

// IProfileValue is the interface implemented by all types whose instances can
// be used as Profile values.
type IProfileValue interface {
	extensions.ITypeChoiceValue
}

// Profile is an identification of a CoRIM profile that defines which of the
// optional parts of CoRIM are required, which are prohibited, and which
// extension points are exercised and how.
type Profile struct {
	Value IProfileValue
}

// NewProfile instantiates a new Profile using the provided value and type
// name.
func NewProfile(val any, typ string) (*Profile, error) {
	factory, ok := profileValueRegister[typ]
	if !ok {
		return nil, fmt.Errorf("unknown profile type: %s", typ)
	}

	return factory(val)
}

// NewProfileFromString creates a new Profile based on the provided string. The
// string must be either a valid absolute OID or a valid absolute URI.
func NewProfileFromString(uriOrOID string) (*Profile, error) {
	oid, err := comid.NewTaggedOID(uriOrOID)
	if err == nil {
		return &Profile{oid}, nil
	}

	return NewURIProfile(uriOrOID)
}

// Type returns the type of the underlying value.
func (o Profile) Type() string {
	return o.Value.Type()
}

// String returns a string representation of the Profile's value.
func (o Profile) String() string {
	return o.Value.String()
}

// Valid returns an error if the underlying profile value is invalid
func (o *Profile) Valid() error {
	if o.IsNil() {
		return errors.New("cannot be nil")
	}

	return o.Value.Valid()
}

// IsNil returns true if the *Profile is nil or if its innver value is nil.
func (o *Profile) IsNil() bool {
	return o == nil || o.Value == nil
}

func (o Profile) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.Value)
}

func (o *Profile) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &o.Value)
}

func (o Profile) MarshalJSON() ([]byte, error) {
	return extensions.TypeChoiceValueMarshalJSON(o.Value)
}

func (o *Profile) UnmarshalJSON(data []byte) error {
	var tnv encoding.TypeAndValue
	if err := json.Unmarshal(data, &tnv); err != nil {
		return err
	}

	decoded, err := NewProfile(nil, tnv.Type)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(tnv.Value, &decoded.Value); err != nil {
		return fmt.Errorf("unmarshalling %s: %w", tnv.Type, err)
	}

	if err := decoded.Value.Valid(); err != nil {
		return fmt.Errorf("invalid %s: %w", tnv.Type, err)
	}

	o.Value = decoded.Value

	return nil
}

// NewURIProfile instantiate a new Profile based on the provided value, which
// must be convertible into a valid URI.
func NewURIProfile(val any) (*Profile, error) {
	uri, err := comid.NewTaggedURI(val)
	if err != nil {
		return nil, err
	}

	return &Profile{uri}, nil
}

// MustNewURIProfile is like MustNewURIProfile except it panics on error.
func MustNewURIProfile(val any) *Profile {
	ret, err := NewURIProfile(val)
	if err != nil {
		panic(err)
	}

	return ret
}

// NewOIDProfile instantiates a new Profile based on the provided value, which
// must be convertible into a valid OID.
func NewOIDProfile(val any) (*Profile, error) {
	oid, err := comid.NewTaggedOID(val)
	if err != nil {
		return nil, err
	}

	return &Profile{oid}, nil
}

// MustNewOIDProfile is like MustNewOIDProfile except it panics on error.
func MustNewOIDProfile(val any) *Profile {
	ret, err := NewOIDProfile(val)
	if err != nil {
		panic(err)
	}

	return ret
}

// IProfileFactory defines the signature for the factory functions that may be
// registered using RegisterProfileType to provide a new implementation of the
// corresponding type choice. The factory function should create a new *Profile
// with the underlying value created based on the provided input. The range of
// valid inputs is up to the specific type choice implementation, however it
// _must_ accept nil as one of the inputs, and return the Zero value for
// implemented type.
// See also https://go.dev/ref/spec#The_zero_value
type IProfileFactory func(any) (*Profile, error)

var profileValueRegister = map[string]IProfileFactory{
	comid.OIDType: NewOIDProfile,
	comid.URIType: NewURIProfile,
}

// RegisterProfileType registers a new IProfileValue implementation (created
// by the provided IProfileFactory) under the specified CBOR tag.
func RegisterProfileType(tag uint64, factory IProfileFactory) error {
	nilVal, err := factory(nil)
	if err != nil {
		return err
	}

	typ := nilVal.Type()
	if _, exists := profileValueRegister[typ]; exists {
		return fmt.Errorf("profile type with name %q already exists", typ)
	}

	if err := registerCORIMTag(tag, nilVal.Value); err != nil {
		return err
	}

	profileValueRegister[typ] = factory

	return nil
}
