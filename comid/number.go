// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

var ErrNumberTypeMismatch = errors.New("attempting to compare Number's of different types")

var (
	// 2^64 - 1
	cborUintMax = new(big.Int).SetUint64(math.MaxUint64)
	// -2^64
	cborNintMin = new(big.Int).Neg(new(big.Int).Lsh(big.NewInt(1), 64))
)

// NumberFromInt creates a new Number from the provided int.
func NumberFromInt(value int) Number {
	return NumberFromInt64(int64(value))
}

// NumberFromInt64 creates a new Number from the provided int64.
func NumberFromInt64(value int64) Number {
	var num big.Int
	num.SetInt64(value)
	return Number{&num}
}

// NumberFromUint creates a new Number from the provided uint.
func NumberFromUint(value uint) Number {
	return NumberFromUint64(uint64(value))
}

// NumberFromUint64 creates a new Number from the provided uint64.
func NumberFromUint64(value uint64) Number {
	var num big.Int
	num.SetUint64(value)
	return Number{&num}
}

// NumberFromBig creates a new Number from the provided big.Int. If the
// provided value does not fit inside a CBOR integer, an error is returned.
func NumberFromBig(value *big.Int) (Number, error) {
	if value == nil {
		return Number{}, errors.New("nil")
	}

	if value.Cmp(cborUintMax) == 1 || value.Cmp(cborNintMin) == -1 {
		return Number{}, errors.New("outside CBOR int range")
	}

	return Number{new(big.Int).Set(value)}, nil
}

// NumberFromFloat64 creates a new Number from the provided float64.
func NumberFromFloat64(value float64) Number {
	return Number{value}
}

// Number encapsulates a numeric value represented as either int64 or float64.
type Number struct {
	value any
}

// Value returns the underlying value.
func (o Number) Value() any {
	return o.value
}

// IsInt returns true if the receiver encapsulates an integer.
func (o Number) IsInt() bool {
	_, ok := o.value.(*big.Int)
	return ok
}

// IsInt64 returns true if the receiver encapsulates an integer that can fit
// inside int64.
func (o Number) IsInt64() bool {
	b, ok := o.value.(*big.Int)
	if !ok {
		return false
	}

	return b.IsInt64()
}

// Uint64 returns the underlying integer value as a uint64. If the value is
// negative or a float, 0 is returned instead.
func (o Number) Uint64() uint64 {
	b, ok := o.value.(*big.Int)
	if !ok || b.Cmp(big.NewInt(0)) <= 0 {
		return 0
	}

	if !b.IsUint64() {
		// should not get here as cannot create a Number from bit.Int
		// that is larger than a uint64.
		panic("too big")
	}

	return b.Uint64()
}

// Int64 returns the underlying value as an int64. If the underlying value is a
// float or does not fit inside int64, 0 is returned instead.
func (o Number) Int64() int64 {
	b, ok := o.value.(*big.Int)
	if !ok || !b.IsInt64() {
		return 0
	}

	return b.Int64()
}

// BigInt returns a pointer to the underlying big.Int value, or the zero value
// for bit.Int if the underlying value is a float.
func (o Number) BigInt() *big.Int {
	ret, _ := o.value.(*big.Int)
	return ret
}

// IsFloat64 returns true if the receiver encapsulates a float64.
func (o Number) IsFloat64() bool {
	_, ok := o.value.(float64)
	return ok
}

// Float64 returns the underlying float64 value. If the underlying value is not
// a float64, 0.0 is returned instead.
func (o Number) Float64() float64 {
	ret, _ := o.value.(float64)
	return ret
}

// Compare returns -1 if the receiver is less than the provided number, 1 if
// the receiver is greater than provided Number, or 0 if the two Number's are
// the same. (The behavior is equivalent to cmp.Compare() behavior for
// cmd.Ordered.)
// For floating point values, Numbers are considered equal if they within
// predefined tolerances (see SetRelativeTolerance() and
// SetAbsoluteTolerance()).
// If one number is and int and the other a float, or if one or more are
// uninitialized, an error is returned.
func (o Number) Compare(other Number) (int, error) {
	switch thisValue := o.value.(type) {
	case *big.Int:
		switch otherValue := other.value.(type) {
		case *big.Int:
			return thisValue.Cmp(otherValue), nil
		default:
			return 0, ErrNumberTypeMismatch
		}
	case float64:
		switch otherValue := other.value.(type) {
		case float64:
			if areClose(thisValue, otherValue) {
				return 0, nil
			} else if thisValue < otherValue {
				return -1, nil
			} else {
				return 1, nil
			}
		default:
			return 0, ErrNumberTypeMismatch
		}
	default:
		return 0, ErrNumberTypeMismatch
	}
}

