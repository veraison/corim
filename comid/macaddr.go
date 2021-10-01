// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"fmt"
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
//
// Supported formats are IEEE 802 MAC-48, EUI-48, EUI-64, e.g.:
//   00:00:5e:00:53:01
//   00-00-5e-00-53-01
//   02:00:5e:10:00:00:00:01
//   02-00-5e-10-00-00-00-01
func (o *MACaddr) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, err := net.ParseMAC(s)
	if err != nil {
		return fmt.Errorf("bad MAC address %w", err)
	}

	*o = MACaddr(val)

	return nil
}

func (o MACaddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(net.HardwareAddr(o).String())
}
