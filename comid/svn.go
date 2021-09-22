// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
)

type TaggedSVN int64
type TaggedMinSVN int64

type SVN struct {
	val interface{}
}

func (o *SVN) SetSVN(val int64) *SVN {
	if o != nil {
		o.val = TaggedSVN(val)
	}
	return o
}

func (o *SVN) SetMinSVN(val int64) *SVN {
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

// Supported formats:
// { "cmp": "==", "value": 123 } -> SVN
// { "cmp": ">=", "value": 123 } -> MinSVN
func (o *SVN) UnmarshalJSON(data []byte) error {
	s := struct {
		Cmp string `json:"cmp"`
		Val int64  `json:"value"`
	}{}

	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("SVN decoding failure: %w", err)
	}

	switch s.Cmp {
	case "==":
		o.SetSVN(s.Val)
	case ">=":
		o.SetMinSVN(s.Val)
	default:
		return fmt.Errorf("unknown comparison operator %s", s.Cmp)
	}

	return nil
}
