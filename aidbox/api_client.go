package aidbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type ApiClient struct {
	URL      string
	Username string
	Password string

	// Username and Password are expected to be the superuser credentials if IsMultibox=true
	IsMultibox bool
}

type AidboxError string

const NotFoundError AidboxError = "Not found"

func (t AidboxError) Error() string {
	return string(t)
}

func NewApiClient(URL, username, password string, isMultibox bool) *ApiClient {
	return &ApiClient{
		URL:        URL,
		Username:   username,
		Password:   password,
		IsMultibox: isMultibox,
	}
}

func isAlright(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusBadRequest
}

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
	default:
		return nil, fmt.Errorf("Unsupported resource type %s", s.ResourceType)
	}
	err = json.Unmarshal(in, &r)
	return r, err
}

func (apiClient *ApiClient) createResource(ctx context.Context, resource Resource, boxId string) (Resource, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(resource)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiClient.URL+"/"+resource.GetResourceName(), &buf)
	if err != nil {
		return nil, err
	}

	err = apiClient.addAuthAndHost(ctx, req, boxId)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code received %d %s", res.StatusCode, res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return parseResource(b)
}

func (apiClient *ApiClient) getResource(ctx context.Context, relativePath, boxId string) (Resource, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiClient.URL+relativePath, nil)
	if err != nil {
		return nil, err
	}
	err = apiClient.addAuthAndHost(ctx, req, boxId)
	if err != nil {
		return nil, err
	}
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

func (apiClient *ApiClient) updateResource(ctx context.Context, resource Resource, boxId string) (Resource, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(resource)
	if err != nil {
		return nil, err
	}
	log.Printf("[TRACE] sending [[ %s ]]", buf.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, apiClient.URL+"/"+resource.GetResourceName()+"/"+resource.GetID(), &buf)
	if err != nil {
		return nil, err
	}
	err = apiClient.addAuthAndHost(ctx, req, boxId)
	if err != nil {
		return nil, err
	}
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

func (apiClient *ApiClient) deleteResource(ctx context.Context, relativePath, boxId string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, apiClient.URL+relativePath, nil)
	if err != nil {
		return err
	}
	err = apiClient.addAuthAndHost(ctx, req, boxId)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if !isAlright(res.StatusCode) {
		return fmt.Errorf("unexpected status code %d %s", res.StatusCode, res.Status)
	}
	return nil
}

func (apiClient *ApiClient) post(ctx context.Context, requestBody interface{}, relativePath, boxId string, responseT interface{}) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiClient.URL+relativePath, &buf)
	if err != nil {
		return err
	}

	err = apiClient.addAuthAndHost(ctx, req, boxId)
	if err != nil {
		return err
	}
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

func (apiClient *ApiClient) get(ctx context.Context, relativePath, boxId string, responseT interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiClient.URL+relativePath, nil)
	if err != nil {
		return err
	}
	err = apiClient.addAuthAndHost(ctx, req, boxId)
	if err != nil {
		return err
	}
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

// Some resources (multibox box management API for instance) are accessible only through the RPC endpoint
// https://docs.aidbox.app/api-1/rpc-api
func (apiClient *ApiClient) rpcRequest(ctx context.Context, method string, request interface{}, responseT interface{}, boxId string) error {
	rm, err := json.Marshal(request)
	if err != nil {
		return err
	}
	wrapper := struct {
		Method string          `json:"method"`
		Params json.RawMessage `json:"params"`
	}{
		Method: method,
		Params: rm,
	}
	tflog.Trace(ctx, "rpcRequest", map[string]interface{}{
		"wrapper": wrapper,
	})
	wr, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", apiClient.URL+"/rpc", bytes.NewBuffer(wr))
	if err != nil {
		return err
	}
	err = apiClient.addAuthAndHost(ctx, req, boxId)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if !isAlright(resp.StatusCode) {
		msg, err := ioutil.ReadAll(resp.Body)
		var msgStr string
		if err != nil {
			msgStr = err.Error()
		} else {
			msgStr = string(msg)
		}
		return fmt.Errorf("unexpected status code from RPC request %d %v [%s]", resp.StatusCode, resp.Status, msgStr)
	}
	defer func() { _ = resp.Body.Close() }()
	var response struct {
		Result json.RawMessage `json:"result,omitempty"`
		Error  json.RawMessage `json:"error,omitempty"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return fmt.Errorf("error response from RPC call %v", string(response.Error))
	}
	return json.Unmarshal(response.Result, responseT)
}

// Adds appropriate Authorization/Cookie & Host header to a request.
// For multibox, we're expected to get the access-token from multibox API and use that.
// and multibox will route our request to the appropriate box based on
func (apiClient *ApiClient) addAuthAndHost(ctx context.Context, req *http.Request, boxId string) error {
	if boxId == "" {
		req.SetBasicAuth(apiClient.Username, apiClient.Password)
	} else {
		box, err := apiClient.getBox(ctx, boxId)
		if err != nil {
			return err
		}
		boxURL, err := url.Parse(box.BoxURL)
		if err != nil {
			return err
		}
		tflog.Info(ctx, "addAuthAndHost", map[string]interface{}{
			"hostname":     boxURL.Hostname(),
			"access_token": box.AccessToken,
		})
		req.Host = boxURL.Hostname()
		req.Header.Set("Cookie", fmt.Sprintf("aidbox-auth-token=%s", box.AccessToken))
	}
	return nil
}

func (apiClient *ApiClient) getBox(ctx context.Context, boxId string) (*Box, error) {
	if !apiClient.IsMultibox {
		return nil, fmt.Errorf("boxId provided to non-multibox client")
	}
	box := Box{}
	err := apiClient.rpcRequest(ctx, "multibox/get-box", struct {
		Id string `json:"id"`
	}{boxId}, &box, "")
	if err != nil {
		return nil, err
	}
	return &box, nil
}
