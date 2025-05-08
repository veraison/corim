// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/corim"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
	"github.com/veraison/swid"
)

// Example_decode_JSON decodes the TDX Measurement Extensions from the given JSON Template
func Example_decode_JSON() {
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
	if err := m.FromJSON([]byte(TDXSeamRefValJSONTemplate)); err != nil {
		panic(err)
	}

	if err := m.Valid(); err != nil {
		panic(err)
	}

	if err := extractRefVals(m); err != nil {
		panic(err)
	}

	// output:
	// OID: 2.16.840.1.113741.1.2.3.4.5
	// Vendor: Intel Corporation
	// Model: TDX SEAM
	// tcbEvalNum Operator: greater_or_equal
	// tcbEvalNum Value: 11
	// IsvProdID: 0303
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// Attributes: f00a0b
	// mrtee Digest Operator: member
	// mrtee Digest Alg: 1
	// mrtee Digest Value: 87428fc522803d31065e7bce3cf03fe475096631e5e07bbd7a0fde60c4cf25c7
	// mrsigner Digest Operator: member
	// mrsigner Digest Alg: 1
	// mrsigner Digest Value: 87428fc522803d31065e7bce3cf03fe475096631e5e07bbd7a0fde60c4cf25c7
	// mrsigner Digest Alg: 8
	// mrsigner Digest Value: a314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6aa314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6a
	// CryptoKey Type: pkix-base64-key
	// CryptoKey Value: -----BEGIN PUBLIC KEY-----
	// MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg==
	// -----END PUBLIC KEY-----
}

func Example_encode_tdx_seam_refval_without_profile() {
	refVal := &comid.ValueTriple{}
	measurement := &comid.Measurement{}
	refVal.Environment = comid.Environment{
		Class: comid.NewClassOID(TestOID).
			SetVendor("Intel Corporation").
			SetModel("TDXSEAM"),
	}

	extMap := extensions.NewMap().
		Add(comid.ExtReferenceValue, &MValExtensions{})

	m := comid.NewComid().
		SetTagIdentity("43BBE37F-2E61-4B33-AED3-53CFF1428B20", 0).
		AddEntity("INTEL", &TestRegID, comid.RoleCreator, comid.RoleTagCreator, comid.RoleMaintainer)

	refVal.Measurements.Add(measurement)
	m.Triples.AddReferenceValue(refVal)
	if err := m.RegisterExtensions(extMap); err != nil {
		panic(err)
	}

	if err := setTDXSeamMvalExtensions(&m.Triples.ReferenceValues.Values[0].Measurements.Values[0].Val); err != nil {
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
		fmt.Printf("unable to format json %s, %s", err.Error(), json)
	}
	// output:
	// a301a1005043bbe37f2e614b33aed353cff1428b200281a30065494e54454c01d8207168747470733a2f2f696e74656c2e636f6d028301000204a1008182a100a300d86f4c6086480186f84d01020304050171496e74656c20436f72706f726174696f6e02675444585345414d81a101a73847c11a6796cc803848d9ea6a82020a38514201013852d9ea7482068182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d753853d9ea7482068282015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7582075830e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f3638544201013855d9ea6a82020b
	// {"tag-identity":{"id":"43bbe37f-2e61-4b33-aed3-53cff1428b20"},"entities":[{"name":"INTEL","regid":"https://intel.com","roles":["creator","tagCreator","maintainer"]}],"triples":{"reference-values":[{"environment":{"class":{"id":{"type":"oid","value":"2.16.840.1.113741.1.2.3.4.5"},"vendor":"Intel Corporation","model":"TDXSEAM"}},"measurements":[{"value":{"tcbdate":"2025-01-27T00:00:00Z","isvsvn":{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":10}}},"attributes":"AQE=","mrtee":{"type":"digest-expression","value":{"set-operator":"member","set-digest":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]}},"mrsigner":{"type":"digest-expression","value":{"set-operator":"member","set-digest":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU=","sha-384;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXXkW3L1wMC1cttNjTq36X82"]}},"isvprodid":{"type":"bytes","value":"AQE="},"tcbevalnum":{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":11}}}}}]}]}}
}

