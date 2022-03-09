package cmd

import "github.com/veraison/apiclient/common"

type ISubmitter interface {
	Run([]byte, string) error
	SetClient(client *common.Client) error
	SetSubmitURI(uri string) error
	SetDeleteSession(session bool)
}
