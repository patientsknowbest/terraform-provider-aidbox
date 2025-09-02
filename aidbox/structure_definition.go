package aidbox

import (
	"context"
	"encoding/json"
)

// StructureDefinition Represents a subset of the FHIR R4 spec "StructureDefinition", this is not a full support
// for the resource. Used only to allow testing and temporarily support defining profiles
type StructureDefinition struct {
	ResourceBase
	ResourceType   string `json:"resourceType,omitempty"`
	Name           string `json:"name"`
	Url            string `json:"url"`
	BaseDefinition string `json:"baseDefinition"`
	Derivation     string `json:"derivation"`
	Abstract       bool   `json:"abstract"`
	Type           string `json:"type"`
	Status         string `json:"status"`
	Kind           string `json:"kind"`
	Version        string `json:"version"`
	// Deliberately not doing more validation than "is json?" or adding a custom type as it's unnecessarily complex
	// to handle Element given the intention here is temporary and partial support. This leaves more chance for user
	// error, however we will print the details about any issues related to this property received from the server, so
	// the user can correct it anyway
	Differential *json.RawMessage `json:"differential"`
}

func (*StructureDefinition) GetResourcePath() string {
	return "fhir/StructureDefinition"
}

func (apiClient *ApiClient) CreateStructureDefinition(ctx context.Context, structureDefinition *StructureDefinition) (*StructureDefinition, error) {
	response := &StructureDefinition{}
	return response, apiClient.createResource(ctx, structureDefinition, response)
}

func (apiClient *ApiClient) GetStructureDefinition(ctx context.Context, id string) (*StructureDefinition, error) {
	response := &StructureDefinition{}
	return response, apiClient.getResource(ctx, id, response)
}

func (apiClient *ApiClient) UpdateStructureDefinition(ctx context.Context, q *StructureDefinition) (*StructureDefinition, error) {
	response := &StructureDefinition{}
	return response, apiClient.updateResource(ctx, q, response)
}

func (apiClient *ApiClient) DeleteStructureDefinition(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, id, &StructureDefinition{})
}
