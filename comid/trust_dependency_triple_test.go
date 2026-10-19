// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrustDependencyTriple_Valid_ok(t *testing.T) {
	triple := TrustDependencyTriple{
		DomainID: Environment{Instance: MustNewUEIDInstance(TestUEID)},
		Trustees: []Environment{
			{Instance: MustNewUUIDInstance(TestUUID)},
		},
	}
	assert.NoError(t, triple.Valid())
}

func TestTrustDependencyTriple_Valid_missing_domain_id(t *testing.T) {
	triple := TrustDependencyTriple{
		Trustees: []Environment{
			{Instance: MustNewUEIDInstance(TestUEID)},
		},
	}
	err := triple.Valid()
	assert.EqualError(t, err, "domain-id: environment must not be empty")
}

func TestTrustDependencyTriple_Valid_no_trustees(t *testing.T) {
	triple := TrustDependencyTriple{
		DomainID: Environment{Instance: MustNewUEIDInstance(TestUEID)},
		Trustees: []Environment{},
	}
	err := triple.Valid()
	assert.EqualError(t, err, "at least one trustee required")
}

func TestTrustDependencyTriple_Valid_cyclic(t *testing.T) {
	env := Environment{Instance: MustNewUEIDInstance(TestUEID)}
	triple := TrustDependencyTriple{
		DomainID: env,
		Trustees: []Environment{env},
	}
	err := triple.Valid()
	assert.EqualError(t, err, "trustees[0]: domain-id must not appear in trustees (acyclic constraint)")
}

func TestTrustDependencyTriple_Valid_duplicate_trustee(t *testing.T) {
	env := Environment{Instance: MustNewUUIDInstance(TestUUID)}
	triple := TrustDependencyTriple{
		DomainID: Environment{Instance: MustNewUEIDInstance(TestUEID)},
		Trustees: []Environment{env, env},
	}
	err := triple.Valid()
	assert.EqualError(t, err, "trustees[1]: duplicate of trustees[0]")
}

func TestTrustDependencyTriple_Valid_invalid_trustee(t *testing.T) {
	triple := TrustDependencyTriple{
		DomainID: Environment{Instance: MustNewUEIDInstance(TestUEID)},
		Trustees: []Environment{{}},
	}
	err := triple.Valid()
	assert.EqualError(t, err, "trustees[0]: environment must not be empty")
}

func TestTrustDependencyTriples_Valid_transitive_cycle(t *testing.T) {
	// A → B → A across two records
	envA := Environment{Instance: MustNewUEIDInstance(TestUEID)}
	envB := Environment{Instance: MustNewUUIDInstance(TestUUID)}

	triples := TrustDependencyTriples{
		{DomainID: envA, Trustees: []Environment{envB}},
		{DomainID: envB, Trustees: []Environment{envA}},
	}
	err := triples.Valid()
	assert.EqualError(t, err, "dependency-triples contain a cycle (acyclic constraint §5.1.11.2)")
}

func TestTrustDependencyTriples_Valid_empty(t *testing.T) {
	triples := TrustDependencyTriples{}
	assert.NoError(t, triples.Valid())
}

func TestTrustDependencyTriples_IsEmpty(t *testing.T) {
	triples := TrustDependencyTriples{}
	assert.True(t, triples.IsEmpty())

	triples = append(triples, TrustDependencyTriple{
		DomainID: Environment{Instance: MustNewUEIDInstance(TestUEID)},
		Trustees: []Environment{{Instance: MustNewUUIDInstance(TestUUID)}},
	})
	assert.False(t, triples.IsEmpty())
}

func TestTriples_AddTrustDependency(t *testing.T) {
	triple := &TrustDependencyTriple{
		DomainID: Environment{Instance: MustNewUEIDInstance(TestUEID)},
		Trustees: []Environment{{Instance: MustNewUUIDInstance(TestUUID)}},
	}

	var triples Triples
	assert.Nil(t, triples.TrustDependencies)

	triples.AddTrustDependency(triple)
	assert.NotNil(t, triples.TrustDependencies)
	assert.False(t, triples.TrustDependencies.IsEmpty())
	assert.Equal(t, *triple, (*triples.TrustDependencies)[0])
}
