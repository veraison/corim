// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ReferenceValue(t *testing.T) {
	rv := ValueTriple{}
	err := rv.Valid()
	assert.EqualError(t, err, "environment validation failed: environment must not be empty")

	id, err := uuid.NewUUID()
	require.NoError(t, err)
	rv.Environment.Instance = MustNewUUIDInstance(id)
	err = rv.Valid()
	assert.EqualError(t, err, "measurements validation failed: no measurement entries")
}
