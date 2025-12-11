[CoRIM
specification](https://datatracker.ietf.org/doc/draft-ietf-rats-corim/04)
may be extended by CoRIM Profiles documented in other specifications at well
defined points identified in the base CoRIM spec using CDDL extension sockets.
CoRIM profiles may

1. Introduce new data fields to certain objects.
2. Further constrain allowed values for existing fields.
3. Introduce new type choices for fields whose values can be one of several
   types.

This implementation, likewise, allows dependent code to register extension
types. This is done via three distinct extension mechanisms:

1. Structures that allow extensions embed an `Extensions` object, which allows
   registering a user-provided struct. The user-provided struct can extend
   their containing structures in two ways:

   - the fields of the user-provided struct become additional fields in
     the containing structure.
   - the user-provided struct may define additional constraints on the
     containing structure by defining an appropriate validation method for it.

   This corresponds to the map-extension sockets in the spec.
2. Some fields may have values chosen from a set of pre-determined types (e.g., an
   instance ID may be either a UUID or UEID). These (mostly) correspond to
   type-choice sockets in the CoRIM spec. The set of allowed types for a field
   may be extended by registering a factory function for a new type, using the
   registration function associated with the type choice.
3. A couple of type-choice sockets (`$tag-rel-type-choice`,
   `$corim-role-type-choice` and `$comid-role-type-choice`) define what, in
   effect, are extensible enums. They allow providing additional values, rather
   than types. This implementation provides registration functions for new
   values for those types.

> [!NOTE]
> CoRIM also "imports" CDDL from the CoSWID spec. Some of
> these CoSWID CDDL definitions also feature extension sockets.
> However, as they are defined in a different spec and are implemented
> in the [`veraison/swid`](https://github.com/veraison/swid) library, they
> cannot be extended using the extension feature provided in the CoRIM library.
> The extension support in CoRIM library is applicable ONLY to CoRIM and CoMID
> maps and type choices


## Map Extensions

Map extensions allow extending CoRIM and CoMID maps with additional keys,
effectively defining new fields for the corresponding structures. In the code
base, these can be identified by the embedded `Extensions` struct. Each
extensible type has a corresponding `extensions.Point`. These are:

|    Extended Type   |                             Extension Point(s)                            |                   Parent Structure                   |                     Where to Call RegisterExtensions()                    |
|:-------------------:|:-------------------------------------------------------------------------:|:----------------------------------------------------:|:-------------------------------------------------------------------------:|
| comid.Comid         | comid.ExtComid                                                            | comid.Comid (the top-level CoMID)                    | On a comid.Comid instance (e.g. myComid.RegisterExtensions(extMap))       |
| comid.Entity        | comid.ExtEntity                                                           | comid.Entity                                         | Usually indirect via myComid.RegisterExtensions(...) (the Comid sees it). |
| comid.Triples       | comid.ExtTriples                                                          | comid.Triples                                        | Typically indirect via myComid.RegisterExtensions(...).                   |
| comid.Mval          | comid.ExtReferenceValue, comid.ExtEndorsedValue, comid.ExtMval            | comid.Mval (measurement-value in reference/endorsed) | Usually indirect via myComid.RegisterExtensions(...).                     |
| comid.FlagsMap      | comid.ExtReferenceValueFlags, comid.ExtEndorsedValueFlags, comid.ExtFlags | comid.FlagsMap                                       | Typically indirect via myComid.RegisterExtensions(...).                   |
| corim.UnsignedCorim | corim.ExtUnsignedCorim                                                    | corim.UnsignedCorim (the top-level CoRIM)            | On a corim.UnsignedCorim instance (e.g. myCorim.RegisterExtensions(...))  |
| corim.Entity        | corim.ExtEntity                                                           | corim.Entity                                         | Usually indirect via myCorim.RegisterExtensions(...).                     |
| corim.Signer        | corim.ExtSigner                                                           | corim.Signer                                         | Usually indirect via myCorim.RegisterExtensions(...).                     |

Note that `comid.Mval` and `comid.FlagsMap` are used for both reference values
and endorsed values, which may be extended separately. This is why there are
two extension points associated with each. Additionally, `comid.ExtMval` and
`comid.ExtFlags` also exist when you want to register extensions with a
`comid.Mval` or `comid.Measurment` (and so don't have the context of whether it
will be going into a reference or an endorsed value).

The diagram below shows a visual representation of where these extension points
originate in the `struct` hierarchy, and which CoRIM object are "aware" of
which extension points:

![map extensions](corim-map-extensions.png?)

(note: the diagram shows the logical relationships between structures and does
not meant to accurately reflect the code -- e.g. it omits the container types,
representing them as slices)

To extend the above types, you need to define a struct containing your
extensions for each extension point that you want to extend. You then need to
create an `extensions.Map` that maps the extension points to a pointer to an
instance of the corresponding struct. Finally, you need to pass  the map to the
`RegisterExtensions()` method of an instance of the top-most type that is being
extended. For example, if you want to extend the comid, reference values, and
entities, you would create a struct defining the extensions for each, call
`extensions.NewMap()` to create a new map, call `Add()` on the map to add
mappings from the three extension points to pointers to empty instances of your
structs, and, finally, pass the map to `comid.Comid.RegisterExtensions()`. This
should be done as early as possible, before any marshaling is performed (see
the example below).

These types can be extended in two ways: by adding additional fields, and by
introducing additional constraints over existing fields.

### Adding new fields

To add new fields, simply add them to your extensions struct, ensuring that
the `cbor` and `json` tags on those fields are set correctly. As CoRIM
mandates integer keys, you must use the `keyasint` option for the `cbor` tag.

To access the values of those fields, you can call the extended type instance's
`Extensions.Get()` passing in the name of the field you want to access. The
name can be either the Go struct's field name, the name specified in the `json`
tag, or (a string containing) the integer specified in the `cbor` tag.

`Get()` returns an `interface{}`. There are equivalent `GetInt()`,
`GetString()`, etc. methods that perform the required conversions, and return
the value of the indicated type, along with possible errors. ("Must" versions
of these also exist, e.g. `MustGetString()`, that do not return an error and
instead call `panic()`).

You can also get the pointer to your extension's instance itself by calling
the extended type instance's `GetExtensions()`. This returns an `interface{}`, so
you will need to type assert to be able to access the fields directly.

### Introducing additional constraints

To introduce new constraints, add a method called `Constrain<TYPE>(v *<TYPE>)`
to your extensions struct, where `<TYPE>` is the name of the type being
extended (one of the ones listed above) -- e.g.
`ConstrainComid(v *comid.Comid)` when extending `comid.Comid`. This method, if
it exists, will be invoked inside the extended type instance's `Valid()`
method, passing itself as the parameter.

You do not need to define this method unless you actually want to enforce some
constraints (i.e., if you just want to define additional fields).

### Unknown extensions caching

When unmarshaled data contains entries that do not correspond to fields inside
a registered extensions struct, their values get cached inside the `Extensions`
objects. If the containing object is later re-marshalled, cached values will be
included, so that unknown extensions are not lost.

Cached extension values are not accessible via `Get*()` methods described
above, however they are available as `Extensions.Cached` map.

If an extensions struct is registered after unmarshalling, it will be populated
with any now-recognized cached values, which will then be removed from the
cache.

### Example

The following example illustrates how to implement a map extension by extending
`comid.Entity` with the following features:

1. an optional "email" field
2. additional constraint on the existing "name" field that it contains a
   valid UUID (note: since `NameEntry` is a type choice extensible, this can
   also be done by defining a new value type for `NameEntry` -- see the
   following section).

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/extensions"
)

// the struct containing the extensions
type EntityExtensions struct {
    // a string field extension
	Email string `cbor:"-1,keyasint,omitempty" json:"email,omitempty"`
}

// custom constraints for Entity
func (o EntityExtensions) ConstrainEntity(val *comid.Entity) error {
	_, err := uuid.Parse(val.EntityName.String())
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	return nil
}

var sampleText = `
{
      "name": "31fb5abf-023e-4992-aa4e-95f9c1503bfa",
      "regid": "https://acme.example",
      "email": "info@acme.com",
      "roles": [
        "tagCreator",
        "creator",
        "maintainer"
      ]
}
`

func main() {
	var entity comid.Entity
	extMap := extensions.NewMap().Add(comid.ExtEntity, &EntityExtensions{})
	entity.RegisterExtensions(extMap)

	if err := json.Unmarshal([]byte(sampleText), &entity); err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}

	if err := entity.Valid(); err != nil {
		log.Fatalf("failed to validate: %s", err.Error())
	} else {
		fmt.Println("validation succeeded")
	}

	// obtain the extension field value via a generic getter
	email := entity.Extensions.MustGetString("email")
	fmt.Printf("entity email: %s\n", email)

	// retrive the extensions struct and get value via its field.
	exts := entity.GetExtensions().(*EntityExtensions)
	fmt.Printf("also entity email: %s\n", exts.Email)
}
```

### Profiles

Map extensions may be grouped into profiles. A profile is registered,
associating an `eat.Profile` with an `extensions.Map`. A registered profile can
be obtained by calling `corim.GetProfileManifest()`, which returns a `corim.ProfileManifest`
object which can be used to obtain `corim.UnsignedCorim`, `corim.SignedCorim`
and `comid.Comid` instances that have the associated extensions registered.
`corim.UnmarshalUnsignedCorimFromCBOR()` will automatically look up a
registered profile based on the `eat.Profile` in the provided data (there are
corresponding functions for JSON and signed CoRIM).

```go
// define extensions
type EntityExtensions struct {
	Address *string `cbor:"-1,keyasint,omitempty" json:"address,omitempty"`
}

