# Corim Command Line Interface

## Installing and configuring

To install the `cocli` command, do:
```
$ go install github.com/veraison/corim/cocli
```

To configure auto-completion, use the `completion` subcommand.  For example, if
`bash` is your shell, you would do something like:
```
$ cocli completion bash > ~/.bash_completion.d/cocli
$ . .bash_completion
```
to get automatic command completion and suggestions using the TAB key.

To get a list of the supported shells, do:
```
$ cocli completion --help
```

## CoMIDs manipulation

The `comid` subcommand allows you to create, display and validate CoMIDs.

### Create

Use the `comid create` subcommand to create a CBOR-encoded CoMID, passing its
JSON representation<sup>[1](#templates-ex)</sup> via the `--template` switch (or
equivalently its `-t` shorthand):
```
$ cocli comid create --template t1.json
```
On success, you should see something like the following printed to stdout:
```
>> created "t1.cbor" from "t1.json"
```

The CBOR-encoded CoMID file is stored in the current working directory with a
name derived from its template.  If you want, you can specify a different
target directory using the `--output-dir` command line switch (abbrev. `-o`)
```
$ cocli comid create --template t1.json --output-dir /tmp
>> created "/tmp/t1.cbor" from "t1.json"
```
Note that the output directory, as well as all its parent directories, MUST
pre-exist.

You can also create multiple CoMIDs in one go.  Suppose all your templates are
stored in the `templates/` folder:
```
$ tree templates/
templates/
├── t1.json
├── t2.json
...
└── tn.json
```
Then, you can use the `--template-dir` (abbrev. `-T`), and let the tool load,
validate, and CBOR-encode the templates one by one:
```
$ cocli comid create --template-dir templates
>> created "t1.cbor" from "templates/t1.json"
>> created "t2.cbor" from "templates/t2.json"
...
>> created "tn.cbor" from "templates/tn.json"
```

You can specify both the `-T` and `-t` switches as many times as needed, and
even combine them in one invocation:
```
$ cocli comid create -T comid-templates/ \
                   -T comid-templates-aux/ \
                   -t extra-comid.json \
                   -t yet-another-comid.json \
                   -o /var/spool/comid
```

**NOTE** that since the output file name is deterministically generated from the
template file name, all the template files (when from different directories)
MUST have different base names.


### Display

Use the `comid display` subcommand to print to stdout one or more CBOR-encoded
CoMIDs in human readable (JSON) format.

You can supply individual files using the `--file` switch (abbrev. `-f`), or
directories that may (or may not) contain CoMID files using the `--dir` switch
(abbrev. `-d`).  Only valid CoMIDs will be displayed, and any decoding or
validation error will be printed alongside the corresponding file name.

For example:
```
$ cocli comid display --file m1.cbor
```
provided the `m1.cbor` file contains valid CoMID, would print something like:
```
>> [m1.cbor]
{
  "lang": "en-GB",
  "tag-identity": {
    "id": "43bbe37f-2e61-4b33-aed3-53cff1428b16"
  },
  "entities": [
    {
      "name": "ACME Ltd.",
      "regid": "https://acme.example",
      "roles": [
        "tagCreator",
        "creator",
        "maintainer"
      ]
    }
[...]
```
While a `comids.d` folder with the following contents:
```
$ tree comids.d/
comids.d/
├── rubbish.cbor
├── valid-comid-1.cbor
└── valid-comid-2.cbor
```
could be inspected in one go using:
```
$ cocli comid display --dir comids.d/
```
which would output something like:
```
>> failed displaying "comids.d/rubbish.cbor": CBOR decoding failed: EOF
>> [comids.d/valid-comid-1.cbor]
{
  "tag-identity": {
    "id": "43bbe37f-2e61-4b33-aed3-53cff1428b16"
  },
[...]
}
>> [comids.d/valid-comid-2.cbor]
{
  "tag-identity": {
    "id": "366d0a0a-5988-45ed-8488-2f2a544f6242"
  },
[...]
}
Error: 1/3 display(s) failed
```

One of more files and directories can be supplied in the same invocation, e.g.:
```
$ cocli comid display -f m1.cbor \
                    -f comids.d/m2.cbor \
                    -d /var/spool/comids \
                    -d yet-another-comid-folder/
```

## CoRIMs manipulation

The `corim` subcommand allows you to create, display, sign and verify CoRIMs.
It also provides a means to extract as-is the embedded CoSWIDs and CoMIDs and save
them as separate files.

### Create

