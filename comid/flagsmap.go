// Copyright 2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

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
}

func NewFlagsMap() *FlagsMap {
	return &FlagsMap{}
}

func (o *FlagsMap) AnySet() bool {
	if o.IsConfigured != nil || o.IsSecure != nil || o.IsRecovery != nil || o.IsDebug != nil ||
		o.IsReplayProtected != nil || o.IsIntegrityProtected != nil ||
		o.IsRuntimeMeasured != nil || o.IsImmutable != nil || o.IsTcb != nil {
		return true
	}

	return false
}

func (o *FlagsMap) SetTrue(flags ...Flag) {
	for _, flag := range flags {
		switch flag {
		case FlagIsConfigured:
			o.IsConfigured = &True
		case FlagIsSecure:
			o.IsSecure = &True
		case FlagIsRecovery:
			o.IsRecovery = &True
		case FlagIsDebug:
			o.IsDebug = &True
		case FlagIsReplayProtected:
			o.IsReplayProtected = &True
		case FlagIsIntegrityProtected:
			o.IsIntegrityProtected = &True
		case FlagIsRuntimeMeasured:
			o.IsRuntimeMeasured = &True
		case FlagIsImmutable:
			o.IsImmutable = &True
		case FlagIsTcb:
			o.IsTcb = &True
		default:
		}
	}
}

func (o *FlagsMap) SetFalse(flags ...Flag) {
	for _, flag := range flags {
		switch flag {
		case FlagIsConfigured:
			o.IsConfigured = &False
		case FlagIsSecure:
			o.IsSecure = &False
		case FlagIsRecovery:
			o.IsRecovery = &False
		case FlagIsDebug:
			o.IsDebug = &False
		case FlagIsReplayProtected:
			o.IsReplayProtected = &False
		case FlagIsIntegrityProtected:
			o.IsIntegrityProtected = &False
		case FlagIsRuntimeMeasured:
			o.IsRuntimeMeasured = &False
		case FlagIsImmutable:
			o.IsImmutable = &False
		case FlagIsTcb:
			o.IsTcb = &False
		default:
		}
	}
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
		default:
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
	default:
		return nil
	}
}

// Valid returns an error if the FlagsMap is invalid.
func (o FlagsMap) Valid() error {
	return nil
}
