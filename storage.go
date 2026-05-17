package goapitosdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const storageSettingsFragment = `
	use_free_cloud_storage
	endpoint
	region
	bucket
	access_key_id
	has_secret_access_key
	public_base_url
	force_path_style
`

// GetProjectStorageSettings reads project storage settings via getProject.
func (c *Client) GetProjectStorageSettings(ctx context.Context, projectID string) (*ProjectStorageSettings, error) {
	if strings.TrimSpace(projectID) == "" {
		return nil, fmt.Errorf("getProjectStorageSettings: project id is required")
	}
	query := `
		query GetProjectStorageSettings($_id: String!) {
			getProject(_id: $_id) {
				storage_settings {
					` + storageSettingsFragment + `
				}
			}
		}
	`
	response, err := c.executeGraphQL(ctx, query, map[string]interface{}{"_id": projectID})
	if err != nil {
		return nil, fmt.Errorf("getProjectStorageSettings: %w", err)
	}
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}
	proj, ok := data["getProject"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected getProject response")
	}
	raw, ok := proj["storage_settings"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected storage_settings response")
	}
	return mapToProjectStorageSettings(raw), nil
}

// UpdateProjectStorageSettings persists project storage settings.
func (c *Client) UpdateProjectStorageSettings(ctx context.Context, input UpdateProjectStorageInput) (*ProjectStorageSettings, error) {
	query := `
		mutation UpdateProjectStorageSettings($input: UpdateProjectStorageInput!) {
			updateProjectStorageSettings(input: $input) {
				storage_settings {
					` + storageSettingsFragment + `
				}
			}
		}
	`
	variables := map[string]interface{}{
		"input": buildUpdateProjectStorageInputMap(input),
	}
	response, err := c.executeGraphQL(ctx, query, variables)
	if err != nil {
		return nil, fmt.Errorf("updateProjectStorageSettings: %w", err)
	}
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}
	payload, ok := data["updateProjectStorageSettings"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected updateProjectStorageSettings response")
	}
	raw, ok := payload["storage_settings"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected storage_settings response")
	}
	settings := mapToProjectStorageSettings(raw)
	return settings, nil
}

func buildUpdateProjectStorageInputMap(input UpdateProjectStorageInput) map[string]interface{} {
	out := map[string]interface{}{}
	if input.UseFreeCloudStorage != nil {
		out["use_free_cloud_storage"] = *input.UseFreeCloudStorage
	}
	if input.Endpoint != nil {
		out["endpoint"] = *input.Endpoint
	}
	if input.Region != nil {
		out["region"] = *input.Region
	}
	if input.Bucket != nil {
		out["bucket"] = *input.Bucket
	}
	if input.AccessKeyID != nil {
		out["access_key_id"] = *input.AccessKeyID
	}
	if input.SecretAccessKey != nil {
		out["secret_access_key"] = *input.SecretAccessKey
	}
	if input.PublicBaseURL != nil {
		out["public_base_url"] = *input.PublicBaseURL
	}
	if input.ForcePathStyle != nil {
		out["force_path_style"] = *input.ForcePathStyle
	}
	return out
}

func mapToProjectStorageSettings(m map[string]interface{}) *ProjectStorageSettings {
	if m == nil {
		return nil
	}
	s := &ProjectStorageSettings{}
	if v, ok := m["use_free_cloud_storage"].(bool); ok {
		s.UseFreeCloudStorage = v
	}
	if v, ok := m["has_secret_access_key"].(bool); ok {
		s.HasSecretAccessKey = v
	}
	s.Endpoint = stringPtrFromMap(m, "endpoint")
	s.Region = stringPtrFromMap(m, "region")
	s.Bucket = stringPtrFromMap(m, "bucket")
	s.AccessKeyID = stringPtrFromMap(m, "access_key_id")
	s.PublicBaseURL = stringPtrFromMap(m, "public_base_url")
	if v, ok := m["force_path_style"].(bool); ok {
		s.ForcePathStyle = &v
	}
	return s
}

func stringPtrFromMap(m map[string]interface{}, key string) *string {
	if v, ok := m[key].(string); ok && strings.TrimSpace(v) != "" {
		s := v
		return &s
	}
	return nil
}

func mapToSystemFile(m map[string]interface{}) SystemFile {
	f := SystemFile{}
	if v, ok := m["id"].(string); ok {
		f.ID = v
	}
	if v, ok := m["file_type"].(string); ok {
		f.FileType = v
	}
	if v, ok := m["file_name"].(string); ok {
		f.FileName = v
	}
	if v, ok := m["file_extension"].(string); ok {
		f.FileExtension = v
	}
	if v, ok := m["content_type"].(string); ok {
		f.ContentType = v
	}
	if v, ok := m["size"].(float64); ok {
		f.Size = int64(v)
	}
	if v, ok := m["size"].(int64); ok {
		f.Size = v
	}
	if v, ok := m["url"].(string); ok {
		f.URL = v
	}
	if v, ok := m["created_by"].(string); ok {
		f.CreatedBy = v
	}
	if v, ok := m["created_at"].(string); ok {
		f.CreatedAt = v
	}
	return f
}

