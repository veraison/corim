// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
)

var PSARefValIDType = "psa.refval-id"

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

func CreatePSARefValID(signerID []byte, label, version string) (*PSARefValID, error) {
	ret, err := NewPSARefValID(signerID)
	if err != nil {
		return nil, err
	}

	ret.SetLabel(label)
	ret.SetVersion(version)

	return ret, nil
}

func MustCreatePSARefValID(signerID []byte, label, version string) *PSARefValID {
	ret, err := CreatePSARefValID(signerID, label, version)

	if err != nil {
		panic(err)
	}

	return ret
}

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

func (o *PSARefValID) SetLabel(label string) *PSARefValID {
	if o != nil {
		o.Label = &label
	}
	return o
}

func (o *PSARefValID) SetVersion(version string) *PSARefValID {
	if o != nil {
		o.Version = &version
	}
	return o
}

type TaggedPSARefValID PSARefValID

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

func (o TaggedPSARefValID) Valid() error {
	return PSARefValID(o).Valid()
}

func (o TaggedPSARefValID) String() string {
	ret, err := json.Marshal(o)
	if err != nil {
		return ""
	}

	return string(ret)
}

func (o TaggedPSARefValID) Type() string {
	return PSARefValIDType
}

func (o TaggedPSARefValID) IsZero() bool {
	return len(o.SignerID) == 0
}
