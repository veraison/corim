// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

// nolint:dupl
package tdx

import (
	"encoding/json"
	"fmt"

	"github.com/veraison/corim/encoding"
)

// TeeTcbStatus, holds the TCB Status. Allowed values are an array of strings OR
// An array of strings expressed in an expression
type TeeTcbStatus struct {
	val interface{}
}

// NewTcbStatusExpr creates a new TeeTcbStatus from the supplied operator and
// an array of strings
func NewTcbStatusExpr(op Operator, val []string) (*TeeTcbStatus, error) {
	if len(val) == 0 {
		return nil, fmt.Errorf("zero len value for TeeTcbStatus")
	}
	switch op {
	case MEM, NMEM:
		set := SetStringExpression{SetOperator: op, SetString: val}
		return &TeeTcbStatus{val: TaggedSetStringExpression(set)}, nil
	default:
		return nil, fmt.Errorf("invalid operator : %d", op)
	}

}

// NewTeeTcbStatusString creates a new TeeTcbStatus from the
// supplied array of strings
func NewTeeTcbStatusString(val []string) (*TeeTcbStatus, error) {
	if len(val) == 0 {
		return nil, fmt.Errorf("zero len value for TeeTcbStatus")
	}
	return &TeeTcbStatus{val: val}, nil
}

// AddTeeTcbStatus add supplied TeeTcbStatus to existing TeeTcbStatus
func (o *TeeTcbStatus) AddTeeTcbStatus(op Operator, val []string) (*TeeTcbStatus, error) {
	if len(val) == 0 {
		return nil, fmt.Errorf("zero len value for TeeTcbStatus")
	}
	switch t := o.val.(type) {
	case []string:
		t = append(t, val...)
		o.val = t
	case TaggedSetStringExpression:
		if t.SetOperator != op {
			return nil, fmt.Errorf("operator mis-match TeeTcbStatus Op: %d, Input Op: %d", t.SetOperator, op)
		}
		for _, value := range val {
			t.SetString = append(t.SetString, value)
		}
		o.val = t
	}
	return o, nil
}

// Valid checks for validity of TeeTcbStatus and
// returns an error, if invalid
func (o TeeTcbStatus) Valid() error {
	if o.val == nil {
		return fmt.Errorf("TeeTcbStatus not set")
	}
	switch t := o.val.(type) {
	case []string:
		if len(t) == 0 {
			return fmt.Errorf("TeeTcbStatus not set")
		}
	case TaggedSetStringExpression:
		if t.SetOperator != MEM && t.SetOperator != NMEM {
			return fmt.Errorf("invalid operator in a TeeTcbStatus: %d", t.SetOperator)
		}
		if len(t.SetString) == 0 {
			return fmt.Errorf("TeeTcbStatus not set")
		}
	default:
		return fmt.Errorf("unknown type %T for TeeTcbStatus", t)
	}
	return nil
}

// IsString returns true if TeeTcbStatus is slice of strings
func (o TeeTcbStatus) IsString() bool {
	return isType[[]string](o.val)
}

// IsStringExpr returns true if TeeTcbStatus is SetStringExpr
func (o TeeTcbStatus) IsStringExpr() bool {
	return isType[TaggedSetStringExpression](o.val)
}

// GetString returns a slice of TeeTcbStatus as strings
func (o TeeTcbStatus) GetString() ([]string, error) {
	switch t := o.val.(type) {
	case []string:
		return t, nil

	default:
		return nil, fmt.Errorf("unknown type %T for TeeAdvisoryIDs", t)
	}
}

// GetStringExpression returns TaggedSetStringExpression from TeeTcbStatus
func (o TeeTcbStatus) GetStringExpression() (TaggedSetStringExpression, error) {
	switch t := o.val.(type) {
	case TaggedSetStringExpression:
		return t, nil

	default:
		return TaggedSetStringExpression{}, fmt.Errorf("unknown type %T for TeeAdvisoryIDs", t)
	}
}

// MarshalJSON marshals the TeeTcbStatus to JSON
func (o TeeTcbStatus) MarshalJSON() ([]byte, error) {
	var (
		v   encoding.TypeAndValue
		b   []byte
		err error
	)
	switch t := o.val.(type) {
	case []string:
		b, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		v = encoding.TypeAndValue{Type: StringType, Value: b}
	case TaggedSetStringExpression:
		b, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		v = encoding.TypeAndValue{Type: StringExprType, Value: b}
	default:
		return nil, fmt.Errorf("unknown type %T for TeeTcbStatus", t)
	}
	return json.Marshal(v)
}

// UnmarshalJSON Unmarshals supplied JSON bytes to TeeTcbStatus
func (o *TeeTcbStatus) UnmarshalJSON(data []byte) error {
	var v encoding.TypeAndValue

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Type {
	case StringType:
		var x []string
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal TeeTcbStatus of type string: %w", err)
		}
		o.val = x
	case StringExprType:
		var x SetStringExpression
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal TeeTcbStatus of type set expression: %w", err)
		}
		o.val = TaggedSetStringExpression(x)
	default:
		return fmt.Errorf("unknown type %s for TeeTcbStatus", v.Type)
	}
	return nil
}

// MarshalCBOR marshals TeeTcbStatus to CBOR
func (o TeeTcbStatus) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

// UnmarshalCBOR Unmarshals supplied CBOR bytes to TeeTcbStatus
func (o *TeeTcbStatus) UnmarshalCBOR(data []byte) error {
	var x []string
	err := dm.Unmarshal(data, &x)
	if err == nil {
		o.val = x
	} else {
		var x TaggedSetStringExpression
		err = dm.Unmarshal(data, &x)
		if err == nil {
			o.val = x
		} else {
			return err
		}
	}
	return nil
}