// UploadSystemFile uploads a file via POST /system/files/upload.
func (c *Client) UploadSystemFile(ctx context.Context, params SystemFileUploadParams) (*SystemFile, error) {
	if len(params.Content) == 0 {
		return nil, fmt.Errorf("uploadSystemFile: file content is required")
	}
	fileName := strings.TrimSpace(params.FileName)
	if fileName == "" {
		fileName = "upload"
	}

	body, status, err := c.executeREST(ctx, restRequest{
		method: http.MethodPost,
		path:   "/files/upload",
		multipartFn: func(w *multipart.Writer) error {
			part, err := w.CreateFormFile("file", fileName)
			if err != nil {
				return err
			}
			if _, err := io.Copy(part, bytes.NewReader(params.Content)); err != nil {
				return err
			}
			if strings.TrimSpace(params.FileType) != "" {
				if err := w.WriteField("file_type", strings.TrimSpace(params.FileType)); err != nil {
					return err
				}
			}
			return nil
		},
	})
	if err != nil {
		return nil, fmt.Errorf("uploadSystemFile: %w", err)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("uploadSystemFile: HTTP %d: %s", status, string(body))
	}
	envelope, err := parseRESTEnvelope(body)
	if err != nil {
		return nil, fmt.Errorf("uploadSystemFile: %w", err)
	}
	raw, ok := envelope["file"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("uploadSystemFile: unexpected response")
	}
	file := mapToSystemFile(raw)
	return &file, nil
}

// ListSystemFiles lists files via GET /system/files/list.
func (c *Client) ListSystemFiles(ctx context.Context, fileType string, limit, offset int) (*SystemFilesListResponse, error) {
	q := url.Values{}
	if strings.TrimSpace(fileType) != "" {
		q.Set("file_type", strings.TrimSpace(fileType))
	}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}

	body, status, err := c.executeREST(ctx, restRequest{
		method: http.MethodGet,
		path:   "/files/list",
		query:  q,
	})
	if err != nil {
		return nil, fmt.Errorf("listSystemFiles: %w", err)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("listSystemFiles: HTTP %d: %s", status, string(body))
	}
	envelope, err := parseRESTEnvelope(body)
	if err != nil {
		return nil, fmt.Errorf("listSystemFiles: %w", err)
	}
	resp := &SystemFilesListResponse{}
	if v, ok := envelope["total"].(float64); ok {
		resp.Total = int(v)
	}
	if v, ok := envelope["total"].(int); ok {
		resp.Total = v
	}
	if arr, ok := envelope["files"].([]interface{}); ok {
		for _, it := range arr {
			if m, ok := it.(map[string]interface{}); ok {
				resp.Files = append(resp.Files, mapToSystemFile(m))
			}
		}
	}
	return resp, nil
}

// DeleteSystemFiles deletes files via POST /system/files/delete.
func (c *Client) DeleteSystemFiles(ctx context.Context, ids []string) (*DeleteSystemFilesResponse, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("deleteSystemFiles: ids are required")
	}
	body, status, err := c.executeREST(ctx, restRequest{
		method:   http.MethodPost,
		path:     "/files/delete",
		jsonBody: map[string]interface{}{"ids": ids},
	})
	if err != nil {
		return nil, fmt.Errorf("deleteSystemFiles: %w", err)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("deleteSystemFiles: HTTP %d: %s", status, string(body))
	}
	var envelope map[string]interface{}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("deleteSystemFiles: %w", err)
	}
	out := &DeleteSystemFilesResponse{}
	if v, ok := envelope["success"].(bool); ok {
		out.Success = v
	}
	if arr, ok := envelope["deleted_ids"].([]interface{}); ok {
		for _, it := range arr {
			if s, ok := it.(string); ok {
				out.DeletedIDs = append(out.DeletedIDs, s)
			}
		}
	}
	if arr, ok := envelope["storage_failed"].([]interface{}); ok {
		for _, it := range arr {
			if s, ok := it.(string); ok {
				out.StorageFailed = append(out.StorageFailed, s)
			}
		}
	}
	if msg, ok := envelope["message"].(string); ok {
		out.Message = msg
	}
	if !out.Success && out.Message != "" {
		return out, fmt.Errorf("%s", out.Message)
	}
	return out, nil
}
