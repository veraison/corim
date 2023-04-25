package cots

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/eat"
)

var (
	ueID = eat.UEID{
		0x01, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	}
	oemID      = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	nonceBytes = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	AcmeInc     = "Acme Inc."
	origination = getACMEStringToURI(AcmeInc)
	secLevel    = eat.SecurityLevel(eat.SecLevelHardware)
	secBoot     = true
	debug       = eat.Debug(eat.DebugDisabled)
	location    = eat.Location{Latitude: 12.34, Longitude: 56.78}
	uptime      = uint(60)
	issuer      = AcmeInc
	subject     = "rr-trap"
	audience    = eat.Audience{origination}
	epoch       = eat.NumericDate(time.Unix(0, 0))

	nonce          = getNonce()
	fatEatCWTClaim = EatCWTClaim{
		Nonce:         &nonce,
		UEID:          &ueID,
		Origination:   &origination,
		OemID:         &oemID,
		SecurityLevel: &secLevel,
		SecureBoot:    &secBoot,
		Debug:         &debug,
		Location:      &location,
		Uptime:        &uptime,

		CWTClaims: eat.CWTClaims{
			Issuer:     &issuer,
			Subject:    &subject,
			Audience:   &audience,
			Expiration: &epoch,
			NotBefore:  &epoch,
			IssuedAt:   &epoch,
			CwtID:      &oemID,
		},
	}
)

func getACMEStringToURI(currStr string) eat.StringOrURI {
	acmeString := eat.StringOrURI{}
	_ = acmeString.FromString(currStr)
	return acmeString
}

func getNonce() eat.Nonce {
	nonce := eat.Nonce{}
	_ = nonce.Add(nonceBytes)
	return nonce
}

func cborRoundTripper(t *testing.T, tv EatCWTClaim, expected []byte) {
	data, err := tv.ToCBOR()

	t.Logf("CBOR: %x", data)

	assert.Nil(t, err)
	assert.Equal(t, expected, data)

	actual := EatCWTClaim{}
	err = actual.FromCBOR(data)

	assert.Nil(t, err)
	assert.Equal(t, tv, actual)
}

func jsonRoundTripper(t *testing.T, tv EatCWTClaim, expected string) {
	data, err := tv.ToJSON()

	t.Logf("JSON: '%s'", string(data))

	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(data))

	actual := EatCWTClaim{}
	err = actual.FromJSON(data)

	assert.Nil(t, err)
	assert.Equal(t, tv, actual)
}

func TestEatCWTClaim_Full_RoundtripCBOR(t *testing.T) {
	tv := fatEatCWTClaim

	expected := []byte{
		0xb0, 0xa, 0x48, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xb,
		0x51, 0x1, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xc, 0x69,
		0x41, 0x63, 0x6d, 0x65, 0x20, 0x49, 0x6e, 0x63, 0x2e, 0xd,
		0x46, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xe, 0x3, 0xf,
		0xf5, 0x10, 0x1, 0x11, 0xa2, 0x1, 0xfb, 0x40, 0x28, 0xae,
		0x14, 0x7a, 0xe1, 0x47, 0xae, 0x2, 0xfb, 0x40, 0x4c, 0x63,
		0xd7, 0xa, 0x3d, 0x70, 0xa4, 0x13, 0x18, 0x3c, 0x1, 0x69,
		0x41, 0x63, 0x6d, 0x65, 0x20, 0x49, 0x6e, 0x63, 0x2e, 0x2,
		0x67, 0x72, 0x72, 0x2d, 0x74, 0x72, 0x61, 0x70, 0x3, 0x69,
		0x41, 0x63, 0x6d, 0x65, 0x20, 0x49, 0x6e, 0x63, 0x2e, 0x4,
		0xc1, 0x0, 0x5, 0xc1, 0x0, 0x6, 0xc1, 0x0, 0x7, 0x46, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff,
	}

	cborRoundTripper(t, tv, expected)
}

func TestEatCWTClaim_Full_RoundtripJSON(t *testing.T) {
	tv := fatEatCWTClaim
	expected := `
{
	"nonce": "AAAAAAAAAAA=",
	"origination": "Acme Inc.",
	"oemid": "////////",
	"security-level": 3,
	"secure-boot": true,
	"debug-disable": 1,
	"location": {
		"lat": 12.34,
		"long": 56.78
	},
	"ueid": "Ad6tvu/erb7v3q2+796tvu8=",
	"uptime": 60,
	"iss": "Acme Inc.",
	"sub": "rr-trap",
	"aud": "Acme Inc.",
	"exp": 0,
	"nbf": 0,
	"iat": 0,
	"cti": "////////"
}`
	jsonRoundTripper(t, tv, expected)
}

func TestEatCWTClaims_Valid_empty_list(t *testing.T) {
	tv := EatCWTClaims{}
	err := tv.Valid()

	assert.EqualError(t, err, "empty EatCWTClaims")
}
