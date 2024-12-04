package aidbox

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
)

type TokenIntrospector struct {
	ResourceBase
	TokenIntrospectionEndpoint *TokenIntrospectionEndpoint `json:"introspection_endpoint,omitempty"`
	JWKSURI                    string                      `json:"jwks_uri,omitempty"`
	TokenIntrospectorJWT       *TokenIntrospectorJWT       `json:"jwt,omitempty"`
	Type                       TokenIntrospectorType       `json:"type"`
}

func (*TokenIntrospector) GetResourceName() string {
	return "TokenIntrospector"
}

type TokenIntrospectionEndpoint struct {
	Authorization string `json:"authorization"`
	URL           string `json:"url"`
}

type TokenIntrospectorJWT struct {
	ISS    string `json:"iss"`
	Secret string `json:"secret,omitempty"`
}

type TokenIntrospectorType int

const (
	TokenIntrospectorTypeJWT TokenIntrospectorType = iota
	TokenIntrospectorTypeOpaque
)

const ErrInvalidTokenIntrospectorType AidboxError = "Invalid token introspector type"

func ParseTokenIntrospectorType(typeString string) (TokenIntrospectorType, error) {
	switch typeString {
	case "opaque":
		return TokenIntrospectorTypeOpaque, nil
	case "jwt":
		return TokenIntrospectorTypeJWT, nil
	default:
		return 0, ErrInvalidTokenIntrospectorType
	}
}

func (g TokenIntrospectorType) ToString() string {
	switch g {
	case TokenIntrospectorTypeOpaque:
		return "opaque"
	case TokenIntrospectorTypeJWT:
		return "jwt"
	}
	// Expect the compiler to have spotted problems before now
	log.Panicf("Unexpected TokenIntrospectorType %d\n", g)
	return ""
}

func (g TokenIntrospectorType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(g.ToString())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (g *TokenIntrospectorType) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	tt, err := ParseTokenIntrospectorType(j)
	if err != nil {
		return err
	}
	*g = tt
	return nil
}

func (apiClient *ApiClient) CreateTokenIntrospector(ctx context.Context, introspector *TokenIntrospector) (*TokenIntrospector, error) {
	rr, err := apiClient.createResource(ctx, introspector)
	if err != nil {
		return nil, err
	}
	return rr.(*TokenIntrospector), nil
}

func (apiClient *ApiClient) GetTokenIntrospector(ctx context.Context, id string) (*TokenIntrospector, error) {
	rr, err := apiClient.getResource(ctx, "/TokenIntrospector/"+id)
	if err != nil {
		return nil, err
	}
	return rr.(*TokenIntrospector), nil
}

func (apiClient *ApiClient) UpdateTokenIntrospector(ctx context.Context, q *TokenIntrospector) (*TokenIntrospector, error) {
	rr, err := apiClient.updateResource(ctx, q)
	if err != nil {
		return nil, err
	}
	return rr.(*TokenIntrospector), nil
}

func (apiClient *ApiClient) DeleteTokenIntrospector(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, "/TokenIntrospector/"+id)
}
