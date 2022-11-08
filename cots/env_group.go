// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cots

import (
	"encoding/json"
	"errors"

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

func (o EnvironmentGroup) GetEnvironment() comid.Environment {
	return *o.Environment
}

func (o *EnvironmentGroup) SetAbbreviatedSwidTag(swidtag AbbreviatedSwidTag) *EnvironmentGroup {
	if o != nil {
		o.SwidTag = &swidtag
	}
	return o
}

func (o EnvironmentGroup) GetAbbreviatedSwidTag() AbbreviatedSwidTag {
	return *o.SwidTag
}

func (o *EnvironmentGroup) SetNamedTaStore(namedtastore string) *EnvironmentGroup {
	if o != nil {
		o.NamedTaStore = &namedtastore
	}
	return o
}

func (o EnvironmentGroup) GetNamedTaStore() string {
	if o.NamedTaStore == nil {
		return ""
	}
	return *o.NamedTaStore
}

// Valid checks the validity (according to the spec) of the target unsigned CoRIM
func (o EnvironmentGroup) Valid() error {
	//TODO validation
	return nil
}

// ToCBOR serializes the target unsigned CoRIM to CBOR
func (o EnvironmentGroup) ToCBOR() ([]byte, error) {
	return em.Marshal(&o)
}

// FromCBOR deserializes a CBOR-encoded unsigned CoRIM into the target EnvironmentGroup
func (o *EnvironmentGroup) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, o)
}

// ToJSON serializes the target Comid to JSON
func (o EnvironmentGroup) ToJSON() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return json.Marshal(&o)
}

// FromJSON deserializes a JSON-encoded unsigned CoRIM into the target EnvironmentGroup
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

func (o EnvironmentGroups) Valid() error {
	if len(o) == 0 {
		return errors.New("empty EnvironmentGroups")
	}
	return nil
}

// FromJSON deserializes a JSON-encoded CoMID into the target Comid
func (o *EnvironmentGroups) FromJSON(data []byte) error {
	return json.Unmarshal(data, o)
}

// ToJSON serializes the target Comid to JSON
func (o EnvironmentGroups) ToJSON() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return json.Marshal(&o)
}
