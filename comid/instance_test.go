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
