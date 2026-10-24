// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"

	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/corim/encoding"
)

const (
	RangeTagNumber = 565
	SetTagNumber   = 566

	UnitType  = "unit"
	SetType   = "set"
	RangeType = "range"
)

// MatcherType indicates the method by which a Matcher matches values.
type MatcherType struct {
	value string
}

// String returns the string representation of the MatcherType.
func (o MatcherType) String() string {
	return o.value
}

var (
	// MatcherTypeUnit matches values by comparing them to a single
	// reference values.
	MatcherTypeUnit = MatcherType{UnitType}
	// MatcherTypeRange matches values by checking wither the value within
	// a specific range.
	MatcherTypeRange = MatcherType{RangeType}
	// MatcherTypeSet matches values by checking whether the value is
	// one of a specific set of values
	MatcherTypeSet = MatcherType{SetType}
)

// Compare returns -1 if lhs is less than rhs, 1 if lhs is greater than rhs,
// and 0 if they are equal. This is a generalization of cmp.Compare[T
// cmp.Ordered].
type Compare[T any] func(lhs T, rhs T) (int, error)

// RegisterMatcherDescriptor registers an equality function for type T. Each
// type must specify two things: a compare function for that type, and which
// matchers are supported by the type (note: MatcherTypeUnit must always be
// specified).
func RegisterMatcherDescriptor[T any](compare Compare[T], supportedVariants ...MatcherType) error {
	key := typeName[T]()
	if _, ok := matcherDescriptorRegistry[key]; ok {
		return fmt.Errorf("equal for %T already registered", key)
	}

	if !slices.Contains(supportedVariants, MatcherTypeUnit) {
		return errors.New("unit type must be in supported variants")
	}

	matcherDescriptorRegistry[key] = matcherDescriptor[T]{
		Compare:           compare,
		SupportedVariants: supportedVariants,
	}

	return nil
}

// TextMatcher matches string values to other string values or sets of string
// values.
type TextMatcher = Matcher[string]

// NewTextMatcher returns a pointer to a TextMatcher containing the
// provided values. If No tail values are specified, a unit matcher is created
// using head, otherwise a set matcher containing head and tail values is
// created.
func NewTextMatcher(head string, tail ...string) *TextMatcher {
	if len(tail) == 0 {
		return &TextMatcher{NewUnitMatcher(head)}
	}

	values := append([]string{head}, tail...)

	return &TextMatcher{MustNewSetMatcher(values...)}
}

// BytesMatcher matches []byte values to other []byte values or sets of []byte
// values.
type BytesMatcher = Matcher[[]byte]

// NewBytesMatcher returns a pointer to a BytesMatcher containing the
// provided values. If No tail values are specified, a unit matcher is created
// using head, otherwise a set matcher containing head and tail values is
// created.
func NewBytesMatcher(head []byte, tail ...[]byte) *BytesMatcher {
	if len(tail) == 0 {
		return &BytesMatcher{NewUnitMatcher(head)}
	}

	values := append([][]byte{head}, tail...)

	return &BytesMatcher{
		value: MustNewSetMatcher(values...),
	}
}

// MatcherVariant is the interface implemented by all Matcher[T] variants.
type MatcherVariant[T any] interface {
	// MatcherType returns the type of the receiver.
	MatcherType() MatcherType
	// MatchesValue returns true if the value is matched by the receiver.
	MatchesValue(value T) bool
	// UnitValue return the underlying value if the matcher is a Unit[T],
	// otherwise returns the zero value for T.
	UnitValue() T
	// TagNumber returns the CBOR tag number associated with this matcher.
	// Unit[T], which does not have an associated tag, returns 0.
	TagNumber() uint64
}

// Matcher allows matching a value to a reference value or set.
type Matcher[T any] struct {
	value MatcherVariant[T]
}

// Type returns the MatcherType for the underlying variant. This panics
// if the receiver is not initialized.
func (o Matcher[T]) Type() MatcherType {
	return o.value.MatcherType()
}

// MatchesValue returns true if the receiver matches the provided value.
func (o Matcher[T]) MatchesValue(value T) bool {
	if o.value == nil {
		return false
	}

	return o.value.MatchesValue(value)
}

