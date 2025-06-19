// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"errors"
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
)

type EvTriples struct {
	EvidenceTriples   *comid.ValueTriples `cbor:"0,keyasint,omitempty" json:"evidence-triples,omitempty"`
	IdentityTriples   *comid.KeyTriples   `cbor:"1,keyasint,omitempty" json:"identity-triples,omitempty"`
	CoSWIDTriples     *CoSWIDTriples      `cbor:"4,keyasint,omitempty" json:"coswid-triples,omitempty"`
	AttestKeysTriples *comid.KeyTriples   `cbor:"5,keyasint,omitempty" json:"attestkey-triples,omitempty"`
	Extensions
}

func NewEvTriples() *EvTriples {
	return &EvTriples{}
}

func (o EvTriples) Valid() error {

	// Check if triples are set ?
	if o.EvidenceTriples == nil &&
		o.IdentityTriples == nil &&
		o.CoSWIDTriples == nil &&
		o.AttestKeysTriples == nil {
		return errors.New("no Triples set inside EvTriples")
	}

	if o.EvidenceTriples != nil {
		if err := o.EvidenceTriples.Valid(); err != nil {
			return fmt.Errorf("invalid EvidenceTriples: %w", err)
		}
	}

	if o.IdentityTriples != nil {
		for i, identity := range *o.IdentityTriples {
			if err := identity.Valid(); err != nil {
				return fmt.Errorf("invalid IdentityTriples at index: %d, %w", i, err)
			}

		}
	}

	if o.CoSWIDTriples != nil {
		for i, swid := range *o.CoSWIDTriples {
			if err := swid.Valid(); err != nil {
				return fmt.Errorf("invalid CoSWIDTriples at index: %d, %w", i, err)
			}

		}
	}

	if o.AttestKeysTriples != nil {
		for i, key := range *o.AttestKeysTriples {
			if err := key.Valid(); err != nil {
				return fmt.Errorf("invalid AttestKeysTriple at index: %d, %w", i, err)
			}
		}
	}

	return nil
}

func (o *EvTriples) AddEvidenceTriple(val *comid.ValueTriple) *EvTriples {
	if o != nil {
		if o.EvidenceTriples == nil {
			o.EvidenceTriples = comid.NewValueTriples()
		}
		o.EvidenceTriples.Add(val)
	}

	return o
}

func (o *EvTriples) AddCoSWIDTriple(val *CoSWIDTriple) *EvTriples {
	if o != nil {
		if o.CoSWIDTriples == nil {
			o.CoSWIDTriples = NewCoSWIDTriples()
		}
		*o.CoSWIDTriples = append(*o.CoSWIDTriples, *val)
	}
	return o
}

func (o *EvTriples) AddIdentityTriple(val *comid.KeyTriple) *EvTriples {
	if o != nil {
		if o.IdentityTriples == nil {
			o.IdentityTriples = comid.NewKeyTriples()
		}
		*o.IdentityTriples = append(*o.IdentityTriples, *val)
	}

	return o
}

func (o *EvTriples) AddAttestKeyTriple(val *comid.KeyTriple) *EvTriples {
	if o != nil {
		if o.AttestKeysTriples == nil {
			o.AttestKeysTriples = comid.NewKeyTriples()
		}
		*o.AttestKeysTriples = append(*o.AttestKeysTriples, *val)
	}

	return o
}

func (o *EvTriples) RegisterExtensions(exts extensions.Map) error {
	EvidenceTriplesExts := extensions.NewMap()
	for p, v := range exts {
		switch p {
		case ExtEvTriples:
			o.Register(v)
		case ExtEvidenceTriples:
			if o.EvidenceTriples == nil {
				o.EvidenceTriples = comid.NewValueTriples()
			}
			EvidenceTriplesExts[comid.ExtMval] = v
		case ExtEvidenceTriplesFlags:
			if o.EvidenceTriples == nil {
				o.EvidenceTriples = comid.NewValueTriples()
			}
			EvidenceTriplesExts[comid.ExtFlags] = v
		default:
			return fmt.Errorf("%w: %q", extensions.ErrUnexpectedPoint, p)
		}
	}
	if len(EvidenceTriplesExts) != 0 {
		return o.EvidenceTriples.RegisterExtensions(EvidenceTriplesExts)
	}
	return nil
}

// GetExtensions returns previously registered extension
func (o *EvTriples) GetExtensions() extensions.IMapValue {
	return o.IMapValue
}

// UnmarshalCBOR deserializes from CBOR
func (o *EvTriples) UnmarshalCBOR(data []byte) error {
	return encoding.PopulateStructFromCBOR(dm, data, o)
}

// MarshalCBOR serializes to CBOR
func (o EvTriples) MarshalCBOR() ([]byte, error) {
	// If extensions have been registered, the collection will exist, but
	// might be empty. If that is the case, set it to nil to avoid
	// marshaling an empty list (and let the marshaller omit the claim
	// instead). Note that since the receiver was passed by value, we do not
	// need to worry about saving the field's value before setting it to
	// nil.
	if o.EvidenceTriples != nil && o.EvidenceTriples.IsEmpty() {
		o.EvidenceTriples = nil
	}
	// as of now there are no further extensions in EvTriples
	return encoding.SerializeStructToCBOR(em, o)
}

// UnmarshalJSON deserializes from JSON
func (o *EvTriples) UnmarshalJSON(data []byte) error {
	return encoding.PopulateStructFromJSON(data, o)
}

// MarshalJSON serializes to JSON
func (o EvTriples) MarshalJSON() ([]byte, error) {
	// If extensions have been registered, the collection will exist, but
	// might be empty. If that is the case, set it to nil to avoid
	// marshaling an empty list (and let the marshaller omit the claim
	// instead). Note that since the receiver was passed by value, we do not
	// need to worry about saving the field's value before setting it to
	// nil.
	if o.EvidenceTriples != nil && o.EvidenceTriples.IsEmpty() {
		o.EvidenceTriples = nil
	}

	return encoding.SerializeStructToJSON(o)
}
