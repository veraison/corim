// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/base64"
	"fmt"

	"github.com/veraison/eat"
)

const UEIDType = "ueid"

// UEID is an Unique Entity Identifier
type UEID eat.UEID

func (o UEID) Empty() bool {
	return len(o) == 0
}

// Valid checks that the target UEID is in one of the defined formats: IMEI, EUI or RAND
func (o UEID) Valid() error {
	if err := eat.UEID(o).Validate(); err != nil {
		return fmt.Errorf("UEID validation failed: %w", err)
	}
	return nil
}

func (o UEID) String() string {
	return base64.StdEncoding.EncodeToString(o)
}

// TaggedUEID is an alias to allow automatic tagging of an UEID type
type TaggedUEID UEID

func NewTaggedUEID(val any) (*TaggedUEID, error) {
	var ret TaggedUEID

	if val == nil {
		return &ret, nil
	}

	switch t := val.(type) {
	case string:
		b, err := base64.StdEncoding.DecodeString(t)
		if err != nil {
			return nil, fmt.Errorf("bad UEID: %w", err)
		}

		ret = TaggedUEID(b)
	case []byte:
		ret = TaggedUEID(t)
	case TaggedUEID:
		ret = append(ret, t...)
	case *TaggedUEID:
		ret = append(ret, *t...)
	case UEID:
		ret = append(ret, t...)
	case *UEID:
		ret = append(ret, *t...)
	case eat.UEID:
		ret = append(ret, t...)
	case *eat.UEID:
		ret = append(ret, *t...)
	default:
		return nil, fmt.Errorf("unexpeted type for UEID: %T", t)
	}

	if err := ret.Valid(); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (o TaggedUEID) Valid() error {
	return UEID(o).Valid()
}

func (o TaggedUEID) String() string {
	return UEID(o).String()
}

func (o TaggedUEID) Type() string {
	return "ueid"
}

func (o TaggedUEID) Bytes() []byte {
	return []byte(o)
}
