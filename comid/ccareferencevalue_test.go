// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/swid"
)

var TestCCASignerID = []byte{
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
	0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
	0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
}

func TestCCARefValID_Valid(t *testing.T) {
	// Valid case with 32 bytes
	refValID := &CCARefValID{
		SignerID: TestCCASignerID,
	}
	err := refValID.Valid()
	assert.NoError(t, err)

	// Valid case with 48 bytes
	refValID = &CCARefValID{
		SignerID: make([]byte, 48),
	}
	err = refValID.Valid()
	assert.NoError(t, err)

	// Valid case with 64 bytes
	refValID = &CCARefValID{
		SignerID: make([]byte, 64),
	}
	err = refValID.Valid()
	assert.NoError(t, err)

	// Invalid case with nil SignerID
	refValID = &CCARefValID{
		SignerID: nil,
	}
	err = refValID.Valid()
	assert.Error(t, err)

	// Invalid case with wrong length
	refValID = &CCARefValID{
		SignerID: []byte{0x01, 0x02, 0x03},
	}
	err = refValID.Valid()
	assert.Error(t, err)
}

func TestNewCCARefValID(t *testing.T) {
	// Test with nil
	refValID, err := NewCCARefValID(nil)
	require.NoError(t, err)
	assert.NotNil(t, refValID)

	// Test with CCARefValID
	original := &CCARefValID{
		SignerID: TestCCASignerID,
	}
	label := "Label"
	version := "1.0.0"
	original.SetLabel(label)
	original.SetVersion(version)

	refValID, err = NewCCARefValID(*original)
	require.NoError(t, err)
	assert.Equal(t, TestCCASignerID, refValID.SignerID)
	assert.Equal(t, &label, refValID.Label)
	assert.Equal(t, &version, refValID.Version)

	// Test with pointer to CCARefValID
	refValID, err = NewCCARefValID(original)
	require.NoError(t, err)
	assert.Equal(t, TestCCASignerID, refValID.SignerID)

	// Test with byte array (signer ID)
	refValID, err = NewCCARefValID(TestCCASignerID)
	require.NoError(t, err)
	assert.Equal(t, TestCCASignerID, refValID.SignerID)

	// Test with byte array of invalid length
	_, err = NewCCARefValID([]byte{0x01, 0x02})
	assert.Error(t, err)

	// Test with invalid type
	_, err = NewCCARefValID(123)
	assert.Error(t, err)
}

func TestCreateCCARefValID(t *testing.T) {
	label := "CCA-BL"
	version := "1.0.0"

	refValID, err := CreateCCARefValID(TestCCASignerID, label, version)
	require.NoError(t, err)

	assert.Equal(t, TestCCASignerID, refValID.SignerID)
	assert.Equal(t, &label, refValID.Label)
	assert.Equal(t, &version, refValID.Version)

	// Invalid signer ID
	_, err = CreateCCARefValID([]byte{0x01}, label, version)
	assert.Error(t, err)
}

func TestMustCreateCCARefValID(t *testing.T) {
	label := "CCA-BL"
	version := "1.0.0"

	refValID := MustCreateCCARefValID(TestCCASignerID, label, version)
	assert.Equal(t, TestCCASignerID, refValID.SignerID)
	assert.Equal(t, &label, refValID.Label)
	assert.Equal(t, &version, refValID.Version)

	// Should panic with invalid signer ID
	assert.Panics(t, func() {
		MustCreateCCARefValID([]byte{0x01}, label, version)
	})
}

func TestCCARefValID_SetLabel(t *testing.T) {
	refValID := &CCARefValID{}
	label := "CCA-BL"

	refValID.SetLabel(label)
	assert.Equal(t, &label, refValID.Label)
}

func TestCCARefValID_SetVersion(t *testing.T) {
	refValID := &CCARefValID{}
	version := "1.0.0"

	refValID.SetVersion(version)
	assert.Equal(t, &version, refValID.Version)
}

func TestTaggedCCARefValID(t *testing.T) {
	// Test Type()
	var tagged TaggedCCARefValID
	assert.Equal(t, CCARefValIDType, tagged.Type())

	// Test Valid()
	tagged = TaggedCCARefValID{
		SignerID: TestCCASignerID,
	}
	err := tagged.Valid()
	assert.NoError(t, err)

	tagged = TaggedCCARefValID{
		SignerID: nil,
	}
	err = tagged.Valid()
	assert.Error(t, err)

	// Test String()
	tagged = TaggedCCARefValID{
		SignerID: TestCCASignerID,
	}
	str := tagged.String()
	assert.NotEmpty(t, str)

	// Test IsZero()
	tagged = TaggedCCARefValID{
		SignerID: TestCCASignerID,
	}
	assert.False(t, tagged.IsZero())

	tagged = TaggedCCARefValID{
		SignerID: nil,
	}
	assert.True(t, tagged.IsZero())
}

