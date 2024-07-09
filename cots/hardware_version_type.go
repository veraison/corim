// Copyright 2023-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package cots

import "github.com/veraison/swid"

type HardwareVersionType struct {
	Version string

	Scheme swid.VersionScheme
}
