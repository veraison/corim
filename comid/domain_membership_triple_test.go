// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DomainMembershipTriple_Valid(t *testing.T) {
	triple := DomainMembershipTriple{}
	err := triple.Valid()
	assert.ErrorContains(t, err, "environment must not be empty")

	triple.DomainID = Environment{
		Instance: MustNewUUIDInstance(TestUUID),
	}

	err = triple.Valid()
	assert.ErrorContains(t, err, "must have at least one member")

	triple.AddMember(Environment{})

	err = triple.Valid()
	assert.ErrorContains(t, err, "member[0]: environment must not be empty")

	triple.Members[0].Instance = MustNewUUIDInstance(TestUUID)

	err = triple.Valid()
	assert.NoError(t, err)
}

func Test_DomainMembershipTriples_Valid(t *testing.T) {
	triples := NewDomainMebershipTriples()
	assert.True(t, triples.IsEmpty())
	err := triples.Valid()
	assert.ErrorContains(t, err, "must not be empty")

	triples.Add(DomainMembershipTriple{
		DomainID: Environment{
			Instance: MustNewUUIDInstance(TestUUID),
		},
	})

	err = triples.Valid()
	assert.ErrorContains(t, err, "triple[0]: must have at least one member")

	(*triples)[0].AddMember(Environment{
		Instance: MustNewUUIDInstance(TestUUID),
	})

	err = triples.Valid()
	assert.NoError(t, err)
}
