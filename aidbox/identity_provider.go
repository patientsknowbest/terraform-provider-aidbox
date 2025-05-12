package aidbox

import (
	"context"
	"encoding/json"
	"log"
)

type IdentityProviderClient struct {
	ID     string `json:"id,omitempty"`
	Secret string `json:"secret,omitempty"`
}

type IdentityProvider struct {
	ResourceBase
	Title             string                         `json:"title,omitempty"`
	System            string                         `json:"system,omitempty"`
	AuthorizeEndpoint string                         `json:"authorize_endpoint,omitempty"`
	TokenEndpoint     string                         `json:"token_endpoint,omitempty"`
	UserinfoSource    IdentityProviderUserinfoSource `json:"userinfo-source,omitempty"`
	UserinfoEndpoint  string                         `json:"userinfo_endpoint,omitempty"`
	Scopes            []string                       `json:"scopes,omitempty"`
	Client            *IdentityProviderClient        `json:"client,omitempty"`
}

func (g *IdentityProvider) GetResourcePath() string {
	return "IdentityProvider/" + g.ID
}

type IdentityProviderUserinfoSource int

const (
	UserinfoSourceIdToken IdentityProviderUserinfoSource = iota
	UserinfoSourceUserinfoEndpoint
)

func (g IdentityProviderUserinfoSource) ToString() string {
	switch g {
	case UserinfoSourceIdToken:
		return "id-token"
	case UserinfoSourceUserinfoEndpoint:
		return "userinfo-endpoint"
	}
	log.Panicf("Unexpected IdentityProviderUserinfoSource %d\n", g)
	return ""
}

const ErrInvalidUserinfoSource AidboxError = "Invalid userinfo-source"

func ParseUserinfoSource(s string) (IdentityProviderUserinfoSource, error) {
	switch s {
	case UserinfoSourceIdToken.ToString():
		return UserinfoSourceIdToken, nil
	case UserinfoSourceUserinfoEndpoint.ToString():
		return UserinfoSourceUserinfoEndpoint, nil
	default:
		return 0, ErrInvalidUserinfoSource
	}
}

func (g IdentityProviderUserinfoSource) MarshalJSON() ([]byte, error) {
	return []byte("\"" + g.ToString() + "\""), nil
}

func (g *IdentityProviderUserinfoSource) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	tt, err := ParseUserinfoSource(j)
	if err != nil {
		return err
	}
	*g = tt
	return nil
}

func (apiClient *ApiClient) CreateIdentityProvider(ctx context.Context, identityProvider *IdentityProvider) (*IdentityProvider, error) {
	response := &IdentityProvider{}
	return response, apiClient.createResource(ctx, identityProvider, response)
}

func (apiClient *ApiClient) GetIdentityProvider(ctx context.Context, id string) (*IdentityProvider, error) {
	response := &IdentityProvider{}
	return response, apiClient.getResource(ctx, id, response)
}

func (apiClient *ApiClient) UpdateIdentityProvider(ctx context.Context, q *IdentityProvider) (*IdentityProvider, error) {
	response := &IdentityProvider{}
	return response, apiClient.updateResource(ctx, q, response)
}

func (apiClient *ApiClient) DeleteIdentityProvider(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, id, &IdentityProvider{})
}
