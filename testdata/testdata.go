// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
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
