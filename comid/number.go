// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
)

var (
	// 2^64 - 1
	cborUintMax = new(big.Int).SetUint64(math.MaxUint64)
	// -2^64
	cborNintMin = new(big.Int).Neg(new(big.Int).Lsh(big.NewInt(1), 64))
)

// Integer encapsulates a CoRIM integer value. This is a value in the range
// [-math.Uint64Max, math.Uint64Max].
type Integer struct {
	value *big.Int
}

// IntegerFromInt creates a new Integer from the provided int.
func IntegerFromInt(value int) Integer {
	return IntegerFromInt64(int64(value))
}

// IntegerFromInt64 creates a new Integer from the provided int64.
func IntegerFromInt64(value int64) Integer {
	var num big.Int
	num.SetInt64(value)
	return Integer{&num}
}

// IntegerFromUint creates a new Integer from the provided uint.
func IntegerFromUint(value uint) Integer {
	return IntegerFromUint64(uint64(value))
}

// IntegerFromUint64 creates a new Integer from the provided uint64.
func IntegerFromUint64(value uint64) Integer {
	var num big.Int
	num.SetUint64(value)
	return Integer{&num}
}

// MustIntegerFromBig creates a new Integer from the provided big.Int. It panics
// if the provided value does not fit inside a CBOR integer.
func MustIntegerFromBig(value *big.Int) Integer {
	ret, err := IntegerFromBig(value)
	if err != nil {
		panic(err)
	}

	return ret
}

// IntegerFromBig creates a new Integer from the provided big.Int. If the
// provided value does not fit inside a CBOR integer, an error is returned.
func IntegerFromBig(value *big.Int) (Integer, error) {
	if value == nil {
		return Integer{}, errors.New("nil")
	}

	if value.Cmp(cborUintMax) == 1 || value.Cmp(cborNintMin) == -1 {
		return Integer{}, errors.New("outside CBOR int range")
	}

	return Integer{new(big.Int).Set(value)}, nil
}

// IsInt64 returns true if the receiver encapsulates an integer that can fit
// inside int64.
func (o Integer) IsInt64() bool {
	if o.value == nil {
		return true
	}

	return o.value.IsInt64()
}

// Uint64 returns the underlying integer value as a uint64. If the value is
// negative, 0 is returned instead.
func (o Integer) Uint64() uint64 {
	if o.value == nil || o.value.Cmp(big.NewInt(0)) < 0 {
		return 0
	}

	if !o.value.IsUint64() {
		// should not get here as cannot create a Integer from bit.Int
		// that is larger than a uint64.
		panic("too big")
	}

	return o.value.Uint64()
}

// Int64 returns the underlying value as an int64. If the underlying value does
// not fit inside int64, 0 is returned instead.
func (o Integer) Int64() int64 {
	if o.value == nil || !o.value.IsInt64() {
		return 0
	}

	return o.value.Int64()
}

// BigInt returns a pointer to the underlying big.Int value.
func (o Integer) BigInt() *big.Int {
	ret := new(big.Int)
	if o.value != nil {
		ret.Set(o.value)
	}

	return ret
}

// Compare returns -1 if the receiver is less than the provided Integer, 1 if
// the receiver is greater than provided Integer, or 0 if the two Integer's are
// the same. (The behavior is equivalent to cmp.Compare() behavior for
// cmd.Ordered.)
func (o Integer) Compare(other Integer) (int, error) {
	return o.BigInt().Cmp(other.BigInt()), nil
}

func (o Integer) MarshalCBOR() ([]byte, error) {
	if o.value != nil {
		return em.Marshal(o.value)
	} else {
		return em.Marshal(0)
	}
}

func (o *Integer) UnmarshalCBOR(data []byte) error {
	if len(data) < 1 {
		return errors.New("truncated input")
	}

	majorType := (data[0] & 0xe0) >> 5

	switch majorType {
	case 0, 1:
		var value big.Int
		if err := dm.Unmarshal(data, &value); err != nil {
			return err
		}

		o.value = &value
	default:
		return fmt.Errorf("unexpected major type %d for Integer", majorType)
	}

	return nil
}

func (o Integer) MarshalJSON() ([]byte, error) {
	if o.value == nil {
		return json.Marshal(0)
	}

	return json.Marshal(o.value)
}

