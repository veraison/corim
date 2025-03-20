// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coserv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
)

func TestCoserv_ToCBOR_rv_class_simple(t *testing.T) {
	class := comid.NewClassBytes([]byte{0x00, 0x11, 0x22, 0x33}).
		SetVendor("Example Vendor").
		SetModel("Example Model")
	require.NotNil(t, class)

	envSelector := NewEnvironmentSelector().
		AddClass(*class)
	require.NotNil(t, envSelector)

	tv, err := NewCoserv(
		ArtifactTypeReferenceValues,
		`tag:example.com,2025:cc-platform#1.0.0`,
		*envSelector,
	)
	require.NoError(t, err)

	actual, err := tv.ToCBOR()
	assert.NoError(t, err)

	// {0: 2, 1: "tag:example.com,2025:cc-platform#1.0.0", 2: {0: [{0: 560(h'00112233'), 1: "Example Vendor", 2: "Example Model"}]}}
	expected := comid.MustHexDecode(t, "a300020178267461673a6578616d706c652e636f6d2c323032353a63632d706c6174666f726d23312e302e3002a10081a300d902304400112233016e4578616d706c652056656e646f72026d4578616d706c65204d6f64656c")

	assert.Equal(t, expected, actual)
}
