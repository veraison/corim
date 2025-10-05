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
	Evidence    CoSWIDEvidence    `json:"coswid-evidence"`
}

func NewCoSWIDTriple() *CoSWIDTriple {
	return &CoSWIDTriple{}
}

func (o *CoSWIDTriple) AddEnvironment(e *comid.Environment) error {
	if e == nil {
		return errors.New("no environment to add")
	}
	if err := e.Valid(); err != nil {
		return fmt.Errorf("environment is not valid: %w", err)
	}

	o.Environment = *e
	return nil
}

func (o *CoSWIDTriple) AddEvidence(e *CoSWIDEvidenceMap) error {

	if e == nil {
		return errors.New("no evidence map to add")
	}

	if len(o.Evidence) == 0 {
		o.Evidence = *NewCoSWIDEvidence()
	}

	o.Evidence = append(o.Evidence, *e)
	return nil
}

func (o CoSWIDTriple) Valid() error {
	if err := o.Environment.Valid(); err != nil {
		return fmt.Errorf("environment validation failed: %w", err)
	}

	if len(o.Evidence) == 0 {
		return errors.New("no evidence entry in the CoSWIDTriple")
	}

	// Validate Evidence entries
	if err := o.Evidence.Valid(); err != nil {
		return fmt.Errorf("evidence validation failed: %w", err)
	}

	return nil
}

type CoSWIDTriples []CoSWIDTriple

func NewCoSWIDTriples() *CoSWIDTriples {
	return &CoSWIDTriples{}
}
