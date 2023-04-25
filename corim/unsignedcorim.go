// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/veraison/corim/cots"

	"github.com/veraison/corim/comid"
	"github.com/veraison/eat"
	"github.com/veraison/swid"
)

// UnsignedCorim is the top-level representation of the unsigned-corim-map with
// CBOR and JSON serialization.
type UnsignedCorim struct {
	ID            swid.TagID     `cbor:"0,keyasint" json:"corim-id"`
	Tags          []Tag          `cbor:"1,keyasint" json:"tags"`
	DependentRims *[]Locator     `cbor:"2,keyasint,omitempty" json:"dependent-rims,omitempty"`
	Profiles      *[]eat.Profile `cbor:"3,keyasint,omitempty" json:"profiles,omitempty"`
	RimValidity   *Validity      `cbor:"4,keyasint,omitempty" json:"validity,omitempty"`
	Entities      *Entities      `cbor:"5,keyasint,omitempty" json:"entities,omitempty"`
}

// NewUnsignedCorim instantiates an empty UnsignedCorim
func NewUnsignedCorim() *UnsignedCorim {
	return &UnsignedCorim{}
}

// SetID sets the corim-id in the unsigned-corim-map to the supplied value.  The
// corim-id can be passed as UUID in string or binary form (i.e., byte array),
// or as a (non-empty) string
func (o *UnsignedCorim) SetID(v interface{}) *UnsignedCorim {
	if o != nil {
		tagID := swid.NewTagID(v)
		if tagID == nil {
			return nil
		}
		o.ID = *tagID
	}
	return o
}

// GetID retrieves the corim-id from the unsigned-corim-map as a string
func (o UnsignedCorim) GetID() string {
	return o.ID.String()
}

// AddComid appends the CBOR encoded (and appropriately tagged) CoMID to the
// tags array of the unsigned-corim-map
func (o *UnsignedCorim) AddComid(c comid.Comid) *UnsignedCorim {
	if o != nil {
		if c.Valid() != nil {
			return nil
		}

		comidCBOR, err := c.ToCBOR()
		if err != nil {
			return nil
		}

		taggedComid := append(ComidTag, comidCBOR...)

		o.Tags = append(o.Tags, taggedComid)
	}
	return o
}

// AddCots appends the CBOR encoded (and appropriately tagged) CoTS to the
// tags array of the unsigned-corim-map
func (o *UnsignedCorim) AddCots(c cots.ConciseTaStore) *UnsignedCorim {
	if o != nil {
		if c.Valid() != nil {
			return nil
		}

		cotsCBOR, err := c.ToCBOR()
		if err != nil {
			return nil
		}

		taggedCots := append(cots.CotsTag, cotsCBOR...)

		o.Tags = append(o.Tags, taggedCots)
	}
	return o
}

// AddCoswid appends the CBOR encoded (and appropriately tagged) CoSWID to the
// tags array of the unsigned-corim-map
func (o *UnsignedCorim) AddCoswid(c swid.SoftwareIdentity) *UnsignedCorim {
	if o != nil {
		// Currently the swid package doesn't offer an interface
		// for validating the supplied CoSWID, so -- for now --
		// we take any input for granted and pass it to the encoder.
		// See also https://github.com/veraison/swid/issues/23.

		coswidCBOR, err := c.ToCBOR()
		if err != nil {
			return nil
		}

		taggedCoswid := append(CoswidTag, coswidCBOR...)

		o.Tags = append(o.Tags, taggedCoswid)
	}
	return o
}

// AddDependentRim creates a corim-locator-map from the supplied arguments and
// appends it to the dependent RIMs in the unsigned-corim-map
func (o *UnsignedCorim) AddDependentRim(href string, thumbprint *swid.HashEntry) *UnsignedCorim {
	if o != nil {
		l := Locator{
			Href:       comid.TaggedURI(href),
			Thumbprint: thumbprint,
		}

		if o.DependentRims == nil {
			o.DependentRims = new([]Locator)
		}

		*o.DependentRims = append(*o.DependentRims, l)

	}
	return o
}

// AddProfile appends the supplied profile identifier (either a URL or OID) to
// the profiles array in the unsigned-corim-map
func (o *UnsignedCorim) AddProfile(urlOrOID string) *UnsignedCorim {
	if o != nil {
		p, err := eat.NewProfile(urlOrOID)
		if err != nil {
			return nil
		}

		if o.Profiles == nil {
			o.Profiles = new([]eat.Profile)
		}

		*o.Profiles = append(*o.Profiles, *p)
	}
	return o
}

