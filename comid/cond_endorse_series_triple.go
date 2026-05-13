// Copyright 2025-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"errors"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/corim/extensions"
)

// CondEndorseSeriesCondition represent the condtion part of a Conditional
// Endorsement Series Triple.
type CondEndorseSeriesCondition struct {
	Environment  Environment  `json:"environment"`
	Measurements Measurements `json:"measurements"`
	AuthorizedBy *CryptoKeys  `json:"authorized-by,omitempty"`
}

func (o *CondEndorseSeriesCondition) RegisterExtensions(exts extensions.Map) error {
	return o.Measurements.RegisterExtensions(exts)
}

func (o *CondEndorseSeriesCondition) GetExtensions() extensions.IMapValue {
	return o.Measurements.GetExtensions()
}

func (o *CondEndorseSeriesCondition) Valid() error {
	if err := o.Environment.Valid(); err != nil {
		return fmt.Errorf("environment validation failed: %w", err)
	}

	if err := o.Measurements.Valid(); err != nil {
		return fmt.Errorf("measurements validation failed: %w", err)
	}

	if o.AuthorizedBy != nil {
		if err := o.AuthorizedBy.Valid(); err != nil {
			return fmt.Errorf("authorized-by validation failed: %w", err)
		}
	}

	return nil
}

func (o *CondEndorseSeriesCondition) MarshalCBOR() ([]byte, error) {
	toMarshal := []any{o.Environment, o.Measurements}
	if o.AuthorizedBy != nil {
		toMarshal = append(toMarshal, o.AuthorizedBy)
	}

	return em.Marshal(toMarshal)
}

func (o *CondEndorseSeriesCondition) UnmarshalCBOR(data []byte) error {
	var raw []cbor.RawMessage
	if err := dm.Unmarshal(data, &raw); err != nil {
		return err
	}

	numElts := len(raw)
	if numElts < 2 || numElts > 3 {
		return fmt.Errorf("expected array between 2 and 3 elements; found %d", numElts)
	}

	if err := dm.Unmarshal(raw[0], &o.Environment); err != nil {
		return fmt.Errorf("environment: %w", err)
	}

	if err := dm.Unmarshal(raw[1], &o.Measurements); err != nil {
		return fmt.Errorf("measurements: %w", err)
	}

	if numElts == 3 {
		if err := dm.Unmarshal(raw[2], &o.AuthorizedBy); err != nil {
			return fmt.Errorf("authorized-by: %w", err)
		}
	}

	return nil
}

// A Conditional Endorsement Series Record, has a series of conditions identified by
// the selection which are matched with the Attester Actual State(from Evidence)
// First successful match terminates matching and corresponding addition are added
// as Endorsements
type CondEndorseSeriesRecord struct {
	_         struct{}     `cbor:",toarray"`
	Selection Measurements `json:"selection"`
	Addition  Measurements `json:"addition"`
}

func (o CondEndorseSeriesRecord) Valid() error {
	if err := o.Selection.Valid(); err != nil {
		return fmt.Errorf("selection validation failed: %w", err)
	}

	if err := o.Addition.Valid(); err != nil {
		return fmt.Errorf("addition validation failed: %w", err)
	}
	return nil
}

// nolint:gocritic
func (o CondEndorseSeriesRecord) GetExtensions() extensions.IMapValue {
	// Extensions are always the same for Selection and Addition
	return o.Selection.GetExtensions()
}

func (o *CondEndorseSeriesRecord) RegisterExtensions(exts extensions.Map) error {
	if err := o.Selection.RegisterExtensions(exts); err != nil {
		return fmt.Errorf("selection: %w", err)
	}
	if err := o.Addition.RegisterExtensions(exts); err != nil {
		return fmt.Errorf("addition: %w", err)
	}

	return nil
}

type CondEndorseSeriesRecords = extensions.Collection[CondEndorseSeriesRecord, *CondEndorseSeriesRecord]

func NewCondEndorseSeriesRecords() *CondEndorseSeriesRecords {
	return extensions.NewCollection[CondEndorseSeriesRecord]()
}

// The Conditional Endorsement Series Triple is used to assert endorsed values based
// on an initial condition match (specified by Condition StatefulEnv) followed by a series
// condition match (specified in selection: inside conditional-series-record).
type CondEndorseSeriesTriple struct {
	_         struct{}                   `cbor:",toarray"`
	Condition CondEndorseSeriesCondition `json:"condition"`
	Series    CondEndorseSeriesRecords   `json:"series"`
}

// RegisterExtensions accepts MVal and MFlag Extension points, that will be registered with
// all Measurements contained within CondEndorseSeriesTriple structure
func (o *CondEndorseSeriesTriple) RegisterExtensions(exts extensions.Map) error {
	if err := o.Condition.RegisterExtensions(exts); err != nil {
		return fmt.Errorf("condition: %w", err)
	}
	if err := o.Series.RegisterExtensions(exts); err != nil {
		return fmt.Errorf("selection: %w", err)
	}

	return nil
}

// nolint:gocritic
func (o CondEndorseSeriesTriple) Valid() error {
	if err := o.Condition.Valid(); err != nil {
		return fmt.Errorf("condition validation failed: %w", err)
	}

	if o.Series.IsEmpty() {
		return errors.New("empty conditional series")
	}

	if err := o.Series.Valid(); err != nil {
		return fmt.Errorf("conditional series validation failed: %w", err)
	}

	return nil
}

func (o *CondEndorseSeriesTriple) GetExtensions() extensions.IMapValue {
	return o.Series.GetExtensions()
}

// CondEndorseSeriesTriples is a container for CondEndorseSeriesTriple instances and their extensions.
// It is a thin wrapper around extensions.Collection.
type CondEndorseSeriesTriples = extensions.Collection[CondEndorseSeriesTriple, *CondEndorseSeriesTriple]

func NewCondEndorseSeriesTriples() *CondEndorseSeriesTriples {
	return extensions.NewCollection[CondEndorseSeriesTriple]()
}
