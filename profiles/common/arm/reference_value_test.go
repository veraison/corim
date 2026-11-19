package arm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/profiles/common/arm/testhelpers"
)

// Test validateCCASignerID
func TestValidateCCASignerID_ValidAndInvalidLengths(t *testing.T) {
	testCases := []struct {
		title           string
		length          int
		shouldError     bool
		expectedMessage string
	}{
		{
			title:       "valid 32 bytes",
			length:      32,
			shouldError: false,
		},
		{
			title:       "valid 48 bytes",
			length:      48,
			shouldError: false,
		},
		{
			title:       "valid 64 bytes",
			length:      64,
			shouldError: false,
		},
		{
			title:           "invalid 0 bytes",
			length:          0,
			shouldError:     true,
			expectedMessage: "must be 32, 48, or 64 bytes",
		},
		{
			title:           "invalid 16 bytes",
			length:          16,
			shouldError:     true,
			expectedMessage: "must be 32, 48, or 64 bytes",
		},
		{
			title:           "invalid 31 bytes",
			length:          31,
			shouldError:     true,
			expectedMessage: "must be 32, 48, or 64 bytes",
		},
		{
			title:           "invalid 33 bytes",
			length:          33,
			shouldError:     true,
			expectedMessage: "must be 32, 48, or 64 bytes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			key := testhelpers.MustNewTaggedBytesCryptoKey(tc.length)
			keys := comid.NewCryptoKeys()
			keys.Add(key)

			err := ValidateSignerID(keys, 0, 0)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCCASignerID_InvalidCases(t *testing.T) {
	testCases := []struct {
		title           string
		keys            *comid.CryptoKeys
		expectedMessage string
	}{
		{
			title:           "nil keys",
			keys:            nil,
			expectedMessage: "cryptokeys (signer-id) is mandatory but not set",
		},
		{
			title:           "empty cryptokeys",
			keys:            comid.NewCryptoKeys(),
			expectedMessage: "cryptokeys must contain exactly one entry",
		},
		{
			title: "multiple entries",
			keys: func() *comid.CryptoKeys {
				keys := comid.NewCryptoKeys()
				keys.Add(testhelpers.MustNewTaggedBytesCryptoKey(32))
				keys.Add(testhelpers.MustNewTaggedBytesCryptoKey(32))
				return keys
			}(),
			expectedMessage: "cryptokeys must contain exactly one entry",
		},
		{
			title: "wrong type (pkix-base64-key instead of bytes)",
			keys: func() *comid.CryptoKeys {
				keys := comid.NewCryptoKeys()
				keys.Add(testhelpers.MustNewPKIXKey())
				return keys
			}(),
			expectedMessage: "must be of type 'bytes'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := ValidateSignerID(tc.keys, 0, 0)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedMessage)
		})
	}
}
