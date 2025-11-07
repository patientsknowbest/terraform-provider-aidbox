package aidbox

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
)

// SearchParameterV2 Represents the FHIR R4 spec "SearchParameter"
type SearchParameterV2 struct {
	ResourceBase
	ResourceType string                `json:"resourceType,omitempty"`
	Name         string                `json:"name"`
	Type         SearchParameterTypeV2 `json:"type"`
	Expression   string                `json:"expression"`
	Description  string                `json:"description"`
	Url          string                `json:"url"`
	Status       string                `json:"status"`
	Code         string                `json:"code"`
	Base         []string              `json:"base"`
}

func (*SearchParameterV2) GetResourcePath() string {
	return "fhir/SearchParameter"
}

type SearchParameterTypeV2 int

const (
	SearchParameterTypeV2String SearchParameterTypeV2 = iota
	SearchParameterTypeV2Number
	SearchParameterTypeV2Date
	SearchParameterTypeV2Token
	SearchParameterTypeV2Quantity
	SearchParameterTypeV2Reference
	SearchParameterTypeV2Uri
	SearchParameterTypeV2Composite
)

func (t SearchParameterTypeV2) ToString() string {
	switch t {
	case SearchParameterTypeV2String:
		return "string"
	case SearchParameterTypeV2Number:
		return "number"
	case SearchParameterTypeV2Date:
		return "date"
	case SearchParameterTypeV2Token:
		return "token"
	case SearchParameterTypeV2Quantity:
		return "quantity"
	case SearchParameterTypeV2Reference:
		return "reference"
	case SearchParameterTypeV2Uri:
		return "uri"
	case SearchParameterTypeV2Composite:
		return "composite"
	}
	log.Panicf("Unexpected SearchParameterTypeV2 %d\n", t)
	return ""
}

const ErrInvalidSearchParameterTypeV2 AidboxError = "Unsupported search parameter type"

func ParseSearchParameterTypeV2(typeString string) (SearchParameterTypeV2, error) {
	switch typeString {
	case "string":
		return SearchParameterTypeV2String, nil
	case "number":
		return SearchParameterTypeV2Number, nil
	case "date":
		return SearchParameterTypeV2Date, nil
	case "token":
		return SearchParameterTypeV2Token, nil
	case "quantity":
		return SearchParameterTypeV2Quantity, nil
	case "reference":
		return SearchParameterTypeV2Reference, nil
	case "uri":
		return SearchParameterTypeV2Uri, nil
	case "composite":
		return SearchParameterTypeV2Composite, nil
	default:
		return 0, ErrInvalidSearchParameterTypeV2
	}
}

func (t SearchParameterTypeV2) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(t.ToString())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (t *SearchParameterTypeV2) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	tt, err := ParseSearchParameterTypeV2(j)
	if err != nil {
		return err
	}
	*t = tt
	return nil
}

func (apiClient *ApiClient) CreateSearchParameterV2(ctx context.Context, searchParameter *SearchParameterV2) (*SearchParameterV2, error) {
	response := &SearchParameterV2{}
	// use PUT to create instead of POST. Aidbox ignores the id in the request body otherwise for this resource
	return response, apiClient.updateResource(ctx, searchParameter, response)
}

func (apiClient *ApiClient) GetSearchParameterV2(ctx context.Context, id string) (*SearchParameterV2, error) {
	response := &SearchParameterV2{}
	return response, apiClient.getResource(ctx, id, response)
}

func (apiClient *ApiClient) UpdateSearchParameterV2(ctx context.Context, q *SearchParameterV2) (*SearchParameterV2, error) {
	response := &SearchParameterV2{}
	return response, apiClient.updateResource(ctx, q, response)
}

func (apiClient *ApiClient) DeleteSearchParameterV2(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, id, &SearchParameterV2{})
}
