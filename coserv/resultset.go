// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"errors"
	"time"

	"github.com/veraison/cmw"
)

type ResultSet struct {
	RVQ *[]RefValQuad      `cbor:"0,keyasint,omitempty"`
	EVQ *[]EndorsedValQuad `cbor:"1,keyasint,omitempty"`
	AKQ *[]AKQuad          `cbor:"3,keyasint,omitempty"`
	// TODO(tho) add CoTS
	Expiry          *time.Time `cbor:"10,keyasint"`
	SourceArtifacts *[]cmw.CMW `cbor:"11,keyasint,omitempty"`
}

// NewResultSet instantiates a new ResultSet
func NewResultSet() *ResultSet {
	return &ResultSet{}
}

// AddReferenceValues adds the supplied ref-val quad to the target ResultSet
func (o *ResultSet) AddReferenceValues(v RefValQuad) *ResultSet {
	if o.RVQ == nil {
		o.RVQ = new([]RefValQuad)
	}

	*o.RVQ = append(*o.RVQ, v)

	return o
}

// AddEndorsedValues adds the supplied endorsed-values quad to the target ResultSet
func (o *ResultSet) AddEndorsedValues(v EndorsedValQuad) *ResultSet {
	if o.EVQ == nil {
		o.EVQ = new([]EndorsedValQuad)
	}

	*o.EVQ = append(*o.EVQ, v)

	return o
}

// AddAttestationKeys adds the supplied ak quad to the target ResultSet
func (o *ResultSet) AddAttestationKeys(v AKQuad) *ResultSet {
	if o.AKQ == nil {
		o.AKQ = new([]AKQuad)
	}

	*o.AKQ = append(*o.AKQ, v)

	return o
}

// AddSourceArtifacts adds the supplied CMW to the target ResultSet
func (o *ResultSet) AddSourceArtifacts(v cmw.CMW) *ResultSet { // nolint:gocritic
	if o.SourceArtifacts == nil {
		o.SourceArtifacts = new([]cmw.CMW)
	}

	*o.SourceArtifacts = append(*o.SourceArtifacts, v)

	return o
}

// SetExpiry sets the Expiry attribute of the target ResultSet to the supplied time
func (o *ResultSet) SetExpiry(exp time.Time) *ResultSet {
	o.Expiry = &exp
	return o
}

// Valid checks that the supplied ResultSet is syntactically correct
func (o ResultSet) Valid() error {
	if o.Expiry == nil {
		return errors.New("missing mandatory expiry")
	}
	// Nothing else to validate structurally here; combinations are checked at Coserv level
	// The coherency between query and results must be checked by the Coserv's
	// Valid()
	return nil
}
