// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	_ "embed"
	"errors"
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/corim"
	"github.com/veraison/eat"
)

// Example_decode_PCE_JSON decodes the TDX Provisioning Certification Enclave Measurement Extensions from the given JSON Template
func Example_decode_PCE_JSON() {
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
	if err := m.FromJSON([]byte(TDXPCERefValTemplate)); err != nil {
		panic(err)
	}

	if err := m.Valid(); err != nil {
		panic(err)
	}

	if err := extractPCERefVals(m); err != nil {
		panic(err)
	}

	// output:
	// OID: 2.16.840.1.113741.1.2.3.4.6
	// Vendor: Intel Corporation
	// Model: 0123456789ABCDEF
	// InstanceID: 11
	// pceID: 0000
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 2
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// CryptoKey Type: pkix-base64-key
	// CryptoKey Value: -----BEGIN PUBLIC KEY-----
	// MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg==
	// -----END PUBLIC KEY-----
}

func extractPCERefVals(c *comid.Comid) error {
	if c.Triples.ReferenceValues == nil {
		return errors.New("no reference values triples")
	}

	for i, rv := range c.Triples.ReferenceValues.Values {
		if err := extractPCERefVal(rv); err != nil {
			return fmt.Errorf("bad PSA reference value at index %d: %w", i, err)
		}
	}

	return nil
}

func extractPCERefVal(rv comid.ValueTriple) error {
	class := rv.Environment.Class

	if err := testextractClassElements(class); err != nil {
		return fmt.Errorf("extracting class: %w", err)
	}

	measurements := rv.Measurements
	if err := extractPCEMeasurements(&measurements); err != nil {
		return fmt.Errorf("extracting measurements: %w", err)
	}

	return nil
}

