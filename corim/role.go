// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Role int64

const (
	RoleManifestCreator Role = iota + 1
)

var (
	stringToRole = map[string]Role{
		"manifestCreator": RoleManifestCreator,
	}
)

type Roles []Role

// Add appends the supplied roles to Roles list.
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

// Valid iterates over the range of individual roles to check for validity
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
