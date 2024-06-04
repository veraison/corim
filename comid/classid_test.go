// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClassID_MarshalCBOR_UUID(t *testing.T) {
	tv := MustNewUUIDClassID(TestUUID)

	// 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')
	// tag(37): d8 25
	expected := MustHexDecode(t, "d8255031fb5abf023e4992aa4e95f9c1503bfa")

	actual, err := tv.MarshalCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClassID_MarshalCBOR_ImplID(t *testing.T) {
	tv := MustNewImplIDClassID(TestImplID)

	// 600 (h'61636D652D696D706C656D656E746174696F6E2D69642D303030303030303031')
	// tag(600): d9 0258
	expected := MustHexDecode(t, "d90258582061636d652d696d706c656d656e746174696f6e2d69642d303030303030303031")

	actual, err := tv.MarshalCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClassID_MarshalCBOR_Empty(t *testing.T) {
	var tv ClassID

	// null (primitive 22)
	expected := MustHexDecode(t, "f6")

	actual, err := tv.MarshalCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClassID_UnmarshalCBOR_UUID_OK(t *testing.T) {
	tv := MustHexDecode(t, "d8255031fb5abf023e4992aa4e95f9c1503bfa")

	var actual ClassID
	err := actual.UnmarshalCBOR(tv)

	assert.Nil(t, err)
	assert.Equal(t, "uuid", actual.Type())
	assert.Equal(t, TestUUIDString, actual.String())
}

func TestClassID_UnmarshalCBOR_ImplID_OK(t *testing.T) {
	tv := MustHexDecode(t, "d90258582061636d652d696d706c656d656e746174696f6e2d69642d303030303030303031")

	expected := b64TestImplID()

	var actual ClassID
	err := actual.UnmarshalCBOR(tv)

	assert.Nil(t, err)
	assert.Equal(t, "psa.impl-id", actual.Type())
	assert.Equal(t, expected, actual.String())
}

func TestClassID_UnmarshalCBOR_badInput(t *testing.T) {
	// raw, no tag
	hex := "582061636d652d696d706c656d656e746174696f6e2d69642d303030303030303031"
	tv := MustHexDecode(t, hex)

	var actual ClassID

	err := actual.UnmarshalCBOR(tv)
	assert.EqualError(t, err, "cbor: cannot unmarshal byte string into Go value of type comid.IClassIDValue")
}

func TestClassID_UnmarshalJSON_UUID(t *testing.T) {
	tvFmt := `{ "type": "uuid", "value": "%s" }`

	tv := fmt.Sprintf(tvFmt, TestUUIDString)

	var actual ClassID
	err := actual.UnmarshalJSON([]byte(tv))

	assert.Nil(t, err)
	assert.Equal(t, "uuid", actual.Type())
	assert.Equal(t, TestUUIDString, actual.String())
}

func TestClassID_UnmarshalJSON_ImplID(t *testing.T) {
	b64ImplID := b64TestImplID()

	tvFmt := `{ "type": "psa.impl-id", "value": "%s" }`
	// in JSON, impl-id is base64 encoded
	tv := fmt.Sprintf(tvFmt, b64ImplID)

	expected := b64ImplID

	var actual ClassID
	err := actual.UnmarshalJSON([]byte(tv))

	assert.Nil(t, err)
	assert.Equal(t, "psa.impl-id", actual.Type())
	// the returned string is the base64 encoding of the stored binary
	assert.Equal(t, expected, actual.String())
}

func TestClassID_UnmarshalJSON_badInput_unknown_type(t *testing.T) {
	tv := `{ "type": "FOOBAR", "value": "1234567890" }`

	var actual ClassID
	err := actual.UnmarshalJSON([]byte(tv))

	assert.EqualError(t, err, "unknown class id type: FOOBAR")
	assert.Equal(t, "", actual.Type())
}

func TestClassID_UnmarshalJSON_badInput_missing_value(t *testing.T) {
	tv := `{ "type": "psa.impl-id" }`

	var actual ClassID
	err := actual.UnmarshalJSON([]byte(tv))

	assert.EqualError(t, err, "class id decoding failure: no value provided for psa.impl-id")
	assert.Equal(t, "", actual.Type())
}

func TestClassID_UnmarshalJSON_badInput_empty_value(t *testing.T) {
	tv := `{ "type": "psa.impl-id", "value": "" }`

	var actual ClassID
	err := actual.UnmarshalJSON([]byte(tv))

	assert.EqualError(t, err, "cannot unmarshal class id: bad psa.impl-id: decoded 0 bytes, want 32")
	assert.Equal(t, "", actual.Type())
}

func TestClassID_UnmarshalJSON_badInput_badly_encoded_ImplID_value(t *testing.T) {
	tv := `{ "type": "psa.impl-id", "value": ";" }`

	var actual ClassID
	err := actual.UnmarshalJSON([]byte(tv))

	assert.EqualError(t, err, "cannot unmarshal class id: illegal base64 data at input byte 0")
	assert.Equal(t, "", actual.Type())
}

func TestClassID_UnmarshalJSON_badInput_badly_encoded_UUID_value(t *testing.T) {
	tv := `{ "type": "uuid", "value": "abcd-1234" }`

	var actual ClassID
	err := actual.UnmarshalJSON([]byte(tv))

	assert.EqualError(t, err, "cannot unmarshal class id: bad UUID: invalid UUID length: 9")
	assert.Equal(t, "", actual.Type())
}

func TestClassID_SetOID_ok(t *testing.T) {
	tvs := []string{
		"1.2.3",
		"1.2.3.4",
		"1.2.3.4.5",
		"1.2.3.4.5.6",
		"1.2.3.4.5.6.7",
		"1.2.3.4.5.6.7.8",
		"1.2.3.4.5.6.7.8.9",
		"1.2.3.4.5.6.7.8.9.10",
		"1.2.3.4.5.6.7.8.9.10.11",
		"1.2.3.4.5.6.7.8.9.10.11.12",
		"1.2.3.4.5.6.7.8.9.10.11.12.13",
		"1.2.3.4.5.6.7.8.9.10.11.12.13.14",
		"1.2.3.4.5.6.7.8.9.10.11.12.13.14.15",
		"1.2.3.4.5.6.7.8.9.10.11.12.13.14.15.16",
		"1.2.3.4.5.6.7.8.9.10.11.12.13.14.15.16.17",
		"1.2.3.4.5.6.7.8.9.10.11.12.13.14.15.16.17.18",
		"1.2.3.4.5.6.7.8.9.10.11.12.13.14.15.16.17.18.19",
		"1.2.3.4.5.6.7.8.9.10.11.12.13.14.15.16.17.18.19.20",
	}

	for _, tv := range tvs {
		c := MustNewOIDClassID(tv)
		assert.Equal(t, tv, c.String())
	}
}

func TestClassID_SetOID_bad(t *testing.T) {
	tvs := []string{
		"",                             // empty
		"1",                            // too little
		"1.2",                          // still too little
		"1.2.-3",                       // negative arc
		".1.2.3",                       // not absolute
		"iso(1) org(3) dod(6) iana(1)", // not dotted decimal
		"1...",
		"a.b.c",
	}

	for _, tv := range tvs {
		c, err := NewOIDClassID(tv)
		assert.NotNil(t, err)
		assert.Nil(t, c)
	}
}

func Test_NewImplIDClassID(t *testing.T) {
	classID, err := NewImplIDClassID(nil)
	expected := [32]byte{}
	require.NoError(t, err)
	assert.Equal(t, expected[:], classID.Bytes())

	taggedImplID := TaggedImplID(TestImplID)

	for _, v := range []any{
		TestImplID,
		&TestImplID,
		taggedImplID,
		&taggedImplID,
		taggedImplID.Bytes(),
	} {
		classID, err = NewImplIDClassID(v)
		require.NoError(t, err)
		assert.Equal(t, taggedImplID.Bytes(), classID.Bytes())
	}

	expected = [32]byte{
		0x61, 0x63, 0x6d, 0x65, 0x2d, 0x69, 0x6d, 0x70,
		0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74,
		0x69, 0x6f, 0x6e, 0x2d, 0x69, 0x64, 0x2d, 0x30,
		0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31,
	}
	classID, err = NewImplIDClassID("YWNtZS1pbXBsZW1lbnRhdGlvbi1pZC0wMDAwMDAwMDE=")
	require.NoError(t, err)
	assert.Equal(t, expected[:], classID.Bytes())

	_, err = NewImplIDClassID(7)
	assert.EqualError(t, err, "unexpected type for psa.impl-id: int")
}

func Test_NewUUIDClassID(t *testing.T) {
	classID, err := NewUUIDClassID(nil)

	expected := [16]byte{}
	require.NoError(t, err)
	assert.Equal(t, expected[:], classID.Bytes())

	taggedUUID := TaggedUUID(TestUUID)

	for _, v := range []any{
		TestUUID,
		&TestUUID,
		taggedUUID,
		&taggedUUID,
		taggedUUID.Bytes(),
	} {
		classID, err = NewUUIDClassID(v)
		require.NoError(t, err)
		assert.Equal(t, taggedUUID.Bytes(), classID.Bytes())
	}

	classID, err = NewUUIDClassID(taggedUUID.String())
	require.NoError(t, err)
	assert.Equal(t, taggedUUID.Bytes(), classID.Bytes())
}

func Test_NewOIDClassID(t *testing.T) {
	classID, err := NewOIDClassID(nil)

	expected := []byte{}
	require.NoError(t, err)
	assert.Equal(t, expected, classID.Bytes())

	var oid OID
	require.NoError(t, oid.FromString(TestOID))
	taggedOID := TaggedOID(oid)

	for _, v := range []any{
		TestOID,
		oid,
		&oid,
		taggedOID,
		&taggedOID,
		taggedOID.Bytes(),
	} {
		classID, err = NewOIDClassID(v)
		require.NoError(t, err)
		expected := taggedOID.Bytes()
		got := classID.Bytes()
		assert.Equal(t, expected, got)
	}

	classID, err = NewOIDClassID(taggedOID.String())
	require.NoError(t, err)
	assert.Equal(t, taggedOID.Bytes(), classID.Bytes())
}

func Test_NewIntClassID(t *testing.T) {
	classID, err := NewIntClassID(nil)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, classID.Bytes())

	testInt := 7
	testInt64 := int64(7)
	testUint64 := uint64(7)

	var testBytes [8]byte
	binary.BigEndian.PutUint64(testBytes[:], testUint64)

	for _, v := range []any{
		testInt,
		&testInt,
		testInt64,
		&testInt64,
		testUint64,
		&testUint64,
		"7",
		testBytes[:],
	} {
		classID, err = NewIntClassID(v)
		require.NoError(t, err)
		got := classID.Bytes()
		assert.Equal(t, testBytes[:], got)
	}
}

func Test_TaggedInt(t *testing.T) {
	val := TaggedInt(7)
	assert.Equal(t, "7", val.String())
	assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07}, val.Bytes())
	assert.Equal(t, "int", val.Type())
	assert.NoError(t, val.Valid())

	classID := ClassID{&val}

	bytes, err := em.Marshal(classID)
	require.NoError(t, err)
	assert.Equal(t, []byte{
		0xd9, 0x02, 0x27, // tag 551
		0x07, // int 7
	}, bytes)

	var out ClassID
	err = dm.Unmarshal(bytes, &out)
	require.NoError(t, err)
	assert.Equal(t, classID, out)

	jsonBytes, err := json.Marshal(classID)
	require.NoError(t, err)
	assert.Equal(t, `{"type":"int","value":7}`, string(jsonBytes))

	out = ClassID{}
	err = json.Unmarshal(jsonBytes, &out)
	require.NoError(t, err)
	assert.Equal(t, classID, out)
}

