
## CoTSs manipulation

The `cots` subcommand allows you to create, display and validate CoTSs.

### Create

Use the `cots create` subcommand to create a CBOR-encoded CoTS. The `environment` switch takes in a JSON template specifiying the environments that are valid for the keys specified and the `tas` switch takes in a directory of trust anchors files:

* Please inspect `data/cots/templates` JSON templates as examples for `environment` and `claims`


```
$ cocli cots create --environment data/cots/env/vendor.json --tafile data/cots/shared_ta.ta
```
On success, you should see something like the following printed to stdout:
```
>> created "vendor.cbor"
```

The CBOR-encoded CoTS file is stored in the current working directory with a
name derived from its environment template.  If you want, you can specify a different
target directory and file name using the `--output` command line switch (abbrev. `-o`)
```
$ cocli cots create --environment data/cots/env/vendor.json --tafile data/cots/shared_ta.ta --output /tmp/myCots.cbor
>> created "/tmp/myCots.cbor"
```
Note that the output directory, as well as all its parent directories, MUST pre-exist.

### Display

Use the `cots display` subcommand to print to stdout one or more CBOR-encoded
CoTSs in human readable (JSON) format.

You can supply individual files using the `--file` switch (abbrev. `-f`), or
directories that may (or may not) contain CoTS files using the `--dir` switch
(abbrev. `-d`).  Only valid CoTSs will be displayed, and any decoding or
validation error will be printed alongside the corresponding file name.

