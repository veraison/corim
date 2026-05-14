// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"errors"
	"fmt"
	"reflect"
)

// DomainMembershipTriple is the CDDL domain-membership-triple-record:
//
//	domain-membership-triple-record = [
//	  domain-id: domain-type,
//	  members: [ + domain-type ]
//	]
//	domain-type = environment-map
//
// A DMT links a domain to its member domains. It allows an endorser to issue
// an authoritative statement about the composition of an Attester as a
// collection of environments.
// (draft-ietf-rats-corim §5.1.11.1).
type DomainMembershipTriple struct {
	_        struct{}      `cbor:",toarray"`
	DomainID Environment   `json:"domain-id"`
	Members  []Environment `json:"members"`
}

// AddMember adds provided Environment as a member to the
// DomainMembershipTriple.
func (o *DomainMembershipTriple) AddMember(env Environment) *DomainMembershipTriple {
	o.Members = append(o.Members, env)
	return o
}

// Valid returns an error if the DomainMembershipTriple does not have any
// members or contains invalid environments.
func (o DomainMembershipTriple) Valid() error {
	if err := o.DomainID.Valid(); err != nil {
		return fmt.Errorf("domain-id: %w", err)
	}

	if len(o.Members) == 0 {
		return errors.New("must have at least one member")
	}

	for i, m := range o.Members {
		if err := m.Valid(); err != nil {
			return fmt.Errorf("member[%d]: %w", i, err)
		}
	}

	return nil
}

// DomainMembershipTriples is a container of DomainMembershipTriple instances.
type DomainMembershipTriples []DomainMembershipTriple

// NewDomainMebershipTriples returns a new empty DomainMebershipTriples.
func NewDomainMebershipTriples() *DomainMembershipTriples {
	return &DomainMembershipTriples{}
}

// Add a triple to the DomainMembershipTriples.
func (o *DomainMembershipTriples) Add(triple DomainMembershipTriple) *DomainMembershipTriples {
	*o = append(*o, triple)
	return o
}

// Valid returns an error if DomainMembershipTriples is empty or contains
// invalid elements.
func (o DomainMembershipTriples) Valid() error {
	if len(o) == 0 {
		return errors.New("must not be empty")
	}

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
			return fmt.Errorf("membership triple[%d]: %w", i, err)
		}
		from := indexOf(r.DomainID)
		for _, t := range r.Members {
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
			return errors.New("membership-triples contain a cycle (acyclic constraint §TODO)")
		}
	}

	return nil
}

// IsEmpty returns true if the DomainMembershipTriples does not contain any
// triples.
func (o DomainMembershipTriples) IsEmpty() bool {
	return len(o) == 0
}
