// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import "github.com/veraison/corim/comid"

type RefValQuad struct {
	Authorities *comid.CryptoKeys  `cbor:"1,keyasint"`
	RVTriple    *comid.ValueTriple `cbor:"2,keyasint"`
}

// EndorsedValQuad represents an endorsed-values result quad as per CoSERV
// It mirrors RefValQuad but carries endorsed values instead of reference values
type EndorsedValQuad struct {
	Authorities *[]comid.CryptoKey `cbor:"1,keyasint"`
	EVTriple    *comid.ValueTriple `cbor:"2,keyasint"`
}

type AKQuad struct {
	Authorities *comid.CryptoKeys `cbor:"1,keyasint"`
	AKTriple    *comid.KeyTriple  `cbor:"2,keyasint"`
}
