package goapitosdk

// TenantUser is a row from the engine system control plane (pro_tenant_users).
type TenantUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	TenantID  string `json:"tenant_id"`
	Provider  string `json:"provider"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// LoginTenantUserParams configures login via system GraphQL loginTenantUser.
// Password path (AuthMethod empty or "general"): set Password plus Email or Phone per project Authentication.
// Google path (AuthMethod "google"): set Code and State from OAuth callback; optionally use TenantGoogleOAuthState first.
type LoginTenantUserParams struct {
	Password   string
	Email      string
	Phone      string
	AuthMethod string // optional; "", "general", or "google"
	Code       string // OAuth authorization code (Google)
	State      string // OAuth state (from TenantGoogleOAuthState or callback)
}

// TenantGoogleOAuthStateResponse is returned by tenantGoogleOAuthState.
type TenantGoogleOAuthStateResponse struct {
	State string
}

// CreateTenantUserParams configures createTenantUser. The engine requires an email or phone
// according to the project's general authentication identifier mode.
type CreateTenantUserParams struct {
	Password string
	Role     string // optional; engine defaults when empty
	Email    string
	Phone    string
}

// UpdateTenantUserParams lists optional fields for updateTenantUser. Nil pointers are omitted from the mutation.
type UpdateTenantUserParams struct {
	Email    *string
	Phone    *string
	Password *string
	Role     *string
}

// TenantLoginResponse is returned by loginTenantUser (general or Google code flow).
type TenantLoginResponse struct {
	Token string      `json:"token"`
	User  *TenantUser `json:"user"`
}

// TenantUsersResponse is returned by searchTenantUsers.
type TenantUsersResponse struct {
	Users []*TenantUser `json:"users"`
	Count int           `json:"count"`
}

// TenantCatalogSearchRow is one catalog tenant row from searchTenantsByDomain.
type TenantCatalogSearchRow struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Domain string `json:"domain"`
	Data   string `json:"data"`
}

// TenantByDomainResponse is returned by searchTenantsByDomain (at most one match per project).
type TenantByDomainResponse struct {
	Tenant *TenantCatalogSearchRow `json:"tenant"`
}
