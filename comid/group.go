// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
)

// Group stores a group identity.  The supported format is UUID.
type Group struct {
	val interface{}
}

// NewGroup instantiates an empty group
func NewGroup() *Group {
	return &Group{}
}

// SetUUID sets the identity of the target group to the supplied UUID
func (o *Group) SetUUID(val UUID) *Group {
	if o != nil {
		o.val = TaggedUUID(val)
	}
	return o
}

// NewGroupUUID instantiates a new group with the supplied UUID identity
func NewGroupUUID(val UUID) *Group {
	return NewGroup().SetUUID(val)
}

// Valid checks for the validity of given group
func (o Group) Valid() error {
	if o.String() == "" {
		return fmt.Errorf("invalid group id")
	}
	return nil
}

// String returns a printable string of the Group value.  UUIDs use the
// canonical 8-4-4-4-12 format, UEIDs are hex encoded.
func (o Group) String() string {
	switch t := o.val.(type) {
	case TaggedUUID:
		return UUID(t).String()
	default:
		return ""
	}
}

// MarshalCBOR serializes the target group to CBOR
func (o Group) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

// UnmarshalCBOR deserializes the supplied CBOR into the target group
func (o *Group) UnmarshalCBOR(data []byte) error {
	var uuid TaggedUUID

	if dm.Unmarshal(data, &uuid) == nil {
		o.val = uuid
		return nil
	}

	return fmt.Errorf("unknown group type (CBOR: %x)", data)
}

// UnmarshalJSON deserializes the supplied JSON type/value object into the Group
// target.  The only supported format is UUID, e.g.:
//
//   {
//     "type": "uuid",
//     "value": "69E027B2-7157-4758-BCB4-D9F167FE49EA"
//   }
func (o *Group) UnmarshalJSON(data []byte) error {
	v := struct {
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	}{}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Type {
	case "uuid":
		var uuid UUID
		if err := jsonDecodeUUID(v.Value, &uuid); err != nil {
			return err
		}
		o.SetUUID(uuid)
	default:
		return fmt.Errorf("unknown type %s for group", v.Type)
	}

	return nil
}
