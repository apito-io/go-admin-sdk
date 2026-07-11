// Project users example. Run from repo root: go run ./examples/users/
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	goapitosdk "github.com/apito-io/go-admin-sdk"
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

	list, err := client.SearchUsers(ctx, projectID, 50, 0, "", "")
	if err != nil {
		log.Fatalf("SearchUsers: %v", err)
	}
	fmt.Printf("Users (count=%d):\n", list.Count)
	for _, u := range list.Users {
		idLabel := strings.TrimSpace(u.Email)
		if idLabel == "" {
			idLabel = strings.TrimSpace(u.Phone)
		}
		if idLabel == "" {
			idLabel = "(no email/phone)"
		}
		fmt.Printf("  - %s (%s) role=%s status=%s\n", idLabel, u.ID, u.Role, u.Status)
	}

	email := strings.TrimSpace(os.Getenv("APITO_TENANT_EMAIL"))
	phone := strings.TrimSpace(os.Getenv("APITO_TENANT_PHONE"))
	pw := os.Getenv("APITO_TENANT_PASSWORD")
	if (email != "" || phone != "") && pw != "" {
		login, err := client.LoginUser(ctx, projectID, goapitosdk.LoginUserParams{
			Password: pw,
			Email:    email,
			Phone:    phone,
		})
		if err != nil {
			log.Fatalf("LoginUser: %v", err)
		}
		fmt.Printf("Login OK, token length=%d\n", len(login.Token))
	}
}
