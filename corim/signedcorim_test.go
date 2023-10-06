// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testES256Key = []byte(`{
		"kty": "EC",
		"crv": "P-256",
		"x": "MKBCTNIcKUSDii11ySs3526iDZ8AiTo7Tu6KPAqv7D4",
		"y": "4Etl6SRW2YiLUrN5vfvVHuhp7x8PxltmWWlbbM4IFyM",
		"d": "870MB6gfuTJ4HtUnUvYMyJpr5eUZNP4Bk43bVdj3eAE",
		"use": "enc",
		"kid": "1"
	  }`)

	testES384Key = []byte(`{
		"kty": "EC",
		"d": "XiZ_ZEDMw3Hr9BjNc_4qbNxMG6VpkFHTN3KcdT1UlOc51pFwS1t6Yg_aFYJTGMBf",
		"use": "sig",
		"crv": "P-384",
		"x": "Ay-c_vlONI_FNQn4PNHXwEswuoxOTqOEHNIQbSKv5OnC_KBLwAbg5uBQRHCRmFnu",
		"y": "mJpRrG-ex0R08heh1qm-osCH7SSTKC1Bjx1SrFpUQZCiYQXdPLIokC0DGRAMYq41",
		"alg": "ES384"
	  }`)

	testES512Key = []byte(`{
		"use": "sig",
		"kty": "EC",
		"kid": "Xt7n2MSHsgErmf1Uq-UZV451DhzlSPVuH75Rj9adAZ0",
		"crv": "P-521",
		"alg": "ES512",
		"x": "AVBBp8Mckn-HYsdx5bMSkFRxGhKH2M7ked49PqK2PzG2A5QEBPc813AHUO3MHoe-_JQjEm-r-E52sNln-zn6OFJL",
		"y": "AcsVxiDaIJpr3MToPmDqSjWnCkg765Rum3DWuFNaTmvietwrY6OYdoW995m4SkWv4GYI0mdchuXoThvPn0CXcDb9",
		"d": "AeSLG30MsuX6wzm-AYpBbTooVPt3GvU_Fl9LesAFZrtJ4HJhPL3QhMLmiDbB3Am0j_IpIR3P9dTJTNpt6B_YSVda"
	  }`)

	testEdDSAKey = []byte(`{
		"use": "sig",
		"kty": "OKP",
		"kid": "RBx2781Ag7Sd1vmuVbxpe0LzWT94pmB3GPtNx6m_gsQ",
		"crv": "Ed25519",
		"alg": "EdDSA",
		"x": "JL3cmVCzN3m3afnctG2agbjb6nrZWFl48A8Feknkpx0",
		"d": "m8LDAfKvGWAZTXWC21tzHeSYLqVSP4YpzI-Z7fL3NEY"
	  }`)

	testPS256Key = []byte(`{
		"use": "sig",
		"kty": "RSA",
		"kid": "jYGw-iPMi7AxzdJPMHYh_gb9YI-BQGAVAvf6hgZndzw",
		"alg": "PS256",
		"n": "r8tDvmXtJjtGOgX34bxDGT3-v2AtfVkP4vhdOl5Wau-XFyaPNpob5u3DtNsYUnHREQFnrPbIp02IeassUqSi6FlT9SZsYX8M5xkfpCuLb6FD5Loqz4ZMhzqtMNoKjUt2_9tdyW-iMfMm-EWLfVRfiXnfXq__o122LZ93-zmR4kEusCp7rUa12-E48pv4Wu5CwKntz08DjP-WB-yR8ZT1_F4IacqK6Uhhdh56TLONoUytyQkJTYi0lvohzVVtuRp7jXDpG9TBMBsyAJ0yj6FvpA-Bs1mkMNUlUr-p6xbSIAsOrv4FBtLXDKtApurRQmnNAtm4LTE4RsuQxI2FSKlOnQ",
		"e": "AQAB",
		"d": "bx6bObUQDISXRYIUSDpKZ6BKcQoIdx1e72dy9rw-_-VmqhmTmT4cuQI-HQoI-8Q6FPfAYxKzjx1xUQckQzESULB2Y5XgGFjI_SNiXtGvl-ZmFiSffwIzSZ-Lbj_FP78d_2jYhcXszooWbgT3wUceBLZmvWGew9MunvQYUVL4pfzktRn7zX0u9ks8GYxNfnwbeB8e8x7ZGrGpPSy5MNJkTHuPpu6XGXR8fJFEEFZZdsyJYd-Ii5Nma2uXyVZfBeRYmqlRIvok5jcNGmFm9wM291v7fieuJycSV71iFQnZfoF48uiNt3mGsdzPNfulCSKjMdR03jk-v1YyyQP34wu8uQ",
		"p": "ww7iRrQKk37YmQP_4xVtdAtOj5-bBWkM6wid2VNDss3u3GbivCqchqY2fQFgw9wKVYN0T7hS8wgErKPgE7ALTImwrK66TdTLZ_ljLScoYcrHRdAnTiqSbK3iyUnCs15ptSzOSXJHUXeVbynK7K9wo6TALz8c7-y05Gc_XpvM4I8",
		"q": "5reZTbuRXJ3O2sIG4vBvmn0UujZ9WbnvzQ36c92vxIqsWZ1MJzzc-9FKv9iG3zHS8tLLLYT4V6InIovJ6ZNgit0HieFyWfGNfc-3rt13OZwcFLhAu5nizZzkh24Mx0lquXoRxQwgc43Fg0Lk64C-xhgWAhW6OeNIxwp3zxpLHBM",
		"dp": "MfOK0M5kcvcl4rGagv3GxNPsb21RFqabT0kqmy_ug0inZbvXTpae9QB1rbd_n0inQNTkIVIzs9cW01s4E_KeQiB0pRQt06at3FeKJVMEzV5Pf7pZhnPygXBaRm_kM2j3KxVpUne4ec1k8E3EkK4w60dSjAbekzaL8H3cRY8ifVM",
		"dq": "oe6PDP3vIqAoRWYVS0cSLc5ItAH2rPlSFAwRky0vZrUmDqfWgVu4ho349vnUf-cKdh_5NvOzEl7fNOIET4p_IjfMSLwRdIuTkZAvDe6m9apaEzjXRlTV2RabV2qoUV94JsJEopbGWBRTYrOa1KhCPes91yzEzkh2Fi2Etblwqj0",
		"qi": "tWSsmBzj2NZ71Phe6wOUSKvdvLFyASRgnF-YZN1VKsyXER-SPujBiSjHC5f0TCcNKg4Y1GZxRjfVin2_ZO_qmlOij3xc6GsxVdXXvyzKg1wamf-Fut0I2K6dGLalJuCY0rp_Q1BGLv-LET2wArEkMgYWtsQLRuZXCktMi9XPsiU"
	  }`)

	testPS384Key = []byte(`{
		"use": "sig",
		"kty": "RSA",
		"kid": "-Q8rUDoeh60B1OorMT6v80KxBTy0SzqeAab2Isi09Ko",
		"alg": "PS384",
		"n": "6aeDsQF3WjbJEBq00Oxne7m9CNXGUn4ANI4DsjYDDfZnewNiZDhnO14g-VvrnF5MoP725Ho9b2VSQZ5Ke0Bp8mtvIXohsCtzodNe7_dcU9ycD-PdLBRZalvloLKt_o8-rVdWEpiJtg8CO4VZHu2ReOX-mzecr4W8RxBiQGRaP6t1XAUwj3bQnkHYcQ9VtD3-RxRS9vr1xDUv-v8HrmlUtrlBTAABIzDUS8YXP5h2hjGfd8qaY1c-4TvKAdNuasTCWUjJjVhj_-6pA0bmxM0F2jEdZo-1yA0vZE0V6hBeayMCTOw08wPdOVC07dupZMFYBiRdH1jL_eGwBaCTp2QIEw",
		"e": "AQAB",
		"d": "6FTnP8Rjd1LujpLfpLbNF1vTOcvHjhM4BQoJZtUKKIIQ12LAHUNwcrngM9NQ7oVd0OB1gy6BlBi9t_27td6Q-roVIMaeZNxv-EODLT1bkw_UJoC_VatOVdHW_PluxaaN_jLPpWID3QIDiEfKHFTBx-N6TcD4jhd-5XLHH5wpmQ_wE6G55PJmq0jONCntK5c4D21CAePzT2FL3oJ4OmpsZqChWAmYm1Lv2OFzSQpaI0HjE6Cchlkx_8ENKhEJrfxkLCNA1yPUfkXGOkpF6GFVpjTSVYjnwyYJKeC5D5KDb6Ln_n56OFejksuzCcKJCkDEKsS7UZmxw02DbYgjjCcsgQ",
		"p": "9ygG1oK4M9_ZJyMdAoDQFIkc4ISekZAm1_y8PW-NR0_l-bZkP2SUtRXBiT0qEtg1ltkVSCLRRtd940HLuOu7ZKOoF0DRvzFlkyrKyqwCwjGSBhqiSdC9PW9YGjYabIRzJaXBvp_F4v6u3KV1aIJDx9hrlAEJcUtzW1EHYB_GWuE",
		"q": "8gPO5_MpV2J6WVXVmFfYJ6TJvPLkM0oHfticCH3wkj3E6NgCygOkWe8SuFEWb2FJJ8AwBXCLp-j7gN2oU-TylqFd_Of2oQb4YjTJAwvub-kJuN8Y434C1DAEdSEqNGyvMAxMKE-q89DtzDmDKd5v5gI9Ci4DCQdl27q0ddks1XM",
		"dp": "17zyqxATpgRBUu5Nhj_WYfaFZF2e5ETGA0azMZVL5vGRNvXEb6lmPOMuupLPRP_BV1lKQFtT_dhgJJzsLRBn1KMeOJ31-EQv-9Qgi-S1y7jlU7qv6mrwpM2qQ8byLcM3l6cmhTSF0WyqSiOLZpw-ehUpYlm9Wk2X9h-2pmtWA0E",
		"dq": "GsirHGZ_28jtS3fBZNPL-080eHHVKYv22mX0lsgBWN33LeHCJUNT7BQWWUm4FumIZBrT9bYn7pRNSUy-tVIwOtVvBm9Rjy6rTIsU9_5ZDA-ZYNln8r1eaMdLpv7doeGpXcLupsNyYvtrZd-zkW2pqqXyxW6kLVqhPjkigaxgVts",
		"qi": "jH5q55Ez6f3yQE3-PEf3a6i7C1PgHfAy6E3sXONoyztbCMm6hSj_agBUj2bt_Q9XeddQf66OmBa2zPGGJqp0lclSDwuJB82iZspHjsAQuUQaBWE2eOfJlJkU4L34ibahQT2DgRxuosvv9gCSIXMOh7fZ2hp71XExrQTuA2BIC_s"
	  }`)

	testPS512Key = []byte(`{
		"use": "sig",
		"kty": "RSA",
		"kid": "umrdA9anNl7Qewgx-QWa7iBeVBHJ_i50NXdAiLLoU80",
		"alg": "PS512",
		"n": "yAdpoRbB0CPVjBFoVr4k4bEvrtpxLYDiqtGBwWxSEBlXxSLfABIgQXpR2Dkhda9Wsk-kBOEXS-8MUs8Slq0f6g0oaG1X_3IbZOcI_QzKDcIiQnoJH5herrY8S0eLkcUWvmAHuQoJ7G6JIiMePuM6eLNagwJzSIr6qq4y7-3ua1ZEieJXUkhsHEL3_ZsIWEIZ1tkI6nrfZnqRc1OMXqxRgeZTI5nqF99UptxUwGoTKGQB-yiH9kVGyYqjgvOQFCVLOyUMbdN76BAKGduQ56vjfNqqRNtZm_E1S7dhyu9F86C5p97BswRcFjkCb499bDEvYofJHY7npBnmJMRBhF1krw",
		"e": "AQAB",
		"d": "Qd5XBUnmJrE2J_qvfij9Iijjx9N9A3v2qEN3VAdkepKt2WfjQTW203kBLI-bmhJUHUGmhEjPEB021KoFuAJoiP0uOj0PhjnAFZkS16l3e9Jaz8M57-KQAz5VWoDD0AuzspsSz_cjT20S0V_5HMJcxdRh0NRkvBWv97aHZYTXRxa6SAZnyElAE2IqkLkKFRyK63V8TXJNFUc9DRzDtKgDszpZ7QflMqqb0nh02VImrcZdyRQtHEyCyCZvUURzL44qMayYWVmlDNwjXsS37iZLP0kValfe6sjNrS8JCq52ccyD_PBGTDv6aDDatrUdg64aEcNFvKzbRdQKxlMY9U8G2Q",
		"p": "_ohOSzzBHPycHTtZUOXR4b-DVfDtdLLxtiPBQUdcWvu_TL9HFkjRdx3Y-YaW6WRFAirF6dDz1VixhQFgnwlbRkYW2U8VB1wlYmVPOnJvDWO7aWVfJ8_QZjfhdfnQya_ATYu-X8o71yo0eP6GLRGNs6a_1L13o_2P8Dtg0FOKVws",
		"q": "yS6ooWwNOuEKLZUmLLc-euxaVxEOabKwRWoxTzbSaIeTR0MfhTDWWBVjWqKp6wDtN-3gFTnC56oUwddsktTxKvpQHj49CHCGPPUpOFROJoqRrjnxi4cd35oPkbbcoujMacSCLygKpNRZb7eqnVe8pW5aFIXUYqN5hjA1MxlEH20",
		"dp": "5xAKM1bV4GCZwBeufzgCjjLzIUNz7Oq9bqGKwJ3tg1LiWOOTvvEf5kicPfkman1yAAOgYyAjGlxH2vxjIDy4NVVPTLrz1hiaf3aEtARKOBd_fLBf755CC2lTLWw5U75Ojpb7na3TIQLZW7WDTMqQnrQTlSbiw2ZeErF0s-oCvf0",
		"dq": "NJQlLkr3CjRWXKNmXrllcurikW67vZQdzYZ7bKB_TSJhs3YvfrfMzSiJ1t48WlbbqIpazjFSZwlkc2TB034jqX_SAJVzjgkajEPmifo-koQUntw17KlbfVzeRM7tywXcpqfc_kYQwhNdbH0r8gNEIlg84rA3WbAvyoo-3SP1UeE",
		"qi": "pCqnuy1DLWTyNXW3pihPkhIXHPMBWBR6kjXGoQ_QG8ig_ukWEs5CVxMXBNid8zOEklzMOCIShK03n8o5U60tAjOztzB_cSoNsKSLLuO2lRQKIOTjfI4I4QY9eY3lvfXCSt40DH6YPXX3fOgy6b52WDNdOdu1BK6AQ5JYgMdLttw"
	  }`)
)

func TestSignedCorim_FromCOSE_ok(t *testing.T) {
	/*
	   18(
	     [
	       / protected h'a10126' / << {
	         / alg / 1: -7, / ECDSA 256 /
	         / content-type / 3: "application/rim+cbor",
	         / issuer-key-id / 4: 'meriadoc.brandybuck@buckland.example',
	         / corim-meta / 8: h'a200a1006941434d45204c74642e01a101c11a5fad2056'
	       } >>,
	       / unprotected / {},
	       / payload / << {
	         0: "test corim id",
	         1: [
	           h'D901FAA40065656E2D474201A1005043BBE37F2E614B33AED353CFF1428B160281A3006941434D45204C74642E01D8207468747470733A2F2F61636D652E6578616D706C65028300010204A1008182A100A300D90258582061636D652D696D706C656D656E746174696F6E2D69642D303030303030303031016441434D45026A526F616452756E6E657283A200D90258A30162424C0465322E312E30055820ACBB11C7E4DA217205523CE4CE1A245AE1A239AE3C6BFD9E7871F7E5D8BAE86B01A102818201582087428FC522803D31065E7BCE3CF03FE475096631E5E07BBD7A0FDE60C4CF25C7A200D90258A3016450526F540465312E332E35055820ACBB11C7E4DA217205523CE4CE1A245AE1A239AE3C6BFD9E7871F7E5D8BAE86B01A10281820158200263829989B6FD954F72BAAF2FC64BC2E2F01D692D4DE72986EA808F6E99813FA200D90258A3016441526F540465302E312E34055820ACBB11C7E4DA217205523CE4CE1A245AE1A239AE3C6BFD9E7871F7E5D8BAE86B01A1028182015820A3A5E715F0CC574A73C3F9BEBB6BC24F32FFD5B67B387244C2C909DA779A1478'
	         ]
	       } >>,
	       / signature / h'deadbeef'
	     ]
	   )
	*/
	tv := []byte{
		0xd2, 0x84, 0x58, 0x59, 0xa4, 0x01, 0x26, 0x03, 0x74, 0x61, 0x70, 0x70,
		0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x72, 0x69, 0x6d,
		0x2b, 0x63, 0x62, 0x6f, 0x72, 0x04, 0x58, 0x24, 0x6d, 0x65, 0x72, 0x69,
		0x61, 0x64, 0x6f, 0x63, 0x2e, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x79, 0x62,
		0x75, 0x63, 0x6b, 0x40, 0x62, 0x75, 0x63, 0x6b, 0x6c, 0x61, 0x6e, 0x64,
		0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x08, 0x57, 0xa2, 0x00,
		0xa1, 0x00, 0x69, 0x41, 0x43, 0x4d, 0x45, 0x20, 0x4c, 0x74, 0x64, 0x2e,
		0x01, 0xa1, 0x01, 0xc1, 0x1a, 0x5f, 0xad, 0x20, 0x56, 0xa0, 0x59, 0x01,
		0xb8, 0xa2, 0x00, 0x6d, 0x74, 0x65, 0x73, 0x74, 0x20, 0x63, 0x6f, 0x72,
		0x69, 0x6d, 0x20, 0x69, 0x64, 0x01, 0x81, 0x59, 0x01, 0xa3, 0xd9, 0x01,
		0xfa, 0xa4, 0x00, 0x65, 0x65, 0x6e, 0x2d, 0x47, 0x42, 0x01, 0xa1, 0x00,
		0x50, 0x43, 0xbb, 0xe3, 0x7f, 0x2e, 0x61, 0x4b, 0x33, 0xae, 0xd3, 0x53,
		0xcf, 0xf1, 0x42, 0x8b, 0x16, 0x02, 0x81, 0xa3, 0x00, 0x69, 0x41, 0x43,
		0x4d, 0x45, 0x20, 0x4c, 0x74, 0x64, 0x2e, 0x01, 0xd8, 0x20, 0x74, 0x68,
		0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x61, 0x63, 0x6d, 0x65, 0x2e,
		0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x02, 0x83, 0x00, 0x01, 0x02,
		0x04, 0xa1, 0x00, 0x81, 0x82, 0xa1, 0x00, 0xa3, 0x00, 0xd9, 0x02, 0x58,
		0x58, 0x20, 0x61, 0x63, 0x6d, 0x65, 0x2d, 0x69, 0x6d, 0x70, 0x6c, 0x65,
		0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2d, 0x69, 0x64,
		0x2d, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x01, 0x64,
		0x41, 0x43, 0x4d, 0x45, 0x02, 0x6a, 0x52, 0x6f, 0x61, 0x64, 0x52, 0x75,
		0x6e, 0x6e, 0x65, 0x72, 0x83, 0xa2, 0x00, 0xd9, 0x02, 0x58, 0xa3, 0x01,
		0x62, 0x42, 0x4c, 0x04, 0x65, 0x32, 0x2e, 0x31, 0x2e, 0x30, 0x05, 0x58,
		0x20, 0xac, 0xbb, 0x11, 0xc7, 0xe4, 0xda, 0x21, 0x72, 0x05, 0x52, 0x3c,
		0xe4, 0xce, 0x1a, 0x24, 0x5a, 0xe1, 0xa2, 0x39, 0xae, 0x3c, 0x6b, 0xfd,
		0x9e, 0x78, 0x71, 0xf7, 0xe5, 0xd8, 0xba, 0xe8, 0x6b, 0x01, 0xa1, 0x02,
		0x81, 0x82, 0x01, 0x58, 0x20, 0x87, 0x42, 0x8f, 0xc5, 0x22, 0x80, 0x3d,
		0x31, 0x06, 0x5e, 0x7b, 0xce, 0x3c, 0xf0, 0x3f, 0xe4, 0x75, 0x09, 0x66,
		0x31, 0xe5, 0xe0, 0x7b, 0xbd, 0x7a, 0x0f, 0xde, 0x60, 0xc4, 0xcf, 0x25,
		0xc7, 0xa2, 0x00, 0xd9, 0x02, 0x58, 0xa3, 0x01, 0x64, 0x50, 0x52, 0x6f,
		0x54, 0x04, 0x65, 0x31, 0x2e, 0x33, 0x2e, 0x35, 0x05, 0x58, 0x20, 0xac,
		0xbb, 0x11, 0xc7, 0xe4, 0xda, 0x21, 0x72, 0x05, 0x52, 0x3c, 0xe4, 0xce,
		0x1a, 0x24, 0x5a, 0xe1, 0xa2, 0x39, 0xae, 0x3c, 0x6b, 0xfd, 0x9e, 0x78,
		0x71, 0xf7, 0xe5, 0xd8, 0xba, 0xe8, 0x6b, 0x01, 0xa1, 0x02, 0x81, 0x82,
		0x01, 0x58, 0x20, 0x02, 0x63, 0x82, 0x99, 0x89, 0xb6, 0xfd, 0x95, 0x4f,
		0x72, 0xba, 0xaf, 0x2f, 0xc6, 0x4b, 0xc2, 0xe2, 0xf0, 0x1d, 0x69, 0x2d,
		0x4d, 0xe7, 0x29, 0x86, 0xea, 0x80, 0x8f, 0x6e, 0x99, 0x81, 0x3f, 0xa2,
		0x00, 0xd9, 0x02, 0x58, 0xa3, 0x01, 0x64, 0x41, 0x52, 0x6f, 0x54, 0x04,
		0x65, 0x30, 0x2e, 0x31, 0x2e, 0x34, 0x05, 0x58, 0x20, 0xac, 0xbb, 0x11,
		0xc7, 0xe4, 0xda, 0x21, 0x72, 0x05, 0x52, 0x3c, 0xe4, 0xce, 0x1a, 0x24,
		0x5a, 0xe1, 0xa2, 0x39, 0xae, 0x3c, 0x6b, 0xfd, 0x9e, 0x78, 0x71, 0xf7,
		0xe5, 0xd8, 0xba, 0xe8, 0x6b, 0x01, 0xa1, 0x02, 0x81, 0x82, 0x01, 0x58,
		0x20, 0xa3, 0xa5, 0xe7, 0x15, 0xf0, 0xcc, 0x57, 0x4a, 0x73, 0xc3, 0xf9,
		0xbe, 0xbb, 0x6b, 0xc2, 0x4f, 0x32, 0xff, 0xd5, 0xb6, 0x7b, 0x38, 0x72,
		0x44, 0xc2, 0xc9, 0x09, 0xda, 0x77, 0x9a, 0x14, 0x78, 0x44, 0xde, 0xad,
		0xbe, 0xef,
	}

	var actual SignedCorim
	err := actual.FromCOSE(tv)

	assert.Nil(t, err)
}

func TestSignedCorim_FromCOSE_fail_no_tag(t *testing.T) {
	// a single null byte is sufficient to test this condition because the tag
	// is the very first thing we stumble upon
	tv := []byte{0xf6}
	var actual SignedCorim
	err := actual.FromCOSE(tv)

	assert.EqualError(t, err, "failed CBOR decoding for COSE-Sign1 signed CoRIM: cbor: invalid COSE_Sign1_Tagged object")
}

func TestSignedCorim_FromCOSE_fail_corim_bad_cbor(t *testing.T) {
	/*
		18(
		  [
		    / protected / << {
		      / alg / 1: -7, / ECDSA 256 /
		      / content-type / 3: "application/rim+cbor",
		      / corim-meta / 8: h'a200a1006941434d45204c74642e01a101c11a5fad2056'
		    } >>,
		    / unprotected / {},
		    / payload / h'badcb030',
		    / signature / h'deadbeef'
		  ]
		)
	*/
	tv := []byte{
		0xd2, 0x84, 0x58, 0x32, 0xa3, 0x01, 0x26, 0x03, 0x74, 0x61, 0x70, 0x70,
		0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x72, 0x69, 0x6d,
		0x2b, 0x63, 0x62, 0x6f, 0x72, 0x08, 0x57, 0xa2, 0x00, 0xa1, 0x00, 0x69,
		0x41, 0x43, 0x4d, 0x45, 0x20, 0x4c, 0x74, 0x64, 0x2e, 0x01, 0xa1, 0x01,
		0xc1, 0x1a, 0x5f, 0xad, 0x20, 0x56, 0xa0, 0x44, 0xba, 0xdc, 0xb0, 0x30,
		0x44, 0xde, 0xad, 0xbe, 0xef,
	}

	var actual SignedCorim
	err := actual.FromCOSE(tv)

	assert.EqualError(t, err, "failed CBOR decoding of unsigned CoRIM: unexpected EOF")
}

func TestSignedCorim_FromCOSE_fail_invalid_corim(t *testing.T) {
	/*
		18(
		  [
		    / protected / << {
		      / alg / 1: -7, / ECDSA 256 /
		      / content-type / 3: "application/rim+cbor",
		      / corim-meta / 8: h'a200a1006941434d45204c74642e01a101c11a5fad2056'
		    } >>,
		    / unprotected / {},
		    / payload / << {
		      0: "invalid corim"
		    } >>,
		    / signature / h'deadbeef'
		  ]
		)
	*/
	tv := []byte{
		0xd2, 0x84, 0x58, 0x32, 0xa3, 0x01, 0x26, 0x03, 0x74, 0x61, 0x70, 0x70,
		0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x72, 0x69, 0x6d,
		0x2b, 0x63, 0x62, 0x6f, 0x72, 0x08, 0x57, 0xa2, 0x00, 0xa1, 0x00, 0x69,
		0x41, 0x43, 0x4d, 0x45, 0x20, 0x4c, 0x74, 0x64, 0x2e, 0x01, 0xa1, 0x01,
		0xc1, 0x1a, 0x5f, 0xad, 0x20, 0x56, 0xa0, 0x50, 0xa1, 0x00, 0x6d, 0x69,
		0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x20, 0x63, 0x6f, 0x72, 0x69, 0x6d,
		0x44, 0xde, 0xad, 0xbe, 0xef,
	}

	var actual SignedCorim
	err := actual.FromCOSE(tv)

	assert.EqualError(t, err, `failed CBOR decoding of unsigned CoRIM: missing mandatory field "Tags" (1)`)
}

func TestSignedCorim_FromCOSE_fail_no_content_type(t *testing.T) {
	/*
	   18(
	     [
	       / protected / << {
	         / alg / 1: -7 / ECDSA 256 /
	       } >>,
	       / unprotected / {},
	       / payload / << {
	         0: "test corim id",
	         1: [ h'cafecafe' ]
	       } >>,
	       / signature / h'deadbeef'
	     ]
	   )
	*/
	tv := []byte{
		0xd2, 0x84, 0x43, 0xa1, 0x01, 0x26, 0xa0, 0x57, 0xa2, 0x00, 0x6d, 0x74,
		0x65, 0x73, 0x74, 0x20, 0x63, 0x6f, 0x72, 0x69, 0x6d, 0x20, 0x69, 0x64,
		0x01, 0x81, 0x44, 0xca, 0xfe, 0xca, 0xfe, 0x44, 0xde, 0xad, 0xbe, 0xef,
	}
	var actual SignedCorim
	err := actual.FromCOSE(tv)

	assert.EqualError(t, err, "processing COSE headers: missing mandatory content type")
}

func TestSignedCorim_FromCOSE_fail_unexpected_content_type(t *testing.T) {
	/*
	   18(
	     [
	       / protected / << {
	         / alg / 1: -7, / ECDSA 256 /
	         / content-type / 3: "application/cbor"
	       } >>,
	       / unprotected / {},
	       / payload / << {
	         0: "test corim id",
	         1: [ h'cafecafe' ]
	       } >>,
	       / signature / h'deadbeef'
	     ]
	   )
	*/
	tv := []byte{
		0xd2, 0x84, 0x55, 0xa2, 0x01, 0x26, 0x03, 0x70, 0x61, 0x70, 0x70, 0x6c,
		0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x63, 0x62, 0x6f, 0x72,
		0xa0, 0x57, 0xa2, 0x00, 0x6d, 0x74, 0x65, 0x73, 0x74, 0x20, 0x63, 0x6f,
		0x72, 0x69, 0x6d, 0x20, 0x69, 0x64, 0x01, 0x81, 0x44, 0xca, 0xfe, 0xca,
		0xfe, 0x44, 0xde, 0xad, 0xbe, 0xef,
	}
	var actual SignedCorim
	err := actual.FromCOSE(tv)

	assert.EqualError(t, err, `processing COSE headers: expecting content type "application/rim+cbor", got "application/cbor" instead`)
}

func unsignedCorimFromCBOR(t *testing.T, cbor []byte) *UnsignedCorim {
	var unsignedCorim UnsignedCorim

	err := unsignedCorim.FromCBOR(cbor)
	require.Nil(t, err)
	require.Nil(t, unsignedCorim.Valid())

	return &unsignedCorim
}

func metaGood(t *testing.T) *Meta {
	var (
		name     = "ACME Ltd."
		notAfter = time.Date(2021, time.October, 0, 0, 0, 0, 0, time.UTC)
	)

	m := NewMeta().
		SetSigner(name, nil).
		SetValidity(notAfter, nil)
	require.NotNil(t, m)

	return m
}

func TestSignedCorim_SignVerify_ok(t *testing.T) {
	for _, key := range [][]byte{
		testES256Key,
		testES384Key,
		testES512Key,
		testEdDSAKey,
		testPS256Key,
		testPS384Key,
		testPS512Key,
	} {
		signer, err := NewSignerFromJWK(key)
		require.NoError(t, err)

		var SignedCorimIn SignedCorim

		SignedCorimIn.UnsignedCorim = *unsignedCorimFromCBOR(t, testGoodUnsignedCorim)
		SignedCorimIn.Meta = *metaGood(t)

		cbor, err := SignedCorimIn.Sign(signer)
		assert.Nil(t, err)

		var SignedCorimOut SignedCorim

		fmt.Printf("signed-corim: %x\n", cbor)

		err = SignedCorimOut.FromCOSE(cbor)
		assert.Nil(t, err)

		pk, err := NewPublicKeyFromJWK(key)
		require.NoError(t, err)

		err = SignedCorimOut.Verify(pk)
		assert.Nil(t, err)
	}
}

func TestSignedCorim_SignVerify_fail_tampered(t *testing.T) {
	signer, err := NewSignerFromJWK(testES256Key)
	require.NoError(t, err)

	var SignedCorimIn SignedCorim

	SignedCorimIn.UnsignedCorim = *unsignedCorimFromCBOR(t, testGoodUnsignedCorim)

	cbor, err := SignedCorimIn.Sign(signer)
	assert.Nil(t, err)

	var SignedCorimOut SignedCorim

	fmt.Printf("signed-corim: %x", cbor)

	// Flip the last byte in the signature field
	cbor[len(cbor)-1] ^= 0xff

	// Since we don't modify the Sign1 payload structurally, decoding the COSE
	// envelope is still OK...
	err = SignedCorimOut.FromCOSE(cbor)
	assert.Nil(t, err)

	pk, err := NewPublicKeyFromJWK(testES256Key)
	require.NoError(t, err)

	// ... but the signature verification fails
	err = SignedCorimOut.Verify(pk)
	assert.EqualError(t, err, "verification error")
}

func TestSignedCorim_Sign_fail_bad_corim(t *testing.T) {
	signer, err := NewSignerFromJWK(testES256Key)
	require.NoError(t, err)

	var SignedCorimIn SignedCorim

	emptyCorim := NewUnsignedCorim()
	require.NotNil(t, emptyCorim)

	SignedCorimIn.UnsignedCorim = *emptyCorim

	_, err = SignedCorimIn.Sign(signer)
	assert.EqualError(t, err, "failed validation of unsigned CoRIM: empty id")
}

func TestSignedCorim_Sign_fail_no_signer(t *testing.T) {
	var SignedCorimIn SignedCorim

	emptyCorim := NewUnsignedCorim()
	require.NotNil(t, emptyCorim)

	SignedCorimIn.UnsignedCorim = *emptyCorim

	_, err := SignedCorimIn.Sign(nil)
	assert.EqualError(t, err, "nil signer")
}
