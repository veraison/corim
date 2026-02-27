// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cca

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper Functions

// mustNewPlatformImplID creates a valid 32-byte Platform Implementation ID
func mustNewPlatformImplID() []byte {
	implID := make([]byte, 32)
	for i := range implID {
		implID[i] = byte(i)
	}
	return implID
}

// mustNewPlatformInstanceID creates a valid 33-byte Platform Instance ID with 0x01 prefix
func mustNewPlatformInstanceID() []byte {
	instanceID := make([]byte, 33)
	instanceID[0] = 0x01 // RAND type
	for i := 1; i < 33; i++ {
		instanceID[i] = byte(i - 1)
	}
	return instanceID
}

// newPlatformImplIDWithLength creates an Implementation ID with specified length
func newPlatformImplIDWithLength(length int) []byte {
	implID := make([]byte, length)
	for i := range implID {
		implID[i] = byte(i % 256)
	}
	return implID
}

// newPlatformInstanceIDWithLength creates an Instance ID with specified length and prefix
func newPlatformInstanceIDWithLength(length int, prefix byte) []byte {
	instanceID := make([]byte, length)
	if length > 0 {
		instanceID[0] = prefix
	}
	for i := 1; i < length; i++ {
		instanceID[i] = byte(i % 256)
	}
	return instanceID
}

// Tests for NewPlatformImplIDClassID

func TestNewPlatformImplIDClassID_ValidAndInvalidLengths(t *testing.T) {
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
			title:           "invalid 0 bytes",
			length:          0,
			shouldError:     true,
			expectedMessage: "got 0 bytes, expected 32",
		},
		{
			title:           "invalid 16 bytes",
			length:          16,
			shouldError:     true,
			expectedMessage: "got 16 bytes, expected 32",
		},
		{
			title:           "invalid 31 bytes",
			length:          31,
			shouldError:     true,
			expectedMessage: "got 31 bytes, expected 32",
		},
		{
			title:           "invalid 33 bytes",
			length:          33,
			shouldError:     true,
			expectedMessage: "got 33 bytes, expected 32",
		},
		{
			title:           "invalid 64 bytes",
			length:          64,
			shouldError:     true,
			expectedMessage: "got 64 bytes, expected 32",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			implID := newPlatformImplIDWithLength(tc.length)
			classID, err := NewPlatformImplIDClassID(implID)

			if tc.shouldError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrWrongImplIDSize)
				assert.Contains(t, err.Error(), tc.expectedMessage)
				assert.Nil(t, classID)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, classID)
				assert.Equal(t, implID, classID.Bytes())
			}
		})
	}
}

func TestMustNewPlatformImplIDClassID_ValidAndPanics(t *testing.T) {
	testCases := []struct {
		title       string
		length      int
		shouldPanic bool
	}{
		{
			title:       "valid does not panic",
			length:      32,
			shouldPanic: false,
		},
		{
			title:       "invalid length panics",
			length:      31,
			shouldPanic: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			implID := newPlatformImplIDWithLength(tc.length)

			if tc.shouldPanic {
				assert.Panics(t, func() {
					MustNewPlatformImplIDClassID(implID)
				})
			} else {
				assert.NotPanics(t, func() {
					classID := MustNewPlatformImplIDClassID(implID)
					assert.NotNil(t, classID)
				})
			}
		})
	}
}

func TestNewClassPlatformImplID_Valid(t *testing.T) {
	implID := mustNewPlatformImplID()
	class, err := NewClassPlatformImplID(implID)
	require.NoError(t, err)
	assert.NotNil(t, class)
	assert.NotNil(t, class.ClassID)
}

// Tests for NewPlatformInstanceID

func TestNewPlatformInstanceID_ValidAndInvalidLengths(t *testing.T) {
	// Test valid case separately using helper
	t.Run("valid 33 bytes with 0x01 prefix", func(t *testing.T) {
		instanceID := mustNewPlatformInstanceID()
		inst, err := NewPlatformInstanceID(instanceID)
		assert.NoError(t, err)
		assert.NotNil(t, inst)
	})

	// Test invalid cases with table
	testCases := []struct {
		title           string
		length          int
		prefix          byte
		expectedError   error
		expectedMessage string
	}{
		{
			title:           "invalid 0 bytes",
			length:          0,
			prefix:          0x01,
			expectedError:   ErrWrongInstanceIDSize,
			expectedMessage: "got 0 bytes, expected 33",
		},
		{
			title:           "invalid 16 bytes",
			length:          16,
			prefix:          0x01,
			expectedError:   ErrWrongInstanceIDSize,
			expectedMessage: "got 16 bytes, expected 33",
		},
		{
			title:           "invalid 32 bytes (too short)",
			length:          32,
			prefix:          0x01,
			expectedError:   ErrWrongInstanceIDSize,
			expectedMessage: "got 32 bytes, expected 33",
		},
		{
			title:           "invalid 34 bytes (too long)",
			length:          34,
			prefix:          0x01,
			expectedError:   ErrWrongInstanceIDSize,
			expectedMessage: "got 34 bytes, expected 33",
		},
		{
			title:           "invalid 64 bytes",
			length:          64,
			prefix:          0x01,
			expectedError:   ErrWrongInstanceIDSize,
			expectedMessage: "got 64 bytes, expected 33",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			instanceID := newPlatformInstanceIDWithLength(tc.length, tc.prefix)
			inst, err := NewPlatformInstanceID(instanceID)

			assert.Error(t, err)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Contains(t, err.Error(), tc.expectedMessage)
			assert.Nil(t, inst)
		})
	}
}

