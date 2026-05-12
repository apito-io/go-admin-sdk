# Tenant catalog users (Pro)

Minimal sample for **Pro** Apito: `SearchTenantUsers` and optional `LoginTenantUser`.

## Environment

| Variable | Required | Description |
|----------|----------|-------------|
| `APITO_API_KEY` | yes | Admin API key (`ak_...`) |
| `APITO_PROJECT_ID` | yes | Project id |
| `APITO_BASE_URL` | no | Defaults to `http://localhost:5050/system/graphql` |
| `APITO_TENANT_USERNAME` | no | If set with `APITO_TENANT_PASSWORD`, runs login after search |

## Run

From the **go-internal-sdk** repository root:

```bash
export APITO_API_KEY='ak_...'
export APITO_PROJECT_ID='your-project-id'
go run ./examples/tenant_users/
```

With login check:

```bash
export APITO_TENANT_USERNAME='admin'
export APITO_TENANT_PASSWORD='your-password'
go run ./examples/tenant_users/
```
