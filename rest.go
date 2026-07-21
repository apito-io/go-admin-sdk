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
	"strings"
)

type restRequest struct {
	method      string
	path        string
	query       url.Values
	jsonBody    interface{}
	multipartFn func(w *multipart.Writer) error
}

func (c *Client) setAuthHeaders(req *http.Request, ctx context.Context) {
	req.Header.Set("Content-Type", "application/json")
	applyAuthCredential(req, c.apiKey)
	if c.projectID != "" {
		req.Header.Set(headerApitoProjectID, c.projectID)
	}
	if ctx != nil && ctx.Value("tenant_id") != nil {
		if tenantID, ok := ctx.Value("tenant_id").(string); ok && strings.TrimSpace(tenantID) != "" {
			req.Header.Set("X-Apito-Tenant-ID", tenantID)
		}
	}
}

func (c *Client) executeREST(ctx context.Context, rr restRequest) ([]byte, int, error) {
	if strings.TrimSpace(c.restBaseURL) == "" {
		return nil, 0, fmt.Errorf("rest base URL is not configured")
	}
	if isRetiredTokenPrefix(c.apiKey) {
		return nil, 0, errTokenFormatRetired
	}

	u, err := url.Parse(strings.TrimSuffix(c.restBaseURL, "/") + rr.path)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid REST URL: %w", err)
	}
	if len(rr.query) > 0 {
		u.RawQuery = rr.query.Encode()
	}

	var body io.Reader
	var contentType string

	if rr.multipartFn != nil {
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)
		if err := rr.multipartFn(w); err != nil {
			return nil, 0, err
		}
		if err := w.Close(); err != nil {
			return nil, 0, err
		}
		body = &buf
		contentType = w.FormDataContentType()
	} else if rr.jsonBody != nil {
		jsonData, err := json.Marshal(rr.jsonBody)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal JSON body: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
		contentType = "application/json"
	}

	req, err := http.NewRequestWithContext(ctx, rr.method, u.String(), body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	if rr.multipartFn != nil {
		req.Header.Set("Content-Type", contentType)
		applyAuthCredential(req, c.apiKey)
		if c.projectID != "" {
			req.Header.Set(headerApitoProjectID, c.projectID)
		}
		if ctx != nil && ctx.Value("tenant_id") != nil {
			if tenantID, ok := ctx.Value("tenant_id").(string); ok && strings.TrimSpace(tenantID) != "" {
				req.Header.Set("X-Apito-Tenant-ID", tenantID)
			}
		}
	} else {
		c.setAuthHeaders(req, ctx)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}
	return respBody, resp.StatusCode, nil
}

func parseRESTEnvelope(body []byte) (map[string]interface{}, error) {
	var out map[string]interface{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("failed to parse REST response: %w", err)
	}
	if success, ok := out["success"].(bool); ok && !success {
		msg, _ := out["message"].(string)
		if msg == "" {
			msg = "request failed"
		}
		return nil, fmt.Errorf("%s", msg)
	}
	return out, nil
}
