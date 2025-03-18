// Copyright 2024-2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/swid"
)

func prepareRegister(t *testing.T, iType string) (*IntegrityRegisters, error) {
	var err error
	val := MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")
	entry := swid.HashEntry{HashAlgID: swid.Sha256, HashValue: val}
	reg := NewIntegrityRegisters()
	for index := 0; index < 5; index++ {
		switch iType {
		case "uint":
			err = reg.AddDigest(uint(index), entry)
			require.NoError(t, err)
		case "text":
			i := fmt.Sprint(index)
			err = reg.AddDigest(i, entry)
			require.NoError(t, err)
		default:
			err = fmt.Errorf("invalid iType = %s", iType)
		}
		if err != nil {
			return nil, err
		}
	}
	return reg, nil
}

func TestIntegrityRegisters_AddDigest_OK(t *testing.T) {
	_, err := prepareRegister(t, "uint")
	require.NoError(t, err)
}

func TestIntegrityRegisters_AddDigest_NOK(t *testing.T) {
	expectedErr := `no register to add digest`
	register := IntegrityRegisters{}
	err := register.AddDigest(uint(0), swid.HashEntry{})
	assert.EqualError(t, err, expectedErr)
	expectedErr = `unexpected type for index: bool`
	var k bool
	reg, err := prepareRegister(t, "uint")
	require.NoError(t, err)
	err = reg.AddDigest(k, swid.HashEntry{})
	assert.EqualError(t, err, expectedErr)
}

func TestIntegrityRegisters_AddDigests_NOK(t *testing.T) {
	expectedErr := `no digests to add`
	register := IntegrityRegisters{}
	err := register.AddDigests(uint(0), []swid.HashEntry{})
	assert.EqualError(t, err, expectedErr)
}

func TestIntegrityRegisters_MarshalCBOR_UIntIndex_OK(t *testing.T) {
	reg, err := prepareRegister(t, "uint")
	// Below is the partial CBOR Pretty notation to highlight index as uint: unsigned(0)

	// A5                                      # map(5)
	//    00                                   # unsigned(0)
	//    81                                   # array(1)
	//       82                                # array(2)
	//          01                             # unsigned(1)
	//          58 20                          # bytes(32)
	expected := MustHexDecode(t, "a5008182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75018182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75028182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75038182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75048182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")
	require.NoError(t, err)
	actual, err := reg.MarshalCBOR()
	require.NoError(t, err)
	fmt.Printf("CBOR Payload = %x\n", actual)
	assert.Equal(t, expected, actual)
}

func TestIntegrityRegisters_UnmarshalCBOR_UIntIndex_OK(t *testing.T) {
	expected := IntegrityRegisters{map[IRegisterIndex]Digests{
		uint64(0): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(1): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(2): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(3): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(4): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
	}}
	bstr := MustHexDecode(nil, `a5008182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75018182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75028182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75038182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75048182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75`)
	actual := IntegrityRegisters{}
	err := actual.UnmarshalCBOR(bstr)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(expected, actual))
}

func TestIntegrityRegisters_MarshalJSON_UIntIndex_OK(t *testing.T) {
	expected := `{
		"0": {
			"key-type": "uint",
			"value": [
				"sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
			]
		},
		"1": {
			"key-type": "uint",
			"value": [
				"sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
			]
		},
		"2": {
			"key-type": "uint",
			"value": [
				"sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
			]
		},
		"3": {
			"key-type": "uint",
			"value": [
				"sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
			]
		},
		"4": {
			"key-type": "uint",
			"value": [
				"sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
			]
		}
	}
	`
	reg, err := prepareRegister(t, "uint")
	require.NoError(t, err)
	actual, err := reg.MarshalJSON()
	require.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
}

func TestIntegrityRegisters_MarshalCBOR_TextIndex_OK(t *testing.T) {
	reg, err := prepareRegister(t, "text")
	// Below is the partial CBOR Pretty notation to highlight index as text: "0"

	// A5                                   # map(5)
	// 61                                   # text(1)
	//    30                                # "0"
	// 81                                   # array(1)
	//    82                                # array(2)
	// 	      01                             # unsigned(1)
	//  	  58 20                          # bytes(32)
	expected := MustHexDecode(t, `a561308182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7561318182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7561328182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7561338182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7561348182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75`)
	require.NoError(t, err)
	actual, err := reg.MarshalCBOR()
	require.NoError(t, err)
	fmt.Printf("CBOR Payload = %x\n", actual)
	assert.Equal(t, expected, actual)
}

