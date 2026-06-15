# Go Apito SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/apito-io/go-admin-sdk.svg)](https://pkg.go.dev/github.com/apito-io/go-admin-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/apito-io/go-admin-sdk)](https://goreportcard.com/report/github.com/apito-io/go-admin-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive Go SDK for communicating with Apito GraphQL API endpoints. This SDK implements the `InjectedDBOperationInterface` and provides both type-safe and flexible interfaces for interacting with Apito's backend services.

## 🚀 Features

- ✅ **Complete SDK Implementation**: Full implementation of `InjectedDBOperationInterface`
- ✅ **Type-Safe Operations**: Generic typed methods for better development experience
- ✅ **GraphQL-Based**: Native GraphQL communication with Apito backend
- ✅ **Authentication Ready**: API key and tenant-based authentication
- ✅ **Context-Aware**: Full context support with timeout and cancellation
- ✅ **Comprehensive Error Handling**: Detailed error responses and GraphQL error support
- ✅ **Plugin-Ready**: Perfect for HashiCorp Go plugins and microservices
- ✅ **Production Ready**: Battle-tested in production environments

## 📦 Installation

```bash
go get github.com/apito-io/go-admin-sdk
```

When building this repository from a checkout that vendors **`github.com/apito-io/types`** via `replace ... => ../types`, keep the **`types`** module cloned as a sibling directory (`../types` relative to this module root). Remove the `replace` line after upgrading **`go.mod`** to a published **`types`** release that includes the matching **`InternalSDKOperation.GenerateTenantToken`** signature.

## 🎯 Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    goapitosdk "github.com/apito-io/go-admin-sdk"
)

func main() {
    // Create a new client
    client := goapitosdk.NewClient(goapitosdk.Config{
        BaseURL: "https://api.apito.io/graphql",
        APIKey:  "your-api-key-here",
        Timeout: 30 * time.Second,
    })

    ctx := context.Background()

    // Create a new todo
    todoData := map[string]interface{}{
        "title":       "Learn Apito SDK",
        "description": "Complete the SDK tutorial",
        "status":      "todo",
        "priority":    "high",
    }

    request := &goapitosdk.CreateAndUpdateRequest{
        Model:   "todos",
        Payload: todoData,
    }

    todo, err := client.CreateNewResource(ctx, request)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created todo: %s\n", todo.ID)
}
```

## ⚙️ Configuration

### Basic Configuration

```go
client := goapitosdk.NewClient(goapitosdk.Config{
    BaseURL: "https://api.apito.io/graphql",  // Your Apito GraphQL endpoint
    APIKey:  "your-api-key-here",             // X-APITO-KEY header value
    Timeout: 30 * time.Second,                // HTTP client timeout
})
```

### Advanced Configuration

```go
// Custom HTTP client with specific settings
customClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

client := goapitosdk.NewClient(goapitosdk.Config{
    BaseURL:    "https://api.apito.io/graphql",
    APIKey:     "your-api-key-here",
    HTTPClient: customClient,
})
```

### Context with Tenant ID

```go
ctx := context.Background()
ctx = context.WithValue(ctx, "tenant_id", "your-tenant-id")

