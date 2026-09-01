// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

// Package cotl implements the Concise Tag List (CoTL) data structure defined
// in Section 6 of draft-ietf-rats-corim-11, including CBOR (tag 508) and JSON
// serialization, validation, temporal appraisal checks, and Veraison-style
// map extensions on the concise-tl-tag extension point.
package cotl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/veraison/corim/corim"
	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/swid"
)

// CotlTag is the CBOR tag number for a tagged-concise-tl-tag
// (draft-ietf-rats-corim-11, Section 6).
const CotlTag uint64 = 508

var (
	em cbor.EncMode
	dm cbor.DecMode
)

func init() {
	encOpt := cbor.EncOptions{
		Sort:          cbor.SortCoreDeterministic,
		IndefLength:   cbor.IndefLengthForbidden,
		NilContainers: cbor.NilContainerAsEmpty,
		TimeTag:       cbor.EncTagRequired,
	}

	m, err := encOpt.EncMode()
	if err != nil {
		panic(err)
	}
	em = m

	d, err := cbor.DecOptions{
		IndefLength: cbor.IndefLengthAllowed,
		TimeTag:     cbor.DecTagOptional,
	}.DecMode()
	if err != nil {
		panic(err)
	}
	dm = d
}

// TagIdentityMap identifies a CoMID or CoSWID tag, either as the identity of
// the CoTL itself or as an entry in its tags-list.
type TagIdentityMap struct {
	TagID      swid.TagID `cbor:"0,keyasint" json:"tag-id"`
	TagVersion *uint      `cbor:"1,keyasint,omitempty" json:"tag-version,omitempty"`
}

// NewTagIdentityMap returns a TagIdentityMap for the supplied identifier,
// which may be a string (text or UUID form), UUID bytes, or a uuid.UUID.
func NewTagIdentityMap(v interface{}) (*TagIdentityMap, error) {
	tagID := swid.NewTagID(v)
	if tagID == nil {
		return nil, fmt.Errorf("invalid tag-id: %T", v)
	}
	return &TagIdentityMap{TagID: *tagID}, nil
}

// Valid returns nil if the tag identity carries a usable identifier.
func (o TagIdentityMap) Valid() error {
	if o.TagID == (swid.TagID{}) {
		return errors.New("empty tag-id")
	}
	return nil
}

// String renders the identifier in human-readable form.
func (o TagIdentityMap) String() string {
	s := o.TagID.String()
	if o.TagVersion != nil {
		return fmt.Sprintf("%s (version %d)", s, *o.TagVersion)
	}
	return s
}

// ConciseTlTag is the concise-tl-tag structure from draft-ietf-rats-corim-11,
// Section 6.1:
//
//	concise-tl-tag = {
//	  &(tag-identity: 0) => tag-identity-map
//	  &(tags-list: 1) => [ + tag-identity-map ],
//	  &(tl-validity: 2) => validity-map
//	}
//
// All tags listed in TagsList are activated atomically: either every tag is
// available to the Verifier, or the entire CoTL is rejected (Section 6).
type ConciseTlTag struct {
	TagIdentity TagIdentityMap   `cbor:"0,keyasint" json:"tag-identity"`
	TagsList    []TagIdentityMap `cbor:"1,keyasint" json:"tags-list"`
	TlValidity  corim.Validity   `cbor:"2,keyasint" json:"tl-validity"`

	extensions.Extensions
}

// Valid performs structural validation per Section 6.1: the tags-list must be
// non-empty (activation is atomic), every entry must have a valid identity,
// and the validity window must not be inverted.
func (o *ConciseTlTag) Valid() error {
	if err := o.TagIdentity.Valid(); err != nil {
		return fmt.Errorf("invalid tag-identity: %w", err)
	}

	if len(o.TagsList) == 0 {
		return ErrEmptyTagsList
	}

	for i, t := range o.TagsList {
		if err := t.Valid(); err != nil {
			return fmt.Errorf("invalid tags-list entry %d: %w", i, err)
		}
	}

	if err := o.TlValidity.Valid(); err != nil {
		return fmt.Errorf("invalid tl-validity: %w", err)
	}

	if err := o.validExtensions(); err != nil {
		return err
	}

	return nil
}

// ToCBOR serializes the CoTL as deterministic CBOR wrapped in CBOR tag 508.
func (o *ConciseTlTag) ToCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	buf, err := encoding.SerializeStructToCBOR(em, o)
	if err != nil {
		return nil, fmt.Errorf("CBOR encoding failed: %w", err)
	}

	var out bytes.Buffer
	if err := cbor.NewEncoder(&out).Encode(cbor.Tag{Number: CotlTag, Content: buf}); err != nil {
		return nil, fmt.Errorf("tag wrapping failed: %w", err)
	}

	return out.Bytes(), nil
}

// FromCBOR parses CBOR bytes into a ConciseTlTag. Both the tagged (#6.508)
// and untagged forms are accepted for interoperability with processors that
// convey the bare map.
func (o *ConciseTlTag) FromCBOR(buf []byte) error {
	content, err := unwrapCotlTag(buf)
	if err != nil {
		return err
	}

	if err := populateSafely(o, content); err != nil {
		return err
	}

	if err := o.Valid(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// unwrapCotlTag verifies that buf carries CBOR tag 508 (if any) and returns
// the enclosed concise-tl-tag map bytes.
func unwrapCotlTag(buf []byte) ([]byte, error) {
	var tag cbor.Tag
	if err := dm.Unmarshal(buf, &tag); err != nil {
		// Not a tagged message: assume bare map.
		return buf, nil
	}

	if tag.Number != CotlTag {
		return nil, fmt.Errorf("expected CBOR tag %d, got %d", CotlTag, tag.Number)
	}

	content, ok := tag.Content.([]byte)
	if !ok {
		return nil, errors.New("tag content is not a byte string")
	}

	return content, nil
}

// populateSafely wraps CBOR struct decoding so that panics on malformed
// input are surfaced as errors instead of crashing the calling process.
func populateSafely(o *ConciseTlTag, content []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("CBOR decoding failed: %v", r)
		}
	}()

	// Gate on RFC 8949 well-formedness first: the corim encoding helpers do
	// not defend against malformed map headers (they can loop indefinitely).
	if err := dm.Wellformed(content); err != nil {
		return fmt.Errorf("malformed CBOR: %w", err)
	}

	// Note: the corim encoding helpers (rather than plain Unmarshal) are
	// deliberate here so that any extension structs previously passed to
	// RegisterExtensions() are resolved through their embedded interface
	// and populated from the wire form.
	if err := encoding.PopulateStructFromCBOR(dm, content, o); err != nil {
		return fmt.Errorf("CBOR decoding failed: %w", err)
	}

	return nil
}

// NewConciseTlTagFromCBOR parses tagged or untagged CoTL CBOR bytes.
func NewConciseTlTagFromCBOR(buf []byte) (*ConciseTlTag, error) {
	cotl := new(ConciseTlTag)
	if err := cotl.FromCBOR(buf); err != nil {
		return nil, err
	}

	return cotl, nil
}

// ToJSON serializes the CoTL to human-readable JSON.
func (o *ConciseTlTag) ToJSON() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	buf, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("JSON encoding failed: %w", err)
	}

	return buf, nil
}

// FromJSON parses human-readable JSON into a ConciseTlTag.
func (o *ConciseTlTag) FromJSON(buf []byte) error {
	if err := json.Unmarshal(buf, o); err != nil {
		return fmt.Errorf("JSON decoding failed: %w", err)
	}

	if err := o.Valid(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
