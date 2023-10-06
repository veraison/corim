// Copyright 2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package corim

import (
	"github.com/veraison/corim/extensions"
)

type IEntityConstrainer interface {
	ConstrainEntity(*Entity) error
}

type ICorimConstrainer interface {
	ConstrainCorim(*UnsignedCorim) error
}

type ISignerConstrainer interface {
	ConstrainSigner(*Signer) error
}

type Extensions struct {
	extensions.Extensions
}

func (o *Extensions) validEntity(entity *Entity) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IExtensionsValue.(IEntityConstrainer)
	if ok {
		if err := ev.ConstrainEntity(entity); err != nil {
			return err
		}
	}

	return nil
}

func (o *Extensions) validCorim(c *UnsignedCorim) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IExtensionsValue.(ICorimConstrainer)
	if ok {
		if err := ev.ConstrainCorim(c); err != nil {
			return err
		}
	}

	return nil
}

func (o *Extensions) validSigner(signer *Signer) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IExtensionsValue.(ISignerConstrainer)
	if ok {
		if err := ev.ConstrainSigner(signer); err != nil {
			return err
		}
	}

	return nil
}
