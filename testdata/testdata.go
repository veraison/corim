// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

// Package testdata holds PUBLIC test-only X.509 material embedded for unit tests.
// Do not use these keys or certificates as production trust material.
package testdata

import (
	_ "embed"
)

var (
	//go:embed endEntity.der
	// EndEntityDer is an intermediate CA-signed certificate.
	EndEntityDer []byte

	//go:embed endEntity.key
	// EndEntityKey is a key used for signing CoRIMs.
	EndEntityKey []byte

	//go:embed intermediateCA.der
	// IntermediateCA is a root-signed intermediate certificate authority key certificate.
	IntermediateCA []byte

	//go:embed rootCA.der
	// RootCA is a self-signed root certificate authority key certificate.
	RootCA []byte
)
