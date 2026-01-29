// Copyright 2021-2026 Contributors to the Veraison project.
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

// Test Marshal and Unmarshal of Cryptokeys(tag 13)
func TestMeasurement_CryptoKeys_RoundTrip(t *testing.T) {
	// Create a new measurement with a valid key
	m := MustNewMeasurement("31fb5abf-023e-4992-aa4e-95f9c1503bfa", UUIDType)

	// Create a CryptoKey (using a dummy PKIX base64 key for testing)
	pkText := `-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAGb9ECWmEzf6FQbrBZ9w7lshQhqowtrbLDFw4rXAxZuE=
-----END PUBLIC KEY-----`
	ck := MustNewPKIXBase64Key(pkText)

	// Add CryptoKey to Measurement
	m.AddCryptoKey(ck)

	// Marshal to CBOR
	// Use em (from cbor.go)
	data, err := em.Marshal(m)
	require.NoError(t, err)

	// Check if tag 13 is present in the CBOR output hex
	// We can try to decode it back to verify.

	// Unmarshal back
	var m2 Measurement
	err = dm.Unmarshal(data, &m2)
	require.NoError(t, err)

	// Verify CryptoKeys are present and correct
	require.NotNil(t, m2.Val.CryptoKeys)
	require.Len(t, *m2.Val.CryptoKeys, 1)
	assert.Equal(t, pkText, (*m2.Val.CryptoKeys)[0].String())
}
