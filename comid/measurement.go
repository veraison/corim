// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
	"github.com/veraison/eat"
	"github.com/veraison/swid"
)

const MaxUint64 = ^uint64(0)

// Mkey stores a $measured-element-type-choice.
// The supported types are UUID, PSA refval-id, CCA platform-config-id and unsigned integer
// TO DO Add tagged OID: see https://github.com/veraison/corim/issues/35
type Mkey struct {
	Value IMKeyValue
}

// NewMkey creates a new Mkey of the specfied type using the provided value.
func NewMkey(val any, typ string) (*Mkey, error) {
	factory, ok := mkeyValueRegister[typ]
	if !ok {
		return nil, fmt.Errorf("unexpected measurement key type: %q", typ)
	}

	return factory(val)
}

// MustNewMkey is like NewMkey, execept it does not return an error, assuming
// that the provided value is valid. It panics if that is not the case.
func MustNewMkey(val any, typ string) *Mkey {
	ret, err := NewMkey(val, typ)
	if err != nil {
		panic(err)
	}

	return ret
}

// IsSet returns true if the value of the Mkey is set.
func (o Mkey) IsSet() bool {
	return o.Value != nil
}

// Type returns the type of Mkey
func (o Mkey) Type() string {
	return o.Value.Type()
}

// Valid returns nil if the Mkey is valid or an error describing the problem,
// if it is not.
func (o Mkey) Valid() error {
	if o.Value == nil {
		return errors.New("Mkey value not set")
	}

	if err := o.Value.Valid(); err != nil {
		return fmt.Errorf("invalid %s: %w", o.Value.Type(), err)
	}

	return nil
}

func (o Mkey) GetPSARefValID() (PSARefValID, error) {
	if !o.IsSet() {
		return PSARefValID{}, errors.New("MKey is not set")
	}
	switch t := o.Value.(type) {
	case *TaggedPSARefValID:
		return PSARefValID(*t), nil
	case TaggedPSARefValID:
		return PSARefValID(t), nil
	default:
		return PSARefValID{}, fmt.Errorf("measurement-key type is: %T", t)
	}
}

func (o Mkey) GetCCAPlatformConfigID() (CCAPlatformConfigID, error) {
	if !o.IsSet() {
		return "", errors.New("MKey is not set")
	}
	switch t := o.Value.(type) {
	case *TaggedCCAPlatformConfigID:
		return CCAPlatformConfigID(*t), nil
	case TaggedCCAPlatformConfigID:
		return CCAPlatformConfigID(t), nil
	default:
		return "", fmt.Errorf("measurement-key type is: %T", t)
	}
}

func (o Mkey) GetKeyUint() (uint64, error) {
	switch t := o.Value.(type) {
	case UintMkey:
		return uint64(t), nil
	case *UintMkey:
		return uint64(*t), nil
	default:
		return MaxUint64, fmt.Errorf("measurement-key type is: %T", t)
	}
}

// UnmarshalJSON deserializes the supplied JSON object into the target MKey
// The key object must have the following shape:
//
//	{
//	  "type": "<MKEY_TYPE>",
//	  "value": <MKEY_JSON_VALUE>
//	}
//
// where <MKEY_TYPE> must be one of the known IMKeyValue implementation
// type names (available in the base implementation: "uuid", "oid",
// "psa.impl-id"), and <MKEY_JSON_VALUE> is the class id value serialized to
// JSON. The exact serialization is <CLASS_ID_TYPE> depenent. For the base
// implementation types it is
//
//	oid: dot-seprated integers, e.g. "1.2.3.4"
//	uuid: standard UUID string representation, e.g. "550e8400-e29b-41d4-a716-446655440000"
//	psa.refval-id: JSON representation of the PSA refval-id
func (o *Mkey) UnmarshalJSON(data []byte) error {
	var tnv encoding.TypeAndValue

	if err := json.Unmarshal(data, &tnv); err != nil {
		return err
	}

	decoded, err := NewMkey(nil, tnv.Type)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(tnv.Value, decoded.Value); err != nil {
		return fmt.Errorf("invalid %s: %w", tnv.Type, err)
	}

	if err := decoded.Value.Valid(); err != nil {
		return fmt.Errorf("invalid %s: %w", tnv.Type, err)
	}

	o.Value = decoded.Value

	return nil
}

