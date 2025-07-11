// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"errors"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/corim/comid"
)

func unmarshalStatefulFirstPass(data []byte) ([]cbor.RawMessage, error) {
	var a []cbor.RawMessage

	if err := cbor.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("CBOR decoding: %w", err)
	}

	alen := len(a)

	if alen < 1 || alen > 2 {
		return nil, fmt.Errorf("wrong number of entries (%d) in the array", alen)
	}

	return a, nil
}

type StatefulClass struct {
	Class        *comid.Class
	Measurements *comid.Measurements
}

func (o StatefulClass) MarshalCBOR() ([]byte, error) {
	if o.Class == nil {
		return nil, errors.New("mandatory field class not set")
	}

	a := []any{o.Class}
	if o.Measurements != nil {
		a = append(a, o.Measurements)
	}

	return cbor.Marshal(a)
}

func (o *StatefulClass) UnmarshalCBOR(data []byte) error {
	a, err := unmarshalStatefulFirstPass(data)
	if err != nil {
		return fmt.Errorf("unmarshaling StatefulClass: %w", err)
	}

	if err := cbor.Unmarshal(a[0], &o.Class); err != nil {
		return fmt.Errorf("unmarshaling StatefulClass Class: %w", err)
	}

	if len(a) == 2 {
		if err := cbor.Unmarshal(a[1], &o.Measurements); err != nil {
			return fmt.Errorf("unmarshaling StatefulClass Measurements: %w", err)
		}
	}

	return nil
}

type StatefulInstance struct {
	Instance     *comid.Instance
	Measurements *comid.Measurements
}

func (o StatefulInstance) MarshalCBOR() ([]byte, error) {
	if o.Instance == nil {
		return nil, errors.New("mandatory field instance not set")
	}

	a := []any{o.Instance}
	if o.Measurements != nil {
		a = append(a, o.Measurements)
	}

	return cbor.Marshal(a)
}

func (o *StatefulInstance) UnmarshalCBOR(data []byte) error {
	a, err := unmarshalStatefulFirstPass(data)
	if err != nil {
		return fmt.Errorf("unmarshaling StatefulInstance: %w", err)
	}

	if err := cbor.Unmarshal(a[0], &o.Instance); err != nil {
		return fmt.Errorf("unmarshaling StatefulInstance Instance: %w", err)
	}

	if len(a) == 2 {
		if err := cbor.Unmarshal(a[1], &o.Measurements); err != nil {
			return fmt.Errorf("unmarshaling StatefulInstance Measurements: %w", err)
		}
	}

	return nil
}

type StatefulGroup struct {
	Group        *comid.Group
	Measurements *comid.Measurements
}

func (o StatefulGroup) MarshalCBOR() ([]byte, error) {
	if o.Group == nil {
		return nil, errors.New("mandatory field group not set")
	}

	a := []any{o.Group}
	if o.Measurements != nil {
		a = append(a, o.Measurements)
	}

	return cbor.Marshal(a)
}

func (o *StatefulGroup) UnmarshalCBOR(data []byte) error {
	a, err := unmarshalStatefulFirstPass(data)
	if err != nil {
		return fmt.Errorf("unmarshaling StatefulGroup: %w", err)
	}

	if err := cbor.Unmarshal(a[0], &o.Group); err != nil {
		return fmt.Errorf("unmarshaling StatefulGroup Group: %w", err)
	}

	if len(a) == 2 {
		if err := cbor.Unmarshal(a[1], &o.Measurements); err != nil {
			return fmt.Errorf("unmarshaling StatefulGroup Measurements: %w", err)
		}
	}

	return nil
}

type EnvironmentSelector struct {
	Classes   *[]StatefulClass    `cbor:"0,keyasint,omitempty"`
	Instances *[]StatefulInstance `cbor:"1,keyasint,omitempty"`
	Groups    *[]StatefulGroup    `cbor:"2,keyasint,omitempty"`
}

// NewEnvironmentSelector creates a new EnvironmentSelector instance
func NewEnvironmentSelector() *EnvironmentSelector {
	return &EnvironmentSelector{}
}

// AddClass adds the supplied CoMID class to the target EnvironmentSelector
func (o *EnvironmentSelector) AddClass(v StatefulClass) *EnvironmentSelector {
	if o.Classes == nil {
		o.Classes = new([]StatefulClass)
	}

	*o.Classes = append(*o.Classes, v)

	return o
}

// AddInstance adds the supplied CoMID instance to the target EnvironmentSelector
func (o *EnvironmentSelector) AddInstance(v StatefulInstance) *EnvironmentSelector {
	if o.Instances == nil {
		o.Instances = new([]StatefulInstance)
	}

	*o.Instances = append(*o.Instances, v)

	return o
}

// AddGroup adds the supplied CoMID group to the target EnvironmentSelector
func (o *EnvironmentSelector) AddGroup(v StatefulGroup) *EnvironmentSelector {
	if o.Groups == nil {
		o.Groups = new([]StatefulGroup)
	}

	*o.Groups = append(*o.Groups, v)

	return o
}

// Valid ensures that the target EnvironmentSelector is correctly populated
func (o EnvironmentSelector) Valid() error {
	if o.Classes == nil && o.Groups == nil && o.Instances == nil {
		return errors.New("non-empty<> constraint violation")
	}

	if o.Classes != nil {
		if o.Instances != nil || o.Groups != nil {
			return errors.New("only one selector type is allowed")
		}
	} else if o.Instances != nil {
		if o.Groups != nil {
			return errors.New("only one selector type is allowed")
		}
	}

	return nil
}
