// Copyright 2021-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"iter"

	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
)

type Triples struct {
	ReferenceValues    *ValueTriples             `cbor:"0,keyasint,omitempty" json:"reference-values,omitempty"`
	EndorsedValues     *ValueTriples             `cbor:"1,keyasint,omitempty" json:"endorsed-values,omitempty"`
	DevIdentityKeys    *KeyTriples               `cbor:"2,keyasint,omitempty" json:"dev-identity-keys,omitempty"`
	AttestVerifKeys    *KeyTriples               `cbor:"3,keyasint,omitempty" json:"attester-verification-keys,omitempty"`
	DomainDependencies *DomainDependencyTriples  `cbor:"4,keyasint,omitempty" json:"dependency-triples,omitempty"`
	DomainMemberships  *DomainMembershipTriples  `cbor:"5,keyasint,omitempty" json:"membership-triples,omitempty"`
	CondEndorseSeries  *CondEndorseSeriesTriples `cbor:"8,keyasint,omitempty" json:"conditional-endorsement-series,omitempty"`
	CondEndorsements   *CondEndorseTriples       `cbor:"10,keyasint,omitempty" json:"conditional-endorsements,omitempty"`
	Extensions
}

// RegisterExtensions registers a struct as a collections of extensions
func (o *Triples) RegisterExtensions(exts extensions.Map) error {
	refValExts := extensions.NewMap()
	endValExts := extensions.NewMap()
	conSeriesExts := extensions.NewMap()
	conExts := extensions.NewMap()

	for p, v := range exts {
		switch p {
		case ExtTriples:
			o.Register(v)
		case ExtReferenceValue:
			refValExts[ExtMval] = v
		case ExtReferenceValueFlags:
			refValExts[ExtFlags] = v
		case ExtEndorsedValue:
			endValExts[ExtMval] = v
		case ExtEndorsedValueFlags:
			endValExts[ExtFlags] = v
		case ExtCondEndorseValue:
			conExts[ExtMval] = v
		case ExtCondEndorseValueFlags:
			conExts[ExtFlags] = v
		case ExtCondEndorseSeriesValue:
			conSeriesExts[ExtMval] = v
		case ExtCondEndorseSeriesValueFlags:
			conSeriesExts[ExtFlags] = v
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

		if err := o.EndorsedValues.RegisterExtensions(endValExts); err != nil {
			return err
		}
	}

	if len(conSeriesExts) != 0 {
		if o.CondEndorseSeries == nil {
			o.CondEndorseSeries = NewCondEndorseSeriesTriples()
		}

		if err := o.CondEndorseSeries.RegisterExtensions(conSeriesExts); err != nil {
			return err
		}
	}

	if len(conExts) != 0 {
		if o.CondEndorsements == nil {
			o.CondEndorsements = NewCondEndorseTriples()
		}

		if err := o.CondEndorsements.RegisterExtensions(conExts); err != nil {
			return err
		}
	}

	return nil
}

// GetExtensions returns previously registered extension
func (o *Triples) GetExtensions() extensions.IMapValue {
	return o.IMapValue
}

// UnmarshalCBOR deserializes from CBOR
func (o *Triples) UnmarshalCBOR(data []byte) error {
	return encoding.PopulateStructFromCBOR(dm, data, o)
}

// nolint:gocritic
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

	if o.CondEndorseSeries != nil && o.CondEndorseSeries.IsEmpty() {
		o.CondEndorseSeries = nil
	}

	if o.CondEndorsements != nil && o.CondEndorsements.IsEmpty() {
		o.CondEndorsements = nil
	}

	if o.DomainDependencies != nil && o.DomainDependencies.IsEmpty() {
		o.DomainDependencies = nil
	}

	if o.DomainMemberships != nil && o.DomainMemberships.IsEmpty() {
		o.DomainMemberships = nil
	}

	return encoding.SerializeStructToCBOR(em, o)
}

// UnmarshalJSON deserializes from JSON
func (o *Triples) UnmarshalJSON(data []byte) error {
	return encoding.PopulateStructFromJSON(data, o)
}

// nolint:gocritic
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

	if o.CondEndorseSeries != nil && o.CondEndorseSeries.IsEmpty() {
		o.CondEndorseSeries = nil
	}

	if o.CondEndorsements != nil && o.CondEndorsements.IsEmpty() {
		o.CondEndorsements = nil
	}

	if o.DomainDependencies != nil && o.DomainDependencies.IsEmpty() {
		o.DomainDependencies = nil
	}

	if o.DomainMemberships != nil && o.DomainMemberships.IsEmpty() {
		o.DomainMemberships = nil
	}

	return encoding.SerializeStructToJSON(o)
}

// IterRefVals provides an iterator over reference value ValueTriple's inside
// the Comid.
func (o *Triples) IterRefVals() iter.Seq[*ValueTriple] {
	seq := func(yield func(*ValueTriple) bool) {
		if o.ReferenceValues != nil {
			for _, vt := range o.ReferenceValues.Values {
				if !yield(&vt) {
					return
				}
			}
		}
	}

	return seq
}

