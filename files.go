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

func mapToFile(m map[string]interface{}) File {
	f := File{}
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

// UploadFile uploads a project file via POST /system/files/upload.
// Metadata is persisted in the project DB; blobs use project-scoped object storage.
func (c *Client) UploadFile(ctx context.Context, params UploadFileParams) (*File, error) {
	if len(params.Content) == 0 {
		return nil, fmt.Errorf("uploadFile: file content is required")
	}
	fileName := strings.TrimSpace(params.FileName)
	if fileName == "" {
		fileName = "upload"
	}

	body, status, err := c.executeREST(ctx, restRequest{
		method: http.MethodPost,
		path:   FilesUploadPath,
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
		return nil, fmt.Errorf("uploadFile: %w", err)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("uploadFile: HTTP %d: %s", status, string(body))
	}
	envelope, err := parseRESTEnvelope(body)
	if err != nil {
		return nil, fmt.Errorf("uploadFile: %w", err)
	}
	raw, ok := envelope["file"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("uploadFile: unexpected response")
	}
	file := mapToFile(raw)
	return &file, nil
}

// ListFiles lists project files via GET /system/files/list.
func (c *Client) ListFiles(ctx context.Context, fileType string, limit, offset int) (*FilesListResponse, error) {
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
		path:   FilesListPath,
		query:  q,
	})
	if err != nil {
		return nil, fmt.Errorf("listFiles: %w", err)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("listFiles: HTTP %d: %s", status, string(body))
	}
	envelope, err := parseRESTEnvelope(body)
	if err != nil {
		return nil, fmt.Errorf("listFiles: %w", err)
	}
	resp := &FilesListResponse{}
	if v, ok := envelope["total"].(float64); ok {
		resp.Total = int(v)
	}
	if v, ok := envelope["total"].(int); ok {
		resp.Total = v
	}
	if arr, ok := envelope["files"].([]interface{}); ok {
		for _, it := range arr {
			if m, ok := it.(map[string]interface{}); ok {
				resp.Files = append(resp.Files, mapToFile(m))
			}
		}
	}
	return resp, nil
}

// DeleteFiles deletes project files via POST /system/files/delete.
func (c *Client) DeleteFiles(ctx context.Context, ids []string) (*DeleteFilesResponse, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("deleteFiles: ids are required")
	}
	body, status, err := c.executeREST(ctx, restRequest{
		method:   http.MethodPost,
		path:     FilesDeletePath,
		jsonBody: map[string]interface{}{"ids": ids},
	})
	if err != nil {
		return nil, fmt.Errorf("deleteFiles: %w", err)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("deleteFiles: HTTP %d: %s", status, string(body))
	}
	var envelope map[string]interface{}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("deleteFiles: %w", err)
	}
	out := &DeleteFilesResponse{}
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
