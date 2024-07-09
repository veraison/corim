// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Rel int64

const (
	RelSupplements Rel = iota
	RelReplaces

	RelUnset = ^Rel(0)
)

var (
	relToString = map[Rel]string{
		RelReplaces:    "replaces",
		RelSupplements: "supplements",
	}

	stringToRel = map[string]Rel{
		"replaces":    RelReplaces,
		"supplements": RelSupplements,
	}
)

// RegisterRel creates a new Rel association between the provided value and
// name. An error is returned if either clashes with any of the existing roles.
func RegisterRel(val int64, name string) error {
	rel := Rel(val)

	if _, ok := relToString[rel]; ok {
		return fmt.Errorf("rel with value %d already exists", val)
	}

	if _, ok := stringToRel[name]; ok {
		return fmt.Errorf("rel with name %q already exists", name)
	}

	relToString[rel] = name
	stringToRel[name] = rel

	return nil
}

func NewRel() *Rel {
	r := RelUnset
	return &r
}

func (o *Rel) Set(r Rel) *Rel {
	if o != nil {
		*o = r
	}
	return o
}

func (o Rel) Get() Rel {
	return o
}

func (o Rel) Valid() error {
	if o == RelUnset {
		return errors.New("rel is unset")
	}

	return nil
}

func (o Rel) String() string {
	ret, ok := relToString[o]
	if ok {
		return ret
	}

	return fmt.Sprintf("rel(%d)", o)
}

func (o Rel) ToCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	data, err := em.Marshal(o)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (o *Rel) FromCBOR(data []byte) error {
	err := dm.Unmarshal(data, o)
	if err != nil {
		return err
	}

	return o.Valid()
}

func (o *Rel) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("cannot unmarshal rel: %w", err)
	}

	if s == "" {
		return fmt.Errorf("empty rel")
	}

	switch s {
	case "supplements":
		*o = RelSupplements
	case "replaces":
		*o = RelReplaces
	default:
		return fmt.Errorf("unknown rel '%s'", s)
	}

	return nil
}

func (o Rel) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}