type RefValExtensions struct {
	Timestamp *int `cbor:"-1,keyasint,omitempty" json:"timestamp,omitempty"`
}


// register profile
func init() {
	profileID, err := eat.NewProfile("http://example.com/example-profile")
	if err != nil {
		panic(err) // will not error, as the hard-coded string above is valid
	}

	extMap := extensions.NewMap().
		Add(comid.ExtEntity, &EntityExtensions{}).
		Add(comid.ExtReferenceValue, &RefValExtensions{})

	if err := RegisterProfile(profileID, extMap); err != nil {
		// will not error, assuming our profile ID is unique, and we've
		// correctly set up the extensions Map above
		panic(err)
	}
}


// use the profile
func main() {
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

    // ...
}

```

Please see [example_profile_test.go](../corim/example_profile_test.go) for the complete
example of creating and using CoRIM profiles.

> [!NOTE]
> Currently, only map extensions can be associated with profiles. Type choice
> and enum value extensions (described below) can only be registered globally.


## Type Choice Extensions

Type Choice extensions allow specifying alternative types for existing CoRIM
fields by defining a type that implements an appropriate interface and
registering it with a CBOR tag.

A type choice struct contains a single field, `Value`, that contains the actual
object represented by the type choice. The `Value` implements an interface
that is specific to the type choice and is derived from `ITypeChoiceValue`:

```go
type ITypeChoiceValue interface {
	// String returns the string representation of the ITypeChoiceValue.
	String() string
	// Valid returns an error if validation of the ITypeChoiceValue fails,
	// or nil if it succeeds.
	Valid() error
	// Type returns the type name of this ITypeChoiceValue implementation.
	Type() string
}
```

The following is the full list of type choice structs:

- `comid.ClassID`
- `comid.CryptoKey`
- `comid.EntityName`
- `comid.Group`
- `comid.Instance`
- `comid.Mkey`
- `comid.SVN`
- `corim.EntityName`

To provide a new value type, the following is required:

1. Define a type that implements the value interface for the type choice you
   want to extend. This interface is called `I<NAME>Value`,  where `<NAME>` is
   the name of the type choice type(e.g. `IClassIDValue`).  These interfaces
   always embed `ITypeChoiceValue` and possibly define additional methods.
2. Create a factory function for your type, with the signature `func (any)
   (*<NAME>, error)`, where `<NAME>` is the name of the type choice type that
   will contain your value. (Note that the function must return a pointer to
   the container type choice struct, _not_ to the value type you define.) This
   function should create an instance of your value type from the provided
   input and return a new type choice struct instance containing it. The range
   of valid inputs is up to you, however it _must_ handle `null`, returning the
   [zero-value](https://go.dev/ref/spec#The_zero_value) for your type in that
   case.
3. Register your factory function with the CBOR tag for your new type by
   passing it to the registration function corresponding to the type choice
   struct. It will have the name `Register<NAME>Type`, where `<NAME>` is the
   name of the type choice struct that will contain your value (e.g.
   `RegisterClassIDType`).

### Example

The following example illustrates how to add a new type choice value
implementation by extending the `CryptoKey` type to support DER values.

```go
package main

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/veraison/corim/comid"
)

