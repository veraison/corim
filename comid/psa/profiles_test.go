// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/eat"
)

func TestPSAProfiles_URIFormat(t *testing.T) {
	// Verify Token Profile ID
	assert.Equal(t, 
		"tag:trustedcomputinggroup.org,2025:psa-token",
		TokenProfileID.String(),
		"TokenProfileID should use tag URI scheme",
	)

	// Verify Endorsements Profile ID
	assert.Equal(t,
		"tag:trustedcomputinggroup.org,2025:psa-endorsements",
		EndorsementsProfileID.String(), 
		"EndorsementsProfileID should use tag URI scheme",
	)
}

func TestPSAProfiles_Validation(t *testing.T) {
	// Test valid tag URIs can be created
	tests := []struct {
		name string
		uri string
	}{
		{
			name: "Token Profile",
			uri: "tag:trustedcomputinggroup.org,2025:psa-token",
		},
		{
			name: "Endorsements Profile",
			uri: "tag:trustedcomputinggroup.org,2025:psa-endorsements",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := eat.NewProfile(tt.uri)
			require.NoError(t, err)
			require.NotNil(t, profile)
			assert.Equal(t, tt.uri, profile.String())
		})
	}
}

func TestPSAProfiles_InvalidURIs(t *testing.T) {
	// Test invalid URIs are rejected
	tests := []struct {
		name string
		uri string
	}{
		{
			name: "HTTP URL instead of tag URI",
			uri: "http://trustedcomputinggroup.org/psa-token",
		},
		{
			name: "Missing date",
			uri: "tag:trustedcomputinggroup.org:psa-token",
		},
		{
			name: "Invalid date",
			uri: "tag:trustedcomputinggroup.org,abcd:psa-token",
		},
		{
			name: "Empty specific part",
			uri: "tag:trustedcomputinggroup.org,2025:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := eat.NewProfile(tt.uri)
			if err == nil {
				t.Errorf("Expected error for invalid URI %q, got nil", tt.uri)
			}
			assert.Nil(t, profile)
		})
	}
}

func TestPSAProfiles_Equality(t *testing.T) {
	// Test profile equality
	token1, err := eat.NewProfile("tag:trustedcomputinggroup.org,2025:psa-token")
	require.NoError(t, err)
	token2, err := eat.NewProfile("tag:trustedcomputinggroup.org,2025:psa-token")
	require.NoError(t, err)
	endorsements, err := eat.NewProfile("tag:trustedcomputinggroup.org,2025:psa-endorsements")
	require.NoError(t, err)

	// Same profile URIs should be equal
	assert.Equal(t, token1, token2)

	// Different profile URIs should not be equal
	assert.NotEqual(t, token1, endorsements)
}