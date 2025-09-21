// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
)

func TestEvTriples_AddDependencyTriple(t *testing.T) {
	evTriples := NewEvTriples()
	
	// Create a valid dependency triple
	dt := NewDependencyTriple()
	domainClass := comid.NewClassUUID(TestUUID)
	domainEnv := comid.Environment{Class: domainClass}
	dt.SetDomain(domainEnv)
	depClass := comid.NewClassUUID(TestUUID2)
	depEnv := comid.Environment{Class: depClass}
	dt.AddDependentDomain(depEnv)
	
	// Add to EvTriples
	result := evTriples.AddDependencyTriple(dt)
	assert.Equal(t, evTriples, result) // Should return self for chaining
	
	require.NotNil(t, evTriples.DependencyTriples)
	require.Len(t, *evTriples.DependencyTriples, 1)
	assert.Equal(t, *dt, (*evTriples.DependencyTriples)[0])
	
	// Test that Valid() passes
	err := evTriples.Valid()
	assert.NoError(t, err)
}

func TestEvTriples_AddMembershipTriple(t *testing.T) {
	evTriples := NewEvTriples()
	
	// Create a valid membership triple
	mt := NewMembershipTriple()
	domainClass := comid.NewClassUUID(TestUUID)
	domainEnv := comid.Environment{Class: domainClass}
	mt.SetDomain(domainEnv)
	memberClass := comid.NewClassUUID(TestUUID2)
	memberEnv := comid.Environment{Class: memberClass}
	mt.AddEnvironment(memberEnv)
	
	// Add to EvTriples
	result := evTriples.AddMembershipTriple(mt)
	assert.Equal(t, evTriples, result) // Should return self for chaining
	
	require.NotNil(t, evTriples.MembershipTriples)
	require.Len(t, *evTriples.MembershipTriples, 1)
	assert.Equal(t, *mt, (*evTriples.MembershipTriples)[0])
	
	// Test that Valid() passes
	err := evTriples.Valid()
	assert.NoError(t, err)
}

func TestEvTriples_Valid_WithNewTriples(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func() *EvTriples
		wantErr string
	}{
		{
			name: "valid with dependency triples only",
			setup: func() *EvTriples {
				evTriples := NewEvTriples()
				dt := NewDependencyTriple()
				domainClass := comid.NewClassUUID(TestUUID)
				domainEnv := comid.Environment{Class: domainClass}
				dt.SetDomain(domainEnv)
				depClass := comid.NewClassUUID(TestUUID2)
				depEnv := comid.Environment{Class: depClass}
				dt.AddDependentDomain(depEnv)
				evTriples.AddDependencyTriple(dt)
				return evTriples
			},
			wantErr: "",
		},
		{
			name: "valid with membership triples only",
			setup: func() *EvTriples {
				evTriples := NewEvTriples()
				mt := NewMembershipTriple()
				domainClass := comid.NewClassUUID(TestUUID)
				domainEnv := comid.Environment{Class: domainClass}
				mt.SetDomain(domainEnv)
				memberClass := comid.NewClassUUID(TestUUID2)
				memberEnv := comid.Environment{Class: memberClass}
				mt.AddEnvironment(memberEnv)
				evTriples.AddMembershipTriple(mt)
				return evTriples
			},
			wantErr: "",
		},
		{
			name: "valid with both dependency and membership triples",
			setup: func() *EvTriples {
				evTriples := NewEvTriples()
				
				// Add dependency triple
				dt := NewDependencyTriple()
				domainClass := comid.NewClassUUID(TestUUID)
				domainEnv := comid.Environment{Class: domainClass}
				dt.SetDomain(domainEnv)
				depClass := comid.NewClassUUID(TestUUID2)
				depEnv := comid.Environment{Class: depClass}
				dt.AddDependentDomain(depEnv)
				evTriples.AddDependencyTriple(dt)
				
				// Add membership triple
				mt := NewMembershipTriple()
				mt.SetDomain(domainEnv)
				memberClass := comid.NewClassUUID(TestUUID3)
				memberEnv := comid.Environment{Class: memberClass}
				mt.AddEnvironment(memberEnv)
				evTriples.AddMembershipTriple(mt)
				
				return evTriples
			},
			wantErr: "",
		},
		{
			name: "invalid dependency triple",
			setup: func() *EvTriples {
				evTriples := NewEvTriples()
				dt := NewDependencyTriple()
				// Don't set domain or dependent domains - invalid
				evTriples.AddDependencyTriple(dt)
				return evTriples
			},
			wantErr: "invalid DependencyTriples",
		},
		{
			name: "invalid membership triple",
			setup: func() *EvTriples {
				evTriples := NewEvTriples()
				mt := NewMembershipTriple()
				// Don't set domain or environments - invalid
				evTriples.AddMembershipTriple(mt)
				return evTriples
			},
			wantErr: "invalid MembershipTriples",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			evTriples := tc.setup()
			err := evTriples.Valid()
			
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			}
		})
	}
}

