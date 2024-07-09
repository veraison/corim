// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package extensions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testTCValue string

func (o testTCValue) String() string {
	return string(o)
}

func (o testTCValue) Valid() error {
	return nil
}

func (o testTCValue) Type() string {
	return "test_type"
}

func Test_TypeChoiceValueMarshalJSON(t *testing.T) {
	buf, err := TypeChoiceValueMarshalJSON(testTCValue("test"))
	assert.NoError(t, err)
	assert.JSONEq(t, `{"type": "test_type", "value": "test"}`, string(buf))
}
