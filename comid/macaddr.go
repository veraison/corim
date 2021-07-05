// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"net"
)

// MACaddr is an HW address (e.g., IEEE 802 MAC-48, EUI-48, EUI-64)
//
// Note: Since TextUnmarshal is not defined on net.HardwareAddr
// (see: https://github.com/golang/go/issues/29678)
// we need to create an alias type with a custom decoder.
type MACaddr net.HardwareAddr

// UnmarshalJSON deserialize a MAC address in textual form into the MACaddr
// target, e.g.:
//   "mac-addr": "00:00:5e:00:53:01"
// or
//   "mac-addr": "02:00:5e:10:00:00:00:01"
func (o *MACaddr) UnmarshalJSON(data []byte) error {
	var s interface{}

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	var i net.HardwareAddr

	if err := jsonDecodeMACaddr(s, &i); err != nil {
		return err
	}

	*o = MACaddr(i)

	return nil
}
