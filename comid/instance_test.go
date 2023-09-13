package comid

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstance_GetUUID_OK(t *testing.T) {
	inst := MustNewUUIDInstance(TestUUID)
	u, ok := inst.Value.(*TaggedUUID)
	assert.True(t, ok)
	assert.EqualValues(t, TestUUID, *u)
}

type testInstance string

func newTestInstance(val any) (*Instance, error) {
	ret := testInstance("test")
	return &Instance{&ret}, nil
}

func (o testInstance) Bytes() []byte {
	return []byte(o)
}

func (o testInstance) Type() string {
	return "test-instance"
}

func (o testInstance) String() string {
	return string(o)
}

func (o testInstance) Valid() error {
	return nil
}

func Test_RegisterInstanceType(t *testing.T) {
	err := RegisterInstanceType(99997, newTestInstance)
	require.NoError(t, err)

	instance, err := newTestInstance(nil)
	require.NoError(t, err)

	data, err := json.Marshal(instance)
	require.NoError(t, err)
	assert.Equal(t, string(data), `{"type":"test-instance","value":"test"}`)

	var out Instance
	err = json.Unmarshal(data, &out)
	require.NoError(t, err)
	assert.Equal(t, instance.Bytes(), out.Bytes())

	data, err = em.Marshal(instance)
	require.NoError(t, err)
	assert.Equal(t, data, []byte{
		0xda, 0x0, 0x1, 0x86, 0x9d, // tag 99997
		0x64,                   // tstr(4)
		0x74, 0x65, 0x73, 0x74, // "test"
	})

	var out2 Instance
	err = dm.Unmarshal(data, &out2)
	require.NoError(t, err)
	assert.Equal(t, instance.Bytes(), out2.Bytes())
}
