// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

// nolint:dupl
package tdx

import "fmt"

type TeeAttributes maskType

func NewTeeAttributes(val []byte) (*TeeAttributes, error) {
	if val == nil {
		return nil, fmt.Errorf("nil TeeAttributes")
	}
	teeAttributes := TeeAttributes(val)
	return &teeAttributes, nil
}

func (o TeeAttributes) Valid() error {
	if o == nil {
		return fmt.Errorf("nil TeeAttributes")
	}
	if len(o) == 0 {
		return fmt.Errorf("zero len TeeAttributes")
	}
	return nil
}
