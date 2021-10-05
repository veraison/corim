// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
)

var (
	minimalCorimTemplate = []byte(`{ "corim-id": "5c57e8f4-46cd-421b-91c9-08cf93e13cfc" }`)
	BadCBOR              = comid.MustHexDecode(nil, "ffff")
	// a "tag-id only" CoMID {1: {0: h'366D0A0A598845ED84882F2A544F6242'}}
	InvalidComid = comid.MustHexDecode(nil, "a101a10050366d0a0a598845ed84882f2a544f6242")
	testComid    = comid.MustHexDecode(nil, "a40065656e2d474201a1005043bbe37f2e614b33aed353cff1428b160281a3006941434d45204c74642e01d8207468747470733a2f2f61636d652e6578616d706c65028300010204a1008182a100a300d90227582061636d652d696d706c656d656e746174696f6e2d69642d303030303030303031016441434d45026a526f616452756e6e657283a200d90258a30162424c0465322e312e30055820acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b01a102818201582087428fc522803d31065e7bce3cf03fe475096631e5e07bbd7a0fde60c4cf25c7a200d90258a3016450526f540465312e332e35055820acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b01a10281820158200263829989b6fd954f72baaf2fc64bc2e2f01d692d4de72986ea808f6e99813fa200d90258a3016441526f540465302e312e34055820acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b01a1028182015820a3a5e715f0cc574a73c3f9bebb6bc24f32ffd5b67b387244c2c909da779a1478")
	testCoswid   = comid.MustHexDecode(nil, "a8007820636f6d2e61636d652e727264323031332d63652d7370312d76342d312d352d300c0001783041434d4520526f616472756e6e6572204465746563746f72203230313320436f796f74652045646974696f6e205350310d65342e312e3505a5182b65747269616c182d6432303133182f66636f796f7465183473526f616472756e6e6572204465746563746f721836637370310282a3181f745468652041434d4520436f72706f726174696f6e18206861636d652e636f6d1821820102a3181f75436f796f74652053657276696365732c20496e632e18206c6d79636f796f74652e636f6d18210404a21826781c7777772e676e752e6f72672f6c6963656e7365732f67706c2e7478741828676c6963656e736506a110a318186a72726465746563746f7218196d2570726f6772616d6461746125181aa111a318186e72726465746563746f722e657865141a000820e80782015820a314fc2dc663ae7a6b6bc6787594057396e6b3f569cd50fd5ddb4d1bbafd2b6a")
)

func Test_CorimCreateCmd_unknown_argument(t *testing.T) {
	cmd := NewCorimCreateCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_CorimCreateCmd_no_templates(t *testing.T) {
	cmd := NewCorimCreateCmd()

	// no args

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoRIM template supplied")
}

func Test_CorimCreateCmd_no_files_found(t *testing.T) {
	cmd := NewCorimCreateCmd()

	args := []string{
		"--tmpl-file=unknown.json",
		"--comid-file=unsure.cbor",
		"--comid-dir=somedir",
		"--coswid-file=what.cbor",
		"--coswid-dir=someotherdir",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoMID or CoSWID files found")
}

func Test_CorimCreateCmd_no_tag_files(t *testing.T) {
	cmd := NewCorimCreateCmd()

	args := []string{
		"--tmpl-file=unknown.json",
		// no --co{m,sw}id-{dir,file}
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CoMID or CoSWID files or folders supplied")
}

func Test_CorimCreateCmd_template_not_found(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "ignored-comid.cbor", []byte{}, 0644)
	require.NoError(t, err)

	args := []string{
		"--tmpl-file=nonexistent.json",
		"--comid-file=ignored-comid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, "error loading template from nonexistent.json: open nonexistent.json: file does not exist")
}

func Test_CorimCreateCmd_template_with_invalid_json(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "invalid.json", []byte("..."), 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "ignored-comid.cbor", []byte{}, 0644)
	require.NoError(t, err)

	args := []string{
		"--tmpl-file=invalid.json",
		"--comid-file=ignored-comid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, "error decoding template from invalid.json: invalid character '.' looking for beginning of value")
}

func Test_CorimCreateCmd_with_a_bad_comid(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "min-tmpl.json", minimalCorimTemplate, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "bad-comid.cbor", BadCBOR, 0644)
	require.NoError(t, err)

	args := []string{
		"--tmpl-file=min-tmpl.json",
		"--comid-file=bad-comid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, `error loading CoMID from bad-comid.cbor: cbor: unexpected "break" code`)
}

func Test_CorimCreateCmd_with_an_invalid_comid(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "min-tmpl.json", minimalCorimTemplate, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "invalid-comid.cbor", InvalidComid, 0644)
	require.NoError(t, err)

	args := []string{
		"--tmpl-file=min-tmpl.json",
		"--comid-file=invalid-comid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, `error adding CoMID from invalid-comid.cbor (check its validity using the "comid validate" sub-command)`)
}

func Test_CorimCreateCmd_with_a_bad_coswid(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "min-tmpl.json", minimalCorimTemplate, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "bad-coswid.cbor", BadCBOR, 0644)
	require.NoError(t, err)

	args := []string{
		"--tmpl-file=min-tmpl.json",
		"--coswid-file=bad-coswid.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.EqualError(t, err, `error loading CoSWID from bad-coswid.cbor: cbor: unexpected "break" code`)
}

func Test_CorimCreateCmd_successful_comid_and_coswid_from_file(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "min-tmpl.json", minimalCorimTemplate, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "coswid.cbor", testCoswid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "comid.cbor", testComid, 0644)
	require.NoError(t, err)

	args := []string{
		"--tmpl-file=min-tmpl.json",
		"--coswid-file=coswid.cbor",
		"--comid-file=comid.cbor",
		"--output-file=corim.cbor",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.NoError(t, err)

	_, err = fs.Stat("corim.cbor")
	assert.NoError(t, err)
}

func Test_CorimCreateCmd_successful_comid_and_coswid_from_dir(t *testing.T) {
	var err error

	cmd := NewCorimCreateCmd()

	fs = afero.NewMemMapFs()
	err = afero.WriteFile(fs, "min-tmpl.json", minimalCorimTemplate, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "coswid/1.cbor", testCoswid, 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "comid/1.cbor", testComid, 0644)
	require.NoError(t, err)

	args := []string{
		"--tmpl-file=min-tmpl.json",
		"--coswid-dir=coswid",
		"--comid-dir=comid",
	}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.NoError(t, err)

	_, err = fs.Stat("5c57e8f4-46cd-421b-91c9-08cf93e13cfc.cbor")
	assert.NoError(t, err)
}
