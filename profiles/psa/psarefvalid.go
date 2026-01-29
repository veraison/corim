// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"encoding/json"
	"fmt"

	"github.com/veraison/corim/comid"
)

const PSARefValIDType = "psa.refval-id"

// PSARefValID stores a PSA refval-id with CBOR and JSON serializations
// (See https://datatracker.ietf.org/doc/html/draft-xyz-rats-psa-endorsements)
type PSARefValID struct {
	Label    *string `cbor:"1,keyasint,omitempty" json:"label,omitempty"`
	Version  *string `cbor:"4,keyasint,omitempty" json:"version,omitempty"`
	SignerID []byte  `cbor:"5,keyasint" json:"signer-id"` // 32, 48 or 64
}

// Valid checks the validity (according to the spec) of the target PSARefValID
func (o PSARefValID) Valid() error {
	if o.SignerID == nil {
		return fmt.Errorf("missing mandatory signer ID")
	}

	switch len(o.SignerID) {
	case 32, 48, 64:
	default:
		return fmt.Errorf("want 32, 48 or 64 bytes, got %d", len(o.SignerID))
	}

	return nil
}

// CreatePSARefValID creates a PSARefValID with signer ID, label, and version
func CreatePSARefValID(signerID []byte, label, version string) (*PSARefValID, error) {
	ret, err := NewPSARefValID(signerID)
	if err != nil {
		return nil, err
	}

	ret.SetLabel(label)
	ret.SetVersion(version)

	return ret, nil
}

// MustCreatePSARefValID is like CreatePSARefValID except it panics on error
func MustCreatePSARefValID(signerID []byte, label, version string) *PSARefValID {
	ret, err := CreatePSARefValID(signerID, label, version)

	if err != nil {
		panic(err)
	}

	return ret
}

// NewPSARefValID creates a new PSARefValID from various input types
func NewPSARefValID(val any) (*PSARefValID, error) {
	var ret PSARefValID

	if val == nil {
		return &ret, nil
	}

	switch t := val.(type) {
	case PSARefValID:
		ret = t
	case *PSARefValID:
		ret = *t
	case string:
		if err := json.Unmarshal([]byte(t), &ret); err != nil {
			return nil, err
		}
	case []byte:
		switch len(t) {
		case 32, 48, 64:
			ret.SignerID = t
		default:
			return nil, fmt.Errorf("invalid PSA RefVal ID length: %d", len(t))
		}
	default:
		return nil, fmt.Errorf("unexpected type for PSA RefVal ID: %T", t)
	}

	return &ret, nil
}

// SetLabel sets the label field
func (o *PSARefValID) SetLabel(label string) *PSARefValID {
	if o != nil {
		o.Label = &label
	}
	return o
}

// SetVersion sets the version field
func (o *PSARefValID) SetVersion(version string) *PSARefValID {
	if o != nil {
		o.Version = &version
	}
	return o
}

// TaggedPSARefValID is the CBOR-tagged version of PSARefValID
type TaggedPSARefValID PSARefValID

// NewTaggedPSARefValID creates a new TaggedPSARefValID
func NewTaggedPSARefValID(val any) (*TaggedPSARefValID, error) {
	var ret TaggedPSARefValID

	switch t := val.(type) {
	case TaggedPSARefValID:
		ret = t
	case *TaggedPSARefValID:
		ret = *t
	default:
		refvalID, err := NewPSARefValID(val)
		if err != nil {
			return nil, err
		}
		ret = TaggedPSARefValID(*refvalID)

	}

	return &ret, nil
}

// Valid validates the TaggedPSARefValID
func (o TaggedPSARefValID) Valid() error {
	return PSARefValID(o).Valid()
}

// String returns the JSON string representation
func (o TaggedPSARefValID) String() string {
	ret, err := json.Marshal(o)
	if err != nil {
		return ""
	}

	return string(ret)
}

// Type returns the type identifier for this measurement key
func (o TaggedPSARefValID) Type() string {
	return PSARefValIDType
}

// IsZero returns true if the PSARefValID is zero-valued
func (o TaggedPSARefValID) IsZero() bool {
	return len(o.SignerID) == 0
}

// UnmarshalJSON deserializes from JSON
func (o *TaggedPSARefValID) UnmarshalJSON(data []byte) error {
	var temp PSARefValID
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	*o = TaggedPSARefValID(temp)
	return nil
}

// MarshalJSON serializes to JSON
func (o TaggedPSARefValID) MarshalJSON() ([]byte, error) {
	return json.Marshal(PSARefValID(o))
}

// NewMkeyPSARefvalID creates a new Mkey with PSA refval-id type
func NewMkeyPSARefvalID(val any) (*comid.Mkey, error) {
	ret, err := NewTaggedPSARefValID(val)
	if err != nil {
		return nil, err
	}

	return &comid.Mkey{
		Value: ret,
	}, nil
}

// NewPSAMeasurement instantiates a new measurement-map with the key set to the
// supplied PSA refval-id. This is a convenience function for PSA profiles.
func NewPSAMeasurement(key any) (*comid.Measurement, error) {
	return comid.NewMeasurement(key, PSARefValIDType)
}

// MustNewPSAMeasurement is like NewPSAMeasurement except it panics on error
func MustNewPSAMeasurement(key any) *comid.Measurement {
	ret, err := NewPSAMeasurement(key)

	if err != nil {
		panic(err)
	}

	return ret
}
