package aidbox

import (
	"context"
)

// https://docs.aidbox.app/multibox/multibox-box-manager-api

type Box struct {
	ResourceBase
	Description string            `json:"description"`
	FhirVersion string            `json:"fhirVersion"`
	AccessToken string            `json:"access-token,omitempty"`
	BoxURL      string            `json:"box-url,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
}

func (*Box) GetResourceName() string {
	return "Box"
}

func (apiClient *ApiClient) CreateBox(ctx context.Context, box *Box) (*Box, error) {
	resultBox := &Box{}
	err := apiClient.rpcRequest(ctx, "multibox/create-box", box, resultBox, "")
	if err != nil {
		return nil, err
	}
	// Not all the data are returned from the create-box operation, so call GetBox again to get the full thing
	resultBox, err = apiClient.GetBox(ctx, resultBox.ID)
	if err != nil {
		return nil, err
	}
	return resultBox, err
}

func (apiClient *ApiClient) GetBox(ctx context.Context, id string) (*Box, error) {
	resultBox := Box{}
	err := apiClient.rpcRequest(ctx, "multibox/get-box", &struct {
		Id string `json:"id"`
	}{id}, &resultBox, "")
	return &resultBox, err
}

func (apiClient *ApiClient) DeleteBox(ctx context.Context, id string) error {
	return apiClient.rpcRequest(ctx, "multibox/delete-box", &struct {
		Id string `json:"id"`
	}{id}, &Box{}, "")
}
