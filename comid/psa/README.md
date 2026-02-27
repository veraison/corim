# PSA (Platform Security Architecture) Profiles

This package defines PSA profile identifiers using the tag URI scheme as specified in [RFC 4151](https://tools.ietf.org/html/rfc4151).

## Profile Identifiers

Two profile identifiers are defined:

1. **PSA Token Profile**
   ```
   tag:trustedcomputinggroup.org,2025:psa-token
   ```
   Used for PSA attestation tokens as defined in [draft-tschofenig-rats-psa-token](https://datatracker.ietf.org/doc/html/draft-tschofenig-rats-psa-token).

2. **PSA Platform Endorsements Profile**
   ```
   tag:trustedcomputinggroup.org,2025:psa-endorsements
   ```
   Used for PSA platform endorsements.

## Usage

```go
import "github.com/veraison/corim/comid/psa"

func example() {
    // Use PSA Token Profile
    tokenProfile := psa.TokenProfileID

    // Use PSA Endorsements Profile
    endorsementsProfile := psa.EndorsementsProfileID
}
```

## Tag URI Format

The tag URIs follow RFC 4151 format:
- Authority: `trustedcomputinggroup.org` - representing the TCG organization
- Date: `2025` - year of profile definition
- Specific ID: Either `psa-token` or `psa-endorsements`

These tag URIs are used instead of HTTP URLs to avoid accidental dereferencing while maintaining unique identification.