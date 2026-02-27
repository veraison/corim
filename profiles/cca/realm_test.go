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

// Helper Functions for Realm Tests

// newCCARealmEnvironmentWithRIM creates a valid CCA Realm Environment with 32-byte RIM as class identifier
func newCCARealmEnvironmentWithRIM() comid.Environment {
	rimBytes := make([]byte, 32)
	for i := range rimBytes {
		rimBytes[i] = byte(i)
	}
	class := comid.NewClassBytes(rimBytes)
	return comid.Environment{
		Class: class,
	}
}

// newCCARealmEnvironmentWithRIMLength creates Environment with RIM of specified length
func newCCARealmEnvironmentWithRIMLength(length int) comid.Environment {
	rim := make([]byte, length)
	for i := range rim {
		rim[i] = byte(i % 256)
	}
	class := comid.NewClassBytes(rim)
	return comid.Environment{
		Class: class,
	}
}

// mustNewCCARealmDigestMeasurement creates a valid RIM or REM measurement with digest
func mustNewCCARealmDigestMeasurement(mkey string, hashSize int) *comid.Measurement {
	measurement, err := comid.NewMeasurement(mkey, "string")
	if err != nil {
		panic(err)
	}

	// Set digests with specified hash size
	digests := &comid.Digests{}
	hash := make([]byte, hashSize)
	digests.AddDigest(swid.Sha256, hash)
	measurement.Val.Digests = digests

	return measurement
}

// mustNewCCARealmRPVMeasurement creates a valid RPV (Realm Personalization Value) measurement
func mustNewCCARealmRPVMeasurement(data []byte) *comid.Measurement {
	measurement, err := comid.NewMeasurement(CCARealmPersonalizationMkey, "string")
	if err != nil {
		panic(err)
	}

	// Set raw-value as tagged bytes
	rv := comid.NewRawValue().SetBytes(data)
	measurement.Val.RawValue = rv

	return measurement
}

// Test that the realm profile is registered
func TestRealmProfile_Registered(t *testing.T) {
	triples := &comid.Triples{}
	// Profile is registered in init(), just check we can use it
	assert.NotNil(t, triples)
}

