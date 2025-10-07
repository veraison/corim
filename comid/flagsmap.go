// Copyright 2023-2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"reflect"

	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
)

var True = true
var False = false

// Flag indicates whether a particular operational mode is active within the
// measured environment.
type Flag int

const (
	FlagIsConfigured Flag = iota
	FlagIsSecure
	FlagIsRecovery
	FlagIsDebug
	FlagIsReplayProtected
	FlagIsIntegrityProtected
	FlagIsRuntimeMeasured
	FlagIsImmutable
	FlagIsTcb
	FlagIsConfidentialityProtected
)

// FlagsMap describes a number of boolean operational modes. If a value is nil,
// then the operational mode is unknown.
type FlagsMap struct {
	// IsConfigured indicates whether the measured environment is fully
	// configured for normal operation.
	IsConfigured *bool `cbor:"0,keyasint,omitempty" json:"is-configured,omitempty"`
	// IsSecure indicates whether the measured environment's configurable
	// security settings are fully enabled.
	IsSecure *bool `cbor:"1,keyasint,omitempty" json:"is-secure,omitempty"`
	// IsRecovery indicates whether the measured environment is in recovery
	// mode.
	IsRecovery *bool `cbor:"2,keyasint,omitempty" json:"is-recovery,omitempty"`
	// IsDebug indicates whether the measured environment is in a debug
	// enabled mode.
	IsDebug *bool `cbor:"3,keyasint,omitempty" json:"is-debug,omitempty"`
	// IsReplayProtected indicates whether the measured environment is
	// protected from replay by a previous image that differs from the
	// current image.
	IsReplayProtected *bool `cbor:"4,keyasint,omitempty" json:"is-replay-protected,omitempty"`
	// IsIntegrityProtected indicates whether the measured environment is
	// protected from unauthorized update.
	IsIntegrityProtected *bool `cbor:"5,keyasint,omitempty" json:"is-integrity-protected,omitempty"`
	// IsRuntimeMeasured indicates whether the measured environment is
	// measured after being loaded into memory.
	IsRuntimeMeasured *bool `cbor:"6,keyasint,omitempty" json:"is-runtime-meas,omitempty"`
	// IsImmutable indicates whether the measured environment is immutable.
	IsImmutable *bool `cbor:"7,keyasint,omitempty" json:"is-immutable,omitempty"`
	// IsTcb indicates whether the measured environment is a trusted
	// computing base.
	IsTcb *bool `cbor:"8,keyasint,omitempty" json:"is-tcb,omitempty"`
	// IsConfidentialityProtected indicates whether the measured environment
	// is confidentiality protected. For example, if the measured environment consists of memory,
	// the sensitive values in memory are encrypted.
	IsConfidentialityProtected *bool `cbor:"9,keyasint,omitempty" json:"is-confidentiality-protected,omitempty"`

	Extensions
}

func NewFlagsMap() *FlagsMap {
	return &FlagsMap{}
}

// nolint:gocritic
func (o FlagsMap) IsEmpty() bool {
	if o.IsConfigured != nil || o.IsSecure != nil || o.IsRecovery != nil ||
		o.IsDebug != nil || o.IsReplayProtected != nil || o.IsIntegrityProtected != nil ||
		o.IsRuntimeMeasured != nil || o.IsImmutable != nil || o.IsTcb != nil ||
		o.IsConfidentialityProtected != nil {
		return false
	}

	return o.IsEmpty()
}

func (o *FlagsMap) AnySet() bool {
	if o.IsConfigured != nil || o.IsSecure != nil || o.IsRecovery != nil || o.IsDebug != nil ||
		o.IsReplayProtected != nil || o.IsIntegrityProtected != nil ||
		o.IsRuntimeMeasured != nil || o.IsImmutable != nil || o.IsTcb != nil ||
		o.IsConfidentialityProtected != nil {
		return true
	}

	return o.anySet()
}

