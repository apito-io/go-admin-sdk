package goapitosdk

// Project files REST paths are relative to Config.RestBaseURL (e.g. http://host:5050/secured).
// Full URLs: POST/GET/POST …/secured/files/{upload|list|delete}.
// Metadata is stored in the project database files table (tenant-scoped on SaaS via X-Apito-Tenant-ID).
const (
	FilesUploadPath = "/files/upload"
	FilesListPath   = "/files/list"
	FilesDeletePath = "/files/delete"
)
