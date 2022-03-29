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
	Link        string             `json:"link,omitempty"`
}

func (*AccessPolicy) GetResourceName() string {
	return "AccessPolicy"
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

func (client *Client) CreateAccessPolicy(ctx context.Context, accessPolicy *AccessPolicy, boxId string) (*AccessPolicy, error) {
	ap, err := client.createResource(ctx, accessPolicy, boxId)
	if err != nil {
		return nil, err
	}
	return ap.(*AccessPolicy), nil
}

func (client *Client) GetAccessPolicy(ctx context.Context, id, boxId string) (*AccessPolicy, error) {
	rr, err := client.getResource(ctx, "/AccessPolicy/"+id, boxId)
	if err != nil {
		return nil, err
	}
	return rr.(*AccessPolicy), nil
}

func (client *Client) UpdateAccessPolicy(ctx context.Context, q *AccessPolicy, boxId string) (*AccessPolicy, error) {
	rr, err := client.updateResource(ctx, q, boxId)
	if err != nil {
		return nil, err
	}
	return rr.(*AccessPolicy), nil
}

func (client *Client) DeleteAccessPolicy(ctx context.Context, id, boxId string) error {
	return client.deleteResource(ctx, "/AccessPolicy/"+id, boxId)
}
