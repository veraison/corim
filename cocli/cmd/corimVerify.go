// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"crypto"
	"errors"
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/veraison/corim/corim"
)

var (
	corimVerifyCorimFile *string
	corimVerifyKeyFile   *string
)

var corimVerifyCmd = NewCorimVerifyCmd()

func NewCorimVerifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "verify a signed CoRIM using the supplied key",
		Long: `verify a signed CoRIM using the supplied key

	Verify the signed CoRIM signed-corim.cbor using the key in JWK format from
	file key.jwk
	
	  cocli corim verify --file=signed-corim.cbor --key=key.jwk
	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkCorimVerifyArgs(); err != nil {
				return err
			}

			// checkCorimVerifyArgs makes sure corimVerifyCorimFile is not nil
			err := verify(*corimVerifyCorimFile, *corimVerifyKeyFile)
			if err != nil {
				return err
			}
			fmt.Printf(">> %q verified\n", *corimVerifyCorimFile)

			return nil
		},
	}

	corimVerifyCorimFile = cmd.Flags().StringP("file", "f", "", "a signed CoRIM file (in CBOR format)")
	corimVerifyKeyFile = cmd.Flags().StringP("key", "k", "", "verification key in JWK format")

	return cmd
}

func checkCorimVerifyArgs() error {
	if corimVerifyCorimFile == nil || *corimVerifyCorimFile == "" {
		return errors.New("no CoRIM supplied")
	}

	if corimVerifyKeyFile == nil || *corimVerifyKeyFile == "" {
		return errors.New("no key supplied")
	}

	return nil
}

func verify(signedCorimFile, keyFile string) error {
	var (
		signedCorimCBOR []byte
		keyJWK          []byte
		err             error
		pkey            crypto.PublicKey
		s               corim.SignedCorim
	)

	if signedCorimCBOR, err = afero.ReadFile(fs, signedCorimFile); err != nil {
		return fmt.Errorf("error loading signed CoRIM from %s: %w", signedCorimFile, err)
	}

	if err = s.FromCOSE(signedCorimCBOR); err != nil {
		return fmt.Errorf("error decoding signed CoRIM from %s: %w", signedCorimFile, err)
	}

	if keyJWK, err = afero.ReadFile(fs, keyFile); err != nil {
		return fmt.Errorf("error loading verifying key from %s: %w", keyFile, err)
	}

	if pkey, err = corim.NewPublicKeyFromJWK(keyJWK); err != nil {
		return fmt.Errorf("error loading verifying key from %s: %w", keyFile, err)
	}

	if err = s.Verify(pkey); err != nil {
		return fmt.Errorf("error verifying %s with key %s: %w", signedCorimFile, keyFile, err)
	}

	return nil
}

func init() {
	corimCmd.AddCommand(corimVerifyCmd)
}
