// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
)

// Environment stores the identifying information about a target or attesting
// environment at the class, instance and group scope.  The Environment type
// has JSON and CBOR serializations.
type Environment struct {
	Class    *Class    `cbor:"0,keyasint,omitempty" json:"class,omitempty"`
	Instance *Instance `cbor:"1,keyasint,omitempty" json:"instance,omitempty"`
	Group    *Group    `cbor:"2,keyasint,omitempty" json:"group,omitempty"`
}

// Valid checks the validity (according to the spec) of the target Environment
func (o Environment) Valid() error {
	// non-empty<>
	if o.Class == nil && o.Instance == nil && o.Group == nil {
		return fmt.Errorf("environment must not be empty")
	}

	if o.Class != nil {
		if err := o.Class.Valid(); err != nil {
			return fmt.Errorf("environment validation failed: %w", err)
		}
	}

	if o.Instance != nil {
		if err := o.Instance.Valid(); err != nil {
			return fmt.Errorf("environment validation failed: %w", err)
		}
	}

	if o.Group != nil {
		if err := o.Group.Valid(); err != nil {
			return fmt.Errorf("environment validation failed: %w", err)
		}
	}

	return nil
}

// ToCBOR serializes the target Environment to CBOR (if the Environment is "valid")
func (o Environment) ToCBOR() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return em.Marshal(&o)
}

// FromCBOR deserializes the supplied CBOR data into the target Environment
func (o *Environment) FromCBOR(data []byte) error {
	if err := dm.Unmarshal(data, o); err != nil {
		return err
	}

	return o.Valid()
}

// FromJSON deserializes the supplied JSON string into the target Environment
func (o *Environment) FromJSON(data []byte) error {
	if err := json.Unmarshal(data, o); err != nil {
		return err
	}

	return o.Valid()
}
