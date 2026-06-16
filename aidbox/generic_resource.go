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
	ResourceTypeAndId string
	ResourceContent   json.RawMessage
}

func (g *GenericResource) MarshalJSON() ([]byte, error) {
	return g.ResourceContent, nil
}

func (g *GenericResource) UnmarshalJSON(b []byte) error {
	resourceTypeAndId, err := GetResourceTypeAndId(b)
	if err != nil {
		return err
	}
	*g = GenericResource{
		resourceTypeAndId,
		b,
	}
	return nil
}

func GetResourceTypeAndId(b []byte) (string, error) {
	h, err := parseToMap(b)
	if err != nil {
		return "", err
	}
	resourceType, ok := h["resourceType"]
	if !ok {
		return "", errors.New("no 'resourceType' field in JSON body")
	}
	idPart, ok := h["id"]
	if !ok {
		return "", errors.New("no 'id' field in JSON body")
	}
	return resourceType.(string) + "/" + idPart.(string), nil
}

func getResourceType(b []byte) (string, error) {
	h, err := parseToMap(b)
	if err != nil {
		return "", err
	}
	resourceType, ok := h["resourceType"]
	if !ok {
		return "", errors.New("no 'resourceType' field in JSON body")
	}
	return resourceType.(string), nil
}

func parseToMap(b []byte) (map[string]any, error) {
	var h map[string]any
	err := json.Unmarshal(b, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// These functions skip some of the high level abstraction, because GenericResource.ResourceTypeAndId is a combination of
// the resourcetype/ID, it doesn't fit.
func (apiClient *ApiClient) CreateGenericResource(ctx context.Context, genericResource *GenericResource) (*GenericResource, error) {
	responseTarget := &GenericResource{}
	resourceType, err := getResourceType(genericResource.ResourceContent)
	if err != nil {
		return nil, err
	}
	resourceTypeAndId := genericResource.ResourceTypeAndId
	if genericResource.ResourceTypeAndId == "" {
		i, err := GetResourceTypeAndId(genericResource.ResourceContent)
		if err == nil {
			resourceTypeAndId = i
		} else {
			// couldn't parse an ID from the resource; don't panic, the user just didn't specify one, we'll use POST
		}
	}
	if resourceTypeAndId != "" {
		err = apiClient.put(ctx, genericResource, path.Join("/", resourceTypeAndId), responseTarget)
	} else {
		err = apiClient.post(ctx, genericResource, path.Join("/", resourceType), responseTarget)
	}

	if err != nil {
		return nil, err
	}
	return responseTarget, nil
}

func (apiClient *ApiClient) GetGenericResource(ctx context.Context, resourceTypeAndId string) (*GenericResource, error) {
	responseTarget := &GenericResource{}
	err := apiClient.get(ctx, path.Join("/", resourceTypeAndId), responseTarget)
	if err != nil {
		return nil, err
	}
	return responseTarget, nil
}

func (apiClient *ApiClient) UpdateGenericResource(ctx context.Context, q *GenericResource) (*GenericResource, error) {
	responseTarget := &GenericResource{}
	err := apiClient.put(ctx, q, path.Join("/", q.ResourceTypeAndId), responseTarget)
	if err != nil {
		return nil, err
	}
	return responseTarget, nil
}

func (apiClient *ApiClient) DeleteGenericResource(ctx context.Context, resourceTypeAndId string) error {
	return apiClient.send(ctx, struct{}{}, path.Join("/", resourceTypeAndId), &struct{}{}, http.MethodDelete)
}