// MarshalJSON serializes the target Mkey into the type'n'value JSON object
func (o Mkey) MarshalJSON() ([]byte, error) {
	valueBytes, err := json.Marshal(o.Value)
	if err != nil {
		return nil, err
	}

	value := encoding.TypeAndValue{
		Type:  o.Value.Type(),
		Value: valueBytes,
	}

	return json.Marshal(value)
}

// MarshalCBOR serializes the taret mkey into  CBOR-encoded bytes.
func (o Mkey) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.Value)
}

// UnmarshalCBOR deserializes the Mkey from the provided CBOR bytes.
func (o *Mkey) UnmarshalCBOR(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty input")
	}

	majorType := (data[0] & 0xe0) >> 5
	if majorType == 6 { // tag
		return dm.Unmarshal(data, &o.Value)
	}

	// untagged value must be a uint

	var val UintMkey
	if err := dm.Unmarshal(data, &val); err != nil {
		return err
	}

	o.Value = &val
	return nil
}

// IMKeyValue is the interface implemented by all Mkey value implementations.
type IMKeyValue interface {
	extensions.ITypeChoiceValue
}

const UintType = "uint"

type UintMkey uint64

func NewUintMkey(val any) (*UintMkey, error) {
	var ret UintMkey

	if val == nil {
		return &ret, nil
	}

	switch t := val.(type) {
	case UintMkey:
		ret = t
	case *UintMkey:
		ret = *t
	case string:
		u, err := strconv.ParseUint(t, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = UintMkey(u)
	case uint64:
		ret = UintMkey(t)
	case uint:
		ret = UintMkey(t)
	default:
		return nil, fmt.Errorf("unexpected type for UintMkey: %T", t)
	}

	return &ret, nil
}

func (o UintMkey) Valid() error {
	return nil
}

func (o UintMkey) String() string {
	return fmt.Sprint(uint64(o))
}

func (o UintMkey) Type() string {
	return UintType
}

func (o *UintMkey) UnmarshalJSON(data []byte) error {
	var tmp uint64

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	*o = UintMkey(tmp)

	return nil
}

func NewMkeyOID(val any) (*Mkey, error) {
	ret, err := NewTaggedOID(val)
	if err != nil {
		return nil, err
	}

	return &Mkey{ret}, nil
}

func NewMkeyUUID(val any) (*Mkey, error) {
	ret, err := NewTaggedUUID(val)
	if err != nil {
		return nil, err
	}

	return &Mkey{ret}, nil
}

func NewMkeyUint(val any) (*Mkey, error) {
	ret, err := NewUintMkey(val)
	if err != nil {
		return nil, err
	}

	return &Mkey{ret}, nil
}

func NewMkeyPSARefvalID(val any) (*Mkey, error) {
	ret, err := NewTaggedPSARefValID(val)
	if err != nil {
		return nil, err
	}

	return &Mkey{ret}, nil
}

func NewMkeyCCAPlatformConfigID(val any) (*Mkey, error) {
	ret, err := NewTaggedCCAPlatformConfigID(val)
	if err != nil {
		return nil, err
	}

	return &Mkey{ret}, nil
}

// IMkeyFactory defines the signature for the factory functions that may be
// registred using RegisterMkeyType to provide a new implementation of the
// corresponding type choice. The factory function should create a new *Mkey
// with the underlying value created based on the provided input. The range of
// valid inputs is up to the specific type choice implementation, however it
// _must_ accept nil as one of the inputs, and return the Zero value for
// implemented type.
// See also https://go.dev/ref/spec#The_zero_value
type IMkeyFactory = func(val any) (*Mkey, error)

var mkeyValueRegister = map[string]IMkeyFactory{
	OIDType:                 NewMkeyOID,
	UUIDType:                NewMkeyUUID,
	UintType:                NewMkeyUint,
	PSARefValIDType:         NewMkeyPSARefvalID,
	CCAPlatformConfigIDType: NewMkeyCCAPlatformConfigID,
}

// RegisterMkeyType registers a new IMKeyValue implementation
// (created by the provided IMKeyFactory) under the specified CBOR tag.
func RegisterMkeyType(tag uint64, factory IMkeyFactory) error {

	nilVal, err := factory(nil)
	if err != nil {
		return err
	}

	typ := nilVal.Value.Type()
	if _, exists := mkeyValueRegister[typ]; exists {
		return fmt.Errorf("measurement key type with name %q already exists", typ)
	}

	if err := registerCOMIDTag(tag, nilVal.Value); err != nil {
		return err
	}

	mkeyValueRegister[typ] = factory

	return nil
}

// Mval stores a measurement-values-map with JSON and CBOR serializations.
type Mval struct {
	Ver                *Version            `cbor:"0,keyasint,omitempty" json:"version,omitempty"`
	SVN                *SVN                `cbor:"1,keyasint,omitempty" json:"svn,omitempty"`
	Digests            *Digests            `cbor:"2,keyasint,omitempty" json:"digests,omitempty"`
	Flags              *FlagsMap           `cbor:"3,keyasint,omitempty" json:"flags,omitempty"`
	RawValue           *RawValue           `cbor:"4,keyasint,omitempty" json:"raw-value,omitempty"`
	RawValueMask       *[]byte             `cbor:"5,keyasint,omitempty" json:"raw-value-mask,omitempty"`
	MACAddr            *MACaddr            `cbor:"6,keyasint,omitempty" json:"mac-addr,omitempty"`
	IPAddr             *net.IP             `cbor:"7,keyasint,omitempty" json:"ip-addr,omitempty"`
	SerialNumber       *string             `cbor:"8,keyasint,omitempty" json:"serial-number,omitempty"`
	UEID               *eat.UEID           `cbor:"9,keyasint,omitempty" json:"ueid,omitempty"`
	UUID               *UUID               `cbor:"10,keyasint,omitempty" json:"uuid,omitempty"`
	Name               *string             `cbor:"11,keyasint,omitempty" json:"name,omitempty"`
	IntegrityRegisters *IntegrityRegisters `cbor:"14,keyasint,omitempty" json:"integrity-registers,omitempty"`
	Extensions
}

// RegisterExtensions registers a struct as a collections of extensions
func (o *Mval) RegisterExtensions(exts extensions.Map) error {
	for p, v := range exts {
		switch p {
		case ExtMval:
			o.Extensions.Register(v)
		case ExtFlags:
			if o.Flags == nil {
				o.Flags = new(FlagsMap)
			}

			o.Flags.Extensions.Register(v)
		default:
			return fmt.Errorf("%w: %q", extensions.ErrUnexpectedPoint, p)
		}
	}

	return nil
}

// GetExtensions returns pervisouosly registered extension
func (o *Mval) GetExtensions() extensions.IMapValue {
	return o.Extensions.IMapValue
}

// UnmarshalCBOR deserializes from CBOR
func (o *Mval) UnmarshalCBOR(data []byte) error {
	return encoding.PopulateStructFromCBOR(dm, data, o)
}

// MarshalCBOR serializes to CBOR
func (o Mval) MarshalCBOR() ([]byte, error) {
	// If extensions have been registered, the collection will exist, but
	// might be empty. If that is the case, set it to nil to avoid
	// marshaling an empty list (and let the marshaller omit the claim
	// instead). Note that since the receiver was passed by value, we do not
	// need to worry about saving the field's value before setting it to
	// nil.
	if o.Flags != nil && o.Flags.IsEmpty() {
		o.Flags = nil
	}

	return encoding.SerializeStructToCBOR(em, o)
}

// UnmarshalJSON deserializes from JSON
func (o *Mval) UnmarshalJSON(data []byte) error {
	return encoding.PopulateStructFromJSON(data, o)
}

// MarshalJSON serializes to JSON
func (o Mval) MarshalJSON() ([]byte, error) {
	// If extensions have been registered, the collection will exist, but
	// might be empty. If that is the case, set it to nil to avoid
	// marshaling an empty list (and let the marshaller omit the claim
	// instead). Note that since the receiver was passed by value, we do not
	// need to worry about saving the field's value before setting it to
	// nil.
	if o.Flags != nil && o.Flags.IsEmpty() {
		o.Flags = nil
	}

	return encoding.SerializeStructToJSON(o)
}

func (o Mval) Valid() error {
	if o.Ver == nil &&
		o.SVN == nil &&
		o.Digests == nil &&
		o.Flags == nil &&
		o.RawValue == nil &&
		o.RawValueMask == nil &&
		o.MACAddr == nil &&
		o.IPAddr == nil &&
		o.SerialNumber == nil &&
		o.UEID == nil &&
		o.UUID == nil &&
		o.Name == nil &&
		o.IntegrityRegisters == nil {
		return fmt.Errorf("no measurement value set")
	}

	if o.Ver != nil {
		if err := o.Ver.Valid(); err != nil {
			return err
		}
	}

	if o.Digests != nil {
		if err := o.Digests.Valid(); err != nil {
			return err
		}
	}

	if o.Flags != nil {
		if err := o.Flags.Valid(); err != nil {
			return err
		}
	}

	// raw value and mask have no specific semantics

	// TODO(tho) MAC addr & friends (see https://github.com/veraison/corim/issues/18)

	return o.Extensions.validMval(&o)
}

// Version stores a version-map with JSON and CBOR serializations.
type Version struct {
	Version string             `cbor:"0,keyasint" json:"value"`
	Scheme  swid.VersionScheme `cbor:"1,keyasint" json:"scheme"`
}

func NewVersion() *Version {
	return &Version{}
}

func (o *Version) SetVersion(v string) *Version {
	if o != nil {
		o.Version = v
	}
	return o
}

func (o *Version) SetScheme(v int64) *Version {
	if o != nil {
		if o.Scheme.SetCode(v) != nil {
			return nil
		}
	}
	return o
}

func (o Version) Valid() error {
	if o.Version == "" {
		return fmt.Errorf("empty version")
	}
	return nil
}

// Measurement stores a measurement-map with CBOR and JSON serializations.
type Measurement struct {
	Key          *Mkey      `cbor:"0,keyasint,omitempty" json:"key,omitempty"`
	Val          Mval       `cbor:"1,keyasint" json:"value"`
	AuthorizedBy *CryptoKey `cbor:"2,keyasint,omitempty" json:"authorized-by,omitempty"`
}

func NewMeasurement(val any, typ string) (*Measurement, error) {
	keyFactory, ok := mkeyValueRegister[typ]
	if !ok {
		return nil, fmt.Errorf("unknown Mkey type: %s", typ)
	}

	key, err := keyFactory(val)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %w", err)
	}

	if err = key.Valid(); err != nil {
		return nil, fmt.Errorf("invalid key: %w", err)
	}

	var ret Measurement
	ret.Key = key

	return &ret, nil
}

