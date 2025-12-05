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

func Example_decode_QE_Evidence_JSON() {
	manifest, found := coev.GetProfileManifest(ProfileID)
	if !found {
		fmt.Printf("Evidence Profile NOT FOUND")
		return
	}
	ce := manifest.GetTaggedConciseEvidence()
	if err := ce.FromJSON([]byte(TDXQECETemplate)); err != nil {
		panic(err)
	}
	if err := ce.Valid(); err != nil {
		panic(err)
	}

	if err := ExtractQeEvidence(ce); err != nil {
		panic(err)
	}

	// OUTPUT:
	// OID: 2.16.840.1.113741.1.15.4.99.1
	// Vendor: Intel Corporation
	// Model: TDX QE TCB
	// tcbEvalNum: 11
	// IsvProdID: 1
	// mrsigner Digest Alg: 1
	// mrsigner Digest Value: e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75
	// mrsigner Digest Alg: 8
	// mrsigner Digest Value: a314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6aa314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6a
	// miscselect: a0b0c0d000000000
	// TEE TCB Status = UpToDate
	// Tee AdvisoryID = INTEL-SA-00078
	// Tee AdvisoryID = INTEL-SA-00079
	// CryptoKeys: [-----BEGIN PUBLIC KEY-----
	// MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg==
	// -----END PUBLIC KEY-----]
}

func Example_encode_tdx_qe_evidence_without_profile() {
	coEv := coev.NewConciseEvidence()
	valTriple := &comid.ValueTriple{}
	measurement := &comid.Measurement{}
	valTriple.Environment = comid.Environment{
		Class: comid.NewClassOID(TestQEOID).
			SetVendor("Intel Corporation").
			SetModel("0123456789ABCDEF"), // From irim-qe-cend.diag, CPUID[0x01].EAX.FMSP & 0x0FFF0FF0
	}

	extMap := extensions.NewMap().
		Add(comid.ExtMval, &tdx.MValExtensions{})

	if err := measurement.Val.RegisterExtensions(extMap); err != nil {
		log.Fatal("could not register mval extensions")
	}
	if err := tdx.SetTdxQeMvalExtensions(tdx.Evidence, &measurement.Val); err != nil {
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
	// d9023ba100a1008182a100a300d86f4c6086480186f84d01020304050171496e74656c20436f72706f726174696f6e02703031323334353637383941424344454681a101a538480a385046c000fbff000038538282015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7582075830e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f3638540138550b
	// {"ev-triples":{"evidence-triples":[{"environment":{"class":{"id":{"type":"oid","value":"2.16.840.1.113741.1.2.3.4.5"},"vendor":"Intel Corporation","model":"0123456789ABCDEF"}},"measurements":[{"value":{"isvsvn":{"type":"uint","value":10},"miscselect":"wAD7/wAA","mrsigner":{"type":"digest","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU=","sha-384;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXXkW3L1wMC1cttNjTq36X82"]},"isvprodid":{"type":"uint","value":1},"tcbevalnum":{"type":"uint","value":11}}}]}]}}
}

func Example_encode_tdx_qe_evidence_with_profile() {
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
		Class: comid.NewClassOID(TestQEOID).
			SetVendor("Intel Corporation").
			SetModel("0123456789ABCDEF"), // From irim-qe-cend.diag, CPUID[0x01].EAX.FMSP & 0x0FFF0FF0"),
	}

	measurement := &comid.Measurement{}
	valTriple.Measurements.Add(measurement)
	coEv.EvTriples.EvidenceTriples.Add(valTriple)

	if err = tdx.SetTdxQeMvalExtensions(tdx.Evidence, &coEv.EvTriples.EvidenceTriples.Values[0].Measurements.Values[0].Val); err != nil {
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
	// d9023ba100a1008182a100a300d86f4c6086480186f84d01020304050171496e74656c20436f72706f726174696f6e02703031323334353637383941424344454681a101a538480a385046c000fbff000038538282015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7582075830e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f3638540138550b
	// {"ev-triples":{"evidence-triples":[{"environment":{"class":{"id":{"type":"oid","value":"2.16.840.1.113741.1.2.3.4.5"},"vendor":"Intel Corporation","model":"0123456789ABCDEF"}},"measurements":[{"value":{"isvsvn":{"type":"uint","value":10},"miscselect":"wAD7/wAA","mrsigner":{"type":"digest","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU=","sha-384;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXXkW3L1wMC1cttNjTq36X82"]},"isvprodid":{"type":"uint","value":1},"tcbevalnum":{"type":"uint","value":11}}}]}]}}
}

var (
	//go:embed testcases/ce-qe-evidence.cbor
	testCeQeEvidence []byte
)

func Example_decode_QE_Evidence_CBOR() {
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

	if err := ce.FromCBOR(testCeQeEvidence); err != nil {
		panic(err)
	}
	if err := ce.Valid(); err != nil {
		panic(err)
	}
	if err := ExtractQeEvidence(ce); err != nil {
		panic(err)
	}

	// OUTPUT:
	// OID: 2.16.840.1.113741.1.15.4.99.1
	// Vendor: Intel Corporation
	// Model: TDX QE TCB
	// tcbEvalNum: 11
	// IsvProdID: 1
	// mrsigner Digest Alg: 1
	// mrsigner Digest Value: e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75
	// mrsigner Digest Alg: 8
	// mrsigner Digest Value: a314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6aa314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6a
	// miscselect: a0b0c0d000000000
	// TEE TCB Status = UpToDate
	// Tee AdvisoryID = INTEL-SA-00078
	// Tee AdvisoryID = INTEL-SA-00079
	// CryptoKeys: [-----BEGIN PUBLIC KEY-----
	// MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg==
	// -----END PUBLIC KEY-----]
}

func ExtractQeEvidence(ce *coev.TaggedConciseEvidence) error {
	if ce.EvTriples.EvidenceTriples == nil {
		return errors.New("no evidence triples")
	}

	for i, ev := range ce.EvTriples.EvidenceTriples.Values {
		if err := tdx.ExtractQeMeas(ev); err != nil {
			return fmt.Errorf("bad evidence at index %d: %w", i, err)
		}
	}
	return nil
}
