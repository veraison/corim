// Copyright 2024-2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRawValue_NewRawValue_ok(t *testing.T) {
	tv := NewRawValue()
	require.NotNil(t, tv)
}

func TestRawValue_Set_Get_Bytes_ok(t *testing.T) {
	tv := RawValue{}
	expected := []byte{0x01, 0x02, 0x03}
	rv := tv.SetBytes([]byte{0x01, 0x02, 0x03})
	require.NotNil(t, rv)
	rval, err := rv.GetBytes()
	assert.NoError(t, err)
	assert.Equal(t, expected, rval)
}

func TestRawValue_Get_Bytes_nok(t *testing.T) {
	rv := RawValue{}
	expectedErr := "raw value is not set"
	_, err := rv.GetBytes()
	assert.EqualError(t, err, expectedErr)
	rv = RawValue{"testraw"}
	expectedErr = "unknown type string for $raw-value-type-choice"
	_, err = rv.GetBytes()
	assert.EqualError(t, err, expectedErr)
}

func TestRawValue_Marshal_UnMarshal_JSON_ok(t *testing.T) {
	tv := RawValue{}
	rv := tv.SetBytes([]byte{0x01, 0x02, 0x03})
	bytes, err := rv.MarshalJSON()
	assert.NoError(t, err)
	sv := RawValue{}
	err = sv.UnmarshalJSON(bytes)
	assert.NoError(t, err)
	assert.Equal(t, *rv, sv)
}

func TestRawValue_Marshal_UnMarshal_CBOR_ok(t *testing.T) {
	tv := RawValue{}
	rv := tv.SetBytes([]byte{0x01, 0x02, 0x03})
	bytes, err := rv.MarshalCBOR()
	assert.NoError(t, err)
	sv := RawValue{}
	err = sv.UnmarshalCBOR(bytes)
	assert.NoError(t, err)
	assert.Equal(t, *rv, sv)
}

func TestRawValue_Equal_True(t *testing.T) {
	claim := RawValue{}
	claim.SetBytes([]byte{0x01, 0x02, 0x03})
	ref := RawValue{}
	ref.SetBytes([]byte{0x01, 0x02, 0x03})

	assert.True(t, claim.Equal(ref))
}

func TestRawValue_Equal_False(t *testing.T) {
	claim := RawValue{}
	claim.SetBytes([]byte{0x01, 0x02, 0x03})
	ref := RawValue{}
	ref.SetBytes([]byte{0x01, 0x02, 0x04})

	assert.False(t, claim.Equal(ref))
}

func TestRawValue_Compare_True(t *testing.T) {
	claim := RawValue{}
	claim.SetBytes([]byte{0x01, 0x02, 0x03})
	ref := []byte{0x01, 0x00, 0x03}
	mask := []byte{0xff, 0x00, 0xff}

	assert.True(t, claim.CompareAgainstReference(ref, &mask))
}

func TestRawValue_Compare_False(t *testing.T) {
	claim := RawValue{}
	claim.SetBytes([]byte{0x01, 0x02, 0x03})
	ref := []byte{0x04, 0x05, 0x06}

	assert.False(t, claim.CompareAgainstReference(ref, nil))
}
