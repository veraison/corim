// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"fmt"

	"github.com/veraison/corim/comid"
)

type SetDigest comid.Digests

type SetString []string

type SetStringExpression struct {
	_           struct{}  `cbor:",toarray"`
	SetOperator Operator  `json:"set-operator"`
	SetString   SetString `json:"set-string"`
}

type SetDigestExpression struct {
	_           struct{}  `cbor:",toarray"`
	SetOperator Operator  `json:"set-operator"`
	SetDigest   SetDigest `json:"set-digest"`
}

type TaggedSetStringExpression SetStringExpression

type TaggedSetDigestExpression SetDigestExpression

// NewTaggedSetStringExpression creates a TaggedSetStringExpression from
// the supplied operator and an array of strings
func NewTaggedSetStringExpression(op uint, val []string) (*TaggedSetStringExpression, error) {
	switch op {
	case MEM, NMEM:
		if len(val) == 0 {
			return nil, fmt.Errorf("zero len string array supplied")
		}
		set := SetStringExpression{SetOperator: Operator(op), SetString: SetString(val)}
		tset := TaggedSetStringExpression(set)
		return &tset, nil
	default:
		return nil, fmt.Errorf("invalid set operator %d", op)
	}
}

// NewTaggedSetDigestExpression creates a TaggedSetStringExpression from
// the supplied operator and Digests
func NewTaggedSetDigestExpression(op uint, val comid.Digests) (*TaggedSetDigestExpression, error) {
	switch op {
	case MEM, NMEM:
		if len(val) == 0 {
			return nil, fmt.Errorf("no Digests supplied")
		}
		set := SetDigestExpression{SetOperator: Operator(op), SetDigest: SetDigest(val)}
		tset := TaggedSetDigestExpression(set)
		return &tset, nil
	default:
		return nil, fmt.Errorf("invalid set operator %d", op)
	}
}
