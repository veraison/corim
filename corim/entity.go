// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"fmt"

	"github.com/veraison/corim/comid"
)

// Entity stores an entity-map capable of CBOR and JSON serializations.
type Entity struct {
	EntityName string           `cbor:"0,keyasint" json:"name"`
	RegID      *comid.TaggedURI `cbor:"1,keyasint,omitempty" json:"regid,omitempty"`
	Roles      Roles            `cbor:"2,keyasint" json:"roles"`
}

func NewEntity() *Entity {
	return &Entity{}
}

// SetEntityName is used to set the EntityName field of Entity using supplied name
func (o *Entity) SetEntityName(name string) *Entity {
	if o != nil {
		if name == "" {
			return nil
		}
		o.EntityName = name
	}
	return o
}

// SetRegID is used to set the RegID field of Entity using supplied uri
func (o *Entity) SetRegID(uri string) *Entity {
	if o != nil {
		if uri == "" {
			return nil
		}

		taggedURI, err := comid.String2URI(&uri)
		if err != nil {
			return nil
		}

		o.RegID = taggedURI
	}
	return o
}

// SetRoles appends the supplied roles to the target entity.
func (o *Entity) SetRoles(roles ...Role) *Entity {
	if o != nil {
		if o.Roles.Add(roles...) == nil {
			return nil
		}
	}
	return o
}

// Valid checks for validity of the fields within each Entity
func (o Entity) Valid() error {
	if o.EntityName == "" {
		return fmt.Errorf("invalid entity: empty entity-name")
	}

	if o.RegID != nil && o.RegID.Empty() {
		return fmt.Errorf("invalid entity: empty reg-id")
	}

	if err := o.Roles.Valid(); err != nil {
		return fmt.Errorf("invalid entity: %w", err)
	}

	return nil
}

// Entities is an array of entity-map's
type Entities []Entity

// NewEntities instantiates an empty entity-map array
func NewEntities() *Entities {
	return new(Entities)
}

// AddEntity adds the supplied entity-map to the target Entities
func (o *Entities) AddEntity(e Entity) *Entities {
	if o != nil {
		*o = append(*o, e)
	}
	return o
}

// Valid iterates over the range of individual entities to check for validity
func (o Entities) Valid() error {
	for i, m := range o {
		if err := m.Valid(); err != nil {
			return fmt.Errorf("entity at index %d: %w", i, err)
		}
	}
	return nil
}
