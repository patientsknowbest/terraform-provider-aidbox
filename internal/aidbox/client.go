package aidbox

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	URL      string
	Username string
	Password string
}

type AidboxError string

func (t AidboxError) Error() string {
	return string(t)
}

func NewClient(URL, username, password string) *Client {
	return &Client{URL: URL, Username: username, Password: password}
}

func isAlright(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusBadRequest
}

/// Horrid double unmarshal business to do discriminators on incoming json objects.
func parseResource(in []byte) (Resource, error) {
	s := struct {
		ResourceType string `json:"resourceType"`
	}{}
	err := json.Unmarshal(in, &s)
	if err != nil {
		return nil, err
	}
	var r Resource
	switch s.ResourceType {
	case "TokenIntrospector":
		r = &TokenIntrospector{}
	default:
		return nil, fmt.Errorf("Unsupported resource type %s", s.ResourceType)
	}
	err = json.Unmarshal(in, &r)
	return r, err
}

///// Horrid double unmarshal business to do discriminators on incoming json objects.
//func parseSearchResponse(in []byte) ([]Resource, error) {
//	lst := &struct {
//		ResourceType string `json:"resourceType"`
//		Entry        []struct {
//			Resource json.RawMessage `json:"resource"`
//		} `json:"entry"`
//	}{}
//	err := json.Unmarshal(in, &lst)
//	if err != nil {
//		return nil, err
//	}
//	// Check we're dealing with a Bundle here
//	if lst.ResourceType != "Bundle" {
//		return nil, fmt.Errorf("Unexpected resource type %s", lst.ResourceType)
//	}
//	vv := make([]Resource, len(lst.Entry))
//	for ix, v := range lst.Entry {
//		res, err := parseResource(v.Resource)
//		if err != nil {
//			return nil, err
//		}
//		vv[ix] = res
//	}
//	return vv, nil
//}
