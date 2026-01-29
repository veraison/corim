// Copyright 2023-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CryptoKeys(t *testing.T) {
	keys := NewCryptoKeys()
	err := keys.Valid()

	assert.EqualError(t, err, "no keys to validate")

	keys.Add(MustNewCOSEKey(TestCOSEKey)).
		Add(MustNewPKIXBase64Key(TestECPubKey)).
		Add(nil)

	err = keys.Valid()
	assert.NoError(t, err)

	badKey := CryptoKey{TaggedPKIXBase64Cert("lol, nope!")}
	keys.Add(&badKey)

	err = keys.Valid()
	assert.ErrorContains(t, err, "invalid key at index 2")
}

func Test_CryptoKeys_String(t *testing.T) {
	// Test empty keys
	keys := NewCryptoKeys()
	assert.Equal(t, "[]", keys.String())

	// Test with keys
	keys.Add(MustNewCOSEKey(TestCOSEKey)).
		Add(MustNewPKIXBase64Key(TestECPubKey))

	expected := "[" + base64.StdEncoding.EncodeToString(TestCOSEKey) + ", " + TestECPubKey + "]"
	assert.Equal(t, expected, keys.String())
}
