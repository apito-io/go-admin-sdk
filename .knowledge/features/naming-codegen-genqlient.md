---
type: feature
title: Naming Codegen Genqlient
description: apito-gen CLI, naming.go, and genqlient pipeline from introspection to Go client
resource: cmd/apito-gen/main.go
tags: [go-admin-sdk, codegen, genqlient, naming]
timestamp: 2026-07-07T00:00:00Z
---

# Naming Codegen Genqlient

## Purpose

Generates per-model `.graphql` operation files and genqlient Go client from checked-in introspection. Naming follows [naming-engine](../../../../.knowledge/features/naming-engine.md) golden vectors in `test/fixtures/naming_vectors.json`.

## Flows

- **Generate ops**: `go run ./cmd/apito-gen` or `make generate` → `codegen/operations/*.graphql`.
- **Genqlient**: `genqlient.yaml` → typed Go client in generated package.
- **Filter models**: `APITO_MODELS=loan,customer` env limits output.
- **Schema input**: `schema/apito_introspection.json` or `APITO_SCHEMA_FILE`.

Pipeline: [introspection-codegen-pipeline](../../../../.knowledge/features/introspection-codegen-pipeline.md).

## Main files

- `cmd/apito-gen/main.go` — CLI entry
- `naming.go`, `document_builder.go`, `schema_reader.go` — naming + doc build
- `generate.go` — `go:generate` directives
- `genqlient.yaml` — genqlient config
- `codegen/operations/` — generated GraphQL documents
- `Makefile` — `generate` target

## Dependencies

- [naming-engine](../../../../.knowledge/features/naming-engine.md)
- genqlient code generator
- Shared vectors synced with JS/Flutter per `CONTRACT.md`

## Invariants

- Never hand-edit genqlient output — regenerate from operations + schema.
- Operation names must pass `naming_test.go` vector assertions.
- Five ops per model: list, get, create, update, delete.

## Common bugs

- Stale `schema/apito_introspection.json` after engine schema change.
- genqlient version mismatch — check `go.mod` after upgrade.
- Forgot to run genqlient after `apito-gen` — compile errors on missing types.

## Tests

- `naming_test.go`
- `make generate && go test ./...`

## Related

- [typed-model-ops](typed-model-ops.md)
- JS: [codegen-cli](../js-admin-sdk/.knowledge/features/codegen-cli.md)
- Global: [introspection-codegen-pipeline](../../../../.knowledge/features/introspection-codegen-pipeline.md)
