// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"errors"
	"time"

	"github.com/veraison/cmw"
)

type ResultSet struct {
	RVQ *[]RefValQuad `cbor:"0,keyasint,omitempty"`
	AKQ *[]AKQuad     `cbor:"3,keyasint,omitempty"`
	TAS *[]CoTSStmt   `cbor:"4,keyasint,omitempty"`
	// TODO(tho) add endorsed values
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

// AddAttestationKeys adds the supplied ak quad to the target ResultSet
func (o *ResultSet) AddAttestationKeys(v AKQuad) *ResultSet {
	if o.AKQ == nil {
		o.AKQ = new([]AKQuad)
	}

	*o.AKQ = append(*o.AKQ, v)

	return o
}

// AddCoTS adds the supplied CoTS statement to the target ResultSet
func (o *ResultSet) AddCoTS(v CoTSStmt) *ResultSet {
	if o.TAS == nil {
		o.TAS = new([]CoTSStmt)
	}

	*o.TAS = append(*o.TAS, v)

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
	// The coherency between query and results must be checked by the Coserv's
	// Valid()
	return nil
}
