package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndorsedValue_Valid_empty_environment(t *testing.T) {
	tv := EndorsedValue{}

	assert.EqualError(t, tv.Valid(), "environment validation failed: environment must not be empty")
}

func TestEndorsedValue_Valid_empty_measurements(t *testing.T) {
	tv := EndorsedValue{
		Environment: Environment{
			Class: NewClassUUID(TestUUID),
		},
	}

	assert.EqualError(t, tv.Valid(), "measurements validation failed: no measurement entries")
}
