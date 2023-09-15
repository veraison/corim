package corim

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

const testParam2 = 0x40

func TestProfileEntities_Valid1_ok(t *testing.T) {
	e := &ProfileXEntity{}
	e.SetEntityName("Intel")
	e.SetRegID("http://intelprofile.com")
	e.SetRoles(RoleManifestCreator)
	e.SetEntityExtension("testParam1", testParam2)
	err := e.Valid()
	require.NoError(t, err)
}

func TestProfileEntities_ToCBOR_ok(t *testing.T) {
	e := &ProfileXEntity{}
	e.SetEntityName("Intel")
	e.SetRegID("http://intelprofile.com")
	e.SetRoles(RoleManifestCreator)
	e.SetEntityExtension("testParam1", testParam2)
	err := e.Valid()
	require.NoError(t, err)
	data, err := e.ToCBOR()
	require.NoError(t, err)
	fmt.Printf("CBOR Encoded Profile = %x", data)
}