// All operations will now include the tenant ID
results, err := client.SearchResources(ctx, "users", filter, false)
```

## 📚 Complete API Reference

### 🔐 Authentication

#### Generate Tenant Token

Generate a new tenant token for multi-tenant operations. Arguments match the engine `generateTenantToken` mutation: **`tenantID`**, **`duration`** (`YYYY-MM-DD`; empty string uses one year ahead in UTC), optional **`role`** (empty uses engine default `admin`). Auth uses the client API key, not a legacy token parameter.

```go
tenantToken, err := client.GenerateTenantToken(ctx, "tenant-catalog-id", "2027-12-31", "")
if err != nil {
    log.Fatal(err)
}
fmt.Println("Generated token:", tenantToken)
```

#### Project users (system GraphQL)

Requires system GraphQL (`/system/graphql`) with an admin API key. Operations are project-first (`project_id`). On Pro/SaaS engines, user payloads include `tenant_id`.

| Method | Description |
|--------|-------------|
| `LoginUser(ctx, projectID, LoginUserParams)` | General: **`Password`** + **`Email`** or **`Phone`**. SaaS per-tenant DB: **`TenantID` required**. Google **code** flow: **`AuthMethod: "google"`**, **`Code`**, **`State`**; use **`GoogleOAuthState`** first for **`State`**. Google login may link a verified email to an existing user instead of creating a duplicate. |
| `GoogleOAuthState(ctx, projectID)` | Returns signed OAuth **`State`** string for building the Google authorize URL. |
| `SearchUsers(ctx, projectID, limit, offset)` | List project end-users (`email`, `phone`, optional `tenant_id`). |
| `SearchTenantsByDomain(ctx, projectID, domain)` | Resolve the single SaaS catalog tenant for an exact domain match in the project (`tenant` null if none). |
| `CreateUser(ctx, projectID, CreateUserParams)` | Create a local-password user; **`Password`**, optional **`Role`**, **`Email`**, **`Phone`**. Engine rejects duplicate email/phone project-wide. |
| `UpdateUser(ctx, userID, UpdateUserParams)` | Update **`email`**, **`phone`**, **`role`** (non-nil **`*string` pointers only**). Engine rejects duplicate email/phone project-wide. |
| `ResetUserPassword(ctx, userID, password)` | Set a new password (admin mutation). |
| `DeleteUser(ctx, userID)` | Hard-delete a user (returns bool from GraphQL). |

#### Files (REST)

REST base URL is derived from `Config.BaseURL` by stripping `/graphql` (e.g. `http://host:5050/system`), or set `Config.RestBaseURL` explicitly.

File **metadata** is stored in the **project database** `files` table (not the system DB). On Pro/SaaS engines, set `tenant_id` on the request context (`ctx` value) so operations target the tenant project DB. Default `RestBaseURL` resolves to `/secured` when GraphQL uses `/system/graphql`; full URLs are `/secured/files/upload`, `/secured/files/list`, and `/secured/files/delete`.

| Method | Description |
|--------|-------------|
| `UploadFile(ctx, params)` | POST `/files/upload` (multipart). |
| `ListFiles(ctx, fileType, limit, offset)` | GET `/files/list`. |
| `DeleteFiles(ctx, ids)` | POST `/files/delete`. |

Path constants: `FilesUploadPath`, `FilesListPath`, `FilesDeletePath`.

On the engine system GraphQL API, `createTenant` accepts an optional `domain` argument; when it is set, the engine requires that domain to be unused in the project (otherwise the mutation fails with a clear error). `updateTenant` validates the same when changing `domain`. Use `executeGraphQL` if you need those catalog mutations from the SDK.

```go
// List users by project
list, err := client.SearchUsers(ctx, "project-id", 50, 0)
if err != nil {
    log.Fatal(err)
}
for _, u := range list.Users {
    label := u.Email
    if label == "" {
        label = u.Phone
    }
    if label == "" {
        label = "(no email/phone)"
    }
    fmt.Println(label, u.Role, u.Status)
}

// Login (returns token + user on success)
login, err := client.LoginUser(ctx, "project-id", goapitosdk.LoginUserParams{
    Password: "secret",
    Email:    "user@example.com", // use Phone: "+1555..." when project uses phone identifier
})
if err != nil {
    log.Fatal(err)
}
if login.Token != "" {
    fmt.Println("tenant token:", login.Token)
}
```

See also: `examples/users/main.go`, `examples/files/main.go`.

### 📝 Resource Management

#### Create New Resource

**Untyped Creation:**

```go
request := &goapitosdk.CreateAndUpdateRequest{
    Model: "users",
    Payload: map[string]interface{}{
        "name":   "John Doe",
        "email":  "john@example.com",
        "active": true,
    },
    Connect: map[string]interface{}{
        "organization_id": "org-123",
    },
}

user, err := client.CreateNewResource(ctx, request)
```

