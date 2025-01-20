// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/swid"
)

// Comid is the top-level representation of a Concise Module IDentifier with
// CBOR and JSON serialization.
type Comid struct {
	Language    *string     `cbor:"0,keyasint,omitempty" json:"lang,omitempty"`
	TagIdentity TagIdentity `cbor:"1,keyasint" json:"tag-identity"`
	Entities    *Entities   `cbor:"2,keyasint,omitempty" json:"entities,omitempty"`
	LinkedTags  *LinkedTags `cbor:"3,keyasint,omitempty" json:"linked-tags,omitempty"`
	Triples     Triples     `cbor:"4,keyasint" json:"triples"`

	Extensions
}

// NewComid instantiates an empty Comid
func NewComid() *Comid {
	return &Comid{}
}

// RegisterExtensions registers a struct as a collections of extensions
func (o *Comid) RegisterExtensions(exts extensions.Map) error {
	triplesExts := extensions.NewMap()

	for p, v := range exts {
		switch p {
		case ExtComid:
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
			triplesExts.Add(p, v)
		}
	}

	return o.Triples.RegisterExtensions(triplesExts)
}

// GetExtensions returns previously registered extension
func (o *Comid) GetExtensions() extensions.IMapValue {
	return o.Extensions.IMapValue
}

// SetLanguage sets the language used in the target Comid to the supplied
// language tag.  See also: BCP 47 and the IANA Language subtag registry.
func (o *Comid) SetLanguage(language string) *Comid {
	if o != nil {
		if language == "" {
			return nil
		}
		o.Language = &language
	}
	return o
}

// SetTagIdentity sets the identifier of the target Comid to the supplied tagID,
// which MUST be of type string or [16]byte.  A tagIDVersion must also be
// supplied to disambiguate between different revisions of the same tag
// identity.  If the tagID is newly minted, use 0.  If the tagID has already
// been associated with a CoMID, pick a tagIDVersion greater than any other
// existing tagIDVersion's associated with that tagID.
func (o *Comid) SetTagIdentity(tagID interface{}, tagIDVersion uint) *Comid {
	if o != nil {
		id := swid.NewTagID(tagID)
		if id == nil {
			return nil
		}
		o.TagIdentity.TagID = *id
		o.TagIdentity.TagVersion = tagIDVersion
	}
	return o
}

func IsAbsoluteURI(s string) error {
	var (
		u   *url.URL
		err error
	)

	if u, err = url.Parse(s); err != nil {
		return fmt.Errorf("%q failed to parse as URI: %w", s, err)
	}

	if !u.IsAbs() {
		return fmt.Errorf("%q is not an absolute URI", s)
	}

	return nil
}

func String2URI(s *string) (*TaggedURI, error) {
	if s == nil {
		return nil, nil
	}

	if err := IsAbsoluteURI(*s); err != nil {
		return nil, fmt.Errorf("expecting an absolute URI: %w", err)
	}

	v := TaggedURI(*s)

	return &v, nil

}

// AddEntity adds an organizational entity, together with the roles this entity
// claims with regards to the CoMID, to the target Comid.  name is the entity
// name, regID is a URI that uniquely identifies the entity, and roles are one
// or more claimed roles chosen from the following: RoleTagCreator, RoleCreator
// and RoleMaintainer.
func (o *Comid) AddEntity(name string, regID *string, roles ...Role) *Comid {
	if o != nil {
		var rs Roles
		if rs.Add(roles...) == nil {
			return nil
		}

		uri, err := String2URI(regID)
		if err != nil {
			return nil
		}

		e := Entity{
			Name:  MustNewStringEntityName(name),
			RegID: uri,
			Roles: rs,
		}

		if o.Entities == nil {
			o.Entities = new(Entities)
		}

		if o.Entities.Add(&e) == nil {
			return nil
		}
	}
	return o
}

