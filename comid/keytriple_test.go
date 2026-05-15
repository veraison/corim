// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyTriple_Valid(t *testing.T) {
	testCases := []struct {
		title  string
		triple *KeyTriple
		err    string
	}{
		{
			title: "ok minimal",
			triple: &KeyTriple{
				Environment: Environment{
					Instance: MustNewBytesInstance([]byte{0x01, 0x02, 0x03}),
				},
				VerifKeys: *NewCryptoKeys().Add(MustNewCryptoKeyTaggedBytes(
					[]byte{0x04, 0x05, 0x06},
				)),
			},
		},
		{
			title: "ok full",
			triple: &KeyTriple{
				Environment: Environment{
					Instance: MustNewBytesInstance([]byte{0x01, 0x02, 0x03}),
				},
				VerifKeys: *NewCryptoKeys().Add(MustNewCryptoKeyTaggedBytes(
					[]byte{0x04, 0x05, 0x06},
				)),
				Conditions: &KeyTripleCondition{
					Mkey: MustNewMkey("foo", "string"),
					AuthorizedBy: NewCryptoKeys().Add(MustNewCryptoKeyTaggedBytes(
						[]byte{0x07, 0x08, 0x09},
					)),
				},
			},
		},
		{
			title:  "bad no environment",
			triple: &KeyTriple{},
			err:    "environment must not be empty",
		},
		{
			title: "bad no key list",
			triple: &KeyTriple{
				Environment: Environment{
					Instance: MustNewBytesInstance([]byte{0x01, 0x02, 0x03}),
				},
			},
			err: "verification-keys: no keys to validate",
		},
		{
			title: "bad empty condition",
			triple: &KeyTriple{
				Environment: Environment{
					Instance: MustNewBytesInstance([]byte{0x01, 0x02, 0x03}),
				},
				VerifKeys: *NewCryptoKeys().Add(MustNewCryptoKeyTaggedBytes(
					[]byte{0x04, 0x05, 0x06},
				)),
				Conditions: &KeyTripleCondition{},
			},
			err: "condition must not be empty",
		},
		{
			title: "bad invalid mkey",
			triple: &KeyTriple{
				Environment: Environment{
					Instance: MustNewBytesInstance([]byte{0x01, 0x02, 0x03}),
				},
				VerifKeys: *NewCryptoKeys().Add(MustNewCryptoKeyTaggedBytes(
					[]byte{0x04, 0x05, 0x06},
				)),
				Conditions: &KeyTripleCondition{
					Mkey: &Mkey{},
					AuthorizedBy: NewCryptoKeys().Add(MustNewCryptoKeyTaggedBytes(
						[]byte{0x07, 0x08, 0x09},
					)),
				},
			},
			err: "Mkey value not set",
		},
		{
			title: "bad invalid authorized-by",
			triple: &KeyTriple{
				Environment: Environment{
					Instance: MustNewBytesInstance([]byte{0x01, 0x02, 0x03}),
				},
				VerifKeys: *NewCryptoKeys().Add(MustNewCryptoKeyTaggedBytes(
					[]byte{0x04, 0x05, 0x06},
				)),
				Conditions: &KeyTripleCondition{
					Mkey:         MustNewMkey("foo", "string"),
					AuthorizedBy: NewCryptoKeys().Add(&CryptoKey{}),
				},
			},
			err: "invalid key at index 0: CryptoKey not set",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := tc.triple.Valid()
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}
		})
	}
}

