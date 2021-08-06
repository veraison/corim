// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDigests_AddDigest_OK(t *testing.T) {
	d := NewDigests()
	require.NotNil(t, d)

	tvs := []struct {
		alg uint64
		val []byte
	}{
		{
			alg: Sha256,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"),
		},
		{
			alg: Sha256_128,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f36"),
		},
		{
			alg: Sha256_120,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f"),
		},
		{
			alg: Sha256_96,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3a"),
		},
		{
			alg: Sha256_64,
			val: MustHexDecode(t, "e45b72f5c0c0b572"),
		},
		{
			alg: Sha256_32,
			val: MustHexDecode(t, "e45b72ab"),
		},
		{
			alg: Sha384,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f36"),
		},
		{
			alg: Sha512,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"),
		},
		{
			alg: Sha3_224,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f36e45b72f5c0c0b572db4d8d3a"),
		},
		{
			alg: Sha3_256,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"),
		},
		{
			alg: Sha3_384,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f36"),
		},
		{
			alg: Sha3_512,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"),
		},
	}

	for _, tv := range tvs {
		assert.NotNil(t, d.AddDigest(tv.alg, tv.val))
		assert.Nil(t, d.Valid())
	}
}

func TestDigests_AddDigest_unknown_algo(t *testing.T) {
	d := NewDigests()
	require.NotNil(t, d)

	assert.Nil(t, d.AddDigest(0, []byte("0 is a reserved value")))
}

func TestDigests_AddDigest_inconsistent_length_for_algo(t *testing.T) {
	d := NewDigests()
	require.NotNil(t, d)

	assert.Nil(t, d.AddDigest(Sha3_512, MustHexDecode(t, "deadbeef")))
}

func TestDigests_MarshalJSON(t *testing.T) {
	d := NewDigests().
		AddDigest(Sha256_32, MustHexDecode(t, "e45b72ab")).
		AddDigest(Sha256_64, MustHexDecode(t, "e45b72f5c0c0b572"))
	require.NotNil(t, d)

	expected := `[ "sha-256-32:5Ftyqw==", "sha-256-64:5Fty9cDAtXI=" ]`

	actual, err := json.Marshal(d)

	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(actual))
}

func TestDigests_MarshalCBOR(t *testing.T) {
	d := NewDigests().
		AddDigest(Sha256_32, MustHexDecode(t, "e45b72ab")).
		AddDigest(Sha256_64, MustHexDecode(t, "e45b72f5c0c0b572"))
	require.NotNil(t, d)

	// [[6, h'E45B72AB'], [5, h'E45B72F5C0C0B572']]
	expected := MustHexDecode(t, "82820644e45b72ab820548e45b72f5c0c0b572")

	actual, err := em.Marshal(d)

	fmt.Printf("CBOR: %x\n", actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestDigests_UnmarshalCBOR(t *testing.T) {
	// [[6, h'E45B72AB'], [5, h'E45B72F5C0C0B572']]
	tv := MustHexDecode(t, "82820644e45b72ab820548e45b72f5c0c0b572")

	var actual Digests

	err := dm.Unmarshal(tv, &actual)

	assert.Nil(t, err)
	assert.Equal(t, Sha256_32, actual[0].HashAlgID)
	assert.Equal(t, MustHexDecode(t, "e45b72ab"), actual[0].HashValue)
	assert.Equal(t, Sha256_64, actual[1].HashAlgID)
	assert.Equal(t, MustHexDecode(t, "e45b72f5c0c0b572"), actual[1].HashValue)
}
