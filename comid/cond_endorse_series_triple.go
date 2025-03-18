// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"

	"github.com/veraison/corim/extensions"
)

// A Stateful Environment is an Environment in a known reference state
type StatefulEnv = ValueTriple

// A Conditional Endorsement Series Record, has a series of conditions identified by
// the selection which are matched with the Attester Actual State(from Evidence)
// First successful match terminates matching and corresponding addition are added
// as Endorsements
type CondEndorseSeriesRecord struct {
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

type CondEndorseSeriesRecords extensions.Collection[CondEndorseSeriesRecord, *CondEndorseSeriesRecord]

func NewCondEndorseSeriesRecords() *CondEndorseSeriesRecords {
	return (*CondEndorseSeriesRecords)(extensions.NewCollection[CondEndorseSeriesRecord]())
}

func (o *CondEndorseSeriesRecords) IsEmpty() bool {
	return (*extensions.Collection[CondEndorseSeriesRecord, *CondEndorseSeriesRecord])(o).IsEmpty()
}

func (o *CondEndorseSeriesRecords) Add(val *CondEndorseSeriesRecord) *CondEndorseSeriesRecords {
	ret := (*extensions.Collection[CondEndorseSeriesRecord, *CondEndorseSeriesRecord])(o).Add(val)
	return (*CondEndorseSeriesRecords)(ret)
}

func (o *CondEndorseSeriesRecords) GetExtensions() extensions.IMapValue {
	return (*extensions.Collection[CondEndorseSeriesRecord, *CondEndorseSeriesRecord])(o).GetExtensions()
}

func (o *CondEndorseSeriesRecords) RegisterExtensions(exts extensions.Map) error {
	return (*extensions.Collection[CondEndorseSeriesRecord, *CondEndorseSeriesRecord])(o).RegisterExtensions(exts)
}

func (o *CondEndorseSeriesRecords) Valid() error {
	return (*extensions.Collection[CondEndorseSeriesRecord, *CondEndorseSeriesRecord])(o).Valid()
}

// The Conditional Endorsement Series Triple is used to assert endorsed values based
// on an initial condition match (specified by Condition StatefulEnv) followed by a series
// condition match (specified in selection: inside conditional-series-record).
type CondEndorseSeriesTriple struct {
	_         struct{}                 `cbor:",toarray"`
	Condition StatefulEnv              `json:"statefulenv"`
	Series    CondEndorseSeriesRecords `json:"series"`
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
		return fmt.Errorf("stateful environment validation failed: %w", err)
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
type CondEndorseSeriesTriples extensions.Collection[CondEndorseSeriesTriple, *CondEndorseSeriesTriple]

func NewCondEndorseSeriesTriples() *CondEndorseSeriesTriples {
	return (*CondEndorseSeriesTriples)(extensions.NewCollection[CondEndorseSeriesTriple]())
}

func (o *CondEndorseSeriesTriples) GetExtensions() extensions.IMapValue {
	return (*extensions.Collection[CondEndorseSeriesTriple, *CondEndorseSeriesTriple])(o).GetExtensions()
}

func (o *CondEndorseSeriesTriples) RegisterExtensions(exts extensions.Map) error {
	return (*extensions.Collection[CondEndorseSeriesTriple, *CondEndorseSeriesTriple])(o).RegisterExtensions(exts)
}

func (o CondEndorseSeriesTriples) Valid() error {
	return (extensions.Collection[CondEndorseSeriesTriple, *CondEndorseSeriesTriple])(o).Valid()
}

func (o *CondEndorseSeriesTriples) IsEmpty() bool {
	return (*extensions.Collection[CondEndorseSeriesTriple, *CondEndorseSeriesTriple])(o).IsEmpty()
}

func (o *CondEndorseSeriesTriples) Add(val *CondEndorseSeriesTriple) *CondEndorseSeriesTriples {
	ret := (*extensions.Collection[CondEndorseSeriesTriple, *CondEndorseSeriesTriple])(o).Add(val)
	return (*CondEndorseSeriesTriples)(ret)
}

func (o CondEndorseSeriesTriples) MarshalCBOR() ([]byte, error) {
	return (extensions.Collection[CondEndorseSeriesTriple, *CondEndorseSeriesTriple])(o).MarshalCBOR()
}

func (o *CondEndorseSeriesTriples) UnmarshalCBOR(data []byte) error {
	return (*extensions.Collection[CondEndorseSeriesTriple, *CondEndorseSeriesTriple])(o).UnmarshalCBOR(data)
}

func (o CondEndorseSeriesTriples) MarshalJSON() ([]byte, error) {
	return (extensions.Collection[CondEndorseSeriesTriple, *CondEndorseSeriesTriple])(o).MarshalJSON()
}

func (o *CondEndorseSeriesTriples) UnmarshalJSON(data []byte) error {
	return (*extensions.Collection[CondEndorseSeriesTriple, *CondEndorseSeriesTriple])(o).UnmarshalJSON(data)
}
