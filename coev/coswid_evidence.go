// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"errors"
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
)

// CoSWIDEvidenceMap is the Map to carry CoSWID Evidence
type CoSWIDEvidenceMap struct {
	TagID        *swid.TagID       `cbor:"0,keyasint,omitempty" json:"tagId,omitempty"`
	Evidence     swid.Evidence     `cbor:"1,keyasint,omitempty" json:"evidence,omitempty"`
	AuthorizedBy *comid.CryptoKeys `cbor:"2,keyasint,omitempty" json:"authorized-by,omitempty"`
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

// Valid validates the CoSWIDEvidenceMap structure
func (o *CoSWIDEvidenceMap) Valid() error {
	// Validate TagID if present
	if o.TagID != nil {
		if err := o.TagID.Valid(); err != nil {
			return fmt.Errorf("tagId validation failed: %w", err)
		}
	}

	// Validate Evidence
	if err := o.Evidence.Valid(); err != nil {
		return fmt.Errorf("evidence validation failed: %w", err)
	}

	return nil
}

// Valid validates all CoSWIDEvidence entries
func (o CoSWIDEvidence) Valid() error {
	if len(o) == 0 {
		return errors.New("must contain at least one entry")
	}

	for i := range o {
		if err := o[i].Valid(); err != nil {
			return fmt.Errorf("evidence[%d] validation failed: %w", i, err)
		}
	}

	return nil
}
