package aidbox

import (
	"context"
)

// GcpServiceAccount is not an official FHIR resource.
// It is a proprietary custom resource used specifically by Aidbox to store
// Google Cloud Platform service account credentials, which can then be referenced
// by other resources like SDCConfig.
type GcpServiceAccount struct {
	ResourceBase
	ResourceType        string `json:"resourceType,omitempty"`
	ServiceAccountEmail string `json:"service-account-email,omitempty"`
	GcloudKey           string `json:"private-key,omitempty"`
}

func (*GcpServiceAccount) GetResourcePath() string {
	return "fhir/GcpServiceAccount"
}

func (apiClient *ApiClient) CreateGcpServiceAccount(ctx context.Context, gcpServiceAccount *GcpServiceAccount) (*GcpServiceAccount, error) {
	response := &GcpServiceAccount{}
	return response, apiClient.createResource(ctx, gcpServiceAccount, response)
}

func (apiClient *ApiClient) GetGcpServiceAccount(ctx context.Context, id string) (*GcpServiceAccount, error) {
	response := &GcpServiceAccount{}
	return response, apiClient.getResource(ctx, id, response)
}

func (apiClient *ApiClient) UpdateGcpServiceAccount(ctx context.Context, q *GcpServiceAccount) (*GcpServiceAccount, error) {
	response := &GcpServiceAccount{}
	return response, apiClient.updateResource(ctx, q, response)
}

func (apiClient *ApiClient) DeleteGcpServiceAccount(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, id, &GcpServiceAccount{})
}
