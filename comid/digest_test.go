// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDigestAlgorithm_conversion(t *testing.T) {
	alg := IntDigestAlgorithm(Sha256)
	assert.True(t, alg.IsInt())
	assert.False(t, alg.IsString())
	assert.Equal(t, Sha256, alg.Int())
	assert.Equal(t, "sha-256", alg.String())

	alg = StringDigestAlgorithm("foo")
	assert.False(t, alg.IsInt())
	assert.True(t, alg.IsString())
	assert.Equal(t, 0, alg.Int())
	assert.Equal(t, "foo", alg.String())
}

func TestDigestAlgorithmFromString(t *testing.T) {
	testCases := []struct {
		title    string
		text     string
		expected DigestAlgorithm
	}{
		{
			title:    "int",
			text:     "-1",
			expected: DigestAlgorithm{-1},
		},
		{
			title:    "known string",
			text:     "sha-256",
			expected: DigestAlgorithm{Sha256},
		},
		{
			title:    "unknown string",
			text:     "foo",
			expected: DigestAlgorithm{"foo"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			alg := DigestAlgorithmFromString(tc.text)
			assert.EqualValues(t, tc.expected, alg)
		})
	}
}

func TestDigestAlgorithmFromAny(t *testing.T) {
	testCases := []struct {
		title    string
		value    any
		expected DigestAlgorithm
		err      string
	}{
		{
			title:    "int",
			value:    Sha256,
			expected: DigestAlgorithm{Sha256},
		},
		{
			title:    "int64",
			value:    int64(-1),
			expected: DigestAlgorithm{-1},
		},
		{
			title:    "float64",
			value:    -1.0,
			expected: DigestAlgorithm{-1},
		},
		{
			title:    "string",
			value:    "foo",
			expected: DigestAlgorithm{"foo"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			alg, err := DigestAlgorithmFromAny(tc.value)
			if tc.err == "" {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expected, alg)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}
		})
	}
}

func TestDigestAlgorithm_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		value        DigestAlgorithm
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title:        "known int",
			value:        IntDigestAlgorithm(Sha256),
			expectedCBOR: []byte{0x01},
			expectedJSON: `1`,
		},
		{
			title:        "nagative int",
			value:        IntDigestAlgorithm(-1),
			expectedCBOR: []byte{0x20},
			expectedJSON: `-1`,
		},
		{
			title: "string",
			value: StringDigestAlgorithm("foo"),
			expectedCBOR: []byte{
				0x63,             // tstr(3)
				0x66, 0x6f, 0x6f, // . "foo"
			},
			expectedJSON: `"foo"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			bytes, err := em.Marshal(tc.value)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, bytes)

			var alg DigestAlgorithm
			err = dm.Unmarshal(bytes, &alg)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.value, alg)

			bytes, err = json.Marshal(tc.value)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(bytes))

			err = json.Unmarshal(bytes, &alg)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.value, alg)
		})
	}

}

func TestDigestAlgorithm_UnmarshalCBOR_bad(t *testing.T) {
	var alg DigestAlgorithm
	err := alg.UnmarshalCBOR([]byte{})
	assert.ErrorContains(t, err, "buffer too short")

	err = alg.UnmarshalCBOR(MustHexDecode(t, "19"))
	assert.ErrorContains(t, err, "unexpected EOF")

	err = alg.UnmarshalCBOR(MustHexDecode(t, "64ffffffff"))
	assert.ErrorContains(t, err, "invalid UTF-8 string")

	err = alg.UnmarshalCBOR(MustHexDecode(t, "f4"))
	assert.ErrorContains(t, err, "unexpected CBOR major type")
}

func TestDigestAlgorithm_UnmarshalJSON_bad(t *testing.T) {
	var alg DigestAlgorithm
	err := alg.UnmarshalJSON([]byte{})
	assert.ErrorContains(t, err, "unexpected end of JSON input")

	err = alg.UnmarshalJSON([]byte("true"))
	assert.ErrorContains(t, err, "unexpected algorithm value: true(bool)")
}

func TestDigestFromString(t *testing.T) {
	digest, err := DigestFromString("sha-256;AQID")
	assert.NoError(t, err)
	assert.EqualValues(t, NewDigestIntAlg(Sha256, []byte{0x01, 0x02, 0x03}), digest)

	_, err = DigestFromString("foo")
	assert.ErrorContains(t, err, `expected exactly two ;-separated parts, got "foo"`)

	_, err = DigestFromString("foo;@@@")
	assert.ErrorContains(t, err, "val: illegal base64 data")
}

func TestDigest_Valid(t *testing.T) {
	err := NewDigestIntAlg(Sha256, []byte{}).Valid()
	assert.ErrorContains(t, err, "zero length value")

	err = NewDigestIntAlg(0, []byte{0x1, 0x2, 0x3}).Valid()
	assert.ErrorContains(t, err, "zero algorithm")

	err = NewDigestIntAlg(Sha256, []byte{0x1, 0x2, 0x3}).Valid()
	assert.ErrorContains(t, err, "length mismatch for hash algorithm sha-256")

	err = NewDigestStringAlg("foo", []byte{0x1, 0x2, 0x3}).Valid()
	assert.NoError(t, err)
}

func TestDigest_round_trip(t *testing.T) {
	bytes := MustHexDecode(t, "deadbeef")
	testCases := []struct {
		title        string
		value        Digest
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title: "int",
			value: NewDigestIntAlg(Sha256, bytes),
			expectedCBOR: []byte{
				0x82, // array(2)
				0x01, // . [0]1 [sha-256]
				0x44, // . [1]bstr(4)
				0xde, 0xad, 0xbe, 0xef,
			},
			expectedJSON: `[1, "3q2-7w"]`,
		},
		{
			title: "string",
			value: NewDigestStringAlg("foo", bytes),
			expectedCBOR: []byte{
				0x82,             // array(2)
				0x63,             // . [0]tstr(3)
				0x66, 0x6f, 0x6f, // . . "foo"
				0x44, //             . [1]bstr(4)
				0xde, 0xad, 0xbe, 0xef,
			},
			expectedJSON: `["foo", "3q2-7w"]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			bytes, err := em.Marshal(tc.value)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, bytes)

			var digest Digest
			err = dm.Unmarshal(bytes, &digest)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.value, digest)

			bytes, err = json.Marshal(tc.value)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(bytes))

			err = json.Unmarshal(bytes, &digest)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.value, digest)
		})
	}

}

func TestDigest_UnmarshalJSON_bad(t *testing.T) {
	var digest Digest

	err := digest.UnmarshalJSON([]byte(""))
	assert.ErrorContains(t, err, "unexpected end of JSON input")

	err = digest.UnmarshalJSON([]byte(`{"alg": 1, "value": "foo"}`))
	assert.ErrorContains(t, err, "cannot unmarshal object into Go value of type []interface {}")

	err = digest.UnmarshalJSON([]byte("[1, 2, 3]"))
	assert.ErrorContains(t, err, "expected array with two elements")

	err = digest.UnmarshalJSON([]byte(`[true, "bar"]`))
	assert.ErrorContains(t, err, "invalid digest algorithm: true(bool)")

	err = digest.UnmarshalJSON([]byte(`[1, "@@@"]`))
	assert.ErrorContains(t, err, "val: illegal base64 data")
}
