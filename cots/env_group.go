// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cots

import (
	"encoding/json"
	"fmt"

	"github.com/veraison/corim/comid"
)

// EnvironmentGroup is the top-level representation of the unsigned-corim-map with
// CBOR and JSON serialization.
type EnvironmentGroup struct {
	Environment  *comid.Environment  `cbor:"1,keyasint,omitempty" json:"environment,omitempty"`
	SwidTag      *AbbreviatedSwidTag `cbor:"2,keyasint,omitempty" json:"swidtag,omitempty"`
	NamedTaStore *string             `cbor:"3,keyasint,omitempty" json:"namedtastore,omitempty"`
}

// NewEnvironmentGroup instantiates an empty EnvironmentGroup
func NewEnvironmentGroup() *EnvironmentGroup {
	return &EnvironmentGroup{}
}

func (o *EnvironmentGroup) SetEnvironment(environment comid.Environment) *EnvironmentGroup {
	if o != nil {
		o.Environment = &environment
	}
	return o
}

func (o *EnvironmentGroup) SetAbbreviatedSwidTag(swidtag *AbbreviatedSwidTag) *EnvironmentGroup {
	if o != nil {
		o.SwidTag = swidtag
	}
	return o
}

func (o *EnvironmentGroup) SetNamedTaStore(namedtastore string) *EnvironmentGroup {
	if o != nil {
		o.NamedTaStore = &namedtastore
	}
	return o
}

// Valid checks the validity of the target EnvironmentGroup
func (o EnvironmentGroup) Valid() error {
	if o.Environment != nil {
		if err := o.Environment.Valid(); err != nil {
			return fmt.Errorf("environment group validation failed: %w", err)
		}
	}

	if o.SwidTag != nil {
		if err := o.SwidTag.Valid(); err != nil {
			return fmt.Errorf("abbreviated swid tag validation failed: %w", err)
		}
	}

	return nil
}

// ToCBOR serializes the target EnvironmentGroup to CBOR
func (o EnvironmentGroup) ToCBOR() ([]byte, error) {
	return em.Marshal(&o)
}

// FromCBOR deserializes a CBOR-encoded data into the target EnvironmentGroup
func (o *EnvironmentGroup) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, o)
}

// ToJSON serializes the target EnvironmentGroup to JSON
func (o EnvironmentGroup) ToJSON() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return json.Marshal(&o)
}

// FromJSON deserializes a JSON-encoded data into the target EnvironmentGroup
func (o *EnvironmentGroup) FromJSON(data []byte) error {
	return json.Unmarshal(data, o)
}

type EnvironmentGroups []EnvironmentGroup

func NewEnvironmentGroups() *EnvironmentGroups {
	return new(EnvironmentGroups)
}

func (o *EnvironmentGroups) AddEnvironmentGroup(e EnvironmentGroup) *EnvironmentGroups {
	if o != nil {
		*o = append(*o, e)
	}
	return o
}

// An empty list signifies all contexts SHOULD be considered as applicable
func (o EnvironmentGroups) Valid() error {
	for i, e := range o {
		if err := e.Valid(); err != nil {
			return fmt.Errorf("bad environment group at index %d: %w", i, err)
		}
	}
	return nil
}

// FromJSON deserializes a JSON-encoded data into the target EnvironmentGroup
func (o *EnvironmentGroups) FromJSON(data []byte) error {
	return json.Unmarshal(data, o)
}

// ToJSON serializes the target EnvironmentGroup to JSON
func (o EnvironmentGroups) ToJSON() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return json.Marshal(&o)
}
