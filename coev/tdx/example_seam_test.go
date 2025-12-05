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

func Example_decode_Seam_Evidence_JSON() {
	manifest, found := coev.GetProfileManifest(ProfileID)
	if !found {
		fmt.Printf("Evidence profile not found")
		return
	}
	ce := manifest.GetTaggedConciseEvidence()
	if err := ce.FromJSON([]byte(TDXSeamCETemplate)); err != nil {
		panic(err)
	}
	if err := ce.Valid(); err != nil {
		panic(err)
	}

	if err := ExtractSeamEvidence(ce); err != nil {
		panic(err)
	}

	// OUTPUT:
	// OID: 2.16.840.1.113741.1.2.3.4.3
	// Vendor: Intel Corporation
	// Model: TDXSEAM
	// tcbEvalNum: 11
	// IsvProdID: 0101
	// ISVSVN: 10
	// Attributes: 0101
	// mrtee Digest Alg: 1
	// mrtee Digest Value: e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75
	// mrsigner Digest Alg: 1
	// mrsigner Digest Value: e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75
	// mrsigner Digest Alg: 7
	// mrsigner Digest Value: e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f36
}

func Example_encode_tdx_seam_evidence_without_profile() {
	coEv := coev.NewConciseEvidence()
	valTriple := &comid.ValueTriple{}
	measurement := &comid.Measurement{}
	valTriple.Environment = comid.Environment{
		Class: comid.NewClassOID(TestSeamOID).
			SetVendor("Intel Corporation").
			SetModel("TDXSEAM"),
	}

	extMap := extensions.NewMap().
		Add(comid.ExtMval, &tdx.MValExtensions{})

	if err := measurement.Val.RegisterExtensions(extMap); err != nil {
		log.Fatal("could not register mval extensions")
	}
	if err := tdx.SetTDXSeamMvalExtensions(tdx.Evidence, &measurement.Val); err != nil {
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
	// d9023ba100a1008182a100a300d86f4c6086480186f84d01020304030171496e74656c20436f72706f726174696f6e02675444585345414d81a101a73847c11a6796cc8038480a385142010138528182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7538538282015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7582075830e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f36385442010138550b
	// {"ev-triples":{"evidence-triples":[{"environment":{"class":{"id":{"type":"oid","value":"2.16.840.1.113741.1.2.3.4.3"},"vendor":"Intel Corporation","model":"TDXSEAM"}},"measurements":[{"value":{"tcbdate":"2025-01-27T00:00:00Z","isvsvn":{"type":"uint","value":10},"attributes":"AQE=","mrtee":{"type":"digest","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]},"mrsigner":{"type":"digest","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU=","sha-384;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXXkW3L1wMC1cttNjTq36X82"]},"isvprodid":{"type":"bytes","value":"AQE="},"tcbevalnum":{"type":"uint","value":11}}}]}]}}
}

var (
	//go:embed testcases/ce-seam-evidence.cbor
	testCeSeamEvidence []byte
)

func Example_encode_tdx_seam_evidence_with_profile() {
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
		Class: comid.NewClassOID(TestSeamOID).
			SetVendor("Intel Corporation").
			SetModel("TDXSEAM"),
	}
	measurement := &comid.Measurement{}
	valTriple.Measurements.Add(measurement)
	coEv.EvTriples.EvidenceTriples.Add(valTriple)

	if err = tdx.SetTDXSeamMvalExtensions(tdx.Evidence, &coEv.EvTriples.EvidenceTriples.Values[0].Measurements.Values[0].Val); err != nil {
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
	// d9023ba100a1008182a100a300d86f4c6086480186f84d01020304030171496e74656c20436f72706f726174696f6e02675444585345414d81a101a73847c11a6796cc8038480a385142010138528182015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7538538282015820e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d7582075830e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f36385442010138550b
	// {"ev-triples":{"evidence-triples":[{"environment":{"class":{"id":{"type":"oid","value":"2.16.840.1.113741.1.2.3.4.3"},"vendor":"Intel Corporation","model":"TDXSEAM"}},"measurements":[{"value":{"tcbdate":"2025-01-27T00:00:00Z","isvsvn":{"type":"uint","value":10},"attributes":"AQE=","mrtee":{"type":"digest","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU="]},"mrsigner":{"type":"digest","value":["sha-256;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXU=","sha-384;5Fty9cDAtXLbTY06t+l/No/3TmI0eoJN7LZ6hOUiTXXkW3L1wMC1cttNjTq36X82"]},"isvprodid":{"type":"bytes","value":"AQE="},"tcbevalnum":{"type":"uint","value":11}}}]}]}}
}

func Example_decode_CBOR() {
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

	if err := ce.FromCBOR(testCeSeamEvidence); err != nil {
		panic(err)
	}
	if err := ce.Valid(); err != nil {
		panic(err)
	}
	if err := ExtractSeamEvidence(ce); err != nil {
		panic(err)
	}

	// OUTPUT:
	// OID: 2.16.840.1.113741.1.15.4.99.1
	// Vendor: Intel Corporation
	// Model: TDX Seam
	// tcbEvalNum: 11
	// IsvProdID: abcd
	// ISVSVN: 6
	// Attributes: 0102
	// mrtee Digest Alg: 1
	// mrtee Digest Value: e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75
	// mrsigner Digest Alg: 1
	// mrsigner Digest Value: e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75
	// mrsigner Digest Alg: 7
	// mrsigner Digest Value: e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f36
	// CryptoKeys: [-----BEGIN PUBLIC KEY-----
	// MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFn0taoAwR3PmrKkYLtAsD9o05KSM6mbgfNCgpuL0g6VpTHkZl73wk5BDxoV7n+Oeee0iIqkW3HMZT3ETiniJdg==
	// -----END PUBLIC KEY-----]
}

func ExtractSeamEvidence(ce *coev.TaggedConciseEvidence) error {
	if ce.EvTriples.EvidenceTriples == nil {
		return errors.New("no evidence triples")
	}

	for i, ev := range ce.EvTriples.EvidenceTriples.Values {
		if err := tdx.ExtractSeamMeas(ev); err != nil {
			return fmt.Errorf("bad evidence at index %d: %w", i, err)
		}
	}
	return nil
}
