---
type: feature
title: Injected DB Interface
description: Go client methods matching engine InjectedDBOperationInterface for plugin and server injection
resource: client.go
tags: [go-admin-sdk, injection, engine, interface]
timestamp: 2026-07-07T00:00:00Z
---

# Injected DB Interface

## Purpose

Engine and plugins inject a DB operation implementation at runtime. Go `Client` exposes the same surface as JS `InjectedDBOperationInterface` so resolvers can delegate CRUD without duplicating HTTP logic.

## Flows

- **Engine injects**: pass `Client` (or thin wrapper) where `interfaces.InjectedDBOperationInterface` expected.
- **List/search**: dynamic model name + filter map → GraphQL search query.
- **CRUD**: create/update/delete with typed payloads from codegen.
- **Cross-SDK**: behavior locked by `CONTRACT.md` and shared naming vectors.

## Main files

- `client.go` — interface-aligned methods
- `CONTRACT.md` — cross-SDK injection contract
- `github.com/apito-io/types/interfaces` — shared Go interface definitions
- `generated_helpers.go` — helpers over genqlient output

## Dependencies

- [admin-sdk-contract](../../../../.knowledge/features/admin-sdk-contract.md)
- [typed-model-ops](typed-model-ops.md)
- Engine resolver injection points

## Invariants

- Method signatures must stay compatible with `apito-io/types` interfaces — breaking changes need coordinated release.
- Injection uses same auth headers as standalone SDK usage.
- Do not bypass client for raw HTTP inside engine resolvers.

## Common bugs

- Interface drift after SDK update — engine compile fails until types bumped.
- Missing context tenant on injected list calls in SaaS mode.
- Using system GraphQL URL in injected client config.

## Tests

- `client_test.go` — core op smoke tests
- Engine integration tests (external repo)

## Related

- JS parity: `js-admin-sdk` [apito-client-core](../js-admin-sdk/.knowledge/features/apito-client-core.md)
- `CONTRACT.md`, `DECISIONS.md`
