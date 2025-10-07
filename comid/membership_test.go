// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
)

func TestMemberVal_SettersAndGetters(t *testing.T) {
	memberVal := &MemberVal{}

	// Test SetGroupID
	result := memberVal.SetGroupID("group-1")
	assert.Equal(t, memberVal, result)
	assert.Equal(t, "group-1", *memberVal.GroupID)

	// Test SetGroupName
	memberVal.SetGroupName("Test Group")
	assert.Equal(t, "Test Group", *memberVal.GroupName)

	// Test SetRole
	memberVal.SetRole("admin")
	assert.Equal(t, "admin", *memberVal.Role)

	// Test SetStatus
	memberVal.SetStatus("active")
	assert.Equal(t, "active", *memberVal.Status)

	// Test SetPermissions
	permissions := []string{"read", "write", "admin"}
	memberVal.SetPermissions(permissions)
	assert.Equal(t, permissions, *memberVal.Permissions)

	// Test SetOrganizationID
	memberVal.SetOrganizationID("org-123")
	assert.Equal(t, "org-123", *memberVal.OrganizationID)

	// Test SetUEID
	testUEID := eat.UEID(TestUEID)
	memberVal.SetUEID(testUEID)
	assert.Equal(t, testUEID, *memberVal.UEID)

	// Test SetUUID
	memberVal.SetUUID(TestUUID)
	assert.Equal(t, TestUUID, *memberVal.UUID)

	// Test SetName
	memberVal.SetName("Test Name")
	assert.Equal(t, "Test Name", *memberVal.Name)
}

func TestMemberVal_Valid_Success(t *testing.T) {
	memberVal := MemberVal{}
	memberVal.SetGroupID("group-1")

	err := memberVal.Valid()
	assert.NoError(t, err)
}

func TestMemberVal_Valid_EmptyValues(t *testing.T) {
	memberVal := MemberVal{}

	err := memberVal.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no membership value set")
}

func TestMemberVal_Valid_WithValidUEID(t *testing.T) {
	memberVal := MemberVal{}
	testUEID := eat.UEID(TestUEID)
	memberVal.SetUEID(testUEID)

	err := memberVal.Valid()
	assert.NoError(t, err)
}

func TestMemberVal_Valid_WithValidUUID(t *testing.T) {
	memberVal := MemberVal{}
	memberVal.SetUUID(TestUUID)

	err := memberVal.Valid()
	assert.NoError(t, err)
}

func TestMemberVal_RegisterExtensions(t *testing.T) {
	memberVal := &MemberVal{}

	extMap := extensions.NewMap().Add(ExtMemberVal, &struct{}{})
	err := memberVal.RegisterExtensions(extMap)
	assert.NoError(t, err)

	exts := memberVal.GetExtensions()
	assert.NotNil(t, exts)
}

func TestMemberVal_CBOR_RoundTrip(t *testing.T) {
	original := MemberVal{}
	original.SetGroupID("group-1").SetRole("admin").SetStatus("active")

	data, err := original.MarshalCBOR()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded MemberVal
	err = decoded.UnmarshalCBOR(data)
	require.NoError(t, err)

	assert.Equal(t, *original.GroupID, *decoded.GroupID)
	assert.Equal(t, *original.Role, *decoded.Role)
	assert.Equal(t, *original.Status, *decoded.Status)
}

func TestMemberVal_JSON_RoundTrip(t *testing.T) {
	original := MemberVal{}
	original.SetGroupID("group-1").SetRole("admin").SetStatus("active")

	data, err := original.MarshalJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded MemberVal
	err = decoded.UnmarshalJSON(data)
	require.NoError(t, err)

	assert.Equal(t, *original.GroupID, *decoded.GroupID)
	assert.Equal(t, *original.Role, *decoded.Role)
	assert.Equal(t, *original.Status, *decoded.Status)
}

func TestMembership_NewMembership_Success(t *testing.T) {
	membership, err := NewMembership(TestUUID, "uuid")
	require.NoError(t, err)
	assert.NotNil(t, membership)
	assert.NotNil(t, membership.Key)
}

func TestMembership_NewMembership_InvalidType(t *testing.T) {
	_, err := NewMembership(TestUUID, "invalid-type")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown Mkey type")
}

func TestMembership_MustNewMembership_Success(t *testing.T) {
	membership := MustNewMembership(TestUUID, "uuid")
	assert.NotNil(t, membership)
	assert.NotNil(t, membership.Key)
}

