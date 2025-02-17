// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
)

func TestISVProdID_NewISVProdID_OK(t *testing.T) {
	tvs := []struct {
		desc  string
		input interface{}
	}{
		{
			desc:  "integer",
			input: TestUIntISVProdID,
		},
		{
			desc:  "byte array",
			input: TestBytesISVProdID,
		},
		{
			desc:  "unsigned integer 64",
			input: uint64(TestUIntISVProdID),
		},
	}

	for _, tv := range tvs {
		_, err := NewTeeISVProdID(tv.input)
		require.Nil(t, err)
	}
}

func TestIsvProdID_SetTeeISVProdID_OK(t *testing.T) {
	id := &TeeISVProdID{}
	err := id.SetTeeISVProdID(TestUIntISVProdID)
	require.NoError(t, err)
	err = id.SetTeeISVProdID(TestBytesISVProdID)
	require.NoError(t, err)
	err = id.SetTeeISVProdID(uint64(1000))
	require.NoError(t, err)
}

func TestIsvProdID_SetTeeISVProdID_NOK(t *testing.T) {
	id := &TeeISVProdID{}
	expectedErr := "unsupported negative TeeISVProdID: -1"
	err := id.SetTeeISVProdID(-1)
	assert.EqualError(t, err, expectedErr)
	expectedErr = "unsupported TeeISVProdID type: float64"
	err = id.SetTeeISVProdID(-1.234)
	assert.EqualError(t, err, expectedErr)
}

func TestIsvProdID_Valid_OK(t *testing.T) {
	id := &TeeISVProdID{TestUIntISVProdID}
	err := id.Valid()
	require.NoError(t, err)
	id = &TeeISVProdID{TestBytesISVProdID}
	err = id.Valid()
	require.NoError(t, err)
}

func TestIsvProdID_Valid_NOK(t *testing.T) {
	tvs := []struct {
		desc        string
		input       interface{}
		expectedErr string
	}{
		{
			desc:        "unsupported type negative integer",
			input:       -1,
			expectedErr: "unsupported negative TeeISVProdID: -1",
		},
		{
			desc:        "non existent TeeISVProdID",
			input:       nil,
			expectedErr: "empty TeeISVProdID",
		},
		{
			desc:        "non existent TeeISVProdID",
			input:       []byte{},
			expectedErr: "empty TeeISVProdID",
		},
		{
			desc:        "unsupported type float64",
			input:       1.234,
			expectedErr: "unsupported TeeISVProdID type: float64",
		},
	}

	for _, tv := range tvs {
		id := &TeeISVProdID{tv.input}
		err := id.Valid()
		assert.EqualError(t, err, tv.expectedErr)
	}
}

func TestIsvProdID_MarshalCBOR_Bytes(t *testing.T) {
	id, err := NewTeeISVProdID(TestBytesISVProdID)
	require.Nil(t, err)
	expected := comid.MustHexDecode(t, "43010203")
	actual, err := id.MarshalCBOR()
	fmt.Printf("CBOR: %x\n", actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestIsvProdID_JSON(t *testing.T) {
	isv := TeeISVProdID{TestBytesISVProdID}
	jisv, err := isv.MarshalJSON()
	assert.Nil(t, err)
	i := &TeeISVProdID{}
	err = i.UnmarshalJSON(jisv)
	assert.Nil(t, err)
}
