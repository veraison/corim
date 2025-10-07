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
