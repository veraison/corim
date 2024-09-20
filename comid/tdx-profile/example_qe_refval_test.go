// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	_ "embed"
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/corim"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
	"github.com/veraison/swid"
)

// Example_decode_QE_JSON decodes the TDX Quoting Enclave Measurement Extensions from the given JSON Template
func Example_decode_QE_JSON() {
	profileID, err := eat.NewProfile("2.16.840.1.113741.1.16.1")
	if err != nil {
		panic(err) // will not error, as the hard-coded string above is valid
	}
	manifest, found := corim.GetProfileManifest(profileID)
	if !found {
		fmt.Printf("CoRIM Profile NOT FOUND")
		return
	}

	m := manifest.GetComid()
	if err := m.FromJSON([]byte(TDXQERefValTemplate)); err != nil {
		panic(err)
	}

	if err := m.Valid(); err != nil {
		panic(err)
	}

	if err := extractQERefVals(m); err != nil {
		panic(err)
	}

	// output:
	// OID: 2.16.840.1.113741.1.2.3.4.1
	// Vendor: Intel Corporation
	// Model: TDX QE TCB
	// miscselect: c0000000fbff0000
	// tcbEvalNum: 11
	// IsvProdID: 0303
	// mrsigner Digest Alg: 1
	// mrsigner Digest Value: 87428fc522803d31065e7bce3cf03fe475096631e5e07bbd7a0fde60c4cf25c7
	// mrsigner Digest Alg: 8
	// mrsigner Digest Value: a314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6aa314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6a
	// CryptoKey Type: pkix-base64-key
	// CryptoKey Value: -----BEGIN PUBLIC KEY-----
	// MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg==
	// -----END PUBLIC KEY-----
}

func extractQERefVals(c *comid.Comid) error {
	if c.Triples.ReferenceValues == nil {
		return fmt.Errorf("no reference values triples")
	}

	for i, rv := range c.Triples.ReferenceValues.Values {
		if err := extractQERefVal(rv); err != nil {
			return fmt.Errorf("bad PSA reference value at index %d: %w", i, err)
		}
	}

	return nil
}

func extractQERefVal(rv comid.ValueTriple) error {
	class := rv.Environment.Class

	if err := extractClassElements(class); err != nil {
		return fmt.Errorf("extracting class: %w", err)
	}

	measurements := rv.Measurements
	if err := extractQEMeasurements(&measurements); err != nil {
		return fmt.Errorf("extracting measurements: %w", err)
	}

	return nil
}

func extractQEMeasurements(meas *comid.Measurements) error {
	if len(meas.Values) == 0 {
		return fmt.Errorf("no measurements")
	}
	for i := range meas.Values {
		m := &meas.Values[0]
		if err := decodeQEMValExtensions(m); err != nil {
			return fmt.Errorf("extracting measurement at index %d: %w", i, err)
		}

		if m.AuthorizedBy != nil {
			err := decodeAuthorisedBy(m)
			if err != nil {
				return fmt.Errorf("extracting measurement at index %d: %w", i, err)
			}
		}
	}
	return nil
}

func decodeQEMValExtensions(m *comid.Measurement) error {
	val, err := m.Val.Extensions.Get("miscselect")
	if err != nil {
		return fmt.Errorf("failed to decode miscselect from measurement extensions")
	}
	f, ok := val.(*TeeMiscSelect)
	if !ok {
		fmt.Printf("val was not pointer to TeeMiscSelect")
	}
	miscselect := *f
	fmt.Printf("\nmiscselect: %x", miscselect)

	val, err = m.Val.Extensions.Get("tcbevalnum")
	if err != nil {
		return fmt.Errorf("failed to decode tcbevalnum from measurement extensions")
	}
	t, ok := val.(*TeeTcbEvalNum)
	if !ok {
		fmt.Printf("val was not pointer to TeeTcbEvalNum")
	}
	tcbValNum := *t
	fmt.Printf("\ntcbEvalNum: %d", tcbValNum)

	val, err = m.Val.Extensions.Get("isvprodid")
	if err != nil {
		return fmt.Errorf("failed to decode isvprodid from measurement extensions")
	}
	tS, ok := val.(*TeeISVProdID)
	if !ok {
		fmt.Printf("val was not pointer to IsvProdID")
	}

	if tS.IsBytes() {
		val, err = tS.GetBytes()
		if err != nil {
			return fmt.Errorf("failed to decode isvprodid: %w", err)
		}
		fmt.Printf("\nIsvProdID: %x", val)
	} else if tS.IsUint() {
		val, err = tS.GetUint()
		if err != nil {
			return fmt.Errorf("failed to decode isvprodid: %w", err)
		}
		fmt.Printf("\nIsvProdID: %d", val)
	} else {
		return fmt.Errorf("isvprodid is neither integer or byte string")
	}

	val, err = m.Val.Extensions.Get("mrsigner")
	if err != nil {
		return fmt.Errorf("failed to decode mrsigner from measurement extensions")
	}

	tD, ok := val.(*TeeDigest)
	if !ok {
		fmt.Printf("val was not pointer to TeeDigest")
	}

	err = extractDigest("mrsigner", tD)
	if err != nil {
		return fmt.Errorf("unable to extract TEE Digest: %w", err)
	}
	return nil
}

