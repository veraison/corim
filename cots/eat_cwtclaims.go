package cots

import (
	"encoding/json"
	"fmt"

	"github.com/veraison/eat"
)

type EatCWTClaim struct {
	Nonce         *eat.Nonce         `cbor:"10,keyasint,omitempty" json:"nonce,omitempty"`
	UEID          *eat.UEID          `cbor:"11,keyasint,omitempty" json:"ueid,omitempty"`
	Origination   *eat.StringOrURI   `cbor:"12,keyasint,omitempty" json:"origination,omitempty"`
	OemID         *[]byte            `cbor:"13,keyasint,omitempty" json:"oemid,omitempty"`
	SecurityLevel *eat.SecurityLevel `cbor:"14,keyasint,omitempty" json:"security-level,omitempty"`
	SecureBoot    *bool              `cbor:"15,keyasint,omitempty" json:"secure-boot,omitempty"`
	Debug         *eat.Debug         `cbor:"16,keyasint,omitempty" json:"debug-disable,omitempty"`
	Location      *eat.Location      `cbor:"17,keyasint,omitempty" json:"location,omitempty"`
	Profile       *eat.Profile       `cbor:"18,keyasint,omitempty" json:"eat-profile,omitempty"`
	Uptime        *uint              `cbor:"19,keyasint,omitempty" json:"uptime,omitempty"`
	Submods       *eat.Submods       `cbor:"20,keyasint,omitempty" json:"submods,omitempty"`

	eat.CWTClaims

	// Partial list of claims defined by draft-ietf-rats-eat-12
	HardwareModelLabel    *[]byte              `cbor:"259,keyasint,omitempty" json:"hwmodel,omitempty"`
	HardwareVersionScheme *HardwareVersionType `cbor:"260,keyasint,omitempty" json:"hwvers,omitempty"`

	// numbers for the next two have not yet been assigned
	SoftwareNameLabel     *string              `cbor:"998,keyasint,omitempty" json:"swname,omitempty"`
	SoftwareVersionScheme *HardwareVersionType `cbor:"999,keyasint,omitempty" json:"swversion,omitempty"`
}

// ToCBOR serializes the target EatCWTClaim to CBOR
func (o EatCWTClaim) ToCBOR() ([]byte, error) {
	return em.Marshal(o)
}

// FromCBOR deserializes a CBOR-encoded data into the target EatCWTClaim
func (o *EatCWTClaim) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, o)
}

// ToJSON serializes the target EatCWTClaim to JSON
func (o EatCWTClaim) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

// FromJSON deserializes a JSON-encoded data into the target EatCWTClaim
func (o *EatCWTClaim) FromJSON(data []byte) error {
	return json.Unmarshal(data, o)
}

type EatCWTClaims []EatCWTClaim

func (o EatCWTClaims) Valid() error {
	if len(o) == 0 {
		return fmt.Errorf("empty EatCWTClaims")
	}

	return nil
}

// ToCBOR serializes the target EatCWTClsim to CBOR
func (o EatCWTClaims) ToCBOR() ([]byte, error) {
	return em.Marshal(&o)
}

// FromCBOR deserializes a CBOR-encoded data into the target EatCWTClaim
func (o *EatCWTClaims) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, o)
}

// ToJSON serializes the target EatCWTClaim to JSON
func (o EatCWTClaims) ToJSON() ([]byte, error) {
	if err := o.Valid(); err != nil {
		return nil, err
	}

	return json.Marshal(&o)
}

// FromJSON deserializes a JSON-encoded data into the target EatCWTClaim
func (o *EatCWTClaims) FromJSON(data []byte) error {
	return json.Unmarshal(data, o)
}
