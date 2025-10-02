// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import "github.com/veraison/corim/comid"

type RefValQuad struct {
	Authorities *comid.CryptoKeys  `cbor:"1,keyasint"`
	RVTriple    *comid.ValueTriple `cbor:"2,keyasint"`
}

type AKQuad struct {
	Authorities *comid.CryptoKeys `cbor:"1,keyasint"`
	AKTriple    *comid.KeyTriple  `cbor:"2,keyasint"`
}
