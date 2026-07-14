# go-admin-sdk — Handoff

## Branch
- Check submodule `git branch` / tags before push (getTenant work may be local)

## Done
- **v2.6.5 (2026-07-14):** `GetTenant(ctx, projectID, tenantID, status)` wraps `SearchTenants` + exact id match
- CONTRACT.md / README.md / CHANGELOG.md + `tenant_catalog_test.go`
- Earlier: `SearchTenants` + catalog row fields (2026-07-13)

## Broken / watch
- GetTenant is exact-id filter on SearchTenants — empty/nil when no match (callers must handle)

## Next
- Publish/tag 2.6.5; confirm Kisti/other Go consumers if any adopt GetTenant

## Do not touch
- Engine GraphQL naming without CONTRACT update across JS/Flutter/Go

## Last Updated
2026-07-14
