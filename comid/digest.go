// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Named Information Hash Algorithm Registry
// https://www.iana.org/assignments/named-information/named-information.xhtml#hash-alg
const (
	Sha256 int = (iota + 1)
	Sha256_128
	Sha256_120
	Sha256_96
	Sha256_64
	Sha256_32
	Sha384
	Sha512
	Sha3_224
	Sha3_256
	Sha3_384
	Sha3_512
)

var (
	algToValueLen = map[int]int{
		Sha256:     32,
		Sha256_128: 16,
		Sha256_120: 15,
		Sha256_96:  12,
		Sha256_64:  8,
		Sha256_32:  4,
		Sha384:     48,
		Sha512:     64,
		Sha3_224:   28,
		Sha3_256:   32,
		Sha3_384:   48,
		Sha3_512:   64,
	}

	algToString = map[int]string{
		Sha256:     "sha-256",
		Sha256_128: "sha-256-128",
		Sha256_120: "sha-256-120",
		Sha256_96:  "sha-256-96",
		Sha256_64:  "sha-256-64",
		Sha256_32:  "sha-256-32",
		Sha384:     "sha-384",
		Sha512:     "sha-512",
		Sha3_224:   "sha3-224",
		Sha3_256:   "sha3-256",
		Sha3_384:   "sha3-384",
		Sha3_512:   "sha3-512",
	}

	stringToAlg = map[string]int{
		"sha-256":     Sha256,
		"sha-256-128": Sha256_128,
		"sha-256-120": Sha256_120,
		"sha-256-96":  Sha256_96,
		"sha-256-64":  Sha256_64,
		"sha-256-32":  Sha256_32,
		"sha-384":     Sha384,
		"sha-512":     Sha512,
		"sha3-224":    Sha3_224,
		"sha3-256":    Sha3_256,
		"sha3-384":    Sha3_384,
		"sha3-512":    Sha3_512,
	}
)

func IntDigestAlgorithm(val int) DigestAlgorithm {
	return DigestAlgorithm{val}
}

func StringDigestAlgorithm(val string) DigestAlgorithm {
	return DigestAlgorithm{val}
}

func DigestAlgorithmFromString(val string) DigestAlgorithm {
	i, ok := stringToAlg[val]
	if ok {
		return DigestAlgorithm{i}
	}

	i, err := strconv.Atoi(val)
	if err == nil {
		return DigestAlgorithm{i}
	}

	return DigestAlgorithm{val}
}

func DigestAlgorithmFromAny(val any) (DigestAlgorithm, error) {
	switch t := val.(type) {
	case int:
		return IntDigestAlgorithm(t), nil
	case int64:
		return IntDigestAlgorithm(int(t)), nil
	case float64:
		return IntDigestAlgorithm(int(t)), nil
	case string:
		return DigestAlgorithmFromString(t), nil
	default:
		return DigestAlgorithm{0}, fmt.Errorf("invalid digest algorithm: %v(%T)", t, t)
	}
}

type DigestAlgorithm struct {
	val any
}

func (o DigestAlgorithm) IsString() bool {
	_, ok := o.val.(string)
	return ok
}

func (o DigestAlgorithm) IsInt() bool {
	_, ok := o.val.(int)
	return ok
}

func (o DigestAlgorithm) String() string {
	switch t := o.val.(type) {
	case string:
		return t
	case int:
		text, ok := algToString[t]
		if !ok {
			text = fmt.Sprintf("%d", t)
		}

		return text
	default:
		return ""
	}
}

func (o DigestAlgorithm) Int() int {
	switch t := o.val.(type) {
	case int:
		return t
	case string:
		val := stringToAlg[t]
		return val
	default:
		return 0
	}
}

func (o DigestAlgorithm) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

func (o *DigestAlgorithm) UnmarshalCBOR(data []byte) error {
	if len(data) == 0 {
		return errors.New("buffer too short")
	}

	majorType := (data[0] & 0xe0) >> 5
	switch majorType {
	case 0, 1:
		var val int
		if err := dm.Unmarshal(data, &val); err != nil {
			return err
		}

		o.val = val
		return nil
	case 3:
		var val string
		if err := dm.Unmarshal(data, &val); err != nil {
			return err
		}

		o.val = val
		return nil
	default:
		return fmt.Errorf("unexpected CBOR major type for DigestAlgID: %d", majorType)
	}
}

func (o DigestAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.val)
}

func (o *DigestAlgorithm) UnmarshalJSON(data []byte) error {
	var val any
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	switch t := val.(type) {
	case float64:
		*o = IntDigestAlgorithm(int(t))
	case string:
		*o = StringDigestAlgorithm(t)
	default:
		return fmt.Errorf("unexpected algorithm value: %v(%T)", t, t)
	}

	return nil
}

func DigestFromString(val string) (Digest, error) {
	parts := strings.Split(val, ";")
	if len(parts) != 2 {
		return Digest{}, fmt.Errorf("expected exactly two ;-separated parts, got %q", val)
	}

	value, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Digest{}, fmt.Errorf("val: %w", err)
	}

	return NewDigestStringAlg(parts[0], value), nil
}

type Digest struct {
	_         struct{} `cbor:",toarray"`
	Algorithm DigestAlgorithm
	Value     []byte
}

func NewDigestIntAlg(alg int, value []byte) Digest {
	return NewDigest(IntDigestAlgorithm(alg), value)
}

func NewDigestStringAlg(alg string, value []byte) Digest {
	return NewDigest(DigestAlgorithmFromString(alg), value)
}

func NewDigest(alg DigestAlgorithm, value []byte) Digest {
	return Digest{Algorithm: alg, Value: value}
}

func (o Digest) String() string {
	return o.Algorithm.String() + ";" + base64.RawURLEncoding.EncodeToString(o.Value)
}

func (o Digest) Valid() error {
	if len(o.Value) == 0 {
		return errors.New("zero length value")
	}

	if o.Algorithm.Int() == 0 && !o.Algorithm.IsString() {
		return errors.New("zero algorithm")
	}

	wantLen, ok := algToValueLen[o.Algorithm.Int()]
	if ok {
		gotLen := len(o.Value)
		if wantLen != gotLen {
			return fmt.Errorf(
				"length mismatch for hash algorithm %s: want %d bytes, got %d",
				o.Algorithm.String(), wantLen, gotLen,
			)
		}
	}

	return nil
}

func (o Digest) PublicKey() (crypto.PublicKey, error) {
	return nil, errors.New("cannot get PublicKey from a digest")
}

func (o Digest) Bytes() []byte {
	ret, _ := em.Marshal(o)
	return ret
}

func (o Digest) MarshalJSON() ([]byte, error) {
	toMarshal := [2]any{
		o.Algorithm,
		base64.RawURLEncoding.EncodeToString(o.Value),
	}

	return json.Marshal(toMarshal)
}

func (o *Digest) UnmarshalJSON(data []byte) error {
	var decoded []any
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	if len(decoded) != 2 {
		return fmt.Errorf("expected array with two elements, got %v", decoded)
	}

	alg, err := DigestAlgorithmFromAny(decoded[0])
	if err != nil {
		return fmt.Errorf("alg: %w", err)
	}

	switch t := decoded[1].(type) {
	case string:
		bytes, err := base64.RawURLEncoding.DecodeString(t)
		if err != nil {
			return fmt.Errorf("val: %w", err)
		}

		o.Algorithm = alg
		o.Value = bytes

		return nil
	default:
		return fmt.Errorf("invalid val: %v(%T)", t, t)
	}

}
