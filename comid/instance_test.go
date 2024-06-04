package comid

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstance_SetUUID_OK(t *testing.T) {
	inst := &Instance{}
	testUUID, err := uuid.Parse(TestUUIDString)
	require.NoError(t, err)
	i := inst.SetUUID(testUUID)
	require.NotNil(t, i)
}

func TestInstance_GetUUID_OK(t *testing.T) {
	inst := MustNewUUIDInstance(TestUUID)
	require.NotNil(t, inst)
	u, err := inst.GetUUID()
	assert.Nil(t, err)
	assert.Equal(t, u, TestUUID)
}

func TestInstance_GetUUID_NOK(t *testing.T) {
	inst := &Instance{}
	expectedErr := "instance-id type is: <nil>"
	_, err := inst.GetUUID()
	assert.EqualError(t, err, expectedErr)
}

func TestInstance_SetGetUEID_OK(t *testing.T) {
	inst := &Instance{}
	inst = inst.SetUEID(TestUEID)
	require.NotNil(t, inst)
	expectedUEID, err := inst.GetUEID()
	require.NoError(t, err)
	assert.Equal(t, TestUEID, expectedUEID)
}

type testInstance string

func newTestInstance(_ any) (*Instance, error) {
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

func Test_NewBytesInstance_OK(t *testing.T) {
	var testBytes = []byte{0x01, 0x02, 0x03, 0x04}

	for _, v := range []any{
		testBytes,
		&testBytes,
		string(testBytes),
	} {
		instance, err := NewBytesInstance(v)
		require.NoError(t, err)
		got := instance.Bytes()
		assert.Equal(t, testBytes[:], got)
	}
}

func Test_NewBytesInstance_NOK(t *testing.T) {
	for _, tv := range []struct {
		Name  string
		Input any
		Err   string
	}{
		{
			Name:  "invalid input integer",
			Input: 7,
			Err:   "unexpected type for bytes: int",
		},
		{
			Name:  "invalid input fixed array",
			Input: [3]byte{0x01, 0x02, 0x03},
			Err:   "unexpected type for bytes: [3]uint8",
		},
	} {
		t.Run(tv.Name, func(t *testing.T) {
			_, err := NewBytesInstance(tv.Input)
			assert.EqualError(t, err, tv.Err)
		})
	}
}

func TestInstance_MarshalCBOR_Bytes(t *testing.T) {
	tv, err := NewBytesInstance(TestBytes)
	require.NoError(t, err)
	// 560 (h'458999786556')
	// tag(560): d9 0230
	expected := MustHexDecode(t, "d90230458999786556")

	actual, err := tv.MarshalCBOR()
	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestInstance_UnmarshalCBOR_Bytes_OK(t *testing.T) {
	tv := MustHexDecode(t, "d90230458999786556")

	var actual Instance
	err := actual.UnmarshalCBOR(tv)

	assert.Nil(t, err)
	assert.Equal(t, "bytes", actual.Type())
	assert.Equal(t, TestBytes, actual.Bytes())
}

func TestInstance_MarshalJSONBytes_OK(t *testing.T) {
	testBytes := []byte{0x01, 0x02, 0x03}
	tb := TaggedBytes(testBytes)
	tv := Instance{&tb}
	jsonBytes, err := tv.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `{"type":"bytes","value":"AQID"}`, string(jsonBytes))

}

func TestInstance_UnmarshalJSON_Bytes_OK(t *testing.T) {
	for _, tv := range []struct {
		Name  string
		Input string
	}{
		{
			Name:  "valid input test 1",
			Input: `{ "type": "bytes", "value": "MTIzNDU2Nzg5" }`,
		},
		{
			Name:  "valid input test 2",
			Input: `{ "type": "bytes", "value": "CgsMDQ4="}`,
		},
	} {
		t.Run(tv.Name, func(t *testing.T) {
			var actual Instance
			err := actual.UnmarshalJSON([]byte(tv.Input))
			require.NoError(t, err)
		})
	}
}

func TestInstance_UnmarshalJSON_Bytes_NOK(t *testing.T) {
	for _, tv := range []struct {
		Name  string
		Input string
		Err   string
	}{
		{
			Name:  "invalid value",
			Input: `{ "type": "bytes", "value": "/0" }`,
			Err:   "cannot unmarshal instance: illegal base64 data at input byte 0",
		},
		{
			Name:  "invalid input",
			Input: `{ "type": "bytes", "value": 10 }`,
			Err:   "cannot unmarshal instance: json: cannot unmarshal number into Go value of type comid.TaggedBytes",
		},
	} {
		t.Run(tv.Name, func(t *testing.T) {
			var actual Instance
			err := actual.UnmarshalJSON([]byte(tv.Input))
			assert.EqualError(t, err, tv.Err)
		})
	}
}
