// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

// nolint:dupl
package tdx

import (
	"encoding/json"
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/encoding"
)

type Digests = comid.Digests

// TeeDigest, holds the digests. Allowed values are an array of digests OR
// a digest expression
type TeeDigest struct {
	val interface{}
}

// NewTeeDigest create a new TeeDigest from the
// supplied Digests and returns a pointer to
// the TeeDigest.
func NewTeeDigest(val Digests) (*TeeDigest, error) {
	if len(val) == 0 {
		return nil, fmt.Errorf("zero len value for TeeDigest")
	}
	return &TeeDigest{val: val}, nil
}

// NewTeeDigestExpr create a new TeeDigest (that holds Digest Expression) from the
// supplied operator and Digests and returns a pointer to
// the TeeDigest.
func NewTeeDigestExpr(op uint, val Digests) (*TeeDigest, error) {
	if len(val) == 0 {
		return nil, fmt.Errorf("zero len value for TeeDigestExpr")
	}
	switch op {
	case MEM, NMEM: // Allowed operators are either Members or Non Members
		set, err := NewTaggedSetDigestExpression(op, val)
		if err != nil {
			return nil, fmt.Errorf("zero len value for TeeDigestExpr")
		}
		return &TeeDigest{val: *set}, nil
	default:
		return nil, fmt.Errorf("invalid operator : %d", op)
	}
}

// AddTeeDigest add supplied Digests to existing TeeDigest
func (o *TeeDigest) AddTeeDigest(op uint, val Digests) (*TeeDigest, error) {
	if len(val) == 0 {
		return nil, fmt.Errorf("zero len value for TeeDigestExpr")
	}
	switch t := o.val.(type) {
	case Digests:
		t = append(t, val...)
		o.val = t
	case TaggedSetDigestExpression:
		if t.SetOperator != Operator(op) {
			return nil, fmt.Errorf("operator mis-match TeeDigest Op: %d, Input Op: %d", t.SetOperator, op)
		}
		for _, value := range val {
			t.SetDigest = append(t.SetDigest, value)
		}
		o.val = t
	}
	return o, nil
}

// Valid checks for validity of TeeDigest and
// returns an error, if invalid
func (o TeeDigest) Valid() error {
	if o.val == nil {
		return fmt.Errorf("TeeDigest not set")
	}
	switch t := o.val.(type) {
	case Digests:
		if len(t) == 0 {
			return fmt.Errorf("TeeDigest not set")
		}
	case TaggedSetDigestExpression:
		if t.SetOperator != MEM && t.SetOperator != NMEM {
			return fmt.Errorf("invalid operator in a TeeDigest: %d", t.SetOperator)
		}
		if len(t.SetDigest) == 0 {
			return fmt.Errorf("TeeDigest not set")
		}
	default:
		return fmt.Errorf("invalid type: %T", t)
	}
	return nil
}

// IsDigests returns true if TeeDigest is TeeDigest
func (o TeeDigest) IsDigests() bool {
	return isType[Digests](o.val)
}

// IsDigestExpr returns true if TeeDigest is DigestExpr
func (o TeeDigest) IsDigestExpr() bool {
	return isType[TaggedSetDigestExpression](o.val)
}

// GetDigestExpr returns a Digest Expression
func (o TeeDigest) GetDigestExpr() (*TaggedSetDigestExpression, error) {
	if err := o.Valid(); err != nil {
		return nil, fmt.Errorf("invalid TEEDigest: %w", err)
	}

	switch t := o.val.(type) {
	case TaggedSetDigestExpression:
		return &t, nil
	default:
		return nil, fmt.Errorf("invalid type: %T", t)
	}
}

func (o TeeDigest) GetDigest() (Digests, error) {
	switch t := o.val.(type) {
	case Digests:
		return t, nil
	default:
		return nil, fmt.Errorf("invalid type: %T", t)
	}
}

func (o TeeDigest) MarshalJSON() ([]byte, error) {
	var (
		v   encoding.TypeAndValue
		b   []byte
		err error
	)
	switch t := o.val.(type) {
	case Digests:
		b, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		v = encoding.TypeAndValue{Type: DigestType, Value: b}
	case TaggedSetDigestExpression:
		b, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		v = encoding.TypeAndValue{Type: DigestExprType, Value: b}
	default:
		return nil, fmt.Errorf("unknown type %T for TeeDigest", t)
	}
	return json.Marshal(v)
}

// UnmarshalJSON UnMarshals supplied JSON bytes to TeeDigest
func (o *TeeDigest) UnmarshalJSON(data []byte) error {
	var v encoding.TypeAndValue

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Type {
	case DigestType:
		var x Digests
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal TeeDigest of type string: %w", err)
		}
		o.val = x
	case DigestExprType:
		var x SetDigestExpression
		if err := json.Unmarshal(v.Value, &x); err != nil {
			return fmt.Errorf(
				"cannot unmarshal TeeDigest of type set expression: %w", err)
		}
		o.val = TaggedSetDigestExpression(x)
	}
	return nil
}

// MarshalCBOR Marshals TeeDigest to CBOR
func (o TeeDigest) MarshalCBOR() ([]byte, error) {

	return em.Marshal(o.val)
}

// UnmarshalCBOR UnMarshals supplied CBOR bytes to TeeDigest
func (o *TeeDigest) UnmarshalCBOR(data []byte) error {
	var x comid.Digests
	err := dm.Unmarshal(data, &x)
	if err == nil {
		o.val = x
		return nil
	} else {
		var y TaggedSetDigestExpression
		err = dm.Unmarshal(data, &y)
		if err == nil {
			o.val = y
			return nil
		} else {
			return err
		}
	}
}
