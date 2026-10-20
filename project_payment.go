package goapitosdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ProjectTenantSubscription is an ApitoProjectPayment ledger row (REST JSON).
type ProjectTenantSubscription struct {
	ID                     string `json:"id,omitempty"`
	ProjectID              string `json:"project_id,omitempty"`
	TenantID               string `json:"tenant_id,omitempty"`
	PlanTier               string `json:"plan_tier,omitempty"`
	Provider               string `json:"provider,omitempty"`
	Status                 string `json:"status,omitempty"`
	BillingContactUserID   string `json:"billing_contact_user_id,omitempty"`
	ExternalSubscriptionID string `json:"external_subscription_id,omitempty"`
	ProductID              string `json:"product_id,omitempty"`
	CurrentPeriodStart     string `json:"current_period_start,omitempty"`
	CurrentPeriodEnd       string `json:"current_period_end,omitempty"`
	CancelAt               string `json:"cancel_at,omitempty"`
	CreatedBy              string `json:"created_by,omitempty"`
	Notes                  string `json:"notes,omitempty"`
	ManageURL              string `json:"manage_url,omitempty"`
	CreatedAt              string `json:"created_at,omitempty"`
	UpdatedAt              string `json:"updated_at,omitempty"`
}

// VerifyProjectPaymentParams configures POST .../providers/:provider/verify.
type VerifyProjectPaymentParams struct {
	ProjectID            string
	TenantID             string
	ProductID            string
	PurchaseToken        string
	PackageName          string
	Provider             string // default google_play
	BillingContactUserID string
}

// CancelTenantSubscriptionParams configures POST .../subscriptions/cancel.
type CancelTenantSubscriptionParams struct {
	ProjectID       string
	TenantID        string
	Mode            string // cancel | revoke (default cancel)
	SubscriptionID  string
	ForceLedgerOnly bool
}

// deriveSystemAPIBaseURL returns the engine origin (no /system/graphql or /secured/graphql).
func deriveSystemAPIBaseURL(graphqlURL string) string {
	u := strings.TrimSuffix(strings.TrimSpace(graphqlURL), "/")
	u = strings.TrimSuffix(u, "/system/graphql")
	u = strings.TrimSuffix(u, "/secured/graphql")
	return strings.TrimSuffix(u, "/")
}

func applyProjectPaymentAuth(req *http.Request, key string) {
	key = strings.TrimSpace(key)
	if key == "" {
		return
	}
	if strings.HasPrefix(key, "ak_") {
		req.Header.Set("X-Apito-Key", key)
		return
	}
	// apt_ access tokens and app-user JWTs use Bearer.
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("X-Use-Cookies", "false")
}

func (c *Client) requireProjectPaymentCredential() error {
	key := ""
	if c != nil {
		key = strings.TrimSpace(c.apiKey)
	}
	if key == "" || strings.HasPrefix(key, "ak_") {
		return fmt.Errorf("OPERATOR_CREDENTIAL_REQUIRED: projectPayment requires an app-user session or an apt_ operator token")
	}
	if isRetiredTokenPrefix(key) {
		return errTokenFormatRetired
	}
	return nil
}

