// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

// nolint:dupl
package tdx

import "fmt"

type TeeTcbStatus setType

func NewTeeTcbStatus(val []any) (*TeeTcbStatus, error) {
	var ts TeeTcbStatus
	if len(val) == 0 {
		return nil, fmt.Errorf("nil value argument")
	}

	for i, v := range val {
		switch t := v.(type) {
		case string:
			ts = append(ts, t)
		default:
			return nil, fmt.Errorf("invalid type: %T for tcb status at index: %d", t, i)
		}
	}
	return &ts, nil
}

func (o *TeeTcbStatus) AddTeeTcbStatus(val []any) error {
	for i, v := range val {
		switch t := v.(type) {
		case string:
			*o = append(*o, t)
		default:
			return fmt.Errorf("invalid type: %T for tcb status at index: %d", t, i)
		}
	}
	return nil
}

func (o TeeTcbStatus) Valid() error {
	if len(o) == 0 {
		return fmt.Errorf("empty tcb status")
	}

	for i, v := range o {
		switch t := v.(type) {
		case string:
			continue
		default:
			return fmt.Errorf("invalid type: %T for tcb status at index: %d", t, i)
		}
	}
	return nil
}
