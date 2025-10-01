// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"fmt"

	"github.com/veraison/corim/comid"
)

// MembershipTriple represents an ev-membership-triple-record as defined in the
// TCG Concise Evidence CDDL. It contains a domain and a list of member environments.
type MembershipTriple struct {
	Domain       comid.Environment   `cbor:"0,keyasint" json:"domain"`
	Environments []comid.Environment `cbor:"1,keyasint" json:"environments"`
}

// NewMembershipTriple creates a new MembershipTriple
func NewMembershipTriple() *MembershipTriple {
	return &MembershipTriple{}
}

// SetDomain sets the domain for this membership triple
func (o *MembershipTriple) SetDomain(domain comid.Environment) *MembershipTriple {
	if o != nil {
		o.Domain = domain
	}
	return o
}

// AddEnvironment adds an environment to the triple
func (o *MembershipTriple) AddEnvironment(env comid.Environment) *MembershipTriple {
	if o != nil {
		o.Environments = append(o.Environments, env)
	}
	return o
}

// Valid checks the validity of the MembershipTriple
func (o MembershipTriple) Valid() error {
	if err := o.Domain.Valid(); err != nil {
		return fmt.Errorf("invalid domain: %w", err)
	}

	if len(o.Environments) == 0 {
		return fmt.Errorf("no environments specified")
	}

	for i, env := range o.Environments {
		if err := env.Valid(); err != nil {
			return fmt.Errorf("invalid environment at index %d: %w", i, err)
		}
	}

	return nil
}

// MembershipTriples is a collection of MembershipTriple
type MembershipTriples []MembershipTriple

// NewMembershipTriples creates a new MembershipTriples collection
func NewMembershipTriples() *MembershipTriples {
	return &MembershipTriples{}
}

// Add appends a MembershipTriple to the collection
func (o *MembershipTriples) Add(mt *MembershipTriple) *MembershipTriples {
	if o != nil && mt != nil {
		*o = append(*o, *mt)
	}
	return o
}

// Valid checks the validity of all MembershipTriples in the collection
func (o MembershipTriples) Valid() error {
	for i, mt := range o {
		if err := mt.Valid(); err != nil {
			return fmt.Errorf("invalid membership triple at index %d: %w", i, err)
		}
	}
	return nil
}