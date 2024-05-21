// Copyright 2021-2023 Contributors to the Veraison project.
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
	"github.com/veraison/swid"
)

var (
	TestUUIDString = "31fb5abf-023e-4992-aa4e-95f9c1503bfa"
	TestUUID       = UUID(uuid.Must(uuid.Parse(TestUUIDString)))
	TestImplID     = ImplID([32]byte{
		0x61, 0x63, 0x6d, 0x65, 0x2d, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65,
		0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2d, 0x69, 0x64, 0x2d, 0x30,
		0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31,
	})
	TestOID               = "2.5.2.8192"
	TestRegID             = "https://acme.example"
	TestMACaddr, _        = net.ParseMAC("02:00:5e:10:00:00:00:01")
	TestIPaddr            = net.ParseIP("2001:db8::68")
	TestBytes             = []byte{0x89, 0x99, 0x78, 0x65, 0x56}
	TestUEIDString        = "02deadbeefdead"
	TestUEID              = eat.UEID(MustHexDecode(nil, TestUEIDString))
	TestSignerID          = MustHexDecode(nil, "acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b")
	TestTagID             = "urn:example:veraison"
	TestMKey       uint64 = 700
	TestCCALabel          = "cca-platform-config"

	// #nosec G101
	TestECPrivKey = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEICAm3+mCCDTMuzKqfZso9NT8ur9U9GjuUQ/lNEJvwRFMoAoGCCqGSM49
AwEHoUQDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8BlLT4MFHOaO+ICTtIvrEeEpr/
sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==
-----END EC PRIVATE KEY-----`

	TestECPubKey = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B
lLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==
-----END PUBLIC KEY-----`

	TestCert = `-----BEGIN CERTIFICATE-----
MIIB4TCCAYegAwIBAgIUGhrA9M3yQIFqckA2v6nQewmF30IwCgYIKoZIzj0EAwIw
RTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGElu
dGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAgFw0yMzA5MDQxMTAxNDhaGA8yMDUxMDEx
OTExMDE0OFowRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAf
BgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDBZMBMGByqGSM49AgEGCCqG
SM49AwEHA0IABFtQb6hfv68vAVmu2RDFNcWGBxEPAZS0+DBRzmjviAk7SL6xHhKa
/7H0wD+uh9oQh3W+RxFyrUSqHekC3DizwNWjUzBRMB0GA1UdDgQWBBQWpNPb6eWD
SM/+jwpbzoO3iHg4LTAfBgNVHSMEGDAWgBQWpNPb6eWDSM/+jwpbzoO3iHg4LTAP
BgNVHRMBAf8EBTADAQH/MAoGCCqGSM49BAMCA0gAMEUCIAayNIF0eCJDZmcrqRjH
f9h8GxeIDUnLqldeIvNfa+9SAiEA9ULBTPjnTUhYle226OAjg2sdhkXtb3Mu0E0F
nuUmsIQ=
-----END CERTIFICATE-----`

	TestCertPath = `-----BEGIN CERTIFICATE-----
MIICejCCAiygAwIBAgIUIpeVwVhN/qYLgtNJlwZHJj+IT/wwBQYDK2VwMDMxMTAv
BgNVBAUTKDdhMDZlZWU0MWI3ODlmNDg2M2Q4NmI4Nzc4YjFhMjAxYTZmZWRkNTYw
IBcNMTgwMzIyMjM1OTU5WhgPOTk5OTEyMzEyMzU5NTlaMDMxMTAvBgNVBAUTKDIy
OTc5NWMxNTg0ZGZlYTYwYjgyZDM0OTk3MDY0NzI2M2Y4ODRmZmMwKjAFBgMrZXAD
IQAVUi7xVynM85UJ6lwVomvpSeOIB6XCbvkoFIfvSuZ7RqOCAU4wggFKMB8GA1Ud
IwQYMBaAFHoG7uQbeJ9IY9hrh3ixogGm/t1WMB0GA1UdDgQWBBQil5XBWE3+pguC
00mXBkcmP4hP/DAOBgNVHQ8BAf8EBAMCAgQwDwYDVR0TAQH/BAUwAwEB/zCB5gYK
KwYBBAHWeQIBGAEB/wSB1DCB0aBCBEAMtveoAJx4v3Qu9K8gI0J2kp7igBv4Vd9B
HEMtlfZOQdy6yDOlcPrQGYjSzwvAjFh8DA46bq+xGj31NFUy6pkho0IEQNCkl/Kf
bESLZ6OhEcdOnAziS5gx5TqJmF22yKCjIvLRIVxNIhZN2EnMAtm4dp1EGuPBLUrA
tzXhlzuZuK1xV6SkQgRASeSoHmnNgLhnnEKTWKzcL2jzPjOAFUQTNRy+iOghluI3
ficG6NB7cMbLAZkfV12lihSV+/7iK3TJ0bUNjQgWpaYDCgEBMAUGAytlcANBAHu6
DtuPNOurcAXc+41QY23hY8KRkBCKCE7phsiIwRfbxKMLldFGN5OytQfROQaWoAcv
IWTqV9JRzGQaGYnlLwE=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIICejCCAiygAwIBAgIUEbae7u1cP7s+G/CDb2Nu6nPRdYkwBQYDK2VwMDMxMTAv
BgNVBAUTKDIyOTc5NWMxNTg0ZGZlYTYwYjgyZDM0OTk3MDY0NzI2M2Y4ODRmZmMw
IBcNMTgwMzIyMjM1OTU5WhgPOTk5OTEyMzEyMzU5NTlaMDMxMTAvBgNVBAUTKDEx
YjY5ZWVlZWQ1YzNmYmIzZTFiZjA4MzZmNjM2ZWVhNzNkMTc1ODkwKjAFBgMrZXAD
IQDzgkTR7uvoP9NBzSEB9gu/lpd+NL38OVYQl0feiWKX+aOCAU4wggFKMB8GA1Ud
IwQYMBaAFCKXlcFYTf6mC4LTSZcGRyY/iE/8MB0GA1UdDgQWBBQRtp7u7Vw/uz4b
8INvY27qc9F1iTAOBgNVHQ8BAf8EBAMCAgQwDwYDVR0TAQH/BAUwAwEB/zCB5gYK
KwYBBAHWeQIBGAEB/wSB1DCB0aBCBEBbsKgcavE+uy1AxkIdl7lN9ifHy3HE+Liu
8lME27CMY9kUtyw/les1H8vpmSyxhO4aTWgwuwQa7Yn9HoGweEHso0IEQHb870HN
1bUn9nFih11SBAj9lobpuJ5GrI/m+g6HwmoQz5Uly0oXMNnxEMA7fL2za01ynGpI
/uz82rUI2vLWSlGkQgRA3XFgIoVImosdAgvuPHVaobv3JGjGl3+ADOT1c6dT6dQE
dnObRNudY8qhzTvfEWR4eS6OJtfyrOeRyXek2OVJh6YDCgEBMAUGAytlcANBAJIj
yFqwdrZCSuYmC4+ZUUcANKQKA1KcRFiIlKcg/ppwKVykPXbAhsn6SCVqWGA7v7Ce
Li5hOrH/VljAQAcdYgc=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIICejCCAiygAwIBAgIUYUXMcfV3Mucpunl193wuXxBBr/gwBQYDK2VwMDMxMTAv
BgNVBAUTKDExYjY5ZWVlZWQ1YzNmYmIzZTFiZjA4MzZmNjM2ZWVhNzNkMTc1ODkw
IBcNMTgwMzIyMjM1OTU5WhgPOTk5OTEyMzEyMzU5NTlaMDMxMTAvBgNVBAUTKDYx
NDVjYzcxZjU3NzMyZTcyOWJhNzk3NWY3N2MyZTVmMTA0MWFmZjgwKjAFBgMrZXAD
IQBJo9PgveHj0ahv8MkWHQUGSxZ/wSTdaNNZbdBZNa1L0aOCAU4wggFKMB8GA1Ud
IwQYMBaAFBG2nu7tXD+7Phvwg29jbupz0XWJMB0GA1UdDgQWBBRhRcxx9Xcy5ym6
eXX3fC5fEEGv+DAOBgNVHQ8BAf8EBAMCAgQwDwYDVR0TAQH/BAUwAwEB/zCB5gYK
KwYBBAHWeQIBGAEB/wSB1DCB0aBCBEA2ansef0SbRN8j76w5hzW5/TCXFIsQcERs
bSKQYNnqug1rjECPnhe3A/8Z6WGxaDK1ehE+nrcvC9BRgrWpU67Jo0IEQIZyRCHK
9HUi/8y6V9P0ZuNEvmdpEdImQ09RU/lNPsXXxyv0VEmi6WDs4eFypmBR9LVXBXud
rCduuvyS6tBWsS6kQgRAbWRTCbXrd/qlLPII85IPB8pZ9uX+XgIHI4sSHf+3F6se
hA/80zUBzSi6Ozc0D+IbYYBYxdrXZEkn8iUWSdQokKYDCgEBMAUGAytlcANBAKlJ
/3VYalZm9XbEGTKrVRaoCVoUxQVH3udMrk9yoqjFowC4e3kdSBlGGf8mYEI7xvsA
ar1kf2bGXT/cEeFGIwM=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIICejCCAiygAwIBAgIUX+ivPHTOmvVktMnQGYQjuNlk/DUwBQYDK2VwMDMxMTAv
BgNVBAUTKDYxNDVjYzcxZjU3NzMyZTcyOWJhNzk3NWY3N2MyZTVmMTA0MWFmZjgw
IBcNMTgwMzIyMjM1OTU5WhgPOTk5OTEyMzEyMzU5NTlaMDMxMTAvBgNVBAUTKDVm
ZThhZjNjNzRjZTlhZjU2NGI0YzlkMDE5ODQyM2I4ZDk2NGZjMzUwKjAFBgMrZXAD
IQC6u3blwE4B1xdPMeUJP657P/m7iSt+HergvGbkkSxMrqOCAU4wggFKMB8GA1Ud
IwQYMBaAFGFFzHH1dzLnKbp5dfd8Ll8QQa/4MB0GA1UdDgQWBBRf6K88dM6a9WS0
ydAZhCO42WT8NTAOBgNVHQ8BAf8EBAMCAgQwDwYDVR0TAQH/BAUwAwEB/zCB5gYK
KwYBBAHWeQIBGAEB/wSB1DCB0aBCBEC4z2juJIx5jD7x6IuMNUi7TUomWxCQf9Qn
CJ91ozXk0vJ9nJO3RdveJvbvZhoPfDQIY8TiZp8UKDx4e+zW0cHko0IEQOhpMJ6G
EXLZgHtRAm81oXXACEF+nev2MCv6COhuRtFypG9B3foRm2rnFUbaVZs0pLfBMG8s
sSRJRcawXCimW4OkQgRA27Fgx7A4212qpqLaxaPd9tI+zpfKWrLYcLx20+DLfcqn
BIIpUCN30SuAu71se4x/ilcKuaWOO0qDg34SJEwFyqYDCgEBMAUGAytlcANBACtT
5Xrx659qGnywmlKHdlHO6Bd7fPboyzyIQhoEtFNuiD3WjDg/Vwz8cNCUkU+thG7f
C+WZhcpAckDldai+PAc=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIICejCCAiygAwIBAgIUXX0DVOylgEip7rzvyaNFegcDTZIwBQYDK2VwMDMxMTAv
BgNVBAUTKDVmZThhZjNjNzRjZTlhZjU2NGI0YzlkMDE5ODQyM2I4ZDk2NGZjMzUw
IBcNMTgwMzIyMjM1OTU5WhgPOTk5OTEyMzEyMzU5NTlaMDMxMTAvBgNVBAUTKDVk
N2QwMzU0ZWNhNTgwNDhhOWVlYmNlZmM5YTM0NTdhMDcwMzRkOTIwKjAFBgMrZXAD
IQCiaC2gHhMO1pbeQbUgLHhSgFBPD/zXNAGwAHsW272+c6OCAU4wggFKMB8GA1Ud
IwQYMBaAFF/orzx0zpr1ZLTJ0BmEI7jZZPw1MB0GA1UdDgQWBBRdfQNU7KWASKnu
vO/Jo0V6BwNNkjAOBgNVHQ8BAf8EBAMCAgQwDwYDVR0TAQH/BAUwAwEB/zCB5gYK
KwYBBAHWeQIBGAEB/wSB1DCB0aBCBEBYXwS/mrrX+D4MqzM8JTmIHC9XHqsJfOGc
b2fqBYPX0UQriLDRl1apHN22q1E+FeaLHWBE2uXda1Q6lYkQAaHio0IEQGKH8EAN
Mv1PMMbWsddZew2G/DR+A9tbSi7H680yBSe9Ce+gtabBarQDHpg9B8LebmoPpdXt
ATv+oSzzk+ZueVKkQgRAztbU2QzaJbcG5twEYjYAgFutCbngpg2t/2ez7QTNn4Nm
r94pOAx8LIpu6Cf/Wzcvd/4kLOvWxSb/buMqbGvrsqYDCgEBMAUGAytlcANBAMdx
fFnk52ru4fV5J1gsFBtlhy5mFbQnIRiGHLxBFaTJk9i+ixO5qFbRjqv7HQ/jGUsI
srUJQ4e2JEnSNakNcgE=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIICejCCAiygAwIBAgIUA1oWoWPwVdC5GO3bQoD3roDe1SowBQYDK2VwMDMxMTAv
BgNVBAUTKDVkN2QwMzU0ZWNhNTgwNDhhOWVlYmNlZmM5YTM0NTdhMDcwMzRkOTIw
IBcNMTgwMzIyMjM1OTU5WhgPOTk5OTEyMzEyMzU5NTlaMDMxMTAvBgNVBAUTKDAz
NWExNmExNjNmMDU1ZDBiOTE4ZWRkYjQyODBmN2FlODBkZWQ1MmEwKjAFBgMrZXAD
IQDSNEY1gbLMNAOC+3eok+RyQ6fhN8F23o2dx61QbsM0TqOCAU4wggFKMB8GA1Ud
IwQYMBaAFF19A1TspYBIqe6878mjRXoHA02SMB0GA1UdDgQWBBQDWhahY/BV0LkY
7dtCgPeugN7VKjAOBgNVHQ8BAf8EBAMCAgQwDwYDVR0TAQH/BAUwAwEB/zCB5gYK
KwYBBAHWeQIBGAEB/wSB1DCB0aBCBEA9mPpAmW+IEOXXOBgSy3ry53I562D9OZHZ
+DG1/M9mWxiUkRA1UciqMpGg6ngyqp38J5OpUIuFsoSVDqFVPyjxo0IEQIG/7h17
Am77GLmQ1nSMBZjtrJ+FrmWTcZjxJ9cX0CPJqu7wugL5Tcj1I8c9nBNqsokFx8pE
tRoqiz7rt6Z52D2kQgRAZrvqFdyj4rVcjtVkJbMlp/8jmfGeaKh/RG64WrK2uNk9
yhKOpkiQR0p5UsTam+XdEvqrxjLa43sr0di/pKEbZqYDCgEBMAUGAytlcANBAOuQ
qXZU521LzDTXXx2EYqVuWCyUZIJZgRl/JGs2RmYPYJCZun0Kj1YIvX5mBZ3pC85w
0fJFmM1B2+ACsp+p6Qg=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIICejCCAiygAwIBAgIUS3VU752sqDUfY60E/hEqSn142AUwBQYDK2VwMDMxMTAv
BgNVBAUTKDAzNWExNmExNjNmMDU1ZDBiOTE4ZWRkYjQyODBmN2FlODBkZWQ1MmEw
IBcNMTgwMzIyMjM1OTU5WhgPOTk5OTEyMzEyMzU5NTlaMDMxMTAvBgNVBAUTKDRi
NzU1NGVmOWRhY2E4MzUxZjYzYWQwNGZlMTEyYTRhN2Q3OGQ4MDUwKjAFBgMrZXAD
IQB/oGXT67ucYx9lpxFZFRYvtgmCyBH22i/LnUN0KF6LsaOCAU4wggFKMB8GA1Ud
IwQYMBaAFANaFqFj8FXQuRjt20KA966A3tUqMB0GA1UdDgQWBBRLdVTvnayoNR9j
rQT+ESpKfXjYBTAOBgNVHQ8BAf8EBAMCAgQwDwYDVR0TAQH/BAUwAwEB/zCB5gYK
KwYBBAHWeQIBGAEB/wSB1DCB0aBCBEDJ6W15aipyC2UvMiq6IC2/wFkvsFc9POrT
1NngZGfke8JlnO78VRUZcsF7uhtqyreyjiq5iZHS9hM0J5vIOxujo0IEQJcZ78al
nCtJiWqCHjTgGZjoW+lQJjJ9Ux50TTxReEp3eEOD9O3t4ygdSH4rTFuiuL6tdlZ8
rC/0KTC4G5vEowGkQgRAgQPIBQemZ1isQoF5rKpPotpHXN8YYxGY5WFQIzk9dz7P
zxInQ1qnGAsjQPSS9+JMywDDAi7XKuFwRf0Wk2T9TaYDCgEBMAUGAytlcANBAFUl
UrTQ5qpCcBfPGeTacXNwl5y3WTFgpjFKr+Mw6qusj+bdZ6l+N3CxvOxJ9m+i96Mx
rpT6kiSnAzk+2zgSiA4=
-----END CERTIFICATE-----`

	TestCOSEKey = MustHexDecode(nil, `a501020258246d65726961646f632e6272616e64796275636b406275636b6c616e642e6578616d706c65200121582065eda5a12577c2bae829437fe338701a10aaa375e1bb5b5de108de439c08551d2258201e52ed75701163f7f9e40ddf9f341b3dc9ba860af7e0ca7ca7e9eecd0084d19c`)

	TestCOSEKeySetOne = MustHexDecode(nil, `81a501020258246d65726961646f632e6272616e64796275636b406275636b6c616e642e6578616d706c65200121582065eda5a12577c2bae829437fe338701a10aaa375e1bb5b5de108de439c08551d2258201e52ed75701163f7f9e40ddf9f341b3dc9ba860af7e0ca7ca7e9eecd0084d19c`)

	TestCOSEKeySetMulti = MustHexDecode(nil, `82a501020258246d65726961646f632e6272616e64796275636b406275636b6c616e642e6578616d706c65200121582065eda5a12577c2bae829437fe338701a10aaa375e1bb5b5de108de439c08551d2258201e52ed75701163f7f9e40ddf9f341b3dc9ba860af7e0ca7ca7e9eecd0084d19ca601010327048202647369676e0543030201200621582015522ef15729ccf39509ea5c15a26be949e38807a5c26ef9281487ef4ae67b46`)

	TestThumbprint = swid.HashEntry{
		HashAlgID: 1,
		HashValue: MustHexDecode(nil, `68e656b251e67e8358bef8483ab0d51c6619f3e7a1a9f0e75838d41ff368f728`),
	}
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
	var implID = TestImplID[:]

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
						"type": "pkix-base64-key",
						"value": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"
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
						"type": "pkix-base64-key",
						"value": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8B\nlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q==\n-----END PUBLIC KEY-----"
					}
				]
			}
		]
	}
}
`
	CCARefValJSONTemplate = `{
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
					},
					{
						"key": {
							"type": "cca.platform-config-id",
							"value": "a non-empty (unique) label"
						},
						"value": {
							"raw-value": {
								"type": "bytes",
								"value": "cmF3dmFsdWUKcmF3dmFsdWUK"
							}
						}
					}
				]
			}
		]
	}
}
`
	CCARealmRefValJSONTemplate = `{
	"lang": "en-GB",
	"tag-identity": {
	"id": "99019224-57AA-44BC-BEF8-D36BDD6BD035",
	"version": 0
},
"entities": [
	{
		"name": "Workload Client Ltd.",
		"regid": "https://workloadclient.example",
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
						"type": "uuid",
						"value": "CD1F0E55-26F9-460D-B9D8-F7FDE171787C"
					},
					"vendor": "Workload Client Ltd"
				},
				"instance": {
					"type": "bytes",
					"value": "QoS1aUymwNLPR4mguVrIAlyBjeUjBDZL580pgbLS7caFsyInfsJYGZYkE9jJssH1"
				}
			},
			"measurements": [
				{
					"value": {
						"integrity-registers": {
							"rim": {
								"key-type": "text",
								"value": [
									"sha-384;QoS1aUymwNLPR4mguVrIAlyBjeUjBDZL580pgbLS7caFsyInfsJYGZYkE9jJssH1"
								]
							},
							"rem0": {
								"key-type": "text",
								"value": [
									"sha-384;IQe752H8pS2VE2oTVNt6TdV7Gya+DT2nHZ6yOYazS6YVq/ZRTPNeWp6lWgMtBop4"
								]
							},
							"rem1": {
								"key-type": "text",
								"value": [
									"sha-384;JQe752H8pS2VE2oTVNt6TdV7Gya+DT2nHZ6yOYazS6YVq/ZRTPNeWp6lWgMtBop4"
								]
							},
							"rem2": {
								"key-type": "text",
								"value": [
									"sha-384;MQe752H8pS2VE2oTVNt6TdV7Gya+DT2nHZ6yOYazS6YVq/ZRTPNeWp6lWgMtBop4"
								]
							},
							"rem3": {
								"key-type": "text",
								"value": [
									"sha-384;NQe752H8pS2VE2oTVNt6TdV7Gya+DT2nHZ6yOYazS6YVq/ZRTPNeWp6lWgMtBop4"
								]
							}
						}
					}
				}
			]
		}
	]
}
}`
)
