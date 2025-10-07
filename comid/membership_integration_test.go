// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComid_AddMembershipTriple_Success(t *testing.T) {
	comid := NewComid()
	comid.SetTagIdentity("test-comid", 0)

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

	result := comid.AddMembershipTriple(triple)
	assert.Equal(t, comid, result)
	assert.NotNil(t, comid.Triples.MembershipTriples)
	assert.False(t, comid.Triples.MembershipTriples.IsEmpty())
}

func TestComid_AddMembershipTriple_Validation(t *testing.T) {
	comid := NewComid()
	comid.SetTagIdentity("test-comid", 0)

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

	comid.AddMembershipTriple(triple)

	err := comid.Valid()
	assert.NoError(t, err)
}

func TestTriples_AddMembershipTriple_Success(t *testing.T) {
	triples := &Triples{}

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

	result := triples.AddMembershipTriple(triple)
	assert.Equal(t, triples, result)
	assert.NotNil(t, triples.MembershipTriples)
	assert.False(t, triples.MembershipTriples.IsEmpty())
}

func TestTriples_Valid_WithMembershipTriples(t *testing.T) {
	triples := &Triples{}

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

	triples.AddMembershipTriple(triple)

	err := triples.Valid()
	assert.NoError(t, err)
}

func TestTriples_CBOR_RoundTrip_WithMembershipTriples(t *testing.T) {
	original := &Triples{}

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

	original.AddMembershipTriple(triple)

	data, err := original.MarshalCBOR()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded Triples
	err = decoded.UnmarshalCBOR(data)
	require.NoError(t, err)

	err = decoded.Valid()
	assert.NoError(t, err)

	// Verify that the membership triple was preserved
	assert.NotNil(t, decoded.MembershipTriples)
	assert.False(t, decoded.MembershipTriples.IsEmpty())
}

func TestComid_Full_Example_WithMembershipTriple(t *testing.T) {
	comid := NewComid().
		SetLanguage("en-US").
		SetTagIdentity("membership-test-comid", 1).
		AddEntity("Test Corp", &TestRegID, RoleCreator, RoleTagCreator)

	// Create membership information
	memberVal := MemberVal{}
	memberVal.SetGroupID("admin-group").
		SetGroupName("Administrator Group").
		SetRole("admin").
		SetStatus("active").
		SetPermissions([]string{"read", "write", "admin"}).
		SetOrganizationID("test-corp")

	membership := MustNewUUIDMembership(TestUUID)
	membership.SetValue(memberVal)

	triple := &MembershipTriple{
		Environment: Environment{
			Class: NewClassUUID(TestUUID).
				SetVendor("Test Vendor").
				SetModel("Test Model").
				SetLayer(1),
			Instance: MustNewUEIDInstance(TestUEID),
		},
		Memberships: *NewMemberships().Add(membership),
	}

	result := comid.AddMembershipTriple(triple)
	assert.Equal(t, comid, result)

	// Validate the full comid
	err := comid.Valid()
	assert.NoError(t, err)

	// Test CBOR serialization
	cborData, err := comid.ToCBOR()
	require.NoError(t, err)
	assert.NotEmpty(t, cborData)

	// Test JSON serialization
	jsonData, err := comid.ToJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Verify that membership triples are included in the JSON
	assert.Contains(t, string(jsonData), "membership-triples")
}
