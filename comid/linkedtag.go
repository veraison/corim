// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"

	"github.com/veraison/swid"
)

// LinkedTag stores one link relation of type Rel between the embedding CoMID
// (the link context) and the referenced CoMID (the link target). The link can
// be viewed as a statement of the form: "$link_context $link_relation_type
// $link_target".
type LinkedTag struct {
	LinkedTagID swid.TagID `cbor:"0,keyasint" json:"target"`
	Rel         Rel        `cbor:"1,keyasint" json:"rel"`
}

func NewLinkedTag() *LinkedTag {
	return &LinkedTag{
		Rel: *NewRel(),
		// LinkedTagID default constructed
	}
}

func (o *LinkedTag) SetLinkedTag(t swid.TagID) *LinkedTag {
	if o != nil {
		o.LinkedTagID = t
	}
	return o
}

func (o *LinkedTag) SetRel(r Rel) *LinkedTag {
	if o != nil {
		o.Rel.Set(r)
	}
	return o
}

func (o LinkedTag) Valid() error {
	if o.LinkedTagID == (swid.TagID{}) {
		return fmt.Errorf("tag-id must be set in linked-tag")
	}

	if err := o.Rel.Valid(); err != nil {
		return fmt.Errorf("rel validation failed: %w", err)
	}

	return nil
}

// LinkedTags is an array of LinkedTag
type LinkedTags []LinkedTag

func NewLinkedTags() *LinkedTags {
	return new(LinkedTags)
}

// AddLinkedTag adds the supplied linked Tag-map to the target Entities
func (o *LinkedTags) AddLinkedTag(lt LinkedTag) *LinkedTags {
	if o != nil {
		*o = append(*o, lt)
	}
	return o
}

func (o LinkedTags) Valid() error {
	for i, l := range o {
		if err := l.Valid(); err != nil {
			return fmt.Errorf("invalid linked-tag entry at index %d: %w", i, err)
		}
	}

	return nil
}
