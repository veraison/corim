// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"encoding/json"
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
)

// Entity stores an entity-map capable of CBOR and JSON serializations.
type Entity struct {
	Name  *EntityName      `cbor:"0,keyasint" json:"name"`
	RegID *comid.TaggedURI `cbor:"1,keyasint,omitempty" json:"regid,omitempty"`
	Roles Roles            `cbor:"2,keyasint" json:"roles"`

	Extensions
}

func NewEntity() *Entity {
	return &Entity{}
}

// RegisterExtensions registers a struct as a collections of extensions
func (o *Entity) RegisterExtensions(exts extensions.Map) error {
	for p, v := range exts {
		switch p {
		case ExtEntity:
			o.Register(v)
		default:
			return fmt.Errorf("%w: %q", extensions.ErrUnexpectedPoint, p)
		}
	}

	return nil
}

// GetExtensions returns pervisouosly registered extension
func (o *Entity) GetExtensions() extensions.IMapValue {
	return o.IMapValue
}

// SetName is used to set the EntityName field of Entity using supplied name
func (o *Entity) SetName(name any) *Entity {

	if o != nil {
		if name == "" {
			return nil
		}
		o.Name = MustNewStringEntityName(name)
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
	if o.Name == nil {
		return fmt.Errorf("invalid entity: empty entity-name")
	}

	if err := o.Name.Valid(); err != nil {
		return fmt.Errorf("invalid entity: %w", err)
	}

	if o.RegID != nil && o.RegID.Empty() {
		return fmt.Errorf("invalid entity: empty reg-id")
	}

	if err := o.Roles.Valid(); err != nil {
		return fmt.Errorf("invalid entity: %w", err)
	}

	return o.validEntity(&o)
}

// UnmarshalCBOR deserializes from CBOR
func (o *Entity) UnmarshalCBOR(data []byte) error {
	return encoding.PopulateStructFromCBOR(dm, data, o)
}

// MarshalCBOR serializes to CBOR
func (o Entity) MarshalCBOR() ([]byte, error) {
	return encoding.SerializeStructToCBOR(em, o)
}

// UnmarshalJSON deserializes from JSON
func (o *Entity) UnmarshalJSON(data []byte) error {
	return encoding.PopulateStructFromJSON(data, o)
}

// MarshalJSON serializes to JSON
func (o Entity) MarshalJSON() ([]byte, error) {
	return encoding.SerializeStructToJSON(o)
}

// Entities is a container for Entity instances and their extensions.
// It is a thin wrapper around extensions.Collection.
type Entities extensions.Collection[Entity, *Entity]

func NewEntities() *Entities {
	return (*Entities)(extensions.NewCollection[Entity]())
}

func (o *Entities) RegisterExtensions(exts extensions.Map) error {
	return (*extensions.Collection[Entity, *Entity])(o).RegisterExtensions(exts)
}

func (o *Entities) GetExtensions() extensions.IMapValue {
	return (*extensions.Collection[Entity, *Entity])(o).GetExtensions()
}

func (o *Entities) Valid() error {
	return (*extensions.Collection[Entity, *Entity])(o).Valid()
}

func (o *Entities) IsEmpty() bool {
	return (*extensions.Collection[Entity, *Entity])(o).IsEmpty()
}

func (o *Entities) Add(val *Entity) *Entities {
	ret := (*extensions.Collection[Entity, *Entity])(o).Add(val)
	return (*Entities)(ret)
}

func (o Entities) MarshalCBOR() ([]byte, error) {
	return (extensions.Collection[Entity, *Entity])(o).MarshalCBOR()
}

func (o *Entities) UnmarshalCBOR(data []byte) error {
	return (*extensions.Collection[Entity, *Entity])(o).UnmarshalCBOR(data)
}

func (o Entities) MarshalJSON() ([]byte, error) {
	return (extensions.Collection[Entity, *Entity])(o).MarshalJSON()
}

func (o *Entities) UnmarshalJSON(data []byte) error {
	return (*extensions.Collection[Entity, *Entity])(o).UnmarshalJSON(data)
}

// EntityName encapsulates the name of the associated Entity. The CoRIM
// specification only allows for text (string) name, but this may be extended
// by other specifications.
type EntityName struct {
	Value IEntityNameValue
}

// NewEntityName creates a new EntityName of the specified type using the
// provided value.
func NewEntityName(val any, typ string) (*EntityName, error) {
	factory, ok := entityNameValueRegister[typ]
	if !ok {
		return nil, fmt.Errorf("unexpected entity name type: %s", typ)
	}

	return factory(val)
}

// MustNewEntityName is like NewEntityName, except it doesn't return an error,
// assuming that the provided value is valid. It panics if that isn't the case.
func MustNewEntityName(val any, typ string) *EntityName {
	ret, err := NewEntityName(val, typ)
	if err != nil {
		panic(err)
	}

	return ret
}

// String returns the string representation of the EntityName
func (o EntityName) String() string {
	return o.Value.String()
}

// Valid returns nil if the underlying EntityName value is valid, or an error
// describing the problem otherwise.
func (o EntityName) Valid() error {
	if o.Value == nil {
		return errors.New("empty entity name")
	}

	return o.Value.Valid()
}

// MarshalCBOR serializes the EntityName into CBOR-encoded bytes.
func (o EntityName) MarshalCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return em.Marshal(o.Value)
}

