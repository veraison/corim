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
			{Instance: MustNewUEIDInstance(TestUEID)},
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

func TestDomainDependencyTriple_Valid_invalid_trustee(t *testing.T) {
	triple := DomainDependencyTriple{
		DomainID: Environment{Instance: MustNewUEIDInstance(TestUEID)},
		Trustees: []Environment{{}},
	}
	err := triple.Valid()
	assert.EqualError(t, err, "trustees[0]: environment must not be empty")
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
		Trustees: []Environment{{Instance: MustNewUEIDInstance(TestUEID)}},
	})
	assert.False(t, triples.IsEmpty())
}

func TestTriples_AddDomainDependency(t *testing.T) {
	triple := &DomainDependencyTriple{
		DomainID: Environment{Instance: MustNewUEIDInstance(TestUEID)},
		Trustees: []Environment{{Instance: MustNewUEIDInstance(TestUEID)}},
	}

	var triples Triples
	assert.Nil(t, triples.DomainDependencies)

	triples.AddDomainDependency(triple)
	assert.NotNil(t, triples.DomainDependencies)
	assert.False(t, triples.DomainDependencies.IsEmpty())
	assert.Equal(t, *triple, (*triples.DomainDependencies)[0])
}
