// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package arm

import (
	"fmt"

	"github.com/veraison/corim/comid"
)

// validateSignerID validates the cryptokeys field (signer-id)
//
//	a) CryptoKeys field is set and contains exactly one entry
//	b) The type of CryptoKey is TaggedBytes (Type() returns "bytes")
//	c) The length of the value is 32, 48, or 64
func validateSignerID(keys *comid.CryptoKeys, tripleIndex, measurementIndex int) error {
	prefix := fmt.Sprintf("reference value at index %d, measurement at index %d", tripleIndex, measurementIndex)

	// a) CryptoKeys field MUST be set (mandatory per any Arm profile)
	if keys == nil {
		return fmt.Errorf("%s: cryptokeys (signer-id) is mandatory but not set", prefix)
	}

	// a) Must contain exactly one entry
	if len(*keys) == 0 {
		return fmt.Errorf("%s: cryptokeys must contain exactly one entry, got 0", prefix)
	}

	if len(*keys) != 1 {
		return fmt.Errorf("%s: cryptokeys must contain exactly one entry, got %d", prefix, len(*keys))
	}

	key := (*keys)[0]
	if key == nil {
		return fmt.Errorf("%s: cryptokeys entry is nil", prefix)
	}

	// b) The CryptoKey must be of type "bytes" (TaggedBytes)
	if key.Type() != "bytes" {
		return fmt.Errorf("%s: cryptokeys (signer-id) must be of type 'bytes', got '%s'", prefix, key.Type())
	}

	// c) Byte length must be 32, 48, or 64 (arm-hash-type: SHA-256, SHA-384, or SHA-512)
	tb, ok := key.Value.(comid.TaggedBytes)
	if !ok {
		tbPtr, ok := key.Value.(*comid.TaggedBytes)
		if !ok {
			return fmt.Errorf("%s: failed to extract TaggedBytes from cryptokeys", prefix)
		}
		tb = *tbPtr
	}

	switch len(tb) {
	case 32, 48, 64:
		return nil
	default:
		return fmt.Errorf("%s: signer-id must be 32, 48, or 64 bytes (got %d)", prefix, len(tb))
	}
}

// validateSWComponentDigests validates the digests field (key 2)
func validateSWComponentDigests(digests *comid.Digests) error {
	// Digests field is mandatory
	if digests == nil {
		return fmt.Errorf("digests field is mandatory but not set")
	}

	// Must have at least one digest
	if len(*digests) == 0 {
		return fmt.Errorf("digests must contain at least one entry")
	}

	// Validate each digest entry
	for i, digest := range *digests {
		// For Arm platform, we accept both integer and text algorithm IDs
		// The hash value size is what matters most for validation

		// Hash value must be 32, 48, or 64 bytes (SHA-256, SHA-384, SHA-512)
		if err := ValidateHashDigestSize(digest.Value); err != nil {
			return fmt.Errorf("digest at index %d: %w", i, err)
		}
	}

	return nil
}

// ValidateHashDigestSize validates that the given bytes represent a valid hash digest.
// Valid hash digest sizes are 32, 48, or 64 bytes (SHA-256, SHA-384, or SHA-512).
func ValidateHashDigestSize(val []byte) error {
	if len(val) != 32 && len(val) != 48 && len(val) != 64 {
		return fmt.Errorf("%w: got %d bytes, expected 32, 48, or 64", ErrWrongHashDigestSize, len(val))
	}
	return nil
}

// ValidatePlatformSoftwareComponent validates an Arm Platform software component measurement
func ValidatePlatformSoftwareComponent(measurement *comid.Measurement) error {
	// Validate digests (key 2) - mandatory
	if err := validateSWComponentDigests(measurement.Val.Digests); err != nil {
		return fmt.Errorf("digests: %w", err)
	}

	// Validate cryptokeys (key 13) - mandatory, exactly one signer-id
	if err := validateSignerID(measurement.Val.CryptoKeys, 0, 0); err != nil {
		return fmt.Errorf("cryptokeys: %w", err)
	}

	// version (key 0) is optional - if present, version-scheme MUST NOT be present
	if measurement.Val.Ver != nil {
		if measurement.Val.Ver.Scheme != nil {
			return fmt.Errorf("version-scheme field MUST NOT be present in cca.software-component")
		}
	}
	return nil
}
