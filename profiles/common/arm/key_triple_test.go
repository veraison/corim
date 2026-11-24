// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package arm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/profiles/common/arm/testhelpers"
)

func TestValidateAttestVerifKey_AllCases(t *testing.T) {
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
					VerifKeys:   comid.CryptoKeys{testhelpers.MustNewPKIXKey()},
				}
			},
			shouldError: false,
		},
		{
			title: "invalid no implementation id",
			setupAVK: func() *comid.KeyTriple {
				return &comid.KeyTriple{
					Environment: testhelpers.NewTestEnvironment(),
					VerifKeys:   comid.CryptoKeys{testhelpers.MustNewPKIXKey()},
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
					VerifKeys:   comid.CryptoKeys{testhelpers.MustNewPKIXKey(), testhelpers.MustNewPKIXKey()},
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
					VerifKeys:   comid.CryptoKeys{testhelpers.MustNewTaggedBytesCryptoKey(32)},
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
					VerifKeys: comid.CryptoKeys{testhelpers.MustNewPKIXKey()},
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
					VerifKeys: comid.CryptoKeys{testhelpers.MustNewPKIXKey()},
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
					VerifKeys: comid.CryptoKeys{testhelpers.MustNewPKIXKey()},
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
					VerifKeys: comid.CryptoKeys{testhelpers.MustNewPKIXKey()},
				}
			},
			shouldError:     true,
			expectedMessage: "must have RAND type (0x01)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			avk := tc.setupAVK()
			err := ValidateArmPlatformAttestVerifKey(avk)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
