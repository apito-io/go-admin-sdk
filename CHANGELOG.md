# Changelog

All notable changes to the Go Apito SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.5.2] - 2026-05-13

### Changed (breaking — module path only)

- **Module path** is now **`github.com/apito-io/go-admin-sdk`** (was `github.com/apito-io/go-internal-sdk`). Update imports and `go get github.com/apito-io/go-admin-sdk@v1.5.2`. The Git remote may still point at a repository named `go-internal-sdk` until it is renamed on GitHub; the **module path** in `go.mod` is what `go get` uses.

### Fixed

- **Examples** and **README** use the **`go-admin-sdk`** import path consistently.

## [1.5.1] - 2026-05-13

### Changed (breaking)

- **`GenerateTenantToken`**: signature is now **`(ctx, tenantID, duration, role string)`**, aligned with engine `generateTenantToken` (`tenant_id`, `duration`, optional `role`). Removed the unused legacy **`token`** first argument. Empty **`duration`** still selects the default one-year-ahead expiry in UTC.
- **`github.com/apito-io/types`** `InternalSDKOperation` updated in lockstep. **`go.mod`** requires **`github.com/apito-io/types v0.1.10`** or newer.

## [1.5.0] - 2026-05-09

### Changed (breaking)

- **`SearchTenantsByDomain`**: signature is now `(ctx, projectID, domain)` only (no pagination). Response type renamed to **`TenantByDomainResponse`** with a single nullable **`Tenant`** field (exact per-project domain match; was list + count).

### Engine parity (documented)

- System GraphQL **`searchTenantsByDomain`** returns a single nullable **`tenant`** (no list/count). **`createTenant`** optional **`domain`** is rejected when that domain is already taken in the project; **`updateTenant`** enforces the same when **`domain`** is set to a non-empty value.

## [1.4.0] - 2026-05-09

### Added

- **Pro tenant catalog search by domain**: `SearchTenantsByDomain`; types `TenantCatalogSearchRow`, `TenantsByDomainResponse` (engine `searchTenantsByDomain` on system GraphQL).

## [1.3.0] - 2026-05-08

### Added

- **Pro tenant catalog users** (Apito Pro system GraphQL): `LoginTenantUser`, `LoginTenantUserGoogle`, `SearchTenantUsers`, `CreateTenantUser`; types `TenantUser`, `TenantLoginResponse`, `TenantUsersResponse`.
- **`examples/tenant_users`**: minimal runnable sample using env `APITO_BASE_URL`, `APITO_API_KEY`, `APITO_PROJECT_ID` (optional `APITO_TENANT_USERNAME` / `APITO_TENANT_PASSWORD` for login).

### Tests

- **`TestTenantUserProIntegration`**: optional live checks when `APITO_PROJECT_ID` is set; skipped otherwise.

## [1.2.0] - 2024-12-30

### Added

- 🎯 **Type-Safe Operations**: Complete generic typed methods for all operations
  - `GetSingleResourceTyped[T]()` for type-safe single resource retrieval
  - `SearchResourcesTyped[T]()` for type-safe search operations
  - `GetRelationDocumentsTyped[T]()` for type-safe relation queries
  - `CreateNewResourceTyped[T]()` for type-safe resource creation
  - `UpdateResourceTyped[T]()` for type-safe resource updates
- 🚀 **Comprehensive Todo Example**: Complete practical example demonstrating all SDK features
  - Authentication & tenant token generation
  - Resource creation (todos, users, categories)
  - Both typed and untyped search operations
  - Single resource retrieval
  - Resource updates with connections
  - Relation document queries
  - Audit logging
  - Debug functionality
  - Resource cleanup
- 📚 **Enhanced Documentation**: Completely rewritten README with comprehensive examples
  - Quick start guide
  - Complete API reference
  - Type system documentation
  - Plugin integration examples
  - Production deployment guides
  - Performance optimization tips
  - Error handling best practices
- 🔧 **Improved Request Structure**: New `CreateAndUpdateRequest` struct for cleaner API
- 📊 **Version Tracking**: Added `version.go` with `GetVersion()` function

### Changed

- 🔄 **Updated Client Interface**: Enhanced all methods to use the new request structure
- 📖 **Documentation**: Complete rewrite with practical examples and comprehensive coverage
- 🎨 **Example Structure**: Replaced basic example with comprehensive todo application

### Fixed

- 🐛 **Type Conversion**: Improved JSON marshaling/unmarshaling for typed operations
- 🔧 **Error Handling**: Enhanced GraphQL and HTTP error reporting

### Technical Details

- All generic functions follow the pattern: `OperationTyped[T](client, ctx, ...params)`
- Backward compatibility maintained for all existing non-typed methods
- Enhanced context support with tenant ID handling
- Improved connection pooling and performance optimizations

## [1.1.3] - Previous Version

- Previous features and bug fixes

## [1.1.2] - Previous Version

- Previous features and bug fixes

## [1.1.1] - Previous Version

- Previous features and bug fixes

## [1.1.0] - Previous Version

- Previous features and bug fixes

## [1.0.0] - Initial Release

- Initial SDK implementation
- Basic GraphQL communication
- API key authentication
- Core CRUD operations
