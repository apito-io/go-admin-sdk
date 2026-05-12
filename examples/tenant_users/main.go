// Tenant catalog users example (Apito Pro). Run from repo root: go run ./examples/tenant_users/
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	goapitosdk "github.com/apito-io/go-internal-sdk"
)

func main() {
	baseURL := os.Getenv("APITO_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:5050/system/graphql"
	}
	apiKey := os.Getenv("APITO_API_KEY")
	if apiKey == "" {
		log.Fatal("APITO_API_KEY is required")
	}
	projectID := os.Getenv("APITO_PROJECT_ID")
	if projectID == "" {
		log.Fatal("APITO_PROJECT_ID is required")
	}

	client := goapitosdk.NewClient(goapitosdk.Config{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Timeout: 30 * time.Second,
	})
	ctx := context.Background()

	list, err := client.SearchTenantUsers(ctx, projectID, 50, 0)
	if err != nil {
		log.Fatalf("SearchTenantUsers: %v", err)
	}
	fmt.Printf("Tenant users (count=%d):\n", list.Count)
	for _, u := range list.Users {
		fmt.Printf("  - %s (%s) role=%s status=%s\n", u.Username, u.ID, u.Role, u.Status)
	}

	user := os.Getenv("APITO_TENANT_USERNAME")
	pw := os.Getenv("APITO_TENANT_PASSWORD")
	if user != "" && pw != "" {
		login, err := client.LoginTenantUser(ctx, projectID, user, pw)
		if err != nil {
			log.Fatalf("LoginTenantUser: %v", err)
		}
		fmt.Printf("Login OK, token length=%d\n", len(login.Token))
	}
}
