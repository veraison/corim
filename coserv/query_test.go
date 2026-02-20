// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery_Valid_empty_query(t *testing.T) {
	tv := &Query{}

	actual := tv.Valid()

	expected := "invalid environment selector: non-empty<> constraint violation"

	assert.EqualError(t, actual, expected)
}

func TestQuery_NewQuery_invalid_selector(t *testing.T) {
	actual, err := NewQuery(ArtifactTypeEndorsedValues, *NewEnvironmentSelector(), ResultTypeBoth)
	assert.Nil(t, actual)
	assert.EqualError(t, err, "invalid environment selector: non-empty<> constraint violation")
}
