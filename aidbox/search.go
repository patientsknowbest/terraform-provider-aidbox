package aidbox

import (
	"context"
)

// Aidbox Search resource: https://docs.aidbox.app/api/rest-api/aidbox-search#search-resource
type Search struct {
	ResourceBase
	Name     string    `json:"name"`
	Module   string    `json:"module,omitempty"`
	Resource Reference `json:"resource"`
	Where    string    `json:"where"`
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
