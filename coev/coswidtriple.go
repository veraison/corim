// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"errors"
	"fmt"

	"github.com/veraison/corim/comid"
)

// CoSWIDTriple stores a CoSWID Evidence
// pertaining to an Environment
type CoSWIDTriple struct {
	_           struct{}          `cbor:",toarray"`
	Environment comid.Environment `json:"environment"`
	Evidences   CoSWIDEvidences   `json:"coswid-evidences"`
}

func NewCoSWIDTriple() *CoSWIDTriple {
	return &CoSWIDTriple{}
}

func (o *CoSWIDTriple) AddEnvironment(e *comid.Environment) (*CoSWIDTriple, error) {
	if e == nil {
		return nil, errors.New("no environment to add")
	}
	if err := e.Valid(); err != nil {
		return nil, fmt.Errorf("environment is not valid: %w", err)
	}

	o.Environment = *e
	return o, nil
}

func (o *CoSWIDTriple) AddEvidence(e *CoSWIDEvidenceMap) (*CoSWIDTriple, error) {
	if len(o.Evidences) == 0 {
		o.Evidences = *NewCoSWIDEvidences()
	}
	if e == nil {
		return nil, errors.New("no evidencemap to add")
	}
	o.Evidences = append(o.Evidences, *e)
	return o, nil
}

func (o CoSWIDTriple) Valid() error {
	if err := o.Environment.Valid(); err != nil {
		return fmt.Errorf("environment validation failed: %w", err)
	}

	if len(o.Evidences) == 0 {
		return errors.New("no evidence entry in the CoSWIDTriple")
	}
	return nil
}

type CoSWIDTriples []CoSWIDTriple

func NewCoSWIDTriples() *CoSWIDTriples {
	return &CoSWIDTriples{}
}
