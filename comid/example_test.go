// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/veraison/swid"
)

func Example_encode() {
	comid := NewComid().
		SetLanguage("en-GB").
		SetTagIdentity("my-ns:acme-roadrunner-supplement", 0).
		AddEntity("ACME Ltd.", &TestRegID, RoleCreator, RoleTagCreator).
		AddEntity("EMCA Ltd.", nil, RoleMaintainer).
		AddLinkedTag("my-ns:acme-roadrunner-base", RelSupplements).
		AddLinkedTag("my-ns:acme-roadrunner-old", RelReplaces).
		AddReferenceValue(
			&ValueTriple{
				Environment: Environment{
					Class: NewClassOID(TestOID).
						SetVendor("ACME Ltd.").
						SetModel("RoadRunner").
						SetLayer(0).
						SetIndex(1),
					Instance: MustNewUEIDInstance(TestUEID),
					Group:    MustNewUUIDGroup(TestUUID),
				},
				Measurements: *NewMeasurements().
					Add(
						MustNewUUIDMeasurement(TestUUID).
							SetRawValueBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{0xff, 0xff, 0xff, 0xff}).
							SetSVN(2).
							AddDigest(swid.Sha256_32, []byte{0xab, 0xcd, 0xef, 0x00}).
							AddDigest(swid.Sha256_32, []byte{0xff, 0xff, 0xff, 0xff}).
							SetFlagsTrue(FlagIsDebug).
							SetFlagsFalse(FlagIsSecure).
							SetSerialNumber("C02X70VHJHD5").
							SetUEID(TestUEID).
							SetUUID(TestUUID).
							SetMACaddr(MACaddr(TestMACaddr)).
							SetIPaddr(TestIPaddr),
					),
			},
		).
		AddEndorsedValue(
			&ValueTriple{
				Environment: Environment{
					Class: NewClassUUID(TestUUID).
						SetVendor("ACME Ltd.").
						SetModel("RoadRunner").
						SetLayer(0).
						SetIndex(1),
					Instance: MustNewUEIDInstance(TestUEID),
					Group:    MustNewUUIDGroup(TestUUID),
				},
				Measurements: *NewMeasurements().
					Add(
						MustNewUUIDMeasurement(TestUUID).
							SetRawValueBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{0xff, 0xff, 0xff, 0xff}).
							SetMinSVN(2).
							AddDigest(swid.Sha256_32, []byte{0xab, 0xcd, 0xef, 0x00}).
							AddDigest(swid.Sha256_32, []byte{0xff, 0xff, 0xff, 0xff}).
							SetFlagsTrue(FlagIsDebug).
							SetFlagsFalse(FlagIsSecure, FlagIsConfigured).
							SetSerialNumber("C02X70VHJHD5").
							SetUEID(TestUEID).
							SetUUID(TestUUID).
							SetMACaddr(MACaddr(TestMACaddr)).
							SetIPaddr(TestIPaddr),
					),
			},
		).
		AddAttestVerifKey(
			&KeyTriple{
				Environment: Environment{
					Instance: MustNewUUIDInstance(uuid.UUID(TestUUID)),
				},
				VerifKeys: *NewCryptoKeys().
					Add(
						MustNewPKIXBase64Key(TestECPubKey),
					),
			},
		).AddDevIdentityKey(
		&KeyTriple{
			Environment: Environment{
				Instance: MustNewUEIDInstance(TestUEID),
			},
			VerifKeys: *NewCryptoKeys().
				Add(
					MustNewPKIXBase64Key(TestECPubKey),
				),
		},
	).
		AddCondEndorseSeries(
			&CondEndorseSeriesTriple{
				Condition: ValueTriple{
					Environment: Environment{
						Class: NewClassOID(TestOID).
							SetVendor("ACME Ltd.").
							SetModel("RoadRunner").
							SetLayer(0).
							SetIndex(1),
						Instance: MustNewUEIDInstance(TestUEID),
						Group:    MustNewUUIDGroup(TestUUID),
					},
					Measurements: *NewMeasurements().
						Add(
							MustNewUUIDMeasurement(TestUUID).
								SetRawValueBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{0xff, 0xff, 0xff, 0xff}).
								SetSVN(2).
								AddDigest(swid.Sha256_32, []byte{0xab, 0xcd, 0xef, 0x00}).
								AddDigest(swid.Sha256_32, []byte{0xff, 0xff, 0xff, 0xff}).
								SetFlagsTrue(FlagIsDebug).
								SetFlagsFalse(FlagIsSecure).
								SetSerialNumber("C02X70VHJHD5").
								SetUEID(TestUEID).
								SetUUID(TestUUID).
								SetMACaddr(MACaddr(TestMACaddr)).
								SetIPaddr(TestIPaddr),
						),
				},
				Series: *NewCondEndorseSeriesRecords().
					Add(
						&CondEndorseSeriesRecord{
							Selection: *NewMeasurements().
								Add(
									MustNewUUIDMeasurement(TestUUID).
										SetRawValueBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{0xff, 0xff, 0xff, 0xff}).
										SetSVN(2).
										AddDigest(swid.Sha256_32, []byte{0xab, 0xcd, 0xef, 0x00}).
										AddDigest(swid.Sha256_32, []byte{0xff, 0xff, 0xff, 0xff}).
										SetFlagsTrue(FlagIsDebug).
										SetFlagsFalse(FlagIsSecure),
								),
							Addition: *NewMeasurements().
								Add(
									MustNewUUIDMeasurement(TestUUID).
										SetUEID(TestUEID).
										SetMACaddr(MACaddr(TestMACaddr)).
										SetIPaddr(TestIPaddr),
								),
						},
					),
			},
		)

	cbor, err := comid.ToCBOR()
	if err == nil {
		fmt.Printf("%x\n", cbor)
	}

	json, err := comid.ToJSON()
	if err == nil {
		fmt.Printf("%s\n", string(json))
	}

	// Output:
	// a50065656e2d474201a10078206d792d6e733a61636d652d726f616472756e6e65722d737570706c656d656e740282a3006941434d45204c74642e01d8207468747470733a2f2f61636d652e6578616d706c6502820100a20069454d4341204c74642e0281020382a200781a6d792d6e733a61636d652d726f616472756e6e65722d626173650100a20078196d792d6e733a61636d652d726f616472756e6e65722d6f6c64010104a5008182a300a500d86f445502c000016941434d45204c74642e026a526f616452756e6e65720300040101d902264702deadbeefdead02d8255031fb5abf023e4992aa4e95f9c1503bfa81a200d8255031fb5abf023e4992aa4e95f9c1503bfa01aa01d90228020282820644abcdef00820644ffffffff03a201f403f504d9023044010203040544ffffffff064802005e1000000001075020010db8000000000000000000000068086c43303258373056484a484435094702deadbeefdead0a5031fb5abf023e4992aa4e95f9c1503bfa018182a300a500d8255031fb5abf023e4992aa4e95f9c1503bfa016941434d45204c74642e026a526f616452756e6e65720300040101d902264702deadbeefdead02d8255031fb5abf023e4992aa4e95f9c1503bfa81a200d8255031fb5abf023e4992aa4e95f9c1503bfa01aa01d90229020282820644abcdef00820644ffffffff03a300f401f403f504d9023044010203040544ffffffff064802005e1000000001075020010db8000000000000000000000068086c43303258373056484a484435094702deadbeefdead0a5031fb5abf023e4992aa4e95f9c1503bfa028182a101d902264702deadbeefdead81d9022a78b12d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0a4d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a304441516344516741455731427671462b2f727938425761375a454d553178595948455138420a6c4c54344d46484f614f2b4943547449767245654570722f7366544150363648326843486462354845584b74524b6f6436514c634f4c504131513d3d0a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d038182a101d8255031fb5abf023e4992aa4e95f9c1503bfa81d9022a78b12d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0a4d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a304441516344516741455731427671462b2f727938425761375a454d553178595948455138420a6c4c54344d46484f614f2b4943547449767245654570722f7366544150363648326843486462354845584b74524b6f6436514c634f4c504131513d3d0a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d08818282a300a500d86f445502c000016941434d45204c74642e026a526f616452756e6e65720300040101d902264702deadbeefdead02d8255031fb5abf023e4992aa4e95f9c1503bfa81a200d8255031fb5abf023e4992aa4e95f9c1503bfa01aa01d90228020282820644abcdef00820644ffffffff03a201f403f504d9023044010203040544ffffffff064802005e1000000001075020010db8000000000000000000000068086c43303258373056484a484435094702deadbeefdead0a5031fb5abf023e4992aa4e95f9c1503bfa818281a200d8255031fb5abf023e4992aa4e95f9c1503bfa01a501d90228020282820644abcdef00820644ffffffff03a201f403f504d9023044010203040544ffffffff81a200d8255031fb5abf023e4992aa4e95f9c1503bfa01a3064802005e1000000001075020010db8000000000000000000000068094702deadbeefdead
	// {"lang":"en-GB","tag-identity":{"id":"my-ns:acme-roadrunner-supplement"},"entities":[{"name":"ACME Ltd.","regid":"https://acme.example","roles":["creator","tagCreator"]},{"name":"EMCA Ltd.","roles":["maintainer"]}],"linked-tags":[{"target":"my-ns:acme-roadrunner-base","rel":"supplements"},{"target":"my-ns:acme-roadrunner-old","rel":"replaces"}],"triples":{"reference-values":[{"environment":{"class":{"id":{"type":"oid","value":"2.5.2.8192"},"vendor":"ACME Ltd.","model":"RoadRunner","layer":0,"index":1},"instance":{"type":"ueid","value":"At6tvu/erQ=="},"group":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}},"measurements":[{"key":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"},"value":{"svn":{"type":"exact-value","value":2},"digests":["sha-256-32;q83vAA==","sha-256-32;/////w=="],"flags":{"is-secure":false,"is-debug":true},"raw-value":{"type":"bytes","value":"AQIDBA=="},"raw-value-mask":"/////w==","mac-addr":"02:00:5e:10:00:00:00:01","ip-addr":"2001:db8::68","serial-number":"C02X70VHJHD5","ueid":"At6tvu/erQ==","uuid":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}}]}],"endorsed-values":[{"environment":{"class":{"id":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"},"vendor":"ACME Ltd.","model":"RoadRunner","layer":0,"index":1},"instance":{"type":"ueid","value":"At6tvu/erQ=="},"group":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}},"measurements":[{"key":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"},"value":{"svn":{"type":"min-value","value":2},"digests":["sha-256-32;q83vAA==","sha-256-32;/////w=="],"flags":{"is-configured":false,"is-secure":false,"is-debug":true},"raw-value":{"type":"bytes","value":"AQIDBA=="},"raw-value-mask":"/////w==","mac-addr":"02:00:5e:10:00:00:00:01","ip-addr":"2001:db8::68","serial-number":"C02X70VHJHD5","ueid":"At6tvu/erQ==","uuid":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}}]}],"dev-identity-keys":[{"environment":{"instance":{"type":"ueid","value":"At6tvu/erQ=="}},"verification-keys":[{"type":"pkix-base64-key","value":"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"}]}],"attester-verification-keys":[{"environment":{"instance":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}},"verification-keys":[{"type":"pkix-base64-key","value":"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"}]}],"conditional-endorsement-series":[{"statefulenv":{"environment":{"class":{"id":{"type":"oid","value":"2.5.2.8192"},"vendor":"ACME Ltd.","model":"RoadRunner","layer":0,"index":1},"instance":{"type":"ueid","value":"At6tvu/erQ=="},"group":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}},"measurements":[{"key":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"},"value":{"svn":{"type":"exact-value","value":2},"digests":["sha-256-32;q83vAA==","sha-256-32;/////w=="],"flags":{"is-secure":false,"is-debug":true},"raw-value":{"type":"bytes","value":"AQIDBA=="},"raw-value-mask":"/////w==","mac-addr":"02:00:5e:10:00:00:00:01","ip-addr":"2001:db8::68","serial-number":"C02X70VHJHD5","ueid":"At6tvu/erQ==","uuid":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"}}]},"series":[{"selection":[{"key":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"},"value":{"svn":{"type":"exact-value","value":2},"digests":["sha-256-32;q83vAA==","sha-256-32;/////w=="],"flags":{"is-secure":false,"is-debug":true},"raw-value":{"type":"bytes","value":"AQIDBA=="},"raw-value-mask":"/////w=="}}],"addition":[{"key":{"type":"uuid","value":"31fb5abf-023e-4992-aa4e-95f9c1503bfa"},"value":{"mac-addr":"02:00:5e:10:00:00:00:01","ip-addr":"2001:db8::68","ueid":"At6tvu/erQ=="}}]}]}]}}
}

