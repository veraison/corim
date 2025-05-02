// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRawInt_NewRawIntInteger_OK(t *testing.T) {
	var val int64 = -25
	rawInt, err := NewRawInt(val, "rawIntInteger")
	assert.NoError(t, err)
	assert.NoError(t, rawInt.Valid())
}

func TestRawInt_NewRawIntRange_OK(t *testing.T) {
	minVal := int64(-25)
	maxVal := int64(-25)
	val := TaggedRawIntRange{Max: &maxVal, Min: &minVal}

	rawInt, err := NewRawInt(val, "rawIntRange")
	assert.NoError(t, err)
	assert.NoError(t, rawInt.Valid())
}

func TestRawInt_NewRawIntRange_Validity(t *testing.T) {
	minVal := int64(-25)
	maxVal := int64(-35)
	val := TaggedRawIntRange{Max: &maxVal, Min: &minVal}

	rawInt, err := NewRawInt(val, "rawIntRange")
	assert.NoError(t, err)
	assert.EqualError(t, rawInt.Valid(), fmt.Sprintf("TaggedRawIntRange: Invalid Range, Min: %d Max: %d", minVal, maxVal))
}

func TestRawInt_NewRawIntInteger_MarshalUnmarshalJSON(t *testing.T) {
	var val int64 = -65

	rawInt, err := NewRawInt(val, "rawIntInteger")
	assert.NoError(t, err)

	json, err := rawInt.MarshalJSON()
	assert.NoError(t, err)

	unmarshaledRawInt := &RawInt{}
	err = unmarshaledRawInt.UnmarshalJSON(json)
	assert.NoError(t, err)

	rawIntInteger, ok := unmarshaledRawInt.Value.(*RawIntInteger)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("%d", val), rawIntInteger.String())
}

func TestRawInt_NewRawIntInteger_MarshalUnmarshalCBOR(t *testing.T) {
	var val int64 = -85

	rawInt, err := NewRawInt(val, "rawIntInteger")
	assert.NoError(t, err)

	rawIntCbor, err := rawInt.MarshalCBOR()
	assert.NoError(t, err)

	unmarshaledRawInt := &RawInt{}
	err = unmarshaledRawInt.UnmarshalCBOR(rawIntCbor)
	assert.NoError(t, err)

	rawIntInteger, ok := unmarshaledRawInt.Value.(*RawIntInteger)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("%d", val), rawIntInteger.String())
}

func TestRawInt_NewRawIntRange_MarshalUnmarshalJSON(t *testing.T) {
	minVal := int64(-25)

	rawInt, err := NewRawInt(TaggedRawIntRange{Max: nil, Min: &minVal}, "rawIntRange")
	assert.NoError(t, err)

	rawIntJSON, err := rawInt.MarshalJSON()
	assert.NoError(t, err)

	unmarshaledRawInt := &RawInt{}
	err = unmarshaledRawInt.UnmarshalJSON(rawIntJSON)
	assert.NoError(t, err)

	rawIntRange, ok := unmarshaledRawInt.Value.(*TaggedRawIntRange)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("[%d:inf)", minVal), rawIntRange.String())
}

func TestRawInt_NewRawIntRange_MarshalUnmarshalCBOR(t *testing.T) {
	maxVal := int64(650)

	rawInt, err := NewRawInt(TaggedRawIntRange{Max: &maxVal, Min: nil}, "rawIntRange")
	assert.NoError(t, err)

	rawIntCBOR, err := rawInt.MarshalCBOR()
	assert.NoError(t, err)

	unmarshaledRawInt := &RawInt{}
	err = unmarshaledRawInt.UnmarshalCBOR(rawIntCBOR)
	assert.NoError(t, err)

	rawIntRange, ok := unmarshaledRawInt.Value.(*TaggedRawIntRange)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("(-inf:%d]", maxVal), rawIntRange.String())
}

