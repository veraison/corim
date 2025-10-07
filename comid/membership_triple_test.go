// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/extensions"
)

func TestMembershipTriple_Valid_Success(t *testing.T) {
	memberVal := MemberVal{}
	memberVal.SetGroupID("group-1").SetRole("admin").SetStatus("active")

	membership := MustNewUUIDMembership(TestUUID)
	membership.SetValue(memberVal)

	triple := &MembershipTriple{
		Environment: Environment{
			Class: NewClassUUID(TestUUID).
				SetVendor("Test Vendor").
				SetModel("Test Model"),
		},
		Memberships: *NewMemberships().Add(membership),
	}

	err := triple.Valid()
	assert.NoError(t, err)
}

func TestMembershipTriple_Valid_EmptyEnvironment(t *testing.T) {
	memberVal := MemberVal{}
	memberVal.SetGroupID("group-1").SetRole("admin")

	membership := MustNewUUIDMembership(TestUUID)
	membership.SetValue(memberVal)

	triple := &MembershipTriple{
		Environment: Environment{}, // Empty environment
		Memberships: *NewMemberships().Add(membership),
	}

	err := triple.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "environment validation failed")
}

func TestMembershipTriple_Valid_EmptyMemberships(t *testing.T) {
	triple := &MembershipTriple{
		Environment: Environment{
			Class: NewClassUUID(TestUUID).
				SetVendor("Test Vendor").
				SetModel("Test Model"),
		},
		Memberships: *NewMemberships(), // Empty memberships
	}

	err := triple.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no membership entries")
}

func TestMembershipTriple_Extensions(t *testing.T) {
	triple := &MembershipTriple{}

	// Test RegisterExtensions
	extMap := extensions.NewMap().Add(ExtMemberVal, &struct{}{})
	err := triple.RegisterExtensions(extMap)
	assert.NoError(t, err)

	// Test GetExtensions
	exts := triple.GetExtensions()
	assert.NotNil(t, exts)
}

func TestMembershipTriples_NewMembershipTriples(t *testing.T) {
	triples := NewMembershipTriples()
	assert.NotNil(t, triples)
	assert.True(t, triples.IsEmpty())
}

func TestMembershipTriples_Add_Success(t *testing.T) {
	triples := NewMembershipTriples()

	memberVal := MemberVal{}
	memberVal.SetGroupID("group-1").SetRole("admin")

	membership := MustNewUUIDMembership(TestUUID)
	membership.SetValue(memberVal)

	triple := &MembershipTriple{
		Environment: Environment{
			Class: NewClassUUID(TestUUID).
				SetVendor("Test Vendor").
				SetModel("Test Model"),
		},
		Memberships: *NewMemberships().Add(membership),
	}

	result := triples.Add(triple)
	assert.Equal(t, triples, result)
	assert.False(t, triples.IsEmpty())
}

func TestMembershipTriples_Valid_Success(t *testing.T) {
	triples := NewMembershipTriples()

	memberVal := MemberVal{}
	memberVal.SetGroupID("group-1").SetRole("admin")

	membership := MustNewUUIDMembership(TestUUID)
	membership.SetValue(memberVal)

	triple := &MembershipTriple{
		Environment: Environment{
			Class: NewClassUUID(TestUUID).
				SetVendor("Test Vendor").
				SetModel("Test Model"),
		},
		Memberships: *NewMemberships().Add(membership),
	}

	triples.Add(triple)

	err := triples.Valid()
	assert.NoError(t, err)
}

func TestMembershipTriples_Valid_Empty(t *testing.T) {
	triples := NewMembershipTriples()

	err := triples.Valid()
	assert.NoError(t, err) // Empty collection is valid
}

func TestMembershipTriples_Valid_InvalidTriple(t *testing.T) {
	triples := NewMembershipTriples()

	// Add invalid triple with empty environment
	triple := &MembershipTriple{
		Environment: Environment{}, // Empty environment
		Memberships: *NewMemberships(),
	}

	triples.Add(triple)

	err := triples.Valid()
	assert.Error(t, err)
}

func TestMembershipTriples_RegisterExtensions(t *testing.T) {
	triples := NewMembershipTriples()

	extMap := extensions.NewMap().Add(ExtMemberVal, &struct{}{})
	err := triples.RegisterExtensions(extMap)
	assert.NoError(t, err)

	exts := triples.GetExtensions()
	assert.NotNil(t, exts)
}

func TestMembershipTriples_CBOR_RoundTrip(t *testing.T) {
	memberVal := MemberVal{}
	memberVal.SetGroupID("group-1").SetRole("admin")

	membership := MustNewUUIDMembership(TestUUID)
	membership.SetValue(memberVal)

	original := NewMembershipTriples()
	original.Add(&MembershipTriple{
		Environment: Environment{
			Class: NewClassUUID(TestUUID).
				SetVendor("Test Vendor").
				SetModel("Test Model"),
		},
		Memberships: *NewMemberships().Add(membership),
	})

	data, err := original.MarshalCBOR()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded MembershipTriples
	err = decoded.UnmarshalCBOR(data)
	require.NoError(t, err)

	err = decoded.Valid()
	assert.NoError(t, err)
}

func TestMembershipTriples_JSON_RoundTrip(t *testing.T) {
	memberVal := MemberVal{}
	memberVal.SetGroupID("group-1").SetRole("admin")

	membership := MustNewUUIDMembership(TestUUID)
	membership.SetValue(memberVal)

	original := NewMembershipTriples()
	original.Add(&MembershipTriple{
		Environment: Environment{
			Class: NewClassUUID(TestUUID).
				SetVendor("Test Vendor").
				SetModel("Test Model"),
		},
		Memberships: *NewMemberships().Add(membership),
	})

	data, err := original.MarshalJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded MembershipTriples
	err = decoded.UnmarshalJSON(data)
	require.NoError(t, err)

	err = decoded.Valid()
	assert.NoError(t, err)
}