func Example_encode_tdx_seam_refval_with_profile() {
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
	if m == nil {
		panic(err)
	}
	m.SetTagIdentity("43BBE37F-2E61-4B33-AED3-53CFF1428B20", 0).
		AddEntity("INTEL", &TestRegID, comid.RoleCreator, comid.RoleTagCreator, comid.RoleMaintainer)

	refVal := &comid.ValueTriple{}
	measurement := &comid.Measurement{}
	refVal.Environment = comid.Environment{
		Class: comid.NewClassOID(TestOID).
			SetVendor("Intel Corporation").
			SetModel("TDXSEAM"),
	}

	refVal.Measurements.Add(measurement)
	m.Triples.AddReferenceValue(refVal)

	err = setTDXSeamMvalExtensions(&m.Triples.ReferenceValues.Values[0].Measurements.Values[0].Val)
	if err != nil {
		fmt.Printf("unable to set extensions :%s", err.Error())
	}

	err = m.Valid()
	if err != nil {
		fmt.Printf("CoMID is not Valid :%s", err.Error())
	}

	cbor, err := m.ToCBOR()
	if err == nil {
		fmt.Printf("%x\n", cbor)
	} else {
		fmt.Printf("\n To CBOR Failed: %s \n", err.Error())
	}

	json, err := m.ToJSON()
	if err == nil {
		fmt.Printf("%s\n", string(json))
	} else {
		fmt.Printf("\n To JSON Failed \n")
	}

	// output:
	// a301a1005043bbe37f2e614b33aed353cff1428b200281a30065494e54454c01d8207168747470733a2f2f696e74656c2e636f6d028301000204a1008182a100a300d86f4c6086480186f84d01020304050171496e74656c20436f72706f726174696f6e02675444585345414d81a101a73847c11a6796cc803848d9ea6a82020a38514201013852d9ea7482068182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d753853d9ea7482068282015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7582075830e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f3638544201013855d9ea6a82020b
	// {"tag-identity":{"id":"43bbe37f-2e61-4b33-aed3-53cff1428b20"},"entities":[{"name":"INTEL","regid":"https://intel.com","roles":["creator","tagCreator","maintainer"]}],"triples":{"reference-values":[{"environment":{"class":{"id":{"type":"oid","value":"2.16.840.1.113741.1.2.3.4.5"},"vendor":"Intel Corporation","model":"TDXSEAM"}},"measurements":[{"value":{"tcbdate":"2025-01-27T00:00:00Z","isvsvn":{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":10}}},"attributes":"AQE=","mrtee":{"type":"digest-expression","value":{"set-operator":"member","set-digest":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]}},"mrsigner":{"type":"digest-expression","value":{"set-operator":"member","set-digest":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU=","sha-384;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXXkW3L1wMC1cttNjTq36X82"]}},"isvprodid":{"type":"bytes","value":"AQE="},"tcbevalnum":{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":11}}}}}]}]}}
}

func Example_encode_tdx_seam_refval_direct() {
	refVal := &comid.ValueTriple{}
	measurement := &comid.Measurement{}
	refVal.Environment = comid.Environment{
		Class: comid.NewClassOID(TestOID).
			SetVendor("Intel Corporation").
			SetModel("TDXSEAM"),
	}

	extMap := extensions.NewMap().Add(comid.ExtMval, &MValExtensions{})
	m := comid.NewComid().
		SetTagIdentity("43BBE37F-2E61-4B33-AED3-53CFF1428B20", 0).
		AddEntity("INTEL", &TestRegID, comid.RoleCreator, comid.RoleTagCreator, comid.RoleMaintainer)

	if err := measurement.Val.RegisterExtensions(extMap); err != nil {
		log.Fatal("could not register mval extensions")
	}

	if err := setTDXSeamMvalExtensions(&measurement.Val); err != nil {
		log.Fatal("could not set mval extensions")
	}

	refVal.Measurements.Add(measurement)
	m.Triples.AddReferenceValue(refVal)

	err := m.Valid()
	if err != nil {
		fmt.Printf("CoMID is not Valid :%s", err.Error())
	}

	cbor, err := m.ToCBOR()
	if err == nil {
		fmt.Printf("%x\n", cbor)
	} else {
		fmt.Printf("\n To CBOR Failed: %s \n", err.Error())
	}

	json, err := m.ToJSON()
	if err == nil {
		fmt.Printf("%s\n", string(json))
	} else {
		fmt.Printf("\n To JSON Failed \n")
	}

	// output:
	// a301a1005043bbe37f2e614b33aed353cff1428b200281a30065494e54454c01d8207168747470733a2f2f696e74656c2e636f6d028301000204a1008182a100a300d86f4c6086480186f84d01020304050171496e74656c20436f72706f726174696f6e02675444585345414d81a101a73847c11a6796cc803848d9ea6a82020a38514201013852d9ea7482068182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d753853d9ea7482068282015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7582075830e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f3638544201013855d9ea6a82020b
	// {"tag-identity":{"id":"43bbe37f-2e61-4b33-aed3-53cff1428b20"},"entities":[{"name":"INTEL","regid":"https://intel.com","roles":["creator","tagCreator","maintainer"]}],"triples":{"reference-values":[{"environment":{"class":{"id":{"type":"oid","value":"2.16.840.1.113741.1.2.3.4.5"},"vendor":"Intel Corporation","model":"TDXSEAM"}},"measurements":[{"value":{"tcbdate":"2025-01-27T00:00:00Z","isvsvn":{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":10}}},"attributes":"AQE=","mrtee":{"type":"digest-expression","value":{"set-operator":"member","set-digest":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]}},"mrsigner":{"type":"digest-expression","value":{"set-operator":"member","set-digest":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU=","sha-384;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXXkW3L1wMC1cttNjTq36X82"]}},"isvprodid":{"type":"bytes","value":"AQE="},"tcbevalnum":{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":11}}}}}]}]}}
}

