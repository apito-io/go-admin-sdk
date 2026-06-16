# Go Admin SDK — Cross-SDK Sync Summary

**Package:** `github.com/apito-io/go-admin-sdk` (v2.6.0)  
**Aligned with:** `flutter_admin_sdk` v0.5.0, `js-admin-sdk` v3.7.0

## Shared contract

See [CONTRACT.md](CONTRACT.md).

## v2.6.0 (2026-06-11)

- **`TenantID` on user CRUD** — `SearchUsers`, `CreateUser`, `UpdateUser` pass optional GraphQL `tenant_id` (pro SaaS)

## v2.5.0 (2026-06-08)

- **`LoginUser` `TenantID`** — optional GraphQL `tenant_id`; required for SaaS per-tenant separate DB projects

## v2.4.0 (2026-06-05)

- **`LoginUser` `google_id_token`** — native mobile sign-in via `IDToken`
- **Secured files REST** — default `RestBaseURL` → `/secured`

## v2.3.0 (2026-06-05)

- Added **naming engine** (`naming.go`) with golden-vector tests
- Added **operation emitter** (`make gen-operations`) → `codegen/operations/*.graphql` + `schema.graphql`
- Added **genqlient** (`genqlient.yaml`, `make gen-types`)
- Added **GraphQLDoer** + **TypedModelOps** context helpers
- Introspection snapshot: `schema/apito_introspection.json`

## Make targets

| Target | Purpose |
|--------|---------|
| `make gen-operations` | Emit .graphql ops from introspection |
| `make gen-types` | Run genqlient |
| `make gen` | Both |

## Prior versions

See [CHANGELOG.md](CHANGELOG.md) for v2.2.0 (files), v2.1.0 (file rename), v2.0.0 (user rename).
