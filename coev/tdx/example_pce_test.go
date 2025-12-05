// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	_ "embed"
	"errors"
	"fmt"
	"log"

	"github.com/veraison/corim/coev"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/corim/profiles/tdx"
	"github.com/veraison/eat"
)

func Example_decode_PCE_Evidence_JSON() {
	manifest, found := coev.GetProfileManifest(ProfileID)
	if !found {
		fmt.Printf("Evidence profile not found")
		return
	}
	ce := manifest.GetTaggedConciseEvidence()
	if err := ce.FromJSON([]byte(TDXPCECETemplate)); err != nil {
		panic(err)
	}
	if err := ce.Valid(); err != nil {
		panic(err)
	}

	if err := ExtractPceEvidence(ce); err != nil {
		panic(err)
	}

	// OUTPUT:
	// OID: 2.16.840.1.113741.1.2.3.4.4
	// Vendor: Intel Corporation
	// Model: TDX PCE TCB
	// InstanceID: 45
	// pceID: PCEID001
	// ISVSVN: 1
	// ISVSVN: 2
	// ISVSVN: 3
	// ISVSVN: 4
	// ISVSVN: 5
	// ISVSVN: 6
	// ISVSVN: 7
	// ISVSVN: 8
	// ISVSVN: 9
	// ISVSVN: 10
	// ISVSVN: 11
	// ISVSVN: 12
	// ISVSVN: 13
	// ISVSVN: 14
	// ISVSVN: 15
	// ISVSVN: 16
}

func Example_encode_tdx_pce_evidence_without_profile() {
	coEv := coev.NewConciseEvidence()
	valTriple := &comid.ValueTriple{}
	measurement := &comid.Measurement{}
	valTriple.Environment = comid.Environment{
		Class: comid.NewClassOID(TestPCEOID).
			SetVendor("Intel Corporation").
			SetModel("TDX PCE TCB"),
	}

	extMap := extensions.NewMap().
		Add(comid.ExtMval, &tdx.MValExtensions{})

	if err := measurement.Val.RegisterExtensions(extMap); err != nil {
		log.Fatal("could not register mval extensions")
	}
	if err := tdx.SetTdxPceMvalExtensions(tdx.Evidence, &measurement.Val); err != nil {
		panic(err)
	}
	valTriple.Measurements.Add(measurement)

	ev := coev.NewEvTriples()
	ev.AddEvidenceTriple(valTriple)
	err := coEv.AddTriples(ev)
	if err != nil {
		panic(err)
	}
	if err = coEv.Valid(); err != nil {
		panic(err)
	}
	te, err := coev.NewTaggedConciseEvidence(coEv)
	if err != nil {
		panic(err)
	}

	cbor, err := te.ToCBOR()
	if err == nil {
		fmt.Printf("%x\n", cbor)
	} else {
		fmt.Printf("To CBOR failed \n")
	}

	json, err := te.ToJSON()
	if err == nil {
		fmt.Printf("%s\n", string(json))
	} else {
		fmt.Printf("unable to format json %s, %s", err.Error(), json)
	}

	// output:
	// d9023ba100a1008182a100a300d86f4c6086480186f84d01020304040171496e74656c20436f72706f726174696f6e026b544458205043452054434281a101a3384c182d384f685043454944303031387c900102030405060708090a0b0c0d0e0f10
	// {"ev-triples":{"evidence-triples":[{"environment":{"class":{"id":{"type":"oid","value":"2.16.840.1.113741.1.2.3.4.4"},"vendor":"Intel Corporation","model":"TDX PCE TCB"}},"measurements":[{"value":{"instanceid":{"type":"uint","value":45},"pceid":"PCEID001","tcbcompsvn":[{"type":"uint","value":1},{"type":"uint","value":2},{"type":"uint","value":3},{"type":"uint","value":4},{"type":"uint","value":5},{"type":"uint","value":6},{"type":"uint","value":7},{"type":"uint","value":8},{"type":"uint","value":9},{"type":"uint","value":10},{"type":"uint","value":11},{"type":"uint","value":12},{"type":"uint","value":13},{"type":"uint","value":14},{"type":"uint","value":15},{"type":"uint","value":16}]}}]}]}}
}

