// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package corim

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
	"github.com/veraison/swid"
)

// ----- profile definition -----
// The following code defines a profile with the following extensions and
// constraints:
// - Entities may contain an address field
// - Reference values may contain a Unix timestamp indicating when the
//   individual measurement was taken.
// - The language claim (CoMID index 0) must be present, and must be "en-GB"

// EntityExtensions will be used for both CoMID and CoRIM entities.
type EntityExtensions struct {
	Address *string `cbor:"-1,keyasint,omitempty" json:"address,omitempty"`
}

type RefValExtensions struct {
	Timestamp *int `cbor:"-1,keyasint,omitempty" json:"timestamp,omitempty"`
}

// We're not defining any additional fields, however we're providing extra
// constraints that will be applied on top of standard CoMID validation.
type ComidExtensions struct{}

func (*ComidExtensions) ConstrainComid(c *comid.Comid) error {
	if c.Language == nil {
		return errors.New("language not specified")
	}

	if *c.Language != "en-GB" {
		return fmt.Errorf(`language must be "en-GB", but found %q`, *c.Language)
	}

	return nil
}

// Registering the profile inside init() in the same file where it is defined
// ensures that the profile will always be available, and you don't need to
// remember to register it at the time you want to use it. The only potential
// danger with that is if the your profile ID clashes with another profile,
// which should not happen if it a registered PEN or a URL containing a domain
// that you own.
func init() {
	profileID, err := eat.NewProfile("http://example.com/example-profile")
	if err != nil {
		panic(err) // will not error, as the hard-coded string above is valid
	}

	extMap := extensions.NewMap().
		Add(ExtEntity, &EntityExtensions{}).
		Add(comid.ExtComid, &ComidExtensions{}).
		Add(comid.ExtEntity, &EntityExtensions{}).
		Add(comid.ExtReferenceValue, &RefValExtensions{})

	if err := RegisterProfile(profileID, extMap); err != nil {
		// will not error, assuming our profile ID is unique, and we've
		// correctly set up the extensions Map above
		panic(err)
	}
}

// ----- end of profile definition -----
// The following code demonstrates how the profile might be used.

func Example_profile_unmarshal() {
	buf, err := os.ReadFile("testcases/unsigned-example-corim.cbor")
	if err != nil {
		log.Fatalf("could not read corim file: %v", err)
	}

	// UnmarshalUnsignedCorimFromCBOR will detect the profile and ensure
	// the correct extensions are loaded before unmarshalling
	extractedCorim, err := UnmarshalUnsignedCorimFromCBOR(buf)
	if err != nil {
		log.Fatalf("could not unmarshal corim: %v", err)
	}

	extractedComid, err := UnmarshalComidFromCBOR(
		extractedCorim.Tags[0].Content,
		extractedCorim.Profile,
	)
	if err != nil {
		log.Fatalf("could not unmarshal corim: %v", err)
	}

	fmt.Printf("Language: %s\n", *extractedComid.Language)
	fmt.Printf("Entity: %s\n", *extractedComid.Entities.Values[0].Name)
	fmt.Printf("        %s\n", extractedComid.Entities.Values[0].MustGetString("Address"))

	fmt.Printf("Measurements:\n")
	for i := range extractedComid.Triples.ReferenceValues.Values[0].Measurements.Values {
		m := &extractedComid.Triples.ReferenceValues.Values[0].Measurements.Values[i]

		val := hex.EncodeToString((*m.Val.Digests)[0].HashValue)
		tsInt := m.Val.MustGetInt64("timestamp")
		ts := time.Unix(tsInt, 0).UTC()

		fmt.Printf("    %v taken at %s\n", val, ts.Format("2006-01-02T15:04:05"))
	}

	// Output:
	// Language: en-GB
	// Entity: ACME Ltd.
	//         123 Fake Street
	// Measurements:
	//     87428fc522803d31065e7bce3cf03fe475096631e5e07bbd7a0fde60c4cf25c7 taken at 2024-07-12T11:03:10
	//     0263829989b6fd954f72baaf2fc64bc2e2f01d692d4de72986ea808f6e99813f taken at 2024-07-12T11:03:10
	//     a3a5e715f0cc574a73c3f9bebb6bc24f32ffd5b67b387244c2c909da779a1478 taken at 2024-07-12T11:03:10
}

