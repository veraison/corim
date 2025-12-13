// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"crypto"
	"fmt"
	"net"
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

	assert.NotNil(t, tv.AddDigest(swid.Sha256, []byte{0xff}))
	assert.ErrorContains(t, tv.Valid(), "digest at index 0")
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

func TestMeasurement_NameMeasurement(t *testing.T) {
	want := "Maureen"
	got := *(&Measurement{}).SetName("Maureen").Val.Name
	assert.Equal(t, want, got)
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
			expected: "unexpected CBOR major type for mkey: 5",
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

func TestNewStringMkey(t *testing.T) {
	testString := "foo"
	testBytes := MustHexDecode(t, "666f6f")
	testBytesBad := MustHexDecode(t, "ff")
	testVal := StringMkey(testString)

	tvs := []struct {
		input    any
		expected StringMkey
		err      string
	}{
		{
			input:    testString,
			expected: testVal,
		},
		{
			input:    &testString,
			expected: testVal,
		},
		{
			input:    testBytes,
			expected: testVal,
		},
		{
			input:    &testBytes,
			expected: testVal,
		},
		{
			input:    testVal,
			expected: testVal,
		},
		{
			input:    &testVal,
			expected: testVal,
		},
		{
			input: testBytesBad,
			err:   "invalid utf-8 string: ff",
		},
		{
			input: &testBytesBad,
			err:   "invalid utf-8 string: ff",
		},
		{
			input: 7,
			err:   "unexpected type for StringMkey: int",
		},
	}

	for _, tv := range tvs {
		out, err := NewStringMkey(tv.input)
		if tv.err != "" {
			assert.Nil(t, out)
			assert.EqualError(t, err, tv.err)
		} else {
			assert.Equal(t, tv.expected, *out)
		}
	}
}

func TestMKey_string_marshaling_round_trip(t *testing.T) {
	tvs := []struct {
		input         *Mkey
		expected_json []byte
		expected_cbor []byte
	}{
		{
			input:         MustNewMkey("foo", "string"),
			expected_json: []byte(`{"type":"string","value":"foo"}`),
			expected_cbor: MustHexDecode(t, "63666f6f"),
		},
	}

	for _, tv := range tvs {
		actual_json, err := tv.input.MarshalJSON()
		assert.Nil(t, err)
		assert.Equal(t, tv.expected_json, actual_json)

		var key Mkey
		err = key.UnmarshalJSON(actual_json)
		assert.Nil(t, err)
		assert.Equal(t, tv.input, &key)

		actual_cbor, err := tv.input.MarshalCBOR()
		assert.Nil(t, err)
		assert.Equal(t, tv.expected_cbor, actual_cbor)

		err = key.UnmarshalCBOR(actual_cbor)
		assert.Nil(t, err)
		assert.Equal(t, tv.input, &key)
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
	assert.EqualError(t, err, `measurement key type with name "uuid" already exists`)

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

func TestMkey_new(t *testing.T) {
	psaRefValID, err := NewPSARefValID(TestSignerID)
	require.NoError(t, err)

	key := MustNewMkey(psaRefValID, PSARefValIDType)
	assert.EqualValues(t, psaRefValID, key.Value)
	assert.Equal(t, PSARefValIDType, key.Type())
}

func TestMkey_UintMkey(t *testing.T) {
	var v uint64 = 7
	key, err := NewMkey(v, UintType)
	assert.NoError(t, err)
	assert.Equal(t, "7", key.Value.String())

	ret, err := key.GetKeyUint()
	assert.NoError(t, err)
	assert.EqualValues(t, 7, ret)
}

func TestMval_Valid(t *testing.T) {
	t.Run("No fields set", func(t *testing.T) {
		mval := Mval{}
		err := mval.Valid()
		assert.EqualError(t, err, "no measurement value set")
	})

	t.Run("All fields nil except Ver, which is valid", func(t *testing.T) {
		var scheme swid.VersionScheme
		_ = scheme.SetCode(swid.VersionSchemeSemVer)
		mval := Mval{
			Ver: &Version{
				Version: "1.0",
				Scheme:  scheme,
			},
		}
		err := mval.Valid()
		assert.NoError(t, err)
	})

	// Test with valid 6-byte MAC
	t.Run("MACAddr valid (6 bytes)", func(t *testing.T) {
		mac := MACaddr([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}) // EUI-48
		mval := Mval{MACAddr: &mac}
		err := mval.Valid()
		assert.NoError(t, err, "6-byte MAC should be valid")
	})

	// Test with valid 8-byte MAC
	t.Run("MACAddr valid (8 bytes)", func(t *testing.T) {
		mac := MACaddr([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}) // EUI-64
		mval := Mval{MACAddr: &mac}
		err := mval.Valid()
		assert.NoError(t, err, "8-byte MAC should be valid")
	})

	// Test with invalid MAC length
	t.Run("MACAddr invalid (too many bytes)", func(t *testing.T) {
		mac := MACaddr([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}) // 7 bytes
		mval := Mval{MACAddr: &mac}
		err := mval.Valid()
		assert.EqualError(t, err, "invalid MAC address length: expected 6 or 8 bytes, got 7")
	})

	// Test with invalid MAC length
	t.Run("MACAddr invalid (too few bytes)", func(t *testing.T) {
		mac := MACaddr([]byte{0x01, 0x02, 0x03, 0x04}) // 4 bytes
		mval := Mval{MACAddr: &mac}
		err := mval.Valid()
		assert.EqualError(t, err, "invalid MAC address length: expected 6 or 8 bytes, got 4")
	})

	t.Run("MACAddr valid (6 bytes)", func(t *testing.T) {
		mac := MACaddr([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
		mval := Mval{MACAddr: &mac}
		err := mval.Valid()
		assert.NoError(t, err)
	})

	t.Run("IPAddr valid (IPv4)", func(t *testing.T) {
		ip := net.IPv4(192, 168, 1, 100)
		mval := Mval{IPAddr: &ip}
		err := mval.Valid()
		assert.NoError(t, err)
	})

	t.Run("IPAddr valid (IPv6)", func(t *testing.T) {
		ip := net.ParseIP("2001:db8::1")
		mval := Mval{IPAddr: &ip}
		err := mval.Valid()
		assert.NoError(t, err)
	})

	t.Run("Digests invalid", func(t *testing.T) {
		ds := NewDigests()
		ds.AddDigest(swid.Sha256, []byte{0xAA, 0xBB})
		mval := Mval{
			Digests: ds,
		}
		err := mval.Valid()
		assert.ErrorContains(t, err, "digest at index 0")
	})

	t.Run("Extensions valid", func(t *testing.T) {
		// Suppose we have some extension data that is considered valid
		ext := Extensions{}
		mval := Mval{
			Extensions: ext,
			// Must also set one non-empty field to pass "no measurement value set"
			Ver: &Version{Version: "1.0"},
		}
		err := mval.Valid()
		assert.NoError(t, err)
	})
}

// Unit tests for MKey with PSA refval-id

func TestMKey_MarshalCBOR_PSARefValID_ok(t *testing.T) {
	signerID32 := MustHexDecode(t, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	signerID48 := MustHexDecode(t, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	signerID64 := MustHexDecode(t, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")

	tvs := []struct {
		name     string
		refvalID *PSARefValID
		expected []byte
	}{
		{
			name:     "PSA RefVal ID with 32-byte signer ID",
			refvalID: MustCreatePSARefValID(signerID32, "BL", "2.1.0"),
			expected: MustHexDecode(t, "d90259a30162424c0465322e312e30055820deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
		},
		{
			name:     "PSA RefVal ID with 48-byte signer ID",
			refvalID: MustCreatePSARefValID(signerID48, "PRoT", "1.3.5"),
			expected: MustHexDecode(t, "d90259a3016450526f540465312e332e35055830deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
		},
		{
			name:     "PSA RefVal ID with 64-byte signer ID",
			refvalID: MustCreatePSARefValID(signerID64, "M1", "5.0.7"),
			expected: MustHexDecode(t, "d90259a301624d310465352e302e37055840deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
		},
	}

	for _, tv := range tvs {
		t.Run(tv.name, func(t *testing.T) {
			mkey, err := NewMkeyPSARefvalID(tv.refvalID)
			require.NoError(t, err)

			actual, err := mkey.MarshalCBOR()
			require.NoError(t, err)
			assert.Equal(t, tv.expected, actual)
		})
	}
}

func TestMKey_UnmarshalCBOR_PSARefValID_ok(t *testing.T) {
	signerID32 := MustHexDecode(t, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	signerID48 := MustHexDecode(t, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	signerID64 := MustHexDecode(t, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")

	tvs := []struct {
		name     string
		input    []byte
		expected *PSARefValID
	}{
		{
			name:     "PSA RefVal ID with 32-byte signer ID",
			input:    MustHexDecode(t, "d90259a30162424c0465322e312e30055820deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
			expected: MustCreatePSARefValID(signerID32, "BL", "2.1.0"),
		},
		{
			name:     "PSA RefVal ID with 48-byte signer ID",
			input:    MustHexDecode(t, "d90259a3016450526f540465312e332e35055830deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
			expected: MustCreatePSARefValID(signerID48, "PRoT", "1.3.5"),
		},
		{
			name:     "PSA RefVal ID with 64-byte signer ID",
			input:    MustHexDecode(t, "d90259a301624d310465352e302e37055840deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
			expected: MustCreatePSARefValID(signerID64, "M1", "5.0.7"),
		},
	}

	for _, tv := range tvs {
		t.Run(tv.name, func(t *testing.T) {
			var mkey Mkey
			err := mkey.UnmarshalCBOR(tv.input)
			require.NoError(t, err)

			actual, ok := mkey.Value.(*TaggedPSARefValID)
			require.True(t, ok)
			assert.Equal(t, tv.expected.Label, (*PSARefValID)(actual).Label)
			assert.Equal(t, tv.expected.Version, (*PSARefValID)(actual).Version)
			assert.Equal(t, tv.expected.SignerID, (*PSARefValID)(actual).SignerID)
		})
	}
}

func TestMKey_UnmarshalCBOR_PSARefValID_nok(t *testing.T) {
	tvs := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Invalid CBOR data - truncated",
			input:    MustHexDecode(t, "d9025ba301624254"),
			expected: "unexpected EOF",
		},
		{
			name:     "Invalid CBOR data - wrong type",
			input:    MustHexDecode(t, "43616263"),
			expected: "unexpected CBOR major type for mkey: 2",
		},
	}

	for _, tv := range tvs {
		t.Run(tv.name, func(t *testing.T) {
			var mkey Mkey
			err := mkey.UnmarshalCBOR(tv.input)
			assert.ErrorContains(t, err, tv.expected)
		})
	}
}

func TestMKey_MarshalJSON_PSARefValID_ok(t *testing.T) {
	signerID32 := MustHexDecode(t, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")

	tvs := []struct {
		name     string
		refvalID *PSARefValID
		expected string
	}{
		{
			name:     "PSA RefVal ID with all fields",
			refvalID: MustCreatePSARefValID(signerID32, "BL", "2.1.0"),
			expected: `{"type":"psa.refval-id","value":{"label":"BL","version":"2.1.0","signer-id":"3q2+796tvu/erb7v3q2+796tvu/erb7v3q2+796tvu8="}}`,
		},
		{
			name:     "PSA RefVal ID with signer ID only",
			refvalID: MustCreatePSARefValID(signerID32, "", ""),
			expected: `{"type":"psa.refval-id","value":{"label":"","version":"","signer-id":"3q2+796tvu/erb7v3q2+796tvu/erb7v3q2+796tvu8="}}`,
		},
	}

	for _, tv := range tvs {
		t.Run(tv.name, func(t *testing.T) {
			mkey, err := NewMkeyPSARefvalID(tv.refvalID)
			require.NoError(t, err)

			actual, err := mkey.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tv.expected, string(actual))
		})
	}
}

func TestMKey_UnmarshalJSON_PSARefValID_ok(t *testing.T) {
	signerID32 := MustHexDecode(t, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")

	tvs := []struct {
		name     string
		input    string
		expected *PSARefValID
	}{
		{
			name:     "PSA RefVal ID with all fields",
			input:    `{"type":"psa.refval-id","value":{"label":"BL","version":"2.1.0","signer-id":"3q2+796tvu/erb7v3q2+796tvu/erb7v3q2+796tvu8="}}`,
			expected: MustCreatePSARefValID(signerID32, "BL", "2.1.0"),
		},
		{
			name:     "PSA RefVal ID with signer ID only",
			input:    `{"type":"psa.refval-id","value":{"signer-id":"3q2+796tvu/erb7v3q2+796tvu/erb7v3q2+796tvu8="}}`,
			expected: MustCreatePSARefValID(signerID32, "", ""),
		},
	}

	for _, tv := range tvs {
		t.Run(tv.name, func(t *testing.T) {
			var mkey Mkey
			err := mkey.UnmarshalJSON([]byte(tv.input))
			require.NoError(t, err)

			actual, ok := mkey.Value.(*TaggedPSARefValID)
			require.True(t, ok)

			if tv.expected.Label != nil && *tv.expected.Label != "" {
				assert.Equal(t, tv.expected.Label, (*PSARefValID)(actual).Label)
			}
			if tv.expected.Version != nil && *tv.expected.Version != "" {
				assert.Equal(t, tv.expected.Version, (*PSARefValID)(actual).Version)
			}
			assert.Equal(t, tv.expected.SignerID, (*PSARefValID)(actual).SignerID)
		})
	}
}

func TestMKey_UnmarshalJSON_PSARefValID_nok(t *testing.T) {
	tvs := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Invalid JSON - missing signer-id",
			input:    `{"type":"psa.refval-id","value":{"label":"BL"}}`,
			expected: "missing mandatory signer ID",
		},
		{
			name:     "Invalid JSON - empty signer-id",
			input:    `{"type":"psa.refval-id","value":{"signer-id":""}}`,
			expected: "want 32, 48 or 64 bytes",
		},
		{
			name:     "Invalid JSON - wrong signer-id length",
			input:    `{"type":"psa.refval-id","value":{"signer-id":"YWJjZA=="}}`,
			expected: "want 32, 48 or 64 bytes",
		},
		{
			name:     "Invalid JSON - malformed JSON",
			input:    `{"type":"psa.refval-id","value":`,
			expected: "unexpected end of JSON input",
		},
	}

	for _, tv := range tvs {
		t.Run(tv.name, func(t *testing.T) {
			var mkey Mkey
			err := mkey.UnmarshalJSON([]byte(tv.input))
			assert.ErrorContains(t, err, tv.expected)
		})
	}
}

func TestMKey_PSARefValID_RoundTrip(t *testing.T) {
	signerID := MustHexDecode(t, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	refvalID := MustCreatePSARefValID(signerID, "BL", "5.0.0")

	t.Run("CBOR round trip", func(t *testing.T) {
		mkey, err := NewMkeyPSARefvalID(refvalID)
		require.NoError(t, err)

		cborData, err := mkey.MarshalCBOR()
		require.NoError(t, err)

		var decoded Mkey
		err = decoded.UnmarshalCBOR(cborData)
		require.NoError(t, err)

		decodedVal, ok := decoded.Value.(*TaggedPSARefValID)
		require.True(t, ok)
		assert.Equal(t, refvalID.SignerID, (*PSARefValID)(decodedVal).SignerID)
		assert.Equal(t, refvalID.Label, (*PSARefValID)(decodedVal).Label)
		assert.Equal(t, refvalID.Version, (*PSARefValID)(decodedVal).Version)
	})

	t.Run("JSON round trip", func(t *testing.T) {
		mkey, err := NewMkeyPSARefvalID(refvalID)
		require.NoError(t, err)

		jsonData, err := mkey.MarshalJSON()
		require.NoError(t, err)

		var decoded Mkey
		err = decoded.UnmarshalJSON(jsonData)
		require.NoError(t, err)

		decodedVal, ok := decoded.Value.(*TaggedPSARefValID)
		require.True(t, ok)
		assert.Equal(t, refvalID.SignerID, (*PSARefValID)(decodedVal).SignerID)
		assert.Equal(t, refvalID.Label, (*PSARefValID)(decodedVal).Label)
		assert.Equal(t, refvalID.Version, (*PSARefValID)(decodedVal).Version)
	})
}

// Unit tests for MKey with UUID

func TestMKey_MarshalCBOR_UUID_ok(t *testing.T) {
	tvs := []struct {
		name     string
		uuid     string
		expected []byte
	}{
		{
			name:     "Valid RFC4122 UUID v4",
			uuid:     "31fb5abf-023e-4992-aa4e-95f9c1503bfa",
			expected: MustHexDecode(t, "d8255031fb5abf023e4992aa4e95f9c1503bfa"),
		},
		{
			name:     "Valid RFC4122 UUID v1",
			uuid:     "f47ac10b-58cc-0372-8567-0e02b2c3d479",
			expected: MustHexDecode(t, "d82550f47ac10b58cc037285670e02b2c3d479"),
		},
		{
			name:     "Another valid RFC4122 UUID",
			uuid:     "550e8400-e29b-41d4-a716-446655440000",
			expected: MustHexDecode(t, "d82550550e8400e29b41d4a716446655440000"),
		},
	}

	for _, tv := range tvs {
		t.Run(tv.name, func(t *testing.T) {
			mkey, err := NewMkeyUUID(tv.uuid)
			require.NoError(t, err)

			actual, err := mkey.MarshalCBOR()
			require.NoError(t, err)
			assert.Equal(t, tv.expected, actual)
		})
	}
}

func TestMKey_UnmarshalCBOR_UUID_ok(t *testing.T) {
	tvs := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Valid RFC4122 UUID v4",
			input:    MustHexDecode(t, "d8255031fb5abf023e4992aa4e95f9c1503bfa"),
			expected: "31fb5abf-023e-4992-aa4e-95f9c1503bfa",
		},
		{
			name:     "Valid RFC4122 UUID v1",
			input:    MustHexDecode(t, "d82550f47ac10b58cc037285670e02b2c3d479"),
			expected: "f47ac10b-58cc-0372-8567-0e02b2c3d479",
		},
		{
			name:     "Another valid RFC4122 UUID",
			input:    MustHexDecode(t, "d82550550e8400e29b41d4a716446655440000"),
			expected: "550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for _, tv := range tvs {
		t.Run(tv.name, func(t *testing.T) {
			var mkey Mkey
			err := mkey.UnmarshalCBOR(tv.input)
			require.NoError(t, err)

			actual, ok := mkey.Value.(*TaggedUUID)
			require.True(t, ok)
			assert.Equal(t, tv.expected, actual.String())
		})
	}
}

func TestMKey_UnmarshalCBOR_UUID_nok(t *testing.T) {
	tvs := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Invalid CBOR - truncated UUID",
			input:    MustHexDecode(t, "d825504831fb5abf023e4992"),
			expected: "unexpected EOF",
		},
		{
			name:     "Invalid CBOR - wrong type",
			input:    MustHexDecode(t, "43616263"),
			expected: "unexpected CBOR major type for mkey: 2",
		},
		{
			name:     "Invalid CBOR - extraneous data",
			input:    MustHexDecode(t, "d82550f47ac10b58cc037285670e02b2c3d47900"),
			expected: "extraneous data",
		},
	}

	for _, tv := range tvs {
		t.Run(tv.name, func(t *testing.T) {
			var mkey Mkey
			err := mkey.UnmarshalCBOR(tv.input)
			assert.ErrorContains(t, err, tv.expected)
		})
	}
}

func TestMKey_MarshalJSON_UUID_ok(t *testing.T) {
	tvs := []struct {
		name     string
		uuid     string
		expected string
	}{
		{
			name:     "Valid RFC4122 UUID v4",
			uuid:     "31fb5abf-023e-4992-aa4e-95f9c1503bfa",
			expected: `{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}`,
		},
		{
			name:     "Valid RFC4122 UUID v1",
			uuid:     "f47ac10b-58cc-0372-8567-0e02b2c3d479",
			expected: `{"type":"uuid","value":"f47ac10b-58cc-0372-8567-0e02b2c3d479"}`,
		},
		{
			name:     "Another valid RFC4122 UUID",
			uuid:     "550e8400-e29b-41d4-a716-446655440000",
			expected: `{"type":"uuid","value":"550e8400-e29b-41d4-a716-446655440000"}`,
		},
	}

	for _, tv := range tvs {
		t.Run(tv.name, func(t *testing.T) {
			mkey, err := NewMkeyUUID(tv.uuid)
			require.NoError(t, err)

			actual, err := mkey.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tv.expected, string(actual))
		})
	}
}

func TestMKey_UnmarshalJSON_UUID_ok(t *testing.T) {
	tvs := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid RFC4122 UUID v4",
			input:    `{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}`,
			expected: "31fb5abf-023e-4992-aa4e-95f9c1503bfa",
		},
		{
			name:     "Valid RFC4122 UUID v1",
			input:    `{"type":"uuid","value":"f47ac10b-58cc-0372-8567-0e02b2c3d479"}`,
			expected: "f47ac10b-58cc-0372-8567-0e02b2c3d479",
		},
		{
			name:     "Another valid RFC4122 UUID",
			input:    `{"type":"uuid","value":"550e8400-e29b-41d4-a716-446655440000"}`,
			expected: "550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for _, tv := range tvs {
		t.Run(tv.name, func(t *testing.T) {
			var mkey Mkey
			err := mkey.UnmarshalJSON([]byte(tv.input))
			require.NoError(t, err)

			actual, ok := mkey.Value.(*TaggedUUID)
			require.True(t, ok)
			assert.Equal(t, tv.expected, actual.String())
		})
	}
}

func TestMKey_UnmarshalJSON_UUID_nok(t *testing.T) {
	tvs := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Invalid UUID format",
			input:    `{"type":"uuid","value":"not-a-uuid"}`,
			expected: "bad UUID",
		},
		{
			name:     "Empty UUID",
			input:    `{"type":"uuid","value":""}`,
			expected: "bad UUID",
		},
		{
			name:     "Non-RFC4122 UUID",
			input:    `{"type":"uuid","value":"f47ac10b-58cc-4372-c567-0e02b2c3d479"}`,
			expected: "expecting RFC4122 UUID",
		},
		{
			name:     "Malformed JSON",
			input:    `{"type":"uuid","value":`,
			expected: "unexpected end of JSON input",
		},
		{
			name:     "Wrong value type - integer",
			input:    `{"type":"uuid","value":123}`,
			expected: "json: cannot unmarshal number",
		},
	}

	for _, tv := range tvs {
		t.Run(tv.name, func(t *testing.T) {
			var mkey Mkey
			err := mkey.UnmarshalJSON([]byte(tv.input))
			assert.ErrorContains(t, err, tv.expected)
		})
	}
}

func TestMKey_UUID_RoundTrip(t *testing.T) {
	uuidStr := "31fb5abf-023e-4992-aa4e-95f9c1503bfa"

	t.Run("CBOR round trip", func(t *testing.T) {
		mkey, err := NewMkeyUUID(uuidStr)
		require.NoError(t, err)

		cborData, err := mkey.MarshalCBOR()
		require.NoError(t, err)

		var decoded Mkey
		err = decoded.UnmarshalCBOR(cborData)
		require.NoError(t, err)

		decodedVal, ok := decoded.Value.(*TaggedUUID)
		require.True(t, ok)
		assert.Equal(t, uuidStr, decodedVal.String())
	})

	t.Run("JSON round trip", func(t *testing.T) {
		mkey, err := NewMkeyUUID(uuidStr)
		require.NoError(t, err)

		jsonData, err := mkey.MarshalJSON()
		require.NoError(t, err)

		var decoded Mkey
		err = decoded.UnmarshalJSON(jsonData)
		require.NoError(t, err)

		decodedVal, ok := decoded.Value.(*TaggedUUID)
		require.True(t, ok)
		assert.Equal(t, uuidStr, decodedVal.String())
	})
}

func TestMKey_UUID_Bytes_Input(t *testing.T) {
	uuidBytes := []byte{
		0x31, 0xfb, 0x5a, 0xbf, 0x02, 0x3e, 0x49, 0x92,
		0xaa, 0x4e, 0x95, 0xf9, 0xc1, 0x50, 0x3b, 0xfa,
	}

	mkey, err := NewMkeyUUID(uuidBytes)
	require.NoError(t, err)

	val, ok := mkey.Value.(*TaggedUUID)
	require.True(t, ok)
	assert.Equal(t, "31fb5abf-023e-4992-aa4e-95f9c1503bfa", val.String())
}

func TestMKey_UUID_Invalid_Byte_Length(t *testing.T) {
	tvs := []struct {
		name  string
		bytes []byte
	}{
		{
			name:  "Too short - 15 bytes",
			bytes: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
		},
		{
			name:  "Too long - 17 bytes",
			bytes: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11},
		},
		{
			name:  "Empty bytes",
			bytes: []byte{},
		},
	}

	for _, tv := range tvs {
		t.Run(tv.name, func(t *testing.T) {
			_, err := NewMkeyUUID(tv.bytes)
			assert.Error(t, err)
		})
	}
}
