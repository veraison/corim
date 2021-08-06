// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "fmt"

type Triples struct {
	ReferenceValues *[]ReferenceValue `cbor:"0,keyasint,omitempty" json:"reference-values,omitempty"`
	EndorsedValues  *[]EndorsedValue  `cbor:"1,keyasint,omitempty" json:"endorsed-values,omitempty"`
	AttestVerifKeys *[]AttestVerifKey `cbor:"2,keyasint,omitempty" json:"attester-verification-keys,omitempty"`
	DevIdentityKeys *[]DevIdentityKey `cbor:"3,keyasint,omitempty" json:"dev-identity-keys,omitempty"`
}

// Valid checks that the Triples is valid as per the specification
func (o Triples) Valid() error {
	// non-empty<>
	if o.ReferenceValues == nil && o.EndorsedValues == nil &&
		o.AttestVerifKeys == nil && o.DevIdentityKeys == nil {
		return fmt.Errorf("triples struct must not be empty")
	}

	if o.ReferenceValues != nil {
		for i, rv := range *o.ReferenceValues {
			if err := rv.Valid(); err != nil {
				return fmt.Errorf("reference value at index %d: %w", i, err)
			}
		}
	}

	if o.EndorsedValues != nil {
		for i, ev := range *o.EndorsedValues {
			if err := ev.Valid(); err != nil {
				return fmt.Errorf("endorsed value at index %d: %w", i, err)
			}
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
	return nil
}

func (o *Triples) AddReferenceValue(val ReferenceValue) *Triples {
	if o != nil {
		*o.ReferenceValues = append(*o.ReferenceValues, val)
	}
	return o
}

func (o *Triples) AddEndorsedValue(val EndorsedValue) *Triples {
	if o != nil {
		*o.EndorsedValues = append(*o.EndorsedValues, val)
	}
	return o
}

func (o *Triples) AddAttestVerifKey(val AttestVerifKey) *Triples {
	if o != nil {
		*o.AttestVerifKeys = append(*o.AttestVerifKeys, val)
	}
	return o
}

func (o *Triples) AddDevIdentityKey(val DevIdentityKey) *Triples {
	if o != nil {
		*o.DevIdentityKeys = append(*o.DevIdentityKeys, val)
	}
	return o
}
