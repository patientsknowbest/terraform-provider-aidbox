package aidbox

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
)

// SearchParameter Aidbox customized representation of FHIR "SearchParameter"
type SearchParameter struct {
	ResourceBase
	Name               string              `json:"name"`
	Module             string              `json:"module,omitempty"`
	Type               SearchParameterType `json:"type"`
	ExpressionElements [][]interface{}     `json:"expression"`
	Resource           Reference           `json:"resource"`
}

func (*SearchParameter) GetResourcePath() string {
	return "SearchParameter"
}

type SearchParameterType int

const (
	SearchParameterTypeString SearchParameterType = iota
	SearchParameterTypeNumber
	SearchParameterTypeDate
	SearchParameterTypeToken
	SearchParameterTypeQuantity
	SearchParameterTypeReference
	SearchParameterTypeUri
	SearchParameterTypeComposite
)

func (t SearchParameterType) ToString() string {
	switch t {
	case SearchParameterTypeString:
		return "string"
	case SearchParameterTypeNumber:
		return "number"
	case SearchParameterTypeDate:
		return "date"
	case SearchParameterTypeToken:
		return "token"
	case SearchParameterTypeQuantity:
		return "quantity"
	case SearchParameterTypeReference:
		return "reference"
	case SearchParameterTypeUri:
		return "uri"
	case SearchParameterTypeComposite:
		return "composite"
	}
	log.Panicf("Unexpected SearchParameterType %d\n", t)
	return ""
}

const ErrInvalidSearchParameterType AidboxError = "Unsupported search parameter type"

func ParseSearchParameterType(typeString string) (SearchParameterType, error) {
	switch typeString {
	case "string":
		return SearchParameterTypeString, nil
	case "number":
		return SearchParameterTypeNumber, nil
	case "date":
		return SearchParameterTypeDate, nil
	case "token":
		return SearchParameterTypeToken, nil
	case "quantity":
		return SearchParameterTypeQuantity, nil
	case "reference":
		return SearchParameterTypeReference, nil
	case "uri":
		return SearchParameterTypeUri, nil
	case "composite":
		return SearchParameterTypeComposite, nil
	default:
		return 0, ErrInvalidSearchParameterType
	}
}

func (t SearchParameterType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(t.ToString())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (t *SearchParameterType) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	tt, err := ParseSearchParameterType(j)
	if err != nil {
		return err
	}
	*t = tt
	return nil
}

func (apiClient *ApiClient) CreateSearchParameter(ctx context.Context, searchParameter *SearchParameter) (*SearchParameter, error) {
	response := &SearchParameter{}
	return response, apiClient.createResource(ctx, searchParameter, response)
}

func (apiClient *ApiClient) GetSearchParameter(ctx context.Context, id string) (*SearchParameter, error) {
	response := &SearchParameter{}
	return response, apiClient.getResource(ctx, id, response)
}

func (apiClient *ApiClient) UpdateSearchParameter(ctx context.Context, q *SearchParameter) (*SearchParameter, error) {
	response := &SearchParameter{}
	return response, apiClient.updateResource(ctx, q, response)
}

func (apiClient *ApiClient) DeleteSearchParameter(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, id, &SearchParameter{})
}
