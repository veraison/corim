// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TaggedBytes_ValidatePSASignerID_Valid(t *testing.T) {
	testCases := []struct {
		name   string
		length int
		desc   string
	}{
		{"SHA-256", 32, "32 bytes"},
		{"SHA-384", 48, "48 bytes"},
		{"SHA-512", 64, "64 bytes"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			signerID := make([]byte, tc.length)
			for i := range signerID {
				signerID[i] = byte(i % 256)
			}

			tb := TaggedBytes(signerID)
			err := tb.ValidatePSASignerID()
			assert.NoError(t, err, "Expected %s to be valid", tc.desc)
		})
	}
}

func Test_TaggedBytes_ValidatePSASignerID_Invalid(t *testing.T) {
	invalidLengths := []int{0, 1, 16, 31, 33, 47, 49, 63, 65, 128}

	for _, length := range invalidLengths {
		t.Run(fmt.Sprintf("%d_bytes", length), func(t *testing.T) {
			signerID := make([]byte, length)
			tb := TaggedBytes(signerID)
			err := tb.ValidatePSASignerID()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "must be 32, 48, or 64 bytes")
		})
	}
}
