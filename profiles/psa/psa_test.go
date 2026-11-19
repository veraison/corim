// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/corim/profiles/common/arm/testhelpers"
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

// newTestEnvironment creates a valid Environment with Instance (no Implementation ID)
func newTestEnvironment() comid.Environment {
	inst := comid.MustNewUUIDInstance(comid.TestUUID)
	return comid.Environment{
		Instance: inst,
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
					Environment: testhelpers.NewPlatformEnvironmentWithImplID(),
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
					Environment: testhelpers.NewPlatformEnvironmentWithImplID(),
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
					Environment: testhelpers.NewPlatformEnvironmentWithImplID(),
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
					Environment: testhelpers.NewPlatformEnvironmentWithImplID(),
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
					Environment: testhelpers.NewPlatformEnvironmentWithImplID(),
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
						Instance: testhelpers.MustNewPlatformInstanceIDWithLength(32, 0x01),
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
						Instance: testhelpers.MustNewPlatformInstanceIDWithLength(33, 0x02), // wrong type
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
					Environment: testhelpers.NewPlatformEnvironmentWithImplID(),
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
					Environment: testhelpers.NewPlatformEnvironmentWithImplID(),
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
					Environment: testhelpers.NewPlatformEnvironmentWithImplID(),
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
					Environment: testhelpers.NewPlatformEnvironmentWithImplID(),
					VerifKeys:   comid.CryptoKeys{mustNewPKIXKey()},
				}
				invalidAVK := comid.KeyTriple{
					Environment: testhelpers.NewPlatformEnvironmentWithImplID(),
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
