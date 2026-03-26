// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainDependencyTriple_Valid_ok(t *testing.T) {
	triple := DomainDependencyTriple{
		DomainID: Environment{Instance: MustNewUEIDInstance(TestUEID)},
		Trustees: []Environment{
			{Instance: MustNewUUIDInstance(TestUUID)},
		},
	}
	assert.NoError(t, triple.Valid())
}

func TestDomainDependencyTriple_Valid_missing_domain_id(t *testing.T) {
	triple := DomainDependencyTriple{
		Trustees: []Environment{
			{Instance: MustNewUEIDInstance(TestUEID)},
		},
	}
	err := triple.Valid()
	assert.EqualError(t, err, "domain-id: environment must not be empty")
}

func TestDomainDependencyTriple_Valid_no_trustees(t *testing.T) {
	triple := DomainDependencyTriple{
		DomainID: Environment{Instance: MustNewUEIDInstance(TestUEID)},
		Trustees: []Environment{},
	}
	err := triple.Valid()
	assert.EqualError(t, err, "at least one trustee required")
}

func TestDomainDependencyTriple_Valid_cyclic(t *testing.T) {
	env := Environment{Instance: MustNewUEIDInstance(TestUEID)}
	triple := DomainDependencyTriple{
		DomainID: env,
		Trustees: []Environment{env},
	}
	err := triple.Valid()
	assert.EqualError(t, err, "trustees[0]: domain-id must not appear in trustees (acyclic constraint)")
}

func TestDomainDependencyTriple_Valid_duplicate_trustee(t *testing.T) {
	env := Environment{Instance: MustNewUUIDInstance(TestUUID)}
	triple := DomainDependencyTriple{
		DomainID: Environment{Instance: MustNewUEIDInstance(TestUEID)},
		Trustees: []Environment{env, env},
	}
	err := triple.Valid()
	assert.EqualError(t, err, "trustees[1]: duplicate of trustees[0]")
}

func TestDomainDependencyTriple_Valid_invalid_trustee(t *testing.T) {
	triple := DomainDependencyTriple{
		DomainID: Environment{Instance: MustNewUEIDInstance(TestUEID)},
		Trustees: []Environment{{}},
	}
	err := triple.Valid()
	assert.EqualError(t, err, "trustees[0]: environment must not be empty")
}

func TestDomainDependencyTriples_Valid_transitive_cycle(t *testing.T) {
	// A → B → A across two records
	envA := Environment{Instance: MustNewUEIDInstance(TestUEID)}
	envB := Environment{Instance: MustNewUUIDInstance(TestUUID)}

	triples := DomainDependencyTriples{
		{DomainID: envA, Trustees: []Environment{envB}},
		{DomainID: envB, Trustees: []Environment{envA}},
	}
	err := triples.Valid()
	assert.EqualError(t, err, "dependency-triples contain a cycle (acyclic constraint §5.1.11.2)")
}

func TestDomainDependencyTriples_Valid_empty(t *testing.T) {
	triples := DomainDependencyTriples{}
	assert.NoError(t, triples.Valid())
}

func TestDomainDependencyTriples_IsEmpty(t *testing.T) {
	triples := DomainDependencyTriples{}
	assert.True(t, triples.IsEmpty())

	triples = append(triples, DomainDependencyTriple{
		DomainID: Environment{Instance: MustNewUEIDInstance(TestUEID)},
		Trustees: []Environment{{Instance: MustNewUUIDInstance(TestUUID)}},
	})
	assert.False(t, triples.IsEmpty())
}

func TestTriples_AddDomainDependency(t *testing.T) {
	triple := &DomainDependencyTriple{
		DomainID: Environment{Instance: MustNewUEIDInstance(TestUEID)},
		Trustees: []Environment{{Instance: MustNewUUIDInstance(TestUUID)}},
	}

	var triples Triples
	assert.Nil(t, triples.DomainDependencies)

	triples.AddDomainDependency(triple)
	assert.NotNil(t, triples.DomainDependencies)
	assert.False(t, triples.DomainDependencies.IsEmpty())
	assert.Equal(t, *triple, (*triples.DomainDependencies)[0])
}
