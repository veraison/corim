// Copyright 2021-2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevIdentityKey_Valid_empty(t *testing.T) {
	invalidKey := CryptoKey{TaggedPKIXBase64Key("")}

	tvs := []struct {
		env      Environment
		verifkey CryptoKeys
		testerr  string
	}{
		{
			env:      Environment{},
			verifkey: CryptoKeys{},
			testerr:  "environment validation failed: environment must not be empty",
		},
		{
			env:      Environment{Instance: NewInstanceUEID(TestUEID)},
			verifkey: CryptoKeys{},
			testerr:  "verification keys validation failed: no keys to validate",
		},
		{
			env:      Environment{Instance: NewInstanceUEID(TestUEID)},
			verifkey: CryptoKeys{&invalidKey},
			testerr:  "verification keys validation failed: invalid key at index 0: key value not set",
		},
	}
	for _, tv := range tvs {
		av := DevIdentityKey{Environment: tv.env, VerifKeys: tv.verifkey}
		err := av.Valid()
		assert.EqualError(t, err, tv.testerr)
	}
}
