// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"math"
	"math/big"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberMatcher_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		input        *NumberMatcher
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title: "ok int64 unit",
			input: NewNumberMatcher(
				NumberFromInt(1),
			),
			expectedCBOR: []byte{0x01},
			expectedJSON: `{"type":"unit","value":{"type":"int","value":1}}`,
		},
		{
			title: "ok int64 set",
			input: NewNumberMatcher(
				NumberFromInt(1),
				NumberFromInt(2),
			),
			expectedCBOR: []byte{
				0xd9, 0x02, 0x36, // tag(566) [tagged-*-set]
				0x82, //             . array(2)
				0x01, //             . . 1
				0x02, //             . . 2
			},
			expectedJSON: `{"type":"set","value":[{"type":"int","value":1},{"type":"int","value":2}]}`,
		},
		{
			title: "ok int64 range",
			input: MustNewNumberRange(
				NumberFromInt(1),
				NumberFromInt(2),
			),
			expectedCBOR: []byte{
				0xd9, 0x02, 0x35, // tag(565) [tagged-number-range]
				0x82, //             . array(2)
				0x01, //             . . 1
				0x02, //             . . 2
			},
			expectedJSON: `{"type":"range","value":[{"type":"int","value":1},{"type":"int","value":2}]}`,
		},
		{
			title: "ok float64 unit",
			input: NewNumberMatcher(
				NumberFromFloat64(1.1),
			),
			expectedCBOR: []byte{
				0xfb, 0x3f, 0xf1, 0x99, 0x99, 0x99, 0x99, 0x99,
				0x9a,
			},
			expectedJSON: `{"type":"unit","value":{"type":"float","value":1.1}}`,
		},
		{
			title: "ok float64 range with single bound",
			input: NewNumberUpperBound(
				NumberFromFloat64(1.1),
			),
			expectedCBOR: []byte{
				0xd9, 0x02, 0x35, //                               tag(565) [tagged-number-range]
				0x82,                                           // . array(2)
				0xf6,                                           // . . null
				0xfb, 0x3f, 0xf1, 0x99, 0x99, 0x99, 0x99, 0x99, // . . 1.1
				0x9a,
			},
			expectedJSON: `{"type":"range","value":[null, {"type":"float","value":1.1}]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			encodedCBOR, err := tc.input.MarshalCBOR()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, encodedCBOR)

			var matcher NumberMatcher
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

func TestNumberMatcher_MatchesValue(t *testing.T) {
	testCases := []struct {
		title    string
		matcher  *NumberMatcher
		value    Number
		expected bool
	}{
		{
			title:    "match unit int64",
			matcher:  NewNumberMatcher(NumberFromInt(1)),
			value:    NumberFromInt(1),
			expected: true,
		},
		{
			title:    "match unit flaot64",
			matcher:  NewNumberMatcher(NumberFromFloat64(1.1)),
			value:    NumberFromFloat64(1.1),
			expected: true,
		},
		{
			title: "match set",
			matcher: NewNumberMatcher(
				NumberFromFloat64(1.0),
				NumberFromFloat64(2.1),
				NumberFromInt(3),
			),
			value:    NumberFromFloat64(1.0),
			expected: true,
		},
		{
			title: "match range",
			matcher: MustNewNumberRange(
				NumberFromFloat64(1.0),
				NumberFromFloat64(3.0),
			),
			value:    NumberFromFloat64(2.0),
			expected: true,
		},
		{
			title: "match range inf",
			matcher: MustNewNumberRange(
				NumberFromFloat64(math.Inf(-1)),
				NumberFromFloat64(math.Inf(1)),
			),
			value:    NumberFromFloat64(1e52),
			expected: true,
		},
		{
			title: "match upper bound",
			matcher: NewNumberUpperBound(
				NumberFromFloat64(3.0),
			),
			value:    NumberFromFloat64(2.0),
			expected: true,
		},
		{
			title: "match lower bound",
			matcher: NewNumberLowerBound(
				NumberFromFloat64(1.0),
			),
			value:    NumberFromFloat64(2.0),
			expected: true,
		},
		{
			title:    "fail unit int64",
			matcher:  NewNumberMatcher(NumberFromInt(1)),
			value:    NumberFromInt(2),
			expected: false,
		},
		{
			title:    "fail unit float64",
			matcher:  NewNumberMatcher(NumberFromFloat64(1.1)),
			value:    NumberFromFloat64(1.0),
			expected: false,
		},
		{
			title:    "fail unit mixed",
			matcher:  NewNumberMatcher(NumberFromFloat64(1.0)),
			value:    NumberFromInt(1),
			expected: false,
		},
		{
			title:    "fail unit NaN",
			matcher:  NewNumberMatcher(NumberFromFloat64(math.NaN())),
			value:    NumberFromFloat64(1.0),
			expected: false,
		},
		{
			title: "fail set mixed",
			matcher: NewNumberMatcher(
				NumberFromFloat64(1.0),
				NumberFromFloat64(2.1),
				NumberFromInt(3),
			),
			value:    NumberFromInt(2.0),
			expected: false,
		},
		{
			title: "fail range",
			matcher: MustNewNumberRange(
				NumberFromFloat64(1.0),
				NumberFromFloat64(3.0),
			),
			value:    NumberFromFloat64(4.0),
			expected: false,
		},
		{
			title: "fail range mixed",
			matcher: MustNewNumberRange(
				NumberFromFloat64(1.0),
				NumberFromFloat64(3.0),
			),
			value:    NumberFromInt(2),
			expected: false,
		},
		{
			title: "fail upper bound mixed",
			matcher: NewNumberUpperBound(
				NumberFromFloat64(1.0),
			),
			value:    NumberFromInt(2),
			expected: false,
		},
		{
			title: "fail lower bound mixed",
			matcher: NewNumberLowerBound(
				NumberFromFloat64(3.0),
			),
			value:    NumberFromInt(2),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			ret := tc.matcher.MatchesValue(tc.value)
			assert.Equal(t, tc.expected, ret)
		})
	}
}

func TestNumberMatcher_IterValues(t *testing.T) {
	testCases := []struct {
		title    string
		matcher  *NumberMatcher
		expected []Number
	}{
		{
			title:    "unit",
			matcher:  NewNumberMatcher(NumberFromInt(1)),
			expected: []Number{NumberFromInt(1)},
		},
		{
			title: "set",
			matcher: NewNumberMatcher(
				NumberFromInt(1),
				NumberFromInt(2),
			),
			expected: []Number{
				NumberFromInt(1),
				NumberFromInt(2),
			},
		},
		{
			title: "range",
			matcher: MustNewNumberRange(
				NumberFromInt(1),
				NumberFromInt(2),
			),
			expected: []Number{
				NumberFromInt(1),
				NumberFromInt(2),
			},
		},
		{
			title: "range lower bound",
			matcher: NewNumberLowerBound(
				NumberFromInt(1),
			),
			expected: []Number{
				NumberFromInt(1),
			},
		},
		{
			title: "range upper bound",
			matcher: NewNumberUpperBound(
				NumberFromInt(2),
			),
			expected: []Number{
				NumberFromInt(2),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result := slices.Collect(tc.matcher.IterValues())
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNumberMatcher_range_bounds(t *testing.T) {
	matcher := MustNewNumberRange(NumberFromInt(1), NumberFromInt(2))
	assert.Equal(t, NumberFromInt(1), *matcher.Variant().(RangeMatcher[Number]).Lower())
	assert.Equal(t, NumberFromInt(2), *matcher.Variant().(RangeMatcher[Number]).Upper())

	matcher = NewNumberLowerBound(NumberFromInt(1))
	assert.Equal(t, NumberFromInt(1), *matcher.Variant().(RangeMatcher[Number]).Lower())
	assert.Equal(t, (*Number)(nil), matcher.Variant().(RangeMatcher[Number]).Upper())

	matcher = NewNumberUpperBound(NumberFromInt(2))
	assert.Equal(t, (*Number)(nil), matcher.Variant().(RangeMatcher[Number]).Lower())
	assert.Equal(t, NumberFromInt(2), *matcher.Variant().(RangeMatcher[Number]).Upper())
}

func TestNumber_zero_value(t *testing.T) {
	var zero Number
	assert.Equal(t, nil, zero.Value())
	assert.Equal(t, (*big.Int)(nil), zero.BigInt())

	_, err := zero.Compare(NumberFromInt(0))
	assert.EqualError(t, err, ErrNumberTypeMismatch.Error())

	_, err = zero.Compare(NumberFromFloat64(0))
	assert.EqualError(t, err, ErrNumberTypeMismatch.Error())

	_, err = NumberFromInt(0).Compare(zero)
	assert.EqualError(t, err, ErrNumberTypeMismatch.Error())

	_, err = NumberFromFloat64(0).Compare(zero)
	assert.EqualError(t, err, ErrNumberTypeMismatch.Error())

	assert.Equal(t, int64(0), zero.Int64())
	assert.Equal(t, 0.0, zero.Float64())

	encodedCBOR, err := zero.MarshalCBOR()
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x00}, encodedCBOR)

	encodedJSON, err := zero.MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `{"type":"int","value":0}`, string(encodedJSON))

	var zero2 Number
	_, err = zero2.Compare(zero)
	assert.EqualError(t, err, ErrNumberTypeMismatch.Error())
}

func TestNumber_misc(t *testing.T) {
	number := NumberFromInt(1)
	assert.True(t, number.IsInt64())
	assert.False(t, number.IsFloat64())
	assert.EqualValues(t, big.NewInt(1), number.Value())
	assert.EqualValues(t, 1, number.Int64())
	assert.Equal(t, 0.0, number.Float64())

	number = NumberFromFloat64(1.0)
	assert.False(t, number.IsInt64())
	assert.True(t, number.IsFloat64())
	assert.Equal(t, int64(0), number.Int64())
	assert.Equal(t, 1.0, number.Float64())
}

func Test_tolerance(t *testing.T) {
	n1 := NumberFromFloat64(1.0)
	n2 := NumberFromFloat64(1.0 + 1e-15)
	result, _ := n1.Compare(n2)
	assert.Equal(t, 0, result)

	err := SetRelativeTolerance(1e-15)
	assert.NoError(t, err)
	result, _ = n1.Compare(n2)
	assert.Equal(t, -1, result)

	// reset back to default
	err = SetRelativeTolerance(1e-9)
	assert.NoError(t, err)

	n2 = NumberFromFloat64(1.01)
	result, _ = n1.Compare(n2)
	assert.Equal(t, -1, result)

	err = SetAbsoluteTolerance(1e-1)
	assert.NoError(t, err)

	result, _ = n1.Compare(n2)
	assert.Equal(t, 0, result)

	err = SetAbsoluteTolerance(-1.0)
	assert.ErrorContains(t, err, "tolerance cannot be negative")

	err = SetRelativeTolerance(-1.0)
	assert.ErrorContains(t, err, "tolerance cannot be negative")
}

func TestNumber_big(t *testing.T) {
	number := NumberFromUint(1)
	assert.EqualValues(t, big.NewInt(1), number.Value())
	assert.EqualValues(t, big.NewInt(1), number.BigInt())

	number, err := NumberFromBig(big.NewInt(1))
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(1), number.Value())

	number = MustNumberFromBig(new(big.Int).Lsh(big.NewInt(1), 63))
	assert.True(t, number.IsInt())
	assert.False(t, number.IsInt64())
	assert.Equal(t, uint64(1)<<63, number.Uint64())

	_, err = NumberFromBig(new(big.Int).Lsh(big.NewInt(1), 64))
	assert.ErrorContains(t, err, "outside CBOR int range")

	assert.Panics(t, func() {
		MustNumberFromBig(new(big.Int).Lsh(big.NewInt(1), 64))
	})

	_, err = NumberFromBig(nil)
	assert.ErrorContains(t, err, "nil")

	assert.Equal(t, uint64(0), NumberFromInt(-1).Uint64())
}

func TestNumber_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		title    string
		input    string
		expected Number
		err      string
	}{
		{
			title:    "ok int",
			input:    `{"type":"int","value":1}`,
			expected: NumberFromInt(1),
		},
		{
			title:    "ok int exp",
			input:    `{"type":"int","value":1e2}`,
			expected: NumberFromInt(100),
		},
		{
			title:    "ok float",
			input:    `{"type":"float","value":1e-2}`,
			expected: NumberFromFloat64(0.01),
		},
		{
			title: "err invalid number",
			input: `{"type":"int","value":"foo"}`,
			err:   "not a JSON number",
		},
		{
			title: "err bad type",
			input: `{"type":"foo","value":1}`,
			err:   `invalid Number value type: "foo"`,
		},
		{
			title: "err invalid number",
			input: `{"type":"int","value":"foo"}`,
			err:   "not a JSON number",
		},
		{
			title: "err int factional",
			input: `{"type":"int","value":1.1}`,
			err:   `JSON number "1.1" has a non-zero fractional part`,
		},
		{
			title: "err int out of range",
			input: `{"type":"int","value":18446744073709551616}`,
			err:   "outside CBOR int range",
		},
		{
			title: "err int factional pathological",
			input: `{"type":"int","value":1.00000000000000000000000000000000000000000000001}`,
			err:   `has a non-zero fractional part`,
		},
		{
			title: "err string with number",
			input: `{"type":"int","value":"42"}`,
			err:   "not a JSON number",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var number Number
			err := number.UnmarshalJSON([]byte(tc.input))

			if tc.err == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, number)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}
		})
	}
}
