// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"crypto"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/eat"
	"github.com/veraison/swid"
)

func TestMeasurement_NewUUIDMeasurement_good_uuid(t *testing.T) {
	_, err := NewUUIDMeasurement(TestUUID)
	assert.NoError(t, err)
}

func TestMeasurement_NewUUIDMeasurement_empty_uuid(t *testing.T) {
	emptyUUID := UUID{}

	_, err := NewUUIDMeasurement(emptyUUID)

	assert.EqualError(t, err,
		"invalid key: expecting RFC4122 UUID, got Reserved instead")
}

func TestMeasurement_NewUIntMeasurement(t *testing.T) {
	var TestUint uint64 = 35

	_, err := NewUintMeasurement(TestUint)

	assert.NoError(t, err)
}

func TestMeasurement_NewPSAMeasurement_empty(t *testing.T) {
	emptyPSARefValID := PSARefValID{}

	_, err := NewPSAMeasurement(emptyPSARefValID)

	assert.EqualError(t, err, "invalid key: invalid psa.refval-id: missing mandatory signer ID")
}

func TestMeasurement_NewPSAMeasurement_no_values(t *testing.T) {
	psaRefValID, err := NewPSARefValID(TestSignerID)
	require.NoError(t, err)
	psaRefValID.SetLabel("PRoT")
	psaRefValID.SetVersion("1.2.3")
	require.NotNil(t, psaRefValID)

	tv, err := NewPSAMeasurement(*psaRefValID)
	assert.NoError(t, err)

	err = tv.Valid()
	assert.EqualError(t, err, "no measurement value set")
}

func TestGetPSARefValID(t *testing.T) {
	psaRefValID, err := NewPSARefValID(TestSignerID)
	require.NoError(t, err)
	psaRefValID.SetLabel("PRoT")
	psaRefValID.SetVersion("1.2.3")
	mkey, err := NewMkeyPSARefvalID(psaRefValID)
	require.NoError(t, err)
	actual, err := mkey.GetPSARefValID()
	require.NoError(t, err)
	assert.Equal(t, *psaRefValID, actual)
}

func TestGetPSARefValID_NOK(t *testing.T) {
	mkey := &Mkey{}
	expected := "MKey is not set"
	_, err := mkey.GetPSARefValID()
	assert.EqualError(t, err, expected)
}

func TestGetPSARefValID_InvalidType(t *testing.T) {
	expected := "measurement-key type is: *comid.TaggedCCAPlatformConfigID"
	mkey, err := NewMkeyCCAPlatformConfigID(TestCCALabel)
	require.NoError(t, err)
	_, err = mkey.GetPSARefValID()
	assert.EqualError(t, err, expected)
}

func TestMeasurement_NewCCAPlatCfgMeasurement_no_values(t *testing.T) {
	ccaplatID := CCAPlatformConfigID(TestCCALabel)

	tv, err := NewCCAPlatCfgMeasurement(ccaplatID)
	assert.NoError(t, err)

	err = tv.Valid()
	assert.EqualError(t, err, "no measurement value set")
}

func TestGetCCAPlatformConfigID(t *testing.T) {
	ccaplatID := CCAPlatformConfigID(TestCCALabel)
	mkey, err := NewMkeyCCAPlatformConfigID(TestCCALabel)
	require.NoError(t, err)
	actual, err := mkey.GetCCAPlatformConfigID()
	require.NoError(t, err)
	assert.Equal(t, ccaplatID, actual)
}

func TestGetCCAPlatformConfigID_NOK(t *testing.T) {
	mkey := &Mkey{}
	expected := "MKey is not set"
	_, err := mkey.GetCCAPlatformConfigID()
	assert.EqualError(t, err, expected)
}

func TestGetCCAPlatformConfigID_InvalidType(t *testing.T) {
	mkey := &Mkey{UintMkey(10)}
	expected := "measurement-key type is: comid.UintMkey"
	_, err := mkey.GetCCAPlatformConfigID()
	assert.EqualError(t, err, expected)
}

func TestMeasurement_NewCCAPlatCfgMeasurement_valid_meas(t *testing.T) {
	ccaplatID := CCAPlatformConfigID(TestCCALabel)

	tv, err := NewCCAPlatCfgMeasurement(ccaplatID)
	assert.NoError(t, err)

	tv.SetRawValueBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{})

	err = tv.Valid()
	assert.NoError(t, err)
}

