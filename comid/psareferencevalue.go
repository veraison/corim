// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
)

// PSARefValID stores a PSA refval-id with CBOR and JSON serializations
// (See https://datatracker.ietf.org/doc/html/draft-xyz-rats-psa-endorsements)
type PSARefValID struct {
	Label    *string `cbor:"1,keyasint,omitempty" json:"label,omitempty"`
	Version  *string `cbor:"4,keyasint,omitempty" json:"version,omitempty"`
	SignerID []byte  `cbor:"5,keyasint" json:"signer-id"` // 32, 48 or 64
}

// Valid checks the validity (according to the spec) of the target PSARefValID
func (o PSARefValID) Valid() error {
	if o.SignerID == nil {
		return fmt.Errorf("missing mandatory signer ID")
	}

	switch len(o.SignerID) {
	case 32, 48, 64:
	default:
		return fmt.Errorf("want 32, 48 or 64 bytes, got %d", len(o.SignerID))
	}

	return nil
}

type TaggedPSARefValID PSARefValID

func NewPSARefValID(signerID []byte) *PSARefValID {
	switch len(signerID) {
	case 32, 48, 64:
	default:
		return nil
	}

	return &PSARefValID{
		SignerID: signerID,
	}
}

func (o *PSARefValID) SetLabel(label string) *PSARefValID {
	if o != nil {
		o.Label = &label
	}
	return o
}

func (o *PSARefValID) SetVersion(version string) *PSARefValID {
	if o != nil {
		o.Version = &version
	}
	return o
}
