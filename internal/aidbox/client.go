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

func (client *Client) createResource(ctx context.Context, resource Resource) (Resource, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(resource)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.URL+"/"+resource.GetResourceName(), &buf)
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
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return parseResource(b)
}

func (client *Client) getResource(ctx context.Context, relativePath string) (Resource, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.URL+relativePath, nil)
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
	return parseResource(b)
}

func (client *Client) updateResource(ctx context.Context, resource Resource) (Resource, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(resource)
	if err != nil {
		return nil, err
	}
	log.Printf("[TRACE] sending [[ %s ]]", buf.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, client.URL+"/"+resource.GetResourceName()+"/"+resource.GetID(), &buf)
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
	return parseResource(b)
}

func (client *Client) deleteResource(ctx context.Context, relativePath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, client.URL+relativePath, nil)
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
