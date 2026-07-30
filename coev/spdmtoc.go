// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	cborfx "github.com/fxamacker/cbor/v2"
	"github.com/veraison/corim/corim"
)

// SpdmTocTag is the CBOR tag 570 prefix for a tagged-spdm-toc.
var SpdmTocTag = []byte{0xd9, 0x02, 0x3a}

// SpdmToc is spdm-toc-map
//
// The mandatory field TaggedEvidence (key 0) holds one or more
// tagged-concise-evidence items (CBOR Tag 571).
// The optional RimLocators field (key 1) carries one or more corim-locator-map
// entries pointing to external CoRIMs.
// The optional Profile field (key 2) identifies the CoEV profile in use.
type SpdmToc struct {
	TaggedEvidence []TaggedConciseEvidence `cbor:"0,keyasint" json:"tagged-evidence"`
	RimLocators    *[]corim.Locator        `cbor:"1,keyasint,omitempty" json:"rim-locators,omitempty"`
	Profile        *corim.Profile          `cbor:"2,keyasint,omitempty" json:"profile,omitempty"`
}

// UnmarshalCBOR implements cbor.Unmarshaler. Key 0 is decoded into raw
// messages so that SPDM extensions can be pre-wired on each
// TaggedConciseEvidence before it is decoded (extensions must be registered
// before PopulateStructFromCBOR runs.)
func (o *SpdmToc) UnmarshalCBOR(data []byte) error {
	var wire struct {
		Items       []cborfx.RawMessage `cbor:"0,keyasint"`
		RimLocators *[]corim.Locator    `cbor:"1,keyasint,omitempty"`
		Profile     *corim.Profile      `cbor:"2,keyasint,omitempty"`
	}
	if err := cborfx.Unmarshal(data, &wire); err != nil {
		return fmt.Errorf("decoding spdm-toc-map: %w", err)
	}
	o.RimLocators = wire.RimLocators
	o.Profile = wire.Profile
	o.TaggedEvidence = make([]TaggedConciseEvidence, 0, len(wire.Items))
	for i, raw := range wire.Items {
		// Register the spdm-indirect baseline extension (CBOR key 12) before
		// decoding, so the field is recognized rather than stashed in the cache.
		// Profiles that specialise spdm-toc provide their own extensions
		// on top of SpdmExtensionMap(); this registers only the baseline.
		var tce TaggedConciseEvidence
		if err := (*ConciseEvidence)(&tce).RegisterExtensions(SpdmExtensionMap()); err != nil {
			return fmt.Errorf("tagged-concise-evidence %d: registering extensions: %w", i, err)
		}
		if err := tce.FromCBOR([]byte(raw)); err != nil {
			return fmt.Errorf("tagged-concise-evidence %d: %w", i, err)
		}
		o.TaggedEvidence = append(o.TaggedEvidence, tce)
	}
	return nil
}

// TaggedSpdmToc is a tagged-spdm-toc: #6.570(spdm-toc-map).
type TaggedSpdmToc struct {
	SpdmToc
}

// ToCBOR serializes to #6.570(spdm-toc-map). Each TaggedConciseEvidence item is
// encoded via its MarshalCBOR method, which prepends CBOR Tag 571
// automatically.
func (o TaggedSpdmToc) ToCBOR() ([]byte, error) {
	inner, err := em.Marshal(&o.SpdmToc)
	if err != nil {
		return nil, fmt.Errorf("encoding spdm-toc-map: %w", err)
	}
	return append(SpdmTocTag, inner...), nil
}

// FromCBOR deserializes a #6.570(spdm-toc-map), delegating the inner map
// decode to SpdmToc.UnmarshalCBOR.
func (o *TaggedSpdmToc) FromCBOR(data []byte) error {
	if !bytes.HasPrefix(data, SpdmTocTag) {
		return errors.New("missing spdm-toc tag 570")
	}
	return cborfx.Unmarshal(data[3:], &o.SpdmToc)
}

// ToJSON serializes to { "tagged-evidence": [...], "rim-locators"?: [...], "profile"?: ... }.
func (o TaggedSpdmToc) ToJSON() ([]byte, error) {
	items := make([]json.RawMessage, 0, len(o.TaggedEvidence))
	for i, tce := range o.TaggedEvidence {
		b, err := tce.ToJSON()
		if err != nil {
			return nil, fmt.Errorf("tagged-concise-evidence %d: %w", i, err)
		}
		items = append(items, b)
	}
	return json.Marshal(struct {
		TaggedEvidence []json.RawMessage `json:"tagged-evidence"`
		RimLocators    *[]corim.Locator  `json:"rim-locators,omitempty"`
		Profile        *corim.Profile    `json:"profile,omitempty"`
	}{items, o.RimLocators, o.Profile})
}

// FromJSON deserializes a spdm-toc-map from JSON, pre-wiring SPDM extensions
// on each ConciseEvidence item before decoding.
func (o *TaggedSpdmToc) FromJSON(data []byte) error {
	var wrapper struct {
		TaggedEvidence []json.RawMessage `json:"tagged-evidence"`
		RimLocators    *[]corim.Locator  `json:"rim-locators,omitempty"`
		Profile        *corim.Profile    `json:"profile,omitempty"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return err
	}
	o.RimLocators = wrapper.RimLocators
	o.Profile = wrapper.Profile
	o.TaggedEvidence = make([]TaggedConciseEvidence, 0, len(wrapper.TaggedEvidence))
	for i, raw := range wrapper.TaggedEvidence {
		var tce TaggedConciseEvidence
		if err := (*ConciseEvidence)(&tce).RegisterExtensions(SpdmExtensionMap()); err != nil {
			return fmt.Errorf("tagged-concise-evidence %d: registering extensions: %w", i, err)
		}
		if err := (*ConciseEvidence)(&tce).FromJSON(raw); err != nil {
			return fmt.Errorf("tagged-concise-evidence %d: %w", i, err)
		}
		o.TaggedEvidence = append(o.TaggedEvidence, tce)
	}
	return nil
}