type testClassID [4]byte

func newTestClassID(_ any) (*ClassID, error) {
	return &ClassID{&testClassID{0x74, 0x65, 0x73, 0x74}}, nil
}

func (o testClassID) Bytes() []byte {
	return o[:]
}

func (o testClassID) Type() string {
	return "test-class-id"
}

func (o testClassID) String() string {
	return "test"
}

func (o testClassID) Valid() error {
	return nil
}

func (o testClassID) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *testClassID) UnmarshalJSON(data []byte) error {
	var out string
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}

	if len(out) != 4 {
		return fmt.Errorf("bad testClassID: decoded %d bytes, want 4", len(out))
	}

	copy((*o)[:], out)

	return nil
}

func Test_RegisterClassIDType(t *testing.T) {
	err := RegisterClassIDType(99999, newTestClassID)
	require.NoError(t, err)

	classID, err := newTestClassID(nil)
	require.NoError(t, err)

	data, err := json.Marshal(classID)
	require.NoError(t, err)
	assert.Equal(t, string(data), `{"type":"test-class-id","value":"test"}`)

	var out ClassID
	err = json.Unmarshal(data, &out)
	require.NoError(t, err)
	assert.Equal(t, classID.Bytes(), out.Bytes())

	data, err = em.Marshal(classID)
	require.NoError(t, err)
	assert.Equal(t, data, []byte{
		0xda, 0x0, 0x1, 0x86, 0x9f, // tag 99999
		0x44,                   // bstr(4)
		0x74, 0x65, 0x73, 0x74, // "test"
	})

	var out2 ClassID
	err = dm.Unmarshal(data, &out2)
	require.NoError(t, err)
	assert.Equal(t, classID.Bytes(), out2.Bytes())
}

