// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package encoding

import (
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PopulateStructFromCBOR_simple(t *testing.T) {
	type SimpleStruct struct {
		FieldOne string `cbor:"0,keyasint,omitempty"`
		FieldTwo int    `cbor:"1,keyasint"`
	}

	var v SimpleStruct

	data := []byte{
		0xa2, // map(2)

		0x00,                   // key 0
		0x64,                   // val tstr(4)
		0x61, 0x63, 0x6d, 0x65, // "acme"

		0x01, // key 1
		0x06, // val 6
	}

	dm, err := cbor.DecOptions{}.DecMode()
	require.NoError(t, err)

	err = PopulateStructFromCBOR(dm, data, &v)
	require.NoError(t, err)
	assert.Equal(t, "acme", v.FieldOne)
	assert.Equal(t, 6, v.FieldTwo)

	data = []byte{
		0xa1, // map(1)

		0x01, // key 1
		0x06, // val 6
	}
	v = SimpleStruct{}

	err = PopulateStructFromCBOR(dm, data, &v)
	require.NoError(t, err)
	assert.Equal(t, "", v.FieldOne)
	assert.Equal(t, 6, v.FieldTwo)

	data = []byte{
		0xa1, // map(1)

		0x02,                   // key 2
		0x64,                   // val tstr(4)
		0x61, 0x63, 0x6d, 0x65, // "acme"
	}
	v = SimpleStruct{}

	err = PopulateStructFromCBOR(dm, data, &v)
	assert.EqualError(t, err, `missing mandatory field "FieldTwo" (1)`)

	err = PopulateStructFromCBOR(dm, []byte{0x01}, &v)
	assert.EqualError(t, err, `cbor: cannot unmarshal positive integer into Go value of type map[int]cbor.RawMessage`)

	type CompositeStruct struct {
		FieldThree string `cbor:"2,keyasint"`
		SimpleStruct
	}

	var c CompositeStruct

	data = []byte{
		0xa3, // map(3)

		0x00,                   // key 0
		0x64,                   // val tstr(4)
		0x61, 0x63, 0x6d, 0x65, // "acme"

		0x01, // key 1
		0x06, // val 6

		0x02,             // key 2
		0x63,             // val tstr(3)
		0x66, 0x6f, 0x6f, // "foo"
	}

	err = PopulateStructFromCBOR(dm, data, &c)
	require.NoError(t, err)
	assert.Equal(t, "acme", c.FieldOne)
	assert.Equal(t, 6, c.FieldTwo)
	assert.Equal(t, "foo", c.FieldThree)

	em, err := cbor.EncOptions{}.EncMode()
	require.NoError(t, err)

	res, err := SerializeStructToCBOR(em, &c)
	require.NoError(t, err)

	var c2 CompositeStruct
	err = PopulateStructFromCBOR(dm, res, &c2)
	require.NoError(t, err)
	assert.EqualValues(t, c, c2)

}
