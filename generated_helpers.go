package goapitosdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// TypedModelOps provides context-aware helpers for secured-endpoint CRUD (Go "hooks" layer).
type TypedModelOps struct {
	client *Client
	doer   *GraphQLDoer
}

// TypedOps returns typed operation helpers for the client.
func (c *Client) TypedOps() *TypedModelOps {
	return &TypedModelOps{client: c, doer: NewGraphQLDoer(c)}
}

// ExecuteRaw runs an arbitrary generated operation document with variables.
func (t *TypedModelOps) ExecuteRaw(
	ctx context.Context,
	document string,
	variables map[string]interface{},
) (map[string]interface{}, error) {
	raw, err := t.doer.Do(ctx, document, variables)
	if err != nil {
		return nil, err
	}
	var out struct {
		Data   map[string]interface{} `json:"data"`
		Errors []interface{}          `json:"errors"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if len(out.Errors) > 0 {
		return nil, fmt.Errorf("graphql errors: %v", out.Errors)
	}
	return out.Data, nil
}

// ListModel executes a generated list query for a model using the standard Apito list document.
func (t *TypedModelOps) ListModel(
	ctx context.Context,
	model string,
	fields []string,
	variables map[string]interface{},
) (map[string]interface{}, error) {
	doc := NewDocumentBuilder(model).BuildListQuery(fields)
	return t.ExecuteRaw(ctx, doc, variables)
}

// GetModel executes a generated get query for a model by id.
func (t *TypedModelOps) GetModel(
	ctx context.Context,
	model string,
	id string,
	fields []string,
) (map[string]interface{}, error) {
	doc := NewDocumentBuilder(model).BuildGetQuery(fields)
	return t.ExecuteRaw(ctx, doc, map[string]interface{}{"id": id})
}
