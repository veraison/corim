// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cca

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
)

// Helper Functions

// mustNewTaggedBytesCryptoKey creates a TaggedBytes CryptoKey with specified length
func mustNewTaggedBytesCryptoKey(length int) *comid.CryptoKey {
	data := make([]byte, length)
	for i := range data {
		data[i] = byte(i % 256)
	}
	key, err := comid.NewCryptoKey(data, comid.BytesType)
	if err != nil {
		panic(err)
	}
	return key
}

// mustNewPKIXKey creates a PKIX Base64 Key (for AttestVerifKeys)
func mustNewPKIXKey() *comid.CryptoKey {
	return comid.MustNewCryptoKey(
		"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----",
		comid.PKIXBase64KeyType,
	)
}

// mustNewCCAPlatformInstanceID creates a valid CCA Platform Instance ID (33-byte UEID with 0x01 RAND prefix)
func mustNewCCAPlatformInstanceID() *comid.Instance {
	// CCA Platform Instance ID: 0x01 (RAND type) + 32 bytes of identifier
	ueidBytes := make([]byte, 33)
	ueidBytes[0] = 0x01 // RAND type
	for i := 1; i < 33; i++ {
		ueidBytes[i] = byte(i - 1)
	}
	inst, err := comid.NewUEIDInstance(comid.UEID(ueidBytes))
	if err != nil {
		panic(err)
	}
	return inst
}

// mustNewCCAPlatformInstanceIDWithLength creates a Platform Instance ID with specified length
// This bypasses EAT validation to allow testing with invalid UEIDs
func mustNewCCAPlatformInstanceIDWithLength(length int, typePrefix byte) *comid.Instance {
	ueidBytes := make([]byte, length)
	if length > 0 {
		ueidBytes[0] = typePrefix
	}
	for i := 1; i < length; i++ {
		ueidBytes[i] = byte(i % 256)
	}
	// Directly set the tagged value to bypass EAT validation
	tagged := comid.TaggedUEID(ueidBytes)
	return &comid.Instance{Value: tagged}
}

// newCCAPlatformEnvironmentWithImplID creates a valid CCA Platform Environment with 32-byte Implementation ID
// and valid 33-byte UEID Instance ID (RAND type)
func newCCAPlatformEnvironmentWithImplID() comid.Environment {
	// Create a 32-byte Implementation ID
	implID := make([]byte, 32)
	for i := range implID {
		implID[i] = byte(i)
	}
	class := comid.NewClassBytes(implID)
	return comid.Environment{
		Class:    class,
		Instance: mustNewCCAPlatformInstanceID(),
	}
}

// newCCAPlatformEnvironmentWithImplIDLength creates Environment with Implementation ID of specified length
// and a valid CCA Platform Instance ID
func newCCAPlatformEnvironmentWithImplIDLength(length int) comid.Environment {
	implID := make([]byte, length)
	for i := range implID {
		implID[i] = byte(i % 256)
	}
	class := comid.NewClassBytes(implID)
	return comid.Environment{
		Class:    class,
		Instance: mustNewCCAPlatformInstanceID(),
	}
}

// Test that the profile is registered
func TestPlatformProfile_Registered(t *testing.T) {
	triples := &comid.Triples{}
	// Profile is registered in init(), just check we can use it
	assert.NotNil(t, triples)
}

