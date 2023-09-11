// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cots

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/veraison/corim/v2/comid"
	"github.com/veraison/swid"
)

type ConciseTaStore struct {
	Language     *string            `cbor:"0,keyasint,omitempty" json:"language,omitempty"`
	TagIdentity  *comid.TagIdentity `cbor:"1,keyasint,omitempty" json:"tag-identity,omitempty"`
	Environments EnvironmentGroups  `cbor:"2,keyasint" json:"environments"`
	Purposes     []string           `cbor:"3,keyasint,omitempty" json:"purposes,omitempty"`
	PermClaims   EatCWTClaims       `cbor:"4,keyasint,omitempty" json:"permclaims,omitempty"`
	ExclClaims   EatCWTClaims       `cbor:"5,keyasint,omitempty" json:"exclclaims,omitempty"`
	Keys         *TasAndCas         `cbor:"6,keyasint" json:"keys"`
}

func NewConciseTaStore() *ConciseTaStore {
	return &ConciseTaStore{}
}

func (o *ConciseTaStore) SetTagIdentity(tagID interface{}, tagIDVersion *uint) *ConciseTaStore {
	if o != nil {
		id := swid.NewTagID(tagID)
		if id == nil {
			return nil
		}
		o.TagIdentity = &comid.TagIdentity{}
		o.TagIdentity.TagID = *id
		if nil != tagIDVersion {
			o.TagIdentity.TagVersion = *tagIDVersion
		}
	}
	return o
}

func (o *ConciseTaStore) SetLanguage(language string) *ConciseTaStore {
	if o != nil {
		o.Language = &language
	}
	return o
}

func (o *ConciseTaStore) AddEnvironmentGroup(eg EnvironmentGroup) *ConciseTaStore {
	if o != nil {
		o.Environments = append(o.Environments, eg)
	}
	return o
}

func (o *ConciseTaStore) AddPurpose(purpose string) *ConciseTaStore {
	if o != nil {
		o.Purposes = append(o.Purposes, purpose)
	}
	return o
}

func (o *ConciseTaStore) AddPermClaims(permclaim EatCWTClaim) *ConciseTaStore {
	if o != nil {
		o.PermClaims = append(o.PermClaims, permclaim)
	}
	return o
}

func (o *ConciseTaStore) AddExclClaims(exclclaim EatCWTClaim) *ConciseTaStore {
	if o != nil {
		o.ExclClaims = append(o.ExclClaims, exclclaim)
	}
	return o
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

// FromCBOR deserializes a CBOR-encoded CoTS into the target ConciseTaStore
func (o *ConciseTaStore) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, o)
}

// Valid iterates over the range of individual entities to check for validity
func (o ConciseTaStore) Valid() error {
	if o.Environments == nil {
		return fmt.Errorf("environmentGroups must be present")
	}
	if len(o.Environments) != 0 {
		if err := o.Environments.Valid(); err != nil {
			return fmt.Errorf("invalid environmentGroups: %w", err)
		}
	}

	if o.TagIdentity != nil {
		if err := o.TagIdentity.Valid(); err != nil {
			return fmt.Errorf("invalid TagIdentity: %w", err)
		}
	}

	if o.Keys == nil || len(o.Keys.Tas) == 0 {
		return fmt.Errorf("empty Keys")
	}

	return nil
}

// FromJSON deserializes a JSON-encoded CoTS into the target ConciseTaStore
func (o *ConciseTaStore) FromJSON(data []byte) error {
	return json.Unmarshal(data, o)
}

// FromJSON deserializes a JSON-encoded CoTS into the target ConsiseTaStore
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

// ToCBOR serializes the target ConciseTaStores to CBOR
func (o ConciseTaStores) ToCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return em.Marshal(&o)
}

// FromCBOR deserializes a CBOR-encoded CoTS into the target ConsiseTaStores
func (o *ConciseTaStores) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, o)
}

// FromJSON deserializes a JSON-encoded CoTS into the target ConsiseTaStores
func (o *ConciseTaStores) FromJSON(data []byte) error {
	return json.Unmarshal(data, o)
}

// ToJSON serializes the target ConsiseTaStore to JSON
func (o ConciseTaStores) ToJSON() ([]byte, error) {
	return json.Marshal(&o)
}

func (o ConciseTaStores) Valid() error {
	if len(o) == 0 {
		return errors.New("empty concise-ta-stores")
	}

	for i, c := range o {
		if err := c.Valid(); err != nil {
			return fmt.Errorf("bad ConciseTaStore group at index %d: %w", i, err)
		}
	}
	return nil
}
