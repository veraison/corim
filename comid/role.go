// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
)

type Role int64

/*
$comid-role-type-choice /= comid.tag-creator
$comid-role-type-choice /= comid.creator
$comid-role-type-choice /= comid.maintainer

comid.tag-creator = 0
comid.creator = 1
comid.maintainer = 2
*/

const (
	RoleTagCreator Role = iota
	RoleCreator
	RoleMaintainer
)

var (
	roleToString = map[Role]string{
		RoleTagCreator: "tagCreator",
		RoleCreator:    "creator",
		RoleMaintainer: "maintainer",
	}

	stringToRole = map[string]Role{
		"tagCreator": RoleTagCreator,
		"creator":    RoleCreator,
		"maintainer": RoleMaintainer,
	}
)

// String returns the string representation of the Role.
func (o Role) String() string {
	text, ok := roleToString[o]
	if ok {
		return text
	}

	return fmt.Sprintf("Role(%d)", o)
}

// RegisterRole creates a new Role association between the provided value and
// name. An error is returned if either clashes with any of the existing roles.
func RegisterRole(val int64, name string) error {
	role := Role(val)

	if _, ok := roleToString[role]; ok {
		return fmt.Errorf("role with value %d already exists", val)
	}

	if _, ok := stringToRole[name]; ok {
		return fmt.Errorf("role with name %q already exists", name)
	}

	roleToString[role] = name
	stringToRole[name] = role

	return nil
}

type Roles []Role

func NewRoles() *Roles {
	return new(Roles)
}

func (o *Roles) Add(roles ...Role) *Roles {
	if o != nil {
		*o = append(*o, roles...)
	}

	return o
}

func (o Roles) Valid() error {
	if len(o) == 0 {
		return fmt.Errorf("empty roles")
	}

	return nil
}

func (o Roles) ToCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	data, err := em.Marshal(o)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (o *Roles) FromCBOR(data []byte) error {
	err := dm.Unmarshal(data, o)
	if err != nil {
		return err
	}

	return o.Valid()
}

func (o *Roles) UnmarshalJSON(data []byte) error {
	var a []string

	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	if len(a) == 0 {
		return fmt.Errorf("no roles found")
	}

	for _, s := range a {
		r, ok := stringToRole[s]
		if !ok {
			return fmt.Errorf("unknown role %q", s)
		}
		o = o.Add(r)
	}

	return nil
}

func (o Roles) MarshalJSON() ([]byte, error) {
	roles := []string{}

	for _, r := range o {
		s, ok := roleToString[r]
		if !ok {
			return nil, fmt.Errorf("unknown role %d", r)
		}
		roles = append(roles, s)
	}

	return json.Marshal(roles)
}