func Example_encode_tdx_pce_evidence_with_profile() {
	profileID, err := eat.NewProfile("2.16.840.1.113741.1.16.1")
	if err != nil {
		panic(err) // will not error, as the hard-coded string above is valid
	}

	manifest, found := coev.GetProfileManifest(profileID)
	if !found {
		fmt.Printf("CoEV Profile NOT FOUND")
		return
	}
	coEv := manifest.GetConciseEvidence()

	valTriple := &comid.ValueTriple{}
	valTriple.Environment = comid.Environment{
		Class: comid.NewClassOID(TestPCEOID).
			SetVendor("Intel Corporation").
			SetModel("TDX PCE TCB"),
	}

	measurement := &comid.Measurement{}
	valTriple.Measurements.Add(measurement)
	coEv.EvTriples.EvidenceTriples.Add(valTriple)

	if err = tdx.SetTdxPceMvalExtensions(tdx.Evidence, &coEv.EvTriples.EvidenceTriples.Values[0].Measurements.Values[0].Val); err != nil {
		panic(err)
	}

	te, err := coev.NewTaggedConciseEvidence(coEv)
	if err != nil {
		panic(err)
	}

	cbor, err := te.ToCBOR()
	if err == nil {
		fmt.Printf("%x\n", cbor)
	} else {
		fmt.Printf("To CBOR failed \n")
	}

	json, err := te.ToJSON()
	if err == nil {
		fmt.Printf("%s\n", string(json))
	} else {
		fmt.Printf("unable to format json %s, %s", err.Error(), json)
	}

	// output:
	// d9023ba100a1008182a100a300d86f4c6086480186f84d01020304040171496e74656c20436f72706f726174696f6e026b544458205043452054434281a101a3384c182d384f685043454944303031387c900102030405060708090a0b0c0d0e0f10
	// {"ev-triples":{"evidence-triples":[{"environment":{"class":{"id":{"type":"oid","value":"2.16.840.1.113741.1.2.3.4.4"},"vendor":"Intel Corporation","model":"TDX PCE TCB"}},"measurements":[{"value":{"instanceid":{"type":"uint","value":45},"pceid":"PCEID001","tcbcompsvn":[{"type":"uint","value":1},{"type":"uint","value":2},{"type":"uint","value":3},{"type":"uint","value":4},{"type":"uint","value":5},{"type":"uint","value":6},{"type":"uint","value":7},{"type":"uint","value":8},{"type":"uint","value":9},{"type":"uint","value":10},{"type":"uint","value":11},{"type":"uint","value":12},{"type":"uint","value":13},{"type":"uint","value":14},{"type":"uint","value":15},{"type":"uint","value":16}]}}]}]}}
}

var (
	//go:embed testcases/ce-pce-evidence.cbor
	testCePceEvidence []byte
)

func Example_decode_PCE_Evidence_CBOR() {
	profileID, err := eat.NewProfile("2.16.840.1.113741.1.16.1")
	if err != nil {
		panic(err) // will not error, as the hard-coded string above is valid
	}
	manifest, found := coev.GetProfileManifest(profileID)
	if !found {
		fmt.Printf("Evidence Profile NOT FOUND")
		return
	}
	ce := manifest.GetTaggedConciseEvidence()

	if err := ce.FromCBOR(testCePceEvidence); err != nil {
		panic(err)
	}
	if err := ce.Valid(); err != nil {
		panic(err)
	}
	if err := ExtractPceEvidence(ce); err != nil {
		panic(err)
	}

	// OUTPUT:
	// OID: 2.16.840.1.113741.1.2.3.4.5
	// Vendor: Intel Corporation
	// Model: TDX PCE TCB
	// InstanceID: 00112233445566778899aabbccddeeff
	// pceID: 0000
	// ISVSVN: 10
	// ISVSVN: 10
	// ISVSVN: 2
	// ISVSVN: 2
	// ISVSVN: 2
	// ISVSVN: 1
	// ISVSVN: 4
	// ISVSVN: 0
	// ISVSVN: 0
	// ISVSVN: 0
	// ISVSVN: 0
	// ISVSVN: 0
	// ISVSVN: 0
	// ISVSVN: 0
	// ISVSVN: 0
	// ISVSVN: 0
	// CryptoKeys: [-----BEGIN PUBLIC KEY-----
	// MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg==
	// -----END PUBLIC KEY-----]
}

func ExtractPceEvidence(ce *coev.TaggedConciseEvidence) error {
	if ce.EvTriples.EvidenceTriples == nil {
		return errors.New("no evidence triples")
	}

	for i, ev := range ce.EvTriples.EvidenceTriples.Values {
		if err := tdx.ExtractPceMeas(ev); err != nil {
			return fmt.Errorf("bad evidence at index %d: %w", i, err)
		}
	}
	return nil
}
