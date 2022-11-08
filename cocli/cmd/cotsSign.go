// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/veraison/corim/corim"
	cose "github.com/veraison/go-cose"
)

var (
	cotsSignCotsFile   *string
	cotsSignKeyFile    *string
	cotsSignOutputFile *string
	cotsSignMetaFile   *string
)

var cotsSignCmd = NewCotsSignCmd()

func NewCotsSignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign",
		Short: "create a signed CoTS from an unsigned, CBOR-encoded CoTS using the supplied key",
		Long: `create a signed CoTS from an unsigned, CBOR-encoded CoTS using the supplied key

	Sign the unsigned CoTS unsigned-cots.cbor using the key in JWK format from
	file key.jwk and save the resulting COSE Sign1 to signed-cots.cbor.  Read
	the relevant CorimMeta information from file meta.json.
	
	  cocli corim sign  --file=unsigned-cots.cbor \
					--key=key.jwk \
					--meta=meta.json \
					--output=signed-cots.cbor
	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkCotsSignArgs(); err != nil {
				return err
			}

			// checkCotsSignArgs makes sure cotsSignCorimFile is not nil
			coseFile, err := signCots(*cotsSignCotsFile, *cotsSignKeyFile,
				*cotsSignMetaFile, cotsSignOutputFile)
			if err != nil {
				return err
			}
			fmt.Printf(">> %q signed and saved to %q\n", *cotsSignCotsFile, coseFile)

			return nil
		},
	}

	cotsSignCotsFile = cmd.Flags().StringP("file", "f", "", "an unsigned CoTS file (in CBOR format)")
	cotsSignMetaFile = cmd.Flags().StringP("meta", "m", "", "CoRIM Meta file (in JSON format)")
	cotsSignKeyFile = cmd.Flags().StringP("key", "k", "", "signing key in JWK format")
	cotsSignOutputFile = cmd.Flags().StringP("output", "o", "", "name of the generated COSE Sign1 file")

	return cmd
}

func checkCotsSignArgs() error {
	if cotsSignCotsFile == nil || *cotsSignCotsFile == "" {
		return errors.New("no CoTS supplied")
	}

	if cotsSignKeyFile == nil || *cotsSignKeyFile == "" {
		return errors.New("no key supplied")
	}

	if cotsSignMetaFile == nil || *cotsSignMetaFile == "" {
		return errors.New("no CoRIM Meta supplied")
	}

	return nil
}

func signCots(unsignedCorimFile, keyFile, metaFile string, outputFile *string) (string, error) {
	var (
		unsignedCorimCBOR []byte
		signedCorimCBOR   []byte
		metaJSON          []byte
		keyJWK            []byte
		err               error
		signedCorimFile   string
		c                 corim.UnsignedCorim
		m                 corim.Meta
		signer            cose.Signer
	)

	if unsignedCorimCBOR, err = afero.ReadFile(fs, unsignedCorimFile); err != nil {
		return "", fmt.Errorf("error loading unsigned CoRIM from %s: %w", unsignedCorimFile, err)
	}

	if err = c.FromCBOR(unsignedCorimCBOR); err != nil {
		return "", fmt.Errorf("error decoding unsigned CoRIM from %s: %w", unsignedCorimFile, err)
	}

	if err = c.Valid(); err != nil {
		return "", fmt.Errorf("error validating CoRIM: %w", err)
	}

	if metaJSON, err = afero.ReadFile(fs, metaFile); err != nil {
		return "", fmt.Errorf("error loading CoRIM Meta from %s: %w", metaFile, err)
	}

	if err = m.FromJSON(metaJSON); err != nil {
		return "", fmt.Errorf("error decoding CoRIM Meta from %s: %w", metaFile, err)
	}

	if err = m.Valid(); err != nil {
		return "", fmt.Errorf("error validating CoRIM Meta: %w", err)
	}

	if keyJWK, err = afero.ReadFile(fs, keyFile); err != nil {
		return "", fmt.Errorf("error loading signing key from %s: %w", keyFile, err)
	}

	if signer, err = corim.NewSignerFromJWK(keyJWK); err != nil {
		return "", fmt.Errorf("error loading signing key from %s: %w", keyFile, err)
	}

	s := corim.SignedCorim{
		UnsignedCorim: c,
		Meta:          m,
	}

	signedCorimCBOR, err = s.Sign(signer)
	if err != nil {
		return "", fmt.Errorf("error signing CoRIM: %w", err)
	}

	if outputFile == nil || *outputFile == "" {
		signedCorimFile = "signed-" + unsignedCorimFile
	} else {
		signedCorimFile = *outputFile
	}

	err = afero.WriteFile(fs, signedCorimFile, signedCorimCBOR, 0644)
	if err != nil {
		return "", fmt.Errorf("error saving signed CoRIM to file %s: %w", signedCorimFile, err)
	}

	return signedCorimFile, nil
}

func init() {
	cotsCmd.AddCommand(cotsSignCmd)
}
