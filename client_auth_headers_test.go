package goapitosdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApplyAuthCredential_AptTokenSendsBearerOnly(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://example.test", nil)
	applyAuthCredential(req, "apt_abc123")

	if got := req.Header.Get("Authorization"); got != "Bearer apt_abc123" {
		t.Fatalf("Authorization = %q, want %q", got, "Bearer apt_abc123")
	}
	if got := req.Header.Get("X-Use-Cookies"); got != "false" {
		t.Fatalf("X-Use-Cookies = %q, want %q", got, "false")
	}
	if got := req.Header.Get("X-Apito-Key"); got != "" {
		t.Fatalf("X-Apito-Key should be empty (no dual header), got %q", got)
	}
}

func TestApplyAuthCredential_LegacyProjectKeyUsesXApitoKey(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://example.test", nil)
	applyAuthCredential(req, "ak_project_key")

	if got := req.Header.Get("X-Apito-Key"); got != "ak_project_key" {
		t.Fatalf("X-Apito-Key = %q, want %q", got, "ak_project_key")
	}
	if got := req.Header.Get("Authorization"); got != "" {
		t.Fatalf("Authorization should be empty for legacy project keys, got %q", got)
	}
}

func TestIsRetiredTokenPrefix(t *testing.T) {
	for _, prefix := range []string{"cli-", "sdk-", "mcp-"} {
		if !isRetiredTokenPrefix(prefix + "legacy") {
			t.Errorf("isRetiredTokenPrefix(%q) = false, want true", prefix+"legacy")
		}
	}
	if isRetiredTokenPrefix("apt_ok") {
		t.Errorf("isRetiredTokenPrefix(apt_ok) = true, want false")
	}
}

func TestExecuteREST_RejectsRetiredTokenPrefixes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("engine should not be called for a retired token prefix")
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL + "/system/graphql", APIKey: "mcp-legacy"})
	_, _, err := client.executeREST(context.Background(), restRequest{method: "GET", path: "/files"})
	if err == nil || !strings.Contains(err.Error(), "TOKEN_FORMAT_RETIRED") {
		t.Fatalf("executeREST: err = %v, want TOKEN_FORMAT_RETIRED", err)
	}
}

func TestExecuteGraphQL_RejectsRetiredTokenPrefixes(t *testing.T) {
	// Server should never be hit for a retired-prefix token; fail the test if it is.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("engine should not be called for a retired token prefix")
	}))
	defer server.Close()

	for _, prefix := range []string{"cli-", "sdk-", "mcp-"} {
		client := NewClient(Config{BaseURL: server.URL, APIKey: prefix + "legacy"})
		_, err := client.executeGraphQL(context.Background(), "query { __typename }", nil)
		if err == nil || !strings.Contains(err.Error(), "TOKEN_FORMAT_RETIRED") {
			t.Fatalf("executeGraphQL with %q prefix: err = %v, want TOKEN_FORMAT_RETIRED", prefix, err)
		}
	}
}

func TestExecuteGraphQL_ConfigProjectAndTenantHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get(headerApitoProjectID); got != "project-config" {
			t.Fatalf("%s = %q, want project-config", headerApitoProjectID, got)
		}
		if got := r.Header.Get("X-Apito-Tenant-ID"); got != "tenant-context" {
			t.Fatalf("X-Apito-Tenant-ID = %q, want tenant-context", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{}})
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL:     server.URL,
		AccessToken: "apt_test",
		ProjectID:   "project-config",
	})
	ctx := context.WithValue(context.Background(), "tenant_id", "tenant-context")
	if _, err := client.executeGraphQL(ctx, "query { __typename }", nil); err != nil {
		t.Fatalf("executeGraphQL: %v", err)
	}
}

func TestSearchUsers_ExplicitProjectOverridesConfigHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get(headerApitoProjectID); got != "project-method" {
			t.Fatalf("%s = %q, want project-method", headerApitoProjectID, got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"searchUsers": map[string]interface{}{"count": 0, "users": []interface{}{}},
			},
		})
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL:     server.URL,
		AccessToken: "apt_test",
		ProjectID:   "project-config",
	})
	if _, err := client.SearchUsers(context.Background(), "project-method", 10, 0, "", ""); err != nil {
		t.Fatalf("SearchUsers: %v", err)
	}
}
