# go-admin-sdk — Current

- **Facebook token (2026-10-16, uncommitted):** `LoginUserParams.AccessToken` uses the canonical `loginUser` query's optional `access_token` argument. Existing login methods remain valid for published clients. No version bump. Ask before commit or `./command.sh save`.

- **Done (2026-10-12, uncommitted):** `cmd/apito-gen` accepts `-schema`, `-out`, `-models`, `-module`, and `-prune`. With no flags it still walks to `go.mod` and writes `codegen/`. `TestApitoGenWritesConsumerModuleAndPrunes` passed. No version bump. Ask before commit.

- **Done (2026-08-12, uncommitted):** ApitoProjectPayment REST —
  `VerifyProjectPayment` / `GetTenantSubscription` / `CancelTenantSubscription`
  + `ProjectTenantSubscription`.
- Prior: OAuthState / access helpers / GetTenant.

## Next

- Confirm commit with engine Play cancel release

## Last Updated

2026-08-12
