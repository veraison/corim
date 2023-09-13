// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
)

// Class represents the class of the (target / attesting) environment.  The only
// required field is the class unique identifier (see ClassID).  Optionally,
// information about the specific brand & product as well as its topological
// coordinates within the wider device can be recorded.
type Class struct {
	ClassID *ClassID `cbor:"0,keyasint,omitempty" json:"id,omitempty"`
	Vendor  *string  `cbor:"1,keyasint,omitempty" json:"vendor,omitempty"`
	Model   *string  `cbor:"2,keyasint,omitempty" json:"model,omitempty"`
	Layer   *uint64  `cbor:"3,keyasint,omitempty" json:"layer,omitempty"`
	Index   *uint64  `cbor:"4,keyasint,omitempty" json:"index,omitempty"`
}

// NewClassUUID instantiates a new Class object with the specified UUID as
// identifier
func NewClassUUID(uuid UUID) *Class {
	classID, err := NewUUIDClassID(uuid)
	if err != nil {
		return nil
	}

	return &Class{ClassID: classID}
}

// NewClassImplID instantiates a new Class object that identifies the specified PSA
// Implementation ID
func NewClassImplID(implID ImplID) *Class {
	classID, err := NewImplIDClassID(implID)
	if err != nil {
		return nil
	}

	return &Class{ClassID: classID}
}

// NewClassOID instantiates a new Class object that identifies the OID
func NewClassOID(oid string) *Class {
	classID, err := NewOIDClassID(oid)
	if err != nil {
		return nil
	}

	return &Class{ClassID: classID}
}

// SetVendor sets the vendor metadata to the supplied string
func (o *Class) SetVendor(vendor string) *Class {
	if o != nil {
		o.Vendor = &vendor
	}
	return o
}

// GetVendor returns the vendor string if it set in the target Class.
// Otherwise, an empty string is returned.
func (o Class) GetVendor() string {
	if o.Vendor == nil {
		return ""
	}
	return *o.Vendor
}

// GetModel returns the model string if it set in the target Class.
// Otherwise, an empty string is returned.
func (o Class) GetModel() string {
	if o.Model == nil {
		return ""
	}
	return *o.Model
}

// GetLayer returns the layer number if it set in the target Class.
// Otherwise, uint64_max is returned.
func (o Class) GetLayer() uint64 {
	if o.Layer == nil {
		return ^uint64(0)
	}
	return *o.Layer
}

// GetIndex returns the index number if it set in the target Class.
// Otherwise, uint64_max is returned.
func (o Class) GetIndex() uint64 {
	if o.Layer == nil {
		return ^uint64(0)
	}
	return *o.Index
}

// SetModel sets the model metadata to the supplied string
func (o *Class) SetModel(model string) *Class {
	if o != nil {
		o.Model = &model
	}
	return o
}

// SetLayer sets the "layer" (i.e., the logical/topological location of the
// environment in the device) as indicated
func (o *Class) SetLayer(layer uint64) *Class {
	if o != nil {
		o.Layer = &layer
	}
	return o
}

// SetIndex sets the "index" (i.e., the identifier of the environment instance
// in a specific layer) as indicated
func (o *Class) SetIndex(index uint64) *Class {
	if o != nil {
		o.Index = &index
	}
	return o
}

// Valid checks the non-empty<> constraint on the map
func (o Class) Valid() error {
	// check non-empty<{ ... }>
	if (o.ClassID == nil || !o.ClassID.IsSet()) &&
		o.Vendor == nil && o.Model == nil && o.Layer == nil && o.Index == nil {
		return fmt.Errorf("class must not be empty")
	}
	return nil
}

// ToCBOR serializes the target Class to CBOR (if the Class is "valid")
func (o Class) ToCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return em.Marshal(&o)
}

// FromCBOR deserializes the supplied CBOR data into the target Class
func (o *Class) FromCBOR(data []byte) error {
	if err := dm.Unmarshal(data, o); err != nil {
		return err
	}

	return o.Valid()
}

// FromJSON deserializes the supplied JSON string into the target Class
func (o *Class) FromJSON(data []byte) error {
	if err := json.Unmarshal(data, o); err != nil {
		return err
	}

	return o.Valid()
}

// ToJSON serializes the target Class to JSON (if the Class is "valid")
func (o Class) ToJSON() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return json.Marshal(&o)
}
