// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
)

type ConciseEvidence struct {
	EvTriples  EvTriples    `cbor:"0,keyasint" json:"ev-triples"`
	EvidenceID *EvidenceID  `cbor:"1,keyasint,omitempty" json:"evidence-id,omitempty"`
	Profile    *eat.Profile `cbor:"2,keyasint,omitempty" json:"profile,omitempty"`
	Extensions
}

// NewConciseEvidence instantiates an empty ConciseEvidence
func NewConciseEvidence() *ConciseEvidence {
	return &ConciseEvidence{}
}

// AddTriples adds Evidence Triples to Concise Evidence
func (o *ConciseEvidence) AddTriples(evTriples *EvTriples) error {
	if o != nil {
		if evTriples == nil {
			return errors.New("no evidence triples")
		}

		if err := evTriples.Valid(); err != nil {
			return fmt.Errorf("invalid evidence triples: %w", err)
		}
		o.EvTriples = *evTriples
	}
	return nil
}

// AddEvidenceID adds EvidenceID to ConciseEvidence
func (o *ConciseEvidence) AddEvidenceID(evidenceID *EvidenceID) error {
	if o != nil {
		if evidenceID == nil {
			return errors.New("no evidence id supplied")
		}
		if err := evidenceID.Valid(); err != nil {
			return fmt.Errorf("invalid EvidenceID: %w", err)
		}
		o.EvidenceID = evidenceID
	}
	return nil
}

// AddProfile adds a chosen profile to ConciseEvidence
func (o *ConciseEvidence) AddProfile(urlOrOID string) error {
	p, err := eat.NewProfile(urlOrOID)
	if err != nil {
		return err
	}
	o.Profile = p
	return nil
}

// nolint:gocritic
func (o ConciseEvidence) Valid() error {
	if err := o.EvTriples.Valid(); err != nil {
		return fmt.Errorf("invalid EvTriples: %w", err)
	}
	if o.EvidenceID != nil {
		if err := o.EvidenceID.Valid(); err != nil {
			return fmt.Errorf("invalid EvidenceID: %w", err)
		}
	}

	return nil
}

// RegisterExtensions registers a struct as a collections of extensions
func (o *ConciseEvidence) RegisterExtensions(exts extensions.Map) error {
	evTriplesMap := extensions.NewMap()
	for p, v := range exts {
		switch p {
		case ExtConciseEvidence:
			o.Register(v)
		default:
			evTriplesMap.Add(p, v)
		}
	}
	return o.EvTriples.RegisterExtensions(evTriplesMap)
}

// GetExtensions returns previously registered extension
func (o *ConciseEvidence) GetExtensions() extensions.IMapValue {
	return o.IMapValue
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
	if err := encoding.PopulateStructFromCBOR(dm, data, o); err != nil {
		return err
	}
	return o.Valid()
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
	if err := encoding.PopulateStructFromJSON(data, o); err != nil {
		return err
	}
	return o.Valid()
}

type TaggedConciseEvidence ConciseEvidence

// NewTagged Concise Evidence creates a Tagged Concise Evidence from a supplied Concise Evidence
func NewTaggedConciseEvidence(ev *ConciseEvidence) (*TaggedConciseEvidence, error) {
	var tce TaggedConciseEvidence
	if ev == nil {
		return nil, errors.New("non existent concise evidence")
	}
	if err := ev.Valid(); err != nil {
		return nil, fmt.Errorf("concise Evidence is not valid: %w", err)
	}
	tce = TaggedConciseEvidence(*ev)
	return &tce, nil
}

// Valid checks the validity of TaggedConciseEvidence
// nolint:gocritic
func (o TaggedConciseEvidence) Valid() error {
	c := ConciseEvidence(o)
	return c.Valid()
}

// ToCBOR serializes the target TaggedConciseEvidence to CBOR
// nolint:gocritic
func (o TaggedConciseEvidence) ToCBOR() ([]byte, error) {
	ce := ConciseEvidence(o)
	if err := ce.Valid(); err != nil {
		return nil, err
	}

	data, err := ce.ToCBOR()
	if err != nil {
		return nil, fmt.Errorf("unable to serialize the data: %w", err)
	}
	return append(ConciseEvidenceTag, data...), nil
}

// FromCBOR deserializes a CBOR-encoded date into the TaggedConciseEvidence
func (o *TaggedConciseEvidence) FromCBOR(data []byte) error {
	if len(data) < 3 {
		return errors.New("input CBOR data too short")
	}
	if !bytes.Equal(data[:3], ConciseEvidenceTag) {
		return errors.New("did not see concise evidence tag")
	}
	if err := encoding.PopulateStructFromCBOR(dm, data[3:], o); err != nil {
		return err
	}
	return o.Valid()

}

// FromJSON deserializes a JSON-encoded TaggedConciseEvidence into the target TaggedConciseEvidence
func (o *TaggedConciseEvidence) FromJSON(data []byte) error {
	return encoding.PopulateStructFromJSON(data, o)
}

// ToJSON serializes the target TaggedConciseEvidence to JSON
// nolint:gocritic
func (o TaggedConciseEvidence) ToJSON() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}
	return encoding.SerializeStructToJSON(&o)
}
