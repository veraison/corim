package comid

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var FlagTestFlag = Flag(-1)

type TestExtension struct {
	TestFlag *bool
}

func (o *TestExtension) ConstrainComid(_ *Comid) error {
	return errors.New("invalid")
}

func (o *TestExtension) ValidTriples(_ *Triples) error {
	return errors.New("invalid")
}

func (o *TestExtension) ConstrainMval(_ *Mval) error {
	return errors.New("invalid")
}

func (o *TestExtension) ConstrainFlagsMap(_ *FlagsMap) error {
	return errors.New("invalid")
}

func (o *TestExtension) ConstrainEntity(_ *Entity) error {
	return errors.New("invalid")
}

func (o *TestExtension) SetTrue(flag Flag) {
	if flag == FlagTestFlag {
		o.TestFlag = &True
	}
}
func (o *TestExtension) SetFalse(flag Flag) {
	if flag == FlagTestFlag {
		o.TestFlag = &False
	}
}

func (o *TestExtension) Clear(flag Flag) {
	if flag == FlagTestFlag {
		o.TestFlag = nil
	}
}

func (o *TestExtension) Get(flag Flag) *bool {
	if flag == FlagTestFlag {
		return o.TestFlag
	}

	return nil
}

func (o *TestExtension) AnySet() bool {
	return o.TestFlag != nil
}

func Test_Extensions(t *testing.T) {
	exts := Extensions{}
	exts.Register(&TestExtension{})

	err := exts.validComid(nil)
	assert.EqualError(t, err, "invalid")

	err = exts.validTriples(nil)
	assert.EqualError(t, err, "invalid")

	err = exts.validMval(nil)
	assert.EqualError(t, err, "invalid")

	err = exts.validEntity(nil)
	assert.EqualError(t, err, "invalid")

	err = exts.validFlagsMap(nil)
	assert.EqualError(t, err, "invalid")

	assert.False(t, exts.anySet())

	exts.setTrue(FlagTestFlag)

	exts.setFalse(FlagTestFlag)
	assert.False(t, *exts.get(FlagTestFlag))

	exts.clear(FlagTestFlag)
	assert.Nil(t, exts.get(FlagTestFlag))
	assert.False(t, exts.anySet())
}