**Type-Safe Creation:**

```go
type User struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Email  string `json:"email"`
    Active bool   `json:"active"`
}

typedUser, err := goapitosdk.CreateNewResourceTyped[User](client, ctx, request)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created user: %s (%s)\n", typedUser.Data.Name, typedUser.Data.Email)
```

#### Update Resource

```go
updateRequest := &goapitosdk.CreateAndUpdateRequest{
    ID:    "user-123",
    Model: "users",
    Payload: map[string]interface{}{
        "name": "Jane Doe Updated",
    },
    Connect: map[string]interface{}{
        "role_id": "role-456",
    },
    Disconnect: map[string]interface{}{
        "old_role_id": "role-123",
    },
    ForceUpdate: false,
}

updatedUser, err := client.UpdateResource(ctx, updateRequest)
```

#### Delete Resource

```go
err := client.DeleteResource(ctx, "users", "user-123")
if err != nil {
    log.Fatal(err)
}
```

### 🔍 Search & Retrieval

#### Search Resources

**Basic Search:**

```go
filter := map[string]interface{}{
    "limit": 10,
    "page":  1,
    "where": map[string]interface{}{
        "status": "active",
        "role":   "admin",
    },
    "search": "john@example.com",
}

results, err := client.SearchResources(ctx, "users", filter, false)
```

**Type-Safe Search:**

```go
typedResults, err := goapitosdk.SearchResourcesTyped[User](client, ctx, "users", filter, false)
if err != nil {
    log.Fatal(err)
}

for _, userDoc := range typedResults.Results {
    fmt.Printf("User: %s (%s)\n", userDoc.Data.Name, userDoc.Data.Email)
}
```

**Advanced Filtering:**

```go
advancedFilter := map[string]interface{}{
    "limit":  20,
    "offset": 10,
    "where": map[string]interface{}{
        "created_at": map[string]interface{}{
            "$gte": "2024-01-01T00:00:00Z",
        },
        "status": map[string]interface{}{
            "$in": []string{"active", "pending"},
        },
    },
    "sort": map[string]interface{}{
        "created_at": -1, // Descending order
    },
}

results, err := client.SearchResources(ctx, "users", advancedFilter, false)
```

#### Get Single Resource

**Untyped Retrieval:**

```go
user, err := client.GetSingleResource(ctx, "users", "user-123", false)
if err != nil {
    log.Fatal(err)
}
```

**Type-Safe Retrieval:**

```go
typedUser, err := goapitosdk.GetSingleResourceTyped[User](client, ctx, "users", "user-123", false)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("User: %s\n", typedUser.Data.Name)
```

#### Get Related Documents

```go
relationConnection := map[string]interface{}{
    "model": "todos",
    "filter": map[string]interface{}{
        "limit": 10,
        "where": map[string]interface{}{
            "status": "pending",
        },
    },
}

// Get todos related to a user
relatedTodos, err := client.GetRelationDocuments(ctx, "user-123", relationConnection)
if err != nil {
    log.Fatal(err)
}

// Type-safe version
typedTodos, err := goapitosdk.GetRelationDocumentsTyped[Todo](client, ctx, "user-123", relationConnection)
```

### 📊 Audit & Debug

#### Send Audit Log

```go
auditData := goapitosdk.AuditData{
    Resource: "users",
    Action:   "create",
    Author: map[string]interface{}{
        "user_id": "admin-123",
        "name":    "Admin User",
    },
    Data: map[string]interface{}{
        "user_id": "user-456",
        "email":   "newuser@example.com",
    },
    Meta: map[string]interface{}{
        "ip_address": "192.168.1.1",
        "user_agent": "Apito-SDK/1.0",
        "timestamp":  time.Now().Format(time.RFC3339),
    },
}

err := client.SendAuditLog(ctx, auditData)
if err != nil {
    log.Fatal(err)
}
```

#### Debug Operations

