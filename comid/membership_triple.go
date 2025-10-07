// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"errors"
	"fmt"

	"github.com/veraison/corim/extensions"
)

// MembershipTriple relates membership information to a target environment,
// essentially forming a subject-predicate-object triple of "memberships-pertain
// to-environment". This structure is used to represent membership-triple-record
// in the CoRIM spec.
type MembershipTriple struct {
	_           struct{}    `cbor:",toarray"`
	Environment Environment `json:"environment"`
	Memberships Memberships `json:"memberships"`
}

func (o *MembershipTriple) RegisterExtensions(exts extensions.Map) error {
	return o.Memberships.RegisterExtensions(exts)
}

func (o *MembershipTriple) GetExtensions() extensions.IMapValue {
	return o.Memberships.GetExtensions()
}

func (o *MembershipTriple) Valid() error {
	if err := o.Environment.Valid(); err != nil {
		return fmt.Errorf("environment validation failed: %w", err)
	}

	if o.Memberships.IsEmpty() {
		return errors.New("memberships validation failed: no membership entries")
	}

	if err := o.Memberships.Valid(); err != nil {
		return fmt.Errorf("memberships validation failed: %w", err)
	}

	return nil
}

// MembershipTriples is a container for MembershipTriple instances and their extensions.
// It is a thin wrapper around extensions.Collection.
type MembershipTriples extensions.Collection[MembershipTriple, *MembershipTriple]

func NewMembershipTriples() *MembershipTriples {
	return (*MembershipTriples)(extensions.NewCollection[MembershipTriple]())
}

func (o *MembershipTriples) RegisterExtensions(exts extensions.Map) error {
	return (*extensions.Collection[MembershipTriple, *MembershipTriple])(o).RegisterExtensions(exts)
}

func (o *MembershipTriples) GetExtensions() extensions.IMapValue {
	return (*extensions.Collection[MembershipTriple, *MembershipTriple])(o).GetExtensions()
}

func (o *MembershipTriples) Valid() error {
	return (*extensions.Collection[MembershipTriple, *MembershipTriple])(o).Valid()
}

func (o *MembershipTriples) IsEmpty() bool {
	return (*extensions.Collection[MembershipTriple, *MembershipTriple])(o).IsEmpty()
}

func (o *MembershipTriples) Add(val *MembershipTriple) *MembershipTriples {
	ret := (*extensions.Collection[MembershipTriple, *MembershipTriple])(o).Add(val)
	return (*MembershipTriples)(ret)
}

func (o *MembershipTriples) MarshalCBOR() ([]byte, error) {
	return (*extensions.Collection[MembershipTriple, *MembershipTriple])(o).MarshalCBOR()
}

func (o *MembershipTriples) UnmarshalCBOR(data []byte) error {
	return (*extensions.Collection[MembershipTriple, *MembershipTriple])(o).UnmarshalCBOR(data)
}

func (o *MembershipTriples) MarshalJSON() ([]byte, error) {
	return (*extensions.Collection[MembershipTriple, *MembershipTriple])(o).MarshalJSON()
}

func (o *MembershipTriples) UnmarshalJSON(data []byte) error {
	return (*extensions.Collection[MembershipTriple, *MembershipTriple])(o).UnmarshalJSON(data)
}
