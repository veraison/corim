# Concise Reference Integrity Manifest and Module Identifiers

The [`corim/corim`](corim) and [`corim/comid`](comid) packages provide a golang API for low-level manipulation of [Concise Reference Integrity Manifest (CoRIM)](https://datatracker.ietf.org/doc/draft-birkholz-rats-corim/) and Concise Module Identifier (CoMID) tags respectively.

The [`corim/cocli`](cocli) package builds on them (as well as the [swid](https://github.com/veraison/swid) package) to provide a command line interface for dealing with CoRIM and CoMID, which allows creating, signing, verifying, visualising, and more.

## Resources

* [Package documentation](https://pkg.go.dev/github.com/veraison/corim)

## Developer tips

Before requesting a PR (and routinely during the dev/test cycle), you are encouraged to run:
```
make presubmit
```
and check its output to make sure your code coverage figures are in line with the set target and that there are no newly introduced lint problems.
