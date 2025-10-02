// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

//nolint:lll
const (
	TestUIntInstance    = 45
	TestInvalidProdID   = -23
	TestUIntISVProdID   = 23
	TestInvalidInstance = -1
	TestPCEID           = "PCEID001"
	TestISVSVN          = uint(10)
	TestTCBEvalNum      = uint(11)
	TestTime            = "2025-01-29T00:00:00Z"
	TestSeamOID         = "2.16.840.1.113741.1.2.3.4.3"
	TestQEOID           = "2.16.840.1.113741.1.2.3.4.5"
	TestPCEOID          = "2.16.840.1.113741.1.2.3.4.4"
)

var (
	TestTeeMiscSelect = []byte{0x0B, 0x0C, 0x0D}
	TestTeeAttributes = []byte{0x01, 0x01}
)

//nolint:lll
var (
	TDXQECETemplate = `{
    "ev-triples": {
        "evidence-triples": [
            {
                "environment": {
                    "class": {
                        "id": {
                            "type": "oid",
                            "value": "2.16.840.1.113741.1.15.4.99.1"
                        },
                        "vendor": "Intel Corporation",
                        "model": "TDX QE TCB"
                    }
                },
                "measurements": [
                    {
                        "value": {
                            "miscselect": "oLDA0AAAAAA=",
                            "mrsigner": {
                                "type": "digest",
                                "value": [
                                    "sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU=",
                                    "sha-512;oxT8LcZjrnpra8Z4dZQFc5bms/VpzVD9XdtNG7r9K2qjFPwtxmOuemtrxnh1lAVzluaz9WnNUP1d200buv0rag=="
                                ]
                            },
                            "isvprodid": {
                                "type": "uint",
                                "value": 1
                            },
                            "tcbevalnum": {
                                "type": "uint",
                                "value": 11
                            },
                            "tcbstatus": {
                                "type": "string",
                                "value": [
                                    "UpToDate"
                                ]
                            },
                            "advisoryids": {
                                "type": "string",
                                "value": [
                                    "INTEL-SA-00078",
                                    "INTEL-SA-00079"
                                ]
                            }
                        },
                        "authorized-by": [
                            {
                                "type": "pkix-base64-key",
                                "value": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg==\n-----END PUBLIC KEY-----"
                            }
                        ]
                    }
                ]
            }
        ]
    }
}`
	TDXPCECETemplate = ` {
    "ev-triples": {
        "evidence-triples": [
            {
                "environment": {
                    "class": {
                        "id": {
                            "type": "oid",
                            "value": "2.16.840.1.113741.1.2.3.4.4"
                        },
                        "vendor": "Intel Corporation",
                        "model": "TDX PCE TCB"
                    }
                },
                "measurements": [
                    {
                        "value": {
                            "instanceid": {
                                "type": "uint",
                                "value": 45
                            },
                            "pceid": "PCEID001",
                            "tcbcompsvn": [
                                {
                                    "type": "uint",
                                    "value": 1
                                },
                                {
                                    "type": "uint",
                                    "value": 2
                                },
                                {
                                    "type": "uint",
                                    "value": 3
                                },
                                {
                                    "type": "uint",
                                    "value": 4
                                },
                                {
                                    "type": "uint",
                                    "value": 5
                                },
                                {
                                    "type": "uint",
                                    "value": 6
                                },
                                {
                                    "type": "uint",
                                    "value": 7
                                },
                                {
                                    "type": "uint",
                                    "value": 8
                                },
                                {
                                    "type": "uint",
                                    "value": 9
                                },
                                {
                                    "type": "uint",
                                    "value": 10
                                },
                                {
                                    "type": "uint",
                                    "value": 11
                                },
                                {
                                    "type": "uint",
                                    "value": 12
                                },
                                {
                                    "type": "uint",
                                    "value": 13
                                },
                                {
                                    "type": "uint",
                                    "value": 14
                                },
                                {
                                    "type": "uint",
                                    "value": 15
                                },
                                {
                                    "type": "uint",
                                    "value": 16
                                }
                            ]
                        }
                    }
                ]
            }
        ]
    }
}`
	TDXSeamCETemplate = `{
    "ev-triples": {
        "evidence-triples": [
            {
                "environment": {
                    "class": {
                        "id": {
                            "type": "oid",
                            "value": "2.16.840.1.113741.1.2.3.4.3"
                        },
                        "vendor": "Intel Corporation",
                        "model": "TDXSEAM"
                    }
                },
                "measurements": [
                    {
                        "value": {
                            "tcbdate": "2025-01-27T00:00:00Z",
                            "isvsvn": {
                                "type": "uint",
                                "value": 10
                            },
                            "attributes": "AQE=",
                            "mrtee": {
                                "type": "digest",
                                "value": [
                                    "sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
                                ]
                            },
                            "mrsigner": {
                                "type": "digest",
                                "value": [
                                    "sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU=",
                                    "sha-384;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXXkW3L1wMC1cttNjTq36X82"
                                ]
                            },
                            "isvprodid": {
                                "type": "bytes",
                                "value": "AQE="
                            },
                            "tcbevalnum": {
                                "type": "uint",
                                "value": 11
                            }
                        }
                    }
                ]
            }
        ]
    }
}
`
)