func TestKeyTriple_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		triple       *KeyTriple
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title: "minimal",
			triple: &KeyTriple{
				Environment: Environment{
					Instance: MustNewBytesInstance([]byte{0x01, 0x02, 0x03}),
				},
				VerifKeys: *NewCryptoKeys().Add(MustNewCryptoKeyTaggedBytes(
					[]byte{0x04, 0x05, 0x06},
				)),
			},
			expectedCBOR: []byte{
				0x82,             // array(2)
				0xa1,             // . [0]map(1) [environment]
				0x01,             // . . key: 1 [instance]
				0xd9, 0x02, 0x30, // . . value: tag(560) [bytes]
				0x43, // . . . bstr(3)
				0x01, 0x02, 0x03,
				0x81,             // . [1]array(1) [crypto-keys]
				0xd9, 0x02, 0x30, // . . value: tag(560) [bytes]
				0x43, // . . . bstr(3)
				0x04, 0x05, 0x06,
			},
			expectedJSON: `
			{
				"environment": {
					"instance": {"type": "bytes", "value": "AQID"}
				},
				"verification-keys": [
					{"type": "bytes", "value": "BAUG"}
				]
			}
			`,
		},
		{
			title: "full",
			triple: &KeyTriple{
				Environment: Environment{
					Instance: MustNewBytesInstance([]byte{0x01, 0x02, 0x03}),
				},
				VerifKeys: *NewCryptoKeys().Add(MustNewCryptoKeyTaggedBytes(
					[]byte{0x04, 0x05, 0x06},
				)),
				Conditions: &KeyTripleCondition{
					Mkey: MustNewMkey("foo", "string"),
					AuthorizedBy: NewCryptoKeys().Add(MustNewCryptoKeyTaggedBytes(
						[]byte{0x07, 0x08, 0x09},
					)),
				},
			},
			expectedCBOR: []byte{
				0x83,             // array(3)
				0xa1,             // . [0]map(1) [environment]
				0x01,             // . . key: 1 [instance]
				0xd9, 0x02, 0x30, // . . value: tag(560) [bytes]
				0x43, //             . . . bstr(3)
				0x01, 0x02, 0x03,
				0x81,             // . [1]array(1) [crypto-keys]
				0xd9, 0x02, 0x30, // . . value: tag(560) [bytes]
				0x43, //             . . . bstr(3)
				0x04, 0x05, 0x06,
				0xa2,             // . [2]map(2) [identy-triple-condition]
				0x00,             // . . key: 0 [mkey]
				0x63,             // . . value: tstr(3)
				0x66, 0x6f, 0x6f, // . . . "foo"
				0x01,             // . . key: 1 [authorized-by]
				0x81,             // . . value: array(1) [crypto-keys]
				0xd9, 0x02, 0x30, // . . . [0]tag(560) [bytes]
				0x43, //             . . . . bstr(3)
				0x07, 0x08, 0x09,
			},
			expectedJSON: `
			{
				"environment": {
					"instance": {"type": "bytes", "value": "AQID"}
				},
				"verification-keys": [
					{"type": "bytes", "value": "BAUG"}
				],
				"conditions": {
					"mkey": {"type": "string", "value": "foo"},
					"authorized-by": [
						{"type": "bytes", "value": "BwgJ"}
					]
				}
			}
			`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			bytes, err := em.Marshal(&tc.triple)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, bytes)

			var decoded KeyTriple
			err = dm.Unmarshal(bytes, &decoded)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.triple.Environment, decoded.Environment)

			bytes, err = json.Marshal(&tc.triple)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(bytes))

			err = json.Unmarshal(bytes, &decoded)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.triple.Environment, decoded.Environment)
		})
	}
}

func TestKeyTriple_UnmarshalCBOR_bad(t *testing.T) {
	testCases := []struct {
		title string
		data  []byte
		err   string
	}{
		{
			title: "invalid CBOR",
			data:  MustHexDecode(t, "81"),
			err:   "unexpected EOF",
		},
		{
			title: "wrong len",
			data:  MustHexDecode(t, "8101"),
			err:   "expected array between 2 and 3 elements",
		},
		{
			title: "bad environment",
			data:  MustHexDecode(t, "820102"),
			err:   "environment: cbor: cannot unmarshal positive integer",
		},
		{
			title: "bad keys",
			data:  MustHexDecode(t, "82a101d902304301020301"),
			err:   "verification-keys: cbor: cannot unmarshal positive integer",
		},
		{
			title: "bad condition",
			data:  MustHexDecode(t, "83a101d902304301020381d902304304050601"),
			err:   "conditions: cbor: cannot unmarshal positive integer",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var triple KeyTriple
			err := triple.UnmarshalCBOR(tc.data)
			assert.ErrorContains(t, err, tc.err)
		})
	}
}

func TestKeyTriples_Add(t *testing.T) {
	triples := NewKeyTriples()
	assert.Len(t, *triples, 0)

	triples.Add(&KeyTriple{})
	assert.Len(t, *triples, 1)
}