// IterRefVals provides an iterator over endorsed value ValueTriple's inside
// the Triples.
func (o *Triples) IterEndVals() iter.Seq[*ValueTriple] {
	seq := func(yield func(*ValueTriple) bool) {
		if o.EndorsedValues != nil {
			for _, vt := range o.EndorsedValues.Values {
				if !yield(&vt) {
					return
				}
			}
		}
	}

	return seq
}

// IterAttestVerifKeys provides an iterator over attest. verif. key KeyTriple's
// inside the Triples.
func (o *Triples) IterAttestVerifKeys() iter.Seq[*KeyTriple] {
	seq := func(yield func(*KeyTriple) bool) {
		if o.AttestVerifKeys != nil {
			for _, kt := range *o.AttestVerifKeys {
				if !yield(&kt) {
					return
				}
			}
		}
	}

	return seq
}

// IterDevIdentityKeys provides an iterator over device identity key
// KeyTriple's inside the Triples.
func (o *Triples) IterDevIdentityKeys() iter.Seq[*KeyTriple] {
	seq := func(yield func(*KeyTriple) bool) {
		if o.DevIdentityKeys != nil {
			for _, kt := range *o.DevIdentityKeys {
				if !yield(&kt) {
					return
				}
			}
		}
	}

	return seq
}

// Valid checks that the Triples is valid as per the specification
// nolint:gocyclo
func (o *Triples) Valid() error {
	// non-empty<>
	if (o.ReferenceValues == nil || o.ReferenceValues.IsEmpty()) &&
		(o.EndorsedValues == nil || o.EndorsedValues.IsEmpty()) &&
		(o.AttestVerifKeys == nil || len(*o.AttestVerifKeys) == 0) &&
		(o.DevIdentityKeys == nil || len(*o.DevIdentityKeys) == 0) &&
		(o.DomainDependencies == nil || o.DomainDependencies.IsEmpty()) &&
		(o.DomainMemberships == nil || o.DomainMemberships.IsEmpty()) &&
		(o.CondEndorseSeries == nil || o.CondEndorseSeries.IsEmpty()) &&
		(o.CondEndorsements == nil || o.CondEndorsements.IsEmpty()) {
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

	if o.DomainDependencies != nil {
		if err := o.DomainDependencies.Valid(); err != nil {
			return fmt.Errorf("dependency triples: %w", err)
		}
	}

	if o.DomainMemberships != nil {
		if err := o.DomainMemberships.Valid(); err != nil {
			return fmt.Errorf("membership triples: %w", err)
		}
	}

	if o.CondEndorseSeries != nil {
		if err := o.CondEndorseSeries.Valid(); err != nil {
			return fmt.Errorf("conditional series: %w", err)
		}
	}

	if o.CondEndorsements != nil {
		if err := o.CondEndorsements.Valid(); err != nil {
			return fmt.Errorf("conditional endorsements: %w", err)
		}
	}

	return o.validTriples(o)
}

func (o *Triples) AddReferenceValue(val *ValueTriple) *Triples {
	if o != nil {
		if o.ReferenceValues == nil {
			o.ReferenceValues = new(ValueTriples)
		}

		o.ReferenceValues.Add(val)
	}

	return o
}

func (o *Triples) AddEndorsedValue(val *ValueTriple) *Triples {
	if o != nil {
		if o.EndorsedValues == nil {
			o.EndorsedValues = new(ValueTriples)
		}

		o.EndorsedValues.Add(val)
	}

	return o
}

func (o *Triples) AddAttestVerifKey(val *KeyTriple) *Triples {
	if o != nil {
		*o.AttestVerifKeys = append(*o.AttestVerifKeys, *val)
	}

	return o
}

func (o *Triples) AddDevIdentityKey(val *KeyTriple) *Triples {
	if o != nil {
		*o.DevIdentityKeys = append(*o.DevIdentityKeys, *val)
	}

	return o
}

// AddDomainDependency appends a domain-dependency-triple to DomainDependencies.
func (o *Triples) AddDomainDependency(val *DomainDependencyTriple) *Triples {
	if o != nil {
		if o.DomainDependencies == nil {
			o.DomainDependencies = new(DomainDependencyTriples)
		}
		*o.DomainDependencies = append(*o.DomainDependencies, *val)
	}

	return o
}

// AddDomainMembership appends a domain-memberbship-triple to DomainMemberships.
func (o *Triples) AddDomainMembership(val *DomainMembershipTriple) *Triples {
	if o != nil {
		if o.DomainMemberships == nil {
			o.DomainMemberships = new(DomainMembershipTriples)
		}
		*o.DomainMemberships = append(*o.DomainMemberships, *val)
	}

	return o
}

// nolint:gocritic
func (o *Triples) AddCondEndorseSeries(val *CondEndorseSeriesTriple) *Triples {
	if o != nil {
		if o.CondEndorseSeries == nil {
			o.CondEndorseSeries = new(CondEndorseSeriesTriples)
		}

		o.CondEndorseSeries.Add(val)
	}

	return o
}

func (o *Triples) AddCondEndorsement(val *CondEndorseTriple) *Triples {
	if o != nil {
		if o.CondEndorsements == nil {
			o.CondEndorsements = NewCondEndorseTriples()
		}

		o.CondEndorsements.Add(val)
	}

	return o
}
