# Handoff — go-admin-sdk 2026-10-18

- **Done:** Package comment, example tests, and README. Remote stays `go-internal-sdk`. Do not rename the GitHub repo in this pass. Do not tag. Leave `.memory` out of the public diff.
- **Verified:** `go test -count=1 -run Example .` passed.

- **Done:** `LoginUser` now keeps one canonical query with optional `access_token`; the Facebook method supplies the variable and existing methods stay compatible for published clients. Added a request-contract test.
- **Verified:** `go test ./... -run TestLoginUserFacebookAccessTokenUsesCanonicalOptionalArgument -count=1` passed.
- **Not done:** Version bump, tag, and commit.
- **Do not:** `./command.sh save` until asked.

- **Done:** `LoginUser` sends `access_token` for `facebook_access_token`.
- **Not done:** Version bump. Commit.
- **Do not:** `./command.sh save` until asked.

# Handoff — go-admin-sdk 2026-10-12

## Consumer codegen flags

- **Done:** Flagged `apito-gen` writes under `-out` for a consumer module and can delete graphql files it did not rewrite.
- **Not done:** Version bump. Commit.
- **Do not:** Run the SDK repo `go generate` from the CLI inside an app folder. `./command.sh save` until asked.

# Handoff — go-admin-sdk 2026-08-12

## Done (uncommitted)

- `project_payment.go` + tests; CONTRACT/CHANGELOG monetization section

## Next

- Ask before commit/tag
