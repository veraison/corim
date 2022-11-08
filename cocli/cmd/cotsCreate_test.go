// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_CotsCreateCmd_unknown_argument(t *testing.T) {
	cmd := NewCotsCreateCmd()

	args := []string{"--unknown-argument=val"}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "unknown flag: --unknown-argument")
}

func Test_CotsCreateCmd_no_templates(t *testing.T) {
	cmd := NewCotsCreateCmd()

	// no args

	err := cmd.Execute()
	assert.EqualError(t, err, "no CTS files or folders supplied")
}

func Test_CotsCreateCmd_no_files_found(t *testing.T) {
	cmd := NewCotsCreateCmd()

	args := []string{
		"--output=output.cbor",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CTS files or folders supplied")
}

func Test_CotsCreateCmd_cts_not_found(t *testing.T) {
	cmd := NewCotsCreateCmd()

	args := []string{
		"--output=output.cbor",
		"--ctsfile=nonexistent.cbor",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.EqualError(t, err, "no CTS files found")
}

func Test_CotsCreateCmd_template_with_invalid_cts(t *testing.T) {
	var err error

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	cmd := NewCotsCreateCmd()

	fs = afero.NewOsFs()

	args := []string{
		"--ctsfile=../data/cots/not.cbor",
	}
	cmd.SetArgs(args)
	fmt.Println(args)

	err = cmd.Execute()
	assert.EqualError(t, err, "failed to parse as CBOR from ../data/cots/not.cbor: cbor: cannot unmarshal negative integer into Go value of type cots.ConciseTaStore")
}
