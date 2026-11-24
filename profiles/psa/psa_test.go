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
					VerifKeys:   comid.CryptoKeys{testhelpers.MustNewPKIXKey()},
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
					VerifKeys:   comid.CryptoKeys{testhelpers.MustNewTaggedBytesCryptoKey(32)},
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
					VerifKeys:   comid.CryptoKeys{testhelpers.MustNewPKIXKey(), testhelpers.MustNewPKIXKey()},
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
					VerifKeys:   comid.CryptoKeys{testhelpers.MustNewPKIXKey()},
				}
				invalidAVK := comid.KeyTriple{
					Environment: testhelpers.NewPlatformEnvironmentWithImplID(),
					VerifKeys:   comid.CryptoKeys{testhelpers.MustNewTaggedBytesCryptoKey(32)},
				}
				return &comid.Triples{
					AttestVerifKeys: &comid.KeyTriples{validAVK, invalidAVK},
				}
			},
			shouldError:     true,
			expectedMessage: "attestation verification key at index 1",
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