func extractPCEMeasurements(meas *comid.Measurements) error {
	if len(meas.Values) == 0 {
		return errors.New("no measurements")
	}
	for i := range meas.Values {
		m := &meas.Values[0]
		if err := decodePCEMValExtensions(m); err != nil {
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

func decodePCEMValExtensions(m *comid.Measurement) error {
	val, err := m.Val.Get("instanceid")
	if err != nil {
		return errors.New("failed to decode instanceid from measurement extensions")
	}
	i, ok := val.(*TeeInstanceID)
	if !ok {
		return errors.New("val was not pointer to teeInstanceID")
	}

	if i.IsBytes() {
		val, err = i.GetBytes()
		if err != nil {
			return fmt.Errorf("failed to decode teeinstanceid: %w", err)
		}
		fmt.Printf("\nInstanceID: %x", val)
	} else if i.IsUint() {
		val, err = i.GetUint()
		if err != nil {
			return fmt.Errorf("failed to decode teeinstanceid: %w", err)
		}
		fmt.Printf("\nInstanceID: %d", val)
	} else {
		return errors.New("teeinstanceid is neither integer or byte string")
	}

	val, err = m.Val.Get("tcbcompsvn")
	if err != nil {
		return errors.New("failed to decode teetcbcompsvn from measurement extensions")
	}

	tcs, ok := val.(*TeeTcbCompSvn)
	if !ok {
		return errors.New("val was not pointer to teetcbcompsvn")
	}
	if err = tcs.Valid(); err != nil {
		return fmt.Errorf("invalid computed SVN: %w", err)
	}

	val, err = m.Val.Get("pceid")
	if err != nil {
		return errors.New("failed to decode tcbevalnum from measurement extensions")
	}
	t, ok := val.(*TeePCEID)
	if !ok {
		fmt.Printf("val was not pointer to TeeTcbEvalNum")
	}
	if err = t.Valid(); err != nil {
		return fmt.Errorf("invalid PCEID: %w", err)
	}
	pceID := *t
	fmt.Printf("\npceID: %s", pceID)

	err = extractCompSvn(tcs)
	if err != nil {
		return fmt.Errorf("unable to extract TeeTcbCompSVN: %w", err)
	}
	return nil
}

func extractCompSvn(s *TeeTcbCompSvn) error {
	if s == nil {
		return errors.New("no TEE TCB Comp SVN")
	}

	if len(*s) > 16 {
		return errors.New("computed SVN cannot be greater than 16")
	}

	for i, teesvn := range *s {
		svn := teesvn // Avoid gosec: Implicit memory aliasing in for loop
		if err := testextractTeeSvn(&svn); err != nil {
			return fmt.Errorf("unable to extract SVN at index %d: %w", i, err)
		}
	}

	return nil
}

var (
	// test cases are based on diag files here:
	// https://github.com/ietf-rats-wg/draft-ietf-rats-corim/tree/main/cddl/examples

	//go:embed testcases/comid_pce_refval.cbor
	testComid3 []byte
)

func Example_decode_PCE_CBOR() {
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

	if err := m.FromCBOR(testComid3); err != nil {
		panic(err)
	}
	if err := m.Valid(); err != nil {
		panic(err)
	}

	if err := extractPCERefVals(m); err != nil {
		panic(err)
	}

	// output:
	// OID: 2.16.840.1.113741.1.2.3.4.5
	// Vendor: Intel Corporation
	// Model: TDX PCE TCB
	// InstanceID: 00112233445566778899aabbccddeeff
	// pceID: 0000
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 10
	// SVN Operator: greater_or_equal
	// SVN Value: 2
	// SVN Operator: greater_or_equal
	// SVN Value: 2
	// SVN Operator: greater_or_equal
	// SVN Value: 2
	// SVN Operator: greater_or_equal
	// SVN Value: 1
	// SVN Operator: greater_or_equal
	// SVN Value: 4
	// SVN Operator: greater_or_equal
	// SVN Value: 0
	// SVN Operator: greater_or_equal
	// SVN Value: 0
	// SVN Operator: greater_or_equal
	// SVN Value: 0
	// SVN Operator: greater_or_equal
	// SVN Value: 0
	// SVN Operator: greater_or_equal
	// SVN Value: 0
	// SVN Operator: greater_or_equal
	// SVN Value: 0
	// SVN Operator: greater_or_equal
	// SVN Value: 0
	// SVN Operator: greater_or_equal
	// SVN Value: 0
	// SVN Operator: greater_or_equal
	// SVN Value: 0
	// CryptoKey Type: pkix-base64-key
	// CryptoKey Value: -----BEGIN PUBLIC KEY-----
	// MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg==
	// -----END PUBLIC KEY-----
}

func Example_encode_tdx_pce_refval_with_profile() {
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
			SetModel("TDX PCE TCB"),
	}

	refVal.Measurements.Add(measurement)
	m.Triples.AddReferenceValue(refVal)

	err = setTDXPCEMvalExtension(&m.Triples.ReferenceValues.Values[0].Measurements.Values[0].Val)
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
	}

	// Output:
	// a301a1005043bbe37f2e614b33aed353cff1428b200281a30065494e54454c01d8207168747470733a2f2f696e74656c2e636f6d028301000204a1008182a100a300d86f4c6086480186f84d01020304050171496e74656c20436f72706f726174696f6e026b544458205043452054434281a101a3384c182d384f685043454944303031387c90d9ea6a820201d9ea6a820202d9ea6a820203d9ea6a820204d9ea6a820205d9ea6a820206d9ea6a820207d9ea6a820208d9ea6a820209d9ea6a82020ad9ea6a82020bd9ea6a82020cd9ea6a82020dd9ea6a82020ed9ea6a82020fd9ea6a820210
	// {"tag-identity":{"id":"43bbe37f-2e61-4b33-aed3-53cff1428b20"},"entities":[{"name":"INTEL","regid":"https://intel.com","roles":["creator","tagCreator","maintainer"]}],"triples":{"reference-values":[{"environment":{"class":{"id":{"type":"oid","value":"2.16.840.1.113741.1.2.3.4.5"},"vendor":"Intel Corporation","model":"TDX PCE TCB"}},"measurements":[{"value":{"instanceid":{"type":"uint","value":45},"pceid":"PCEID001","tcbcompsvn":[{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":1}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":2}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":3}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":4}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":5}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":6}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":7}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":8}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":9}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":10}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":11}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":12}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":13}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":14}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":15}}},{"type":"numeric-expression","value":{"numeric-operator":"greater_or_equal","numeric-type":{"type":"uint","value":16}}}]}}]}]}}
}

func setTDXPCEMvalExtension(val *comid.Mval) error {
	instanceID, err := NewTeeInstanceID(TestUIntInstance)
	if err != nil {
		return fmt.Errorf("unable to get teeinstanceID: %w", err)
	}
	err = val.Set("instanceid", instanceID)
	if err != nil {
		return fmt.Errorf("unable to set teeinstanceID: %w", err)
	}

	p, err := NewTeePCEID(TestPCEID)
	if err != nil {
		return fmt.Errorf("unable to get NewTeepceID: %w", err)
	}
	err = val.Set("pceid", p)
	if err != nil {
		return fmt.Errorf("unable to set teepceID: %w", err)
	}

	c, err := NewTeeTcbCompSvnExpression(TestCompSvn)
	if err != nil {
		return fmt.Errorf("failed to get TeeTcbCompSvn: %w", err)
	}

	err = val.Set("tcbcompsvn", c)
	if err != nil {
		return fmt.Errorf("unable to set teetcbcompsvn: %w", err)
	}

	return nil
}
