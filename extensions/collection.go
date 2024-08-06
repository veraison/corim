// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package extensions

import (
	"encoding/json"
	"errors"
	"fmt"

	cbor "github.com/fxamacker/cbor/v2"
)

var ErrUnexpectedPoint = errors.New("unexpected extension point")

// Point defines an extension point -- a place where a struct defining
// extensions may be attached. There is at least one Point corresponding to
// each struct that embeds Extensions.
type Point string

// Map is a mapping of extension Point's to IMapValue's (pointers to structs that
// define extension fields.
type Map map[Point]IMapValue

// NewMap instantiates an empty Map
func NewMap() Map {
	return make(Map)
}

// Add a Point-IMapValue entry to the map. If there is already an IMapValue
// associated with the Point, it will be replaced.
func (o Map) Add(p Point, v IMapValue) Map {
	o[p] = v
	return o
}

// IExtensible defines an interface for extensible objects (those that can
// register extensions.
type IExtensible[P any] interface {
	RegisterExtensions(exts Map) error
	GetExtensions() IMapValue
	Valid() error

	*P // this interface must be implemented by a pointer to P
}

// Collection is a generic container for objects who's pointers implement
// IExtensible. In addition of containing instances of extensible objects, it
// is capable of registering extensions for those instances.
type Collection[P any, I IExtensible[P]] struct {
	Values []P

	valueExtensions cache
}

// NewCollection returns a pointer to a new Collections instance.
func NewCollection[P any, I IExtensible[P]]() *Collection[P, I] {
	return &Collection[P, I]{Values: []P{}, valueExtensions: cache{make(Map)}}
}

// Add a new element to the collection.
func (o *Collection[P, I]) Add(p *P) *Collection[P, I] {
	if o != nil && p != nil {
		if !o.valueExtensions.IsEmpty() {
			var m I = p // #nosec G601 -- not an issue in Go 1.22
			// only register cached extensions if the object does
			// not have extensions of its own.
			if m.GetExtensions() == nil {
				exts := o.valueExtensions.Get()
				if err := m.RegisterExtensions(exts); err != nil {
					// o.valueExtensions have been validated when
					// they were set inside o.RegisterExtensions(),
					// and as the field is not exported, could not
					// have been modified since. Therefore, this
					// cannot fail.
					panic(err)
				}
			}
		}
		o.Values = append(o.Values, *p)
	}

	return o
}

// Clear empties the collection of its values. This does _not_ unregister
// extensions.
func (o *Collection[P, I]) Clear() {
	// Setting Values to an empty array rather then nil to ensure correct
	// marshaling.
	o.Values = []P{}
}

// GetExtensions returns the extensions IMapValue that has been registered with
// the collection.
func (o *Collection[P, I]) GetExtensions() IMapValue {
	if o.valueExtensions.IsEmpty() {
		return nil
	}

	return o.valueExtensions.Get()
}

// Valid returns an error if the collection is invalid, i.e. if it is empty or
// if any of its contents are invalid.
func (o Collection[P, I]) Valid() error {
	for i, p := range o.Values {
		var m I = &p // #nosec G601 -- not an issue in Go 1.22
		if err := m.Valid(); err != nil {
			return fmt.Errorf("error at index %d: %w", i, err)
		}
	}

	return nil
}

// IsEmpty return true if the collection is empty, i.e. if it does not contain
// any values. This does not indicate whether or not any extensions have been
// registered with the collection.
func (o Collection[P, I]) IsEmpty() bool {
	return len(o.Values) == 0
}

// RegisterExtensions register extensions for the collection's values. An error
// is returned if the provided map contains extension points not supported by
// the collection's values.
func (o *Collection[P, I]) RegisterExtensions(exts Map) error {
	// validate the provided extensions by creating a new instanced of the
	// contained object and attempting to register them with it.
	var m I = new(P)
	if err := m.RegisterExtensions(exts); err != nil {
		return err
	}

	o.valueExtensions.Set(exts)

	for i := 0; i < len(o.Values); i++ {
		var vi I = &o.Values[i]
		if vi.GetExtensions() == nil {
			if err := vi.RegisterExtensions(exts); err != nil {
				return fmt.Errorf("error at index %d: %w", i, err)
			}
		}
	}

	return nil
}

func (o Collection[P, I]) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.Values)
}

func (o *Collection[P, I]) UnmarshalCBOR(data []byte) error {
	var rawVals []cbor.RawMessage

	if err := cbor.Unmarshal(data, &rawVals); err != nil {
		return err
	}

	vals := make([]P, len(rawVals))

	for i, rv := range rawVals {
		var m I = new(P)

		if !o.valueExtensions.IsEmpty() {
			if err := m.RegisterExtensions(o.valueExtensions.Get()); err != nil {
				return err
			}
		}

		if err := cbor.Unmarshal(rv, m); err != nil {
			return fmt.Errorf("error at index %d: %w", i, err)
		}

		vals[i] = *m
	}

	o.Values = vals

	return nil
}

func (o Collection[P, I]) MarshalJSON() ([]byte, error) {
	ret, err := json.Marshal(o.Values)
	return ret, err
}

func (o *Collection[P, I]) UnmarshalJSON(data []byte) error {
	var rawVals []json.RawMessage

	if err := json.Unmarshal(data, &rawVals); err != nil {
		return err
	}

	vals := make([]P, len(rawVals))

	for i, rv := range rawVals {
		var m I = new(P)

		if !o.valueExtensions.IsEmpty() {
			if err := m.RegisterExtensions(o.valueExtensions.Get()); err != nil {
				return fmt.Errorf("could not register extensions: %w", err)
			}
		}

		if err := json.Unmarshal(rv, m); err != nil {
			return fmt.Errorf("error at index %d: %w", i, err)
		}

		vals[i] = *m
	}

	o.Values = vals

	return nil
}

type cache struct {
	extensions Map
}

func (o *cache) Set(exts Map) {
	o.extensions = exts
}

func (o cache) Get() Map {
	res := make(Map)

	for p, v := range o.extensions {
		res[p] = newIMapValue(v)
	}

	return res
}

func (o cache) IsEmpty() bool {
	return len(o.extensions) == 0
}
