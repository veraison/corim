// Copyright 2023-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package cots

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/veraison/swid"
)

// This is copied from the SoftwareIdentity implementation in github.com/veraison/swid with most fields made
// optional here.
type AbbreviatedSwidTag struct {
	XMLName xml.Name `cbor:"-" json:"-"`

	swid.CoSWIDExtension

	swid.GlobalAttributes

	// A 16 byte binary string or textual identifier uniquely referencing a
	// software component. The tag identifier MUST be globally unique. If
	// represented as a 16 byte binary string, the identifier MUST be a valid
	// universally unique identifier as defined by [RFC4122]. There are no
	// strict guidelines on how this identifier is structured, but examples
	// include a 16 byte GUID (e.g. class 4 UUID) [RFC4122], or a text string
	// appended to a DNS domain name to ensure uniqueness across organizations.
	TagID *swid.TagID `cbor:"0,keyasint,omitempty" json:"tag-id,omitempty" xml:"tagId,attr,omitempty"`

	// An integer value that indicate the specific release revision of the tag.
	// Typically, the initial value of this field is set to 0 and the value is
	// monotonically increased for subsequent tags produced for the same
	// software component release. This value allows a CoSWID tag producer to
	// correct an incorrect tag previously released without indicating a change
	// to the underlying software component the tag represents. For example, the
	// tag version could be changed to add new metadata, to correct a broken
	// link, to add a missing payload entry, etc. When producing a revised tag,
	// the new tag-version value MUST be greater than the old tag-version value.
	TagVersion int `cbor:"12,keyasint,omitempty" json:"tag-version,omitempty" xml:"tagVersion,attr,omitempty"`

	// A boolean value that indicates if the tag identifies and describes an
	// installable software component in its pre-installation state. Installable
	// software includes a installation package or installer for a software
	// component, a software update, or a patch. If the CoSWID tag represents
	// installable software, the corpus item MUST be set to "true". If not
	// provided, the default value MUST be considered "false"
	Corpus bool `cbor:"8,keyasint,omitempty" json:"corpus,omitempty" xml:"corpus,attr,omitempty"`

	// A boolean value that indicates if the tag identifies and describes an
	// installed patch that has made incremental changes to a software component
	// installed on an endpoint. Typically, an installed patch has made a set of
	// file modifications to pre-installed software and does not alter the
	// version number or the descriptive metadata of an installed software
	// component. If a CoSWID tag is for a patch, the patch item MUST be set to
	// "true". If not provided, the default value MUST be considered "false".
	//
	// Note: If the software component's version number is modified, then the
	// correct course of action would be to replace the previous primary tag for
	// the component with a new primary tag that reflected this new version. In
	// such a case, the new tag would have a patch item value of "false" or
	// would omit this item completely.
	Patch bool `cbor:"9,keyasint,omitempty" json:"patch,omitempty" xml:"patch,attr,omitempty"`

	// A boolean value that indicates if the tag is providing additional
	// information to be associated with another referenced SWID or CoSWID tag.
	// This allows tools and users to record their own metadata about a software
	// component without modifying SWID primary or patch tags created by a
	// software provider. If a CoSWID tag is a supplemental tag, the
	// supplemental item MUST be set to "true". If not provided, the default
	// value MUST be considered "false".
	Supplemental bool `cbor:"11,keyasint,omitempty" json:"supplemental,omitempty" xml:"supplemental,attr,omitempty"`

	// This textual item provides the software component's name. This name is
	// likely the same name that would appear in a package management tool.
	SoftwareName string `cbor:"1,keyasint,omitempty" json:"software-name,omitempty" xml:"name,attr,omitempty"`

	// A textual value representing the specific release or development version
	// of the software component.
	SoftwareVersion string `cbor:"13,keyasint,omitempty" json:"software-version,omitempty" xml:"version,attr,omitempty"`

	// An integer or textual value representing the versioning scheme used for
	// the software-version item. If an integer value is used it MUST be an
	// index value in the range -256 to 65535. Integer values in the range -256
	// to -1 are reserved for testing and use in closed environments (see
	// section Section 5.2.2). Integer values in the range 0 to 65535 correspond
	// to registered entries in the IANA "SWID/CoSWID Version Scheme Value"
	// registry (see section Section 5.2.4. If a string value is used it MUST be
	// a private use name as defined in section Section 5.2.2. String values
	// based on a Version Scheme Name from the IANA "SWID/CoSWID Version Scheme
	// Value" registry MUST NOT be used, as these values are less concise than
	// their index value equivalent.
	VersionScheme *swid.VersionScheme `cbor:"14,keyasint,omitempty" json:"version-scheme,omitempty" xml:"versionScheme,attr,omitempty"`

	// This text value is a hint to the tag consumer to understand what target
	// platform this tag applies to. This item represents a query as defined by
	// the W3C Media Queries Recommendation (see
	// http://www.w3.org/TR/2012/REC-css3-mediaqueries-20120619)
	Media string `cbor:"10,keyasint,omitempty" json:"media,omitempty" xml:"media,attr,omitempty"`

	// An open-ended map of key/value data pairs. A number of predefined keys
	// can be used within this item providing for common usage and semantics
	// across the industry. Use of this map allows any additional attribute to
	// be included in the tag. It is expected that industry groups will use a
	// common set of attribute names to allow for interoperability within their
	// communities.
	SoftwareMetas *swid.SoftwareMetas `cbor:"5,keyasint,omitempty" json:"software-meta,omitempty" xml:"Meta,omitempty"`

	// Provides information about one or more organizations responsible for
	// producing the CoSWID tag, and producing or releasing the software
	// component referenced by this CoSWID tag.
	Entities swid.Entities `cbor:"2,keyasint" json:"entity" xml:"Entity"`

	// Provides a means to establish relationship arcs between the tag and
	// another items. A given link can be used to establish the relationship
	// between tags or to reference another resource that is related to the
	// CoSWID tag, e.g. vulnerability database association, ROLIE feed
	// [RFC8322], MUD resource [RFC8520], software download location, etc). This
	// is modeled after the HTML "link" element.
	Links *swid.Links `cbor:"4,keyasint,omitempty" json:"link,omitempty" xml:"Link,omitempty"`

	// This item represents a collection of software artifacts (described by
	// child items) that compose the target software. For example, these
	// artifacts could be the files included with an installer for a corpus tag
	// or installed on an endpoint when the software component is installed for
	// a primary or patch tag. The artifacts listed in a payload may be a
	// superset of the software artifacts that are actually installed. Based on
	// user selections at install time, an installation might not include every
	// artifact that could be created or executed on the endpoint when the
	// software component is installed or run.
	Payload *swid.Payload `cbor:"6,keyasint,omitempty" json:"payload,omitempty" xml:"Payload,omitempty"`

	// This item can be used to record the results of a software discovery
	// process used to identify untagged software on an endpoint or to represent
	// indicators for why software is believed to be installed on the endpoint.
	// In either case, a CoSWID tag can be created by the tool performing an
	// analysis of the software components installed on the endpoint.
	Evidence *swid.Evidence `cbor:"3,keyasint,omitempty" json:"evidence,omitempty" xml:"Evidence,omitempty"`
}

