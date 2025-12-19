// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
)

func TestPSASoftwareComponentKey_Valid(t *testing.T) {
	tests := []struct {
		name    string
		key     PSASoftwareComponentKeyType
		wantErr string
	}{
		{
			name: "valid key",
			key:  PSASoftwareComponentKeyType(PSASoftwareComponentType),
		},
		{
			name:    "invalid key",
			key:     PSASoftwareComponentKeyType("invalid"),
			wantErr: `invalid PSA software component key: expected "psa.software-component", got "invalid"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.key.Valid()
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPSASoftwareComponentKey_String(t *testing.T) {
	key := PSASoftwareComponentKeyType(PSASoftwareComponentType)
	assert.Equal(t, PSASoftwareComponentType, key.String())
}

func TestPSASoftwareComponentKey_Type(t *testing.T) {
	key := PSASoftwareComponentKeyType(PSASoftwareComponentType)
	assert.Equal(t, PSASoftwareComponentType, key.Type())
}

func TestPSASoftwareComponentKey_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    PSASoftwareComponentKeyType
		wantErr bool
	}{
		{
			name: "valid json",
			json: `"psa.software-component"`,
			want: PSASoftwareComponentKeyType(PSASoftwareComponentType),
		},
		{
			name:    "invalid json",
			json:    `123`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var key PSASoftwareComponentKeyType
			err := json.Unmarshal([]byte(tt.json), &key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, key)
			}
		})
	}
}

func TestNewPSASoftwareComponentKey(t *testing.T) {
	tests := []struct {
		name    string
		val     any
		want    PSASoftwareComponentKeyType
		wantErr bool
	}{
		{
			name: "nil value",
			val:  nil,
			want: PSASoftwareComponentKeyType(""),
		},
		{
			name: "string value",
			val:  PSASoftwareComponentType,
			want: PSASoftwareComponentKeyType(PSASoftwareComponentType),
		},
		{
			name: "PSASoftwareComponentKeyType value",
			val:  PSASoftwareComponentKeyType(PSASoftwareComponentType),
			want: PSASoftwareComponentKeyType(PSASoftwareComponentType),
		},
		{
			name: "pointer to PSASoftwareComponentKeyType",
			val:  func() *PSASoftwareComponentKeyType { k := PSASoftwareComponentKeyType(PSASoftwareComponentType); return &k }(),
			want: PSASoftwareComponentKeyType(PSASoftwareComponentType),
		},
		{
			name:    "invalid type",
			val:     123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := NewPSASoftwareComponentKey(tt.val)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, key)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, *key)
			}
		})
	}
}

func TestTaggedPSASoftwareComponentKey(t *testing.T) {
	key := TaggedPSASoftwareComponentKey(PSASoftwareComponentType)
	
	assert.Equal(t, PSASoftwareComponentType, key.String())
	assert.Equal(t, PSASoftwareComponentType, key.Type())
	assert.NoError(t, key.Valid())
}

func TestNewTaggedPSASoftwareComponentKey(t *testing.T) {
	key, err := NewTaggedPSASoftwareComponentKey(PSASoftwareComponentType)
	require.NoError(t, err)
	assert.Equal(t, TaggedPSASoftwareComponentKey(PSASoftwareComponentType), *key)
}

func TestNewMkeyPSASoftwareComponent(t *testing.T) {
	// Test the factory function
	mkey, err := newMkeyPSASoftwareComponent(PSASoftwareComponentType)
	require.NoError(t, err)
	assert.NotNil(t, mkey)
	assert.Equal(t, PSASoftwareComponentType, mkey.Value.Type())
}

func TestPSASoftwareComponentKeyIntegration(t *testing.T) {
	// Ensure PSA profile has been initialized by referencing it
	_ = ProfileID
	
	// Test that we can create a measurement key using the PSA software component type
	key, err := comid.NewMkey(PSASoftwareComponentType, PSASoftwareComponentType)
	require.NoError(t, err)
	assert.NotNil(t, key)
	assert.Equal(t, PSASoftwareComponentType, key.Value.Type())
}
