package aidbox

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
)

type Client struct {
	ResourceBase
	Secret     string      `json:"secret"`
	GrantTypes []GrantType `json:"grant_types"`
}

func (*Client) GetResourcePath() string {
	return "Client"
}

type GrantType int

const (
	GrantTypeBasic GrantType = iota
	//GrantTypeAuthorizationCode
	//GrantTypeCode
	//GrantTypePassword
	//GrantTypeClientCredentials
	//GrantTypeImplicit
	//GrantTypeRefreshToken
)

func (g GrantType) ToString() string {
	switch g {
	case GrantTypeBasic:
		return "basic"
		//case GrantTypeAuthorizationCode:
		//	return "authorization_code"
		//case GrantTypeCode:
		//	return "code"
		//case GrantTypePassword:
		//	return "password"
		//case GrantTypeClientCredentials:
		//	return "client_credentials"
		//case GrantTypeImplicit:
		//	return "implicit"
		//case GrantTypeRefreshToken:
		//	return "refresh_token"
	}
	log.Panicf("Unexpected GrantType %d\n", g)
	return ""
}

const ErrInvalidGrantType AidboxError = "Unsupported grant type"

func ParseGrantType(typeString string) (GrantType, error) {
	switch typeString {
	case "basic":
		return GrantTypeBasic, nil
	//case "authorization_code":
	//	return GrantTypeAuthorizationCode, nil
	//case "code":
	//	return GrantTypeCode, nil
	//case "password":
	//	return GrantTypePassword, nil
	//case "client_credentials":
	//	return GrantTypeClientCredentials, nil
	//case "implicit":
	//	return GrantTypeImplicit, nil
	//case "refresh_token":
	//	return GrantTypeRefreshToken, nil
	default:
		return 0, ErrInvalidGrantType
	}
}

func (g GrantType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(g.ToString())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (g *GrantType) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	tt, err := ParseGrantType(j)
	if err != nil {
		return err
	}
	*g = tt
	return nil
}

func (apiClient *ApiClient) CreateClient(ctx context.Context, client *Client) (*Client, error) {
	response := &Client{}
	return response, apiClient.createResource(ctx, client, response)
}

func (apiClient *ApiClient) GetClient(ctx context.Context, id string) (*Client, error) {
	response := &Client{}
	return response, apiClient.getResource(ctx, id, response)
}

func (apiClient *ApiClient) UpdateClient(ctx context.Context, q *Client) (*Client, error) {
	response := &Client{}
	return response, apiClient.updateResource(ctx, q, response)
}

func (apiClient *ApiClient) DeleteClient(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, id, &Client{})
}
