// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	comidDisplayFiles []string
	comidDisplayDirs  []string
)

var comidDisplayCmd = NewComidDisplayCmd()

func NewComidDisplayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "display",
		Short: "display one or more CBOR-encoded CoMID(s) in human readable (JSON) format",
		Long: `display one or more CBOR-encoded CoMID(s) in human readable (JSON) format.
	You can supply individual CoMID files or directories containing CoMID files.

	Display CoMID in file c.cbor.

	  cocli comid display --file=c.cbor

	Display CoMIDs in files c1.cbor, c2.cbor and any cbor file in the comids/
	directory.
	
	  cocli comid display --file=c1.cbor --file=c2.cbor --dir=comids
	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkComidDisplayArgs(); err != nil {
				return err
			}

			filesList := filesList(comidDisplayFiles, comidDisplayDirs, ".cbor")
			if len(filesList) == 0 {
				return errors.New("no files found")
			}

			errs := 0
			for _, file := range filesList {
				if err := displayComidFile(file); err != nil {
					fmt.Printf(">> failed displaying %q: %v\n", file, err)
					errs++
					continue
				}
			}

			if errs != 0 {
				return fmt.Errorf("%d/%d display(s) failed", errs, len(filesList))
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(
		&comidDisplayFiles, "file", "f", []string{}, "a CoMID file (in CBOR format)",
	)

	cmd.Flags().StringArrayVarP(
		&comidDisplayDirs, "dir", "d", []string{}, "a directory containing CoMID files (in CBOR format)",
	)

	return cmd
}

func displayComidFile(file string) error {
	var (
		data []byte
		err  error
	)

	if data, err = afero.ReadFile(fs, file); err != nil {
		return fmt.Errorf("error loading CoMID from %s: %w", file, err)
	}

	// use file name as heading
	return printComid(data, ">> ["+file+"]")
}

func checkComidDisplayArgs() error {
	if len(comidDisplayFiles) == 0 && len(comidDisplayDirs) == 0 {
		return errors.New("no files supplied")
	}
	return nil
}

func init() {
	comidCmd.AddCommand(comidDisplayCmd)
}
