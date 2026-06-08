package goapitosdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GraphQLDoer adapts Client for generated genqlient operations.
type GraphQLDoer struct {
	Client *Client
}

// Do executes a GraphQL request.
func (d *GraphQLDoer) Do(ctx context.Context, query string, variables map[string]interface{}) ([]byte, error) {
	if d.Client == nil {
		return nil, fmt.Errorf("graphql doer: client is nil")
	}
	payload := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.Client.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	d.Client.setAuthHeaders(req, ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := d.Client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("graphql HTTP %d: %s", resp.StatusCode, buf.String())
	}
	return buf.Bytes(), nil
}

// NewGraphQLDoer returns a doer wrapping the admin client.
func NewGraphQLDoer(c *Client) *GraphQLDoer {
	return &GraphQLDoer{Client: c}
}