// AddLinkedTag adds a link relationship of type rel between the target Comid
// and another CoMID identified by its tagID.  The rel parameter can be one of
// RelSupplements or RelReplaces.
func (o *Comid) AddLinkedTag(tagID interface{}, rel Rel) *Comid {
	if o != nil {
		id := swid.NewTagID(tagID)
		if id == nil {
			return nil
		}

		lt := LinkedTag{
			LinkedTagID: *id,
			Rel:         rel,
		}

		if o.LinkedTags == nil {
			o.LinkedTags = new(LinkedTags)
		}

		if o.LinkedTags.AddLinkedTag(lt) == nil {
			return nil
		}
	}
	return o
}

// AddReferenceValue adds the supplied reference value to the
// reference-triples list of the target Comid.
func (o *Comid) AddReferenceValue(val ValueTriple) *Comid {
	if o != nil {
		if o.Triples.ReferenceValues == nil {
			o.Triples.ReferenceValues = NewValueTriples()
		}

		if o.Triples.AddReferenceValue(val) == nil {
			return nil
		}
	}
	return o
}

// AddEndorsedValue adds the supplied endorsed value to the
// endorsed-triples list of the target Comid.
func (o *Comid) AddEndorsedValue(val ValueTriple) *Comid {
	if o != nil {
		if o.Triples.EndorsedValues == nil {
			o.Triples.EndorsedValues = NewValueTriples()
		}

		if o.Triples.AddEndorsedValue(val) == nil {
			return nil
		}
	}
	return o
}

// AddAttestVerifKey adds the supplied verification key to the
// attest-key-triples list of the target Comid.
func (o *Comid) AddAttestVerifKey(val KeyTriple) *Comid {
	if o != nil {
		if o.Triples.AttestVerifKeys == nil {
			o.Triples.AttestVerifKeys = NewKeyTriples()
		}

		if o.Triples.AddAttestVerifKey(val) == nil {
			return nil
		}
	}
	return o
}

// AddDevIdentityKey adds the supplied identity key to the
// identity-triples list of the target Comid.
func (o *Comid) AddDevIdentityKey(val KeyTriple) *Comid {
	if o != nil {
		if o.Triples.DevIdentityKeys == nil {
			o.Triples.DevIdentityKeys = NewKeyTriples()
		}

		if o.Triples.AddDevIdentityKey(val) == nil {
			return nil
		}
	}
	return o
}

// nolint:gocritic
func (o Comid) Valid() error {
    if err := o.TagIdentity.Valid(); err != nil {
        return fmt.Errorf("tag-identity validation failed: %v", err) // Changed %w to %v
    }

    if o.Entities != nil {
        if err := o.Entities.Valid(); err != nil {
            return fmt.Errorf("entities validation failed: %v", err) // Changed %w to %v
        }
    }

    if o.LinkedTags != nil {
        if err := o.LinkedTags.Valid(); err != nil {
            return fmt.Errorf("linked-tags validation failed: %v", err) // Changed %w to %v
        }
    }

    if err := o.Triples.Valid(); err != nil {
        return fmt.Errorf("triples validation failed: %v", err) // Changed %w to %v
    }

    return o.Extensions.validComid(&o)
}

// ToCBOR serializes the target Comid to CBOR
// nolint:gocritic
func (o Comid) ToCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	// If extensions have been registered, the collection will exist, but
	// might be empty. If that is the case, set it to nil to avoid
	// marshaling an empty list (and let the marshaller omit the claim
	// instead). Note that since the receiver was passed by value, we do not
	// need to worry about saving the field's value before setting it to
	// nil.
	if o.Entities != nil && o.Entities.IsEmpty() {
		o.Entities = nil
	}

	return encoding.SerializeStructToCBOR(em, &o)
}

// FromCBOR deserializes a CBOR-encoded CoMID into the target Comid
func (o *Comid) FromCBOR(data []byte) error {
	return encoding.PopulateStructFromCBOR(dm, data, o)
}

