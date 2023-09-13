// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"reflect"

	cbor "github.com/fxamacker/cbor/v2"
)

var (
	em, emError = initCBOREncMode()
	dm, dmError = initCBORDecMode()

	comidTagsMap = map[uint64]interface{}{
		32:  TaggedURI(""),
		37:  TaggedUUID{},
		111: TaggedOID{},
		// CoMID tags
		550: TaggedUEID{},
		551: TaggedInt(0),
		552: TaggedSVN(0),
		553: TaggedMinSVN(0),
		554: TaggedPKIXBase64Key(""),
		555: TaggedPKIXBase64Cert(""),
		556: TaggedPKIXBase64CertPath(""),
		557: TaggedThumbprint{},
		558: TaggedCOSEKey{},
		559: TaggedCertThumbprint{},
		560: TaggedRawValueBytes{},
		561: TaggedCertPathThumbprint{},
		// PSA profile tags
		600: TaggedImplID{},
		601: TaggedPSARefValID{},
		602: TaggedCCAPlatformConfigID(""),
	}
)

func comidTags() cbor.TagSet {
	opts := cbor.TagOptions{
		EncTag: cbor.EncTagRequired,
		DecTag: cbor.DecTagRequired,
	}

	tags := cbor.NewTagSet()

	for tag, typ := range comidTagsMap {
		if err := tags.Add(opts, reflect.TypeOf(typ), tag); err != nil {
			panic(err)
		}
	}

	return tags
}

func initCBOREncMode() (en cbor.EncMode, err error) {
	encOpt := cbor.EncOptions{
		IndefLength: cbor.IndefLengthForbidden,
		TimeTag:     cbor.EncTagRequired,
	}
	return encOpt.EncModeWithTags(comidTags())
}

func initCBORDecMode() (dm cbor.DecMode, err error) {
	decOpt := cbor.DecOptions{
		IndefLength: cbor.IndefLengthForbidden,
	}
	return decOpt.DecModeWithTags(comidTags())
}

func registerCOMIDTag(tag uint64, t interface{}) error {
	if _, exists := comidTagsMap[tag]; exists {
		return fmt.Errorf("tag %d is already registered", tag)
	}

	comidTagsMap[tag] = t

	var err error

	em, err = initCBOREncMode()
	if err != nil {
		return err
	}

	dm, err = initCBORDecMode()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	if emError != nil {
		panic(emError)
	}
	if dmError != nil {
		panic(dmError)
	}
}
