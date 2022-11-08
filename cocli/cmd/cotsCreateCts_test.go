// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_CotsCreateCtsCmd_unknown_argument(t *testing.T) {
	cmd := NewCotsCreateCtsCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_CotsCreateCtsCmd_no_templates(t *testing.T) {
	cmd := NewCotsCreateCtsCmd()

	// no args

	err := cmd.Execute()
	assert.EqualError(t, err, "no environment template supplied")
}

func Test_CotsCreateCtsCmd_no_files_found(t *testing.T) {
	cmd := NewCotsCreateCtsCmd()

	args := []string{
		"--output=output.cbor",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no environment template supplied")
}

func Test_CotsCreateCtsCmd_env_not_found_no_tas(t *testing.T) {
	cmd := NewCotsCreateCtsCmd()

	args := []string{
		"--output=output.cbor",
		"--environment=nonexistent.cbor",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no TA files or folders supplied")
}

func Test_CotsCreateCtsCmd_env_not_found(t *testing.T) {
	cmd := NewCotsCreateCtsCmd()

	args := []string{
		"--output=output.cbor",
		"--environment=nonexistent.cbor",
		"--tafile=../data/cots/shared_ta.der",
	}
	cmd.SetArgs(args)
	fs = afero.NewOsFs()
	err := cmd.Execute()
	assert.EqualError(t, err, "no TA files found")
}

func Test_CotsCreateCtsCmd_too_many_ids(t *testing.T) {
	cmd := NewCotsCreateCtsCmd()

	args := []string{
		"--output=output.cbor",
		"--uuid",
		"--id=some_tag_identity",
		"--environment=../data/cots/env/vendor.json",
		"--tafile=../data/cots/shared_ta.der",
	}
	cmd.SetArgs(args)
	fs = afero.NewOsFs()
	err := cmd.Execute()
	assert.EqualError(t, err, "only one of --uuid, --uuid-str and --id can be used at the same time")
}

func Test_CotsCreateCtsCmd_invalid_uuid(t *testing.T) {
	cmd := NewCotsCreateCtsCmd()

	args := []string{
		"--output=output.cbor",
		"--uuid-str=NotAUuid",
		"--environment=../data/cots/env/vendor.json",
		"--tafile=../data/cots/shared_ta.der",
	}
	cmd.SetArgs(args)
	fs = afero.NewOsFs()
	err := cmd.Execute()
	assert.EqualError(t, err, "--uuid-str does not contain a valid UUID")
}

// TODO add more test cases
