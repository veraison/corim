package cmd

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock_deps "github.com/veraison/corim/cocli/cmd/mocks"
)

func Test_CorimSubmitCmd_bad_server_url(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_deps.NewMockISubmitter(ctrl)
	cmd := NewCorimSubmitCmd(ms)

	args := []string{
		"--corim-file=corim.cbor",
		"--api-server=http://www.example.com:80index",
		"--media-type=application/corim-unsigned+cbor; profile=http://arm.com/psa/iot/1",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "corim.cbor", testSignedCorimValid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "malformed API server URL")
}

func Test_CorimSubmitCmd_missing_server_url(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_deps.NewMockISubmitter(ctrl)
	cmd := NewCorimSubmitCmd(ms)

	args := []string{
		"--corim-file=corim.cbor",
		"--media-type=application/corim-unsigned+cbor; profile=http://arm.com/psa/iot/1",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "corim.cbor", testSignedCorimValid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "no API server supplied")
}

func Test_CorimSubmitCmd_missing_media_type(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_deps.NewMockISubmitter(ctrl)
	cmd := NewCorimSubmitCmd(ms)

	args := []string{
		"--corim-file=corim.cbor",
		"--api-server=http://www.example.com:8080",
		"--media-type=",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "corim.cbor", testSignedCorimValid, 0644)
	require.NoError(t, err)

	err = cmd.Execute()
	assert.EqualError(t, err, "no media type supplied")

}

func Test_CorimSubmitCmd_missing_corim_file(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_deps.NewMockISubmitter(ctrl)
	cmd := NewCorimSubmitCmd(ms)

	args := []string{
		"--corim-file=",
		"--api-server=http://www.example.com:8080",
		"--media-type=application/corim-unsigned+cbor; profile=http://arm.com/psa/iot/1",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoRIM input file supplied")

}

func Test_CorimSubmitCmd_non_existent_corim_file(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_deps.NewMockISubmitter(ctrl)
	cmd := NewCorimSubmitCmd(ms)

	args := []string{
		"--corim-file=bad.cbor",
		"--api-server=http://www.example.com:8080",
		"--media-type=application/corim-unsigned+cbor; profile=http://arm.com/psa/iot/1",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "corim payload read failed: open bad.cbor: file does not exist")
}

func Test_CorimSubmitCmd_submit_ok(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_deps.NewMockISubmitter(ctrl)
	cmd := NewCorimSubmitCmd(ms)

	args := []string{
		"--corim-file=corim.cbor",
		"--api-server=http://veraison.example/endorsement-provisioning/v1/submit",
		"--media-type=application/corim-unsigned+cbor; profile=http://arm.com/psa/iot/1",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "corim.cbor", testSignedCorimValid, 0644)
	require.NoError(t, err)
	ms.EXPECT().SetSubmitURI("http://veraison.example/endorsement-provisioning/v1/submit").Return(nil)
	ms.EXPECT().SetDeleteSession(true)
	ms.EXPECT().Run(testSignedCorimValid, "application/corim-unsigned+cbor; profile=http://arm.com/psa/iot/1").Return(nil)
	err = cmd.Execute()
	assert.NoError(t, err)
}

func Test_CorimSubmitCmd_submit_not_ok(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_deps.NewMockISubmitter(ctrl)
	cmd := NewCorimSubmitCmd(ms)

	args := []string{
		"--corim-file=corim.cbor",
		"--api-server=http://veraison.example/endorsement-provisioning/v1/submit",
		"--media-type=application/corim-unsigned+cbor; profile=http://arm.com/psa/iot/1",
	}
	cmd.SetArgs(args)

	fs = afero.NewMemMapFs()
	err := afero.WriteFile(fs, "corim.cbor", testSignedCorimValid, 0644)
	require.NoError(t, err)
	ms.EXPECT().SetSubmitURI("http://veraison.example/endorsement-provisioning/v1/submit").Return(nil)
	ms.EXPECT().SetDeleteSession(true)
	err = errors.New(`unexpected HTTP response code 404`)

	ms.EXPECT().Run(testSignedCorimValid, "application/corim-unsigned+cbor; profile=http://arm.com/psa/iot/1").Return(err)
	err = cmd.Execute()
	assert.EqualError(t, err, "corim submit failed reason: run failed: unexpected HTTP response code 404")
}
