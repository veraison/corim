// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

// Package coserv provides an implementation of draft-howard-rats-coserv
package coserv

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/eat"
)

// Coserv is the internal representation of a CoSERV data item
type Coserv struct {
	ArtifactType        ArtifactType        `cbor:"0,keyasint"`
	Profile             eat.Profile         `cbor:"1,keyasint"`
	EnvironmentSelector EnvironmentSelector `cbor:"2,keyasint"`
}

func NewCoserv(artifactType ArtifactType, profile string, envSelector EnvironmentSelector) (*Coserv, error) {
	p, err := eat.NewProfile(profile)
	if err != nil {
		return nil, fmt.Errorf("invalid profile: %w", err)
	}

	if err := envSelector.Valid(); err != nil {
		return nil, fmt.Errorf("invalid environment selector: %w", err)
	}

	return &Coserv{
		ArtifactType:        artifactType,
		Profile:             *p,
		EnvironmentSelector: envSelector,
	}, nil
}

// Valid ensures that the Coserv target is correctly populated
func (o Coserv) Valid() error {
	// TBC:
	// * artifact type mismatch should be caught on decoding
	// * ditto for profile syntax errors
	if err := o.EnvironmentSelector.Valid(); err != nil {
		return fmt.Errorf("invalid environment selector: %w", err)
	}
	return nil
}

// ToCBOR validates and serializes to CBOR the target Coserv
// An error is returned if either validation or encoding of the Coserv target fails
func (o Coserv) ToCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return cbor.Marshal(o)
}

// FromCBOR deserializes from CBOR into the target Coserv
// An error is returned if either decoding or validation of the CoSERV payload fails
func (o *Coserv) FromCBOR(data []byte) error {
	if err := cbor.Unmarshal(data, o); err != nil {
		return fmt.Errorf("decoding CoSERV: %w", err)
	}

	if err := o.Valid(); err != nil {
		return fmt.Errorf("validating CoSERV: %w", err)
	}

	return nil
}
