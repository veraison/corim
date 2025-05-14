// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"encoding/json"
	"fmt"
)

type Operator uint

const (
	EQ   = iota // EQual
	GT          // greater than
	GE          // greater than or equal to
	LT          // less than
	LE          // less than or equal to
	NOP         // undefined
	MEM         // member
	NMEM        // no member
	SUB         // sub-set
	SUP         // super set
	DIS         // dis-joint
)

var (
	StringToNumericOperator = map[string]Operator{
		"equal":            EQ,
		"greater_than":     GT,
		"greater_or_equal": GE,
		"less_than":        LT,
		"less_or equal":    LE,
		"nop":              NOP,
		"member":           MEM,
		"non_member":       NMEM,
		"subset":           SUB,
		"superset":         SUP,
		"disjoint":         DIS,
	}
	NumericOperatorToString = map[Operator]string{
		EQ:   "equal",
		GT:   "greater_than",
		GE:   "greater_or_equal",
		LT:   "less_than",
		LE:   "less_or equal",
		NOP:  "nop",
		MEM:  "member",
		NMEM: "non_member",
		SUB:  "subset",
		SUP:  "superset",
		DIS:  "disjoint",
	}
)

func (o Operator) MarshalJSON() ([]byte, error) {
	return json.Marshal(NumericOperatorToString[o])
}

func (o *Operator) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return fmt.Errorf("unable to unmarshal operator: %w", err)
	}
	*o = StringToNumericOperator[str]
	return nil
}
