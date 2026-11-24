// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package arm

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/profiles/common/arm/testhelpers"
)

// Test Implementation ID 32 bytes with identifiable pattern
var TestImplIDBytes = [32]byte{
	0x61, 0x63, 0x6d, 0x65, 0x2d, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2d, 0x69, 0x64, 0x2d, 0x30,
	0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31,
}

// TestImplIDBase64 is the base64-encoded string of TestImplIDBytes
var TestImplIDBase64 = base64.StdEncoding.EncodeToString(TestImplIDBytes[:])

func TestValidatePlatformImplementationID_AllCases(t *testing.T) {
	testCases := []struct {
		title           string
		setupEnv        func() comid.Environment
		shouldError     bool
		expectedMessage string
	}{
		{
			title:       "valid 32 bytes",
			setupEnv:    testhelpers.NewPlatformEnvironmentWithImplID,
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
			setupEnv:        func() comid.Environment { return testhelpers.NewPlatformEnvironmentWithImplIDLength(0) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "invalid 1 byte",
			setupEnv:        func() comid.Environment { return testhelpers.NewPlatformEnvironmentWithImplIDLength(1) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "invalid 16 bytes",
			setupEnv:        func() comid.Environment { return testhelpers.NewPlatformEnvironmentWithImplIDLength(16) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "invalid 31 bytes",
			setupEnv:        func() comid.Environment { return testhelpers.NewPlatformEnvironmentWithImplIDLength(31) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "invalid 33 bytes",
			setupEnv:        func() comid.Environment { return testhelpers.NewPlatformEnvironmentWithImplIDLength(33) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "invalid 48 bytes",
			setupEnv:        func() comid.Environment { return testhelpers.NewPlatformEnvironmentWithImplIDLength(48) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "invalid 64 bytes",
			setupEnv:        func() comid.Environment { return testhelpers.NewPlatformEnvironmentWithImplIDLength(64) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
		{
			title:           "invalid 128 bytes",
			setupEnv:        func() comid.Environment { return testhelpers.NewPlatformEnvironmentWithImplIDLength(128) },
			shouldError:     true,
			expectedMessage: "must be exactly 32 bytes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			env := tc.setupEnv()
			err := ValidateImplementationID(&env, "test")
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_NewImplIDClassID(t *testing.T) {
	bytes := TestImplIDBytes[:]
	taggedBytes := comid.TaggedBytes(TestImplIDBytes[:])

	for _, v := range []any{
		bytes,
		TestImplIDBase64,
		taggedBytes,
		&taggedBytes,
	} {
		t.Run(typeName(v), func(t *testing.T) {
			ret, err := NewImplIDClassID(v)
			require.NoError(t, err)
			require.NotNil(t, ret)
			require.NotNil(t, ret.Value)

			// Verify the value is TaggedBytes with the correct content
			tb, ok := ret.Value.(*comid.TaggedBytes)
			require.True(t, ok, "expected *comid.TaggedBytes, got %T", ret.Value)
			assert.Equal(t, TestImplIDBytes[:], tb.Bytes())

			// Verify that Type() returns "bytes" (TaggedBytes type)
			assert.Equal(t, "bytes", ret.Type())
		})
	}
}

func Test_NewImplIDClassID_NilValue(t *testing.T) {
	ret, err := NewImplIDClassID(nil)
	require.NoError(t, err)
	require.NotNil(t, ret)
	require.NotNil(t, ret.Value)

	// Should return zero-value 32-byte Implementation ID
	tb, ok := ret.Value.(*comid.TaggedBytes)
	require.True(t, ok)
	assert.Equal(t, make([]byte, 32), tb.Bytes())
}

func Test_NewImplIDClassID_InvalidLength(t *testing.T) {
	testCases := []struct {
		name        string
		input       any
		expectError string
	}{
		{
			name:        "bytes too short",
			input:       make([]byte, 16),
			expectError: "bad platform.impl-id: got 16 bytes, want 32",
		},
		{
			name:        "bytes too long",
			input:       make([]byte, 64),
			expectError: "bad platform.impl-id: got 64 bytes, want 32",
		},
		{
			name:        "TaggedBytes too short",
			input:       comid.TaggedBytes(make([]byte, 16)),
			expectError: "bad platform.impl-id: got 16 bytes, want 32",
		},
		{
			name:        "TaggedBytes too long",
			input:       comid.TaggedBytes(make([]byte, 64)),
			expectError: "bad platform.impl-id: got 64 bytes, want 32",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret, err := NewImplIDClassID(tc.input)
			require.Error(t, err)
			assert.Nil(t, ret)
			assert.Contains(t, err.Error(), tc.expectError)
		})
	}
}

func Test_NewImplIDClassID_InvalidBase64(t *testing.T) {
	ret, err := NewImplIDClassID("not-valid-base64!!!")
	require.Error(t, err)
	assert.Nil(t, ret)
	assert.Contains(t, err.Error(), "bad platform.impl-id")
}

func Test_NewImplIDClassID_InvalidBase64Length(t *testing.T) {
	// Valid base64 but wrong decoded length (16 bytes)
	shortBase64 := base64.StdEncoding.EncodeToString(make([]byte, 16))
	ret, err := NewImplIDClassID(shortBase64)
	require.Error(t, err)
	assert.Nil(t, ret)
	assert.Contains(t, err.Error(), "decoded 16 bytes, want 32")
}

func Test_NewImplIDClassID_UnsupportedType(t *testing.T) {
	ret, err := NewImplIDClassID(12345)
	require.Error(t, err)
	assert.Nil(t, ret)
	assert.Contains(t, err.Error(), "unexpected type for platform.impl-id: int")
}

func Test_MustNewImplIDClassID_Success(t *testing.T) {
	// Should not panic with valid input
	ret := MustNewImplIDClassID(TestImplIDBytes[:])
	require.NotNil(t, ret)
	require.NotNil(t, ret.Value)

	tb, ok := ret.Value.(*comid.TaggedBytes)
	require.True(t, ok)
	assert.Equal(t, TestImplIDBytes[:], tb.Bytes())
}

func Test_MustNewImplIDClassID_Panic(t *testing.T) {
	// Should panic with invalid input
	assert.Panics(t, func() {
		MustNewImplIDClassID(make([]byte, 16)) // Invalid length
	})

	assert.Panics(t, func() {
		MustNewImplIDClassID("invalid-base64!!!")
	})

	assert.Panics(t, func() {
		MustNewImplIDClassID(12345) // Unsupported type
	})
}

func Test_NewClassImplID(t *testing.T) {
	class := NewClassImplID(TestImplIDBytes[:])
	require.NotNil(t, class)
	require.NotNil(t, class.ClassID)
	require.NotNil(t, class.ClassID.Value)

	tb, ok := class.ClassID.Value.(*comid.TaggedBytes)
	require.True(t, ok)
	assert.Equal(t, TestImplIDBytes[:], tb.Bytes())
}

func Test_NewClassImplID_ZeroValue(t *testing.T) {
	zeroImplID := make([]byte, 32)

	class := NewClassImplID(zeroImplID)
	require.NotNil(t, class)
	require.NotNil(t, class.ClassID)
	require.NotNil(t, class.ClassID.Value)

	tb, ok := class.ClassID.Value.(*comid.TaggedBytes)
	require.True(t, ok)
	assert.Equal(t, make([]byte, 32), tb.Bytes())
}

func Test_NewClassImplID_InvalidLength(t *testing.T) {
	// Should return nil for invalid length
	class := NewClassImplID(make([]byte, 16))
	assert.Nil(t, class)

	class = NewClassImplID(make([]byte, 64))
	assert.Nil(t, class)
}

// Helper function to get type name for test case naming
func typeName(v any) string {
	switch v.(type) {
	case []byte:
		return "[]byte"
	case string:
		return "string(base64)"
	case comid.TaggedBytes:
		return "comid.TaggedBytes"
	case *comid.TaggedBytes:
		return "*comid.TaggedBytes"
	default:
		return "unknown"
	}
}
