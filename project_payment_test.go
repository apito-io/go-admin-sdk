package goapitosdk

import (
	"context"
	"strings"
	"testing"
)

func TestDeriveSystemAPIBaseURL(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"http://localhost:5050/system/graphql", "http://localhost:5050"},
		{"http://localhost:5050/secured/graphql/", "http://localhost:5050"},
		{"http://localhost:5050/system/graphql/", "http://localhost:5050"},
	}
	for _, tc := range cases {
		if got := deriveSystemAPIBaseURL(tc.in); got != tc.want {
			t.Fatalf("deriveSystemAPIBaseURL(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestProjectPaymentAllowsAppUserSession(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "http://127.0.0.1:1/system/graphql",
		APIKey:  "session-jwt",
	})
	_, err := client.GetTenantSubscription(context.Background(), "project", "tenant")
	if err == nil || strings.Contains(err.Error(), "OPERATOR_CREDENTIAL_REQUIRED") {
		t.Fatalf("app-user session should pass the credential gate, got %v", err)
	}
}
