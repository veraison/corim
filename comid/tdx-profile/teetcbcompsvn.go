// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import "fmt"

// MaxSVNCount is the maximum SVN count in TeeTcbCompSvn
const MaxSVNCount = 16

type TeeTcbCompSvn [MaxSVNCount]TeeSVN

func NewTeeTcbCompSVN(val []uint) (*TeeTcbCompSvn, error) {
	if len(val) > MaxSVNCount {
		return nil, fmt.Errorf("invalid length %d for TeeTcbCompSVN", len(val))
	} else if len(val) == 0 {
		return nil, fmt.Errorf("no value supplied for TeeTcbCompSVN")
	}

	TeeTcbCompSVN := make([]TeeSVN, MaxSVNCount)
	for i, value := range val {
		TeeTcbCompSVN[i] = TeeSVN(value)
	}
	return (*TeeTcbCompSvn)(TeeTcbCompSVN), nil
}

// nolint:gocritic
func (o TeeTcbCompSvn) Valid() error {
	if len(o) == 0 {
		return fmt.Errorf("empty TeeTcbCompSVN")
	}
	if len(o) > MaxSVNCount {
		return fmt.Errorf("invalid length: %d for TeeTcbCompSVN", len(o))
	}
	return nil
}
