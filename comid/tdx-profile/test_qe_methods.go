// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"errors"
	"fmt"

	"github.com/veraison/corim/comid"
	"github.com/veraison/swid"
)

//nolint:funlen // reason: this function is long but readability is fine
func SetTdxQeMvalExtensions(typ MessageType, val *comid.Mval) error {
	var svn *TeeSVN
	var err error
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

	teeMiscSel := TeeMiscSelect(TestTeeMiscSelect) // Taken from irim-qe-ref.diag
	err = val.Set("miscselect", &teeMiscSel)
	if err != nil {
		return fmt.Errorf("unable to set miscselect: %w", err)
	}

	var ts *TeeDigest
	// assign mrSignerDigest
	dSign := comid.NewDigests()
	dSign.AddDigest(swid.Sha256,
		comid.MustHexDecode(nil, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75"))
	dSign.AddDigest(swid.Sha384,
		comid.MustHexDecode(nil, "e45b72f5c0c0b572db4d8d3ab7e97f368ff74e62347a824decb67a84e5224d75e45b72f5c0c0b572db4d8d3ab7e97f36"))

	// assign mrsignerDigest
	switch typ {
	case ReferenceValue:
		ts, err = NewTeeDigestExpr(MEM, *dSign)
		if err != nil {
			return fmt.Errorf("unable to get TeeDigest: %w", err)
		}
	case Evidence:
		ts, err = NewTeeDigest(*dSign)
		if err != nil {
			return fmt.Errorf("unable to get TeeDigest: %w", err)
		}
	}

	err = val.Set("mrsigner", ts)
	if err != nil {
		return fmt.Errorf("unable to set mrsigner %w", err)
	}

	// Taken below from irim-qe-ref.diag
	isvProdID, err := NewTeeISVProdID(TestUIntISVProdID)
	if err != nil {
		return fmt.Errorf("unable to get isvprodid: %w", err)
	}

	err = val.Set("isvprodid", isvProdID)
	if err != nil {
		return fmt.Errorf("unable to set isvprodid: %w", err)
	}

	return nil
}

func extractQeMeasurements(meas *comid.Measurements) error {
	if len(meas.Values) == 0 {
		return errors.New("no measurements")
	}
	for i := range meas.Values {
		m := &meas.Values[0]
		if err := decodeQeMValExtensions(m); err != nil {
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

func decodeQeMValExtensions(m *comid.Measurement) error {
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

	val, err = m.Val.Get("mrsigner")
	if err != nil {
		return fmt.Errorf("failed to decode mrsigner from measurement extensions: %w", err)
	}

	tD, ok := val.(*TeeDigest)
	if !ok {
		return errors.New("val was not pointer to TeeDigest")
	}

	if err = extractTeeDigest("mrsigner", tD); err != nil {
		return fmt.Errorf("failed to extarct mrsigner digest: %w", err)
	}

	val, err = m.Val.Get("miscselect")
	if err != nil {
		return errors.New("failed to decode miscselect from measurement extensions")
	}
	tm, ok := val.(*TeeMiscSelect)
	if !ok {
		return errors.New("val was not pointer to TeeMiscSelect")
	}
	miscselect := *tm
	fmt.Printf("\nmiscselect: %x", miscselect)

	tst, err := m.Val.Get("tcbstatus")
	if err != nil {
		return errors.New("failed to decode tcb status from measurement extensions")
	}
	ts, ok := tst.(*TeeTcbStatus)
	if !ok {
		return errors.New("val was not pointer to TeeTcbStatus")
	}
	if err = extractTeeTCBStatus(ts); err != nil {
		return fmt.Errorf("failed to extract tee tcb status: %w", err)
	}

	ta, err := m.Val.Get("advisoryids")
	if err != nil {
		return errors.New("failed to decode tee advisory ids from measurement extensions")
	}
	tas, ok := ta.(*TeeAdvisoryIDs)
	if !ok {
		return errors.New("val was not pointer to TeeAdvisoryIDs")
	}
	if err = extractTeeAdvisoryIDs(tas); err != nil {
		return fmt.Errorf("failed to extarct tee advisory ids: %w", err)
	}

	return nil
}

func ExtractQeMeas(rv comid.ValueTriple) error {
	class := rv.Environment.Class

	if err := extractClassElements(class); err != nil {
		return fmt.Errorf("extracting class: %w", err)
	}

	measurements := rv.Measurements
	if err := extractQeMeasurements(&measurements); err != nil {
		return fmt.Errorf("extracting measurements: %w", err)
	}

	return nil
}

// nolint:dupl
func extractTeeTCBStatus(tcb *TeeTcbStatus) error {
	if tcb == nil {
		return errors.New("teeadv is nil")
	}

	if tcb.IsStringExpr() {
		texp, err := tcb.GetStringExpression()
		if err != nil {
			return fmt.Errorf("unable to get Tee TCB StatusExpression: %w", err)
		} else if len(texp.SetString) == 0 {
			return errors.New("zero len Tee TCB Status String")
		}

		fmt.Printf("\nTEE TCB Status Operator: %s", NumericOperatorToString[texp.SetOperator])
		for _, ta := range texp.SetString {
			fmt.Printf("\nTEE TCB Status = %s", ta)
		}
	} else if tcb.IsString() {
		ts, err := tcb.GetString()
		if err != nil {
			return fmt.Errorf("unable to get Tee Tcb Status String: %w", err)
		} else if len(ts) == 0 {
			return errors.New("zero len Tee Tcb Status String")
		}
		for _, t := range ts {
			fmt.Printf("\nTEE TCB Status = %s", t)
		}
	} else {
		return errors.New("tcb status is neither a StringExpr nor a String array")
	}
	return nil
}
