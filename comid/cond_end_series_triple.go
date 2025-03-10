// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"

	"github.com/veraison/corim/extensions"
)

// A Stateful Environment is an Environment in a known reference state
type StatefulEnv = ValueTriple

// A Conditional Series Record, has a series of conditions identified by
// the selection which are matched with the Attester Actual State(from Evidence)
// First successful match terminates matching and corresponding addition are added
// as Endorsements
type CondSeriesRecord struct {
	Selection Measurements `json:"selection"`
	Addition  Measurements `json:"addition"`
}

func (o CondSeriesRecord) Valid() error {
	if err := o.Selection.Valid(); err != nil {
		return fmt.Errorf("conditional series record selection validation failed: %w", err)
	}

	if err := o.Addition.Valid(); err != nil {
		return fmt.Errorf("conditional series record addition validation failed: %w", err)
	}
	return nil
}

// nolint:gocritic
func (o CondSeriesRecord) GetExtensions() extensions.IMapValue {
	// Extensions are always the same for Selection and Addition
	return o.Selection.GetExtensions()
}

func (o *CondSeriesRecord) RegisterExtensions(exts extensions.Map) error {
	if err := o.Selection.RegisterExtensions(exts); err != nil {
		return fmt.Errorf("selection: %w", err)
	}
	if err := o.Addition.RegisterExtensions(exts); err != nil {
		return fmt.Errorf("selection: %w", err)
	}

	return nil
}

type CondSeriesRecords extensions.Collection[CondSeriesRecord, *CondSeriesRecord]

func (o *CondSeriesRecords) GetExtensions() extensions.IMapValue {
	return (*extensions.Collection[CondSeriesRecord, *CondSeriesRecord])(o).GetExtensions()
}

func (o *CondSeriesRecords) RegisterExtensions(exts extensions.Map) error {
	return (*extensions.Collection[CondSeriesRecord, *CondSeriesRecord])(o).RegisterExtensions(exts)
}

func (o *CondSeriesRecords) Valid() error {
	return (*extensions.Collection[CondSeriesRecord, *CondSeriesRecord])(o).Valid()
}

// The Conditional Endorsement Series Triple is used to assert endorsed values based
// on an initial condition match (specified in condition:) followed by a series
// condition match (specified in selection: inside conditional-series-record).
type CondEndSeriesTriple struct {
	_         struct{}          `cbor:",toarray"`
	Condition StatefulEnv       `json:"statefulenv"`
	Series    CondSeriesRecords `json:"series"`
}

// RegisterExtensions accepts MVal and MFlag Extension points, that will be registered with
// all Measurements contained within CondEndSeriesTriple structure
func (o *CondEndSeriesTriple) RegisterExtensions(exts extensions.Map) error {
	if err := o.Condition.RegisterExtensions(exts); err != nil {
		return fmt.Errorf("selection: %w", err)
	}
	if err := o.Series.RegisterExtensions(exts); err != nil {
		return fmt.Errorf("selection: %w", err)
	}

	return nil
}

// nolint:gocritic
func (o CondEndSeriesTriple) Valid() error {
	fmt.Printf("Yogesh: Valid Called")
	if err := o.Condition.Valid(); err != nil {
		return fmt.Errorf("stateful environment validation failed: %w", err)
	}
	if err := o.Series.Valid(); err != nil {
		return fmt.Errorf("conditional series validation failed: %w", err)
	}

	return nil
}

func (o *CondEndSeriesTriple) GetExtensions() extensions.IMapValue {
	return o.Series.GetExtensions()
}

// CondEndSeriesTriples is a container for CondEndSeriesTriple instances and their extensions.
// It is a thin wrapper around extensions.Collection.
type CondEndSeriesTriples extensions.Collection[CondEndSeriesTriple, *CondEndSeriesTriple]

func NewCondEndSeriesTriples() *CondEndSeriesTriples {
	return (*CondEndSeriesTriples)(extensions.NewCollection[CondEndSeriesTriple]())
}

func (o *CondEndSeriesTriples) GetExtensions() extensions.IMapValue {
	return (*extensions.Collection[CondEndSeriesTriple, *CondEndSeriesTriple])(o).GetExtensions()
}

func (o *CondEndSeriesTriples) RegisterExtensions(exts extensions.Map) error {
	return (*extensions.Collection[CondEndSeriesTriple, *CondEndSeriesTriple])(o).RegisterExtensions(exts)
}

func (o CondEndSeriesTriples) Valid() error {
	return (extensions.Collection[CondEndSeriesTriple, *CondEndSeriesTriple])(o).Valid()
}

func (o *CondEndSeriesTriples) IsEmpty() bool {
	return (*extensions.Collection[CondEndSeriesTriple, *CondEndSeriesTriple])(o).IsEmpty()
}

func (o *CondEndSeriesTriples) Add(val *CondEndSeriesTriple) *CondEndSeriesTriples {
	ret := (*extensions.Collection[CondEndSeriesTriple, *CondEndSeriesTriple])(o).Add(val)
	return (*CondEndSeriesTriples)(ret)
}

func (o CondEndSeriesTriples) MarshalCBOR() ([]byte, error) {
	return (extensions.Collection[CondEndSeriesTriple, *CondEndSeriesTriple])(o).MarshalCBOR()
}

func (o *CondEndSeriesTriples) UnmarshalCBOR(data []byte) error {
	return (*extensions.Collection[CondEndSeriesTriple, *CondEndSeriesTriple])(o).UnmarshalCBOR(data)
}

func (o CondEndSeriesTriples) MarshalJSON() ([]byte, error) {
	return (extensions.Collection[CondEndSeriesTriple, *CondEndSeriesTriple])(o).MarshalJSON()
}

func (o *CondEndSeriesTriples) UnmarshalJSON(data []byte) error {
	return (*extensions.Collection[CondEndSeriesTriple, *CondEndSeriesTriple])(o).UnmarshalJSON(data)
}
