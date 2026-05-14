// Copyright 2024-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRawValue_NewRawValue(t *testing.T) {
	testCases := []struct {
		title    string
		val      any
		typ      string
		expected *RawValue
		err      string
	}{
		{
			title: "ok bytes from []byte",
			val:   []byte{0x01, 0x02, 0x03},
			typ:   "bytes",
			expected: &RawValue{
				Value: &TaggedBytes{0x01, 0x02, 0x03},
			},
		},
		{
			title: "ok bytes from TaggedBytes",
			val:   TaggedBytes{0x01, 0x02, 0x03},
			typ:   "bytes",
			expected: &RawValue{
				Value: &TaggedBytes{0x01, 0x02, 0x03},
			},
		},
		{
			title: "ok bytes from *TaggedBytes",
			val:   &TaggedBytes{0x01, 0x02, 0x03},
			typ:   "bytes",
			expected: &RawValue{
				Value: &TaggedBytes{0x01, 0x02, 0x03},
			},
		},
		{
			title: "ok masked from []byte",
			val:   []byte{0x01, 0x02, 0x03},
			typ:   "masked",
			expected: &RawValue{
				Value: &TaggedMaskedRawValue{
					Value:     []byte{0x01, 0x02, 0x03},
					MaskBytes: []byte{0xFF, 0xFF, 0xFF},
				},
			},
		},
		{
			title: "ok masked from TaggedBytes",
			val:   TaggedBytes{0x01, 0x02, 0x03},
			typ:   "masked",
			expected: &RawValue{
				Value: &TaggedMaskedRawValue{
					Value:     []byte{0x01, 0x02, 0x03},
					MaskBytes: []byte{0xFF, 0xFF, 0xFF},
				},
			},
		},
		{
			title: "ok masked from *TaggedBytes",
			val:   &TaggedBytes{0x01, 0x02, 0x03},
			typ:   "masked",
			expected: &RawValue{
				Value: &TaggedMaskedRawValue{
					Value:     []byte{0x01, 0x02, 0x03},
					MaskBytes: []byte{0xFF, 0xFF, 0xFF},
				},
			},
		},
		{
			title: "ok masked from [2][]byte",
			val:   [2][]byte{{0x01, 0x02, 0x03}, {0x04, 0x05, 0x06}},
			typ:   "masked",
			expected: &RawValue{
				Value: &TaggedMaskedRawValue{
					Value:     []byte{0x01, 0x02, 0x03},
					MaskBytes: []byte{0x04, 0x05, 0x06},
				},
			},
		},
		{
			title: "ok masked from [][]byte (len 2)",
			val:   [][]byte{{0x01, 0x02, 0x03}, {0x04, 0x05, 0x06}},
			typ:   "masked",
			expected: &RawValue{
				Value: &TaggedMaskedRawValue{
					Value:     []byte{0x01, 0x02, 0x03},
					MaskBytes: []byte{0x04, 0x05, 0x06},
				},
			},
		},
		{
			title: "bad bytes from string",
			val:   "foo",
			typ:   "bytes",
			err:   "value must be a byte slice",
		},
		{
			title: "bad masked from string",
			val:   "foo",
			typ:   "masked",
			err:   "value must be a byte slice",
		},
		{
			title: "bad masked from [][]byte (len 1)",
			val:   [][]byte{{0x01, 0x02, 0x03}},
			typ:   "masked",
			err:   "[][]byte must contain exactly two elements",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			ret, err := NewRawValue(tc.val, tc.typ)
			if tc.err == "" {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expected, ret)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}

		})
	}
}

func TestRawValue_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		val          *RawValue
		expectedJSON string
		expectedCBOR []byte
	}{
		{
			title:        "bytes",
			val:          NewRawValueFromBytes([]byte{0x01, 0x02, 0x03}),
			expectedJSON: `{"type":"bytes","value":"AQID"}`,
			expectedCBOR: []byte{
				0xd9, 0x02, 0x30, // tag(560) [bytes]
				0x43, // . bstr(3)
				0x01, 0x02, 0x03,
			},
		},
		{
			title: "masked",
			val: MustNewRawValueWithMask(
				[]byte{0x01, 0x02, 0x03},
				[]byte{0x04, 0x05, 0x06},
			),
			expectedJSON: `{"type":"masked","value":{"value":"AQID","mask":"BAUG"}}`,
			expectedCBOR: []byte{
				0xd9, 0x02, 0x33, // tag(563) [masked-raw-value]
				0x82, // . array(2) [masked-raw-value]
				0x43, // . . [0]bstr(3)
				0x01, 0x02, 0x03,
				0x43, // . . [1]bstr(3)
				0x04, 0x05, 0x06,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			bytes, err := tc.val.MarshalJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(bytes))

			var decoded RawValue
			err = decoded.UnmarshalJSON(bytes)
			assert.NoError(t, err)
			assert.Equal(t, tc.val, &decoded)

			bytes, err = tc.val.MarshalCBOR()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, bytes)

			err = decoded.UnmarshalCBOR(bytes)
			assert.NoError(t, err)
			assert.Equal(t, tc.val, &decoded)
		})
	}
}

