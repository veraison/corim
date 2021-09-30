# corim
Concise Reference Integrity Manifest (CoRIM) package provides a golang API for manipulating 
CoRIM and Concise Module Identifier (CoMID) as per [Concise Reference Integrity Manifest](https://datatracker.ietf.org/doc/draft-birkholz-rats-corim/)

Specifically library supports following functions 
* To instantiate a CoMID, set desired values within a CoMID
* CBOR Encoding/Decoding from/to a CoMID
* A user friendly interface to populate a CoMID using a JSON byte stream
* Facility to add multiple CoMIDs and/or multiple CoSWID's as an array, in an Unsigned CoRIM
* CBOR Encoding/Decoding from/To an Unsigned CoRIM
* A user friendly interface to populate an unsigned CoRIM using a JSON byte stream
* Take an Unsigned CoRIM, sign it using a supplied COSE signer to generate a Signed CoRIM Message
* Verify a Signed CoRIM Message using a supplied public key
* Decode a signed COSE buffer (containing a Signed CoRIM) to provide an Unsigned CoRIM structure