For example:
```
$ cocli cots display --file vendor.cbor
```
provided the `vendor.cbor` file contains valid CoTS, would print something like:
```
>> [vendor.cbor]
{
  "environments": [
    {
      "environment": {
        "class": {
          "vendor": "Zesty Hands, Inc."
        }
      }
    }
  ],
  "keys": {
    "tas": [
      {
        "format": 1,
        "data": "ooICejCCAnYwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAATjUaoQOSQHpL0DfKC8EVTQ5wHwZ085yyxPkhBpLOu+7B0nl33FYWV1Hg4je/37FTbpmohFkUKWYd81z8C/K1DMBBQBXEXJrLBGKnFd1xCgeMAVSfEBPzCCAgEwPjELMAkGA1UEBgwCVVMxEDAOBgNVBAoMB0V4YW1wbGUxHTAbBgNVBAMMFEV4YW1wbGUgVHJ1c3QgQW5jaG9yoIIBvTCCAWSgAwIBAgIVANCdkL89UlzHc9Ui7XfVniK7pFuIMAoGCCqGSM49BAMCMD4xCzAJBgNVBAYMAlVTMRAwDgYDVQQKDAdFeGFtcGxlMR0wGwYDVQQDDBRFeGFtcGxlIFRydXN0IEFuY2hvcjAeFw0yMjA1MTkxNTEzMDdaFw0zMjA1MTYxNTEzMDdaMD4xCzAJBgNVBAYMAlVTMRAwDgYDVQQKDAdFeGFtcGxlMR0wGwYDVQQDDBRFeGFtcGxlIFRydXN0IEFuY2hvcjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABONRqhA5JAekvQN8oLwRVNDnAfBnTznLLE+SEGks677sHSeXfcVhZXUeDiN7/fsVNumaiEWRQpZh3zXPwL8rUMyjPzA9MB0GA1UdDgQWBBQBXEXJrLBGKnFd1xCgeMAVSfEBPzALBgNVHQ8EBAMCAoQwDwYDVR0TAQH/BAUwAwEB/zAKBggqhkjOPQQDAgNHADBEAiALBidABsfpzG0lTL9Eh9b6AUbqnzF+koEZbgvppvvt9QIgVoE+bhEN0j6wSPzePjLrEdD+PEgyjHJ5rbA11SPq/1M="
      }
    ]
  }
}

```
While a `data/cots` folder with the following contents:
```
$ tree cots/
cots/
├── rubbish.cbor
├── namedtastore.cbor
├── vendor.cbor
```
could be inspected in one go using:
```
$ cocli cots display --dir data/cots/
```
which would output something like:
```
>> [data/cots/namedtastore.cbor]
{
  "environments": [
    {
      "namedtastore": "Miscellaneous TA Store"
    }
  ],
  "keys": {
    "tas": [
      {
        "format": 1,
        "data": "ooIC1TCCAtEwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAATN0f5kzywEzZOYbaV23O3N8cku39JoLNjlHPwECbXDDWp0LpAO1z248/hoy6UW/TZMTPPR/93XwHsG16mSFy8XBBSKhM/5gJWjvDbW7qUY1peNm9cfYDCCAlwwXDELMAkGA1UEBgwCVVMxHzAdBgNVBAoMFlNub2JiaXNoIEFwcGFyZWwsIEluYy4xLDAqBgNVBAMMI1Nub2JiaXNoIEFwcGFyZWwsIEluYy4gVHJ1c3QgQW5jaG9yoIIB+jCCAZ+gAwIBAgIUEBuTRGXAEEVEHhu4xafAnqm+qYgwCgYIKoZIzj0EAwIwXDELMAkGA1UEBgwCVVMxHzAdBgNVBAoMFlNub2JiaXNoIEFwcGFyZWwsIEluYy4xLDAqBgNVBAMMI1Nub2JiaXNoIEFwcGFyZWwsIEluYy4gVHJ1c3QgQW5jaG9yMB4XDTIyMDUxOTE1MTMwOFoXDTMyMDUxNjE1MTMwOFowXDELMAkGA1UEBgwCVVMxHzAdBgNVBAoMFlNub2JiaXNoIEFwcGFyZWwsIEluYy4xLDAqBgNVBAMMI1Nub2JiaXNoIEFwcGFyZWwsIEluYy4gVHJ1c3QgQW5jaG9yMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEzdH+ZM8sBM2TmG2ldtztzfHJLt/SaCzY5Rz8BAm1ww1qdC6QDtc9uPP4aMulFv02TEzz0f/d18B7BtepkhcvF6M/MD0wHQYDVR0OBBYEFIqEz/mAlaO8NtbupRjWl42b1x9gMAsGA1UdDwQEAwIChDAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49BAMCA0kAMEYCIQC2cf43f3PPlCO6/dxv40ftIgxxToKHF72UzENv7+y4ygIhAIGtC/r6SGaFMaP7zD2EloBuIXTtyWu8Hwl+YGdXRY93"
      }
    ]
  }
}
>> failed displaying "data/cots/rubbish.cbor": CBOR decoding failed: cbor: cannot unmarshal primitives into Go value of type cots.ConciseTaStore
>> [data/cots/vendor.cbor]
{
  "environments": [
    {
      "environment": {
        "class": {
          "vendor": "Zesty Hands, Inc."
        }
      }
    }
  ],
  "keys": {
    "tas": [
      {
        "format": 1,
        "data": "ooICejCCAnYwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAATjUaoQOSQHpL0DfKC8EVTQ5wHwZ085yyxPkhBpLOu+7B0nl33FYWV1Hg4je/37FTbpmohFkUKWYd81z8C/K1DMBBQBXEXJrLBGKnFd1xCgeMAVSfEBPzCCAgEwPjELMAkGA1UEBgwCVVMxEDAOBgNVBAoMB0V4YW1wbGUxHTAbBgNVBAMMFEV4YW1wbGUgVHJ1c3QgQW5jaG9yoIIBvTCCAWSgAwIBAgIVANCdkL89UlzHc9Ui7XfVniK7pFuIMAoGCCqGSM49BAMCMD4xCzAJBgNVBAYMAlVTMRAwDgYDVQQKDAdFeGFtcGxlMR0wGwYDVQQDDBRFeGFtcGxlIFRydXN0IEFuY2hvcjAeFw0yMjA1MTkxNTEzMDdaFw0zMjA1MTYxNTEzMDdaMD4xCzAJBgNVBAYMAlVTMRAwDgYDVQQKDAdFeGFtcGxlMR0wGwYDVQQDDBRFeGFtcGxlIFRydXN0IEFuY2hvcjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABONRqhA5JAekvQN8oLwRVNDnAfBnTznLLE+SEGks677sHSeXfcVhZXUeDiN7/fsVNumaiEWRQpZh3zXPwL8rUMyjPzA9MB0GA1UdDgQWBBQBXEXJrLBGKnFd1xCgeMAVSfEBPzALBgNVHQ8EBAMCAoQwDwYDVR0TAQH/BAUwAwEB/zAKBggqhkjOPQQDAgNHADBEAiALBidABsfpzG0lTL9Eh9b6AUbqnzF+koEZbgvppvvt9QIgVoE+bhEN0j6wSPzePjLrEdD+PEgyjHJ5rbA11SPq/1M="
      }
    ]
  }
}

Note: One of more files and directories can be supplied in the same invocation, using -f and -d directive:

```
