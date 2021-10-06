// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/base64"
	"encoding/hex"
	"net"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/veraison/eat"
)

var (
	TestUUIDString = "31fb5abf-023e-4992-aa4e-95f9c1503bfa"
	TestUUID       = UUID(uuid.Must(uuid.Parse(TestUUIDString)))
	TestImplID     = ImplID([32]byte{
		0x61, 0x63, 0x6d, 0x65, 0x2d, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65,
		0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2d, 0x69, 0x64, 0x2d, 0x30,
		0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31,
	})
	TestOID        = "2.5.2.8192"
	TestRegID      = "https://acme.example"
	TestMACaddr, _ = net.ParseMAC("02:00:5e:10:00:00:00:01")
	TestIPaddr     = net.ParseIP("2001:db8::68")
	TestUEIDString = "02deadbeefdead"
	TestUEID       = eat.UEID(MustHexDecode(nil, TestUEIDString))
	TestSignerID   = MustHexDecode(nil, "acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b")
	TestTagID      = "urn:example:veraison"
)

func MustHexDecode(t *testing.T, s string) []byte {
	// allow long hex string to be split over multiple lines (with soft or hard
	// tab indentation)
	m := regexp.MustCompile("[ \t\n]")
	s = m.ReplaceAllString(s, "")

	data, err := hex.DecodeString(s)
	if t != nil {
		require.Nil(t, err)
	} else if err != nil {
		panic(err)
	}
	return data
}

func b64TestImplID() string {
	var implID []byte = TestImplID[:]

	return base64.StdEncoding.EncodeToString(implID)
}

var (
	PSARefValJSONTemplate = `{
	"lang": "en-GB",
	"tag-identity": {
		"id": "43BBE37F-2E61-4B33-AED3-53CFF1428B16",
		"version": 0
	},
	"entities": [
		{
			"name": "ACME Ltd.",
			"regid": "https://acme.example",
			"roles": [ "tagCreator", "creator", "maintainer" ]
		}
	],
	"triples": {
		"reference-values": [
			{
				"environment": {
					"class": {
						"id": {
							"type": "psa.impl-id",
							"value": "YWNtZS1pbXBsZW1lbnRhdGlvbi1pZC0wMDAwMDAwMDE="
						},
						"vendor": "ACME",
						"model": "RoadRunner"
					}
				},
				"measurements": [
					{
						"key": {
							"type": "psa.refval-id",
							"value": {
								"label": "BL",
								"version": "2.1.0",
								"signer-id": "rLsRx+TaIXIFUjzkzhokWuGiOa48a/2eeHH35di66Gs="
							}
						},
						"value": {
							"digests": [
								"sha-256:h0KPxSKAPTEGXnvOPPA/5HUJZjHl4Hu9eg/eYMTPJcc="
							]
						}
					},
					{
						"key": {
							"type": "psa.refval-id",
							"value": {
								"label": "PRoT",
								"version": "1.3.5",
								"signer-id": "rLsRx+TaIXIFUjzkzhokWuGiOa48a/2eeHH35di66Gs="
							}
						},
						"value": {
							"digests": [
								"sha-256:AmOCmYm2/ZVPcrqvL8ZLwuLwHWktTecphuqAj26ZgT8="
							]
						}
					},
					{
						"key": {
							"type": "psa.refval-id",
							"value": {
								"label": "ARoT",
								"version": "0.1.4",
								"signer-id": "rLsRx+TaIXIFUjzkzhokWuGiOa48a/2eeHH35di66Gs="
							}
						},
						"value": {
							"digests": [
								"sha-256:o6XnFfDMV0pzw/m+u2vCTzL/1bZ7OHJEwskJ2neaFHg="
							]
						}
					}
				]
			}
		]
	}
}
`
	PSAKeysJSONTemplate = `{
	"lang": "en-GB",
	"tag-identity": {
		"id": "366D0A0A-5988-45ED-8488-2F2A544F6242",
		"version": 0
	},
	"entities": [
		{
			"name": "ACME Ltd.",
			"regid": "https://acme.example",
			"roles": [ "tagCreator", "creator", "maintainer" ]
		}
	],
	"triples": {
		"attester-verification-keys": [
			{
				"environment": {
					"class": {
						"id": {
							"type": "psa.impl-id",
							"value": "YWNtZS1pbXBsZW1lbnRhdGlvbi1pZC0wMDAwMDAwMDE="
						},
						"vendor": "ACME",
						"model": "RoadRunner"
					},
					"instance": {
						"type": "ueid",
						"value": "Ac7rrnuJJ6MiflMDz14PH3s0u1Qq1yUKwD+83jbsLxUI"
					}
				},
				"verification-keys": [
					{
						"key": "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg=="
					}
				]
			},
			{
				"environment": {
					"class": {
						"id": {
							"type": "psa.impl-id",
							"value": "YWNtZS1pbXBsZW1lbnRhdGlvbi1pZC0wMDAwMDAwMDE="
						},
						"vendor": "ACME",
						"model": "RoadRunner"
					},
					"instance": {
						"type": "ueid",
						"value": "AUyj5PUL8kjDl4cCDWj/0FyIdndRvyZFypI/V6mL7NKW"
					}
				},
				"verification-keys": [
					{
						"key": "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE6Vwqe7hy3O8Ypa+BUETLUjBNU3rEXVUyt9XHR7HJWLG7XTKQd9i1kVRXeBPDLFnfYru1/euxRnJM7H9UoFDLdA=="
					}
				]
			}
		]
	}
}
`
)
