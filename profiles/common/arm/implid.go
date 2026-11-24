// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package arm

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/veraison/corim/comid"
)

// Arm Platform Implementation ID size constant
const ImplIDSize = 32

const (
	// PlatformImplIDSize is the expected size in bytes of a Arm Platform Implementation ID
	PlatformImplIDSize = 32
)

var (
	ErrWrongImplIDSize = errors.New("wrong Implementation ID size")
)

// NewImplIDClassID creates a new ClassID containing a Arm Platform Implementation ID.
// The Implementation ID is represented as TaggedBytes (CBOR tag 560) and must be exactly 32 bytes.
// Validation of the 32-byte requirement is performed here; additional profile-level validation
// is done via ValidateImplementationID().
//
// Accepted input types:
//   - []byte: raw 32-byte implementation ID
//   - string: base64-encoded implementation ID
//   - comid.TaggedBytes or *comid.TaggedBytes: existing tagged bytes
//   - nil: returns a zero-value 32-byte implementation ID
func NewImplIDClassID(val any) (*comid.ClassID, error) {
	implIDBytes := make([]byte, ImplIDSize)

	if val == nil {
		// Return a zero-value 32-byte Implementation ID as TaggedBytes
		tb := comid.TaggedBytes(implIDBytes)
		return &comid.ClassID{
			Value: &tb,
		}, nil
	}

	switch t := val.(type) {
	case []byte:
		if nb := len(t); nb != ImplIDSize {
			return nil, fmt.Errorf("bad platform.impl-id: got %d bytes, want %d", nb, ImplIDSize)
		}
		copy(implIDBytes, t)
	case string:
		v, err := base64.StdEncoding.DecodeString(t)
		if err != nil {
			return nil, fmt.Errorf("bad platform.impl-id: %w", err)
		}

		if nb := len(v); nb != ImplIDSize {
			return nil, fmt.Errorf("bad platform.impl-id: decoded %d bytes, want %d", nb, ImplIDSize)
		}
		copy(implIDBytes, v)
	case comid.TaggedBytes:
		if nb := len(t); nb != ImplIDSize {
			return nil, fmt.Errorf("bad platform.impl-id: got %d bytes, want %d", nb, ImplIDSize)
		}
		copy(implIDBytes, t)
	case *comid.TaggedBytes:
		if nb := len(*t); nb != ImplIDSize {
			return nil, fmt.Errorf("bad platform.impl-id: got %d bytes, want %d", nb, ImplIDSize)
		}
		copy(implIDBytes, *t)
	default:
		return nil, fmt.Errorf("unexpected type for platform.impl-id: %T", t)
	}

	// Create TaggedBytes which applies CBOR tag 560 during serialization
	tb, err := comid.NewBytes(implIDBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to create impl-ID: %w", err)
	}
	return &comid.ClassID{
		Value: tb,
	}, nil
}

// MustNewImplIDClassID is like NewImplIDClassID except it panics on error
func MustNewImplIDClassID(val any) *comid.ClassID {
	ret, err := NewImplIDClassID(val)
	if err != nil {
		panic(err)
	}

	return ret
}

// NewClassImplID instantiates a new Class object with the specified Arm Platform Implementation ID.
// The implID must be exactly 32 bytes. Returns nil if the input is invalid.
func NewClassImplID(implID []byte) *comid.Class {
	classID, err := NewImplIDClassID(implID)
	if err != nil {
		return nil
	}

	return &comid.Class{ClassID: classID}
}

// ValidateImplementationID validates the Implementation ID
// The Implementation ID must be:
//   - A tagged-bytes type (CBOR tag 560)
//   - Exactly 32 bytes in length
func ValidateImplementationID(env *comid.Environment, prefix string) error {
	// Implementation ID is in Environment.Class.ClassID
	if env.Class == nil {
		return fmt.Errorf("%s: environment.class is required", prefix)
	}

	if env.Class.ClassID == nil {
		return fmt.Errorf("%s: environment.class.id (implementation-id) is required", prefix)
	}

	classID := env.Class.ClassID

	// Must be of type "bytes" (tagged-bytes, CBOR tag 560)
	if classID.Type() != "bytes" {
		return fmt.Errorf("%s: implementation-id must be of type 'bytes', got '%s'", prefix, classID.Type())
	}

	// Must be exactly 32 bytes
	idBytes := classID.Bytes()
	if len(idBytes) != 32 {
		return fmt.Errorf("%s: implementation-id must be exactly 32 bytes (got %d)", prefix, len(idBytes))
	}

	return nil
}
