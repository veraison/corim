// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"errors"

	"github.com/veraison/corim/comid"
)

type EnvironmentSelector struct {
	Classes   *[]comid.Class    `cbor:"0,keyasint,omitempty"`
	Instances *[]comid.Instance `cbor:"1,keyasint,omitempty"`
	Groups    *[]comid.Group    `cbor:"2,keyasint,omitempty"`
}

// NewEnvironmentSelector creates a new EnvironmentSelector instance
func NewEnvironmentSelector() *EnvironmentSelector {
	return &EnvironmentSelector{}
}

// AddClass adds the supplied CoMID class to the target EnvironmentSelector
func (o *EnvironmentSelector) AddClass(v comid.Class) *EnvironmentSelector {
	if o.Classes == nil {
		o.Classes = new([]comid.Class)
	}

	*o.Classes = append(*o.Classes, v)

	return o
}

// AddInstance adds the supplied CoMID instance to the target EnvironmentSelector
func (o *EnvironmentSelector) AddInstance(v comid.Instance) *EnvironmentSelector {
	if o.Instances == nil {
		o.Instances = new([]comid.Instance)
	}

	*o.Instances = append(*o.Instances, v)

	return o
}

// AddGroup adds the supplied CoMID group to the target EnvironmentSelector
func (o *EnvironmentSelector) AddGroup(v comid.Group) *EnvironmentSelector {
	if o.Groups == nil {
		o.Groups = new([]comid.Group)
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
