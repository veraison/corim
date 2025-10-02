// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"encoding/json"
	"fmt"
	
	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/corim/comid"
)

// PSASwRelType defines the type of software relationship (updates or patches)
type PSASwRelType int

const (
	PSAUpdates PSASwRelType = 1
	PSAPatches PSASwRelType = 2
)

// PSASwRel defines a software relationship with type and security criticality
type PSASwRel struct {
	Type             PSASwRelType `cbor:"0,keyasint" json:"type"`
	SecurityCritical bool         `cbor:"1,keyasint" json:"security-critical"`
}

// Valid validates the PSA software relationship
func (o PSASwRel) Valid() error {
	switch o.Type {
	case PSAUpdates, PSAPatches:
		return nil
	default:
		return fmt.Errorf("invalid PSA software relationship type: %d (must be 1 for updates or 2 for patches)", o.Type)
	}
}

// PSASwRelationship represents a complete software relationship record
type PSASwRelationship struct {
	New      *comid.Measurement `cbor:"0,keyasint" json:"new"`
	Relation *PSASwRel          `cbor:"1,keyasint" json:"rel"`
	Old      *comid.Measurement `cbor:"2,keyasint" json:"old"`
}

// Valid validates the PSA software relationship record
func (o PSASwRelationship) Valid() error {
	if o.New == nil {
		return fmt.Errorf("new measurement is required")
	}
	if o.Relation == nil {
		return fmt.Errorf("relationship definition is required")
	}
	if o.Old == nil {
		return fmt.Errorf("old measurement is required")
	}

	if err := o.Relation.Valid(); err != nil {
		return fmt.Errorf("invalid relationship: %w", err)
	}

	return nil
}

// PSASwRelTriple represents a PSA software relationship triple
type PSASwRelTriple struct {
	Environment  comid.Environment  `cbor:"0,keyasint" json:"environment"`
	Relationship *PSASwRelationship `cbor:"1,keyasint" json:"relationship"`
}

// Valid validates the PSA software relationship triple
func (o PSASwRelTriple) Valid() error {
	if err := o.Environment.Valid(); err != nil {
		return fmt.Errorf("invalid environment: %w", err)
	}

	if o.Relationship == nil {
		return fmt.Errorf("relationship is required")
	}

	if err := o.Relationship.Valid(); err != nil {
		return fmt.Errorf("invalid relationship: %w", err)
	}

	return nil
}

// PSASwRelTriples represents a collection of PSA software relationship triples
type PSASwRelTriples struct {
	Values []PSASwRelTriple `cbor:"-" json:"-"`
}

// MarshalCBOR implements CBOR marshaling for PSASwRelTriples
func (o PSASwRelTriples) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(o.Values)
}

// UnmarshalCBOR implements CBOR unmarshaling for PSASwRelTriples
func (o *PSASwRelTriples) UnmarshalCBOR(data []byte) error {
	return cbor.Unmarshal(data, &o.Values)
}

// MarshalJSON implements JSON marshaling for PSASwRelTriples
func (o PSASwRelTriples) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Values)
}

// UnmarshalJSON implements JSON unmarshaling for PSASwRelTriples
func (o *PSASwRelTriples) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &o.Values)
}

// Add adds a PSA software relationship triple to the collection
func (o *PSASwRelTriples) Add(triple PSASwRelTriple) *PSASwRelTriples {
	if o != nil {
		o.Values = append(o.Values, triple)
	}
	return o
}

// Valid validates all PSA software relationship triples in the collection
func (o PSASwRelTriples) Valid() error {
	if len(o.Values) == 0 {
		return fmt.Errorf("at least one PSA software relationship triple is required")
	}

	for i, triple := range o.Values {
		if err := triple.Valid(); err != nil {
			return fmt.Errorf("invalid PSA software relationship triple at index %d: %w", i, err)
		}
	}

	return nil
}

// NewPSASwRelTriples creates a new PSA software relationship triples collection
func NewPSASwRelTriples() *PSASwRelTriples {
	return &PSASwRelTriples{
		Values: []PSASwRelTriple{},
	}
}

// Helper functions for creating software relationships

// NewPSAUpdateRelationship creates a new PSA update relationship
func NewPSAUpdateRelationship(newMeasurement, oldMeasurement *comid.Measurement, securityCritical bool) (*PSASwRelationship, error) {
	relation := &PSASwRel{
		Type:             PSAUpdates,
		SecurityCritical: securityCritical,
	}

	rel := &PSASwRelationship{
		New:      newMeasurement,
		Relation: relation,
		Old:      oldMeasurement,
	}

	if err := rel.Valid(); err != nil {
		return nil, err
	}

	return rel, nil
}

// NewPSAPatchRelationship creates a new PSA patch relationship
func NewPSAPatchRelationship(newMeasurement, oldMeasurement *comid.Measurement, securityCritical bool) (*PSASwRelationship, error) {
	relation := &PSASwRel{
		Type:             PSAPatches,
		SecurityCritical: securityCritical,
	}

	rel := &PSASwRelationship{
		New:      newMeasurement,
		Relation: relation,
		Old:      oldMeasurement,
	}

	if err := rel.Valid(); err != nil {
		return nil, err
	}

	return rel, nil
}
