// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClassID_MarshalCBOR_UUID(t *testing.T) {
	var tv ClassID

	require.NotNil(t, tv.SetUUID(TestUUID))

	// 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')
	// tag(37): d8 25
	expected := MustHexDecode(t, "d8255031fb5abf023e4992aa4e95f9c1503bfa")

	actual, err := tv.MarshalCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClassID_MarshalCBOR_ImplID(t *testing.T) {
	var tv ClassID

	require.NotNil(t, tv.SetImplID(TestImplID))

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
	assert.Equal(t, ClassIDTypeUUID, actual.Type())
	assert.Equal(t, TestUUIDString, actual.String())
}

func TestClassID_UnmarshalCBOR_ImplID_OK(t *testing.T) {
	tv := MustHexDecode(t, "d90258582061636d652d696d706c656d656e746174696f6e2d69642d303030303030303031")

	expected := b64TestImplID()

	var actual ClassID
	err := actual.UnmarshalCBOR(tv)

	assert.Nil(t, err)
	assert.Equal(t, ClassIDTypeImplID, actual.Type())
	assert.Equal(t, expected, actual.String())
}

func TestClassID_UnmarshalCBOR_badInput(t *testing.T) {
	// raw, no tag
	hex := "582061636d652d696d706c656d656e746174696f6e2d69642d303030303030303031"
	tv := MustHexDecode(t, hex)

	expectedError := fmt.Sprintf("unknown class id (CBOR: %s)", hex)

	var actual ClassID

	err := actual.UnmarshalCBOR(tv)
	assert.EqualError(t, err, expectedError)
}

func TestClassID_UnmarshalJSON_UUID(t *testing.T) {
	tvFmt := `{ "type": "uuid", "value": "%s" }`

	tv := fmt.Sprintf(tvFmt, TestUUIDString)

	var actual ClassID
	err := actual.UnmarshalJSON([]byte(tv))

	assert.Nil(t, err)
	assert.Equal(t, ClassIDTypeUUID, actual.Type())
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
	assert.Equal(t, ClassIDTypeImplID, actual.Type())
	// the returned string is the base64 encoding of the stored binary
	assert.Equal(t, expected, actual.String())
}

func TestClassID_UnmarshalJSON_badInput_unknown_type(t *testing.T) {
	tv := `{ "type": "FOOBAR", "value": "1234567890" }`

	var actual ClassID
	err := actual.UnmarshalJSON([]byte(tv))

	assert.EqualError(t, err, "unknown type 'FOOBAR' for class id")
	assert.Equal(t, ClassIDTypeUnknown, actual.Type())
}

func TestClassID_UnmarshalJSON_badInput_missing_value(t *testing.T) {
	tv := `{ "type": "psa.impl-id" }`

	var actual ClassID
	err := actual.UnmarshalJSON([]byte(tv))

	assert.EqualError(t, err, "bad ImplID: unexpected end of JSON input")
	assert.Equal(t, ClassIDTypeUnknown, actual.Type())
}

func TestClassID_UnmarshalJSON_badInput_empty_value(t *testing.T) {
	tv := `{ "type": "psa.impl-id", "value": "" }`

	var actual ClassID
	err := actual.UnmarshalJSON([]byte(tv))

	assert.EqualError(t, err, "bad ImplID format: got 0 bytes, want 32")
	assert.Equal(t, ClassIDTypeUnknown, actual.Type())
}

func TestClassID_UnmarshalJSON_badInput_badly_encoded_ImplID_value(t *testing.T) {
	tv := `{ "type": "psa.impl-id", "value": ";" }`

	var actual ClassID
	err := actual.UnmarshalJSON([]byte(tv))

	assert.EqualError(t, err, "bad ImplID: illegal base64 data at input byte 0")
	assert.Equal(t, ClassIDTypeUnknown, actual.Type())
}

func TestClassID_UnmarshalJSON_badInput_badly_encoded_UUID_value(t *testing.T) {
	tv := `{ "type": "uuid", "value": "abcd-1234" }`

	var actual ClassID
	err := actual.UnmarshalJSON([]byte(tv))

	assert.EqualError(t, err, "bad UUID: invalid UUID length: 9")
	assert.Equal(t, ClassIDTypeUnknown, actual.Type())
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
		c := ClassID{}
		assert.NotNil(t, c.SetOID(tv))
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
		c := ClassID{}
		assert.Nil(t, c.SetOID(tv))
	}
}
