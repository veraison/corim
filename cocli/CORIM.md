# CoRIMs manipulation

The `corim` subcommand allows you to create, display, sign, verify CoRIMs or submit
a CoRIM using the [Veraison provisioning API](https://github.com/veraison/docs/tree/main/api/endorsement-provisioning).
It also provides a means to extract as-is the embedded CoSWIDs, CoMIDs and CoTSs and save
them as separate files.

## Create

Use the `corim create` subcommand to create a CBOR-encoded, unsigned CoRIM, by
passing its JSON representation via the `--template` switch (or equivalently its `-t` shorthand)
together with the CBOR-encoded CoMIDs, CoSWIDs and/or CoTS to be embedded.

* Please inspect `corim` JSON templates as examples under `data/corim/templates` `corim-*.json`

```
$ cocli corim create --template data/corim/templates/corim-full.json --comid data/comid/comid-dice-refval.cbor --coswid data/coswid/1.cbor --cots data/cots/vendor.cbor
```
On success, you should see something like the following printed to stdout:
```
>> created "corim-full.cbor" from "corim-full.json"
```

The CBOR-encoded CoRIM file is stored in the current working directory with a
name derived from its template.  If you want, you can specify a different
file name using the `--output` command line switch (abbrev. `-o`):
```
$ cocli corim create -t data/corim/templates/corim-full.json -m data/comid/comid-dice-refval.cbor -s data/coswid/1.cbor -c data/cots/c1.cbor -o unsigned-corim.cbor
>> created "unsigned-corim.cbor" from "corim-full.json"
```

CoMIDs, CoSWIDs and CoTSs can be either supplied as individual files, using the
`--comid` (abbrev. `-m`), `--coswid` (abbrev. `-s`) and `--cots` (abbrev. `-c`) switches respectively, or
as "per-folder" blocks using the `--comid-dir` (abbrev. `-M`), `--coswid-dir` and `--cots-dir`
(abbrev. `-C`) switch.  For example:
```
$ cocli corim create --template data/corim/templates/corim-full.json --comid-dir data/comid/cbor/
```

Creation will fail if *any* of the inputs is non conformant.  For example, if
`data/comid/cbor/` contains an invalid CoMID file `rubbish.cbor`, an attempt to create a
CoRIM:
```
$ cocli corim create -t data/corim/templates/corim-full.json -M data/comid/cbor/
```
will fail with:
```
Error: error loading CoMID from data/comid/cbor/rubbish.cbor: EOF
```

## Sign

Use the `corim sign` subcommand to cryptographically seal the unsigned CoRIM
supplied via the `--file` switch (abbrev. `-f`).  The signature is produced
using the key supplied via the `--key` switch (abbrev. `-k`), which is expected
to be in [JWK](https://www.rfc-editor.org/rfc/rfc7517) format.  On success, the
resulting COSE Sign1 payload is saved to file whose name can be controlled using
the `--output` switch (abbrev. `-o`).  A CoRIM Meta template in JSON format must 
also be provided using the `--meta` switch (abbrev.`-m`).

* Please inspect the `data/corim/templates` directory for `meta` JSON templates.

For example, with the default output file:
```
$ cocli corim sign --file corim.cbor --key ec-p256.jwk --meta meta.json
>> "corim.cbor" signed and saved to "signed-corim.cbor"
```
Or, the same but with a custom output file:
```
$ cocli corim sign --file data/corim/corim-full.cbor \
                 --key data/keys/ec-p256.jwk \
                 --meta data/corim/templates/meta-full.json \
                 --output /var/spool/signed-corim.cbor
>> "corim-full.cbor" signed and saved to "/var/spool/signed-corim.cbor"
```

## Verify

Use the `corim verify` subcommand to cryptographically verify the signed CoRIM
supplied via the `--file` switch (abbrev. `-f`).  The signature is checked
using the key supplied via the `--key` switch (abbrev. `-k`), which is expected
to be in [JWK](https://www.rfc-editor.org/rfc/rfc7517) format.  For example:
```
$ cocli corim verify --file data/corim/signed-corim.cbor --key data/keys/ec-p256.jwk
>> "signed-corim.cbor" verified
```

Verification can fail either because the cryptographic processing fails or
because the signed payload or protected headers are themselves invalid.  For example:
```
$ cocli corim verify --file data/corim/signed-corim-bad-signature.cbor --key data/keys/ec-p256.jwk
```
will give
```
Error: error verifying signed-corim-bad-signature.cbor with key ec-p256.jwk: verification failed ecdsa.Verify
```

## Display

Use the `corim display` subcommand to print to stdout a signed CoRIM in human
readable (JSON) format.

You must supply the file you want to display using the `--file` switch (abbrev.
`-f`).  Only a valid CoRIM will be displayed, and any occurring decoding or
validation errors will be printed instead.

The output has two logical sections: one for Meta and one for the (unsigned)
CoRIM:
```
$ cocli corim display --file data/corim/signed-corim.cbor
Meta:
{
  "signer": {
    "name": "ACME Ltd signing key",
    "uri": "https://acme.example/signing-key.pub"
  },
[...]
}
Corim:
{
  "corim-id": "5c57e8f4-46cd-421b-91c9-08cf93e13cfc",
  "tags": [
    "2QH...",
[...]
  ]
}
```

By default, the embedded CoMID, CoSWID and CoTS tags are not expanded, and what you
will see is the base64 encoding of their CBOR serialisation.  If you want to
peek at the tags' content, supply the `--show-tags` (abbrev. `-v`) switch, which
will add a further Tags section with one entry per each expanded tag:
```
$ cocli corim display --file data/corim/signed-corim.cbor --show-tags
Meta:
{
[...]
}
Corim:
{
[...]
}
Tags:
>> [ 0 ]
{
  "tag-identity": {
    "id": "366d0a0a-5988-45ed-8488-2f2a544f6242"
  },
[...]
}
>> [ 1 ]
{
  "tag-identity": {
    "id": "43bbe37f-2e61-4b33-aed3-53cff1428b16"
  },
[...]
}
>> [ 2 ]
{
  "tag-id": "com.acme.rrd2013-ce-sp1-v4-1-5-0",
[...]
}
```

## Extract CoSWIDs, CoMIDs and CoTSs

Use the `corim extract` subcommand to extract the embedded CoMIDs, CoSWIDs and CoTSs
from a signed CoRIM.

You must supply a signed CoRIM file using the `--file` switch (abbrev. `-f`) and
an optional output folder (default is the current working directory) using the
`--output-dir` switch (abbrev. `-o`).  Make sure that the output directory as
well as any parent folder exists prior to issuing the command.

On success, the found CoMIDs, CoSWIDs, CoTS are saved in CBOR format:
```
$ cocli corim extract --file data/corim/signed-corim.cbor --output-dir output.d/
$ tree output.d/
output.d/
├── 000000-comid.cbor
├── 000001-comid.cbor
├── 000002-coswid.cbor
└── 000003-cots.cbor
```