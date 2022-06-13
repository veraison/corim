// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cots

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ConciseTaStore struct {
	Language     *string            `cbor:"0,keyasint,omitempty" json:"language,omitempty"`
	Environments EnvironmentGroups  `cbor:"1,keyasint" json:"environments"`
	Purposes     []string           `cbor:"2,keyasint,omitempty" json:"purposes,omitempty"`
	PermClaims   EatCWTClaims     	`cbor:"3,keyasint,omitempty" json:"permclaims,omitempty"`
	ExclClaims   EatCWTClaims       	`cbor:"4,keyasint,omitempty" json:"exclclaims,omitempty"`
	Keys         *TasAndCas         `cbor:"5,keyasint" json:"keys"`
}

func NewConciseTaStore() *ConciseTaStore {
	return &ConciseTaStore{}
}

func (o *ConciseTaStore) SetLanguage(language string) *ConciseTaStore {
	if o != nil {
		o.Language = &language
	}
	return o
}

func (o ConciseTaStore) GetLanguage() string {
	if o.Language == nil {
		return ""
	}
	return *o.Language
}

func (o *ConciseTaStore) AddEnvironmentGroup(eg EnvironmentGroup) *ConciseTaStore {
	if o != nil {
		o.Environments = append(o.Environments, eg)
	}
	return o
}

func (o ConciseTaStore) GetEnvironments() []EnvironmentGroup {
	return o.Environments
}

func (o *ConciseTaStore) AddPurpose(purpose string) *ConciseTaStore {
	if o != nil {
		o.Purposes = append(o.Purposes, purpose)
	}
	return o
}

func (o ConciseTaStore) GetPurposes() []string {
	return o.Purposes
}

func (o *ConciseTaStore) AddPermClaims(permclaim EatCWTClaim) *ConciseTaStore {
	if o != nil {
		o.PermClaims = append(o.PermClaims, permclaim)
	}
	return o
}

func (o ConciseTaStore) GetPermClaims() EatCWTClaims {
	return o.PermClaims
}

func (o *ConciseTaStore) AddExclClaims(exclclaim EatCWTClaim) *ConciseTaStore {
	if o != nil {
		o.ExclClaims = append(o.ExclClaims, exclclaim)
	}
	return o
}

func (o ConciseTaStore) GetExclClaims() EatCWTClaims {
	return o.ExclClaims
}

func (o ConciseTaStore) GetKeys() TasAndCas {
	return *o.Keys
}

func (o *ConciseTaStore) SetKeys(keys TasAndCas) *ConciseTaStore {
	if o != nil {
		o.Keys = &keys
	}
	return o
}

// ToCBOR serializes the target ConciseTaStore to CBOR
func (o ConciseTaStore) ToCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return em.Marshal(&o)
}

// FromCBOR deserializes a CBOR-encoded CoMID into the target Comid
func (o *ConciseTaStore) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, o)
}

// Valid iterates over the range of individual entities to check for validity
func (o ConciseTaStore) Valid() error {
	if nil == o.Keys || len(o.Keys.Tas) == 0 {
		return fmt.Errorf("empty Keys")
	}

	return nil
}

// FromJSON deserializes the supplied JSON data into the target Meta
func (o *ConciseTaStore) FromJSON(data []byte) error {
	return json.Unmarshal(data, o)
}

// ToJSON serializes the target Meta to JSON
func (o ConciseTaStore) ToJSON() ([]byte, error) {
	return json.Marshal(&o)
}

type ConciseTaStores []ConciseTaStore

func NewConciseTaStores() *ConciseTaStores {
	return new(ConciseTaStores)
}

func (o *ConciseTaStores) AddConciseTaStores(cts ConciseTaStore) *ConciseTaStores {
	if o != nil {
		if cts.Valid() != nil {
			return nil
		}

		*o = append(*o, cts)
	}
	return o
}

func (o ConciseTaStores) ToCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return em.Marshal(&o)
}

// FromCBOR deserializes a CBOR-encoded ConciseTaStores into the target ConciseTaStores
func (o *ConciseTaStores) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, o)
}

// FromJSON deserializes the supplied JSON data into the target deserializes a
func (o *ConciseTaStores) FromJSON(data []byte) error {
	return json.Unmarshal(data, o)
}

// ToJSON serializes the target deserializes a  to JSON
func (o ConciseTaStores) ToJSON() ([]byte, error) {
	return json.Marshal(&o)
}

func (o ConciseTaStores) Valid() error {
	if len(o) == 0 {
		return errors.New("empty concise-ta-stores")
	}
	return nil
}
