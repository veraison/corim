// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlags_UnmarshalJSON_skip_unknown(t *testing.T) {
	tv := []byte(`[ "notSecure", "mysteriousFlagWhichWillBeIgnored" ]`)

	flags := NewOpFlags().SetOpFlags(OpFlagNotSecure)
	require.NotNil(t, flags)
	expected := *flags

	var actual OpFlags
	err := actual.UnmarshalJSON(tv)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
	assert.True(t, actual.IsSet(OpFlagNotSecure))
}

func TestFlags_UnmarshalJSON_all_known(t *testing.T) {
	tv := []byte(`[ "notSecure", "notConfigured", "recovery", "debug" ]`)

	flags := NewOpFlags().
		SetOpFlags(OpFlagNotSecure).
		SetOpFlags(OpFlagNotConfigured).
		SetOpFlags(OpFlagRecovery).
		SetOpFlags(OpFlagDebug)
	require.NotNil(t, flags)
	expected := *flags

	var actual OpFlags
	err := actual.UnmarshalJSON(tv)

	fmt.Printf("CBOR: %02x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
	assert.True(t, actual.IsSet(OpFlagNotSecure))
	assert.True(t, actual.IsSet(OpFlagRecovery))
	assert.True(t, actual.IsSet(OpFlagDebug))
	assert.True(t, actual.IsSet(OpFlagNotConfigured))
}

func TestFlags_UnmarshalJSON_empty(t *testing.T) {
	tv := []byte(`[ ]`)

	flags := NewOpFlags()
	require.NotNil(t, flags)
	expected := *flags

	var actual OpFlags
	err := actual.UnmarshalJSON(tv)

	fmt.Printf("%02x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
	assert.False(t, actual.IsSet(OpFlagNotSecure))
	assert.False(t, actual.IsSet(OpFlagRecovery))
	assert.False(t, actual.IsSet(OpFlagDebug))
	assert.False(t, actual.IsSet(OpFlagNotConfigured))
}

func TestFlags_Valid_ok(t *testing.T) {
	// all valid flags combinations
	for i := 1; i <= 15; i++ {
		tv := OpFlags(i)

		assert.Nil(t, tv.Valid())
	}
}

func TestFlags_Valid_bad_combos(t *testing.T) {
	for i := 1; i <= 15; i++ {
		for j := 1; j <= 15; j++ {
			tv := OpFlags(i<<4 | j)

			expectedErr := fmt.Sprintf("op-flags has unknown bits asserted: %02x", tv)

			assert.EqualError(t, tv.Valid(), expectedErr)
		}
	}
}
