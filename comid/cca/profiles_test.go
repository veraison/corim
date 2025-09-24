// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cca

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/eat"
)

func TestCCAProfiles_URIFormat(t *testing.T) {
	// Verify Token Profile ID
	assert.Equal(t,
		"tag:arm.com,2025:cca-token",
		TokenProfileID.String(),
		"TokenProfileID should use tag URI scheme",
	)

	// Verify Platform Endorsements Profile ID
	assert.Equal(t,
		"tag:arm.com,2025:cca-endorsements",
		EndorsementsProfileID.String(),
		"EndorsementsProfileID should use tag URI scheme",
	)

	// Verify Realm Endorsements Profile ID
	assert.Equal(t,
		"tag:arm.com,2025:cca-realm-endorsements",
		RealmEndorsementsProfileID.String(),
		"RealmEndorsementsProfileID should use tag URI scheme",
	)
}

func TestCCAProfiles_Validation(t *testing.T) {
	// Test valid tag URIs can be created
	tests := []struct {
		name string
		uri string
	}{
		{
			name: "Token Profile",
			uri: "tag:arm.com,2025:cca-token",
		},
		{
			name: "Platform Endorsements Profile",
			uri: "tag:arm.com,2025:cca-endorsements",
		},
		{
			name: "Realm Endorsements Profile",
			uri: "tag:arm.com,2025:cca-realm-endorsements",
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

func TestCCAProfiles_InvalidURIs(t *testing.T) {
	// Test invalid URIs are rejected
	tests := []struct {
		name string
		uri string
	}{
		{
			name: "HTTP URL instead of tag URI",
			uri: "http://arm.com/cca-token",
		},
		{
			name: "Missing date",
			uri: "tag:arm.com:cca-token",
		},
		{
			name: "Invalid date",
			uri: "tag:arm.com,abcd:cca-token",
		},
		{
			name: "Empty specific part",
			uri: "tag:arm.com,2025:",
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

func TestCCAProfiles_Equality(t *testing.T) {
	// Test profile equality
	token1, err := eat.NewProfile("tag:arm.com,2025:cca-token")
	require.NoError(t, err)
	token2, err := eat.NewProfile("tag:arm.com,2025:cca-token")
	require.NoError(t, err)
	endorsements, err := eat.NewProfile("tag:arm.com,2025:cca-endorsements")
	require.NoError(t, err)

	// Same profile URIs should be equal
	assert.Equal(t, token1, token2)

	// Different profile URIs should not be equal
	assert.NotEqual(t, token1, endorsements)
}