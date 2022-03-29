package aidbox

import (
	"context"
)

// https://docs.aidbox.app/multibox/multibox-box-manager-api

type Box struct {
	ResourceBase
	Description string `json:"description"`
	FhirVersion string `json:"fhirVersion"`
	AccessToken string `json:"access-token,omitempty"`
	BoxURL      string `json:"box-url,omitempty"`
}

func (*Box) GetResourceName() string {
	return "Box"
}

func (client *Client) CreateBox(ctx context.Context, box *Box) (*Box, error) {
	resultBox := Box{}
	err := client.rpcRequest(ctx, "multibox/create-box", box, &resultBox, "")
	return &resultBox, err
}

func (client *Client) GetBox(ctx context.Context, id string) (*Box, error) {
	resultBox := Box{}
	err := client.rpcRequest(ctx, "multibox/get-box", &struct {
		Id string `json:"id"`
	}{id}, &resultBox, "")
	return &resultBox, err
}

func (client *Client) DeleteBox(ctx context.Context, id string) error {
	return client.rpcRequest(ctx, "multibox/delete-box", &struct {
		Id string `json:"id"`
	}{id}, &Box{}, "")
}
