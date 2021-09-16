// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRel_NewRel_default_value(t *testing.T) {
	tv := NewRel()
	require.NotNil(t, tv)

	// default is unset
	assert.Equal(t, tv.Get(), RelUnset)
	assert.EqualError(t, tv.Valid(), "rel is unset")
}

func TestRel_NewRel_set_and_reset(t *testing.T) {
	// initialize and set
	tv := NewRel().Set(RelReplaces)
	require.NotNil(t, tv)

	assert.Nil(t, tv.Valid())
	assert.Equal(t, RelReplaces, tv.Get())

	// reset
	tv = tv.Set(RelSupplements)

	assert.Nil(t, tv.Valid())
	assert.Equal(t, RelSupplements, tv.Get())
}

func TestRel_UnmarshalJSON_ok(t *testing.T) {
	tvs := []struct {
		json     string
		expected Rel
	}{
		{
			json:     `"supplements"`,
			expected: RelSupplements,
		},
		{
			json:     `"replaces"`,
			expected: RelReplaces,
		},
	}

	for _, tv := range tvs {
		var actual Rel

		err := actual.UnmarshalJSON([]byte(tv.json))

		assert.Nil(t, err)
		assert.Equal(t, tv.expected, actual.Get())
	}
}

func TestRel_UnmarshalJSON_fail(t *testing.T) {
	tvs := []struct {
		json        string
		expectedErr string
	}{
		{
			json:        `""`,
			expectedErr: "empty rel",
		},
		{
			json:        `"blabla"`,
			expectedErr: "unknown rel 'blabla'",
		},
		{
			json:        `"unterminated strin`,
			expectedErr: "cannot unmarshal rel: unexpected end of JSON input",
		},
		{
			json:        `0`,
			expectedErr: "cannot unmarshal rel: json: cannot unmarshal number into Go value of type string",
		},
	}

	for _, tv := range tvs {
		var actual Rel

		err := actual.UnmarshalJSON([]byte(tv.json))

		assert.EqualError(t, err, tv.expectedErr)
	}
}

func TestRel_FromCBOR_ok(t *testing.T) {
	tvs := []struct {
		cbor     []byte
		expected Rel
	}{
		{
			// 00 => unsigned(0)
			cbor:     MustHexDecode(t, "00"),
			expected: RelSupplements,
		},
		{
			// 01 => unsigned(1)
			cbor:     MustHexDecode(t, "01"),
			expected: RelReplaces,
		},
	}

	for _, tv := range tvs {
		var actual Rel

		err := actual.FromCBOR(tv.cbor)

		assert.Nil(t, err)
		assert.Equal(t, tv.expected, actual.Get())
	}
}

func TestRel_ToCBOR_ok(t *testing.T) {
	tvs := []struct {
		rel      Rel
		expected []byte
	}{
		{
			rel:      RelSupplements,
			expected: MustHexDecode(t, "00"),
		},
		{
			rel:      RelReplaces,
			expected: MustHexDecode(t, "01"),
		},
		{
			// it is possible to force non-registered values
			rel:      Rel(6),
			expected: MustHexDecode(t, "06"),
		},
	}

	r := NewRel()
	require.NotNil(t, r)

	for _, tv := range tvs {
		require.NotNil(t, r.Set(tv.rel))

		actual, err := r.ToCBOR()

		fmt.Printf("CBOR: %x\n", actual)

		assert.Nil(t, err)
		assert.Equal(t, tv.expected, actual)

	}
}

func TestRel_ToCBOR_fail_unset(t *testing.T) {
	r := NewRel()
	require.NotNil(t, r)

	_, err := r.ToCBOR()

	assert.EqualError(t, err, "rel is unset")
}
