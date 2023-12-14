// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const UUIDType = "uuid"

// UUID represents an Universally Unique Identifier (UUID, see RFC4122)
type UUID uuid.UUID

// ParseUUID parses the supplied string into a UUID
func ParseUUID(s string) (UUID, error) {
	v, err := uuid.Parse(s)

	return UUID(v), err
}

// String returns a string representation of the binary UUID
func (o UUID) String() string {
	return uuid.UUID(o).String()
}

func (o UUID) Empty() bool {
	return o == (UUID{})
}

// Valid checks that the target UUID is formatted as per RFC4122
func (o UUID) Valid() error {
	if variant := uuid.UUID(o).Variant(); variant != uuid.RFC4122 {
		return fmt.Errorf("expecting RFC4122 UUID, got %s instead", variant)
	}
	return nil
}

// UnmarshalJSON deserializes the supplied string into the UUID target
// The UUID string in expected to be in canonical 8-4-4-4-12 format
func (o *UUID) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	u, err := ParseUUID(s)
	if err != nil {
		return fmt.Errorf("bad UUID: %w", err)
	}

	*o = u

	return nil
}

// MarshalJSON serialize the target UUID to a JSON string in canonical
// 8-4-4-4-12 format
func (o UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

// TaggedUUID is an alias to allow automatic tagging of a UUID type
type TaggedUUID UUID

func NewTaggedUUID(val any) (*TaggedUUID, error) {
	var ret TaggedUUID

	if val == nil {
		return &ret, nil
	}

	switch t := val.(type) {
	case string:
		u, err := ParseUUID(t)
		if err != nil {
			return nil, fmt.Errorf("bad UUID: %w", err)
		}
		ret = TaggedUUID(u)
	case []byte:
		if len(t) != 16 {
			return nil, fmt.Errorf(
				"unexpected size for UUID: expected 16 bytes, found %d",
				len(t),
			)
		}

		copy(ret[:], t)
	case TaggedUUID:
		copy(ret[:], t[:])
	case *TaggedUUID:
		copy(ret[:], (*t)[:])
	case UUID:
		copy(ret[:], t[:])
	case *UUID:
		copy(ret[:], (*t)[:])
	case uuid.UUID:
		copy(ret[:], t[:])
	case *uuid.UUID:
		copy(ret[:], (*t)[:])
	default:
		return nil, fmt.Errorf("unexpected type for UUID: %T", t)
	}

	if err := ret.Valid(); err != nil {
		return nil, err
	}

	return &ret, nil
}

// String returns a string representation of the binary UUID
func (o TaggedUUID) String() string {
	return UUID(o).String()
}

func (o TaggedUUID) Valid() error {
	return UUID(o).Valid()
}

// Type returns a string containing type name. This is part of the
// ITypeChoiceValue implementation.
func (o TaggedUUID) Type() string {
	return UUIDType
}

// Bytes returns a []byte containing the raw UUID bytes
func (o TaggedUUID) Bytes() []byte {
	return o[:]
}

func (o TaggedUUID) MarshalJSON() ([]byte, error) {
	temp := o.String()
	return json.Marshal(temp)
}

func (o *TaggedUUID) UnmarshalJSON(data []byte) error {
	var temp string
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	u, err := ParseUUID(temp)
	if err != nil {
		return fmt.Errorf("bad UUID: %w", err)
	}

	*o = TaggedUUID(u)

	return nil
}
