// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
)

// OpFlags implements the flags-type, mapping to DiceTcbInfo.flags via the
// operational flags not-configured, not-secure, recovery and debug.
// If the flags field is omitted, all flags are assumed to be 0.
type OpFlags uint8

const (
	OpFlagNotConfigured OpFlags = 1 << iota
	OpFlagNotSecure
	OpFlagRecovery
	OpFlagDebug
)

func NewOpFlags() *OpFlags {
	return new(OpFlags)
}

func (o OpFlags) Strings() []string {
	var a []string

	if o&OpFlagNotConfigured != 0 {
		a = append(a, "notConfigured")
	}

	if o&OpFlagNotSecure != 0 {
		a = append(a, "notSecure")
	}

	if o&OpFlagRecovery != 0 {
		a = append(a, "recovery")
	}

	if o&OpFlagDebug != 0 {
		a = append(a, "debug")
	}

	return a
}

func (o OpFlags) Valid() error {
	// While any combination in the lower half-byte is acceptable, the most
	// significant nibble must be all zeroes.
	if o&0xf0 != 0 {
		return fmt.Errorf("op-flags has unknown bits asserted: %02x", o)
	}

	return nil
}

// SetFlags sets the target object as specified.  As many flags as necessary can
// be specified in one call.
func (o *OpFlags) SetOpFlags(flags ...OpFlags) *OpFlags {
	if o != nil {
		for _, flag := range flags {
			*o |= flag
		}
	}
	return o
}

func (o OpFlags) IsSet(flag OpFlags) bool {
	return o&flag != 0
}

// UnmarshalJSON provides a custom deserializer for the OpFlags type that uses an
// array of identifiers rather than a bit set, e.g.:
//
//	"op-flags": [
//	  "notSecure",
//	  "debug"
//	]
func (o *OpFlags) UnmarshalJSON(data []byte) error {
	var a []string

	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	if len(a) == 0 {
		*o = 0
		return nil
	}

	for _, s := range a {
		switch s {
		case "notSecure":
			*o |= OpFlagNotSecure
		case "notConfigured":
			*o |= OpFlagNotConfigured
		case "recovery":
			*o |= OpFlagRecovery
		case "debug":
			*o |= OpFlagDebug
		default:
			// ignore unknown opflags
			continue
		}
	}

	return nil
}

func (o OpFlags) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Strings())
}
