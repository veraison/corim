// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package cca

import (
	"errors"
	"fmt"

	"github.com/veraison/corim/comid"
)

const (
	// PlatformImplIDSize is the expected size in bytes of a CCA Platform Implementation ID
	PlatformImplIDSize = 32
	// PlatformInstanceIDSize is the expected size in bytes of a CCA Platform Instance ID
	PlatformInstanceIDSize = 33
	// RealmRIMSize is the expected size in bytes of a CCA Realm Initial Measurement (RIM)
	RealmRIMSize = 32
)

var (
	ErrWrongImplIDSize     = errors.New("wrong Implementation ID size")
	ErrWrongInstanceIDSize = errors.New("wrong Instance ID size")
	ErrWrongInstancePrefix = errors.New("instance ID must start with 0x01")
	ErrWrongRIMSize        = errors.New("wrong Realm Initial Measurement (RIM) size")
)

// NewPlatformImplIDClassID creates a new CCA Platform Implementation ID as a ClassID.
// The Implementation ID MUST be exactly 32 bytes.
func NewPlatformImplIDClassID(val []byte) (*comid.ClassID, error) {
	if len(val) != PlatformImplIDSize {
		return nil, fmt.Errorf("%w: got %d bytes, expected %d", ErrWrongImplIDSize, len(val), PlatformImplIDSize)
	}

	return comid.NewClassID(val, comid.BytesType)
}

// MustNewPlatformImplIDClassID is like NewPlatformImplIDClassID but panics on error.
func MustNewPlatformImplIDClassID(val []byte) *comid.ClassID {
	c, err := NewPlatformImplIDClassID(val)
	if err != nil {
		panic(err)
	}
	return c
}

// NewClassPlatformImplID creates a new CCA Platform Implementation ID as a Class.
// The Implementation ID MUST be exactly 32 bytes.
func NewClassPlatformImplID(val []byte) (*comid.Class, error) {
	if len(val) != PlatformImplIDSize {
		return nil, fmt.Errorf("%w: got %d bytes, expected %d", ErrWrongImplIDSize, len(val), PlatformImplIDSize)
	}

	return comid.NewClassBytes(val), nil
}

// MustNewClassPlatformImplID is like NewClassPlatformImplID but panics on error.
func MustNewClassPlatformImplID(val []byte) *comid.Class {
	c, err := NewClassPlatformImplID(val)
	if err != nil {
		panic(err)
	}
	return c
}

// NewPlatformInstanceID creates a new CCA Platform Instance ID.
// The Instance ID MUST be exactly 33 bytes and start with 0x01.
func NewPlatformInstanceID(val []byte) (*comid.Instance, error) {
	if len(val) != PlatformInstanceIDSize {
		return nil, fmt.Errorf("%w: got %d bytes, expected %d", ErrWrongInstanceIDSize, len(val), PlatformInstanceIDSize)
	}

	if val[0] != 0x01 {
		return nil, fmt.Errorf("%w: got 0x%02x", ErrWrongInstancePrefix, val[0])
	}

	return comid.NewUEIDInstance(comid.UEID(val))
}

// MustNewPlatformInstanceID is like NewPlatformInstanceID but panics on error.
func MustNewPlatformInstanceID(val []byte) *comid.Instance {
	i, err := NewPlatformInstanceID(val)
	if err != nil {
		panic(err)
	}
	return i
}

// NewInstancePlatformInstanceID creates a new CCA Platform Instance ID as an Instance.
// The Instance ID MUST be exactly 33 bytes and start with 0x01.
// See section 3.2.1 of the IETF draft.
func NewInstancePlatformInstanceID(val []byte) (*comid.Instance, error) {
	return NewPlatformInstanceID(val)
}

// MustNewInstancePlatformInstanceID is like NewInstancePlatformInstanceID but panics on error.
func MustNewInstancePlatformInstanceID(val []byte) *comid.Instance {
	i, err := NewInstancePlatformInstanceID(val)
	if err != nil {
		panic(err)
	}
	return i
}

// ValidatePlatformImplID validates that the given bytes represent a valid CCA Platform Implementation ID.
func ValidatePlatformImplID(val []byte) error {
	if len(val) != PlatformImplIDSize {
		return fmt.Errorf("%w: got %d bytes, expected %d", ErrWrongImplIDSize, len(val), PlatformImplIDSize)
	}
	return nil
}

// ValidatePlatformInstanceID validates that the given bytes represent a valid CCA Platform Instance ID.
func ValidatePlatformInstanceID(val []byte) error {
	if len(val) != PlatformInstanceIDSize {
		return fmt.Errorf("%w: got %d bytes, expected %d", ErrWrongInstanceIDSize, len(val), PlatformInstanceIDSize)
	}

	if val[0] != 0x01 {
		return fmt.Errorf("%w: got 0x%02x", ErrWrongInstancePrefix, val[0])
	}

	return nil
}

// NewRealmRIMClassID creates a new CCA Realm Initial Measurement (RIM) as a ClassID.
// The RIM can be 32, 48, or 64 bytes (SHA-256, SHA-384, or SHA-512 hash digest).
// RIM uniquely identifies a Realm Target Environment.
func NewRealmRIMClassID(val []byte) (*comid.ClassID, error) {
	if len(val) != 32 && len(val) != 48 && len(val) != 64 {
		return nil, fmt.Errorf("%w: got %d bytes, expected 32, 48, or 64", ErrWrongRIMSize, len(val))
	}

	return comid.NewClassID(val, comid.BytesType)
}

// MustNewRealmRIMClassID is like NewRealmRIMClassID but panics on error.
func MustNewRealmRIMClassID(val []byte) *comid.ClassID {
	c, err := NewRealmRIMClassID(val)
	if err != nil {
		panic(err)
	}
	return c
}

// NewClassRealmRIM creates a new CCA Realm Initial Measurement (RIM) as a Class.
// The RIM can be 32, 48, or 64 bytes (SHA-256, SHA-384, or SHA-512 hash digest).
// RIM uniquely identifies a Realm Target Environment.
func NewClassRealmRIM(val []byte) (*comid.Class, error) {
	if len(val) != 32 && len(val) != 48 && len(val) != 64 {
		return nil, fmt.Errorf("%w: got %d bytes, expected 32, 48, or 64", ErrWrongRIMSize, len(val))
	}

	return comid.NewClassBytes(val), nil
}

// MustNewClassRealmRIM is like NewClassRealmRIM but panics on error.
func MustNewClassRealmRIM(val []byte) *comid.Class {
	c, err := NewClassRealmRIM(val)
	if err != nil {
		panic(err)
	}
	return c
}

// ValidateRealmRIM validates that the given bytes represent a valid CCA Realm Initial Measurement.
// The RIM can be 32, 48, or 64 bytes (SHA-256, SHA-384, or SHA-512 hash digest).
func ValidateRealmRIM(val []byte) error {
	if len(val) != 32 && len(val) != 48 && len(val) != 64 {
		return fmt.Errorf("%w: got %d bytes, expected 32, 48, or 64", ErrWrongRIMSize, len(val))
	}
	return nil
}
