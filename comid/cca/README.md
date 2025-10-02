# CCA (Confidential Computing Architecture) Profiles

This package defines CCA profile identifiers using the tag URI scheme as specified in [RFC 4151](https://tools.ietf.org/html/rfc4151).

## Profile Identifiers

Three profile identifiers are defined:

1. **CCA Token Profile**
   ```
   tag:arm.com,2025:cca-token
   ```
   Used for CCA attestation tokens.

2. **CCA Platform Endorsements Profile**
   ```
   tag:arm.com,2025:cca-endorsements
   ```
   Used for CCA platform endorsements.

3. **CCA Realm Endorsements Profile**
   ```
   tag:arm.com,2025:cca-realm-endorsements
   ```
   Used for CCA realm endorsements.

## Usage

```go
import "github.com/veraison/corim/comid/cca"

func example() {
    // Use CCA Token Profile
    tokenProfile := cca.TokenProfileID

    // Use CCA Platform Endorsements Profile
    platformProfile := cca.EndorsementsProfileID

    // Use CCA Realm Endorsements Profile
    realmProfile := cca.RealmEndorsementsProfileID
}
```

## Tag URI Format

The tag URIs follow RFC 4151 format:
- Authority: `arm.com` - representing Arm Limited
- Date: `2025` - year of profile definition
- Specific ID: One of `cca-token`, `cca-endorsements`, or `cca-realm-endorsements`

These tag URIs are used instead of HTTP URLs to avoid accidental dereferencing while maintaining unique identification.