// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/eat"
)

func TestPSAProfile_Registration(t *testing.T) {
	// Test that the PSA profile is properly registered
	expectedURI := "tag:arm.com,2025:psa#1.0.0"
	
	_, err := eat.NewProfile(expectedURI)
	require.NoError(t, err)
	
	// Check that ProfileID is initialized correctly
	assert.NotNil(t, ProfileID)
	actualURI, err := ProfileID.Get()
	require.NoError(t, err)
	assert.Equal(t, expectedURI, actualURI)
}

func TestPSASwComponentMeasurementValues_Valid(t *testing.T) {
	tests := []struct {
		name        string
		values      PSASwComponentMeasurementValues
		expectError string
	}{
		{
			name: "valid minimal component",
			values: PSASwComponentMeasurementValues{
				Digests: []PSADigest{
					{
						Algorithm: "sha-256",
						Value:     make([]byte, 32), // 32 bytes for SHA-256
					},
				},
				CryptoKeys: [][]byte{make([]byte, 32)}, // 32-byte signer ID
			},
			expectError: "",
		},
		{
			name: "valid component with version and name",
			values: PSASwComponentMeasurementValues{
				Version: &PSASwComponentVersion{Version: "1.3.5"},
				Digests: []PSADigest{
					{
						Algorithm: "sha-256",
						Value:     make([]byte, 32),
					},
				},
				Name:       &[]string{"PRoT"}[0],
				CryptoKeys: [][]byte{make([]byte, 32)},
			},
			expectError: "",
		},
		{
			name: "multiple digests with different algorithms",
			values: PSASwComponentMeasurementValues{
				Digests: []PSADigest{
					{
						Algorithm: "sha-256",
						Value:     make([]byte, 32),
					},
					{
						Algorithm: "sha-384",
						Value:     make([]byte, 48),
					},
				},
				CryptoKeys: [][]byte{make([]byte, 32)},
			},
			expectError: "",
		},
		{
			name: "missing digests",
			values: PSASwComponentMeasurementValues{
				Digests:    []PSADigest{},
				CryptoKeys: [][]byte{make([]byte, 32)},
			},
			expectError: "digests field is mandatory and must contain at least one entry",
		},
		{
			name: "duplicate digest algorithms",
			values: PSASwComponentMeasurementValues{
				Digests: []PSADigest{
					{
						Algorithm: "sha-256",
						Value:     make([]byte, 32),
					},
					{
						Algorithm: "sha-256", // Duplicate
						Value:     make([]byte, 32),
					},
				},
				CryptoKeys: [][]byte{make([]byte, 32)},
			},
			expectError: "duplicate digest algorithm: sha-256",
		},
		{
			name: "missing crypto keys",
			values: PSASwComponentMeasurementValues{
				Digests: []PSADigest{
					{
						Algorithm: "sha-256",
						Value:     make([]byte, 32),
					},
				},
				CryptoKeys: [][]byte{},
			},
			expectError: "cryptokeys field is mandatory and must contain exactly one entry",
		},
		{
			name: "too many crypto keys",
			values: PSASwComponentMeasurementValues{
				Digests: []PSADigest{
					{
						Algorithm: "sha-256",
						Value:     make([]byte, 32),
					},
				},
				CryptoKeys: [][]byte{make([]byte, 32), make([]byte, 32)},
			},
			expectError: "cryptokeys field is mandatory and must contain exactly one entry",
		},
		{
			name: "invalid signer ID length",
			values: PSASwComponentMeasurementValues{
				Digests: []PSADigest{
					{
						Algorithm: "sha-256",
						Value:     make([]byte, 32),
					},
				},
				CryptoKeys: [][]byte{make([]byte, 31)}, // Invalid length
			},
			expectError: "signer-id must be 32, 48, or 64 bytes, got 31",
		},
		{
			name: "invalid digest length for sha-256",
			values: PSASwComponentMeasurementValues{
				Digests: []PSADigest{
					{
						Algorithm: "sha-256",
						Value:     make([]byte, 31), // Should be 32
					},
				},
				CryptoKeys: [][]byte{make([]byte, 32)},
			},
			expectError: "invalid hash length for sha-256: expected 32 bytes, got 31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.values.Valid()
			if tt.expectError != "" {
				assert.ErrorContains(t, err, tt.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPSACertNumType_Valid(t *testing.T) {
	tests := []struct {
		name        string
		certNum     PSACertNumType
		expectError string
	}{
		{
			name:        "valid certificate number",
			certNum:     PSACertNumType("1234567890123 - 12345"),
			expectError: "",
		},
		{
			name:        "invalid format - missing dash",
			certNum:     PSACertNumType("1234567890123 12345"),
			expectError: "invalid PSA certificate number format",
		},
		{
			name:        "invalid format - wrong first part length",
			certNum:     PSACertNumType("123456789012 - 12345"),
			expectError: "invalid PSA certificate number format",
		},
		{
			name:        "invalid format - wrong second part length",
			certNum:     PSACertNumType("1234567890123 - 1234"),
			expectError: "invalid PSA certificate number format",
		},
		{
			name:        "invalid format - contains letters",
			certNum:     PSACertNumType("123456789012a - 12345"),
			expectError: "invalid PSA certificate number format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.certNum.Valid()
			if tt.expectError != "" {
				assert.ErrorContains(t, err, tt.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPSASwRel_Valid(t *testing.T) {
	tests := []struct {
		name        string
		rel         PSASwRel
		expectError string
	}{
		{
			name: "valid updates relationship",
			rel: PSASwRel{
				Type:             PSAUpdates,
				SecurityCritical: true,
			},
			expectError: "",
		},
		{
			name: "valid patches relationship",
			rel: PSASwRel{
				Type:             PSAPatches,
				SecurityCritical: false,
			},
			expectError: "",
		},
		{
			name: "invalid relationship type",
			rel: PSASwRel{
				Type:             PSASwRelType(99),
				SecurityCritical: false,
			},
			expectError: "invalid PSA software relationship type: 99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rel.Valid()
			if tt.expectError != "" {
				assert.ErrorContains(t, err, tt.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPSASwRelationship_Valid(t *testing.T) {
	// Create valid measurements for testing
	newMeasurement := createTestMeasurement(t, "new-component", "1.3.0")
	oldMeasurement := createTestMeasurement(t, "old-component", "1.2.5")

	tests := []struct {
		name        string
		rel         PSASwRelationship
		expectError string
	}{
		{
			name: "valid relationship",
			rel: PSASwRelationship{
				New: newMeasurement,
				Relation: &PSASwRel{
					Type:             PSAUpdates,
					SecurityCritical: true,
				},
				Old: oldMeasurement,
			},
			expectError: "",
		},
		{
			name: "missing new measurement",
			rel: PSASwRelationship{
				New: nil,
				Relation: &PSASwRel{
					Type:             PSAUpdates,
					SecurityCritical: true,
				},
				Old: oldMeasurement,
			},
			expectError: "new measurement is required",
		},
		{
			name: "missing relationship",
			rel: PSASwRelationship{
				New:      newMeasurement,
				Relation: nil,
				Old:      oldMeasurement,
			},
			expectError: "relationship definition is required",
		},
		{
			name: "missing old measurement",
			rel: PSASwRelationship{
				New: newMeasurement,
				Relation: &PSASwRel{
					Type:             PSAUpdates,
					SecurityCritical: true,
				},
				Old: nil,
			},
			expectError: "old measurement is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rel.Valid()
			if tt.expectError != "" {
				assert.ErrorContains(t, err, tt.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewPSAUpdateRelationship(t *testing.T) {
	newMeasurement := createTestMeasurement(t, "component", "1.3.0")
	oldMeasurement := createTestMeasurement(t, "component", "1.2.5")

	rel, err := NewPSAUpdateRelationship(newMeasurement, oldMeasurement, true)
	require.NoError(t, err)
	assert.NotNil(t, rel)
	assert.Equal(t, PSAUpdates, rel.Relation.Type)
	assert.True(t, rel.Relation.SecurityCritical)
	assert.Equal(t, newMeasurement, rel.New)
	assert.Equal(t, oldMeasurement, rel.Old)
}

func TestNewPSAPatchRelationship(t *testing.T) {
	newMeasurement := createTestMeasurement(t, "component", "1.2.6")
	oldMeasurement := createTestMeasurement(t, "component", "1.2.5")

	rel, err := NewPSAPatchRelationship(newMeasurement, oldMeasurement, false)
	require.NoError(t, err)
	assert.NotNil(t, rel)
	assert.Equal(t, PSAPatches, rel.Relation.Type)
	assert.False(t, rel.Relation.SecurityCritical)
	assert.Equal(t, newMeasurement, rel.New)
	assert.Equal(t, oldMeasurement, rel.Old)
}

// Helper functions

func createTestMeasurement(t *testing.T, name, version string) *comid.Measurement {
	// Use a uint key as a placeholder since "text" type is not supported
	// In a full implementation, we would register "psa.software-component" 
	// as a new measurement key type
	mkey, err := comid.NewMkey(uint64(123), "uint")
	require.NoError(t, err)

	measurement := &comid.Measurement{
		Key: mkey,
		Val: comid.Mval{
			// In real usage, this would contain the PSA software component values
		},
	}

	return measurement
}
