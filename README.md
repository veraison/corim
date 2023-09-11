# Concise Reference Integrity Manifest and Module Identifiers
[![ci](https://github.com/veraison/corim/actions/workflows/ci.yml/badge.svg)](https://github.com/veraison/corim/actions/workflows/ci.yml)
[![cover ≥82%](https://github.com/veraison/corim/actions/workflows/ci-go-cover.yml/badge.svg)](https://github.com/veraison/corim/actions/workflows/ci-go-cover.yml)
[![linters](https://github.com/veraison/corim/actions/workflows/linters.yml/badge.svg)](https://github.com/veraison/corim/actions/workflows/linters.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/veraison/corim.svg)](https://pkg.go.dev/github.com/veraison/corim)

The [`corim/v2/corim`](corim) and [`corim/v2/comid`](comid) packages provide a golang API for low-level manipulation of [Concise Reference Integrity Manifest (CoRIM)](https://datatracker.ietf.org/doc/draft-ietf-rats-corim/) and Concise Module Identifier (CoMID) tags respectively.

> [!NOTE]
> These API are still in active development (as is the underlying CoRIM spec).
> They are **subject to change** in the future.

The [`corim/v2/cocli`](cocli) package uses the API above (as well as the API from [`veraison/swid`](https://github.com/veraison/swid) package) to provide a user friendly command line interface for working with CoRIM, CoMID, CoSWID and CoTS.  Specifically it allows creating, signing, verifying, displaying, uploading, and more.  See [`cocli/README.md`](cocli/README.md) for further details.

## Developer tips

Before requesting a PR (and routinely during the dev/test cycle), you are encouraged to run:
```
make presubmit
```
and check its output to make sure your code coverage figures are in line with the set target and that there are no newly introduced lint problems.
