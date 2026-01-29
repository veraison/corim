// Copyright 2021-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"encoding/base64"
	"fmt"

	"github.com/veraison/corim/comid"
)

// PSA Implementation ID size constant
const ImplIDSize = 32

// NewImplIDClassID creates a new ClassID containing a PSA Implementation ID.
// The Implementation ID is represented as TaggedBytes (CBOR tag 560) and must be exactly 32 bytes.
// Validation of the 32-byte requirement is performed here; additional profile-level validation
// is done via validatePSAImplementationID().
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
			return nil, fmt.Errorf("bad psa.impl-id: got %d bytes, want %d", nb, ImplIDSize)
		}
		copy(implIDBytes, t)
	case string:
		v, err := base64.StdEncoding.DecodeString(t)
		if err != nil {
			return nil, fmt.Errorf("bad psa.impl-id: %w", err)
		}

		if nb := len(v); nb != ImplIDSize {
			return nil, fmt.Errorf("bad psa.impl-id: decoded %d bytes, want %d", nb, ImplIDSize)
		}
		copy(implIDBytes, v)
	case comid.TaggedBytes:
		if nb := len(t); nb != ImplIDSize {
			return nil, fmt.Errorf("bad psa.impl-id: got %d bytes, want %d", nb, ImplIDSize)
		}
		copy(implIDBytes, t)
	case *comid.TaggedBytes:
		if nb := len(*t); nb != ImplIDSize {
			return nil, fmt.Errorf("bad psa.impl-id: got %d bytes, want %d", nb, ImplIDSize)
		}
		copy(implIDBytes, *t)
	default:
		return nil, fmt.Errorf("unexpected type for psa.impl-id: %T", t)
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

// NewClassImplID instantiates a new Class object with the specified PSA Implementation ID.
// The implID must be exactly 32 bytes. Returns nil if the input is invalid.
func NewClassImplID(implID []byte) *comid.Class {
	classID, err := NewImplIDClassID(implID)
	if err != nil {
		return nil
	}

	return &comid.Class{ClassID: classID}
}
