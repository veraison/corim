// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"fmt"
)

// Query is the internal representation of a Query data item
type Query struct {
	ArtifactType        ArtifactType        `cbor:"0,keyasint"`
	EnvironmentSelector EnvironmentSelector `cbor:"1,keyasint"`
	ResultType          ResultType          `cbor:"2,keyasint"`
}

// NewQuery creates a new Query instance.
// An error is returned if the supplied environment selector is invalid.
func NewQuery(
	artifactType ArtifactType,
	envSelector EnvironmentSelector,
	resultType ResultType,
) (*Query, error) {
	if err := envSelector.Valid(); err != nil {
		return nil, fmt.Errorf("invalid environment selector: %w", err)
	}

	return &Query{
		ArtifactType:        artifactType,
		EnvironmentSelector: envSelector,
		ResultType:          resultType,
	}, nil
}

// Valid ensures that the Query target is correctly populated
func (o Query) Valid() error {
	// TODO(tho) add tests for these two:
	// * artifact and result type mismatch should be caught on decoding
	// * ditto for profile syntax errors

	if err := o.EnvironmentSelector.Valid(); err != nil {
		return fmt.Errorf("invalid environment selector: %w", err)
	}

	return nil
}
