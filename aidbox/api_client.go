package aidbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type ApiClient struct {
	URL      string
	Username string
	Password string
}

type AidboxError string

const NotFoundError AidboxError = "Not found"

func (t AidboxError) Error() string {
	return string(t)
}

func NewApiClient(URL, username, password string) *ApiClient {
	return &ApiClient{
		URL:        URL,
		Username:   username,
		Password:   password,
	}
}

func isAlright(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusBadRequest
}

// TODO (AS) why the resource kind is not passed here instead?
// TODO (AS) find a place to put debug response body. e.g. wrap(err) and wrap adds the body onto the error or whatever so tf prints it out
// Horrid double unmarshal business to do discriminators on incoming json objects.
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
	case "AccessPolicy":
		r = &AccessPolicy{}
	case "Client":
		r = &Client{}
	case "SearchParameter":
		r = &SearchParameter{}
	case "AidboxSubscriptionTopic":
		r = &AidboxSubscriptionTopic{}
	case "AidboxTopicDestination":
		r = &AidboxTopicDestination{}
	default:
		return nil, fmt.Errorf("unsupported resource type %s, please add implementation", s.ResourceType)
	}
	err = json.Unmarshal(in, &r)
	return r, err
}

func (apiClient *ApiClient) createResource(ctx context.Context, resource Resource) (Resource, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(resource)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiClient.URL+"/"+resource.GetResourcePath(), &buf)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(apiClient.Username, apiClient.Password)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// TODO (AS) do the same for all methods
	if res.StatusCode != http.StatusCreated {
		json.NewEncoder(&buf).Encode(resource)
		return nil, fmt.Errorf("unexpected status code (%d) received: %s, response body: %s\nrequest sent:%s", res.StatusCode, res.Status, body, buf.String())
	}
	if err != nil {
		return nil, err
	}
	return parseResource(body)
}

func (apiClient *ApiClient) getResource(ctx context.Context, relativePath string) (Resource, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiClient.URL+relativePath, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(apiClient.Username, apiClient.Password)
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, NotFoundError
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d %s", res.StatusCode, res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return parseResource(b)
}

func (apiClient *ApiClient) updateResource(ctx context.Context, resource Resource) (Resource, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(resource)
	if err != nil {
		return nil, err
	}
	log.Printf("[TRACE] sending [[ %s ]]", buf.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, apiClient.URL+"/"+resource.GetResourcePath()+"/"+resource.GetID(), &buf)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(apiClient.Username, apiClient.Password)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if !isAlright(res.StatusCode) {
		return nil, fmt.Errorf("unexpected status code %d %s", res.StatusCode, res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return parseResource(b)
}

func (apiClient *ApiClient) deleteResource(ctx context.Context, relativePath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, apiClient.URL+relativePath, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(apiClient.Username, apiClient.Password)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if !isAlright(res.StatusCode) {
		return fmt.Errorf("unexpected status code %d %s", res.StatusCode, res.Status)
	}
	return nil
}

func (apiClient *ApiClient) post(ctx context.Context, requestBody interface{}, relativePath string, responseT interface{}) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiClient.URL+relativePath, &buf)
	if err != nil {
		return err
	}

	req.SetBasicAuth(apiClient.Username, apiClient.Password)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if !isAlright(res.StatusCode) {
		return fmt.Errorf("unexpected status code received %d %s", res.StatusCode, res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, responseT)
}

func (apiClient *ApiClient) get(ctx context.Context, relativePath string, responseT interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiClient.URL+relativePath, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(apiClient.Username, apiClient.Password)
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusNotFound {
		return NotFoundError
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d %s", res.StatusCode, res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, responseT)
}
