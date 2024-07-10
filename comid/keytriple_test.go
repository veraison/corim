// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerificationKeys_Valid_empty(t *testing.T) {
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
			env:      Environment{Instance: MustNewUEIDInstance(TestUEID)},
			verifkey: CryptoKeys{},
			testerr:  "verification keys validation failed: no keys to validate",
		},
		{
			env:      Environment{Instance: MustNewUEIDInstance(TestUEID)},
			verifkey: CryptoKeys{&invalidKey},
			testerr:  "verification keys validation failed: invalid key at index 0: key value not set",
		},
	}
	for _, tv := range tvs {
		av := KeyTriple{Environment: tv.env, VerifKeys: tv.verifkey}
		err := av.Valid()
		assert.EqualError(t, err, tv.testerr)
	}
}
