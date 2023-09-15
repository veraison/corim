// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"encoding/json"
	"fmt"
)

type EntityExtension struct {
	Param1 string `cbor:"0,keyasint" json:"param1"`
	Param2 byte   `cbor:"1,keyasint,omitempty" json:"param2,omitempty"`
}

type ProfileXEntity struct {
	Entity
	Extension EntityExtension `cbor:"8,keyasint" json:"extension"`
}

func NewProfileXEntity() *ProfileXEntity {
	return &ProfileXEntity{}
}

func (e *ProfileXEntity) SetEntityExtension(p1 string, p2 byte) {

	e.Extension.Param1 = p1
	e.Extension.Param2 = p2
}

func (o *ProfileXEntity) FromCBOR(data []byte) error {
	if err := o.Valid(); err != nil {
		return fmt.Errorf("invalid Profile XEntity %w", err)
	}
	return dm.Unmarshal(data, o)
}

func (o *ProfileXEntity) ToCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, fmt.Errorf("invalid Profile XEntity %w", err)
	}
	return em.Marshal(o)
}

func (o *ProfileXEntity) FromJSON(data []byte) error {
	if err := o.Valid(); err != nil {
		return err
	}
	return json.Unmarshal(data, o)
}

func (o *ProfileXEntity) ToJSON() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}
	return json.Marshal(o)
}
