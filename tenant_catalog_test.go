package goapitosdk

import (
	"strings"
	"testing"
)

func TestCreateTenantParamsValidation(t *testing.T) {
	c := NewClient(Config{BaseURL: "http://localhost/system/graphql", APIKey: "test"})
	_, err := c.CreateTenant(t.Context(), CreateTenantParams{Name: ""})
	if err == nil || !strings.Contains(err.Error(), "name is required") {
		t.Fatalf("expected name required error, got %v", err)
	}
}

func TestUpdateTenantRequiresField(t *testing.T) {
	c := NewClient(Config{BaseURL: "http://localhost/system/graphql", APIKey: "test"})
	_, err := c.UpdateTenant(t.Context(), "tid", UpdateTenantParams{})
	if err == nil || !strings.Contains(err.Error(), "at least one field") {
		t.Fatalf("expected at least one field error, got %v", err)
	}
}

func TestDeleteTenantRequiresID(t *testing.T) {
	c := NewClient(Config{BaseURL: "http://localhost/system/graphql", APIKey: "test"})
	_, err := c.DeleteTenant(t.Context(), "")
	if err == nil || !strings.Contains(err.Error(), "tenantID is required") {
		t.Fatalf("expected tenantID required error, got %v", err)
	}
}

func TestMapToTenantCatalogSearchRow(t *testing.T) {
	row := mapToTenantCatalogSearchRow(map[string]interface{}{
		"id":         "01ABC",
		"name":       "Acme",
		"status":     "active",
		"domain":     "acme.example.com",
		"data":       `{"owner_uid":"u1"}`,
		"icon":       "https://cdn/logo.png",
		"created_at": "2026-07-12T00:00:00Z",
	})
	if row == nil || row.ID != "01ABC" || row.Icon == "" || row.CreatedAt == "" {
		t.Fatalf("unexpected row: %+v", row)
	}
}

func TestSearchTenantsRequiresProjectID(t *testing.T) {
	c := NewClient(Config{BaseURL: "http://localhost/system/graphql", APIKey: "test"})
	_, err := c.SearchTenants(t.Context(), "", 10, 0, "", "")
	if err == nil || !strings.Contains(err.Error(), "projectID is required") {
		t.Fatalf("expected projectID required error, got %v", err)
	}
}

func TestGetTenantRequiresIDs(t *testing.T) {
	c := NewClient(Config{BaseURL: "http://localhost/system/graphql", APIKey: "test"})
	_, err := c.GetTenant(t.Context(), "", "tid", "")
	if err == nil || !strings.Contains(err.Error(), "projectID is required") {
		t.Fatalf("expected projectID required error, got %v", err)
	}
	_, err = c.GetTenant(t.Context(), "proj", "", "")
	if err == nil || !strings.Contains(err.Error(), "tenantID is required") {
		t.Fatalf("expected tenantID required error, got %v", err)
	}
}

func TestMapToTenantCatalogListItem(t *testing.T) {
	row := mapToTenantCatalogListItem(map[string]interface{}{
		"id":     "01ABC",
		"name":   "user@example.com",
		"domain": "shop.example.com",
		"icon":   "https://cdn/logo.png",
		"data":   "{}",
	})
	if row == nil || row.ID != "01ABC" || row.Name != "user@example.com" {
		t.Fatalf("unexpected row: %+v", row)
	}
}
