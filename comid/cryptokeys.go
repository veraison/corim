// Copyright 2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import "fmt"

// CryptoKeys is an array of *CryptoKey
type CryptoKeys []*CryptoKey

// NewCryptoKeys instantiates an empty CryptoKeys
func NewCryptoKeys() *CryptoKeys {
	return new(CryptoKeys)
}

// Add the supplied *CryptoKey to the CryptoKeys
func (o *CryptoKeys) Add(v *CryptoKey) *CryptoKeys {
	if o != nil && v != nil {
		*o = append(*o, v)
	}
	return o
}

// Valid returns an error if any of the contained keys fail to validate, or if
// CryptoKeys is empty
func (o CryptoKeys) Valid() error {
	if len(o) == 0 {
		return fmt.Errorf("no keys to validate")
	}

	for i, vk := range o {
		if err := vk.Valid(); err != nil {
			return fmt.Errorf("invalid key at index %d: %w", i, err)
		}
	}
	return nil
}

// String returns a string representation of all CryptoKeys
func (o CryptoKeys) String() string {
	if len(o) == 0 {
		return "[]"
	}

	result := "["
	for i, key := range o {
		if i > 0 {
			result += ", "
		}
		result += key.String()
	}
	result += "]"
	return result
}
