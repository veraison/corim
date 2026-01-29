// Copyright 2021-2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/swid"
)

func TestDigests_AddDigest_OK(t *testing.T) {
	d := NewDigests()
	require.NotNil(t, d)

	tvs := []struct {
		alg uint64
		val []byte
	}{
		{
			alg: swid.Sha256,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"),
		},
		{
			alg: swid.Sha256_128,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f36"),
		},
		{
			alg: swid.Sha256_120,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f"),
		},
		{
			alg: swid.Sha256_96,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3a"),
		},
		{
			alg: swid.Sha256_64,
			val: MustHexDecode(t, "e45b72f5c0c0b572"),
		},
		{
			alg: swid.Sha256_32,
			val: MustHexDecode(t, "e45b72ab"),
		},
		{
			alg: swid.Sha384,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f36"),
		},
		{
			alg: swid.Sha512,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"),
		},
		{
			alg: swid.Sha3_224,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f36e45b72f5c0c0b572db4d8d3a"),
		},
		{
			alg: swid.Sha3_256,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"),
		},
		{
			alg: swid.Sha3_384,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f36"),
		},
		{
			alg: swid.Sha3_512,
			val: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"),
		},
	}

	for _, tv := range tvs {
		assert.NotNil(t, d.AddDigest(tv.alg, tv.val))
		assert.Nil(t, d.Valid())
	}
}
func TestDigests_Valid_empty(t *testing.T) {
	d := NewDigests()
	require.NotNil(t, d)

	// simulate evil CBOR
	*d = append(*d, swid.HashEntry{
		HashAlgID: 666,
		HashValue: []byte{0x66, 0x66, 0x06},
	})

	assert.EqualError(t, d.Valid(), "digest at index 0: unknown hash algorithm 666")
}

func TestDigests_AddDigest_unknown_algo(t *testing.T) {
	d := NewDigests()
	require.NotNil(t, d)

	assert.NotNil(t, d.AddDigest(0, []byte("0 is a reserved value")))

	err := d.Valid()
	assert.ErrorContains(t, err, "digest at index 0: unknown hash algorithm")
}

func TestDigests_AddDigest_inconsistent_length_for_algo(t *testing.T) {
	d := NewDigests()
	require.NotNil(t, d)

	assert.NotNil(t, d.AddDigest(swid.Sha3_512, MustHexDecode(t, "deadbeef")))

	err := d.Valid()
	assert.ErrorContains(t, err, "digest at index 0: length mismatch")
}

func TestDigests_MarshalJSON(t *testing.T) {
	d := NewDigests().
		AddDigest(swid.Sha256_32, MustHexDecode(t, "e45b72ab")).
		AddDigest(swid.Sha256_64, MustHexDecode(t, "e45b72f5c0c0b572"))
	require.NotNil(t, d)

	expected := `[ "sha-256-32;5Ftyqw==", "sha-256-64;5Fty9cDAtXI=" ]`

	actual, err := json.Marshal(d)

	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(actual))
}

func TestDigests_MarshalCBOR(t *testing.T) {
	d := NewDigests().
		AddDigest(swid.Sha256_32, MustHexDecode(t, "e45b72ab")).
		AddDigest(swid.Sha256_64, MustHexDecode(t, "e45b72f5c0c0b572"))
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
	assert.Equal(t, swid.Sha256_32, actual[0].HashAlgID)
	assert.Equal(t, MustHexDecode(t, "e45b72ab"), actual[0].HashValue)
	assert.Equal(t, swid.Sha256_64, actual[1].HashAlgID)
	assert.Equal(t, MustHexDecode(t, "e45b72f5c0c0b572"), actual[1].HashValue)
}

