package corim

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"reflect"

	"github.com/lestrrat-go/jwx/v2/jwk"
	cose "github.com/veraison/go-cose"
)

func getAlgAndKeyFromJWK(j string) (cose.Algorithm, crypto.Signer, error) {
	const noAlg = cose.Algorithm(-65537)
	var (
		err error
		ks  jwk.Set
		k   jwk.Key
		ok  bool
		crv elliptic.Curve
		alg cose.Algorithm
	)

	ks, err = jwk.ParseString(j)
	if err != nil {
		return noAlg, nil, err
	}

	k, ok = ks.Key(0)
	if !ok {
		return noAlg, nil, errors.New("key extraction failed")
	}

	var key crypto.Signer

	err = k.Raw(&key)
	if err != nil {
		return noAlg, nil, err
	}

	switch v := key.(type) {
	case *ecdsa.PrivateKey:
		crv = v.Curve
		if crv == elliptic.P256() {
			alg = cose.AlgorithmES256
			break
		}
		return noAlg, nil, fmt.Errorf("unknown elliptic curve %v", crv)
	default:
		return noAlg, nil, fmt.Errorf("unknown private key type %v", reflect.TypeOf(key))
	}

	return alg, key, nil
}

func NewSignerFromJWK(j []byte) (cose.Signer, error) {
	alg, key, err := getAlgAndKeyFromJWK(string(j))
	if err != nil {
		return nil, err
	}

	return cose.NewSigner(alg, key)
}

func NewPublicKeyFromJWK(j []byte) (crypto.PublicKey, error) {
	_, key, err := getAlgAndKeyFromJWK(string(j))
	if err != nil {
		return nil, err
	}

	return key.Public(), nil
}
