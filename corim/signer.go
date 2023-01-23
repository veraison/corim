package corim

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"reflect"

	"github.com/lestrrat-go/jwx/v2/jwk"
	cose "github.com/veraison/go-cose"
)

func getAlgAndKeyFromJWK(j []byte) (cose.Algorithm, crypto.Signer, error) {
	const noAlg = cose.Algorithm(-65537)
	var (
		err error
		k   jwk.Key
		crv elliptic.Curve
		alg cose.Algorithm
	)

	k, err = jwk.ParseKey(j)
	if err != nil {
		return noAlg, nil, err
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
	alg, key, err := getAlgAndKeyFromJWK(j)
	if err != nil {
		return nil, err
	}

	return cose.NewSigner(alg, key)
}

func NewPublicKeyFromJWK(j []byte) (crypto.PublicKey, error) {
	_, key, err := getAlgAndKeyFromJWK(j)
	if err != nil {
		return nil, err
	}

	return key.Public(), nil
}
