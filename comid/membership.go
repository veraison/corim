// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"

	"github.com/veraison/corim/extensions"
)

// Membership represents a membership record that associates an identifier with membership information.
// It contains a key identifying the membership target and a value containing the membership details.
type Membership struct {
	Key *Mkey     `cbor:"0,keyasint,omitempty" json:"key,omitempty"`
	Val MemberVal `cbor:"1,keyasint" json:"value"`
}

// NewMembership creates a new Membership with the specified key type and value.
func NewMembership(val any, typ string) (*Membership, error) {
	keyFactory, ok := mkeyValueRegister[typ]
	if !ok {
		return nil, fmt.Errorf("unknown Mkey type: %s", typ)
	}

	key, err := keyFactory(val)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %w", err)
	}

	if err = key.Valid(); err != nil {
		return nil, fmt.Errorf("invalid key: %w", err)
	}

	var ret Membership
	ret.Key = key

	return &ret, nil
}

// MustNewMembership is like NewMembership but panics on error.
func MustNewMembership(val any, typ string) *Membership {
	ret, err := NewMembership(val, typ)
	if err != nil {
		panic(err)
	}
	return ret
}

// MustNewUUIDMembership creates a new Membership with a UUID key.
func MustNewUUIDMembership(uuid UUID) *Membership {
	return MustNewMembership(uuid, "uuid")
}

// MustNewUintMembership creates a new Membership with a uint key.
func MustNewUintMembership(u uint64) *Membership {
	return MustNewMembership(u, UintType)
}

// SetValue sets the membership value.
func (o *Membership) SetValue(val *MemberVal) *Membership {
	if o != nil {
		o.Val = *val
	}
	return o
}

func (o *Membership) RegisterExtensions(exts extensions.Map) error {
	return o.Val.RegisterExtensions(exts)
}

func (o *Membership) GetExtensions() extensions.IMapValue {
	return o.Val.GetExtensions()
}

// Valid validates the Membership.
func (o *Membership) Valid() error {
	if o.Key != nil {
		if err := o.Key.Valid(); err != nil {
			return fmt.Errorf("invalid measurement key: %w", err)
		}
	}

	return o.Val.Valid()
}

// Memberships is a container for Membership instances and their extensions.
// It is a thin wrapper around extensions.Collection.
type Memberships extensions.Collection[Membership, *Membership]

func NewMemberships() *Memberships {
	return (*Memberships)(extensions.NewCollection[Membership]())
}

func (o *Memberships) RegisterExtensions(exts extensions.Map) error {
	return (*extensions.Collection[Membership, *Membership])(o).RegisterExtensions(exts)
}

func (o *Memberships) GetExtensions() extensions.IMapValue {
	return (*extensions.Collection[Membership, *Membership])(o).GetExtensions()
}

func (o *Memberships) Valid() error {
	return (*extensions.Collection[Membership, *Membership])(o).Valid()
}

func (o *Memberships) IsEmpty() bool {
	return (*extensions.Collection[Membership, *Membership])(o).IsEmpty()
}

func (o *Memberships) Add(val *Membership) *Memberships {
	ret := (*extensions.Collection[Membership, *Membership])(o).Add(val)
	return (*Memberships)(ret)
}

func (o *Memberships) MarshalCBOR() ([]byte, error) {
	return (*extensions.Collection[Membership, *Membership])(o).MarshalCBOR()
}

func (o *Memberships) UnmarshalCBOR(data []byte) error {
	return (*extensions.Collection[Membership, *Membership])(o).UnmarshalCBOR(data)
}

func (o *Memberships) MarshalJSON() ([]byte, error) {
	return (*extensions.Collection[Membership, *Membership])(o).MarshalJSON()
}

func (o *Memberships) UnmarshalJSON(data []byte) error {
	return (*extensions.Collection[Membership, *Membership])(o).UnmarshalJSON(data)
}
