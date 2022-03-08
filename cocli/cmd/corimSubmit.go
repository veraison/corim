// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/veraison/apiclient/provisioning"
)

var (
	corimFile      *string
	apiServer      *string
	mediaType      *string
	submitter      ISubmitter = &provisioning.SubmitConfig{}
	corimSubmitCmd            = NewCorimSubmitCmd(submitter)
)

func NewCorimSubmitCmd(submitter ISubmitter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit",
		Short: "submit a CBOR-encoded CoRIM payload",
		Long: `submit a CBOR-encoded CoRIM payload with supplied api server and media type as inputs

	Submit a CBOR encoded CoRIM payload from unsigned-corim.cbor to Veraison provisioning server
	"https://veraison.example/endorsement-provisioning/v1" using 
	"application/corim-unsigned+cbor; profile=http://arm.com/psa/iot/1" as the media type
	

	cocli corim submit \
			--corim-file = unsigned-corim.cbor \
			--api-server = "https://veraison.example/endorsement-provisioning/v1"
			--media-type = "application/corim-unsigned+cbor; profile=http://arm.com/psa/iot/1"
	`,

		RunE: func(cmd *cobra.Command, args []string) error {

			if err := checkSubmitArgs(); err != nil {
				return err
			}

			// Load the data from the CBOR File
			data, err := readCorimData(*corimFile)
			if err != nil {
				return fmt.Errorf("corim payload read failed: %w", err)
			}

			if err = provisionData(data, submitter, *apiServer, *mediaType); err != nil {
				return fmt.Errorf("corim submit failed reason: %w", err)
			}
			return nil
		},
	}

	corimFile = cmd.Flags().StringP("corim-file", "f", "", "name of the CoRIM file in CBOR format")
	apiServer = cmd.Flags().StringP("api-server", "s", "", "API server where to submit the corim file")
	mediaType = cmd.Flags().StringP("media-type", "m", "", "media type of the file")

	return cmd
}

func checkSubmitArgs() error {
	if corimFile == nil || *corimFile == "" {
		return errors.New("no CoRIM input file supplied")
	}
	if apiServer == nil || *apiServer == "" {
		return errors.New("no API server in the argument")
	}
	u, err := url.Parse(*apiServer)
	if err != nil || !u.IsAbs() {
		return fmt.Errorf("malformed API server URL")
	}

	if mediaType == nil || *mediaType == "" {
		return errors.New("no media type supplied")
	}

	return nil
}

func provisionData(data []byte, submitter ISubmitter, uri string, mediaType string) error {
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
	data, err := afero.ReadFile(fs, file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func init() {
	corimCmd.AddCommand(corimSubmitCmd)
}
