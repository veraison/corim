// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"errors"
	"fmt"
)

type Role int64

const (
	RoleManifestCreator Role = iota + 1
)

type Roles []Role

func (o *Roles) Add(roles ...Role) *Roles {
	if o != nil {
		for _, r := range roles {
			if !isRole(r) {
				return nil
			}
			*o = append(*o, r)
		}
	}
	return o
}

func isRole(r Role) bool {
	return r == RoleManifestCreator
}

func (o Roles) Valid() error {
	if len(o) == 0 {
		return errors.New("empty roles")
	}

	for i, r := range o {
		if !isRole(r) {
			return fmt.Errorf("unknown role %d at index %d", r, i)
		}
	}

	return nil
}
