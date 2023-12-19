# Corim Command Line Interface

## Installing and configuring

To install the `cocli` command, do:
```
$ go install github.com/veraison/corim/cocli@latest
```

To configure auto-completion, use the `completion` subcommand.  For example, if
`bash` is your shell, you would do something like:
```
$ cocli completion bash > ~/.bash_completion.d/cocli
$ . ~/.bash_completion
```
to get automatic command completion and suggestions using the TAB key.

To get a list of the supported shells, do:
```
$ cocli completion --help
```
# CoRIM Handling
This document provides step-by-step instructions for how to use the `cocli` tool to manipulate CoRIMs, CoMIDs and CoTS.

``` mermaid
flowchart TD
    subgraph CORIM["<b>CoRIM</b>"]
      subgraph CoMID["\n"]
        CM1["CoMID-1"]
        CM2["CoMID-2"]

        CM3["CoMID-N"]
        CM4["<b>COMID COMMANDS</b> \n cocli comid create \n cocli comid display"]
        CM1  -.- CM2
        CM2  -.- CM3
        CM3  -.- CM4
    end
    subgraph CoMID["Blank1"]
        CSW1["CoSWID-1"]
        CSW2["CoSWID-2"]
        CSW3["CoSWID-N"]

        CSW1  -.- CSW2
        CSW2  -.- CSW3
   
    end

    subgraph CoMID["Blank3"]
        CS1["CoTS-1"]
        CS2["CoTS-2"]
       
        CS3["CoTS-N"]
        CS4["<b>COTS COMMANDS</b> \n cocli cots create \n cocli cots display"]
        CS1  -.- CS2

        CS2  -.- CS3
        CS3 -.- CS4
    end
end
CORIM ---> CMD
subgraph CMD["<b>CORIM COMMANDS</b> \n
 1.cocli corim create \n 2.cocli corim display \n 3.cocli corim sign \n4.cocli corim verify\n5.cocli corim extract\n 6.cocli corim submit"]
end

```

## CoMIDs manipulation
The instructions to manipulate CoMIDs are documented [here](COMID.md)

## CoTSs manipulation
The instructions to manipulate CoTSs are documented [here](COTS.md)

