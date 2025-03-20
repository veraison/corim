// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
)

func TestEnvironmentSelector_class_ToSQL(t *testing.T) {
	class0 := comid.NewClassBytes([]byte{0x00, 0x11, 0x22, 0x33}).
		SetVendor("Example Vendor").
		SetModel("Example Model")
	require.NotNil(t, class0)

	class1 := comid.NewClassImplID(comid.TestImplID)
	require.NotNil(t, class1)

	tv := NewEnvironmentSelector().
		AddClass(*class0).
		AddClass(*class1)
	require.NotNil(t, tv)

	expected := `( class-id = "ABEiMw==" AND class-vendor = "Example Vendor" AND class-model = "Example Model" ) OR ( class-id = "YWNtZS1pbXBsZW1lbnRhdGlvbi1pZC0wMDAwMDAwMDE=" )`

	actual, err := tv.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestEnvironmentSelector_instance_ToSQL(t *testing.T) {
	instance0, err := comid.NewUUIDInstance("5c10e890-edd5-432d-9853-e6a0e68f3a3a")
	require.NoError(t, err)

	instance1, err := comid.NewBytesInstance([]byte{0x00, 0x11, 0x22, 0x33})
	require.NoError(t, err)

	tv := NewEnvironmentSelector().
		AddInstance(*instance0).
		AddInstance(*instance1)
	require.NotNil(t, tv)

	expected := `( instance-id = "5c10e890-edd5-432d-9853-e6a0e68f3a3a" ) OR ( instance-id = "ABEiMw==" )`

	actual, err := tv.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	fmt.Println(actual)
}
