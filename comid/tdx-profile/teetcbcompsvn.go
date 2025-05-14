// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"fmt"
)

// MaxSVNCount is the maximum SVN count in TeeTcbCompSvn
const MaxSVNCount = 16

type TeeTcbCompSvn [MaxSVNCount]TeeSVN

// NewTeeTcbCompSvnExpression creates a new TeeTcbCompSvn with an array of SVN as Numeric Expression
// from the supplied array of integers
func NewTeeTcbCompSvnExpression(val []uint) (*TeeTcbCompSvn, error) {
	if len(val) > MaxSVNCount {
		return nil, fmt.Errorf("invalid length %d for TeeTcbCompSVN", len(val))
	} else if len(val) == 0 {
		return nil, fmt.Errorf("no value supplied for TeeTcbCompSVN")
	}

	TeeTcbCompSVN := make([]TeeSVN, MaxSVNCount)
	for i, value := range val {
		svn, err := NewSvnExpression(value)
		if err != nil {
			return nil, fmt.Errorf("unable to get New SVN Numeric: %w", err)
		}
		TeeTcbCompSVN[i] = *svn
	}
	return (*TeeTcbCompSvn)(TeeTcbCompSVN), nil
}

// NewTeeTcbSVNUint creates a new TeeTcbCompSvn with an array of SVN of type unsigned integers
func NewTeeTcbCompSvnUint(val []uint) (*TeeTcbCompSvn, error) {
	if len(val) > MaxSVNCount {
		return nil, fmt.Errorf("invalid length %d for TeeTcbCompSVN", len(val))
	} else if len(val) == 0 {
		return nil, fmt.Errorf("no value supplied for TeeTcbCompSVN")
	}

	TeeTcbCompSVN := make([]TeeSVN, MaxSVNCount)
	for i, value := range val {
		TeeTcbCompSVN[i].val = value
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

	for i, svn := range o {
		if err := svn.Valid(); err != nil {
			return fmt.Errorf("invalid TeeSvn at index %d, %w", i, err)
		}
	}
	return nil
}
