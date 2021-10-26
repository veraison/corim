# Concise Reference Integrity Manifest and Module Identifiers

The [`corim/corim`](corim) and [`corim/comid`](comid) packages provide a golang API for low-level manipulation of [Concise Reference Integrity Manifest (CoRIM)](https://datatracker.ietf.org/doc/draft-birkholz-rats-corim/) and Concise Module Identifier (CoMID) tags respectively.

The [`corim/cocli`](cocli) package uses the API above (as well as the API from [`veraison/swid`](https://github.com/veraison/swid) package) to provide a user friendly command line interface for working with CoRIM, CoMID and CoSWID.  Specifically it allows creating, signing, verifying, displaying, and more.  See [`cocli/README.md`](cocli/README.md) for further details.

## Resources

* [Package documentation](https://pkg.go.dev/github.com/veraison/corim)

## Developer tips

Before requesting a PR (and routinely during the dev/test cycle), you are encouraged to run:
```
make presubmit
```
and check its output to make sure your code coverage figures are in line with the set target and that there are no newly introduced lint problems.
