// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"errors"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// KeyTripleCondition is the Conditions part of IdentityTriple, together
// with the Environment field, it is used to match an IdentityTriple to a
// Target Environment.
type KeyTripleCondition struct {
	Mkey         *Mkey       `cbor:"0,keyasint,omitempty" json:"mkey,omitempty"`
	AuthorizedBy *CryptoKeys `cbor:"1,keyasint,omitempty" json:"authorized-by,omitempty"`
}

func (o *KeyTripleCondition) Valid() error {
	if o.Mkey == nil && o.AuthorizedBy == nil {
		return errors.New("condition must not be empty")
	}

	if o.Mkey != nil {
		if err := o.Mkey.Valid(); err != nil {
			return fmt.Errorf("mkey: %w", err)
		}
	}

	if o.AuthorizedBy != nil {
		if err := o.AuthorizedBy.Valid(); err != nil {
			return fmt.Errorf("authorized-by: %w", err)
		}
	}

	return nil
}

// KeyTriple endorses that contained keys were securely provisioned to the
// named Target Environment.
//
//	attest-key-triple-record = [
//	  environment: environment-map
//	  verification-keys: [ + $crypto-key-type-choice ]
//	  ? conditions: non-empty<{
//	    ? &(mkey: 0) => $measured-element-type-choice,
//	    ? &(authorized-by: 1) => [ + $crypto-key-type-choice ]
//	  }>
//	]
type KeyTriple struct {
	Environment Environment         `json:"environment"`
	VerifKeys   CryptoKeys          `json:"verification-keys"`
	Conditions  *KeyTripleCondition `json:"conditions,omitempty"`
}

func (o *KeyTriple) Valid() error {
	if err := o.Environment.Valid(); err != nil {
		return fmt.Errorf("environment: %w", err)
	}

	if err := o.VerifKeys.Valid(); err != nil {
		return fmt.Errorf("verification-keys: %w", err)
	}

	if o.Conditions != nil {
		if err := o.Conditions.Valid(); err != nil {
			return fmt.Errorf("conditions: %w", err)
		}
	}

	return nil
}

func (o *KeyTriple) MarshalCBOR() ([]byte, error) {
	toMarshal := []any{o.Environment, o.VerifKeys}
	if o.Conditions != nil {
		toMarshal = append(toMarshal, o.Conditions)
	}

	return em.Marshal(toMarshal)
}

// nolint:dupl
func (o *KeyTriple) UnmarshalCBOR(data []byte) error {
	var raw []cbor.RawMessage
	if err := dm.Unmarshal(data, &raw); err != nil {
		return err
	}

	numElts := len(raw)
	if numElts < 2 || numElts > 3 {
		return fmt.Errorf("expected array between 2 and 3 elements; found %d", numElts)
	}

	if err := dm.Unmarshal(raw[0], &o.Environment); err != nil {
		return fmt.Errorf("environment: %w", err)
	}

	if err := dm.Unmarshal(raw[1], &o.VerifKeys); err != nil {
		return fmt.Errorf("verification-keys: %w", err)
	}

	if numElts == 3 {
		if err := dm.Unmarshal(raw[2], &o.Conditions); err != nil {
			return fmt.Errorf("conditions: %w", err)
		}
	}

	return nil
}

type KeyTriples []KeyTriple

func NewKeyTriples() *KeyTriples {
	return &KeyTriples{}
}

func (o *KeyTriples) Add(it *KeyTriple) *KeyTriples {
	*o = append(*o, *it)
	return o
}
