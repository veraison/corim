// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/base64"
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

func NewBytesFromBase64(val any) (*TaggedBytes, error) {
	var ret TaggedBytes

	if val == nil {
		return &ret, nil
	}

	switch t := val.(type) {
	case string:
		b, err := base64.StdEncoding.DecodeString(t)
		if err != nil {
			return nil, err
		}
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
	return base64.StdEncoding.EncodeToString(o[:])
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