func TestEvTriples_NilReceiver(t *testing.T) {
	var evTriples *EvTriples
	
	// Test dependency triple methods
	dt := NewDependencyTriple()
	result := evTriples.AddDependencyTriple(dt)
	assert.Nil(t, result)
	
	// Test membership triple methods
	mt := NewMembershipTriple()
	result = evTriples.AddMembershipTriple(mt)
	assert.Nil(t, result)
}

func TestEvTriples_MarshalCBOR_EmptyCollections(t *testing.T) {
	evTriples := NewEvTriples()
	
	// Initialize empty collections
	evTriples.DependencyTriples = NewDependencyTriples()
	evTriples.MembershipTriples = NewMembershipTriples()
	
	// Add a valid evidence triple to make EvTriples valid
	env := comid.Environment{
		Class: comid.NewClassUUID(TestUUID),
	}
	measurements := comid.NewMeasurements().Add(
		comid.MustNewUUIDMeasurement(TestUUID).
			SetRawValueBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{0xff, 0xff, 0xff, 0xff}),
	)
	triple := comid.ValueTriple{
		Environment:  env,
		Measurements: *measurements,
	}
	evTriples.AddEvidenceTriple(&triple)
	
	// Test CBOR marshaling - empty collections should be omitted
	data, err := evTriples.MarshalCBOR()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	
	// Unmarshal and verify
	var unmarshaled EvTriples
	err = unmarshaled.UnmarshalCBOR(data)
	assert.NoError(t, err)
	
	// Empty collections should be nil after marshaling/unmarshaling
	assert.Nil(t, unmarshaled.DependencyTriples)
	assert.Nil(t, unmarshaled.MembershipTriples)
}

func TestEvTriples_MarshalJSON_EmptyCollections(t *testing.T) {
	evTriples := NewEvTriples()
	
	// Initialize empty collections
	evTriples.DependencyTriples = NewDependencyTriples()
	evTriples.MembershipTriples = NewMembershipTriples()
	
	// Add a valid evidence triple to make EvTriples valid
	env := comid.Environment{
		Class: comid.NewClassUUID(TestUUID),
	}
	measurements := comid.NewMeasurements().Add(
		comid.MustNewUUIDMeasurement(TestUUID).
			SetRawValueBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{0xff, 0xff, 0xff, 0xff}),
	)
	triple := comid.ValueTriple{
		Environment:  env,
		Measurements: *measurements,
	}
	evTriples.AddEvidenceTriple(&triple)
	
	// Test JSON marshaling - empty collections should be omitted
	data, err := evTriples.MarshalJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	
	// Verify that empty collections are not in JSON
	jsonStr := string(data)
	assert.NotContains(t, jsonStr, "dependency-triples")
	assert.NotContains(t, jsonStr, "membership-triples")
}