func MustNewMeasurement(val any, typ string) *Measurement {
	ret, err := NewMeasurement(val, typ)

	if err != nil {
		panic(err)
	}

	return ret
}

// NewPSAMeasurement instantiates a new measurement-map with the key set to the
// supplied PSA refval-id
func NewPSAMeasurement(key any) (*Measurement, error) {
	return NewMeasurement(key, PSARefValIDType)
}

func MustNewPSAMeasurement(key any) *Measurement {
	ret, err := NewPSAMeasurement(key)

	if err != nil {
		panic(err)
	}

	return ret
}

// NewCCAPlatCfgMeasurement instantiates a new measurement-map with the key set to the
// supplied CCA platform-config-id
func NewCCAPlatCfgMeasurement(key any) (*Measurement, error) {
	return NewMeasurement(key, CCAPlatformConfigIDType)
}

func MustNewCCAPlatCfgMeasurement(key any) *Measurement {
	ret, err := NewCCAPlatCfgMeasurement(key)

	if err != nil {
		panic(err)
	}

	return ret
}

// NewUUIDMeasurement instantiates a new measurement-map with the key set to the
// supplied UUID
func NewUUIDMeasurement(key any) (*Measurement, error) {
	return NewMeasurement(key, UUIDType)
}

func MustNewUUIDMeasurement(key any) *Measurement {
	ret, err := NewUUIDMeasurement(key)

	if err != nil {
		panic(err)
	}

	return ret
}

