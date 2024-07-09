// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "fmt"

type ReferenceValue struct {
	_            struct{}     `cbor:",toarray"`
	Environment  Environment  `json:"environment"`
	Measurements Measurements `json:"measurements"`
}

func (o ReferenceValue) Valid() error {
	if err := o.Environment.Valid(); err != nil {
		return fmt.Errorf("environment validation failed: %w", err)
	}

	if err := o.Measurements.Valid(); err != nil {
		return fmt.Errorf("measurements validation failed: %w", err)
	}

	return nil
}