## CoSWID manipulation
Tooling to manipulate `CoSWID` is not currently available under Project Veraison.
However CoSWID can be part of CoRIM by constructing CoSWID CBOR by other indistry available
tools such as [swid-tools](https://github.com/usnistgov/swid-tools) and including them
as mentioned under [CORIM Construction](CORIM.md)

## CoRIMs manipulation
The instructions to manipulate CoRIMs are documented [here](CORIM.md)

## CoRIM Submission to Veraison

Use the `corim submit` subcommand to upload a CoRIM using the Veraison provisioning API.
The CoRIM file containing the CoRIM data in CBOR format is supplied via the
`--corim-file` switch (abbrev. `-f`). The server URL where to upload the CoRIM
payload is supplied via the `--api-server` switch (abbrev. `-s`).
Further, it is required to supply the media type of the content via the
`--media-type` switch (abbrev. `-m`)
```
$ cocli corim submit \
    --corim-file data/corim/unsigned-corim.cbor \
    --api-server "https://veraison.example/endorsement-provisioning/v1/submit" \
    --media-type "application/corim-unsigned+cbor; profile=http://arm.com/psa/iot/1"

>> "unsigned-corim.cbor" submit ok
```

#### Remote Service Authentication

The above will work if the remote service does not authenticate
endorsement-provisioning API calls. If the service does authenticate, then
cocli must be configured appropriately. This can be done using a `config.yaml`
file located in the current working directory, or in the standard config
path (usually `~/.config/cocli/config.yaml` on XDG-compliant systems). Please
see `./data/config/example-config.yaml` file for details of the configuration
that needs to be provided.

## Visual Synopsis of the Available Commands

```mermaid
graph LR
    OEM[(OEM/ODM \n DB)]
    JSONTmplCoMID[["JSON \n template \n (CoMID)"]]

    JSONTmplCoSWID[["JSON \n template \n (CoSWID)"]]
    style JSONTmplCoSWID fill:#71797E
    click JSONTmplCoSWID "https://github.com/veraison/corim/issues/81"

    JSONTmplCoRIM[["JSON \n template \n (CoRIM)"]]
    JSONTmplMeta[["JSON \n template \n (Meta)"]]
    key((key))

    %% Cots nodes
    environments[["Environments"]]
    tas(("Trust \n anchors"))
    cas(("CA \n certificates"))
    permClaims[["Permanant claims"]]
    exclClaims[["Excluded claims"]]

    cliComidCreate($ cocli comid create)
    cliComidDisplay($ cocli comid display)
    style cliComidCreate fill:#00758f
    style cliComidDisplay fill:#00758f

    cliCotsCreate($ cocli cots create)
    cliCotsDisplay($ cocli cots display)
    style cliCotsCreate fill:#00758f
    style cliCotsDisplay fill:#00758f

    cliCoswidCreate($ cocli coswid create)
    cliCoswidDisplay($ cocli coswid display)
    style cliCoswidCreate fill:#71797E
    style cliCoswidDisplay fill:#71797E


    cliCorimCreate($ cocli corim create)
    cliCorimSign($ cocli corim sign)
    cliCorimVerify($ cocli corim verify)
    cliCorimExtract($ cocli corim extract)
    cliCorimDisplay($ cocli corim display)
    cliCorimSubmit($ cocli corim submit)
    style cliCorimCreate fill:#00758f
    style cliCorimSign fill:#00758f
    style cliCorimVerify fill:#00758f
    style cliCorimExtract fill:#00758f
    style cliCorimDisplay fill:#00758f
    style cliCorimSubmit fill:#00758f

    provisioningEndpoint{{Veraison \n Provisioning \n Service}}

    CBORComid1((CBOR <br /> CoMID))
    CBORSwid1((CBOR <br /> SWID))
    CBORCots1((CBOR <br /> CoTS))

    CBORComid2((CBOR <br /> CoMID))
    CBORSwid2((CBOR <br /> SWID))
    CBORCots2((CBOR <br /> CoTS))

    CBORCorim((CBOR CoRIM))
    CoseSign1((COSE Sign1 CoRIM))
    signBool((T/F))

    OEM --> JSONTmplCoMID
    OEM --> JSONTmplCoSWID

    %% Cots items provisioning
    OEM --> environments
    OEM --> tas
    OEM --> cas
    OEM --> permClaims
    OEM --> exclClaims

    OEM --> JSONTmplCoRIM
    OEM --> JSONTmplMeta
    OEM --> key

    %% Cots individual items
    environments --> cliCotsCreate
    tas --> cliCotsCreate
    cas --> cliCotsCreate
    permClaims --> cliCotsCreate
    exclClaims --> cliCotsCreate


    JSONTmplCoMID --> cliComidCreate
    JSONTmplCoSWID --> cliCoswidCreate
    JSONTmplCoRIM --> cliCorimCreate
    JSONTmplMeta --> cliCorimSign
    key --> cliCorimSign
    key --> cliCorimVerify

    cliComidCreate --> CBORComid1
    cliCotsCreate --> CBORCots1
    cliCoswidCreate --> CBORSwid1

    cliCorimCreate --> CBORCorim
    cliCorimSign --> CoseSign1
    cliCorimVerify --> signBool
    cliCorimSubmit -- to--> provisioningEndpoint

    CBORComid1 --> cliComidDisplay
    CBORComid1 --> cliCorimCreate

    CBORCots1 --> cliCorimCreate
    CBORCots1  --> cliCotsDisplay

    CBORSwid1 --> cliCoswidDisplay
    CBORSwid1 --> cliCorimCreate

    CBORCorim --> cliCorimSubmit
    CBORCorim --> cliCorimSign
    CoseSign1 --> cliCorimExtract
    CoseSign1 --> cliCorimVerify
    CoseSign1 --> cliCorimDisplay

    cliCorimExtract --> CBORComid2
    cliCorimExtract --> CBORSwid2
    cliCorimExtract --> CBORCots2
```
