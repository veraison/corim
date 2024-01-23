// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
)

const BytesType = "bytes"

type TaggedBytes []byte

func NewBytes(val any) (*TaggedBytes, error) {
	var ret TaggedBytes

	if val == nil {
		return &ret, nil
	}

	switch t := val.(type) {
	case string:
		b := []byte(t)
		ret = TaggedBytes(b)
	case []byte:
		ret = TaggedBytes(t)
	case *[]byte:
		ret = TaggedBytes(*t)
	default:
		return nil, fmt.Errorf("unexpected type for bytes: %T", t)
	}
	return &ret, nil
}

func (o TaggedBytes) String() string {
	return string(o)
}

func (o TaggedBytes) Valid() error {
	return nil
}

func (o TaggedBytes) Type() string {
	return "bytes"
}

func (o TaggedBytes) Bytes() []byte {

	return o
}
