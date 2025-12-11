// Copyright 2021-2025 Contributors to the Veraison project.
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
		FieldOne     string `cbor:"0,keyasint,omitempty"`
		FieldTwo     int    `cbor:"1,keyasint"`
		IgnoredField int    `cbor:"-"`
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
	assert.EqualError(t, err, `expected map (CBOR Major Type 5), found Major Type 0`)

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

func Test_SerializeStructToCBOR_cde_ordering(t *testing.T) {
	val := struct {
		Field8  int `cbor:"8,keyasint"`
		FieldN3 int `cbor:"-3,keyasint"`
		Field0  int `cbor:"0,keyasint"`
		FieldN1 int `cbor:"-1,keyasint"`
	}{
		Field8:  1,
		FieldN3: 3,
		Field0:  0,
		FieldN1: 2,
	}

	expected := []byte{
		0xa4, // map(4)

		0x00, // key: 0
		0x00, // value: 0

		0x08, // key: 8
		0x01, // value: 1

		0x20, // key: -1
		0x02, // value: 2

		0x22, // key: -3
		0x03, // value: 3
	}

	em := mustInitEncMode()
	data, err := SerializeStructToCBOR(em, val)
	assert.NoError(t, err)
	assert.Equal(t, expected, data)
}

func Test_structFieldsCBOR_CRUD(t *testing.T) {
	sf := newStructFieldsCBOR()

	err := sf.Add(2, cbor.RawMessage{0x02})
	assert.NoError(t, err)

	err = sf.Add(1, cbor.RawMessage{0x01})
	assert.NoError(t, err)

	err = sf.Add(3, cbor.RawMessage{0x03})
	assert.NoError(t, err)

	assert.Equal(t, []int{2, 1, 3}, sf.Keys)
	assert.True(t, sf.Has(3))
	assert.False(t, sf.Has(4))

	val, ok := sf.Get(2)
	assert.True(t, ok)
	assert.Equal(t, cbor.RawMessage{0x2}, val)

	_, ok = sf.Get(4)
	assert.False(t, ok)

	sf.Delete(2)
	_, ok = sf.Get(2)
	assert.False(t, ok)

	err = sf.Add(1, cbor.RawMessage{0x11})
	assert.EqualError(t, err, "duplicate cbor key: 1")
}

func Test_structFieldsCBOR_CBOR_roundtrip(t *testing.T) {
	em, err := cbor.EncOptions{}.EncMode()
	require.NoError(t, err)
	dm, err := cbor.DecOptions{}.DecMode()
	require.NoError(t, err)

	sf := newStructFieldsCBOR()

	data, err := sf.ToCBOR(em)
	require.NoError(t, err)
	assert.Equal(t, data, []byte{0xa0}) // empty map

	for i := 0; i < 5; i++ {
		err = sf.Add(i, cbor.RawMessage{0x00})
		require.NoError(t, err)
	}

	data, err = sf.ToCBOR(em)
	require.NoError(t, err)
	assert.Equal(t, data, []byte{
		0xa5, // map 5
		0x00, 0x00,
		0x01, 0x00,
		0x02, 0x00,
		0x03, 0x00,
		0x04, 0x00,
	})

	sfOut := newStructFieldsCBOR()
	err = sfOut.FromCBOR(dm, data)
	require.NoError(t, err)
	assert.Equal(t, sf, sfOut)

	for i := 5; i < 200; i++ {
		err = sf.Add(i, cbor.RawMessage{0x00})
		require.NoError(t, err)
	}

	data, err = sf.ToCBOR(em)
	require.NoError(t, err)
	assert.Equal(t, data[:2], []byte{
		0xb8, 0xc8, // map 200
	})

	sfOut = newStructFieldsCBOR()
	err = sfOut.FromCBOR(dm, data)
	require.NoError(t, err)
	assert.Equal(t, sf, sfOut)

	for i := 200; i < 2048; i++ {
		err = sf.Add(i, cbor.RawMessage{0x00})
		require.NoError(t, err)
	}

	data, err = sf.ToCBOR(em)
	require.NoError(t, err)
	assert.Equal(t, data[:3], []byte{
		0xb9, 0x08, 0x00, // map 2048
	})

	sfOut = newStructFieldsCBOR()
	err = sfOut.FromCBOR(dm, data)
	require.NoError(t, err)
	assert.Equal(t, sf, sfOut)
}

