// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuery_Valid_empty_query(t *testing.T) {
	tv := &Query{}

	actual := tv.Valid()

	expected := "invalid environment selector: non-empty<> constraint violation"

	assert.EqualError(t, actual, expected)
}

func TestQuery_Valid_invalid_timestamp(t *testing.T) {
	tv := exampleClassQuery(t)

	require.NotNil(t, tv.SetTimestamp(time.Time{}))

	actual := tv.Valid()

	expected := "timestamp not set"

	assert.EqualError(t, actual, expected)
}
