# Concise Reference Integrity Manifest and Module Identifiers
[![ci](https://github.com/veraison/corim/actions/workflows/ci.yml/badge.svg)](https://github.com/veraison/corim/actions/workflows/ci.yml)
[![cover ≥82%](https://github.com/veraison/corim/actions/workflows/ci-go-cover.yml/badge.svg)](https://github.com/veraison/corim/actions/workflows/ci-go-cover.yml)
[![linters](https://github.com/veraison/corim/actions/workflows/linters.yml/badge.svg)](https://github.com/veraison/corim/actions/workflows/linters.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/veraison/corim.svg)](https://pkg.go.dev/github.com/veraison/corim)


The [`corim/corim`](corim) and [`corim/comid`](comid) packages provide a golang API for low-level manipulation of [Concise Reference Integrity Manifest (CoRIM)](https://datatracker.ietf.org/doc/draft-ietf-rats-corim/) and Concise Module Identifier (CoMID) tags respectively.
The [`corim/coev`](coev) package provides a golang implementation of TCG Concise Evidence CDDL as documented [here](https://github.com/TrustedComputingGroup/dice-coev/blob/main/concise-evidence.cddl).
The [`corim/coserv`](coserv) package provides a golang API for working with [Concise Selector for Endorsements and Reference Values](https://datatracker.ietf.org/doc/draft-howard-rats-coserv).
The [`corim/cotl`](cotl) package provides a golang implementation of the Concise Tag List (CoTL) data structure defined in Section 6 of [draft-ietf-rats-corim-11](https://datatracker.ietf.org/doc/draft-ietf-rats-corim/).

> [!NOTE]
> These API are still in active development (as is the underlying CoRIM spec).
> They are **subject to change** in the future.

## Compatibility

This implementation is fully compliant with revision 11 of the specification
draft. There is partial support for features from subsequent revisions:

Feature | Revision | Supported?
--|--|--
scalar matchers (mval code points 16-20) | 12 | ✅
"binary" version scheme | 12 | ✅
multiple signatures | 12 | :x:

## Profile documents

The [TPM 2.0 CoRIM Reference Values Profile](docs/profiles/tpm-v1.md) is a
draft specification for TPM and vTPM reference measurements. It includes the
measurement mapping, verification rules, and a worked example.

## Required Tools

Ensure you have the following tools installed with the specified versions on your machine to ensure everything works properly:

- **Go**: Version 1.25

## Developer tips

Before requesting a PR (and routinely during the dev/test cycle), you are encouraged to run:
```
make presubmit
```
and check its output to make sure your code coverage figures are in line with the set target and that there are no newly introduced lint problems.

## x5chain trust verification

Signed CoRIM messages may carry an X.509 chain in the COSE `x5chain` protected header.
Use [`SignedCorim.VerifyWithX5Chain`](https://pkg.go.dev/github.com/veraison/corim/corim#SignedCorim.VerifyWithX5Chain)
after [`FromCOSE`](https://pkg.go.dev/github.com/veraison/corim/corim#SignedCorim.FromCOSE) (or after a successful [`Sign`](https://pkg.go.dev/github.com/veraison/corim/corim#SignedCorim.Sign)) to validate PKIX trust, CRL revocation, and the COSE signature.

PKIX/CRL behavior is implemented by [`github.com/veraison/go-cose`](https://pkg.go.dev/github.com/veraison/go-cose).
Configure trust with [`cose.TrustAnchors`](https://pkg.go.dev/github.com/veraison/go-cose#TrustAnchors) and optional [`cose.X5ChainVerifyOptions`](https://pkg.go.dev/github.com/veraison/go-cose#X5ChainVerifyOptions); load files with [`cose.LoadTrustAnchors`](https://pkg.go.dev/github.com/veraison/go-cose#LoadTrustAnchors).
See the go-cose docs for PEM/DER rules, CRL matching, and [`CRLPolicy`](https://pkg.go.dev/github.com/veraison/go-cose#CRLPolicy).

Defaults that commonly trip callers:
- Verification fails closed for a zero-value `TrustAnchors`, or one loaded with no anchor paths: there is no OS-store fallback. Set `Anchors` and/or `UseSystemRoots`.
- Default [`RevocationFullChain`](https://pkg.go.dev/github.com/veraison/go-cose#RevocationFullChain) requires CRLs; set [`RevocationDisabled`](https://pkg.go.dev/github.com/veraison/go-cose#RevocationDisabled) when revocation is not required.

The `x5chain` header must be leaf-first (RFC 9360 §2). `FromCOSE` and `Sign` reject other orders, including when the caller only uses [`Verify`](https://pkg.go.dev/github.com/veraison/corim/corim#SignedCorim.Verify) with an external key.
Use `errors.Is` with `ErrX5ChainNoTrust`, `ErrX5ChainCRLMissing`, `ErrX5ChainRevoked`, or `ErrX5ChainSignature` to distinguish failures.

For external-key verification without PKIX path validation, use [`SignedCorim.Verify`](https://pkg.go.dev/github.com/veraison/corim/corim#SignedCorim.Verify) instead.

## Extending CoRIM/CoMID

The CoRIM specification provides a mechanism for adding extensions to the base
CoRIM schema. The `corim` and `comid` structs which can be extended, embed an
`Extensions` object  that allows registering a wrapper structure defining
extension fields. For field types that can be extended, i.e. `type choice`,
extensions can be implemented by calling an appropriate registration function
and giving it a new type or a value (for enums).

Please see [extensions documentation](extensions/README.md) for details.


