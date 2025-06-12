// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
)

type ConciseEvidence struct {
	CoEvTriples CoEvTriples  `cbor:"0,keyasint" json:"coev-triples"`
	EvidenceID  *EvidenceID  `cbor:"1,keyasint,omitempty" json:"ev-identity,omitempty"`
	Profile     *eat.Profile `cbor:"2,keyasint,omitempty" json:"profile,omitempty"`
	Extensions
}

// NewConciseEvidence instantiates an empty ConciseEvidence
func NewConciseEvidence() *ConciseEvidence {
	return &ConciseEvidence{}
}

func (o *ConciseEvidence) AddEvidenceTriples(coEvTriples *CoEvTriples) error {
	if coEvTriples == nil {
		return errors.New("no evidence triples")
	}

	if err := coEvTriples.Valid(); err != nil {
		return fmt.Errorf("invalid evidence triples: %w", err)
	}
	o.CoEvTriples = *coEvTriples

	return nil
}

func (o *ConciseEvidence) AddEvidenceID(evidenceID *EvidenceID) error {
	if evidenceID == nil {
		return errors.New("no evidence id supplied")
	}
	if err := evidenceID.Valid(); err != nil {
		return fmt.Errorf("invalid evidenceID: %w", err)
	}
	o.EvidenceID = evidenceID
	return nil
}

func (o *ConciseEvidence) AddProfile(profile *eat.Profile) error {
	if profile == nil {
		return errors.New("no profile supplied")
	}
	if !profile.IsOID() && !profile.IsURI() {
		return errors.New("profile should be OID or URI")
	}
	o.Profile = profile
	return nil
}

func (o ConciseEvidence) Valid() error {
	if err := o.CoEvTriples.Valid(); err != nil {
		return fmt.Errorf("invalid CoEvTriples: %w", err)
	}
	if o.EvidenceID != nil {
		if err := o.EvidenceID.Valid(); err != nil {
			return fmt.Errorf("invalid EvidenceID: %w", err)
		}
	}
	if o.Profile != nil {
		p := o.Profile
		if !p.IsOID() && !p.IsURI() {
			return errors.New("profile should be OID or URI")
		}
	}
	return nil
}

// RegisterExtensions registers a struct as a collections of extensions
func (o *ConciseEvidence) RegisterExtensions(exts extensions.Map) error {
	coEvMap := extensions.NewMap()
	for p, v := range exts {
		switch p {
		case ExtConciseEvidence:
			o.Extensions.Register(v)
		default:
			coEvMap.Add(p, v)
		}
	}
	return o.CoEvTriples.RegisterExtensions(coEvMap)
}

// ToCBOR serializes the target ConciseEvidence to CBOR
// nolint:gocritic
func (o ConciseEvidence) ToCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return encoding.SerializeStructToCBOR(em, &o)
}

// FromCBOR deserializes a CBOR-encoded ConciseEvidence into the target ConciseEvidence
func (o *ConciseEvidence) FromCBOR(data []byte) error {
	return encoding.PopulateStructFromCBOR(dm, data, o)
}

// ToJSON serializes the target ConciseEvidence to JSON
// nolint:gocritic
func (o ConciseEvidence) ToJSON() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return encoding.SerializeStructToJSON(&o)
}

// FromJSON deserializes a JSON-encoded ConciseEvidence into the target ConciseEvidence
func (o *ConciseEvidence) FromJSON(data []byte) error {
	return encoding.PopulateStructFromJSON(data, o)
}

// nolint:gocritic
func (o ConciseEvidence) ToJSONPretty(indent string) ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return json.MarshalIndent(&o, "", indent)
}
