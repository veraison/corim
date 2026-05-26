// Copyright 2025-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"errors"
	"fmt"

	"github.com/veraison/swid"
)

// Query is the internal representation of a Query data item
type Query struct {
	ArtifactType        *ArtifactType        `cbor:"0,keyasint,omitempty"`
	EnvironmentSelector *EnvironmentSelector `cbor:"1,keyasint,omitempty"`
	ResultType          *ResultType          `cbor:"2,keyasint,omitempty"`
	RimSelector         *RimSelectorIDs      `cbor:"3,keyasint,omitempty"`
}

// NewEnvironmentQuery creates a new environment Query instance.
// An error is returned if the supplied environment selector is invalid.
func NewEnvironmentQuery(
	artifactType ArtifactType,
	envSelector EnvironmentSelector,
	resultType ResultType,
) (*Query, error) {
	if err := envSelector.Valid(); err != nil {
		return nil, fmt.Errorf("invalid environment selector: %w", err)
	}

	return &Query{
		ArtifactType:        &artifactType,
		EnvironmentSelector: &envSelector,
		ResultType:          &resultType,
	}, nil
}

// NewRimQuery creates a new RIM Query instance. An error is returned if the
// supplied arguments are invalid.
func NewRimQuery(typ RimSelectorType, tagID swid.TagID) (*Query, error) {
	selector, err := NewRimSelectorID(typ, tagID)
	if err != nil {
		return nil, err
	}

	return &Query{RimSelector: NewRimSelectorIDs().Add(selector)}, nil
}

// Valid ensures that the Query target is correctly populated
func (o Query) Valid() error {
	if o.EnvironmentSelector != nil {
		if o.RimSelector != nil {
			return errors.New("environment and RIM selectors cannot be specified at the same time")
		}

		if o.ArtifactType == nil {
			return errors.New("artifact type must be specified with an environment selector")
		}

		if o.ResultType == nil {
			return errors.New("result type must be specified with an environment selector")
		}

		// TODO(tho) add tests for these two:
		// * artifact and result type mismatch should be caught on decoding
		// * ditto for profile syntax errors

		if err := o.EnvironmentSelector.Valid(); err != nil {
			return fmt.Errorf("invalid environment selector: %w", err)
		}

		return nil
	} else if o.RimSelector != nil {
		if o.EnvironmentSelector != nil {
			return errors.New("environment and RIM selectors cannot be specified at the same time")
		}

		if o.ArtifactType != nil {
			return errors.New("artifact type cannot be specified with a RIM selector")
		}

		if o.ResultType != nil {
			return errors.New("result type cannot be specified with a RIM selector")
		}

		if err := o.RimSelector.Valid(); err != nil {
			return fmt.Errorf("invalid RIM selector: %w", err)
		}

		return nil
	} else {
		return errors.New("no selector specified")
	}
}
