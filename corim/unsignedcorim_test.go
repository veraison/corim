// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"fmt"
	"testing"

	"github.com/veraison/corim/comid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnsignedCorim_id_string(t *testing.T) {
	testIDString := "test string"

	tv := NewUnsignedCorim().SetIDString(testIDString)
	require.NotNil(t, tv)

	actual, err := tv.GetIDString()
	assert.Nil(t, err)
	assert.Equal(t, testIDString, actual)
}

func TestUnsignedCorim_id_string_empty(t *testing.T) {
	emptyString := ""

	tv := NewUnsignedCorim()
	require.NotNil(t, tv)

	assert.Nil(t, tv.SetIDString(emptyString))
}

func TestUnsignedCorim_id_uuid(t *testing.T) {
	tv := NewUnsignedCorim().SetIDUUID(comid.TestUUID)
	require.NotNil(t, tv)

	actual, err := tv.GetIDUUID()
	assert.Nil(t, err)
	assert.Equal(t, comid.TestUUID, actual)
}

func TestUnsignedCorim_id_uuid_empty(t *testing.T) {
	emptyUUID := comid.UUID{}

	tv := NewUnsignedCorim()
	require.NotNil(t, tv)

	assert.Nil(t, tv.SetIDUUID(emptyUUID))
}

func TestUnsignedCorim_AddComid_and_marshal(t *testing.T) {
	c := comid.Comid{}
	err := c.FromJSON([]byte(comid.PSARefValJSONTemplate))
	require.Nil(t, err)

	tv := NewUnsignedCorim().SetIDString("test corim id")
	require.NotNil(t, tv)

	assert.NotNil(t, tv.AddComid(TaggedComid(c)))

	actual, err := em.Marshal(tv)
	assert.Nil(t, err)

	fmt.Printf("CBOR: %x", actual)

	expected := comid.MustHexDecode(t, "a2006d7465737420636f72696d20696401815901b8d901faa40065656e2d474201a100782434334242453337462d324536312d344233332d414544332d3533434646313432384231360281a3006941434d45204c74642e01d8207468747470733a2f2f61636d652e6578616d706c65028300010204a1008182a100a300d90227582061636d652d696d706c656d656e746174696f6e2d69642d303030303030303031016441434d45026a526f616452756e6e657283a200d90258a30162424c0465322e312e30055820acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b01a102818201582087428fc522803d31065e7bce3cf03fe475096631e5e07bbd7a0fde60c4cf25c7a200d90258a3016450526f540465312e332e35055820acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b01a10281820158200263829989b6fd954f72baaf2fc64bc2e2f01d692d4de72986ea808f6e99813fa200d90258a3016441526f540465302e312e34055820acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b01a1028182015820a3a5e715f0cc574a73c3f9bebb6bc24f32ffd5b67b387244c2c909da779a1478")

	assert.Equal(t, expected, actual)
}

func TestUnsignedCorim_unmarshal(t *testing.T) {
	tv := comid.MustHexDecode(t, "a2006d7465737420636f72696d20696401815901b8d901faa40065656e2d474201a100782434334242453337462d324536312d344233332d414544332d3533434646313432384231360281a3006941434d45204c74642e01d8207468747470733a2f2f61636d652e6578616d706c65028300010204a1008182a100a300d90227582061636d652d696d706c656d656e746174696f6e2d69642d303030303030303031016441434d45026a526f616452756e6e657283a200d90258a30162424c0465322e312e30055820acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b01a102818201582087428fc522803d31065e7bce3cf03fe475096631e5e07bbd7a0fde60c4cf25c7a200d90258a3016450526f540465312e332e35055820acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b01a10281820158200263829989b6fd954f72baaf2fc64bc2e2f01d692d4de72986ea808f6e99813fa200d90258a3016441526f540465302e312e34055820acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b01a1028182015820a3a5e715f0cc574a73c3f9bebb6bc24f32ffd5b67b387244c2c909da779a1478")

	var corim UnsignedCorim

	err := dm.Unmarshal(tv, &corim)
	assert.Nil(t, err)

	expectedID, err := corim.GetIDString()
	assert.Nil(t, err)
	assert.Equal(t, expectedID, "test corim id")

	assert.NotNil(t, corim.Tags)
	assert.Equal(t, 1, len(corim.Tags))

	// unmarshal the embedded Comid
	comidTag := corim.Tags[0]
	var comid comid.Comid
	err = dm.Unmarshal(comidTag, &comid)
	assert.Nil(t, err)
}

func TestID_Unmarshal_unknown(t *testing.T) {
	garbage := comid.MustHexDecode(t, "8363676172636261676165")

	var id ID
	err := id.UnmarshalCBOR(garbage)

	assert.EqualError(t, err, "unknown corim-id type (CBOR: 8363676172636261676165)")
}
