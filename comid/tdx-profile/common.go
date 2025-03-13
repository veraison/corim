// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

type numOperator uint

type taggednumOperator numOperator

const (
	gt = iota + 1
	ge
	lt
	le
)

func isType[T any](v any) bool {
	_, ok := v.(T)
	return ok
}

var (
	stringTonumOperator = map[string]numOperator{
		"gt": gt,
		"ge": ge,
		"lt": lt,
		"le": le,
	}
	numOperatorToString = map[numOperator]string{
		gt: "gt",
		ge: "ge",
		lt: "lt",
		le: "le",
	}
)