func TestMembership_MustNewMembership_Panic(t *testing.T) {
	assert.Panics(t, func() {
		MustNewMembership(TestUUID, "invalid-type")
	})
}

func TestMembership_MustNewUUIDMembership(t *testing.T) {
	membership := MustNewUUIDMembership(TestUUID)
	assert.NotNil(t, membership)
	assert.NotNil(t, membership.Key)
}

func TestMembership_SetValue(t *testing.T) {
	membership := MustNewUUIDMembership(TestUUID)

	memberVal := MemberVal{}
	memberVal.SetGroupID("group-1")

	result := membership.SetValue(memberVal)
	assert.Equal(t, membership, result)
	assert.Equal(t, memberVal, membership.Val)
}

func TestMembership_Valid_Success(t *testing.T) {
	membership := MustNewUUIDMembership(TestUUID)

	memberVal := MemberVal{}
	memberVal.SetGroupID("group-1")
	membership.SetValue(memberVal)

	err := membership.Valid()
	assert.NoError(t, err)
}

func TestMembership_Valid_InvalidValue(t *testing.T) {
	membership := MustNewUUIDMembership(TestUUID)

	// Empty MemberVal should be invalid
	memberVal := MemberVal{}
	membership.SetValue(memberVal)

	err := membership.Valid()
	assert.Error(t, err)
}

func TestMemberships_NewMemberships(t *testing.T) {
	memberships := NewMemberships()
	assert.NotNil(t, memberships)
	assert.True(t, memberships.IsEmpty())
}

func TestMemberships_Add_Success(t *testing.T) {
	memberships := NewMemberships()

	membership := MustNewUUIDMembership(TestUUID)
	memberVal := MemberVal{}
	memberVal.SetGroupID("group-1")
	membership.SetValue(memberVal)

	result := memberships.Add(membership)
	assert.Equal(t, memberships, result)
	assert.False(t, memberships.IsEmpty())
}

func TestMemberships_Valid_Success(t *testing.T) {
	memberships := NewMemberships()

	membership := MustNewUUIDMembership(TestUUID)
	memberVal := MemberVal{}
	memberVal.SetGroupID("group-1")
	membership.SetValue(memberVal)

	memberships.Add(membership)

	err := memberships.Valid()
	assert.NoError(t, err)
}

func TestMemberships_Valid_Empty(t *testing.T) {
	memberships := NewMemberships()

	err := memberships.Valid()
	assert.NoError(t, err) // Empty collection is valid
}

func TestMemberships_Valid_InvalidMembership(t *testing.T) {
	memberships := NewMemberships()

	membership := MustNewUUIDMembership(TestUUID)
	// Add membership with empty value (invalid)
	memberVal := MemberVal{}
	membership.SetValue(memberVal)

	memberships.Add(membership)

	err := memberships.Valid()
	assert.Error(t, err)
}

func TestMemberships_RegisterExtensions(t *testing.T) {
	memberships := NewMemberships()

	extMap := extensions.NewMap().Add(ExtMemberVal, &struct{}{})
	err := memberships.RegisterExtensions(extMap)
	assert.NoError(t, err)

	exts := memberships.GetExtensions()
	assert.NotNil(t, exts)
}

func TestMemberships_CBOR_RoundTrip(t *testing.T) {
	original := NewMemberships()

	membership := MustNewUUIDMembership(TestUUID)
	memberVal := MemberVal{}
	memberVal.SetGroupID("group-1").SetRole("admin")
	membership.SetValue(memberVal)

	original.Add(membership)

	data, err := original.MarshalCBOR()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded Memberships
	err = decoded.UnmarshalCBOR(data)
	require.NoError(t, err)

	err = decoded.Valid()
	assert.NoError(t, err)
}

func TestMemberships_JSON_RoundTrip(t *testing.T) {
	original := NewMemberships()

	membership := MustNewUUIDMembership(TestUUID)
	memberVal := MemberVal{}
	memberVal.SetGroupID("group-1").SetRole("admin")
	membership.SetValue(memberVal)

	original.Add(membership)

	data, err := original.MarshalJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded Memberships
	err = decoded.UnmarshalJSON(data)
	require.NoError(t, err)

	err = decoded.Valid()
	assert.NoError(t, err)
}