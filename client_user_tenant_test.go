package goapitosdk

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserCRUDTenantIDWiring(t *testing.T) {
	var lastVars map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var payload struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		lastVars = payload.Variables
		if !strings.Contains(payload.Query, "tenant_id") {
			t.Fatalf("query missing tenant_id: %s", payload.Query)
		}

		var data map[string]interface{}
		switch {
		case strings.Contains(payload.Query, "searchUsers"):
			data = map[string]interface{}{
				"searchUsers": map[string]interface{}{"count": 0, "users": []interface{}{}},
			}
		case strings.Contains(payload.Query, "createUser"):
			data = map[string]interface{}{
				"createUser": map[string]interface{}{"id": "u1", "role": "none", "tenant_id": "t1"},
			}
		case strings.Contains(payload.Query, "updateUser"):
			data = map[string]interface{}{
				"updateUser": map[string]interface{}{"id": "u1", "role": "vendor", "tenant_id": "t1"},
			}
		default:
			t.Fatalf("unexpected query: %s", payload.Query)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
	}))
	defer srv.Close()

	client := NewClient(Config{BaseURL: srv.URL, APIKey: "ak_test"})
	ctx := context.Background()

	t.Run("SearchUsers", func(t *testing.T) {
		lastVars = nil
		_, err := client.SearchUsers(ctx, "proj", 10, 0, "tenant-abc")
		if err != nil {
			t.Fatalf("SearchUsers: %v", err)
		}
		if lastVars["tenant_id"] != "tenant-abc" {
			t.Fatalf("tenant_id = %v", lastVars["tenant_id"])
		}
	})

	t.Run("SearchUsers omits empty tenant", func(t *testing.T) {
		lastVars = nil
		_, err := client.SearchUsers(ctx, "proj", 10, 0, "")
		if err != nil {
			t.Fatalf("SearchUsers: %v", err)
		}
		if _, ok := lastVars["tenant_id"]; ok {
			t.Fatalf("expected no tenant_id, got %v", lastVars["tenant_id"])
		}
	})

	t.Run("CreateUser", func(t *testing.T) {
		lastVars = nil
		_, err := client.CreateUser(ctx, "proj", CreateUserParams{
			Password: "secret",
			Email:    "a@b.com",
			TenantID: "tenant-xyz",
		})
		if err != nil {
			t.Fatalf("CreateUser: %v", err)
		}
		if lastVars["tenant_id"] != "tenant-xyz" {
			t.Fatalf("tenant_id = %v", lastVars["tenant_id"])
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		lastVars = nil
		tid := "tenant-xyz"
		_, err := client.UpdateUser(ctx, "u1", UpdateUserParams{TenantID: &tid})
		if err != nil {
			t.Fatalf("UpdateUser: %v", err)
		}
		if lastVars["tenant_id"] != "tenant-xyz" {
			t.Fatalf("tenant_id = %v", lastVars["tenant_id"])
		}
	})

	t.Run("UpdateUser requires field", func(t *testing.T) {
		_, err := client.UpdateUser(ctx, "u1", UpdateUserParams{})
		if err == nil || !strings.Contains(err.Error(), "at least one field") {
			t.Fatalf("expected validation error, got %v", err)
		}
	})
}
