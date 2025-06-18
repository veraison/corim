// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPSARefValID_Valid_SignerID_range(t *testing.T) {
	signerID := []byte{}

	for i := 1; i <= 100; i++ {
		signerID = append(signerID, byte(0xff))

		tv, err := NewPSARefValID(signerID)

		switch i {
		case 32, 48, 64:
			assert.NotNil(t, tv)
			assert.Nil(t, tv.Valid())
		default:
			assert.Nil(t, tv)
			assert.EqualError(
				t,
				err,
				fmt.Sprintf("invalid PSA RefVal ID length: %d", i),
			)
		}
	}
}

func TestPSARefValID_String(t *testing.T) {
	signerID := MustHexDecode(t, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	refvalID, err := NewTaggedPSARefValID(signerID)
	require.NoError(t, err)

	assert.Equal(t, `{"signer-id":"3q2+796tvu/erb7v3q2+796tvu/erb7v3q2+796tvu8="}`, refvalID.String())
}