func (o *FlagsMap) setFlag(value *bool, flags ...Flag) {
	for _, flag := range flags {
		switch flag {
		case FlagIsConfigured:
			o.IsConfigured = value
		case FlagIsSecure:
			o.IsSecure = value
		case FlagIsRecovery:
			o.IsRecovery = value
		case FlagIsDebug:
			o.IsDebug = value
		case FlagIsReplayProtected:
			o.IsReplayProtected = value
		case FlagIsIntegrityProtected:
			o.IsIntegrityProtected = value
		case FlagIsRuntimeMeasured:
			o.IsRuntimeMeasured = value
		case FlagIsImmutable:
			o.IsImmutable = value
		case FlagIsTcb:
			o.IsTcb = value
		case FlagIsConfidentialityProtected:
			o.IsConfidentialityProtected = value
		default:
			if value == &True {
				o.setTrue(flag)
			} else {
				o.setFalse(flag)
			}
		}
	}
}

func (o *FlagsMap) SetTrue(flags ...Flag) {
	o.setFlag(&True, flags...)
}

func (o *FlagsMap) SetFalse(flags ...Flag) {
	o.setFlag(&False, flags...)
}

func (o *FlagsMap) Clear(flags ...Flag) {
	for _, flag := range flags {
		switch flag {
		case FlagIsConfigured:
			o.IsConfigured = nil
		case FlagIsSecure:
			o.IsSecure = nil
		case FlagIsRecovery:
			o.IsRecovery = nil
		case FlagIsDebug:
			o.IsDebug = nil
		case FlagIsReplayProtected:
			o.IsReplayProtected = nil
		case FlagIsIntegrityProtected:
			o.IsIntegrityProtected = nil
		case FlagIsRuntimeMeasured:
			o.IsRuntimeMeasured = nil
		case FlagIsImmutable:
			o.IsImmutable = nil
		case FlagIsTcb:
			o.IsTcb = nil
		case FlagIsConfidentialityProtected:
			o.IsConfidentialityProtected = nil
		default:
			o.clear(flag)
		}
	}
}

func (o *FlagsMap) Get(flag Flag) *bool {
	switch flag {
	case FlagIsConfigured:
		return o.IsConfigured
	case FlagIsSecure:
		return o.IsSecure
	case FlagIsRecovery:
		return o.IsRecovery
	case FlagIsDebug:
		return o.IsDebug
	case FlagIsReplayProtected:
		return o.IsReplayProtected
	case FlagIsIntegrityProtected:
		return o.IsIntegrityProtected
	case FlagIsRuntimeMeasured:
		return o.IsRuntimeMeasured
	case FlagIsImmutable:
		return o.IsImmutable
	case FlagIsTcb:
		return o.IsTcb
	case FlagIsConfidentialityProtected:
		return o.IsConfidentialityProtected
	default:
		return o.get(flag)
	}
}

func (o FlagsMap) Equal(r FlagsMap) bool { //nolint:gocritic
	return reflect.DeepEqual(o, r)
}

func (o FlagsMap) CompareAgainstReference(r FlagsMap) bool { //nolint:gocritic
	return o.Equal(r)
}

// RegisterExtensions registers a struct as a collections of extensions
func (o *FlagsMap) RegisterExtensions(exts extensions.Map) error {
	for p, v := range exts {
		switch p {
		case ExtFlags:
			o.Register(v)
		default:
			return fmt.Errorf("%w: %q", extensions.ErrUnexpectedPoint, p)
		}
	}

	return nil
}

// GetExtensions returns previously registered extension
func (o *FlagsMap) GetExtensions() extensions.IMapValue {
	return o.IMapValue
}

// UnmarshalCBOR deserializes from CBOR
func (o *FlagsMap) UnmarshalCBOR(data []byte) error {
	return encoding.PopulateStructFromCBOR(dm, data, o)
}

// MarshalCBOR serializes to CBOR
// nolint:gocritic
func (o FlagsMap) MarshalCBOR() ([]byte, error) {
	return encoding.SerializeStructToCBOR(em, o)
}

// UnmarshalJSON deserializes from JSON
func (o *FlagsMap) UnmarshalJSON(data []byte) error {
	return encoding.PopulateStructFromJSON(data, o)
}

// MarshalJSON serializes to JSON
// nolint:gocritic
func (o FlagsMap) MarshalJSON() ([]byte, error) {
	return encoding.SerializeStructToJSON(o)
}

// Valid returns an error if the FlagsMap is invalid.
// nolint:gocritic
func (o FlagsMap) Valid() error {
	return o.validFlagsMap(&o)
}
