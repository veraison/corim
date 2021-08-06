// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "fmt"

type TaggedURI string

// Entity stores an entity-map capable of CBOR and JSON serializations.
type Entity struct {
	EntityName string     `cbor:"0,keyasint" json:"name"`
	RegID      *TaggedURI `cbor:"1,keyasint,omitempty" json:"regid,omitempty"`
	Roles      Roles      `cbor:"2,keyasint" json:"roles"`
}

func (o *Entity) SetEntityName(name string) *Entity {
	if o != nil {
		if name == "" {
			return nil
		}
		o.EntityName = name
	}
	return o
}

func (o *Entity) SetRegID(uri string) *Entity {
	if o != nil {
		if uri == "" {
			return nil
		}
		taggedURI := TaggedURI(uri)
		o.RegID = &taggedURI
	}
	return o
}

func (o *Entity) SetRoles(roles ...Role) *Entity {
	if o != nil {
		o.Roles.Add(roles...)
	}
	return o
}

func (o Entity) Valid() error {
	if o.EntityName == "" {
		return fmt.Errorf("invalid entity: empty entity-name")
	}

	if o.RegID != nil && *o.RegID == "" {
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

func (o Entities) Valid() error {
	for i, m := range o {
		if err := m.Valid(); err != nil {
			return fmt.Errorf("entity at index %d: %w", i, err)
		}
	}
	return nil
}
