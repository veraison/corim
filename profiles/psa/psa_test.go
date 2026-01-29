// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/extensions"
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

// mustNewPSAInstanceID creates a valid PSA Instance ID (33-byte UEID with 0x01 RAND prefix)
func mustNewPSAInstanceID() *comid.Instance {
	// PSA Instance ID: 0x01 (RAND type) + 32 bytes of identifier
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

// mustNewPSAInstanceIDWithLength creates a PSA Instance ID with specified length
// Note: For EAT UEID validation, RAND type (0x01) must be followed by 16, 24, or 32 bytes
// So valid total lengths with RAND type are: 17, 25, or 33 bytes
// This function bypasses validation to allow testing with invalid UEIDs
func mustNewPSAInstanceIDWithLength(length int, typePrefix byte) *comid.Instance {
	ueidBytes := make([]byte, length)
	if length > 0 {
		ueidBytes[0] = typePrefix
	}
	for i := 1; i < length; i++ {
		ueidBytes[i] = byte(i % 256)
	}
	// Directly set the tagged value to bypass EAT validation
	// This allows us to test with intentionally invalid UEIDs
	tagged := comid.TaggedUEID(ueidBytes)
	return &comid.Instance{Value: tagged}
}

// newTestEnvironment creates a valid Environment with Instance (no Implementation ID)
func newTestEnvironment() comid.Environment {
	inst := comid.MustNewUUIDInstance(comid.TestUUID)
	return comid.Environment{
		Instance: inst,
	}
}

// newPSAEnvironmentWithImplID creates a valid PSA Environment with 32-byte Implementation ID
// and valid 33-byte UEID Instance ID (RAND type)
func newPSAEnvironmentWithImplID() comid.Environment {
	// Create a 32-byte Implementation ID
	implID := make([]byte, 32)
	for i := range implID {
		implID[i] = byte(i)
	}
	class := comid.NewClassBytes(implID)
	return comid.Environment{
		Class:    class,
		Instance: mustNewPSAInstanceID(),
	}
}

// newPSAEnvironmentWithImplIDLength creates Environment with Implementation ID of specified length
// and a valid PSA Instance ID
func newPSAEnvironmentWithImplIDLength(length int) comid.Environment {
	implID := make([]byte, length)
	for i := range implID {
		implID[i] = byte(i % 256)
	}
	class := comid.NewClassBytes(implID)
	return comid.Environment{
		Class:    class,
		Instance: mustNewPSAInstanceID(),
	}
}

func TestValidatePSASignerID_ValidAndInvalidLengths(t *testing.T) {
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
			title:           "invalid 1 byte",
			length:          1,
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
		{
			title:           "invalid 47 bytes",
			length:          47,
			shouldError:     true,
			expectedMessage: "must be 32, 48, or 64 bytes",
		},
		{
			title:           "invalid 49 bytes",
			length:          49,
			shouldError:     true,
			expectedMessage: "must be 32, 48, or 64 bytes",
		},
		{
			title:           "invalid 63 bytes",
			length:          63,
			shouldError:     true,
			expectedMessage: "must be 32, 48, or 64 bytes",
		},
		{
			title:           "invalid 65 bytes",
			length:          65,
			shouldError:     true,
			expectedMessage: "must be 32, 48, or 64 bytes",
		},
		{
			title:           "invalid 128 bytes",
			length:          128,
			shouldError:     true,
			expectedMessage: "must be 32, 48, or 64 bytes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			keys := &comid.CryptoKeys{mustNewTaggedBytesCryptoKey(tc.length)}
			err := validatePSASignerID(keys, 0, 0)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePSASignerID_InvalidCases(t *testing.T) {
	testCases := []struct {
		title           string
		keys            *comid.CryptoKeys
		expectedMessage string
	}{
		{
			title:           "nil keys",
			keys:            nil,
			expectedMessage: "cryptokeys (signer-id) is mandatory",
		},
		{
			title:           "empty cryptokeys",
			keys:            &comid.CryptoKeys{},
			expectedMessage: "must contain exactly one entry, got 0",
		},
		{
			title: "multiple entries",
			keys: &comid.CryptoKeys{
				mustNewTaggedBytesCryptoKey(32),
				mustNewTaggedBytesCryptoKey(32),
			},
			expectedMessage: "must contain exactly one entry, got 2",
		},
		{
			title:           "nil entry",
			keys:            &comid.CryptoKeys{nil},
			expectedMessage: "entry is nil",
		},
		{
			title:           "wrong type (pkix-base64-key instead of bytes)",
			keys:            &comid.CryptoKeys{mustNewPKIXKey()},
			expectedMessage: "must be of type 'bytes'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := validatePSASignerID(tc.keys, 0, 0)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedMessage)
		})
	}
}

func TestValidatePSAImplementationID_AllCases(t *testing.T) {
	testCases := []struct {
		title           string
		setupEnv        func() comid.Environment
		shouldError     bool
		expectedMessage string
	}{
		{
			title:       "valid 32 bytes",
			setupEnv:    newPSAEnvironmentWithImplID,
			shouldError: false,
		},
		{
			title: "invalid no class",
			setupEnv: func() comid.Environment {
				return comid.Environment{
					Instance: comid.MustNewUUIDInstance(comid.TestUUID),
				}
			},
			shouldError:     true,
			expectedMessage: "environment.class is required",
		},
		{
			title: "invalid no classid",
			setupEnv: func() comid.Environment {
				return comid.Environment{
					Class:    &comid.Class{},
					Instance: comid.MustNewUUIDInstance(comid.TestUUID),
				}
			},
			shouldError:     true,
			expectedMessage: "implementation-id) is required",
		},
		{
			title: "invalid wrong type (uuid instead of bytes)",
			setupEnv: func() comid.Environment {
				return comid.Environment{
					Class:    comid.NewClassUUID(comid.TestUUID),
					Instance: comid.MustNewUUIDInstance(comid.TestUUID),
				}
			},
			shouldError:     true,
			expectedMessage: "must be of type 'bytes'",
		},
		{
			title:           "invalid 0 bytes",
			setupEnv:        func() comid.Environment { return newPSAEnvironmentWithImplIDLength(0) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "invalid 1 byte",
			setupEnv:        func() comid.Environment { return newPSAEnvironmentWithImplIDLength(1) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "invalid 16 bytes",
			setupEnv:        func() comid.Environment { return newPSAEnvironmentWithImplIDLength(16) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "invalid 31 bytes",
			setupEnv:        func() comid.Environment { return newPSAEnvironmentWithImplIDLength(31) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "invalid 33 bytes",
			setupEnv:        func() comid.Environment { return newPSAEnvironmentWithImplIDLength(33) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "invalid 48 bytes",
			setupEnv:        func() comid.Environment { return newPSAEnvironmentWithImplIDLength(48) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "invalid 64 bytes",
			setupEnv:        func() comid.Environment { return newPSAEnvironmentWithImplIDLength(64) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "invalid 128 bytes",
			setupEnv:        func() comid.Environment { return newPSAEnvironmentWithImplIDLength(128) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			env := tc.setupEnv()
			err := validatePSAImplementationID(&env, "test")
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePSAInstanceID_AllCases(t *testing.T) {
	testCases := []struct {
		title           string
		setupEnv        func() comid.Environment
		shouldError     bool
		expectedMessage string
	}{
		{
			title: "valid 33-byte UEID with RAND type",
			setupEnv: func() comid.Environment {
				return comid.Environment{
					Instance: mustNewPSAInstanceID(),
				}
			},
			shouldError: false,
		},
		{
			title: "invalid no instance",
			setupEnv: func() comid.Environment {
				return comid.Environment{}
			},
			shouldError:     true,
			expectedMessage: "instance-id) is required",
		},
		{
			title: "invalid wrong type (uuid instead of ueid)",
			setupEnv: func() comid.Environment {
				return comid.Environment{
					Instance: comid.MustNewUUIDInstance(comid.TestUUID),
				}
			},
			shouldError:     true,
			expectedMessage: "must be of type 'ueid'",
		},
		{
			title: "invalid 32 bytes (too short)",
			setupEnv: func() comid.Environment {
				return comid.Environment{
					Instance: mustNewPSAInstanceIDWithLength(32, 0x01),
				}
			},
			shouldError:     true,
			expectedMessage: "must be exactly 33 bytes",
		},
		{
			title: "invalid 34 bytes (too long)",
			setupEnv: func() comid.Environment {
				return comid.Environment{
					Instance: mustNewPSAInstanceIDWithLength(34, 0x01),
				}
			},
			shouldError:     true,
			expectedMessage: "must be exactly 33 bytes",
		},
		{
			title: "invalid wrong type prefix (0x00 instead of 0x01)",
			setupEnv: func() comid.Environment {
				return comid.Environment{
					Instance: mustNewPSAInstanceIDWithLength(33, 0x00),
				}
			},
			shouldError:     true,
			expectedMessage: "must have RAND type (0x01)",
		},
		{
			title: "invalid wrong type prefix (0x02 instead of 0x01)",
			setupEnv: func() comid.Environment {
				return comid.Environment{
					Instance: mustNewPSAInstanceIDWithLength(33, 0x02),
				}
			},
			shouldError:     true,
			expectedMessage: "must have RAND type (0x01)",
		},
		{
			title: "invalid wrong type prefix (0xFF instead of 0x01)",
			setupEnv: func() comid.Environment {
				return comid.Environment{
					Instance: mustNewPSAInstanceIDWithLength(33, 0xFF),
				}
			},
			shouldError:     true,
			expectedMessage: "must have RAND type (0x01)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			env := tc.setupEnv()
			err := validatePSAInstanceID(&env, "test")
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePSAAttestVerifKey_AllCases(t *testing.T) {
	testCases := []struct {
		title           string
		setupAVK        func() *comid.KeyTriple
		shouldError     bool
		expectedMessage string
	}{
		{
			title: "valid attestation verification key",
			setupAVK: func() *comid.KeyTriple {
				return &comid.KeyTriple{
					Environment: newPSAEnvironmentWithImplID(),
					VerifKeys:   comid.CryptoKeys{mustNewPKIXKey()},
				}
			},
			shouldError: false,
		},
		{
			title: "invalid no implementation id",
			setupAVK: func() *comid.KeyTriple {
				return &comid.KeyTriple{
					Environment: newTestEnvironment(),
					VerifKeys:   comid.CryptoKeys{mustNewPKIXKey()},
				}
			},
			shouldError:     true,
			expectedMessage: "environment.class is required",
		},
		{
			title: "invalid empty verif keys",
			setupAVK: func() *comid.KeyTriple {
				return &comid.KeyTriple{
					Environment: newPSAEnvironmentWithImplID(),
					VerifKeys:   comid.CryptoKeys{},
				}
			},
			shouldError:     true,
			expectedMessage: "must contain exactly one entry, got 0",
		},
		{
			title: "invalid multiple keys",
			setupAVK: func() *comid.KeyTriple {
				return &comid.KeyTriple{
					Environment: newPSAEnvironmentWithImplID(),
					VerifKeys:   comid.CryptoKeys{mustNewPKIXKey(), mustNewPKIXKey()},
				}
			},
			shouldError:     true,
			expectedMessage: "must contain exactly one entry, got 2",
		},
		{
			title: "invalid nil entry",
			setupAVK: func() *comid.KeyTriple {
				return &comid.KeyTriple{
					Environment: newPSAEnvironmentWithImplID(),
					VerifKeys:   comid.CryptoKeys{nil},
				}
			},
			shouldError:     true,
			expectedMessage: "entry is nil",
		},
		{
			title: "invalid wrong type (bytes instead of pkix-base64-key)",
			setupAVK: func() *comid.KeyTriple {
				return &comid.KeyTriple{
					Environment: newPSAEnvironmentWithImplID(),
					VerifKeys:   comid.CryptoKeys{mustNewTaggedBytesCryptoKey(32)},
				}
			},
			shouldError:     true,
			expectedMessage: "must be of type 'pkix-base64-key'",
		},
		{
			title: "invalid no instance id",
			setupAVK: func() *comid.KeyTriple {
				// Create env with only Implementation ID, no Instance ID
				implID := make([]byte, 32)
				for i := range implID {
					implID[i] = byte(i)
				}
				class := comid.NewClassBytes(implID)
				return &comid.KeyTriple{
					Environment: comid.Environment{
						Class: class,
						// No Instance
					},
					VerifKeys: comid.CryptoKeys{mustNewPKIXKey()},
				}
			},
			shouldError:     true,
			expectedMessage: "instance-id) is required",
		},
		{
			title: "invalid instance id wrong type (uuid instead of ueid)",
			setupAVK: func() *comid.KeyTriple {
				implID := make([]byte, 32)
				for i := range implID {
					implID[i] = byte(i)
				}
				class := comid.NewClassBytes(implID)
				return &comid.KeyTriple{
					Environment: comid.Environment{
						Class:    class,
						Instance: comid.MustNewUUIDInstance(comid.TestUUID),
					},
					VerifKeys: comid.CryptoKeys{mustNewPKIXKey()},
				}
			},
			shouldError:     true,
			expectedMessage: "must be of type 'ueid'",
		},
		{
			title: "invalid instance id wrong length (32 bytes instead of 33)",
			setupAVK: func() *comid.KeyTriple {
				implID := make([]byte, 32)
				for i := range implID {
					implID[i] = byte(i)
				}
				class := comid.NewClassBytes(implID)
				return &comid.KeyTriple{
					Environment: comid.Environment{
						Class:    class,
						Instance: mustNewPSAInstanceIDWithLength(32, 0x01),
					},
					VerifKeys: comid.CryptoKeys{mustNewPKIXKey()},
				}
			},
			shouldError:     true,
			expectedMessage: "must be exactly 33 bytes",
		},
		{
			title: "invalid instance id wrong type prefix",
			setupAVK: func() *comid.KeyTriple {
				implID := make([]byte, 32)
				for i := range implID {
					implID[i] = byte(i)
				}
				class := comid.NewClassBytes(implID)
				return &comid.KeyTriple{
					Environment: comid.Environment{
						Class:    class,
						Instance: mustNewPSAInstanceIDWithLength(33, 0x02), // wrong type
					},
					VerifKeys: comid.CryptoKeys{mustNewPKIXKey()},
				}
			},
			shouldError:     true,
			expectedMessage: "must have RAND type (0x01)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			avk := tc.setupAVK()
			err := validatePSAAttestVerifKey(avk, 0)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidTriples_AllCases(t *testing.T) {
	testCases := []struct {
		title           string
		setupTriples    func() *comid.Triples
		shouldError     bool
		expectedMessage string
	}{
		{
			title: "no triples",
			setupTriples: func() *comid.Triples {
				return &comid.Triples{}
			},
			shouldError: false,
		},
		{
			title: "valid attestation verification key",
			setupTriples: func() *comid.Triples {
				avk := comid.KeyTriple{
					Environment: newPSAEnvironmentWithImplID(),
					VerifKeys:   comid.CryptoKeys{mustNewPKIXKey()},
				}
				return &comid.Triples{
					AttestVerifKeys: &comid.KeyTriples{avk},
				}
			},
			shouldError: false,
		},
		{
			title: "invalid attestation verification key wrong type",
			setupTriples: func() *comid.Triples {
				avk := comid.KeyTriple{
					Environment: newPSAEnvironmentWithImplID(),
					VerifKeys:   comid.CryptoKeys{mustNewTaggedBytesCryptoKey(32)},
				}
				return &comid.Triples{
					AttestVerifKeys: &comid.KeyTriples{avk},
				}
			},
			shouldError:     true,
			expectedMessage: "must be of type 'pkix-base64-key'",
		},
		{
			title: "invalid attestation verification key multiple keys",
			setupTriples: func() *comid.Triples {
				avk := comid.KeyTriple{
					Environment: newPSAEnvironmentWithImplID(),
					VerifKeys:   comid.CryptoKeys{mustNewPKIXKey(), mustNewPKIXKey()},
				}
				return &comid.Triples{
					AttestVerifKeys: &comid.KeyTriples{avk},
				}
			},
			shouldError:     true,
			expectedMessage: "must contain exactly one entry, got 2",
		},
		{
			title: "invalid multiple attestation verification keys second invalid",
			setupTriples: func() *comid.Triples {
				validAVK := comid.KeyTriple{
					Environment: newPSAEnvironmentWithImplID(),
					VerifKeys:   comid.CryptoKeys{mustNewPKIXKey()},
				}
				invalidAVK := comid.KeyTriple{
					Environment: newPSAEnvironmentWithImplID(),
					VerifKeys:   comid.CryptoKeys{mustNewTaggedBytesCryptoKey(32)},
				}
				return &comid.Triples{
					AttestVerifKeys: &comid.KeyTriples{validAVK, invalidAVK},
				}
			},
			shouldError:     true,
			expectedMessage: "attester verification key at index 1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			ext := TriplesExtensions{}
			triples := tc.setupTriples()

			// Register extensions if triples has AttestVerifKeys for valid case
			if triples.AttestVerifKeys != nil && len(*triples.AttestVerifKeys) > 0 {
				extMap := extensions.NewMap().Add(comid.ExtTriples, &ext)
				err := triples.RegisterExtensions(extMap)
				require.NoError(t, err)
			}

			err := ext.ValidTriples(triples)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
