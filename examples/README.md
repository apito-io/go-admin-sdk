# Go Apito SDK - Todo Example

This is a comprehensive example demonstrating all features of the Go Apito SDK through a practical todo application.

## 🚀 Features Demonstrated

- 🔐 Authentication & tenant token generation
- 📝 Creating resources (todos, users, categories)
- 🔍 Searching with both typed and untyped methods
- 📄 Getting single resources
- ✏️ Updating resources
- 🔗 Getting related documents
- 📊 Audit logging
- 🐛 Debug functionality
- 🗑️ Resource cleanup

## 🛠️ Setup

### Environment Variables

Set the following environment variables before running the example:

```bash
# Required
export APITO_BASE_URL="https://api.apito.io/graphql"
export APITO_API_KEY="your-api-key-here"

# Optional (for multi-tenant features)
export APITO_TENANT_ID="your-tenant-id"
export APITO_AUTH_TOKEN="your-auth-token"
```

### Build and Run

```bash
# Build the example
go build -o todo-example main.go

# Run the example
./todo-example

# Or run directly
go run main.go
```

## 📋 Example Output

When you run the example, you'll see output like this:

```
🚀 Apito SDK Comprehensive Todo Example
========================================

🔐 1. Authentication & Tenant Token Generation
✅ Generated tenant token: abcd1234567890...

📝 2. Creating Resources
✅ Created category: cat_123
✅ Created user: user_456
✅ Created todo: todo_789 (Implement user authentication)
✅ Created todo: todo_790 (Write unit tests)
✅ Created todo: todo_791 (Update documentation)

🔍 3. Searching Resources
✅ Found 3 todos (untyped search)
✅ Found 3 todos (typed search)
   - todo_789: Implement user authentication (Status: todo, Priority: high)
   - todo_790: Write unit tests (Status: in_progress, Priority: medium)
   - todo_791: Update documentation (Status: todo, Priority: low)
✅ Found 1 active users
   - user_456: John Doe (john.doe@example.com)

📄 4. Getting Single Resources
✅ Retrieved todo (untyped): todo_789
✅ Retrieved todo (typed): todo_789 - Implement user authentication

✏️  5. Updating Resources
✅ Updated todo status: todo_789
✅ Updated todo (typed): todo_789 - Status: in_progress

🔗 6. Getting Related Documents
✅ Found 3 todos related to user user_456
✅ Found 3 todos related to user (typed)
   - todo_789: Implement user authentication
   - todo_790: Write unit tests
   - todo_791: Update documentation

📊 7. Audit Logging
✅ Audit log sent successfully

🐛 8. Debug Functionality
✅ Debug info sent: map[message:Debug received data:...]

🗑️  9. Cleanup (Deleting Resources)
✅ Deleted todo: todo_789
✅ Deleted todo: todo_790
✅ Deleted todo: todo_791
✅ Deleted user: user_456
✅ Deleted category: cat_123

🎉 Todo Example Completed Successfully!
=====================================

This example demonstrated:
• Authentication & tenant token generation
• Creating resources (todos, users, categories)
• Searching with both typed and untyped methods
• Getting single resources
• Updating resources
• Getting related documents
• Audit logging
• Debug functionality
• Resource cleanup
```

## 🔧 Code Structure

The example demonstrates both **untyped** and **type-safe** operations:

### Type-Safe Operations (Recommended)

```go
// Define your data structures
type Todo struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    Priority    string    `json:"priority"`
    // ... more fields
}

// Use type-safe operations
typedResults, err := goapitosdk.SearchResourcesTyped[Todo](client, ctx, "todos", filter, false)
if err != nil {
    log.Fatal(err)
}

for _, todoDoc := range typedResults.Results {
    fmt.Printf("Todo: %s (Status: %s)\n", todoDoc.Data.Title, todoDoc.Data.Status)
}
```

### Untyped Operations (Flexible)

```go
// Use untyped operations for dynamic data
results, err := client.SearchResources(ctx, "todos", filter, false)
if err != nil {
    log.Fatal(err)
}

for _, todo := range results.Results {
    fmt.Printf("Todo ID: %s, Data: %v\n", todo.ID, todo.Data)
}
```

## 🧪 Testing with Mock Data

If you don't have a real Apito backend, the example will handle errors gracefully and still demonstrate the API structure. The example includes comprehensive error handling and will show you what operations are being attempted.

## 🔗 Related Documentation

- [Main SDK Documentation](../../README.md)
- [Go Apito SDK API Reference](https://pkg.go.dev/github.com/apito-io/go-admin-sdk)
- [Apito Platform Documentation](https://docs.apito.io)

## 💡 Next Steps

1. Modify the example to work with your specific data models
2. Explore the type-safe operations for better development experience
3. Check out the plugin integration examples in the main documentation
4. Review the production deployment guides
