// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPSARefValID_Valid_SignerID_range(t *testing.T) {
	signerID := []byte{}

	for i := 1; i <= 100; i++ {
		signerID = append(signerID, byte(0xff))

		tv := NewPSARefValID(signerID)
		switch i {
		case 32, 48, 64:
			assert.NotNil(t, tv)
			assert.Nil(t, tv.Valid())
		default:
			assert.Nil(t, tv)
		}
	}
}
