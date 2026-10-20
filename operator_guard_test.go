package goapitosdk

import (
	"context"
	"strings"
	"testing"

	"github.com/apito-io/types"
)

func TestBootstrapCredentialRejectsPrivilegedSDKMethods(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "http://127.0.0.1:1/system/graphql",
		APIKey:  "ak_bootstrap",
	})
	ctx := context.Background()

	cases := []struct {
		name string
		call func() error
	}{
		{
			name: "search users",
			call: func() error {
				_, err := client.SearchUsers(ctx, "project", 1, 0, "", "")
				return err
			},
		},
		{
			name: "tenant catalog",
			call: func() error {
				_, err := client.SearchTenants(ctx, "project", 1, 0, "", "")
				return err
			},
		},
		{
			name: "user create",
			call: func() error {
				_, err := client.CreateUser(ctx, "project", CreateUserParams{Password: "password"})
				return err
			},
		},
		{
			name: "raw read",
			call: func() error {
				_, err := client.GetSingleResource(ctx, "model", "row", false)
				return err
			},
		},
		{
			name: "raw write",
			call: func() error {
				_, err := client.CreateNewResource(ctx, &types.CreateAndUpdateRequest{
					Model:   "model",
					Payload: map[string]interface{}{"value": "x"},
				})
				return err
			},
		},
		{
			name: "project payment",
			call: func() error {
				_, err := client.GetTenantSubscription(ctx, "project", "tenant")
				return err
			},
		},
		{
			name: "debug",
			call: func() error {
				_, err := client.Debug(ctx, "stage")
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call()
			if err == nil || !strings.Contains(err.Error(), "OPERATOR_CREDENTIAL_REQUIRED") {
				t.Fatalf("expected operator credential error, got %v", err)
			}
		})
	}
}

func TestOperatorCredentialStillReachesPrivilegedRequest(t *testing.T) {
	client := NewClient(Config{
		BaseURL:     "http://127.0.0.1:1/system/graphql",
		AccessToken: "apt_operator",
	})
	if err := client.requireOperatorCredential("test"); err != nil {
		t.Fatalf("apt_ credential rejected: %v", err)
	}
}
