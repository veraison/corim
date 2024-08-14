// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"

	"github.com/veraison/corim/extensions"
)

// ValueTriple relates a measurement to a target environment, essentially
// forming a subject-predicate-object triple of
// "measurement-pertains-to-environment". This structure is used to represent
// both reference-triple-record and endorsed-triple-record in the CoRIM spec
// (as of rev. 04).
type ValueTriple struct {
	_           struct{}    `cbor:",toarray"`
	Environment Environment `json:"environment"`
	Measurement Measurement `json:"measurement"`
}

func (o *ValueTriple) RegisterExtensions(exts extensions.Map) error {
	return o.Measurement.RegisterExtensions(exts)
}

func (o *ValueTriple) GetExtensions() extensions.IMapValue {
	return o.Measurement.GetExtensions()
}

func (o ValueTriple) Valid() error {
	if err := o.Environment.Valid(); err != nil {
		return fmt.Errorf("environment validation failed: %w", err)
	}

	if err := o.Measurement.Valid(); err != nil {
		return fmt.Errorf("measurement validation failed: %w", err)
	}

	return nil
}

// ValueTriples is a container for ValueTriple instances and their extensions.
// It is a thin wrapper around extensions.Collection.
type ValueTriples extensions.Collection[ValueTriple, *ValueTriple]

func NewValueTriples() *ValueTriples {
	return (*ValueTriples)(extensions.NewCollection[ValueTriple]())
}

func (o *ValueTriples) RegisterExtensions(exts extensions.Map) error {
	return (*extensions.Collection[ValueTriple, *ValueTriple])(o).RegisterExtensions(exts)
}

func (o *ValueTriples) GetExtensions() extensions.IMapValue {
	return (*extensions.Collection[ValueTriple, *ValueTriple])(o).GetExtensions()
}

func (o ValueTriples) Valid() error {
	return (extensions.Collection[ValueTriple, *ValueTriple])(o).Valid()
}

func (o *ValueTriples) IsEmpty() bool {
	return (*extensions.Collection[ValueTriple, *ValueTriple])(o).IsEmpty()
}

func (o *ValueTriples) Add(val *ValueTriple) *ValueTriples {
	ret := (*extensions.Collection[ValueTriple, *ValueTriple])(o).Add(val)
	return (*ValueTriples)(ret)
}

func (o ValueTriples) MarshalCBOR() ([]byte, error) {
	return (extensions.Collection[ValueTriple, *ValueTriple])(o).MarshalCBOR()
}

func (o *ValueTriples) UnmarshalCBOR(data []byte) error {
	return (*extensions.Collection[ValueTriple, *ValueTriple])(o).UnmarshalCBOR(data)
}

func (o ValueTriples) MarshalJSON() ([]byte, error) {
	return (extensions.Collection[ValueTriple, *ValueTriple])(o).MarshalJSON()
}

func (o *ValueTriples) UnmarshalJSON(data []byte) error {
	return (*extensions.Collection[ValueTriple, *ValueTriple])(o).UnmarshalJSON(data)
}
