// Copyright 2024-2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FlagsMap(t *testing.T) {
	fm := NewFlagsMap()
	assert.False(t, fm.AnySet())

	for _, flag := range []Flag{
		FlagIsConfigured,
		FlagIsSecure,
		FlagIsRecovery,
		FlagIsDebug,
		FlagIsReplayProtected,
		FlagIsIntegrityProtected,
		FlagIsRuntimeMeasured,
		FlagIsImmutable,
		FlagIsTcb,
		FlagIsConfidentialityProtected,
		FlagIsRuntimeUpdatable,
	} {
		fm.SetTrue(flag)
		assert.True(t, fm.AnySet())
		assert.Equal(t, true, *fm.Get(flag))

		fm.SetFalse(flag)
		assert.True(t, fm.AnySet())
		assert.Equal(t, false, *fm.Get(flag))

		fm.Clear(flag)
		assert.False(t, fm.AnySet())
		assert.Equal(t, (*bool)(nil), fm.Get(flag))
	}

	fm.SetTrue(Flag(-1))
	fm.SetFalse(Flag(-1))
	assert.False(t, fm.AnySet())
	assert.Equal(t, (*bool)(nil), fm.Get(Flag(-1)))
}

func Test_FlagsMap_Equal_True(t *testing.T) {
	claim := NewFlagsMap()
	ref := NewFlagsMap()

	claim.SetTrue(FlagIsSecure)
	claim.SetTrue(FlagIsRuntimeMeasured)

	ref.SetTrue(FlagIsSecure)
	ref.SetTrue(FlagIsRuntimeMeasured)

	assert.True(t, claim.Equal(*ref))
}

func Test_FlagsMap_Equal_False(t *testing.T) {
	claim := NewFlagsMap()
	ref := NewFlagsMap()

	claim.SetTrue(FlagIsSecure)
	claim.SetTrue(FlagIsRuntimeMeasured)

	ref.SetTrue(FlagIsSecure)

	assert.False(t, claim.Equal(*ref))
}

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