// SetRimValidity can be used to set the validity period of the CoRIM.
// The caller must supply a "not-after" timestamp and optionally a "not-before"
// timestamp.
func (o *UnsignedCorim) SetRimValidity(notAfter time.Time, notBefore *time.Time) *UnsignedCorim {
	if o != nil {
		v := NewValidity().Set(notAfter, notBefore)
		if v == nil {
			return nil
		}

		o.RimValidity = v
	}
	return o
}

// AddEntity adds an organizational entity, together with the roles this entity
// claims with regards to the CoRIM, to the target UnsignerCorim.  name is the entity
// name, regID is a URI that uniquely identifies the entity.  For the moment, roles
// can only be RoleManifestCreator.
func (o *UnsignedCorim) AddEntity(name string, regID *string, roles ...Role) *UnsignedCorim {
	if o != nil {
		e := NewEntity().
			SetEntityName(name).
			SetRoles(roles...)

		if regID != nil {
			e = e.SetRegID(*regID)
		}

		if e == nil {
			return nil
		}

		if o.Entities == nil {
			o.Entities = new(Entities)
		}

		if o.Entities.AddEntity(*e) == nil {
			return nil
		}
	}
	return o
}

// Valid checks the validity (according to the spec) of the target unsigned CoRIM
func (o UnsignedCorim) Valid() error {
	if o.ID == (swid.TagID{}) {
		return fmt.Errorf("empty id")
	}

	if len(o.Tags) == 0 {
		return errors.New("tags validation failed: no tags")
	}

	for i, t := range o.Tags {
		if err := t.Valid(); err != nil {
			return fmt.Errorf("tag validation failed at pos %d: %w", i, err)
		}
	}

	if o.DependentRims != nil {
		for i, r := range *o.DependentRims {
			if err := r.Valid(); err != nil {
				return fmt.Errorf("dependent RIM validation failed at pos %d: %w", i, err)
			}
		}
	}

	if o.Profiles != nil {
		for i, p := range *o.Profiles {
			if err := ValidProfile(p); err != nil {
				return fmt.Errorf("profile validation failed at pos %d: %w", i, err)
			}
		}
	}

	if o.RimValidity != nil {
		if err := o.RimValidity.Valid(); err != nil {
			return fmt.Errorf("RIM validity validation failed: %w", err)
		}
	}

	if o.Entities != nil {
		for i, e := range *o.Entities {
			if err := e.Valid(); err != nil {
				return fmt.Errorf("entity validation failed at pos %d: %w", i, err)
			}
		}
	}

	return nil
}

// ToCBOR serializes the target unsigned CoRIM to CBOR
func (o UnsignedCorim) ToCBOR() ([]byte, error) {
	return em.Marshal(&o)
}

// FromCBOR deserializes a CBOR-encoded unsigned CoRIM into the target UnsignedCorim
func (o *UnsignedCorim) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, o)
}

// FromJSON deserializes a JSON-encoded unsigned CoRIM into the target UnsignedCorim
func (o *UnsignedCorim) FromJSON(data []byte) error {
	return json.Unmarshal(data, o)
}

// Tag is either a CBOR-encoded CoMID, CoSWID or CoTS
type Tag []byte

func (o Tag) Valid() error {
	// there is no much we can check here, except making sure that the tag is
	// not zero-length
	if len(o) == 0 {
		return errors.New("empty tag")
	}
	return nil
}

// Locator is the internal representation of the corim-locator-map with CBOR and
// JSON serialization.
type Locator struct {
	Href       comid.TaggedURI `cbor:"0,keyasint" json:"href"`
	Thumbprint *swid.HashEntry `cbor:"1,keyasint,omitempty" json:"thumbprint,omitempty"`
}

func (o Locator) Valid() error {
	if o.Href.Empty() {
		return errors.New("empty href")
	}

	if tp := o.Thumbprint; tp != nil {
		if err := swid.ValidHashEntry(tp.HashAlgID, tp.HashValue); err != nil {
			return fmt.Errorf("invalid locator thumbprint: %w", err)
		}
	}

	return nil
}

// ValidProfile checks that the supplied profile is in one of the supported
// formats (i.e., URI or OID)
func ValidProfile(p eat.Profile) error {
	if !p.IsOID() && !p.IsURI() {
		return errors.New("profile should be OID or URI")
	}
	return nil
}
