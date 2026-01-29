// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClass_MarshalCBOR_UUID_full(t *testing.T) {
	tv := NewClassUUID(TestUUID).
		SetVendor("ACME Ltd").
		SetModel("Roadrunner").
		SetLayer(1).
		SetIndex(2)
	require.NotNil(t, tv)

	// {0: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA'), 1: "ACME Ltd", 2: "Roadrunner", 3: 1, 4: 2}
	expected := MustHexDecode(t, "a500d8255031fb5abf023e4992aa4e95f9c1503bfa016841434d45204c7464026a526f616472756e6e657203010402")

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClass_MarshalCBOR_OID_full(t *testing.T) {
	tv := NewClassOID(TestOID).
		SetVendor("EMCA Ltd").
		SetModel("Rennurdaor").
		SetLayer(2).
		SetIndex(1)
	require.NotNil(t, tv)

	expected := MustHexDecode(t, "a500d86f445502c0000168454d4341204c7464026a52656e6e757264616f7203020401")

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClass_MarshalCBOR_empty(t *testing.T) {
	var tv Class

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.EqualError(t, err, "class must not be empty")
}

func TestClass_MarshalCBOR_ClassID_only(t *testing.T) {
	tv := NewClassUUID(TestUUID)
	require.NotNil(t, tv)

	// {0: 37(h'31FB5ABF023E4992AA4E95F9C1503BFA')}
	expected := MustHexDecode(t, "a100d8255031fb5abf023e4992aa4e95f9c1503bfa")

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClass_MarshalCBOR_Vendor_only(t *testing.T) {
	var tv Class

	require.NotNil(t, tv.SetVendor("ACME Ltd."))

	// {1: "ACME Ltd."}
	expected := MustHexDecode(t, "a1016941434d45204c74642e")

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClass_MarshalCBOR_Model_only(t *testing.T) {
	var tv Class

	require.NotNil(t, tv.SetModel("Roadrunner"))

	// {2: "Roadrunner"}
	expected := MustHexDecode(t, "a1026a526f616472756e6e6572")

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClass_MarshalCBOR_Layer_only(t *testing.T) {
	var tv Class

	require.NotNil(t, tv.SetLayer(5))

	// {3: 5}
	expected := MustHexDecode(t, "a10305")

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClass_MarshalCBOR_Index_only(t *testing.T) {
	var tv Class

	require.NotNil(t, tv.SetIndex(3))

	// {4: 3}
	expected := MustHexDecode(t, "a10403")

	actual, err := tv.ToCBOR()

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClassID_UnmarshalCBOR_full(t *testing.T) {
	// {0: 560(h'61636D652D696D706C656D656E746174696F6E2D69642D303030303030303031'), 1: "EMCA Ltd", 2: "Rennurdaor", 3: 2, 4: 1}
	tv := MustHexDecode(t, "a500d90230582061636d652d696d706c656d656e746174696f6e2d69642d3030303030303030310168454d4341204c7464026a52656e6e757264616f7203020401")

	var actual Class
	err := actual.FromCBOR(tv)

	assert.Nil(t, err)
	assert.NotNil(t, actual.ClassID)
	assert.Equal(t, string(TestImplID[:]), actual.ClassID.String())
	assert.NotNil(t, actual.Vendor)
	assert.Equal(t, "EMCA Ltd", actual.GetVendor())
	assert.NotNil(t, actual.Model)
	assert.Equal(t, "Rennurdaor", actual.GetModel())
	assert.NotNil(t, actual.Layer)
	assert.Equal(t, uint64(2), actual.GetLayer())
	assert.NotNil(t, actual.Index)
	assert.Equal(t, uint64(1), actual.GetIndex())
}

// TODO(tho) optional fields missing

func TestClass_UnmarshalJSON_empty(t *testing.T) {
	tv := `{}`

	var actual Class
	err := actual.FromJSON([]byte(tv))

	assert.EqualError(t, err, "class must not be empty")
}

func TestClass_UnmarshalJSON_spurious_map(t *testing.T) {
	tv := `{ "What is this?": "I have no idea" }`

	var actual Class
	err := actual.FromJSON([]byte(tv))

	assert.EqualError(t, err, "class must not be empty")
}

func TestClass_UnmarshalJSON_spurious_array(t *testing.T) {
	tv := `[ "x", 1, 5.123 ]`

	var actual Class
	err := actual.FromJSON([]byte(tv))

	assert.EqualError(t, err, "json: cannot unmarshal array into Go value of type comid.Class")
}

func TestClass_UnmarshalJSON_full(t *testing.T) {
	tv := `
{
	"id": {
		"type": "uuid",
		"value": "83294297-97eb-42ef-8a72-ae9fea002750"
	},
	"vendor": "ACME Ltd.",
	"model": "RoadRunner Boot ROM",
	"layer": 4,
	"index": 2
}
`
	var actual Class
	err := actual.FromJSON([]byte(tv))

	assert.Nil(t, err)
	assert.NotNil(t, actual.ClassID)
	assert.Equal(t, "83294297-97eb-42ef-8a72-ae9fea002750", actual.ClassID.String())
	assert.NotNil(t, actual.Vendor)
	assert.Equal(t, "ACME Ltd.", actual.GetVendor())
	assert.NotNil(t, actual.Model)
	assert.Equal(t, "RoadRunner Boot ROM", actual.GetModel())
	assert.NotNil(t, actual.Layer)
	assert.Equal(t, uint64(4), actual.GetLayer())
	assert.NotNil(t, actual.Index)
	assert.Equal(t, uint64(2), actual.GetIndex())
}

func TestClass_NewClassBytes(t *testing.T) {
	// Valid bytes
	c := NewClassBytes([]byte{0x01, 0x02, 0x03})
	assert.NotNil(t, c)

	// Invalid type returns nil
	c = NewClassBytes(12345)
	assert.Nil(t, c)
}

func TestClass_ToJSON(t *testing.T) {
	// Valid class
	c := NewClassUUID(TestUUID).SetVendor("ACME Ltd")
	require.NotNil(t, c)

	data, err := c.ToJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Invalid/empty class
	var empty Class
	_, err = empty.ToJSON()
	assert.Error(t, err)
}
