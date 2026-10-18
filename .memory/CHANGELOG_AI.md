# go-admin-sdk — AI Changelog

Not git history — the *reasoning* behind changes. Newest on top.
Format per entry: date, **Changed**, **Why**, **Affected**.

---
## 2026-10-16 — Canonical Facebook token login query

- **Changed:** `LoginUser` uses one canonical GraphQL document with optional `access_token`; added a request-contract test. Removed the temporary per-auth query fork.
- **Why:** The engine schema evolves additively for published mobile clients. Old binaries omit the argument; SDK/app releases using it follow the engine release.
- **Affected:** `client.go`, `models.go`, `client_user_tenant_test.go`. No version bump. Ask before commit, tag, or `./command.sh save`.

## 2026-10-15 — Facebook access token on LoginUser

- **Changed:** `LoginUserParams.AccessToken` and the `facebook_access_token` branch.
- **Why:** Native Facebook login has a user token, not a code and state.
- **Affected:** `client.go`, `models.go`. Ask before commit or `./command.sh save`.

## 2026-10-12 — consumer apito-gen flags

- **Changed:** `-schema`, `-out`, `-models`, `-module`, and `-prune`. Empty args keep the SDK repo layout.
- **Why:** `apito gen` has to write into an app module without walking up to the SDK `go.mod` or leaving stale operations.
- **Affected:** `cmd/apito-gen/main.go`. Test passed. No version bump. Ask before commit.

---
## 2026-08-12 — ApitoProjectPayment REST

- **Changed:** `VerifyProjectPayment` / `GetTenantSubscription` /
  `CancelTenantSubscription` + `ProjectTenantSubscription`;
  `deriveSystemAPIBaseURL`; CONTRACT/CHANGELOG.
- **Why:** Parity with Flutter/JS for engine-owned Play cancel/verify.
- **Affected:** `project_payment.go`, tests, CONTRACT. Ask before commit.

## 2026-08-06 — v2.6.8 OAuthState + multi-provider LoginUser

- **Changed:** `OAuthState`; LoginUser auth_method facebook/github/x/linkedin.
- **Why:** Match engine multi-provider app-user OAuth.
- **Affected:** `client.go`, `version.go` **2.6.8**, CHANGELOG.

## 2026-07-21
- **Changed:** Added `Config.ProjectID`, canonical scoped GraphQL/REST project
  headers, and explicit method project overrides with tenant context preserved.
- **Why:** `apt_` authorization must receive project scope in HTTP headers.
- **Affected:** `client.go`, `rest.go`, header tests, docs/knowledge/changelog.

## 2026-07-14
- **Changed:** v2.6.5 — `GetTenant(ctx, projectID, tenantID, status)` via SearchTenants + exact id match; CONTRACT/README/CHANGELOG + tests.
- **Why:** Close getTenant parity with JS/Flutter for single-tenant catalog fetch without consumer-side search loops.
- **Affected:** `client.go`, `tenant_catalog_test.go`, `CHANGELOG.md`, `CONTRACT.md`, `README.md`, `version.go`

## 2026-07-13
- **Changed:** `SearchTenants(ctx, projectID, limit, offset, q)` + `SearchTenantsResponse`; `TenantCatalogSearchRow` gains `icon`, `created_at`. `/sync-sdk-all apply go` after tenant user-parity.
- **Why:** Only catalog search gap vs JS/Flutter/engine; Kisti and Console use `searchTenants` for catalog counts.
- **Affected:** `client.go`, `models.go`, `tenant_catalog_test.go`, `CHANGELOG.md`, `CONTRACT.md`, `README.md`

## 2026-07-06
- **Changed:** Bootstrapped knowledge system for this repo.
- **Why:** Cross-LLM durable knowledge + working memory.
- **Affected:** this repo only.

Last Updated: 2026-07-21
