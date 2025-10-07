// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package psa

import (
	"encoding/json"
	"fmt"

	"github.com/veraison/corim/comid"
)

// PSASoftwareComponentType is the type identifier for PSA software component keys
const PSASoftwareComponentType = "psa.software-component"

// PSASoftwareComponentKeyType represents a PSA software component key that implements IMKeyValue
type PSASoftwareComponentKeyType string

// NewPSASoftwareComponentKey creates a new PSA software component key
func NewPSASoftwareComponentKey(val any) (*PSASoftwareComponentKeyType, error) {
	var ret PSASoftwareComponentKeyType

	if val == nil {
		return &ret, nil
	}

	switch t := val.(type) {
	case PSASoftwareComponentKeyType:
		ret = t
	case *PSASoftwareComponentKeyType:
		ret = *t
	case string:
		ret = PSASoftwareComponentKeyType(t)
	default:
		return nil, fmt.Errorf("unexpected type for PSASoftwareComponentKeyType: %T", t)
	}

	return &ret, nil
}

// Valid validates the PSA software component key
func (o PSASoftwareComponentKeyType) Valid() error {
	if string(o) != PSASoftwareComponentType {
		return fmt.Errorf("invalid PSA software component key: expected %q, got %q", PSASoftwareComponentType, string(o))
	}
	return nil
}

// String returns the string representation
func (o PSASoftwareComponentKeyType) String() string {
	return string(o)
}

// Type returns the type identifier
func (o PSASoftwareComponentKeyType) Type() string {
	return PSASoftwareComponentType
}

// UnmarshalJSON unmarshals from JSON
func (o *PSASoftwareComponentKeyType) UnmarshalJSON(data []byte) error {
	var tmp string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	*o = PSASoftwareComponentKeyType(tmp)
	return nil
}

// TaggedPSASoftwareComponentKey is the CBOR-tagged version
type TaggedPSASoftwareComponentKey PSASoftwareComponentKeyType

// Valid validates the tagged PSA software component key
func (o TaggedPSASoftwareComponentKey) Valid() error {
	return PSASoftwareComponentKeyType(o).Valid()
}

// String returns the string representation
func (o TaggedPSASoftwareComponentKey) String() string {
	return PSASoftwareComponentKeyType(o).String()
}

// Type returns the type identifier
func (o TaggedPSASoftwareComponentKey) Type() string {
	return PSASoftwareComponentType
}

// UnmarshalJSON unmarshals from JSON
func (o *TaggedPSASoftwareComponentKey) UnmarshalJSON(data []byte) error {
	var tmp PSASoftwareComponentKeyType
	if err := tmp.UnmarshalJSON(data); err != nil {
		return err
	}
	*o = TaggedPSASoftwareComponentKey(tmp)
	return nil
}

// NewTaggedPSASoftwareComponentKey creates a new tagged PSA software component key
func NewTaggedPSASoftwareComponentKey(val any) (*TaggedPSASoftwareComponentKey, error) {
	key, err := NewPSASoftwareComponentKey(val)
	if err != nil {
		return nil, err
	}

	ret := TaggedPSASoftwareComponentKey(*key)
	return &ret, nil
}

// Factory function for creating PSA software component measurement keys
func newMkeyPSASoftwareComponent(val any) (*comid.Mkey, error) {
	ret, err := NewTaggedPSASoftwareComponentKey(val)
	if err != nil {
		return nil, err
	}

	return &comid.Mkey{Value: ret}, nil
}
