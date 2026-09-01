// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cotl

import (
	"fmt"

	"github.com/veraison/corim/extensions"
)

// ExtConciseTlTag is the extension point on the top-level concise-tl-tag map.
const ExtConciseTlTag extensions.Point = "ConciseTlTag"

// RegisterExtensions registers a struct as a collection of extensions on this
// CoTL. The struct is a pointer to a user-defined type whose exported fields
// carry cbor/json key tags, following the mechanism used by
// github.com/veraison/corim.
//
// An extension struct may additionally implement
//
//	ConstrainConciseTlTag(*ConciseTlTag) error
//
// in which case it participates in structural validation.
func (o *ConciseTlTag) RegisterExtensions(exts extensions.Map) error {
	for p, v := range exts {
		switch p {
		case ExtConciseTlTag:
			o.Register(v)
		default:
			return fmt.Errorf("%w: %q", extensions.ErrUnexpectedPoint, p)
		}
	}

	return nil
}

// GetExtensions returns the registered extension struct, or nil.
func (o *ConciseTlTag) GetExtensions() extensions.IMapValue {
	return o.IMapValue
}

// validExtensions invokes any registered constrainer.
func (o *ConciseTlTag) validExtensions() error {
	if !o.HaveExtensions() {
		return nil
	}

	if ev, ok := o.IMapValue.(IConciseTlTagConstrainer); ok {
		if err := ev.ConstrainConciseTlTag(o); err != nil {
			return err
		}
	}

	return nil
}

// IConciseTlTagConstrainer may be implemented by extension structs to
// participate in structural validation of the enclosing CoTL.
type IConciseTlTagConstrainer interface {
	ConstrainConciseTlTag(*ConciseTlTag) error
}
