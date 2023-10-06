package comid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ReferenceValue(t *testing.T) {
	rv := ReferenceValue{}
	err := rv.Valid()
	assert.EqualError(t, err, "environment validation failed: environment must not be empty")

	rv.Environment.Instance = NewInstance()
	id, err := uuid.NewUUID()
	require.NoError(t, err)
	rv.Environment.Instance.SetUUID(id)
	err = rv.Valid()
	assert.EqualError(t, err, "measurements validation failed: no measurement entries")
}