func TestDigests_Equal_True(t *testing.T) {
	ref := NewDigests().
		AddDigest(swid.Sha256_64, MustHexDecode(t, "e45b72f5c0c0b572")).
		AddDigest(swid.Sha256_32, MustHexDecode(t, "e45b72ab")).
		AddDigest(swid.Sha256_64, MustHexDecode(t, "e45b72f5c0c0b572"))

	claim := NewDigests().
		AddDigest(swid.Sha256_64, MustHexDecode(t, "e45b72f5c0c0b572")).
		AddDigest(swid.Sha256_64, MustHexDecode(t, "e45b72f5c0c0b572")).
		AddDigest(swid.Sha256_32, MustHexDecode(t, "e45b72ab"))

	assert.True(t, claim.Equal(*ref))
}

func TestDigests_Equal_False_Length(t *testing.T) {
	ref := NewDigests().
		AddDigest(swid.Sha256_32, MustHexDecode(t, "e45b72ab")).
		AddDigest(swid.Sha256_64, MustHexDecode(t, "e45b72f5c0c0b572"))

	claim := NewDigests().
		AddDigest(swid.Sha256_32, MustHexDecode(t, "e45b72ab"))

	assert.False(t, claim.Equal(*ref))
}

func TestDigests_Equal_False_Mismatch(t *testing.T) {
	ref := NewDigests().
		AddDigest(swid.Sha256_32, MustHexDecode(t, "e45b72ab")).
		AddDigest(swid.Sha256_64, MustHexDecode(t, "e45b72f5c0c0b572"))

	claim := NewDigests().
		AddDigest(swid.Sha256_32, MustHexDecode(t, "e45b72ab")).
		AddDigest(swid.Sha256_64, MustHexDecode(t, "a26c83e2d0c0b572"))

	assert.False(t, claim.Equal(*ref))
}

func TestDigests_Compare_True(t *testing.T) {
	ref := NewDigests().
		AddDigest(swid.Sha256_32, MustHexDecode(t, "e45b72ab")).
		AddDigest(swid.Sha384, MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f36"))

	claim := NewDigests().
		AddDigest(swid.Sha384, MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f36"))

	assert.True(t, claim.CompareAgainstReference(*ref))
}

func TestDigests_Compare_False(t *testing.T) {
	ref := NewDigests().
		AddDigest(swid.Sha256_32, MustHexDecode(t, "e45b72ab")).
		AddDigest(swid.Sha256_64, MustHexDecode(t, "e45b72f5c0c0b572"))

	claim := NewDigests().
		AddDigest(swid.Sha256_32, MustHexDecode(t, "f39a61fe"))

	assert.False(t, claim.CompareAgainstReference(*ref))
}

func TestDigests_Compare_False_DuplicateIDs(t *testing.T) {
	ref := NewDigests().
		AddDigest(swid.Sha256_32, MustHexDecode(t, "e45b72ab")).
		AddDigest(swid.Sha256_32, MustHexDecode(t, "f34a51de"))

	claim := NewDigests().
		AddDigest(swid.Sha256_32, MustHexDecode(t, "e45b72ab"))

	assert.False(t, claim.CompareAgainstReference(*ref))
}

func TestDigests_Compare_False_PartialMatch(t *testing.T) {
	ref := NewDigests().
		AddDigest(swid.Sha256_32, MustHexDecode(t, "e45b72ab")).
		AddDigest(swid.Sha256_64, MustHexDecode(t, "e45b72f5c0c0b572"))

	claim := NewDigests().
		AddDigest(swid.Sha256_32, MustHexDecode(t, "e45b72ab")).
		AddDigest(swid.Sha256_64, MustHexDecode(t, "f39c2473a0c0f592"))

	assert.False(t, claim.CompareAgainstReference(*ref))
}

func TestNewHashEntry(t *testing.T) {
	// Valid hash entry
	he := NewHashEntry(swid.Sha256, MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"))
	assert.NotNil(t, he)

	// Invalid hash entry - wrong length for algorithm
	he = NewHashEntry(swid.Sha256, []byte{0x01, 0x02})
	assert.Nil(t, he)
}
