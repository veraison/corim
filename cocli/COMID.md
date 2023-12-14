
# CoMIDs manipulation

The `comid` subcommand allows you to create, display and validate CoMIDs.

## Create

Use the `comid create` subcommand to create a CBOR-encoded CoMID, passing its
JSON representation via the `--template` switch (or equivalently its `-t` shorthand):

* Please inspect `comid` JSON templates as examples under `data/comid/templates` `comid-*.json`

```
$ cocli comid create --template data/comid/templates/comid-dice-refval.json
```
On success, you should see something like the following printed to stdout:
```
>> created "comid-dice-refval.cbor" from "comid-dice-refval.json"
```

The CBOR-encoded CoMID file is stored in the current working directory with a
name derived from its template.  If you want, you can specify a different
target directory using the `--output-dir` command line switch (abbrev. `-o`)
```
$ cocli comid create --template data/comid/templates/comid-dice-refval.json --output-dir /tmp
>> created "/tmp/comid-dice-refval.cbor" from "comid-dice-refval.json"
```
Note that the output directory, as well as all its parent directories, MUST
pre-exist.

You can also create multiple CoMIDs in one go.  Suppose all your templates are
stored in the `templates/` folder:
```
$ tree templates/
templates/
├── comid-dice-refval1.json
├── comid-dice-refval2.json
...
└── comid-dice-refvaln.json
```
Then, you can use the `--template-dir` (abbrev. `-T`), and let the tool load,
validate, and CBOR-encode the templates one by one:
```
$ cocli comid create --template-dir templates
>> created "comid-dice-refval1.cbor" from "templates/comid-dice-refval1.json"
>> created "comid-dice-refval2.cbor" from "templates/comid-dice-refval2.json"
...
>> created "comid-dice-refvaln.cbor" from "templates/comid-dice-refvaln.json"
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


## Display

Use the `comid display` subcommand to print to stdout one or more CBOR-encoded
CoMIDs in human readable (JSON) format.

You can supply individual files using the `--file` switch (abbrev. `-f`), or
directories that may (or may not) contain CoMID files using the `--dir` switch
(abbrev. `-d`).  Only valid CoMIDs will be displayed, and any decoding or
validation error will be printed alongside the corresponding file name.

For example:
```
$ cocli comid display --file data/comid/comid-dice-refval.cbor
```
provided the `comid-dice-refval.cbor` file contains valid CoMID, would print something like:
```
>> [comid-dice-refval.cbor]
{
  "tag-identity": {
    "id": "1d5a8c7c-1c70-4c56-937e-3c5713ae5a83"
  },
  "triples": {}
[...]
}
```
While a `data/comid/` folder with the following contents:
```
$ tree data/comid/
data/comid/
├── rubbish.cbor
├── 1.cbor
└── 2.cbor
```
could be inspected in one go using:
```
$ cocli comid display --dir data/comid/
```
which would output something like:
```
>> failed displaying "comids.d/rubbish.cbor": CBOR decoding failed: EOF
>> [data/comid/1.cbor]
{
  "tag-identity": {
    "id": "43bbe37f-2e61-4b33-aed3-53cff1428b16"
  },
[...]
}
>> [data/comid/2.cbor]
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