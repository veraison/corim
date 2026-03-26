// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"errors"
	"fmt"
	"reflect"
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

// Valid checks that the record has a domain-id, at least one valid trustee,
// and that the domain-id does not directly appear in the trustees list
// (§5.1.11.2). Transitive cycle detection across multiple records is enforced
// by DomainDependencyTriples.Valid().
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
		if reflect.DeepEqual(o.DomainID, t) {
			return fmt.Errorf("trustees[%d]: domain-id must not appear in trustees (acyclic constraint)", i)
		}
		for j, prev := range o.Trustees[:i] {
			if reflect.DeepEqual(t, prev) {
				return fmt.Errorf("trustees[%d]: duplicate of trustees[%d]", i, j)
			}
		}
	}
	return nil
}

// DomainDependencyTriples is an array of domain-dependency-triple
// (triples-map key 4: dependency-triples).
type DomainDependencyTriples []DomainDependencyTriple

// Valid checks each record and enforces the transitive acyclic constraint
// across all records (§5.1.11.2: trust dependency graphs MUST be acyclic).
func (o DomainDependencyTriples) Valid() error {
	// Assign a stable integer index to each unique Environment.
	var envs []Environment
	indexOf := func(e Environment) int {
		for i, known := range envs {
			if reflect.DeepEqual(e, known) {
				return i
			}
		}
		envs = append(envs, e)
		return len(envs) - 1
	}

	adj := make(map[int][]int)
	for i, r := range o {
		if err := r.Valid(); err != nil {
			return fmt.Errorf("dependency triple [%d]: %w", i, err)
		}
		from := indexOf(r.DomainID)
		for _, t := range r.Trustees {
			adj[from] = append(adj[from], indexOf(t))
		}
	}

	const (
		unvisited = 0
		inStack   = 1
		done      = 2
	)
	state := make([]int, len(envs))

	var dfs func(node int) bool
	dfs = func(node int) bool {
		state[node] = inStack
		for _, neighbor := range adj[node] {
			if state[neighbor] == inStack || (state[neighbor] == unvisited && dfs(neighbor)) {
				return true
			}
		}
		state[node] = done
		return false
	}

	for i := range envs {
		if state[i] == unvisited && dfs(i) {
			return errors.New("dependency-triples contain a cycle (acyclic constraint §5.1.11.2)")
		}
	}
	return nil
}

// IsEmpty returns true if there are no records.
func (o DomainDependencyTriples) IsEmpty() bool {
	return len(o) == 0
}
