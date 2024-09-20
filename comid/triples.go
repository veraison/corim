// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"

	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
)

type Triples struct {
	ReferenceValues *ValueTriples `cbor:"0,keyasint,omitempty" json:"reference-values,omitempty"`
	EndorsedValues  *ValueTriples `cbor:"1,keyasint,omitempty" json:"endorsed-values,omitempty"`
	DevIdentityKeys *KeyTriples   `cbor:"2,keyasint,omitempty" json:"dev-identity-keys,omitempty"`
	AttestVerifKeys *KeyTriples   `cbor:"3,keyasint,omitempty" json:"attester-verification-keys,omitempty"`

	Extensions
}

// RegisterExtensions registers a struct as a collections of extensions
func (o *Triples) RegisterExtensions(exts extensions.Map) error {
	refValExts := extensions.NewMap()
	endValExts := extensions.NewMap()

	for p, v := range exts {
		switch p {
		case ExtTriples:
			o.Extensions.Register(v)
		case ExtReferenceValue:
			refValExts[ExtMval] = v
		case ExtReferenceValueFlags:
			refValExts[ExtFlags] = v
		case ExtEndorsedValue:
			endValExts[ExtMval] = v
		case ExtEndorsedValueFlags:
			endValExts[ExtFlags] = v
		default:
			return fmt.Errorf("%w: %q", extensions.ErrUnexpectedPoint, p)
		}
	}

	if len(refValExts) != 0 {
		if o.ReferenceValues == nil {
			o.ReferenceValues = NewValueTriples()
		}

		if err := o.ReferenceValues.RegisterExtensions(refValExts); err != nil {
			return err
		}
	}

	if len(endValExts) != 0 {
		if o.EndorsedValues == nil {
			o.EndorsedValues = NewValueTriples()
		}

		if err := o.EndorsedValues.RegisterExtensions(refValExts); err != nil {
			return err
		}
	}

	return nil
}

// GetExtensions returns previously registered extension
func (o *Triples) GetExtensions() extensions.IMapValue {
	return o.Extensions.IMapValue
}

// UnmarshalCBOR deserializes from CBOR
func (o *Triples) UnmarshalCBOR(data []byte) error {
	return encoding.PopulateStructFromCBOR(dm, data, o)
}

// MarshalCBOR serializes to CBOR
func (o Triples) MarshalCBOR() ([]byte, error) {
	// If extensions have been registered, the collection will exist, but
	// might be empty. If that is the case, set it to nil to avoid
	// marshaling an empty list (and let the marshaller omit the claim
	// instead). Note that since the receiver was passed by value, we do not
	// need to worry about saving the field's value before setting it to
	// nil.
	if o.ReferenceValues != nil && o.ReferenceValues.IsEmpty() {
		o.ReferenceValues = nil
	}

	if o.EndorsedValues != nil && o.EndorsedValues.IsEmpty() {
		o.EndorsedValues = nil
	}

	return encoding.SerializeStructToCBOR(em, o)
}

// UnmarshalJSON deserializes from JSON
func (o *Triples) UnmarshalJSON(data []byte) error {
	return encoding.PopulateStructFromJSON(data, o)
}

// MarshalJSON serializes to JSON
func (o Triples) MarshalJSON() ([]byte, error) {
	// If extensions have been registered, the collection will exist, but
	// might be empty. If that is the case, set it to nil to avoid
	// marshaling an empty list (and let the marshaller omit the claim
	// instead). Note that since the receiver was passed by value, we do not
	// need to worry about saving the field's value before setting it to
	// nil.
	if o.ReferenceValues != nil && o.ReferenceValues.IsEmpty() {
		o.ReferenceValues = nil
	}

	if o.EndorsedValues != nil && o.EndorsedValues.IsEmpty() {
		o.EndorsedValues = nil
	}

	return encoding.SerializeStructToJSON(o)
}

// Valid checks that the Triples is valid as per the specification
func (o Triples) Valid() error {
	// non-empty<>
	if (o.ReferenceValues == nil || o.ReferenceValues.IsEmpty()) &&
		(o.EndorsedValues == nil || o.EndorsedValues.IsEmpty()) &&
		(o.AttestVerifKeys == nil || len(*o.AttestVerifKeys) == 0) &&
		(o.DevIdentityKeys == nil || len(*o.DevIdentityKeys) == 0) {
		return fmt.Errorf("triples struct must not be empty")
	}

	if o.ReferenceValues != nil {
		if err := o.ReferenceValues.Valid(); err != nil {
			return fmt.Errorf("reference values: %w", err)
		}
	}

	if o.EndorsedValues != nil {
		if err := o.EndorsedValues.Valid(); err != nil {
			return fmt.Errorf("endorsed values: %w", err)
		}
	}

	if o.AttestVerifKeys != nil {
		for i, ak := range *o.AttestVerifKeys {
			if err := ak.Valid(); err != nil {
				return fmt.Errorf("attestation verification key at index %d: %w", i, err)
			}
		}
	}

	if o.DevIdentityKeys != nil {
		for i, dk := range *o.DevIdentityKeys {
			if err := dk.Valid(); err != nil {
				return fmt.Errorf("device identity key at index %d: %w", i, err)
			}
		}
	}

	return o.Extensions.validTriples(&o)
}

func (o *Triples) AddReferenceValue(val ValueTriple) *Triples {
	if o != nil {
		if o.ReferenceValues == nil {
			o.ReferenceValues = new(ValueTriples)
		}

		o.ReferenceValues.Add(&val)
	}

	return o
}

func (o *Triples) AddEndorsedValue(val ValueTriple) *Triples {
	if o != nil {
		if o.EndorsedValues == nil {
			o.EndorsedValues = new(ValueTriples)
		}

		o.EndorsedValues.Add(&val)
	}

	return o
}

func (o *Triples) AddAttestVerifKey(val KeyTriple) *Triples {
	if o != nil {
		*o.AttestVerifKeys = append(*o.AttestVerifKeys, val)
	}

	return o
}

func (o *Triples) AddDevIdentityKey(val KeyTriple) *Triples {
	if o != nil {
		*o.DevIdentityKeys = append(*o.DevIdentityKeys, val)
	}

	return o
}
