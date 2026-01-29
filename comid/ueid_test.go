// Copyright 2024-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewTaggedUEID(t *testing.T) {
	ueid := UEID(TestUEID)
	tagged := TaggedUEID(TestUEID)
	bytes := MustHexDecode(t, TestUEIDString)

	for _, v := range []any{
		TestUEID,
		&TestUEID,
		ueid,
		&ueid,
		tagged,
		&tagged,
		bytes,
		base64.StdEncoding.EncodeToString(bytes),
	} {
		ret, err := NewTaggedUEID(v)
		require.NoError(t, err)
		assert.Equal(t, []byte(TestUEID), ret.Bytes())
	}
}

func TestUEID_Empty(t *testing.T) {
	var empty UEID
	assert.True(t, empty.Empty())

	nonEmpty := UEID(TestUEID)
	assert.False(t, nonEmpty.Empty())
}
