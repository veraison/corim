// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	cotsDisplayFiles []string
	cotsDisplayDirs  []string
)

var cotsDisplayCmd = NewCotsDisplayCmd()

func NewCotsDisplayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "display",
		Short: "display one or more CBOR-encoded CoTS(s) in human readable (JSON) format",
		Long: `display one or more CBOR-encoded CoTS(s) in human readable (JSON) format.
		You can supply individual CoTS files or directories containing CoTS files

	Display CoTS in cots.cbor 
	
	  cocli cots display --file=cots.cbor

	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkCotsDisplayArgs(); err != nil {
				return err
			}

			filesList := filesList(cotsDisplayFiles, cotsDisplayDirs, ".cbor")
			if len(filesList) == 0 {
				return errors.New("no files found")
			}

			errs := 0
			for _, file := range filesList {
				if err := displayCotsFile(file); err != nil {
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
		&cotsDisplayFiles, "file", "f", []string{}, "a CoTS file (in CBOR format)",
	)

	cmd.Flags().StringArrayVarP(
		&cotsDisplayDirs, "dir", "d", []string{}, "a directory containing CoTS files (in CBOR format)",
	)

	return cmd
}

func displayCotsFile(file string) error {
	var (
		data []byte
		err  error
	)

	if data, err = afero.ReadFile(fs, file); err != nil {
		return fmt.Errorf("error loading CoTS from %s: %w", file, err)
	}

	// use file name as heading
	return printCots(data, ">> ["+file+"]")

}

func checkCotsDisplayArgs() error {
	if len(cotsDisplayFiles) == 0 && len(cotsDisplayDirs) == 0 {
		return errors.New("no files supplied")
	}

	return nil
}

func init() {
	cotsCmd.AddCommand(cotsDisplayCmd)
}
