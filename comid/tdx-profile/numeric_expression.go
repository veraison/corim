// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

// Define a structure for NumericExpression
type NumericExpression struct {
	Operator NumericOperator
	Type     NumericType
}

type NumericExpressions []NumericExpression
