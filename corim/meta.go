// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/veraison/corim/extensions"
)

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

func (o *Meta) RegisterExtensions(exts extensions.Map) error {
	for p, v := range exts {
		switch p {
		case ExtSigner:
			o.Signer.Register(v)
		default:
			return fmt.Errorf("%w: %q", extensions.ErrUnexpectedPoint, p)
		}
	}

	return nil
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
		return fmt.Errorf("invalid signer: %w", err)
	}

	if o.Validity != nil {
		if err := o.Validity.Valid(); err != nil {
			return fmt.Errorf("invalid validity: %w", err)
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

// FromJSON deserializes the supplied JSON data into the target Meta
func (o *Meta) FromJSON(data []byte) error {
	return json.Unmarshal(data, o)
}

// ToJSON serializes the target Meta to JSON
func (o Meta) ToJSON() ([]byte, error) {
	return json.Marshal(&o)
}
