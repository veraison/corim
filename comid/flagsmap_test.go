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
