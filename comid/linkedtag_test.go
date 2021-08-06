// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/swid"
)

func TestLinkedTag_Valid_default_value(t *testing.T) {
	actual := NewLinkedTag()
	require.NotNil(t, actual)

	expected := "tag-id must be set in linked-tag"

	assert.EqualError(t, actual.Valid(), expected)
}

func TestLinkedTag_Valid_no_rel(t *testing.T) {
	tagID := swid.NewTagID(TestTagID)
	require.NotNil(t, tagID)

	actual := NewLinkedTag().SetLinkedTag(*tagID)
	require.NotNil(t, actual)

	expected := "rel validation failed: rel is unset"

	assert.EqualError(t, actual.Valid(), expected)
}

func TestLinkedTag_Valid_all_set(t *testing.T) {
	tagID := swid.NewTagID(TestTagID)
	require.NotNil(t, tagID)

	actual := NewLinkedTag().
		SetLinkedTag(*tagID).
		SetRel(RelReplaces)
	require.NotNil(t, actual)

	assert.Nil(t, actual.Valid())
}

func TestLinkedTags_Valid_bad_entry(t *testing.T) {
	emptyLT := NewLinkedTag()
	require.NotNil(t, emptyLT)

	actual := NewLinkedTags().
		AddLinkedTag(*emptyLT)
	require.NotNil(t, actual)

	expected := "invalid linked-tag entry at index 0: tag-id must be set in linked-tag"

	assert.EqualError(t, actual.Valid(), expected)
}
