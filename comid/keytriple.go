// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "fmt"

// KeyTriple stores a cryptographic key triple record (identity-triple-record
// or attest-key-triple-record) with CBOR and JSON serializations.  Note that
// the CBOR serialization packs the structure into an array.  Instead, when
// serializing to JSON, the structure is converted into an object.
type KeyTriple struct {
	_           struct{}    `cbor:",toarray"`
	Environment Environment `json:"environment"`
	VerifKeys   CryptoKeys  `json:"verification-keys"`
}

func (o KeyTriple) Valid() error {
	if err := o.Environment.Valid(); err != nil {
		return fmt.Errorf("environment validation failed: %w", err)
	}

	if err := o.VerifKeys.Valid(); err != nil {
		return fmt.Errorf("verification keys validation failed: %w", err)
	}
	return nil
}

type KeyTriples []KeyTriple

func NewKeyTriples() *KeyTriples {
	return &KeyTriples{}
}
