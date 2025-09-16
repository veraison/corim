// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package tdx

import (
	"errors"
	"fmt"

	"github.com/veraison/corim/comid"
)

func SetTdxPceMvalExtensions(typ MessageType, val *comid.Mval) error {
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
	var c *TeeTcbCompSvn
	switch typ {
	case Evidence:
		c, err = NewTeeTcbCompSvnUint(TestCompSvn)
		if err != nil {
			return fmt.Errorf("failed to get TeeTcbCompSvn: %w", err)
		}
	case ReferenceValue:
		c, err = NewTeeTcbCompSvnExpression(TestCompSvn)
		if err != nil {
			return fmt.Errorf("failed to get TeeTcbCompSvn: %w", err)
		}
	}

	err = val.Set("tcbcompsvn", c)
	if err != nil {
		return fmt.Errorf("unable to set teetcbcompsvn: %w", err)
	}

	return nil
}

func extractPceMeasurements(meas *comid.Measurements) error {
	if len(meas.Values) == 0 {
		return errors.New("no measurements")
	}
	for i := range meas.Values {
		m := &meas.Values[0]
		if err := decodePCEMValExtensions(m); err != nil {
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

	err = extractCompSVN(tcs)
	if err != nil {
		return fmt.Errorf("unable to extract TeeTcbCompSVN: %w", err)
	}
	return nil
}

func extractCompSVN(s *TeeTcbCompSvn) error {
	if s == nil {
		return errors.New("no TEE TCB Comp SVN")
	}

	if len(*s) > 16 {
		return errors.New("computed SVN cannot be greater than 16")
	}

	for i, teesvn := range *s {
		svn := teesvn // Avoid gosec: Implicit memory aliasing in for loop
		if err := extractTeeSvn(&svn); err != nil {
			return fmt.Errorf("unable to extract SVN at index %d: %w", i, err)
		}
	}

	return nil
}

func ExtractPceMeas(rv comid.ValueTriple) error {
	class := rv.Environment.Class

	if err := extractClassElements(class); err != nil {
		return fmt.Errorf("extracting class: %w", err)
	}

	measurements := rv.Measurements
	if err := extractPceMeasurements(&measurements); err != nil {
		return fmt.Errorf("extracting measurements: %w", err)
	}

	return nil
}
