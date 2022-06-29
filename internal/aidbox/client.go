package aidbox

import (
	"context"
)

type Client struct {
	ResourceBase
	Secret     string        `json:"secret"`
	GrantTypes []interface{} `json:"grant_types"`
}

func (*Client) GetResourceName() string {
	return "Client"
}

func (apiClient *ApiClient) CreateClient(ctx context.Context, client *Client, boxId string) (*Client, error) {
	c, err := apiClient.createResource(ctx, client, boxId)
	if err != nil {
		return nil, err
	}
	return c.(*Client), nil
}

func (apiClient *ApiClient) GetClient(ctx context.Context, id, boxId string) (*Client, error) {
	rr, err := apiClient.getResource(ctx, "/Client/"+id, boxId)
	if err != nil {
		return nil, err
	}
	return rr.(*Client), nil
}

func (apiClient *ApiClient) UpdateClient(ctx context.Context, q *Client, boxId string) (*Client, error) {
	rr, err := apiClient.updateResource(ctx, q, boxId)
	if err != nil {
		return nil, err
	}
	return rr.(*Client), nil
}

func (apiClient *ApiClient) DeleteClient(ctx context.Context, id, boxId string) error {
	return apiClient.deleteResource(ctx, "/Client/"+id, boxId)
}