func TestNewPlatformInstanceID_WrongPrefix(t *testing.T) {
	testCases := []struct {
		title           string
		prefix          byte
		expectedMessage string
	}{
		{
			title:           "prefix 0x00",
			prefix:          0x00,
			expectedMessage: "got 0x00",
		},
		{
			title:           "prefix 0x02",
			prefix:          0x02,
			expectedMessage: "got 0x02",
		},
		{
			title:           "prefix 0xFF",
			prefix:          0xFF,
			expectedMessage: "got 0xff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			instanceID := newPlatformInstanceIDWithLength(33, tc.prefix)
			_, err := NewPlatformInstanceID(instanceID)

			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrWrongInstancePrefix)
			assert.Contains(t, err.Error(), tc.expectedMessage)
		})
	}
}

func TestMustNewPlatformInstanceID_ValidAndPanics(t *testing.T) {
	// Test valid case using helper
	t.Run("valid does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			instanceID := mustNewPlatformInstanceID()
			inst := MustNewPlatformInstanceID(instanceID)
			assert.NotNil(t, inst)
		})
	})

	// Test invalid cases with table
	testCases := []struct {
		title  string
		length int
		prefix byte
	}{
		{
			title:  "invalid length panics",
			length: 32,
			prefix: 0x01,
		},
		{
			title:  "invalid prefix panics",
			length: 33,
			prefix: 0x02,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			instanceID := newPlatformInstanceIDWithLength(tc.length, tc.prefix)
			assert.Panics(t, func() {
				MustNewPlatformInstanceID(instanceID)
			})
		})
	}
}

// Tests for Validate functions

func TestValidateCCAPlatformImplID_AllCases(t *testing.T) {
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
			title:           "invalid 0 bytes",
			length:          0,
			shouldError:     true,
			expectedMessage: "got 0 bytes, expected 32",
		},
		{
			title:           "invalid 31 bytes",
			length:          31,
			shouldError:     true,
			expectedMessage: "got 31 bytes, expected 32",
		},
		{
			title:           "invalid 33 bytes",
			length:          33,
			shouldError:     true,
			expectedMessage: "got 33 bytes, expected 32",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			implID := newPlatformImplIDWithLength(tc.length)
			err := ValidatePlatformImplID(implID)

			if tc.shouldError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrWrongImplIDSize)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCCAPlatformInstanceID_AllCases(t *testing.T) {
	// Test valid case using helper
	t.Run("valid 33 bytes with 0x01", func(t *testing.T) {
		instanceID := mustNewPlatformInstanceID()
		err := ValidatePlatformInstanceID(instanceID)
		assert.NoError(t, err)
	})

	// Test invalid cases with table
	testCases := []struct {
		title           string
		length          int
		prefix          byte
		expectedError   error
		expectedMessage string
	}{
		{
			title:           "invalid size 32 bytes",
			length:          32,
			prefix:          0x01,
			expectedError:   ErrWrongInstanceIDSize,
			expectedMessage: "got 32 bytes, expected 33",
		},
		{
			title:           "invalid size 34 bytes",
			length:          34,
			prefix:          0x01,
			expectedError:   ErrWrongInstanceIDSize,
			expectedMessage: "got 34 bytes, expected 33",
		},
		{
			title:           "invalid prefix 0x00",
			length:          33,
			prefix:          0x00,
			expectedError:   ErrWrongInstancePrefix,
			expectedMessage: "got 0x00",
		},
		{
			title:           "invalid prefix 0x02",
			length:          33,
			prefix:          0x02,
			expectedError:   ErrWrongInstancePrefix,
			expectedMessage: "got 0x02",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			instanceID := newPlatformInstanceIDWithLength(tc.length, tc.prefix)
			err := ValidatePlatformInstanceID(instanceID)

			assert.Error(t, err)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Contains(t, err.Error(), tc.expectedMessage)
		})
	}
}
