// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentSelector_Valid_mixed_fail(t *testing.T) {
	tv := badExampleMixedSelector(t)

	err := tv.Valid()
	assert.EqualError(t, err, "only one selector type is allowed")
}

func TestEnvironmentSelector_Valid_empty_fail(t *testing.T) {
	tv := badExampleEmptySelector(t)

	err := tv.Valid()
	assert.EqualError(t, err, "non-empty<> constraint violation")
}
