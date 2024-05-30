// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/veraison/apiclient/auth"
	"github.com/veraison/apiclient/common"
)

type ISubmitter interface {
	Run([]byte, string) error
	SetClient(client *common.Client) error
	SetAuth(a auth.IAuthenticator)
	SetSubmitURI(uri string) error
	SetDeleteSession(session bool)
	SetIsInsecure(v bool)
	SetCerts(paths []string)
}
