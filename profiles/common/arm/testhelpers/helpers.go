// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import "github.com/veraison/corim/comid"

// newPlatformEnvironmentWithImplID creates a valid Platform Environment with 32-byte Implementation ID
// and valid 33-byte UEID Instance ID (RAND type)
func NewPlatformEnvironmentWithImplID() comid.Environment {
	// Create a 32-byte Implementation ID
	implID := make([]byte, 32)
	for i := range implID {
		implID[i] = byte(i)
	}
	class := comid.NewClassBytes(implID)
	return comid.Environment{
		Class:    class,
		Instance: MustNewPlatformInstanceID(),
	}
}

// NewPlatformEnvironmentWithImplIDLength creates Environment with Implementation ID of specified length
// and a valid Platform Instance ID
func NewPlatformEnvironmentWithImplIDLength(length int) comid.Environment {
	implID := make([]byte, length)
	for i := range implID {
		implID[i] = byte(i % 256)
	}
	class := comid.NewClassBytes(implID)
	return comid.Environment{
		Class:    class,
		Instance: MustNewPlatformInstanceID(),
	}
}

// mustNewPlatformInstanceID creates a valid Platform Instance ID (33-byte UEID with 0x01 RAND prefix)
func MustNewPlatformInstanceID() *comid.Instance {
	// Platform Instance ID: 0x01 (RAND type) + 32 bytes of identifier
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

// mustNewPlatformInstanceIDWithLength creates a Platform Instance ID with specified length
// Note: For EAT UEID validation, RAND type (0x01) must be followed by 16, 24, or 32 bytes
// So valid total lengths with RAND type are: 17, 25, or 33 bytes
// This function bypasses validation to allow testing with invalid UEIDs
func MustNewPlatformInstanceIDWithLength(length int, typePrefix byte) *comid.Instance {
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

// MustNewPKIXKey creates a PKIX Base64 Key (for AttestVerifKeys)
//
//nolint:lll
func MustNewPKIXKey() *comid.CryptoKey {
	return comid.MustNewCryptoKey(
		"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----",
		comid.PKIXBase64KeyType,
	)
}

// mustNewTaggedBytesCryptoKey creates a TaggedBytes CryptoKey with specified length
func MustNewTaggedBytesCryptoKey(length int) *comid.CryptoKey {
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
