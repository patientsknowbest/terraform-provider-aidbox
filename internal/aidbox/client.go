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
)

type Client struct {
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

func NewClient(URL, username, password string, isMultibox bool) *Client {
	return &Client{
		URL:        URL,
		Username:   username,
		Password:   password,
		IsMultibox: isMultibox,
	}
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
	case "AccessPolicy":
		r = &AccessPolicy{}
	default:
		return nil, fmt.Errorf("Unsupported resource type %s", s.ResourceType)
	}
	err = json.Unmarshal(in, &r)
	return r, err
}

func (client *Client) createResource(ctx context.Context, resource Resource, boxId string) (Resource, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(resource)
	if err != nil {
		return nil, err
	}
	url, err := client.getURL(ctx, boxId)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/"+resource.GetResourceName(), &buf)
	if err != nil {
		return nil, err
	}

	err = client.addAuth(ctx, req, boxId)
	if err != nil {
		return nil, err
	}
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

func (client *Client) getResource(ctx context.Context, relativePath, boxId string) (Resource, error) {
	url, err := client.getURL(ctx, boxId)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+relativePath, nil)
	if err != nil {
		return nil, err
	}
	err = client.addAuth(ctx, req, boxId)
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
		return nil, fmt.Errorf("Unexpected status code %d %s", res.StatusCode, res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return parseResource(b)
}

func (client *Client) updateResource(ctx context.Context, resource Resource, boxId string) (Resource, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(resource)
	if err != nil {
		return nil, err
	}
	log.Printf("[TRACE] sending [[ %s ]]", buf.String())
	url, err := client.getURL(ctx, boxId)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url+"/"+resource.GetResourceName()+"/"+resource.GetID(), &buf)
	if err != nil {
		return nil, err
	}
	err = client.addAuth(ctx, req, boxId)
	if err != nil {
		return nil, err
	}
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

func (client *Client) deleteResource(ctx context.Context, relativePath, boxId string) error {
	url, err := client.getURL(ctx, boxId)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url+relativePath, nil)
	if err != nil {
		return err
	}
	err = client.addAuth(ctx, req, boxId)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if !isAlright(res.StatusCode) {
		return fmt.Errorf("Unexpected status code %d %s", res.StatusCode, res.Status)
	}
	return nil
}

/// Some resources (multibox box management API for instance) are accessible only through the RPC endpoint
/// https://docs.aidbox.app/api-1/rpc-api
func (client *Client) rpcRequest(ctx context.Context, method string, request interface{}, responseT interface{}, boxId string) error {
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
	url, err := client.getURL(ctx, boxId)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url+"/rpc", bytes.NewBuffer(wr))
	if err != nil {
		return err
	}
	err = client.addAuth(ctx, req, boxId)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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

/// Adds appropriate authentication header to a request.
/// For multibox, we're expected to get the access-token from multibox API and use that.
func (client *Client) addAuth(ctx context.Context, req *http.Request, boxId string) error {
	if boxId == "" {
		req.SetBasicAuth(client.Username, client.Password)
	} else {
		box, err := client.getBox(ctx, boxId)
		if err != nil {
			return err
		}
		req.Header.Set("Cookie", fmt.Sprintf("aidbox-auth-token=%s", box.AccessToken))
	}
	return nil
}

/// Get the URL for an API call for the given box.
/// boxId may be empty if we're not using multibox
func (client *Client) getURL(ctx context.Context, boxId string) (string, error) {
	if boxId == "" {
		return client.URL, nil
	}
	box, err := client.getBox(ctx, boxId)
	if err != nil {
		return "", err
	}
	return box.BoxURL, nil
}

func (client *Client) getBox(ctx context.Context, boxId string) (*Box, error) {
	if !client.IsMultibox {
		return nil, fmt.Errorf("boxId provided to non-multibox client")
	}
	box := Box{}
	err := client.rpcRequest(ctx, "multibox/get-box", struct {
		Id string `json:"id"`
	}{boxId}, &box, "")
	if err != nil {
		return nil, err
	}
	return &box, nil
}
