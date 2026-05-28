# Changelog

All notable changes to the Go Apito SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.2.0] - 2026-05-28

### Changed

- **Project files storage** — Document exported path constants (`FilesUploadPath`, `FilesListPath`, `FilesDeletePath`). File metadata is persisted in the **project DB** `files` table (engine migration from system DB); REST URLs remain `/system/files/*`. SaaS callers should set `tenant_id` on the request context.

## [2.1.0] - 2026-05-17

### Changed (breaking)

- **Removed storage settings GraphQL** — `GetProjectStorageSettings` and `UpdateProjectStorageSettings` dropped from the SDK (configure storage in the console or via raw GraphQL).
- **File API renamed** — action-oriented names aligned with the User API: `UploadFile`, `ListFiles`, `DeleteFiles` (was `*SystemFile*`).

### Migration

| v2.0.0 | v2.1.0 |
|--------|--------|
| `GetProjectStorageSettings` | removed |
| `UpdateProjectStorageSettings` | removed |
| `UploadSystemFile` | `UploadFile` |
| `ListSystemFiles` | `ListFiles` |
| `DeleteSystemFiles` | `DeleteFiles` |
| `SystemFile` | `File` |
| `SystemFileUploadParams` | `UploadFileParams` |
| `SystemFilesListResponse` | `FilesListResponse` |
| `DeleteSystemFilesResponse` | `DeleteFilesResponse` |

## [2.0.0] - 2026-05-17

### Changed (breaking)

- **Tenant-user GraphQL renamed to User API** — aligned with engine open-core migration (`users` table, `UserItem` type). All `*TenantUser*` types and methods renamed to `*User*` (e.g. `LoginUser`, `SearchUsers`, `CreateUser`, `UpdateUser`, `DeleteUser`).
- **`googleOAuthState`** replaces `tenantGoogleOAuthState`.
- **`UpdateUser`** no longer accepts `password`; use **`ResetUserPassword`**.

### Added

- **`ResetUserPassword(ctx, userID, password)`** — admin password reset mutation.
- **`GetProjectStorageSettings`**, **`UpdateProjectStorageSettings`** — project S3/storage settings GraphQL.
- **`UploadSystemFile`**, **`ListSystemFiles`**, **`DeleteSystemFiles`** — `/system/files` REST API (`Config.RestBaseURL` optional; derived from GraphQL base URL).
- Examples: `examples/users/`, `examples/system_files/` (replaces `examples/tenant_users/`).

### Migration

| v1.7.x | v2.0.0 |
|--------|--------|
| `LoginTenantUser` | `LoginUser` |
| `TenantGoogleOAuthState` | `GoogleOAuthState` |
| `SearchTenantUsers` | `SearchUsers` |
| `CreateTenantUser` | `CreateUser` |
| `UpdateTenantUser` (+ `Password`) | `UpdateUser` + `ResetUserPassword` |
| `DeleteTenantUser` | `DeleteUser` |
| `TenantUser` | `User` |

## [1.7.0] - 2026-05-08

### Changed (breaking)

- **`LoginTenantUserGoogle` removed** — engine dropped **`loginTenantUserGoogle`**. Use **`LoginTenantUser`** with **`AuthMethod: "google"`**, **`Code`**, **`State`**.

### Added

- **`TenantGoogleOAuthState(ctx, projectID)`** → **`TenantGoogleOAuthStateResponse`** (**`State`**) for the Google authorize redirect.

### `LoginTenantUserParams`

- **`Code`**, **`State`** for Google flow.
- **`Password`** / **`Email`** / **`Phone`** only required when **`AuthMethod`** is empty or **`general`**.

## [1.6.0] - 2026-05-14

### Changed (breaking)

- **Tenant catalog users** aligned with engine Pro GraphQL: **`TenantUser`** now has **`Phone`** (no **`Username`**). **`LoginTenantUser`** is now **`LoginTenantUser(ctx, projectID, LoginTenantUserParams)`** with **`Password`**, optional **`Email`** / **`Phone`**, optional **`AuthMethod`**. **`CreateTenantUser`** is **`CreateTenantUser(ctx, projectID, CreateTenantUserParams)`** (**`Password`**, optional **`Role`**, **`Email`**, **`Phone`**).
- Added **`UpdateTenantUser`** and **`DeleteTenantUser`** (arguments match system GraphQL; project scope comes from the API key).

### Migration

- Replace `LoginTenantUser(ctx, pid, username, password)` with  
  `LoginTenantUser(ctx, pid, LoginTenantUserParams{Password: password, Email: "..."})` or `Phone: "..."` per project sign-in mode.
- Replace `CreateTenantUser(ctx, pid, username, email, password, role)` with  
  `CreateTenantUser(ctx, pid, CreateTenantUserParams{Password: password, Role: role, Email: email, Phone: phone})`.

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