func (c *Client) executeProjectPayment(ctx context.Context, method, path string, jsonBody interface{}, projectID, tenantID string) (*ProjectTenantSubscription, error) {
	raw, err := c.executeProjectPaymentRaw(ctx, method, path, jsonBody, projectID, tenantID)
	if err != nil {
		return nil, err
	}
	var out ProjectTenantSubscription
	if len(raw) == 0 {
		return &out, nil
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) executeProjectPaymentRaw(ctx context.Context, method, path string, jsonBody interface{}, projectID, tenantID string) ([]byte, error) {
	if err := c.requireProjectPaymentCredential(); err != nil {
		return nil, err
	}
	if isRetiredTokenPrefix(c.apiKey) {
		return nil, errTokenFormatRetired
	}
	base := deriveSystemAPIBaseURL(c.baseURL)
	u, err := url.Parse(strings.TrimSuffix(base, "/") + path)
	if err != nil {
		return nil, fmt.Errorf("invalid project-payment URL: %w", err)
	}
	var body io.Reader
	if jsonBody != nil {
		b, merr := json.Marshal(jsonBody)
		if merr != nil {
			return nil, merr
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	applyProjectPaymentAuth(req, c.apiKey)
	pid := strings.TrimSpace(projectID)
	if pid == "" {
		pid = c.projectID
	}
	if pid != "" {
		req.Header.Set(headerApitoProjectID, pid)
	}
	tid := strings.TrimSpace(tenantID)
	if tid == "" && ctx != nil {
		if v, ok := ctx.Value("tenant_id").(string); ok {
			tid = strings.TrimSpace(v)
		}
	}
	if tid != "" {
		req.Header.Set("X-Apito-Tenant-ID", tid)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("project payment request failed (%d): %s", res.StatusCode, string(raw))
	}
	return raw, nil
}

// VerifyProjectPayment verifies a store purchase and upserts the ledger row.
func (c *Client) VerifyProjectPayment(ctx context.Context, params VerifyProjectPaymentParams) (*ProjectTenantSubscription, error) {
	provider := strings.TrimSpace(params.Provider)
	if provider == "" {
		provider = "google_play"
	}
	projectID := strings.TrimSpace(params.ProjectID)
	if projectID == "" {
		return nil, fmt.Errorf("project_id required")
	}
	path := fmt.Sprintf("/system/v1/project-payment/projects/%s/providers/%s/verify", projectID, provider)
	body := map[string]interface{}{
		"tenant_id":      strings.TrimSpace(params.TenantID),
		"product_id":     strings.TrimSpace(params.ProductID),
		"purchase_token": strings.TrimSpace(params.PurchaseToken),
		"package_name":   strings.TrimSpace(params.PackageName),
	}
	if bc := strings.TrimSpace(params.BillingContactUserID); bc != "" {
		body["billing_contact_user_id"] = bc
	}
	return c.executeProjectPayment(ctx, http.MethodPost, path, body, projectID, params.TenantID)
}

// GetTenantSubscription returns the latest subscription for a tenant.
func (c *Client) GetTenantSubscription(ctx context.Context, projectID, tenantID string) (*ProjectTenantSubscription, error) {
	projectID = strings.TrimSpace(projectID)
	tenantID = strings.TrimSpace(tenantID)
	if projectID == "" || tenantID == "" {
		return nil, fmt.Errorf("project_id and tenant_id required")
	}
	path := fmt.Sprintf("/system/v1/project-payment/projects/%s/tenants/%s/subscription", projectID, tenantID)
	return c.executeProjectPayment(ctx, http.MethodGet, path, nil, projectID, tenantID)
}

// CancelTenantSubscription cancels via provider (Play first), then finalizes the ledger.
func (c *Client) CancelTenantSubscription(ctx context.Context, params CancelTenantSubscriptionParams) (*ProjectTenantSubscription, error) {
	projectID := strings.TrimSpace(params.ProjectID)
	tenantID := strings.TrimSpace(params.TenantID)
	if projectID == "" || tenantID == "" {
		return nil, fmt.Errorf("project_id and tenant_id required")
	}
	mode := strings.TrimSpace(params.Mode)
	if mode == "" {
		mode = "cancel"
	}
	path := fmt.Sprintf("/system/v1/project-payment/projects/%s/subscriptions/cancel", projectID)
	body := map[string]interface{}{
		"tenant_id": tenantID,
		"mode":      mode,
	}
	if sid := strings.TrimSpace(params.SubscriptionID); sid != "" {
		body["subscription_id"] = sid
	}
	if params.ForceLedgerOnly {
		body["force_ledger_only"] = true
	}
	return c.executeProjectPayment(ctx, http.MethodPost, path, body, projectID, tenantID)
}

// SubmitCustomPaymentParams is the signed-in Custom provider payload.
type SubmitCustomPaymentParams struct {
	ProjectID        string
	TenantID         string
	Email            string
	Plan             string
	Amount           float64
	Currency         string
	PeriodMonths     int
	PaymentReference string
}

// ProjectTenantInvoice is one row from the signed-in invoice list.
type ProjectTenantInvoice struct {
	ID             string  `json:"id,omitempty"`
	SubscriptionID string  `json:"subscription_id,omitempty"`
	TenantID       string  `json:"tenant_id,omitempty"`
	Amount         float64 `json:"amount,omitempty"`
	Currency       string  `json:"currency,omitempty"`
	Status         string  `json:"status,omitempty"`
	Provider       string  `json:"provider,omitempty"`
	IssuedAt       string  `json:"issued_at,omitempty"`
	PaidAt         string  `json:"paid_at,omitempty"`
	CreatedAt      string  `json:"created_at,omitempty"`
}

// SubmitCustomPayment posts a Custom payment for the signed-in user.
// A project that does not create subscribers returns a nil subscription.
func (c *Client) SubmitCustomPayment(ctx context.Context, params SubmitCustomPaymentParams) (*ProjectTenantSubscription, error) {
	projectID := strings.TrimSpace(params.ProjectID)
	tenantID := strings.TrimSpace(params.TenantID)
	if projectID == "" || tenantID == "" {
		return nil, fmt.Errorf("project_id and tenant_id required")
	}
	path := fmt.Sprintf("/system/v1/project-payment/projects/%s/providers/custom/submit", projectID)
	body := map[string]interface{}{
		"email":             strings.TrimSpace(params.Email),
		"plan":              strings.TrimSpace(params.Plan),
		"amount":            params.Amount,
		"currency":          strings.TrimSpace(params.Currency),
		"period_months":     params.PeriodMonths,
		"tenant_id":         tenantID,
		"payment_reference": strings.TrimSpace(params.PaymentReference),
	}
	raw, err := c.executeProjectPaymentRaw(ctx, http.MethodPost, path, body, projectID, tenantID)
	if err != nil {
		return nil, err
	}
	var probe map[string]interface{}
	if err := json.Unmarshal(raw, &probe); err != nil {
		return nil, err
	}
	id, _ := probe["id"].(string)
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}
	var out ProjectTenantSubscription
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListSignedInTenantInvoices lists invoices for the signed-in tenant.
func (c *Client) ListSignedInTenantInvoices(ctx context.Context, projectID, tenantID string) ([]ProjectTenantInvoice, error) {
	projectID = strings.TrimSpace(projectID)
	tenantID = strings.TrimSpace(tenantID)
	if projectID == "" || tenantID == "" {
		return nil, fmt.Errorf("project_id and tenant_id required")
	}
	path := fmt.Sprintf("/system/v1/project-payment/projects/%s/tenants/%s/invoices", projectID, tenantID)
	raw, err := c.executeProjectPaymentRaw(ctx, http.MethodGet, path, nil, projectID, tenantID)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Invoices []ProjectTenantInvoice `json:"invoices"`
	}
	if len(raw) == 0 {
		return nil, nil
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	return payload.Invoices, nil
}
