// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery_class_ToSQL(t *testing.T) {
	tv := exampleClassQuery(t)

	expected := `SELECT * FROM reference-values WHERE ( class-id = "iZl4ZVY=" AND class-vendor = "Example Vendor" AND class-model = "Example Model" ) OR ( class-id = "31fb5abf-023e-4992-aa4e-95f9c1503bfa" )`

	actual, err := tv.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestQuery_instance_ToSQL(t *testing.T) {
	tv := exampleInstanceQuery(t)

	expected := `SELECT * FROM endorsed-values WHERE ( instance-id = "At6tvu/erQ==" ) OR ( instance-id = "iZl4ZVY=" )`

	actual, err := tv.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	fmt.Println(actual)
}

func TestQuery_group_ToSQL(t *testing.T) {
	tv := exampleGroupQuery(t)

	expected := `SELECT * FROM trust-anchors WHERE ( group-id = "iZl4ZVY=" ) OR ( group-id = "31fb5abf-023e-4992-aa4e-95f9c1503bfa" )`

	actual, err := tv.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	fmt.Println(actual)
}
