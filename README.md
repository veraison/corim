# corim
Concise Reference Integrity Manifest (CoRIM) package provides a golang API for manipulating 
CoRIM and Concise Module Identifier (CoMID) as per [Concise Reference Integrity Manifest](https://datatracker.ietf.org/doc/draft-birkholz-rats-corim/)

Current Work, provides a facility to encode and decode a signed-corim message (i.e. a COSE Sign1 Wrapped CoRIM) as well as an Unsigned CoRIM that contains an array of CoMID tags.
