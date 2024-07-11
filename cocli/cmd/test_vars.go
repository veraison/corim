// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	_ "embed"

	"github.com/veraison/corim/comid"
)

var (
	minimalCorimTemplate = []byte(`{
		"corim-id": "5c57e8f4-46cd-421b-91c9-08cf93e13cfc"
	}`)
	badCBOR = comid.MustHexDecode(nil, "ffff")
	// a "tag-id only" CoMID {1: {0: h'366D0A0A598845ED84882F2A544F6242'}}
	invalidComid = comid.MustHexDecode(nil,
		"a101a10050366d0a0a598845ed84882f2a544f6242",
	)
	//
	invalidCots = comid.MustHexDecode(nil, "a2028006a100f6")
	// note: embedded CoSWIDs are not validated {0: h'5C57E8F446CD421B91C908CF93E13CFC', 1: [505(h'deadbeef')]}
	testCorimValid = comid.MustHexDecode(nil,
		"a200505c57e8f446cd421b91c908cf93e13cfc0181d901f944deadbeef",
	)
	// {0: h'5C57E8F446CD421B91C908CF93E13CFC'}
	testCorimInvalid = comid.MustHexDecode(nil,
		"a100505c57e8f446cd421b91c908cf93e13cfc",
	)
	testMetaInvalid = []byte("{}")
	testMetaValid   = []byte(`{
		"signer": {
			"name": "ACME Ltd signing key",
			"uri": "https://acme.example"
		},
		"validity": {
			"not-before": "2021-12-31T00:00:00Z",
			"not-after": "2025-12-31T00:00:00Z"
		}
	}`)

	//go:embed testcases/ec-p256.jwk
	testECKey []byte

	//go:embed testcases/signed-corim-valid.cbor
	testSignedCorimValid []byte

	//go:embed testcases/signed-corim-invalid.cbor
	testSignedCorimInvalid []byte

	//go:embed testcases/signed-corim-valid-with-cots.cbor
	testSignedCorimValidWithCots []byte

	//go:embed testcases/psa-refval.cbor
	PSARefValCBOR []byte

	//go:embed testcases/test-comid.cbor
	testComid []byte

	//go:embed testcases/test-coswid.cbor
	testCoswid []byte

	//go:embed testcases/test-cots.cbor
	testCots []byte
)