// CompareAgainstReference returns true if the receiver matches provided
// reference. If the reference is a unit, then the receiver matches if its
// value is equal to the reference value. If the reference value is a set, the
// receiver matches if its value is inside the reference set. if the reference
// value is a range, the receiver matches if it is greater or equal to the
// lower bound and less than or equal to the upper bound of that range. If
// receiver is not a unit, false is always returned.
func (o Matcher[T]) CompareAgainstReference(r *Matcher[T]) bool {
	if r == nil || o.value == nil || o.Type() != MatcherTypeUnit {
		// value to be compared to reference (and was obtained from
		// evidence) is expected to be a unit
		return false
	}

	return r.MatchesValue(o.value.UnitValue())
}

func (o Matcher[T]) MarshalCBOR() ([]byte, error) {
	if o.value == nil {
		return nil, errors.New("uninitialized matcher")
	}

	tagNumber := o.value.TagNumber()
	if tagNumber == 0 {
		return em.Marshal(o.value)
	}

	toMarshal := cbor.Tag{
		Number:  tagNumber,
		Content: o.value,
	}

	return em.Marshal(toMarshal)
}

func (o *Matcher[T]) UnmarshalCBOR(data []byte) error {
	if len(data) < 1 {
		return errors.New("truncated input")
	}

	var value MatcherVariant[T]

	majorType := (data[0] & 0xe0) >> 5
	if majorType == 6 {
		var tag cbor.RawTag
		if err := dm.Unmarshal(data, &tag); err != nil {
			return err
		}

		switch tag.Number {
		case SetTagNumber:
			var set SetMatcher[T]
			if err := dm.Unmarshal(tag.Content, &set); err != nil {
				return err
			}

			value = set
		case RangeTagNumber:
			var rang RangeMatcher[T]
			if err := dm.Unmarshal(tag.Content, &rang); err != nil {
				return err
			}

			value = rang
		default:
			// if tag number is unrecognized, it's possible that T
			// itself is tagged.
			var unit UnitMatcher[T]
			if err := dm.Unmarshal(data, &unit); err != nil {
				return err
			}

			value = unit
		}
	} else {
		// not a set -> unit
		var unit UnitMatcher[T]
		if err := dm.Unmarshal(data, &unit); err != nil {
			return err
		}

		value = unit
	}

	typ := value.MatcherType()
	if !isMatcherSupported[T](typ) {
		return fmt.Errorf("invalid matcher type %q for %T", typ.String(), value)
	}

	o.value = value

	return nil
}

func (o Matcher[T]) MarshalJSON() ([]byte, error) {
	if o.value == nil {
		return nil, errors.New("uninitialized matcher")
	}

	encodedValue, err := json.Marshal(o.value)
	if err != nil {
		return nil, err
	}

	toMarshal := encoding.TypeAndValue{
		Type:  o.value.MatcherType().String(),
		Value: encodedValue,
	}

	return json.Marshal(toMarshal)
}

func (o *Matcher[T]) UnmarshalJSON(data []byte) error {
	var tnv encoding.TypeAndValue
	if err := json.Unmarshal(data, &tnv); err != nil {
		return err
	}

	var value MatcherVariant[T]

	switch tnv.Type {
	case SetType:
		var set SetMatcher[T]
		if err := json.Unmarshal(tnv.Value, &set); err != nil {
			return err
		}

		value = set
	case RangeType:
		var rang RangeMatcher[T]
		if err := json.Unmarshal(tnv.Value, &rang); err != nil {
			return err
		}

		value = rang
	case UnitType:
		var unit UnitMatcher[T]
		if err := json.Unmarshal(tnv.Value, &unit); err != nil {
			return err
		}

		value = unit
	default:
		return fmt.Errorf("unexpected matcher type: %q", tnv.Type)
	}

	typ := value.MatcherType()
	if !isMatcherSupported[T](typ) {
		return fmt.Errorf("invalid matcher type %q for %T", typ.String(), value)
	}

	o.value = value

	return nil
}

// UnitMatcher is a matcher variant that matches one specific value.
type UnitMatcher[T any] struct {
	value T
}

// NewUnitMatcher returns a new Unit that matches the provided value.
func NewUnitMatcher[T any](value T) UnitMatcher[T] {
	return UnitMatcher[T]{value}
}

