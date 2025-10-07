// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// CCA class type identifier
var CCAImplIDType = "cca.impl-id"

type CCAImplID [32]byte

func (o CCAImplID) String() string {
	return base64.StdEncoding.EncodeToString(o[:])
}

// how do we validate  cca id beyond lenght check?
func (o CCAImplID) Valid() error {
	return nil
}

// A CCAImplID with CBOR tagging
type TaggedCCAImplID CCAImplID

// Creates a new ClassID with a CCA Implementation ID value
func NewCCAImplIDClassID(val any) (*ClassID, error) {
	var ret TaggedCCAImplID

	if val == nil {
		return &ClassID{&TaggedCCAImplID{}}, nil
	}

	switch t := val.(type) {
	case []byte:
		if nb := len(t); nb != 32 {
			return nil, fmt.Errorf("bad cca.impl-id: got %d bytes, want 32", nb)
		}

		copy(ret[:], t)
	case string:
		v, err := base64.StdEncoding.DecodeString(t)
		if err != nil {
			return nil, fmt.Errorf("bad cca.impl-id: %w", err)
		}

		if nb := len(v); nb != 32 {
			return nil, fmt.Errorf("bad cca.impl-id: decoded %d bytes, want 32", nb)
		}

		copy(ret[:], v)
	case TaggedCCAImplID:
		copy(ret[:], t[:])
	case *TaggedCCAImplID:
		copy(ret[:], (*t)[:])
	case CCAImplID:
		copy(ret[:], t[:])
	case *CCAImplID:
		copy(ret[:], (*t)[:])
	default:
		return nil, fmt.Errorf("unexpected type for cca.impl-id: %T", t)
	}

	return &ClassID{&ret}, nil
}

// like above method but panics on error
func MustNewCCAImplIDClassID(val any) *ClassID {
	ret, err := NewCCAImplIDClassID(val)
	if err != nil {
		panic(err)
	}

	return ret
}

func (o TaggedCCAImplID) Valid() error {
	return CCAImplID(o).Valid()
}

func (o TaggedCCAImplID) String() string {
	return CCAImplID(o).String()
}

func (o TaggedCCAImplID) Type() string {
	return CCAImplIDType
}

func (o TaggedCCAImplID) Bytes() []byte {
	return o[:]
}

func (o TaggedCCAImplID) MarshalJSON() ([]byte, error) {
	return json.Marshal((o)[:])
}

func (o *TaggedCCAImplID) UnmarshalJSON(data []byte) error {
	var out []byte
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}

	if len(out) != 32 {
		return fmt.Errorf("bad cca.impl-id: decoded %d bytes, want 32", len(out))
	}

	copy((*o)[:], out)

	return nil
}

func NewCCAImplID(id [32]byte) CCAImplID {
	return CCAImplID(id)
}

func (o *ClassID) SetCCAImplID(implID CCAImplID) *ClassID {
	if o != nil {
		o.Value = TaggedCCAImplID(implID)
	}
	return o
}

func (o ClassID) GetCCAImplID() (CCAImplID, error) {
	switch t := o.Value.(type) {
	case *TaggedCCAImplID:
		return CCAImplID(*t), nil
	case TaggedCCAImplID:
		return CCAImplID(t), nil
	default:
		return CCAImplID{}, fmt.Errorf("class-id type is: %T", t)
	}
}

func NewClassCCAImplID(implID CCAImplID) *ClassID {
	return new(ClassID).SetCCAImplID(implID)
}

func init() {
	// Register the CCA Implementation ID type
	classIDValueRegister[CCAImplIDType] = NewCCAImplIDClassID
}
