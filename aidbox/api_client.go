package aidbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
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
		URL:      URL,
		Username: username,
		Password: password,
	}
}

func isAlright(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusBadRequest
}

func (apiClient *ApiClient) createResource(ctx context.Context, resource Resource, responseTarget any) error {
	return apiClient.post(ctx, resource, path.Join("/", resource.GetResourcePath()), responseTarget)
}

func (apiClient *ApiClient) getResource(ctx context.Context, id string, responseTarget Resource) error {
	return apiClient.get(ctx, path.Join("/", responseTarget.GetResourcePath(), id), responseTarget)
}

func (apiClient *ApiClient) updateResource(ctx context.Context, resource Resource, responseTarget any) error {
	return apiClient.put(ctx, resource, path.Join("/", resource.GetResourcePath(), "/", resource.GetID()), responseTarget)
}

func (apiClient *ApiClient) deleteResource(ctx context.Context, id string, responseTarget Resource) error {
	return apiClient.send(ctx, struct{}{}, path.Join("/", responseTarget.GetResourcePath(), id), &struct{}{}, http.MethodDelete)
}

func (apiClient *ApiClient) put(ctx context.Context, requestBody interface{}, relativePath string, responseT interface{}) error {
	return apiClient.send(ctx, requestBody, relativePath, responseT, http.MethodPut)
}

func (apiClient *ApiClient) post(ctx context.Context, requestBody interface{}, relativePath string, responseT interface{}) error {
	return apiClient.send(ctx, requestBody, relativePath, responseT, http.MethodPost)
}

func (apiClient *ApiClient) send(ctx context.Context, requestBody interface{}, relativePath string, responseT interface{}, httpMethod string) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, httpMethod, apiClient.URL+relativePath, &buf)
	if err != nil {
		return err
	}

	apiClient.addAuthAndHost(req)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if !isAlright(res.StatusCode) {
		return errorToTerraform(req, res, requestBody, body)
	}
	// Deletes in general return the resource you deleted in the response body, but sometimes not (e.g. SearchParameter)
	if httpMethod == http.MethodDelete && len(body) == 0 {
		return nil
	} else {
		return json.Unmarshal(body, responseT)
	}
}

func (apiClient *ApiClient) get(ctx context.Context, relativePath string, responseT interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiClient.URL+relativePath, nil)
	if err != nil {
		return err
	}
	apiClient.addAuthAndHost(req)
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusNotFound {
		return NotFoundError
	}

	if !isAlright(res.StatusCode) {
		return errorToTerraform(req, res, struct{}{}, body)
	}
	return json.Unmarshal(body, responseT)
}

func (apiClient *ApiClient) addAuthAndHost(req *http.Request) {
	req.SetBasicAuth(apiClient.Username, apiClient.Password)
}

// errorToTerraform pretty prints the response body. Very often you get a 422 with useful details in the response body
// only. Print this into an error so terraform can show it to the user. Request details are also printed during testing.
func errorToTerraform(request *http.Request, response *http.Response, requestBody interface{}, responseBody []byte) error {
	var sensitiveDetails = ""
	var prettyResponse bytes.Buffer
	jsonParseErr := json.Indent(&prettyResponse, responseBody, "", "  ")
	if jsonParseErr != nil {
		return fmt.Errorf("unexpected status code (%d) received: %s\n\n"+
			"===== %s %s =====\n\n"+
			"===== RESPONSE BODY =====\n"+
			"%s\n",
			response.StatusCode,
			response.Status,
			request.Method, request.URL.String(),
			string(responseBody))
	}

	if os.Getenv("TF_ACC") == "1" {
		prettyRequest, err := json.MarshalIndent(requestBody, "", "  ")
		if err != nil {
			panic(err)
		}
		prettyHeaders, err := json.MarshalIndent(request.Header, "", "  ")
		if err != nil {
			panic(err)
		}
		sensitiveDetails = fmt.Sprintf("\n\n===== REQUEST HEADERS =====\n"+
			"%s\n\n"+
			"===== REQUEST BODY =====\n"+
			"%s",
			prettyHeaders,
			prettyRequest)
	}

	return fmt.Errorf("unexpected status code (%d) received: %s\n\n"+
		"===== %s %s =====%s\n\n"+
		"===== RESPONSE BODY =====\n"+
		"%s\n",
		response.StatusCode,
		response.Status,
		request.Method, request.URL.String(),
		sensitiveDetails,
		prettyResponse.String())
}
