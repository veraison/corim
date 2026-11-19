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

func validateQuad(
	authorities *comid.CryptoKeys,
	triple interface{ Valid() error },
	missing bool,
	name string,
) error {
	if authorities == nil {
		return errors.New("missing authorities")
	}

	if err := authorities.Valid(); err != nil {
		return fmt.Errorf("invalid authorities: %w", err)
	}

	if missing {
		return fmt.Errorf("missing %s", name)
	}

	if err := triple.Valid(); err != nil {
		return fmt.Errorf("invalid %s: %w", name, err)
	}

	return nil
}

// Valid checks that the RefValQuad has authorities and a valid triple
func (o RefValQuad) Valid() error {
	return validateQuad(o.Authorities, o.RVTriple, o.RVTriple == nil, "reference value triple")
}

// Valid checks that the EndValQuad has authorities and a valid triple
func (o EndValQuad) Valid() error {
	return validateQuad(o.Authorities, o.EVTriple, o.EVTriple == nil, "endorsed value triple")
}

// Valid checks that the CondEndValQuad has authorities and a valid triple
func (o CondEndValQuad) Valid() error {
	return validateQuad(o.Authorities, o.CETriple, o.CETriple == nil, "conditional endorsement triple")
}

// Valid checks that the AKQuad has authorities and a valid triple
func (o AKQuad) Valid() error {
	return validateQuad(o.Authorities, o.AKTriple, o.AKTriple == nil, "attestation key triple")
}

// Valid checks that the CoTSStmt has authorities and a valid CoTS
func (o CoTSStmt) Valid() error {
	return validateQuad(o.Authorities, o.CoTS, o.CoTS == nil, "CoTS")
}
