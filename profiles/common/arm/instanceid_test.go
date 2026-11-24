// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package arm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/profiles/common/arm/testhelpers"
)

func TestValidatePlatformInstanceID_AllCases(t *testing.T) {
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
					Instance: testhelpers.MustNewPlatformInstanceID(),
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
					Instance: testhelpers.MustNewPlatformInstanceIDWithLength(32, 0x01),
				}
			},
			shouldError:     true,
			expectedMessage: "must be exactly 33 bytes",
		},
		{
			title: "invalid 34 bytes (too long)",
			setupEnv: func() comid.Environment {
				return comid.Environment{
					Instance: testhelpers.MustNewPlatformInstanceIDWithLength(34, 0x01),
				}
			},
			shouldError:     true,
			expectedMessage: "must be exactly 33 bytes",
		},
		{
			title: "invalid wrong type prefix (0x00 instead of 0x01)",
			setupEnv: func() comid.Environment {
				return comid.Environment{
					Instance: testhelpers.MustNewPlatformInstanceIDWithLength(33, 0x00),
				}
			},
			shouldError:     true,
			expectedMessage: "must have RAND type (0x01)",
		},
		{
			title: "invalid wrong type prefix (0x02 instead of 0x01)",
			setupEnv: func() comid.Environment {
				return comid.Environment{
					Instance: testhelpers.MustNewPlatformInstanceIDWithLength(33, 0x02),
				}
			},
			shouldError:     true,
			expectedMessage: "must have RAND type (0x01)",
		},
		{
			title: "invalid wrong type prefix (0xFF instead of 0x01)",
			setupEnv: func() comid.Environment {
				return comid.Environment{
					Instance: testhelpers.MustNewPlatformInstanceIDWithLength(33, 0xFF),
				}
			},
			shouldError:     true,
			expectedMessage: "must have RAND type (0x01)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			env := tc.setupEnv()
			err := ValidateInstanceID(&env, "test")
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
