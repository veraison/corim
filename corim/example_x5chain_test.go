// Copyright 2021-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"fmt"
	"log"
	"time"

	"github.com/veraison/corim/testdata"
)

func x5chainExampleMeta() *Meta {
	notAfter := time.Date(2021, time.October, 0, 0, 0, 0, 0, time.UTC)

	return NewMeta().
		SetSigner("ACME Ltd.", nil).
		SetValidity(notAfter, nil)
}

func exampleSignedCorimWithX5Chain() ([]byte, error) {
	signer, err := NewSignerFromJWK(testEndEntityKey)
	if err != nil {
		return nil, err
	}

	var unsigned UnsignedCorim
	if err := unsigned.FromCBOR(testGoodUnsignedCorimCBOR); err != nil {
		return nil, err
	}

	var signed SignedCorim
	signed.UnsignedCorim = unsigned
	signed.Meta = *x5chainExampleMeta()

	if err := signed.AddSigningCert(testdata.EndEntityDer); err != nil {
		return nil, err
	}

	intermediates := make([]byte, len(testdata.IntermediateCA)+len(testdata.RootCA))
	copy(intermediates, testdata.IntermediateCA)
	copy(intermediates[len(testdata.IntermediateCA):], testdata.RootCA)

	if err := signed.AddIntermediateCerts(intermediates); err != nil {
		return nil, err
	}

	return signed.Sign(signer)
}

func ExampleSignedCorim_VerifyWithX5Chain() {
	cbor, err := exampleSignedCorimWithX5Chain()
	if err != nil {
		log.Fatal(err)
	}

	anchors, err := LoadTrustAnchors(func(path string) ([]byte, error) {
		if path != "anchor.der" {
			return nil, fmt.Errorf("unknown path %q", path)
		}

		return testdata.RootCA, nil
	}, []string{"anchor.der"}, nil)
	if err != nil {
		log.Fatal(err)
	}

	var signed SignedCorim
	if err := signed.FromCOSE(cbor); err != nil {
		log.Fatal(err)
	}

	if err := signed.VerifyWithX5Chain(anchors); err != nil {
		log.Fatal(err)
	}
	// Output:
}

func ExampleLoadTrustAnchors() {
	anchors, err := LoadTrustAnchors(func(path string) ([]byte, error) {
		switch path {
		case "anchor.der":
			return testdata.RootCA, nil
		default:
			return nil, fmt.Errorf("unknown path %q", path)
		}
	}, []string{"anchor.der"}, nil)
	if err != nil {
		log.Fatal(err)
	}

	if anchors.Pool == nil {
		log.Fatal("expected explicit trust-anchor pool")
	}
	// Output:
}
