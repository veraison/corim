// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var comidCmd = &cobra.Command{
	Use:   "comid",
	Short: "CoMID specific manipulation",
	Long: `CoMID specific manipulation
	
	Create CoMID from template t1.json and save the result in the current working directory.

	cli comid create --file=t1.json

	`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help() // nolint: errcheck
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(comidCmd)
}