// NewUintMeasurement instantiates a new measurement-map with the key set to the
// supplied Uint
func NewUintMeasurement(key any) (*Measurement, error) {
	return NewMeasurement(key, UintType)
}

func MustNewUintMeasurement(key any) *Measurement {
	ret, err := NewUintMeasurement(key)

	if err != nil {
		panic(err)
	}

	return ret
}

// NewOIDMeasurement instantiates a new measurement-map with the key set to the
// supplied OID
func NewOIDMeasurement(key any) (*Measurement, error) {
	return NewMeasurement(key, OIDType)
}

func (o *Measurement) RegisterExtensions(exts extensions.Map) error {
	return o.Val.RegisterExtensions(exts)
}

func (o Measurement) GetExtensions() extensions.IMapValue {
	return o.Val.GetExtensions()
}

func (o *Measurement) SetVersion(ver string, scheme int64) *Measurement {
	if o != nil {
		v := NewVersion().SetVersion(ver).SetScheme(scheme)
		if v == nil {
			return nil
		}

		o.Val.Ver = v
	}
	return o
}

// SetRawValueBytes sets the supplied raw-value and its mask in the
// measurement-values-map of the target measurement
func (o *Measurement) SetRawValueBytes(rawValue, rawValueMask []byte) *Measurement {
	if o != nil {
		o.Val.RawValue = NewRawValue().SetBytes(rawValue)
		if len(rawValueMask) != 0 {
			o.Val.RawValueMask = &rawValueMask
		}
	}
	return o
}

