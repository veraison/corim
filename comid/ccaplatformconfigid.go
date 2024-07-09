// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"errors"
	"fmt"
	"unicode/utf8"
)

var CCAPlatformConfigIDType = "cca.platform-config-id"

type CCAPlatformConfigID string

func (o CCAPlatformConfigID) Empty() bool {
	return o == ""
}

func (o *CCAPlatformConfigID) Set(v string) error {
	if v == "" {
		return fmt.Errorf("empty input string")
	}
	*o = CCAPlatformConfigID(v)
	return nil
}

func (o CCAPlatformConfigID) Get() (CCAPlatformConfigID, error) {
	if o == "" {
		return "", fmt.Errorf("empty CCA platform config ID")
	}
	return o, nil
}

type TaggedCCAPlatformConfigID CCAPlatformConfigID

func NewTaggedCCAPlatformConfigID(val any) (*TaggedCCAPlatformConfigID, error) {
	var ret TaggedCCAPlatformConfigID

	if val == nil {
		return &ret, nil
	}

	switch t := val.(type) {
	case TaggedCCAPlatformConfigID:
		ret = t
	case *TaggedCCAPlatformConfigID:
		ret = *t
	case CCAPlatformConfigID:
		ret = TaggedCCAPlatformConfigID(t)
	case *CCAPlatformConfigID:
		ret = TaggedCCAPlatformConfigID(*t)
	case string:
		ret = TaggedCCAPlatformConfigID(t)
	case []byte:
		if !utf8.Valid(t) {
			return nil, errors.New("bytes do not form a valid UTF-8 string")
		}
		ret = TaggedCCAPlatformConfigID(t)
	default:
		return nil, fmt.Errorf("unexpected type for CCA platform-config-id: %T", t)
	}

	return &ret, nil
}

func (o TaggedCCAPlatformConfigID) Valid() error {
	if o == "" {
		return errors.New("empty value")
	}

	return nil
}

func (o TaggedCCAPlatformConfigID) String() string {
	return string(o)
}

func (o TaggedCCAPlatformConfigID) Type() string {
	return CCAPlatformConfigIDType
}

func (o TaggedCCAPlatformConfigID) IsZero() bool {
	return len(o) == 0
}

func (o *TaggedCCAPlatformConfigID) UnmarshalJSON(data []byte) error {
	var temp string

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*o = TaggedCCAPlatformConfigID(temp)

	return nil
}
