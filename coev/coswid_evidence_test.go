// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/swid"
)

func TestCoSWIDEvidenceMap_Valid_Success(t *testing.T) {
	validDate := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)

	evidenceMap := CoSWIDEvidenceMap{
		Evidence: swid.Evidence{
			DeviceID: "test-device-123",
			Date:     validDate,
		},
	}

	err := evidenceMap.Valid()
	assert.NoError(t, err, "Valid evidence map should pass validation")
}

func TestCoSWIDEvidenceMap_Valid_WithTagID(t *testing.T) {
	validDate := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)

	evidenceMap := CoSWIDEvidenceMap{
		TagID: swid.NewTagID("test-tag-id"),
		Evidence: swid.Evidence{
			DeviceID: "test-device-123",
			Date:     validDate,
		},
	}

	err := evidenceMap.Valid()
	assert.NoError(t, err, "Valid evidence map with TagID should pass validation")
}

func TestCoSWIDEvidenceMap_Valid_InvalidEvidence(t *testing.T) {
	evidenceMap := CoSWIDEvidenceMap{
		Evidence: swid.Evidence{
			// Missing required DeviceID and Date
		},
	}

	err := evidenceMap.Valid()
	assert.Error(t, err, "Invalid evidence should fail validation")
	assert.Contains(t, err.Error(), "evidence validation failed")
}

func TestCoSWIDEvidenceMap_Valid_InvalidTagID(t *testing.T) {
	validDate := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)
	emptyTagID := &swid.TagID{} // Empty TagID - should be invalid

	evidenceMap := CoSWIDEvidenceMap{
		TagID: emptyTagID,
		Evidence: swid.Evidence{
			DeviceID: "test-device-123",
			Date:     validDate,
		},
	}

	err := evidenceMap.Valid()
	assert.Error(t, err, "Invalid TagID should fail validation")
	assert.Contains(t, err.Error(), "tagId validation failed")
}

func TestCoSWIDEvidence_Valid_Success(t *testing.T) {
	validDate := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)

	evidence := CoSWIDEvidence{
		CoSWIDEvidenceMap{
			Evidence: swid.Evidence{
				DeviceID: "test-device-1",
				Date:     validDate,
			},
		},
		CoSWIDEvidenceMap{
			Evidence: swid.Evidence{
				DeviceID: "test-device-2",
				Date:     validDate,
			},
		},
	}

	err := evidence.Valid()
	assert.NoError(t, err, "Valid evidence slice should pass validation")
}

func TestCoSWIDEvidence_Valid_EmptySlice(t *testing.T) {
	evidence := CoSWIDEvidence{}

	err := evidence.Valid()
	assert.Error(t, err, "Empty evidence slice should fail validation")
	assert.Contains(t, err.Error(), "must contain at least one entry")
}

func TestCoSWIDEvidence_Valid_InvalidEntry(t *testing.T) {
	validDate := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)

	evidence := CoSWIDEvidence{
		CoSWIDEvidenceMap{
			Evidence: swid.Evidence{
				DeviceID: "test-device-1",
				Date:     validDate,
			},
		},
		CoSWIDEvidenceMap{
			Evidence: swid.Evidence{
				// Missing required DeviceID - should fail
				Date: validDate,
			},
		},
	}

	err := evidence.Valid()
	assert.Error(t, err, "Evidence slice with invalid entry should fail validation")
	assert.Contains(t, err.Error(), "evidence[1] validation failed")
}
