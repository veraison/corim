// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// minimalist unsigned-corim that embeds comid test template
	//go:embed testcases/unsigned-good-corim.cbor
	testGoodUnsignedCorimCBOR []byte

	// comid entity and unsigned-corim are extended
	//go:embed testcases/unsigned-corim-with-extensions.cbor
	testUnsignedCorimWithExtensionsCBOR []byte

	// comid entity and unsigned-corim are extended
	//go:embed testcases/signed-good-corim.cbor
	testGoodSignedCorimCBOR []byte

	// comid entity and unsigned-corim are extended
	//go:embed testcases/signed-corim-with-extensions.cbor
	testSignedCorimWithExtensionsCBOR []byte

	//go:embed testcases/corim.json
	testUnsignedCorimJSON []byte

	//go:embed testcases/corim-ext.json
	testUnsignedCorimWithExtensionsJSON []byte

	//go:embed testcases/comid.json
	testComidJSON []byte

	//go:embed testcases/comid-ext.json
	testComidWithExtensionsJSON []byte
)

func assertCoRIMEq(t *testing.T, expected []byte, actual []byte, msgAndArgs ...interface{}) bool {
	var expectedCoRIM, actualCoRIM *UnsignedCorim

	if err := dm.Unmarshal(expected, &expectedCoRIM); err != nil {
		return assert.Fail(t, fmt.Sprintf(
			"Expected value ('%s') is not valid UnsignedCorim: '%s'",
			expected, err.Error()), msgAndArgs...)
	}

	if err := dm.Unmarshal(actual, &actualCoRIM); err != nil {
		return assert.Fail(t, fmt.Sprintf(
			"actual value ('%s') is not valid UnsignedCorim: '%s'",
			actual, err.Error()), msgAndArgs...)
	}

	if !assert.EqualValues(t, expectedCoRIM.ID, actualCoRIM.ID, msgAndArgs...) {
		return false
	}

	if !assert.EqualValues(t, expectedCoRIM.DependentRims,
		actualCoRIM.DependentRims, msgAndArgs...) {
		return false
	}

	if !assert.EqualValues(t, expectedCoRIM.Profile, actualCoRIM.Profile, msgAndArgs...) {
		return false
	}

	if !assert.EqualValues(t, expectedCoRIM.RimValidity,
		actualCoRIM.RimValidity, msgAndArgs...) {
		return false
	}

	if !assert.EqualValues(t, expectedCoRIM.Entities, actualCoRIM.Entities, msgAndArgs...) {
		return false
	}

	if len(expectedCoRIM.Tags) != len(actualCoRIM.Tags) {
		allMsgAndArgs := []interface{}{len(expectedCoRIM.Tags), len(actualCoRIM.Tags)}
		allMsgAndArgs = append(allMsgAndArgs, msgAndArgs...)
		return assert.Fail(t, fmt.Sprintf(
			"Unexpected number of Tags: expected %d, actual %d", allMsgAndArgs...))
	}

	for i, expectedTag := range expectedCoRIM.Tags {
		actualTag := actualCoRIM.Tags[i]

		if !assertCBOREq(t, expectedTag.Content, actualTag.Content, msgAndArgs...) {
			return false
		}
	}

	return true
}

func assertCBOREq(t *testing.T, expected []byte, actual []byte, msgAndArgs ...interface{}) bool {
	var expectedCBOR, actualCBOR interface{}

	if err := dm.Unmarshal(expected, &expectedCBOR); err != nil {
		return assert.Fail(t, fmt.Sprintf(
			"Expected value ('%s') is not valid cbor.\nCBOR parsing error: '%s'",
			expected, err.Error()), msgAndArgs...)
	}

	if err := dm.Unmarshal(actual, &actualCBOR); err != nil {
		return assert.Fail(t, fmt.Sprintf(
			"Input ('%s') needs to be valid cbor.\nCBOR parsing error: '%s'",
			actual, err.Error()), msgAndArgs...)
	}

	return assert.Equal(t, expectedCBOR, actualCBOR, msgAndArgs...)
}
