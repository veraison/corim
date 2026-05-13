// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"errors"
	"fmt"

	"github.com/veraison/corim/extensions"
)

// StatefulEnvironment describes the state of the target environment that is
// being matched by a CondEndorseTriple.
type StatefulEnvironment = ValueTriple

// StatefulEnvironments is a container for StatefulEnvironment instances and their extensions.
// It is a thin wrapper around extensions.Collection.
type StatefulEnvironments = ValueTriples

func NewStatefulEnvironments() *StatefulEnvironments {
	return NewValueTriples()
}

// CondEndorseTriple declares one or more conditions that, once matched,
// results in augmenting the Attester's actual state with the Endorsement
// Claims. It corresponds to CDDL conditional-endorsement-triple-record:
//
//	conditional-endorsement-triple-record = [
//	  conditions: [ + stateful-environment-record ]
//	  endorsements: [ + endorsed-triple-record ]
//	]
//
//	stateful-environment-record = [
//	  environment: environment-map,
//	  claims-list: [ + measurement-map ]
//	]
//
// (draft-ietf-rats-corim §5.1.7)
type CondEndorseTriple struct {
	_            struct{}             `cbor:",toarray"`
	Conditions   StatefulEnvironments `json:"conditions"`
	Endorsements ValueTriples         `json:"endorsements"`
}

func (o *CondEndorseTriple) RegisterExtensions(exts extensions.Map) error {
	if err := o.Conditions.RegisterExtensions(exts); err != nil {
		return err
	}

	return o.Endorsements.RegisterExtensions(exts)
}

func (o *CondEndorseTriple) GetExtensions() extensions.IMapValue {
	return o.Endorsements.GetExtensions()
}

func (o CondEndorseTriple) Valid() error {
	if o.Conditions.IsEmpty() {
		return errors.New("conditions validation failed: no condition entries")
	}

	if err := o.Conditions.Valid(); err != nil {
		return fmt.Errorf("conditions validation failed: %w", err)
	}

	if o.Endorsements.IsEmpty() {
		return errors.New("endorsements validation failed: no endorsement entries")
	}

	if err := o.Endorsements.Valid(); err != nil {
		return fmt.Errorf("endorsements validation failed: %w", err)
	}

	return nil
}

// CondEndorseTriples is a container for CondEndorseTriple instances and their extensions.
// It is a thin wrapper around extensions.Collection.
type CondEndorseTriples = extensions.Collection[CondEndorseTriple, *CondEndorseTriple]

func NewCondEndorseTriples() *CondEndorseTriples {
	return extensions.NewCollection[CondEndorseTriple]()
}
