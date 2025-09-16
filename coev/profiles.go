// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package coev

import (
	"fmt"
	"reflect"

	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
)

// ProfileManifest associates an EAT profile ID with a set of extensions. It allows
// obtaining new Concise Evidence structure that had associated extensions
// registered.
type ProfileManifest struct {
	ID            *eat.Profile
	MapExtensions extensions.Map
}

// CoevMapExtensionPoints is a list of extensions.Point's valid for a coev.ConciseEvidence.
var CoevMapExtensionPoints = []extensions.Point{
	ExtConciseEvidence,
	ExtEvTriples,
	ExtEvidenceTriples,
	ExtEvidenceTriplesFlags,
}

var profilesRegister = make(map[string]ProfileManifest)

// AllExtensionPoints is a list of all valid extension.Point's
var AllExtensionPoints = make(map[extensions.Point]bool) // populated inside init() below

type iextensible interface {
	RegisterExtensions(exts extensions.Map) error
}

// GetCoev returns a pointer to a new coev.ConciseEvidence that had the target ProfileManifest's
// extensions (if any) registered.
func (o *ProfileManifest) GetConciseEvidence() *ConciseEvidence {
	ret := NewConciseEvidence()
	o.registerExtensions(ret, CoevMapExtensionPoints)
	return ret
}

// GetTaggedConciseEvidence returns a pointer to a new TaggedConciseEvidence that had the target
// ProfileManifest's extensions (if any) registered.
func (o *ProfileManifest) GetTaggedConciseEvidence() *TaggedConciseEvidence {
	r := o.GetConciseEvidence()
	ret, err := NewTaggedConciseEvidence(r)
	if err != nil {
		panic(err)
	}
	return ret
}

func (o *ProfileManifest) registerExtensions(e iextensible, points []extensions.Point) {
	exts := extensions.NewMap()
	for _, p := range points {
		if v, ok := o.MapExtensions[p]; ok {
			exts[p] = v
		}
	}

	if err := e.RegisterExtensions(exts); err != nil {
		// exts is a subset of o.MapExtensions which have been
		// validated when the profile was registered, so we should never
		// get here.
		panic(err)
	}
}

// RegisterProfile registers a set of extensions with the specified profile. If
// the profile has already been registered, or if the extensions are invalid,
// an error is returned.
func RegisterProfile(id *eat.Profile, exts extensions.Map) error {
	strID, err := id.Get()
	if err != nil {
		return err
	}

	if _, ok := profilesRegister[strID]; ok {
		return fmt.Errorf("profile with id %q already registered", strID)
	}

	for p, v := range exts {
		if _, ok := AllExtensionPoints[p]; !ok {
			return fmt.Errorf("%w: %q", extensions.ErrUnexpectedPoint, p)
		}

		if reflect.TypeOf(v).Kind() != reflect.Pointer {
			return fmt.Errorf("attempting to register a non-pointer IMapValue for %q", p)
		}
	}

	profilesRegister[strID] = ProfileManifest{ID: id, MapExtensions: exts}

	return nil
}

// GetProfileManifest returns the ProfileManifest associated with the specified ID, or an empty
// profileManifest if no ProfileManifest has been registered for the ID. The second return
// value indicates whether a profileManifest for the ID has been found.
func GetProfileManifest(id *eat.Profile) (*ProfileManifest, bool) {
	if id == nil {
		return nil, false
	}

	strID, err := id.Get()
	if err != nil {
		return nil, false
	}

	prof, ok := profilesRegister[strID]
	return &prof, ok
}

// UnregisterProfile ensures there are no extensions registered for the
// specified profile ID. Returns true if extensions were previously registered
// and have been removed, and false otherwise.
func UnregisterProfile(id *eat.Profile) bool {
	if id == nil {
		return false
	}

	strID, err := id.Get()
	if err != nil {
		return false
	}

	if _, ok := profilesRegister[strID]; ok {
		delete(profilesRegister, strID)
		return true
	}

	return false
}

// UnmarshalConciseEvidenceFromCBOR unmarshals a ConciseEvidence from provided CBOR data. If
// there are extensions associated with the profile specified by the data, they
// will be registered with the coev.ConciseEvidence before it is unmarshaled.
func UnmarshalConciseEvidenceFromCBOR(buf []byte, profileID *eat.Profile) (*ConciseEvidence, error) {
	var ret *ConciseEvidence

	profileManifest, ok := GetProfileManifest(profileID)
	if ok {
		ret = profileManifest.GetConciseEvidence()
	} else {
		ret = NewConciseEvidence()
	}

	if err := ret.FromCBOR(buf); err != nil {
		return nil, err
	}

	return ret, nil
}

// GetConciseEvidence returns a pointer to a new ConciseEvidence instance. If there
// are extensions associated with the provided profileID, they will be
// registered with the instance.
func GetConciseEvidence(profileID *eat.Profile) *ConciseEvidence {
	var ret *ConciseEvidence

	if profileID == nil {
		ret = NewConciseEvidence()
	} else {
		profileManifest, ok := GetProfileManifest(profileID)
		if !ok {
			// unknown profile -- treat here like an unprofiled
			// ConciseEvidence. While the ConciseEvidence spec states that unknown
			// profiles should be rejected, we're not actually
			// validating the profile here, just trying to identify
			// any extensions we may need to load. Profile
			// validation is left up to the calling code, as a
			// profile only needs to be registered here if it
			// defines extensions. Profiles that do not add any
			// additional fields may not be registered.
			ret = NewConciseEvidence()
		} else {
			ret = profileManifest.GetConciseEvidence()
		}
	}

	return ret
}

func init() {
	for _, p := range CoevMapExtensionPoints {
		AllExtensionPoints[p] = true
	}
}
