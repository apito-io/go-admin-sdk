---
type: feature
title: Auth Tenant REST Admin
description: Login, tenant catalog, and project user admin operations on Go Client
resource: models.go
tags: [go-admin-sdk, auth, tenant, admin]
timestamp: 2026-07-07T00:00:00Z
---

# Auth Tenant REST Admin

## Purpose

Go bindings for [auth-tenant-admin](../../../../.knowledge/features/auth-tenant-admin.md): end-user login, SaaS tenant catalog, and project user CRUD. Mirrors JS `ApitoClient` auth methods.

## Flows

- **Login**: `LoginUser(ctx, LoginUserParams)` — password, Google OAuth, or `google_id_token`.
- **OAuth state**: `GoogleOAuthState` before redirect flow.
- **Tenants**: `GetTenants`, `CreateTenant`, `UpdateTenant`, `DeleteTenant`, `GenerateTenantToken`.
- **Users**: `SearchUsers`, `CreateUser`, `UpdateUser`, `DeleteUser`, `ResetUserPassword` with optional `TenantID`.
- **Domain lookup**: `SearchTenantsByDomain` for SaaS routing.

## Main files

- `client.go` — auth/tenant/user method implementations
- `models.go` — `User`, `LoginUserParams`, `CreateUserParams`, tenant types
- `tenant_catalog_test.go` — tenant catalog tests

## Dependencies

- Global: [auth-tenant-admin](../../../../.knowledge/features/auth-tenant-admin.md)
- System GraphQL operations (some admin ops)
- Cloudflare Workers v1 limitations per `CONTRACT.md`

## Invariants

- SaaS user ops require `TenantID` in params or context when engine expects it.
- Google login unavailable on Workers v1 — handle GraphQL error string.
- Password vs Google auth paths are mutually exclusive per `AuthMethod`.

## Common bugs

- `google_id_token` sent without server client id configured on engine.
- Duplicate email/phone errors not mapped to user-friendly messages.
- Tenant token generation attempted on Workers v1.

## Tests

- `client_user_tenant_test.go`
- `tenant_catalog_test.go`

## Related

- [client-config-graphql](client-config-graphql.md)
- Global: [auth-tenant-admin](../../../../.knowledge/features/auth-tenant-admin.md)
