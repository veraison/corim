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

func TestTeeInstanceID_NewTeeInstanceID_OK(t *testing.T) {
	tvs := []struct {
		desc  string
		input any
	}{
		{
			desc:  "integer",
			input: TestUIntInstance,
		},
		{
			desc:  "byte array",
			input: TestByteInstance,
		},
		{
			desc:  "unsigned integer 64",
			input: uint64(TestUIntInstance),
		},
	}

	for _, tv := range tvs {
		_, err := NewTeeInstanceID(tv.input)
		require.Nil(t, err)
	}
}

func TestTeeInstanceID_SetTeeInstanceID_OK(t *testing.T) {
	inst := &TeeInstanceID{}
	err := inst.SetTeeInstanceID(TestUIntInstance)
	require.NoError(t, err)
	err = inst.SetTeeInstanceID(TestByteInstance)
	require.NoError(t, err)
	err = inst.SetTeeInstanceID(uint64(1000))
	require.NoError(t, err)
}

func TestTeeInstanceID_SetTeeInstanceID_NOK(t *testing.T) {
	inst := &TeeInstanceID{}
	expectedErr := "unsupported negative TeeInstanceID: -1"
	err := inst.SetTeeInstanceID(-1)
	assert.EqualError(t, err, expectedErr)
	expectedErr = "unsupported TeeInstanceID type: float64"
	err = inst.SetTeeInstanceID(-1.234)
	assert.EqualError(t, err, expectedErr)
}

func TestTeeInstanceID_Valid_OK(t *testing.T) {
	inst := &TeeInstanceID{TestUIntInstance}
	err := inst.Valid()
	require.NoError(t, err)
}

func TestTeeInstanceID_Valid_NOK(t *testing.T) {
	tvs := []struct {
		desc        string
		input       interface{}
		expectedErr string
	}{
		{
			desc:        "unsupported type negative integer",
			input:       -1,
			expectedErr: "unsupported negative TeeInstanceID: -1",
		},
		{
			desc:        "non existent TeeInstanceID",
			input:       nil,
			expectedErr: "empty TeeInstanceID",
		},
		{
			desc:        "non existent TeeInstanceID",
			input:       []byte{},
			expectedErr: "empty TeeInstanceID",
		},
		{
			desc:        "unsupported type float64",
			input:       1.234,
			expectedErr: "unsupported TeeInstanceID type: float64",
		},
	}

	for _, tv := range tvs {
		inst := &TeeInstanceID{tv.input}
		err := inst.Valid()
		assert.EqualError(t, err, tv.expectedErr)
	}
}

func TestTeeInstanceID_MarshalCBOR_Bytes(t *testing.T) {
	inst, err := NewTeeInstanceID(TestByteInstance)
	require.Nil(t, err)
	expected := comid.MustHexDecode(t, "43454647")
	actual, err := inst.MarshalCBOR()
	fmt.Printf("CBOR: %x\n", actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestTeeInstanceID_JSON(t *testing.T) {
	inst := TeeInstanceID{TestByteInstance}
	ji, err := inst.MarshalJSON()
	assert.Nil(t, err)
	i := &TeeInstanceID{}
	err = i.UnmarshalJSON(ji)
	assert.Nil(t, err)
}
