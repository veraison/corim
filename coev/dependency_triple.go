// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"fmt"

	"github.com/veraison/corim/comid"
)

// DependencyTriple represents an ev-dependency-triple-record as defined in the
// TCG Concise Evidence CDDL. It contains a domain and a list of dependent domains.
// For now, we use Environment as the domain type, but this may be extended in the future.
type DependencyTriple struct {
	Domain           comid.Environment   `cbor:"0,keyasint" json:"domain"`
	DependentDomains []comid.Environment `cbor:"1,keyasint" json:"dependent-domains"`
}

// NewDependencyTriple creates a new DependencyTriple
func NewDependencyTriple() *DependencyTriple {
	return &DependencyTriple{}
}

// SetDomain sets the domain for this dependency triple
func (o *DependencyTriple) SetDomain(domain comid.Environment) *DependencyTriple {
	if o != nil {
		o.Domain = domain
	}
	return o
}

// AddDependentDomain adds a dependent domain to the triple
func (o *DependencyTriple) AddDependentDomain(domain comid.Environment) *DependencyTriple {
	if o != nil {
		o.DependentDomains = append(o.DependentDomains, domain)
	}
	return o
}

// Valid checks the validity of the DependencyTriple
func (o DependencyTriple) Valid() error {
	if err := o.Domain.Valid(); err != nil {
		return fmt.Errorf("invalid domain: %w", err)
	}

	if len(o.DependentDomains) == 0 {
		return fmt.Errorf("no dependent domains specified")
	}

	for i, domain := range o.DependentDomains {
		if err := domain.Valid(); err != nil {
			return fmt.Errorf("invalid dependent domain at index %d: %w", i, err)
		}
	}

	return nil
}

// DependencyTriples is a collection of DependencyTriple
type DependencyTriples []DependencyTriple

// NewDependencyTriples creates a new DependencyTriples collection
func NewDependencyTriples() *DependencyTriples {
	return &DependencyTriples{}
}

// Add appends a DependencyTriple to the collection
func (o *DependencyTriples) Add(dt *DependencyTriple) *DependencyTriples {
	if o != nil && dt != nil {
		*o = append(*o, *dt)
	}
	return o
}

// Valid checks the validity of all DependencyTriples in the collection
func (o DependencyTriples) Valid() error {
	for i, dt := range o {
		if err := dt.Valid(); err != nil {
			return fmt.Errorf("invalid dependency triple at index %d: %w", i, err)
		}
	}
	return nil
}