func TestRawValue_unmarshal_bad(t *testing.T) {
	var rv RawValue

	err := rv.UnmarshalJSON([]byte(`"foo"`))
	assert.ErrorContains(t, err, "cannot unmarshal string into Go value of type struct")

	err = rv.UnmarshalJSON([]byte(`{"type":"foo", "value":"bar"}`))
	assert.ErrorContains(t, err, "unexpected RawValue type: foo")

	err = rv.UnmarshalJSON([]byte(`{"type":"bytes", "value":2}`))
	assert.ErrorContains(t, err, "cannot unmarshal number into Go value of type comid.TaggedBytes")

	err = rv.UnmarshalJSON([]byte(`{"type":"masked", "value":"AQID"}`))
	assert.ErrorContains(t, err, "cannot unmarshal string into Go value of type comid.TaggedMaskedRawValue")

	err = rv.UnmarshalCBOR(MustHexDecode(t, "ff"))
	assert.ErrorContains(t, err, "unexpected \"break\" code")

	err = rv.UnmarshalCBOR(MustHexDecode(t, "d9023001"))
	assert.ErrorContains(t, err, "cannot unmarshal positive integer into Go value of type comid.TaggedBytes")

	rv.Value = nil
	err = rv.UnmarshalCBOR(MustHexDecode(t, "d9023343010203"))
	assert.ErrorContains(t, err,
		"cannot unmarshal byte string into Go value of type comid.TaggedMaskedRawValue")

	err = rv.UnmarshalCBOR(MustHexDecode(t, "d90233824301020301"))
	assert.ErrorContains(t, err,
		"cannot unmarshal positive integer into Go struct field *comid.TaggedMaskedRawValue.mask")
}

func TestRawValue_Equal(t *testing.T) {
	brv := NewRawValueFromBytes([]byte{0x01, 0x02, 0x03})
	mrv := MustNewMaskedRawValue([]byte{0x01, 0x02, 0x03})

	assert.True(t, brv.Equal(mrv))
	assert.True(t, mrv.Equal(brv))

	mrv = MustNewMaskedRawValue([]byte{0x11, 0x12, 0x13})

	assert.False(t, brv.Equal(mrv))
	assert.False(t, mrv.Equal(brv))

	mrv = MustNewRawValueWithMask([]byte{0x11, 0x12, 0x13}, []byte{0x0F, 0x0F, 0x0F})

	assert.True(t, brv.Equal(mrv))
	assert.True(t, mrv.Equal(brv))

	mrv2 := MustNewRawValueWithMask([]byte{0xA1, 0xB2, 0xC3}, []byte{0x0F, 0x0F, 0x0F})

	assert.True(t, mrv2.Equal(mrv))
	assert.True(t, mrv.Equal(mrv2))

	mrv2 = MustNewRawValueWithMask([]byte{0x11, 0x12, 0x13}, []byte{0x07, 0x07, 0x07})

	assert.False(t, mrv2.Equal(mrv))
	assert.False(t, mrv.Equal(mrv2))
}

func TestRawValue_CompareAgainstReference(t *testing.T) {
	brv := NewRawValueFromBytes([]byte{0x01, 0x02, 0x03})

	assert.True(t, brv.CompareAgainstReference([]byte{0x01, 0x02, 0x03}, nil))
	assert.False(t, brv.CompareAgainstReference([]byte{0x11, 0x12, 0x13}, nil))
	assert.True(t, brv.CompareAgainstReference([]byte{0x11, 0x12, 0x13}, []byte{0x0F, 0x0F, 0x0F}))

	mrv := MustNewRawValueWithMask([]byte{0x01, 0x02, 0x03}, []byte{0x0F, 0x0F, 0x0F})
	assert.True(t, mrv.CompareAgainstReference([]byte{0x11, 0x12, 0x13}, []byte{0x0F, 0x0F, 0x0F}))
	assert.False(t, mrv.CompareAgainstReference([]byte{0x11, 0x12, 0x13}, []byte{0x07, 0x07, 0x07}))
}
