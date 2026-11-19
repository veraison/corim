// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package cca

import (
	"errors"
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/profiles/common/arm"
)

const (

	// RealmRIMSize is the expected size in bytes of a CCA Realm Initial Measurement (RIM)
	RealmRIMSize = 32
)

var (
	ErrWrongRIMSize = errors.New("wrong Realm Initial Measurement (RIM) size")
)

// NewRealmRIMClassID creates a new CCA Realm Initial Measurement (RIM) as a ClassID.
// The RIM can be 32, 48, or 64 bytes (SHA-256, SHA-384, or SHA-512 hash digest).
// RIM uniquely identifies a Realm Target Environment.
func NewRealmRIMClassID(val []byte) (*comid.ClassID, error) {
	if err := ValidateRealmRIM(val); err != nil {
		return nil, err
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
	if err := ValidateRealmRIM(val); err != nil {
		return nil, err
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
	if err := ValidateHashDigestSize(val); err != nil {
		return fmt.Errorf("%w: %w", ErrWrongRIMSize, err)
	}
	return nil
}

// ValidateHashDigestSize validates that the given bytes represent a valid hash digest.
// Valid hash digest sizes are 32, 48, or 64 bytes (SHA-256, SHA-384, or SHA-512).
func ValidateHashDigestSize(val []byte) error {
	if len(val) != 32 && len(val) != 48 && len(val) != 64 {
		return fmt.Errorf("%w: got %d bytes, expected 32, 48, or 64", arm.ErrWrongHashDigestSize, len(val))
	}
	return nil
}
