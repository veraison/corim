// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"errors"
	"fmt"

	"github.com/veraison/swid"
)

// RimSelectorType specifies the type the RIM identifier inside a RimSelector.
type RimSelectorType uint

const (
	RimSelectorTypeComid RimSelectorType = iota
	RimSelectorTypeCoswid
	RimSelectorTypeCorim
)

// RimSelectorID specifies the type and ID of RIM elements (CoRIMs, CoSWIDs, or
// CoMIDs) that will be selected by a CoSERV query. it is mutually exclusive
// with EnvironmentSelector.
type RimSelectorID struct {
	_     struct{}        `cbor:",toarray"`
	Type  RimSelectorType `json:"type"`
	TagID swid.TagID      `json:"tag-id"`
}

// NewRimSelectorID creates a new RimSelectorID with the specified type and ID
// values.
func NewRimSelectorID(typ RimSelectorType, tagID swid.TagID) (*RimSelectorID, error) {
	ret := &RimSelectorID{Type: typ, TagID: tagID}
	if err := ret.Valid(); err != nil {
		return nil, err
	}

	return ret, nil
}

// Valid returns an error if the RimSelectorID is invalid.
func (o RimSelectorID) Valid() error {
	if err := o.TagID.Valid(); err != nil {
		return fmt.Errorf("tag-id: %w", err)
	}

	if o.Type > RimSelectorTypeCorim {
		return fmt.Errorf("invalid type: %d", o.Type)
	}

	return nil
}

// RimSelectorIDs is a collection of RimSelectorID.
type RimSelectorIDs []*RimSelectorID

// NewRimSelectorIDs creates a new empty NewRimSelectorIDs.
func NewRimSelectorIDs() *RimSelectorIDs {
	return &RimSelectorIDs{}
}

// Add the specified RimSelectorID to the collection.
func (o *RimSelectorIDs) Add(val *RimSelectorID) *RimSelectorIDs {
	*o = append(*o, val)
	return o
}

// Valid returns an error if the RimSelectorIDs is invalid.
func (o RimSelectorIDs) Valid() error {
	if len(o) == 0 {
		return errors.New("empty")
	}

	for i, rsid := range o {
		if err := rsid.Valid(); err != nil {
			return fmt.Errorf("RIM selector ID[%d]: %w", i, err)
		}
	}

	return nil
}