// NewTag instantiates a new SWID tag with the supplied tag identifier and
// software name and version
func NewTag(tagID interface{}, softwareName, softwareVersion string) (*AbbreviatedSwidTag, error) {
	t := AbbreviatedSwidTag{
		XMLName: xml.Name{
			Space: "http://standards.iso.org/iso/19770/-2/2015/schema.xsd",
			Local: "AbbreviatedSwidTag",
		},
		SoftwareName:    softwareName,
		SoftwareVersion: softwareVersion,
	}

	if err := t.setTagID(tagID); err != nil {
		return nil, err
	}

	return &t, nil
}

// nolint:gocritic
func (t AbbreviatedSwidTag) Valid() error {
	if len(t.Entities) == 0 || t.Entities == nil {
		return fmt.Errorf("no entities present, must have at least 1 entity")
	}

	// Validate Evidence field if present
	if t.Evidence != nil {
		if err := t.Evidence.Valid(); err != nil {
			return fmt.Errorf("evidence validation failed: %w", err)
		}
	}

	return nil
}

// ToXML serializes the receiver AbbreviatedSwidTag to SWID
// nolint:gocritic
func (t AbbreviatedSwidTag) ToXML() ([]byte, error) {
	return xml.Marshal(t)
}

// ToJSON serializes the receiver AbbreviatedSwidTag to CoSWID using the JSON
// formatter.
// nolint:gocritic
func (t AbbreviatedSwidTag) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

// ToCBOR serializes the receiver AbbreviatedSwidTag to CoSWID
// nolint:gocritic
func (t AbbreviatedSwidTag) ToCBOR() ([]byte, error) {
	return em.Marshal(t)
}

// FromXML deserializes the supplied XML encoded CoSWID into the receiver
// AbbreviatedSwidTag
func (t *AbbreviatedSwidTag) FromXML(data []byte) error {
	return xml.Unmarshal(data, t)
}

// FromJSON deserializes the supplied JSON encoded CoSWID into the receiver
// AbbreviatedSwidTag
func (t *AbbreviatedSwidTag) FromJSON(data []byte) error {
	return json.Unmarshal(data, t)
}

// FromCBOR deserializes the supplied CBOR encoded CoSWID into the receiver
// AbbreviatedSwidTag
func (t *AbbreviatedSwidTag) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, t)
}

func (t *AbbreviatedSwidTag) setTagID(v interface{}) error {
	tagID := swid.NewTagID(v)
	if tagID == nil {
		return errors.New("bad type for TagID: expecting string or [16]byte")
	}

	t.TagID = tagID

	return nil
}

// AddEntity adds the supplied Entity to the receiver AbbreviatedSwidTag
func (t *AbbreviatedSwidTag) AddEntity(e *swid.Entity) error {
	t.Entities = append(t.Entities, *e)

	return nil
}

// AddLink adds the supplied Link to the receiver AbbreviatedSwidTag
func (t *AbbreviatedSwidTag) AddLink(l *swid.Link) error {
	if t.Links == nil {
		t.Links = new(swid.Links)
	}

	*t.Links = append(*t.Links, *l)

	return nil
}

// AddSoftwareMeta adds the supplied SoftwareMeta to the receiver
// AbbreviatedSwidTag
func (t *AbbreviatedSwidTag) AddSoftwareMeta(m *swid.SoftwareMeta) error {
	if t.SoftwareMetas == nil {
		t.SoftwareMetas = new(swid.SoftwareMetas)
	}

	*t.SoftwareMetas = append(*t.SoftwareMetas, *m)

	return nil
}
