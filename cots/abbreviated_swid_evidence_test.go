// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cots

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/swid"
)

func TestAbbreviatedSwidTag_Valid_WithEvidence_Success(t *testing.T) {
	validDate := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)

	tag, err := NewTag("test-tag-id", "Test Software", "1.0.0")
	assert.NoError(t, err)

	// Add required entity
	entity := swid.Entity{
		EntityName: "Test Inc.",
	}
	err = entity.SetRoles(swid.RoleTagCreator)
	assert.NoError(t, err)
	tag.Entities = append(tag.Entities, entity)

	// Add valid Evidence
	evidence := &swid.Evidence{
		DeviceID: "test-device-123",
		Date:     validDate,
	}
	tag.Evidence = evidence

	err = tag.Valid()
	assert.NoError(t, err, "Tag with valid Evidence should pass validation")
}

func TestAbbreviatedSwidTag_Valid_WithInvalidEvidence(t *testing.T) {
	tag, err := NewTag("test-tag-id", "Test Software", "1.0.0")
	assert.NoError(t, err)

	// Add required entity
	entity := swid.Entity{
		EntityName: "Test Inc.",
	}
	err = entity.SetRoles(swid.RoleTagCreator)
	assert.NoError(t, err)
	tag.Entities = append(tag.Entities, entity)

	// Add invalid Evidence (missing required fields)
	evidence := &swid.Evidence{
		// Missing DeviceID and Date
	}
	tag.Evidence = evidence

	err = tag.Valid()
	assert.Error(t, err, "Tag with invalid Evidence should fail validation")
	assert.Contains(t, err.Error(), "evidence validation failed")
}

func TestAbbreviatedSwidTag_Valid_WithoutEvidence(t *testing.T) {
	tag, err := NewTag("test-tag-id", "Test Software", "1.0.0")
	assert.NoError(t, err)

	// Add required entity
	entity := swid.Entity{
		EntityName: "Test Inc.",
	}
	err = entity.SetRoles(swid.RoleTagCreator)
	assert.NoError(t, err)
	tag.Entities = append(tag.Entities, entity)

	// Evidence is nil - should still pass validation
	err = tag.Valid()
	assert.NoError(t, err, "Tag without Evidence should pass validation")
}