// SetSVN sets the supplied svn in the measurement-values-map of the target
// measurement
func (o *Measurement) SetSVN(svn uint64) *Measurement {
	o.Val.SVN = MustNewTaggedSVN(svn)
	return o
}

// SetMinSVN sets the supplied min-svn in the measurement-values-map of the
// target measurement
func (o *Measurement) SetMinSVN(svn uint64) *Measurement {
	o.Val.SVN = MustNewTaggedMinSVN(svn)
	return o
}

// AddDigest add the supplied digest - comprising the digest itself together
// with the hash algorithm used to obtain it - to the measurement-values-map of
// the target measurement
func (o *Measurement) AddDigest(algID uint64, digest []byte) *Measurement {
	if o != nil {
		ds := o.Val.Digests
		if ds == nil {
			ds = NewDigests()
		}
		if ds.AddDigest(algID, digest) == nil {
			return nil
		}
		o.Val.Digests = ds
	}

	return o
}

// SetFlagsTrue sets the supplied operational flags to true in the
// measurement-values-map of the target measurement
func (o *Measurement) SetFlagsTrue(flags ...Flag) *Measurement {
	if o != nil {
		if o.Val.Flags == nil {
			o.Val.Flags = NewFlagsMap()
		}
		o.Val.Flags.SetTrue(flags...)
	}

	return o
}

// SetFlagsFalse sets the supplied operational flags to true in the
// measurement-values-map of the target measurement
func (o *Measurement) SetFlagsFalse(flags ...Flag) *Measurement {
	if o != nil {
		if o.Val.Flags == nil {
			o.Val.Flags = NewFlagsMap()
		}
		o.Val.Flags.SetFalse(flags...)
	}

	return o
}

