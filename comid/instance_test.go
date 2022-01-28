package comid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInstance_GetUUID_OK(t *testing.T) {
	inst := NewInstanceUUID(uuid.UUID(TestUUID))
	assert.NotNil(t, inst)
	uuid, err := inst.GetUUID()
	assert.Nil(t, err)
	assert.Equal(t, uuid, TestUUID.String())
}

func TestInstance_GetUUID_NOK(t *testing.T) {
	inst := &Instance{}
	expectedErr := "instance-id type is: <nil>"
	_, err := inst.GetUUID()
	assert.EqualError(t, err, expectedErr)
}
