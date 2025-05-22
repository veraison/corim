// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import "github.com/veraison/eat"

type ConciseEvidence struct {
	EvTriples  EvTriplesMap `cbor:"0,keyasint" json:"ev-triples"`
	EvidenceID EvidenceID   `cbor:"1,keyasint,omitempty" json:"ev-identity,omitempty"`
	Profile    *eat.Profile `cbor:"2,keyasint,omitempty" json:"profile,omitempty"`
}