// ToJSON serializes the target Comid to JSON
// nolint:gocritic
func (o Comid) ToJSON() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	// If extensions have been registered, the collection will exist, but
	// might be empty. If that is the case, set it to nil to avoid
	// marshaling an empty list (and let the marshaller omit the claim
	// instead). Note that since the receiver was passed by value, we do not
	// need to worry about saving the field's value before setting it to
	// nil.
	if o.Entities != nil && o.Entities.IsEmpty() {
		o.Entities = nil
	}

	return encoding.SerializeStructToJSON(&o)
}

// FromJSON deserializes a JSON-encoded CoMID into the target Comid
func (o *Comid) FromJSON(data []byte) error {
	return encoding.PopulateStructFromJSON(data, o)
}

// nolint:gocritic
func (o Comid) ToJSONPretty(indent string) ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return json.MarshalIndent(&o, "", indent)
}

// AddSimpleReferenceValue adds a reference value with a single measurement
func (o *Comid) AddSimpleReferenceValue(env Environment, measurement *Measurement) error {
    if err := env.Valid(); err != nil {
        return fmt.Errorf("invalid environment: %w", err)
    }
    
    if measurement == nil {
        return fmt.Errorf("measurement cannot be nil")
    }

    if o.Triples.ReferenceValues == nil {
        o.Triples.ReferenceValues = NewValueTriples()
    }
    
    builder := NewReferenceValueBuilder().
        WithEnvironment(env).
        WithMeasurement(measurement)
        
    triple, err := builder.Build()
    if err != nil {
        return fmt.Errorf("building reference value: %w", err)
    }
    
    if res := o.AddReferenceValue(*triple); res == nil {
		return fmt.Errorf("failed to add reference value")
	}
    
    return nil
}

func (o *Comid) AddDigestReferenceValue(env Environment, alg string, digest []byte) error {
    if len(digest) == 0 {
        return fmt.Errorf("digest cannot be empty")
    }
    hashAlg := HashAlgFromString(alg)
    if !hashAlg.Valid() {
        return fmt.Errorf("unrecognized algorithm %q", alg)
    }
    m := &Measurement{
        Val: Mval{
            Digests: NewDigests(),
        },
    }
    if m.Val.Digests.AddDigest(hashAlg.ToUint64(), digest) == nil {
        return fmt.Errorf("failed to create hash entry")
    }
    return o.AddSimpleReferenceValue(env, m)
}

// AddRawReferenceValue adds a reference value with raw measurement data
func (o *Comid) AddRawReferenceValue(env Environment, raw []byte) error {
    if len(raw) == 0 {
        return fmt.Errorf("raw value cannot be empty")
    }
    
    m := &Measurement{
        Val: Mval{
            RawValue: NewRawValue().SetBytes(raw),
        },
    }
    
    return o.AddSimpleReferenceValue(env, m)
}

// AddReferenceValueDirect adds a reference value directly to the reference-triples list without creating instances for Measurement and ValueTriples.
func (o *Comid) AddReferenceValueDirect(environment Environment, measurements Measurements) *Comid {
    if o != nil {
        val := ValueTriple{
            Environment:  environment,
            Measurements: measurements,
        }
        if o.Triples.ReferenceValues == nil {
            o.Triples.ReferenceValues = NewValueTriples()
        }

        if o.Triples.AddReferenceValue(val) == nil {
            return nil
        }
    }
    return o
}

// AddEndorsedValueDirect adds an endorsed value directly to the endorsed-triples list without creating instances for Measurement and ValueTriples.
func (o *Comid) AddEndorsedValueDirect(environment Environment, measurements Measurements) *Comid {
    if o != nil {
        val := ValueTriple{
            Environment:  environment,
            Measurements: measurements,
        }
        if o.Triples.EndorsedValues == nil {
            o.Triples.EndorsedValues = NewValueTriples()
        }

		if o.Triples.AddEndorsedValue(val) == nil {
            return nil
        }
    }
    return o
}