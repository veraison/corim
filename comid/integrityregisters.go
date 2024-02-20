// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/veraison/swid"
)

const TextType = "text"

// IRegisterIndex is the interface to hold register index
// Supported index types are uint and text
type IRegisterIndex interface{}

// IntegrityRegisters holds the Integrity Registers
type IntegrityRegisters struct {
	m map[IRegisterIndex]Digests
}

func NewIntegrityRegisters() *IntegrityRegisters {
	return &IntegrityRegisters{m: make(map[IRegisterIndex]Digests)}
}

// AddDigests allows inserting an array of digests at a specific index
// Supported index types are uint and text
func (i *IntegrityRegisters) AddDigests(index IRegisterIndex, digests Digests) error {
	if len(digests) == 0 {
		return fmt.Errorf("no digests to add")
	}
	for _, digest := range digests {
		if err := i.AddDigest(index, digest); err != nil {
			return fmt.Errorf("unable to add Digest: %w", err)
		}
	}
	return nil
}

// AddDigest allows inserting a digest at a specific index
// Supported index types are uint and text
func (i *IntegrityRegisters) AddDigest(index IRegisterIndex, digest swid.HashEntry) error {
	if i.m == nil {
		return fmt.Errorf("no register to add digest")
	}
	switch t := index.(type) {
	case string, uint, uint64:
		i.m[t] = append(i.m[t], digest)
	default:
		return fmt.Errorf("unexpected type for index: %T", t)
	}
	return nil
}

func (i IntegrityRegisters) MarshalCBOR() ([]byte, error) {
	return em.Marshal(i.m)
}

func (i *IntegrityRegisters) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &i.m)
}

type keyTypeandVal struct {
	KeyType string `json:"key_type"`
	Value   json.RawMessage
}

func (i IntegrityRegisters) MarshalJSON() ([]byte, error) {
	jmap := make(map[string]json.RawMessage)
	var newkey string
	for key, val := range i.m {
		var ktv keyTypeandVal
		switch t := key.(type) {
		case uint, uint64:
			ktv.KeyType = UintType
			newkey = fmt.Sprintf("%v", key)
		case string:
			ktv.KeyType = TextType
			newkey = key.(string)
		default:
			return nil, fmt.Errorf("unknown type %T for index-type-choice", t)
		}

		newval, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		ktv.Value = newval
		Value, err := json.Marshal(ktv)
		if err != nil {
			return nil, err
		}
		jmap[newkey] = Value
	}
	return json.Marshal(jmap)
}

func (i *IntegrityRegisters) UnmarshalJSON(data []byte) error {
	if i.m == nil {
		i.m = make(map[IRegisterIndex]Digests)
	}
	jmap := make(map[string]json.RawMessage)
	var index IRegisterIndex
	if err := json.Unmarshal(data, &jmap); err != nil {
		return fmt.Errorf("register map decoding failure: %w", err)
	}
	for key, val := range jmap {
		var ktv keyTypeandVal
		var d Digests

		if err := json.Unmarshal(val, &ktv); err != nil {
			return fmt.Errorf("unable to unmarshal keyTypeAndValue: %w", err)
		}
		if err := json.Unmarshal(ktv.Value, &d); err != nil {
			return fmt.Errorf("unable to unmarshal Digests: %w", err)
		}
		switch ktv.KeyType {
		case UintType:
			u, err := strconv.Atoi(key)
			if err != nil {
				return fmt.Errorf("unable to convert key to uint: %w", err)
			} else if u < 0 {
				return fmt.Errorf("invalid negative integer key")
			}
			index = uint(u)
		case TextType:
			index = key
		default:
			return fmt.Errorf("unexpected key type for index: %s", ktv.KeyType)
		}
		if err := i.AddDigests(index, d); err != nil {
			return fmt.Errorf("unable to add digests into register set: %w", err)
		}
	}
	return nil
}
