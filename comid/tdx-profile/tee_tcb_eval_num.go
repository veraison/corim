// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

// nolint:dupl
package tdx

import (
	"encoding/json"
	"fmt"

	"github.com/veraison/corim/encoding"
)

// TeeTcbEvalNumber holds a TeeTcbEvalNumber, supported formats are uint and TaggedNumericExpression
type TeeTcbEvalNumber struct {
	val interface{}
}

func NewTeeTcbEvalNumberNumeric(val uint) (*TeeTcbEvalNumber, error) {
	tnum, err := NewTaggedNumericExpression(GE, val)
	if err != nil {
		return nil, err
	}
	return &TeeTcbEvalNumber{val: *tnum}, nil
}

func NewTeeTcbEvalNumberUint(val uint) (*TeeTcbEvalNumber, error) {
	return &TeeTcbEvalNumber{val: val}, nil
}

func (o TeeTcbEvalNumber) Valid() error {
	switch t := o.val.(type) {
	case uint, uint64:
		return nil
	case TaggedNumericExpression:
		if t.NumericOperator != GE {
			return fmt.Errorf("unknown operator %d for TeeTcbEvalNumber", t.NumericOperator)
		}
		return nil
	default:
		return fmt.Errorf("unknown type %T for TeeTcbEvalNumber", t)
	}
}

func (o TeeTcbEvalNumber) IsExpression() bool {
	return isType[TaggedNumericExpression](o.val)
}

func (o TeeTcbEvalNumber) IsUint() bool {
	return isType[uint](o.val) || isType[uint64](o.val)
}

func (o TeeTcbEvalNumber) GetUint() (uint, error) {
	switch t := o.val.(type) {
	case uint:
		return t, nil
	case uint64:
		return uint(t), nil
	default:
		return 0, fmt.Errorf("unknown type %T for TeeTcbEvalNumber", t)
	}
}

func (o TeeTcbEvalNumber) GetNumericExpression() (TaggedNumericExpression, error) {
	switch t := o.val.(type) {
	case TaggedNumericExpression:
		return t, nil
	default:
		return TaggedNumericExpression{}, fmt.Errorf("unknown type %T for TeeTcbEvalNumber", t)
	}
}

func (o *TeeTcbEvalNumber) MarshalJSON() ([]byte, error) {
	var (
		v   encoding.TypeAndValue
		b   []byte
		err error
	)
	switch t := o.val.(type) {
	case uint, uint64:
		b, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		v = encoding.TypeAndValue{Type: UintType, Value: b}
	case TaggedNumericExpression:
		b, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		v = encoding.TypeAndValue{Type: NumericExprType, Value: b}
	default:
		return nil, fmt.Errorf("unknown type %T for TeeTcbEvalNumber", t)
	}
	return json.Marshal(v)
}

// UnmarshalJSON UnMarshals supplied JSON bytes to TeeTcbEvalNumber
func (o *TeeTcbEvalNumber) UnmarshalJSON(data []byte) error {
	var v encoding.TypeAndValue

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Type {
	case UintType:
		var x uint
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal TeeTcbEvalNumber of type uint: %w", err)
		}
		o.val = x
	case NumericExprType:
		var x NumericExpression
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal TeeTcbEvalNumber of type numeric: %w", err)
		}
		o.val = TaggedNumericExpression(x)
	}
	return nil
}

// MarshalCBOR Marshals TeeTcbEvalNumber to CBOR
func (o TeeTcbEvalNumber) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

// UnmarshalCBOR UnMarshals supplied CBOR bytes to TeeTcbEvalNumber
func (o *TeeTcbEvalNumber) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &o.val)
}
