// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"errors"
	"time"

	"github.com/veraison/corim/comid"
)

type ResultSet struct {
	Expiry          *time.Time           `cbor:"10,keyasint"`
	ReferenceValues *[]comid.ValueTriple `cbor:"0,keyasint,omitempty"`
	AttestationKeys *[]comid.KeyTriple   `cbor:"3,keyasint,omitempty"`
	// TODO(tho) all other supported types
}

func NewResultSet() *ResultSet {
	return &ResultSet{}
}

func (o *ResultSet) AddReferenceValues(v comid.ValueTriple) *ResultSet {
	if o.ReferenceValues == nil {
		o.ReferenceValues = new([]comid.ValueTriple)
	}

	*o.ReferenceValues = append(*o.ReferenceValues, v)

	return o
}

func (o *ResultSet) AddAttestationKeys(v comid.KeyTriple) *ResultSet {
	if o.AttestationKeys == nil {
		o.AttestationKeys = new([]comid.KeyTriple)
	}

	*o.AttestationKeys = append(*o.AttestationKeys, v)

	return o
}

func (o *ResultSet) SetExpiry(exp time.Time) *ResultSet {
	o.Expiry = &exp
	return o
}

func (o ResultSet) Valid() error {
	if o.Expiry == nil {
		return errors.New("missing mandatory expiry")
	}

	// TODO(tho)
	return nil
}
