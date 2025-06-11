// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
)

// Reference Measurements within a CoSWID is contained in the payload-entry inside Evidence
// while the Evidence is carried in the evidence-entry field inside Evidence
// The CoSWIDEvidenceMap MAY contain a concise-swid-tag-id
// and zero or more $crypto-key-type-choice values that identify
// entities authorized to provide Reference Measurements
type CoSWIDEvidenceMap struct {
	TagID        *swid.TagID      `cbor:"0,keyasint,omitempty" json:"tagId,omitempty"`
	Evidence     swid.Evidence    `cbor:"1,keyasint,omitempty" json:"evidence,omitempty"`
	AuthorizedBy *comid.CryptoKey `cbor:"2,keyasint,omitempty" json:"authorized-by,omitempty"`
}

type CoSWIDEvidence []CoSWIDEvidenceMap
