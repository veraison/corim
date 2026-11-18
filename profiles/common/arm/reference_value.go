package arm

import (
	"fmt"

	"github.com/veraison/corim/comid"
)

// ValidateSignerID validates the cryptokeys field (signer-id)
//
//	a) CryptoKeys field is set and contains exactly one entry
//	b) The type of CryptoKey is TaggedBytes (Type() returns "bytes")
//	c) The length of the value is 32, 48, or 64
func ValidateSignerID(keys *comid.CryptoKeys, tripleIndex, measurementIndex int) error {
	prefix := fmt.Sprintf("reference value at index %d, measurement at index %d", tripleIndex, measurementIndex)

	// a) CryptoKeys field MUST be set (mandatory per PSA profile)
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

	// c) Byte length must be 32, 48, or 64 (psa-hash-type: SHA-256, SHA-384, or SHA-512)
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