func (o Number) MarshalCBOR() ([]byte, error) {
	if o.value != nil {
		return em.Marshal(o.value)
	} else {
		return em.Marshal(0)
	}
}

func (o *Number) UnmarshalCBOR(data []byte) error {
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
	case 7:
		var value float64
		if err := dm.Unmarshal(data, &value); err != nil {
			return err
		}

		o.value = value
	default:
		return fmt.Errorf("unexpected major type %d for Number", majorType)
	}

	return nil
}

func (o Number) MarshalJSON() ([]byte, error) {
	if o.value != nil {
		return json.Marshal(o.value)
	} else {
		return json.Marshal(0)
	}
}

func (o *Number) UnmarshalJSON(data []byte) error {
	var value json.Number
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	b, ok := new(big.Int).SetString(value.String(), 10)
	if ok && b.Cmp(cborUintMax) <= 0 && b.Cmp(cborNintMin) >= 0 {
		o.value = b
		return nil
	}

	f, err := strconv.ParseFloat(value.String(), 64)
	if err != nil {
		return err
	}

	o.value = f
	return nil
}

// NumberMatcher matches Number values to other Number values, sets of Number
// values, or Number ranges.
type NumberMatcher = Matcher[Number]

// NewNumberMatcher returns a pointer to a NumberMatcher containing the
// provided values. If No tail values are specified, a unit matcher is created
// using head, otherwise a set matcher containing head and tail values is
// created.
func NewNumberMatcher(head Number, tail ...Number) *NumberMatcher {
	if len(tail) == 0 {
		return &NumberMatcher{NewUnitMatcher(head)}
	}

	values := append([]Number{head}, tail...)

	return &NumberMatcher{MustNewSetMatcher(values...)}
}

// MustNewNumberRange returns a pointer to a NumberMatcher that matches the
// Number range [lower, upper] (inclusive on both ends). It panics if the upper
// bound is less than the lower bound.
func MustNewNumberRange(lower, upper Number) *NumberMatcher {
	matcher, err := NewNumberRange(lower, upper)
	if err != nil {
		panic(err)
	}

	return matcher
}

// NewNumberRange returns a pointer to a NumberMatcher that matches the Number
// range [lower, upper] (inclusive on both ends). An error is returned if the
// upper bound is less than the lower bound.
func NewNumberRange(lower, upper Number) (*NumberMatcher, error) {
	rang, err := NewRangeMatcher(&lower, &upper)
	if err != nil {
		return nil, err
	}

	return &NumberMatcher{rang}, nil
}

// NewNumberUpperBound returns a pointer to a NumberMatcher that matches any
// Number less than or equal to the specified upper bound.
func NewNumberUpperBound(upper Number) *NumberMatcher {
	rang, err := NewRangeMatcher(nil, &upper)
	if err != nil {
		// won't get here as one of the bounds is nil
		panic(err)
	}

	return &NumberMatcher{rang}
}

// NewNumberLowerBound returns a pointer to a NumberMatcher that matches any
// Number greater than or equal to the specified lower bound.
func NewNumberLowerBound(lower Number) *NumberMatcher {
	rang, err := NewRangeMatcher(&lower, nil)
	if err != nil {
		// won't get here as one of the bounds is nil
		panic(err)
	}

	return &NumberMatcher{rang}
}

// NumberCompare returns -1 if lhs < rhs, 1 if lhs > rhs, and 0 if lhs == rhs.
func NumberCompare(lhs, rhs Number) (int, error) {
	return lhs.Compare(rhs)
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

func init() {
	if err := RegisterMatcherDescriptor(
		NumberCompare,
		MatcherTypeUnit,
		MatcherTypeSet,
		MatcherTypeRange,
	); err != nil {
		panic(err)
	}
}
