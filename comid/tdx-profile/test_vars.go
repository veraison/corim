// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

//nolint:lll
var (
	TestRegID              = "https://intel.com"
	TestOID                = "2.16.840.1.113741.1.2.3.4.5"
	TestUIntInstance       = 45
	TestByteInstance       = []byte{0x45, 0x46, 0x47}
	TestInvalidProdID      = -23
	TestUIntISVProdID      = 23
	TestBytesISVProdID     = []byte{0x01, 0x02, 0x03}
	TestInvalidInstance    = -1
	TestTeeAttributes      = []byte{0x01, 0x01}
	TestTeeMiscSelect      = []byte{0x0B, 0x0C, 0x0D}
	TestPCEID              = "PCEID001"
	TestCompSVN            = []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	TestTCBStatus          = []string{"OutOfDate", "ConfigurationNeeded", "UpToDate"}
	TestInvalidTCBStatus   = []int{1, 2, 3}
	TestAdvisoryIDs        = []string{"SA-00078", "SA-00077", "SA-00079"}
	TestInvalidAdvisoryIDs = []float64{1.234, 2.567}
	TestISVSVN             = 10
	TestTCBEvalNum         = 11
	TestTime               = "2025-01-29T00:00:00Z"
	TDXPCERefValTemplate   = `{
  "lang": "en-GB",
  "tag-identity": {
    "id": "43BBE37F-2E61-4B33-AED3-53CFF1428B17",
    "version": 0
  },
  "entities": [
    {
      "name": "INTEL",
      "regid": "https://intel.com",
      "roles": [
        "tagCreator",
        "creator",
        "maintainer"
      ]
    }
  ],
  "triples": {
    "reference-values": [
      {
        "environment": {
          "class": {
            "id": {
              "type": "oid",
              "value": "2.16.840.1.113741.1.2.3.4.6"
            },
            "vendor": "Intel Corporation",
            "model": "0123456789ABCDEF"
          }
        },
        "measurements": [
          {
            "value": {
              "instanceid": {
               "type": "uint",
               "value": 11
              },
              "tcbcompsvn": [10, 10, 2, 2, 2, 1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0],
              "pceid": "0000"
            },
            "authorized-by": {
              "type": "pkix-base64-key",
              "value": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg==\n-----END PUBLIC KEY-----"
            }
          }
        ]
      }
    ]
  }
}
`
	TDXQERefValTemplate = `{
  "lang": "en-GB",
  "tag-identity": {
    "id": "43BBE37F-2E61-4B33-AED3-53CFF1428B16",
    "version": 0
  },
  "entities": [
    {
      "name": "INTEL",
      "regid": "https://intel.com",
      "roles": [
        "tagCreator",
        "creator",
        "maintainer"
      ]
    }
  ],
  "triples": {
    "reference-values": [
      {
        "environment": {
          "class": {
            "id": {
              "type": "oid",
              "value": "2.16.840.1.113741.1.2.3.4.1"
            },
            "vendor": "Intel Corporation",
            "model": "TDX QE TCB"
          }
        },
        "measurements": [
          {
            "value": {
              "miscselect": "wAAAAPv/AAA=",
              "tcbevalnum": 11, 
              "mrsigner": [
                "sha-256:h0KPxSKAPTEGXnvOPPA/5HUJZjHl4Hu9eg/eYMTPJcc=",
                "sha-512:oxT8LcZjrnpra8Z4dZQFc5bms/VpzVD9XdtNG7r9K2qjFPwtxmOuemtrxnh1lAVzluaz9WnNUP1d200buv0rag=="
              ],
            "isvprodid": {
              "type": "bytes",
              "value": "AwM="
              }
            },
            "authorized-by": {
              "type": "pkix-base64-key",
              "value": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg==\n-----END PUBLIC KEY-----"
            }
          }
        ]
      }
    ]
  }
}	
`
	TDXSeamRefValJSONTemplate = `{
  "lang": "en-GB",
  "tag-identity": {
    "id": "43BBE37F-2E61-4B33-AED3-53CFF1428B20",
    "version": 0
  },
  "entities": [
    {
      "name": "INTEL",
      "regid": "https://intel.com",
      "roles": [
        "tagCreator",
        "creator",
        "maintainer"
      ]
    }
  ],
  "triples": {
    "reference-values": [
      {
        "environment": {
          "class": {
            "id": {
              "type": "oid",
              "value": "2.16.840.1.113741.1.2.3.4.5"
            },
            "vendor": "Intel Corporation",
            "model": "TDX SEAM"
          }
        },
        "measurements": [
          {
            "value": {
              "isvprodid": {
              "type": "bytes",
              "value": "AwM="
              },
              "isvsvn": 10,
              "attributes": "8AoL",
              "tcbevalnum": 11,
              "mrtee" : ["sha-256:h0KPxSKAPTEGXnvOPPA/5HUJZjHl4Hu9eg/eYMTPJcc="],
              "mrsigner": [
                "sha-256:h0KPxSKAPTEGXnvOPPA/5HUJZjHl4Hu9eg/eYMTPJcc=",
                "sha-512:oxT8LcZjrnpra8Z4dZQFc5bms/VpzVD9XdtNG7r9K2qjFPwtxmOuemtrxnh1lAVzluaz9WnNUP1d200buv0rag=="
              ]
            },
            "authorized-by": {
              "type": "pkix-base64-key",
              "value": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg==\n-----END PUBLIC KEY-----"
            }
          }
        ]
      }
    ]
  }
}
`
)
