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

// ProjectStorageSettings is the read shape of project storage settings (secrets never returned).
type ProjectStorageSettings struct {
	UseFreeCloudStorage bool    `json:"use_free_cloud_storage"`
	Endpoint            *string `json:"endpoint,omitempty"`
	Region              *string `json:"region,omitempty"`
	Bucket              *string `json:"bucket,omitempty"`
	AccessKeyID         *string `json:"access_key_id,omitempty"`
	HasSecretAccessKey  bool    `json:"has_secret_access_key"`
	PublicBaseURL       *string `json:"public_base_url,omitempty"`
	ForcePathStyle      *bool   `json:"force_path_style,omitempty"`
}

// UpdateProjectStorageInput is the write shape for updateProjectStorageSettings.
type UpdateProjectStorageInput struct {
	UseFreeCloudStorage *bool   `json:"use_free_cloud_storage,omitempty"`
	Endpoint            *string `json:"endpoint,omitempty"`
	Region              *string `json:"region,omitempty"`
	Bucket              *string `json:"bucket,omitempty"`
	AccessKeyID         *string `json:"access_key_id,omitempty"`
	SecretAccessKey     *string `json:"secret_access_key,omitempty"`
	PublicBaseURL       *string `json:"public_base_url,omitempty"`
	ForcePathStyle      *bool   `json:"force_path_style,omitempty"`
}

// ProjectStorageSettingsPayload wraps storage_settings after update.
type ProjectStorageSettingsPayload struct {
	StorageSettings ProjectStorageSettings `json:"storage_settings"`
}

// SystemFile is metadata for a file in the system files REST API.
type SystemFile struct {
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

// SystemFilesListResponse is returned by listSystemFiles.
type SystemFilesListResponse struct {
	Files []SystemFile `json:"files"`
	Total int          `json:"total"`
}

// SystemFileUploadParams configures uploadSystemFile.
type SystemFileUploadParams struct {
	FileName string
	Content  []byte
	FileType string // optional; inferred from content type when empty
}

// DeleteSystemFilesResponse is returned by deleteSystemFiles.
type DeleteSystemFilesResponse struct {
	Success        bool     `json:"success"`
	DeletedIDs     []string `json:"deleted_ids"`
	StorageFailed  []string `json:"storage_failed,omitempty"`
	Message        string   `json:"message,omitempty"`
}
