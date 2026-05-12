package goapitosdk

// TenantUser is a row from the engine system control plane (pro_tenant_users).
type TenantUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TenantID  string `json:"tenant_id"`
	Provider  string `json:"provider"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// TenantLoginResponse is returned by loginTenantUser and loginTenantUserGoogle system operations.
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
