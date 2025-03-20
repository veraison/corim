// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

type NumericOperator uint

type taggedNumericOperator NumericOperator

const (
	gt = iota + 1
	ge
	lt
	le
)

var (
	stringToNumericOperator = map[string]NumericOperator{
		"gt": gt,
		"ge": ge,
		"lt": lt,
		"le": le,
	}
	NumericOperatorToString = map[NumericOperator]string{
		gt: "gt",
		ge: "ge",
		lt: "lt",
		le: "le",
	}
)
