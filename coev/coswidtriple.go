// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"fmt"

	"github.com/veraison/corim/comid"
)

// CoSWIDTriple stores a cryptographic key triple record (identity-triple-record
// or attest-key-triple-record) with CBOR and JSON serializations.  Note that
// the CBOR serialization packs the structure into an array.  Instead, when
// serializing to JSON, the structure is converted into an object.
type CoSWIDTriple struct {
	_           struct{}          `cbor:",toarray"`
	Environment comid.Environment `json:"environment"`
	Evidence    CoSWIDEvidence    `json:"coswid-evidence"`
}

func (o CoSWIDTriple) Valid() error {
	if err := o.Environment.Valid(); err != nil {
		return fmt.Errorf("environment validation failed: %w", err)
	}

	/*
		if err := o.Evidence.Valid(); err != nil {
			return fmt.Errorf("verification keys validation failed: %w", err)
		}
	*/
	return nil
}

type CoSWIDTriples []CoSWIDTriple

func NewCoSWIDTriples() *CoSWIDTriples {
	return &CoSWIDTriples{}
}