func Example_encode_PSA_attestation_verification() {
	comid := NewComid().
		SetTagIdentity("my-ns:acme-roadrunner-supplement", 0).
		AddEntity("ACME Ltd.", &TestRegID, RoleCreator, RoleTagCreator, RoleMaintainer).
		AddAttestVerifKey(
			&KeyTriple{
				Environment: Environment{
					Instance: MustNewUEIDInstance(TestUEID),
				},
				VerifKeys: *NewCryptoKeys().
					Add(
						MustNewPKIXBase64Key(TestECPubKey),
					),
			},
		)

	cbor, err := comid.ToCBOR()
	if err == nil {
		fmt.Printf("%x\n", cbor)
	}

	json, err := comid.ToJSON()
	if err == nil {
		fmt.Printf("%s", string(json))
	}

	// Output:
	// a301a10078206d792d6e733a61636d652d726f616472756e6e65722d737570706c656d656e740281a3006941434d45204c74642e01d8207468747470733a2f2f61636d652e6578616d706c65028301000204a1038182a101d902264702deadbeefdead81d9022a78b12d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0a4d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a304441516344516741455731427671462b2f727938425761375a454d553178595948455138420a6c4c54344d46484f614f2b4943547449767245654570722f7366544150363648326843486462354845584b74524b6f6436514c634f4c504131513d3d0a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d
	// {"tag-identity":{"id":"my-ns:acme-roadrunner-supplement"},"entities":[{"name":"ACME Ltd.","regid":"https://acme.example","roles":["creator","tagCreator","maintainer"]}],"triples":{"attester-verification-keys":[{"environment":{"instance":{"type":"ueid","value":"At6tvu/erQ=="}},"verification-keys":[{"type":"pkix-base64-key","value":"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"}]}]}}
}

