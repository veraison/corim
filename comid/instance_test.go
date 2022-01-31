package comid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstance_GetUUID_OK(t *testing.T) {
	inst := NewInstanceUUID(uuid.UUID(TestUUID))
	require.NotNil(t, inst)
	uuid, err := inst.GetUUID()
	assert.Nil(t, err)
	assert.Equal(t, uuid, TestUUID)
}

func TestInstance_GetUUID_NOK(t *testing.T) {
	inst := &Instance{}
	expectedErr := "instance-id type is: <nil>"
	_, err := inst.GetUUID()
	assert.EqualError(t, err, expectedErr)
}