func (o *Integer) UnmarshalJSON(data []byte) error {
	value, err := unmarshalBigIntFromJSON(data)
	if err != nil {
		return err
	}

	if value.Cmp(cborUintMax) == 1 || value.Cmp(cborNintMin) == -1 {
		return errors.New("outside CBOR int range")
	}

	o.value = value

	return nil
}

// IntegerMatcher matches Integer values to other Integer values, sets of Integer
// values, or Integer ranges.
type IntegerMatcher = Matcher[Integer]

// NewIntegerMatcher returns a pointer to a IntegerMatcher containing the
// provided values. If No tail values are specified, a unit matcher is created
// using head, otherwise a set matcher containing head and tail values is
// created.
func NewIntegerMatcher(head Integer, tail ...Integer) *IntegerMatcher {
	if len(tail) == 0 {
		return &IntegerMatcher{MustNewUnitMatcher(head)}
	}

	values := append([]Integer{head}, tail...)

	return &IntegerMatcher{MustNewSetMatcher(values...)}
}

// MustNewIntegerRange returns a pointer to a IntegerMatcher that matches the
// Integer range [lower, upper] (inclusive on both ends). It panics if the upper
// bound is less than the lower bound.
func MustNewIntegerRange(lower, upper Integer) *IntegerMatcher {
	matcher, err := NewIntegerRange(lower, upper)
	if err != nil {
		panic(err)
	}

	return matcher
}

// NewIntegerRange returns a pointer to a IntegerMatcher that matches the Integer
// range [lower, upper] (inclusive on both ends). An error is returned if the
// upper bound is less than the lower bound.
func NewIntegerRange(lower, upper Integer) (*IntegerMatcher, error) {
	rang, err := NewRangeMatcher(&lower, &upper)
	if err != nil {
		return nil, err
	}

	return &IntegerMatcher{rang}, nil
}

// NewIntegerUpperBound returns a pointer to a IntegerMatcher that matches any
// Integer less than or equal to the specified upper bound.
func NewIntegerUpperBound(upper Integer) *IntegerMatcher {
	rang, err := NewRangeMatcher(nil, &upper)
	if err != nil {
		// won't get here as one of the bounds is nil
		panic(err)
	}

	return &IntegerMatcher{rang}
}

// NewIntegerLowerBound returns a pointer to a IntegerMatcher that matches any
// Integer greater than or equal to the specified lower bound.
func NewIntegerLowerBound(lower Integer) *IntegerMatcher {
	rang, err := NewRangeMatcher(&lower, nil)
	if err != nil {
		// won't get here as one of the bounds is nil
		panic(err)
	}

	return &IntegerMatcher{rang}
}

// IntegerCompare returns -1 if lhs < rhs, 1 if lhs > rhs, and 0 if lhs == rhs.
func IntegerCompare(lhs, rhs Integer) (int, error) {
	return lhs.Compare(rhs)
}

// FloatMatcher matches Float values to other Float values, sets of Float
// values, or Float ranges.
type FloatMatcher = Matcher[float64]

// NewFloatMatcher returns a pointer to a FloatMatcher containing the
// provided values. If No tail values are specified, a unit matcher is created
// using head, otherwise a set matcher containing head and tail values is
// created. It panics is any of the values are NaN.
func NewFloatMatcher(head float64, tail ...float64) *FloatMatcher {
	if len(tail) == 0 {
		return &FloatMatcher{MustNewUnitMatcher(head)}
	}

	values := append([]float64{head}, tail...)

	return &FloatMatcher{MustNewSetMatcher(values...)}
}

// MustNewFloatRange returns a pointer to a FloatMatcher that matches the
// Float range [lower, upper] (inclusive on both ends). It panics if the upper
// bound is less than the lower bound or either of the values is NaN.
func MustNewFloatRange(lower, upper float64) *FloatMatcher {
	matcher, err := NewFloatRange(lower, upper)
	if err != nil {
		panic(err)
	}

	return matcher
}

// NewFloatRange returns a pointer to a FloatMatcher that matches the Float
// range [lower, upper] (inclusive on both ends). An error is returned if the
// upper bound is less than the lower bound or either of the values is NaN.
func NewFloatRange(lower, upper float64) (*FloatMatcher, error) {
	rang, err := NewRangeMatcher(&lower, &upper)
	if err != nil {
		return nil, err
	}

	return &FloatMatcher{rang}, nil
}

