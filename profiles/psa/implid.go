// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/veraison/corim/comid"
)

const ImplIDType = "psa.impl-id"

// TaggedImplID is the PSA Implementation ID type (32 bytes) as a ClassID value.
// It implements IClassIDValue interface with Type() returning "psa.impl-id".
// See Section 3.2.2 of draft-tschofenig-rats-psa-token
type TaggedImplID [32]byte

// String returns the base64-encoded string representation of the ImplID
func (o TaggedImplID) String() string {
	return base64.StdEncoding.EncodeToString(o[:])
}

// Valid validates the ImplID (always returns nil as any 32-byte value is valid)
func (o TaggedImplID) Valid() error {
	return nil
}

// Type returns the type identifier for PSA Implementation ID
func (o TaggedImplID) Type() string {
	return ImplIDType
}

// Bytes returns the raw bytes of the ImplID
func (o TaggedImplID) Bytes() []byte {
	return o[:]
}

// MarshalJSON serializes the TaggedImplID to JSON
func (o TaggedImplID) MarshalJSON() ([]byte, error) {
	return json.Marshal(o[:])
}

// UnmarshalJSON deserializes JSON into TaggedImplID
func (o *TaggedImplID) UnmarshalJSON(data []byte) error {
	var b []byte
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if len(b) != 32 {
		return fmt.Errorf("bad psa.impl-id: got %d bytes, want 32", len(b))
	}
	copy(o[:], b)
	return nil
}

// ImplID is an alias for backward compatibility
type ImplID = TaggedImplID

// NewImplIDClassID creates a new ClassID of type psa.impl-id
func NewImplIDClassID(val any) (*comid.ClassID, error) {
	var implID TaggedImplID

	if val == nil {
		return &comid.ClassID{
			Value: &implID,
		}, nil
	}

	switch t := val.(type) {
	case []byte:
		if nb := len(t); nb != 32 {
			return nil, fmt.Errorf("bad psa.impl-id: got %d bytes, want 32", nb)
		}

		copy(implID[:], t)
	case string:
		v, err := base64.StdEncoding.DecodeString(t)
		if err != nil {
			return nil, fmt.Errorf("bad psa.impl-id: %w", err)
		}

		if nb := len(v); nb != 32 {
			return nil, fmt.Errorf("bad psa.impl-id: decoded %d bytes, want 32", nb)
		}

		copy(implID[:], v)
	case TaggedImplID:
		copy(implID[:], t[:])
	case *TaggedImplID:
		copy(implID[:], (*t)[:])
	default:
		return nil, fmt.Errorf("unexpected type for psa.impl-id: %T", t)
	}

	return &comid.ClassID{
		Value: &implID,
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

// NewClassImplID instantiates a new Class object with the specified PSA Implementation ID
// This is a convenience function for use in PSA profiles only
func NewClassImplID(implID TaggedImplID) *comid.Class {
	classID, err := NewImplIDClassID(implID)
	if err != nil {
		return nil
	}

	return &comid.Class{ClassID: classID}
}
