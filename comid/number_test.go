// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"math/big"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegerMatcher_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		input        *IntegerMatcher
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title: "ok int64 unit",
			input: NewIntegerMatcher(
				IntegerFromInt(1),
			),
			expectedCBOR: []byte{0x01},
			expectedJSON: `{"type":"unit","value":1}`,
		},
		{
			title: "ok int64 set",
			input: NewIntegerMatcher(
				IntegerFromInt(1),
				IntegerFromInt(2),
			),
			expectedCBOR: []byte{
				0xd9, 0x02, 0x36, // tag(566) [tagged-*-set]
				0x82, //             . array(2)
				0x01, //             . . 1
				0x02, //             . . 2
			},
			expectedJSON: `{"type":"set","value":[1,2]}`,
		},
		{
			title: "ok int64 range",
			input: MustNewIntegerRange(
				IntegerFromInt(1),
				IntegerFromInt(2),
			),
			expectedCBOR: []byte{
				0xd9, 0x02, 0x35, // tag(565) [tagged-number-range]
				0x82, //             . array(2)
				0x01, //             . . 1
				0x02, //             . . 2
			},
			expectedJSON: `{"type":"range","value":[1,2]}`,
		},
		{
			title: "ok int64 range with single bound",
			input: NewIntegerUpperBound(
				IntegerFromInt(2),
			),
			expectedCBOR: []byte{
				0xd9, 0x02, 0x35, // tag(565) [tagged-number-range]
				0x82, //             . array(2)
				0xf6, //             . . nil
				0x02, //             . . 2
			},
			expectedJSON: `{"type":"range","value":[null,2]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			encodedCBOR, err := tc.input.MarshalCBOR()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, encodedCBOR)

			var matcher IntegerMatcher
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

func TestIntegerMatcher_MatchesValue(t *testing.T) {
	testCases := []struct {
		title    string
		matcher  *IntegerMatcher
		value    Integer
		expected bool
	}{
		{
			title:    "match unit",
			matcher:  NewIntegerMatcher(IntegerFromInt(1)),
			value:    IntegerFromInt(1),
			expected: true,
		},
		{
			title: "match set",
			matcher: NewIntegerMatcher(
				IntegerFromInt(1),
				IntegerFromInt(2),
				IntegerFromInt(3),
			),
			value:    IntegerFromInt(1),
			expected: true,
		},
		{
			title: "match range",
			matcher: MustNewIntegerRange(
				IntegerFromInt(1),
				IntegerFromInt(3),
			),
			value:    IntegerFromInt(2),
			expected: true,
		},
		{
			title: "match upper bound",
			matcher: NewIntegerUpperBound(
				IntegerFromInt(3),
			),
			value:    IntegerFromInt(2),
			expected: true,
		},
		{
			title: "match lower bound",
			matcher: NewIntegerLowerBound(
				IntegerFromInt(1),
			),
			value:    IntegerFromInt(2),
			expected: true,
		},
		{
			title:    "fail unit",
			matcher:  NewIntegerMatcher(IntegerFromInt(1)),
			value:    IntegerFromInt(2),
			expected: false,
		},
		{
			title: "fail set",
			matcher: NewIntegerMatcher(
				IntegerFromInt(1),
				IntegerFromInt(2),
				IntegerFromInt(3),
			),
			value:    IntegerFromInt(4),
			expected: false,
		},
		{
			title: "fail range",
			matcher: MustNewIntegerRange(
				IntegerFromInt(1),
				IntegerFromInt(3),
			),
			value:    IntegerFromInt(4),
			expected: false,
		},
		{
			title: "fail upper bound",
			matcher: NewIntegerUpperBound(
				IntegerFromInt(1),
			),
			value:    IntegerFromInt(2),
			expected: false,
		},
		{
			title: "fail lower bound",
			matcher: NewIntegerLowerBound(
				IntegerFromInt(3),
			),
			value:    IntegerFromInt(2),
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

func TestIntegerMatcher_IterValues(t *testing.T) {
	testCases := []struct {
		title    string
		matcher  *IntegerMatcher
		expected []Integer
	}{
		{
			title:    "unit",
			matcher:  NewIntegerMatcher(IntegerFromInt(1)),
			expected: []Integer{IntegerFromInt(1)},
		},
		{
			title: "set",
			matcher: NewIntegerMatcher(
				IntegerFromInt(1),
				IntegerFromInt(2),
			),
			expected: []Integer{
				IntegerFromInt(1),
				IntegerFromInt(2),
			},
		},
		{
			title: "range",
			matcher: MustNewIntegerRange(
				IntegerFromInt(1),
				IntegerFromInt(2),
			),
			expected: []Integer{
				IntegerFromInt(1),
				IntegerFromInt(2),
			},
		},
		{
			title: "range lower bound",
			matcher: NewIntegerLowerBound(
				IntegerFromInt(1),
			),
			expected: []Integer{
				IntegerFromInt(1),
			},
		},
		{
			title: "range upper bound",
			matcher: NewIntegerUpperBound(
				IntegerFromInt(2),
			),
			expected: []Integer{
				IntegerFromInt(2),
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

func TestIntegerMatcher_range_bounds(t *testing.T) {
	matcher := MustNewIntegerRange(IntegerFromInt(1), IntegerFromInt(2))
	assert.Equal(t, IntegerFromInt(1), *matcher.Variant().(RangeMatcher[Integer]).Lower())
	assert.Equal(t, IntegerFromInt(2), *matcher.Variant().(RangeMatcher[Integer]).Upper())

	matcher = NewIntegerLowerBound(IntegerFromInt(1))
	assert.Equal(t, IntegerFromInt(1), *matcher.Variant().(RangeMatcher[Integer]).Lower())
	assert.Equal(t, (*Integer)(nil), matcher.Variant().(RangeMatcher[Integer]).Upper())

	matcher = NewIntegerUpperBound(IntegerFromInt(2))
	assert.Equal(t, (*Integer)(nil), matcher.Variant().(RangeMatcher[Integer]).Lower())
	assert.Equal(t, IntegerFromInt(2), *matcher.Variant().(RangeMatcher[Integer]).Upper())
}

func TestInteger_zero_value(t *testing.T) {
	var zero Integer
	assert.Equal(t, new(big.Int), zero.BigInt())
	assert.True(t, zero.IsInt64())

	cmp, err := zero.Compare(IntegerFromInt(0))
	assert.NoError(t, err)
	assert.Equal(t, 0, cmp)

	cmp, err = IntegerFromInt(0).Compare(zero)
	assert.NoError(t, err)
	assert.Equal(t, 0, cmp)

	assert.Equal(t, int64(0), zero.Int64())

	encodedCBOR, err := zero.MarshalCBOR()
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x00}, encodedCBOR)

	encodedJSON, err := zero.MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, "0", string(encodedJSON))

	var zero2 Integer
	cmp, err = zero2.Compare(zero)
	assert.NoError(t, err)
	assert.Equal(t, 0, cmp)
}

func TestInteger_misc(t *testing.T) {
	integer := IntegerFromInt(1)
	assert.True(t, integer.IsInt64())
	assert.EqualValues(t, int64(1), integer.Int64())
}

func TestInteger_big(t *testing.T) {
	number := IntegerFromUint(1)
	assert.EqualValues(t, big.NewInt(1), number.BigInt())

	number, err := IntegerFromBig(big.NewInt(1))
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(1), number.BigInt())

	number = MustIntegerFromBig(new(big.Int).Lsh(big.NewInt(1), 63))
	assert.False(t, number.IsInt64())
	assert.Equal(t, uint64(1)<<63, number.Uint64())

	_, err = IntegerFromBig(new(big.Int).Lsh(big.NewInt(1), 64))
	assert.ErrorContains(t, err, "outside CBOR int range")

	assert.Panics(t, func() {
		MustIntegerFromBig(new(big.Int).Lsh(big.NewInt(1), 64))
	})

	_, err = IntegerFromBig(nil)
	assert.ErrorContains(t, err, "nil")

	assert.Equal(t, uint64(0), IntegerFromInt(-1).Uint64())
}

func TestInteger_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		title    string
		input    string
		expected Integer
		err      string
	}{
		{
			title:    "ok",
			input:    `1`,
			expected: IntegerFromInt(1),
		},
		{
			title:    "ok exp",
			input:    `1e2`,
			expected: IntegerFromInt(100),
		},
		{
			title: "err invalid number",
			input: `"foo"`,
			err:   "not a JSON number",
		},
		{
			title: "err factional",
			input: `1.1`,
			err:   `JSON number "1.1" has a non-zero fractional part`,
		},
		{
			title: "err out of range",
			input: `18446744073709551616`,
			err:   "outside CBOR int range",
		},
		{
			title: "err factional pathological",
			input: `1.00000000000000000000000000000000000000000000001`,
			err:   `has a non-zero fractional part`,
		},
		{
			title: "err string with number",
			input: `"42"`,
			err:   "not a JSON number",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var number Integer
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

func TestFloatMatcher_MatchesValue(t *testing.T) {
	testCases := []struct {
		title    string
		matcher  *FloatMatcher
		value    float64
		expected bool
	}{
		{
			title:    "match unit",
			matcher:  NewFloatMatcher(1.1),
			value:    1.1,
			expected: true,
		},
		{
			title:    "match set",
			matcher:  NewFloatMatcher(1.1, 2.2, 3.3),
			value:    2.2,
			expected: true,
		},
		{
			title:    "match range",
			matcher:  MustNewFloatRange(1.0, 3.0),
			value:    2.0,
			expected: true,
		},
		{
			title:    "match upper bound",
			matcher:  NewFloatUpperBound(3.0),
			value:    2.0,
			expected: true,
		},
		{
			title:    "match lower bound",
			matcher:  NewFloatLowerBound(1.0),
			value:    2.0,
			expected: true,
		},
		{
			title:    "fail unit",
			matcher:  NewFloatMatcher(1.1),
			value:    1.0,
			expected: false,
		},
		{
			title:    "fail set mixed",
			matcher:  NewFloatMatcher(1.0, 2.1, 3.0),
			value:    2.0,
			expected: false,
		},
		{
			title:    "fail range",
			matcher:  MustNewFloatRange(1.0, 3.0),
			value:    4.0,
			expected: false,
		},
		{
			title:    "fail upper bound",
			matcher:  NewFloatUpperBound(1.0),
			value:    2.0,
			expected: false,
		},
		{
			title:    "fail lower bound",
			matcher:  NewFloatLowerBound(3.0),
			value:    2.0,
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

func Test_tolerance(t *testing.T) {
	n1 := 1.0
	n2 := 1.0 + 1e-15
	result, _ := Float64Compare(n1, n2)
	assert.Equal(t, 0, result)

	err := SetRelativeTolerance(1e-15)
	assert.NoError(t, err)
	result, _ = Float64Compare(n1, n2)
	assert.Equal(t, -1, result)

	// reset back to default
	err = SetRelativeTolerance(1e-9)
	assert.NoError(t, err)

	n2 = 1.01
	result, _ = Float64Compare(n1, n2)
	assert.Equal(t, -1, result)

	err = SetAbsoluteTolerance(1e-1)
	assert.NoError(t, err)

	result, _ = Float64Compare(n1, n2)
	assert.Equal(t, 0, result)

	err = SetAbsoluteTolerance(-1.0)
	assert.ErrorContains(t, err, "tolerance cannot be negative")

	err = SetRelativeTolerance(-1.0)
	assert.ErrorContains(t, err, "tolerance cannot be negative")
}

func TestFloatMatcher_UnmarshalCBOR_negative(t *testing.T) {
	testCases := []struct {
		title string
		bytes []byte
		err   string
	}{
		{
			title: "err range bound both null",
			bytes: []byte{
				0xd9, 0x02, 0x35, // tag(565) [tagged-*-range]
				0x82, //             . array(2)
				0xf6, //             . . null
				0xf6, //             . . null
			},
			err: "both range bounds cannot be null",
		},
		{
			title: "err range lower higher than upper",
			bytes: []byte{
				0xd9, 0x02, 0x35, // tag(565) [tagged-*-range]
				0x82,             // . array(2)
				0xf9, 0x40, 0x00, // . . 2.0
				0xf9, 0x3c, 0x00, // . . 1.0
			},
			err: "upper bound cannot be less than lower bound",
		},
		{
			title: "err range contains NaN",
			bytes: []byte{
				0xd9, 0x02, 0x35, // tag(565) [tagged-*-range]
				0x82,             // . array(2)
				0xf9, 0x7e, 0x00, // . . NaN
				0xf9, 0x3c, 0x00, // . . 1.0
			},
			err: "lower bound: NaN",
		},
		{
			title: "err range contains Inf",
			bytes: []byte{
				0xd9, 0x02, 0x35, // tag(565) [tagged-*-range]
				0x82,             // . array(2)
				0xf9, 0x7c, 0x00, // . . Inf
				0xf9, 0x3c, 0x00, // . . 1.0
			},
			err: "lower bound: Inf",
		},
		{
			title: "err set with duplicates",
			bytes: []byte{
				0xd9, 0x02, 0x36, // tag(566) [tagged-*-set]
				0x82,             // . array(2)
				0xf9, 0x3c, 0x00, // . . 1.0
				0xf9, 0x3c, 0x00, // . . 1.0
			},
			err: "duplicate: 1",
		},
		{
			title: "err set with NaN",
			bytes: []byte{
				0xd9, 0x02, 0x36, // tag(566) [tagged-*-set]
				0x82,             // . array(2)
				0xf9, 0x7e, 0x00, // . . NaN
				0xf9, 0x3c, 0x00, // . . 1.0
			},
			err: "value[0]: NaN",
		},
		{
			title: "err set with -Inf",
			bytes: []byte{
				0xd9, 0x02, 0x36, // tag(566) [tagged-*-set]
				0x82,             // . array(2)
				0xf9, 0xfc, 0x00, // . . -Inf
				0xf9, 0x3c, 0x00, // . . 1.0
			},
			err: "value[0]: Inf",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var matcher FloatMatcher
			err := matcher.UnmarshalCBOR(tc.bytes)
			assert.ErrorContains(t, err, tc.err)
		})
	}

}

func TestFloatMatcher_UnmarshalJSON_negative(t *testing.T) {
	testCases := []struct {
		title string
		text  string
		err   string
	}{
		{
			title: "err range bound both null",
			text:  `{"type":"range","value":[null,null]}`,
			err:   "both range bounds cannot be null",
		},
		{
			title: "err range lower higher than upper",
			text:  `{"type":"range","value":[2,1]}`,
			err:   "upper bound cannot be less than lower bound",
		},
		{
			title: "err set with duplicates",
			text:  `{"type":"set","value":[1,1]}`,
			err:   "duplicate: 1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var matcher FloatMatcher
			err := matcher.UnmarshalJSON([]byte(tc.text))
			assert.ErrorContains(t, err, tc.err)
		})
	}

}
