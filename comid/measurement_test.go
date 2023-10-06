// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/eat"
	"github.com/veraison/swid"
)

func TestMeasurement_NewUUIDMeasurement_good_uuid(t *testing.T) {
	tv := NewUUIDMeasurement(TestUUID)
	assert.NotNil(t, tv)
}

func TestMeasurement_NewUUIDMeasurement_empty_uuid(t *testing.T) {
	emptyUUID := UUID{}

	tv := NewUUIDMeasurement(emptyUUID)

	assert.Nil(t, tv)
}

func TestMeasurement_NewUIntMeasurement(t *testing.T) {
	var TestUint uint64 = 35

	tv := NewUintMeasurement(TestUint)

	assert.NotNil(t, tv)
}

func TestMeasurement_NewPSAMeasurement_empty(t *testing.T) {
	emptyPSARefValID := PSARefValID{}

	tv := NewPSAMeasurement(emptyPSARefValID)
	assert.Nil(t, tv)
}

func TestMeasurement_NewPSAMeasurement_no_values(t *testing.T) {
	psaRefValID :=
		NewPSARefValID(TestSignerID).
			SetLabel("PRoT").
			SetVersion("1.2.3")
	require.NotNil(t, psaRefValID)

	tv := NewPSAMeasurement(*psaRefValID)
	assert.NotNil(t, tv)

	err := tv.Valid()
	assert.EqualError(t, err, "no measurement value set")
}

func TestMeasurement_NewCCAPlatCfgMeasurement_no_values(t *testing.T) {
	ccaplatID := CCAPlatformConfigID(TestCCALabel)

	tv := NewCCAPlatCfgMeasurement(ccaplatID)
	assert.NotNil(t, tv)

	err := tv.Valid()
	assert.EqualError(t, err, "no measurement value set")
}

func TestMeasurement_NewCCAPlatCfgMeasurement_valid_meas(t *testing.T) {
	ccaplatID := CCAPlatformConfigID(TestCCALabel)

	tv := NewCCAPlatCfgMeasurement(ccaplatID).SetRawValueBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{})
	assert.NotNil(t, tv)

	err := tv.Valid()
	assert.Nil(t, err)
}

func TestMeasurement_NewPSAMeasurement_one_value(t *testing.T) {
	psaRefValID :=
		NewPSARefValID(TestSignerID).
			SetLabel("PRoT").
			SetVersion("1.2.3")
	require.NotNil(t, psaRefValID)

	tv := NewPSAMeasurement(*psaRefValID).SetIPaddr(TestIPaddr)
	assert.NotNil(t, tv)

	err := tv.Valid()
	assert.Nil(t, err)
}

func TestMeasurement_NewUUIDMeasurement_no_values(t *testing.T) {
	tv := NewUUIDMeasurement(TestUUID)
	require.NotNil(t, tv)

	err := tv.Valid()
	assert.EqualError(t, err, "no measurement value set")
}

func TestMeasurement_NewUUIDMeasurement_some_value(t *testing.T) {
	var vs swid.VersionScheme
	require.NoError(t, vs.SetCode(swid.VersionSchemeSemVer))

	tv := NewUUIDMeasurement(TestUUID).
		SetMinSVN(2).
		SetFlagsTrue(FlagIsDebug).
		SetVersion("1.2.3", swid.VersionSchemeSemVer)
	require.NotNil(t, tv)

	err := tv.Valid()
	assert.Nil(t, err)
}

func TestMeasurement_NewUUIDMeasurement_bad_digest(t *testing.T) {
	tv := NewUUIDMeasurement(TestUUID)
	require.NotNil(t, tv)

	assert.Nil(t, tv.AddDigest(swid.Sha256, []byte{0xff}))
}

func TestMeasurement_NewUUIDMeasurement_bad_ueid(t *testing.T) {
	tv := NewUUIDMeasurement(TestUUID)
	require.NotNil(t, tv)

	badUEID := eat.UEID{
		0xFF, // Invalid
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	}

	assert.Nil(t, tv.SetUEID(badUEID))
}