func Example_decode_JSON() {
	j := `
{
    "lang": "en-GB",
    "tag-identity": {
        "id": "43BBE37F-2E61-4B33-AED3-53CFF1428B16",
        "version": 1
    },
    "entities": [
        {
            "name": "ACME Ltd.",
            "regid": "https://acme.example",
            "roles": [
                "tagCreator"
            ]
        },
        {
            "name": "EMCA Ltd.",
            "regid": "https://emca.example",
            "roles": [
                "maintainer",
                "creator"
            ]
        }
    ],
    "linked-tags": [
        {
            "target": "6F7D8D2F-EAEC-4A15-BB46-1E4DCB85DDFF",
            "rel": "replaces"
        }
    ],
    "triples": {
        "reference-values": [
            {
                "environment": {
                    "class": {
                        "id": {
                            "type": "uuid",
                            "value": "83294297-97EB-42EF-8A72-AE9FEA002750"
                        },
                        "vendor": "ACME",
                        "model": "RoadRunner Boot ROM",
                        "layer": 0,
                        "index": 0
                    },
                    "instance": {
                        "type": "ueid",
                        "value": "Ad6tvu/erb7v3q2+796tvu8="
                    }
                },
                "measurements": [
                    {
                        "value": {
                            "digests": [
                                "sha-256:5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
                            ]
                        }
                    }
                ]
            },
            {
                "environment": {
                    "class": {
                        "id": {
                            "type": "bytes",
                            "value": "YWNtZS1pbXBsZW1lbnRhdGlvbi1pZC0wMDAwMDAwMDE="
                        },
                        "vendor": "Generic-X",
                        "model": "Turbo PRoT"
                    }
                },
                "measurements": [
                    {
                        "key": {
                            "type": "uuid",
                            "value": "31fb5abf-023e-4992-aa4e-95f9c1503bfa"
                        },
                        "value": {
                            "digests": [
                                "sha-256:5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
                            ],
                            "svn": {
                                "type": "exact-value",
                                "value": 1
                            },
                            "mac-addr": "00:00:5e:00:53:01"
                        }
                    }
                ]
            }
        ],
        "endorsed-values": [
            {
                "environment": {
                    "class": {
                        "id": {
                            "type": "oid",
                            "value": "2.16.840.1.101.3.4.2.1"
                        }
                    },
                    "instance": {
                        "type": "uuid",
                        "value": "9090B8D3-3B17-474C-A0B9-6F54731CAB72"
                    }
                },
                "measurements": [
                    {
                        "value": {
                            "mac-addr": "00:00:5e:00:53:01",
                            "ip-addr": "2001:4860:0:2001::68",
                            "serial-number": "C02X70VHJHD5",
                            "ueid": "Ad6tvu/erb7v3q2+796tvu8=",
                            "uuid": "9090B8D3-3B17-474C-A0B9-6F54731CAB72",
                            "raw-value": {
                                "type": "bytes",
                                "value": "cmF3dmFsdWUKcmF3dmFsdWUK"
                            },
                            "raw-value-mask": "qg==",
                            "op-flags": [
                                "notSecure"
                            ],
                            "digests": [
                                "sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU=",
                                "sha-384;S1bPoH+usqtX3pIeSpfWVRRLVGRw66qrb3HA21GN31tKX7KPsq0bSTQmRCTrHlqG"
                            ],
                            "version": {
                                "scheme": "semaver",
                                "value": "1.2.3beta4"
                            },
                            "svn": {
                                "type": "min-value",
                                "value": 10
                            }
                        }
                    }
                ]
            }
        ],
        "attester-verification-keys": [
            {
                "environment": {
                    "group": {
                        "type": "uuid",
                        "value": "83294297-97EB-42EF-8A72-AE9FEA002750"
                    }
                },
                "verification-keys": [
                    {
                        "type": "pkix-base64-key",
                        "value": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"
                    }
                ]
            }
        ],
        "dev-identity-keys": [
            {
                "environment": {
                    "instance": {
                        "type": "uuid",
                        "value": "4ECCE47C-85F2-4FD9-9EC6-00DEB72DA707"
                    }
                },
                "verification-keys": [
                    {
                        "type": "pkix-base64-key",
                        "value": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"
                    },
                    {
                        "type": "pkix-base64-key",
                        "value": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"
                    }
                ]
            }
        ],
        "conditional-endorsement-series": [
            {
                "statefulenv": {
                    "environment": {
                        "class": {
                            "id": {
                                "type": "oid",
                                "value": "2.5.2.8192"
                            },
                            "vendor": "ACME Ltd.",
                            "model": "RoadRunner",
                            "layer": 0,
                            "index": 1
                        },
                        "instance": {
                            "type": "ueid",
                            "value": "At6tvu/erQ=="
                        },
                        "group": {
                            "type": "uuid",
                            "value": "31fb5abf-023e-4992-aa4e-95f9c1503bfa"
                        }
                    },
                    "measurements": [
                        {
                            "key": {
                                "type": "uuid",
                                "value": "31fb5abf-023e-4992-aa4e-95f9c1503bfa"
                            },
                            "value": {
                                "svn": {
                                    "type": "exact-value",
                                    "value": 1
                                },
                                "digests": [
                                    "sha-256-32;q83vAA==",
                                    "sha-256-32;/////w=="
                                ],
                                "flags": {
                                    "is-secure": false,
                                    "is-debug": true
                                },
                                "raw-value": {
                                    "type": "bytes",
                                    "value": "AQIDBA=="
                                },
                                "raw-value-mask": "/////w==",
                                "serial-number": "C02X70VHJHD5",
                                "uuid": "31fb5abf-023e-4992-aa4e-95f9c1503bfa"
                            }
                        }
                    ]
                },
                "series": [
                    {
                        "selection": [
                            {
                                "key": {
                                    "type": "uuid",
                                    "value": "31fb5abf-023e-4992-aa4e-95f9c1503bfa"
                                },
                                "value": {
                                    "svn": {
                                        "type": "exact-value",
                                        "value": 2
                                    },
                                    "version": {
                                        "value": "2.0.0",
                                        "scheme": "semver"
                                    }
                                }
                            }
                        ],
                        "addition": [
                            {
                                "key": {
                                    "type": "uuid",
                                    "value": "31fb5abf-023e-4992-aa4e-95f9c1503bfa"
                                },
                                "value": {
                                    "mac-addr": "02:00:5e:10:00:00:00:01",
                                    "ip-addr": "2001:db8::68",
                                    "ueid": "At6tvu/erQ=="
                                }
                            }
                        ]
                    },
                    {
                        "selection": [
                            {
                                "key": {
                                    "type": "uuid",
                                    "value": "31fb5abf-023e-4992-aa4e-95f9c1503bfa"
                                },
                                "value": {
                                    "svn": {
                                        "type": "exact-value",
                                        "value": 3
                                    },
                                    "version": {
                                        "value": "3.0.1",
                                        "scheme": "semver"
                                    }
                                }
                            }
                        ],
                        "addition": [
                            {
                                "key": {
                                    "type": "uuid",
                                    "value": "31fb5abf-023e-4992-aa4e-95f9c1503bfa"
                                },
                                "value": {
                                    "mac-addr": "02:00:5e:10:00:00:00:02",
                                    "ip-addr": "2001:db8::69",
                                    "ueid": "At6tvu/erQ=="
                                }
                            }
                        ]
                    }
                ]
            }
        ]
    }
}
`
	comid := Comid{}
	err := comid.FromJSON([]byte(j))

	if err != nil {
		fmt.Printf("FAIL: %v", err)
	} else {
		fmt.Println("OK")
	}

	// Output: OK
}