// the CBOR tag to be used for the new type
var DerKeyTag = uint64(9999)

// new implementation of ICryptoKeyValue type
type TaggedDerKey []byte

// The factory function for the new type
func NewTaggedDerKey(k any) (*comid.CryptoKey, error) {
	var b []byte
	var err error

	if k == nil {
		k = *new([]byte)
	}

	switch t := k.(type) {
	case []byte:
		b = t
	case string:
		b, err = base64.StdEncoding.DecodeString(t)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("value must be a []byte; found %T", k)
	}

	key := TaggedDerKey(b)

	return &comid.CryptoKey{Value: key}, nil
}

func (o TaggedDerKey) String() string {
	return base64.StdEncoding.EncodeToString(o)
}

func (o TaggedDerKey) Valid() error {
	_, err := o.PublicKey()
	return err
}

func (o TaggedDerKey) Type() string {
	return "pkix-der-key"
}

func (o TaggedDerKey) PublicKey() (crypto.PublicKey, error) {
	if len(o) == 0 {
		return nil, errors.New("key value not set")
	}

	key, err := x509.ParsePKIXPublicKey(o)
	if err != nil {
		return nil, fmt.Errorf("unable to parse public key: %w", err)
	}

	return key, nil
}

var testKeyJSON = `
{
	"type": "pkix-der-key",
	"value": "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEW1BvqF+/ry8BWa7ZEMU1xYYHEQ8BlLT4MFHOaO+ICTtIvrEeEpr/sfTAP66H2hCHdb5HEXKtRKod6QLcOLPA1Q=="
}
`