func TestMeasurement_NewUUIDMeasurement_bad_uuid(t *testing.T) {
	tv := NewUUIDMeasurement(TestUUID)
	require.NotNil(t, tv)

	nonRFC4122UUID, err := ParseUUID("f47ac10b-58cc-4372-c567-0e02b2c3d479")
	require.Nil(t, err)

	assert.Nil(t, tv.SetUUID(nonRFC4122UUID))
}

var (
	testMKeyUintMin uint64
	testMKeyUintMax = ^uint64(0)
)

func TestMkey_Valid_no_value(t *testing.T) {
	mkey := &Mkey{}
	expectedErr := "unknown measurement key type: <nil>"
	err := mkey.Valid()
	assert.EqualError(t, err, expectedErr)
}

func TestMKey_MarshalCBOR_CCAPlatformConfigID_ok(t *testing.T) {
	tvs := []struct {
		mkey     CCAPlatformConfigID
		expected []byte
	}{
		{
			mkey:     CCAPlatformConfigID(TestCCALabel),
			expected: MustHexDecode(t, "d9025a736363612d706c6174666f726d2d636f6e666967"),
		},
		{
			mkey:     CCAPlatformConfigID("mytestplatformfig"),
			expected: MustHexDecode(t, "d9025a716d7974657374706c6174666f726d666967"),
		},
		{
			mkey:     CCAPlatformConfigID("mytestlabel2"),
			expected: MustHexDecode(t, "d9025a6c6d79746573746c6162656c32"),
		},
	}

	for _, tv := range tvs {
		mkey := &Mkey{TaggedCCAPlatformConfigID(tv.mkey)}
		actual, err := mkey.MarshalCBOR()
		assert.Nil(t, err)
		assert.Equal(t, tv.expected, actual)
		fmt.Printf("CBOR: %x\n", actual)
	}
}

func TestMKey_UnmarshalCBOR_CCAPlatformConfigID_ok(t *testing.T) {
	tvs := []struct {
		input    []byte
		expected CCAPlatformConfigID
	}{
		{
			input:    MustHexDecode(t, "d9025a736363612d706c6174666f726d2d636f6e666967"),
			expected: CCAPlatformConfigID(TestCCALabel),
		},
		{
			input:    MustHexDecode(t, "d9025a716d7974657374706c6174666f726d666967"),
			expected: CCAPlatformConfigID("mytestplatformfig"),
		},
		{
			input:    MustHexDecode(t, "d9025a6c6d79746573746c6162656c32"),
			expected: CCAPlatformConfigID("mytestlabel2"),
		},
	}

	for _, tv := range tvs {
		mkey := &Mkey{}
		err := mkey.UnmarshalCBOR(tv.input)
		assert.Nil(t, err)
		actual, err := mkey.GetCCAPlatformConfigID()
		assert.Nil(t, err)
		assert.Equal(t, tv.expected, actual)
		fmt.Printf("CBOR: %x\n", actual)
	}
}

func TestMKey_MarshalCBOR_uint_ok(t *testing.T) {
	tvs := []struct {
		mkey     uint64
		expected []byte
	}{
		{
			mkey:     testMKeyUintMin,
			expected: MustHexDecode(t, "00"),
		},
		{
			mkey:     TestMKey,
			expected: MustHexDecode(t, "1902BC"),
		},
		{
			mkey:     testMKeyUintMax,
			expected: MustHexDecode(t, "1BFFFFFFFFFFFFFFFF"),
		},
	}

	for _, tv := range tvs {
		mkey := &Mkey{tv.mkey}
		actual, err := mkey.MarshalCBOR()
		assert.Nil(t, err)
		assert.Equal(t, tv.expected, actual)
		fmt.Printf("CBOR: %x\n", actual)
	}
}
func TestMkey_MarshalCBOR_uint_not_ok(t *testing.T) {
	tvs := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    123.456,
			expected: "unknown measurement key type: float64",
		},
		{
			input:    "sample",
			expected: "unknown measurement key type: string",
		},
	}

	for _, tv := range tvs {
		mkey := &Mkey{tv.input}
		_, err := mkey.MarshalCBOR()
		assert.Nil(t, err)
		err = mkey.Valid()
		assert.EqualError(t, err, tv.expected)
	}
}

