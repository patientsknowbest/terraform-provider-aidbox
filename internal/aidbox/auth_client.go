package aidbox

import (
	"context"
)

type AuthClient struct {
	ResourceBase
	Secret     string        `json:"secret"`
	GrantTypes []interface{} `json:"grant_types"`
}

func (*AuthClient) GetResourceName() string {
	return "Client"
}

func (client *Client) CreateAuthClient(ctx context.Context, basicClient *AuthClient, boxId string) (*AuthClient, error) {
	ap, err := client.createResource(ctx, basicClient, boxId)
	if err != nil {
		return nil, err
	}
	return ap.(*AuthClient), nil
}

func (client *Client) GetAuthClient(ctx context.Context, id, boxId string) (*AuthClient, error) {
	rr, err := client.getResource(ctx, "/Client/"+id, boxId)
	if err != nil {
		return nil, err
	}
	return rr.(*AuthClient), nil
}

func (client *Client) UpdateAuthClient(ctx context.Context, q *AuthClient, boxId string) (*AuthClient, error) {
	rr, err := client.updateResource(ctx, q, boxId)
	if err != nil {
		return nil, err
	}
	return rr.(*AuthClient), nil
}

func (client *Client) DeleteAuthClient(ctx context.Context, id, boxId string) error {
	return client.deleteResource(ctx, "/Client/"+id, boxId)
}
