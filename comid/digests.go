// Copyright 2021-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"bytes"
	"fmt"
)

// Digests is an alias for an array of Digest
type Digests []Digest

// NewDigests instantiates an empty array of Digests
func NewDigests() *Digests {
	return new(Digests)
}

// AddDigest create a new digest from the supplied arguments and appends it to
// the (already instantiated) Digests target.  The method is a no-op if it is
// invoked on a nil target and will refuse to add inconsistent algo/value
// combinations.
func (o *Digests) AddDigest(algID int, value []byte) *Digests {
	if o != nil {
		*o = append(*o, *NewDigestIntAlg(algID, value))
	}
	return o
}

func (o Digests) Valid() error {
	for i, m := range o {
		if err := m.Valid(); err != nil {
			return fmt.Errorf("digest at index %d: %w", i, err)
		}
	}
	return nil
}

// Equal confirms if the Digests instances are equal
//
// Two digests are considered to be equal if they meet the following criteria:
//   - They contain the same number of elements
//   - All the elements that use the same algorithm have the same value,
//     though the elements could be in any order
func (o Digests) Equal(r Digests) bool {
	om := make(map[DigestAlgorithm][][]byte)
	for _, oe := range o {
		vs, ok := om[oe.Algorithm]
		if ok {
			om[oe.Algorithm] = append(vs, oe.Value)
		} else {
			om[oe.Algorithm] = [][]byte{oe.Value}
		}
	}

outer:
	for _, re := range r {
		ovs, ok := om[re.Algorithm]
		if !ok {
			return false
		}

		for _, ov := range ovs {
			if bytes.Equal(ov, re.Value) {
				continue outer
			}
		}

		return false
	}

	return true
}

// CompareAgainstReference checks if digests object matches with a reference
//
//	See the following CoRIM spec for rules to compare
//	digests against a reference:
//	https://ietf-rats-wg.github.io/draft-ietf-rats-corim/draft-ietf-rats-corim.html#section-8.5.6.1.3
func (o Digests) CompareAgainstReference(r Digests) bool {
	result := false

	if len(r) == 0 {
		return false
	}

	// Insert the reference values into a map
	ref := make(map[DigestAlgorithm][]byte)
	for _, digest := range r {
		val, ok := ref[digest.Algorithm]
		if ok && !bytes.Equal(digest.Value, val) {
			// If two entries with the same hashing algorithm have different
			// values, that's an automatic false.
			return false
		}

		ref[digest.Algorithm] = digest.Value
	}

	// Check the object against the reference value map
	for _, digest := range o {
		val, ok := ref[digest.Algorithm]
		if !ok {
			continue
		}

		if !bytes.Equal(digest.Value, val) {
			// All hash values must be equal if a claim has the same
			// digest represented using multiple algorithms.
			return false
		}

		result = true
	}

	return result
}
