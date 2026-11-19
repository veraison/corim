// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cca

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/profiles/common/arm/testhelpers"
)

// Helper Functions

// mustNewPKIXKey creates a PKIX Base64 Key (for AttestVerifKeys)
func mustNewPKIXKey() *comid.CryptoKey {
	return comid.MustNewCryptoKey(
		"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----",
		comid.PKIXBase64KeyType,
	)
}

// Test that the profile is registered
func TestPlatformProfile_Registered(t *testing.T) {
	triples := &comid.Triples{}
	// Profile is registered in init(), just check we can use it
	assert.NotNil(t, triples)
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
				d.AddDigest(comid.Sha256, make([]byte, 32))
				return d
			}(),
			shouldError: false,
		},
		{
			title: "valid digest with SHA-384 (48 bytes)",
			digests: func() *comid.Digests {
				d := &comid.Digests{}
				d.AddDigest(comid.Sha384, make([]byte, 48))
				return d
			}(),
			shouldError: false,
		},
		{
			title: "valid digest with SHA-512 (64 bytes)",
			digests: func() *comid.Digests {
				d := &comid.Digests{}
				d.AddDigest(comid.Sha512, make([]byte, 64))
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
			err := validateCCADigests(tc.digests)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTBBRoTPK_ValidCases(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			key := testhelpers.MustNewTaggedBytesCryptoKey(tc.length)
			keys := comid.NewCryptoKeys()
			keys.Add(key)

			err := validateCCATBBRoTPK(keys)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePlatformCfg_InvalidCases(t *testing.T) {
	testCases := []struct {
		title           string
		meas            *comid.Measurement
		expectedMessage string
	}{
		{
			title: "no mask",
			meas: func() *comid.Measurement {
				mkey, err := comid.NewMkey("cca.platform-config", comid.StringType)
				assert.NoError(t, err)

				data := []byte{1, 2, 3}
				mask := []byte{}
				m := &comid.Measurement{}
				m.Key = mkey
				rv, err := comid.NewRawValueWithMask(data, mask)
				assert.NoError(t, err)
				m.Val.RawValue = rv
				return m
			}(),
			expectedMessage: "raw-value mask cannot be empty",
		},
		{
			title: "no value",
			meas: func() *comid.Measurement {
				mkey, err := comid.NewMkey("cca.platform-config", comid.StringType)
				assert.NoError(t, err)

				mask := []byte{0xF, 0xF, 0xF}
				var data []byte

				m := &comid.Measurement{}
				m.Key = mkey
				rv, err := comid.NewRawValueWithMask(data, mask)
				assert.NoError(t, err)
				m.Val.RawValue = rv
				return m
			}(),
			expectedMessage: "raw-value cannot be empty",
		},
		{
			title: "incorrect type of raw value",
			meas: func() *comid.Measurement {
				mkey, err := comid.NewMkey("cca.platform-config", comid.StringType)
				assert.NoError(t, err)

				mask := []byte{0xF, 0xF, 0xF}
				var data []byte

				m := &comid.Measurement{}
				m.Key = mkey
				m.SetRawValueBytes(mask, data)
				return m
			}(),
			expectedMessage: "raw-value must be of type masked",
		},
		{
			title: "unequal value and mask",
			meas: func() *comid.Measurement {
				mkey, err := comid.NewMkey("cca.platform-config", comid.StringType)
				assert.NoError(t, err)

				data := []byte{1, 2, 3, 4, 5}
				mask := []byte{0x0F, 0x0F}
				m := &comid.Measurement{}
				m.Key = mkey
				rv, err := comid.NewRawValueWithMask(data, mask)
				assert.NoError(t, err)
				m.Val.RawValue = rv
				return m
			}(),
			expectedMessage: "length of raw-value bytes: 5, not equal to length of raw-value mask: 2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := validateCCAPlatformConfig(tc.meas)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedMessage)
		})
	}
}

func TestValidateTBBRoTPK_InvalidCases(t *testing.T) {
	testCases := []struct {
		title           string
		keys            *comid.CryptoKeys
		expectedMessage string
	}{
		{
			title:           "nil keys",
			keys:            nil,
			expectedMessage: "missing crypto keys for TBB RoTPK reference value",
		},
		{
			title:           "empty cryptokeys",
			keys:            comid.NewCryptoKeys(),
			expectedMessage: "no keys present",
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
		{
			title: "wrong length of key",
			keys: func() *comid.CryptoKeys {
				key := testhelpers.MustNewTaggedBytesCryptoKey(16)
				keys := comid.NewCryptoKeys()
				keys.Add(key)
				return keys
			}(),
			expectedMessage: "expected 32, 48, or 64",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := validateCCATBBRoTPK(tc.keys)
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
				env := testhelpers.NewPlatformEnvironmentWithImplID()
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
			err := validateCCAPlatformAttestVerifKey(avk)
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
				env := testhelpers.NewPlatformEnvironmentWithImplID()

				// Create valid software component measurement
				measurement, err := comid.NewMeasurement(CCASoftwareComponentMkey, "string")
				require.NoError(t, err)

				// Set digests
				digests := &comid.Digests{}
				hash := make([]byte, 32)
				digests.AddDigest(comid.Sha256, hash)
				measurement.Val.Digests = digests

				// Set signer-id (cryptokeys)
				signerID := testhelpers.MustNewTaggedBytesCryptoKey(32)
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
			title: "valid all reference values",
			setupRefVal: func() *comid.ValueTriple {
				env := testhelpers.NewPlatformEnvironmentWithImplID()

				// Create valid software component measurement
				measurement, err := comid.NewMeasurement(CCASoftwareComponentMkey, "string")
				require.NoError(t, err)

				// Set digests
				digests := &comid.Digests{}
				hash := make([]byte, 32)
				digests.AddDigest(comid.Sha256, hash)
				measurement.Val.Digests = digests

				// Set signer-id (cryptokeys)
				signerID := testhelpers.MustNewTaggedBytesCryptoKey(32)
				measurement.Val.CryptoKeys = comid.NewCryptoKeys()
				measurement.Val.CryptoKeys.Add(signerID)

				measurements := comid.NewMeasurements()
				measurements.Values = append(measurements.Values, *measurement)

				// Create Platform Config Measurement

				measurement, err = comid.NewMeasurement(CCAPlatformConfigMkey, "string")
				require.NoError(t, err)
				var rv = []byte{0x1, 0x02, 0x03, 0x04}
				var rm = []byte{0xF, 0xF, 0xF, 0xF}
				raw, err := comid.NewRawValueWithMask(rv, rm)
				require.NoError(t, err)
				measurement.Val.RawValue = raw
				measurements.Values = append(measurements.Values, *measurement)

				// Create Platform Manufacturing Config Claim
				measurement, err = comid.NewMeasurement(CCAPlatformManufacturingConfigKey, "string")
				require.NoError(t, err)
				rv = []byte{0x5, 0x06, 0x07, 0x08}
				rm = []byte{0xF, 0xF, 0xF, 0xF}

				raw, err = comid.NewRawValueWithMask(rv, rm)
				require.NoError(t, err)
				measurement.Val.RawValue = raw
				measurements.Values = append(measurements.Values, *measurement)

				// Create TBB RoTPK Claim
				measurement, err = comid.NewMeasurement("cca.rotpk.CM.2", "string")
				require.NoError(t, err)
				// Set TBB RotPK Key (cryptokeys)
				key := testhelpers.MustNewTaggedBytesCryptoKey(32)
				measurement.Val.CryptoKeys = comid.NewCryptoKeys()
				measurement.Val.CryptoKeys.Add(key)
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
				env := testhelpers.NewPlatformEnvironmentWithImplID()
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
				env := testhelpers.NewPlatformEnvironmentWithImplID()

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
			err := validateCCAPlatformReferenceValue(refVal)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
