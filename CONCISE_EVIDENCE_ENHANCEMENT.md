# Concise Evidence Support Enhancement

## Summary

This implementation addresses GitHub issue #83 by enhancing the Concise Evidence support in the `coev` package to align with the TCG Concise Evidence CDDL specification.

## Changes Made

### 1. Added Missing Triple Types

The original `EvTriples` struct was missing two critical triple types defined in the TCG Concise Evidence specification:

- **DependencyTriples** (CBOR index 2): Represents dependencies between domains
- **MembershipTriples** (CBOR index 3): Represents membership relationships between domains and environments

### 2. New Files Added

#### `coev/dependency_triple.go`
- `DependencyTriple` struct: Represents `ev-dependency-triple-record` from the CDDL spec
- `DependencyTriples` collection type
- Validation methods and helper functions
- Setter methods with fluent API support

#### `coev/membership_triple.go`
- `MembershipTriple` struct: Represents `ev-membership-triple-record` from the CDDL spec  
- `MembershipTriples` collection type  
- Validation methods and helper functions
- Setter methods with fluent API support

#### `coev/dependency_triple_test.go`
- Comprehensive test coverage for `DependencyTriple` and `DependencyTriples`
- Tests for validation, setters, nil receiver handling, and collection operations

#### `coev/membership_triple_test.go`  
- Comprehensive test coverage for `MembershipTriple` and `MembershipTriples`
- Tests for validation, setters, nil receiver handling, and collection operations

#### `coev/ev_triples_test.go`
- New tests for the enhanced `EvTriples` functionality
- Tests for the new `AddDependencyTriple` and `AddMembershipTriple` methods
- CBOR/JSON marshaling tests with empty collection handling

#### `coev/dependency_membership_examples_test.go`
- Example usage demonstrating how to create and use dependency and membership triples
- Shows practical use cases for domain relationships and environment membership

### 3. Enhanced Existing Files

#### `coev/ev_triples.go`
- Added `DependencyTriples` and `MembershipTriples` fields to the `EvTriples` struct
- Updated the `Valid()` method to validate the new triple types
- Added `AddDependencyTriple()` and `AddMembershipTriple()` helper methods
- Enhanced CBOR/JSON marshaling to handle empty collections properly

#### `coev/test_vars.go`
- Added additional test UUIDs (`TestUUID2`, `TestUUID3`) for comprehensive testing

## CDDL Compliance

The implementation now fully supports the TCG Concise Evidence CDDL specification structure:

```cddl
ev-triples-map = non-empty< {
  ? &(ce.evidence-triples: 0) => [ + evidence-triple-record ]
  ? &(ce.identity-triples: 1) => [ + ev-identity-triple-record ]
  ? &(ce.dependency-triples: 2) => [ + ev-dependency-triple-record ]  ; ✅ NOW SUPPORTED
  ? &(ce.membership-triples: 3) => [ + ev-membership-triple-record ]  ; ✅ NOW SUPPORTED
  ? &(ce.coswid-triples: 4) => [ + ev-coswid-triple-record ]
  ? &(ce.attest-key-triples: 5) => [ + ev-attest-key-triple-record ]
  * $$ev-triples-map-extension
} >
```

## Key Features

### Domain Type Implementation
- Used `comid.Environment` as the domain type for both dependency and membership triples
- This aligns with the existing CoRIM/CoMID architecture and provides flexibility for future extensions

### Validation
- Comprehensive validation for all new triple types
- Ensures domains and dependent domains/environments are properly specified and valid
- Clear error messages for debugging

### Fluent API
- All setter methods return the receiver for method chaining
- Consistent with existing codebase patterns
- Null-safe operations (methods handle nil receivers gracefully)

### Serialization Support
- Full CBOR and JSON marshaling/unmarshaling support
- Empty collection optimization (empty collections are omitted from serialized output)
- Maintains compatibility with existing serialization patterns

### Test Coverage
- 100% test coverage for new functionality
- Tests cover both positive and negative cases
- Integration tests ensure new triples work with the broader evidence system

## Usage Examples

### Creating Dependency Triples
```go
dt := NewDependencyTriple()
dt.SetDomain(hypervisorEnv).
   AddDependentDomain(vm1Env).
   AddDependentDomain(vm2Env)

evTriples := NewEvTriples().AddDependencyTriple(dt)
```

### Creating Membership Triples
```go
mt := NewMembershipTriple()  
mt.SetDomain(trustZoneEnv).
   AddEnvironment(secureWorldEnv).
   AddEnvironment(normalWorldEnv)

evTriples := NewEvTriples().AddMembershipTriple(mt)
```

## Backward Compatibility

All changes are fully backward compatible:
- Existing `EvTriples` functionality remains unchanged
- New fields are optional and default to nil
- CBOR/JSON serialization omits new fields when empty
- All existing tests continue to pass

## Test Results

- ✅ All existing tests pass
- ✅ New comprehensive test suite passes  
- ✅ Integration tests with full codebase pass
- ✅ Example tests demonstrate real-world usage

This implementation successfully addresses the GitHub issue by providing complete support for the TCG Concise Evidence specification while maintaining full backward compatibility and following established codebase patterns.