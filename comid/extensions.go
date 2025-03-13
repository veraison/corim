// Copyright 2023-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"github.com/veraison/corim/extensions"
)

const (
	ExtComid                       extensions.Point = "Comid"
	ExtEntity                      extensions.Point = "ComidEntity"
	ExtTriples                     extensions.Point = "Triples"
	ExtReferenceValue              extensions.Point = "ReferenceValue"
	ExtReferenceValueFlags         extensions.Point = "ReferenceValueFlags"
	ExtEndorsedValue               extensions.Point = "EndorsedValue"
	ExtEndorsedValueFlags          extensions.Point = "EndorsedValueFlags"
	ExtCondEndorseSeriesValue      extensions.Point = "CondEndorseSeriesValue"
	ExtCondEndorseSeriesValueFlags extensions.Point = "CondEndorseSeriesValueFlags"
	ExtMval                        extensions.Point = "Mval"
	ExtFlags                       extensions.Point = "Flags"
)

type IComidConstrainer interface {
	ConstrainComid(*Comid) error
}

type ITriplesConstrainer interface {
	ValidTriples(*Triples) error
}

type IMvalConstrainer interface {
	ConstrainMval(*Mval) error
}

type IEntityConstrainer interface {
	ConstrainEntity(*Entity) error
}

type IFlagsMapConstrainer interface {
	ConstrainFlagsMap(*FlagsMap) error
}

type IFlagSetter interface {
	AnySet() bool
	SetTrue(Flag)
	SetFalse(Flag)
	Clear(Flag)
	Get(Flag) *bool
}

type Extensions struct {
	extensions.Extensions
}

func (o *Extensions) validComid(comid *Comid) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IMapValue.(IComidConstrainer)
	if ok {
		if err := ev.ConstrainComid(comid); err != nil {
			return err
		}
	}

	return nil
}

func (o *Extensions) validTriples(triples *Triples) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IMapValue.(ITriplesConstrainer)
	if ok {
		if err := ev.ValidTriples(triples); err != nil {
			return err
		}
	}

	return nil
}

func (o *Extensions) validMval(triples *Mval) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IMapValue.(IMvalConstrainer)
	if ok {
		if err := ev.ConstrainMval(triples); err != nil {
			return err
		}
	}

	return nil
}

func (o *Extensions) validEntity(triples *Entity) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IMapValue.(IEntityConstrainer)
	if ok {
		if err := ev.ConstrainEntity(triples); err != nil {
			return err
		}
	}

	return nil
}

func (o *Extensions) validFlagsMap(triples *FlagsMap) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IMapValue.(IFlagsMapConstrainer)
	if ok {
		if err := ev.ConstrainFlagsMap(triples); err != nil {
			return err
		}
	}

	return nil
}

func (o *Extensions) setTrue(flag Flag) {
	if !o.HaveExtensions() {
		return
	}

	ev, ok := o.IMapValue.(IFlagSetter)
	if ok {
		ev.SetTrue(flag)
	}
}

func (o *Extensions) setFalse(flag Flag) {
	if !o.HaveExtensions() {
		return
	}

	ev, ok := o.IMapValue.(IFlagSetter)
	if ok {
		ev.SetFalse(flag)
	}
}

func (o *Extensions) clear(flag Flag) {
	if !o.HaveExtensions() {
		return
	}

	ev, ok := o.IMapValue.(IFlagSetter)
	if ok {
		ev.Clear(flag)
	}
}

func (o *Extensions) get(flag Flag) *bool {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IMapValue.(IFlagSetter)
	if ok {
		return ev.Get(flag)
	}

	return nil
}

func (o *Extensions) anySet() bool {
	if !o.HaveExtensions() {
		return false
	}

	ev, ok := o.IMapValue.(IFlagSetter)
	if ok {
		return ev.AnySet()
	}

	return false
}
