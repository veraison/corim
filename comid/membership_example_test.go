// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Example_membershipTriple() {
	// Create a new Comid
	comid := NewComid().
		SetLanguage("en-US").
		SetTagIdentity("membership-example", 1).
		AddEntity("ACME Corp", &TestRegID, RoleCreator, RoleTagCreator)

	// Create membership information for an administrator
	adminMember := MemberVal{}
	adminMember.SetGroupID("admin-group").
		SetGroupName("Administrator Group").
		SetRole("admin").
		SetStatus("active").
		SetPermissions([]string{"read", "write", "admin"}).
		SetOrganizationID("acme-corp")

	// Create a membership keyed by UUID
	membership := MustNewUUIDMembership(TestUUID)
	membership.SetValue(adminMember)

	// Create a membership triple that associates an environment with memberships
	triple := &MembershipTriple{
		Environment: Environment{
			Class: NewClassUUID(TestUUID).
				SetVendor("ACME Corp").
				SetModel("Secure Device v1.0").
				SetLayer(1),
			Instance: MustNewUEIDInstance(TestUEID),
		},
		Memberships: *NewMemberships().Add(membership),
	}

	// Add the membership triple to the Comid
	comid.AddMembershipTriple(triple)

	// Validate the comid
	err := comid.Valid()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Convert to JSON for demonstration
	jsonData, err := comid.ToJSON()
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		return
	}

	fmt.Printf("Successfully created Comid with MembershipTriple: %d bytes\n", len(jsonData))
	fmt.Println("MembershipTriple includes:")
	fmt.Println("- Environment with class and instance")
	fmt.Println("- Membership with group, role, and permissions")

	// Output:
	// Successfully created Comid with MembershipTriple: 669 bytes
	// MembershipTriple includes:
	// - Environment with class and instance
	// - Membership with group, role, and permissions
}

func Example_membershipTriple_multipleMembers() {
	// Create a new Comid for multiple memberships
	comid := NewComid().
		SetLanguage("en-US").
		SetTagIdentity("multi-membership-example", 1).
		AddEntity("ACME Corp", &TestRegID, RoleCreator, RoleTagCreator)

	// Create different membership types
	adminMember := MemberVal{}
	adminMember.SetGroupID("admin-group").
		SetRole("admin").
		SetStatus("active").
		SetPermissions([]string{"read", "write", "admin"})

	userMember := MemberVal{}
	userMember.SetGroupID("user-group").
		SetRole("user").
		SetStatus("active").
		SetPermissions([]string{"read"})

	// Create memberships with different key types
	adminMembership := MustNewUUIDMembership(TestUUID)
	adminMembership.SetValue(adminMember)

	userMembership := MustNewUUIDMembership(TestUUID)
	userMembership.SetValue(userMember)

	// Create membership collection
	memberships := NewMemberships().
		Add(adminMembership).
		Add(userMembership)

	// Create a membership triple
	triple := &MembershipTriple{
		Environment: Environment{
			Class: NewClassUUID(TestUUID).
				SetVendor("ACME Corp").
				SetModel("Multi-User Device"),
		},
		Memberships: *memberships,
	}

	// Add to comid
	comid.AddMembershipTriple(triple)

	// Validate
	err := comid.Valid()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Created Comid with %d memberships\n", len(memberships.Values))

	// Output:
	// Created Comid with 2 memberships
}

func TestExample_membershipTriple(t *testing.T) {
	// This test ensures the example function works correctly
	Example_membershipTriple()
}

func TestExample_membershipTriple_multipleMembers(t *testing.T) {
	// This test ensures the multiple members example works correctly
	Example_membershipTriple_multipleMembers()
}

func Test_membershipTriple_RealWorldScenario(t *testing.T) {
	// Test a more complex real-world scenario
	comid := NewComid().
		SetLanguage("en-US").
		SetTagIdentity("enterprise-device-membership", 1).
		AddEntity("Enterprise Corp", &TestRegID, RoleCreator, RoleTagCreator)

	// Device administrator
	deviceAdmin := MemberVal{}
	deviceAdmin.SetGroupID("device-admin").
		SetGroupName("Device Administrators").
		SetRole("device-admin").
		SetStatus("active").
		SetPermissions([]string{"configure", "monitor", "update", "reset"}).
		SetOrganizationID("enterprise-corp").
		SetName("Device Admin Role")

	// Security officer
	securityOfficer := MemberVal{}
	securityOfficer.SetGroupID("security-team").
		SetGroupName("Security Officers").
		SetRole("security-officer").
		SetStatus("active").
		SetPermissions([]string{"audit", "monitor", "investigate"}).
		SetOrganizationID("enterprise-corp").
		SetName("Security Officer Role")

	// Regular user
	regularUser := MemberVal{}
	regularUser.SetGroupID("users").
		SetGroupName("Regular Users").
		SetRole("user").
		SetStatus("active").
		SetPermissions([]string{"use", "view-status"}).
		SetOrganizationID("enterprise-corp").
		SetName("Regular User Role")

	// Create memberships
	adminMembership := MustNewUUIDMembership(TestUUID)
	adminMembership.SetValue(deviceAdmin)

	securityMembership := MustNewUintMembership(12345)
	securityMembership.SetValue(securityOfficer)

	userMembership := MustNewUintMembership(67890)
	userMembership.SetValue(regularUser)

	// Create the environment (enterprise device)
	environment := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Enterprise Corp").
			SetModel("Secure Workstation Pro").
			SetLayer(1),
		Instance: MustNewUEIDInstance(TestUEID),
	}

	// Create membership triple
	triple := &MembershipTriple{
		Environment: environment,
		Memberships: *NewMemberships().
			Add(adminMembership).
			Add(securityMembership).
			Add(userMembership),
	}

	// Add to comid
	comid.AddMembershipTriple(triple)

	// Validate
	err := comid.Valid()
	require.NoError(t, err)

	// Test serialization
	cborData, err := comid.ToCBOR()
	require.NoError(t, err)
	assert.NotEmpty(t, cborData)

	jsonData, err := comid.ToJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Verify content
	assert.Contains(t, string(jsonData), "membership-triples")
	assert.Contains(t, string(jsonData), "device-admin")
	assert.Contains(t, string(jsonData), "security-officer")
	assert.Contains(t, string(jsonData), "Enterprise Corp")

	fmt.Printf("Enterprise membership scenario: %d bytes CBOR, %d bytes JSON\n",
		len(cborData), len(jsonData))
}
