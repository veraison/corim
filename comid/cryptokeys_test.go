// Copyright 2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
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
