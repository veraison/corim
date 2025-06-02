// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
)

// Query is the internal representation of a Query data item
type Query struct {
	ArtifactType        ArtifactType        `cbor:"0,keyasint"`
	EnvironmentSelector EnvironmentSelector `cbor:"1,keyasint"`
	Timestamp           time.Time           `cbor:"2,keyasint"`
}

// NewQuery creates a new Query instance with the timestamp set to instantiation time.
// (If needed, the timestamp can be changed using SetTimestamp.)
// An error is returned if the supplied environment selector is invalid.
func NewQuery(artifactType ArtifactType, envSelector EnvironmentSelector) (*Query, error) {
	if err := envSelector.Valid(); err != nil {
		return nil, fmt.Errorf("invalid environment selector: %w", err)
	}

	return &Query{
		ArtifactType:        artifactType,
		EnvironmentSelector: envSelector,
		Timestamp:           time.Now(),
	}, nil
}

// SetTimestamp allows setting an explicit timestamp for the target Query object
func (o *Query) SetTimestamp(ts time.Time) *Query {
	o.Timestamp = ts
	return o
}

// Valid ensures that the Query target is correctly populated
func (o Query) Valid() error {
	// TBC:
	// * artifact type mismatch should be caught on decoding
	// * ditto for profile syntax errors
	if err := o.EnvironmentSelector.Valid(); err != nil {
		return fmt.Errorf("invalid environment selector: %w", err)
	}

	zeroTime := time.Time{}
	if o.Timestamp == zeroTime {
		return fmt.Errorf("timestamp not set")
	}

	return nil
}

// ToCBOR validates and serializes to CBOR the target Query.
// An error is returned if either validation or encoding of the Query target fails
func (o Query) ToCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, fmt.Errorf("validating Query: %w", err)
	}

	data, err := cbor.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("encoding Query to CBOR: %w", err)
	}

	return data, nil
}

// FromCBOR deserializes from CBOR into the target Query.
// An error is returned if either decoding or validation of the Query payload fails
func (o *Query) FromCBOR(data []byte) error {
	if err := cbor.Unmarshal(data, o); err != nil {
		return fmt.Errorf("decoding Query from CBOR: %w", err)
	}

	if err := o.Valid(); err != nil {
		return fmt.Errorf("validating Query: %w", err)
	}

	return nil
}

func (o Query) ToSQL() (string, error) {
	condition, err := o.EnvironmentSelector.ToSQL()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("SELECT * FROM %s WHERE %s", o.ArtifactType, condition), nil
}