// ClearFlags clears the supplied operational flags in the
// measurement-values-map of the target measurement
func (o *Measurement) ClearFlags(flags ...Flag) *Measurement {
	if o != nil {
		if o.Val.Flags == nil {
			return o
		}

		o.Val.Flags.Clear(flags...)

		if !o.Val.Flags.AnySet() {
			o.Val.Flags = nil
		}
	}

	return o
}

// SetIPaddr sets the supplied IP (v4 or v6) address in the
// measurement-values-map of the target measurement
func (o *Measurement) SetIPaddr(a net.IP) *Measurement {
	if o != nil {
		o.Val.IPAddr = &a
	}
	return o
}

// SetMACaddr sets the supplied MAC address in the measurement-values-map of the
// target measurement
func (o *Measurement) SetMACaddr(a MACaddr) *Measurement {
	if o != nil {
		o.Val.MACAddr = &a
	}
	return o
}

// SetSerialNumber sets the supplied serial number in the measurement-values-map
// of the target measurement
func (o *Measurement) SetSerialNumber(sn string) *Measurement {
	if o != nil {
		o.Val.SerialNumber = &sn
	}
	return o
}

// SetUEID sets the supplied ueid in the measurement-values-map
// of the target measurement
func (o *Measurement) SetUEID(ueid eat.UEID) *Measurement {
	if o != nil {
		if ueid.Validate() != nil {
			return nil
		}
		o.Val.UEID = &ueid
	}
	return o
}

// SetUUID sets the supplied uuid in the measurement-values-map
// of the target measurement
func (o *Measurement) SetUUID(u UUID) *Measurement {
	if o != nil {
		if u.Valid() != nil {
			return nil
		}
		o.Val.UUID = &u
	}
	return o
}

// SetName sets the supplied name string in the measurement-values-map of the
// target measurement
func (o *Measurement) SetName(name string) *Measurement {
	if o != nil {
		o.Val.Name = &name
	}
	return o
}

func (o Measurement) Valid() error {
	if o.Key != nil && o.Key.IsSet() {
		if err := o.Key.Valid(); err != nil {
			return err
		}
	}

	return o.Val.Valid()
}

// Measurements is a container for Measurement instances and their extensions.
// It is a thin wrapper around extensions.Collection.
type Measurements extensions.Collection[Measurement, *Measurement]

func NewMeasurements() *Measurements {
	return (*Measurements)(extensions.NewCollection[Measurement]())
}

func (o *Measurements) RegisterExtensions(exts extensions.Map) error {
	return (*extensions.Collection[Measurement, *Measurement])(o).RegisterExtensions(exts)
}

func (o *Measurements) GetExtensions() extensions.IMapValue {
	return (*extensions.Collection[Measurement, *Measurement])(o).GetExtensions()
}

func (o *Measurements) Valid() error {
	return (*extensions.Collection[Measurement, *Measurement])(o).Valid()
}

func (o *Measurements) IsEmpty() bool {
	return (*extensions.Collection[Measurement, *Measurement])(o).IsEmpty()
}

func (o *Measurements) Add(val *Measurement) *Measurements {
	ret := (*extensions.Collection[Measurement, *Measurement])(o).Add(val)
	return (*Measurements)(ret)
}

func (o Measurements) MarshalCBOR() ([]byte, error) {
	return (extensions.Collection[Measurement, *Measurement])(o).MarshalCBOR()
}

func (o *Measurements) UnmarshalCBOR(data []byte) error {
	return (*extensions.Collection[Measurement, *Measurement])(o).UnmarshalCBOR(data)
}
func (o Measurements) MarshalJSON() ([]byte, error) {
	return (extensions.Collection[Measurement, *Measurement])(o).MarshalJSON()
}

func (o *Measurements) UnmarshalJSON(data []byte) error {
	return (*extensions.Collection[Measurement, *Measurement])(o).UnmarshalJSON(data)
}
