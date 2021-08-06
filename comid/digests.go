// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"

	"github.com/veraison/swid"
)

// Digests is an alias for an array of SWID HashEntry
type Digests []swid.HashEntry

// NewDigests instantiates an empty array of Digests
func NewDigests() *Digests {
	return new(Digests)
}

// TODO(tho) move this to veraison/swid
//
// Named Information Hash Algorithm Registry
// https://www.iana.org/assignments/named-information/named-information.xhtml#hash-alg
const (
	Sha256 uint64 = (iota + 1)
	Sha256_128
	Sha256_120
	Sha256_96
	Sha256_64
	Sha256_32
	Sha384
	Sha512
	Sha3_224
	Sha3_256
	Sha3_384
	Sha3_512
)

var algToValueLen = map[uint64]int{
	Sha256:     32,
	Sha256_128: 16,
	Sha256_120: 15,
	Sha256_96:  12,
	Sha256_64:  8,
	Sha256_32:  4,
	Sha384:     48,
	Sha512:     64,
	Sha3_224:   28,
	Sha3_256:   32,
	Sha3_384:   48,
	Sha3_512:   64,
}

// AddDigest create a new digest from the supplied arguments and appends it to
// the (already instantiated) Digests target.  The method is a no-op if it is
// invoked on a nil target and will refuse to add inconsistent algo/value
// combinations.
func (o *Digests) AddDigest(algID uint64, value []byte) *Digests {
	if o != nil {
		he := NewHashEntry(algID, value)
		if he == nil {
			return nil
		}
		*o = append(*o, *he)
	}
	return o
}

func (o Digests) Valid() error {
	for i, m := range o {
		if err := ValidHashEntry(m.HashAlgID, m.HashValue); err != nil {
			return fmt.Errorf("digest at index %d: %w", i, err)
		}
	}
	return nil
}

func NewHashEntry(algID uint64, value []byte) *swid.HashEntry {
	if ValidHashEntry(algID, value) != nil {
		return nil
	}

	return &swid.HashEntry{
		HashAlgID: algID,
		HashValue: value,
	}
}

func ValidHashEntry(algID uint64, value []byte) error {
	wantLen, ok := algToValueLen[algID]
	if !ok {
		return fmt.Errorf("unknown hash algorithm %d", algID)
	}

	gotLen := len(value)

	if wantLen != gotLen {
		return fmt.Errorf(
			"length mismatch for hash algorithm %d: want %d bytes, got %d",
			algID, wantLen, gotLen,
		)
	}

	return nil
}
