// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUUID_JSON(t *testing.T) {
	val := TaggedUUID(TestUUID)
	expected := fmt.Sprintf("%q", val.String())

	out, err := val.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, expected, string(out))

	var outUUID TaggedUUID

	err = outUUID.UnmarshalJSON(out)
	require.NoError(t, err)
	assert.Equal(t, val, outUUID)
}
