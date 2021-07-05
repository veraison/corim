// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/eat"
	"github.com/veraison/swid"
)

type UnsignedCorim struct {
	ID            ID             `cbor:"0,keyasint" json:"corim-id"`
	Tags          []Tag          `cbor:"1,keyasint" json:"tags"`
	DependentRims *[]Locator     `cbor:"2,keyasint,omitempty" json:"dependent-rims,omitempty"`
	Profiles      *[]eat.Profile `cbor:"3,keyasint,omitempty" json:"profiles,omitempty"`
}

type ID struct {
	val interface{}
}

func (o ID) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

func (o *ID) UnmarshalCBOR(data []byte) error {
	var s string

	if dm.Unmarshal(data, &s) == nil {
		o.val = s
		return nil
	}

	var u comid.UUID

	if dm.Unmarshal(data, &u) == nil {
		o.val = u
		return nil
	}

	return fmt.Errorf("unknown corim-id type (CBOR: %x)", data)
}

type Tag []byte

type TaggedComid comid.Comid

type Locator struct {
	Href       comid.TaggedURI `cbor:"0,keyasint" json:"href"`
	Thumbprint swid.HashEntry  `cbor:"1,keyasint,omitempty" json:"thumbprint,omitempty"`
}

func NewUnsignedCorim() *UnsignedCorim {
	return &UnsignedCorim{}
}

func (o *UnsignedCorim) SetIDString(s string) *UnsignedCorim {
	if o != nil {
		if s == "" {
			return nil
		}
		o.ID.val = s
	}
	return o
}

func (o *UnsignedCorim) SetIDUUID(u comid.UUID) *UnsignedCorim {
	if o != nil {
		if u == (comid.UUID{}) {
			return nil
		}
		o.ID.val = u
	}
	return o
}

func (o UnsignedCorim) GetIDString() (string, error) {
	switch t := o.ID.val.(type) {
	case string:
		return t, nil
	default:
		return "", fmt.Errorf("corim-id type is: %T", t)
	}
}

func (o UnsignedCorim) GetIDUUID() (comid.UUID, error) {
	switch t := o.ID.val.(type) {
	case comid.UUID:
		return t, nil
	default:
		return comid.UUID{}, fmt.Errorf("corim-id type is: %T", t)
	}
}

func (o *UnsignedCorim) AddComid(c TaggedComid) *UnsignedCorim {
	comid, err := em.Marshal(c)
	if err != nil {
		return nil
	}

	o.Tags = append(o.Tags, comid)

	return o
}