func main() {
    // register the factory function under the CBOR tag.
	if err := comid.RegisterCryptoKeyType(DerKeyTag, NewTaggedDerKey); err != nil {
		log.Fatal(err)
	}

	var key comid.CryptoKey

	if err := json.Unmarshal([]byte(testKeyJSON), &key); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded DER key: %x\n", key)
}
```


## Enum extensions

The following enum types may be extended with additional values:

- `comid.Rel`
- `comid.Role`
- `corim.Role`

This can be done by calling `RegisterRel` or `RegisterRole`, as appropriate,
and providing a new `uint64` value and corresponding `string` name.

### Example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/veraison/corim/comid"
)

var sampleText = `
{
      "name": "Acme Ltd.",
      "regid": "https://acme.example",
      "roles": [
        "tagCreator",
        "owner"
      ]
}
`

func main() {
    // associate role value 4 with the name "owner"
	comid.RegisterRole(4, "owner")

	var entity comid.Entity

	if err := json.Unmarshal([]byte(sampleText), &entity); err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}

	if err := entity.Valid(); err != nil {
		log.Fatalf("failed to validate: %s", err.Error())
	} else {
		fmt.Println("validation succeeded")
	}

	fmt.Println("roles:")
	for _, role := range entity.Roles {
		fmt.Printf("\t%s\n", role.String())
	}
}
```