```go
debugData := map[string]interface{}{
    "operation": "user_creation",
    "duration":  "150ms",
    "success":   true,
}

result, err := client.Debug(ctx, "user_management", debugData)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Debug result: %+v\n", result)
```

## 🎯 Complete Todo Example

The SDK includes a comprehensive todo application example that demonstrates all features:

```bash
# Set environment variables
export APITO_BASE_URL="https://api.apito.io/graphql"
export APITO_API_KEY="your-api-key"
export APITO_TENANT_ID="your-tenant-id"  # Optional (for tenant-token generation only)
export APITO_AUTH_TOKEN="your-auth-token"  # Optional for token generation

# Run the example
cd examples/basic
go run main.go
```

The example demonstrates:

- 🔐 Authentication & tenant token generation
- 📝 Creating resources (todos, users, categories)
- 🔍 Searching with both typed and untyped methods
- 📄 Getting single resources
- ✏️ Updating resources
- 🔗 Getting related documents
- 📊 Audit logging
- 🐛 Debug functionality
- 🗑️ Resource cleanup

### Pro: tenant catalog users

```bash
export APITO_API_KEY="your-api-key"
export APITO_PROJECT_ID="your-project-id"
# Optional: export APITO_TENANT_EMAIL and/or APITO_TENANT_PHONE / APITO_TENANT_PASSWORD to test login
go run ./examples/tenant_users/
```

Details: `examples/tenant_users/README.md`.

## 🏗️ Type System

### Defining Custom Types

```go
// Define your data structures
type Todo struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    Priority    string    `json:"priority"`
    DueDate     time.Time `json:"due_date"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type User struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    Active   bool   `json:"active"`
}
```

### Type-Safe Operations

All operations have type-safe counterparts:

```go
// Type-safe alternatives
GetSingleResourceTyped[T](client, ctx, model, id, singlePageData)
SearchResourcesTyped[T](client, ctx, model, filter, aggregate)
GetRelationDocumentsTyped[T](client, ctx, id, connection)
CreateNewResourceTyped[T](client, ctx, request)
UpdateResourceTyped[T](client, ctx, request)
```

## 🔌 Plugin Integration

### HashiCorp Go Plugin Usage

```go
// In your plugin
type MyPlugin struct {
    client goapitosdk.InjectedDBOperationInterface
}

func (p *MyPlugin) Initialize(client goapitosdk.InjectedDBOperationInterface) {
    p.client = client
}

func (p *MyPlugin) ProcessData(ctx context.Context) error {
    // Use the client for database operations
    results, err := p.client.SearchResources(ctx, "data", filter, false)
    if err != nil {
        return err
    }

    // Process results...
    return nil
}
```

### Microservice Integration

```go
// In your microservice
type UserService struct {
    apitoClient *goapitosdk.Client
}

func NewUserService(config goapitosdk.Config) *UserService {
    return &UserService{
        apitoClient: goapitosdk.NewClient(config),
    }
}

func (s *UserService) CreateUser(ctx context.Context, userData User) (*User, error) {
    request := &goapitosdk.CreateAndUpdateRequest{
        Model:   "users",
        Payload: structToMap(userData),
    }

    result, err := goapitosdk.CreateNewResourceTyped[User](s.apitoClient, ctx, request)
    if err != nil {
        return nil, err
    }

    return &result.Data, nil
}
```

## 🔧 Error Handling

### GraphQL Errors

```go
results, err := client.SearchResources(ctx, "users", filter, false)
if err != nil {
    // Check if it's a GraphQL error
    if graphqlErr, ok := err.(*goapitosdk.GraphQLError); ok {
        fmt.Printf("GraphQL Error: %s\n", graphqlErr.Message)
        fmt.Printf("Path: %v\n", graphqlErr.Path)
        fmt.Printf("Extensions: %v\n", graphqlErr.Extensions)
    } else {
        // Handle other errors (HTTP, network, etc.)
        fmt.Printf("Error: %v\n", err)
    }
}
```

### HTTP Errors

```go
// Handle HTTP-level errors
client := goapitosdk.NewClient(goapitosdk.Config{
    BaseURL: "https://api.apito.io/graphql",
    APIKey:  "invalid-key",
    Timeout: 5 * time.Second,
})

