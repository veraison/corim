// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"errors"
	"fmt"

	"github.com/veraison/corim/comid"
)

func extractClassElements(c *comid.Class) error {
	if c == nil {
		return errors.New("no class")
	}

	classID := c.ClassID

	if classID == nil {
		return errors.New("no class-id")
	}

	if classID.Type() != comid.OIDType {
		return errors.New("class id is not an oid")
	}

	fmt.Printf("OID: %s", classID.Value.String())

	if c.Vendor == nil {
		return errors.New("no Vendor")
	}
	fmt.Printf("\nVendor: %s", *c.Vendor)

	if c.Model == nil {
		return fmt.Errorf("no Model")
	}
	fmt.Printf("\nModel: %s", *c.Model)

	return nil
}

func extractDigest(typ string, d *Digests) error {
	if d == nil {
		return errors.New("no digest")
	}

	for _, digest := range *d {
		fmt.Printf("\n%s Digest Alg: %d", typ, digest.HashAlgID)
		fmt.Printf("\n%s Digest Value: %x", typ, digest.HashValue)
	}

	return nil
}

func extractTeeDigest(typ string, d *TeeDigest) error {
	if d == nil {
		return errors.New("no TEE digest")
	}

	if typ != "mrsigner" && typ != "mrtee" {
		return fmt.Errorf("invalid type for TEE digest: %s", typ)
	}
	if err := d.Valid(); err != nil {
		return fmt.Errorf("invalid TEE Digest: %w", err)
	}

	if d.IsDigests() {
		dg, err := d.GetDigest()
		if err != nil {
			return fmt.Errorf("unable to extract TEE Digest: %w", err)
		}
		err = extractDigest(typ, &dg)
		if err != nil {
			return fmt.Errorf("unable to extract %s Digest: %w", typ, err)
		}
	} else if d.IsDigestExpr() {
		de, err := d.GetDigestExpr()
		if err != nil {
			return fmt.Errorf("unable to extract TEE Digest Expression: %w", err)
		}
		fmt.Printf("\n%s Digest Operator: %s", typ, NumericOperatorToString[de.SetOperator])
		dg := comid.Digests(de.SetDigest)
		err = extractDigest(typ, &dg)
		if err != nil {
			return fmt.Errorf("unable to extract %s Digest: %w", typ, err)
		}
	} else {
		return fmt.Errorf("teedigest neither a valid digest or a digest expression for type: %s", typ)
	}
	return nil
}

func extractTeeISVProdID(isvprodID *TeeISVProdID) error {
	if isvprodID == nil {
		return errors.New("isvprodID is nil")
	}

	if isvprodID.IsBytes() {
		val, err := isvprodID.GetBytes()
		if err != nil {
			return fmt.Errorf("failed to decode isvprodid: %w", err)
		}
		fmt.Printf("\nIsvProdID: %x", val)
	} else if isvprodID.IsUint() {
		val, err := isvprodID.GetUint()
		if err != nil {
			return fmt.Errorf("failed to decode isvprodid: %w", err)
		}
		fmt.Printf("\nIsvProdID: %d", val)
	} else {
		return errors.New("isvprodid is neither integer or byte string")
	}
	return nil
}

func extractTeeTcbEvalNum(tcbEvalNum *TeeTcbEvalNumber) error {
	if tcbEvalNum == nil {
		return errors.New("tcbevalnum is nil")
	}
	if tcbEvalNum.IsExpression() {
		ne, err1 := tcbEvalNum.GetNumericExpression()
		if err1 != nil {
			return fmt.Errorf("failed to get tcbEvalNum numeric expression: %w", err1)
		}
		fmt.Printf("\ntcbEvalNum Operator: %s", NumericOperatorToString[ne.NumericOperator])
		fmt.Printf("\ntcbEvalNum Value: %d", ne.NumericType.val)
	} else if tcbEvalNum.IsUint() {
		nv, err1 := tcbEvalNum.GetUint()
		if err1 != nil {
			return fmt.Errorf("failed to get tcbEvalNum uint: %w", err1)
		}
		fmt.Printf("\ntcbEvalNum: %d", nv)
	}
	return nil
}

func decodeAuthorisedBy(m *comid.Measurement) error {
	if m.AuthorizedBy == nil {
		return fmt.Errorf("no authorized-by keys")
	}
	if err := m.AuthorizedBy.Valid(); err != nil {
		return fmt.Errorf("invalid cryptokeys: %w", err)
	}
	fmt.Printf("\nCryptoKeys: %s", m.AuthorizedBy.String())
	return nil
}

func extractTeeSvn(teesvn *TeeSVN) error {
	if teesvn == nil {
		return errors.New("teesvn is nil")
	}
	if teesvn.IsUint() {
		svn, err := teesvn.GetUint()
		if err != nil {
			return fmt.Errorf("unable to get Uint SVN at index: %w", err)
		}
		fmt.Printf("\nISVSVN: %d", svn)
	} else if teesvn.IsExpression() {
		svn, err := teesvn.GetNumericExpression()
		if err != nil {
			return fmt.Errorf("unable to get SVN Expression: %w", err)
		}
		fmt.Printf("\nSVN Operator: %s", NumericOperatorToString[svn.NumericOperator])
		fmt.Printf("\nSVN Value: %d", svn.NumericType.val)
	} else {
		return fmt.Errorf("teesvn, is neither uint or numeric")
	}
	return nil
}

// nolint:dupl
func extractTeeAdvisoryIDs(teeadv *TeeAdvisoryIDs) error {
	if teeadv == nil {
		return errors.New("teeadv is nil")
	}

	if teeadv.IsStringExpr() {
		texp, err := teeadv.GetStringExpression()
		if err != nil {
			return fmt.Errorf("unable to get Tee Advisory Expression: %w", err)
		} else if len(texp.SetString) == 0 {
			return errors.New("zero len Tee Advisory ID Strings")
		}

		fmt.Printf("\nTeeAdvisory Operator: %s", NumericOperatorToString[texp.SetOperator])
		for _, ta := range texp.SetString {
			fmt.Printf("\nTee AdvisoryID = %s", ta)
		}
	} else if teeadv.IsString() {
		tadv, err := teeadv.GetString()
		if err != nil {
			return fmt.Errorf("unable to get Tee Advisory String: %w", err)
		} else if len(tadv) == 0 {
			return errors.New("zero len Tee Advisory ID Strings")
		}
		for _, ta := range tadv {
			fmt.Printf("\nTee AdvisoryID = %s", ta)
		}
	} else {
		return errors.New("teeadv is neither a stringExpr nor a String array")
	}
	return nil
}
