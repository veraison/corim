package corim

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"reflect"

	"github.com/lestrrat-go/jwx/jwk"
	cose "github.com/veraison/go-cose"
)

func SignerFromJWK(j []byte) (*cose.Signer, error) {
	var (
		err  error
		ks   jwk.Set
		k    jwk.Key
		ok   bool
		pkey crypto.PrivateKey
		crv  elliptic.Curve
		alg  *cose.Algorithm
	)

	if ks, err = jwk.ParseString(string(j)); err != nil {
		return nil, err
	}

	if k, ok = ks.Get(0); !ok {
		return nil, errors.New("no key found at slot 0")
	}

	if err = k.Raw(&pkey); err != nil {
		return nil, err
	}

	switch v := pkey.(type) {
	case *ecdsa.PrivateKey:
		crv = v.Curve
		if crv == elliptic.P256() {
			alg = cose.ES256
			break
		}
		return nil, fmt.Errorf("unknown elliptic curve %v", crv)
	default:
		return nil, fmt.Errorf("unknown private key type %v", reflect.TypeOf(pkey))
	}

	return cose.NewSignerFromKey(alg, pkey)
}
