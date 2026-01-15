package aidbox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// StructureDefinition Represents a subset of the FHIR R4 spec "StructureDefinition", this is not a full support
// for the resource.
// This supports two ways of interacting with the StructureDefinition
// - the "regular" methods as seen in other resource implementations only parse and write a subset of the resource,
// this is for setting up custom profiles in the "FHIR-way"
// - the "byUrl" methods, these parse and write the entirety of the resource but provide no type checking, and used
// for overriding core FHIR spec StructureDefinitions, which is an aidbox specific feature
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

func (apiClient *ApiClient) GetStructureDefinitionByUrl(ctx context.Context, canonicalUrl string) (*map[string]interface{}, error) {
	response := &Bundle{}
	err := apiClient.get(ctx, fmt.Sprintf("/fhir/StructureDefinition?url=%s", url.QueryEscape(canonicalUrl)), response)
	if err != nil {
		return nil, err
	}
	if len(response.Entry) == 0 {
		return nil, fmt.Errorf("StructureDefinition with canonical url '%s' does not exist", canonicalUrl)
	}
	if len(response.Entry) > 1 {
		return nil, fmt.Errorf("found %d StructureDefinition entries for canonical url '%s', expected 1", len(response.Entry), canonicalUrl)
	}
	resource := response.Entry[0].Resource
	var structureDefinition = map[string]interface{}{}
	err = json.Unmarshal(resource, &structureDefinition)
	if err != nil {
		return nil, err
	}
	return &structureDefinition, nil
}

func (apiClient *ApiClient) UpdateStructureDefinitionByUrl(ctx context.Context, sd *map[string]interface{}, canonicalUrl string) (*map[string]interface{}, error) {
	response := &map[string]interface{}{}
	return response, apiClient.put(ctx, sd, fmt.Sprintf("/fhir/StructureDefinition?url=%s", url.QueryEscape(canonicalUrl)), response)
}
