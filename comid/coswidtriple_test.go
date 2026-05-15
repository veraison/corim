// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/swid"
)

func TestCoswidTriple_Valid(t *testing.T) {
	testCasses := []struct {
		title  string
		triple *CoswidTriple
		err    string
	}{
		{
			title: "ok",
			triple: &CoswidTriple{
				Environment: Environment{
					Instance: MustNewBytesInstance([]byte{0x01, 0x02, 0x03}),
				},
				TagIDs: CoswidTagIDs{*swid.NewTagID("1.2.3.4")},
			},
		},
		{
			title:  "bad empty environment",
			triple: &CoswidTriple{},
			err:    "environment must not be empty",
		},
		{
			title: "bad empty tag IDs",
			triple: &CoswidTriple{
				Environment: Environment{
					Instance: MustNewBytesInstance([]byte{0x01, 0x02, 0x03}),
				},
			},
			err: "tag-ids: must not be empty",
		},
		{
			title: "bad invalid tag",
			triple: &CoswidTriple{
				Environment: Environment{
					Instance: MustNewBytesInstance([]byte{0x01, 0x02, 0x03}),
				},
				TagIDs: CoswidTagIDs{swid.TagID{}},
			},
			err: "tag-id[0]: tag-id value is nil",
		},
	}

	for _, tc := range testCasses {
		t.Run(tc.title, func(t *testing.T) {
			err := tc.triple.Valid()
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}
		})
	}
}

func TestCoswidTagIDs_Add(t *testing.T) {
	tagIDs := CoswidTagIDs{}
	tagIDs.Add(swid.NewTagID("1.2.3.4"))
	assert.Len(t, tagIDs, 1)
}

func TestCoswidTriples(t *testing.T) {
	triples := NewCoswidTriples()
	assert.True(t, triples.IsEmpty())

	triples.Add(&CoswidTriple{})
	assert.False(t, triples.IsEmpty())

	err := triples.Valid()
	assert.ErrorContains(t, err, "triple[0]: environment: environment must not be empty")

	(*triples)[0] = CoswidTriple{
		Environment: Environment{
			Instance: MustNewBytesInstance([]byte{0x01, 0x02, 0x03}),
		},
		TagIDs: CoswidTagIDs{*swid.NewTagID("1.2.3.4")},
	}

	err = triples.Valid()
	assert.NoError(t, err)
}