func Test_NewBytesClassID_OK(t *testing.T) {
	var testBytes = []byte{0x01, 0x02, 0x03, 0x04}

	for _, v := range []any{
		testBytes,
		&testBytes,
		string(testBytes),
	} {
		classID, err := NewBytesClassID(v)
		require.NoError(t, err)
		got := classID.Bytes()
		assert.Equal(t, testBytes, got)
	}
}

func Test_NewBytesClassID_NOK(t *testing.T) {
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
			_, err := NewBytesClassID(tv.Input)
			assert.EqualError(t, err, tv.Err)
		})
	}
}

func TestClassID_MarshalCBOR_Bytes(t *testing.T) {
	tv, err := NewBytesClassID(TestBytes)
	require.NoError(t, err)
	// 560 (h'458999786556')
	// tag(560): d9 0230
	expected := MustHexDecode(t, "d90230458999786556")

	actual, err := tv.MarshalCBOR()
	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClassID_UnmarshalCBOR_Bytes_OK(t *testing.T) {
	tv := MustHexDecode(t, "d90230458999786556")

	var actual ClassID
	err := actual.UnmarshalCBOR(tv)

	assert.Nil(t, err)
	assert.Equal(t, "bytes", actual.Type())
	assert.Equal(t, TestBytes, actual.Bytes())
}

func TestClassID_MarshalJSONBytes_OK(t *testing.T) {
	testBytes := []byte{0x01, 0x02, 0x03}
	tb := TaggedBytes(testBytes)
	tv := ClassID{&tb}
	jsonBytes, err := tv.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `{"type":"bytes","value":"AQID"}`, string(jsonBytes))

}

func TestClassID_UnmarshalJSON_Bytes_OK(t *testing.T) {
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
			var actual ClassID
			err := actual.UnmarshalJSON([]byte(tv.Input))
			require.NoError(t, err)
		})
	}
}

func TestClassID_UnmarshalJSON_Bytes_NOK(t *testing.T) {
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
