// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"errors"
	"fmt"
	"time"

	"github.com/veraison/corim/comid"
)

type Signer struct {
	Name string           `cbor:"0,keyasint" json:"name"`
	URI  *comid.TaggedURI `cbor:"1,keyasint,omitempty" json:"uri,omitempty"`
}

func NewSigner() *Signer {
	return &Signer{}
}

// SetName sets the target Signer's name to the supplied value
func (o *Signer) SetName(name string) *Signer {
	if o != nil {
		if name == "" {
			return nil
		}
		o.Name = name
	}
	return o
}

// SetURI sets the target Signer's URI to the supplied value
func (o *Signer) SetURI(uri string) *Signer {
	if o != nil {
		if uri == "" {
			return nil
		}

		taggedURI, err := comid.String2URI(&uri)
		if err != nil {
			return nil
		}

		o.URI = taggedURI
	}
	return o
}

// Valid checks the validity of individual fields within Signer
func (o Signer) Valid() error {
	if o.Name == "" {
		return errors.New("empty name")
	}

	if o.URI != nil {
		if err := comid.IsAbsoluteURI(string(*o.URI)); err != nil {
			return fmt.Errorf("invalid URI: %w", err)
		}
	}

	return nil
}

// Meta stores a corim-meta-map with JSON and CBOR serializations.  It carries
// information about the CoRIM signer and, optionally, a validity period
// associated with the signed assertion.  A corim-meta-map is serialized to CBOR
// and added to the protected header structure in the signed-corim as a byte string
type Meta struct {
	Signer   Signer    `cbor:"0,keyasint" json:"signer"`
	Validity *Validity `cbor:"1,keyasint,omitempty" json:"validity,omitempty"`
}

func NewMeta() *Meta {
	return &Meta{}
}

// SetSigner populates the Signer element in the target Meta with the supplied
// name and optional URI
func (o *Meta) SetSigner(name string, uri *string) *Meta {
	if o != nil {
		s := NewSigner().SetName(name)

		if uri != nil {
			s = s.SetURI(*uri)
		}

		if s == nil {
			return nil
		}

		o.Signer = *s
	}
	return o
}

// SetValidity sets the validity period of the target Meta to the supplied time
// range
func (o *Meta) SetValidity(notAfter time.Time, notBefore *time.Time) *Meta {
	if o != nil {
		v := NewValidity().Set(notAfter, notBefore)
		if v == nil {
			return nil
		}

		o.Validity = v
	}
	return o
}

// Valid checks for validity of the fields within Meta
func (o Meta) Valid() error {
	if err := o.Signer.Valid(); err != nil {
		return fmt.Errorf("invalid meta: %w", err)
	}

	if o.Validity != nil {
		if err := o.Validity.Valid(); err != nil {
			return fmt.Errorf("invalid meta: %w", err)
		}
	}

	return nil
}

// ToCBOR serializes the target Meta to CBOR
func (o Meta) ToCBOR() ([]byte, error) {
	return em.Marshal(&o)
}

// FromCBOR deserializes the supplied CBOR data into the target Meta
func (o *Meta) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, o)
}