func TestIntegrityRegisters_UnmarshalCBOR_TextIndex_OK(t *testing.T) {
	expected := IntegrityRegisters{map[IRegisterIndex]Digests{
		"0": []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		"1": []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		"2": []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		"3": []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		"4": []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
	}}
	bstr := MustHexDecode(t, `a561308182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7561318182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7561328182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7561338182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7561348182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75`)
	actual := IntegrityRegisters{}
	err := actual.UnmarshalCBOR(bstr)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(expected, actual))
}

func TestIntegrityRegisters_MarshalJSON_TextIndex_OK(t *testing.T) {
	expected := `{
		"0": {
			"key-type": "text",
			"value": [
				"sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
			]
		},
		"1": {
			"key-type": "text",
			"value": [
				"sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
			]
		},
		"2": {
			"key-type": "text",
			"value": [
				"sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
			]
		},
		"3": {
			"key-type": "text",
			"value": [
				"sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
			]
		},
		"4": {
			"key-type": "text",
			"value": [
				"sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="
			]
		}
	}
	`
	reg, err := prepareRegister(t, "text")
	require.NoError(t, err)
	actual, err := reg.MarshalJSON()
	require.NoError(t, err)
	fmt.Printf("JSON Payload = %s", actual)
	assert.JSONEq(t, expected, string(actual))
}

func TestIntegrityRegisters_UnmarshalJSON_TextIndex_OK(t *testing.T) {
	j := `{"abcd":{"key-type":"text","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]}}`
	expected := IntegrityRegisters{map[IRegisterIndex]Digests{
		"abcd": []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
	}}
	actual := IntegrityRegisters{}
	err := actual.UnmarshalJSON([]byte(j))
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(expected, actual))
}

func TestIntegrityRegisters_UnmarshalJSON_UIntIndex_OK(t *testing.T) {
	j := `{
	"0":{"key-type":"uint","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]},
	"1":{"key-type":"uint","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]},
	"2":{"key-type":"uint","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]},
	"3":{"key-type":"uint","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]},
	"4":{"key-type":"uint","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]}
	}`
	expected := IntegrityRegisters{map[IRegisterIndex]Digests{
		uint(0): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint(1): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint(2): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint(3): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint(4): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
	}}
	actual := IntegrityRegisters{}
	err := actual.UnmarshalJSON([]byte(j))
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(expected, actual))
}

func TestIntegrityRegisters_UnmarshalJSON_TextUInt_Index_OK(t *testing.T) {
	j := `{
		"0":{"key-type":"uint","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]},
		"1":{"key-type":"uint","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]},
		"2":{"key-type":"uint","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]},
		"3":{"key-type":"text","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]},
		"4":{"key-type":"text","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]}
		}`
	expected := IntegrityRegisters{
		map[IRegisterIndex]Digests{
			uint(0): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
			uint(1): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
			uint(2): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
			"3":     []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
			"4":     []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		}}
	actual := IntegrityRegisters{}
	err := actual.UnmarshalJSON([]byte(j))
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(expected, actual))
}

func TestIntegrityRegisters_UnmarshalJSON_NOK(t *testing.T) {
	for _, tv := range []struct {
		Name  string
		Input string
		Err   string
	}{
		{
			Name:  "invalid input integer",
			Input: `{"0":{"key-type":"int","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]}}`,
			Err:   "unexpected key type for index: int",
		},
		{
			Name:  "negative index",
			Input: `{"-1":{"key-type":"uint","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]}}`,
			Err:   `invalid negative integer key`,
		},
		{
			Name:  "not an integer",
			Input: `{"0.2345":{"key-type":"uint","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]}}`,
			Err:   `unable to convert key to uint: strconv.Atoi: parsing "0.2345": invalid syntax`,
		},
		{
			Name:  "invalid digest",
			Input: `{"1":{"key-type":"uint","value":["sha-256;5Fty9cDAtXLbTY06t+l/3TmI0eoJN7LZ6hOUiTXU="]}}`,
			Err:   `unable to unmarshal Digests: illegal base64 data at input byte 40`,
		},
	} {
		t.Run(tv.Name, func(t *testing.T) {
			reg := IntegrityRegisters{}
			err := reg.UnmarshalJSON([]byte(tv.Input))
			assert.EqualError(t, err, tv.Err)
		})
	}
}

