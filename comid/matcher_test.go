// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesMatcher_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		input        *BytesMatcher
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title: "ok unit",
			input: NewBytesMatcher(
				[]byte{0xde, 0xad, 0xbe, 0xef},
			),
			expectedCBOR: []byte{
				0x44, // bstr(4)
				0xde, 0xad, 0xbe, 0xef,
			},
			expectedJSON: `{"type":"unit","value":"3q2+7w=="}`,
		},
		{
			title: "ok set",
			input: NewBytesMatcher(
				[]byte{0xde, 0xad, 0xbe, 0xef},
				[]byte{0xc0, 0xde, 0xca, 0xfe},
			),
			expectedCBOR: []byte{
				0xd9, 0x02, 0x36, // tag(566) [tagged-*-set]
				0x82, //             . array(2)
				0x44, //             . . bstr(4)
				0xde, 0xad, 0xbe, 0xef,
				0x44, //             . . bstr(4)
				0xc0, 0xde, 0xca, 0xfe,
			},
			expectedJSON: `{"type":"set","value":["3q2+7w==","wN7K/g=="]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			encodedCBOR, err := tc.input.MarshalCBOR()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, encodedCBOR)

			var matcher BytesMatcher
			err = matcher.UnmarshalCBOR(encodedCBOR)
			assert.NoError(t, err)
			assert.Equal(t, tc.input, &matcher)

			encodedJSON, err := tc.input.MarshalJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(encodedJSON))

			err = matcher.UnmarshalJSON(encodedJSON)
			assert.NoError(t, err)
			assert.Equal(t, tc.input, &matcher)
		})
	}
}

func TestTextMatcher_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		input        *TextMatcher
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title: "ok unit",
			input: NewTextMatcher(
				"foo",
			),
			expectedCBOR: []byte{
				0x63,             // tstr(3)
				0x66, 0x6f, 0x6f, // . "foo"
			},
			expectedJSON: `{"type":"unit","value":"foo"}`,
		},
		{
			title: "ok set",
			input: NewTextMatcher(
				"foo",
				"bar",
			),
			expectedCBOR: []byte{
				0xd9, 0x02, 0x36, // tag(566) [tagged-*-set]
				0x82,             //             . array(2)
				0x63,             //             . tstr(3)
				0x66, 0x6f, 0x6f, // . . "foo"
				0x63,             //             . tstr(3)
				0x62, 0x61, 0x72, // . . "bar"
			},
			expectedJSON: `{"type":"set","value":["foo","bar"]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			encodedCBOR, err := tc.input.MarshalCBOR()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, encodedCBOR)

			var matcher TextMatcher
			err = matcher.UnmarshalCBOR(encodedCBOR)
			assert.NoError(t, err)
			assert.Equal(t, tc.input, &matcher)

			encodedJSON, err := tc.input.MarshalJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(encodedJSON))

			err = matcher.UnmarshalJSON(encodedJSON)
			assert.NoError(t, err)
			assert.Equal(t, tc.input, &matcher)
		})
	}
}

func TestTextMatcher_CompareAgainstReference(t *testing.T) {
	testCases := []struct {
		title    string
		value    *TextMatcher
		ref      *TextMatcher
		expected bool
	}{
		{
			title:    "ok unit to unit",
			value:    NewTextMatcher("foo"),
			ref:      NewTextMatcher("foo"),
			expected: true,
		},
		{
			title:    "ok unit to set",
			value:    NewTextMatcher("foo"),
			ref:      NewTextMatcher("foo", "bar"),
			expected: true,
		},
		{
			title:    "nok unit to unit",
			value:    NewTextMatcher("foo"),
			ref:      NewTextMatcher("bar"),
			expected: false,
		},
		{
			title:    "nok unit to set",
			value:    NewTextMatcher("foo"),
			ref:      NewTextMatcher("qux", "bar"),
			expected: false,
		},
		{
			title:    "nok set to unit",
			value:    NewTextMatcher("foo", "bar"),
			ref:      NewTextMatcher("foo"),
			expected: false,
		},
		{
			title:    "nok set to set",
			value:    NewTextMatcher("foo", "bar"),
			ref:      NewTextMatcher("foo", "bar"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result := tc.value.CompareAgainstReference(tc.ref)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBytesMatcher_CompareAgainstReference(t *testing.T) {
	testCases := []struct {
		title    string
		value    *BytesMatcher
		ref      *BytesMatcher
		expected bool
	}{
		{
			title:    "ok unit to unit",
			value:    NewBytesMatcher([]byte{0xde, 0xad, 0xbe, 0xef}),
			ref:      NewBytesMatcher([]byte{0xde, 0xad, 0xbe, 0xef}),
			expected: true,
		},
		{
			title: "ok unit to set",
			value: NewBytesMatcher([]byte{0xde, 0xad, 0xbe, 0xef}),
			ref: NewBytesMatcher(
				[]byte{0xde, 0xad, 0xbe, 0xef},
				[]byte{0xc0, 0xde, 0xca, 0xfe},
			),
			expected: true,
		},
		{
			title:    "nok unit to unit",
			value:    NewBytesMatcher([]byte{0xde, 0xad, 0xbe, 0xef}),
			ref:      NewBytesMatcher([]byte{0xc0, 0xde, 0xca, 0xfe}),
			expected: false,
		},
		{
			title: "nok unit to set",
			value: NewBytesMatcher([]byte{0xde, 0xad, 0xbe, 0xef}),
			ref: NewBytesMatcher(
				[]byte{0xba, 0xad, 0xf0, 0x0d},
				[]byte{0xc0, 0xde, 0xca, 0xfe},
			),
			expected: false,
		},
		{
			title: "nok set to unit",
			value: NewBytesMatcher(
				[]byte{0xde, 0xad, 0xbe, 0xef},
				[]byte{0xc0, 0xde, 0xca, 0xfe},
			),
			ref:      NewBytesMatcher([]byte{0xde, 0xad, 0xbe, 0xef}),
			expected: false,
		},
		{
			title: "nok set to set",
			value: NewBytesMatcher(
				[]byte{0xde, 0xad, 0xbe, 0xef},
				[]byte{0xc0, 0xde, 0xca, 0xfe},
			),
			ref: NewBytesMatcher(
				[]byte{0xde, 0xad, 0xbe, 0xef},
				[]byte{0xc0, 0xde, 0xca, 0xfe},
			),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result := tc.value.CompareAgainstReference(tc.ref)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNewSet_negative(t *testing.T) {
	_, err := NewSetMatcher("foo")
	assert.ErrorContains(t, err, "must specify at least two distinct elements")
}

func TestSet_Add(t *testing.T) {
	set := MustNewSetMatcher("foo", "bar")
	assert.True(t, set.Add("qux"))
	assert.False(t, set.Add("foo"))
	assert.Equal(t, 3, set.Len())
	assert.Equal(t, []string{"foo", "bar", "qux"}, set.Values())
}
