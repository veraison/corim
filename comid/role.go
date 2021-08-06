// Copyright 2021 Contributors to the Veraison project.
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

	var r Role

	for _, s := range a {
		switch s {
		case "tagCreator":
			r = RoleTagCreator
		case "creator":
			r = RoleCreator
		case "maintainer":
			r = RoleMaintainer
		default:
			return fmt.Errorf("unknown role '%s'", s)
		}
		o = o.Add(r)
	}

	return nil
}
