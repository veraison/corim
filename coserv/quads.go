// Copyright 2025-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"errors"
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/cots"
)

type RefValQuad struct {
	Authorities *comid.CryptoKeys  `cbor:"1,keyasint"`
	RVTriple    *comid.ValueTriple `cbor:"2,keyasint"`
}

type EndValQuad struct {
	Authorities *comid.CryptoKeys  `cbor:"1,keyasint"`
	EVTriple    *comid.ValueTriple `cbor:"2,keyasint"`
}

type CondEndValQuad struct {
	Authorities *comid.CryptoKeys        `cbor:"1,keyasint"`
	CETriple    *comid.CondEndorseTriple `cbor:"2,keyasint"`
}

type AKQuad struct {
	Authorities *comid.CryptoKeys `cbor:"1,keyasint"`
	AKTriple    *comid.KeyTriple  `cbor:"2,keyasint"`
}

type CoTSStmt struct {
	Authorities *comid.CryptoKeys    `cbor:"1,keyasint"`
	CoTS        *cots.ConciseTaStore `cbor:"2,keyasint"`
}

func validateAuthorities(authorities *comid.CryptoKeys) error {
	if authorities == nil {
		return errors.New("missing authorities")
	}

	if err := authorities.Valid(); err != nil {
		return fmt.Errorf("invalid authorities: %w", err)
	}

	return nil
}

// Valid checks that the RefValQuad has authorities and a valid triple
func (o RefValQuad) Valid() error {
	if err := validateAuthorities(o.Authorities); err != nil {
		return err
	}

	if o.RVTriple == nil {
		return errors.New("missing reference value triple")
	}

	if err := o.RVTriple.Valid(); err != nil {
		return fmt.Errorf("invalid reference value triple: %w", err)
	}

	return nil
}

// Valid checks that the EndValQuad has authorities and a valid triple
func (o EndValQuad) Valid() error {
	if err := validateAuthorities(o.Authorities); err != nil {
		return err
	}

	if o.EVTriple == nil {
		return errors.New("missing endorsed value triple")
	}

	if err := o.EVTriple.Valid(); err != nil {
		return fmt.Errorf("invalid endorsed value triple: %w", err)
	}

	return nil
}

// Valid checks that the CondEndValQuad has authorities and a valid triple
func (o CondEndValQuad) Valid() error {
	if err := validateAuthorities(o.Authorities); err != nil {
		return err
	}

	if o.CETriple == nil {
		return errors.New("missing conditional endorsement triple")
	}

	if err := o.CETriple.Valid(); err != nil {
		return fmt.Errorf("invalid conditional endorsement triple: %w", err)
	}

	return nil
}

// Valid checks that the AKQuad has authorities and a valid triple
func (o AKQuad) Valid() error {
	if err := validateAuthorities(o.Authorities); err != nil {
		return err
	}

	if o.AKTriple == nil {
		return errors.New("missing attestation key triple")
	}

	if err := o.AKTriple.Valid(); err != nil {
		return fmt.Errorf("invalid attestation key triple: %w", err)
	}

	return nil
}

// Valid checks that the CoTSStmt has authorities and a valid CoTS
func (o CoTSStmt) Valid() error {
	if err := validateAuthorities(o.Authorities); err != nil {
		return err
	}

	if o.CoTS == nil {
		return errors.New("missing CoTS")
	}

	if err := o.CoTS.Valid(); err != nil {
		return fmt.Errorf("invalid CoTS: %w", err)
	}

	return nil
}
