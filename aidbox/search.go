package aidbox

import (
	"context"
)

// Aidbox Search resource: https://docs.aidbox.app/api/rest-api/aidbox-search#search-resource
type Search struct {
	ResourceBase
	Name        string    `json:"name"`
	ParamParser string    `json:"param-parser,omitempty"`
	Module      string    `json:"module,omitempty"`
	Resource    Reference `json:"resource"`
	Where       string    `json:"where"`
	TokenSql    *TokenSql `json:"token-sql,omitempty"`
}

// TokenSql holds the SQL fragments Aidbox uses for param-parser: token
// search parameters. See https://docs.aidbox.app/api/rest-api/aidbox-search#token-search
type TokenSql struct {
	// SQL template when only code is provided
	OnlyCode string `json:"only-code,omitempty"`
	// SQL template when only system is provided
	OnlySystem string `json:"only-system,omitempty"`
	// SQL template when no system is provided
	NoSystem string `json:"no-system,omitempty"`
	// SQL template when both system and code are provided.
	Both string `json:"both,omitempty"`
	// SQL template for text search.
	Text string `json:"text,omitempty"`
	// Format for text search.
	TextFormat string `json:"text-format,omitempty"`
}

func (*Search) GetResourcePath() string {
	return "Search"
}

func (apiClient *ApiClient) CreateSearch(ctx context.Context, searchParameter *Search) (*Search, error) {
	response := &Search{}
	return response, apiClient.createResource(ctx, searchParameter, response)
}

func (apiClient *ApiClient) GetSearch(ctx context.Context, id string) (*Search, error) {
	response := &Search{}
	return response, apiClient.getResource(ctx, id, response)
}

func (apiClient *ApiClient) UpdateSearch(ctx context.Context, q *Search) (*Search, error) {
	response := &Search{}
	return response, apiClient.updateResource(ctx, q, response)
}

func (apiClient *ApiClient) DeleteSearch(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, id, &Search{})
}
