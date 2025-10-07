// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FlagsMap_IsConfidentialityProtected(t *testing.T) {
	fm := NewFlagsMap()

	// Test initial state - should be nil (unset)
	assert.Nil(t, fm.Get(FlagIsConfidentialityProtected))
	assert.False(t, fm.AnySet())

	// Test setting to true
	fm.SetTrue(FlagIsConfidentialityProtected)
	assert.True(t, fm.AnySet())
	assert.Equal(t, true, *fm.Get(FlagIsConfidentialityProtected))

	// Test setting to false
	fm.SetFalse(FlagIsConfidentialityProtected)
	assert.True(t, fm.AnySet())
	assert.Equal(t, false, *fm.Get(FlagIsConfidentialityProtected))

	// Test clearing
	fm.Clear(FlagIsConfidentialityProtected)
	assert.False(t, fm.AnySet())
	assert.Nil(t, fm.Get(FlagIsConfidentialityProtected))
}

func Test_FlagsMap_IsConfidentialityProtected_Serialization(t *testing.T) {
	fm := NewFlagsMap()
	fm.SetTrue(FlagIsConfidentialityProtected)

	// Test CBOR serialization
	cbor, err := fm.MarshalCBOR()
	assert.NoError(t, err)
	assert.NotNil(t, cbor)

	var fm2 FlagsMap
	err = fm2.UnmarshalCBOR(cbor)
	assert.NoError(t, err)
	assert.Equal(t, true, *fm2.Get(FlagIsConfidentialityProtected))

	// Test JSON serialization
	json, err := fm.MarshalJSON()
	assert.NoError(t, err)
	assert.NotNil(t, json)
	assert.Contains(t, string(json), "is-confidentiality-protected")

	var fm3 FlagsMap
	err = fm3.UnmarshalJSON(json)
	assert.NoError(t, err)
	assert.Equal(t, true, *fm3.Get(FlagIsConfidentialityProtected))
}