func (o UnitMatcher[T]) MatcherType() MatcherType {
	return MatcherTypeUnit
}

func (o UnitMatcher[T]) MatchesValue(value T) bool {
	return equal[T]()(o.value, value)
}

func (o UnitMatcher[T]) UnitValue() T {
	return o.value
}

func (o UnitMatcher[T]) TagNumber() uint64 {
	return 0
}

func (o UnitMatcher[T]) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.value)
}

func (o *UnitMatcher[T]) UnmarshalCBOR(data []byte) error {
	var value T
	if err := dm.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = value
	return nil
}

func (o UnitMatcher[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

func (o *UnitMatcher[T]) UnmarshalJSON(data []byte) error {
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = value
	return nil
}

// SetMatcher represents a set of values of any type. Each element of a set is unique.
type SetMatcher[T any] struct {
	values []T
}

// MustNewSetMatcher is a version of NewSet that panics instead of returning an error.
func MustNewSetMatcher[T any](values ...T) SetMatcher[T] {
	ret, err := NewSetMatcher(values...)
	if err != nil {
		panic(err)
	}

	return ret
}

// NewSetMatcher returns a new set containing the provided values. An error is
// returned if fewer than two values are provided, or if there are duplicates
// among the values.
func NewSetMatcher[T any](values ...T) (SetMatcher[T], error) {
	var ret SetMatcher[T]
	for _, value := range values {
		if !ret.Add(value) {
			return SetMatcher[T]{}, fmt.Errorf("duplicate: %v", value)
		}
	}

	if ret.Len() < 2 {
		return SetMatcher[T]{}, errors.New("must specify at least two distinct elements")
	}

	return ret, nil
}

// Contains returns true of the receiver contains the provided value.
func (o SetMatcher[T]) Contains(value T) bool {
	equal := equal[T]()

	for _, setValue := range o.values {
		if equal(setValue, value) {
			return true
		}
	}

	return false
}

// Add ensures that the receiver contains the provided value. If the value is
// added to the receiver, true is returned. If the receiver already contains
// the value, false is returned.
func (o *SetMatcher[T]) Add(value T) bool {
	if o.Contains(value) {
		return false
	}

	o.values = append(o.values, value)
	return true
}

// Values returns a slice of the values contained in the set.
func (o SetMatcher[T]) Values() []T {
	return o.values
}

// Len returns the number of values in the set
func (o SetMatcher[T]) Len() int {
	return len(o.values)
}

func (o SetMatcher[T]) MatcherType() MatcherType {
	return MatcherTypeSet
}

func (o SetMatcher[T]) MatchesValue(value T) bool {
	return o.Contains(value)
}

func (o SetMatcher[T]) UnitValue() T {
	var zero T
	return zero
}

func (o SetMatcher[T]) TagNumber() uint64 {
	return SetTagNumber
}

func (o SetMatcher[T]) MarshalCBOR() ([]byte, error) {
	if o.Len() < 2 {
		return nil, errors.New("set must contain at least two entries")
	}

	return em.Marshal(o.values)
}

func (o *SetMatcher[T]) UnmarshalCBOR(data []byte) error {
	var values []T
	if err := dm.Unmarshal(data, &values); err != nil {
		return err
	}

	if len(values) < 2 {
		return errors.New("set must contain at least two entries")
	}

	for _, value := range values {
		if !o.Add(value) {
			return fmt.Errorf("duplicate: %v", value)
		}
	}

	return nil
}

func (o SetMatcher[T]) MarshalJSON() ([]byte, error) {
	if o.Len() < 2 {
		return nil, errors.New("set must contain at least two entries")
	}

	return json.Marshal(o.values)
}

func (o *SetMatcher[T]) UnmarshalJSON(data []byte) error {
	var values []T
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}

	if len(values) < 2 {
		return errors.New("set must contain at least two entries")
	}

	for _, value := range values {
		if !o.Add(value) {
			return fmt.Errorf("duplicate: %v", value)
		}
	}

	return nil
}

// RangeMatcher represents an inclusive range of values defined by an upper and
// a lower bound.
type RangeMatcher[T any] struct {
	lower *T
	upper *T
}

// MustNewRangeMatcher returns a RangeMatcher for the range [lower, upper]. It
// panics if lower is greater than upper.
func MustNewRangeMatcher[T any](lower, upper *T) RangeMatcher[T] {
	ret, err := NewRangeMatcher(lower, upper)
	if err != nil {
		panic(err)
	}

	return ret
}

// NewRangeMatcher returns a RangeMatcher for the range [lower, upper]. If
// lower is greater than upper, an error is returned.
func NewRangeMatcher[T any](lower, upper *T) (RangeMatcher[T], error) {
	if lower != nil && upper != nil && lessThan[T]()(*upper, *lower) {
		return RangeMatcher[T]{}, errors.New("upper bound cannot be less than lower bound")
	}

	return RangeMatcher[T]{lower: lower, upper: upper}, nil
}

func (o RangeMatcher[T]) MatcherType() MatcherType {
	return MatcherTypeRange
}

func (o RangeMatcher[T]) UnitValue() T {
	var zero T
	return zero
}

func (o RangeMatcher[T]) MatchesValue(value T) bool {
	equal := equal[T]()
	lessThan := lessThan[T]()

	if o.lower != nil && !equal(*o.lower, value) && !lessThan(*o.lower, value) {
		return false
	}

	if o.upper != nil && !equal(*o.upper, value) && !lessThan(value, *o.upper) {
		return false
	}

	return true
}

func (o RangeMatcher[T]) TagNumber() uint64 {
	return RangeTagNumber
}

func (o RangeMatcher[T]) MarshalCBOR() ([]byte, error) {
	return em.Marshal([]*T{o.lower, o.upper})
}

func (o *RangeMatcher[T]) UnmarshalCBOR(data []byte) error {
	var bounds [2]*T
	if err := dm.Unmarshal(data, &bounds); err != nil {
		return err
	}

	o.lower = bounds[0]
	o.upper = bounds[1]

	return nil
}

func (o RangeMatcher[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal([]*T{o.lower, o.upper})
}

func (o *RangeMatcher[T]) UnmarshalJSON(data []byte) error {
	var bounds [2]*T
	if err := json.Unmarshal(data, &bounds); err != nil {
		return err
	}

	o.lower = bounds[0]
	o.upper = bounds[1]

	return nil
}

type matcherDescriptor[T any] struct {
	Compare           Compare[T]
	SupportedVariants []MatcherType
}

func typeName[T any]() string {
	var zero T
	t := reflect.TypeOf(&zero).Elem()
	return t.String()
}

type matchTest[T any] func(lhs, rhs T) bool

func equal[T any]() matchTest[T] {
	key := typeName[T]()
	desc, ok := matcherDescriptorRegistry[key]
	if ok {
		return func(lhs, rhs T) bool {
			result, err := desc.(matcherDescriptor[T]).Compare(lhs, rhs)
			if err != nil {
				return false
			}
			return result == 0
		}
	}

	return func(lhs, rhs T) bool {
		return false
	}
}

func lessThan[T any]() matchTest[T] {
	key := typeName[T]()
	desc, ok := matcherDescriptorRegistry[key]
	if ok {
		return func(lhs, rhs T) bool {
			result, err := desc.(matcherDescriptor[T]).Compare(lhs, rhs)
			if err != nil {
				return false
			}
			return result < 0
		}
	}

	return func(lhs, rhs T) bool {
		return false
	}
}

func isMatcherSupported[T any](typ MatcherType) bool {
	key := typeName[T]()
	desc, ok := matcherDescriptorRegistry[key]
	if ok {
		return slices.Contains(desc.(matcherDescriptor[T]).SupportedVariants, typ)
	}

	return false
}

var matcherDescriptorRegistry = map[string]any{
	"string": matcherDescriptor[string]{
		SupportedVariants: []MatcherType{
			MatcherTypeUnit,
			MatcherTypeSet,
		},
		Compare: func(lhs, rhs string) (int, error) {
			return cmp.Compare(lhs, rhs), nil
		},
	},
	"[]uint8": matcherDescriptor[[]byte]{
		SupportedVariants: []MatcherType{
			MatcherTypeUnit,
			MatcherTypeSet,
		},
		Compare: func(lhs, rhs []byte) (int, error) {
			return bytes.Compare(lhs, rhs), nil
		},
	},
}