func Example_encode_tdx_QE_refval_without_profile() {
	refVal := &comid.ValueTriple{}
	measurement := &comid.Measurement{}
	refVal.Environment = comid.Environment{
		Class: comid.NewClassOID(TestOID).
			SetVendor("Intel Corporation").
			SetModel("0123456789ABCDEF"), // From irim-qe-cend.diag, CPUID[0x01].EAX.FMSP & 0x0FFF0FF0
	}

	extMap := extensions.NewMap().
		Add(comid.ExtReferenceValue, &MValExtensions{})

	m := comid.NewComid().
		SetTagIdentity("43BBE37F-2E61-4B33-AED3-53CFF1428B20", 0).
		AddEntity("INTEL", &TestRegID, comid.RoleCreator, comid.RoleTagCreator, comid.RoleMaintainer)

	refVal.Measurements.Add(measurement)
	m.Triples.AddReferenceValue(*refVal)
	if err := m.RegisterExtensions(extMap); err != nil {
		panic(err)
	}

	if err := setTDXQEMvalExtensions(&m.Triples.ReferenceValues.Values[0].Measurements.Values[0].Val); err != nil {
		panic(err)
	}
	if err := m.Valid(); err != nil {
		panic(err)
	}

	cbor, err := m.ToCBOR()
	if err == nil {
		fmt.Printf("%x\n", cbor)
	} else {
		fmt.Printf("To CBOR failed \n")
	}

	json, err := m.ToJSON()
	if err == nil {
		fmt.Printf("%s\n", string(json))
	} else {
		fmt.Printf("To JSON failed \n")
	}

	// Output:
	// a301a1005043bbe37f2e614b33aed353cff1428b200281a30065494e54454c01d8207168747470733a2f2f696e74656c2e636f6d028301000204a1008182a100a300d86f4c6086480186f84d01020304050171496e74656c20436f72706f726174696f6e02703031323334353637383941424344454681a101a538480a385046c000fbff000038538282015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7582075830e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f3638540138550b
	// {"tag-identity":{"id":"43bbe37f-2e61-4b33-aed3-53cff1428b20"},"entities":[{"name":"INTEL","regid":"https://intel.com","roles":["creator","tagCreator","maintainer"]}],"triples":{"reference-values":[{"environment":{"class":{"id":{"type":"oid","value":"2.16.840.1.113741.1.2.3.4.5"},"vendor":"Intel Corporation","model":"0123456789ABCDEF"}},"measurements":[{"value":{"isvsvn":10,"miscselect":"wAD7/wAA","mrsigner":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU=","sha-384;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXXkW3L1wMC1cttNjTq36X82"],"isvprodid":{"type":"uint","value":1},"tcbevalnum":11}}]}]}}
}

func setTDXQEMvalExtensions(val *comid.Mval) error {
	svn := TeeSVN(10)
	teeTcbEvNum := TeeTcbEvalNum(11)
	teeMiscSel := TeeMiscSelect([]byte{0xC0, 0x00, 0xFB, 0xFF, 0x00, 0x00}) // Taken from irim-qe-ref.diag
	// Taken below from irim-qe-ref.diag
	r := 1
	isvProdID, err := NewTeeISVProdID(r)
	if err != nil {
		return fmt.Errorf("unable to get isvprodid %w", err)
	}

	err = val.Extensions.Set("isvprodid", isvProdID)
	if err != nil {
		return fmt.Errorf("unable to set isvprodid %w", err)
	}
	err = val.Extensions.Set("isvsvn", &svn)
	if err != nil {
		return fmt.Errorf("unable to set isvsvn %w", err)
	}
	err = val.Extensions.Set("tcbevalnum", &teeTcbEvNum)
	if err != nil {
		return fmt.Errorf("unable to set tcbevalnum %w", err)
	}
	err = val.Extensions.Set("miscselect", &teeMiscSel)
	if err != nil {
		return fmt.Errorf("unable to set miscselect %w", err)
	}

	d := comid.NewDigests()
	d.AddDigest(swid.Sha256, comid.MustHexDecode(nil, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"))
	d.AddDigest(swid.Sha384, comid.MustHexDecode(nil, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f36"))

	err = val.Extensions.Set("mrsigner", d)
	if err != nil {
		return fmt.Errorf("unable to set mrsigner %w", err)
	}
	return nil
}

var (
	// test cases are based on diag files here:
	// https://github.com/ietf-rats-wg/draft-ietf-rats-corim/tree/main/cddl/examples

	//go:embed testcases/comid_qe_refval.cbor
	testComid2 []byte
)

func Example_decode_QE_CBOR() {
	profileID, err := eat.NewProfile("2.16.840.1.113741.1.16.1")
	if err != nil {
		panic(err) // will not error, as the hard-coded string above is valid
	}
	manifest, found := corim.GetProfileManifest(profileID)
	if !found {
		fmt.Printf("CoRIM Profile NOT FOUND")
		return
	}

	m := manifest.GetComid()

	if err := m.FromCBOR(testComid2); err != nil {
		panic(err)
	}
	if err := m.Valid(); err != nil {
		panic(err)
	}

	if err := extractQERefVals(m); err != nil {
		panic(err)
	}

	// output:
	// OID: 2.16.840.1.113741.1.2.3.4.1
	// Vendor: Intel Corporation
	// Model: SGX QE TCB
	// miscselect: a0b0c0d000000000
	// tcbEvalNum: 11
	// IsvProdID: 1
	// mrsigner Digest Alg: 1
	// mrsigner Digest Value: a314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6a
	// mrsigner Digest Alg: 8
	// mrsigner Digest Value: a314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6aa314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6a
	// CryptoKey Type: pkix-base64-key
	// CryptoKey Value: -----BEGIN PUBLIC KEY-----
	// MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg==
	// -----END PUBLIC KEY-----
}
