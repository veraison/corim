# Concise Reference Integrity Manifest and Module Identifiers
[![ci](https://github.com/veraison/corim/actions/workflows/ci.yml/badge.svg)](https://github.com/veraison/corim/actions/workflows/ci.yml)
[![cover ≥82%](https://github.com/veraison/corim/actions/workflows/ci-go-cover.yml/badge.svg)](https://github.com/veraison/corim/actions/workflows/ci-go-cover.yml)
[![linters](https://github.com/veraison/corim/actions/workflows/linters.yml/badge.svg)](https://github.com/veraison/corim/actions/workflows/linters.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/veraison/corim.svg)](https://pkg.go.dev/github.com/veraison/corim)


The [`corim/corim`](corim) and [`corim/comid`](comid) packages provide a golang API for low-level manipulation of [Concise Reference Integrity Manifest (CoRIM)](https://datatracker.ietf.org/doc/draft-ietf-rats-corim/) and Concise Module Identifier (CoMID) tags respectively.
The [`corim/coev`](coev) package provides a minimal golang implementation of TCG Concise Evidence CDDL as documented [here](https://github.com/TrustedComputingGroup/dice-coev/blob/main/concise-evidence.cddl)
The [`corim/coserv`](coserv) package provides a golang API for working with [Concise Selector for Endorsements and Reference Values](https://datatracker.ietf.org/doc/draft-howard-rats-coserv).

> [!NOTE]
> These API are still in active development (as is the underlying CoRIM spec).
> They are **subject to change** in the future.

## Required Tools

Ensure you have the following tools installed with the specified versions on your machine to ensure everything works properly:

- **Go**: Version 1.22
- **golangci-lint**: Version 1.54.2

## Developer tips

Before requesting a PR (and routinely during the dev/test cycle), you are encouraged to run:
```
make presubmit
```
and check its output to make sure your code coverage figures are in line with the set target and that there are no newly introduced lint problems.

## x5chain trust verification

Signed CoRIM messages may carry an X.509 chain in the COSE `x5chain` protected header.
Use [`SignedCorim.VerifyWithX5Chain`](https://pkg.go.dev/github.com/veraison/corim/corim#SignedCorim.VerifyWithX5Chain)
after [`FromCOSE`](https://pkg.go.dev/github.com/veraison/corim/corim#SignedCorim.FromCOSE) to validate PKIX trust, optional CRL revocation, and the COSE signature.

Load trust material with [`LoadTrustAnchors`](https://pkg.go.dev/github.com/veraison/corim/corim#LoadTrustAnchors).
When no trust-anchor paths are supplied, verification uses the OS certificate store; for production deployments, pass explicit anchors.
When no CRL paths are supplied, revocation checks are skipped; when CRLs are loaded, [`CrlPolicyStrict`](https://pkg.go.dev/github.com/veraison/corim/corim#CrlPolicyStrict) is the default.

For external-key verification without PKIX path validation, use [`SignedCorim.Verify`](https://pkg.go.dev/github.com/veraison/corim/corim#SignedCorim.Verify) instead.

## Extending CoRIM/CoMID

The CoRIM specification provides a mechanism for adding extensions to the base
CoRIM schema. The `corim` and `comid` structs which can be extended, embed an
`Extensions` object  that allows registering a wrapper structure defining
extension fields. For field types that can be extended, i.e. `type choice`,
extensions can be implemented by calling an appropriate registration function
and giving it a new type or a value (for enums).

Please see [extensions documentation](extensions/README.md) for details.


