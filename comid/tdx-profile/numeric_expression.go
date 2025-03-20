// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

type NumericExpression []interface{}

func NewExpression(a NumericOperator, b NumericType) (*NumericExpression, error) {
	// check for validity of a
	// check for validity of b
	expression := NumericExpression{
		a, b,
	}
	return &expression, nil
}
