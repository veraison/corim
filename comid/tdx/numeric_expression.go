// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import "fmt"

type NumericExpression struct {
	_               struct{}    `cbor:",toarray"`
	NumericOperator Operator    `json:"numeric-operator"`
	NumericType     NumericType `json:"numeric-type"`
}

type TaggedNumericExpression NumericExpression

func (o NumericExpression) Valid() error {
	switch o.NumericOperator {
	case GT, GE, LT, LE:
	default:
		return fmt.Errorf("invalid Operator %d", o.NumericOperator)
	}
	return nil
}

// NewTaggedNumericExpression creates a new TaggedNumericExpression from the
// supplied operator and an interface value. Allowed values are int, uint and float
func NewTaggedNumericExpression(op uint, val any) (*TaggedNumericExpression, error) {
	var tnum TaggedNumericExpression
	switch op {
	case GT, GE, LT, LE:
		numType, err := NewNumericType(val)
		if err != nil {
			return nil, fmt.Errorf("invalid NumericType type: %w", err)
		}
		numeric := NumericExpression{NumericOperator: Operator(op), NumericType: *numType}
		tnum = TaggedNumericExpression(numeric)
		return &tnum, nil
	default:
		return nil, fmt.Errorf("invalid numeric operator %d", op)
	}
}
