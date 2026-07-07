---
type: feature
title: Secured Files REST
description: Go multipart file upload, list, and delete against secured REST paths
resource: files.go
tags: [go-admin-sdk, files, rest, upload]
timestamp: 2026-07-07T00:00:00Z
---

# Secured Files REST

## Purpose

Go implementation of [secured-files-rest](../../../../.knowledge/features/secured-files-rest.md). Handles multipart upload and JSON list/delete on `/secured` REST base derived from GraphQL URL.

## Flows

- **Upload**: `UploadFile(ctx, params)` — multipart form with file bytes + metadata.
- **List**: `ListFiles(ctx, query params)` — paginated secured file index.
- **Delete**: `DeleteFiles(ctx, ids)` — bulk delete by file id.
- **Paths**: `files_paths.go` constants align with JS `filesPaths.ts`.

## Main files

- `files.go` — upload/list/delete implementations
- `files_paths.go` — `FILES_UPLOAD_PATH`, etc.
- `rest.go` — `executeREST`, multipart builder, auth headers

## Dependencies

- [client-config-graphql](client-config-graphql.md) for `restBaseURL`
- Global: [secured-files-rest](../../../../.knowledge/features/secured-files-rest.md)

## Invariants

- Upload uses multipart — not JSON GraphQL.
- Auth headers on REST match GraphQL key rules (`setAuthHeaders`).
- REST base must be configured — error if empty after derivation.

## Common bugs

- JSON Content-Type on upload — must use multipart writer in `rest.go`.
- Listing wrong folder/prefix — check query params against engine API.
- File ids from GraphQL media fields passed to delete without URL decode.

## Tests

- `client_rest_base_test.go`
- `examples/files/main.go` — manual integration example

## Related

- JS: [media-upload-headless](../js-admin-sdk/.knowledge/features/media-upload-headless.md)
- Global: [secured-files-rest](../../../../.knowledge/features/secured-files-rest.md)