func Test_structFieldsCBOR_CBOR_decode_tagged(t *testing.T) {
	data := []byte{
		0xc1, // tag 1
		0xa5, // map 5
		0x00, 0x00,
		0x01, 0x00,
		0x02, 0x00,
		0x03, 0x00,
		0x04, 0x00,
	}

	dm, err := cbor.DecOptions{}.DecMode()
	require.NoError(t, err)

	sfOut := newStructFieldsCBOR()
	err = sfOut.FromCBOR(dm, data)
	require.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, sfOut.Keys)
}

func Test_structFieldsCBOR_CBOR_decode_indefinite(t *testing.T) {
	data := []byte{
		0xbf, // indefinite map
		0x00, 0x00,
		0x01, 0x00,
		0x02, 0x00,
		0x03, 0x00,
		0x04, 0x00,
		0xff, // break
	}

	dm, err := cbor.DecOptions{}.DecMode()
	require.NoError(t, err)

	sfOut := newStructFieldsCBOR()
	err = sfOut.FromCBOR(dm, data)
	require.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, sfOut.Keys)
}

func Test_structFieldsCBOR_CBOR_decode_negative(t *testing.T) {
	dm, err := cbor.DecOptions{}.DecMode()
	require.NoError(t, err)

	sfOut := newStructFieldsCBOR()
	err = sfOut.FromCBOR(dm, []byte{0xa1, 0xff, 0x00})
	assert.EqualError(t, err, `map item 0: could not unmarshal key: cbor: unexpected "break" code`)
	err = sfOut.FromCBOR(dm, []byte{0xbf, 0x00, 0x00})
	assert.EqualError(t, err, `unexpected EOF`)
	err = sfOut.FromCBOR(dm, []byte{0xa1, 0x00, 0xff})
	assert.EqualError(t, err, `map item 0: could not unmarshal value: cbor: unexpected "break" code`)

	err = sfOut.FromCBOR(dm, []byte{0x00})
	assert.EqualError(t, err, `expected map (CBOR Major Type 5), found Major Type 0`)
}

func Test_processAdditionalInfo(t *testing.T) {
	addInfo := byte(26)
	data := []byte{0x00, 0x00, 0x00, 0x01}

	val, rest, err := processAdditionalInfo(addInfo, data)
	require.NoError(t, err)
	assert.Equal(t, 1, val)
	assert.Equal(t, []byte{}, rest)

	_, _, err = processAdditionalInfo(byte(27), data)
	assert.EqualError(t, err, "cbor: cannot decode length value of 8 bytes")

	_, _, err = processAdditionalInfo(byte(28), data)
	assert.EqualError(t, err, "cbor: unexpected additional information value 28")

	_, _, err = processAdditionalInfo(addInfo, []byte{})
	assert.EqualError(t, err, "unexpected EOF")
}

func Test_lexSort(t *testing.T) {
	test_cases := []struct {
		title    string
		input    []int
		expected []int
	}{
		{
			title:    "non-negative",
			input:    []int{1, 4, 0, 2, 3},
			expected: []int{0, 1, 2, 3, 4},
		},
		{
			title:    "negative",
			input:    []int{-1, -4, -2, -3},
			expected: []int{-1, -2, -3, -4},
		},
		{
			title:    "mixed",
			input:    []int{-1, 0, 3, 1, -2},
			expected: []int{0, 1, 3, -1, -2},
		},
		{
			title:    "already sorted",
			input:    []int{0, 1, 3, -1, -2},
			expected: []int{0, 1, 3, -1, -2},
		},
		{
			title:    "different length encoding",
			input:    []int{65535, 256},
			expected: []int{256, 65535},
		},
	}

	for _, tc := range test_cases {
		em := mustInitEncMode()
		t.Run(tc.title, func(t *testing.T) {
			lexSort(em, tc.input)
			assert.Equal(t, tc.expected, tc.input)
		})
	}
}
