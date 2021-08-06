// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "fmt"

// VerifKeys is an array of VerifKey
type VerifKeys []VerifKey

// NewVerifKeys instantiates an empty VerifKeys array
func NewVerifKeys() *VerifKeys {
	return new(VerifKeys)
}

// AddVerifKey adds the supplied VerifKey to the target VerifKeys array
func (o *VerifKeys) AddVerifKey(v *VerifKey) *VerifKeys {
	if o != nil && v != nil {
		*o = append(*o, *v)
	}
	return o
}

func (o VerifKeys) Valid() error {
	if len(o) == 0 {
		return fmt.Errorf("no verification key to validate")
	}

	for i, vk := range o {
		if err := vk.Valid(); err != nil {
			return fmt.Errorf("invalid verification key at index %d: %w", i, err)
		}
	}
	return nil
}