// note: this example is rather verbose as we're going to be constructing a
// CoMID by hand. In practice, you would typically write a JSON document and
// then unmarshal that into a CoRIM before marshaling it into CBOR (in which
// case, extensions will work as with unmarshaling example above).
func Example_profile_marshal() {
	profileID, err := eat.NewProfile("http://example.com/example-profile")
	if err != nil {
		panic(err)
	}

	profileManifest, ok := GetProfileManifest(profileID)
	if !ok {
		log.Fatalf("profile %v not found", profileID)
	}

	myComid := profileManifest.GetComid().
		SetLanguage("en-GB").
		SetTagIdentity("example", 0).
		// Adding an entity to the Entities collection also registers
		// profile's extensions
		AddEntity("ACME Ltd.", &comid.TestRegID, comid.RoleCreator)

	address := "123 Fake Street"
	err = myComid.Entities.Values[0].Set("Address", &address)
	if err != nil {
		log.Fatalf("could not set entity Address: %v", err)
	}

	// Use generic UUID-based class ID instead of PSA-specific impl-id
	refVal := comid.ValueTriple{
		Environment: comid.Environment{
			Class: comid.NewClassUUID(comid.TestUUID).
				SetVendor("ACME Ltd.").
				SetModel("RoadRunner 2.0"),
		},
		Measurements: *comid.NewMeasurements(),
	}

	// Use generic UUID measurement key instead of PSA-specific refval-id
	measurement, err := comid.NewUUIDMeasurement(comid.TestUUID)
	if err != nil {
		log.Fatalf("could not create measurement: %v", err)
	}
	measurement.AddDigest(swid.Sha256_32, []byte{0xab, 0xcd, 0xef, 0x00})

	// alternatively, we can add extensions to individual value before
	// adding it to the collection. Note that because we're adding the
	// extension directly to the measurement, we're using a different
	// extension point, comid.ExtMval rather than comid.ExtReferenceValue,
	// as a measurement doesn't know that its going to be part of reference
	// value, and so is unaware of reference value extension points.
	extMap := extensions.NewMap().Add(comid.ExtMval, &RefValExtensions{})
	if err = measurement.Val.RegisterExtensions(extMap); err != nil {
		log.Fatal("could not register refval extensions")
	}

	refVal.Measurements.Add(measurement)
	myComid.Triples.AddReferenceValue(&refVal)

	err = myComid.Valid()
	if err != nil {
		log.Fatalf("comid validity: %v", err)
	}

	myCorim := profileManifest.GetUnsignedCorim()
	myCorim.SetID("foo")
	myCorim.AddComid(myComid)

	buf, err := myCorim.ToCBOR()
	if err != nil {
		log.Fatalf("could not encode CoRIM: %v", err)
	}

	fmt.Printf("corim: %v", hex.EncodeToString(buf))

	// output:
	// corim: d901f5a30063666f6f0181d901fa58a5a40065656e2d474201a100676578616d706c650281a4006941434d45204c74642e01d8207468747470733a2f2f61636d652e6578616d706c65028101206f3132332046616b652053747265657404a1008182a100a300d8255031fb5abf023e4992aa4e95f9c1503bfa016941434d45204c74642e026e526f616452756e6e657220322e3081a200d8255031fb5abf023e4992aa4e95f9c1503bfa01a10281820644abcdef00037822687474703a2f2f6578616d706c652e636f6d2f6578616d706c652d70726f66696c65
}
