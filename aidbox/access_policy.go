package aidbox

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
)

type AccessPolicy struct {
	ResourceBase
	Description string             `json:"description,omitempty"`
	Engine      AccessPolicyEngine `json:"engine"`
	Schema      json.RawMessage    `json:"schema,omitempty"`
	Link        []Reference        `json:"link,omitempty"`
}

func (*AccessPolicy) GetResourcePath() string {
	return "AccessPolicy"
}

type Reference struct {
	ResourceId   string `json:"id"`
	ResourceType string `json:"resourceType"`
}

type AccessPolicyEngine int

const (
	AccessPolicyEngineJsonSchema AccessPolicyEngine = iota
	AccessPolicyEngineAllow
	//AccessPolicyEngineSql
	//AccessPolicyEngineComplex
	//AccessPolicyEngineMatcho
	//AccessPolicyEngineClj
)

func (g AccessPolicyEngine) ToString() string {
	switch g {
	case AccessPolicyEngineJsonSchema:
		return "json-schema"
	case AccessPolicyEngineAllow:
		return "allow"
		//case AccessPolicyEngineSql:
		//	return "sql"
		//case AccessPolicyEngineComplex:
		//	return "complex"
		//case AccessPolicyEngineMatcho:
		//	return "matcho"
		//case AccessPolicyEngineClj:
		//	return "clj"
	}
	log.Panicf("Unexpected AccessPolicyEngine %d\n", g)
	return ""
}

const ErrInvalidAccessPolicyEngine AidboxError = "Invalid access policy engine type"

func ParseAccessPolicyEngine(s string) (AccessPolicyEngine, error) {
	switch s {
	case "json-schema":
		return AccessPolicyEngineJsonSchema, nil
	case "allow":
		return AccessPolicyEngineAllow, nil
	//case "sql":
	//	return AccessPolicyEngineSql, nil
	//case "complex":
	//	return AccessPolicyEngineComplex, nil
	//case "matcho":
	//	return AccessPolicyEngineMatcho, nil
	//case "clj":
	//	return AccessPolicyEngineClj, nil
	default:
		return 0, ErrInvalidAccessPolicyEngine
	}
}

func (g AccessPolicyEngine) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(g.ToString())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (g *AccessPolicyEngine) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	tt, err := ParseAccessPolicyEngine(j)
	if err != nil {
		return err
	}
	*g = tt
	return nil
}

func (apiClient *ApiClient) CreateAccessPolicy(ctx context.Context, accessPolicy *AccessPolicy) (*AccessPolicy, error) {
	response := &AccessPolicy{}
	return response, apiClient.createResource(ctx, accessPolicy, response)
}

func (apiClient *ApiClient) GetAccessPolicy(ctx context.Context, id string) (*AccessPolicy, error) {
	response := &AccessPolicy{}
	return response, apiClient.getResource(ctx, id, response)
}

func (apiClient *ApiClient) UpdateAccessPolicy(ctx context.Context, q *AccessPolicy) (*AccessPolicy, error) {
	response := &AccessPolicy{}
	return response, apiClient.updateResource(ctx, q, response)
}

func (apiClient *ApiClient) DeleteAccessPolicy(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, id, &AccessPolicy{})
}