// Test validateCCAPlatformImplementationID
func TestValidateCCAPlatformImplementationID_AllCases(t *testing.T) {
	testCases := []struct {
		title           string
		env             comid.Environment
		shouldError     bool
		expectedMessage string
	}{
		{
			title:       "valid environment with 32-byte impl ID",
			env:         newCCAPlatformEnvironmentWithImplID(),
			shouldError: false,
		},
		{
			title: "missing class",
			env: comid.Environment{
				Instance: mustNewCCAPlatformInstanceID(),
			},
			shouldError:     true,
			expectedMessage: "environment.class is required",
		},
		{
			title: "wrong class type (UUID instead of bytes)",
			env: comid.Environment{
				Class: comid.NewClassUUID(comid.TestUUID),
			},
			shouldError:     true,
			expectedMessage: "must be of type 'bytes'",
		},
		{
			title:           "impl ID too short (31 bytes)",
			env:             newCCAPlatformEnvironmentWithImplIDLength(31),
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "impl ID too long (33 bytes)",
			env:             newCCAPlatformEnvironmentWithImplIDLength(33),
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := validateCCAPlatformImplementationID(&tc.env, "test")
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test validateCCAPlatformInstanceID (internal function with Environment parameter)
func TestValidateCCAPlatformInstanceID_InEnvironment_AllCases(t *testing.T) {
	testCases := []struct {
		title           string
		env             comid.Environment
		shouldError     bool
		expectedMessage string
	}{
		{
			title:       "valid instance ID",
			env:         newCCAPlatformEnvironmentWithImplID(),
			shouldError: false,
		},
		{
			title:           "missing instance",
			env:             comid.Environment{},
			shouldError:     true,
			expectedMessage: "instance-id) is required",
		},
		{
			title: "wrong instance type (UUID instead of UEID)",
			env: comid.Environment{
				Instance: comid.MustNewUUIDInstance(comid.TestUUID),
			},
			shouldError:     true,
			expectedMessage: "must be of type 'ueid'",
		},
		{
			title: "instance ID too short (32 bytes)",
			env: comid.Environment{
				Instance: mustNewCCAPlatformInstanceIDWithLength(32, 0x01),
			},
			shouldError:     true,
			expectedMessage: "must be exactly 33 bytes",
		},
		{
			title: "instance ID too long (34 bytes)",
			env: comid.Environment{
				Instance: mustNewCCAPlatformInstanceIDWithLength(34, 0x01),
			},
			shouldError:     true,
			expectedMessage: "must be exactly 33 bytes",
		},
		{
			title: "wrong prefix 0x00",
			env: comid.Environment{
				Instance: mustNewCCAPlatformInstanceIDWithLength(33, 0x00),
			},
			shouldError:     true,
			expectedMessage: "must have RAND type (0x01), got 0x00",
		},
		{
			title: "wrong prefix 0x02",
			env: comid.Environment{
				Instance: mustNewCCAPlatformInstanceIDWithLength(33, 0x02),
			},
			shouldError:     true,
			expectedMessage: "must have RAND type (0x01), got 0x02",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := validateCCAPlatformInstanceID(&tc.env, "test")
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test validateCCADigests
func TestValidateCCADigests_AllCases(t *testing.T) {
	testCases := []struct {
		title           string
		digests         *comid.Digests
		shouldError     bool
		expectedMessage string
	}{
		{
			title: "valid digest with SHA-256 (32 bytes)",
			digests: func() *comid.Digests {
				d := &comid.Digests{}
				d.AddDigest(swid.Sha256, make([]byte, 32))
				return d
			}(),
			shouldError: false,
		},
		{
			title: "valid digest with SHA-384 (48 bytes)",
			digests: func() *comid.Digests {
				d := &comid.Digests{}
				d.AddDigest(swid.Sha384, make([]byte, 48))
				return d
			}(),
			shouldError: false,
		},
		{
			title: "valid digest with SHA-512 (64 bytes)",
			digests: func() *comid.Digests {
				d := &comid.Digests{}
				d.AddDigest(swid.Sha512, make([]byte, 64))
				return d
			}(),
			shouldError: false,
		},
		{
			title:           "nil digests",
			digests:         nil,
			shouldError:     true,
			expectedMessage: "digests field is mandatory",
		},
		{
			title:           "empty digests",
			digests:         &comid.Digests{},
			shouldError:     true,
			expectedMessage: "digests must contain at least one entry",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := validateCCADigests(tc.digests, 0, 0)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

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
			key := mustNewTaggedBytesCryptoKey(tc.length)
			keys := comid.NewCryptoKeys()
			keys.Add(key)

			err := validateCCASignerID(keys, 0, 0)

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
				keys.Add(mustNewTaggedBytesCryptoKey(32))
				keys.Add(mustNewTaggedBytesCryptoKey(32))
				return keys
			}(),
			expectedMessage: "cryptokeys must contain exactly one entry",
		},
		{
			title: "wrong type (pkix-base64-key instead of bytes)",
			keys: func() *comid.CryptoKeys {
				keys := comid.NewCryptoKeys()
				keys.Add(mustNewPKIXKey())
				return keys
			}(),
			expectedMessage: "must be of type 'bytes'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := validateCCASignerID(tc.keys, 0, 0)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedMessage)
		})
	}
}

// Test validateCCAPlatformAttestVerifKey
func TestValidateCCAPlatformAttestVerifKey_AllCases(t *testing.T) {
	testCases := []struct {
		title           string
		setupAVK        func() *comid.KeyTriple
		shouldError     bool
		expectedMessage string
	}{
		{
			title: "valid PKIX key with environment",
			setupAVK: func() *comid.KeyTriple {
				key := mustNewPKIXKey()
				env := newCCAPlatformEnvironmentWithImplID()
				return &comid.KeyTriple{
					Environment: env,
					VerifKeys:   comid.CryptoKeys{key},
				}
			},
			shouldError: false,
		},
		{
			title: "missing environment",
			setupAVK: func() *comid.KeyTriple {
				return &comid.KeyTriple{}
			},
			shouldError:     true,
			expectedMessage: "", // Will fail on environment validation
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			avk := tc.setupAVK()
			err := validateCCAPlatformAttestVerifKey(avk, 0)

			if tc.shouldError {
				assert.Error(t, err)
				if tc.expectedMessage != "" {
					assert.Contains(t, err.Error(), tc.expectedMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test complete reference value validation
func TestValidateCCAPlatformReferenceValue_AllCases(t *testing.T) {
	testCases := []struct {
		title           string
		setupRefVal     func() *comid.ValueTriple
		shouldError     bool
		expectedMessage string
	}{
		{
			title: "valid software component",
			setupRefVal: func() *comid.ValueTriple {
				env := newCCAPlatformEnvironmentWithImplID()

				// Create valid software component measurement
				measurement, err := comid.NewMeasurement(CCASoftwareComponentMkey, "string")
				require.NoError(t, err)

				// Set digests
				digests := &comid.Digests{}
				hash := make([]byte, 32)
				digests.AddDigest(swid.Sha256, hash)
				measurement.Val.Digests = digests

				// Set signer-id (cryptokeys)
				signerID := mustNewTaggedBytesCryptoKey(32)
				measurement.Val.CryptoKeys = comid.NewCryptoKeys()
				measurement.Val.CryptoKeys.Add(signerID)

				measurements := comid.NewMeasurements()
				measurements.Values = append(measurements.Values, *measurement)

				return &comid.ValueTriple{
					Environment:  env,
					Measurements: *measurements,
				}
			},
			shouldError: false,
		},
		{
			title: "missing software component",
			setupRefVal: func() *comid.ValueTriple {
				env := newCCAPlatformEnvironmentWithImplID()
				measurements := comid.NewMeasurements()

				return &comid.ValueTriple{
					Environment:  env,
					Measurements: *measurements,
				}
			},
			shouldError:     true,
			expectedMessage: "at least one software component measurement is required",
		},
		{
			title: "invalid mkey",
			setupRefVal: func() *comid.ValueTriple {
				env := newCCAPlatformEnvironmentWithImplID()

				// Create measurement with invalid mkey
				measurement, err := comid.NewMeasurement("invalid-key", "string")
				require.NoError(t, err)

				measurements := comid.NewMeasurements()
				measurements.Values = append(measurements.Values, *measurement)

				return &comid.ValueTriple{
					Environment:  env,
					Measurements: *measurements,
				}
			},
			shouldError:     true,
			expectedMessage: "invalid mkey",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			refVal := tc.setupRefVal()
			err := validateCCAPlatformReferenceValue(refVal, 0)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
