// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"errors"
	"fmt"
	"time"

	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
)

//nolint:funlen // reason: this function is long but readability is fine
func SetTDXSeamMvalExtensions(typ MessageType, val *comid.Mval) error {
	tcbDate, _ := time.Parse(time.RFC3339, "2025-01-27T00:00:00Z")
	err := val.Set("tcbdate", &tcbDate)
	if err != nil {
		return fmt.Errorf("unable to set tcbDate: %w", err)
	}
	r := []byte{0x01, 0x01}
	isvProdID, err := NewTeeISVProdID(r)
	if err != nil {
		return fmt.Errorf("unable to get isvprodid: %w", err)
	}
	err = val.Set("isvprodid", isvProdID)
	if err != nil {
		return fmt.Errorf("unable to set isvprodid: %w", err)
	}
	var svn *TeeSVN
	var teeTcbEvNum *TeeTcbEvalNumber
	switch typ {
	case ReferenceValue:
		svn, err = NewSvnExpression(TestISVSVN)
		if err != nil {
			return fmt.Errorf("unable to get isvsvn numeric: %w", err)
		}
		teeTcbEvNum, err = NewTeeTcbEvalNumberNumeric(TestTCBEvalNum)
		if err != nil {
			return fmt.Errorf("unable to get tcbevalnum numeric: %w", err)
		}
	case Evidence:
		svn, err = NewSvnUint(TestISVSVN)
		if err != nil {
			return fmt.Errorf("unable to get isvsvn uint: %w", err)
		}
		teeTcbEvNum, err = NewTeeTcbEvalNumberUint(TestTCBEvalNum)
		if err != nil {
			return fmt.Errorf("unable to get tcbevalnum uint: %w", err)
		}
	default:
		return fmt.Errorf("unknonw typ: %d", typ)
	}

	// set the populated svn variable
	err = val.Set("isvsvn", svn)
	if err != nil {
		return fmt.Errorf("unable to set isvsvn: %w", err)
	}
	// set the populated teeTcbEvNum
	err = val.Set("tcbevalnum", teeTcbEvNum)
	if err != nil {
		return fmt.Errorf("unable to set tcbevalnum: %w", err)
	}

	teeAttr, err := NewTeeAttributes(TestTeeAttributes)
	if err != nil {
		return fmt.Errorf("unable to get teeAttributes: %w", err)
	}
	err = val.Set("attributes", teeAttr)
	if err != nil {
		return fmt.Errorf("unable to set attributes: %w", err)
	}

	var td, ts *TeeDigest
	// assign mrteeDigest
	dTee := comid.NewDigests()
	dTee.AddDigest(swid.Sha256, comid.MustHexDecode(nil, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"))

	// assign mrSignerDigest
	dSign := comid.NewDigests()
	dSign.AddDigest(swid.Sha256,
		comid.MustHexDecode(nil, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"))
	dSign.AddDigest(swid.Sha384,
		comid.MustHexDecode(nil, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f36"))

	// assign mrsignerDigest
	switch typ {
	case ReferenceValue:
		td, err = NewTeeDigestExpr(MEM, *dTee)
		if err != nil {
			return fmt.Errorf("unable to get TeeDigest Expression: %w", err)
		}

		ts, err = NewTeeDigestExpr(MEM, *dSign)
		if err != nil {
			return fmt.Errorf("unable to get TeeDigest: %w", err)
		}
	case Evidence:
		td, err = NewTeeDigest(*dTee)
		if err != nil {
			return fmt.Errorf("unable to get TeeDigest: %w", err)
		}
		ts, err = NewTeeDigest(*dSign)
		if err != nil {
			return fmt.Errorf("unable to get TeeDigest: %w", err)
		}
	}
	err = val.Set("mrtee", td)
	if err != nil {
		return fmt.Errorf("unable to set mrtee: %w", err)
	}
	err = val.Set("mrsigner", ts)
	if err != nil {
		return fmt.Errorf("unable to set mrsigner %w", err)
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
			err := decodeAuthorisedBy(m)
			if err != nil {
				return fmt.Errorf("extracting measurement at index %d: %w", i, err)
			}
		}
	}
	return nil
}

//nolint:funlen // reason: this function is long but readability is fine
func decodeMValExtensions(m *comid.Measurement) error {
	val, err := m.Val.Get("tcbevalnum")
	if err != nil {
		return fmt.Errorf("failed to decode tcbevalnum from measurement extensions: %w", err)
	}
	f, ok := val.(*TeeTcbEvalNumber)
	if !ok {
		return fmt.Errorf("val pointer to TeeTcbEvalNum NOT OK")
	}
	tcbValNum := *f
	if err = extractTeeTcbEvalNum(&tcbValNum); err != nil {
		return fmt.Errorf("failed to extract tcbevalnum: %w", err)
	}

	val, err = m.Val.Get("isvprodid")
	if err != nil {
		return errors.New("failed to decode isvprodid from measurement extensions")
	}
	tS, ok := val.(*TeeISVProdID)
	if !ok {
		fmt.Printf("val was not pointer to IsvProdID")
	}
	if err = extractTeeISVProdID(tS); err != nil {
		return fmt.Errorf("failed to decode teeISVProdID from measurement extensions: %w", err)
	}

	val, err = m.Val.Get("isvsvn")
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

	err = extractTeeSvn(teesvn)
	if err != nil {
		return fmt.Errorf("unable to extract tee svn: %w", err)
	}
	val, err = m.Val.Get("attributes")
	if err != nil {
		return errors.New("failed to decode attributes from measurement extensions")
	}

	tA, ok := val.(*TeeAttributes)
	if !ok {
		fmt.Printf("val was not pointer to teeAttributes")
	}
	fmt.Printf("\nAttributes: %x", *tA)

	val, err = m.Val.Get("mrtee")
	if err != nil {
		return errors.New("failed to decode mrtee from measurement extensions")
	}

	tD, ok := val.(*TeeDigest)
	if !ok {
		fmt.Printf("val was not pointer to TeeDigest")
	}

	if err = extractTeeDigest("mrtee", tD); err != nil {
		return fmt.Errorf("failed to decode mrtee from digest: %w", err)
	}

	val, err = m.Val.Get("mrsigner")
	if err != nil {
		return fmt.Errorf("failed to decode mrsigner from measurement extensions: %w", err)
	}

	tD, ok = val.(*TeeDigest)
	if !ok {
		return errors.New("val was not pointer to TeeDigest")
	}

	if err := extractTeeDigest("mrsigner", tD); err != nil {
		return fmt.Errorf("failed to extract mrsigner digest: %w", err)
	}

	return nil
}

func ExtractSeamMeas(rv comid.ValueTriple) error {
	class := rv.Environment.Class

	if err := extractClassElements(class); err != nil {
		return fmt.Errorf("extracting class: %w", err)
	}

	measurements := rv.Measurements
	if err := extractSeamMeasurements(&measurements); err != nil {
		return fmt.Errorf("extracting measurements: %w", err)
	}

	return nil
}
