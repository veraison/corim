// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cots

var (
	ConciseTaStoreTemplateSingleOrg = `{
		"tag-identity": {
			"id": "ab0f44b1-bfdc-4604-ab4a-30f80407ebcc",
			"version": 5
		},
		"environments": [
			{
				"environment": {
					"class": {
						"vendor": "Worthless Sea, Inc."
					}
				}
			}
		],
		"keys": {
			"tas": [
				{
					"format": 2,
					"data": "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAErYoMAdqe2gJT3CvCcifZxyE9+N8T6Jy5zbeo5LYtnOipmi1wXA9/gNtlwAbRCRQitH/GEcvUaGlzPZxIOITV/g=="
				}
			]
		}
	}`

	ConciseTaStoreTemplateMultipleOrgs = `{
		"tag-identity": {
		  "id": "some_tag_identity"
		},
		"environments": [
		  	{
				"namedtastore": "Miscellaneous TA Store"
			}
	  	],
		"keys": {
			"tas": [
				{
					"format": 0,
					"data": "MIIBvTCCAWSgAwIBAgIVANCdkL89UlzHc9Ui7XfVniK7pFuIMAoGCCqGSM49BAMCMD4xCzAJBgNVBAYMAlVTMRAwDgYDVQQKDAdFeGFtcGxlMR0wGwYDVQQDDBRFeGFtcGxlIFRydXN0IEFuY2hvcjAeFw0yMjA1MTkxNTEzMDdaFw0zMjA1MTYxNTEzMDdaMD4xCzAJBgNVBAYMAlVTMRAwDgYDVQQKDAdFeGFtcGxlMR0wGwYDVQQDDBRFeGFtcGxlIFRydXN0IEFuY2hvcjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABONRqhA5JAekvQN8oLwRVNDnAfBnTznLLE+SEGks677sHSeXfcVhZXUeDiN7/fsVNumaiEWRQpZh3zXPwL8rUMyjPzA9MB0GA1UdDgQWBBQBXEXJrLBGKnFd1xCgeMAVSfEBPzALBgNVHQ8EBAMCAoQwDwYDVR0TAQH/BAUwAwEB/zAKBggqhkjOPQQDAgNHADBEAiALBidABsfpzG0lTL9Eh9b6AUbqnzF+koEZbgvppvvt9QIgVoE+bhEN0j6wSPzePjLrEdD+PEgyjHJ5rbA11SPq/1M="
				},
				{
					"format": 1,
           			"data": "ooICtjCCArIwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASXz21w12owQAx58euratYWiHEkhxDU9MEgetrvAtGYZxNnkfLCsp9vLcw8ISXC8tL97k9ZCUtnr0MzLw37XKRABBT22tHlEou/DenpU0Ozccb3/+fibjCCAj0wUjELMAkGA1UEBgwCVVMxGjAYBgNVBAoMEVplc3R5IEhhbmRzLCBJbmMuMScwJQYDVQQDDB5aZXN0eSBIYW5kcywgSW5jLiBUcnVzdCBBbmNob3KgggHlMIIBi6ADAgECAhQL3EqgUXlQPljyddVSRnNHvK+1MzAKBggqhkjOPQQDAjBSMQswCQYDVQQGDAJVUzEaMBgGA1UECgwRWmVzdHkgSGFuZHMsIEluYy4xJzAlBgNVBAMMHlplc3R5IEhhbmRzLCBJbmMuIFRydXN0IEFuY2hvcjAeFw0yMjA1MTkxNTEzMDdaFw0zMjA1MTYxNTEzMDdaMFIxCzAJBgNVBAYMAlVTMRowGAYDVQQKDBFaZXN0eSBIYW5kcywgSW5jLjEnMCUGA1UEAwweWmVzdHkgSGFuZHMsIEluYy4gVHJ1c3QgQW5jaG9yMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEl89tcNdqMEAMefHrq2rWFohxJIcQ1PTBIHra7wLRmGcTZ5HywrKfby3MPCElwvLS/e5PWQlLZ69DMy8N+1ykQKM/MD0wHQYDVR0OBBYEFPba0eUSi78N6elTQ7Nxxvf/5+JuMAsGA1UdDwQEAwIChDAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49BAMCA0gAMEUCIB2li+f6RCxs2EnvNWciSpIDwiUViWayGv1A8xks80eYAiEAmCez4KGrolFKOZT6bvqf1sYQuJBfvtk/y1JQdUvoqlg="
				},
				{
					"format": 1,
					"data": "ooIC1TCCAtEwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAATN0f5kzywEzZOYbaV23O3N8cku39JoLNjlHPwECbXDDWp0LpAO1z248/hoy6UW/TZMTPPR/93XwHsG16mSFy8XBBSKhM/5gJWjvDbW7qUY1peNm9cfYDCCAlwwXDELMAkGA1UEBgwCVVMxHzAdBgNVBAoMFlNub2JiaXNoIEFwcGFyZWwsIEluYy4xLDAqBgNVBAMMI1Nub2JiaXNoIEFwcGFyZWwsIEluYy4gVHJ1c3QgQW5jaG9yoIIB+jCCAZ+gAwIBAgIUEBuTRGXAEEVEHhu4xafAnqm+qYgwCgYIKoZIzj0EAwIwXDELMAkGA1UEBgwCVVMxHzAdBgNVBAoMFlNub2JiaXNoIEFwcGFyZWwsIEluYy4xLDAqBgNVBAMMI1Nub2JiaXNoIEFwcGFyZWwsIEluYy4gVHJ1c3QgQW5jaG9yMB4XDTIyMDUxOTE1MTMwOFoXDTMyMDUxNjE1MTMwOFowXDELMAkGA1UEBgwCVVMxHzAdBgNVBAoMFlNub2JiaXNoIEFwcGFyZWwsIEluYy4xLDAqBgNVBAMMI1Nub2JiaXNoIEFwcGFyZWwsIEluYy4gVHJ1c3QgQW5jaG9yMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEzdH+ZM8sBM2TmG2ldtztzfHJLt/SaCzY5Rz8BAm1ww1qdC6QDtc9uPP4aMulFv02TEzz0f/d18B7BtepkhcvF6M/MD0wHQYDVR0OBBYEFIqEz/mAlaO8NtbupRjWl42b1x9gMAsGA1UdDwQEAwIChDAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49BAMCA0kAMEYCIQC2cf43f3PPlCO6/dxv40ftIgxxToKHF72UzENv7+y4ygIhAIGtC/r6SGaFMaP7zD2EloBuIXTtyWu8Hwl+YGdXRY93"
				}
			]
		}
	}`

	ConciseTaStoreTemplateEnvSWID = `
	{
		"environments": [
		  {
			"swidtag": {
			  "entity": [
				{
				  "entity-name": "Zesty Hands, Inc.",
				  "role": "softwareCreator"
				}
			  ]
			}
		  }
		],
		"permclaims": [
		  {
			"swname": "Bitter Paper"
		  }
		],
		"keys": {
			"tas": [
				{
					"format": 0,
					"data":"MIIB5TCCAYugAwIBAgIUC9xKoFF5UD5Y8nXVUkZzR7yvtTMwCgYIKoZIzj0EAwIwUjELMAkGA1UEBgwCVVMxGjAYBgNVBAoMEVplc3R5IEhhbmRzLCBJbmMuMScwJQYDVQQDDB5aZXN0eSBIYW5kcywgSW5jLiBUcnVzdCBBbmNob3IwHhcNMjIwNTE5MTUxMzA3WhcNMzIwNTE2MTUxMzA3WjBSMQswCQYDVQQGDAJVUzEaMBgGA1UECgwRWmVzdHkgSGFuZHMsIEluYy4xJzAlBgNVBAMMHlplc3R5IEhhbmRzLCBJbmMuIFRydXN0IEFuY2hvcjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABJfPbXDXajBADHnx66tq1haIcSSHENT0wSB62u8C0ZhnE2eR8sKyn28tzDwhJcLy0v3uT1kJS2evQzMvDftcpECjPzA9MB0GA1UdDgQWBBT22tHlEou/DenpU0Ozccb3/+fibjALBgNVHQ8EBAMCAoQwDwYDVR0TAQH/BAUwAwEB/zAKBggqhkjOPQQDAgNIADBFAiAdpYvn+kQsbNhJ7zVnIkqSA8IlFYlmshr9QPMZLPNHmAIhAJgns+Chq6JRSjmU+m76n9bGELiQX77ZP8tSUHVL6KpY"
				}
			]
		}
	}`

	cotsCBOR = []byte{0xa3, 0x1, 0xa2, 0x0, 0x50, 0xab, 0xf, 0x44, 0xb1, 0xbf,
		0xdc, 0x46, 0x4, 0xab, 0x4a, 0x30, 0xf8, 0x4, 0x7, 0xeb, 0xcc, 0x1,
		0x5, 0x2, 0x81, 0xa1, 0x1, 0xa1, 0x0, 0xa1, 0x1, 0x73, 0x57, 0x6f,
		0x72, 0x74, 0x68, 0x6c, 0x65, 0x73, 0x73, 0x20, 0x53, 0x65, 0x61,
		0x2c, 0x20, 0x49, 0x6e, 0x63, 0x2e, 0x6, 0xa1, 0x0, 0x81, 0x82, 0x2,
		0x58, 0x5b, 0x30, 0x59, 0x30, 0x13, 0x6, 0x7, 0x2a, 0x86, 0x48, 0xce,
		0x3d, 0x2, 0x1, 0x6, 0x8, 0x2a, 0x86, 0x48, 0xce, 0x3d, 0x3, 0x1, 0x7,
		0x3, 0x42, 0x0, 0x4, 0xad, 0x8a, 0xc, 0x1, 0xda, 0x9e, 0xda, 0x2, 0x53,
		0xdc, 0x2b, 0xc2, 0x72, 0x27, 0xd9, 0xc7, 0x21, 0x3d, 0xf8, 0xdf, 0x13,
		0xe8, 0x9c, 0xb9, 0xcd, 0xb7, 0xa8, 0xe4, 0xb6, 0x2d, 0x9c, 0xe8, 0xa9,
		0x9a, 0x2d, 0x70, 0x5c, 0xf, 0x7f, 0x80, 0xdb, 0x65, 0xc0, 0x6, 0xd1,
		0x9, 0x14, 0x22, 0xb4, 0x7f, 0xc6, 0x11, 0xcb, 0xd4, 0x68, 0x69, 0x73,
		0x3d, 0x9c, 0x48, 0x38, 0x84, 0xd5, 0xfe}
)