func TestMeasurement_NewPSAMeasurement_one_value(t *testing.T) {
	tv, err := NewPSAMeasurement(MustCreatePSARefValID(TestSignerID, "PRoT", "1.2.3"))
	require.NoError(t, err)

	tv.SetIPaddr(TestIPaddr)

	err = tv.Valid()
	assert.Nil(t, err)
}

func TestMeasurement_NewUUIDMeasurement_no_values(t *testing.T) {
	tv, err := NewUUIDMeasurement(TestUUID)
	require.NoError(t, err)

	err = tv.Valid()
	assert.EqualError(t, err, "no measurement value set")
}

func TestMeasurement_NewUUIDMeasurement_some_value(t *testing.T) {
	var vs swid.VersionScheme
	require.NoError(t, vs.SetCode(swid.VersionSchemeSemVer))

	tv, err := NewUUIDMeasurement(TestUUID)
	require.NoError(t, err)

	tv.SetMinSVN(2).
		SetFlagsTrue(FlagIsDebug).
		SetVersion("1.2.3", swid.VersionSchemeSemVer)

	err = tv.Valid()
	assert.Nil(t, err)
}

func TestMeasurement_NewUUIDMeasurement_bad_digest(t *testing.T) {
	tv, err := NewUUIDMeasurement(TestUUID)
	require.NoError(t, err)

	assert.Nil(t, tv.AddDigest(swid.Sha256, []byte{0xff}))
}

func TestMeasurement_NewUUIDMeasurement_bad_ueid(t *testing.T) {
	tv, err := NewUUIDMeasurement(TestUUID)
	require.NoError(t, err)

	badUEID := eat.UEID{
		0xFF, // Invalid
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	}

	assert.Nil(t, tv.SetUEID(badUEID))
}

func TestMeasurement_NewUUIDMeasurement_bad_uuid(t *testing.T) {
	tv, err := NewUUIDMeasurement(TestUUID)
	require.NoError(t, err)

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
	expectedErr := "Mkey value not set"
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
		expected TaggedCCAPlatformConfigID
	}{
		{
			input:    MustHexDecode(t, "d9025a736363612d706c6174666f726d2d636f6e666967"),
			expected: TaggedCCAPlatformConfigID(TestCCALabel),
		},
		{
			input:    MustHexDecode(t, "d9025a716d7974657374706c6174666f726d666967"),
			expected: TaggedCCAPlatformConfigID("mytestplatformfig"),
		},
		{
			input:    MustHexDecode(t, "d9025a6c6d79746573746c6162656c32"),
			expected: TaggedCCAPlatformConfigID("mytestlabel2"),
		},
	}

	for _, tv := range tvs {
		mkey := &Mkey{}
		err := mkey.UnmarshalCBOR(tv.input)
		assert.Nil(t, err)
		actual, ok := mkey.Value.(*TaggedCCAPlatformConfigID)
		assert.True(t, ok)
		assert.Equal(t, tv.expected, *actual)
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
		mkey := &Mkey{UintMkey(tv.mkey)}
		actual, err := mkey.MarshalCBOR()
		assert.Nil(t, err)
		assert.Equal(t, tv.expected, actual)
		fmt.Printf("CBOR: %x\n", actual)
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
		require.NoError(t, err)
		actual, ok := mKey.Value.(*UintMkey)
		require.True(t, ok)
		assert.Equal(t, tv.expected, uint64(*actual))
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

func TestMKey_MarshalJSON_CCAPlatformConfigID_ok(t *testing.T) {
	refval := TestCCALabel
	mkey := &Mkey{Value: TaggedCCAPlatformConfigID(refval)}

	expected := `{"type":"cca.platform-config-id","value":"cca-platform-config"}`

	actual, err := mkey.MarshalJSON()
	assert.Nil(t, err)

	assert.JSONEq(t, expected, string(actual))
	fmt.Printf("JSON: %x\n", actual)
}

func TestMKey_UnMarshalJSON_CCAPlatformConfigID_ok(t *testing.T) {
	input := []byte(`{"type":"cca.platform-config-id","value":"cca-platform-config"}`)
	expected := TaggedCCAPlatformConfigID(TestCCALabel)

	mKey := &Mkey{}

	err := mKey.UnmarshalJSON(input)
	assert.Nil(t, err)
	actual, ok := mKey.Value.(*TaggedCCAPlatformConfigID)
	assert.True(t, ok)
	assert.Equal(t, expected, *actual)

}

func TestMKey_UnMarshalJSON_CCAPlatformConfigID_not_ok(t *testing.T) {
	input := []byte(`{"type":"cca.platform-config-id","value":""}`)
	expected := "invalid cca.platform-config-id: empty value"

	mKey := &Mkey{}

	err := mKey.UnmarshalJSON(input)

	assert.EqualError(t, err, expected)
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

		mkey := &Mkey{UintMkey(tv.mkey)}

		actual, err := mkey.MarshalJSON()
		assert.Nil(t, err)

		assert.JSONEq(t, tv.expected, string(actual))
		fmt.Printf("JSON: %x\n", actual)
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
		actual, ok := mKey.Value.(*UintMkey)
		assert.True(t, ok)
		assert.Equal(t, tv.expected, uint64(*actual))
	}
}

