// Copyright 2024-2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewSVN(t *testing.T) {
	for _, tv := range []struct {
		Name     string
		Input    any
		Expected uint64
		Err      string
	}{
		{
			Name:     "string ok",
			Input:    "7",
			Expected: 7,
			Err:      "",
		},
		{
			Name:     "string err",
			Input:    "test",
			Expected: 0,
			Err:      `strconv.ParseUint: parsing "test": invalid syntax`,
		},
		{
			Name:     "uint",
			Input:    uint(7),
			Expected: 7,
			Err:      "",
		},
		{
			Name:     "uint64",
			Input:    uint64(7),
			Expected: 7,
			Err:      "",
		},
		{
			Name:     "int ok",
			Input:    7,
			Expected: 7,
			Err:      "",
		},
		{
			Name:     "int not ok",
			Input:    -7,
			Expected: 0,
			Err:      "SVN cannot be negative: -7",
		},
		{
			Name:     "int64 ok",
			Input:    int64(7),
			Expected: 7,
			Err:      "",
		},
		{
			Name:     "int64 not ok",
			Input:    int64(-7),
			Expected: 0,
			Err:      "SVN cannot be negative: -7",
		},
		{
			Name:     "nil",
			Input:    nil,
			Expected: 0,
			Err:      "",
		},
	} {
		t.Run(tv.Name, func(t *testing.T) {
			ret, err := NewSVN(tv.Input, "exact-value")
			exact := TaggedSVN(tv.Expected)
			expected := SVN{&exact}

			if tv.Err != "" {
				assert.EqualError(t, err, tv.Err)
			} else {
				assert.Equal(t, &expected, ret)
			}

			retMin, err := NewSVN(tv.Input, "min-value")
			min := TaggedMinSVN(tv.Expected)
			expected = SVN{&min}

			if tv.Err != "" {
				assert.EqualError(t, err, tv.Err)
			} else {
				assert.Equal(t, &expected, retMin)
			}
		})
	}

	in := TaggedSVN(7)

	_, err := NewSVN(in, "exact-value")
	assert.NoError(t, err)

	_, err = NewSVN(&in, "exact-value")
	assert.NoError(t, err)

	_, err = NewSVN(true, "exact-value")
	assert.EqualError(t, err, "unexpected type for SVN: bool")

	inMin := TaggedMinSVN(7)

	_, err = NewSVN(inMin, "min-value")
	assert.NoError(t, err)

	_, err = NewSVN(&inMin, "min-value")
	assert.NoError(t, err)

	_, err = NewSVN(true, "min-value")
	assert.EqualError(t, err, "unexpected type for SVN: bool")

	_, err = NewSVN(true, "test")
	assert.EqualError(t, err, "unknown SVN type: test")

	ret := MustNewSVN(7, "exact-value")
	assert.NotNil(t, ret)

	assert.Panics(t, func() { MustNewSVN(true, "exact-value") })
}

func TestSVN_JSON(t *testing.T) {
	var v SVN

	err := v.UnmarshalJSON([]byte(`{"type":"exact-value","value":2.3}`))
	assert.EqualError(t, err, "invalid SVN exact-value: json: cannot unmarshal number 2.3 into Go value of type comid.TaggedSVN")

	err = v.UnmarshalJSON([]byte(`{"type":"test","value":7}`))
	assert.EqualError(t, err, "unknown SVN type: test")

	err = v.UnmarshalJSON([]byte(`@@@`))
	assert.EqualError(t, err, "SVN decoding failure: invalid character '@' looking for beginning of value")

}

func Test_TaggedSVN_Equal_True(t *testing.T) {
	claim := TaggedSVN(7)
	ref := TaggedSVN(7)

	assert.True(t, claim.Equal(ref))
}

func Test_TaggedSVN_Equal_False(t *testing.T) {
	claim := TaggedSVN(7)
	ref := TaggedSVN(8)

	assert.False(t, claim.Equal(ref))
}

func Test_TaggedSVN_CompareAgainstRefMinSVN_True(t *testing.T) {
	claim := TaggedSVN(8)
	ref := TaggedMinSVN(7)

	assert.True(t, claim.CompareAgainstRefMinSVN(ref))
}

func Test_TaggedSVN_CompareAgainstRefMinSVN_False(t *testing.T) {
	claim := TaggedSVN(7)
	ref := TaggedMinSVN(8)

	assert.False(t, claim.CompareAgainstRefMinSVN(ref))
}

func Test_TaggedMinSVN_Equal_True(t *testing.T) {
	claim := TaggedMinSVN(8)
	ref := TaggedMinSVN(8)

	assert.True(t, claim.Equal(ref))
}

func Test_TaggedMinSVN_Equal_False(t *testing.T) {
	claim := TaggedMinSVN(7)
	ref := TaggedMinSVN(8)

	assert.False(t, claim.Equal(ref))
}
