// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/corim/comid"
)

func TestEnvironmentSelector_Valid_mixed_fail(t *testing.T) {
	tv := badExampleMixedSelector(t)

	err := tv.Valid()
	assert.EqualError(t, err, "only one selector type is allowed")
}

func TestEnvironmentSelector_Valid_empty_fail(t *testing.T) {
	tv := badExampleEmptySelector(t)

	err := tv.Valid()
	assert.EqualError(t, err, "non-empty<> constraint violation")
}

func TestStatefulClass_MarshalCBOR_invalid_missing_mandatory(t *testing.T) {
	tv := StatefulClass{}
	_, err := tv.MarshalCBOR()
	assert.EqualError(t, err, "mandatory field class not set")
}

func TestStatefulClass_MarshalCBOR_valid_full(t *testing.T) {
	tv := StatefulClass{
		Class:        comid.NewClassUUID(comid.TestUUID),
		Measurements: comid.NewMeasurements().Add(comid.MustNewUintMeasurement(uint(1))),
	}
	_, err := tv.MarshalCBOR()
	assert.NoError(t, err)
}

func TestStatefulClass_UnmarshalCBOR_invalid_eof(t *testing.T) {
	tv := comid.MustHexDecode(t, "")
	var actual StatefulClass
	err := actual.UnmarshalCBOR(tv)
	assert.EqualError(t, err, "unmarshaling StatefulClass: CBOR decoding: EOF")
}

func TestStatefulClass_UnmarshalCBOR_invalid_empty_array(t *testing.T) {
	tv := comid.MustHexDecode(t, "80")
	var actual StatefulClass
	err := actual.UnmarshalCBOR(tv)
	assert.EqualError(t, err, "unmarshaling StatefulClass: wrong number of entries (0) in the array")
}
