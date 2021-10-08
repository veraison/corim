# Corim Command Line Interface

## CoMIDs manipulation

The `comid` subcommand allows you to create, display and validate CoMIDs.

### Create

Use the `comid create` subcommand to create a CBOR-encoded CoMID, passing its
JSON representation<sup>[1](#comid-templates-ex)</sup> via the `--template`
switch (or equivalently its `-t` shorthand):
```
$ cli comid create --template t1.json
```
On success, you should see something like the following printed to stdout:
```
>> created "t1.cbor" from "t1.json"
```

The CBOR-encoded CoMID file is stored in the current working directory with a
name derived from its template.  If you want, you can specify a different
target directory using the `--output-dir` command line switch (abbrev. `-o`)
```
$ cli comid create --template t1.json --output-dir /tmp
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
$ cli comid create --template-dir templates
>> created "t1.cbor" from "templates/t1.json"
>> created "t2.cbor" from "templates/t2.json"
...
>> created "tn.cbor" from "templates/tn.json"
```

You can specify both the `-T` and `-t` switches as many times as needed, and
even combine them in one invocation:
```
$ cli comid create -T comid-templates/ \
                   -T comid-templates-aux/ \
                   -t extra-comid.json \
                   -t yet-another-comid.json \
                   -o /var/spool/comid
```

**NOTE** that since the output file name is deterministically generated from the
template file name, all the template files (when from different directories)
MUST have different base names.

<a name="comid-templates-ex">1</a>: A few examples of CoMID JSON templates,
can be found in the [data/templates](data/templates) folder.

### Display

Use the `comid display` subcommand to print to stdout one or more CBOR-encoded
CoMIDs in human readable JSON format.

You can supply individual files using the `--file` switch (abbrev. `-f`), or
directories that may (or may not) contain CoMID files using the `--dir` switch
(abbrev. `-d`).  Only valid CoMIDs will be displayed, and any decoding or
validation error will be printed alongside the corresponding file name.

For example:
```
$ cli comid display --file m1.cbor
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
$ cli comid display --dir comids.d/
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
$ cli comid display -f m1.cbor \
                    -f comids.d/m2.cbor \
                    -d /var/spool/comids \
                    -d yet-another-comid-folder/
```