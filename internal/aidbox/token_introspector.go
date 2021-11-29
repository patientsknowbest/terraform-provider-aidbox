package aidbox

import (
	"bytes"
	"encoding/json"
)

type TokenIntrospectorError string

func (t TokenIntrospectorError) Error() string {
	return string(t)
}

const ErrInvalidTokenIntrospectorType TokenIntrospectorError = "Invalid token introspector type"

type TokenIntrospectorType int

const (
	TokenIntrospectorTypeJWT TokenIntrospectorType = iota
	TokenIntrospectorTypeOpaque
)

func (g TokenIntrospectorType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	switch g {
	case TokenIntrospectorTypeOpaque:
		buffer.WriteString("opaque")
	case TokenIntrospectorTypeJWT:
		buffer.WriteString("jwt")
	default:
		return nil, ErrInvalidTokenIntrospectorType
	}
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (g *TokenIntrospectorType) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	switch j {
	case "opaque":
		*g = TokenIntrospectorTypeOpaque
	case "jwt":
		*g = TokenIntrospectorTypeJWT
	default:
		return ErrInvalidTokenIntrospectorType
	}
	return nil
}

type TokenIntrospector struct {
	ResourceBase
	TokenIntrospectionEndpoint *TokenIntrospectionEndpoint `json:"introspection_endpoint,omitempty"`
	JWKSURI                    string `json:"jwks_uri,omitempty"`
	TokenIntrospectorJWT       *TokenIntrospectorJWT `json:"jwt,omitempty"`
	Type                       TokenIntrospectorType `json:"type"`
}

type TokenIntrospectionEndpoint struct {
	Authorization string `json:"authorization"`
	URL           string `json:"url"`
}

type TokenIntrospectorJWT struct {
	ISS    string `json:"iss"`
	Secret string `json:"secret,omitempty"`
}
