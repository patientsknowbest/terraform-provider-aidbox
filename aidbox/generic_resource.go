package aidbox

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"path"
)

// GenericResource
// Provides an 'escape-hatch' to allow using the terraform provider to control arbitrary resource type
// Similar to kubernetes_manifest
type GenericResource struct {
	// this is the combined ResourceType/ID, otherwise terraform doesn't know how to destroy this thing.
	ID              string
	ResourceContent json.RawMessage
}

func (g *GenericResource) MarshalJSON() ([]byte, error) {
	return g.ResourceContent, nil
}

func (g *GenericResource) UnmarshalJSON(b []byte) error {
	id, err := parseId(b)
	if err != nil {
		return err
	}
	*g = GenericResource{
		id,
		b,
	}
	return nil
}

func parseId(b []byte) (string, error) {
	var h map[string]any
	err := json.Unmarshal(b, &h)
	if err != nil {
		return "", err
	}
	idPart, ok := h["id"].(string)
	if !ok {
		return "", errors.New("no 'id' field in JSON body")
	}
	resourceType, ok := h["resourceType"].(string)
	if !ok {
		return "", errors.New("no 'resourceType' field in JSON body")
	}
	return resourceType + "/" + idPart, nil
}

// These functions skip some of the high level abstraction, because GenericResource.ID is a combination of
// the resourcetype/ID, it doesn't fit.
func (apiClient *ApiClient) CreateGenericResource(ctx context.Context, genericResource *GenericResource) (*GenericResource, error) {
	responseTarget := &GenericResource{}
	id := genericResource.ID
	if genericResource.ID == "" {
		i, err := parseId(genericResource.ResourceContent)
		if err != nil {
			return nil, err
		}
		id = i
	}
	err := apiClient.put(ctx, genericResource, path.Join("/", id), responseTarget)
	if err != nil {
		return nil, err
	}
	return responseTarget, nil
}

func (apiClient *ApiClient) GetGenericResource(ctx context.Context, id string) (*GenericResource, error) {
	responseTarget := &GenericResource{}
	err := apiClient.get(ctx, path.Join("/", id), responseTarget)
	if err != nil {
		return nil, err
	}
	return responseTarget, nil
}

func (apiClient *ApiClient) UpdateGenericResource(ctx context.Context, q *GenericResource) (*GenericResource, error) {
	responseTarget := &GenericResource{}
	err := apiClient.put(ctx, q, path.Join("/", q.ID), responseTarget)
	if err != nil {
		return nil, err
	}
	return responseTarget, nil
}

func (apiClient *ApiClient) DeleteGenericResource(ctx context.Context, id string) error {
	return apiClient.send(ctx, struct{}{}, path.Join("/", id), &struct{}{}, http.MethodDelete)
}