_, err := client.SearchResources(ctx, "users", nil, false)
if err != nil {
    if strings.Contains(err.Error(), "HTTP error 401") {
        fmt.Println("Authentication failed - check your API key")
    } else if strings.Contains(err.Error(), "HTTP error 403") {
        fmt.Println("Authorization failed - check your permissions")
    }
}
```

## 🧪 Testing

### Mock Client

```go
// For testing, you can implement the interface
type MockClient struct{}

func (m *MockClient) SearchResources(ctx context.Context, model string, filter map[string]interface{}, aggregate bool) (*goapitosdk.SearchResult, error) {
    // Return mock data
    return &goapitosdk.SearchResult{
        Results: []*shared.DefaultDocumentStructure{
            {ID: "test-1", Data: map[string]interface{}{"name": "Test User"}},
        },
        Count: 1,
    }, nil
}

// Use in tests
func TestUserService(t *testing.T) {
    service := &UserService{apitoClient: &MockClient{}}
    // Test your service...
}
```

## 📈 Performance Tips

### Connection Pooling

```go
// Configure HTTP client for better performance
client := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:          100,
        MaxIdleConnsPerHost:   10,
        IdleConnTimeout:       90 * time.Second,
        TLSHandshakeTimeout:   10 * time.Second,
        ExpectContinueTimeout: 1 * time.Second,
    },
    Timeout: 30 * time.Second,
}

apitoClient := goapitosdk.NewClient(goapitosdk.Config{
    BaseURL:    "https://api.apito.io/graphql",
    APIKey:     "your-api-key",
    HTTPClient: client,
})
```

### Batch Operations

```go
// Instead of multiple individual requests, batch them
var wg sync.WaitGroup
results := make(chan *goapitosdk.SearchResult, 10)

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(page int) {
        defer wg.Done()
        filter := map[string]interface{}{"page": page, "limit": 100}
        result, err := client.SearchResources(ctx, "users", filter, false)
        if err == nil {
            results <- result
        }
    }(i)
}

go func() {
    wg.Wait()
    close(results)
}()

// Process results as they come in
for result := range results {
    // Process each batch...
}
```

## 🚀 Production Deployment

### Environment Variables

```bash
# Required
APITO_BASE_URL=https://api.apito.io/graphql
APITO_API_KEY=your-production-api-key

# Optional
APITO_TENANT_ID=your-tenant-id
APITO_AUTH_TOKEN=your-auth-token
APITO_TIMEOUT=30s
```

### Docker Configuration

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .
CMD ["./app"]
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: apito-sdk-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: apito-sdk-app
  template:
    metadata:
      labels:
        app: apito-sdk-app
    spec:
      containers:
        - name: app
          image: your-app:latest
          env:
            - name: APITO_BASE_URL
              value: "https://api.apito.io/graphql"
            - name: APITO_API_KEY
              valueFrom:
                secretKeyRef:
                  name: apito-secrets
                  key: api-key
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

```bash
git clone https://github.com/apito-io/go-admin-sdk.git
cd go-apito-sdk
go mod download
go test ./...
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Apito Documentation](https://docs.apito.io)
- [API Reference](https://pkg.go.dev/github.com/apito-io/go-admin-sdk)
- [GitHub Repository](https://github.com/apito-io/go-admin-sdk)
- [Issues](https://github.com/apito-io/go-admin-sdk/issues)

## 🆘 Support

- 📧 Email: support@apito.io
- 💬 Discord: [Join our community](https://discord.gg/apito)
- 📖 Documentation: [docs.apito.io](https://docs.apito.io)
- 🐛 Bug Reports: [GitHub Issues](https://github.com/apito-io/go-admin-sdk/issues)
