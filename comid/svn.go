// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
)

type TaggedSVN uint64
type TaggedMinSVN uint64

type SVN struct {
	val interface{}
}

func (o *SVN) SetSVN(val uint64) *SVN {
	if o != nil {
		o.val = TaggedSVN(val)
	}
	return o
}

func (o *SVN) SetMinSVN(val uint64) *SVN {
	if o != nil {
		o.val = TaggedMinSVN(val)
	}
	return o
}

func (o SVN) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

func (o *SVN) UnmarshalCBOR(data []byte) error {
	var svn TaggedSVN

	if dm.Unmarshal(data, &svn) == nil {
		o.val = svn
		return nil
	}

	var minsvn TaggedMinSVN

	if dm.Unmarshal(data, &minsvn) == nil {
		o.val = svn
		return nil
	}

	return fmt.Errorf("unknown SVN (CBOR: %x)", data)
}

type svnJSONRepr tnv

// Supported formats:
// { "type": "exact-value", "value": 123 } -> SVN
// { "type": "min-value", "value": 123 } -> MinSVN
func (o *SVN) UnmarshalJSON(data []byte) error {
	var s svnJSONRepr

	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("SVN decoding failure: %w", err)
	}

	var x uint64
	if err := json.Unmarshal(s.Value, &x); err != nil {
		return fmt.Errorf(
			"cannot unmarshal svn or min-svn: %w",
			err,
		)
	}

	switch s.Type {
	case "exact-value":
		o.val = TaggedSVN(x)
	case "min-value":
		o.val = TaggedMinSVN(x)
	default:
		return fmt.Errorf("unknown comparison operator %s", s.Type)
	}

	return nil
}

func (o SVN) MarshalJSON() ([]byte, error) {
	var (
		v   svnJSONRepr
		b   []byte
		err error
	)

	b, err = json.Marshal(o.val)
	if err != nil {
		return nil, err
	}
	switch t := o.val.(type) {
	case TaggedSVN:
		v = svnJSONRepr{Type: "exact-value", Value: b}
	case TaggedMinSVN:
		v = svnJSONRepr{Type: "min-value", Value: b}
	default:
		return nil, fmt.Errorf("unknown SVN type: %T", t)
	}

	return json.Marshal(v)
}
