// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
    "encoding/json"
    "errors"
    "fmt"
)

var CCARefValIDType = "cca.refval-id"

type CCARefValID struct {
    Label    *string `cbor:"1,keyasint,omitempty" json:"label,omitempty"`
    Version  *string `cbor:"4,keyasint,omitempty" json:"version,omitempty"`
    SignerID []byte  `cbor:"5,keyasint" json:"signer-id"` 
}

func (o CCARefValID) Valid() error {
    if o.SignerID == nil {
        return errors.New("missing mandatory signer ID")
    }

    switch len(o.SignerID) {
    case 32, 48, 64:
        // Valid lengths
    default:
        return fmt.Errorf("want 32, 48 or 64 bytes, got %d", len(o.SignerID))
    }

    return nil
}

func CreateCCARefValID(signerID []byte, label, version string) (*CCARefValID, error) {
    ret, err := NewCCARefValID(signerID)
    if err != nil {
        return nil, err
    }

    ret.SetLabel(label)
    ret.SetVersion(version)

    return ret, nil
}

// MustCreateCCARefValID is like CreateCCARefValID but panics on error
func MustCreateCCARefValID(signerID []byte, label, version string) *CCARefValID {
    ret, err := CreateCCARefValID(signerID, label, version)

    if err != nil {
        panic(err)
    }

    return ret
}

func NewCCARefValID(val any) (*CCARefValID, error) {
    var ret CCARefValID

    if val == nil {
        return &ret, nil
    }

    switch t := val.(type) {
    case CCARefValID:
        ret = t
    case *CCARefValID:
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
            return nil, fmt.Errorf("invalid CCA RefVal ID length: %d", len(t))
        }
    default:
        return nil, fmt.Errorf("unexpected type for CCA RefVal ID: %T", t)
    }

    return &ret, nil
}

func (o *CCARefValID) SetLabel(label string) *CCARefValID {
    if o != nil {
        o.Label = &label
    }
    return o
}

func (o *CCARefValID) SetVersion(version string) *CCARefValID {
    if o != nil {
        o.Version = &version
    }
    return o
}

type TaggedCCARefValID CCARefValID

func NewTaggedCCARefValID(val any) (*TaggedCCARefValID, error) {
    var ret TaggedCCARefValID

    switch t := val.(type) {
    case TaggedCCARefValID:
        ret = t
    case *TaggedCCARefValID:
        ret = *t
    default:
        refvalID, err := NewCCARefValID(val)
        if err != nil {
            return nil, err
        }
        ret = TaggedCCARefValID(*refvalID)
    }

    return &ret, nil
}

func (o TaggedCCARefValID) Valid() error {
    return CCARefValID(o).Valid()
}

func (o TaggedCCARefValID) String() string {
    ret, err := json.Marshal(o)
    if err != nil {
        return ""
    }

    return string(ret)
}

func (o TaggedCCARefValID) Type() string {
    return CCARefValIDType
}

func (o TaggedCCARefValID) IsZero() bool {
    return len(o.SignerID) == 0
}

func NewMkeyCCARefValID(val any) (*Mkey, error) {
    ret, err := NewTaggedCCARefValID(val)
    if err != nil {
        return nil, err
    }

    return &Mkey{ret}, nil
}

