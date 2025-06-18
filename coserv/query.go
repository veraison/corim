// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"fmt"
	"time"
)

// Query is the internal representation of a Query data item
type Query struct {
	ArtifactType        ArtifactType        `cbor:"0,keyasint"`
	EnvironmentSelector EnvironmentSelector `cbor:"1,keyasint"`
	Timestamp           time.Time           `cbor:"2,keyasint"`
	ResultType          ResultType          `cbor:"3,keyasint"`
}

// NewQuery creates a new Query instance with the timestamp set to instantiation time.
// (If needed, the timestamp can be changed using SetTimestamp.)
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
		Timestamp:           time.Now(),
		ResultType:          resultType,
	}, nil
}

// SetTimestamp allows setting an explicit timestamp for the target Query object
func (o *Query) SetTimestamp(ts time.Time) *Query {
	o.Timestamp = ts
	return o
}

// Valid ensures that the Query target is correctly populated
func (o Query) Valid() error {
	// TODO(tho) add tests for these two:
	// * artifact and result type mismatch should be caught on decoding
	// * ditto for profile syntax errors

	if err := o.EnvironmentSelector.Valid(); err != nil {
		return fmt.Errorf("invalid environment selector: %w", err)
	}

	zeroTime := time.Time{}
	if o.Timestamp.Equal(zeroTime) {
		return fmt.Errorf("timestamp not set")
	}

	return nil
}
