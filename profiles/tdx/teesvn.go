// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"encoding/json"
	"fmt"

	"github.com/veraison/corim/encoding"
)

// TeeSVN holds a TeeSVN, supported formats are uint and TaggedNumericExpression
type TeeSVN struct {
	val interface{}
}

// NewSvnExpression creates a new TeeSVN which contains an SVN of type Numeric Expression
func NewSvnExpression(val uint) (*TeeSVN, error) {
	tnum, err := NewTaggedNumericExpression(GE, val)
	if err != nil {
		return nil, err
	}
	return &TeeSVN{val: *tnum}, nil
}

// NewSvnUint creates a new TeeSVN which contains an SVN of type uint
func NewSvnUint(val uint) (*TeeSVN, error) {
	return &TeeSVN{val: val}, nil
}

func (o TeeSVN) Valid() error {
	switch t := o.val.(type) {
	case uint, uint64:
		return nil
	case TaggedNumericExpression:
		if t.NumericOperator != GE {
			return fmt.Errorf("unknown operator %s for Numeric TeeSVN", NumericOperatorToString[t.NumericOperator])
		}
		exp := t.NumericType
		switch k := exp.val.(type) {
		case uint, uint64:
			return nil
		default:
			return fmt.Errorf("unknown type %T for Numeric TeeSVN", k)
		}
	default:
		return fmt.Errorf("unknown type %T for TeeSVN", t)
	}
}

func (o TeeSVN) IsExpression() bool {
	return isType[TaggedNumericExpression](o.val)
}

func (o TeeSVN) IsUint() bool {
	return isType[uint](o.val) || isType[uint64](o.val)
}

func (o TeeSVN) GetUint() (uint, error) {
	switch t := o.val.(type) {
	case uint:
		return t, nil
	case uint64:
		return uint(t), nil
	default:
		return 0, fmt.Errorf("unknown type %T for TeeSVN", t)
	}
}

func (o TeeSVN) GetNumericExpression() (TaggedNumericExpression, error) {
	switch t := o.val.(type) {
	case TaggedNumericExpression:
		return t, nil
	default:
		return TaggedNumericExpression{}, fmt.Errorf("unknown type %T for TeeSVN", t)
	}
}

// MarshalJSON Marshals TeeSVN to JSON bytes
func (o TeeSVN) MarshalJSON() ([]byte, error) {
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
		return nil, fmt.Errorf("unknown type %T for TeeSVN", t)
	}
	return json.Marshal(v)
}

// UnmarshalJSON UnMarshals supplied JSON bytes to TeeSVN
func (o *TeeSVN) UnmarshalJSON(data []byte) error {
	var v encoding.TypeAndValue

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Type {
	case UintType:
		var x uint
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal TeeInstanceID of type uint: %w", err)
		}
		o.val = x
	case NumericExprType:
		var x NumericExpression
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal TeeSVN of type numeric-expression: %w", err)
		}
		o.val = TaggedNumericExpression(x)
	}
	return o.Valid()
}

// MarshalCBOR Marshals TeeSVN to CBOR
func (o TeeSVN) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

// UnmarshalCBOR UnMarshals supplied CBOR bytes to TeeSVN
func (o *TeeSVN) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &o.val)
}
