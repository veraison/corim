// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"fmt"

	"github.com/veraison/corim/comid"
)

func TestextractClassElements(c *comid.Class) error {
	if c == nil {
		return fmt.Errorf("no class")
	}

	classID := c.ClassID

	if classID == nil {
		return fmt.Errorf("no class-id")
	}

	if classID.Type() != comid.OIDType {
		return fmt.Errorf("class id is not an oid")
	}

	fmt.Printf("OID: %s", classID.Value.String())

	if c.Vendor == nil {
		return fmt.Errorf("no Vendor")
	}
	fmt.Printf("\nVendor: %s", *c.Vendor)

	if c.Model == nil {
		return fmt.Errorf("no Model")
	}
	fmt.Printf("\nModel: %s", *c.Model)

	return nil
}

func TestextractDigest(typ string, d *Digests) error {
	if d == nil {
		return fmt.Errorf("no digest")
	}

	for _, digest := range *d {
		fmt.Printf("\n%s Digest Alg: %d", typ, digest.HashAlgID)
		fmt.Printf("\n%s Digest Value: %x", typ, digest.HashValue)
	}

	return nil
}

func TestextractTeeDigest(typ string, d *TeeDigest) error {
	if d == nil {
		return fmt.Errorf("no TEE digest")
	}

	if typ != "mrsigner" && typ != "mrtee" {
		return fmt.Errorf("invalid type for TEE digest: %s", typ)
	}

	if d.IsDigests() {
		dg, err := d.GetDigest()
		if err != nil {
			return fmt.Errorf("unable to extract TEE Digest: %w", err)
		}
		err = TestextractDigest(typ, &dg)
		if err != nil {
			return fmt.Errorf("unable to extract %s Digest: %w", typ, err)
		}
	} else {
		de, err := d.GetDigestExpr()
		if err != nil {
			return fmt.Errorf("unable to extract TEE Digest Expression: %w", err)
		}
		fmt.Printf("\n%s Digest Operator: %s", typ, NumericOperatorToString[de.SetOperator])
		dg := comid.Digests(de.SetDigest)
		err = TestextractDigest(typ, &dg)
		if err != nil {
			return fmt.Errorf("unable to extract %s Digest: %w", typ, err)
		}
	}
	return nil
}

func TestextractTeeSvn(teesvn *TeeSVN) error {
	if teesvn == nil {
		return fmt.Errorf("teesvn is nil")
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

func TestextractTeeISVProdID(isvprodID *TeeISVProdID) error {
	if isvprodID == nil {
		return fmt.Errorf("isvprodID is nil")
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
		return fmt.Errorf("isvprodid is neither integer or byte string")
	}
	return nil
}

func TestextractTeeTcbEvalNum(tcbEvalNum *TeeTcbEvalNumber) error {
	if tcbEvalNum == nil {
		return fmt.Errorf("tcbevalnum is nil")
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

func TestdecodeAuthorisedBy(m *comid.Measurement) error {
	if err := m.AuthorizedBy.Valid(); err != nil {
		return fmt.Errorf("invalid cryptokey: %w", err)
	}
	fmt.Printf("\nCryptoKey Type: %s", m.AuthorizedBy.Type())
	fmt.Printf("\nCryptoKey Value: %s", m.AuthorizedBy.String())
	return nil
}
