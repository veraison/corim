package comid

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testGroup uint64

func newTestGroup(_ any) (*Group, error) {
	v := testGroup(7)
	return &Group{&v}, nil
}

func (o testGroup) Type() string {
	return "test-value"
}

func (o testGroup) String() string {
	return "test"
}
func (o testGroup) Bytes() []byte {
	var ret [8]byte
	binary.BigEndian.PutUint64(ret[:], uint64(o))
	return ret[:]
}

func (o testGroup) Valid() error {
	return nil
}

type testGroupBadType struct {
	testGroup
}

func newTestGroupBadType(_ any) (*Group, error) {
	v := testGroupBadType{testGroup(7)}
	return &Group{&v}, nil
}

func (o testGroupBadType) Type() string {
	return "uuid"
}

func Test_RegisterGroupType(t *testing.T) {
	err := RegisterGroupType(32, newTestGroup)
	assert.EqualError(t, err, "tag 32 is already registered")

	err = RegisterGroupType(99993, newTestGroupBadType)
	assert.EqualError(t, err, `Group type with name "uuid" already exists`)

	err = RegisterGroupType(99993, newTestGroup)
	require.NoError(t, err)

}

func TestGroup_UmarshalJSON(t *testing.T) {
	var group Group

	err := group.UnmarshalJSON([]byte(`{`))
	assert.EqualError(t, err, "group decoding failure: unexpected end of JSON input")

	err = group.UnmarshalJSON([]byte(`{"type":"uuid","value":"aaaa"}`))
	assert.EqualError(t, err, "cannot unmarshal group: bad UUID: invalid UUID length: 4")
}

func Test_NewBytesGroup_OK(t *testing.T) {
	var testBytes = []byte{0x01, 0x02, 0x03, 0x04}

	for _, v := range []any{
		testBytes,
		&testBytes,
		string(testBytes),
	} {
		group, err := NewBytesGroup(v)
		require.NoError(t, err)
		got := group.Bytes()
		assert.Equal(t, testBytes, got)
	}
}

func Test_NewBytesGroup_NOK(t *testing.T) {
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
			_, err := NewBytesGroup(tv.Input)
			assert.EqualError(t, err, tv.Err)
		})
	}
}

func TestGroup_MarshalCBOR_Bytes(t *testing.T) {
	tv, err := NewBytesGroup(TestBytes)
	require.NoError(t, err)
	// 560 (h'458999786556')
	// tag(560): d9 0230
	expected := MustHexDecode(t, "d90230458999786556")

	actual, err := tv.MarshalCBOR()
	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestGroup_UnmarshalCBOR_Bytes_OK(t *testing.T) {
	tv := MustHexDecode(t, "d90230458999786556")

	var actual Group
	err := actual.UnmarshalCBOR(tv)

	assert.Nil(t, err)
	assert.Equal(t, "bytes", actual.Type())
	assert.Equal(t, TestBytes, actual.Bytes())
}

func TestGroup_MarshalJSONBytes_OK(t *testing.T) {
	testBytes := []byte{0x01, 0x02, 0x03}
	tb := TaggedBytes(testBytes)
	tv := Group{&tb}
	jsonBytes, err := tv.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `{"type":"bytes","value":"AQID"}`, string(jsonBytes))
}

func TestGroup_UnmarshalJSON_Bytes_OK(t *testing.T) {
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
			Input: `{ "type": "bytes", "value": "deadbeef"}`,
		},
	} {
		t.Run(tv.Name, func(t *testing.T) {
			var actual Group
			err := actual.UnmarshalJSON([]byte(tv.Input))
			require.NoError(t, err)
		})
	}
}

func TestGroup_UnmarshalJSON_Bytes_NOK(t *testing.T) {
	for _, tv := range []struct {
		Name  string
		Input string
		Err   string
	}{
		{
			Name:  "invalid value",
			Input: `{ "type": "bytes", "value": "/0" }`,
			Err:   "cannot unmarshal class id: illegal base64 data at input byte 0",
		},
		{
			Name:  "invalid input",
			Input: `{ "type": "bytes", "value": 10 }`,
			Err:   "cannot unmarshal class id: json: cannot unmarshal number into Go value of type comid.TaggedBytes",
		},
	} {
		t.Run(tv.Name, func(t *testing.T) {
			var actual ClassID
			err := actual.UnmarshalJSON([]byte(tv.Input))
			assert.EqualError(t, err, tv.Err)
		})
	}
}
