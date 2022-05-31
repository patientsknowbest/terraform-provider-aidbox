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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.URL+"/"+resource.GetResourceName(), &buf)
	if err != nil {
		return nil, err
	}

	err = client.addAuthAndHost(ctx, req, boxId)
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
	err = client.clearCache(ctx)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return parseResource(b)
}

func (client *Client) getResource(ctx context.Context, relativePath, boxId string) (Resource, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.URL+relativePath, nil)
	if err != nil {
		return nil, err
	}
	err = client.addAuthAndHost(ctx, req, boxId)
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

func (client *Client) updateResource(ctx context.Context, resource Resource, boxId string) (Resource, error) {
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
	err = client.addAuthAndHost(ctx, req, boxId)
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
	err = client.clearCache(ctx)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return parseResource(b)
}

func (client *Client) deleteResource(ctx context.Context, relativePath, boxId string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, client.URL+relativePath, nil)
	if err != nil {
		return err
	}
	err = client.addAuthAndHost(ctx, req, boxId)
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
	return client.clearCache(ctx)
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
	req, err := http.NewRequestWithContext(ctx, "POST", client.URL+"/rpc", bytes.NewBuffer(wr))
	if err != nil {
		return err
	}
	err = client.addAuthAndHost(ctx, req, boxId)
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

/// Adds appropriate Authorization/Cookie & Host header to a request.
/// For multibox, we're expected to get the access-token from multibox API and use that.
/// and multibox will route our request to the appropriate box based on
func (client *Client) addAuthAndHost(ctx context.Context, req *http.Request, boxId string) error {
	if boxId == "" {
		req.SetBasicAuth(client.Username, client.Password)
	} else {
		box, err := client.getBox(ctx, boxId)
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

/// Some objects don't invalidate cache properly in multibox.
/// perform a cache reset to work around it.
/// https://github.com/Aidbox/Issues/issues/501
/// https://docs.aidbox.app/multibox/multibox-box-manager-api#multibox-drop-box-caches
func (client *Client) clearCache(ctx context.Context) error {
	if !client.IsMultibox {
		return nil
	}
	var response string
	err := client.rpcRequest(ctx, "multibox/drop-box-caches", struct{}{}, &response, "")
	if err != nil {
		return err
	}
	if response != "ok" {
		return fmt.Errorf("unexpected response to multibox/drop-box-caches: %s", response)
	}
	return nil
}