func TestMkey_UnmarshalJSON_notok(t *testing.T) {
	tvs := []struct {
		input    []byte
		expected string
	}{
		{
			input:    []byte(`{"type":"uint","value":"abcdefg"}`),
			expected: `invalid uint: json: cannot unmarshal string into Go value of type uint64`,
		},
		{
			input:    []byte(`{"type":"uint","value":123.456}`),
			expected: "invalid uint: json: cannot unmarshal number 123.456 into Go value of type uint64",
		},
	}

	for _, tv := range tvs {
		mKey := &Mkey{}

		err := mKey.UnmarshalJSON(tv.input)

		assert.EqualError(t, err, tv.expected)
	}
}

func TestNewUintMkey(t *testing.T) {
	testVal := UintMkey(7)

	tvs := []struct {
		input    any
		expected UintMkey
		err      string
	}{
		{
			input:    testVal,
			expected: testVal,
		},
		{
			input:    &testVal,
			expected: testVal,
		},
		{
			input:    uint(7),
			expected: testVal,
		},
		{
			input:    uint64(7),
			expected: testVal,
		},
		{
			input:    "7",
			expected: testVal,
		},
		{
			input: true,
			err:   "unexpected type for UintMkey: bool",
		},
	}

	for _, tv := range tvs {
		out, err := NewUintMkey(tv.input)
		if tv.err != "" {
			assert.Nil(t, out)
			assert.EqualError(t, err, tv.err)
		} else {
			assert.Equal(t, tv.expected, *out)
		}
	}
}

func TestNewMkeyOID(t *testing.T) {
	var expectedOID OID
	require.NoError(t, expectedOID.FromString(TestOID))
	expected := TaggedOID(expectedOID)

	out, err := NewMkeyOID(TestOID)
	require.NoError(t, err)
	assert.Equal(t, &expected, out.Value)
}

type testMkey [4]byte

func newTestMkey(_ any) (*Mkey, error) {
	return &Mkey{&testMkey{0x74, 0x64, 0x73, 0x74}}, nil
}

func (o testMkey) PublicKey() (crypto.PublicKey, error) {
	return crypto.PublicKey(o[:]), nil
}

func (o testMkey) Type() string {
	return "test-mkey"
}

func (o testMkey) String() string {
	return "test"
}

func (o testMkey) Valid() error {
	return nil
}

type badMkey struct {
	testMkey
}

func (o badMkey) Type() string {
	return "uuid"
}

func newBadMkey(_ any) (*Mkey, error) {
	return &Mkey{&badMkey{testMkey{0x74, 0x64, 0x73, 0x74}}}, nil
}

func TestRegisterMkeyType(t *testing.T) {
	err := RegisterMkeyType(32, newTestMkey)
	assert.EqualError(t, err, "tag 32 is already registered")

	err = RegisterMkeyType(99996, newBadMkey)
	assert.EqualError(t, err, `mesurement key type with name "uuid" already exists`)

	err = RegisterMkeyType(99996, newTestMkey)
	assert.NoError(t, err)
}

func TestMkey_UnmarshalJSON_regression_issue_100(t *testing.T) {
	u := `31fb5abf-023e-4992-aa4e-95f9c1503bfa`

	tv := []byte(fmt.Sprintf(`{ "type": "uuid", "value": %q }`, u))

	expected, err := NewMkeyUUID(u)
	require.NoError(t, err)

	actual := &Mkey{}
	err = actual.UnmarshalJSON(tv)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