var (
	// test cases are based on diag files here:
	// https://github.com/ietf-rats-wg/draft-ietf-rats-corim/tree/main/cddl/examples

	//go:embed testcases/comid-1.cbor
	testComid1 []byte

	//go:embed testcases/comid-2.cbor
	testComid2 []byte

	//go:embed testcases/comid-design-cd.cbor
	testComidDesignCD []byte

	//go:embed testcases/comid-firmware-cd.cbor
	testComidFirmwareCD []byte

	//go:embed testcases/comid-3.cbor
	testComid3 []byte

	//go:embed testcases/comid-4.cbor
	testComid4 []byte

	//go:embed testcases/comid-5.cbor
	testComid5 []byte

	//go:embed testcases/comid-cond-endorse-series.cbor
	testComidCondEndorseSeries []byte
)

func TestExample_decode_CBOR(_ *testing.T) {
	tvs := []struct {
		descr string
		inp   []byte
	}{
		{
			descr: "Test with CoMID-1 Diag",
			inp:   testComid1,
		},
		{
			descr: "Test with CoMID-2 Diag",
			inp:   testComid2,
		},
		{
			descr: "Test with CoMID-Design-CD Diag",
			inp:   testComidDesignCD,
		},
		{
			descr: "Test with Firmware-CD Diag",
			inp:   testComidFirmwareCD,
		},
		{
			descr: "Test with CoMID-3 Diag",
			inp:   testComid3,
		},
		{
			descr: "Test with CoMID-4 Diag",
			inp:   testComid4,
		},
		{
			descr: "Test with CoMID-5 Diag",
			inp:   testComid5,
		},
		{
			descr: "Test with CoMID-cond-endorse-series Diag",
			inp:   testComidCondEndorseSeries,
		},
	}
	for _, tv := range tvs {
		comid := Comid{}
		err := comid.FromCBOR(tv.inp)
		if err != nil {
			fmt.Printf("FAIL: %v", err)
		} else {
			fmt.Println("OK")
		}
		// Output: OK
	}
}