func TestRawInt_Compare_IntegerClaimVsIntegerRef_Pass(t *testing.T) {
	claim := RawIntInteger(65)
	ref := RawIntInteger(65)

	assert.True(t, claim.CompareAgainstRefInteger(ref))
}

func TestRawInt_Compare_IntegerClaimVsIntegerRef_Fail(t *testing.T) {
	claim := RawIntInteger(65)
	ref := RawIntInteger(75)

	assert.False(t, claim.CompareAgainstRefInteger(ref))
}

func TestRawInt_Compare_IntegerClaimVsRangeRef_Pass(t *testing.T) {
	claim := RawIntInteger(65)
	refMinVal := int64(-25)
	refMaxVal := int64(75)
	ref := TaggedRawIntRange{Min: &refMinVal, Max: &refMaxVal}

	assert.True(t, claim.CompareAgainstRefRange(ref))
}

func TestRawInt_Compare_IntegerClaimVsRangeRef_Fail(t *testing.T) {
	claim := RawIntInteger(85)
	refMinVal := int64(-25)
	refMaxVal := int64(75)
	ref := TaggedRawIntRange{Min: &refMinVal, Max: &refMaxVal}

	assert.False(t, claim.CompareAgainstRefRange(ref))
}

func TestRawInt_Compare_RangeClaimVsIntegerRef_Pass(t *testing.T) {
	claimMinVal := int64(25)
	claimMaxVal := int64(25)
	claim := TaggedRawIntRange{Min: &claimMinVal, Max: &claimMaxVal}
	ref := RawIntInteger(25)

	assert.True(t, claim.CompareAgainstRefInteger(ref))
}

func TestRawInt_Compare_RangeClaimVsIntegerRef_Fail(t *testing.T) {
	claimMaxVal := int64(25)
	claim := TaggedRawIntRange{Min: nil, Max: &claimMaxVal}
	ref := RawIntInteger(25)

	assert.False(t, claim.CompareAgainstRefInteger(ref))
}

func TestRawInt_Compare_RangeClaimVsRangeRef_Pass(t *testing.T) {
	refMinVal := int64(-100)
	refMaxVal := int64(100)
	ref := TaggedRawIntRange{Min: &refMinVal, Max: &refMaxVal}

	claimMinVal := int64(-25)
	claimMaxVal := int64(75)
	claim := TaggedRawIntRange{Min: &claimMinVal, Max: &claimMaxVal}

	assert.True(t, claim.CompareAgainstRefRange(ref))
}

func TestRawInt_Compare_RangeClaimVsRangeRef_Pass_RefMinInf(t *testing.T) {
	refMaxVal := int64(100)
	ref := TaggedRawIntRange{Min: nil, Max: &refMaxVal}

	claimMinVal := int64(-25)
	claimMaxVal := int64(75)
	claim := TaggedRawIntRange{Min: &claimMinVal, Max: &claimMaxVal}

	assert.True(t, claim.CompareAgainstRefRange(ref))
}

func TestRawInt_Compare_RangeClaimVsRangeRef_Pass_RefMinInf_ClaimMinInf(t *testing.T) {
	refMaxVal := int64(100)
	ref := TaggedRawIntRange{Min: nil, Max: &refMaxVal}

	claimMaxVal := int64(75)
	claim := TaggedRawIntRange{Min: nil, Max: &claimMaxVal}

	assert.True(t, claim.CompareAgainstRefRange(ref))
}

func TestRawInt_Compare_RangeClaimVsRangeRef_Fail(t *testing.T) {
	refMinVal := int64(-100)
	refMaxVal := int64(100)
	ref := TaggedRawIntRange{Min: &refMinVal, Max: &refMaxVal}

	claimMinVal := int64(-101)
	claimMaxVal := int64(75)
	claim := TaggedRawIntRange{Min: &claimMinVal, Max: &claimMaxVal}

	assert.False(t, claim.CompareAgainstRefRange(ref))
}