func TestIntegrityRegisters_Equal_True(t *testing.T) {
	claim := IntegrityRegisters{map[IRegisterIndex]Digests{
		uint64(0): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(1): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(2): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(3): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(4): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
	}}

	ref := IntegrityRegisters{map[IRegisterIndex]Digests{
		uint64(0): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(1): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(2): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(3): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(4): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
	}}

	assert.True(t, claim.Equal(ref))
}

func TestIntegrityRegisters_Equal_False(t *testing.T) {
	claim := IntegrityRegisters{map[IRegisterIndex]Digests{
		uint64(0): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(1): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(2): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(3): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(4): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
	}}

	ref := IntegrityRegisters{map[IRegisterIndex]Digests{
		uint64(0): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(1): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(2): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(3): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
	}}

	assert.False(t, claim.Equal(ref))
}

func TestIntegrityRegisters_Compare_True(t *testing.T) {
	claim := IntegrityRegisters{map[IRegisterIndex]Digests{
		uint64(0): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(1): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "34b3bd704b13febb14eca0a3174114cea735e0c92e70c3d0f5cd78d653e5678b")}},
		uint64(2): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "58af0069d43712309b37d645e6729eca3e5aee9d22bdb595c31b59ee6e2d3750")}},
		uint64(3): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "408d1344f60ec4a06a610406c84cee1d9a5c524b0ddd1264719cc347f4b15a08")}},
		uint64(4): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "9aca8354b65a9b4815cf471a6fe9ca9629389691c4183831e63c37a744b2d8ec")}},
	}}

	ref := IntegrityRegisters{map[IRegisterIndex]Digests{
		uint64(2): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "58af0069d43712309b37d645e6729eca3e5aee9d22bdb595c31b59ee6e2d3750")}},
		uint64(3): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "408d1344f60ec4a06a610406c84cee1d9a5c524b0ddd1264719cc347f4b15a08")}},
	}}

	assert.True(t, claim.CompareAgainstReference(ref))
}

func TestIntegrityRegisters_Compare_False(t *testing.T) {
	claim := IntegrityRegisters{map[IRegisterIndex]Digests{
		uint64(0): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(1): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "34b3bd704b13febb14eca0a3174114cea735e0c92e70c3d0f5cd78d653e5678b")}},
		uint64(2): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "58af0069d43712309b37d645e6729eca3e5aee9d22bdb595c31b59ee6e2d3750")}},
		uint64(3): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "408d1344f60ec4a06a610406c84cee1d9a5c524b0ddd1264719cc347f4b15a08")}},
		uint64(4): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "9aca8354b65a9b4815cf471a6fe9ca9629389691c4183831e63c37a744b2d8ec")}},
	}}

	ref := IntegrityRegisters{map[IRegisterIndex]Digests{
		uint64(4): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
	}}

	assert.False(t, claim.CompareAgainstReference(ref))
}

func TestIntegrityRegisters_Compare_False_MissingEntry(t *testing.T) {
	claim := IntegrityRegisters{map[IRegisterIndex]Digests{
		uint64(0): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
		uint64(1): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "34b3bd704b13febb14eca0a3174114cea735e0c92e70c3d0f5cd78d653e5678b")}},
		uint64(2): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "58af0069d43712309b37d645e6729eca3e5aee9d22bdb595c31b59ee6e2d3750")}},
		uint64(3): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "408d1344f60ec4a06a610406c84cee1d9a5c524b0ddd1264719cc347f4b15a08")}},
	}}

	ref := IntegrityRegisters{map[IRegisterIndex]Digests{
		uint64(4): []swid.HashEntry{{HashAlgID: swid.Sha256, HashValue: MustHexDecode(t, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75")}},
	}}

	assert.False(t, claim.CompareAgainstReference(ref))
}