// NewFloatUpperBound returns a pointer to a FloatMatcher that matches any
// Float less than or equal to the specified upper bound. It panics if the
// provided value is NaN.
func NewFloatUpperBound(upper float64) *FloatMatcher {
	rang, err := NewRangeMatcher(nil, &upper)
	if err != nil {
		// won't get here as one of the bounds is nil
		panic(err)
	}

	return &FloatMatcher{rang}
}

// NewFloatLowerBound returns a pointer to a FloatMatcher that matches any
// Float greater than or equal to the specified lower bound. It panics if
// the provided value is NaN.
func NewFloatLowerBound(lower float64) *FloatMatcher {
	rang, err := NewRangeMatcher(&lower, nil)
	if err != nil {
		// won't get here as one of the bounds is nil
		panic(err)
	}

	return &FloatMatcher{rang}
}

var relativeTolerance = 0.1e-09
var absoluteTolerance = 0.0

// SetRelativeTolerance sets the maximum difference for two float64's being
// considered "close", relative to the magnitude of their values. The default
// tolerance is 1e-09.
func SetRelativeTolerance(value float64) error {
	if value < 0.0 {
		return errors.New("tolerance cannot be negative")
	}

	relativeTolerance = value
	return nil
}

// SetRelativeTolerance sets the maximum difference for two float64's being
// considered "close", regardless of the magnitude of their values. The default
// tolerance is 0.0 (i.e. absolute tolerance is not used by default).
func SetAbsoluteTolerance(value float64) error {
	if value < 0.0 {
		return errors.New("tolerance cannot be negative")
	}

	absoluteTolerance = value
	return nil
}

// Float64Compare returns 0 if the two values are within set tolerances (see
// SetRelativeTolerance and SetAbsoluteTolerance) of each other, otherwise -1
// if lhs < rhs, and 1 if lhs > rhs.
func Float64Compare(lhs, rhs float64) (int, error) {
	if areClose(lhs, rhs) {
		return 0, nil
	} else if lhs < rhs {
		return -1, nil
	}

	return 1, nil
}

func areClose(a, b float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}

	// this also handles infinities
	if a == b {
		return true
	}

	diff := math.Abs(a - b)

	return (diff <= absoluteTolerance ||
		diff <= relativeTolerance*math.Abs(a) ||
		diff <= relativeTolerance*math.Abs(b))
}

func unmarshalBigIntFromJSON(data []byte) (*big.Int, error) {
	// for some reason, json.Number will happily accept a string that
	// contains a numeric value, so we need to explicitly guard for strings
	if strings.Contains(string(data), `"`) {
		return nil, errors.New("not a JSON number: string")
	}

	var jsonNumber json.Number
	if err := json.Unmarshal(data, &jsonNumber); err != nil {
		return nil, fmt.Errorf("not a JSON number: %w", err)
	}

	numberText := jsonNumber.String()
	if bigInt, ok := new(big.Int).SetString(numberText, 10); ok {
		return bigInt, nil
	}

	bigRat, ok := new(big.Rat).SetString(numberText)
	if !ok {
		// note: I don't know if it is actually possible to hit this,
		// since we've already parsed json.Number.
		return nil, fmt.Errorf("invalid JSON number: %q", numberText)
	}
	if !bigRat.IsInt() {
		return nil, fmt.Errorf("JSON number %q has a non-zero fractional part", numberText)
	}

	return bigRat.Num(), nil
}

func init() {
	if err := RegisterMatcherDescriptor(
		IntegerCompare,
		func(Integer) error { return nil },
		MatcherTypeUnit,
		MatcherTypeSet,
		MatcherTypeRange,
	); err != nil {
		panic(err)
	}

	if err := RegisterMatcherDescriptor(
		Float64Compare,
		func(value float64) error {
			if math.IsNaN(value) {
				return errors.New("NaN")
			}

			if math.Inf(1) == value || math.Inf(-1) == value {
				return errors.New("Inf")
			}

			return nil
		},
		MatcherTypeUnit,
		MatcherTypeSet,
		MatcherTypeRange,
	); err != nil {
		panic(err)
	}
}
