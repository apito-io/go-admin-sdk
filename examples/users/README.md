# Project users example

Lists project end-users and optionally logs in with email/phone + password.

## Environment

| Variable | Required | Description |
|----------|----------|-------------|
| `APITO_API_KEY` | yes | Admin API key |
| `APITO_PROJECT_ID` | yes | Project ID |
| `APITO_BASE_URL` | no | Defaults to `http://localhost:5050/system/graphql` |
| `APITO_TENANT_EMAIL` | no | Login email (with password) |
| `APITO_TENANT_PHONE` | no | Login phone (with password) |
| `APITO_TENANT_PASSWORD` | no | Login password |

```bash
go run ./examples/users/
```