func setTDXSeamMvalExtensions(val *comid.Mval) error {
	tcbDate, _ := time.Parse(time.RFC3339, "2025-01-27T00:00:00Z")
	err := val.Extensions.Set("tcbdate", &tcbDate)
	if err != nil {
		return fmt.Errorf("unable to set tcbDate: %w", err)
	}
	r := []byte{0x01, 0x01}
	isvProdID, err := NewTeeISVProdID(r)
	if err != nil {
		return fmt.Errorf("unable to get isvprodid: %w", err)
	}
	err = val.Extensions.Set("isvprodid", isvProdID)
	if err != nil {
		return fmt.Errorf("unable to set isvprodid: %w", err)
	}
	svn, err := NewSvnExpression(TestISVSVN)
	if err != nil {
		return fmt.Errorf("unable to get isvsvn numeric: %w", err)
	}
	err = val.Extensions.Set("isvsvn", svn)
	if err != nil {
		return fmt.Errorf("unable to set isvsvn: %w", err)
	}
	teeTcbEvNum, err := NewTeeTcbEvalNumberNumeric(TestTCBEvalNum)
	if err != nil {
		return fmt.Errorf("unable to get tcbevalnum numeric: %w", err)
	}
	err = val.Extensions.Set("tcbevalnum", teeTcbEvNum)
	if err != nil {
		return fmt.Errorf("unable to set tcbevalnum: %w", err)
	}

	teeAttr, err := NewTeeAttributes(TestTeeAttributes)
	if err != nil {
		return fmt.Errorf("unable to get teeAttributes: %w", err)
	}
	err = val.Extensions.Set("attributes", teeAttr)
	if err != nil {
		return fmt.Errorf("unable to set attributes: %w", err)
	}

	d := comid.NewDigests()
	d.AddDigest(swid.Sha256, comid.MustHexDecode(nil, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"))
	td, err := NewTeeDigestExpr(MEM, *d)
	if err != nil {
		return fmt.Errorf("unable to get TeeDigest: %w", err)
	}

	err = val.Extensions.Set("mrtee", td)
	if err != nil {
		return fmt.Errorf("unable to set mrtee: %w", err)
	}

	d = comid.NewDigests()
	d.AddDigest(swid.Sha256, comid.MustHexDecode(nil, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"))
	d.AddDigest(swid.Sha384, comid.MustHexDecode(nil, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f36"))

	td, err = NewTeeDigestExpr(MEM, *d)
	if err != nil {
		return fmt.Errorf("unable to get TeeDigest: %w", err)
	}

	err = val.Extensions.Set("mrsigner", td)
	if err != nil {
		return fmt.Errorf("unable to set mrsigner %w", err)
	}
	return nil
}

func decodeMValExtensions(m *comid.Measurement) error {
	val, err := m.Val.Extensions.Get("tcbevalnum")
	if err != nil {
		return fmt.Errorf("failed to decode tcbevalnum from measurement extensions: %w", err)
	}
	f, ok := val.(*TeeTcbEvalNumber)
	if !ok {
		fmt.Printf("val was not pointer to TeeTcbEvalNum")
	}
	tcbValNum := *f
	if err = TestextractTeeTcbEvalNum(&tcbValNum); err != nil {
		return fmt.Errorf("failed to extract tcbevalnum: %w", err)
	}

	val, err = m.Val.Extensions.Get("isvprodid")
	if err != nil {
		return errors.New("failed to decode isvprodid from measurement extensions")
	}
	tS, ok := val.(*TeeISVProdID)
	if !ok {
		fmt.Printf("val was not pointer to IsvProdID")
	}
	if err = TestextractTeeISVProdID(tS); err != nil {
		return fmt.Errorf("failed to decode teeISVProdID from measurement extensions: %w", err)
	}

	val, err = m.Val.Extensions.Get("isvsvn")
	if err != nil {
		return errors.New("failed to decode isvsvn from measurement extensions")
	}
	teesvn, ok := val.(*TeeSVN)
	if !ok {
		return errors.New("val was not pointer to tee svn")
	}
	err = teesvn.Valid()
	if err != nil {
		return fmt.Errorf("invalid tee svn: %w", err)
	}

	err = TestextractTeeSvn(teesvn)
	if err != nil {
		return fmt.Errorf("unable to extract tee svn: %w", err)
	}
	val, err = m.Val.Extensions.Get("attributes")
	if err != nil {
		return errors.New("failed to decode attributes from measurement extensions")
	}

	tA, ok := val.(*TeeAttributes)
	if !ok {
		fmt.Printf("val was not pointer to teeAttributes")
	}
	fmt.Printf("\nAttributes: %x", *tA)

	val, err = m.Val.Extensions.Get("mrtee")
	if err != nil {
		return errors.New("failed to decode mrtee from measurement extensions")
	}

	tD, ok := val.(*TeeDigest)
	if !ok {
		fmt.Printf("val was not pointer to TeeDigest")
	}

	if err = TestextractTeeDigest("mrtee", tD); err != nil {
		return fmt.Errorf("failed to decode mrtee from digest: %w", err)
	}

	val, err = m.Val.Extensions.Get("mrsigner")
	if err != nil {
		return fmt.Errorf("failed to decode mrsigner from measurement extensions: %w", err)
	}

	tD, ok = val.(*TeeDigest)
	if !ok {
		return errors.New("val was not pointer to TeeDigest")
	}

	if err := TestextractTeeDigest("mrsigner", tD); err != nil {
		return fmt.Errorf("failed to extarct mrsigner digest: %w", err)
	}
	return nil
}

var (
	// test cases are based on diag files here:
	// https://github.com/ietf-rats-wg/draft-ietf-rats-corim/tree/main/cddl/examples

	//go:embed testcases/comid_seam_refval.cbor
	testComid1 []byte
)

func Example_decode_CBOR() {
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

	if err := m.FromCBOR(testComid1); err != nil {
		panic(err)
	}
	if err := m.Valid(); err != nil {
		panic(err)
	}

	if err := extractRefVals(m); err != nil {
		panic(err)
	}

	// output:
	// OID: 2.16.840.1.113741.1.2.3.4.3
	// Vendor: Intel Corporation
	// Model: TDX SEAM
	// tcbEvalNum Operator: greater_or_equal
	// tcbEvalNum Value: 11
	// IsvProdID: abcd
	// SVN Operator: greater_or_equal
	// SVN Value: 6
	// Attributes: 0102
	// mrtee Digest Operator: member
	// mrtee Digest Alg: 1
	// mrtee Digest Value: a314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6a
	// mrsigner Digest Operator: member
	// mrsigner Digest Alg: 1
	// mrsigner Digest Value: a314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6a
	// mrsigner Digest Alg: 8
	// mrsigner Digest Value: a314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6aa314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6a
	// CryptoKey Type: pkix-base64-key
	// CryptoKey Value: -----BEGIN PUBLIC KEY-----
	// MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg==
	// -----END PUBLIC KEY-----

}

func extractRefVals(c *comid.Comid) error {
	if c.Triples.ReferenceValues == nil {
		return errors.New("no reference values triples")
	}

	for i, rv := range c.Triples.ReferenceValues.Values {
		if err := extractSeamRefVal(rv); err != nil {
			return fmt.Errorf("bad reference value at index %d: %w", i, err)
		}
	}
	return nil
}

func extractSeamRefVal(rv comid.ValueTriple) error {
	class := rv.Environment.Class

	if err := TestextractClassElements(class); err != nil {
		return fmt.Errorf("extracting class: %w", err)
	}

	measurements := rv.Measurements
	if err := extractSeamMeasurements(&measurements); err != nil {
		return fmt.Errorf("extracting measurements: %w", err)
	}

	return nil
}

func extractSeamMeasurements(meas *comid.Measurements) error {
	if len(meas.Values) == 0 {
		return errors.New("no measurements")
	}
	for i := range meas.Values {
		m := &meas.Values[0]
		if err := decodeMValExtensions(m); err != nil {
			return fmt.Errorf("extracting measurement at index %d: %w", i, err)
		}

		if m.AuthorizedBy != nil {
			err := TestdecodeAuthorisedBy(m)
			if err != nil {
				return fmt.Errorf("extracting measurement at index %d: %w", i, err)
			}
		}
	}
	return nil
}
