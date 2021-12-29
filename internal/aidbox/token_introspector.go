package aidbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type TokenIntrospector struct {
	ResourceBase
	TokenIntrospectionEndpoint *TokenIntrospectionEndpoint `json:"introspection_endpoint,omitempty"`
	JWKSURI                    string                      `json:"jwks_uri,omitempty"`
	TokenIntrospectorJWT       *TokenIntrospectorJWT       `json:"jwt,omitempty"`
	Type                       TokenIntrospectorType       `json:"type"`
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
	tt, err := ParseTokenIntrospectorType(j)
	if err != nil {
		return err
	}
	*g = tt
	return nil
}

func (client *Client) CreateTokenIntrospector(ctx context.Context, introspector *TokenIntrospector) (*TokenIntrospector, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(introspector)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.URL+"/TokenIntrospector", &buf)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(client.Username, client.Password)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Unexpected status code received %d %s", res.StatusCode, res.Status)
	}
	result := &TokenIntrospector{}
	err = json.NewDecoder(res.Body).Decode(result)
	return result, err
}

func (client *Client) GetTokenIntrospector(ctx context.Context, id string) (*TokenIntrospector, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.URL+"/TokenIntrospector/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(client.Username, client.Password)
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code %d %s", res.StatusCode, res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	rr, err := parseResource(b)
	if err != nil {
		return nil, err
	}
	return rr.(*TokenIntrospector), nil
}

//func (client *Client) GetTokenIntrospectors(ctx context.Context) ([]*TokenIntrospector, error) {
//	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.URL+"/TokenIntrospector", nil)
//	if err != nil {
//		return nil, err
//	}
//	req.SetBasicAuth(client.Username, client.Password)
//	req.Header.Set("Accept", "application/json")
//	res, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	if res.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("Unexpected status code %d %s", res.StatusCode, res.Status)
//	}
//	b, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return nil, err
//	}
//	rr, err := parseSearchResponse(b)
//	if err != nil {
//		return nil, err
//	}
//	tis := make([]*TokenIntrospector, len(rr))
//	for i, r := range rr {
//		tis[i] = r.(*TokenIntrospector)
//	}
//	return tis, err
//}

func (client *Client) UpdateTokenIntrospector(ctx context.Context, q *TokenIntrospector) (*TokenIntrospector, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(q)
	if err != nil {
		return nil, err
	}
	log.Printf("[TRACE] sending [[ %s ]]", buf.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, client.URL+"/TokenIntrospector/"+q.ID, &buf)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(client.Username, client.Password)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if !isAlright(res.StatusCode) {
		return nil, fmt.Errorf("Unexpected status code %d %s", res.StatusCode, res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	rr, err := parseResource(b)
	if err != nil {
		return nil, err
	}
	return rr.(*TokenIntrospector), nil
}

func (client *Client) DeleteTokenIntrospector(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, client.URL+"/TokenIntrospector/"+id, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(client.Username, client.Password)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if !isAlright(res.StatusCode) {
		return fmt.Errorf("Unexpected status code %d %s", res.StatusCode, res.Status)
	}
	return nil
}
