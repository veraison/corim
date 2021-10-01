# corim
Concise Reference Integrity Manifest (CoRIM) package provides a golang API for manipulating 
CoRIM and Concise Module Identifier (CoMID) as per [Concise Reference Integrity Manifest](https://datatracker.ietf.org/doc/draft-birkholz-rats-corim/)

Specifically, the library supports following functions:
* APIs to set and get individual fields within a CoMID
* CBOR encoding/decoding to/from a CoMID
* A user friendly interface to populate a CoMID using an equivalent JSON representation
* A facility to add multiple CoMIDs and/or multiple CoSWIDs to an Unsigned CoRIM
* CBOR encoding/decoding to/from an Unsigned CoRIM
* A user friendly interface to populate an unsigned CoRIM using an equivalent JSON representation
* Signing an unsigned CoRIM with a private key to obtain a signed CoRIM message
* Verifying a signed CoRIM using a public key
* Extracting an unsigned CoRIM and CoRIM Meta structures from a serialized signed CoRIM