package arm

import (
	"errors"
	"fmt"

	"github.com/veraison/corim/comid"
)

const (
	// PlatformInstanceIDSize is the expected size in bytes of a Platform Instance ID
	PlatformInstanceIDSize = 33
)

var (
	ErrWrongInstancePrefix = errors.New("instance ID must start with 0x01")
	ErrWrongHashDigestSize = errors.New("wrong hash digest size")
	ErrWrongInstanceIDSize = errors.New("wrong Instance ID size")
)

// NewInstancePlatformInstanceID creates a new Platform Instance ID as an Instance.
// The Instance ID MUST be exactly 33 bytes and start with 0x01.
// See section 3.2.1 of the IETF draft.
func NewInstancePlatformInstanceID(val []byte) (*comid.Instance, error) {
	if err := ValidatePlatformInstanceID(val); err != nil {
		return nil, err
	}

	return comid.NewUEIDInstance(comid.UEID(val))
}

// ValidatePlatformInstanceID validates that the given bytes represent a valid Platform Instance ID.
func ValidatePlatformInstanceID(val []byte) error {
	if len(val) != PlatformInstanceIDSize {
		return fmt.Errorf("%w: got %d bytes, expected %d", ErrWrongInstanceIDSize, len(val), PlatformInstanceIDSize)
	}

	if val[0] != 0x01 {
		return fmt.Errorf("%w: got 0x%02x", ErrWrongInstancePrefix, val[0])
	}

	return nil
}

// ValidateInstanceID validates the Instance ID (UEID) per Platform. Endorsements Section 3.2
// The Instance ID must be:
//   - A tagged-ueid-type (CBOR tag 550)
//   - The first byte MUST be 0x01 (RAND type)
//   - Followed by exactly 32 bytes (total 33 bytes)
func ValidateInstanceID(env *comid.Environment, prefix string) error {
	// Instance ID is in Environment.Instance
	if env.Instance == nil {
		return fmt.Errorf("%s: environment.instance (instance-id) is required", prefix)
	}

	instance := env.Instance

	// Must be of type "ueid" (tagged-ueid-type, CBOR tag 550)
	if instance.Type() != comid.UEIDType {
		return fmt.Errorf("%s: instance-id must be of type '%s', got '%s'",
			prefix, comid.UEIDType, instance.Type())
	}

	// Get the UEID bytes
	ueidBytes := instance.Bytes()

	// Platform Instance ID must be exactly 33 bytes (1 byte RAND type + 32 bytes identifier)
	if len(ueidBytes) != 33 {
		return fmt.Errorf("%s: instance-id must be exactly 33 bytes (got %d)", prefix, len(ueidBytes))
	}

	// The first byte MUST be 0x01 (RAND type per EAT UEID specification)
	if ueidBytes[0] != 0x01 {
		return fmt.Errorf("%s: instance-id must have RAND type (0x01), got 0x%02x", prefix, ueidBytes[0])
	}

	return nil
}
