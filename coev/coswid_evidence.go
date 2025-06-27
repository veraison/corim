// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
)

// CoSWIDEvidenceMap is the Map to carry CoSWID Evidence
type CoSWIDEvidenceMap struct {
	TagID        *swid.TagID      `cbor:"0,keyasint,omitempty" json:"tagId,omitempty"`
	Evidence     swid.Evidence    `cbor:"1,keyasint,omitempty" json:"evidence,omitempty"`
	AuthorizedBy *comid.CryptoKey `cbor:"2,keyasint,omitempty" json:"authorized-by,omitempty"`
}

type CoSWIDEvidence []CoSWIDEvidenceMap

func NewCoSWIDEvidence() *CoSWIDEvidence {
	return &CoSWIDEvidence{}
}

func (o *CoSWIDEvidence) AddCoSWIDEvidenceMap(e *CoSWIDEvidenceMap) *CoSWIDEvidence {
	if o == nil {
		o = NewCoSWIDEvidence()
	}
	*o = append(*o, *e)
	return o
}
