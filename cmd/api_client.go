package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// API client for Access Analyzer
type APIClient struct {
	BaseURL string
	Client  *http.Client
}

// SourceType represents a scanner/source type from the API
type SourceType struct {
	SourceTypeID     string `json:"sourceTypeId"`
	TypeName         string `json:"typeName"`
	DisplayName      string `json:"displayName"`
	Description      string `json:"description"`
	Version          string `json:"version"`
	ScannerImage     string `json:"scannerImage"`
	IsActive         bool   `json:"isActive"`
	IsBuiltIn        bool   `json:"isBuiltIn"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
	SupportedScans   []string `json:"supportedScanTypes,omitempty"`
	Icon             string `json:"icon,omitempty"`
}

// SourceTypeListResponse represents the API response for listing source types
type SourceTypeListResponse struct {
	Data       []SourceType       `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// PaginationMetadata represents pagination information
type PaginationMetadata struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// NewAPIClient creates a new API client
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetSourceTypes fetches all source types from the API
func (c *APIClient) GetSourceTypes() (*SourceTypeListResponse, error) {
	// Build URL with pagination
	u, err := url.Parse(c.BaseURL + "/source-types")
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	params := url.Values{}
	params.Set("page", "1")
	params.Set("pageSize", "100") // Get all scanners in one request
	u.RawQuery = params.Encode()

	// Make HTTP request
	resp, err := c.Client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result SourceTypeListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	return &result, nil
}

// TestConnection tests the connection to the API
func (c *APIClient) TestConnection() error {
	// Try to get source types as a health check
	resp, err := c.Client.Get(c.BaseURL + "/source-types?page=1&pageSize=1")
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Helper function to get API client with configured endpoint
func getAPIClient() (*APIClient, error) {
	endpoint, err := getAAEndpoint()
	if err != nil {
		return nil, err
	}

	if endpoint == "" {
		return nil, fmt.Errorf("no endpoint configured - use 'nwx aa config --endpoint=\"<url>\"'")
	}

	return NewAPIClient(endpoint), nil
}