func TestNewTaggedCCARefValID(t *testing.T) {
	// Test with nil
	tagged, err := NewTaggedCCARefValID(nil)
	require.NoError(t, err)
	assert.NotNil(t, tagged)

	// Test with CCARefValID
	refValID := &CCARefValID{
		SignerID: TestCCASignerID,
	}
	tagged, err = NewTaggedCCARefValID(refValID)
	require.NoError(t, err)
	assert.Equal(t, TestCCASignerID, tagged.SignerID)

	// Test with TaggedCCARefValID
	original := TaggedCCARefValID{
		SignerID: TestCCASignerID,
	}
	tagged, err = NewTaggedCCARefValID(original)
	require.NoError(t, err)
	assert.Equal(t, TestCCASignerID, tagged.SignerID)

	// Test with pointer to TaggedCCARefValID
	tagged, err = NewTaggedCCARefValID(&original)
	require.NoError(t, err)
	assert.Equal(t, TestCCASignerID, tagged.SignerID)

	// Test with byte array (signer ID)
	tagged, err = NewTaggedCCARefValID(TestCCASignerID)
	require.NoError(t, err)
	assert.Equal(t, TestCCASignerID, tagged.SignerID)

	// Test with byte array of invalid length
	_, err = NewTaggedCCARefValID([]byte{0x01, 0x02})
	assert.Error(t, err)
}

func TestNewMkeyCCARefValID(t *testing.T) {
	refValID := &CCARefValID{
		SignerID: TestCCASignerID,
	}

	mkey, err := NewMkeyCCARefValID(refValID)
	require.NoError(t, err)

	assert.Equal(t, CCARefValIDType, mkey.Type())
	assert.True(t, mkey.IsSet())

	// Test with invalid input
	_, err = NewMkeyCCARefValID([]byte{0x01, 0x02})
	assert.Error(t, err)
}

func TestCCARefValID_JSON(t *testing.T) {
	// Create a test CCARefValID
	label := "CCA-BL"
	version := "1.0.0"
	refValID := MustCreateCCARefValID(TestCCASignerID, label, version)

	// Marshal to JSON
	data, err := json.Marshal(refValID)
	require.NoError(t, err)

	// Unmarshal from JSON
	var unmarshalled CCARefValID
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	// Check values
	assert.Equal(t, label, *unmarshalled.Label)
	assert.Equal(t, version, *unmarshalled.Version)
	assert.Equal(t, TestCCASignerID, unmarshalled.SignerID)
}

func TestCCAMeasurement(t *testing.T) {
	// Test NewCCAMeasurement
	refValID := MustCreateCCARefValID(TestCCASignerID, "CCA-BL", "1.0.0")

	measurement, err := NewCCAMeasurement(refValID)
	require.NoError(t, err)

	assert.Equal(t, CCARefValIDType, measurement.Key.Type())

	// Test MustNewCCAMeasurement
	measurement = MustNewCCAMeasurement(refValID)
	assert.Equal(t, CCARefValIDType, measurement.Key.Type())

	// Test MustNewCCAMeasurement with invalid input
	invalidRefValID := &CCARefValID{
		SignerID: []byte{0x01, 0x02}, // Invalid length
	}
	assert.Panics(t, func() {
		MustNewCCAMeasurement(invalidRefValID)
	})
}

func TestCCARefValIDRegistration(t *testing.T) {
	// Check if CCARefValIDType is registered in mkeyValueRegister
	factory, ok := mkeyValueRegister[CCARefValIDType]
	assert.True(t, ok, "CCA RefVal ID type should be registered")
	assert.NotNil(t, factory, "Factory function should not be nil")

	// Test the factory function
	refValID := MustCreateCCARefValID(TestCCASignerID, "CCA-BL", "1.0.0")
	mkey, err := factory(refValID)
	require.NoError(t, err)
	assert.Equal(t, CCARefValIDType, mkey.Type())
}

func TestCCARefValInMeasurement(t *testing.T) {
	refValID := MustCreateCCARefValID(TestCCASignerID, "CCA-BL", "1.0.0")

	measurement := MustNewCCAMeasurement(refValID)

	digest := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	measurement.AddDigest(swid.Sha256, digest)

	// Validate the measurement
	err := measurement.Valid()
	assert.NoError(t, err)
}
