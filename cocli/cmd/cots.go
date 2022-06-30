// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cotsCmd = &cobra.Command{
	Use:   "cots",
	Short: "CoTS manipulation",
	Long: "To prepare a signed CoTS file, several commands must be executed sequentially. Use the createStore command" +
		"to create concise-ta-store-map files. Use the create command to assemble these into a concise-ta-stores file. " +
		"Use the createCorim command to prepare an unsigned CoRIM file then use the sign command to prepare the final " +
		"signed CoRIM file containing CoTS data.",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help() // nolint: errcheck
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(cotsCmd)
}