// UnmarshalCBOR deserializes the EntityName from CBOR-encoded bytes.
func (o *EntityName) UnmarshalCBOR(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty")
	}

	majorType := (data[0] & 0xe0) >> 5
	if majorType == 3 { // text string
		var text string

		if err := dm.Unmarshal(data, &text); err != nil {
			return err
		}

		name := StringEntityName(text)
		o.Value = &name

		return nil
	}

	return dm.Unmarshal(data, &o.Value)
}

// MarshalJSON serializes the EntityName into a JSON object.
func (o EntityName) MarshalJSON() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	if o.Value.Type() == extensions.StringType {
		return json.Marshal(o.Value.String())
	}

	valueBytes, err := json.Marshal(o.Value)
	if err != nil {
		return nil, err
	}

	value := encoding.TypeAndValue{
		Type:  o.Value.Type(),
		Value: valueBytes,
	}

	return json.Marshal(value)
}

// UnmarshalJSON deserializes EntityName from the provided JSON object.
func (o *EntityName) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		*o = *MustNewStringEntityName(text)
		return nil
	}

	var tnv encoding.TypeAndValue

	if err := json.Unmarshal(data, &tnv); err != nil {
		return fmt.Errorf("entity name decoding failure: %w", err)
	}

	decoded, err := NewEntityName(nil, tnv.Type)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(tnv.Value, &decoded.Value); err != nil {
		return fmt.Errorf(
			"cannot unmarshal entity name: %w",
			err,
		)
	}

	if err := decoded.Value.Valid(); err != nil {
		return fmt.Errorf("invalid %s: %w", tnv.Type, err)
	}

	o.Value = decoded.Value

	return nil
}

// IEntityNameValue is the interface implemented by all EntityName value types.
type IEntityNameValue interface {
	extensions.ITypeChoiceValue
}

// StringEntityName is a text string EntityName with no other contraints. This
// is the only EntityName value type defined by the CoRIM specification itself.
type StringEntityName string

func NewStringEntityName(val any) (*EntityName, error) {
	var ret StringEntityName

	if val == nil {
		ret = StringEntityName("")
		return &EntityName{&ret}, nil
	}

	switch t := val.(type) {
	case string:
		ret = StringEntityName(t)
	case []byte:
		if !utf8.Valid(t) {
			return nil, errors.New("bytes do not form a valid UTF-8 string")
		}

		ret = StringEntityName(t)
	default:
		return nil, fmt.Errorf("unexpected type for string entity name: %T", t)
	}

	return &EntityName{&ret}, nil
}

func MustNewStringEntityName(val any) *EntityName {
	ret, err := NewStringEntityName(val)
	if err != nil {
		panic(err)
	}

	return ret
}

func (o StringEntityName) String() string {
	return string(o)
}

func (o StringEntityName) Type() string {
	return extensions.StringType
}

func (o StringEntityName) Valid() error {
	if o == "" {
		return errors.New("empty entity-name")
	}

	return nil
}

// IEntityNameFactory defines the signature for the factory functions that may
// be registred using RegisterEntityNameType to provide a new implementation of
// the corresponding type choice. The factory function should create a new
// *EntityName with the underlying value created based on the provided input.
// The range of valid inputs is up to the specific type choice implementation,
// however it _must_ accept nil as one of the inputs, and return the Zero value
// for implemented type.
// See also https://go.dev/ref/spec#The_zero_value
type IEntityNameFactory func(any) (*EntityName, error)

var entityNameValueRegister = map[string]IEntityNameFactory{
	extensions.StringType: NewStringEntityName,
}

// RegisterEntityNameType registers a new IEntityNameValue implementation
// (created by the provided IEntityNameFactory) under the specified type name
// and CBOR tag.
func RegisterEntityNameType(tag uint64, factory IEntityNameFactory) error {

	nilVal, err := factory(nil)
	if err != nil {
		return err
	}

	typ := nilVal.Value.Type()
	if _, exists := entityNameValueRegister[typ]; exists {
		return fmt.Errorf("entity name type with name %q already exists", typ)
	}

	if err := registerCORIMTag(tag, nilVal.Value); err != nil {
		return err
	}

	entityNameValueRegister[typ] = factory

	return nil
}
