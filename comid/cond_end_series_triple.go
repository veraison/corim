// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "fmt"

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

// The Conditional Endorsement Series Triple is used to assert endorsed values based
// on an initial condition match (specified in condition:) followed by a series
// condition match (specified in selection: inside conditional-series-record).
type CondEndSeriesTriple struct {
	_         struct{}           `cbor:",toarray"`
	Condition StatefulEnv        `json:"statefulenv"`
	Series    []CondSeriesRecord `json:"series"`
}

func (o CondEndSeriesTriple) Valid() error {
	if err := o.Condition.Valid(); err != nil {
		return fmt.Errorf("stateful environment validation failed: %w", err)
	}

	for i, elem := range o.Series {
		if err := elem.Valid(); err != nil {
			return fmt.Errorf("series validation at index: %d failed: %w", i, err)
		}
	}
	return nil
}

type CondEndSeriesTriples []CondEndSeriesTriple

func NewCondEndSeriesTriples() *CondEndSeriesTriples {
	return &CondEndSeriesTriples{}
}