func TestMkey_UnmarshalCBOR_uint_ok(t *testing.T) {
	tvs := []struct {
		mkey     []byte
		expected uint64
	}{
		{
			mkey:     MustHexDecode(t, "00"),
			expected: testMKeyUintMin,
		},
		{
			mkey:     MustHexDecode(t, "1902BC"),
			expected: TestMKey,
		},
		{
			mkey:     MustHexDecode(t, "1BFFFFFFFFFFFFFFFF"),
			expected: testMKeyUintMax,
		},
	}

	for _, tv := range tvs {
		mKey := &Mkey{}

		err := mKey.UnmarshalCBOR(tv.mkey)
		assert.Nil(t, err)
		actual, err := mKey.GetKeyUint()
		assert.Nil(t, err)
		assert.Equal(t, tv.expected, actual)
	}
}

func TestMkey_UnmarshalCBOR_not_ok(t *testing.T) {
	tvs := []struct {
		input    []byte
		expected string
	}{
		{
			input:    []byte{0xAB, 0xCD},
			expected: "unexpected EOF",
		},
		{
			input:    []byte{0xCC, 0xDD, 0xFF},
			expected: "cbor: invalid additional information 29 for type tag",
		},
	}

	for _, tv := range tvs {
		mKey := &Mkey{}

		err := mKey.UnmarshalCBOR(tv.input)

		assert.EqualError(t, err, tv.expected)
	}
}

func TestMkey_UnmarshalCBOR_uint_not_ok(t *testing.T) {
	tvs := []struct {
		input    []byte
		expected string
	}{
		{
			input: []byte{0xd8, 0x25, 0x50, 0x31, 0xfb, 0x5a, 0xbf, 0x02,
				0x3e, 0x49, 0x92, 0xaa, 0x4e, 0x95, 0xf9, 0xc1,
				0x50, 0x3b, 0xfa},
			expected: "measurement-key type is: comid.TaggedUUID",
		},
		{
			input: []byte{0xd8, 0x21, 0x50, 0x31, 0xfb, 0x5a, 0xff, 0x12,
				0xFF, 0xFF, 0x92, 0xaa, 0x4e, 0x95, 0xf9, 0xc1,
				0x50, 0x3b, 0xfa},
			expected: "measurement-key type is: cbor.Tag",
		},
	}

	for _, tv := range tvs {
		mKey := &Mkey{}

		err := mKey.UnmarshalCBOR(tv.input)
		assert.Nil(t, err)
		_, err = mKey.GetKeyUint()
		assert.EqualError(t, err, tv.expected)
	}
}

func TestMKey_MarshalJSON_CCAPlatformConfigID_ok(t *testing.T) {
	refval := TestCCALabel
	mkey := &Mkey{val: TaggedCCAPlatformConfigID(refval)}

	expected := `{"type":"cca.platform-config-id","value":"cca-platform-config"}`

	actual, err := mkey.MarshalJSON()
	assert.Nil(t, err)

	assert.JSONEq(t, expected, string(actual))
	fmt.Printf("JSON: %x\n", actual)
}

func TestMKey_UnMarshalJSON_CCAPlatformConfigID_ok(t *testing.T) {
	input := []byte(`{"type":"cca.platform-config-id","value":"cca-platform-config"}`)
	expected := CCAPlatformConfigID(TestCCALabel)

	mKey := &Mkey{}

	err := mKey.UnmarshalJSON(input)
	assert.Nil(t, err)
	actual, err := mKey.GetCCAPlatformConfigID()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

}

