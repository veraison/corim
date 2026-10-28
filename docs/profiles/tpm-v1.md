# TPM 2.0 CoRIM Reference Values Profile

**Status:** Draft v1 for review · 2026-10-28

This document proposes a CoRIM profile for TPM 2.0 quotes, including quotes
from virtual TPMs. It defines a shared measurement mapping for reference-value
providers and verifiers.

## Scope and identification

Each CoMID reference triple describes one approved combination of TPM firmware
version, PCR selection, and aggregate PCR digest. Individual PCR values are
optional. This profile uses existing CoRIM types; its measurement-key assignments
are profile-specific and are not TPM or IANA allocations.

The design targets [CoRIM draft-11](https://datatracker.ietf.org/doc/html/draft-ietf-rats-corim-11).

The proposed publication location is:

```text
https://github.com/veraison/corim/blob/main/docs/profiles/tpm-v1.md
```

The CoRIM `profile` field MUST contain this URL as the profile identifier.
Verifiers MUST reject manifests with a missing or unsupported profile identifier.
The URL remains proposed until maintainer approval; once finalized, changes to
this profile MUST use a new identifier. See [CoRIM §4.1.4](https://datatracker.ietf.org/doc/html/draft-ietf-rats-corim-11#section-4.1.4).

## Environment

Each reference triple MUST identify the platform or component to which its
measurements apply. Deployments define their environment class identifiers and
how evidence is associated with them. An environment class identifies a family
of comparable systems, such as an OCI shape and OS release. Each class can have
multiple approved reference configurations, for example, to accommodate firmware
or kernel updates. Each approved configuration is represented by a separate
reference triple.

For example, an OCI deployment can use:

| Class field | Example |
| --- | --- |
| `vendor` | `Oracle Cloud Infrastructure` |
| `model` | `VM.Standard.E5.Flex` |

This OCI convention is illustrative, not required for other deployments.

## Reference measurements

Each triple MUST contain keys 0, 1, and 2 exactly once; key 3 MAY occur once.
Unknown or duplicate measurement keys MUST be rejected for this profile
version. Each measurement's `mval` MUST contain only the value field listed below.

| `mkey` (unsigned integer) | Claim | `mval` field | Encoding                              |
| --- | --- | --- |---------------------------------------|
| 0 | `firmwareVersion` | `raw-value` | Eight-byte unsigned big-endian value  |
| 1 | Complete `pcrSelect` | `raw-value` | Serialized `TPML_PCR_SELECTION`       |
| 2 | `pcrDigest` | `digests` | Exactly one aggregate hash algorithm/digest pair |
| 3 | Individual PCR values (optional) | `integrity-registers` | Unsigned PCR index → bank algorithm/digest pairs |

For keys 0 and 1, `raw-value` uses CoRIM's `tagged-bytes` representation
(CBOR tag 560); the bytes inside that wrapper use the encodings above.

**Firmware.** Key 0 identifies the TPM/vTPM implementation version, not the host
firmware, guest UEFI, or kernel. Comparison is exact equality. TPM attestation
can obfuscate this field depending on the signing-key hierarchy, so providers
must establish its comparability before publishing shared baselines. See the
[TPM attestation implementation](https://github.com/microsoft/ms-tpm-20-ref/blob/main/TPMCmd/tpm/src/command/Attestation/Attest_spt.c).

**Selection.** Key 1 includes the UINT32 count, followed by each entry's UINT16
bank algorithm, UINT8 bitmap length, and bitmap bytes. Multi-byte integers are
big-endian. Validate count, lengths, supported banks, and PCR indices; reject
empty selections, duplicate banks, and trailing bytes. Preserve selection-list
order and compare the complete encoding exactly.

### Understanding the PCR selection

A **PCR bank** is a set of PCRs using one hash algorithm. For example,
SHA-256 PCR 0 and SHA-384 PCR 0 are separate values, containing 32 and 48 bytes
respectively.
See [PCR banks](https://learn.microsoft.com/en-us/windows/security/hardware-security/tpm/switch-pcr-banks-on-tpm-2-0-devices)
and the [TPM platform profile](https://trustedcomputinggroup.org/wp-content/uploads/TCG-PC-Client-Platform-TPM-Profile-for-TPM-2.0-Version-1.06-Revision-32_5April24.pdf).

`TPML_PCR_SELECTION` is a list of `TPMS_PCR_SELECTION` entries, one per selected
bank in this profile. The supplied quote has this structure:

```text
TPML_PCR_SELECTION
  count = 1
  pcrSelections[0]: TPMS_PCR_SELECTION
    hash         = 0x000b       (SHA-256 bank)
    sizeofSelect = 3            (bitmap bytes, not PCR-value bytes)
    pcrSelect    = ff 07 00
                   |  |  +---- PCRs 16–23: none selected
                   |  +------- PCRs 8–15: 8, 9, 10 selected
                   +---------- PCRs 0–7: all selected
```

| Field | Meaning in this example |
| --- | --- |
| `count` | One selection entry, not eleven PCRs or all banks supported by the TPM |
| `hash` | Hash algorithm identifying the bank |
| `sizeofSelect` | Three bitmap bytes can represent 24 PCR indices |
| `pcrSelect` | Eleven set bits select PCRs 0–10; bit `i % 8` of byte `i // 8` selects PCR `i` (`//` is integer division) |

The serialized structure is `00000001 000b 03 ff0700`. The TPM field name is
`sizeofSelect` (lowercase `of`). See the [TPM structure definitions](https://trustedcomputinggroup.org/wp-content/uploads/Trusted-Platform-Module-Library-Family-2.0-Level-00-Revision-1.81-Part-2-Structures_8December2023.pdf).

### Digest and optional PCRs across banks

**Aggregate digest.** Key 2 uses the hash algorithm from the quote's signing
scheme, which can differ from the PCR-bank algorithm. Please refer to the
[TCG algorithm registry](https://www.trustedcomputinggroup.org/resources/documents/tcg-algorithm-registry) for a list of all available algorithms. Translate TPM algorithm
identifiers into the selected CoRIM digest representation; their numeric values
are not interchangeable. See the [TPM quote implementation](https://github.com/microsoft/ms-tpm-20-ref/blob/main/TPMCmd/tpm/src/command/Attestation/Quote.c).

**Optional PCRs.** Key 3 MUST contain exactly the PCR/bank pairs selected by key
1, with one value per pair. Hash the raw PCR bytes in selection-list order,
then ascending PCR index within each bank, using key 2's algorithm. The result
MUST equal key 2's digest; otherwise reject the reference triple.

A quote contains **one aggregate `pcrDigest`**, even when `count` is greater
than one. For example, if the selection list is SHA-256 PCRs 0 and 1 followed
by SHA-384 PCRs 0 and 2, then `count = 2` and:

```text
pcrDigest = Hquote(
    PCR_SHA256[0] || PCR_SHA256[1] ||
    PCR_SHA384[0] || PCR_SHA384[2]
)
```

There is no intermediate per-bank hash in this construction. `Hquote` is the
signing scheme's hash algorithm. See [PCRComputeCurrentDigest](https://github.com/microsoft/ms-tpm-20-ref/blob/main/TPMCmd/tpm/src/subsystem/PCR.c).

One `integrity-registers` map can represent these bank-specific values because
each index maps to algorithm/digest pairs. For this example, the optional key 3
would contain the following schematic values (placeholders, not literal CBOR):

```text
integrity-registers:
  0: [[sha-256, <32-byte value>], [sha-384, <48-byte value>]]
  1: [[sha-256, <32-byte value>]]
  2: [[sha-384, <48-byte value>]]
```

Here the algorithm identifies the bank; no separate map or aggregate digest is
needed per bank. The selection and consistency rules above require all selected
PCR/bank pairs, regardless of the map's ordering.

## Quote appraisal

The verifier MUST:

1. Authenticate the reference manifest and ensure the baseline is currently
   approved for the target workload and environment.
2. Verify the quote using a trusted attestation key, require
   `magic = 0xff544347` and `type = 0x8018`, and validate `extraData` against the
   expected session challenge/nonce or protocol-defined challenge binding.
3. Enforce policy for permitted algorithms and required PCR coverage.
4. Match firmware version, encoded PCR selection, and aggregate algorithm/digest
   against the same approved reference triple. Never mix claims from different
   triples.

Separate PCR values from the attester are unnecessary for aggregate matching.
`qualifiedSigner` and `clockInfo` remain evidence for identity and runtime-policy
checks; neither is a shared reference measurement in this profile.

## Example

The following CoMID, in CBOR diagnostic notation, illustrates the three required
measurements. The shape is illustrative and SHA-256 is assumed for the aggregate
hash; neither is established by the supplied quote text. No individual PCR
values are available, so key 3 is absent.

```cbor-diag
{
  / tag-identity / 1: { 0: "example-tpm-baseline-1" },
  / triples / 4: {
    / reference-triples / 0: [[
      / environment / {
        / class / 0: {
          / vendor / 1: "Oracle Cloud Infrastructure",
          / model / 2: "VM.Standard.E5.Flex"
        }
      },
      / measurements / [
        { 0: 0, 1: { 4: 560(h'0000000020032520') } },
        { 0: 1, 1: { 4: 560(h'00000001000b03ff0700') } },
        { 0: 2, 1: { 2: [["sha-256",
          h'3fb83052a48a4289014c8407976ef311d9049eae74a2b58cbcd83392b5ee432c'
        ]] } }
      ]
    ]]
  }
}
```

The selection bytes decode as `00000001` (one bank), `000b` (SHA-256), `03`
(three bitmap bytes), and `ff0700` (PCRs 0–10). Embed the encoded CoMID in a
CoRIM carrying the profile identifier. CBOR is the interchange format;
library-specific JSON is an authoring/display format, not the normative schema.
The example is not an approved baseline.
