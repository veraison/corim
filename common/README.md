# Common Package

This package contains shared types and utilities used across the `corim`, `comid`, and related packages.

## Purpose

Created to address issue #11: "Factor out common code"

## Analysis of Common Code

After analyzing the codebase, the following areas of duplication were identified:

### 1. Entity Name Implementation
- **Location**: `comid/entity.go` and `corim/entity.go`
- **Duplication**: 
  - `EntityName` struct and methods (~200 lines duplicated)
  - `StringEntityName` type and methods
  - `IEntityNameValue` interface
  - `IEntityNameFactory` type
  - Factory functions: `NewEntityName`, `MustNewEntityName`, `NewStringEntityName`, `MustNewStringEntityName`
  - CBOR/JSON marshaling logic for EntityName
  - Entity name registry and `RegisterEntityNameType`

**Recommendation**: The EntityName implementation is nearly identical between comid and corim packages. However, each package requires different CBOR tag registries. A refactoring would need to:
- Extract common base types to `common` package
- Keep package-specific CBOR/JSON marshaling in each package
- Or use a more sophisticated registry pattern

### 2. Role/Roles Implementation
- **Location**: `comid/role.go` and `corim/role.go`
- **Duplication**:
  - `Role` type (int64)
  - `Roles` slice type
  - `RegisterRole` function  
  - `String()` method
  - `Valid()` method for Roles
  - JSON marshaling/unmarshaling logic
  - roleToString and stringToRole maps

**Differences**:
- Different role constants (comid: TagCreator, Creator, Maintainer; corim: ManifestCreator)
- comid's Roles.Add() doesn't validate, corim's does
- corim has additional ToJSON/FromJSON helper methods
- comid has ToCBOR/FromCBOR methods

**Recommendation**: Could extract a base `Role` type and registration system, but each package would keep its own role constants and any package-specific behavior.

### 3. TaggedURI
- **Location**: `comid/entity.go`
- **Usage**: Used by both comid and corim Entity types
- Simple string wrapper with `Empty()` method

**Recommendation**: Easy candidate for extraction - it's a simple utility type with no dependencies.

### 4. Common Validation/Marshaling Patterns
- Many types implement similar `Valid()`, `MarshalCBOR()`, `UnmarshalCBOR()`, `MarshalJSON()`, `UnmarshalJSON()` patterns
- Uses `encoding.PopulateStructFromCBOR`, `encoding.SerializeStructToCBOR`, etc.

**Recommendation**: These patterns are already well-factored through the `encoding` package. No further action needed.

## Implementation Considerations

### Challenges
1. **CBOR Tag Registration**: comid and corim have separate CBOR tag registries. Moving types to `common` would require either:
   - Maintaining separate encoders/decoders in each package
   - Creating a more complex registry system
   - Keeping marshaling logic in each package

2. **Package Dependencies**: Need to avoid circular dependencies between common, comid, and corim

3. **Backward Compatibility**: Any refactoring must maintain the existing public API

4. **Test Coverage**: Extensive test suites exist for current code - all must continue to pass

### Recommended Approach

Given the complexity of CBOR tag handling and the relatively small size of the codebase, the most practical approach is:

1. **Phase 1** (this PR): Document common code patterns and create the `common` package structure
2. **Phase 2** (future PR): Extract truly independent utilities (like TaggedURI)
3. **Phase 3** (future PR): Consider more complex refactoring of Role/EntityName if warranted

## Status

**Current**: Analysis and documentation complete.  
**Next Steps**: Team decision on whether to proceed with extraction or keep as-is with documentation.

Note: The `cocli` package mentioned in the original issue has been moved to a separate repository, so cocli/cmd/common.go is no longer relevant to this codebase.
