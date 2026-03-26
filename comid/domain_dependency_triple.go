// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"errors"
	"fmt"
)

// DomainDependencyTriple is the CDDL domain-dependency-triple-record:
//
//	domain-dependency-triple-record = [
//	  domain-id: domain-type,
//	  trustees: [ + domain-type ]
//	]
//	domain-type = environment-map
//
// A DDT links a domain to a set of trustee domains. The domain-id's
// trustworthiness depends on the trustees having been appraised first
// (draft-ietf-rats-corim §5.1.11.2).
type DomainDependencyTriple struct {
	_        struct{}      `cbor:",toarray"`
	DomainID Environment   `json:"domain-id"`
	Trustees []Environment `json:"trustees"`
}

// Valid checks that the record has a domain-id and at least one trustee.
func (o DomainDependencyTriple) Valid() error {
	if err := o.DomainID.Valid(); err != nil {
		return fmt.Errorf("domain-id: %w", err)
	}
	if len(o.Trustees) == 0 {
		return errors.New("at least one trustee required")
	}
	for i, t := range o.Trustees {
		if err := t.Valid(); err != nil {
			return fmt.Errorf("trustees[%d]: %w", i, err)
		}
	}
	return nil
}

// DomainDependencyTriples is an array of domain-dependency-triple
// (triples-map key 4: dependency-triples).
type DomainDependencyTriples []DomainDependencyTriple

// Valid checks each record.
func (o DomainDependencyTriples) Valid() error {
	for i, r := range o {
		if err := r.Valid(); err != nil {
			return fmt.Errorf("dependency triple [%d]: %w", i, err)
		}
	}
	return nil
}

// IsEmpty returns true if there are no records.
func (o DomainDependencyTriples) IsEmpty() bool {
	return len(o) == 0
}
