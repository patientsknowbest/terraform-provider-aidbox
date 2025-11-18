package aidbox

import (
	"context"
	"encoding/json"
)

// SDCConfig is not an official FHIR resource.
// It is a proprietary custom resource used specifically by Aidbox to configure its Structured Data Capture (SDC) module
type SDCConfig struct {
	ResourceBase
	ResourceType string `json:"resourceType,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Default      bool   `json:"default,omitempty"`
	// Deliberately not doing more validation than "is json?" or adding a custom type as it's unnecessarily complex
	// to handle Element given the intention here is temporary and partial support. This leaves more chance for user
	// error, however we will print the details about any issues related to this property received from the server, so
	// the user can correct it anyway
	Storage *json.RawMessage `json:"storage"`
}

func (*SDCConfig) GetResourcePath() string {
	return "/SDCConfig"
}

func (apiClient *ApiClient) CreateSDCConfig(ctx context.Context, sDCConfig *SDCConfig) (*SDCConfig, error) {
	response := &SDCConfig{}
	return response, apiClient.createResource(ctx, sDCConfig, response)
}

func (apiClient *ApiClient) GetSDCConfig(ctx context.Context, id string) (*SDCConfig, error) {
	response := &SDCConfig{}
	return response, apiClient.getResource(ctx, id, response)
}

func (apiClient *ApiClient) UpdateSDCConfig(ctx context.Context, q *SDCConfig) (*SDCConfig, error) {
	response := &SDCConfig{}
	return response, apiClient.updateResource(ctx, q, response)
}

func (apiClient *ApiClient) DeleteSDCConfig(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, id, &SDCConfig{})
}
