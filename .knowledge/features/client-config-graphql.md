---
type: feature
title: Client Config GraphQL
description: Go Client and Config with GraphQL POST execution and REST base derivation
resource: client.go
tags: [go-admin-sdk, graphql, client, config]
timestamp: 2026-07-07T00:00:00Z
---

# Client Config GraphQL

## Purpose

Go implementation of [admin-sdk-contract](../../../../.knowledge/features/admin-sdk-contract.md). `Client` executes GraphQL against project public API with API key auth, optional tenant context, and derived REST base for file ops.

## Flows

- **Create**: `NewClient(Config{ BaseURL, APIKey, RestBaseURL?, Timeout? })`.
- **GraphQL**: `executeGraphQL(ctx, query, variables)` — sets `X-Apito-Key` / sync key headers.
- **Tenant**: pass `tenant_id` via `context.Context` value → `X-Apito-Tenant-ID` header.
- **Model ops**: generated genqlient methods or hand-built queries via document builder.

## Main files

- `client.go` — `Client`, `Config`, `NewClient`, GraphQL dispatch
- `graphql_doer.go` — HTTP doer abstraction for genqlient
- `types` from `github.com/apito-io/types` — shared GraphQL response shapes

## Dependencies

- Engine public GraphQL endpoint
- [naming-engine](../../../../.knowledge/features/naming-engine.md)
- Standard library `net/http`

## Invariants

- `RestBaseURL` derived from `BaseURL` when empty — same rules as JS SDK (`/system/graphql` → `/secured`).
- `cli-*` / `sdk-*` keys use `X-Apito-Sync-Key`; others use `X-Apito-Key`.
- Always use `context.Context` for cancellation and tenant injection.

## Common bugs

- Tenant not in context → SaaS queries return wrong tenant data.
- Timeout zero — defaults to 30s; long exports need custom `HTTPClient`.
- Wrong endpoint (system vs public GraphQL).

## Tests

- `client_test.go`
- `client_rest_base_test.go`
- `client_user_tenant_test.go`

## Related

- [typed-model-ops](typed-model-ops.md), [secured-files-rest](secured-files-rest.md)
- Global: [admin-sdk-contract](../../../../.knowledge/features/admin-sdk-contract.md)
