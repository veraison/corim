// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
)

// CCARefValID stores a CCA refval-id with CBOR and JSON serializations
type CCARefValID struct {
	Label *string `cbor:"1,keyasint" json:"label"`
}

// Valid checks the validity (according to the spec) of the target CCARefValID
func (o CCARefValID) Valid() error {
	if o.Label == nil {
		return fmt.Errorf("missing mandatory Label")
	}
	if *o.Label == "" {
		return fmt.Errorf("mandatory Label is empty")
	}
	return nil
}

type TaggedCCARefValID CCARefValID

func NewCCARefValID(label string) *CCARefValID {
	if label == "" {
		return nil
	}

	return &CCARefValID{
		Label: &label,
	}
}

func (o *CCARefValID) SetLabel(label string) error {
	if label == "" {
		return fmt.Errorf("no label supplied")
	}
	if o != nil {
		o.Label = &label
	}
	return nil
}