Use the `corim create` subcommand to create a CBOR-encoded, unsigned CoRIM, by
passing its JSON representation<sup>[1](#templates-ex)</sup> via the
`--template` switch (or equivalently its `-t` shorthand) together with the
CBOR-encoded CoMIDs and/or CoSWIDs to be embedded.  For example:
```
$ cocli corim create --template c1.json --comid m1.cbor --coswid s1.cbor
```
On success, you should see something like the following printed to stdout:
```
>> created "c1.cbor" from "c1.json"
```

The CBOR-encoded CoRIM file is stored in the current working directory with a
name derived from its template.  If you want, you can specify a different
file name using the `--output` command line switch (abbrev. `-o`):
```
$ cocli corim create -t c1.json -m m1.cbor -s s1.cbor -o my.cbor
>> created "my.cbor" from "c1.json"
```

CoMIDs and CoSWIDs can be either supplied as individual files, using the
`--comid` (abbrev. `-m`) and `--coswid` (abbrev. `-s`) switches respectively, or
as "per-folder" blocks using the `--comid-dir` (abbrev. `-M`) and `--coswid-dir`
(abbrev. `-S`) switch.  For example:
```
$ cocli corim create --template c1.json --comid-dir comids.d/
```

Creation will fail if *any* of the inputs is non conformant.  For example, if
`comids.d` contains an invalid CoMID file `rubbish.cbor`, an attempt to create a
CoRIM:
```
$ cocli corim create -t c1.json -M comids.d/
```
will fail with:
```
Error: error loading CoMID from comids.d/rubbish.cbor: EOF
```

### Sign

Use the `corim sign` subcommand to cryptographically seal the unsigned CoRIM
supplied via the `--file` switch (abbrev. `-f`).  The signature is produced
using the key supplied via the `--key` switch (abbrev. `-k`), which is expected
to be in [JWK](https://www.rfc-editor.org/rfc/rfc7517) format.  On success, the
resulting COSE Sign1 payload is saved to file whose name can be controlled using
the `--output` switch (abbrev. `-o`).  A CoRIM Meta<sup>[1](#templates-ex)</sup>
template in JSON format must also be provided using the `--meta` switch (abbrev.
`-m`).  For example, with the default output file:
```
$ cocli corim sign --file corim.cbor --key ec-p256.jwk --meta meta.json
>> "corim.cbor" signed and saved to "signed-corim.cbor"
```
Or, the same but with a custom output file:
```
$ cocli corim sign --file corim.cbor \
                 --key ec-p256.jwk \
                 --meta meta.json \
                 --output /var/spool/signed-corim.cbor
>> "corim.cbor" signed and saved to "/var/spool/signed-corim.cbor"
```

### Verify

Use the `corim verify` subcommand to cryptographically verify the signed CoRIM
supplied via the `--file` switch (abbrev. `-f`).  The signature is checked
using the key supplied via the `--key` switch (abbrev. `-k`), which is expected
to be in [JWK](https://www.rfc-editor.org/rfc/rfc7517) format.  For example:
```
$ cocli corim verify --file signed-corim.cbor --key ec-p256.jwk
>> "corim.cbor" verified
```

Verification can fail either because the cryptographic processing fails or
because the signed payload or protected headers are themselves invalid.  For example:
```
$ cocli corim verify --file signed-corim-bad-signature.cbor --key ec-p256.jwk
```
will give
```
Error: error verifying signed-corim-bad-signature.cbor with key ec-p256.jwk: verification failed ecdsa.Verify
```

### Display

Use the `corim display` subcommand to print to stdout a signed CoRIM in human
readable (JSON) format.

You must supply the file you want to display using the `--file` switch (abbrev.
`-f`).  Only a valid CoRIM will be displayed, and any occurring decoding or
validation errors will be printed instead.

The output has two logical sections: one for Meta and one for the (unsigned)
CoRIM:
```
$ cocli corim display --file signed-corim.cbor
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

By default, the embedded CoMID and CoSWID tags are not expanded, and what you
will see is the base64 encoding of their CBOR serialisation.  If you want to
peek at the tags' content, supply the `--show-tags` (abbrev. `-v`) switch, which
will add a further Tags section with one entry per each expanded tag:
```
$ cocli corim display --file signed-corim.cbor --show-tags
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

### Extract CoSWIDs and CoMIDs

Use the `corim extract` subcommand to extract the embedded CoMIDs and CoSWIDs
from a signed CoRIM.

You must supply a signed CoRIM file using the `--file` switch (abbrev. `-f`) and
an optional output folder (default is the current working directory) using the
`--output-dir` switch (abbrev. `-o`).  Make sure that the output directory as
well as any parent folder exists prior to issuing the command.

On success, the found CoMIDs and CoSWIDs are saved in CBOR format:
```
$ cocli corim extract --file signed-corim.cbor --output-dir output.d/
$ tree output.d/
output.d/
├── 000000-comid.cbor
├── 000001-comid.cbor
└── 000002-coswid.cbor
```



<a name="templates-ex">1</a>: A few examples of CoMID, CoRIM, and Meta JSON
templates can be found in the [data/templates](data/templates) folder.