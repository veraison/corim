// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

// nolint:dupl
package tdx

import (
	"encoding/json"
	"fmt"

	"github.com/veraison/corim/encoding"
)

// TeeAdvisoryIDs, holds the Advisory IDs. Allowed values are an array of strings OR
// An array of strings expressed in an expression
type TeeAdvisoryIDs struct {
	val interface{}
}

// NewTeeAvisoryIDsExpr create a new TeeAvisoryIDs from the
// supplied operator and an array of strings
func NewTeeAdvisoryIDsExpr(op uint, val []string) (*TeeAdvisoryIDs, error) {
	if len(val) == 0 {
		return nil, fmt.Errorf("zero len value for TeeAdvisoryID")
	}
	switch op {
	case MEM, NMEM:
		set, err := NewTaggedSetStringExpression(op, val)
		if err != nil {
			return nil, fmt.Errorf("unable to get NewTeeAdvisoryIDExpr %w", err)
		}
		return &TeeAdvisoryIDs{val: *set}, nil
	default:
		return nil, fmt.Errorf("invalid operator : %d", op)
	}
}

// NewTeeAdvisoryIDsString create a new TeeAvisoryIDs from the
// supplied array of strings
func NewTeeAdvisoryIDsString(val []string) (*TeeAdvisoryIDs, error) {
	if len(val) == 0 {
		return nil, fmt.Errorf("zero len value for TeeAdvisoryID")
	}
	return &TeeAdvisoryIDs{val: val}, nil
}

// AddTeeAdvisoryIDs add supplied AvisoryIDs to existing AdvisoryIDs
func (o *TeeAdvisoryIDs) AddTeeAdvisoryIDs(op uint, val []string) (*TeeAdvisoryIDs, error) {
	if len(val) == 0 {
		return nil, fmt.Errorf("zero len value for TeeAdvisoryID")
	}
	switch t := o.val.(type) {
	case []string:
		t = append(t, val...)
		o.val = t
	case TaggedSetStringExpression:
		if t.SetOperator != Operator(op) {
			return nil, fmt.Errorf("operator mis-match Advisory Op: %d, Input Op: %d", t.SetOperator, op)
		}
		t.SetString = append(t.SetString, val...)
		o.val = t
	default:
		return nil, fmt.Errorf("unknown type %T for TeeAdvisoryIDs", t)
	}
	return o, nil
}

// Valid checks for validity of TeeAdvisoryIDs and
// returns an error, if invalid
func (o TeeAdvisoryIDs) Valid() error {
	if o.val == nil {
		return fmt.Errorf("TeeAdvisoryID not set")
	}
	switch t := o.val.(type) {
	case []string:
		if len(t) == 0 {
			return fmt.Errorf("TeeAdvisoryID not set")
		}
	case TaggedSetStringExpression:
		if t.SetOperator != MEM && t.SetOperator != NMEM {
			return fmt.Errorf("invalid operator in a TeeAdvisoryID: %d", t.SetOperator)
		}
		if len(t.SetString) == 0 {
			return fmt.Errorf("TeeAdvisoryID not set")
		}
	default:
		return fmt.Errorf("unknown type %T for TeeAdvisoryIDs", t)

	}
	return nil
}

// IsString returns true if TeeAdvisoryIDs is slice of strings
func (o TeeAdvisoryIDs) IsString() bool {
	return isType[[]string](o.val)
}

// IsStringExpr returns true if IsStringExpr is SetStringExpr
func (o TeeAdvisoryIDs) IsStringExpr() bool {
	return isType[TaggedSetStringExpression](o.val)
}

// GetString returns a slice of TeeAdvisoryIDs as strings
func (o TeeAdvisoryIDs) GetString() ([]string, error) {
	switch t := o.val.(type) {
	case []string:
		return t, nil

	default:
		return nil, fmt.Errorf("unknown type %T for TeeAdvisoryIDs", t)
	}
}

// GetStringExpression returns TaggedSetStringExpression from TeeAdvisoryIDs
func (o TeeAdvisoryIDs) GetStringExpression() (TaggedSetStringExpression, error) {
	switch t := o.val.(type) {
	case TaggedSetStringExpression:
		return t, nil

	default:
		return TaggedSetStringExpression{}, fmt.Errorf("unknown type %T for TeeAdvisoryIDs", t)
	}
}

func (o TeeAdvisoryIDs) MarshalJSON() ([]byte, error) {
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
		return nil, fmt.Errorf("unknown type %T for TeeAdvisoryIDs", t)
	}
	return json.Marshal(v)
}

// UnmarshalJSON UnMarshals supplied JSON bytes to TeeAdvisoryIDs
func (o *TeeAdvisoryIDs) UnmarshalJSON(data []byte) error {
	var v encoding.TypeAndValue

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Type {
	case StringType:
		var x []string
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal TeeAdvisoryIDs of type string: %w", err)
		}
		o.val = x
	case StringExprType:
		var x SetStringExpression
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal TeeAdvisoryIDs of type set expression: %w", err)
		}
		o.val = TaggedSetStringExpression(x)
	default:
		return fmt.Errorf("unknown type %T for TeeAdvisoryIDs", v.Type)
	}
	return nil
}

// MarshalCBOR Marshals TeeAdvisoryIDs to CBOR
func (o TeeAdvisoryIDs) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

// UnmarshalCBOR UnMarshals supplied CBOR bytes to TeeAdvisoryIDs
func (o *TeeAdvisoryIDs) UnmarshalCBOR(data []byte) error {
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
