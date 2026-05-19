// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package corim

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/corim/comid"
)

func Test_OneOrMore_serilize_round_trip(t *testing.T) {
	stringCases := []struct {
		title        string
		oom          OneOrMore[string]
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title: "one string",
			oom:   OneOrMore[string]{"foo"},
			expectedCBOR: []byte{
				0x63,             // tstr(3)
				0x66, 0x6f, 0x6f, // . "foo"
			},
			expectedJSON: `"foo"`,
		},
		{
			title: "more strings",
			oom:   OneOrMore[string]{"foo", "bar"},
			expectedCBOR: []byte{
				0x82,             // array(2)
				0x63,             // . [0]tstr(3)
				0x66, 0x6f, 0x6f, // . . "foo"
				0x63,             // . [1]tstr(3)
				0x62, 0x61, 0x72, // . . "bar"
			},
			expectedJSON: `["foo","bar"]`,
		},
	}

	for _, tc := range stringCases {
		t.Run(tc.title, func(t *testing.T) {
			encoded, err := tc.oom.MarshalCBOR()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, encoded)

			var decoded OneOrMore[string]
			err = decoded.UnmarshalCBOR(encoded)
			assert.NoError(t, err)
			assert.Equal(t, tc.oom, decoded)

			encoded, err = tc.oom.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedJSON, string(encoded))

			decoded = OneOrMore[string]{}
			err = decoded.UnmarshalJSON(encoded)
			assert.NoError(t, err)
			assert.Equal(t, tc.oom, decoded)
		})
	}

	hash1 := comid.NewDigestIntAlg(
		comid.Sha256,
		comid.MustHexDecode(t, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
	)
	hash2 := comid.NewDigestIntAlg(
		comid.Sha256,
		comid.MustHexDecode(t, "c0decafec0decafec0decafec0decafec0decafec0decafec0decafec0decafe"),
	)

	digestCases := []struct {
		title        string
		oom          OneOrMore[comid.Digest]
		expectdCBOR  []byte
		expectedJSON string
	}{
		{
			title: "one digest",
			oom:   OneOrMore[comid.Digest]{hash1},
			expectdCBOR: []byte{
				0x82,       // array(2)
				0x01,       // . [0]1 [sha-256]
				0x58, 0x20, // . [1]bstr(32)
				0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
				0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
				0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
				0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, // 32
			},
			expectedJSON: `[1, "3q2-796tvu_erb7v3q2-796tvu_erb7v3q2-796tvu8"]`,
		},
		{
			title: "more digests",
			oom:   OneOrMore[comid.Digest]{hash1, hash2},
			expectdCBOR: []byte{
				0x82,       // array(2)
				0x82,       // . [0]array(2)
				0x01,       // . . [0]1 [sha-256]
				0x58, 0x20, // . . [1]bstr(32)
				0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
				0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
				0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
				0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, // 32
				0x82,       // . [1]array(2)
				0x01,       // . . [0]1 [sha-256]
				0x58, 0x20, // . . [1]bstr(32)
				0xc0, 0xde, 0xca, 0xfe, 0xc0, 0xde, 0xca, 0xfe,
				0xc0, 0xde, 0xca, 0xfe, 0xc0, 0xde, 0xca, 0xfe,
				0xc0, 0xde, 0xca, 0xfe, 0xc0, 0xde, 0xca, 0xfe,
				0xc0, 0xde, 0xca, 0xfe, 0xc0, 0xde, 0xca, 0xfe, // 32
			},
			expectedJSON: `[[1, "3q2-796tvu_erb7v3q2-796tvu_erb7v3q2-796tvu8"],[1, "wN7K_sDeyv7A3sr-wN7K_sDeyv7A3sr-wN7K_sDeyv4"]]`,
		},
	}

	for _, tc := range digestCases {
		t.Run(tc.title, func(t *testing.T) {
			encoded, err := tc.oom.MarshalCBOR()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectdCBOR, encoded)

			var decoded OneOrMore[comid.Digest]
			err = decoded.UnmarshalCBOR(encoded)
			assert.NoError(t, err)
			assert.Equal(t, tc.oom, decoded)

			encoded, err = tc.oom.MarshalJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(encoded))

			decoded = OneOrMore[comid.Digest]{}
			err = decoded.UnmarshalJSON(encoded)
			assert.NoError(t, err)
			assert.Equal(t, tc.oom, decoded)
		})
	}
}

func Test_OneOrMore_Valid(t *testing.T) {
	oom := OneOrMore[int]{}
	err := oom.Valid()
	assert.ErrorContains(t, err, "must have at least one")

	oom.Add(1)
	err = oom.Valid()
	assert.NoError(t, err)
}
