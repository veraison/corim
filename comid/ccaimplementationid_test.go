// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
    "encoding/base64"
    "encoding/json"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

var TestCCAImplID = CCAImplID{
    0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
    0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
    0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
    0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
}

func TestCCAImplID_Valid(t *testing.T) {
    implID := TestCCAImplID
    err := implID.Valid()
    assert.NoError(t, err)
}

func TestCCAImplID_String(t *testing.T) {
    implID := TestCCAImplID
    expected := base64.StdEncoding.EncodeToString(implID[:])
    assert.Equal(t, expected, implID.String())
}

func TestTaggedCCAImplID_Valid(t *testing.T) {
    var implID TaggedCCAImplID
    copy(implID[:], TestCCAImplID[:])
    err := implID.Valid()
    assert.NoError(t, err)
}

func TestTaggedCCAImplID_Type(t *testing.T) {
    var implID TaggedCCAImplID
    assert.Equal(t, CCAImplIDType, implID.Type())
}

func TestTaggedCCAImplID_String(t *testing.T) {
    var implID TaggedCCAImplID
    copy(implID[:], TestCCAImplID[:])
    expected := base64.StdEncoding.EncodeToString(implID[:])
    assert.Equal(t, expected, implID.String())
}

func TestNewCCAImplIDClassID(t *testing.T) {
    // Test with nil
    classID, err := NewCCAImplIDClassID(nil)
    require.NoError(t, err)
    assert.Equal(t, CCAImplIDType, classID.Type())
    
    // Test with byte array
    classID, err = NewCCAImplIDClassID(TestCCAImplID[:])
    require.NoError(t, err)
    assert.Equal(t, CCAImplIDType, classID.Type())

    // Test with base64 string
    base64Str := base64.StdEncoding.EncodeToString(TestCCAImplID[:])
    classID, err = NewCCAImplIDClassID(base64Str)
    require.NoError(t, err)
    assert.Equal(t, CCAImplIDType, classID.Type())

    // Test with CCAImplID
    classID, err = NewCCAImplIDClassID(TestCCAImplID)
    require.NoError(t, err)
    assert.Equal(t, CCAImplIDType, classID.Type())

    // Test with TaggedCCAImplID
    var taggedID TaggedCCAImplID
    copy(taggedID[:], TestCCAImplID[:])
    classID, err = NewCCAImplIDClassID(taggedID)
    require.NoError(t, err)
    assert.Equal(t, CCAImplIDType, classID.Type())

    // Test with pointer to TaggedCCAImplID
    classID, err = NewCCAImplIDClassID(&taggedID)
    require.NoError(t, err)
    assert.Equal(t, CCAImplIDType, classID.Type())

    // Test with pointer to CCAImplID
    ptr := &TestCCAImplID
    classID, err = NewCCAImplIDClassID(ptr)
    require.NoError(t, err)
    assert.Equal(t, CCAImplIDType, classID.Type())

    // Test with invalid length
    _, err = NewCCAImplIDClassID([]byte{0x01, 0x02, 0x03})
    assert.Error(t, err)

    // Test with invalid type
    _, err = NewCCAImplIDClassID(123)
    assert.Error(t, err)
}

func TestMustNewCCAImplIDClassID(t *testing.T) {
    classID := MustNewCCAImplIDClassID(TestCCAImplID)
    assert.Equal(t, CCAImplIDType, classID.Type())

    assert.Panics(t, func() {
        MustNewCCAImplIDClassID([]byte{0x01, 0x02, 0x03})
    })
}

func TestClassID_SetGetCCAImplID(t *testing.T) {
    classID := new(ClassID).SetCCAImplID(TestCCAImplID)
    
    implID, err := classID.GetCCAImplID()
    require.NoError(t, err)
    assert.Equal(t, TestCCAImplID, implID)

    // Test with wrong type
    classID, err = NewClassID(TestUUID, UUIDType)
    _, err = classID.GetCCAImplID()
    assert.Error(t, err)
}

func TestNewClassCCAImplID(t *testing.T) {
    classID := NewClassCCAImplID(TestCCAImplID)
    
    implID, err := classID.GetCCAImplID()
    require.NoError(t, err)
    assert.Equal(t, TestCCAImplID, implID)
    assert.Equal(t, CCAImplIDType, classID.Type())
}

func TestTaggedCCAImplID_JSON(t *testing.T) {
    var implID TaggedCCAImplID
    copy(implID[:], TestCCAImplID[:])

    // Test marshaling
    bytes, err := json.Marshal(implID)
    require.NoError(t, err)

    // Test unmarshaling
    var unmarshalledID TaggedCCAImplID
    err = json.Unmarshal(bytes, &unmarshalledID)
    require.NoError(t, err)

    assert.Equal(t, implID, unmarshalledID)

    // Test unmarshaling with invalid data
    invalidData := []byte(`"not valid base64"`)
    err = json.Unmarshal(invalidData, &unmarshalledID)
    assert.Error(t, err)

    // Test unmarshaling with invalid length
    shortData := []byte(`[1,2,3]`)
    err = json.Unmarshal(shortData, &unmarshalledID)
    assert.Error(t, err)
}

func TestCCAImplIDClassIDRegistration(t *testing.T) {
    factory, ok := classIDValueRegister[CCAImplIDType]
    assert.True(t, ok, "CCA Impl ID type should be registered")
    assert.NotNil(t, factory, "Factory function should not be nil")
    
    // Test the factory function
    classID, err := factory(TestCCAImplID)
    require.NoError(t, err)
    assert.Equal(t, CCAImplIDType, classID.Type())
}