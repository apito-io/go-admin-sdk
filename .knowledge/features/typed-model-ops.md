---
type: feature
title: Typed Model Ops
description: Generated and helper-wrapped per-model CRUD operations on Go Client
resource: generated_helpers.go
tags: [go-admin-sdk, crud, models, genqlient]
timestamp: 2026-07-07T00:00:00Z
---

# Typed Model Ops

## Purpose

Typed Go methods for each schema model's five operations. Combines genqlient-generated clients with `generated_helpers.go` convenience wrappers matching [admin-sdk-contract](../../../../.knowledge/features/admin-sdk-contract.md).

## Flows

- **List**: generated `Get{Model}List` with where/sort variables.
- **Get one**: `Get{Model}(id)` with relation selections from document builder.
- **Mutations**: `Create{Model}`, `Update{Model}`, `Delete{Model}` with typed payloads.
- **Examples**: `examples/users/main.go`, `examples/files/main.go` demonstrate usage.

## Main files

- `generated_helpers.go` — wrapper helpers over codegen
- `codegen/operations/*.graphql` — per-model operation source
- genqlient output package (generated — do not edit)
- `document_builder.go` — field selection for get/list documents
- `examples/` — runnable samples

## Dependencies

- [naming-codegen-genqlient](naming-codegen-genqlient.md)
- [client-config-graphql](client-config-graphql.md)
- Filter variable shapes align with [list-filters-and-relations](../../../../.knowledge/features/list-filters-and-relations.md)

## Invariants

- Payload type names follow `{Model}_Create_Payload` / `{Model}_Update_Payload` convention.
- List where input uses `{Model}List_Input_Where_Payload` composed names.
- Regenerate after any engine model or field change.

## Common bugs

- Manual GraphQL strings diverge from generated ops — prefer genqlient calls.
- Update mutation sends full record — send only changed fields per engine rules.
- Missing relation fields in get query — nested data nil in apps.

## Tests

- `client_test.go`
- Example programs under `examples/`

## Related

- [injected-db-interface](injected-db-interface.md)
- Global: [admin-sdk-contract](../../../../.knowledge/features/admin-sdk-contract.md), [mutation-connect-relations](../../../../.knowledge/features/mutation-connect-relations.md)
