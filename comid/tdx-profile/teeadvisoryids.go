// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

// nolint:dupl
package tdx

import "fmt"

type TeeAdvisoryIDs setType

// NewTeeAvisoryIDs create a new TeeAvisoryIDs from the
// supplied interface array and returns a pointer to
// the AdvisoryIDs. Only
// Advisory IDs of string type are supported
func NewTeeAvisoryIDs(val []any) (*TeeAdvisoryIDs, error) {
	var adv TeeAdvisoryIDs
	if len(val) == 0 {
		return nil, fmt.Errorf("zero len TeeAdvisoryIDs")
	}

	for i, v := range val {
		switch t := v.(type) {
		case string:
			adv = append(adv, t)
		default:
			return nil, fmt.Errorf("invalid type: %T for AdvisoryIDs at index: %d", t, i)
		}
	}
	return &adv, nil
}

// AddTeeAdvisoryIDs add supplied AvisoryIDs to existing AdvisoryIDs
func (o *TeeAdvisoryIDs) AddTeeAdvisoryIDs(val []any) error {
	for i, v := range val {
		switch t := v.(type) {
		case string:
			*o = append(*o, t)
		default:
			return fmt.Errorf("invalid type: %T for AdvisoryIDs at index: %d", t, i)
		}
	}
	return nil
}

// Valid checks for validity of TeeAdvisoryIDs and
// returns an error, if invalid
func (o TeeAdvisoryIDs) Valid() error {
	if len(o) == 0 {
		return fmt.Errorf("empty AdvisoryIDs")

	}
	for i, v := range o {
		switch t := v.(type) {
		case string:
			continue
		default:
			return fmt.Errorf("invalid type: %T for AdvisoryIDs at index: %d", t, i)
		}
	}
	return nil
}
