// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"fmt"
)

type TeePCEID string

func NewTeePCEID(val string) (*TeePCEID, error) {
	var pceID TeePCEID
	if val == "" {
		return nil, fmt.Errorf("null string for TeePCEID")
	}
	pceID = TeePCEID(val)
	return &pceID, nil
}

func (o TeePCEID) Valid() error {
	if o == "" {
		return fmt.Errorf("nil TeePCEID")
	}
	return nil
}
