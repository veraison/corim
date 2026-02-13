package comid

import (
	"bytes"
	"encoding/json"
	"net"
	"testing"
)

func TestMACaddr_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   MACaddr
		shouldFail bool
	}{
		{
			name:     "Valid MAC-48 with colon",
			input:    `"00:00:5e:00:53:01"`,
			expected: MACaddr(net.HardwareAddr{0x00, 0x00, 0x5e, 0x00, 0x53, 0x01}),
		},
		{
			name:     "Valid MAC-48 with dash",
			input:    `"00-00-5e-00-53-01"`,
			expected: MACaddr(net.HardwareAddr{0x00, 0x00, 0x5e, 0x00, 0x53, 0x01}),
		},
		{
			name:     "Valid EUI-64 with colon",
			input:    `"02:00:5e:10:00:00:00:01"`,
			expected: MACaddr(net.HardwareAddr{0x02, 0x00, 0x5e, 0x10, 0x00, 0x00, 0x00, 0x01}),
		},
		{
			name:     "Valid EUI-64 with dash",
			input:    `"02-00-5e-10-00-00-00-01"`,
			expected: MACaddr(net.HardwareAddr{0x02, 0x00, 0x5e, 0x10, 0x00, 0x00, 0x00, 0x01}),
		},
		{
			name:       "Invalid MAC address format",
			input:      `"invalid-mac"`,
			shouldFail: true,
		},
		{
			name:       "Invalid JSON type",
			input:      `12345`,
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mac MACaddr
			err := json.Unmarshal([]byte(test.input), &mac)
			if test.shouldFail {
				if err == nil {
					t.Errorf("expected failure but got no error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// Compare the MAC address slices using bytes.Equal
				if !bytes.Equal([]byte(mac), []byte(test.expected)) {
					t.Errorf("expected %v, got %v", test.expected, mac)
				}
			}
		})
	}
}

func TestMACaddr_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    MACaddr
		expected string
	}{
		{
			name:     "Valid MAC-48",
			input:    MACaddr(net.HardwareAddr{0x00, 0x00, 0x5e, 0x00, 0x53, 0x01}),
			expected: `"00:00:5e:00:53:01"`,
		},
		{
			name:     "Valid EUI-64",
			input:    MACaddr(net.HardwareAddr{0x02, 0x00, 0x5e, 0x10, 0x00, 0x00, 0x00, 0x01}),
			expected: `"02:00:5e:10:00:00:00:01"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(data) != test.expected {
				t.Errorf("expected %s, got %s", test.expected, string(data))
			}
		})
	}
}
