// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import "github.com/veraison/corim/extensions"

// IndirectMap models the spdm-indirect-map: { 0: [+ uint] }.
// It carries SPDM measurement index numbers that point into the SPDM
// measurement block rather than embedding measurement hashes directly.
type IndirectMap struct {
	Index []uint64 `cbor:"0,keyasint" json:"index"`
}

// MValExtensions extends a CoEV measurement value with the spdm-indirect
// field at CBOR integer key 12, as defined in the TCG DICE CoEV
// comid-extensions.cddl.
type MValExtensions struct {
	SPDMIndirect *IndirectMap `cbor:"12,keyasint,omitempty" json:"spdm-indirect,omitempty"`
}

// SpdmExtensionMap returns an extensions.Map that registers the SPDM
// measurement value extension on a ConciseEvidence, enabling recognition
// of the spdm-indirect field during encoding and decoding.
func SpdmExtensionMap() extensions.Map {
	return extensions.NewMap().Add(ExtEvidenceTriples, &MValExtensions{})
}
