// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
)

func TestMembershipTriple_Valid(t *testing.T) {
	testCases := []struct {
		name    string
		mt      *MembershipTriple
		wantErr string
	}{
		{
			name: "valid membership triple",
			mt: func() *MembershipTriple {
				mt := NewMembershipTriple()

				// Create a valid domain environment
				domainClass := comid.NewClassUUID(TestUUID)
				domainEnv := comid.Environment{Class: domainClass}
				mt.SetDomain(domainEnv)

				// Add a member environment
				memberClass := comid.NewClassUUID(TestUUID2)
				memberEnv := comid.Environment{Class: memberClass}
				mt.AddEnvironment(memberEnv)

				return mt
			}(),
			wantErr: "",
		},
		{
			name: "invalid domain",
			mt: func() *MembershipTriple {
				mt := NewMembershipTriple()
				// Empty environment is invalid
				mt.SetDomain(comid.Environment{})
				memberClass := comid.NewClassUUID(TestUUID2)
				memberEnv := comid.Environment{Class: memberClass}
				mt.AddEnvironment(memberEnv)
				return mt
			}(),
			wantErr: "invalid domain",
		},
		{
			name: "no environments",
			mt: func() *MembershipTriple {
				mt := NewMembershipTriple()
				domainClass := comid.NewClassUUID(TestUUID)
				domainEnv := comid.Environment{Class: domainClass}
				mt.SetDomain(domainEnv)
				// No environments added
				return mt
			}(),
			wantErr: "no environments specified",
		},
		{
			name: "invalid environment",
			mt: func() *MembershipTriple {
				mt := NewMembershipTriple()
				domainClass := comid.NewClassUUID(TestUUID)
				domainEnv := comid.Environment{Class: domainClass}
				mt.SetDomain(domainEnv)
				// Add invalid environment (empty environment)
				mt.AddEnvironment(comid.Environment{})
				return mt
			}(),
			wantErr: "invalid environment at index 0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.mt.Valid()
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			}
		})
	}
}

func TestMembershipTriple_Setters(t *testing.T) {
	mt := NewMembershipTriple()
	
	// Test SetDomain
	domainClass := comid.NewClassUUID(TestUUID)
	domainEnv := comid.Environment{Class: domainClass}
	result := mt.SetDomain(domainEnv)
	
	assert.Equal(t, mt, result) // Should return self for chaining
	assert.Equal(t, domainEnv, mt.Domain)
	
	// Test AddEnvironment
	memberClass := comid.NewClassUUID(TestUUID2)
	memberEnv := comid.Environment{Class: memberClass}
	result = mt.AddEnvironment(memberEnv)
	
	assert.Equal(t, mt, result) // Should return self for chaining
	require.Len(t, mt.Environments, 1)
	assert.Equal(t, memberEnv, mt.Environments[0])
	
	// Add another environment
	member2Class := comid.NewClassUUID(TestUUID3)
	member2Env := comid.Environment{Class: member2Class}
	mt.AddEnvironment(member2Env)
	
	require.Len(t, mt.Environments, 2)
	assert.Equal(t, memberEnv, mt.Environments[0])
	assert.Equal(t, member2Env, mt.Environments[1])
}

func TestMembershipTriple_NilReceiver(t *testing.T) {
	var mt *MembershipTriple
	
	// Test that methods handle nil receiver gracefully
	domainClass := comid.NewClassUUID(TestUUID)
	domainEnv := comid.Environment{Class: domainClass}
	
	result := mt.SetDomain(domainEnv)
	assert.Nil(t, result)
	
	result = mt.AddEnvironment(domainEnv)
	assert.Nil(t, result)
}

func TestMembershipTriples_Collection(t *testing.T) {
	mts := NewMembershipTriples()
	assert.NotNil(t, mts)
	assert.Len(t, *mts, 0)
	
	// Create a valid membership triple
	mt := NewMembershipTriple()
	domainClass := comid.NewClassUUID(TestUUID)
	domainEnv := comid.Environment{Class: domainClass}
	mt.SetDomain(domainEnv)
	memberClass := comid.NewClassUUID(TestUUID2)
	memberEnv := comid.Environment{Class: memberClass}
	mt.AddEnvironment(memberEnv)
	
	// Add to collection
	result := mts.Add(mt)
	assert.Equal(t, mts, result) // Should return self for chaining
	assert.Len(t, *mts, 1)
	
	// Test Valid() on collection
	err := mts.Valid()
	assert.NoError(t, err)
	
	// Add an invalid triple
	invalidMt := NewMembershipTriple()
	invalidMt.SetDomain(comid.Environment{}) // Invalid empty environment
	mts.Add(invalidMt)
	
	err = mts.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid membership triple at index 1")
}

func TestMembershipTriples_NilReceiver(t *testing.T) {
	var mts *MembershipTriples
	
	mt := NewMembershipTriple()
	result := mts.Add(mt)
	assert.Nil(t, result)
}