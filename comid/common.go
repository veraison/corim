// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/base64"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/veraison/eat"
)

// UUID string in canonical 8-4-4-4-12 format
func jsonDecodeUUID(from interface{}, to *UUID) error {
	s, ok := from.(string)
	if !ok {
		return fmt.Errorf("UUID must be string")
	}

	u, err := uuid.Parse(s)
	if err != nil {
		return fmt.Errorf("bad UUID: %w", err)
	}

	*to = UUID(u)

	return nil
}

// (absolute) OID string in dotted decimal notation
func jsonDecodeOID(from interface{}, to *OID) error {
	s, ok := from.(string)
	if !ok {
		return fmt.Errorf("OID must be string")
	}

	var oid OID

	err := oid.FromString(s)
	if err != nil {
		return fmt.Errorf("decoding %s: %w", s, err)
	}

	*to = oid

	return nil
}

// Implementation ID as base64-encoded string
func jsonDecodeImplID(from interface{}, to *ImplID) error {
	s, ok := from.(string)
	if !ok {
		return fmt.Errorf("ImplID must be string") // nolint: golint
	}

	val, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("bad ImplID: %w", err)
	}

	if nb := len(val); nb != 32 {
		return fmt.Errorf("bad ImplID format: got %d bytes, want 32", nb)
	}

	copy(to[:], val)

	return nil
}

// UEID (bstr .size (7..33)) as base64-encoded string
func jsonDecodeUEID(from interface{}, to *eat.UEID) error {
	s, ok := from.(string)
	if !ok {
		return fmt.Errorf("UEID must be string")
	}

	val, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("bad UUID: %w", err)
	}

	ueid := eat.UEID(val)

	if err := ueid.Validate(); err != nil {
		return err
	}

	*to = ueid

	return nil
}

// Supported formats are IEEE 802 MAC-48, EUI-48, EUI-64, e.g.:
//   00:00:5e:00:53:01
//   00-00-5e-00-53-01
//   02:00:5e:10:00:00:00:01
//   02-00-5e-10-00-00-00-01
func jsonDecodeMACaddr(from interface{}, to *net.HardwareAddr) error {
	s, ok := from.(string)
	if !ok {
		return fmt.Errorf("MAC address must be string")
	}

	val, err := net.ParseMAC(s)
	if err != nil {
		return fmt.Errorf("bad MAC address %w", err)
	}

	*to = val

	return nil
}
