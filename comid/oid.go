// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// BER-encoded absolute OID
type OID []byte

func asn1OIDFromDottedDecimalString(s string) (asn1.ObjectIdentifier, error) {
	if s == "" {
		return nil, fmt.Errorf("empty OID")
	}

	if s[0] == '.' {
		return nil, fmt.Errorf("OID must be absolute")
	}

	var asn1OID asn1.ObjectIdentifier

	for _, s := range strings.Split(s, ".") {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid OID: %w", err)
		}
		if n < 0 {
			return nil, fmt.Errorf("invalid OID: negative arc %d not allowed", n)
		}
		asn1OID = append(asn1OID, n)
	}

	if len(asn1OID) < 3 {
		return nil,
			fmt.Errorf(
				"invalid OID: got %d arcs, expecting at least %d",
				len(asn1OID), 3,
			)
	}

	return asn1OID, nil
}

func (o *OID) FromString(s string) error {
	// decode absolute OID in dotted-decimal format to internal representation
	asn1OID, err := asn1OIDFromDottedDecimalString(s)
	if err != nil {
		return err
	}

	// encode internal representation to BER
	berOID, err := asn1.Marshal(asn1OID)
	if err != nil {
		return err
	}

	// drop T&L and keep the value
	oidVal, err := extractBERValue(berOID)
	if err != nil {
		return err
	}

	*o = OID(oidVal)

	return nil
}

func (o OID) String() string {
	var asn1OID asn1.ObjectIdentifier

	tlv, err := constructBERFromVal(o)
	if err != nil {
		return ""
	}

	if _, err := asn1.Unmarshal(tlv, &asn1OID); err != nil {
		return ""
	}

	return asn1OID.String()
}

const (
	asn1AbsoluteOIDType = 0x06
	asn1LongLenMask     = 0x80
	asn1LenBytesMask    = 0x7F
)

func extractBERValue(asn1OID []byte) ([]byte, error) {
	if asn1OID[0] != asn1AbsoluteOIDType {
		return nil, fmt.Errorf("the supplied value is not an ASN.1 OID")
	}

	byteOffset := 2

	if asn1OID[1]&asn1LongLenMask != 0 {
		byteOffset += int(asn1OID[1] & asn1LenBytesMask)
	}

	return asn1OID[byteOffset:], nil
}

const (
	// MaxASN1OIDLen is the maximum OID length accepted by the implementation
	MaxASN1OIDLen = 255
	// MinNumOIDArcs represents the minimum required arcs for a valid OID
	MinNumOIDArcs = 3
)

func constructBERFromVal(val []byte) ([]byte, error) {
	const maxTLOffset = 3
	var OID [MaxASN1OIDLen + maxTLOffset]byte

	berOID := OID[:2]
	berOID[0] = asn1AbsoluteOIDType

	if len(val) < 127 {
		berOID[1] = byte(len(val))
	} else if len(val) <= MaxASN1OIDLen {
		berOID[1] = 1
		berOID[1] |= asn1LongLenMask
		berOID = append(berOID, byte(len(val)))
	} else {
		return nil, fmt.Errorf("OIDs greater than %d bytes are not accepted", MaxASN1OIDLen)
	}

	berOID = append(berOID, val...)

	return berOID, nil
}

func (o *OID) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if err := o.FromString(s); err != nil {
		return err
	}

	return nil
}

func (o OID) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}
