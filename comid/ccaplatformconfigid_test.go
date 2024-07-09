// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCCAPlatformConfigID_Set_ok(t *testing.T) {
	var cca CCAPlatformConfigID

	err := cca.Set(TestCCALabel)
	assert.NoError(t, err)
}

func TestCCAPlatformConfigID_Set_nok(t *testing.T) {
	var cca CCAPlatformConfigID
	expectedErr := "empty input string"
	err := cca.Set("")
	assert.EqualError(t, err, expectedErr)
}

func TestCCAPlatformConfigID_Get_nok(t *testing.T) {
	var cca CCAPlatformConfigID
	expectedErr := "empty CCA platform config ID"
	_, err := cca.Get()
	assert.EqualError(t, err, expectedErr)
}

func TestNewTaggedCCAPlatformConfigID(t *testing.T) {
	testID := TaggedCCAPlatformConfigID("test")
	untagged := CCAPlatformConfigID("test")

	for _, tv := range []struct {
		Name     string
		Input    any
		Err      string
		Expected TaggedCCAPlatformConfigID
	}{
		{
			Name:     "TaggedCCAPlatformConfigID ok",
			Input:    testID,
			Expected: testID,
		},
		{
			Name:     "*TaggedCCAPlatformConfigID ok",
			Input:    &testID,
			Expected: testID,
		},
		{
			Name:     "CCAPlatformConfigID ok",
			Input:    untagged,
			Expected: testID,
		},
		{
			Name:     "*CCAPlatformConfigID ok",
			Input:    &untagged,
			Expected: testID,
		},
		{
			Name:     "string ok",
			Input:    "test",
			Expected: testID,
		},
		{
			Name:     "[]byte ok",
			Input:    []byte{0x74, 0x65, 0x73, 0x74},
			Expected: testID,
		},
		{
			Name:  "[]byte not ok",
			Input: []byte{0x80, 0x65, 0x73, 0x74},
			Err:   "bytes do not form a valid UTF-8 string",
		},
		{
			Name:  "bad type",
			Input: 7,
			Err:   "unexpected type for CCA platform-config-id: int",
		},
	} {
		t.Run(tv.Name, func(t *testing.T) {
			out, err := NewTaggedCCAPlatformConfigID(tv.Input)

			if tv.Err != "" {
				assert.Nil(t, out)
				assert.EqualError(t, err, tv.Err)
			} else {
				assert.Equal(t, tv.Expected, *out)
			}
		})
	}
}
