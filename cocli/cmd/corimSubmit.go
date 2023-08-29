// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veraison/apiclient/provisioning"
)

var (
	corimFile *string
	mediaType *string
	apiServer string
)

var (
	submitter      ISubmitter = &provisioning.SubmitConfig{}
	corimSubmitCmd            = NewCorimSubmitCmd(submitter)
)

func NewCorimSubmitCmd(submitter ISubmitter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit",
		Short: "submit a CBOR-encoded CoRIM payload",
		Long: `submit a CBOR-encoded CoRIM payload with supplied media type to the given API Server

	To submit the CBOR-encoded CoRIM from file "unsigned-corim.cbor" with media type
	"application/corim-unsigned+cbor; profile=http://arm.com/psa/iot/1" to the Veraison
	provisioning API endpoint "https://veraison.example/endorsement-provisioning/v1", do:


	cocli corim submit \
			--corim-file=unsigned-corim.cbor \
			--api-server="https://veraison.example/endorsement-provisioning/v1/submit" \
			--media-type="application/corim-unsigned+cbor; profile=http://arm.com/psa/iot/1"
	`,

		RunE: func(cmd *cobra.Command, args []string) error {

			if err := checkSubmitArgs(); err != nil {
				return err
			}

			// Load the data from the CBOR File
			data, err := readCorimData(*corimFile)
			if err != nil {
				return fmt.Errorf("read CoRIM payload failed: %w", err)
			}

			if err = provisionData(data, submitter, apiServer, *mediaType); err != nil {
				return fmt.Errorf("submit CoRIM payload failed reason: %w", err)
			}
			return nil
		},
	}

	corimFile = cmd.Flags().StringP("corim-file", "f", "", "name of the CoRIM file in CBOR format")
	mediaType = cmd.Flags().StringP("media-type", "m", "", "media type of the CoRIM file")

	cmd.Flags().StringP("api-server", "s", "", "API server where to submit the corim file")
	cmd.Flags().VarP(&authMethod, "auth", "a",
		`authentication method, must be one of "none"/"passthrough", "basic", "oauth2"`)
	cmd.Flags().StringP("client-id", "C", "", "OAuth2 client ID")
	cmd.Flags().StringP("client-secret", "S", "", "OAuth2 client secret")
	cmd.Flags().StringP("token-url", "T", "", "token URL of the OAuth2 service")
	cmd.Flags().StringP("username", "U", "", "service username")
	cmd.Flags().StringP("password", "P", "", "service password")

	err := viper.BindPFlag("api_server", cmd.Flags().Lookup("api-server"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("auth", cmd.Flags().Lookup("auth"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("client_id", cmd.Flags().Lookup("client-id"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("client_secret", cmd.Flags().Lookup("client-secret"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("username", cmd.Flags().Lookup("username"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("password", cmd.Flags().Lookup("password"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("token_url", cmd.Flags().Lookup("token-url"))
	cobra.CheckErr(err)

	return cmd
}

func checkSubmitArgs() error {
	if corimFile == nil || *corimFile == "" {
		return errors.New("no CoRIM input file supplied")
	}

	apiServer = viper.GetString("api_server")
	if apiServer == "" {
		return errors.New("no API server supplied")
	}
	u, err := url.Parse(apiServer)
	if err != nil || !u.IsAbs() {
		return fmt.Errorf("malformed API server URL")
	}

	if mediaType == nil || *mediaType == "" {
		return errors.New("no media type supplied")
	}

	return nil
}

func provisionData(data []byte, submitter ISubmitter, uri string, mediaType string) error {
	submitter.SetAuth(cliConfig.Auth)

	if err := submitter.SetSubmitURI(uri); err != nil {
		return fmt.Errorf("unable to set submit URI: %w", err)
	}

	submitter.SetDeleteSession(true)
	if err := submitter.Run(data, mediaType); err != nil {
		return fmt.Errorf("run failed: %w", err)
	}

	return nil
}

func readCorimData(file string) ([]byte, error) {
	return afero.ReadFile(fs, file)
}

func init() {
	corimCmd.AddCommand(corimSubmitCmd)
}
