// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
)

func TestDependencyTriple_Valid(t *testing.T) {
	testCases := []struct {
		name    string
		dt      *DependencyTriple
		wantErr string
	}{
		{
			name: "valid dependency triple",
			dt: func() *DependencyTriple {
				dt := NewDependencyTriple()

				// Create a valid domain environment
				domainClass := comid.NewClassUUID(TestUUID)
				domainEnv := comid.Environment{Class: domainClass}
				dt.SetDomain(domainEnv)

				// Add a dependent domain
				depClass := comid.NewClassUUID(TestUUID2)
				depEnv := comid.Environment{Class: depClass}
				dt.AddDependentDomain(depEnv)

				return dt
			}(),
			wantErr: "",
		},
		{
			name: "invalid domain",
			dt: func() *DependencyTriple {
				dt := NewDependencyTriple()
				// Empty environment is invalid
				dt.SetDomain(comid.Environment{})
				depClass := comid.NewClassUUID(TestUUID2)
				depEnv := comid.Environment{Class: depClass}
				dt.AddDependentDomain(depEnv)
				return dt
			}(),
			wantErr: "invalid domain",
		},
		{
			name: "no dependent domains",
			dt: func() *DependencyTriple {
				dt := NewDependencyTriple()
				domainClass := comid.NewClassUUID(TestUUID)
				domainEnv := comid.Environment{Class: domainClass}
				dt.SetDomain(domainEnv)
				// No dependent domains added
				return dt
			}(),
			wantErr: "no dependent domains specified",
		},
		{
			name: "invalid dependent domain",
			dt: func() *DependencyTriple {
				dt := NewDependencyTriple()
				domainClass := comid.NewClassUUID(TestUUID)
				domainEnv := comid.Environment{Class: domainClass}
				dt.SetDomain(domainEnv)
				// Add invalid dependent domain (empty environment)
				dt.AddDependentDomain(comid.Environment{})
				return dt
			}(),
			wantErr: "invalid dependent domain at index 0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.dt.Valid()
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			}
		})
	}
}

func TestDependencyTriple_Setters(t *testing.T) {
	dt := NewDependencyTriple()
	
	// Test SetDomain
	domainClass := comid.NewClassUUID(TestUUID)
	domainEnv := comid.Environment{Class: domainClass}
	result := dt.SetDomain(domainEnv)
	
	assert.Equal(t, dt, result) // Should return self for chaining
	assert.Equal(t, domainEnv, dt.Domain)
	
	// Test AddDependentDomain
	depClass := comid.NewClassUUID(TestUUID2)
	depEnv := comid.Environment{Class: depClass}
	result = dt.AddDependentDomain(depEnv)
	
	assert.Equal(t, dt, result) // Should return self for chaining
	require.Len(t, dt.DependentDomains, 1)
	assert.Equal(t, depEnv, dt.DependentDomains[0])
	
	// Add another dependent domain
	dep2Class := comid.NewClassUUID(TestUUID3)
	dep2Env := comid.Environment{Class: dep2Class}
	dt.AddDependentDomain(dep2Env)
	
	require.Len(t, dt.DependentDomains, 2)
	assert.Equal(t, depEnv, dt.DependentDomains[0])
	assert.Equal(t, dep2Env, dt.DependentDomains[1])
}

func TestDependencyTriple_NilReceiver(t *testing.T) {
	var dt *DependencyTriple
	
	// Test that methods handle nil receiver gracefully
	domainClass := comid.NewClassUUID(TestUUID)
	domainEnv := comid.Environment{Class: domainClass}
	
	result := dt.SetDomain(domainEnv)
	assert.Nil(t, result)
	
	result = dt.AddDependentDomain(domainEnv)
	assert.Nil(t, result)
}

func TestDependencyTriples_Collection(t *testing.T) {
	dts := NewDependencyTriples()
	assert.NotNil(t, dts)
	assert.Len(t, *dts, 0)
	
	// Create a valid dependency triple
	dt := NewDependencyTriple()
	domainClass := comid.NewClassUUID(TestUUID)
	domainEnv := comid.Environment{Class: domainClass}
	dt.SetDomain(domainEnv)
	depClass := comid.NewClassUUID(TestUUID2)
	depEnv := comid.Environment{Class: depClass}
	dt.AddDependentDomain(depEnv)
	
	// Add to collection
	result := dts.Add(dt)
	assert.Equal(t, dts, result) // Should return self for chaining
	assert.Len(t, *dts, 1)
	
	// Test Valid() on collection
	err := dts.Valid()
	assert.NoError(t, err)
	
	// Add an invalid triple
	invalidDt := NewDependencyTriple()
	invalidDt.SetDomain(comid.Environment{}) // Invalid empty environment
	dts.Add(invalidDt)
	
	err = dts.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid dependency triple at index 1")
}

func TestDependencyTriples_NilReceiver(t *testing.T) {
	var dts *DependencyTriples
	
	dt := NewDependencyTriple()
	result := dts.Add(dt)
	assert.Nil(t, result)
}