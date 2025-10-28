package aidbox

import (
	"context"
)

// SDCConfig is not an official FHIR resource.
// It is a proprietary custom resource used specifically by Aidbox to configure its Structured Data Capture (SDC) module
// We need to implement it so we can support file storage from Aidbox in GCP storages
type SDCConfig struct {
	ResourceBase
	ResourceType string      `json:"resourceType,omitempty"`
	Name         string      `json:"name,omitempty"`
	Description  string      `json:"description,omitempty"`
	Default      bool        `json:"default,omitempty"`
	Storage      *SDCStorage `json:"storage,omitempty"`
}

type SDCStorage struct {
	Bucket  string      `json:"bucket,omitempty"`
	Account *SDCAccount `json:"account,omitempty"`
}

type SDCAccount struct {
	Reference string `json:"reference,omitempty"`
}

func (*SDCConfig) GetResourcePath() string {
	return "fhir/SDCConfig"
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
