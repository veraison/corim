// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"errors"
	"fmt"
	"time"

	"github.com/veraison/corim/cots"
	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"

	"github.com/veraison/corim/comid"
	"github.com/veraison/eat"
	"github.com/veraison/swid"
)

// UnsignedCorim is the top-level representation of the unsigned-corim-map with
// CBOR and JSON serialization.
type UnsignedCorim struct {
	ID swid.TagID `cbor:"0,keyasint" json:"corim-id"`
	// note: even though tags are mandatory for CoRIM, we allow omitting
	// them in our JSON templates for cocli (the min template just has
	// corim-id). Since we're never writing JSON (so far), this normally
	// wouldn't matter, however the custom serialization code we use to
	// handle embedded structs relies on the omitempty entry to determine
	// if a field is optional, so we use it during unmarshaling as well as
	// marshaling. Hence omitempty is present for the json tag, but not
	// cbor.
	Tags          []Tag        `cbor:"1,keyasint" json:"tags,omitempty"`
	DependentRims *[]Locator   `cbor:"2,keyasint,omitempty" json:"dependent-rims,omitempty"`
	Profile       *eat.Profile `cbor:"3,keyasint,omitempty" json:"profile,omitempty"`
	RimValidity   *Validity    `cbor:"4,keyasint,omitempty" json:"validity,omitempty"`
	Entities      *Entities    `cbor:"5,keyasint,omitempty" json:"entities,omitempty"`

	Extensions
}

type TaggedUnsignedCorim UnsignedCorim

// NewUnsignedCorim instantiates an empty UnsignedCorim
func NewUnsignedCorim() *UnsignedCorim {
	return &UnsignedCorim{}
}

// RegisterExtensions registers a struct as a collections of extensions
func (o *UnsignedCorim) RegisterExtensions(exts extensions.Map) error {
	for p, v := range exts {
		switch p {
		case ExtUnsignedCorim:
			o.Extensions.Register(v)
		case ExtEntity:
			if o.Entities == nil {
				o.Entities = NewEntities()
			}

			entMap := extensions.NewMap().Add(ExtEntity, v)
			if err := o.Entities.RegisterExtensions(entMap); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%w: %q", extensions.ErrUnexpectedPoint, p)
		}
	}

	return nil
}

// GetExtensions returns pervisouosly registered extension
func (o *UnsignedCorim) GetExtensions() extensions.IMapValue {
	return o.Extensions.IMapValue
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
// nolint:gocritic
func (o UnsignedCorim) GetID() string {
	return o.ID.String()
}

// AddComid appends the CBOR encoded (and appropriately tagged) CoMID to the
// tags array of the unsigned-corim-map
func (o *UnsignedCorim) AddComid(c *comid.Comid) *UnsignedCorim {
	if o != nil {
		if c.Valid() != nil {
			return nil
		}

		comidCBOR, err := c.ToCBOR()
		if err != nil {
			return nil
		}

		taggedComid := append(ComidTag, comidCBOR...) //nolint:gocritic

		o.Tags = append(o.Tags, taggedComid)
	}
	return o
}

// AddCots appends the CBOR encoded (and appropriately tagged) CoTS to the
// tags array of the unsigned-corim-map
func (o *UnsignedCorim) AddCots(c *cots.ConciseTaStore) *UnsignedCorim {
	if o != nil {
		if c.Valid() != nil {
			return nil
		}

		cotsCBOR, err := c.ToCBOR()
		if err != nil {
			return nil
		}

		taggedCots := append(cots.CotsTag, cotsCBOR...) //nolint:gocritic

		o.Tags = append(o.Tags, taggedCots)
	}
	return o
}

// AddCoswid appends the CBOR encoded (and appropriately tagged) CoSWID to the
// tags array of the unsigned-corim-map
func (o *UnsignedCorim) AddCoswid(c *swid.SoftwareIdentity) *UnsignedCorim {
	if o != nil {
		// Currently the swid package doesn't offer an interface
		// for validating the supplied CoSWID, so -- for now --
		// we take any input for granted and pass it to the encoder.
		// See also https://github.com/veraison/swid/issues/23.

		coswidCBOR, err := c.ToCBOR()
		if err != nil {
			return nil
		}

		taggedCoswid := append(CoswidTag, coswidCBOR...) //nolint:gocritic

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

// SetProfile sets the supplied profile identifier (either a URL or OID) as
// the profile in the unsigned-corim-map
func (o *UnsignedCorim) SetProfile(urlOrOID string) *UnsignedCorim {
	if o != nil {
		p, err := eat.NewProfile(urlOrOID)
		if err != nil {
			return nil
		}

		o.Profile = p

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
			SetName(name).
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

		if o.Entities.Add(e) == nil {
			return nil
		}
	}
	return o
}

// Valid checks the validity (according to the spec) of the target unsigned CoRIM
// nolint:gocritic
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

	if o.Profile != nil {
		if err := ValidProfile(*o.Profile); err != nil {
			return fmt.Errorf("profile validation failed: %w", err)
		}
	}

	if o.RimValidity != nil {
		if err := o.RimValidity.Valid(); err != nil {
			return fmt.Errorf("RIM validity validation failed: %w", err)
		}
	}

	if o.Entities != nil {
		for i, e := range o.Entities.Values {
			if err := e.Valid(); err != nil {
				return fmt.Errorf("entity validation failed at pos %d: %w", i, err)
			}
		}
	}

	return o.Extensions.validCorim(&o)
}

// ToCBOR serializes the target unsigned CoRIM to CBOR
// nolint:gocritic
func (o UnsignedCorim) ToCBOR() ([]byte, error) {
	// If extensions have been registered, the collection will exist, but
	// might be empty. If that is the case, set it to nil to avoid
	// marshaling an empty list (and let the marshaller omit the claim
	// instead). Note that since the receiver was passed by value, we do not
	// need to worry about saving the field's value before setting it to
	// nil.
	if o.Entities != nil && o.Entities.IsEmpty() {
		o.Entities = nil
	}

	return encoding.SerializeStructToCBOR(em, o)
}

// FromCBOR deserializes a CBOR-encoded unsigned CoRIM into the target UnsignedCorim
func (o *UnsignedCorim) FromCBOR(data []byte) error {
	return encoding.PopulateStructFromCBOR(dm, data, o)
}

// ToJSON serializes the target unsigned CoRIM to JSON
// nolint:gocritic
func (o UnsignedCorim) ToJSON() ([]byte, error) {
	// If extensions have been registered, the collection will exist, but
	// might be empty. If that is the case, set it to nil to avoid
	// marshaling an empty list (and let the marshaller omit the claim
	// instead). Note that since the receiver was passed by value, we do not
	// need to worry about saving the field's value before setting it to
	// nil.
	if o.Entities != nil && o.Entities.IsEmpty() {
		o.Entities = nil
	}

	return encoding.SerializeStructToJSON(o)
}

// FromJSON deserializes a JSON-encoded unsigned CoRIM into the target UnsignedCorim
func (o *UnsignedCorim) FromJSON(data []byte) error {
	return encoding.PopulateStructFromJSON(data, o)
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
