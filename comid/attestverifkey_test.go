package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttestVerifKey_Valid_empty(t *testing.T) {

	tvs := []struct {
		env      Environment
		verifkey VerifKeys
		testerr  string
	}{
		{
			env:      Environment{},
			verifkey: VerifKeys{},
			testerr:  "environment validation failed: environment must not be empty",
		},
		{
			env:      Environment{Instance: NewInstanceUEID(TestUEID)},
			verifkey: VerifKeys{},
			testerr:  "verification keys validation failed: no verification key to validate",
		},
		{
			env:      Environment{Instance: NewInstanceUEID(TestUEID)},
			verifkey: VerifKeys{{Key: ""}},
			testerr:  "verification keys validation failed: invalid verification key at index 0: verification key not set",
		},
	}
	for _, tv := range tvs {
		av := AttestVerifKey{Environment: tv.env, VerifKeys: tv.verifkey}
		err := av.Valid()
		assert.EqualError(t, err, tv.testerr)
	}
}