func TestMKey_UnMarshalJSON_CCAPlatformConfigID_not_ok(t *testing.T) {
	input := []byte(`{"type":"cca.platform-config-id","value":""}`)
	expected := "cannot unmarshal $measured-element-type-choice of type CCAPlatformConfigID: empty label"

	mKey := &Mkey{}

	err := mKey.UnmarshalJSON(input)

	assert.NotNil(t, err)
	assert.Equal(t, expected, err.Error())

}
func TestMkey_MarshalJSON_uint_ok(t *testing.T) {
	tvs := []struct {
		mkey     uint64
		expected string
	}{
		{
			mkey:     testMKeyUintMin,
			expected: `{"type":"uint","value":0}`,
		},
		{
			mkey:     TestMKey,
			expected: `{"type":"uint","value":700}`,
		},
		{
			mkey:     testMKeyUintMax,
			expected: `{"type":"uint","value":18446744073709551615}`,
		},
	}

	for _, tv := range tvs {

		mkey := &Mkey{tv.mkey}

		actual, err := mkey.MarshalJSON()
		assert.Nil(t, err)

		assert.JSONEq(t, tv.expected, string(actual))
		fmt.Printf("JSON: %x\n", actual)
	}
}

func TestMkey_MarshalJSON_uint_not_ok(t *testing.T) {
	tvs := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    123.456,
			expected: "unknown type float64 for mkey",
		},
		{
			input:    "sample",
			expected: "unknown type string for mkey",
		},
	}

	for _, tv := range tvs {

		mkey := &Mkey{tv.input}

		_, err := mkey.MarshalJSON()

		assert.EqualError(t, err, tv.expected)
	}
}

func TestMkey_UnmarshalJSON_uint_ok(t *testing.T) {
	tvs := []struct {
		input    []byte
		expected uint64
	}{
		{
			input:    []byte(`{"type":"uint","value":0}`),
			expected: testMKeyUintMin,
		},
		{
			input:    []byte(`{"type":"uint","value":700}`),
			expected: TestMKey,
		},
		{
			input:    []byte(`{"type":"uint","value":18446744073709551615}`),
			expected: testMKeyUintMax,
		},
	}

	for _, tv := range tvs {
		mKey := &Mkey{}

		err := mKey.UnmarshalJSON(tv.input)
		assert.Nil(t, err)
		actual, err := mKey.GetKeyUint()
		assert.Nil(t, err)
		assert.Equal(t, tv.expected, actual)
	}
}

func TestMkey_UnmarshalJSON_notok(t *testing.T) {
	tvs := []struct {
		input    []byte
		expected string
	}{
		{
			input:    []byte(`{"type":"uint","value":"abcdefg"}`),
			expected: "cannot unmarshal $measured-element-type-choice of type uint: json: cannot unmarshal string into Go value of type uint64",
		},
		{
			input:    []byte(`{"type":"uint","value":123.456}`),
			expected: "cannot unmarshal $measured-element-type-choice of type uint: json: cannot unmarshal number 123.456 into Go value of type uint64",
		},
	}

	for _, tv := range tvs {
		mKey := &Mkey{}

		err := mKey.UnmarshalJSON(tv.input)

		assert.EqualError(t, err, tv.expected)
	}
}

func TestMkey_UnmarshalJSON_uint_notok(t *testing.T) {
	tvs := []struct {
		input    []byte
		expected string
	}{
		{
			input:    []byte(`{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}`),
			expected: "measurement-key type is: comid.TaggedUUID",
		},
		{
			input:    []byte(`{"type":"psa.refval-id","value":{"label": "BL","version": "2.1.0","signer-id": "rLsRx+TaIXIFUjzkzhokWuGiOa48a/2eeHH35di66Gs="}}`),
			expected: "measurement-key type is: comid.TaggedPSARefValID",
		},
	}

	for _, tv := range tvs {
		mKey := &Mkey{}

		err := mKey.UnmarshalJSON(tv.input)
		assert.Nil(t, err)
		_, err = mKey.GetKeyUint()
		assert.EqualError(t, err, tv.expected)
	}
}
