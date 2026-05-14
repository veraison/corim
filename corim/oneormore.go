// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package corim

import (
	"encoding/json"
	"errors"
)

// OneOrMore is a slice that serializes as either a single item when it
// contains only one element, or an array of items if it contains more than one
// element. It is invalid for OneOrMore to contain no elements.
type OneOrMore[T any] []T

// NewOneOrMore returns a new one OneOrMore contained the provided value.
func NewOneOrMore[T any](v T) *OneOrMore[T] {
	return &OneOrMore[T]{v}
}

// Add the specified values to the OneOrMore.
func (o *OneOrMore[T]) Add(v ...T) *OneOrMore[T] {
	*o = append(*o, v...)
	return o
}

// Valid returns an error if the OneOrMore or more is invalid (i.e. contains no
// elements).
func (o OneOrMore[T]) Valid() error {
	if len(o) == 0 {
		return errors.New("must have at least one")
	}

	return nil
}

func (o OneOrMore[T]) MarshalCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	if len(o) == 1 {
		return em.Marshal(o[0])
	} else {
		return em.Marshal([]T(o))
	}
}

func (o *OneOrMore[T]) UnmarshalCBOR(data []byte) error {
	var vals []T
	if err := dm.Unmarshal(data, &vals); err == nil {
		*o = vals
		return nil
	}

	var one T
	if err := dm.Unmarshal(data, &one); err != nil {
		return err
	}

	*o = []T{one}

	return nil
}

func (o OneOrMore[T]) MarshalJSON() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	if len(o) == 1 {
		return json.Marshal(o[0])
	} else {
		return json.Marshal([]T(o))
	}
}

func (o *OneOrMore[T]) UnmarshalJSON(data []byte) error {
	var vals []T
	if err := json.Unmarshal(data, &vals); err == nil {
		*o = vals
		return nil
	}

	var one T
	if err := json.Unmarshal(data, &one); err != nil {
		return err
	}

	*o = []T{one}

	return nil
}
