package goapitosdk

// User is a project end-user from the engine system DB (table project_users).
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	TenantID  string `json:"tenant_id,omitempty"`
	Provider  string `json:"provider"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// LoginUserParams configures login via system GraphQL loginUser.
// Password path (AuthMethod empty or "general"): set Password plus Email or Phone per project Authentication.
// Google path (AuthMethod "google"): set Code and State from OAuth callback; optionally use GoogleOAuthState first.
type LoginUserParams struct {
	Password   string
	Email      string
	Phone      string
	AuthMethod string // optional; "", "general", or "google"
	Code       string // OAuth authorization code (Google)
	State      string // OAuth state (from GoogleOAuthState or callback)
}

// GoogleOAuthStateResponse is returned by googleOAuthState.
type GoogleOAuthStateResponse struct {
	State string
}

// CreateUserParams configures createUser. The engine requires an email or phone
// according to the project's general authentication identifier mode.
type CreateUserParams struct {
	Password string
	Role     string // optional; engine defaults when empty
	Email    string
	Phone    string
}

// UpdateUserParams lists optional fields for updateUser. Nil pointers are omitted from the mutation.
type UpdateUserParams struct {
	Email *string
	Phone *string
	Role  *string
}

// LoginUserResponse is returned by loginUser (general or Google code flow).
type LoginUserResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// UsersResponse is returned by searchUsers.
type UsersResponse struct {
	Users []*User `json:"users"`
	Count int     `json:"count"`
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

// File is metadata for a file uploaded via the system files REST API.
type File struct {
	ID            string `json:"id"`
	FileType      string `json:"file_type"`
	FileName      string `json:"file_name"`
	FileExtension string `json:"file_extension,omitempty"`
	ContentType   string `json:"content_type,omitempty"`
	Size          int64  `json:"size"`
	URL           string `json:"url"`
	CreatedBy     string `json:"created_by,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
}

// FilesListResponse is returned by ListFiles.
type FilesListResponse struct {
	Files []File `json:"files"`
	Total int    `json:"total"`
}

// UploadFileParams configures UploadFile.
type UploadFileParams struct {
	FileName string
	Content  []byte
	FileType string // optional; inferred from content type when empty
}

// DeleteFilesResponse is returned by DeleteFiles.
type DeleteFilesResponse struct {
	Success        bool     `json:"success"`
	DeletedIDs     []string `json:"deleted_ids"`
	StorageFailed  []string `json:"storage_failed,omitempty"`
	Message        string   `json:"message,omitempty"`
}