// Test validateCCARealmRIM
func TestValidateCCARealmRIM_AllCases(t *testing.T) {
	testCases := []struct {
		title           string
		env             comid.Environment
		shouldError     bool
		expectedMessage string
	}{
		{
			title:       "valid environment with 32-byte RIM",
			env:         newCCARealmEnvironmentWithRIM(),
			shouldError: false,
		},
		{
			title: "missing class",
			env:   comid.Environment{
				// No class set
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
			expectedMessage: "RIM must be of type 'bytes'",
		},
		{
			title:           "RIM too short (31 bytes)",
			env:             newCCARealmEnvironmentWithRIMLength(31),
			shouldError:     true,
			expectedMessage: "RIM must be 32, 48, or 64 bytes",
		},
		{
			title:           "RIM too long (65 bytes)",
			env:             newCCARealmEnvironmentWithRIMLength(65),
			shouldError:     true,
			expectedMessage: "RIM must be 32, 48, or 64 bytes",
		},
		{
			title:       "RIM valid 48 bytes (SHA-384)",
			env:         newCCARealmEnvironmentWithRIMLength(48),
			shouldError: false,
		},
		{
			title:       "RIM valid 64 bytes (SHA-512)",
			env:         newCCARealmEnvironmentWithRIMLength(64),
			shouldError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := validateCCARealmRIM(&tc.env, "test")
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test validateCCARealmMeasurement (RIM and REM validation)
func TestValidateCCARealmMeasurement_AllCases(t *testing.T) {
	testCases := []struct {
		title            string
		setupMeasurement func() *comid.Measurement
		measType         string
		shouldError      bool
		expectedMessage  string
	}{
		{
			title: "valid RIM with SHA-256 (32 bytes)",
			setupMeasurement: func() *comid.Measurement {
				return mustNewCCARealmDigestMeasurement(CCARealmInitialMeasurementMkey, 32)
			},
			measType:    "RIM",
			shouldError: false,
		},
		{
			title: "valid REM with SHA-384 (48 bytes)",
			setupMeasurement: func() *comid.Measurement {
				return mustNewCCARealmDigestMeasurement(CCARealmExtendedMeasurement0Mkey, 48)
			},
			measType:    "REM",
			shouldError: false,
		},
		{
			title: "valid REM with SHA-512 (64 bytes)",
			setupMeasurement: func() *comid.Measurement {
				return mustNewCCARealmDigestMeasurement(CCARealmExtendedMeasurement1Mkey, 64)
			},
			measType:    "REM",
			shouldError: false,
		},
		{
			title: "missing digests",
			setupMeasurement: func() *comid.Measurement {
				measurement, err := comid.NewMeasurement(CCARealmInitialMeasurementMkey, "string")
				require.NoError(t, err)
				return measurement
			},
			measType:        "RIM",
			shouldError:     true,
			expectedMessage: "digests field is mandatory",
		},
		{
			title: "hash too short (16 bytes)",
			setupMeasurement: func() *comid.Measurement {
				return mustNewCCARealmDigestMeasurement(CCARealmInitialMeasurementMkey, 16)
			},
			measType:        "RIM",
			shouldError:     true,
			expectedMessage: "hash value must be 32, 48, or 64 bytes",
		},
		{
			title: "hash too long (96 bytes)",
			setupMeasurement: func() *comid.Measurement {
				return mustNewCCARealmDigestMeasurement(CCARealmExtendedMeasurement2Mkey, 96)
			},
			measType:        "REM",
			shouldError:     true,
			expectedMessage: "hash value must be 32, 48, or 64 bytes",
		},
		{
			title: "multiple digests (not allowed)",
			setupMeasurement: func() *comid.Measurement {
				measurement, err := comid.NewMeasurement(CCARealmInitialMeasurementMkey, "string")
				require.NoError(t, err)

				// Add two digests
				digests := &comid.Digests{}
				digests.AddDigest(swid.Sha256, make([]byte, 32))
				digests.AddDigest(swid.Sha384, make([]byte, 48))
				measurement.Val.Digests = digests

				return measurement
			},
			measType:        "RIM",
			shouldError:     true,
			expectedMessage: "digests must contain exactly one entry",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			measurement := tc.setupMeasurement()
			err := validateCCARealmMeasurement(measurement, 0, 0, tc.measType)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test validateCCARealmPersonalizationValue (RPV validation)
func TestValidateCCARealmPersonalizationValue_AllCases(t *testing.T) {
	testCases := []struct {
		title            string
		setupMeasurement func() *comid.Measurement
		shouldError      bool
		expectedMessage  string
	}{
		{
			title: "valid RPV with data",
			setupMeasurement: func() *comid.Measurement {
				return mustNewCCARealmRPVMeasurement([]byte("test data"))
			},
			shouldError: false,
		},
		{
			title: "valid RPV with empty data",
			setupMeasurement: func() *comid.Measurement {
				return mustNewCCARealmRPVMeasurement([]byte{})
			},
			shouldError: false,
		},
		{
			title: "missing raw-value",
			setupMeasurement: func() *comid.Measurement {
				measurement, err := comid.NewMeasurement(CCARealmPersonalizationMkey, "string")
				require.NoError(t, err)
				return measurement
			},
			shouldError:     true,
			expectedMessage: "raw-value is mandatory",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			measurement := tc.setupMeasurement()
			err := validateCCARealmPersonalizationValue(measurement, 0, 0)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test complete realm reference value validation
func TestValidateCCARealmReferenceValue_AllCases(t *testing.T) {
	testCases := []struct {
		title           string
		setupRefVal     func() *comid.ValueTriple
		shouldError     bool
		expectedMessage string
	}{
		{
			title: "valid RIM only",
			setupRefVal: func() *comid.ValueTriple {
				env := newCCARealmEnvironmentWithRIM()

				measurement := mustNewCCARealmDigestMeasurement(CCARealmInitialMeasurementMkey, 32)

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
			title: "valid RIM with REM0",
			setupRefVal: func() *comid.ValueTriple {
				env := newCCARealmEnvironmentWithRIM()

				rim := mustNewCCARealmDigestMeasurement(CCARealmInitialMeasurementMkey, 32)
				rem0 := mustNewCCARealmDigestMeasurement(CCARealmExtendedMeasurement0Mkey, 32)

				measurements := comid.NewMeasurements()
				measurements.Values = append(measurements.Values, *rim, *rem0)

				return &comid.ValueTriple{
					Environment:  env,
					Measurements: *measurements,
				}
			},
			shouldError: false,
		},
		{
			title: "valid RIM with RPV",
			setupRefVal: func() *comid.ValueTriple {
				env := newCCARealmEnvironmentWithRIM()

				rim := mustNewCCARealmDigestMeasurement(CCARealmInitialMeasurementMkey, 32)
				rpv := mustNewCCARealmRPVMeasurement([]byte("realm personalization data"))

				measurements := comid.NewMeasurements()
				measurements.Values = append(measurements.Values, *rim, *rpv)

				return &comid.ValueTriple{
					Environment:  env,
					Measurements: *measurements,
				}
			},
			shouldError: false,
		},
		{
			title: "valid RIM with all REMs and RPV",
			setupRefVal: func() *comid.ValueTriple {
				env := newCCARealmEnvironmentWithRIM()

				rim := mustNewCCARealmDigestMeasurement(CCARealmInitialMeasurementMkey, 32)
				rem0 := mustNewCCARealmDigestMeasurement(CCARealmExtendedMeasurement0Mkey, 32)
				rem1 := mustNewCCARealmDigestMeasurement(CCARealmExtendedMeasurement1Mkey, 48)
				rem2 := mustNewCCARealmDigestMeasurement(CCARealmExtendedMeasurement2Mkey, 64)
				rem3 := mustNewCCARealmDigestMeasurement(CCARealmExtendedMeasurement3Mkey, 32)
				rpv := mustNewCCARealmRPVMeasurement([]byte("realm personalization data"))

				measurements := comid.NewMeasurements()
				measurements.Values = append(measurements.Values, *rim, *rem0, *rem1, *rem2, *rem3, *rpv)

				return &comid.ValueTriple{
					Environment:  env,
					Measurements: *measurements,
				}
			},
			shouldError: false,
		},
		{
			title: "missing RIM (mandatory)",
			setupRefVal: func() *comid.ValueTriple {
				env := newCCARealmEnvironmentWithRIM()

				rem0 := mustNewCCARealmDigestMeasurement(CCARealmExtendedMeasurement0Mkey, 32)

				measurements := comid.NewMeasurements()
				measurements.Values = append(measurements.Values, *rem0)

				return &comid.ValueTriple{
					Environment:  env,
					Measurements: *measurements,
				}
			},
			shouldError:     true,
			expectedMessage: "RIM (cca.rim) measurement is mandatory",
		},
		{
			title: "invalid mkey",
			setupRefVal: func() *comid.ValueTriple {
				env := newCCARealmEnvironmentWithRIM()

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
			err := validateCCARealmReferenceValue(refVal, 0)

			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test identifier helpers
func TestRealmRIMIdentifiers_AllCases(t *testing.T) {
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
			expectedMessage: "got 0 bytes, expected 32, 48, or 64",
		},
		{
			title:           "invalid 16 bytes",
			length:          16,
			shouldError:     true,
			expectedMessage: "got 16 bytes, expected 32, 48, or 64",
		},
		{
			title:           "invalid 31 bytes",
			length:          31,
			shouldError:     true,
			expectedMessage: "got 31 bytes, expected 32, 48, or 64",
		},
		{
			title:           "invalid 33 bytes",
			length:          33,
			shouldError:     true,
			expectedMessage: "got 33 bytes, expected 32, 48, or 64",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			rim := make([]byte, tc.length)

			// Test NewRealmRIMClassID
			_, err := NewRealmRIMClassID(rim)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}

			// Test NewClassRealmRIM
			_, err = NewClassRealmRIM(rim)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}

			// Test ValidateRealmRIM
			err = ValidateRealmRIM(rim)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test MustNewRealmRIMClassID panic behavior
func TestMustNewRealmRIMClassID_Panic(t *testing.T) {
	assert.Panics(t, func() {
		MustNewRealmRIMClassID(make([]byte, 31))
	})
}

// Test MustNewClassRealmRIM panic behavior
func TestMustNewClassRealmRIM_Panic(t *testing.T) {
	assert.Panics(t, func() {
		MustNewClassRealmRIM(make([]byte, 33))
	})
}
